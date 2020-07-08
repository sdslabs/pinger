// Copyright (c) 2020 SDSLabs
// Use of this source code is governed by an MIT license
// details of which can be found in the LICENSE file.

package config

import (
	"github.com/sdslabs/status/internal/alert"
)

// AlertsProvider represents the configuration of an Alerts exporter.
//
// Implements the Alerts.AlertSystem interface.
type AlertsProvider struct {
	Service string `mapstructure:"service" json:"service"`
	Webhook string `mapstructure:"webhook" json:"webhook"`
}

// GetService returns the preferred method for alert.
func (a *AlertsProvider) GetService() string {
	return a.Service
}

// GetWebhook returns the webhook URL of the selected alert method.
func (a *AlertsProvider) GetWebhook() string {
	return a.Webhook
}

// Interface guards.
var (
	_ alert.Provider = (*AlertsProvider)(nil)
)
