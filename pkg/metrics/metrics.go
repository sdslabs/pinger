package metrics

import (
	"time"
)

// ProviderConfig is the config for metrics database provider.
type ProviderConfig struct {
	PType ProviderType

	Host string
	Port int

	Username string
	Password string

	Interval time.Duration
}

// ProviderType is the type of database.
type ProviderType string

// Various metrics providers.
var (
	PrometheusProviderType ProviderType = "Prometheus"
	TimeScaleProviderType  ProviderType = "TimeScale"
	EmptyProviderType      ProviderType = "None"
)
