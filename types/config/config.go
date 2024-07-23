package config

import (
	"strings"

	"gopkg.in/yaml.v3"

	databaseconfig "github.com/zenrocklabs/juno/database/config"
	loggingconfig "github.com/zenrocklabs/juno/logging/config"
	nodeconfig "github.com/zenrocklabs/juno/node/config"
	parserconfig "github.com/zenrocklabs/juno/parser/config"
)

var (
	// Cfg represents the configuration to be used during the execution
	Cfg Config
)

// Config defines all necessary juno configuration parameters.
type Config struct {
	bytes []byte

	Chain    ChainConfig           `yaml:"chain"`
	Node     nodeconfig.Config     `yaml:"node"`
	Parser   parserconfig.Config   `yaml:"parsing"`
	Database databaseconfig.Config `yaml:"database"`
	Logging  loggingconfig.Config  `yaml:"logging"`
}

// NewConfig builds a new Config instance
func NewConfig(
	nodeCfg nodeconfig.Config,
	chainCfg ChainConfig, dbConfig databaseconfig.Config,
	parserConfig parserconfig.Config, loggingConfig loggingconfig.Config,
) Config {
	return Config{
		Node:     nodeCfg,
		Chain:    chainCfg,
		Database: dbConfig,
		Parser:   parserConfig,
		Logging:  loggingConfig,
	}
}

func DefaultConfig() Config {
	cfg := NewConfig(
		nodeconfig.DefaultConfig(),
		DefaultChainConfig(), databaseconfig.DefaultDatabaseConfig(),
		parserconfig.DefaultParsingConfig(), loggingconfig.DefaultLoggingConfig(),
	)

	bz, err := yaml.Marshal(cfg)
	if err != nil {
		panic(err)
	}

	cfg.bytes = bz
	return cfg
}

func (c Config) GetBytes() ([]byte, error) {
	return c.bytes, nil
}

// ---------------------------------------------------------------------------------------------------------------------

type ChainConfig struct {
	Bech32Prefix string   `yaml:"bech32_prefix"`
	Modules      []string `yaml:"modules"`
}

// NewChainConfig returns a new ChainConfig instance
func NewChainConfig(bech32Prefix string, modules []string) ChainConfig {
	return ChainConfig{
		Bech32Prefix: bech32Prefix,
		Modules:      modules,
	}
}

// DefaultChainConfig returns the default instance of ChainConfig
func DefaultChainConfig() ChainConfig {
	return NewChainConfig("cosmos", nil)
}

func (cfg ChainConfig) IsModuleEnabled(moduleName string) bool {
	for _, module := range cfg.Modules {
		if strings.EqualFold(module, moduleName) {
			return true
		}
	}

	return false
}
