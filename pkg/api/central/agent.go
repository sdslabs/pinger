package central

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"

	"github.com/sdslabs/status/pkg/agent/proto"

	log "github.com/sirupsen/logrus"
)

// Agent is the type for agent that runs the GRPC server.
type Agent struct {
	Host string
	Port int

	Timeout time.Duration
}

var grpcOpts = []grpc.DialOption{grpc.WithInsecure()}

// PushCheckToAgent takes a agentProto.Check and pushes it to the agent.
func (a *Agent) PushCheckToAgent(check *proto.Check) error {
	agentAddr := getAddr(a.Host, a.Port)

	conn, err := grpc.Dial(agentAddr, grpcOpts...)
	if err != nil {
		return fmt.Errorf("cannot dial rpc at %s: %v", agentAddr, err)
	}
	defer conn.Close() //nolint:errcheck

	client := proto.NewAgentServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), a.Timeout)
	defer cancel()

	_, err = client.PushCheck(ctx, check)
	log.Error(err)

	return err
}

// RemoveCheckFromAgent takes a agentProto.CheckMeta and removes it from the agent.
func (a *Agent) RemoveCheckFromAgent(checkMeta *proto.CheckMeta) error {
	agentAddr := getAddr(a.Host, a.Port)

	conn, err := grpc.Dial(agentAddr, grpcOpts...)
	if err != nil {
		return fmt.Errorf("cannot dial rpc at %s: %v", agentAddr, err)
	}
	defer conn.Close() //nolint:errcheck

	client := proto.NewAgentServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), a.Timeout)
	defer cancel()

	_, err = client.RemoveCheck(ctx, checkMeta)
	log.Error(err)

	return err
}
