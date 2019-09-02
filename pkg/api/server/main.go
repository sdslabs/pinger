package server

import (
	"github.com/sdslabs/status/pkg/api/agent/proto"

	log "github.com/sirupsen/logrus"
)

type ApiServer struct {
	Agents []*StatusAgent

	Host string
	Port int64
}

func (s *ApiServer) Run() error {
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

func NewApiServer() ApiServer {
	agents := []*StatusAgent{
		NewStatusAgent("0.0.0.0", 9009),
	}

	return ApiServer{
		Agents: agents,

		Host: "0.0.0.0",
		Port: 9000,
	}
}
