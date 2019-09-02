package main

import (
	"os"

	"github.com/sdslabs/status/pkg/api/agent"
	"github.com/sdslabs/status/pkg/defaults"
	"github.com/sdslabs/status/pkg/metrics"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	prometheusMetricsPort int
	agentRunPort          int
	timescaleMetricsHost  string
)

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Run status page agent.",
	Long:  "Run status page agent on the given host, this agent will expose a GRPC api to accept checks to perform and execute that.",

	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Trying to run agent for the status page.")

		if prometheusMetricsPort == agentRunPort {
			log.Error("Cannot run prometheus metrics and status agent on the same port.")
			os.Exit(1)
		}

		if prometheusMetricsPort > 0 && timescaleMetricsHost != "" {
			log.Error("Status page agnet does not yet support both prometheus and timescale metrics simultaneously, specify only one")
			os.Exit(1)
		}

		var config *metrics.ProviderConfig
		if prometheusMetricsPort > 0 {
			config = &metrics.ProviderConfig{
				PType: metrics.PrometheusProviderType,

				Host: "0.0.0.0",
				Port: prometheusMetricsPort,
			}
		} else if timescaleMetricsHost != "" {
			config = &metrics.ProviderConfig{
				PType: metrics.TimeScaleProviderType,

				Host: timescaleMetricsHost,
			}
		} else {
			config = &metrics.ProviderConfig{
				PType: metrics.EmptyProviderType,
			}
		}

		agent.RunGrpcServer(agentRunPort, config)
	},
}

func init() {
	agentCmd.PersistentFlags().IntVarP(&prometheusMetricsPort, "metrics-port", "pm", defaults.DefaultAgentPrometheusMetricsPort, "Port to host prometheus metrics on.")
	agentCmd.PersistentFlags().IntVarP(&agentRunPort, "port", "p", defaults.DefaultAgentPort, "Port to run the agent on grpc server on.")

	agentCmd.PersistentFlags().StringVarP(&timescaleMetricsHost, "ts-metrics", "-m", "", "Run the agent with push timescale metrics, provide timescale host information int this string.")
}
