// Copyright (c) 2020 SDSLabs
// Use of this source code is governed by an MIT license
// details of which can be found in the LICENSE file.

package metrics

import "time"

// Provider is anything that can be used to configure and create a metrics
// exporter.
type Provider interface {
	GetBackend() string // Returns the provider backend name.

	GetHost() string     // Returns the host.
	GetPort() uint16     // Returns the port.
	GetDBName() string   // Returns the database name.
	GetUsername() string // Returns the username.
	GetPassword() string // Returns the password.
	IsSSLMode() bool     // Tells if connection is through SSL mode.

	GetInterval() time.Duration // Returns the interval after which metrics are exported.
}
