package agent

import (
	context "context"
	"fmt"
	"net"

	"google.golang.org/grpc"

	log "github.com/sirupsen/logrus"
)

const (
	AGENT_SERVER_PORT = 9009
	AGENT_GRPC_HOST   = "0.0.0.0"
)

type agentServer struct{}

func (a agentServer) PushCheck(ctx context.Context, check *Check) (*PushStatus, error) {
	return nil, nil
}

func (a agentServer) GetManagerStats(context.Context, *None) (*ManagerStats, error) {
	return nil, nil
}

// RunGrpcServer starts a GRPC server at the specified port.
func RunGrpcServer(notify chan struct{}) {
	listner, err := net.Listen("tcp", fmt.Sprintf("%s:%d", AGENT_GRPC_HOST, AGENT_SERVER_PORT))
	if err != nil {
		log.Errorf("Error while starting listner : %s", err)
		notify <- struct{}{}
		return
	}

	grpcServer := grpc.NewServer()
	server := agentServer{}

	RegisterAgentServiceServer(grpcServer, server)

	log.Infof("Starting new server at port : %d", AGENT_SERVER_PORT)
	grpcServer.Serve(listner)
}
