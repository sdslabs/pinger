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

// NewICMPChecker creates a checker for ICMP requests.
func NewICMPChecker(agentCheck config.Check) (*ICMPChecker, error) {
	err := validateICMPCheck(agentCheck)
	if err != nil {
		return nil, fmt.Errorf("VALIDATION_ERROR: %s", err.Error())
	}

	address := agentCheck.GetTarget().GetValue()
	timeout := time.Duration(agentCheck.GetTimeout())
	if timeout <= 0 {
		timeout = defaults.ICMPProbeTimeout
	}

	return &ICMPChecker{
		Address: address,
		Timeout: timeout,

		ICMPOutput: &TVComponent{
			Type:  agentCheck.GetOutput().GetType(),
			Value: agentCheck.GetOutput().GetValue(),
		},
	}, nil
}

func validateICMPCheck(agentCheck config.Check) error {
	if err := validateICMPInput(agentCheck.GetInput()); err != nil {
		return err
	}

	if err := validateICMPOutput(agentCheck.GetOutput()); err != nil {
		return err
	}

	if err := validateICMPTarget(agentCheck.GetTarget()); err != nil {
		return err
	}

	return validateICMPPayload(agentCheck.GetPayloads())
}

func validateICMPInput(input config.Component) error {
	val := input.GetValue()
	if val != "PING" && val != "ECHO" && val != "" { // all of these mean the same
		return fmt.Errorf("for ICMP input provided method (%s) is not supported", val)
	}
	return nil
}

func validateICMPOutput(output config.Component) error {
	typ := output.GetType()
	if typ != "timeout" {
		return fmt.Errorf("provided output type (%s) is not supported", typ)
	}
	// for ICMP Echo request, we can only check if output is timeout or not, so we don't
	// have to validate the value of output.
	return nil
}

func validateICMPTarget(target config.Component) error {
	// We don't check the value of type for the target here
	// as for ICMP Check the target is always a address and we can check it that way only.
	// Just for consistency of types not being nil, we check if it's equal to "address"
	typ := target.GetType()
	if typ != "address" {
		return fmt.Errorf("target type %s is not supported", typ)
	}

	if _, err := net.ResolveIPAddr("ip", target.GetValue()); err != nil {
		return err
	}

	return nil
}

func validateICMPPayload(payload []config.Component) error {
	return nil // no payload required for ICMP check
}

// ICMPChecker represents an ICMP check we can deploy for a given target.
type ICMPChecker struct {
	Address string
	Timeout time.Duration

	ICMPOutput *TVComponent
}

// Type returns the type of checker, i.e., ICMP here.
func (c *ICMPChecker) Type() string {
	return string(ICMPInputType)
}

// ExecuteCheck starts the check with ICMP probe and validates if the given output is desired.
func (c *ICMPChecker) ExecuteCheck(ctx context.Context) (controller.FunctionResult, error) {
	prober, err := probes.NewICMPProbe(c.Address, c.Timeout)
	if err != nil {
		return nil, fmt.Errorf("ICMP probe error: %s", err.Error())
	}

	res, err := prober.Probe()
	if err != nil {
		return nil, fmt.Errorf("ICMP probe error: %s", err.Error())
	}

	return Stats{
		Successful: !res.Timeout,
		Timeout:    res.Timeout,
		StartTime:  res.StartTime,
		Duration:   res.Duration,
	}, nil
}
