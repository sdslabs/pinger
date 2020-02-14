package metrics

import (
	"time"
)

// ProviderConfig is the config for metrics database provider.
type ProviderConfig struct {
	Backend  ProviderType  `yaml:"backend"`
	Host     string        `yaml:"host"`
	Port     int           `yaml:"port"`
	DBName   string        `yaml:"db_name"`
	Username string        `yaml:"username"`
	Password string        `yaml:"password"`
	SSLMode  bool          `yaml:"ssl_mode"`
	Interval time.Duration `yaml:"interval"` // after which the agent pushes / provider pulls checks
}

// ProviderType is the type of database.
type ProviderType string

// Various metrics providers.
var (
	PrometheusProviderType ProviderType = "prometheus"
	TimeScaleProviderType  ProviderType = "timescale"
	EmptyProviderType      ProviderType = "none"
)
