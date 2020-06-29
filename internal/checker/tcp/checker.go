// Copyright (c) 2020 SDSLabs
// Use of this source code is governed by an MIT license
// details of which can be found in the LICENSE file.

package tcp

import (
	"context"
	"fmt"
	"net"

	"github.com/sdslabs/status/internal/checker"
)

const (
	// checkerName is the name of the checker.
	checkerName = "TCP"

	//messageType is the type of output for checking message
	messageType = "MESSAGE"

	// timeoutType is the type of output for checking timeout.
	timeoutType = "TIMEOUT"
)

func init() {
	checker.Register(checkerName, func() checker.Checker { return new(Checker) })
}

// Checker sends an TCP ECHO request and checks if a reply is returned.
type Checker struct {
	prober      *Prober
	outputType  string
	outputValue string
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
		return fmt.Errorf("input: %w", err)
	}

	validatePayloadMap := validationMap{
		messageType: validatePayloadMessage,
	}
	for i, p := range check.GetPayloads() {
		if err := checker.ValidateComponent(p, validatePayloadMap); err != nil {
			return fmt.Errorf("payload %d: %w", i, err)
		}
	}

	validateTargetMap := validationMap{"ADDRESS": validateTarget}
	if err := checker.ValidateComponent(check.GetTarget(), validateTargetMap); err != nil {
		return fmt.Errorf("input: %w", err)
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
	c.outputValue = check.GetOutput().GetValue()
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

	if !probeResult.Timeout {
		result.Successful = true
		result.Timeout = false
	}

	if probeResult.Timeout {
		return result, nil
	}

	result.Timeout = false

	switch c.outputType {

	case timeoutType:
		result.Successful = true

		// case messageType:
		// message, ok = c.outputValue
		// if !ok {
		// 	return nil, fmt.Errorf("internal error: outputValue of body not a string")
		// }

		// if probeResult.Body == body {
		// 	result.Successful = true
		// }
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
	_, err := net.ResolveTCPAddr("tcp4", val)
	if err != nil {
		fmt.Errorf("value is not a valid TCP address: %w", err)
	}
	return nil
}

// validatePayloadMessage validates if the target value is a URL.
func validatePayloadMessage(val string) error {
	if len(val) != 0 {
		fmt.Sprintf("Message Cannot be Empty")
	}
	return nil
}

// validateNil doesn't validate anything.
func validateNil(string) error { return nil }
