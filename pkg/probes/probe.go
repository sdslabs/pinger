package probes

import (
	"time"
)

type ProbeResult interface {
	GetDuration() time.Duration
}
