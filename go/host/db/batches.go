package db

import (
	"bytes"
	"math/big"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/log"
)

// DB methods relating to batches.

// headerKey = batchHeaderPrefix  + hash
func batchHeaderKey(hash gethcommon.Hash) []byte {
	return append(batchHeaderPrefix, hash.Bytes()...)
}

// GetHeadBatchHeader returns the header of the node's current head batch, or (nil, false) if no such header is found.
func (db *DB) GetHeadBatchHeader() (*common.Header, bool) {
	headBatchHash, found := db.readHeadBatchHash(db.kvStore)
	if !found {
		return nil, false
	}
	return db.readBatchHeader(db.kvStore, *headBatchHash)
}

// AddBatchHeader adds a batch's header to the known headers
func (db *DB) AddBatchHeader(header *common.Header, txHashes []common.TxHash) {
	b := db.kvStore.NewBatch()
	db.writeBatchHeader(b, header)

	// TODO - #718 - Store the batch txs, batch hash, and batch number per transaction hash, if needed (see `AddRollupHeader`).

	// There's a potential race here, but absolute accuracy of the number of transactions is not required.
	currentTotal := db.readTotalTransactions(db.kvStore)
	newTotal := big.NewInt(0).Add(currentTotal, big.NewInt(int64(len(txHashes))))
	db.writeTotalTransactions(b, newTotal)

	// update the head if the new height is greater than the existing one
	headBatchHeader, found := db.GetHeadBatchHeader()
	if !found || headBatchHeader.Number.Int64() <= header.Number.Int64() {
		db.writeHeadBatchHash(b, header.Hash())
	}

	if err := b.Write(); err != nil {
		db.logger.Crit("Could not write batch.", log.ErrKey, err)
	}
}

// Retrieves the batch header corresponding to the hash, or (nil, false) if no such header is found.
func (db *DB) readBatchHeader(r ethdb.KeyValueReader, hash gethcommon.Hash) (*common.Header, bool) {
	f, err := r.Has(batchHeaderKey(hash))
	if err != nil {
		db.logger.Crit("could not retrieve batch header.", log.ErrKey, err)
	}
	if !f {
		return nil, false
	}
	data, err := r.Get(batchHeaderKey(hash))
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
func (db *DB) readHeadBatchHash(r ethdb.KeyValueReader) (*gethcommon.Hash, bool) {
	f, err := r.Has(headBatch)
	if err != nil {
		db.logger.Crit("could not retrieve head batch.", log.ErrKey, err)
	}
	if !f {
		return nil, false
	}
	value, err := r.Get(headBatch)
	if err != nil {
		db.logger.Crit("could not retrieve head batch.", log.ErrKey, err)
	}
	h := gethcommon.BytesToHash(value)
	return &h, true
}

// Stores a batch header into the database.
func (db *DB) writeBatchHeader(w ethdb.KeyValueWriter, header *common.Header) {
	// Write the encoded header
	data, err := rlp.EncodeToBytes(header)
	if err != nil {
		db.logger.Crit("could not encode batch header.", log.ErrKey, err)
	}
	key := batchHeaderKey(header.Hash())
	if err := w.Put(key, data); err != nil {
		db.logger.Crit("could not put batch header in DB.", log.ErrKey, err)
	}
}

// Stores the head batch header hash into the database.
func (db *DB) writeHeadBatchHash(w ethdb.KeyValueWriter, val gethcommon.Hash) {
	err := w.Put(headBatch, val.Bytes())
	if err != nil {
		db.logger.Crit("could not put head batch hash in DB.", log.ErrKey, err)
	}
}
