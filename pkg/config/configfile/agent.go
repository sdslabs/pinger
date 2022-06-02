package configfile

import (
	"time"

	"github.com/sdslabs/pinger/pkg/config"
)

// AgentPage defines the configuration of the stand alone page deployed by
// the agent.
type AgentPage struct {
	Deploy         bool     `mapstructure:"deploy" json:"deploy"`
	Port           uint16   `mapstructure:"port" json:"port"`
	AllowedOrigins []string `mapstructure:"allowed_origins" json:"allowed_origins"`
	Name           string   `mapstructure:"name" json:"name"`
	Media          string   `mapstructure:"media" json:"media"`
	Logo           string   `mapstructure:"logo" json:"logo"`
	Favicon        string   `mapstructure:"favicon" json:"favicon"`
	Website        string   `mapstructure:"website" json:"website"`
}

// Agent represents the configuration for an agent.
type Agent struct {
	Standalone         bool                   `mapstructure:"standalone" json:"standalone"`
	CentralServerRedis config.DBConn          `mapstructure:"central_server_redis" json:"central_server_redis"`
	NetInterface       string                 `mapstructure:"net_interface" json:"net_interface"`
	Port               uint16                 `mapstructure:"port" json:"port"`
	Metrics            config.MetricsProvider `mapstructure:"metrics" json:"metrics"`
	Alerts             []config.AlertProvider `mapstructure:"alerts" json:"alerts"`
	Interval           time.Duration          `mapstructure:"interval" json:"interval"`
	Checks             []config.Check         `mapstructure:"checks" json:"checks"`
	Page               AgentPage              `mapstructure:"page" json:"page"`
}
