package configfile

import (
	"time"

	"github.com/sdslabs/pinger/pkg/config"
)

// Agent represents the configuration for an agent.
type Agent struct {
	Standalone bool                   `mapstructure:"standalone" json:"standalone"`
	Port       uint16                 `mapstructure:"port" json:"port"`
	Metrics    config.MetricsProvider `mapstructure:"metrics" json:"metrics"`
	Alerts     []config.AlertProvider `mapstructure:"alerts" json:"alerts"`
	Interval   time.Duration          `mapstructure:"interval" json:"interval"`
	Checks     []config.Check         `mapstructure:"checks" json:"checks"`
}
