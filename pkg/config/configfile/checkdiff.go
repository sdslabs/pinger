package configfile

import (
	"github.com/sdslabs/pinger/pkg/config"
)

type CheckDiff struct {
	Additions []config.Check `mapstructure:"additions" json:"additions"`
	Removals  []config.Check `mapstructure:"removals" json:"removals"`
}
