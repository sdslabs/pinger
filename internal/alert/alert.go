// Copyright (c) 2020 SDSLabs
// Use of this source code is governed by an MIT license
// details of which can be found in the LICENSE file.

package alert

import (
	"time"
)

// Alert represents the list of checks that have changed status.
//
// Implements the checker.Alert interface.
type Alert struct {
	CheckID    uint
	CheckName  string
	Successful bool
	Timeout    bool
	StartTime  time.Time
	Duration   time.Duration
}

// GetCheckID returns the ID of the check for which the Alert is.
func (a *Alert) GetCheckID() uint {
	return a.CheckID
}

// GetCheckName returns the name of the check.
func (a *Alert) GetCheckName() string {
	return a.CheckName
}

// IsSuccessful tells if the check was successful.
func (a *Alert) IsSuccessful() bool {
	return a.Successful
}

// IsTimeout tells if the check timed-out.
func (a *Alert) IsTimeout() bool {
	return a.Timeout
}

// GetStartTime returns the start-time of the check.
func (a *Alert) GetStartTime() time.Time {
	return a.StartTime
}

// GetDuration returns the duration that check took to run.
func (a *Alert) GetDuration() time.Duration {
	return a.Duration
}

// AlertsProvider represents the configuration of an Alerts exporter.
//
// Implements the Alerts.AlertSystem interface.
type AlertsProvider struct {
	Service     string `mapstructure:"service" json:"service"`
	Webhook		string `mapstructure:"webhook" json:"webhook"`
}

// GetService returns the preferred method for alert.
func (a *AlertsProvider) GetService() string {
	return a.Service
}

// GetWebhook returns the webhook URL of the selected alert method.
func (a *AlertsProvider) GetWebhook() string {
	return a.Webhook
}

// Interface Guards?