package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// AgentConfig is the configuration structure for agent, mostly used in the case of
// standalone status page agent.
type AgentConfig struct {
	Port int `yaml:"port"`

	PrometheusMetrics     bool `yaml:"prometheus_metrics"`
	PrometheusMetricsPort int  `yaml:"prometheus_metrics_port"`

	Checks []*CheckConfig `yaml:"checks"`
}

func (a *AgentConfig) Parse(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, a)
}

func NewAgentConfig(path string) (*AgentConfig, error) {
	c := &AgentConfig{}
	err := c.Parse(path)
	if err != nil {
		return &AgentConfig{}, err
	}

	return c, nil
}
