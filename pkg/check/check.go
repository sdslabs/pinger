package check

import (
	"context"
	"fmt"
	"time"

	"github.com/sdslabs/status/pkg/config"
	"github.com/sdslabs/status/pkg/controller"
)

// Types of inputs accepted for checks.
const (
	HTTPInputType      InputType = "HTTP"
	TCPInputType       InputType = "TCP"
	WebsocketInputType InputType = "Websocket"
	ICMPInputType      InputType = "ICMP"
	UDPInputType       InputType = "UDP"
)

// NewChecker returns a new Checker for the provided agent Check.
// This first validates the agent Check provided in the argument and if the Check
// is validated then returns a new Checker instance for the provided Check configuration.
func NewChecker(agentCheck config.Check) (Checker, error) {
	switch InputType(agentCheck.GetInput().GetType()) {
	case HTTPInputType:
		return NewHTTPChecker(agentCheck)
	case TCPInputType:
		return NewTCPChecker(agentCheck)
	case WebsocketInputType:
		return NewWSChecker(agentCheck)
	case ICMPInputType:
		return NewICMPChecker(agentCheck)
	case UDPInputType:
		return NewUDPChecker(agentCheck)
	default:
		return nil, fmt.Errorf(
			"invalid check input type: %s",
			agentCheck.GetInput().GetType(),
		)
	}
}

// InputType represents the kind of check.
// Can be HTTP, ICMP, Websocket etc.
type InputType string

// Checker is the main Check interface we are going to use for the status page.
// Each check must implement this interface.
type Checker interface {
	// ExecuteCheck executes the check for the provided checker.
	ExecuteCheck(context.Context) (controller.FunctionResult, error)

	// Type returns the type of the checker, this can be http, icmp, websockets etc.
	Type() string
}

// TVComponent is a key-value pair component with it's key as the type of component.
// This is used to represent payloads or any other type-value pair.
//
// Example of header payload:
//     TVComponent{
//	       Type:  "header",
//         Value: "Authorization=Bearer xyz"
//     }
type TVComponent struct {
	Type  string
	Value string
}

// Stats is the struct with concerned statistics collected by a check.
type Stats struct {
	Successful bool
	Timeout    bool
	StartTime  time.Time
	Duration   time.Duration
}

// GetDuration returns the duration taken by check to execute.
func (cd Stats) GetDuration() time.Duration {
	return cd.Duration
}

// GetStartTime returns the time when the check started to execute.
func (cd Stats) GetStartTime() time.Time {
	return cd.StartTime
}

// IsSuccessful returns true the check executed with desired output
// or false when the check failed to produce the required output.
func (cd Stats) IsSuccessful() bool {
	return cd.Successful
}

// IsTimeout returns true when the check failed due to timeout or not.
func (cd Stats) IsTimeout() bool {
	return cd.Timeout
}
