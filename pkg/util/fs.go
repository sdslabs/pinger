package util

import (
	"io/fs"
	"os"
	"path/filepath"
)

// fsys is the implementation of FileSystem interface.
type fsys struct {
	prefix string
}

// getPath returns the path using the prefix provided while creating the file
// system instance.
func (f *fsys) getPath(name string) string {
	return filepath.Join(f.prefix, name)
}

// Open implements fs.FS.
func (f *fsys) Open(name string) (fs.File, error) {
	return os.Open(f.getPath(name))
}

// NewFS creates a raw file system from OS.
func NewFS(prefix string) fs.FS {
	return &fsys{prefix: prefix}
}
