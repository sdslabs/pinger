// Copyright (c) 2020 SDSLabs
// Use of this source code is governed by an MIT license
// details of which can be found in the LICENSE file.

package alert

// System is anything that can be used to configure and create a alert
// exporter.
type System interface {
	GetService() string // Returns the service name.
	GetWebhook() string // Returns the webhook URL of the service.
}
