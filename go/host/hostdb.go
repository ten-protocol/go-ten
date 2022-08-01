package host

import (
	"bytes"
	"math/big"
	"os"

	"github.com/obscuronet/go-obscuro/go/common/log"

	"github.com/ethereum/go-ethereum/ethdb/leveldb"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/ethdb/memorydb"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/go-obscuro/go/common"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

// DB allows to access the nodes public nodeDB
type DB struct {
	kvStore ethdb.KeyValueStore
}

// NewInMemoryDB returns a new instance of the Node DB
func NewInMemoryDB() *DB {
	return &DB{
		kvStore: memorydb.New(),
	}
}

func NewLevelDBBackedDB() *DB {
	// todo, all these should be configs
	f, err := os.MkdirTemp("", "leveldb_*")
	if err != nil {
		log.Panic("Could not creat temp leveldb directory. Cause %s", err)
	}
	cache := 128
	handles := 128
	db, err := leveldb.New(f, cache, handles, "obscuro_host", false)
	if err != nil {
		log.Panic("Could not create leveldb. Cause: %s", err)
	}
	return &DB{
		kvStore: db,
	}
}

// GetCurrentBlockHead returns the current block header (head) of the Node
func (db *DB) GetCurrentBlockHead() *types.Header {
	head := readHeadBlock(db.kvStore)
	if head == nil {
		return nil
	}
	return readBlockHeader(db.kvStore, *head)
}

// GetBlockHeader returns the block header given the Hash
func (db *DB) GetBlockHeader(hash gethcommon.Hash) *types.Header {
	return readBlockHeader(db.kvStore, hash)
}

// AddBlockHeader adds a types.Header to the known headers
func (db *DB) AddBlockHeader(header *types.Header) {
	b := db.kvStore.NewBatch()
	writeBlockHeader(b, header)

	// update the head if the new height is greater than the existing one
	currentBlockHead := db.GetCurrentBlockHead()
	if currentBlockHead == nil || currentBlockHead.Number.Int64() <= header.Number.Int64() {
		writeHeadBlock(b, header.Hash())
	}

	if err := b.Write(); err != nil {
		log.Panic("Could not write rollup . Cause %s", err)
	}
}

// GetCurrentRollupHead returns the current rollup header (head) of the Node
func (db *DB) GetCurrentRollupHead() *common.Header {
	head := readHeadRollup(db.kvStore)
	if head == nil {
		return nil
	}
	return readRollupHeader(db.kvStore, *head)
}

// GetRollupHeader returns the rollup header given the Hash
func (db *DB) GetRollupHeader(hash gethcommon.Hash) *common.Header {
	return readRollupHeader(db.kvStore, hash)
}

// AddRollupHeader adds a RollupHeader to the known headers
func (db *DB) AddRollupHeader(header *common.Header, txHashes []gethcommon.Hash) {
	b := db.kvStore.NewBatch()
	writeRollupHeader(b, header)
	writeRollupHashByNumber(b, header)
	for _, txHash := range txHashes {
		writeRollupHashByTxHash(b, txHash, header)
	}

	// update the head if the new height is greater than the existing one
	currentRollupHead := db.GetCurrentRollupHead()
	if currentRollupHead == nil || currentRollupHead.Number.Int64() <= header.Number.Int64() {
		writeHeadRollup(b, header.Hash())
	}

	if err := b.Write(); err != nil {
		log.Panic("Could not write rollup . Cause %s", err)
	}
}

func (db *DB) AddSubmittedRollup(hash gethcommon.Hash) {
	err := db.kvStore.Put(submittedRollupHeaderKey(hash), []byte{})
	if err != nil {
		log.Panic("Could not save submitted rollup. Cause: %s", err)
	}
}

func (db *DB) WasSubmitted(hash gethcommon.Hash) bool {
	f, err := db.kvStore.Has(submittedRollupHeaderKey(hash))
	if err != nil {
		log.Panic("Could not retrieve submitted rollup. Cause: %s", err)
	}
	return f
}

// GetRollupHashByNumber returns the hash of a rollup given its number
func (db *DB) GetRollupHashByNumber(number *big.Int) *gethcommon.Hash {
	return readRollupHashByNumber(db.kvStore, number)
}

// GetRollupHashByTxHash returns the hash of a rollup given the hash of a transaction it contains
func (db *DB) GetRollupHashByTxHash(txHash gethcommon.Hash) *gethcommon.Hash {
	return readRollupHashByTxHash(db.kvStore, txHash)
}

// schema
var (
	blockHeaderPrefix      = []byte("b")
	rollupHeaderPrefix     = []byte("r")
	headBlock              = []byte("hb")
	headRollup             = []byte("hr")
	submittedRollupPrefix  = []byte("s")
	rollupHashNumberPrefix = []byte("rhn")
	rollupHashTxHashPrefix = []byte("rht")
)

// headerKey = rollupHeaderPrefix  + hash
func rollupHeaderKey(hash gethcommon.Hash) []byte {
	return append(rollupHeaderPrefix, hash.Bytes()...)
}

// headerKey = blockHeaderPrefix  + hash
func blockHeaderKey(hash gethcommon.Hash) []byte {
	return append(blockHeaderPrefix, hash.Bytes()...)
}

// headerKey = submittedRollupPrefix  + hash
func submittedRollupHeaderKey(hash gethcommon.Hash) []byte {
	return append(submittedRollupPrefix, hash.Bytes()...)
}

// headerKey = rollupHashNumberPrefix + number
func rollupHashNumberKey(num *big.Int) []byte {
	return append(rollupHashNumberPrefix, []byte(num.String())...)
}

// headerKey = rollupHashTxHashPrefix + hash
func rollupHashTxHashKey(txHash gethcommon.Hash) []byte {
	return append(rollupHashTxHashPrefix, txHash.Bytes()...)
}

// WriteBlockHeader stores a block header into the database
func writeBlockHeader(db ethdb.KeyValueWriter, header *types.Header) {
	// Write the encoded header
	data, err := rlp.EncodeToBytes(header)
	if err != nil {
		log.Panic("could not encode block header. Cause: %s", err)
	}
	key := blockHeaderKey(header.Hash())
	if err := db.Put(key, data); err != nil {
		log.Panic("could not put header in DB. Cause: %s", err)
	}
}

// ReadBlockHeader retrieves the rollup header corresponding to the hash.
func readBlockHeader(db ethdb.KeyValueReader, hash gethcommon.Hash) *types.Header {
	f, err := db.Has(blockHeaderKey(hash))
	if err != nil {
		log.Panic("could not retrieve block header. Cause: %s", err)
	}
	if !f {
		return nil
	}
	data, err := db.Get(blockHeaderKey(hash))
	if err != nil {
		log.Panic("could not retrieve block header. Cause: %s", err)
	}
	if len(data) == 0 {
		return nil
	}
	header := new(types.Header)
	if err := rlp.Decode(bytes.NewReader(data), header); err != nil {
		log.Panic("could not decode block header. Cause: %s", err)
	}
	return header
}

// WriteRollupHeader stores a rollup header into the database
func writeRollupHeader(db ethdb.KeyValueWriter, header *common.Header) {
	// Write the encoded header
	data, err := rlp.EncodeToBytes(header)
	if err != nil {
		log.Panic("could not encode rollup header. Cause: %s", err)
	}
	key := rollupHeaderKey(header.Hash())
	if err := db.Put(key, data); err != nil {
		log.Panic("could not put header in DB. Cause: %s", err)
	}
}

// ReadRollupHeader retrieves the rollup header corresponding to the hash.
func readRollupHeader(db ethdb.KeyValueReader, hash gethcommon.Hash) *common.Header {
	f, err := db.Has(rollupHeaderKey(hash))
	if err != nil {
		log.Panic("could not retrieve rollup header. Cause: %s", err)
	}
	if !f {
		return nil
	}
	data, err := db.Get(rollupHeaderKey(hash))
	if err != nil {
		log.Panic("could not retrieve rollup header. Cause: %s", err)
	}
	if len(data) == 0 {
		return nil
	}
	header := new(common.Header)
	if err := rlp.Decode(bytes.NewReader(data), header); err != nil {
		log.Panic("could not decode rollup header. Cause: %s", err)
	}
	return header
}

func readHeadBlock(db ethdb.KeyValueReader) *gethcommon.Hash {
	f, err := db.Has(headBlock)
	if err != nil {
		log.Panic("could not retrieve head block. Cause: %s", err)
	}
	if !f {
		return nil
	}
	value, err := db.Get(headBlock)
	if err != nil {
		log.Panic("could not retrieve head block. Cause: %s", err)
	}
	h := gethcommon.BytesToHash(value)
	return &h
}

func writeHeadBlock(db ethdb.KeyValueWriter, val gethcommon.Hash) {
	err := db.Put(headBlock, val.Bytes())
	if err != nil {
		log.Panic("could not write head block. Cause: %s", err)
	}
}

func readHeadRollup(db ethdb.KeyValueReader) *gethcommon.Hash {
	f, err := db.Has(headRollup)
	if err != nil {
		log.Panic("could not retrieve head rollup. Cause: %s", err)
	}
	if !f {
		return nil
	}
	value, err := db.Get(headRollup)
	if err != nil {
		log.Panic("could not retrieve head rollup. Cause: %s", err)
	}
	h := gethcommon.BytesToHash(value)
	return &h
}

func writeHeadRollup(db ethdb.KeyValueWriter, val gethcommon.Hash) {
	err := db.Put(headRollup, val.Bytes())
	if err != nil {
		log.Panic("could not write head rollup. Cause: %s", err)
	}
}

// Stores the hash of a rollup into the database, keyed by the rollup's number
func writeRollupHashByNumber(db ethdb.KeyValueWriter, header *common.Header) {
	key := rollupHashNumberKey(header.Number)
	if err := db.Put(key, header.Hash().Bytes()); err != nil {
		log.Panic("could not put header in DB. Cause: %s", err)
	}
}

// Stores the hash of a rollup into the database, keyed by the hashes of the transactions in the rollup
func writeRollupHashByTxHash(db ethdb.KeyValueWriter, txHash gethcommon.Hash, header *common.Header) {
	key := rollupHashTxHashKey(txHash)
	if err := db.Put(key, header.Hash().Bytes()); err != nil {
		log.Panic("could not put header in DB. Cause: %s", err)
	}
}

// Retrieves the hash for the rollup with the given number.
func readRollupHashByNumber(db ethdb.KeyValueReader, number *big.Int) *gethcommon.Hash {
	f, err := db.Has(rollupHashNumberKey(number))
	if err != nil {
		log.Panic("could not retrieve rollup hash. Cause: %s", err)
	}
	if !f {
		return nil
	}
	data, err := db.Get(rollupHashNumberKey(number))
	if err != nil {
		log.Panic("could not retrieve rollup hash. Cause: %s", err)
	}
	if len(data) == 0 {
		return nil
	}
	hash := gethcommon.BytesToHash(data)
	return &hash
}

// Retrieves the hash for the rollup containing the transaction with the given hash.
func readRollupHashByTxHash(db ethdb.KeyValueReader, txHash gethcommon.Hash) *gethcommon.Hash {
	f, err := db.Has(rollupHashTxHashKey(txHash))
	if err != nil {
		log.Panic("could not retrieve rollup hash. Cause: %s", err)
	}
	if !f {
		return nil
	}
	data, err := db.Get(rollupHashTxHashKey(txHash))
	if err != nil {
		log.Panic("could not retrieve rollup hash. Cause: %s", err)
	}
	if len(data) == 0 {
		return nil
	}
	hash := gethcommon.BytesToHash(data)
	return &hash
}
