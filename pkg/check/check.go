package check

import (
	"context"
	"fmt"
	"time"

	"github.com/sdslabs/status/pkg/config"
	"github.com/sdslabs/status/pkg/controller"
	"github.com/sdslabs/status/pkg/defaults"
)

// Checker is the main Check interface we are going to use for the status page.
// Each check must implment this interface.
type Checker interface {
	// This function executes the check for the provided checker.
	ExecuteCheck(context.Context) (controller.ControllerFunctionResult, error)

	// Returns the type of the checker, this can be http, icmp, websockets etc.
	Type() string
}

type CheckStats struct {
	Successful bool

	StartTime time.Time
	Duration  time.Duration
}

func (cd CheckStats) GetDuration() time.Duration {
	return cd.Duration
}

func (cd CheckStats) GetStartTime() time.Time {
	return cd.StartTime
}

func (cd CheckStats) IsSuccessful() bool {
	return cd.Successful
}

type TVComponent struct {
	Type  string
	Value string
}

type InputType string

const (
	HTTPInputType      InputType = "HTTP"
	TCPInputType       InputType = "TCP"
	WebsocketInputType InputType = "Websocket"
	ICMPInputType      InputType = "ICMP"
)

// NewChecker returns a new Checker for the provided agent Check.
// This first validates the agent Check provided in the argument and if the Check
// is validated then returns a new Checker instance for the provided Check configuration.
func NewChecker(agentCheck config.Check) (Checker, error) {
	err := validateCheck(agentCheck)
	if err != nil {
		return nil, fmt.Errorf("Error while validating check: %s", err)
	}

	switch InputType(agentCheck.GetInput().GetType()) {
	case HTTPInputType:
		return NewHTTPChecker(agentCheck)
	case TCPInputType:
	case WebsocketInputType:
	case ICMPInputType:
	}

	return nil, nil
}

// Validates the agent check provided in the argument, returns an error
// if the provided check is not valid.
func validateCheck(agentCheck config.Check) error {
	if agentCheck.GetInput() == nil || agentCheck.GetOutput() == nil || agentCheck.GetTarget() == nil {
		return fmt.Errorf("Input, Ouput and Target are required for the check")
	}

	if agentCheck.GetInterval() < int64(defaults.MinControllerRetryInterval.Seconds()) {
		return fmt.Errorf("interval provided is less than the minimum value allowed controller retry interval")
	}

	if agentCheck.GetTimeout() < int64(defaults.MinControllerTimeout.Seconds()) {
		return fmt.Errorf("timeout for the check is less than the minimum value allowed for controller timeout")
	}

	switch InputType(agentCheck.GetInput().GetType()) {
	case HTTPInputType:
		return validateHTTPCheck(agentCheck)
	case TCPInputType:
	case WebsocketInputType:
	case ICMPInputType:
	default:
		return fmt.Errorf("provided input type is not valid: %s", agentCheck.GetInput().GetType())
	}

	return fmt.Errorf("not a valid input type")
}
