package config

import (
	"github.com/sdslabs/status/pkg/metrics"
)

// AgentConfig is the configuration structure for agent.
type AgentConfig struct {
	Standalone bool                   `mapstructure:"standalone" json:"standalone" yaml:"standalone" toml:"standalone"`
	Port       int                    `mapstructure:"port" json:"port" yaml:"port" toml:"port"`
	Metrics    metrics.ProviderConfig `mapstructure:"metrics" json:"metrics" yaml:"metrics" toml:"metrics"`
	Checks     []*CheckConf           `mapstructure:"checks" json:"checks" yaml:"checks" toml:"checks"`
}
