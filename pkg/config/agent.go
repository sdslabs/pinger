package config

import (
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// AgentConfig is the configuration structure for agent, mostly used in the case of
// standalone status page agent.
type AgentConfig struct {
	Port int `yaml:"port"`

	PrometheusMetrics     bool `yaml:"prometheus_metrics"`
	PrometheusMetricsPort int  `yaml:"prometheus_metrics_port"`

	Checks []*Config `yaml:"checks"`
}

// Parse takes the path of agent config file and structurizes data into `AgentConfig`.
func (a *AgentConfig) Parse(path string) error {
	data, err := ioutil.ReadFile(filepath.Clean(path))
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, a)
}

// NewAgentConfig returns an `AgentConfig` from an agent config yaml file.
func NewAgentConfig(path string) (*AgentConfig, error) {
	c := &AgentConfig{}
	err := c.Parse(path)
	if err != nil {
		return &AgentConfig{}, err
	}

	return c, nil
}
