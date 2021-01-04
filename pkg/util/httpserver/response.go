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

// PageCheckMetricsResponse is the JSON response for all the metrics related
// to a particular check.
type PageCheckMetricsResponse struct {
	Metrics     []MetricResponse `json:"metrics"`
	Uptime      int              `json:"uptime"`
	Operational bool             `json:"operational"`
}

// PageMetricsResponse is the JSON response for returning all the metrics
// related information which is fetched once the page is loaded.
type PageMetricsResponse struct {
	ChecksDown int                                 `json:"checks_down"`
	Checks     map[string]PageCheckMetricsResponse `json:"checks"`
}

// PageResponse is the data passed into the template.
type PageResponse struct {
	Name       string
	Checks     map[string]string
	StaticURL  string
	MetricsURL string
	LogoURL    string
	FaviconURL string
	WebsiteURL string
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
