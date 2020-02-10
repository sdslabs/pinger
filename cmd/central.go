package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/sdslabs/status/pkg/api/central"
	"github.com/sdslabs/status/pkg/defaults"
)

var (
	centralAPIPort int

	centralAPICmd = &cobra.Command{
		Use:   "central",
		Short: "Run status page central server.",
		Long:  "Run status page server on the given host, this server will use agents as workers to schedule periodic checks for status.",

		Run: func(cmd *cobra.Command, args []string) {
			log.Info("Trying to run central server for the status page.")

			apiServer := central.NewAPIServer(centralAPIPort)
			if err := apiServer.Run(); err != nil {
				log.Fatalf("Error while running API server: %s", err.Error())
			}
		},
	}
)

func init() {
	centralAPICmd.Flags().IntVarP(
		&centralAPIPort,
		"port",
		"p",
		defaults.CentralAPIPort,
		"Port to run central server on.")
}
