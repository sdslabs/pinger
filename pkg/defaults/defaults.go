package defaults

import (
	"time"
)

var (
	// Minimum value of the retry interval that the controller accepts.
	MinControllerRetryInterval time.Duration = time.Second * 4

	// ControllerRetryInterval is the default value of the retry interval used by
	// the controller.
	ControllerRetryInterval time.Duration = time.Second * 30

	// MinControllerTimeout is the minimum value of the permissible timeout
	// that the controller accepts.
	MinControllerTimeout time.Duration = time.Second * 5

	// DefaultHTTPMethod is the default method to use for HTTP Input
	DefaultHTTPMethod = "GET"
)
