package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/sdslabs/pinger/pkg/appcontext"
)

// contains version information.
var version string

func newVersionCmd(*appcontext.Context, *viper.Viper) (*cobra.Command, error) {
	return &cobra.Command{
		Use:   "version",
		Short: "Prints the version of application",
		Run: func(*cobra.Command, []string) {
			fmt.Println("pinger", version)
		},
	}, nil
}
