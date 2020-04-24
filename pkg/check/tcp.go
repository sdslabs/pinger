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

// NewTCPChecker creates a checker for TCP requests.
func NewTCPChecker(agentCheck config.Check) (*TCPChecker, error) {
	if err := validateTCPCheck(agentCheck); err != nil {
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
		timeout = defaults.TCPProbeTimeout
	}

	return &TCPChecker{
		Address: address,
		Timeout: timeout,
		Message: message,
		TCPOutput: &TVComponent{
			Type:  agentCheck.GetOutput().GetType(),
			Value: agentCheck.GetOutput().GetValue(),
		},
	}, nil
}

func validateTCPCheck(agentCheck config.Check) error {
	if err := validateTCPInput(agentCheck.GetInput()); err != nil {
		return err
	}

	if err := validateTCPOutput(agentCheck.GetOutput()); err != nil {
		return err
	}

	if err := validateTCPTarget(agentCheck.GetTarget()); err != nil {
		return err
	}

	if agentCheck.GetOutput().GetType() == keyMessage && len(agentCheck.GetPayloads()) == 0 {
		return fmt.Errorf("no payload for expected message output")
	}

	return validateTCPPayload(agentCheck.GetPayloads())
}

func validateTCPInput(input config.Component) error {
	val := input.GetValue()
	if val != "ECHO" && val != "" {
		return fmt.Errorf("for TCP input provided: (%s), is not supported", val)
	}

	return nil
}

func validateTCPOutput(output config.Component) error {
	typ := output.GetType()
	if typ != keyMessage && typ != keyTimeout {
		return fmt.Errorf("provided output type (%s) is not supported", typ)
	}
	// Validates the expected Response message on the TCP Connection on sending the
	// input message. The user may only want to check whether the connection times out
	// by providing an output of the type `timeout`.
	return nil
}

func validateTCPTarget(target config.Component) error {
	// We don't check the value of type for the target here,
	// as for TCP Check the target is always a address and we can check it that way only.
	// Just for consistency of types not being nil, we check if it's equal to "address".
	typ := target.GetType()
	if typ != keyAddress {
		return fmt.Errorf("target type %s is not supported", typ)
	}

	if _, err := net.ResolveTCPAddr("tcp", target.GetValue()); err != nil {
		return err
	}

	return nil
}

func validateTCPPayload(payloads []config.Component) error {
	for _, payload := range payloads {
		typ := payload.GetType()
		if typ != keyMessage {
			return fmt.Errorf("for TCP payload provided: (%s), is not supported", typ)
		}
	}

	// Checks whether each payload provided by the user is of the type message.
	// We combine all the payloads in a single message along with the main input
	// and it to the server on TCP Connection.
	return nil
}

// TCPChecker represents an TCP check we can deploy for a given target.
type TCPChecker struct {
	Address   string
	Message   string
	Timeout   time.Duration
	TCPOutput *TVComponent
}

// Type returns the type of checker, i.e., TCP here.
func (c *TCPChecker) Type() string {
	return string(TCPInputType)
}

// ExecuteCheck starts the check with TCP probe and validates if the given output is desired.
func (c *TCPChecker) ExecuteCheck(ctx context.Context) (controller.FunctionResult, error) {
	prober, err := probes.NewTCPProbe(c.Address, c.Message, c.Timeout)
	if err != nil {
		return Stats{Successful: false}, nil
	}

	res, err := prober.Probe()
	if err != nil {
		return Stats{Successful: false}, nil
	}

	checkSuccessful := false

	if !res.Timeout {
		switch c.TCPOutput.Type {
		case keyMessage:
			checkSuccessful = res.Response == c.TCPOutput.Value
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
