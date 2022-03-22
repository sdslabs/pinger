package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/sdslabs/pinger/pkg/util/appcontext"
)

// NewRootCmd creates a new root command for pinger.
func NewRootCmd(ctx *appcontext.Context) (*cobra.Command, error) {
	v := viper.New()

	cmd := &cobra.Command{
		Use:   "pinger",
		Short: "Entry-point for pinger CLI",
		Long: `
Command pinger is used to start the app-server, central-server or the
agent to execute checks. It is used as the entry point for all components
of the application.`,
		Run: func(cmd *cobra.Command, _ []string) {
			if err := cmd.Help(); err != nil {
				ctx.Logger().
					WithError(err).
					Fatalln("cannot execute pinger command")
			}
		},
	}

	// add various commands
	if err := addCommands(ctx, v, cmd,
		// Add commands here
		newCentralCmd,
		newAgentCmd,
		newVersionCmd,
		newListCommand,
	); err != nil {
		return nil, err
	}

	return cmd, nil
}
