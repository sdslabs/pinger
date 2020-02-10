package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/sdslabs/status/pkg/api/application"
	"github.com/sdslabs/status/pkg/defaults"
)

var (
	applicationAPIPort int

	applicationAPICmd = &cobra.Command{
		Use:   "application",
		Short: "Run status page application server.",
		Long:  "Run status page server on the given host, this server is exposed to public to use the web API.",

		Run: func(cmd *cobra.Command, args []string) {
			log.Info("Trying to run the application server for the status page.")

			if err := application.Serve(applicationAPIPort); err != nil {
				log.Fatalf("Error while running the application server: %s", err.Error())
			}
		},
	}
)

func init() {
	applicationAPICmd.Flags().IntVarP(
		&applicationAPIPort,
		"port",
		"p",
		defaults.ApplicationAPIPort,
		"Port to run application server on.")
}
