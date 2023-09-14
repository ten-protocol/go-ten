package db

import (
	"bytes"
	"fmt"

	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/errutil"
	"github.com/pkg/errors"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

// DB methods relating to rollup transactions.

// AddRollupHeader adds a rollup to the DB
func (db *DB) AddRollupHeader(rollup *common.ExtRollup, block *common.L1Block) error {
	// Check if the Header is already stored
	_, err := db.GetRollupHeader(rollup.Hash())
	if err != nil && !errors.Is(err, errutil.ErrNotFound) {
		return fmt.Errorf("could not retrieve rollup header. Cause: %w", err)
	}
	if err == nil {
		// The rollup is already stored, so we return early.
		return errutil.ErrAlreadyExists
	}

	b := db.kvStore.NewBatch()

	if err := db.writeRollupHeader(b, rollup.Header); err != nil {
		return fmt.Errorf("could not write rollup header. Cause: %w", err)
	}

	if err := db.writeRollupByBlockHash(b, rollup.Header, block.Hash()); err != nil {
		return fmt.Errorf("could not write rollup block. Cause: %w", err)
	}

	// Update the tip if the new height is greater than the existing one.
	tipRollupHeader, err := db.GetTipRollupHeader()
	if err != nil && !errors.Is(err, errutil.ErrNotFound) {
		return fmt.Errorf("could not retrieve rollup header at tip. Cause: %w", err)
	}
	if tipRollupHeader == nil || tipRollupHeader.LastBatchSeqNo < rollup.Header.LastBatchSeqNo {
		err = db.writeTipRollupHeader(b, rollup.Hash())
		if err != nil {
			return fmt.Errorf("could not write new rollup hash at tip. Cause: %w", err)
		}
	}

	if err = b.Write(); err != nil {
		return fmt.Errorf("could not write batch to DB. Cause: %w", err)
	}
	return nil
}

// GetRollupHeader returns the rollup with the given hash.
func (db *DB) GetRollupHeader(hash gethcommon.Hash) (*common.RollupHeader, error) {
	return db.readRollupHeader(rollupHashKey(hash))
}

// GetTipRollupHeader returns the header of the node's current tip rollup.
func (db *DB) GetTipRollupHeader() (*common.RollupHeader, error) {
	headBatchHash, err := db.readTipRollupHash()
	if err != nil {
		return nil, err
	}
	return db.readRollupHeader(rollupHashKey(*headBatchHash))
}

// GetRollupHeaderByBlock returns the rollup for the given block
func (db *DB) GetRollupHeaderByBlock(blockHash gethcommon.Hash) (*common.RollupHeader, error) {
	return db.readRollupHeader(rollupBlockKey(blockHash))
}

// Retrieves the rollup corresponding to the hash.
func (db *DB) readRollupHeader(key []byte) (*common.RollupHeader, error) {
	data, err := db.kvStore.Get(key)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errutil.ErrNotFound
	}
	rollupHeader := new(common.RollupHeader)
	if err := rlp.Decode(bytes.NewReader(data), rollupHeader); err != nil {
		return nil, err
	}
	return rollupHeader, nil
}

// Stores a rollup header into the database.
func (db *DB) writeRollupHeader(w ethdb.KeyValueWriter, header *common.RollupHeader) error {
	// Write the encoded header
	data, err := rlp.EncodeToBytes(header)
	if err != nil {
		return err
	}
	key := rollupHashKey(header.Hash())

	return w.Put(key, data)
}

// Retrieves the hash of the rollup at tip
func (db *DB) readTipRollupHash() (*gethcommon.Hash, error) {
	value, err := db.kvStore.Get(tipRollupHash)
	if err != nil {
		return nil, err
	}
	h := gethcommon.BytesToHash(value)
	return &h, nil
}

// Stores the tip rollup header hash into the database
func (db *DB) writeTipRollupHeader(w ethdb.KeyValueWriter, val gethcommon.Hash) error {
	err := w.Put(tipRollupHash, val.Bytes())
	if err != nil {
		return err
	}
	return nil
}

// Stores if a rollup is in a block
func (db *DB) writeRollupByBlockHash(w ethdb.KeyValueWriter, header *common.RollupHeader, blockHash gethcommon.Hash) error {
	// Write the encoded header
	data, err := rlp.EncodeToBytes(header)
	if err != nil {
		return err
	}
	key := rollupBlockKey(blockHash)

	return w.Put(key, data)
}

// rollupHashKey = rollupHeaderPrefix  + hash
func rollupHashKey(hash gethcommon.Hash) []byte {
	return append(rollupHeaderPrefix, hash.Bytes()...)
}

// rollupBlockKey = rollupHeaderBlockPrefix  + hash
func rollupBlockKey(hash gethcommon.Hash) []byte {
	return append(rollupHeaderBlockPrefix, hash.Bytes()...)
}
