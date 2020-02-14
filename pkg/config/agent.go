package config

import (
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"github.com/sdslabs/status/pkg/metrics"
)

// AgentConfig is the configuration structure for agent, mostly used in the case of
// standalone status page agent.
type AgentConfig struct {
	Metrics metrics.ProviderConfig `yaml:"metrics"`
	Checks  []*CheckConf           `yaml:"checks"`
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
