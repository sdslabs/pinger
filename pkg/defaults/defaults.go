// Package defaults contains the default values for various elements.
package defaults

import (
	"os"
	"path/filepath"
	"time"
)

var (
	// MinControllerRetryInterval is the minimum value of the retry interval that the controller accepts.
	MinControllerRetryInterval = time.Second * 4

	// ControllerRetryInterval is the default value of the retry interval used by
	// the controller.
	ControllerRetryInterval = time.Second * 30

	// MinControllerTimeout is the minimum value of the permissible timeout
	// that the controller accepts.
	MinControllerTimeout = time.Second * 5

	// HTTPMethod is the default method to use for HTTP Input
	HTTPMethod = "GET"

	// GRPCRequestTimeout is the default timeout for the GRPC request we
	// are making between the server and the agent and vice versa.
	GRPCRequestTimeout = time.Second * 10

	// HTTPProbeTimeout is the http request timeout for HTTP prober.
	HTTPProbeTimeout = time.Second * 10

	// ICMPProbeTimeout is the icmp request timeout for ICMP prober.
	ICMPProbeTimeout = time.Second * 10

	// WSProbeTimeout is the websocket request timeout for WS prober.
	WSProbeTimeout = time.Second * 10

	// TCPProbeTimeout is the tcp request timeout for TCP prober.
	TCPProbeTimeout = time.Second * 10

	// UDPProbeTimeout is the tcp request timeout for TCP prober.
	UDPProbeTimeout = time.Second * 10

	// DNSProbeTimeout is the dns request timeout for DNS prober.
	DNSProbeTimeout = time.Second * 10

	// JWTExpireInterval is interval for which JWT is valid
	JWTExpireInterval = time.Hour // 1 hour

	// JWTRefreshInterval is the interval for which an expired token can be refreshed
	// Total refresh time = JWTRefreshInterval + JWTExpireInterval
	JWTRefreshInterval = time.Hour * 24 * 28 // 4 weeks by default

	// JWTAuthType is the token type of the authorization header
	JWTAuthType = "Bearer"

	// AgentMetricsBackend is the default metrics to be used with the agent.
	AgentMetricsBackend = "timescale"

	// AgentMetricsHost is the default host to run agent metrics on.
	AgentMetricsHost = "127.0.0.1"

	// AgentMetricsPort is the default port for the agent to host prometheus metrics enpoint on.
	// This corresponds to the default metrics, i.e., timescale db (postgres).
	AgentMetricsPort = 5432

	// AgentMetricsInterval is the default interval after which metrics will be pushed into the db.
	// In case of prometheus this is the interval after which memory is cleaned up.
	AgentMetricsInterval = 2 * time.Minute

	// AgentPort is the default value of port to run the status agent on
	AgentPort = 9009

	// CentralAPIPort is default the central api server port.
	CentralAPIPort = 9010

	// AppAPIPort is the default application api server port.
	AppAPIPort = 8080

	// AppAPIDBHost is the default host for app db.
	AppAPIDBHost = "127.0.0.1"

	// AppAPIDBPort is the default port for app db.
	AppAPIDBPort = 5432

	// ControllerType is the default type of the controller, other types can be
	// specified based on the type of checks/probes that the controller executes
	ControllerType = "default"

	// InvalidDuration is the default value of time duration which will be considered
	// as a failed controller execution
	InvalidDuration = time.Second * 360

	// AppConfigPath is the default path of the config file for the app api server.
	AppConfigPath = filepath.Join(os.Getenv("HOME"), "/.status-app.yml")

	// AgentConfigPath is the default path of the config file for the status page
	// agent. This is mostly used in case when we run the agent in a standalone mode.
	AgentConfigPath = filepath.Join(os.Getenv("HOME"), "/.status-agent.yml")

	// StandaloneUserEmail is the email of the user corresponding to which the standalone
	// agent metrics will be saved in the timescale db.
	StandaloneUserEmail = "standalone@status.agent"

	// StandaloneUserName is the name of the user corresponding to which the standalone
	// agent metrics will be saved in the timescale db.
	StandaloneUserName = "Standalone Status Agent"
)
