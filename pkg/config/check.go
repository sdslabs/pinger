// Copyright (c) 2020 SDSLabs
// Use of this source code is governed by an MIT license
// details of which can be found in the LICENSE file.

package config

import (
	"time"

	"github.com/sdslabs/pinger/pkg/checker"
	"github.com/sdslabs/pinger/pkg/components/agent/proto"
)

// Check is the configuration of each check associated with each check.
//
// Implements checker.Check interface.
type Check struct {
	ID       string        `mapstructure:"id" json:"id"`
	Name     string        `mapstructure:"name" json:"name"`
	Interval time.Duration `mapstructure:"interval" json:"interval"`
	Timeout  time.Duration `mapstructure:"timeout" json:"timeout"`
	Input    Component     `mapstructure:"input" json:"input"`
	Output   Component     `mapstructure:"output" json:"output"`
	Target   Component     `mapstructure:"target" json:"target"`
	Payloads []Component   `mapstructure:"payloads" json:"payloads"`
	Alerts   []Alert       `mapstructure:"alerts" json:"alerts"`
}

// GetID returns the ID for the check.
func (c *Check) GetID() string {
	return c.ID
}

// GetName returns the name of the check.
func (c *Check) GetName() string {
	return c.Name
}

// GetTimeout returns the timeout of the check.
func (c *Check) GetTimeout() time.Duration {
	return c.Timeout
}

// GetInterval returns the interval of the check.
func (c *Check) GetInterval() time.Duration {
	return c.Interval
}

// GetInput returns the input of the check.
func (c *Check) GetInput() checker.Component {
	return &c.Input
}

// GetOutput returns the output of the check.
func (c *Check) GetOutput() checker.Component {
	return &c.Output
}

// GetTarget returns the target of the check.
func (c *Check) GetTarget() checker.Component {
	return &c.Target
}

// GetPayloads returns the payloads for the check.
func (c *Check) GetPayloads() []checker.Component {
	payloads := make([]checker.Component, len(c.Payloads))
	for i := range c.Payloads {
		payloads[i] = &c.Payloads[i]
	}
	return payloads
}

// Component is a key-value pair.
//
// Implements Component interface.
type Component struct {
	Type  string `mapstructure:"type" json:"type"`
	Value string `mapstructure:"value" json:"value"`
}

// GetType returns the type of the component.
func (c *Component) GetType() string {
	return c.Type
}

// GetValue returns the value of the component.
func (c *Component) GetValue() string {
	return c.Value
}

// ProtoToCheck converts a proto.Check into checker.Check.
func ProtoToCheck(check *proto.Check) Check {
	payloads := make([]Component, len(check.Payloads))
	for i := range check.Payloads {
		payloads[i] = Component{
			Type:  check.Payloads[i].Type,
			Value: check.Payloads[i].Value,
		}
	}

	alerts := make([]Alert, len(check.Alerts))
	for i := range check.Alerts {
		alerts[i] = Alert{
			Service: check.Alerts[i].Service,
		}
	}

	return Check{
		ID:       check.ID,
		Name:     check.Name,
		Interval: time.Duration(check.Interval),
		Timeout:  time.Duration(check.Timeout),
		Input: Component{
			Type:  check.Input.Type,
			Value: check.Input.Value,
		},
		Output: Component{
			Type:  check.Output.Type,
			Value: check.Output.Value,
		},
		Target: Component{
			Type:  check.Target.Type,
			Value: check.Target.Value,
		},
		Payloads: payloads,
		Alerts:   alerts,
	}
}

// Interface guards.
var (
	_ checker.Check     = (*Check)(nil)
	_ checker.Component = (*Component)(nil)
)
