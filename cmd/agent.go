package cmd

import (
	"errors"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/sdslabs/pinger/pkg/components/agent"
	"github.com/sdslabs/pinger/pkg/config/configfile"
	"github.com/sdslabs/pinger/pkg/util/appcontext"
)

// agent defaults.
const (
	defaultAgentConfigPath             = "agent.yml"
	defaultAgentExporter               = "timescale"
	defaultAgentInterval               = 2 * time.Minute
	defaultAgentPort            uint16 = 9009
	defaultAgentStandaloneMode         = false
	defaultAgentMetricsHost            = "127.0.0.1"
	defaultAgentMetricsPort     uint16 = 0
	defaultAgentMetricsUsername        = ""
	defaultAgentMetricsPassword        = ""
	defaultAgentMetricsOrgName         = ""
	defaultAgentMetricsDBName          = ""
	defaultAgentMetricsSSLMode         = true
	defaultAgentPageDeploy             = false
	defaultAgentPagePort        uint16 = 9010
	defaultAgentPageName               = "Status Page"
	defaultAgentPageMedia              = ""
	defaultAgentPageLogo               = ""
	defaultAgentPageFavicon            = ""
	defaultAgentPageWebsite            = "/"
)

// non-const agent defaults.
var defaultAgentPageAllowedOrigins []string = nil

// config keys and flags for agent.
const (
	keyAgentConfigPort                = "port"
	flagAgentConfigPort               = "port"
	keyAgentConfigStandalone          = "standalone"
	flagAgentConfigStandalone         = "standalone"
	keyAgentConfigInterval            = "interval"
	flagAgentConfigInterval           = "interval"
	keyAgentConfigMetricsBackend      = "metrics.backend"
	flagAgentConfigMetricsBackend     = "metrics-backend"
	keyAgentConfigMetricsHost         = "metrics.host"
	flagAgentConfigMetricsHost        = "metrics-host"
	keyAgentConfigMetricsPort         = "metrics.port"
	flagAgentConfigMetricsPort        = "metrics-port"
	keyAgentConfigMetricsUsername     = "metrics.username"
	flagAgentConfigMetricsUsername    = "metrics-username"
	keyAgentConfigMetricsPassword     = "metrics.password"
	flagAgentConfigMetricsPassword    = "metrics-password"
	keyAgentConfigMetricsOrgName      = "metrics.org_name"
	flagAgentConfigMetricsOrgName     = "metrics-org-name"
	keyAgentConfigMetricsDBName       = "metrics.db_name"
	flagAgentConfigMetricsDBName      = "metrics-db-name"
	keyAgentConfigMetricsSSLMode      = "metrics.ssl_mode"
	flagAgentConfigMetricsSSLMode     = "metrics-ssl-mode"
	keyAgentConfigPageDeploy          = "page.deploy"
	flagAgentConfigPageDeploy         = "page-deploy"
	keyAgentConfigPagePort            = "page.port"
	flagAgentConfigPagePort           = "page-port"
	keyAgentConfigPageAllowedOrigins  = "page.allowed_origins"
	flagAgentConfigPageAllowedOrigins = "page-allowed-origins"
	keyAgentConfigPageName            = "page.name"
	flagAgentConfigPageName           = "page-name"
	keyAgentConfigPageMedia           = "page.media"
	flagAgentConfigPageMedia          = "page-media"
	keyAgentConfigPageLogo            = "page.logo"
	flagAgentConfigPageLogo           = "page-logo"
	keyAgentConfigPageFavicon         = "page.favicon"
	flagAgentConfigPageFavicon        = "page-favicon"
	keyAgentConfigPageWebsite         = "page.website"
	flagAgentConfigPageWebsite        = "page-website"
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
	cmd.Flags().String(flagAgentConfigMetricsOrgName, defaultAgentMetricsOrgName, "organization name for metrics")
	cmd.Flags().String(flagAgentConfigMetricsDBName, defaultAgentMetricsDBName, "database name for metrics")
	cmd.Flags().Bool(flagAgentConfigMetricsSSLMode, defaultAgentMetricsSSLMode, "whether to run metrics with SSL")
	cmd.Flags().Duration(
		flagAgentConfigInterval, defaultAgentInterval, "interval after which metrics are pushed/pulled")
	cmd.Flags().Bool(flagAgentConfigPageDeploy, defaultAgentPageDeploy, "whether to deploy agent-only status page")
	cmd.Flags().Uint16(flagAgentConfigPagePort, defaultAgentPagePort, "port to deploy status page on")
	cmd.Flags().StringSlice(
		flagAgentConfigPageAllowedOrigins, defaultAgentPageAllowedOrigins, "allowed origins which can request page")
	cmd.Flags().String(flagAgentConfigPageName, defaultAgentPageName, "name/title of status page")
	cmd.Flags().String(flagAgentConfigPageMedia, defaultAgentPageMedia, "directory to serve page media from")
	cmd.Flags().String(flagAgentConfigPageLogo, defaultAgentPageLogo, "filename for logo in media directory")
	cmd.Flags().String(flagAgentConfigPageFavicon, defaultAgentPageFavicon, "filename for favicon in media directory")
	cmd.Flags().String(flagAgentConfigPageWebsite, defaultAgentPageWebsite, "website url for the page")

	mapKeysToFlags := map[string]string{
		keyAgentConfigPort:               flagAgentConfigPort,
		keyAgentConfigStandalone:         flagAgentConfigStandalone,
		keyAgentConfigMetricsBackend:     flagAgentConfigMetricsBackend,
		keyAgentConfigMetricsHost:        flagAgentConfigMetricsHost,
		keyAgentConfigMetricsPort:        flagAgentConfigMetricsPort,
		keyAgentConfigMetricsUsername:    flagAgentConfigMetricsUsername,
		keyAgentConfigMetricsPassword:    flagAgentConfigMetricsPassword,
		keyAgentConfigMetricsOrgName:     flagAgentConfigMetricsOrgName,
		keyAgentConfigMetricsDBName:      flagAgentConfigMetricsDBName,
		keyAgentConfigMetricsSSLMode:     flagAgentConfigMetricsSSLMode,
		keyAgentConfigInterval:           flagAgentConfigInterval,
		keyAgentConfigPageDeploy:         flagAgentConfigPageDeploy,
		keyAgentConfigPagePort:           flagAgentConfigPagePort,
		keyAgentConfigPageAllowedOrigins: flagAgentConfigPageAllowedOrigins,
		keyAgentConfigPageName:           flagAgentConfigPageName,
		keyAgentConfigPageMedia:          flagAgentConfigPageMedia,
		keyAgentConfigPageLogo:           flagAgentConfigPageLogo,
		keyAgentConfigPageFavicon:        flagAgentConfigPageFavicon,
		keyAgentConfigPageWebsite:        flagAgentConfigPageWebsite,
	}

	if err := bindFlagsToViper(v, cmd, mapKeysToFlags); err != nil {
		return nil, err
	}

	return cmd, nil
}
