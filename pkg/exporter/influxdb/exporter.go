package influxdb

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	client "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"

	"github.com/sirupsen/logrus"

	"github.com/sdslabs/pinger/pkg/checker"
	"github.com/sdslabs/pinger/pkg/config"
	"github.com/sdslabs/pinger/pkg/exporter"
	"github.com/sdslabs/pinger/pkg/util/appcontext"
)

const exporterName = "influxdb"

// metric keys
const (
	keyCheckID      = "check_id"
	keyCheckName    = "check_name"
	keyIsSuccessful = "is_successful"
	keyIsTimeout    = "is_timeout"
	keyStartTime    = "start_time"
	keyDuration     = "duration"
)

// Exporter for exporting metrics to influxdb.
type Exporter struct {
	writeAPI api.WriteAPIBlocking
	queryAPI api.QueryAPI
	dbname   string
	log      *logrus.Logger
}

func init() {
	exporter.Register(exporterName, func() exporter.Exporter { return new(Exporter) })
}

// newClient creates a new client for the database.
func newClient(_ *appcontext.Context, provider exporter.Provider) (client.Client, error) {
	protocol := "http"
	if provider.IsSSLMode() {
		protocol = "https"
	}
	addStr := fmt.Sprintf("%s://%s",
		protocol,
		net.JoinHostPort(
			provider.GetHost(),
			strconv.Itoa(int(provider.GetPort())),
		),
	)

	authStr := fmt.Sprintf("%s:%s", provider.GetUsername(), provider.GetPassword())
	if provider.GetUsername() == "" {
		authStr = provider.GetPassword() // in case using token authentication
	}

	c := client.NewClientWithOptions(
		addStr,
		authStr,
		client.DefaultOptions(),
	)

	return c, nil
}

// Export exports the metrics to the exporter.
func (e *Exporter) Export(ctx context.Context, metrics []checker.Metric) error {
	if len(metrics) == 0 {
		return nil
	}

	points := make([]*write.Point, 0, len(metrics))
	for _, metric := range metrics {
		tags := map[string]string{
			keyCheckID: metric.GetCheckID(),
		}
		fields := map[string]interface{}{
			keyStartTime:    metric.GetStartTime(),
			keyDuration:     metric.GetDuration(),
			keyCheckName:    metric.GetCheckName(),
			keyIsSuccessful: metric.IsSuccessful(),
			keyIsTimeout:    metric.IsTimeout(),
		}
		p := client.NewPoint("metrics", tags, fields, metric.GetStartTime())
		points = append(points, p)
	}
	return e.writeAPI.WritePoint(ctx, points...)
}

// getMetricsByChecksAndDuration fetches metrics from the metrics hypertable
// for the given check IDs. It accepts a `duration` parameter that fetches
// metrics for the check in the past `duration time.Duration`.
func (e *Exporter) getMetricsByChecksAndDuration(
	ctx context.Context,
	bucketName string,
	checkIDs []string,
	duration time.Duration,
) (map[string][]checker.Metric, error) {
	startTime := time.Now().Add(-1 * duration).Format(time.RFC3339Nano)

	query := fmt.Sprintf(`
from(bucket:%q)
	|> range(start: %s) 
	|> filter(fn: (r) => r._measurement == "metrics" and r.check_id =~ %s )
	|> sort(columns:["_time"], desc: true)
	|> pivot(rowKey:["_time"], columnKey:["_field"], valueColumn:"_value")
	|> drop(columns:["_time", "_start", "_stop"])
`, bucketName, startTime, regexMatchIDs(checkIDs))
	metrics := map[string][]checker.Metric{}
	result, err := e.queryAPI.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	for result.Next() {
		// NB: Duration and time are parsed separately because influx converts
		// time.Time to string hence manually parsing and assigning is done here.
		parsedDuration, err := time.ParseDuration(result.Record().ValueByKey(keyDuration).(string))
		if err != nil {
			return nil, err
		}
		parsedTime, err := time.Parse(time.RFC3339Nano, result.Record().ValueByKey(keyStartTime).(string))
		if err != nil {
			return nil, err
		}
		metric := config.Metric{
			CheckID:    result.Record().ValueByKey(keyCheckID).(string),
			CheckName:  result.Record().ValueByKey(keyCheckName).(string),
			StartTime:  parsedTime,
			Duration:   parsedDuration,
			Timeout:    result.Record().ValueByKey(keyIsTimeout).(bool),
			Successful: result.Record().ValueByKey(keyIsSuccessful).(bool),
		}

		if _, ok := metrics[metric.CheckID]; !ok {
			metrics[metric.CheckID] = make([]checker.Metric, 0)
		}
		metrics[metric.CheckID] = append(metrics[metric.CheckID], &metric)
	}

	return metrics, nil
}

// regexMatchIDs creates a regex which matches all the given check IDs.
func regexMatchIDs(checkIDs []string) string {
	fmtStr := "/^(%s)$/"
	sanitizedIDs := make([]string, 0, len(checkIDs))
	for i := range checkIDs {
		sanitizedIDs = append(sanitizedIDs, sanitizeID(checkIDs[i]))
	}
	return fmt.Sprintf(fmtStr, strings.Join(sanitizedIDs, "|"))
}

// sanitizeID escapes characters from the string so it can be safely used
// in the regex.
func sanitizeID(id string) string {
	res := ""
	for _, c := range id {
		switch c {
		case '/', '\\', '(', ')', '|':
			res += "\\"
		default:
		}
		res += string(c)
	}
	return res
}

// GetMetrics get the metrics of the given checks.
func (e *Exporter) GetMetrics(
	ctx context.Context,
	time time.Duration,
	checkIDs ...string,
) (map[string][]checker.Metric, error) {
	if len(checkIDs) == 0 {
		return nil, nil
	}

	return e.getMetricsByChecksAndDuration(ctx, e.dbname, checkIDs, time)
}

// Provision sets e's configuration.
func (e *Exporter) Provision(ctx *appcontext.Context, provider exporter.Provider) error {
	if provider.GetBackend() != exporterName {
		return fmt.Errorf(
			"invalid exporter name: expected '%s'; got '%s'",
			exporterName,
			provider.GetBackend(),
		)
	}

	cli, err := newClient(ctx, provider)
	if err != nil {
		return err
	}
	e.log = ctx.Logger()
	e.writeAPI = cli.WriteAPIBlocking(provider.GetOrgName(), provider.GetDBName())
	e.queryAPI = cli.QueryAPI(provider.GetOrgName())
	e.dbname = provider.GetDBName()

	return nil
}

// Interface guard.
var _ exporter.Exporter = (*Exporter)(nil)
