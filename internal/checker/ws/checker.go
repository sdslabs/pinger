// Copyright (c) 2020 SDSLabs
// Use of this source code is governed by an MIT license
// details of which can be found in the LICENSE file.

package ws

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"strings"

	"github.com/sdslabs/status/internal/checker"
)

const (
	// checkerName is the name of the checker.
	checkerName = "WS"

	// messageType is the type of output for checking message
	messageType = "MESSAGE"

	// timeoutType is the type of output for checking timeout.
	timeoutType = "TIMEOUT"

	// headerType is the type of output/payload for HTTP header.
	headerType = "HEADER"

	// statusCodeType is the type of output for a status code.
	statusCodeType = "STATUSCODE"

	// bodyType is the type of output for response body.
	bodyType = "BODY"

	// splitDelim separates the messages in output value.
	splitDelim = "\n---\n"
)

func init() {
	checker.Register(checkerName, func() checker.Checker { return new(Checker) })
}

// Checker sends an WS ECHO request and checks if a reply is returned.
type Checker struct {
	prober      *Prober
	outputType  string
	outputValue interface{}
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
		timeoutType:    validateNil,
		statusCodeType: validateOutputStatusCode,
		bodyType:       validateNil,
		headerType:     validatePayloadHeader,
		messageType:    validateNil,
	}
	if err := checker.ValidateComponent(check.GetOutput(), validateOutputMap); err != nil {
		return fmt.Errorf("input: %w", err)
	}

	validateTargetMap := validationMap{"URL": validateTarget}
	if err := checker.ValidateComponent(check.GetTarget(), validateTargetMap); err != nil {
		return fmt.Errorf("target: %w", err)
	}

	validatePayloadMap := validationMap{
		headerType:  validatePayloadHeader,
		messageType: validateNil,
	}
	for i, p := range check.GetPayloads() {
		if err := checker.ValidateComponent(p, validatePayloadMap); err != nil {
			return fmt.Errorf("payload %d: %w", i, err)
		}
	}
	return nil
}

// Provision initializes required fields for c's execution.
func (c *Checker) Provision(check checker.Check) error {
	targetURL := check.GetTarget().GetValue()
	timeout := check.GetTimeout()
	headers := map[string]string{}
	var messages []string
	for _, p := range check.GetPayloads() {
		switch p.GetType() {
		case headerType:
			k, v, err := extractKVPair(p.GetValue())
			if err != nil {
				return err
			}

			headers[k] = v

		case messageType:
			messages = append(messages, p.GetValue())

		default:
		}
	}
	prober, err := NewProber(targetURL, headers, messages, timeout)
	if err != nil {
		return err
	}

	var outputValue interface{}
	outputType := check.GetOutput().GetType()
	switch outputType {
	case headerType:
		k, v, err := extractKVPair(check.GetOutput().GetValue())
		if err != nil {
			return err
		}

		outputValue = kvPair{k: k, v: v}

	case statusCodeType:
		st, err := extractStatusCode(check.GetOutput().GetValue())
		if err != nil {
			return err
		}

		outputValue = st

	case bodyType:
		outputValue = check.GetOutput().GetValue()

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
		messages := strings.Split(check.GetOutput().GetValue(), splitDelim)
		outputValue = messages

	default:
	}

	c.outputType = outputType
	c.outputValue = outputValue
	c.prober = prober
	return nil
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

	case bodyType:
		body, ok := c.outputValue.(string)
		if !ok {
			return nil, fmt.Errorf("internal error: outputValue of body not a string")
		}
		buf := new(strings.Builder)
		_, err := io.Copy(buf, probeResult.Body)
		if err != nil {
			return nil, fmt.Errorf("invalid buffer :%s", err)
		}
		resultBody := buf.String()
		if resultBody == body {
			result.Successful = true
		}

	case statusCodeType:
		st, ok := c.outputValue.(int)
		if !ok {
			return nil, fmt.Errorf("internal error: outputValue of status code not an int")
		}

		if st == probeResult.StatusCode {
			result.Successful = true
		}
	case messageType:
		msgs, ok := c.outputValue.([]string)
		if !ok {
			return nil, fmt.Errorf("internal error: outputValue of messages not []string")
		}

		if len(msgs) == len(probeResult.Response) {
			validateMessages(msgs, probeResult, result)
		}
	case headerType:
		header, ok := c.outputValue.(kvPair)
		if !ok {
			return nil, fmt.Errorf("internal error: outputValue of header not a kvPair")
		}

		if probeResult.Headers.Get(header.k) == header.v {
			result.Successful = true
		}
	default:
	}

	return result, nil
}

// validationMap is an alias of map used for validating components.
type validationMap = map[string]func(string) error

// validateMessages validates if prober response is same as desired output
func validateMessages(msgs []string, probeResult *ProbeResult, result *checker.Result) {
	allEq := true
	for i := range msgs {
		if msgs[i] != probeResult.Response[i] {
			allEq = false
			break
		}
	}
	if allEq {
		result.Successful = true
	}
}

// validateTarget validates if the target value is a URL.
func validateTarget(val string) error {
	u, err := url.Parse(val)
	if err != nil {
		return fmt.Errorf("value not a valid URL: %w", err)
	}

	if !strings.EqualFold(u.Scheme, "ws") && !strings.EqualFold(u.Scheme, "wss") {
		return fmt.Errorf("not a valid ws(s) url: %s", val)
	}

	return nil
}

// func validateInout validates the check input
func validateInput(val string) error {
	switch val {
	case "", "PING", "ECHO":
	default:
		return fmt.Errorf("invalid value: %s", val)
	}
	return nil
}

// func validateNil doesn't validate anything
func validateNil(string) error { return nil }

// validateOutputStatusCode validates if the output value is a status code.
func validateOutputStatusCode(val string) error {
	_, err := extractStatusCode(val)
	return err
}

// func extractStatusCode validates and extracts status code from a string value.
func extractStatusCode(val string) (int, error) {
	intVal, err := strconv.Atoi(val)
	if err != nil {
		return 0, fmt.Errorf("invalid status code: %v", err)
	}
	if intVal > 599 || intVal < 100 {
		return 0, fmt.Errorf("status code should lie between 100 and 600")
	}
	return intVal, err
}

// func validatePayloadHeader validates the header type payload.
func validatePayloadHeader(val string) error {
	_, _, err := extractKVPair(val)
	return err
}

// kvPair is a string-string key value pair.
type kvPair struct{ k, v string }

// extractsKVPair validates if the value is of the form K=V and extracts the
// key and value.
func extractKVPair(val string) (k, v string, err error) {
	kv := strings.SplitN(val, "=", 2)
	if len(kv) != 2 {
		return "", "", fmt.Errorf("value should be of the format K=V")
	}
	k = kv[0]
	v = kv[1]

	if k == "" {
		return "", "", fmt.Errorf("key cannot be empty")
	}

	return
}

// Interface guard
var _ checker.Checker = (*Checker)(nil)
