// +build !dev

package static

import (
	"embed"
	"io/fs"
)

//go:embed *.png *.js *.css
var fsys embed.FS

// FS is the filesystem with embedded static content.
var FS fs.FS = fsys
