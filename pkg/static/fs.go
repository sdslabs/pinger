package static

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/phogolabs/parcello"
)

// debug env constants.
const (
	debugEnv = "DEBUG"
	trueVal  = "on"
	baseDir  = "./static" // assuming commands will be run from root dir.
)

// fsImpl implements the parcello.FileSystemManager interface for actual FS.
type fsImpl string

// getPath returns the complete path of the file or directory.
func (f fsImpl) getPath(name string) (string, error) {
	if filepath.IsAbs(filepath.Clean(name)) {
		return "", fmt.Errorf("path cannot be absolute: %q", name)
	}

	return filepath.Join(baseDir, string(f), name), nil
}

// Open implements the filepath.Open function for http.FileSystem.
func (f fsImpl) Open(name string) (http.File, error) {
	return f.OpenFile(name, os.O_RDONLY, 0 /* perm */)
}

// Walk implements the filepath.Walk function.
func (f fsImpl) Walk(dir string, fn filepath.WalkFunc) error {
	actualPath, err := f.getPath(dir)
	if err != nil {
		return err
	}

	return filepath.Walk(actualPath, fn)
}

// OpenFile implements the os.OpenFile method.
func (f fsImpl) OpenFile(name string, flag int, perm os.FileMode) (parcello.File, error) {
	actualPath, err := f.getPath(name)
	if err != nil {
		return nil, err
	}

	return os.OpenFile(actualPath, flag, perm) // nolint:gosec
}

// debug tells if the application is run in debug mode.
func debug() bool {
	env := strings.TrimSpace(os.Getenv(debugEnv))
	return env == trueVal
}

// NewFS returns the file system manager for the given path based on the
// environment, i.e., is application running in debug or not.
func NewFS(path string) (parcello.FileSystem, error) {
	fpath := filepath.Clean(path)
	if filepath.IsAbs(fpath) {
		return nil, fmt.Errorf(
			"path cannot be absolute; should be relative to static dir: %s", path,
		)
	}

	if !debug() {
		return parcello.Manager.Dir(path)
	}

	return fsImpl(fpath), nil
}

// Interface guard.
var _ parcello.FileSystem = (*fsImpl)(nil)
