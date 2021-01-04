package util

import (
	"time"

	"github.com/sdslabs/pinger/pkg/checker"
	"github.com/sdslabs/pinger/pkg/util/httpserver"
)

// SerializeMetrics breaks the metrics into multiple batches and retains one
// metric from each batch.
//
// The following rules are applied to each batch:
// 	- Failed metric is prioritized over successful.
// 	- Metric with highest latency is considered.
//  - If number of batches are more than number of metrics, this is probably
// 	  recent addition of check. In this case, The first metric should be
// 	  appended at the front of list.
//  - The first (latest) metric is remains the same.
//
// Minimum number of batches accepted is 2 since only 1 would mean just
// getting the latest metric which doesn't reflect any history.
func SerializeMetrics(
	batches int, metrics []checker.Metric,
) (serialized []httpserver.MetricResponse, uptime int) {
	if batches < 2 || len(metrics) == 0 {
		return
	}

	batches-- // since the first batch will essentially be the first metric

	serialized = make([]httpserver.MetricResponse, 0, batches)
	numEachBatch := (len(metrics) / batches) + 1
	var upNum int

	if metrics[0].IsSuccessful() {
		upNum++
	}
	serialized = append(serialized, httpserver.MetricResponse{
		Successful: metrics[0].IsSuccessful(),
		Timeout:    metrics[0].IsTimeout(),
		StartTime:  metrics[0].GetStartTime(),
		Duration:   metrics[0].GetDuration(),
	})

	for i := 1; i < len(metrics); i += numEachBatch {
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

// PrepareMetricsResponse creates a page metrics response for each of the
// check after serializing the metrics.
func PrepareMetricsResponse(
	batches int, metrics map[string][]checker.Metric,
) httpserver.PageMetricsResponse {
	resp := map[string]httpserver.PageCheckMetricsResponse{}
	var checksDown int
	for cid := range metrics {
		// NB: This shouldn't take that long but since this request is long enough
		// in general, one optimization can be to serialize the metrics in different
		// goroutines. For now this works and we need benchmarks to prove if making
		// this change would really help or not.
		serialized, uptime := SerializeMetrics(batches, metrics[cid])
		if len(serialized) == 0 || len(metrics[cid]) == 0 {
			continue
		}
		operational := serialized[0].Successful
		resp[cid] = httpserver.PageCheckMetricsResponse{
			Metrics:     serialized,
			Uptime:      uptime,
			Operational: operational,
		}
		if !operational {
			checksDown++
		}
	}
	return httpserver.PageMetricsResponse{
		ChecksDown: checksDown,
		Checks:     resp,
	}
}
