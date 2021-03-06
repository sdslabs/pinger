package agent

import (
	"context"

	"github.com/sdslabs/pinger/pkg/components/agent/proto"
	"github.com/sdslabs/pinger/pkg/config"
	"github.com/sdslabs/pinger/pkg/util/controller"
)

// server is the GRPC server that exposes the API so that central server
// can interact with the agent.
type server struct {
	m *controller.Manager
	a *alertMap
	// Unimplemented agent server for "forward compatibility".
	proto.UnimplementedAgentServer
}

// ListChecks fetches a list of checks registered.
func (s *server) ListChecks(context.Context, *proto.Nil) (*proto.CheckList, error) {
	checksMap := s.m.ListControllers()

	checks := make([]*proto.CheckID, len(checksMap))

	i := 0
	for id := range checksMap {
		checks[i] = &proto.CheckID{ID: id}
		i++
	}

	return &proto.CheckList{Checks: checks}, nil
}

// PushCheck creates a new check. If the check already exists it simply
// updates the check.
func (s *server) PushCheck(_ context.Context, check *proto.Check) (*proto.BoolResponse, error) {
	c := config.ProtoToCheck(check)
	if err := addCheckToManager(s.m, s.a, &c); err != nil {
		return &proto.BoolResponse{
			Successful: false,
			Error:      err.Error(),
		}, nil
	}

	return &proto.BoolResponse{Successful: true}, nil
}

// RemoveCheck removes the check.
func (s *server) RemoveCheck(_ context.Context, cid *proto.CheckID) (*proto.BoolResponse, error) {
	s.m.RemoveController(cid.ID)
	return &proto.BoolResponse{Successful: true}, nil
}

// Interface guard.
var _ proto.AgentServer = (*server)(nil)
