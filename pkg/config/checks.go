package config

import (
	"github.com/sdslabs/status/pkg/agent/proto"
)

// Check is the interface which every check that needs to be processed here
// should implement.
type Check interface {
	GetInput() Component
	GetOutput() Component
	GetTarget() Component

	GetPayloads() []Component
	GetInterval() int64
	GetTimeout() int64

	GetName() string
}

// Component is the Type Value component for check components like Input, Output, Target etc.
type Component interface {
	GetType() string
	GetValue() string
}

// CheckConf Associated with each check.
type CheckConf struct {
	Input  *ComponentConfig `yaml:"input"`
	Output *ComponentConfig `yaml:"output"`
	Target *ComponentConfig `yaml:"target"`

	Payloads []*ComponentConfig `yaml:"payloads"`

	Name     string `yaml:"name"`
	Timeout  int64  `yaml:"timeout"`
	Interval int64  `yaml:"interval"`
}

// GetInput returns the input of the check.
func (m *CheckConf) GetInput() Component {
	return m.Input
}

// GetTarget returns the target of the check.
func (m *CheckConf) GetTarget() Component {
	return m.Target
}

// GetOutput returns the output of the check.
func (m *CheckConf) GetOutput() Component {
	return m.Output
}

// GetPayloads returns the payloads of the check.
func (m *CheckConf) GetPayloads() []Component {
	payloads := make([]Component, len(m.Payloads))
	for i, payload := range m.Payloads {
		payloads[i] = payload
	}

	return payloads
}

// GetInterval returns the time interval between indivudal checks with this config.
func (m *CheckConf) GetInterval() int64 {
	return m.Interval
}

// GetTimeout returns the timeout interval of the check.
func (m *CheckConf) GetTimeout() int64 {
	return m.Timeout
}

// GetName returns the name of the check.
func (m *CheckConf) GetName() string {
	return m.Name
}

// ComponentConfig is the config of the TypeValue component of the check
// config. It stores a Type Value pair used within CheckConfig
type ComponentConfig struct {
	Type  string `yaml:"type"`
	Value string `yaml:"value"`
}

// GetType returns the type of check component.
func (c *ComponentConfig) GetType() string {
	return c.Type
}

// GetValue returns the value of check component.
func (c *ComponentConfig) GetValue() string {
	return c.Value
}

// GetCheckFromCheckProto returns a `Check` from proto.
func GetCheckFromCheckProto(agentCheck *proto.Check) Check {
	payloads := []*ComponentConfig{}
	for _, payload := range agentCheck.GetPayloads() {
		payloads = append(payloads, &ComponentConfig{
			Type:  payload.Type,
			Value: payload.Value,
		})
	}

	var input, output, target *ComponentConfig
	if agentCheck.GetInput() != nil {
		input = &ComponentConfig{
			Type:  agentCheck.Input.Type,
			Value: agentCheck.Input.Value,
		}
	}

	if agentCheck.GetOutput() != nil {
		output = &ComponentConfig{
			Type:  agentCheck.Output.Type,
			Value: agentCheck.Output.Value,
		}
	}

	if agentCheck.GetTarget() != nil {
		target = &ComponentConfig{
			Type:  agentCheck.Target.Type,
			Value: agentCheck.Target.Value,
		}
	}

	return &CheckConf{
		Input:  input,
		Output: output,
		Target: target,

		Payloads: payloads,
		Name:     agentCheck.GetName(),
		Interval: agentCheck.GetInterval(),
		Timeout:  agentCheck.GetTimeout(),
	}
}
