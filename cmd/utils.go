package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/sdslabs/pinger/pkg/util/appcontext"
)

// Errors that occur while init-ing viper.
var (
	errReadConfig      = fmt.Errorf("cannot read config file")
	errUnmarshalConfig = fmt.Errorf("cannot unmarshal config file")
	errBindFlags       = fmt.Errorf("cannot bind flags to viper")
)

// initConfig reads and unmarshals the config.
func initConfig(ctx *appcontext.Context, v *viper.Viper, confPath, defaultPath string, resolveTo interface{}) error {
	if confPath != "" {
		v.SetConfigFile(confPath)
	} else {
		v.SetConfigFile(defaultPath)
	}

	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("%w: %v", errReadConfig, err)
	}

	if err := v.Unmarshal(resolveTo); err != nil {
		return fmt.Errorf("%w: %v", errUnmarshalConfig, err)
	}

	return nil
}

// bindFlagsToViper binds the values of flags to viper config. This is just
// to override the config values set from flags.
func bindFlagsToViper(v *viper.Viper, cmd *cobra.Command, keyFlagMap map[string]string) error {
	for key, flag := range keyFlagMap {
		if err := v.BindPFlag(key, cmd.Flags().Lookup(flag)); err != nil {
			return fmt.Errorf("%w: %v", errBindFlags, err)
		}
	}

	return nil
}

// newCmdFunc is the function used to create a new subcommand to the root
// command.
type newCmdFunc = func(*appcontext.Context, *viper.Viper) (*cobra.Command, error)

// addCommands adds the commands to the root specified.
func addCommands(ctx *appcontext.Context, v *viper.Viper, root *cobra.Command, fns ...newCmdFunc) error {
	for _, fn := range fns {
		cmd, err := fn(ctx, v)
		if err != nil {
			return err
		}

		root.AddCommand(cmd)
	}

	return nil
}
