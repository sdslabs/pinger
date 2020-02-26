// Tool status can be used to execute various api servers (application and central)
// and expose the agent API inside an agent.
package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func initConfig(confPath, defaultConfPath string, resolveTo interface{}) {
	if confPath != "" {
		viper.SetConfigFile(confPath)
	} else {
		viper.SetConfigFile(defaultConfPath)
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Cannot read config file: %s", err.Error())
	}

	if err := viper.Unmarshal(resolveTo); err != nil {
		log.Fatalf("Cannot resolve config file: %s", err.Error())
	}
}

func viperErr(err error) {
	log.Errorf("Cannot bind flag with viper: %s", err.Error())
}

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
	rootCmd.AddCommand(centralAPICmd)
	rootCmd.AddCommand(appAPICmd)
}
