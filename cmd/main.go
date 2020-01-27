// Tool sp can be used to run the API server inside an agent.
package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(agentCmd)
	rootCmd.AddCommand(serverCmd)
}
