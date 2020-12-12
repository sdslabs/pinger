package httpserver

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/sdslabs/pinger/pkg/util/appcontext"
)

// ErrorResponse is the standard response for errors.
type ErrorResponse struct {
	Error string `json:"error"`
}

// MetricResponse is the JSON response for returning metrics.
type MetricResponse struct {
	Successful bool          `json:"successful"`
	Timeout    bool          `json:"timeout"`
	StartTime  time.Time     `json:"start_time"`
	Duration   time.Duration `json:"duration"`
}

// RespondError writes an error response to the request.
func RespondError(ctx *appcontext.Context, c *gin.Context, statusCode int, err error) {
	resp := err.Error()
	if statusCode == http.StatusInternalServerError && !ctx.Debug() {
		ctx.Logger().WithError(err).WithField("path", c.Request.URL.Path).Errorln()
		resp = http.StatusText(http.StatusInternalServerError)
	}
	c.JSON(statusCode, ErrorResponse{Error: resp})
}

// RespondErrorInternalServer is same as calling RespondError(ctx, c, 500, err).
func RespondErrorInternalServer(ctx *appcontext.Context, c *gin.Context, err error) {
	RespondError(ctx, c, http.StatusInternalServerError, err)
}

// RespondErrorNotFound is same as calling RespondError(ctx, c, 404, error).
func RespondErrorNotFound(ctx *appcontext.Context, c *gin.Context, err error) {
	RespondError(ctx, c, http.StatusNotFound, err)
}

// RespondOK writes an OK response to the request.
func RespondOK(_ *appcontext.Context, c *gin.Context, resp interface{}) {
	c.JSON(http.StatusOK, resp)
}
