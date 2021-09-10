package static

import (
	"net/http"
	"os"

	"github.com/phogolabs/parcello"

	"github.com/sdslabs/pinger/pkg/util/appcontext"
)

type httpFS struct {
	ctx *appcontext.Context
	fs  parcello.FileSystem
}

// Open implements the http.FileSystem interface.
func (h *httpFS) Open(name string) (http.File, error) {
	f, err := h.fs.Open(name)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if err != nil {
		return nil, err
	}

	if s.IsDir() && !h.ctx.Debug() {
		return nil, os.ErrNotExist
	}

	return f, nil
}

// NewHTTPFS returns the HTTP file system that blocks dir listings when run
// in production mode.
func NewHTTPFS(ctx *appcontext.Context, path string) (http.FileSystem, error) {
	fs, err := NewFS(ctx, path)
	if err != nil {
		return nil, err
	}

	return &httpFS{ctx: ctx, fs: fs}, nil
}

// Interface guard.
var _ http.FileSystem = (*httpFS)(nil)
