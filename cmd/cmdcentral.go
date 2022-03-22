package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/sdslabs/pinger/pkg/components/central"
	"github.com/sdslabs/pinger/pkg/util/appcontext"
)

func newCentralCmd(ctx *appcontext.Context, _ *viper.Viper) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "central",
		Short: "Run pinger central server.",
		Run: func(*cobra.Command, []string) {
			if err := central.Run(ctx); err != nil {
				ctx.Logger().
					WithError(err).
					Fatalln("cannot run central server")
			}
		},
	}

	return cmd, nil
}
