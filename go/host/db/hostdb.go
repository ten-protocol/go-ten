package db

import (
	"fmt"
	"os"

	hostcommon "github.com/obscuronet/go-obscuro/go/common/host"
	"github.com/obscuronet/go-obscuro/go/common/log"

	"github.com/obscuronet/go-obscuro/go/config"

	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/ethdb/leveldb"
	gethlog "github.com/ethereum/go-ethereum/log"
	gethmetrics "github.com/ethereum/go-ethereum/metrics"
	"github.com/obscuronet/go-obscuro/go/common/gethdb"
)

// Schema keys, in alphabetical order.
var (
	blockHeaderPrefix       = []byte("b")
	batchHeaderPrefix       = []byte("ba")
	batchHashPrefix         = []byte("bh")
	batchNumberPrefix       = []byte("bn")
	batchPrefix             = []byte("bp")
	batchHashForSeqNoPrefix = []byte("bs")
	batchTxHashesPrefix     = []byte("bt")
	headBatch               = []byte("hb")
	totalTransactionsKey    = []byte("t")
	rollupHeaderPrefix      = []byte("rh")
	tipRollupHash           = []byte("tr")
)

// DB allows to access the nodes public nodeDB
type DB struct {
	kvStore     ethdb.KeyValueStore
	logger      gethlog.Logger
	batchWrites gethmetrics.Gauge
	batchReads  gethmetrics.Gauge
	blockWrites gethmetrics.Gauge
	blockReads  gethmetrics.Gauge
}

func (db *DB) Start() error {
	return nil
}

// Stop is especially important for graceful shutdown of LevelDB as it may flush data to disk that is currently in cache
func (db *DB) Stop() {
	db.logger.Info("Closing the host DB.")
	err := db.kvStore.Close()
	if err != nil {
		db.logger.Error("Error closing the host DB.", log.ErrKey, err)
	}
}

func (db *DB) HealthStatus() hostcommon.HealthStatus {
	// always healthy for now, satisfies Service interface
	return &hostcommon.BasicErrHealthStatus{ErrMsg: ""}
}

// InitialiseMetrics registers the host DB metrics with the given registry (this is called by the metrics service when it starts)
func (db *DB) InitialiseMetrics(registry gethmetrics.Registry) {
	// register db metrics
	db.batchWrites = gethmetrics.NewRegisteredGauge("host/db/batch/writes", registry)
	db.batchReads = gethmetrics.NewRegisteredGauge("host/db/batch/reads", registry)
	db.blockWrites = gethmetrics.NewRegisteredGauge("host/db/block/writes", registry)
	db.blockReads = gethmetrics.NewRegisteredGauge("host/db/block/reads", registry)
}

func CreateDBFromConfig(cfg *config.HostConfig, logger gethlog.Logger) (*DB, error) {
	if err := validateDBConf(cfg); err != nil {
		return nil, err
	}
	if cfg.UseInMemoryDB {
		logger.Info("UseInMemoryDB flag is true, data will not be persisted. Creating in-memory database...")
		return NewInMemoryDB(logger), nil
	}
	return NewLevelDBBackedDB(cfg.LevelDBPath, logger)
}

func validateDBConf(cfg *config.HostConfig) error {
	if cfg.UseInMemoryDB && cfg.LevelDBPath != "" {
		return fmt.Errorf("useInMemoryDB=true so levelDB will not be used and no path is needed, but levelDBPath=%s", cfg.LevelDBPath)
	}
	return nil
}

// NewInMemoryDB returns a new instance of the Node DB
func NewInMemoryDB(logger gethlog.Logger) *DB {
	return newDB(gethdb.NewMemDB(), logger)
}

// NewLevelDBBackedDB creates a persistent DB for the host, if dbPath == "" it will generate a temp file
func NewLevelDBBackedDB(dbPath string, logger gethlog.Logger) (*DB, error) {
	var err error
	if dbPath == "" {
		// todo (#1618) - we should remove this option before prod, if you want a temp DB it should be wired in via the config
		dbPath, err = os.MkdirTemp("", "leveldb_*")
		if err != nil {
			return nil, fmt.Errorf("could not create temp leveldb directory - %w", err)
		}
		logger.Warn("dbPath was empty, created temp dir for persistence", "dbPath", dbPath)
	}
	// determine if a db file already exists, we don't want to overwrite it
	_, err = os.Stat(dbPath)
	dbDesc := "new"
	if err == nil {
		dbDesc = "existing"
	}

	// todo (#1618) - these should be configs
	cache := 128
	handles := 128
	db, err := leveldb.New(dbPath, cache, handles, "host", false)
	if err != nil {
		return nil, fmt.Errorf("could not create leveldb - %w", err)
	}
	logger.Info(fmt.Sprintf("Opened %s level db dir at %s", dbDesc, dbPath))
	return newDB(&ObscuroLevelDB{db: db}, logger), nil
}

func newDB(kvStore ethdb.KeyValueStore, logger gethlog.Logger) *DB {
	return &DB{
		kvStore: kvStore,
		logger:  logger,
	}
}
