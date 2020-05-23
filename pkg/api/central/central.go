package central

import (
	"context"
	"fmt"
	"net"
	"time"

	"google.golang.org/grpc"

	"github.com/sdslabs/status/pkg/api/central/proto"
	"github.com/sdslabs/status/pkg/config"
	"github.com/sdslabs/status/pkg/defaults"

	log "github.com/sirupsen/logrus"
)

// Server runs the central-server that listens for requests from the app server
// and pushes the checks to the agents.
type Server struct {
	agents map[string]*Agent
	port   int
}

// newServer returns an instance of central server.
func newServer(conf *config.CentralServerConfig) *Server {
	s := &Server{
		port:   conf.Port,
		agents: map[string]*Agent{},
	}

	for _, agent := range conf.Agents {
		if err := s.addAgent(agent.Host, agent.Port, agent.Timeout); err != nil {
			log.WithError(err).Errorln("cannot add agent")
			continue
		}
	}

	return s
}

// run starts the central server.
func (s *Server) run() error {
	listner, err := net.Listen("tcp", getAddr("0.0.0.0", s.port))
	if err != nil {
		return fmt.Errorf("error starting listener: %v", err)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterCentralServiceServer(grpcServer, s)

	if err = grpcServer.Serve(listner); err != nil {
		return fmt.Errorf("error starting server: %v", err)
	}

	return nil
}

// PushCheck pushes a check to the central service to be further assigned
// to a registered agent.
func (s *Server) PushCheck(context.Context, *proto.Check) (*proto.PushStatus, error) {
	return &proto.PushStatus{}, nil
}

// RemoveCheck removes a check from the manager managing the controller for the checks.
// Only reuired filed for the CheckMeta is `ID` rest everything can be ignored.
func (s *Server) RemoveCheck(context.Context, *proto.CheckMeta) (*proto.RemoveStatus, error) {
	return &proto.RemoveStatus{}, nil
}

// addAgent adds a new agent to the central server.
func (s *Server) addAgent(host string, port int, timeout time.Duration) error {
	addr := getAddr(host, port)
	if _, ok := s.agents[addr]; ok {
		return fmt.Errorf("already registered agent: %s", addr)
	}

	tmout := defaults.GRPCRequestTimeout
	if timeout > 0 {
		tmout = timeout
	}

	s.agents[addr] = &Agent{
		Host:    host,
		Port:    port,
		Timeout: tmout,
	}

	return nil
}

// removeAgent removes an agent from central server.
func (s *Server) removeAgent(host string, port int) error {
	addr := getAddr(host, port)

	if _, ok := s.agents[addr]; !ok {
		return fmt.Errorf("agent not registered: %s", addr)
	}

	delete(s.agents, addr)
	return nil
}

// Serve starts the central server with the provided conf.
func Serve(conf *config.CentralServerConfig) error {
	server := newServer(conf)
	return server.run()
}
