package main

import (
	"github.com/sdslabs/status/pkg/api/agent"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Run status page agent.",
	Long:  "Run status page agent on the given host, this agent will expose a GRPC api to accept checks to perform and execute that.",

	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Trying to run agent for the status page.")

		agent.RunGrpcServer(nil)
	},
}
