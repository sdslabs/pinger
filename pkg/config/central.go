package config

import "time"

type agent struct {
	Host    string        `mapstructure:"host" json:"host" yaml:"host" toml:"host"`
	Port    int           `mapstructure:"port" json:"port" yaml:"port" toml:"port"`
	Timeout time.Duration `mapstructure:"timeout" json:"timeout" yaml:"timeout" toml:"timeout"`
}

// CentralServerConfig is the configuration for central server.
type CentralServerConfig struct {
	Port   int      `mapstructure:"port" json:"port" yaml:"port" toml:"port"`
	Agents []*agent `mapstructure:"agents" json:"agents" yaml:"agents" toml:"agents"`
}
