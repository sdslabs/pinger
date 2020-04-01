package main

import (
	"github.com/spf13/cobra"
)

var centralAPICmd = &cobra.Command{
	Use:   "central",
	Short: "Run status page central server.",
	Long:  "Run status page server on the given host, this server will use agents as workers to schedule periodic checks for status.",

	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
}
