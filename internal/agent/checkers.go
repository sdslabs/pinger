// Copyright (c) 2020 SDSLabs
// Use of this source code is governed by an MIT license
// details of which can be found in the LICENSE file.

package agent

import (
	// Register all the checkers here.
	_ "github.com/sdslabs/status/internal/checker/dns"
	_ "github.com/sdslabs/status/internal/checker/http"
	_ "github.com/sdslabs/status/internal/checker/icmp"
	_ "github.com/sdslabs/status/internal/checker/tcp"
	_ "github.com/sdslabs/status/internal/checker/udp"
	_ "github.com/sdslabs/status/internal/checker/ws"
)
