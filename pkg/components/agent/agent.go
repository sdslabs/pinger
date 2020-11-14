package agent

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/sdslabs/kiwi"
	"github.com/sdslabs/kiwi/values/hash"

	"google.golang.org/grpc"

	"github.com/sdslabs/kiwi/stdkiwi"

	"github.com/sdslabs/pinger/pkg/alerter"
	"github.com/sdslabs/pinger/pkg/checker"
	"github.com/sdslabs/pinger/pkg/components/agent/proto"
	"github.com/sdslabs/pinger/pkg/config"
	"github.com/sdslabs/pinger/pkg/config/configfile"
	"github.com/sdslabs/pinger/pkg/exporter"
	"github.com/sdslabs/pinger/pkg/util/appcontext"
	"github.com/sdslabs/pinger/pkg/util/controller"
)

type alertMap struct {
	a  map[string]map[string]alerter.Alert
	mu sync.RWMutex
}

// Run starts the agent.
//
// It either starts the agent in standalone mode where the manager waits for
// it's execution to end or it starts the GRPC API server for the central
// server to interact with.
func Run(ctx *appcontext.Context, conf *configfile.Agent) error {
	if conf.Interval <= 0 {
		return fmt.Errorf("interval should be > 0")
	}

	manager := controller.NewManager(ctx)

	export, err := exporter.Initialize(ctx, &conf.Metrics)
	if err != nil {
		return fmt.Errorf("cannot initialize exporter: %w", err)
	}

	// alertPrevState stores the last state of the check for a particular CheckID in a has map.
	const alertPrevStateKey = "alertPrevState"
	store, err := stdkiwi.NewStoreFromSchema(kiwi.Schema{
		alertPrevStateKey: hash.Type,
	})
	if err != nil {
		return fmt.Errorf("cannot initialize store: %w", err)
	}
	alertPrevState := store.Hash(alertPrevStateKey)

	aMap := alertMap{a: map[string]map[string]alerter.Alert{}}
	alertFuncs := map[string]alerter.AlertFunc{}
	for i := range conf.Alerts {
		ap := conf.Alerts[i]

		if _, ok := alertFuncs[ap.Service]; ok {
			return fmt.Errorf("alerter %q already configured", ap.Service)
		}

		alert, er := alerter.Initialize(ctx, &ap)
		if er != nil {
			return fmt.Errorf("cannot initialize alerter: %w", er)
		}

		alertFuncs[ap.Service] = alert
		aMap.a[ap.Service] = map[string]alerter.Alert{}
	}

	err = initExportAndAlerts(ctx, conf.Interval, manager, export, alertFuncs, &aMap, alertPrevState)
	if err != nil {
		return fmt.Errorf("cannot initialize exporter: %w", err)
	}

	// These are the checks provided through config. This essentially implies
	// that the checks will be run always irrespective of the fact that agent
	// running in standalone mode or not.
	for i := range conf.Checks {
		if err := addCheckToManager(manager, &aMap, &conf.Checks[i]); err != nil {
			return fmt.Errorf("check %d: cannot add to manager: %w", i, err)
		}
	}

	if conf.Standalone {
		// for standalone mode we just need to wait for the controller to end.
		manager.Wait()
		return nil
	}

	return runGRPCServer(manager, &aMap, conf.Port)
}

// initExportAndAlerts initializes the controller for exporting and alerting
// the metrics.
func initExportAndAlerts(
	ctx *appcontext.Context,
	interval time.Duration,
	manager *controller.Manager,
	exportFunc exporter.ExportFunc,
	alertFuncs map[string]alerter.AlertFunc,
	aMap *alertMap,
	alertPrevState *stdkiwi.Hash,
) error {
	ctrl, err := controller.NewController(ctx, &controller.Opts{
		Name:     "metrics_export_and_alert",
		Interval: interval,
		Func: func(c context.Context) (interface{}, error) {
			stats := manager.PullAllStats()
			var (
				exportMetrics []checker.Metric
				alertMetrics  []checker.Metric
			)
			for _, stat := range stats {
				// stat for each check with latest timestamp will be alerted.
				lastTimestamp := time.Time{}
				for _, s := range stat {
					if s == nil {
						continue
					}

					var metric config.Metric
					if s.Err != nil {
						// on error record the failed metric.
						metric = config.Metric{
							CheckID:   s.ID,
							CheckName: s.Name,
						}
					} else {
						res, ok := s.Res.(*checker.Result)
						if !ok {
							ctx.Logger().
								WithField("check_id", s.ID).
								Warnln("unexpected error: check result not checker.Result")
							continue
						}

						metric = config.Metric{
							CheckID:    s.ID,
							CheckName:  s.Name,
							Successful: res.Successful,
							Timeout:    res.Timeout,
							StartTime:  res.StartTime,
							Duration:   res.Duration,
						}
					}

					exportMetrics = append(exportMetrics, &metric)

					shouldAlert, err := shouldUpdateAlert(alertPrevState, &lastTimestamp, &metric)
					if err != nil {
						ctx.Logger().
							WithField("check_id", s.ID).WithError(err).
							Warnln("could not alert")
						continue
					}

					if shouldAlert {
						alertMetrics = append(alertMetrics, &metric)
					}
				}
			}

			// Export metrics into the database
			if er := exportFunc(ctx, exportMetrics); er != nil {
				ctx.Logger().
					WithError(er).
					Errorln("error exporting metrics")
				return nil, er
			}

			// Alert metrics from the corresponding services
			for alertService, alertFunc := range alertFuncs {
				aMap.mu.RLock()
				serviceAlertMap, ok := (aMap.a)[alertService]
				aMap.mu.RUnlock()
				if !ok {
					er := fmt.Errorf("could not find alerts for service %q", alertService)
					ctx.Logger().WithError(er).Errorf("could not alert")
					return nil, er
				}

				if er := alertFunc(ctx, alertMetrics, serviceAlertMap); er != nil {
					ctx.Logger().
						WithError(er).
						Errorln("error alerting exportMetrics")
					return nil, er
				}
			}

			return nil, nil
		},
	})
	if err != nil {
		return err
	}

	ctrl.Start()
	return nil
}

// shouldUpdateAlert tells if an alert should be sent for the particular metric.
func shouldUpdateAlert(
	alertPrevState *stdkiwi.Hash,
	lastTimestamp *time.Time,
	metric checker.Metric,
) (update bool, _ error) {
	if lastTimestamp.After(metric.GetStartTime()) {
		// If the last time stamp is after the metric, we have a stale metric
		// and do not want to consider it for alerting.
		return false, nil
	}

	prevState := boolNil // if nothing is fetched, first alert will be sent
	newState := newBoolean(metric.IsSuccessful())

	has, err := alertPrevState.Has(metric.GetCheckID())
	if err != nil {
		return false, err
	}
	if has {
		v, er := alertPrevState.Get(metric.GetCheckID())
		if er != nil {
			return false, er
		}
		prevState = boolean(v[0])
	}

	if newState == prevState {
		return false, nil
	}

	err = alertPrevState.Insert(metric.GetCheckID(), string(newState))
	if err != nil {
		return false, err
	}

	// update the last time stamp only when there is no error or else let the
	// last metric to be the same as before.
	*lastTimestamp = metric.GetStartTime()
	return true, nil
}

// runGRPCServer starts the GRPC server that exposes an API for the central
// to contact the agent.
func runGRPCServer(manager *controller.Manager, aMap *alertMap, port uint16) error {
	addr := net.JoinHostPort("0.0.0.0", fmt.Sprint(port))

	lst, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("unable to start listener: %v", err)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterAgentServer(grpcServer, &server{
		m: manager,
		a: aMap,
	})

	err = grpcServer.Serve(lst)
	if err != nil {
		return fmt.Errorf("unable to start serer: %v", err)
	}

	return nil
}

// addCheckToManager adds a new check to the manager.
func addCheckToManager(
	manager *controller.Manager,
	aMap *alertMap,
	check *config.Check,
) error {
	ctrlOpts, err := checker.NewControllerOpts(check)
	if err != nil {
		return err
	}

	for i := range check.Alerts {
		alt := check.Alerts[i]

		aMap.mu.RLock()
		_, ok := (aMap.a)[alt.Service]
		aMap.mu.RUnlock()
		if !ok {
			return fmt.Errorf("invalid alerter %q", alt.Service)
		}

		aMap.mu.Lock()
		(aMap.a)[alt.Service][check.ID] = &alt
		aMap.mu.Unlock()
	}

	return manager.UpdateController(ctrlOpts)
}
