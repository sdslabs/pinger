package server

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"

	"github.com/sdslabs/status/pkg/agent/proto"
	"github.com/sdslabs/status/pkg/defaults"

	log "github.com/sirupsen/logrus"
)

var grpcOpts = []grpc.DialOption{grpc.WithInsecure()}

// StatusAgent is the type for agent that runs the GRPC server.
type StatusAgent struct {
	Host string
	Port int64

	Timeout time.Duration
}

// NewStatusAgent creates a new status agent from host and port.
func NewStatusAgent(host string, port int64) *StatusAgent {
	return &StatusAgent{
		Host: host,
		Port: port,

		Timeout: defaults.DefaultGRPCRequestTimeout,
	}
}

// PushCheckToAgent takes a proto.Check and pushes it to the agent.
func (a *StatusAgent) PushCheckToAgent(check *proto.Check) error {
	agentAddr := fmt.Sprintf("%s:%d", a.Host, a.Port)
	log.Debugf("Pushing check to the agent: %s", agentAddr)

	conn, err := grpc.Dial(agentAddr, grpcOpts...)
	if err != nil {
		return fmt.Errorf("ERROR while dailing RPC for agent at %s : %s", agentAddr, err)
	}
	defer conn.Close() //nolint:errcheck

	client := proto.NewAgentServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), a.Timeout)
	defer cancel()

	_, err = client.PushCheck(ctx, check)
	log.Error(err)

	return err
}
