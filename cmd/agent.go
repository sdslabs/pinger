package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/sdslabs/status/pkg/agent"
	"github.com/sdslabs/status/pkg/config"
	"github.com/sdslabs/status/pkg/defaults"
)

const (
	keyAgentConfigPort  = "port"
	flagAgentConfigPort = "port"

	keyAgentConfigStandalone  = "standalone"
	flagAgentConfigStandalone = "standalone"

	keyAgentConfigMetricsBackend  = "metrics.backend"
	flagAgentConfigMetricsBackend = "metrics-backend"

	keyAgentConfigMetricsHost  = "metrics.host"
	flagAgentConfigMetricsHost = "metrics-host"

	keyAgentConfigMetricsPort  = "metrics.port"
	flagAgentConfigMetricsPort = "metrics-port"

	keyAgentConfigMetricsUsername  = "metrics.username"
	flagAgentConfigMetricsUsername = "metrics-username"

	keyAgentConfigMetricsPassword  = "metrics.password"
	flagAgentConfigMetricsPassword = "metrics-password"

	keyAgentConfigMetricsDBName  = "metrics.db_name"
	flagAgentConfigMetricsDBName = "metrics-db-name"

	keyAgentConfigMetricsSSLMode  = "metrics.ssl_mode"
	flagAgentConfigMetricsSSLMode = "metrics-ssl-mode"

	keyAgentConfigMetricsInterval  = "metrics.interval"
	flagAgentConfigMetricsInterval = "metrics-interval"
)

var (
	agentConfigPath string
	agentConf       config.AgentConfig
)

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Run status page agent.",
	Long: `Run status page agent on the given host, this agent will expose
a GRPC api to accept checks to perform and execute that. The
agent also has the ability to run in a standalone mode where it
does not run any GRPC server.`,

	PreRun: func(*cobra.Command, []string) {
		initConfig(agentConfigPath, defaults.AgentConfigPath, &agentConf)
	},

	Run: func(*cobra.Command, []string) {
		a, err := agent.NewAgent(&agentConf)
		if err != nil {
			log.WithError(err).Fatalln("cannot create agent")
		}

		log.WithField("standalone", agentConf.Standalone).Infoln("running agent")
		log.WithError(a.Run()).Fatalln("stopped running agent")
	},
}

func init() {
	agentCmd.Flags().StringVarP(&agentConfigPath, "config", "c", defaults.AgentConfigPath, "Config file path for agent")

	agentCmd.Flags().IntP(flagAgentConfigPort, "p", defaults.AgentPort, "Port to expose agent API on")
	agentCmd.Flags().BoolP(flagAgentConfigStandalone, "s", false, "Should agent run in standalone mode")

	agentCmd.Flags().String(flagAgentConfigMetricsBackend, defaults.AgentMetricsBackend, "Backend service to store metrics")
	agentCmd.Flags().String(flagAgentConfigMetricsHost, defaults.AgentMetricsHost, "Host to run metrics server")
	agentCmd.Flags().Int(flagAgentConfigMetricsPort, defaults.AgentMetricsPort, "Port to run metrics server on")
	agentCmd.Flags().String(flagAgentConfigMetricsUsername, "", "Username credential for metrics")
	agentCmd.Flags().String(flagAgentConfigMetricsPassword, "", "Password credential for metrics")
	agentCmd.Flags().String(flagAgentConfigMetricsDBName, "", "Database name for metrics")
	agentCmd.Flags().Bool(flagAgentConfigMetricsSSLMode, true, "Whether to run metrics with SSL")
	agentCmd.Flags().Duration(flagAgentConfigMetricsInterval, defaults.AgentMetricsInterval, "Interval after which metrics are pushed/pulled")

	if err := viper.BindPFlag(keyAgentConfigPort, agentCmd.Flags().Lookup(flagAgentConfigPort)); err != nil {
		viperErr(err)
	}
	if err := viper.BindPFlag(keyAgentConfigStandalone, agentCmd.Flags().Lookup(flagAgentConfigStandalone)); err != nil {
		viperErr(err)
	}

	if err := viper.BindPFlag(keyAgentConfigMetricsBackend, agentCmd.Flags().Lookup(flagAgentConfigMetricsBackend)); err != nil {
		viperErr(err)
	}
	if err := viper.BindPFlag(keyAgentConfigMetricsHost, agentCmd.Flags().Lookup(flagAgentConfigMetricsHost)); err != nil {
		viperErr(err)
	}
	if err := viper.BindPFlag(keyAgentConfigMetricsPort, agentCmd.Flags().Lookup(flagAgentConfigMetricsPort)); err != nil {
		viperErr(err)
	}
	if err := viper.BindPFlag(keyAgentConfigMetricsDBName, agentCmd.Flags().Lookup(flagAgentConfigMetricsDBName)); err != nil {
		viperErr(err)
	}
	if err := viper.BindPFlag(keyAgentConfigMetricsUsername, agentCmd.Flags().Lookup(flagAgentConfigMetricsPassword)); err != nil {
		viperErr(err)
	}
	if err := viper.BindPFlag(keyAgentConfigMetricsPassword, agentCmd.Flags().Lookup(flagAgentConfigMetricsPassword)); err != nil {
		viperErr(err)
	}
	if err := viper.BindPFlag(keyAgentConfigMetricsSSLMode, agentCmd.Flags().Lookup(flagAgentConfigMetricsSSLMode)); err != nil {
		viperErr(err)
	}
	if err := viper.BindPFlag(keyAgentConfigMetricsInterval, agentCmd.Flags().Lookup(flagAgentConfigMetricsInterval)); err != nil {
		viperErr(err)
	}
}
