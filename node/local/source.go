package local

import (
	"fmt"
	"os"
	"path"
	"reflect"
	"unsafe"

	storetypes "cosmossdk.io/store/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	storemetrics "cosmossdk.io/store/metrics"
	cfg "github.com/cometbft/cometbft/config"
	tmlog "github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmstore "github.com/cometbft/cometbft/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/viper"

	"github.com/zenrocklabs/juno/node"
)

var (
	_ node.Source = &Source{}
)

// Source represents the Source interface implementation that reads the data from a local node
type Source struct {
	Initialized bool

	StoreDB dbm.DB

	Codec codec.Codec

	BlockStore *tmstore.BlockStore
	Logger     log.Logger
	Cms        store.CommitMultiStore
}

// NewSource returns a new Source instance
func NewSource(home string, codec codec.Codec) (*Source, error) {
	levelDB, err := dbm.NewGoLevelDB("application", path.Join(home, "data"), nil)
	if err != nil {
		return nil, err
	}

	tmCfg, err := parseConfig(home)
	if err != nil {
		return nil, err
	}

	blockStoreDB, err := cfg.DefaultDBProvider(&cfg.DBContext{ID: "blockstore", Config: tmCfg})
	if err != nil {
		return nil, err
	}

	logger := log.NewLogger(tmlog.NewSyncWriter(os.Stdout)).With("module", "explorer")

	return &Source{
		StoreDB: levelDB,

		Codec: codec,

		BlockStore: tmstore.NewBlockStore(blockStoreDB),
		Logger:     logger,
		Cms:        store.NewCommitMultiStore(levelDB, logger, storemetrics.NewNoOpMetrics()),
	}, nil
}

func parseConfig(home string) (*cfg.Config, error) {
	viper.SetConfigFile(path.Join(home, "config", "config.toml"))

	conf := cfg.DefaultConfig()
	err := viper.Unmarshal(conf)
	if err != nil {
		return nil, err
	}
	conf.SetRoot(conf.RootDir)

	err = conf.ValidateBasic()
	if err != nil {
		return nil, fmt.Errorf("error in config file: %v", err)
	}

	return conf, nil
}

// Type implements keeper.Source
func (k Source) Type() string {
	return node.LocalKeeper
}

func getFieldUsingReflection(app interface{}, fieldName string) interface{} {
	fv := reflect.ValueOf(app).Elem().FieldByName(fieldName)
	return reflect.NewAt(fv.Type(), unsafe.Pointer(fv.UnsafeAddr())).Elem().Interface()
}

// MountKVStores allows to register the KV stores using the same KVStoreKey instances
// that are used inside the given app. To do so, this method uses the reflection to access
// the field with the specified name inside the given app. Such field must be of type
// map[string]*sdk.KVStoreKey and is commonly named something similar to "keys"
func (k Source) MountKVStores(app interface{}, fieldName string) error {
	keys, ok := getFieldUsingReflection(app, fieldName).(map[string]*storetypes.KVStoreKey)
	if !ok {
		return fmt.Errorf("error while getting keys")
	}

	for _, key := range keys {
		k.Cms.MountStoreWithDB(key, storetypes.StoreTypeIAVL, nil)
	}

	return nil
}

// MountTransientStores allows to register the Transient stores using the same TransientStoreKey instances
// that are used inside the given app. To do so, this method uses the reflection to access
// the field with the specified name inside the given app. Such field must be of type
// map[string]*sdk.TransientStoreKey and is commonly named something similar to "tkeys"
func (k Source) MountTransientStores(app interface{}, fieldName string) error {
	tkeys, ok := getFieldUsingReflection(app, fieldName).(map[string]*storetypes.TransientStoreKey)
	if !ok {
		return fmt.Errorf("error while getting transient keys")
	}

	for _, key := range tkeys {
		k.Cms.MountStoreWithDB(key, storetypes.StoreTypeTransient, nil)
	}

	return nil
}

// MountMemoryStores allows to register the Memory stores using the same MemoryStoreKey instances
// that are used inside the given app. To do so, this method uses the reflection to access
// the field with the specified name inside the given app. Such field must be of type
// map[string]*sdk.MemoryStoreKey and is commonly named something similar to "memkeys"
func (k Source) MountMemoryStores(app interface{}, fieldName string) error {
	memKeys, ok := getFieldUsingReflection(app, fieldName).(map[string]*storetypes.MemoryStoreKey)
	if !ok {
		return fmt.Errorf("error while getting memory keys")
	}

	for _, key := range memKeys {
		k.Cms.MountStoreWithDB(key, storetypes.StoreTypeMemory, nil)
	}

	return nil
}

// InitStores initializes the stores by mounting the various keys that have been specified.
// This method MUST be called before using any method that relies on the local storage somehow.
func (k Source) InitStores() error {
	return k.Cms.LoadLatestVersion()
}

// LoadHeight loads the given height from the store.
// It returns a new Context that can be used to query the data, or an error if something wrong happens.
func (k Source) LoadHeight(height int64) (sdk.Context, error) {
	var err error
	var cms store.CacheMultiStore
	if height > 0 {
		cms, err = k.Cms.CacheMultiStoreWithVersion(height)
		if err != nil {
			return sdk.Context{}, err
		}
	} else {
		cms, err = k.Cms.CacheMultiStoreWithVersion(k.BlockStore.Height())
		if err != nil {
			return sdk.Context{}, err
		}
	}

	return sdk.NewContext(cms, tmproto.Header{}, false, k.Logger), nil
}
