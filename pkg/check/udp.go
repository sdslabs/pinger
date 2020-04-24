package check

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/sdslabs/status/pkg/config"
	"github.com/sdslabs/status/pkg/controller"
	"github.com/sdslabs/status/pkg/defaults"
	"github.com/sdslabs/status/pkg/probes"
)

// NewUDPChecker creates a checker for UDP requests.
func NewUDPChecker(agentCheck config.Check) (*UDPChecker, error) {
	if err := validateUDPCheck(agentCheck); err != nil {
		return nil, fmt.Errorf("VALIDATION_ERROR: %s", err.Error())
	}

	// Join all messages with "\n".
	message := ""
	for i, payload := range agentCheck.GetPayloads() {
		message += payload.GetValue()
		if i != len(agentCheck.GetPayloads())-1 {
			message += "\n"
		}
	}

	address := agentCheck.GetTarget().GetValue()
	timeout := time.Duration(agentCheck.GetTimeout())
	if timeout <= 0 {
		timeout = defaults.UDPProbeTimeout
	}

	return &UDPChecker{
		Address: address,
		Timeout: timeout,
		Message: message,
		UDPOutput: &TVComponent{
			Type:  agentCheck.GetOutput().GetType(),
			Value: agentCheck.GetOutput().GetValue(),
		},
	}, nil
}

func validateUDPCheck(agentCheck config.Check) error {
	if err := validateUDPInput(agentCheck.GetInput()); err != nil {
		return err
	}

	if err := validateUDPOutput(agentCheck.GetOutput()); err != nil {
		return err
	}

	if err := validateUDPTarget(agentCheck.GetTarget()); err != nil {
		return err
	}

	if len(agentCheck.GetPayloads()) == 0 {
		return fmt.Errorf("no payload for expected message output")
	}

	return validateUDPPayload(agentCheck.GetPayloads())
}

func validateUDPInput(input config.Component) error {
	val := input.GetValue()
	if val != "ECHO" && val != "" {
		return fmt.Errorf("for UDP input provided: (%s), is not supported", val)
	}

	return nil
}

func validateUDPOutput(output config.Component) error {
	typ := output.GetType()
	if typ != keyMessage && typ != keyTimeout {
		return fmt.Errorf("provided output type (%s) is not supported", typ)
	}
	// Validates the expected Response message on the UDP Connection on sending the
	// input message. The user may only want to check whether the connection times out
	// by providing an output of the type `timeout`.
	return nil
}

func validateUDPTarget(target config.Component) error {
	// We don't check the value of type for the target here,
	// as for UDP Check the target is always a address and we can check it that way only.
	// Just for consistency of types not being nil, we check if it's equal to "address".
	typ := target.GetType()
	if typ != keyAddress {
		return fmt.Errorf("target type %s is not supported", typ)
	}

	if _, err := net.ResolveUDPAddr("udp", target.GetValue()); err != nil {
		return err
	}

	return nil
}

func validateUDPPayload(payloads []config.Component) error {
	for _, payload := range payloads {
		typ := payload.GetType()
		if typ != keyMessage {
			return fmt.Errorf("for UDP payload provided: (%s), is not supported", typ)
		}
	}

	// Checks whether each payload provided by the user is of the type message.
	// We combine all the payloads in a single message along with the main input
	// and it to the server on UDP Connection.
	return nil
}

// UDPChecker represents an UDP check we can deploy for a given target.
type UDPChecker struct {
	Address   string
	Message   string
	Timeout   time.Duration
	UDPOutput *TVComponent
}

// Type returns the type of checker, i.e., UDP here.
func (c *UDPChecker) Type() string {
	return string(UDPInputType)
}

// ExecuteCheck starts the check with UDP probe and validates if the given output is desired.
func (c *UDPChecker) ExecuteCheck(ctx context.Context) (controller.FunctionResult, error) {
	prober, err := probes.NewUDPProbe(c.Address, c.Message, c.Timeout)
	if err != nil {
		return nil, fmt.Errorf("UDP probe error: %s", err.Error())
	}

	res, err := prober.Probe()
	if err != nil {
		return nil, fmt.Errorf("UDP probe error: %s", err.Error())
	}

	checkSuccessful := false

	if !res.Timeout {
		switch c.UDPOutput.Type {
		case keyMessage:
			checkSuccessful = res.Response == c.UDPOutput.Value
		case keyTimeout:
			checkSuccessful = true
		}
	}

	return Stats{
		Successful: checkSuccessful,
		Timeout:    res.Timeout,
		StartTime:  res.StartTime,
		Duration:   res.Duration,
	}, nil
}
