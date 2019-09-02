package defaults

import (
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

	// DefaultHTTPMethod is the default method to use for HTTP Input
	DefaultHTTPMethod = "GET"

	// DefaultGRPCRequestTimeout is the default timeout for the GRPC request we
	// are making between the server and the agent and vice versa.
	DefaultGRPCRequestTimeout = time.Second * 10

	// DefaultHTTPProbeTimeout is the http request timeout for HTTP prober.
	DefaultHTTPProbeTimeout = time.Second * 10

	// JWTExpireInterval is interval for which JWT is valid
	JWTExpireInterval = time.Hour * 24
)
