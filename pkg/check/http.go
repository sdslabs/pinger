package check

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/sdslabs/status/pkg/api/agent/proto"
	"github.com/sdslabs/status/pkg/defaults"

	log "github.com/sirupsen/logrus"
)

const HEADER_DELIMITER = "="

// HTTPChecker represents a HTTP check we can employ for a given target
type HTTPChecker struct {
	Method string

	URL     string
	Payload map[string]string
	Headers map[string]string

	HTTPOutput *CheckComponent
}

func NewHTTPChecker(agentCheck *proto.Check) (*HTTPChecker, error) {
	err := validateHTTPCheck(agentCheck)
	if err != nil {
		return nil, fmt.Errorf("VALIDATION_ERROR: %s", err)
	}

	method := agentCheck.Input.Value
	if method == "" {
		method = defaults.DefaultHTTPMethod
	}

	params := make(map[string]string)
	headers := make(map[string]string)

	for _, payload := range agentCheck.Payloads {
		kv := strings.SplitN(payload.Payload, HEADER_DELIMITER, 2)

		switch payload.PayloadType {
		case "header":
			headers[kv[0]] = kv[1]
		case "parameter":
			params[kv[0]] = kv[1]
		}
	}

	return &HTTPChecker{
		Method: method,

		URL:     agentCheck.Target.Value,
		Payload: params,
		Headers: headers,

		HTTPOutput: &CheckComponent{
			Type:  agentCheck.Output.Type,
			Value: agentCheck.Output.Value,
		},
	}, nil
}

func (c *HTTPChecker) ExecuteCheck(ctx context.Context) error {
	log.Debug("Executing HTTP check.")
	return nil
}

type validationFunction func(string) error

var validHttpOutputTypes map[string]validationFunction = map[string]validationFunction{
	"status_code": validateStatusCode,
	"body":        validateBody,
	"header":      validateKVPair,
}

func validateStatusCode(val string) error {
	intVal, err := strconv.Atoi(val)
	if err == nil {
		return fmt.Errorf("status code prvoided is parsable: %s", err)
	}
	if intVal > 599 || intVal < 100 {
		return fmt.Errorf("status code is not valid(expected between: 100 - 599) got %s")
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

func validateHTTPCheck(agentCheck *proto.Check) error {
	err := validateHTTPInput(agentCheck.Input)
	if err != nil {
		return err
	}

	err = validateHTTPOutput(agentCheck.Output)
	if err != nil {
		return err
	}

	err = validateHTTPTarget(agentCheck.Target)
	if err != nil {
		return err
	}

	return validateHTTPPayload(agentCheck.Payloads)
}

func validateHTTPInput(input *proto.Check_Component) error {
	inputVal := input.Value
	if inputVal != "GET" || inputVal != "POST" || inputVal != "" {
		return fmt.Errorf("for HTTP input the provided method(%s) is not supported", inputVal)
	}

	return nil
}

func validateHTTPOutput(output *proto.Check_Component) error {
	validateFunc, ok := validHttpOutputTypes[output.Type]
	if !ok {
		return fmt.Errorf("provided Output Type(%s) is not valid for HTTP input", output.Type)
	}

	return validateFunc(output.Value)
}

func validateHTTPTarget(target *proto.Check_Component) error {
	// We don't check the value of type for the target here
	// as for HTTP Check the target is always a URL and we check it that way only.
	_, err := url.Parse(target.Value)
	if err != nil {
		return fmt.Errorf("not a valid target, error while parsing as url: %s", err)
	}

	return nil
}

func validateHTTPPayload(payloads []*proto.Check_Payloads) error {
	for _, payload := range payloads {
		if payload.PayloadType == "header" || payload.PayloadType == "parameter" {
			err := validateKVPair(payload.Payload)

			if err != nil {
				return fmt.Errorf("payload (%s) is not valid: %s", payload.Payload, err)
			}
		}

		return fmt.Errorf("payload type %s is not valid", payload.PayloadType)
	}

	return nil
}
