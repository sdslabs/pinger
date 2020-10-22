// Copyright (c) 2020 SDSLabs
// Use of this source code is governed by an MIT license
// details of which can be found in the LICENSE file.

package agent

import (
	// Register all the metrics alerters here.
	_ "github.com/sdslabs/pinger/pkg/alerter/discord"
	_ "github.com/sdslabs/pinger/pkg/alerter/slack"
)
