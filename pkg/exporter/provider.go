package exporter

// Provider is anything that can be used to configure and create a metrics
// exporter.
type Provider interface {
	GetBackend() string // Returns the provider backend name.

	GetHost() string     // Returns the host.
	GetPort() uint16     // Returns the port.
	GetOrgName() string  // Returns the organization name.
	GetDBName() string   // Returns the database name.
	GetUsername() string // Returns the username.
	GetPassword() string // Returns the password.
	IsSSLMode() bool     // Tells if connection is through SSL mode.
}
