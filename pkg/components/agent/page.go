package agent

import (
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/sdslabs/pinger/pkg/checker"

	"github.com/gin-gonic/gin"

	"github.com/sdslabs/pinger/pkg/config/configfile"
	"github.com/sdslabs/pinger/pkg/exporter"
	"github.com/sdslabs/pinger/pkg/util/appcontext"
	"github.com/sdslabs/pinger/pkg/util/controller"
	"github.com/sdslabs/pinger/pkg/util/httpserver"
	"github.com/sdslabs/pinger/pkg/util/static"
)

const (
	// agentFilePath is the directory where static content is stored specific
	// to the agent.
	agentFilePath = "/agent"

	// staticFilesPath is the directory where static content is stored which
	// can be publicly accessed.
	staticFilesPath = agentFilePath + "/public"

	// templateName is the name of the template for agent's status page.
	templateName = "status_page.tmpl"

	// maxMetricsDuration is the maximum duration for which the metrics are
	// fetched from the database for displaying on status page.
	maxMetricsDuration = 7 * (24 * time.Hour) // 1 week

	// maxMetricsBatches is the maximum possible number of batches in which the
	// metrics are to be distributed when the client requests.
	maxMetricsBatches = 50

	// routeStatic is the route for all static content.
	routeStatic = "/static"

	// routeMetrics is the route for fetching metrics.
	routeMetrics = "/metrics"
)

// serveStatusPage starts a HTTP server that responds with the status page
// of checks listed in the agent's config.
func serveStatusPage(
	ctx *appcontext.Context,
	conf *configfile.AgentPage,
	manager *controller.Manager,
	getMetrics exporter.GetterFunc,
) error {
	if !conf.Deploy {
		return errors.New("cannot deploy agent with AgentPage.Deploy=false")
	}

	router := httpserver.NewRouter(ctx, httpserver.RouterOpts{
		AllowedOrigins: conf.AllowedOrigins,
		AllowedMethods: []string{http.MethodGet},
	})

	if err := addBaseRoute(ctx, router, manager, conf.Name); err != nil {
		return err
	}

	addMetricsRoute(ctx, router, manager, getMetrics)

	staticFS, err := static.NewHTTPFS(ctx, staticFilesPath)
	if err != nil {
		return err
	}
	router.StaticFS(routeStatic, staticFS)

	go func() {
		ctx.Logger().
			WithField("address", fmt.Sprintf(":%d", conf.Port)).
			Infof("serving status page")
		if err := httpserver.ListenAndServe(ctx, conf.Port, router); err != nil {
			ctx.Logger().WithError(err).Errorln("server exited unexpectedly")
		}
	}()

	return nil
}

// addBaseRoute adds the route that returns template for status page.
func addBaseRoute(
	ctx *appcontext.Context,
	router *gin.Engine,
	manager *controller.Manager,
	name string,
) error {
	agentFS, err := static.NewFS(ctx, agentFilePath)
	if err != nil {
		return err
	}

	tmplFile, err := agentFS.Open(templateName)
	if err != nil {
		return err
	}
	defer tmplFile.Close() // nolint:errcheck

	templateContent, err := ioutil.ReadAll(tmplFile)
	if err != nil {
		return err
	}

	compiledTemplate, err := template.New(templateName).Parse(string(templateContent))
	if err != nil {
		return err
	}

	type resp struct {
		Name       string
		Checks     map[string]string
		StaticURL  string
		MetricsURL string
	}

	router.SetHTMLTemplate(compiledTemplate)
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, templateName, resp{
			Name:       name,
			Checks:     manager.ListControllers(),
			StaticURL:  routeStatic,
			MetricsURL: routeMetrics,
		})
	})

	return nil
}

// addMetricsRoute adds the route that fetches metrics for all the checks
// running on the agent.
func addMetricsRoute(
	ctx *appcontext.Context,
	router *gin.Engine,
	manager *controller.Manager,
	getMetrics exporter.GetterFunc,
) {
	type result struct {
		Metrics     []httpserver.MetricResponse `json:"metrics"`
		Uptime      int                         `json:"uptime"`
		Operational bool                        `json:"operational"`
	}

	// `/metrics?duration=1000000000?batches=30`
	router.GET(routeMetrics, func(c *gin.Context) {
		durationStr := c.Query("duration")
		durationInt, err := strconv.Atoi(durationStr)
		duration := time.Duration(durationInt)
		if err != nil || duration <= 0 || duration > maxMetricsDuration {
			duration = maxMetricsDuration
		}

		batchesStr := c.Query("batches")
		batches, err := strconv.Atoi(batchesStr)
		if err != nil || batches <= 0 || batches > maxMetricsBatches {
			batches = maxMetricsBatches
		}

		checksMap := manager.ListControllers()
		checkIDs := make([]string, 0, len(checksMap))
		for checkID := range checksMap {
			checkIDs = append(checkIDs, checkID)
		}
		metrics, err := getMetrics(ctx, duration, checkIDs...)
		if err != nil {
			httpserver.RespondErrorInternalServer(ctx, c, err)
			return
		}

		resp := map[string]result{}
		for cid := range metrics {
			// NB: This shouldn't take that long but since this request is long enough
			// in general, one optimization can be to serialize the metrics in different
			// goroutines. For now this works and we need benchmarks to prove if making
			// this change would really help or not.
			serialized, uptime := serializeMetrics(batches, metrics[cid])
			if len(serialized) == 0 || len(metrics[cid]) == 0 {
				continue
			}
			operational := serialized[0].Successful
			resp[cid] = result{
				Metrics:     serialized,
				Uptime:      uptime,
				Operational: operational,
			}
		}

		httpserver.RespondOK(ctx, c, resp)
	})
}

// serializeMetrics breaks the metrics into multiple batches and retains one
// metric from each batch.
//
// The following rules are applied to each batch:
// 	- Failed metric is prioritized over successful.
// 	- Metric with highest latency is considered.
//  - If number of batches are more than number of metrics, this is probably
// 	  recent addition of check. In this case, The first metric should be
// 	  appended at the front of list.
//
func serializeMetrics(
	batches int, metrics []checker.Metric,
) (serialized []httpserver.MetricResponse, uptime int) {
	if batches == 0 || len(metrics) == 0 {
		return
	}

	serialized = make([]httpserver.MetricResponse, 0, batches)
	numEachBatch := (len(metrics) / batches) + 1
	var upNum int

	for i := 0; i < len(metrics); i += numEachBatch {
		var (
			metric  checker.Metric
			latency time.Duration
			failed  bool
		)

		for j := i; j < i+numEachBatch; j++ {
			if j >= len(metrics) {
				break
			}

			m := metrics[j]

			if !m.IsSuccessful() && !failed {
				metric = m
				failed = true
			}

			if m.IsSuccessful() {
				upNum++
			}

			if failed {
				continue // don't break because we need to calculate uptime
			}

			if latency < m.GetDuration() {
				metric = m
			}
		}

		if metric == nil {
			break
		}

		serialized = append(serialized, httpserver.MetricResponse{
			Successful: metric.IsSuccessful(),
			Timeout:    metric.IsTimeout(),
			StartTime:  metric.GetStartTime(),
			Duration:   metric.GetDuration(),
		})
	}

	if len(serialized) > 0 {
		// Since metrics are ordered in descending order of their start times we need
		// to replicate the last metric so length of serialized equals the number of
		// batches we need to divide the data in.
		lastMetric := serialized[len(serialized)-1]
		for len(serialized) < batches {
			serialized = append(serialized, lastMetric)
		}
	}

	uptime = (upNum * 100) / len(metrics) // percentage

	return
}
