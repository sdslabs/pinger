package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/sdslabs/pinger/pkg/util/appcontext"
)

// version of the binary.
var version string

func newVersionCmd(ctx *appcontext.Context, _ *viper.Viper) (*cobra.Command, error) {
	logger := ctx.Logger()
	return &cobra.Command{
		Use:   "version",
		Short: "Prints the version of application",
		Run: func(*cobra.Command, []string) {
			fmt.Printf("Pinger %s\n", version)
			logger.Debug("Debugging is enabled")
			logger.Trace("Tracing is enabled")
		},
	}, nil
}
