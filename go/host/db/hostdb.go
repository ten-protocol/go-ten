package db

import (
	"os"

	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/ethdb/leveldb"
	"github.com/obscuronet/go-obscuro/go/common/gethdb"
	"github.com/obscuronet/go-obscuro/go/common/log"

	gethlog "github.com/ethereum/go-ethereum/log"
	gethmetrics "github.com/ethereum/go-ethereum/metrics"
)

// Schema keys, in alphabetical order.
var (
	blockHeaderPrefix    = []byte("b")
	batchHeaderPrefix    = []byte("ba")
	batchHashPrefix      = []byte("bh")
	batchNumberPrefix    = []byte("bn")
	batchPrefix          = []byte("bp")
	batchTxHashesPrefix  = []byte("bt")
	headBatch            = []byte("hb")
	totalTransactionsKey = []byte("t")
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

// NewInMemoryDB returns a new instance of the Node DB
func NewInMemoryDB(regMetrics gethmetrics.Registry) *DB {
	return &DB{
		kvStore:     gethdb.NewMemDB(),
		batchWrites: gethmetrics.NewRegisteredGauge("host/db/batch/writes", regMetrics),
		batchReads:  gethmetrics.NewRegisteredGauge("host/db/batch/reads", regMetrics),
		blockWrites: gethmetrics.NewRegisteredGauge("host/db/block/writes", regMetrics),
		blockReads:  gethmetrics.NewRegisteredGauge("host/db/block/reads", regMetrics),
	}
}

func NewLevelDBBackedDB(logger gethlog.Logger) *DB {
	// todo, all these should be configs
	f, err := os.MkdirTemp("", "leveldb_*")
	if err != nil {
		logger.Crit("Could not creat temp leveldb directory.", log.ErrKey, err)
	}
	cache := 128
	handles := 128
	db, err := leveldb.New(f, cache, handles, "host", false)
	if err != nil {
		logger.Crit("Could not create leveldb.", log.ErrKey, err)
	}

	return &DB{
		kvStore: db,
		logger:  logger,
	}
}
