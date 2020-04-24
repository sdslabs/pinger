// Tool status can be used to execute various api servers (application and central)
// and expose the agent API inside an agent.
package main

import (
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
		// Just warn while reading the config file since configuration
		// can also be mostly passed via flags.
		log.WithError(err).Warnln("cannot read config file")
	}

	if err := viper.Unmarshal(resolveTo); err != nil {
		log.WithError(err).Fatalln("cannot resolve config file")
	}
}

func viperErr(err error) {
	log.WithError(err).Fatalln("error binding flags with viper")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.WithError(err).Fatalln("cannot start cmd")
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
