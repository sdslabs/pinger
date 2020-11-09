package dns

import (
	"context"
	"fmt"
	"net"

	"github.com/sdslabs/pinger/pkg/checker"
)

// checkerName is the name of the checker.
const checkerName = "DNS"

func init() {
	checker.Register(checkerName, func() checker.Checker { return new(Checker) })
}

// Checker runs a DNS check on the given host.
type Checker struct {
	prober *Prober

	outputType  string
	outputValue string
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

	validateOutputMap := validationMap{
		"TIMEOUT": validateNil,
		"IP":      validateOutputAddress,
		"ADDRESS": validateOutputAddress,
	}
	if err := checker.ValidateComponent(check.GetOutput(), validateOutputMap); err != nil {
		return fmt.Errorf("output: %w", err)
	}

	validateTargetMap := validationMap{
		"HOST":     validateTarget,
		"HOSTNAME": validateTarget,
		"DNSNAME":  validateTarget,
	}
	if err := checker.ValidateComponent(check.GetTarget(), validateTargetMap); err != nil {
		return fmt.Errorf("target: %w", err)
	}

	return nil
}

// Provision initializes required fields for c's execution.
func (c *Checker) Provision(check checker.Check) (err error) {
	host := check.GetTarget().GetValue()
	timeout := check.GetTimeout()

	c.prober, err = NewProber(host, timeout)
	if err != nil {
		return
	}

	c.outputType = check.GetOutput().GetType()
	c.outputValue = check.GetOutput().GetValue()
	return nil
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

	if probeResult.Timeout {
		return result, nil
	}

	result.Timeout = false

	if c.outputType == "TIMEOUT" {
		// we need to check if the address did resolve to atleast one of IP address.
		if len(probeResult.ResolvedTo) > 0 {
			result.Successful = true
		}

		return result, nil
	}

	// the only other type of output type is "IP" (or "ADDRESS")
	for _, addr := range probeResult.ResolvedTo {
		if addr == c.outputValue {
			result.Successful = true
			return result, nil
		}
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

// validateOutputAddress validates output for timeout type.
func validateOutputAddress(val string) error {
	// output value should be a valid IPv4 or IPv6 address.
	ip := net.ParseIP(val)
	if ip == nil {
		return fmt.Errorf("value is not a valid IP: %s", val)
	}

	return nil
}

// validateTarget validates the target for a check.
func validateTarget(val string) error {
	if !isDNSName(val) {
		return fmt.Errorf("value is not a valid DNS name: %s", val)
	}

	return nil
}

// isDNSName checks if the string is a valid DNS name.
func isDNSName(s string) bool {
	// dns name shouldn't be an IP address
	if net.ParseIP(s) != nil {
		return false
	}

	return checker.AddressRegex.MatchString(s)
}

// Interface guard.
var _ checker.Checker = (*Checker)(nil)
