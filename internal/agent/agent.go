// Copyright (c) 2020 SDSLabs
// Use of this source code is governed by an MIT license
// details of which can be found in the LICENSE file.

package agent

import (
	"fmt"
	"net"

	"google.golang.org/grpc"

	"github.com/sdslabs/status/internal/appcontext"
	"github.com/sdslabs/status/internal/checker"
	"github.com/sdslabs/status/internal/config"
	"github.com/sdslabs/status/internal/config/configfile"
	"github.com/sdslabs/status/internal/controller"
	"github.com/sdslabs/status/internal/metrics"
	"github.com/sdslabs/status/pkg/proto"
)

// Run starts the agent.
//
// It either starts the agent in standalone mode where the manager waits for
// it's execution to end or it starts the GRPC API server for the central
// server to interact with.
func Run(ctx *appcontext.Context, conf *configfile.Agent) error {
	manager := controller.NewManager(ctx)
	checks := config.CheckListToInterface(conf.Checks)

	if err := metrics.Initialize(ctx, manager, &conf.Metrics, checks); err != nil {
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
