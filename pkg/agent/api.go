package agent

import (
	context "context"
	"fmt"
	"net"
	"time"

	"google.golang.org/grpc"

	"github.com/sdslabs/status/pkg/agent/proto"
	"github.com/sdslabs/status/pkg/check"
	"github.com/sdslabs/status/pkg/config"
	"github.com/sdslabs/status/pkg/controller"

	log "github.com/sirupsen/logrus"
)

const (
	// AgentGRPCHost is the host to run the GRPC server on.
	AgentGRPCHost = "0.0.0.0"
)

type agentServer struct{}

func (a agentServer) PushCheck(ctx context.Context, agentCheck *proto.Check) (*proto.PushStatus, error) {
	log.Debug("Received the push for a new check.")
	cfg := config.GetCheckFromCheckProto(agentCheck)
	checker, err := check.NewChecker(cfg)
	if err != nil {
		log.Errorf("Error while creating new checker: %s", err)
		return &proto.PushStatus{
			Pushed: false,
			Reason: err.Error(),
		}, err
	}

	cFunc, err := controller.NewControllerFunction(checker.ExecuteCheck)
	if err != nil {
		sErr := fmt.Errorf("error while creating controller function: %s", err)
		log.Error(sErr)
		return &proto.PushStatus{
			Pushed: false,
			Reason: sErr.Error(),
		}, sErr
	}

	executor := controller.Internal{
		DoFunc:      cFunc,
		RunInterval: time.Second * time.Duration(agentCheck.Interval),
	}
	err = ControllerManager.UpdateController(agentCheck.Name, checker.Type(), executor)
	if err != nil {
		log.Errorf("Error while creating controller: %s", err)
	}

	return &proto.PushStatus{
		Pushed: true,
		Reason: "Push Successful",
	}, nil
}

func (a agentServer) GetManagerStats(context.Context, *proto.None) (*proto.ManagerStats, error) {
	log.Debug("Inside get controller manager stats function")

	stats := ControllerManager.GetStats()

	mStats := []*proto.ManagerStats_ControllerStatus{}
	for _, stat := range stats {
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

func (a agentServer) RemoveCheck(ctx context.Context, agentCheck *proto.CheckMeta) (*proto.RemoveStatus, error) {
	log.Debugf("Removing check from check controller manager: %s", agentCheck.Name)

	err := ControllerManager.RemoveControllerAndWait(agentCheck.Name)
	if err != nil {
		return &proto.RemoveStatus{
			Removed: false,
			Message: fmt.Sprintf("Error while removing: %s", err),
		}, err
	}

	return nil, nil
}

func (a agentServer) ListChecks(context.Context, *proto.None) (*proto.ChecksList, error) {
	log.Debugf("Listing checks managed by agent's default registered manager")

	ctrls := ControllerManager.GetAllControllers()
	checksList := []*proto.CheckMeta{}

	for _, ctrl := range ctrls {
		checksList = append(checksList, &proto.CheckMeta{Name: ctrl})
	}

	return &proto.ChecksList{Checks: checksList}, nil
}

// RunGRPCServer starts a GRPC server at the specified port.
// This also initializes the controller manager instance, which is used further
// to interact with the controllers.
func RunGRPCServer(port int) {
	listner, err := net.Listen("tcp", fmt.Sprintf("%s:%d", AgentGRPCHost, port))
	if err != nil {
		log.Errorf("Error while starting listner : %s", err)
		return
	}

	grpcServer := grpc.NewServer()
	server := agentServer{}

	proto.RegisterAgentServiceServer(grpcServer, server)

	ControllerManager = controller.NewManager()

	log.Infof("Starting new server at port : %d", port)

	if err = grpcServer.Serve(listner); err != nil {
		log.Fatalf("Error starting the grpc server: %s", err.Error())
	}
}
