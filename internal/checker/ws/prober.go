// Copyright (c) 2020 SDSLabs
// Use of this source code is governed by an MIT license
// details of which can be found in the LICENSE file.

package ws

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/websocket"

	"github.com/sdslabs/status/internal/checker"
)

// Prober establishes a websocket connection with the server and sends and
// receives messages.
type Prober struct {
	url      string
	headers  http.Header
	messages []string
	timeout  time.Duration
}

// NewProber creates a new websocket prober.
func NewProber(
	targetURL string,
	headers map[string]string,
	messages []string,
	timeout time.Duration,
) (*Prober, error) {
	if timeout <= 0 {
		return nil, fmt.Errorf("timeout should be > 0")
	}

	parsedHeaders := make(http.Header)
	for key, val := range headers {
		parsedHeaders.Add(key, val)
	}

	return &Prober{
		url:      targetURL,
		headers:  parsedHeaders,
		messages: messages,
		timeout:  timeout,
	}, nil
}

// Probe probes the websocket server and exchanges messages.
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

	conn, resp, err := p.dialWS(probeCtx)
	if err != nil {
		if checker.ErrIsTimeout(err) {
			return timeoutResult, nil
		}

		return nil, err
	}

	defer conn.Close()      //nolint:errcheck
	defer resp.Body.Close() //nolint:errcheck

	if len(p.messages) == 0 {
		return &ProbeResult{
			Timeout:    false,
			StartTime:  startTime,
			Duration:   time.Since(startTime),
			StatusCode: resp.StatusCode,
			Body:       resp.Body,
			Headers:    resp.Header,
		}, nil
	}

	response := make([]string, len(p.messages))

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

		resp, err := p.receiveMessage(probeCtx, conn, deadline)
		if err != nil && err != io.EOF {
			if checker.ErrIsTimeout(err) {
				return timeoutResult, nil
			}

			return nil, err
		}

		response[i] = resp
	}

	return &ProbeResult{
		Timeout:    false,
		StartTime:  startTime,
		Duration:   time.Since(startTime),
		StatusCode: resp.StatusCode,
		Body:       resp.Body,
		Headers:    resp.Header,
		Response:   response,
	}, nil
}

// dialWS establishes a websocket connection between the client and the
// server.
func (p *Prober) dialWS(ctx context.Context) (*websocket.Conn, *http.Response, error) {
	dialer := &websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: p.timeout,
	}

	return dialer.DialContext(ctx, p.url, p.headers)
}

// sendMessage sends a message to the address.
func (p *Prober) sendMessage(
	ctx context.Context,
	conn *websocket.Conn,
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
	conn *websocket.Conn,
	deadline time.Time,
) (string, error) {
	if err := conn.SetReadDeadline(deadline); err != nil {
		return "", err
	}

	packChan := make(chan packet)
	go receive(conn, packChan)

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case pack := <-packChan:
		return pack.message, pack.err
	}
}

// send writes a message to the connection.
func send(conn *websocket.Conn, msg string, stream chan<- error) {
	stream <- conn.WriteMessage(websocket.TextMessage, []byte(msg))
}

// packet represents the TCP response.
type packet struct {
	message string
	err     error
}

// receive receives a message from the connection.
func receive(conn *websocket.Conn, stream chan<- packet) {
	_, recv, err := conn.ReadMessage()
	if err != nil {
		stream <- packet{err: err}
		return
	}

	final := bytes.Trim(recv, "\x00")

	stream <- packet{message: string(final)}
}

// ProbeResult is the result of websocket probe.
type ProbeResult struct {
	Timeout    bool
	StartTime  time.Time
	Duration   time.Duration
	Response   []string
	StatusCode int
	Body       io.ReadCloser
	Headers    http.Header
}
