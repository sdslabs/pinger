// Copyright (c) 2020 SDSLabs
// Use of this source code is governed by an MIT license
// details of which can be found in the LICENSE file.

package configfile

import "github.com/sdslabs/status/internal/config"

// Agent represents the configuration for an agent.
type Agent struct {
	Standalone bool                   `mapstructure:"standalone" json:"standalone"`
	Port       uint16                 `mapstructure:"port" json:"port"`
	Metrics    config.MetricsProvider `mapstructure:"metrics" json:"metrics"`
	Checks     []config.Check         `mapstructure:"checks" json:"checks"`
}
