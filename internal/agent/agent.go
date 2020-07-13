// Copyright (c) 2020 SDSLabs
// Use of this source code is governed by an MIT license
// details of which can be found in the LICENSE file.

package agent

import (
	"context"
	"fmt"
	"net"
	"time"

	"google.golang.org/grpc"

	"github.com/sdslabs/status/internal/appcontext"
	"github.com/sdslabs/status/internal/checker"
	"github.com/sdslabs/status/internal/config"
	"github.com/sdslabs/status/internal/config/configfile"
	"github.com/sdslabs/status/internal/controller"
	"github.com/sdslabs/status/internal/exporter"
	"github.com/sdslabs/status/pkg/proto"
)

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
	checks := config.CheckListToInterface(conf.Checks)

	// add unique ID for config-only checks now and then let the exporter handle
	// the change in IDs
	for i := range checks {
		checks[i].SetID(uint(i))
	}

	export, err := exporter.Initialize(ctx, &conf.Metrics, checks)
	if err != nil {
		return fmt.Errorf("cannot initialize exporter: %w", err)
	}

	err = initExportAndAlerts(ctx, conf.Interval, manager, export)
	if err != nil {
		return fmt.Errorf("cannot initialize exporter: %w", err)
	}

	// These are the checks provided through config. This essentially implies
	// that the checks will be run always irrespective of the fact that agent
	// running in standalone mode or not.
	for i := range conf.Checks {
		if err := addCheckToManager(manager, &conf.Checks[i]); err != nil {
			return fmt.Errorf("check %d: cannot add to manager: %w", i, err)
		}
	}

	if conf.Standalone {
		// for standalone mode we just need to wait for the controller to end.
		manager.Wait()
		return nil
	}

	return runGRPCServer(manager, conf.Port)
}

// doFunc is the function that does the exporting and alerting.
type doFunc = func(context.Context, []checker.Metric) error

// initExportAndAlerts initializes the controller for exporting and alerting
// the metrics.
func initExportAndAlerts(
	ctx *appcontext.Context,
	interval time.Duration,
	manager *controller.Manager,
	export /*, alert*/ doFunc,
) error {
	ctrl, err := controller.NewController(ctx, &controller.Opts{
		Name:     "metrics_export_and_alert",
		Interval: interval,
		Func: func(c context.Context) (interface{}, error) {
			stats := manager.PullAllStats()
			metrics := []checker.Metric{}
			for _, stat := range stats {
				for _, s := range stat {
					if s == nil {
						continue
					}

					if s.Err != nil {
						// on error record the failed metric
						metric := &config.Metric{
							CheckID:   s.ID,
							CheckName: s.Name,
						}
						metrics = append(metrics, metric)
						continue
					}

					res, ok := s.Res.(*checker.Result)
					if !ok {
						ctx.Logger().
							WithField("check_id", s.ID).
							Warnln("unexpected error: check result not checker.Result")
						continue
					}

					metric := &config.Metric{
						CheckID:    s.ID,
						CheckName:  s.Name,
						Successful: res.Successful,
						Timeout:    res.Timeout,
						StartTime:  res.StartTime,
						Duration:   res.Duration,
					}
					metrics = append(metrics, metric)
				}
			}

			// Export metrics into the database
			if er := export(ctx, metrics); er != nil {
				ctx.Logger().
					WithError(er).
					Errorln("error exporting metrics")
				return nil, er
			}

			// TODO(shreyaa-sharmaa): Alert metrics (in a separate thread than this)

			return nil, nil
		},
	})
	if err != nil {
		return err
	}

	ctrl.Start()
	return nil
}

// runGRPCServer starts the GRPC server that exposes an API for the central
// to contact the agent.
func runGRPCServer(manager *controller.Manager, port uint16) error {
	addr := net.JoinHostPort("0.0.0.0", fmt.Sprint(port))

	lst, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("unable to start listener: %v", err)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterAgentServer(grpcServer, &server{m: manager})

	err = grpcServer.Serve(lst)
	if err != nil {
		return fmt.Errorf("unable to start serer: %v", err)
	}

	return nil
}

// addCheckToManager adds a new check to the manager.
func addCheckToManager(manager *controller.Manager, check checker.Check) error {
	ctrlOpts, err := checker.NewControllerOpts(check)
	if err != nil {
		return err
	}

	return manager.UpdateController(ctrlOpts)
}
