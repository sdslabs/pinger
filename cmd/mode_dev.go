//go:build dev
// +build dev

package cmd

// dev flag is set to true when binary is built in development mode.
const dev = true

func init() {
	version = "(dev)"
}
