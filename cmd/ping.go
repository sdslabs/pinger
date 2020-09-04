// Copyright (c) 2020 SDSLabs
// Use of this source code is governed by an MIT license
// details of which can be found in the LICENSE file.

package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/sdslabs/pinger/pkg/appcontext"
)

func newPingCmd(ctx *appcontext.Context, v *viper.Viper) (*cobra.Command, error) {
	return &cobra.Command{
		Use:   "ping",
		Short: "Pongs if the binary is built fine",
		Run: func(*cobra.Command, []string) {
			logrus.Infoln("pong")
		},
	}, nil
}
