package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/sdslabs/status/version"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Displays the version of the current build of status page",
	Long: `Displays the version of the current build of status page, this information
include Version, Revision, Git-Branch, BuildUser, BuildDate, go-version`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf(version.VersionStr,
			version.Info["version"],
			version.Info["revision"],
			version.Info["branch"],
			version.Info["buildUser"],
			version.Info["buildDate"],
			version.Info["goVersion"])
		os.Exit(0)
	},
}
