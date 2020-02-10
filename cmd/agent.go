package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/sdslabs/status/pkg/agent"
	"github.com/sdslabs/status/pkg/config"
	"github.com/sdslabs/status/pkg/defaults"
)

var (
	agentConfigPath string
	agentPort       int
	standaloneMode  bool
)

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Run status page agent.",
	Long: "Run status page agent on the given host, " +
		"this agent will expose a GRPC api to accept checks to perform and execute that. " +
		"The agent also has the ability to run in a standalone mode where it does not run any GRPC server.",

	Run: func(cmd *cobra.Command, args []string) {
		if !standaloneMode {
			runDefault()
			return
		}

		conf, err := config.NewAgentConfig(agentConfigPath)
		if err != nil {
			log.Fatalf("Cannot read agent config: %s", err.Error())
			return
		}

		runStandalone(conf)
	},
}

func runDefault() {
	agent.RunGRPCServer(agentPort)
}

func runStandalone(conf *config.AgentConfig) {
	agent.RunStandaloneAgent(conf)
}

func init() {
	agentCmd.Flags().StringVarP(&agentConfigPath, "config", "c", defaults.AgentConfigPath, "Config file path for agent")
	agentCmd.Flags().IntVarP(&agentPort, "port", "p", defaults.AgentPort, "Port to expose agent API on")
	agentCmd.Flags().BoolVarP(&standaloneMode, "standalone", "s", false, "Should agent run in standalone mode")
}
