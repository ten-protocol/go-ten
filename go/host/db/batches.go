package db

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/obscuronet/go-obscuro/go/common/errutil"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/go-obscuro/go/common"
)

// DB methods relating to batches.

// GetHeadBatchHeader returns the header of the node's current head batch.
func (db *DB) GetHeadBatchHeader() (*common.Header, error) {
	headBatchHash, err := db.readHeadBatchHash()
	if err != nil {
		return nil, err
	}
	return db.readBatchHeader(*headBatchHash)
}

// AddBatchHeader adds a batch's header to the known headers
func (db *DB) AddBatchHeader(header *common.Header, txHashes []common.TxHash) error {
	b := db.kvStore.NewBatch()
	err := db.writeBatchHeader(header)
	if err != nil {
		return fmt.Errorf("could not write batch header. Cause: %w", err)
	}

	// TODO - #718 - Store the batch txs, batch hash, and batch number per transaction hash, if needed (see `AddRollupHeader`).

	// TODO - #718 - Update the total transactions, once we no longer do this in `AddRollupHeader`.

	// update the head if the new height is greater than the existing one
	headBatchHeader, err := db.GetHeadBatchHeader()
	if err != nil && !errors.Is(err, errutil.ErrNotFound) {
		return fmt.Errorf("could not retrieve head batch header. Cause: %w", err)
	}
	if errors.Is(err, errutil.ErrNotFound) || headBatchHeader.Number.Int64() <= header.Number.Int64() {
		err = db.writeHeadBatchHash(header.Hash())
		if err != nil {
			return fmt.Errorf("could not write new head batch hash. Cause: %w", err)
		}
	}

	if err = b.Write(); err != nil {
		return fmt.Errorf("could not write batch to DB. Cause: %w", err)
	}

	return nil
}

// headerKey = batchHeaderPrefix  + hash
func batchHeaderKey(hash gethcommon.Hash) []byte {
	return append(batchHeaderPrefix, hash.Bytes()...)
}

// Retrieves the batch header corresponding to the hash.
func (db *DB) readBatchHeader(hash gethcommon.Hash) (*common.Header, error) {
	// TODO - #1208 - Analyse this weird Has/Get pattern, here and in other part of the `db` package.
	f, err := db.kvStore.Has(batchHeaderKey(hash))
	if err != nil {
		return nil, err
	}
	if !f {
		return nil, errutil.ErrNotFound
	}
	data, err := db.kvStore.Get(batchHeaderKey(hash))
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errutil.ErrNotFound
	}
	header := new(common.Header)
	if err := rlp.Decode(bytes.NewReader(data), header); err != nil {
		return nil, err
	}
	return header, nil
}

// Retrieves the hash of the head batch.
func (db *DB) readHeadBatchHash() (*gethcommon.Hash, error) {
	f, err := db.kvStore.Has(headBatch)
	if err != nil {
		return nil, err
	}
	if !f {
		return nil, errutil.ErrNotFound
	}
	value, err := db.kvStore.Get(headBatch)
	if err != nil {
		return nil, err
	}
	h := gethcommon.BytesToHash(value)
	return &h, nil
}

// Stores a batch header into the database.
func (db *DB) writeBatchHeader(header *common.Header) error {
	// Write the encoded header
	data, err := rlp.EncodeToBytes(header)
	if err != nil {
		return err
	}
	key := batchHeaderKey(header.Hash())
	if err := db.kvStore.Put(key, data); err != nil {
		return err
	}
	return nil
}

// Stores the head batch header hash into the database.
func (db *DB) writeHeadBatchHash(val gethcommon.Hash) error {
	err := db.kvStore.Put(headBatch, val.Bytes())
	if err != nil {
		return err
	}
	return nil
}
