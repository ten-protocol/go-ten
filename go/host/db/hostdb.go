package db

import (
	"os"

	gethlog "github.com/ethereum/go-ethereum/log"

	"github.com/obscuronet/go-obscuro/go/common/log"

	"github.com/ethereum/go-ethereum/ethdb/leveldb"

	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/ethdb/memorydb"
)

// Schema keys, in alphabetical order.
var (
	headBlock            = []byte("hb")
	headRollup           = []byte("hr")
	headBatch            = []byte("hba")
	blockHeaderPrefix    = []byte("b")
	batchHeaderPrefix    = []byte("ba")
	batchHashPrefix      = []byte("bh")
	batchNumberPrefix    = []byte("bn")
	batchTxHashesPrefix  = []byte("bt")
	rollupHeaderPrefix   = []byte("r")
	rollupTxHashesPrefix = []byte("rt")
	rollupHashPrefix     = []byte("rh")
	rollupNumberPrefix   = []byte("rn")
	totalTransactionsKey = []byte("t")
)

// DB allows to access the nodes public nodeDB
type DB struct {
	kvStore ethdb.KeyValueStore
	logger  gethlog.Logger
}

// NewInMemoryDB returns a new instance of the Node DB
func NewInMemoryDB() *DB {
	return &DB{
		kvStore: memorydb.New(),
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
	db, err := leveldb.New(f, cache, handles, "obscuro_host", false)
	if err != nil {
		logger.Crit("Could not create leveldb.", log.ErrKey, err)
	}

	return &DB{
		kvStore: db,
		logger:  logger,
	}
}
