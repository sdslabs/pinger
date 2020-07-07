// Copyright (c) 2020 SDSLabs
// Use of this source code is governed by an MIT license
// details of which can be found in the LICENSE file.

package alert

// Alertsystem is anything that can be used to configure and create a alert
// exporter.
type AlertSystem interface {
	GetService() string // Returns the service name.
	GetWebhook() string // Returns the webhook URL of the service.
}

// similar to inetrnal/exporter/provider we create an interface for storing the webhook url - do we need any other info?
// declare alert struct in config/ - alert.go
// should use controller.PullLatestStats for the last state of check if last state is empty
// should receive its values from internal/agent/agent line 122 - implemented in exporter.go?
