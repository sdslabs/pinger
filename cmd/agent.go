package cmd

import (
	"errors"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/sdslabs/pinger/pkg/appcontext"
	"github.com/sdslabs/pinger/pkg/components/agent"
	"github.com/sdslabs/pinger/pkg/config/configfile"
)

// agent defaults.
const (
	defaultAgentConfigPath             = "agent.yml"
	defaultAgentExporter               = "timescale"
	defaultAgentInterval               = 2 * time.Minute
	defaultAgentPort            uint16 = 9009
	defaultAgentStandaloneMode         = false
	defaultAgentMetricsHost            = "127.0.0.1"
	defaultAgentMetricsPort     uint16 = 5432
	defaultAgentMetricsUsername        = ""
	defaultAgentMetricsPassword        = ""
	defaultAgentMetricsDBName          = ""
	defaultAgentMetricsSSLMode         = true
)

// config keys and flags for agent.
const (
	keyAgentConfigPort             = "port"
	flagAgentConfigPort            = "port"
	keyAgentConfigStandalone       = "standalone"
	flagAgentConfigStandalone      = "standalone"
	keyAgentConfigInterval         = "interval"
	flagAgentConfigInterval        = "interval"
	keyAgentConfigMetricsBackend   = "metrics.backend"
	flagAgentConfigMetricsBackend  = "metrics-backend"
	keyAgentConfigMetricsHost      = "metrics.host"
	flagAgentConfigMetricsHost     = "metrics-host"
	keyAgentConfigMetricsPort      = "metrics.port"
	flagAgentConfigMetricsPort     = "metrics-port"
	keyAgentConfigMetricsUsername  = "metrics.username"
	flagAgentConfigMetricsUsername = "metrics-username"
	keyAgentConfigMetricsPassword  = "metrics.password"
	flagAgentConfigMetricsPassword = "metrics-password"
	keyAgentConfigMetricsDBName    = "metrics.db_name"
	flagAgentConfigMetricsDBName   = "metrics-db-name"
	keyAgentConfigMetricsSSLMode   = "metrics.ssl_mode"
	flagAgentConfigMetricsSSLMode  = "metrics-ssl-mode"
)

func newAgentCmd(ctx *appcontext.Context, v *viper.Viper) (*cobra.Command, error) {
	conf := configfile.Agent{}
	var confPath string

	cmd := &cobra.Command{
		Use:   "agent",
		Short: "Run pinger agent.",
		Long: `
Run pinger agent on the given host, this agent will expose a gRPC API
to accept checks to perform and execute that. The agent also has the ability
to run in a standalone mode where it does not run any GRPC server.`,
		PreRun: func(*cobra.Command, []string) {
			if err := initConfig(ctx, v, confPath, defaultAgentConfigPath, &conf); err != nil {
				if errors.Unwrap(err) == errReadConfig {
					ctx.Logger().
						WithError(err).
						Warnln("continuing without config file")
				} else {
					ctx.Logger().
						WithError(err).
						Fatalln("invalid config")
					return
				}
			}

			if conf.Interval <= 0 {
				conf.Interval = defaultAgentInterval
			}
		},
		Run: func(*cobra.Command, []string) {
			if err := agent.Run(ctx, &conf); err != nil {
				ctx.Logger().
					WithError(err).
					Fatalln("cannot run agent")
			}
		},
	}

	cmd.Flags().StringVarP(&confPath, "config", "c", defaultAgentConfigPath, "config file path for agent")

	cmd.Flags().Uint16P(flagAgentConfigPort, "p", defaultAgentPort, "port to expose agent API on")
	cmd.Flags().BoolP(flagAgentConfigStandalone, "s", defaultAgentStandaloneMode, "should agent run in standalone mode")
	cmd.Flags().String(flagAgentConfigMetricsBackend, defaultAgentExporter, "backend service to store metrics")
	cmd.Flags().String(flagAgentConfigMetricsHost, defaultAgentMetricsHost, "host to run metrics server")
	cmd.Flags().Uint16(flagAgentConfigMetricsPort, defaultAgentMetricsPort, "port to run metrics server on")
	cmd.Flags().String(flagAgentConfigMetricsUsername, defaultAgentMetricsUsername, "username credential for metrics")
	cmd.Flags().String(flagAgentConfigMetricsPassword, defaultAgentMetricsPassword, "password credential for metrics")
	cmd.Flags().String(flagAgentConfigMetricsDBName, defaultAgentMetricsDBName, "database name for metrics")
	cmd.Flags().Bool(flagAgentConfigMetricsSSLMode, defaultAgentMetricsSSLMode, "whether to run metrics with SSL")
	cmd.Flags().Duration(
		flagAgentConfigInterval,
		defaultAgentInterval,
		"interval after which metrics are pushed/pulled",
	)

	mapKeysToFlags := map[string]string{
		keyAgentConfigPort:            flagAgentConfigPort,
		keyAgentConfigStandalone:      flagAgentConfigStandalone,
		keyAgentConfigMetricsBackend:  flagAgentConfigMetricsBackend,
		keyAgentConfigMetricsHost:     flagAgentConfigMetricsHost,
		keyAgentConfigMetricsPort:     flagAgentConfigMetricsPort,
		keyAgentConfigMetricsUsername: flagAgentConfigMetricsUsername,
		keyAgentConfigMetricsPassword: flagAgentConfigMetricsPassword,
		keyAgentConfigMetricsDBName:   flagAgentConfigMetricsDBName,
		keyAgentConfigMetricsSSLMode:  flagAgentConfigMetricsSSLMode,
		keyAgentConfigInterval:        flagAgentConfigInterval,
	}

	if err := bindFlagsToViper(v, cmd, mapKeysToFlags); err != nil {
		return nil, err
	}

	return cmd, nil
}
