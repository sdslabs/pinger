package plugins

import (
	// Register all the metrics exporters here.
	_ "github.com/sdslabs/pinger/pkg/exporter/influxdb"
	_ "github.com/sdslabs/pinger/pkg/exporter/log"
	_ "github.com/sdslabs/pinger/pkg/exporter/questdb"
	_ "github.com/sdslabs/pinger/pkg/exporter/timescale"
)
