// Copyright (c) 2020 SDSLabs
// Use of this source code is governed by an MIT license
// details of which can be found in the LICENSE file.

package tcp

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/sdslabs/status/internal/checker"
)

const (
	// checkerName is the name of the checker.
	checkerName = "TCP"

	// messageType is the type of output for checking message
	messageType = "MESSAGE"

	// timeoutType is the type of output for checking timeout.
	timeoutType = "TIMEOUT"

	// splitDelim separates the messages in output value.
	splitDelim = "\n---\n"
)

func init() {
	checker.Register(checkerName, func() checker.Checker { return new(Checker) })
}

// Checker sends an TCP ECHO request and checks if a reply is returned.
type Checker struct {
	prober      *Prober
	outputType  string
	outputValue []string
}

// Validate validates the check configuration.
func (c *Checker) Validate(check checker.Check) error {
	if check.GetTimeout() <= 0 {
		return fmt.Errorf("timeout should be > 0")
	}

	validateInputMap := validationMap{checkerName: validateInput}
	if err := checker.ValidateComponent(check.GetInput(), validateInputMap); err != nil {
		return fmt.Errorf("input: %w", err)
	}

	validateOutputMap := validationMap{
		timeoutType: validateNil,
		messageType: validateNil,
	}
	if err := checker.ValidateComponent(check.GetOutput(), validateOutputMap); err != nil {
		return fmt.Errorf("output: %w", err)
	}

	validateTargetMap := validationMap{"ADDRESS": validateTarget}
	if err := checker.ValidateComponent(check.GetTarget(), validateTargetMap); err != nil {
		return fmt.Errorf("target: %w", err)
	}

	validatePayloadMap := validationMap{
		messageType: validatePayloadMessage,
	}
	for i, p := range check.GetPayloads() {
		if err := checker.ValidateComponent(p, validatePayloadMap); err != nil {
			return fmt.Errorf("payload %d: %w", i, err)
		}
	}

	return nil
}

// Provision initializes required fields for c's execution.
func (c *Checker) Provision(check checker.Check) (err error) {
	address := check.GetTarget().GetValue()
	timeout := check.GetTimeout()

	var messages []string
	for _, p := range check.GetPayloads() {
		messages = append(messages, p.GetValue())
	}

	c.prober, err = NewProber(address, messages, timeout)
	c.outputType = check.GetOutput().GetType()
	switch c.outputType {
	case messageType:
		// messages are split via "\n---\n", so multiple messages which are to be
		// verified should be separated via the same.
		// For example, if the response should have two messages -- "hello" and say
		// "world", output should be:
		// 	"hello\n---\nworld"
		// OR
		// 	hello
		// 	---
		// 	world
		c.outputValue = strings.Split(check.GetOutput().GetValue(), splitDelim)
	default:
	}
	return
}

// Execute executes the check.
func (c *Checker) Execute(ctx context.Context) (*checker.Result, error) {
	probeResult, err := c.prober.Probe(ctx)
	if err != nil {
		return nil, err
	}

	result := &checker.Result{
		Timeout:   true,
		StartTime: probeResult.StartTime,
		Duration:  probeResult.Duration,
	}

	if probeResult.Timeout {
		return result, nil
	}

	result.Timeout = false

	switch c.outputType {
	case timeoutType:
		result.Successful = true

	case messageType:

		if len(c.outputValue) == len(probeResult.Response) {
			allEq := true
			for i := range c.outputValue {
				if c.outputValue[i] != probeResult.Response[i] {
					allEq = false
					break
				}
			}

			if allEq {
				result.Successful = true
			}
		}
	}

	return result, nil
}

// validationMap is an alias of map used for validating components.
type validationMap = map[string]func(string) error

// validateInput validates the check input.
func validateInput(val string) error {
	switch val {
	case "", "PING", "ECHO":
	default:
		return fmt.Errorf("invalid value: %s", val)
	}

	return nil
}

// validateTarget validates if the target value is a URL.
func validateTarget(val string) error {
	if !isValidTCPAddr(val) {
		return fmt.Errorf("target is not a valid address: %s", val)
	}

	return nil
}

// validatePayloadMessage validates if the target value is a URL.
func validatePayloadMessage(val string) error {
	if val == "" {
		return fmt.Errorf("message cannot be empty")
	}

	return nil
}

// validateOutputMessages validates if the output value messages is valid.
func validateOutputMessages(val string) error {
	msgs := strings.Split(val, splitDelim)
	if len(msgs) == 0 {
		return fmt.Errorf("messages cannot be empty")
	}

	for i := range msgs {
		if msgs[i] == "" {
			return fmt.Errorf("messages cannot be empty")
		}
	}

	return nil
}

// isValidTCPAddr tells if s is a valid TCP address or not.
func isValidTCPAddr(s string) bool {
	host, port, err := net.SplitHostPort(s)
	if err != nil {
		return false
	}

	if !isValidAddress(host) {
		return false
	}

	portNo, err := strconv.ParseUint(port, 10, 0)
	if err != nil {
		return false
	}

	if portNo > 65535 {
		return false
	}

	return true
}

// isValidAddress tells if s is a valid address or not.
func isValidAddress(s string) bool {
	// an IP address is a valid address.
	if net.ParseIP(s) != nil {
		return true
	}

	return checker.AddressRegex.MatchString(s)
}

// validateNil doesn't validate anything.
func validateNil(string) error { return nil }
