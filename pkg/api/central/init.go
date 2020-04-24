// Package central contains the central API server for status agent.
package central

import (
	"fmt"
	"net"
)

// getAddr returns the address combining host and port.
func getAddr(host string, port int) string {
	portStr := fmt.Sprintf("%d", port)
	return net.JoinHostPort(host, portStr)
}
