package configfile

import "github.com/sdslabs/pinger/pkg/config"

// App is the configuration file for app server.
type App struct {
	// Port to run the app server on.
	Port uint16 `mapstructure:"port" json:"port"`

	// Secret for signing tokens.
	Secret string `mapstructure:"secret" json:"secret"`

	// Oauth providers with configuration.
	Oauth []config.OauthProvider `mapstructure:"oauth" json:"oauth"`

	// Database connection.
	Database config.DBConn `mapstructure:"database" json:"database"`
}
