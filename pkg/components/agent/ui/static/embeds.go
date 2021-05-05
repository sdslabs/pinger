package static

import "embed"

// StaticFS is the filesystem with embedded static content.
//go:embed *.png *.js *.css
var StaticFS embed.FS
