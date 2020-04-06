package probes

import (
	"bytes"
	"fmt"
	"net"
	"time"
)

// NewUDPProbe creates a UDP Prober for given address and exits if it takes longer
// than the timeout. `message` is the request message to be sent to the address.
func NewUDPProbe(address, message string, timeout time.Duration) (*UDPProber, error) {
	pr := &UDPProber{}
	if err := pr.SetAddress(address); err != nil {
		return nil, err
	}
	pr.SetTimeout(timeout)
	pr.SetMessage(message)
	return pr, nil
}

// UDPProber probes a UDP address by sending a message as request.
type UDPProber struct {
	timeout time.Duration
	address string
	message string
}

// SetAddress sets the address to probe.
func (pr *UDPProber) SetAddress(address string) error {
	_, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return err
	}

	pr.address = address
	return nil
}

// SetTimeout sets the timeout after which the probe exits.
func (pr *UDPProber) SetTimeout(timeout time.Duration) {
	pr.timeout = timeout
}

// SetMessage sets the message to be sent as the request.
func (pr *UDPProber) SetMessage(message string) {
	pr.message = message
}

// Probe method is used to execute the prober. It sends the given message to the
// given address and returns a `UDPProbeResult` which has the response received.
func (pr *UDPProber) Probe() (*UDPProbeResult, error) {
	startTime := time.Now()

	timeoutResult := &UDPProbeResult{
		Timeout:   true,
		StartTime: startTime,
		Duration:  pr.timeout,
	}

	conn, err := dialUDPTimeout(pr.address, pr.timeout)
	if err != nil {
		if errIsTimeout(err) {
			return timeoutResult, nil
		}

		return nil, err
	}
	defer conn.Close() //nolint:errcheck

	if pr.message == "" {
		return nil, fmt.Errorf("For a UDP Connection, empty messages are not supported")
	}

	err = conn.SetDeadline(startTime.Add(pr.timeout))
	if err != nil {
		return nil, err
	}

	_, err = conn.Write([]byte(pr.message))
	if err != nil {
		if errIsTimeout(err) {
			return timeoutResult, nil
		}

		return nil, err
	}

	buf := make([]byte, maxBufferCapacity)
	_, err = conn.Read(buf)
	if err != nil {
		if errIsTimeout(err) {
			return timeoutResult, nil
		}

		return nil, err
	}

	// Trim null characters from the result
	response := string(bytes.Trim(buf, "\x00"))

	return &UDPProbeResult{
		Timeout:   false,
		StartTime: startTime,
		Duration:  time.Since(startTime),
		Response:  response,
	}, nil
}

// dialUDPTimeout is same as DialUDP but terminates after the `timeout`.
func dialUDPTimeout(address string, timeout time.Duration) (*net.UDPConn, error) {
	conn, err := net.DialTimeout("udp", address, timeout)
	if err != nil {
		return nil, err
	}

	UDPConn, ok := conn.(*net.UDPConn)
	if !ok {
		return nil, fmt.Errorf("not a valid UDP connection")
	}

	return UDPConn, nil
}

// UDPProbeResult stores the result of UDP probe. It contains the response
// received from the address.
type UDPProbeResult struct {
	Timeout   bool
	StartTime time.Time
	Duration  time.Duration
	Response  string
}
