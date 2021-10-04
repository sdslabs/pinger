package cmd

import (
	"fmt"
	"strings"

	"github.com/sdslabs/pinger/pkg/alerter"
	"github.com/sdslabs/pinger/pkg/checker"
	"github.com/sdslabs/pinger/pkg/exporter"
	"github.com/sdslabs/pinger/pkg/oauther"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/sdslabs/pinger/pkg/util/appcontext"
)

// different keys for plugins.
const (
	keyPluginAlerters  = "alerters"
	keyPluginChckers   = "checkers"
	keyPluginExporters = "exporters"
	keyPluginOauthers  = "oauthers"
)

func newListCommand(ctx *appcontext.Context, v *viper.Viper) (*cobra.Command, error) {
	var (
		alerters  = alerter.List()
		checkers  = checker.List()
		exporters = exporter.List()
		oauthers  = oauther.List()
	)

	printPlugins := func(plugin string, list []string) {
		text := fmt.Sprintf("\033[1m%s\033[0m\n", plugin) // makes text bold
		for _, p := range list {
			text += fmt.Sprintf("  %s\n", p)
		}
		fmt.Println(text)
	}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all the installed plugins",
		Run: func(*cobra.Command, []string) {
			printPlugins(keyPluginAlerters, alerters)
			printPlugins(keyPluginChckers, checkers)
			printPlugins(keyPluginExporters, exporters)
			printPlugins(keyPluginOauthers, oauthers)
		},
	}

	if err := addCommands(ctx, v, cmd,
		// different kinds of plugins
		newListSubCommandCreator(keyPluginAlerters, alerters),
		newListSubCommandCreator(keyPluginChckers, checkers),
		newListSubCommandCreator(keyPluginExporters, exporters),
		newListSubCommandCreator(keyPluginOauthers, oauthers),
	); err != nil {
		return nil, err
	}

	return cmd, nil
}

// newListSubCommandCreator creates a command to list specific plugins.
func newListSubCommandCreator(plugin string, list []string) newCmdFunc {
	return func(*appcontext.Context, *viper.Viper) (*cobra.Command, error) {
		return &cobra.Command{
			Use:   plugin,
			Short: "Lists all the installed " + plugin,
			Run: func(*cobra.Command, []string) {
				fmt.Println(strings.Join(list, "\n"))
			},
		}, nil
	}
}
