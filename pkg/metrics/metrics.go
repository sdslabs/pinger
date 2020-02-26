package metrics

import (
	"time"
)

// ProviderConfig is the config for metrics database provider.
type ProviderConfig struct {
	Backend  ProviderType `mapstructure:"backend" json:"backend" yaml:"backend" toml:"backend"`
	Host     string       `mapstructure:"host" json:"host" yaml:"host" toml:"host"`
	Port     int          `mapstructure:"port" json:"port" yaml:"port" toml:"port"`
	DBName   string       `mapstructure:"db_name" json:"db_name" yaml:"db_name" toml:"db_name"`
	Username string       `mapstructure:"username" json:"username" yaml:"username" toml:"username"`
	Password string       `mapstructure:"password" json:"password" yaml:"password" toml:"password"`
	SSLMode  bool         `mapstructure:"ssl_mode" json:"ssl_mode" yaml:"ssl_mode" toml:"ssl_mode"`
	// Interval after which the agent pushes / provider pulls checks
	Interval time.Duration `mapstructure:"interval" json:"interval" yaml:"interval" toml:"interval"`
}

// ProviderType is the type of database.
type ProviderType string

// Various metrics providers.
var (
	PrometheusProviderType ProviderType = "prometheus"
	TimeScaleProviderType  ProviderType = "timescale"
	EmptyProviderType      ProviderType = "none"
)
