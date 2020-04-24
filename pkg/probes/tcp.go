package probes

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"time"
)

// NewTCPProbe creates a TCP Prober for given address and exits if it takes longer
// than the timeout. `message` is the request message to be sent to the address.
func NewTCPProbe(address, message string, timeout time.Duration) (*TCPProber, error) {
	pr := &TCPProber{}
	if err := pr.SetAddress(address); err != nil {
		return nil, err
	}

	pr.SetTimeout(timeout)
	pr.SetMessage(message)
	return pr, nil
}

// TCPProber probes a TCP address by sending a message as request.
type TCPProber struct {
	timeout time.Duration
	address string
	message string
}

// SetAddress sets the address to probe.
func (pr *TCPProber) SetAddress(address string) error {
	_, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return err
	}

	pr.address = address
	return nil
}

// SetTimeout sets the timeout after which the probe exits.
func (pr *TCPProber) SetTimeout(timeout time.Duration) {
	pr.timeout = timeout
}

// SetMessage sets the message to be sent as the request.
func (pr *TCPProber) SetMessage(message string) {
	pr.message = message
}

// Probe method is used to execute the prober. It sends the given message to the
// given address and returns a `TCPProbeResult` which has the response received.
func (pr *TCPProber) Probe() (*TCPProbeResult, error) {
	startTime := time.Now()

	timeoutResult := &TCPProbeResult{
		Timeout:   true,
		StartTime: startTime,
		Duration:  pr.timeout,
	}

	conn, err := dialTCPTimeout(pr.address, pr.timeout)
	if err != nil {
		if errIsTimeout(err) {
			return timeoutResult, nil
		}

		return nil, err
	}
	defer conn.Close() //nolint:errcheck
	if pr.message == "" {
		return &TCPProbeResult{
			Timeout:   false,
			StartTime: startTime,
			Duration:  time.Since(startTime),
			Response:  "",
		}, nil
	}
	if err := conn.SetDeadline(startTime.Add(pr.timeout)); err != nil {
		return nil, err
	}
	if _, err := conn.Write([]byte(pr.message)); err != nil {
		if errIsTimeout(err) {
			return timeoutResult, nil
		}

		return nil, err
	}

	if err := conn.CloseWrite(); err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, conn); err != nil {
		if errIsTimeout(err) {
			return timeoutResult, nil
		}

		return nil, err
	}

	// Trim null characters from the buffer
	response := string(bytes.Trim(buf.Bytes(), "\x00"))

	return &TCPProbeResult{
		Timeout:   false,
		StartTime: startTime,
		Duration:  time.Since(startTime),
		Response:  response,
	}, nil
}

// dialTCPTimeout is same as DialTCP but terminates after the `timeout`.
func dialTCPTimeout(address string, timeout time.Duration) (*net.TCPConn, error) {
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return nil, err
	}

	tcpConn, ok := conn.(*net.TCPConn)
	if !ok {
		return nil, fmt.Errorf("not a valid TCP connection")
	}

	return tcpConn, nil
}

// TCPProbeResult stores the result of TCP probe. It contains the response
// received from the address.
type TCPProbeResult struct {
	Timeout   bool
	StartTime time.Time
	Duration  time.Duration
	Response  string
}
