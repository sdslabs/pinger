// Copyright (c) 2020 SDSLabs
// Use of this source code is governed by an MIT license
// details of which can be found in the LICENSE file.

package tcp

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/sdslabs/status/internal/checker"
)

// Prober creates a TCP connection with the address and sends and receives
// messages to through the connection.
type Prober struct {
	address  string
	timeout  time.Duration
	messages []string
}

// NewProber creates a prober that establishes a TCP connection.
func NewProber(address string, messages []string, timeout time.Duration) (*Prober, error) {
	if timeout <= 0 {
		return nil, fmt.Errorf("timeout should be > 0")
	}

	return &Prober{
		address:  address,
		timeout:  timeout,
		messages: messages,
	}, nil
}

// Probe establishes a TCP connection and exchanges the messages from client
// to the server and vice-versa.
func (p *Prober) Probe(ctx context.Context) (*ProbeResult, error) {
	startTime := time.Now()
	deadline := startTime.Add(p.timeout)

	baseCtx := ctx
	if ctx == nil {
		baseCtx = context.Background()
	}

	probeCtx, cancel := context.WithTimeout(baseCtx, p.timeout)
	defer cancel()

	timeoutResult := &ProbeResult{
		Timeout:   true,
		StartTime: startTime,
		Duration:  p.timeout,
	}

	conn, err := p.dialTCP(probeCtx)
	if err != nil {
		if checker.ErrIsTimeout(err) {
			return timeoutResult, nil
		}

		return nil, err
	}
	defer conn.Close() // nolint:errcheck

	if len(p.messages) == 0 {
		return &ProbeResult{
			Timeout:   false,
			StartTime: startTime,
			Duration:  time.Since(startTime),
		}, nil
	}

	response := make([]string, len(p.messages))
	buf := make([]byte, checker.MaxMessageSize)

	for i := range p.messages {
		select {
		case <-probeCtx.Done():
			err = probeCtx.Err()
			if checker.ErrIsTimeout(err) {
				return timeoutResult, nil
			}

			return nil, err

		default:
		}

		err = p.sendMessage(probeCtx, conn, deadline, i)
		if err != nil {
			if checker.ErrIsTimeout(err) {
				return timeoutResult, nil
			}

			return nil, err
		}

		resp, err := p.receiveMessage(probeCtx, conn, deadline, buf)
		if err != nil && err != io.EOF {
			if checker.ErrIsTimeout(err) {
				return timeoutResult, nil
			}

			return nil, err
		}

		response[i] = resp

		if err == io.EOF {
			// break if EOF is reached even if all the messages are not received.
			break
		}
	}

	return &ProbeResult{
		Timeout:   false,
		StartTime: startTime,
		Duration:  time.Since(startTime),
		Response:  response,
	}, nil
}

// dialTCP establishes a TCP connection between the client and the server.
func (p *Prober) dialTCP(ctx context.Context) (*net.TCPConn, error) {
	dialer := net.Dialer{Timeout: p.timeout}
	conn, err := dialer.DialContext(ctx, "tcp", p.address)
	if err != nil {
		return nil, err
	}

	tcpConn, ok := conn.(*net.TCPConn)
	if !ok {
		return nil, fmt.Errorf("not a valid TCP connection")
	}

	return tcpConn, nil
}

// sendMessage sends a message to the address.
func (p *Prober) sendMessage(
	ctx context.Context,
	conn net.Conn,
	deadline time.Time,
	index int,
) error {
	if err := conn.SetWriteDeadline(deadline); err != nil {
		return err
	}

	errChan := make(chan error)
	go send(conn, p.messages[index], errChan)

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errChan:
		return err
	}
}

// receiveMessage receives the response from the address.
func (p *Prober) receiveMessage(
	ctx context.Context,
	conn net.Conn,
	deadline time.Time,
	buf []byte,
) (string, error) {
	if err := conn.SetReadDeadline(deadline); err != nil {
		return "", err
	}

	packChan := make(chan packet)
	go receive(conn, buf, packChan)

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case pack := <-packChan:
		return pack.message, pack.err
	}
}

// send writes a message to the connection.
func send(conn io.Writer, msg string, stream chan<- error) {
	_, err := conn.Write([]byte(msg))
	stream <- err
}

// packet represents the TCP response.
type packet struct {
	message string
	err     error
}

// receive receives a message from the connection.
func receive(conn io.Reader, buf []byte, stream chan<- packet) {
	readLen, err := conn.Read(buf)
	if err != nil {
		stream <- packet{err: err}
		return
	}

	final := bytes.Trim(buf[:readLen], "\x00")

	stream <- packet{message: string(final)}
}

// ProbeResult is the result for TCP probe.
type ProbeResult struct {
	Timeout   bool
	StartTime time.Time
	Duration  time.Duration
	Response  []string
}
