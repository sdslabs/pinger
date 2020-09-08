// Copyright (c) 2020 SDSLabs
// Use of this source code is governed by an MIT license
// details of which can be found in the LICENSE file.

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
