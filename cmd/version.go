package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/sdslabs/pinger/pkg/util/appcontext"
)

// Various ldflags.
var (
	// version of the binary.
	version string
	// debug flag, if set to true, builds the binary with developer friendly
	// settings, like, debugging logs, etc.
	debug string
)

// debugMode is the parsed value of the "debug" ldflag as a boolean.
var debugMode = func() bool {
	mode, err := strconv.ParseBool(debug)
	if err != nil {
		return false
	}
	return mode
}()

// IsDebug returns true if command line is built in debug mode.
func IsDebug() bool {
	return debugMode
}

func newVersionCmd(*appcontext.Context, *viper.Viper) (*cobra.Command, error) {
	return &cobra.Command{
		Use:   "version",
		Short: "Prints the version of application",
		Run: func(*cobra.Command, []string) {
			fmt.Printf("Pinger %s\n", version)
			if IsDebug() {
				fmt.Println("Build: debug")
			}
		},
	}, nil
}
