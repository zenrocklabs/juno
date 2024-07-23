package v3

import (
	loggingconfig "github.com/zenrocklabs/juno/logging/config"
	"github.com/zenrocklabs/juno/modules/pruning"
	"github.com/zenrocklabs/juno/modules/telemetry"
	nodeconfig "github.com/zenrocklabs/juno/node/config"
	parserconfig "github.com/zenrocklabs/juno/parser/config"
	pricefeedconfig "github.com/zenrocklabs/juno/pricefeed"
	"github.com/zenrocklabs/juno/types/config"
)

type Config struct {
	Chain    config.ChainConfig   `yaml:"chain"`
	Node     nodeconfig.Config    `yaml:"node"`
	Parser   parserconfig.Config  `yaml:"parsing"`
	Database DatabaseConfig       `yaml:"database"`
	Logging  loggingconfig.Config `yaml:"logging"`

	// The following are there to support modules which config are present if they are enabled

	Telemetry *telemetry.Config       `yaml:"telemetry,omitempty"`
	Pruning   *pruning.Config         `yaml:"pruning,omitempty"`
	PriceFeed *pricefeedconfig.Config `yaml:"pricefeed,omitempty"`
}

type DatabaseConfig struct {
	Name               string `yaml:"name"`
	Host               string `yaml:"host"`
	Port               int64  `yaml:"port"`
	User               string `yaml:"user"`
	Password           string `yaml:"password"`
	SSLMode            string `yaml:"ssl_mode,omitempty"`
	Schema             string `yaml:"schema,omitempty"`
	MaxOpenConnections int    `yaml:"max_open_connections"`
	MaxIdleConnections int    `yaml:"max_idle_connections"`
	PartitionSize      int64  `yaml:"partition_size,omitempty"`
	PartitionBatchSize int64  `yaml:"partition_batch,omitempty"`
}
