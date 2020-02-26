package config

// OauthProviderConfig provides configuration settings for an OAuth Provider.
type OauthProviderConfig struct {
	ClientID     string   `mapstructure:"client_id" json:"client_id" yaml:"client_id" toml:"client_id"`
	ClientSecret string   `mapstructure:"client_secret" json:"client_secret" yaml:"client_secret" toml:"client_secret"`
	RedirectURL  string   `mapstructure:"redirect_url" json:"redirect_url" yaml:"redirect_url" toml:"redirect_url"`
	Scopes       []string `mapstructure:"scopes" json:"scopes" yaml:"scopes" toml:"scopes"`
}

type oauth = map[string]OauthProviderConfig

type database struct {
	Host     string `mapstructure:"host" json:"host" yaml:"host" toml:"host"`
	Port     int    `mapstructure:"port" json:"port" yaml:"port" toml:"port"`
	Username string `mapstructure:"username" json:"username" yaml:"username" toml:"username"`
	Password string `mapstructure:"password" json:"password" yaml:"password" toml:"password"`
	Name     string `mapstructure:"name" json:"name" yaml:"name" toml:"name"`
	SSLMode  bool   `mapstructure:"ssl_mode" json:"ssl_mode" yaml:"ssl_mode" toml:"ssl_mode"`
}

// AppConfig for `config.yml`
type AppConfig struct {
	Port int `mapstructure:"port" json:"port" yaml:"port" toml:"port"`

	// Use (*AppConfig).Secret() to access secret value.
	SecretVal string `mapstructure:"secret" json:"secret" yaml:"secret" toml:"secret"`

	Oauth    oauth    `mapstructure:"oauth" json:"oauth" yaml:"oauth" toml:"oauth"`
	Database database `mapstructure:"database" json:"database" yaml:"database" toml:"database"`
}

// Secret returns the secret key to encrypt tokens.
func (c *AppConfig) Secret() []byte {
	return []byte(c.SecretVal)
}
