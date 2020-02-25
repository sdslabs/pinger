package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/sdslabs/status/pkg/api/app"
	"github.com/sdslabs/status/pkg/config"
	"github.com/sdslabs/status/pkg/defaults"
)

var (
	appAPIPort int

	appAPICmd = &cobra.Command{
		Use:   "app",
		Short: "Run status page application server.",
		Long:  "Run status page server on the given host, this server is exposed to public to use the web API.",

		Run: func(cmd *cobra.Command, args []string) {
			log.Info("Trying to run the application server for the status page.")

			conf, err := config.GetStatusConfig(statusConfigPath)
			if err != nil {
				log.Fatalf("Error reading config file '%s': %s", statusConfigPath, err.Error())
				return
			}

			if err := app.Serve(&conf, appAPIPort); err != nil {
				log.Fatalf("Error while running the application server: %s", err.Error())
			}
		},
	}
)

func init() {
	appAPICmd.Flags().StringVarP(
		&statusConfigPath, "config", "c", defaults.StatusConfigPath, "Config file path of status API")

	appAPICmd.Flags().IntVarP(
		&appAPIPort,
		"port",
		"p",
		defaults.AppAPIPort,
		"Port to run application server on.")
}
