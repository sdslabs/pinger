package check

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/sdslabs/status/pkg/config"
	"github.com/sdslabs/status/pkg/controller"
	"github.com/sdslabs/status/pkg/defaults"
	"github.com/sdslabs/status/pkg/probes"
)

const headerDelimeter = "="

var validHttpOutputTypes map[string]validationFunction = map[string]validationFunction{
	"status_code": validateStatusCode,
	"body":        validateBody,
	"header":      validateKVPair,
}

// NewHTTPChecker creates a new Checker for HTTP requests.
func NewHTTPChecker(agentCheck config.Check) (*HTTPChecker, error) {
	err := validateHTTPCheck(agentCheck)
	if err != nil {
		return nil, fmt.Errorf("VALIDATION_ERROR: %s", err)
	}

	method := agentCheck.GetInput().GetValue()
	if method == "" {
		method = defaults.DefaultHTTPMethod
	}

	params := make(map[string]string)
	headers := make(map[string]string)

	for _, payload := range agentCheck.GetPayloads() {
		kv := strings.SplitN(payload.GetValue(), headerDelimeter, 2)

		switch payload.GetType() {
		case "header":
			headers[kv[0]] = kv[1]
		case "parameter":
			params[kv[0]] = kv[1]
		}
	}

	timeout := time.Duration(agentCheck.GetTimeout()) * time.Second
	if timeout <= 0 {
		timeout = defaults.DefaultHTTPProbeTimeout
	}

	return &HTTPChecker{
		Method: method,

		URL:     agentCheck.GetTarget().GetValue(),
		Payload: params,
		Headers: headers,

		HTTPOutput: &TVComponent{
			Type:  agentCheck.GetOutput().GetType(),
			Value: agentCheck.GetOutput().GetValue(),
		},

		Timeout: timeout,
	}, nil
}

func validateStatusCode(val string) error {
	intVal, err := strconv.Atoi(val)
	if err != nil {
		return fmt.Errorf("status code prvoided is not parsable: %s", err)
	}
	if intVal > 599 || intVal < 100 {
		return fmt.Errorf("status code is not valid (expected between: 100 - 599) got %d", intVal)
	}

	return nil
}

func validateBody(val string) error {
	if val == "" {
		return fmt.Errorf("cannot have empty body as http check output")
	}

	return nil
}

func validateKVPair(val string) error {
	kv := strings.SplitN(val, headerDelimeter, 2)
	if len(kv) != 2 {
		return fmt.Errorf("header value is not valid, must have format: HEADER=<Header value>")
	}

	return nil
}

func validateHTTPCheck(agentCheck config.Check) error {
	if err := validateHTTPInput(agentCheck.GetInput()); err != nil {
		return err
	}

	if err := validateHTTPOutput(agentCheck.GetOutput()); err != nil {
		return err
	}

	if err := validateHTTPTarget(agentCheck.GetTarget()); err != nil {
		return err
	}

	return validateHTTPPayload(agentCheck.GetPayloads())
}

func validateHTTPInput(input config.CheckComponent) error {
	inputVal := input.GetValue()
	if inputVal != "GET" && inputVal != "POST" && inputVal != "" {
		return fmt.Errorf("for HTTP input the provided method(%s) is not supported", inputVal)
	}

	return nil
}

func validateHTTPOutput(output config.CheckComponent) error {
	validateFunc, ok := validHttpOutputTypes[output.GetType()]
	if !ok {
		return fmt.Errorf("provided Output Type(%s) is not valid for HTTP input", output.GetType())
	}

	return validateFunc(output.GetValue())
}

func validateHTTPTarget(target config.CheckComponent) error {
	// We don't check the value of type for the target here
	// as for HTTP Check the target is always a URL and we can check it that way only.
	// Just for consistency of types not being nil, we check if it's equal to "url"
	typ := target.GetType()
	if typ != "url" {
		return fmt.Errorf("target type %s is not supported", typ)
	}
	url, err := url.Parse(target.GetValue())
	if err != nil {
		return fmt.Errorf("not a valid target, error while parsing as url: %s", err)
	}

	switch url.Scheme {
	case "http", "https":
	default:
		return fmt.Errorf("not a valid target, requires http(s) url got %s", url.Scheme)
	}

	return nil
}

func validateHTTPPayload(payloads []config.CheckComponent) error {
	for _, payload := range payloads {
		switch payload.GetType() {
		case "header", "parameter":
			if err := validateKVPair(payload.GetValue()); err != nil {
				return fmt.Errorf("payload (%s) is not valid: %s", payload.GetValue(), err)
			}
		default:
			return fmt.Errorf("payload type %s is not valid", payload.GetType())
		}
	}
	return nil
}

type validationFunction func(string) error

// HTTPChecker represents a HTTP check we can employ for a given target
type HTTPChecker struct {
	Method string

	URL     string
	Payload map[string]string
	Headers map[string]string

	HTTPOutput *TVComponent

	Timeout time.Duration
}

// Type returns "HTTP" for a HTTPChecker.
func (c *HTTPChecker) Type() string {
	return string(HTTPInputType)
}

// ExecuteCheck runs the check for given HTTPChecker.
func (c *HTTPChecker) ExecuteCheck(ctx context.Context) (controller.ControllerFunctionResult, error) {
	prober := probes.NewHTTPProber()

	result, err := prober.Probe(c.Method, c.URL, c.Headers, c.Payload, c.Timeout)
	if err != nil {
		return nil, fmt.Errorf("HTTP Probe error: %s", err)
	}

	checkSuccessful := false

	switch c.HTTPOutput.Type {
	case "status_code":
		reqStatusCode, _ := strconv.Atoi(c.HTTPOutput.Value)
		if result.StatusCode == reqStatusCode {
			checkSuccessful = true
		}

	case "body":
		buf := new(bytes.Buffer)
		buf.ReadFrom(result.Body)
		if c.HTTPOutput.Value == buf.String() {
			checkSuccessful = true
		}

	case "header":
		kv := strings.SplitN(c.HTTPOutput.Value, headerDelimeter, 2)

		if kv[1] == result.Headers.Get(kv[0]) {
			checkSuccessful = true
		}
	}

	return CheckStats{
		Successful: checkSuccessful,

		StartTime: result.StartTime,
		Duration:  result.Duration,
	}, nil
}
