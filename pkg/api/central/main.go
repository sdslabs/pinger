package central

import (
	"github.com/sdslabs/status/pkg/agent/proto"

	log "github.com/sirupsen/logrus"
)

// APIServer runs the server that pushes checks to agents.
type APIServer struct {
	Agents []*StatusAgent

	Host string
	Port int64
}

// Run starts the API server.
func (s *APIServer) Run() error {
	check := &proto.Check{
		Name: "http-test-check",

		Input: &proto.Check_Component{
			Type:  "HTTP",
			Value: "GET",
		},
		Output: &proto.Check_Component{
			Type:  "status_code",
			Value: "200",
		},
		Target: &proto.Check_Component{
			Value: "https://google.com",
		},

		Timeout:  10,
		Interval: 6,
	}

	for _, agent := range s.Agents {
		err := agent.PushCheckToAgent(check)
		if err != nil {
			log.Error(err)
		}
	}

	return nil
}

// NewAPIServer returns a default API server.
func NewAPIServer(port int) APIServer {
	agents := []*StatusAgent{
		NewStatusAgent("0.0.0.0", int64(port)),
	}

	return APIServer{
		Agents: agents,

		Host: "0.0.0.0",
		Port: 9000,
	}
}
