package httpserver

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/sdslabs/pinger/pkg/util/appcontext"
)

// ListenAndServe listens on :port for requests and serves them using the
// HTTP handler. It honors context and hence exits gracefully when the
// context is canceled.
func ListenAndServe(ctx *appcontext.Context, port uint16, h http.Handler) error {
	addr := net.JoinHostPort("0.0.0.0", fmt.Sprint(port))
	server := http.Server{
		Addr:    addr,
		Handler: h,
	}

	errChan := make(chan error)
	go func() {
		errChan <- server.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		shutDownCtx, shutDownCancel := context.WithTimeout(
			context.Background(),
			10*time.Second, // wait for 10 seconds else force shutdown
		)
		defer shutDownCancel()
		if err := server.Shutdown(shutDownCtx); err != nil {
			return server.Close()
		}
		return ctx.Err()
	case err := <-errChan:
		return err
	}
}
