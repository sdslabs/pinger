//go:build dev
// +build dev

package static

import (
	"path"

	"github.com/sdslabs/pinger/pkg/util"
)

// FS is the filesystem with embedded static content.
var FS = util.NewFS(path.Join("pkg", "components", "agent", "ui", "static"))
