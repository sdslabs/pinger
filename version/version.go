// Package version populates the version information at the time of build
// and allocates it to the `VersionStr` format string. This is then used
// to print and display the when `version` command is run.
package version

// Build information. Populated at build-time by the build script.
var (
	Version   string
	Revision  string
	Branch    string
	BuildUser string
	BuildDate string
	GoVersion string
)

// Info provides the iterable version information.
var Info = map[string]string{
	"version":   Version,
	"revision":  Revision,
	"branch":    Branch,
	"buildUser": BuildUser,
	"buildDate": BuildDate,
	"goVersion": GoVersion,
}

// VersionStr is the format of string used for version message.
var VersionStr = `************** SDS Status ****************
	Version    : %s
	Revision   : %s
	Branch     : %s
	Build-User : %s
	Build-Date : %s
	Go-Version : %s
******************************************
`
