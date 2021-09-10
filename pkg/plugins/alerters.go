package plugins

import (
	// Register all the metrics alerters here.
	_ "github.com/sdslabs/pinger/pkg/alerter/discord"
	_ "github.com/sdslabs/pinger/pkg/alerter/mail"
	_ "github.com/sdslabs/pinger/pkg/alerter/slack"
)
