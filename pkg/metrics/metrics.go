package metrics

import (
	"time"
)

type ProviderConfig struct {
	PType ProviderType

	Host string
	Port int

	Username string
	Password string

	Interval time.Duration
}

type ProviderType string

var (
	PrometheusProviderType ProviderType = "Prometheus"
	TimeScaleProviderType  ProviderType = "TimeScale"
	EmptyProviderType      ProviderType = "None"
)
