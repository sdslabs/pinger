package dns

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/sdslabs/pinger/pkg/checker"
)

// Prober probes and resolves the host.
type Prober struct {
	host    string
	timeout time.Duration
}

// NewProber creates a prober to resolve the host.
func NewProber(host string, timeout time.Duration) (*Prober, error) {
	if timeout <= 0 {
		return nil, fmt.Errorf("timeout should be > 0")
	}

	return &Prober{
		host:    host,
		timeout: timeout,
	}, nil
}

// Probe probes to resolves the host.
func (p *Prober) Probe(ctx context.Context) (*ProbeResult, error) {
	startTime := time.Now()

	baseCtx := ctx
	if ctx == nil {
		baseCtx = context.Background()
	}

	probeCtx, cancel := context.WithTimeout(baseCtx, p.timeout)
	defer cancel()

	resolver := &net.Resolver{}
	addrs, err := resolver.LookupHost(probeCtx, p.host)
	if err != nil {
		if checker.ErrIsTimeout(err) {
			return &ProbeResult{
				StartTime: startTime,
				Duration:  p.timeout,
				Timeout:   true,
			}, nil
		}

		return nil, err
	}

	return &ProbeResult{
		StartTime:  startTime,
		Duration:   time.Since(startTime),
		Timeout:    false,
		ResolvedTo: addrs,
	}, nil
}

// ProbeResult is the result of a DNS Probe.
type ProbeResult struct {
	Timeout    bool
	StartTime  time.Time
	Duration   time.Duration
	ResolvedTo []string
}
