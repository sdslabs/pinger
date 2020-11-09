package icmp

import (
	"context"
	"fmt"
	"net"

	"github.com/sdslabs/pinger/pkg/checker"
)

// checkerName is the name of the checker.
const checkerName = "ICMP"

func init() {
	checker.Register(checkerName, func() checker.Checker { return new(Checker) })
}

// Checker sends an ICMP ECHO request and checks if a reply is returned.
type Checker struct {
	prober *Prober
}

// Validate validates the check configuration.
func (c *Checker) Validate(check checker.Check) error {
	if check.GetTimeout() <= 0 {
		return fmt.Errorf("timeout should be > 0")
	}

	validateInputMap := validationMap{checkerName: validateInput}
	if err := checker.ValidateComponent(check.GetInput(), validateInputMap); err != nil {
		return fmt.Errorf("input: %w", err)
	}

	validateOutputMap := validationMap{"TIMEOUT": validateNil}
	if err := checker.ValidateComponent(check.GetOutput(), validateOutputMap); err != nil {
		return fmt.Errorf("output: %w", err)
	}

	validateTargetMap := validationMap{"ADDRESS": validateTarget}
	if err := checker.ValidateComponent(check.GetTarget(), validateTargetMap); err != nil {
		return fmt.Errorf("target: %w", err)
	}

	return nil
}

// Provision initializes required fields for c's execution.
func (c *Checker) Provision(check checker.Check) (err error) {
	addr := check.GetTarget().GetValue()
	timeout := check.GetTimeout()

	c.prober, err = NewProber(addr, timeout)
	return
}

// Execute executes the check.
func (c *Checker) Execute(ctx context.Context) (*checker.Result, error) {
	probeResult, err := c.prober.Probe(ctx)
	if err != nil {
		return nil, err
	}

	result := &checker.Result{
		Timeout:   true,
		StartTime: probeResult.StartTime,
		Duration:  probeResult.Duration,
	}

	if !probeResult.Timeout {
		result.Successful = true
		result.Timeout = false
	}

	return result, nil
}

// validationMap is an alias of map used for validating components.
type validationMap = map[string]func(string) error

// validateNil doesn't validate anything.
func validateNil(string) error { return nil }

// validateInput validates the check input.
func validateInput(val string) error {
	switch val {
	case "", "PING", "ECHO":
	default:
		return fmt.Errorf("invalid value: %s", val)
	}

	return nil
}

// validateTarget validates the target for a check.
func validateTarget(val string) error {
	if !isValidAddress(val) {
		return fmt.Errorf("target is not a valid address: %s", val)
	}

	return nil
}

// isValidAddress tells if s is a valid address or not.
func isValidAddress(s string) bool {
	// an IP address is a valid address.
	if net.ParseIP(s) != nil {
		return true
	}

	return checker.AddressRegex.MatchString(s)
}

// Interface guard.
var _ checker.Checker = (*Checker)(nil)
