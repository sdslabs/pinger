package config

import (
	"strings"
	"time"

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
	GetId() uint32 //nolint:golint
}

// Component is the Type Value component for check components like Input, Output, Target etc.
type Component interface {
	GetType() string
	GetValue() string
}

// CheckConf Associated with each check.
type CheckConf struct {
	Input  *ComponentConfig `mapstructure:"input" json:"input" yaml:"input" toml:"input"`
	Output *ComponentConfig `mapstructure:"output" json:"output" yaml:"output" toml:"output"`
	Target *ComponentConfig `mapstructure:"target" json:"target" yaml:"target" toml:"target"`

	Payloads []*ComponentConfig `mapstructure:"payloads" json:"payloads" yaml:"payloads" toml:"payloads"`

	ID       uint          `mapstructure:"id" json:"id" yaml:"id" toml:"id"`
	Name     string        `mapstructure:"name" json:"name" yaml:"name" toml:"name"`
	Timeout  time.Duration `mapstructure:"timeout" json:"timeout" yaml:"timeout" toml:"timeout"`
	Interval time.Duration `mapstructure:"interval" json:"interval" yaml:"interval" toml:"interval"`
}

// GetLabel returns a slug that is unique for checks deployed with the manager on the agent.
func (m *CheckConf) GetLabel() string {
	name := strings.ToLower(m.Name)
	return strings.ReplaceAll(name, " ", "-")
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
	return int64(m.Interval)
}

// GetTimeout returns the timeout interval of the check.
func (m *CheckConf) GetTimeout() int64 {
	return int64(m.Timeout)
}

// GetName returns the name of the check.
func (m *CheckConf) GetName() string {
	return m.Name
}

// GetId returns the ID of the check.
func (m *CheckConf) GetId() uint32 { //nolint:golint
	return uint32(m.ID)
}

// ComponentConfig is the config of the TypeValue component of the check
// config. It stores a Type Value pair used within CheckConfig
type ComponentConfig struct {
	Type  string `mapstructure:"type" json:"type" yaml:"type" toml:"type"`
	Value string `mapstructure:"value" json:"value" yaml:"value" toml:"value"`
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
		Interval: time.Duration(agentCheck.GetInterval()),
		Timeout:  time.Duration(agentCheck.GetTimeout()),
	}
}
