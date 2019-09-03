package config

import (
	"github.com/sdslabs/status/pkg/agent/proto"
)

// Check is the interface which every check that needs to be processed here
// should implement.
type Check interface {
	GetInput() CheckComponent
	GetOutput() CheckComponent
	GetTarget() CheckComponent

	GetPayloads() []CheckComponent
	GetInterval() int64
	GetTimeout() int64

	GetName() string
}

// CheckComponent is the Type Value component for check components like Input, Output, Target etc.
type CheckComponent interface {
	GetType() string
	GetValue() string
}

// Config Associated with each check.
type CheckConfig struct {
	Input  *CheckComponentConfig `yaml:"input"`
	Output *CheckComponentConfig `yaml:"output"`
	Target *CheckComponentConfig `yaml:"target"`

	Payloads []*CheckComponentConfig `yaml:"payloads"`

	Name     string `yaml:"name"`
	Timeout  int64  `yaml:"timeout"`
	Interval int64  `yaml:"interval"`
}

func (m *CheckConfig) GetInput() CheckComponent {
	return m.Input
}

func (m *CheckConfig) GetTarget() CheckComponent {
	return m.Target
}

func (m *CheckConfig) GetOutput() CheckComponent {
	return m.Output
}

func (m *CheckConfig) GetPayloads() []CheckComponent {
	payloads := make([]CheckComponent, len(m.Payloads))
	for i, payload := range m.Payloads {
		payloads[i] = payload
	}

	return payloads
}

func (m *CheckConfig) GetInterval() int64 {
	return m.Interval
}

func (m *CheckConfig) GetTimeout() int64 {
	return m.Timeout
}

func (m *CheckConfig) GetName() string {
	return m.Name
}

// CheckComponentConfig is the config of the TypeValue component of the check
// config. It stores a Type Value pair used within CheckConfig
type CheckComponentConfig struct {
	Type  string `yaml:"type"`
	Value string `yaml:"value"`
}

func (c *CheckComponentConfig) GetType() string {
	return c.Type
}

func (c *CheckComponentConfig) GetValue() string {
	return c.Value
}

func GetCheckFromCheckProto(agentCheck *proto.Check) Check {
	payloads := []*CheckComponentConfig{}
	for _, payload := range agentCheck.GetPayloads() {
		payloads = append(payloads, &CheckComponentConfig{
			Type:  payload.Type,
			Value: payload.Value,
		})
	}

	var input, output, target *CheckComponentConfig
	if agentCheck.GetInput() != nil {
		input = &CheckComponentConfig{
			Type:  agentCheck.Input.Type,
			Value: agentCheck.Input.Value,
		}
	}

	if agentCheck.GetOutput() != nil {
		output = &CheckComponentConfig{
			Type:  agentCheck.Output.Type,
			Value: agentCheck.Output.Value,
		}
	}

	if agentCheck.GetTarget() != nil {
		target = &CheckComponentConfig{
			Type:  agentCheck.Target.Type,
			Value: agentCheck.Target.Value,
		}
	}

	return &CheckConfig{
		Input:  input,
		Output: output,
		Target: target,

		Payloads: payloads,
		Name:     agentCheck.GetName(),
		Interval: agentCheck.GetInterval(),
		Timeout:  agentCheck.GetTimeout(),
	}
}
