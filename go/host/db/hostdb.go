package db

import (
	"bytes"
	"math/big"
	"os"

	gethlog "github.com/ethereum/go-ethereum/log"

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

// GetCurrentBlockHead returns the current block header (head) of the Node
func (db *DB) GetCurrentBlockHead() *types.Header {
	head := db.readHeadBlock(db.kvStore)
	if head == nil {
		return nil
	}
	return db.readBlockHeader(db.kvStore, *head)
}

// GetBlockHeader returns the block header given the Hash
func (db *DB) GetBlockHeader(hash gethcommon.Hash) *types.Header {
	return db.readBlockHeader(db.kvStore, hash)
}

// AddBlockHeader adds a types.Header to the known headers
func (db *DB) AddBlockHeader(header *types.Header) {
	b := db.kvStore.NewBatch()
	db.writeBlockHeader(b, header)

	// update the head if the new height is greater than the existing one
	currentBlockHead := db.GetCurrentBlockHead()
	if currentBlockHead == nil || currentBlockHead.Number.Int64() <= header.Number.Int64() {
		db.writeHeadBlock(b, header.Hash())
	}

	if err := b.Write(); err != nil {
		db.logger.Crit("Could not write rollup .", log.ErrKey, err)
	}
}

// GetCurrentRollupHead returns the current rollup header (head) of the Node
func (db *DB) GetCurrentRollupHead() *common.HeaderWithTxHashes {
	head := db.readHeadRollup(db.kvStore)
	if head == nil {
		return nil
	}
	return db.readRollupHeader(db.kvStore, *head)
}

// GetRollupHeader returns the rollup header given the Hash
func (db *DB) GetRollupHeader(hash gethcommon.Hash) *common.HeaderWithTxHashes {
	return db.readRollupHeader(db.kvStore, hash)
}

// AddRollupHeader adds a RollupHeader to the known headers
func (db *DB) AddRollupHeader(headerWithHashes *common.HeaderWithTxHashes) {
	b := db.kvStore.NewBatch()
	db.writeRollupHeader(b, headerWithHashes)
	db.writeRollupHash(b, headerWithHashes.Header)
	for _, txHash := range headerWithHashes.TxHashes {
		db.writeRollupNumber(b, txHash, headerWithHashes.Header.Number)
	}
	// There's a potential race here, but absolute accuracy of the number of transactions is not required.
	currentTotal := db.readTotalTransactions(db.kvStore)
	newTotal := big.NewInt(0).Add(currentTotal, big.NewInt(int64(len(headerWithHashes.TxHashes))))
	db.writeTotalTransactions(b, newTotal)

	// update the head if the new height is greater than the existing one
	currentRollupHeaderWithHashes := db.GetCurrentRollupHead()
	if currentRollupHeaderWithHashes == nil ||
		currentRollupHeaderWithHashes.Header.Number.Int64() <= headerWithHashes.Header.Number.Int64() {
		db.writeHeadRollup(b, headerWithHashes.Header.Hash())
	}

	if err := b.Write(); err != nil {
		db.logger.Crit("Could not write rollup .", log.ErrKey, err)
	}
}

func (db *DB) AddSubmittedRollup(hash gethcommon.Hash) {
	err := db.kvStore.Put(submittedRollupHeaderKey(hash), []byte{})
	if err != nil {
		db.logger.Crit("Could not save submitted rollup.", log.ErrKey, err)
	}
}

func (db *DB) WasSubmitted(hash gethcommon.Hash) bool {
	f, err := db.kvStore.Has(submittedRollupHeaderKey(hash))
	if err != nil {
		db.logger.Crit("Could not retrieve submitted rollup.", log.ErrKey, err)
	}
	return f
}

// GetRollupHash returns the hash of a rollup given its number
func (db *DB) GetRollupHash(number *big.Int) *gethcommon.Hash {
	return db.readRollupHash(db.kvStore, number)
}

// GetRollupNumber returns the number of the rollup containing the given transaction hash
func (db *DB) GetRollupNumber(txHash gethcommon.Hash) *big.Int {
	return db.readRollupNumber(db.kvStore, txHash)
}

// GetTotalTransactions returns the total number of rolled-up transactions.
func (db *DB) GetTotalTransactions() *big.Int {
	return db.readTotalTransactions(db.kvStore)
}

// schema
var (
	blockHeaderPrefix     = []byte("b")
	rollupHeaderPrefix    = []byte("r")
	headBlock             = []byte("hb")
	headRollup            = []byte("hr")
	submittedRollupPrefix = []byte("s")
	rollupHashPrefix      = []byte("rh")
	rollupNumberPrefix    = []byte("rn")
	totalTransactionsKey  = []byte("t")
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

// headerKey = rollupHashPrefix + number
func rollupHashKey(num *big.Int) []byte {
	return append(rollupHashPrefix, []byte(num.String())...)
}

// headerKey = rollupNumberPrefix + hash
func rollupNumberKey(txHash gethcommon.Hash) []byte {
	return append(rollupNumberPrefix, txHash.Bytes()...)
}

// WriteBlockHeader stores a block header into the database
func (db *DB) writeBlockHeader(w ethdb.KeyValueWriter, header *types.Header) {
	// Write the encoded header
	data, err := rlp.EncodeToBytes(header)
	if err != nil {
		db.logger.Crit("could not encode block header.", log.ErrKey, err)
	}
	key := blockHeaderKey(header.Hash())
	if err := w.Put(key, data); err != nil {
		db.logger.Crit("could not put header in DB.", log.ErrKey, err)
	}
}

// ReadBlockHeader retrieves the rollup header corresponding to the hash.
func (db *DB) readBlockHeader(r ethdb.KeyValueReader, hash gethcommon.Hash) *types.Header {
	f, err := r.Has(blockHeaderKey(hash))
	if err != nil {
		db.logger.Crit("could not retrieve block header.", log.ErrKey, err)
	}
	if !f {
		return nil
	}
	data, err := r.Get(blockHeaderKey(hash))
	if err != nil {
		db.logger.Crit("could not retrieve block header.", log.ErrKey, err)
	}
	if len(data) == 0 {
		return nil
	}
	header := new(types.Header)
	if err := rlp.Decode(bytes.NewReader(data), header); err != nil {
		db.logger.Crit("could not decode block header.", log.ErrKey, err)
	}
	return header
}

// WriteRollupHeader stores a rollup header into the database
func (db *DB) writeRollupHeader(w ethdb.KeyValueWriter, headerWithHashes *common.HeaderWithTxHashes) {
	// Write the encoded header
	data, err := rlp.EncodeToBytes(headerWithHashes)
	if err != nil {
		db.logger.Crit("could not encode rollup header.", log.ErrKey, err)
	}
	key := rollupHeaderKey(headerWithHashes.Header.Hash())
	if err := w.Put(key, data); err != nil {
		db.logger.Crit("could not put header in DB.", log.ErrKey, err)
	}
}

// ReadRollupHeader retrieves the rollup header corresponding to the hash.
func (db *DB) readRollupHeader(r ethdb.KeyValueReader, hash gethcommon.Hash) *common.HeaderWithTxHashes {
	f, err := r.Has(rollupHeaderKey(hash))
	if err != nil {
		db.logger.Crit("could not retrieve rollup header.", log.ErrKey, err)
	}
	if !f {
		return nil
	}
	data, err := r.Get(rollupHeaderKey(hash))
	if err != nil {
		db.logger.Crit("could not retrieve rollup header.", log.ErrKey, err)
	}
	if len(data) == 0 {
		return nil
	}
	header := new(common.HeaderWithTxHashes)
	if err := rlp.Decode(bytes.NewReader(data), header); err != nil {
		db.logger.Crit("could not decode rollup header.", log.ErrKey, err)
	}
	return header
}

func (db *DB) readHeadBlock(r ethdb.KeyValueReader) *gethcommon.Hash {
	f, err := r.Has(headBlock)
	if err != nil {
		db.logger.Crit("could not retrieve head block.", log.ErrKey, err)
	}
	if !f {
		return nil
	}
	value, err := r.Get(headBlock)
	if err != nil {
		db.logger.Crit("could not retrieve head block.", log.ErrKey, err)
	}
	h := gethcommon.BytesToHash(value)
	return &h
}

func (db *DB) writeHeadBlock(w ethdb.KeyValueWriter, val gethcommon.Hash) {
	err := w.Put(headBlock, val.Bytes())
	if err != nil {
		db.logger.Crit("could not write head block.", log.ErrKey, err)
	}
}

func (db *DB) readHeadRollup(r ethdb.KeyValueReader) *gethcommon.Hash {
	f, err := r.Has(headRollup)
	if err != nil {
		db.logger.Crit("could not retrieve head rollup.", log.ErrKey, err)
	}
	if !f {
		return nil
	}
	value, err := r.Get(headRollup)
	if err != nil {
		db.logger.Crit("could not retrieve head rollup.", log.ErrKey, err)
	}
	h := gethcommon.BytesToHash(value)
	return &h
}

func (db *DB) writeHeadRollup(w ethdb.KeyValueWriter, val gethcommon.Hash) {
	err := w.Put(headRollup, val.Bytes())
	if err != nil {
		db.logger.Crit("could not write head rollup.", log.ErrKey, err)
	}
}

// Stores the hash of a rollup into the database, keyed by the rollup's number
func (db *DB) writeRollupHash(w ethdb.KeyValueWriter, header *common.Header) {
	key := rollupHashKey(header.Number)
	if err := w.Put(key, header.Hash().Bytes()); err != nil {
		db.logger.Crit("could not put header in DB.", log.ErrKey, err)
	}
}

// Stores the hash of a rollup into the database, keyed by the hashes of the transactions in the rollup
func (db *DB) writeRollupNumber(w ethdb.KeyValueWriter, txHash gethcommon.Hash, rollupNumber *big.Int) {
	key := rollupNumberKey(txHash)
	// TODO - Investigate this off-by-one issue. The tx hashes that are in the `BlockSubmissionResponse` for rollup #1
	//  are actually the transactions for rollup #2.
	number := big.NewInt(0).Add(rollupNumber, big.NewInt(1))
	if err := w.Put(key, number.Bytes()); err != nil {
		db.logger.Crit("could not put rollup number in DB.", log.ErrKey, err)
	}
}

func (db *DB) writeTotalTransactions(w ethdb.KeyValueWriter, newTotal *big.Int) {
	err := w.Put(totalTransactionsKey, newTotal.Bytes())
	if err != nil {
		db.logger.Crit("Could not save total transactions.", log.ErrKey, err)
	}
}

// Retrieves the hash for the rollup with the given number.
func (db *DB) readRollupHash(r ethdb.KeyValueReader, number *big.Int) *gethcommon.Hash {
	f, err := r.Has(rollupHashKey(number))
	if err != nil {
		db.logger.Crit("could not retrieve rollup hash.", log.ErrKey, err)
	}
	if !f {
		return nil
	}
	data, err := r.Get(rollupHashKey(number))
	if err != nil {
		db.logger.Crit("could not retrieve rollup hash.", log.ErrKey, err)
	}
	if len(data) == 0 {
		return nil
	}
	hash := gethcommon.BytesToHash(data)
	return &hash
}

// Retrieves the number of the rollup containing the transaction with the given hash.
func (db *DB) readRollupNumber(r ethdb.KeyValueReader, txHash gethcommon.Hash) *big.Int {
	f, err := r.Has(rollupNumberKey(txHash))
	if err != nil {
		db.logger.Crit("could not retrieve rollup number.", log.ErrKey, err)
	}
	if !f {
		return nil
	}
	data, err := r.Get(rollupNumberKey(txHash))
	if err != nil {
		db.logger.Crit("could not retrieve rollup number.", log.ErrKey, err)
	}
	if len(data) == 0 {
		return nil
	}
	return big.NewInt(0).SetBytes(data)
}

// Retrieves the total number of rolled-up transactions.
func (db *DB) readTotalTransactions(r ethdb.KeyValueReader) *big.Int {
	f, err := r.Has(totalTransactionsKey)
	if err != nil {
		db.logger.Crit("could not retrieve total transactions.", log.ErrKey, err)
	}
	if !f {
		return big.NewInt(0)
	}
	data, err := r.Get(totalTransactionsKey)
	if err != nil {
		db.logger.Crit("could not retrieve total transactions.", log.ErrKey, err)
	}
	if len(data) == 0 {
		return big.NewInt(0)
	}
	return big.NewInt(0).SetBytes(data)
}
