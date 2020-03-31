// Package probes defines multiple probes used in the various kinds of checks
// that can be deployed. These can be ICMP, HTTP, Websocket etc.
package probes

import "net"

// errIsTimeout returns true for a timeout error.
func errIsTimeout(err error) bool {
	netErr, ok := err.(net.Error)
	return ok && netErr.Timeout()
}
