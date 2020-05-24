// Copyright (c) 2020 SDSLabs
// Use of this source code is governed by an MIT license
// details of which can be found in the LICENSE file.

package config

import (
	"time"

	"github.com/sdslabs/status/internal/metrics"
)

// MetricsProvider represents the configuration of a metrics exporter.
//
// Implements the metrics.Provider interface.
type MetricsProvider struct {
	Backend  string        `mapstructure:"backend" json:"backend"`
	Host     string        `mapstructure:"host" json:"host"`
	Port     int           `mapstructure:"port" json:"port"`
	DBName   string        `mapstructure:"db_name" json:"db_name"`
	Username string        `mapstructure:"username" json:"username"`
	Password string        `mapstructure:"password" json:"password"`
	SSLMode  bool          `mapstructure:"ssl_mode" json:"ssl_mode"`
	Interval time.Duration `mapstructure:"interval" json:"interval"`
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
	return uint16(m.Port)
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

// GetInterval returns the interval after which metrics are exported.
func (m *MetricsProvider) GetInterval() time.Duration {
	return m.Interval
}

// Interface guard.
var _ metrics.Provider = (*MetricsProvider)(nil)
