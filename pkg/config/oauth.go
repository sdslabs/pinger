package config

import "github.com/sdslabs/pinger/pkg/oauther"

// OauthProvider is configures the service that provides authentication
// through OAuth2.
type OauthProvider struct {
	Provider     string   `mapstructure:"provider" json:"provider"`
	ClientID     string   `mapstructure:"client_id" json:"client_id"`
	ClientSecret string   `mapstructure:"client_secret" json:"client_secret"`
	RedirectURL  string   `mapstructure:"redirect_url" json:"redirect_url"`
	Scopes       []string `mapstructure:"scopes" json:"scopes"`
}

// GetProvider returns the provider name.
func (o *OauthProvider) GetProvider() string {
	return o.Provider
}

// GetClientID returns the client ID of the provider.
func (o *OauthProvider) GetClientID() string {
	return o.ClientID
}

// GetClientSecret returns the client secret of the provider.
func (o *OauthProvider) GetClientSecret() string {
	return o.ClientSecret
}

// GetRedirectURL returns the redirect URL.
func (o *OauthProvider) GetRedirectURL() string {
	return o.RedirectURL
}

// GetScopes returns the scopes.
func (o *OauthProvider) GetScopes() []string {
	return o.Scopes
}

// Interface guard.
var _ oauther.Provider = (*OauthProvider)(nil)
