package check

import (
	"context"
	"fmt"

	"github.com/sdslabs/status/pkg/api/agent/proto"
	"github.com/sdslabs/status/pkg/defaults"
)

// Checker is the main Check interface we are going to use for the status page.
// Each check must implment this interface.
type Checker interface {
	ExecuteCheck(context.Context) error
}

type CheckComponent struct {
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
func NewChecker(agentCheck *proto.Check) (Checker, error) {
	err := validateCheck(agentCheck)
	if err != nil {
		return nil, fmt.Errorf("Error while validating check: %s", err)
	}

	switch InputType(agentCheck.Input.Type) {
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
func validateCheck(agentCheck *proto.Check) error {
	if agentCheck.Input == nil || agentCheck.Output == nil || agentCheck.Target == nil {
		return fmt.Errorf("Input, Ouput and Target are required for the check")
	}

	if agentCheck.Interval < int64(defaults.MinControllerRetryInterval.Seconds()) {
		return fmt.Errorf("interval provided is less than the minimum value allowed controller retry interval")
	}

	if agentCheck.Timeout < int64(defaults.MinControllerTimeout.Seconds()) {
		return fmt.Errorf("timeout for the check is less than the minimum value allowed for controller timeout")
	}

	switch InputType(agentCheck.Input.Type) {
	case HTTPInputType:
		return validateHTTPCheck(agentCheck)
	case TCPInputType:
	case WebsocketInputType:
	case ICMPInputType:
	default:
		return fmt.Errorf("provided input type is not valid: %s", agentCheck.Input.Type)
	}

	return fmt.Errorf("not a valid input type")
}
