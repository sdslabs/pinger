// Package check creates check for various types of input such as
// HTTP, ICMP and Websocket. It probes the check target and parses
// the output into a common struct `Stats`.
package check

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	splitDelimeter = "="
	keyHeader      = "header"
	keyParameter   = "parameter"
	keyMessage     = "message"
	keyBody        = "body"
	keyStatusCode  = "status_code"
	keyTimeout     = "timeout"
	keyURL         = "url"
	keyAddress     = "address"
	keyECHO        = "ECHO"
)

type validationFunction func(string) error

func validateStatusCode(val string) error {
	intVal, err := strconv.Atoi(val)
	if err != nil {
		return fmt.Errorf("status code prvoided is not parsable: %s", err)
	}
	if intVal > 599 || intVal < 100 {
		return fmt.Errorf("status code is not valid (expected between: 100 - 599) got %d", intVal)
	}

	return nil
}

func validateBody(val string) error {
	if val == "" {
		return fmt.Errorf("cannot have empty body as check output")
	}

	return nil
}

func validateMessages(val string) error {
	if val == "" {
		return fmt.Errorf("cannot have empty messages as check output")
	}

	return nil
}

func validateKVPair(val string) error {
	kv := strings.SplitN(val, splitDelimeter, 2)
	if len(kv) != 2 {
		return fmt.Errorf("payload value is not valid, must have format: <KEY>=<VALUE>")
	}

	return nil
}
