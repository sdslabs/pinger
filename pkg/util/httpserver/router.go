package httpserver

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/sdslabs/pinger/pkg/util/appcontext"
)

func init() {
	gin.SetMode(gin.ReleaseMode) // Not really helpful debug messages
}

// RouterOpts are the options used in creating a new router.
type RouterOpts struct {
	AllowedOrigins []string // use "*" to allow all origins.
	AllowedMethods []string // GET and POST are allowed by default.

	allowAllOrigins bool
}

// defaults updates the options with their default values.
func (r *RouterOpts) defaults(ctx *appcontext.Context) {
	if ctx.Debug() {
		r.allowAllOrigins = true
		r.AllowedOrigins = nil
	}

	if len(r.AllowedMethods) == 0 {
		r.AllowedMethods = append(r.AllowedMethods, http.MethodGet, http.MethodPost)
	}
}

// NewRouter returns an instance of the gin router with the required
// configurations and setup for the application.
func NewRouter(ctx *appcontext.Context, opts RouterOpts) *gin.Engine {
	opts.defaults(ctx)

	corsConf := cors.Config{
		AllowAllOrigins: opts.allowAllOrigins,
		AllowMethods:    opts.AllowedMethods,
		AllowOrigins:    opts.AllowedOrigins,
		AllowWildcard:   true,
	}

	router := gin.Default()
	router.Use(cors.New(corsConf))

	return router
}
