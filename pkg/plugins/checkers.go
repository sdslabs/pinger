package plugins

import (
	// Register all the checkers here.
	_ "github.com/sdslabs/pinger/pkg/checker/dns"
	_ "github.com/sdslabs/pinger/pkg/checker/http"
	_ "github.com/sdslabs/pinger/pkg/checker/icmp"
	_ "github.com/sdslabs/pinger/pkg/checker/tcp"
	_ "github.com/sdslabs/pinger/pkg/checker/udp"
	_ "github.com/sdslabs/pinger/pkg/checker/ws"
)
