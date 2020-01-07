package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/sdslabs/status/pkg/api/server"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run status page central server.",
	Long:  "Run status page server on the given host, this server will use agents as workers to schedule periodic checks for status.",

	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Trying to run server for the status page.")

		apiServer := server.NewAPIServer()
		if err := apiServer.Run(); err != nil {
			log.Fatalf("Error while running API server: %s", err.Error())
		}
	},
}
