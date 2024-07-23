package v4

import (
	databaseconfig "github.com/zenrocklabs/juno/database/config"
	loggingconfig "github.com/zenrocklabs/juno/logging/config"
	"github.com/zenrocklabs/juno/modules/pruning"
	"github.com/zenrocklabs/juno/modules/telemetry"
	nodeconfig "github.com/zenrocklabs/juno/node/config"
	parserconfig "github.com/zenrocklabs/juno/parser/config"
	pricefeedconfig "github.com/zenrocklabs/juno/pricefeed"
	"github.com/zenrocklabs/juno/types/config"
)

// Config defines all necessary juno configuration parameters.
type Config struct {
	Chain    config.ChainConfig    `yaml:"chain"`
	Node     nodeconfig.Config     `yaml:"node"`
	Parser   parserconfig.Config   `yaml:"parsing"`
	Database databaseconfig.Config `yaml:"database"`
	Logging  loggingconfig.Config  `yaml:"logging"`

	// The following are there to support modules which config are present if they are enabled

	Telemetry *telemetry.Config       `yaml:"telemetry,omitempty"`
	Pruning   *pruning.Config         `yaml:"pruning,omitempty"`
	PriceFeed *pricefeedconfig.Config `yaml:"pricefeed,omitempty"`
}
