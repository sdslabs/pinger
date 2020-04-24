package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/sdslabs/status/pkg/api/app"
	"github.com/sdslabs/status/pkg/config"
	"github.com/sdslabs/status/pkg/defaults"
	"github.com/sdslabs/status/pkg/utils"
)

const (
	keyAppConfigPort  = "port"
	flagAppConfigPort = "port"

	keyAppConfigSecret  = "secret"
	flagAppConfigSecret = "secret"

	keyAppConfigDBHost  = "database.host"
	flagAppConfigDBHost = "database-host"

	keyAppConfigDBPort  = "database.port"
	flagAppConfigDBPort = "database-port"

	keyAppConfigDBUsername  = "database.username"
	flagAppConfigDBUsername = "database-username"

	keyAppConfigDBPassword  = "database.password"
	flagAppConfigDBPassword = "database-password"

	keyAppConfigDBName  = "database.name"
	flagAppConfigDBName = "database-name"

	keyAppConfigDBSSLMode  = "database.ssl_mode"
	flagAppConfigDBSSLMode = "database-ssl-mode"
)

var (
	appConfigPath string
	appConf       config.AppConfig
)

var appAPICmd = &cobra.Command{
	Use:   "app",
	Short: "Run status page application server.",
	Long:  "Run status page server on the given host, this server is exposed to public to use the web API.",

	PreRun: func(*cobra.Command, []string) {
		initConfig(appConfigPath, defaults.AppConfigPath, &appConf)
	},

	Run: func(*cobra.Command, []string) {
		log.WithField("port", appConf.Port).Infoln("starting app api server")
		log.WithField("app_secret", appConf.SecretVal).Warnf("save this secret in case passed via flags")

		if err := app.Serve(&appConf); err != nil {
			log.WithError(err).Fatalln("cannot start app server")
		}
	},
}

func init() {
	appAPICmd.Flags().StringVarP(&appConfigPath, "config", "c", defaults.AppConfigPath, "Config file path for app api server")

	appAPICmd.Flags().IntP(flagAppConfigPort, "p", defaults.AppAPIPort, "Port to server app api server on")
	appAPICmd.Flags().StringP(flagAppConfigSecret, "s", utils.RandomToken(), "Application secret, used to encrypt tokens")

	appAPICmd.Flags().String(flagAppConfigDBHost, defaults.AppAPIDBHost, "Database host")
	appAPICmd.Flags().Int(flagAppConfigDBPort, defaults.AppAPIDBPort, "Database port")
	appAPICmd.Flags().String(flagAppConfigDBUsername, "", "Database username")
	appAPICmd.Flags().String(flagAppConfigDBPassword, "", "Database password")
	appAPICmd.Flags().String(flagAppConfigDBName, "", "Database name")
	appAPICmd.Flags().Bool(flagAppConfigDBSSLMode, true, "Should use ssl to connect with DB?")

	if err := viper.BindPFlag(keyAppConfigPort, appAPICmd.Flags().Lookup(flagAppConfigPort)); err != nil {
		viperErr(err)
	}
	if err := viper.BindPFlag(keyAppConfigSecret, appAPICmd.Flags().Lookup(flagAppConfigSecret)); err != nil {
		viperErr(err)
	}

	if err := viper.BindPFlag(keyAppConfigDBHost, appAPICmd.Flags().Lookup(flagAppConfigDBHost)); err != nil {
		viperErr(err)
	}
	if err := viper.BindPFlag(keyAppConfigDBPort, appAPICmd.Flags().Lookup(flagAppConfigDBPort)); err != nil {
		viperErr(err)
	}
	if err := viper.BindPFlag(keyAppConfigDBUsername, appAPICmd.Flags().Lookup(flagAppConfigDBUsername)); err != nil {
		viperErr(err)
	}
	if err := viper.BindPFlag(keyAppConfigDBPassword, appAPICmd.Flags().Lookup(flagAppConfigDBPassword)); err != nil {
		viperErr(err)
	}
	if err := viper.BindPFlag(keyAppConfigDBName, appAPICmd.Flags().Lookup(flagAppConfigDBName)); err != nil {
		viperErr(err)
	}
	if err := viper.BindPFlag(keyAppConfigDBSSLMode, appAPICmd.Flags().Lookup(flagAppConfigDBSSLMode)); err != nil {
		viperErr(err)
	}
}
