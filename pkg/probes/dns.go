package probes

import (
	"context"
	"net"
	"time"
)

// DNSProbeResult is the result of a DNS Probe.
type DNSProbeResult struct {
	Successful bool
	Timeout    bool
	StartTime  time.Time
	Duration   time.Duration
}

// ProbeDNS tries to resolve the host.
func ProbeDNS(host string, timeout time.Duration) *DNSProbeResult {
	startTime := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	result := &DNSProbeResult{
		StartTime:  startTime,
		Timeout:    false,
		Successful: false,
	}

	if _, err := net.DefaultResolver.LookupHost(ctx, host); err != nil {
		if errIsTimeout(err) {
			result.Timeout = true
			result.Duration = timeout
			return result
		}

		result.Duration = time.Since(startTime)
		return result
	}

	result.Duration = time.Since(startTime)
	result.Successful = true
	return result
}
