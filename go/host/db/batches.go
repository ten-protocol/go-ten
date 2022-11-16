package db

import (
	"bytes"
	"math/big"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/log"
)

// DB methods relating to batches.

// GetHeadBatchHeader returns the header of the node's current head batch, or (nil, false) if no such header is found.
func (db *DB) GetHeadBatchHeader() (*common.Header, bool) {
	headBatchHash, found := db.readHeadBatchHash()
	if !found {
		return nil, false
	}
	return db.readBatchHeader(*headBatchHash)
}

// AddBatchHeader adds a batch's header to the known headers
func (db *DB) AddBatchHeader(header *common.Header, txHashes []common.TxHash) {
	b := db.kvStore.NewBatch()
	db.writeBatchHeader(header)

	// TODO - #718 - Store the batch txs, batch hash, and batch number per transaction hash, if needed (see `AddRollupHeader`).

	// There's a potential race here, but absolute accuracy of the number of transactions is not required.
	currentTotal := db.readTotalTransactions()
	newTotal := big.NewInt(0).Add(currentTotal, big.NewInt(int64(len(txHashes))))
	db.writeTotalTransactions(b, newTotal)

	// update the head if the new height is greater than the existing one
	headBatchHeader, found := db.GetHeadBatchHeader()
	if !found || headBatchHeader.Number.Int64() <= header.Number.Int64() {
		db.writeHeadBatchHash(header.Hash())
	}

	if err := b.Write(); err != nil {
		db.logger.Crit("Could not write batch.", log.ErrKey, err)
	}
}

// headerKey = batchHeaderPrefix  + hash
func batchHeaderKey(hash gethcommon.Hash) []byte {
	return append(batchHeaderPrefix, hash.Bytes()...)
}

// Retrieves the batch header corresponding to the hash, or (nil, false) if no such header is found.
func (db *DB) readBatchHeader(hash gethcommon.Hash) (*common.Header, bool) {
	f, err := db.kvStore.Has(batchHeaderKey(hash))
	if err != nil {
		db.logger.Crit("could not retrieve batch header.", log.ErrKey, err)
	}
	if !f {
		return nil, false
	}
	data, err := db.kvStore.Get(batchHeaderKey(hash))
	if err != nil {
		db.logger.Crit("could not retrieve batch header.", log.ErrKey, err)
	}
	if len(data) == 0 {
		return nil, false
	}
	header := new(common.Header)
	if err := rlp.Decode(bytes.NewReader(data), header); err != nil {
		db.logger.Crit("could not decode batch header.", log.ErrKey, err)
	}
	return header, true
}

// Retrieves the hash of the head batch.
func (db *DB) readHeadBatchHash() (*gethcommon.Hash, bool) {
	f, err := db.kvStore.Has(headBatch)
	if err != nil {
		db.logger.Crit("could not retrieve head batch.", log.ErrKey, err)
	}
	if !f {
		return nil, false
	}
	value, err := db.kvStore.Get(headBatch)
	if err != nil {
		db.logger.Crit("could not retrieve head batch.", log.ErrKey, err)
	}
	h := gethcommon.BytesToHash(value)
	return &h, true
}

// Stores a batch header into the database.
func (db *DB) writeBatchHeader(header *common.Header) {
	// Write the encoded header
	data, err := rlp.EncodeToBytes(header)
	if err != nil {
		db.logger.Crit("could not encode batch header.", log.ErrKey, err)
	}
	key := batchHeaderKey(header.Hash())
	if err := db.kvStore.Put(key, data); err != nil {
		db.logger.Crit("could not put batch header in DB.", log.ErrKey, err)
	}
}

// Stores the head batch header hash into the database.
func (db *DB) writeHeadBatchHash(val gethcommon.Hash) {
	err := db.kvStore.Put(headBatch, val.Bytes())
	if err != nil {
		db.logger.Crit("could not put head batch hash in DB.", log.ErrKey, err)
	}
}
