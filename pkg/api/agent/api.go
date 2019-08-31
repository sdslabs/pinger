package agent

import (
	context "context"
	"fmt"
	"net"
	"time"

	"google.golang.org/grpc"

	"github.com/sdslabs/status/pkg/api/agent/proto"
	"github.com/sdslabs/status/pkg/check"
	"github.com/sdslabs/status/pkg/controller"

	log "github.com/sirupsen/logrus"
)

const (
	// AGENT_SERVER_PORT: Port for running the GRPC server
	AGENT_SERVER_PORT = 9009

	// AGENT_GRPC_HOST: Host to run the GRPC server on.
	AGENT_GRPC_HOST = "0.0.0.0"
)

// ControllerManager is the global manager for the controller that comes with the
// agent. It is initialized when we run the GRPC servers.
var ControllerManager controller.Manager

type agentServer struct{}

func (a agentServer) PushCheck(ctx context.Context, agentCheck *proto.Check) (*proto.PushStatus, error) {
	log.Debug("Recieved the push for a new check.")
	checker, err := check.NewChecker(agentCheck)
	if err != nil {
		log.Errorf("Error while creating new checker: %s", err)
		return &proto.PushStatus{
			Pushed: false,
			Reason: err.Error(),
		}, err
	}

	cFunc, err := controller.NewControllerFunction(checker.ExecuteCheck)
	if err != nil {
		sErr := fmt.Errorf("Error while creating controller function: %s", err)
		log.Error(sErr)
		return &proto.PushStatus{
			Pushed: false,
			Reason: sErr.Error(),
		}, sErr
	}

	executor := controller.ControllerInternal{
		DoFunc:      cFunc,
		RunInterval: time.Second * time.Duration(agentCheck.Interval),
	}
	err = ControllerManager.UpdateController(agentCheck.Name, executor)
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
	return nil, nil
}

// RunGrpcServer starts a GRPC server at the specified port.
// This also initializes the controller manager instance, which is used further
// to interact with the controllers.
func RunGrpcServer(notify chan struct{}) {
	listner, err := net.Listen("tcp", fmt.Sprintf("%s:%d", AGENT_GRPC_HOST, AGENT_SERVER_PORT))
	if err != nil {
		log.Errorf("Error while starting listner : %s", err)
		notify <- struct{}{}
		return
	}

	grpcServer := grpc.NewServer()
	server := agentServer{}

	proto.RegisterAgentServiceServer(grpcServer, server)

	ControllerManager = controller.NewManager()

	log.Infof("Starting new server at port : %d", AGENT_SERVER_PORT)
	grpcServer.Serve(listner)
}
