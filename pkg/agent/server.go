// Copyright (c) 2020 SDSLabs
// Use of this source code is governed by an MIT license
// details of which can be found in the LICENSE file.

package agent

import (
	"context"

	"github.com/sdslabs/pinger/pkg/config"
	"github.com/sdslabs/pinger/pkg/controller"
	"github.com/sdslabs/pinger/pkg/proto"
)

// server is the GRPC server that exposes the API so that central server
// can interact with the agent.
type server struct {
	m *controller.Manager
}

// ListChecks fetches a list of checks registered.
func (s *server) ListChecks(context.Context, *proto.Nil) (*proto.CheckList, error) {
	checksMap := s.m.ListControllers()

	checks := make([]*proto.CheckID, len(checksMap))

	i := 0
	for id := range checksMap {
		checks[i] = &proto.CheckID{ID: uint32(id)}
		i++
	}

	return &proto.CheckList{Checks: checks}, nil
}

// PushCheck creates a new check. If the check already exists it simply
// updates the check.
func (s *server) PushCheck(_ context.Context, check *proto.Check) (*proto.BoolResponse, error) {
	if err := addCheckToManager(s.m, config.ProtoToCheck(check)); err != nil {
		return &proto.BoolResponse{
			Successful: false,
			Error:      err.Error(),
		}, nil
	}

	return &proto.BoolResponse{Successful: true}, nil
}

// RemoveCheck removes the check.
func (s *server) RemoveCheck(_ context.Context, cid *proto.CheckID) (*proto.BoolResponse, error) {
	s.m.RemoveController(uint(cid.ID))
	return &proto.BoolResponse{Successful: true}, nil
}

// Interface guard.
var _ proto.AgentServer = (*server)(nil)
