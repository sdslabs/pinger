package check

import (
	"context"
	"fmt"
	"time"

	"github.com/sdslabs/status/pkg/config"
	"github.com/sdslabs/status/pkg/controller"
	"github.com/sdslabs/status/pkg/defaults"
	"github.com/sdslabs/status/pkg/probes"
)

// DNSChecker looks up the host, i.e., resolves the DNS name.
type DNSChecker struct {
	Host    string
	Timeout time.Duration
}

// NewDNSChecker creates a DNS resolver checker for the given conf.
func NewDNSChecker(agentCheck config.Check) (*DNSChecker, error) {
	if err := validateDNSCheck(agentCheck); err != nil {
		return nil, err
	}

	timeout := time.Duration(agentCheck.GetTimeout())
	if timeout <= 0 {
		timeout = defaults.DNSProbeTimeout
	}

	host := agentCheck.GetTarget().GetValue()

	return &DNSChecker{Host: host, Timeout: timeout}, nil
}

// Type returns the type of checker, i.e., DNS here.
func (dc *DNSChecker) Type() string {
	return string(DNSInputType)
}

// ExecuteCheck resolves the host to check if the host name is valid.
func (dc *DNSChecker) ExecuteCheck(context.Context) (controller.FunctionResult, error) {
	result := probes.ProbeDNS(dc.Host, dc.Timeout)

	return Stats{
		Successful: result.Successful,
		Timeout:    result.Timeout,
		StartTime:  result.StartTime,
		Duration:   result.Duration,
	}, nil
}

func validateDNSCheck(agentCheck config.Check) error {
	// input value can only be empty or "ECHO" (standard)
	inputValue := agentCheck.GetInput().GetValue()
	if inputValue != "" && inputValue != "ECHO" {
		return fmt.Errorf("provided input value (%s) is not supported", inputValue)
	}

	// target type can only be host
	targetType := agentCheck.GetTarget().GetType()
	if targetType != "host" {
		return fmt.Errorf("provided target type (%s) is not supported", targetType)
	}

	// For DNS checks all we need to know is that it resolves within the time limit
	// hence there is no point checking for output and there is also no extra payload.
	//
	// Since the check is basically to resolve the host, we don't validate the host,
	// i.e., target value. Rather the "ExecuteCheck" does that for us.

	return nil
}
