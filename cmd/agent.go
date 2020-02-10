package main

import (
	"os"

	"github.com/sdslabs/status/pkg/agent"
	"github.com/sdslabs/status/pkg/config"
	"github.com/sdslabs/status/pkg/defaults"
	"github.com/sdslabs/status/pkg/metrics"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var (
	prometheusMetricsPort int
	shouldRunPrometheus   bool
	agentRunPort          int
	timescaleMetricsHost  string
	standaloneMode        bool
	agentConfigPath       string
)

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Run status page agent.",
	Long: "Run status page agent on the given host, " +
		"this agent will expose a GRPC api to accept checks to perform and execute that. " +
		"The agent also has the ability to run in a standalone mode where it does not run any GRPC server.",

	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Trying to run agent for the status page.")
		var cfg *config.AgentConfig
		var err error
		if standaloneMode {
			log.Info("Running status page agent in standalone mode.")
			if agentConfigPath == "" {
				agentConfigPath = defaults.StatusConfigPath
			}
			cfg, err = config.NewAgentConfig(agentConfigPath)
			if err != nil {
				log.Errorf("Error while parsing status page config file: %s", err)
				os.Exit(1)
			}

			if agentRunPort == defaults.AgentPort {
				agentRunPort = cfg.Port
			}

			if !shouldRunPrometheus && cfg.PrometheusMetrics {
				shouldRunPrometheus = true
			}

			if cfg.PrometheusMetricsPort != 0 {
				prometheusMetricsPort = cfg.PrometheusMetricsPort
			}

		}

		if shouldRunPrometheus && prometheusMetricsPort == 0 {
			prometheusMetricsPort = defaults.AgentPrometheusMetricsPort
		}

		if standaloneMode {
			agent.RunStandaloneAgent(cfg, getMetricsProviderConfig())
		} else {
			if prometheusMetricsPort == agentRunPort {
				log.Error("Cannot run prometheus metrics and status agent on the same port.")
				os.Exit(1)
			}

			agent.RunGRPCServer(agentRunPort, getMetricsProviderConfig())
		}
	},
}

func getMetricsProviderConfig() *metrics.ProviderConfig {
	if prometheusMetricsPort > 0 && timescaleMetricsHost != "" {
		log.Error("Status page agent does not yet support both prometheus and timescale metrics simultaneously, specify only one")
		os.Exit(1)
	}

	var metricsConfig *metrics.ProviderConfig

	if prometheusMetricsPort > 0 {
		metricsConfig = &metrics.ProviderConfig{
			PType: metrics.PrometheusProviderType,

			Host: "0.0.0.0",
			Port: prometheusMetricsPort,
		}
	} else if timescaleMetricsHost != "" {
		metricsConfig = &metrics.ProviderConfig{
			PType: metrics.TimeScaleProviderType,

			Host: timescaleMetricsHost,
		}
	} else {
		metricsConfig = &metrics.ProviderConfig{
			PType: metrics.EmptyProviderType,
		}
	}

	return metricsConfig
}

func init() {
	agentCmd.PersistentFlags().IntVar(
		&prometheusMetricsPort,
		"prometheus-port",
		defaults.AgentPrometheusMetricsPort,
		"Port to host prometheus metrics on.")

	agentCmd.PersistentFlags().BoolVar(
		&shouldRunPrometheus,
		"prometheus",
		false,
		"Should we expose metrics using prometheus.")

	agentCmd.PersistentFlags().IntVarP(
		&agentRunPort,
		"port",
		"p",
		defaults.AgentPort,
		"Port to run the agent on grpc server on.")

	agentCmd.PersistentFlags().StringVarP(
		&timescaleMetricsHost,
		"ts-metrics",
		"t",
		"",
		"Run the agent with push timescale metrics, provide timescale host information int this string.")

	agentCmd.PersistentFlags().BoolVarP(
		&standaloneMode,
		"standalone",
		"s",
		false,
		"Run agent in the standalone mode, "+
			"this does not expose any grpc server to collect the work, "+
			"it takes that data from a config file mentioned in another argument.")

	agentCmd.PersistentFlags().StringVarP(
		&agentConfigPath,
		"config",
		"c",
		defaults.StatusConfigPath,
		"Path to where find the config for the agent in standalone mode, "+
			"this file contains all information including the hosts to ping using which checks")
}
