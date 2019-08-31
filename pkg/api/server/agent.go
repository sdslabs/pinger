package server

import (
	"context"
	"fmt"
	"time"

	"github.com/sdslabs/status/pkg/api/agent/proto"
	"github.com/sdslabs/status/pkg/defaults"
	"google.golang.org/grpc"

	log "github.com/sirupsen/logrus"
)

var grpcOpts = []grpc.DialOption{grpc.WithInsecure()}

type StatusAgent struct {
	Host string
	Port int64

	Timeout time.Duration
}

func NewStatusAgent(host string, port int64) *StatusAgent {
	return &StatusAgent{
		Host: host,
		Port: port,

		Timeout: defaults.DefaultGRPCRequestTimeout,
	}
}

func (a *StatusAgent) PushCheckToAgent(check *proto.Check) error {
	agentAddr := fmt.Sprintf("%s:%d", a.Host, a.Port)
	log.Debugf("Pushing check to the agent: %s", agentAddr)

	conn, err := grpc.Dial(agentAddr, grpcOpts...)
	if err != nil {
		return fmt.Errorf("ERROR while dailing RPC for agent at %s : %s", agentAddr, err)
	}
	defer conn.Close()

	client := proto.NewAgentServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), a.Timeout)
	defer cancel()

	_, err = client.PushCheck(ctx, check)
	log.Error(err)

	return err
}
