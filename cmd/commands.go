package main

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sp",
	Short: "Status page is an open source implementation of ping based status page.",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}
