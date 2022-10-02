package cmd

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/sdslabs/pinger/pkg/components/central"
	"github.com/sdslabs/pinger/pkg/config/configfile"
	"github.com/sdslabs/pinger/pkg/util/appcontext"
)

const (
	defaultCheckDiffPath = "checkdiff.yml"
)

func newCentralCmd(ctx *appcontext.Context, v *viper.Viper) (*cobra.Command, error) {
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

	if err := addCommands(ctx, v, cmd,
		// Sub-commands for central server
		newCentralSubCmdApply,
	); err != nil {
		return nil, err
	}

	return cmd, nil
}

func newCentralSubCmdApply(ctx *appcontext.Context, v *viper.Viper) (*cobra.Command, error) {
	checkdiff := configfile.CheckDiff{}
	var checkdiffPath string

	cmd := &cobra.Command{
		Use:   "apply",
		Short: "Apply the given check diff",
		PreRun: func(*cobra.Command, []string) {
			if err := initConfig(ctx, v, checkdiffPath, defaultCheckDiffPath, &checkdiff); err != nil {
				if errors.Unwrap(err) == errReadConfig {
					ctx.Logger().
						WithError(err).
						Fatalln("invalid check diff")
					return
				}
			}
		},
		Run: func(*cobra.Command, []string) {
			if err := central.RunApply(ctx, &checkdiff); err != nil {
				ctx.Logger().
					WithError(err).
					Fatalln("could not apply the check diff")
			}
		},
	}

	cmd.Flags().StringVarP(&checkdiffPath, "checkdiff", "c", defaultCheckDiffPath, "check diff file to apply")

	return cmd, nil
}
