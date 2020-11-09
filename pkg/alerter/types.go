package alerter

// Provider is anything that can be used to configure and create an alerter.
type Provider interface {
	GetService() string // Name of the alert service.

	GetHost() string
	GetPort() uint16
	GetUser() string   // Username or Email.
	GetSecret() string // Password or token.
}

// Alert is anything that tells the alerter where to send the alert.
type Alert interface {
	GetService() string // Name of the alert service.

	GetTarget() string // Target of alert.
}
