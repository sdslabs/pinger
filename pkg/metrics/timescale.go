package metrics

import (
	"time"

	"github.com/sdslabs/status/pkg/controller"
)

// TimescaleExporter for exporting metrics to timescale db.
type TimescaleExporter struct {
	*controller.Manager
	Quit     chan bool
	Interval time.Duration
}

// Rest of the implementation is in @/pkg/database to use the same schema as the web
// application. The metrics data can later be easily exported from the timescale db.
