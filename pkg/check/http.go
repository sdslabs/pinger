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

	log "github.com/sirupsen/logrus"
)

const HEADER_DELIMITER = "="

// HTTPChecker represents a HTTP check we can employ for a given target
type HTTPChecker struct {
	Method string

	URL     string
	Payload map[string]string
	Headers map[string]string

	HTTPOutput *TVComponent

	Timeout time.Duration
}

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
		kv := strings.SplitN(payload.GetValue(), HEADER_DELIMITER, 2)

		switch payload.GetType() {
		case "header":
			headers[kv[0]] = kv[1]
		case "parameter":
			params[kv[0]] = kv[1]
		}
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

		Timeout: defaults.DefaultHTTPProbeTimeout,
	}, nil
}

func (c *HTTPChecker) Type() string {
	return "http"
}

func (c *HTTPChecker) ExecuteCheck(ctx context.Context) (controller.ControllerFunctionResult, error) {
	log.Debug("Executing HTTP check.")
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
			log.Info("Check successful")
			checkSuccessful = true
		} else {
			log.Warnf("Check unsuccessful, status(req %d) does not match with %d", reqStatusCode, result.StatusCode)
		}

	case "body":
		buf := new(bytes.Buffer)
		buf.ReadFrom(result.Body)
		if c.HTTPOutput.Value == buf.String() {
			log.Info("Check Successful")
			checkSuccessful = true
		} else {
			log.Warnf("Check Unsuccessful")
		}

	case "header":
		kv := strings.SplitN(c.HTTPOutput.Value, HEADER_DELIMITER, 2)

		if kv[1] == result.Headers.Get(kv[0]) {
			log.Info("Check Successful")
			checkSuccessful = true
		} else {
			log.Warn("Check Unsuccessful")
		}
	}

	return CheckStats{
		Successful: checkSuccessful,

		StartTime: result.StartTime,
		Duration:  result.Duration,
	}, nil
}

type validationFunction func(string) error

var validHttpOutputTypes map[string]validationFunction = map[string]validationFunction{
	"status_code": validateStatusCode,
	"body":        validateBody,
	"header":      validateKVPair,
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
	kv := strings.SplitN(val, HEADER_DELIMITER, 2)
	if len(kv) != 2 {
		return fmt.Errorf("header value is not valid, must have format: HEADER=<Header value>")
	}

	return nil
}

func validateHTTPCheck(agentCheck config.Check) error {
	err := validateHTTPInput(agentCheck.GetInput())
	if err != nil {
		return err
	}

	err = validateHTTPOutput(agentCheck.GetOutput())
	if err != nil {
		return err
	}

	err = validateHTTPTarget(agentCheck.GetTarget())
	if err != nil {
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
	// as for HTTP Check the target is always a URL and we check it that way only.
	_, err := url.Parse(target.GetValue())
	if err != nil {
		return fmt.Errorf("not a valid target, error while parsing as url: %s", err)
	}

	return nil
}

func validateHTTPPayload(payloads []config.CheckComponent) error {
	for _, payload := range payloads {
		if payload.GetType() == "header" && payload.GetType() == "parameter" {
			err := validateKVPair(payload.GetValue())

			if err != nil {
				return fmt.Errorf("payload (%s) is not valid: %s", payload.GetValue(), err)
			}
		}

		return fmt.Errorf("payload type %s is not valid", payload.GetType())
	}

	return nil
}
