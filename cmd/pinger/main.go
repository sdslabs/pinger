// Copyright (c) 2020 SDSLabs
// Use of this source code is governed by an MIT license
// details of which can be found in the LICENSE file.

// Command pinger is used to start the app-server, central-server or the
// agent to execute checks. It is used as the entry point for all components
// of the application.
package main

import (
	"os"

	"github.com/sdslabs/status/cmd/commands"
	"github.com/sdslabs/status/internal/appcontext"
)

func main() {
	// Parent context for the application.
	ctx, cancel := appcontext.WithSignals(
		appcontext.Background(),
		os.Interrupt, os.Kill,
	)
	defer cancel()

	cmd := commands.NewRootCmd(ctx)
	if err := cmd.Execute(); err != nil {
		ctx.Logger().
			WithError(err).
			Fatalln("cannot execute pinger command")
	}
}
