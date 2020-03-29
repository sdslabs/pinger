package metrics

import (
	"context"
	"time"

	"github.com/sdslabs/status/pkg/controller"
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
	LogProviderType        ProviderType = "log"
)

type controllerFunc = func(context.Context) (controller.FunctionResult, error)

// MetricsFunctionResult implements controller.FunctionResult to create controllers
// for fetching metrics.
type MetricsFunctionResult struct {
	Duration  time.Duration
	StartTime time.Time
	Success   bool
	Timeout   bool
}

func (m MetricsFunctionResult) GetDuration() time.Duration {
	return m.Duration
}

func (m MetricsFunctionResult) GetStartTime() time.Time {
	return m.StartTime
}

func (m MetricsFunctionResult) IsSuccessful() bool {
	return m.Success
}

func (m MetricsFunctionResult) IsTimeout() bool {
	return m.Timeout
}
