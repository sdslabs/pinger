package check

import (
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

var validWSOutputTypes map[string]validationFunction = map[string]validationFunction{
	keyStatusCode: validateStatusCode,
	keyMessage:    validateMessages,
	keyHeader:     validateKVPair,
	keyTimeout:    func(string) error { return nil },
}

// WSChecker represents a Websocket check we can deploy for a given target.
type WSChecker struct {
	URL      string
	Headers  map[string]string
	Messages []string

	Output *TVComponent

	Timeout time.Duration
}

// Type returns "Websocket" for a WSChecker.
func (c *WSChecker) Type() string {
	return string(WebsocketInputType)
}

// ExecuteCheck runs the check for given WSChecker.
func (c *WSChecker) ExecuteCheck(ctx context.Context) (controller.FunctionResult, error) {
	prober, err := probes.NewWSProber(c.URL, c.Messages, c.Headers, c.Timeout)
	if err != nil {
		return nil, err
	}

	result, err := prober.Probe()
	if err != nil {
		return nil, err
	}

	checkSuccessful := false

	if !result.Timeout {
		switch c.Output.Type {
		case keyStatusCode:
			var reqStatusCode int
			reqStatusCode, err = strconv.Atoi(c.Output.Value)
			if err != nil {
				return nil, err
			}
			if result.StatusCode == reqStatusCode {
				checkSuccessful = true
			}

		case keyMessage:
			// Messages in output are separated by `"\n---\n"`, i.e., if the output is expected to
			// have the messages 'hello world' and 'bye world', the output value is supposed to be
			// of the form:
			//     hello world
			//     ---
			//     bye world
			messages := strings.Split(c.Output.Value, "\n---\n")
			if len(messages) == len(result.Messages) {
				messagesSame := true
				for i := 0; i < len(messages); i++ {
					if messages[i] != result.Messages[i] {
						messagesSame = false
						break
					}
				}
				checkSuccessful = messagesSame
			}

		case keyHeader:
			kv := strings.SplitN(c.Output.Value, splitDelimeter, 2)

			if kv[1] == result.Headers.Get(kv[0]) {
				checkSuccessful = true
			}
		case keyTimeout:
			checkSuccessful = true
		}
	}

	return Stats{
		Successful: checkSuccessful,
		Timeout:    result.Timeout,
		StartTime:  result.StartTime,
		Duration:   result.Duration,
	}, nil
}

// NewWSChecker returns a new websocket cheker from the check config.
func NewWSChecker(agentCheck config.Check) (*WSChecker, error) {
	if err := validateWSCheck(agentCheck); err != nil {
		return nil, fmt.Errorf("VALIDATION_ERROR: %s", err)
	}

	messages := []string{}
	headers := map[string]string{}

	for _, payload := range agentCheck.GetPayloads() {
		switch payload.GetType() {
		case keyHeader:
			kv := strings.SplitN(payload.GetValue(), splitDelimeter, 2)
			headers[kv[0]] = kv[1]
		case keyMessage:
			messages = append(messages, payload.GetValue())
		}
	}

	timeout := time.Duration(agentCheck.GetTimeout())
	if timeout <= 0 {
		timeout = defaults.WSProbeTimeout
	}

	return &WSChecker{
		URL:      agentCheck.GetTarget().GetValue(),
		Headers:  headers,
		Messages: messages,

		Output: &TVComponent{
			Type:  agentCheck.GetOutput().GetType(),
			Value: agentCheck.GetOutput().GetValue(),
		},

		Timeout: timeout,
	}, nil
}

func validateWSCheck(agentCheck config.Check) error {
	if err := validateWSInput(agentCheck.GetInput()); err != nil {
		return err
	}

	if err := validateWSOutput(agentCheck.GetOutput()); err != nil {
		return err
	}

	if err := validateWSTarget(agentCheck.GetTarget()); err != nil {
		return err
	}

	return validateWSPayload(agentCheck.GetPayloads())
}

func validateWSInput(input config.Component) error {
	if val := input.GetValue(); val != "PING" && val != "" {
		return fmt.Errorf("input value provided (%s) is not valid", val)
	}

	return nil
}

func validateWSOutput(output config.Component) error {
	validateFunc, ok := validWSOutputTypes[output.GetType()]
	if !ok {
		return fmt.Errorf("provided output type (%s) is not valid for WS input", output.GetType())
	}

	return validateFunc(output.GetValue())
}

func validateWSTarget(target config.Component) error {
	// We don't check the value of type for the target here
	// as for WS Check the target is always a URL and we can check it that way only.
	// Just for consistency of types not being nil, we check if it's equal to "url"
	typ := target.GetType()
	if typ != keyURL {
		return fmt.Errorf("target type %s is not supported", typ)
	}
	u, err := url.Parse(target.GetValue())
	if err != nil {
		return fmt.Errorf("not a valid target, error while parsing as url: %s", err.Error())
	}

	switch u.Scheme {
	case "ws", "wss":
	default:
		return fmt.Errorf("not a valid target, requires ws(s) url got %s", u.Scheme)
	}

	return nil
}

func validateWSPayload(payloads []config.Component) error {
	for _, payload := range payloads {
		switch payload.GetType() {
		case keyHeader:
			if err := validateKVPair(payload.GetValue()); err != nil {
				return err
			}
		case keyMessage:
			if payload.GetValue() == "" {
				return fmt.Errorf("payload message cannot be empty")
			}
		default:
			return fmt.Errorf("payload type %s is not valid", payload.GetType())
		}
	}

	return nil
}
