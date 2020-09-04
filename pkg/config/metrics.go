// Copyright (c) 2020 SDSLabs
// Use of this source code is governed by an MIT license
// details of which can be found in the LICENSE file.

package config

import (
	"time"

	"github.com/sdslabs/pinger/pkg/checker"
	"github.com/sdslabs/pinger/pkg/exporter"
)

// Metric represents the result of a check.
//
// Implements the checker.Metric interface.
type Metric struct {
	CheckID    uint
	CheckName  string
	Successful bool
	Timeout    bool
	StartTime  time.Time
	Duration   time.Duration
}

// GetCheckID returns the ID of the check for which the metric is.
func (m *Metric) GetCheckID() uint {
	return m.CheckID
}

// GetCheckName returns the name of the check.
func (m *Metric) GetCheckName() string {
	return m.CheckName
}

// IsSuccessful tells if the check was successful.
func (m *Metric) IsSuccessful() bool {
	return m.Successful
}

// IsTimeout tells if the check timed-out.
func (m *Metric) IsTimeout() bool {
	return m.Timeout
}

// GetStartTime returns the start-time of the check.
func (m *Metric) GetStartTime() time.Time {
	return m.StartTime
}

// GetDuration returns the duration that check took to run.
func (m *Metric) GetDuration() time.Duration {
	return m.Duration
}

// MetricsProvider represents the configuration of a metrics exporter.
//
// Implements the metrics.Provider interface.
type MetricsProvider struct {
	Backend  string `mapstructure:"backend" json:"backend"`
	Host     string `mapstructure:"host" json:"host"`
	Port     uint16 `mapstructure:"port" json:"port"`
	DBName   string `mapstructure:"db_name" json:"db_name"`
	Username string `mapstructure:"username" json:"username"`
	Password string `mapstructure:"password" json:"password"`
	SSLMode  bool   `mapstructure:"ssl_mode" json:"ssl_mode"`
}

// GetBackend returns the backend of the provider.
func (m *MetricsProvider) GetBackend() string {
	return m.Backend
}

// GetHost returns the host of the database provider.
func (m *MetricsProvider) GetHost() string {
	return m.Host
}

// GetPort returns the port of the database provider.
func (m *MetricsProvider) GetPort() uint16 {
	return m.Port
}

// GetDBName returns the database name of the provider.
func (m *MetricsProvider) GetDBName() string {
	return m.DBName
}

// GetUsername returns the username of the database provider.
func (m *MetricsProvider) GetUsername() string {
	return m.Username
}

// GetPassword returns the password of the database provider.
func (m *MetricsProvider) GetPassword() string {
	return m.Password
}

// IsSSLMode tells if the connection with the provider is to be established
// through SSL.
func (m *MetricsProvider) IsSSLMode() bool {
	return m.SSLMode
}

// Interface guards.
var (
	_ checker.Metric    = (*Metric)(nil)
	_ exporter.Provider = (*MetricsProvider)(nil)
)
