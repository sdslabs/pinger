package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/sdslabs/pinger/pkg/util/appcontext"
)

var (
	// debug variable contains the information about if the binary is built
	// with in debug mode.
	debug string
	// debugMode stores the debug variable information as a boolean.
	debugMode = isDebug()
)

// IsDebug returns true if command line is built in debug mode.
func IsDebug() bool {
	return debugMode
}

func isDebug() bool {
	mode, err := strconv.ParseBool(debug)
	if err != nil {
		return false
	}
	return mode
}

func newModeCmd(*appcontext.Context, *viper.Viper) (*cobra.Command, error) {
	return &cobra.Command{
		Use:   "mode",
		Short: "Prints which mode the binary is compiled in",
		Run: func(*cobra.Command, []string) {
			mode := "release"
			if IsDebug() {
				mode = "debug"
			}
			fmt.Println(mode)
		},
	}, nil
}
