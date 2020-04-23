// Package agent contains the agent that runs a GRPC server to interact with checks.
package agent

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/sdslabs/status/pkg/agent/proto"
	"github.com/sdslabs/status/pkg/check"
	"github.com/sdslabs/status/pkg/config"
	"github.com/sdslabs/status/pkg/controller"
	"github.com/sdslabs/status/pkg/database"
	"github.com/sdslabs/status/pkg/metrics"
)

// Agent runs a bunch of checks assigned to its manager. It updates the metrics
// depending upon the configuration of metrics provider. An agent can run in
// two modes - standalone and with the central API server.
type Agent struct {
	manager    *controller.Manager
	metrics    *metrics.ProviderConfig
	standalone bool
	port       int

	confChecks []*config.CheckConf // just the ones from the conf
}

// NewAgent creates a new agent with it's manager and config.
func NewAgent(conf *config.AgentConfig) (*Agent, error) {
	a := &Agent{
		manager:    controller.NewManager(),
		metrics:    &conf.Metrics,
		standalone: conf.Standalone,
		port:       conf.Port,

		confChecks: conf.Checks,
	}

	if err := a.setupMetrics(); err != nil {
		return nil, err
	}

	return a, nil
}

// Run updates the checks with the manager and starts the execution of checks.
func (a *Agent) Run() {
	// Register the already added checks with the agent.
	for _, c := range a.confChecks {
		if err := a.registerCheck(c); err != nil {
			logrus.WithField(
				"check", c.GetName()).WithError(err).Errorln("error creating check")
		}
	}

	if a.standalone {
		// Just keep waiting for termination in standalone mode.
		a.manager.Wait()
	} else {
		// Start the GRPC server if running with a central api server.
		a.runGRPCServer()
	}
}

// PushCheck pushes the check in the agent.
func (a *Agent) PushCheck(ctx context.Context, agentCheck *proto.Check) (*proto.PushStatus, error) {
	cfg := config.GetCheckFromCheckProto(agentCheck)

	if err := a.registerCheck(cfg); err != nil {
		return &proto.PushStatus{
			Pushed: false,
			Reason: err.Error(),
		}, err
	}

	return &proto.PushStatus{Pushed: true}, nil
}

// GetManagerStats fetches stats about the checks from the manager.
func (a *Agent) GetManagerStats(context.Context, *proto.None) (*proto.ManagerStats, error) {
	mStats := []*proto.ManagerStats_ControllerStatus{}

	for _, stat := range a.fetchStats() {
		mStats = append(mStats, &proto.ManagerStats_ControllerStatus{
			Name: stat.Name,
			ConfigStatus: &proto.ManagerStats_ControllerConfigurationStatus{
				ErrorRetry:    stat.Configuration.ErrorRetry,
				ShouldBackOff: stat.Configuration.ShouldBackOff,
				Interval:      stat.Configuration.Interval,
			},
			RunStatus: &proto.ManagerStats_ControllerRunStatus{
				SuccessCount:       stat.Status.SuccessCount,
				FailureCount:       stat.Status.FailureCount,
				ConsecFailureCount: stat.Status.ConsecutiveFailureCount,
				LastSuccessTime:    stat.Status.LastSuccessStamp,
				LastFailureTime:    stat.Status.LastFailureStamp,
			},
		})
	}

	return &proto.ManagerStats{ControllerStatus: mStats}, nil
}

// RemoveCheck removes the check from the agent.
func (a *Agent) RemoveCheck(ctx context.Context, agentCheck *proto.CheckMeta) (*proto.RemoveStatus, error) {
	err := a.removeCheck(agentCheck.Name)
	if err != nil {
		return &proto.RemoveStatus{
			Removed: false,
			Message: err.Error(),
		}, err
	}

	return &proto.RemoveStatus{Removed: true}, nil
}

// ListChecks returns a list of checks with the agent.
func (a *Agent) ListChecks(context.Context, *proto.None) (*proto.ChecksList, error) {
	checksList := []*proto.CheckMeta{}

	for _, ctrl := range a.fetchChecks() {
		checksList = append(checksList, &proto.CheckMeta{Name: ctrl})
	}

	return &proto.ChecksList{Checks: checksList}, nil
}

// runGRPCServer starts a GRPC server at the specified port. This also initializes
// the controller manager instance, which is used further to interact with the controllers.
func (a *Agent) runGRPCServer() error {
	listner, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", a.port))
	if err != nil {
		return fmt.Errorf("error starting listener: %v", err)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterAgentServiceServer(grpcServer, a)

	if err = grpcServer.Serve(listner); err != nil {
		return fmt.Errorf("error starting server: %v", err)
	}

	return nil
}

// setupMetrics paves the way for an agent to update metrics in the configured DB.
func (a *Agent) setupMetrics() error {
	if !a.standalone {
		// For non-standalone mode only use the main database,
		// i.e., the timescale-postgres database.
		database.SetupMetrics(a.metrics, a.manager)
		return nil
	}

	// For standalone mode, metrics providers have options.
	switch a.metrics.Backend {
	case metrics.PrometheusProviderType:
		metrics.SetupPrometheusMetrics(a.metrics, a.manager)

	case metrics.TimeScaleProviderType:
		database.SetupMetrics(a.metrics, a.manager)

	case metrics.LogProviderType:
		metrics.SetupLogMetrics(a.metrics, a.manager)

	default:
		return fmt.Errorf("invalid metrics provider: %v", a.metrics.Backend)
	}

	return nil
}

// register check updates the manager with the check.
func (a *Agent) registerCheck(c *config.CheckConf) error {
	// Only add checks to the database if the backend is timescale and
	// the agent runs in standalone mode. The standalone mode is important
	// because for checks that come through the api-server will already be
	// inserted into the database.
	if a.standalone && a.metrics.Backend == metrics.TimeScaleProviderType {
		err := database.AddCheckToDB(c)
		if err != nil {
			return err
		}
	}

	checker, err := check.NewChecker(c)
	if err != nil {
		return err
	}

	cFunc, err := controller.NewControllerFunction(checker.ExecuteCheck)
	if err != nil {
		return err
	}

	executor := controller.Internal{
		DoFunc:      cFunc,
		RunInterval: time.Duration(c.GetInterval()),
	}

	controllerName := fmt.Sprint(c.GetId())
	if a.metrics.Backend != metrics.TimeScaleProviderType {
		controllerName = c.GetLabel()
	}

	err = a.manager.UpdateController(controllerName, checker.Type(), executor)
	if err != nil {
		return err
	}

	return nil
}

// fetchChecks fetches checks from the agent manager.
func (a *Agent) fetchChecks() []string {
	return a.manager.GetAllControllers()
}

// removeCheck removes a check from the agent manager.
func (a *Agent) removeCheck(id string) error {
	return a.manager.RemoveControllerAndWait(id)
}

// fetchStats fetches stats from the agent manager.
func (a *Agent) fetchStats() []*controller.Status {
	return a.manager.GetStats()
}

// Interface guard
var _ proto.AgentServiceServer = (*Agent)(nil)
