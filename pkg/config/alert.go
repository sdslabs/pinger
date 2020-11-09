package config

import "github.com/sdslabs/pinger/pkg/alerter"

// AlertProvider configures the alerter.
//
// This implements the alerter.Provider interface.
type AlertProvider struct {
	Service string `json:"service" mapstructure:"service"`
	Host    string `json:"host" mapstructure:"host"`
	Port    uint16 `json:"port" mapstructure:"port"`
	User    string `json:"user" mapstructure:"user"`
	Secret  string `json:"secret" mapstructure:"secret"`
}

// GetService returns the service name of the provider.
func (a *AlertProvider) GetService() string {
	return a.Service
}

// GetHost returns the host name of the service provider.
func (a *AlertProvider) GetHost() string {
	return a.Host
}

// GetPort returns the port number of the service provider.
func (a *AlertProvider) GetPort() uint16 {
	return a.Port
}

// GetUser returns the user identifier of the service provider.
func (a *AlertProvider) GetUser() string {
	return a.User
}

// GetSecret returns the secret key for configuring provider.
func (a *AlertProvider) GetSecret() string {
	return a.Secret
}

// Alert configures where alert is to be sent.
type Alert struct {
	Service string `json:"service" mapstructure:"service"`
	Target  string `json:"target" mapstructure:"target"`
}

// GetService returns the service name of the alert.
func (a *Alert) GetService() string {
	return a.Service
}

// GetTarget returns the target of the alert.
func (a *Alert) GetTarget() string {
	return a.Target
}

// Interface guards.
var (
	_ alerter.Provider = (*AlertProvider)(nil)
	_ alerter.Alert    = (*Alert)(nil)
)
