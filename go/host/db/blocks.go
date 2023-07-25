package db

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/go-obscuro/go/common/errutil"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

// DB methods relating to blocks.

// GetBlockHeader returns the block header given the hash.
func (db *DB) GetBlockHeader(hash gethcommon.Hash) (*types.Header, error) {
	return db.readBlockHeader(db.kvStore, hash)
}

// AddBlockHeader adds a types.Header to the known headers
func (db *DB) AddBlockHeader(header *types.Header) error {
	b := db.kvStore.NewBatch()
	err := db.writeBlockHeader(header)
	if err != nil {
		return fmt.Errorf("could not write block header. Cause: %w", err)
	}

	// Update the head if the new height is greater than the existing one.
	tipBlockHeader, err := db.GetTipBlockHeader()
	if err != nil && !errors.Is(err, errutil.ErrNotFound) {
		return fmt.Errorf("could not retrieve the tip block header. Cause: %w", err)
	}
	if tipBlockHeader == nil || tipBlockHeader.Number.Cmp(header.Number) == -1 {
		err = db.writeTipBlockHash(b, header.Hash())
		if err != nil {
			return fmt.Errorf("could not write new head block hash. Cause: %w", err)
		}
	}

	if err = b.Write(); err != nil {
		return fmt.Errorf("could not write batch to DB. Cause: %w", err)
	}

	return nil
}

// GetTipBlockHash returns the hash of the node's current head block.
func (db *DB) GetTipBlockHash() (*gethcommon.Hash, error) {
	return db.readTipBlockHash()
}

// GetTipBlockHeader returns the header of the node's current head block.
func (db *DB) GetTipBlockHeader() (*types.Header, error) {
	headBatchHash, err := db.GetTipBlockHash()
	if err != nil {
		return nil, err
	}
	return db.readBlockHeader(db.kvStore, *headBatchHash)
}

// headerKey = blockHeaderPrefix  + hash
func blockHeaderKey(hash gethcommon.Hash) []byte {
	return append(blockHeaderPrefix, hash.Bytes()...)
}

// Stores a block header into the database
func (db *DB) writeBlockHeader(header *types.Header) error {
	// Write the encoded header
	data, err := rlp.EncodeToBytes(header)
	if err != nil {
		return err
	}
	key := blockHeaderKey(header.Hash())
	if err := db.kvStore.Put(key, data); err != nil {
		return err
	}
	db.blockWrites.Inc(1)
	return nil
}

// Retrieves the block header corresponding to the hash.
func (db *DB) readBlockHeader(r ethdb.KeyValueReader, hash gethcommon.Hash) (*types.Header, error) {
	data, err := r.Get(blockHeaderKey(hash))
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errutil.ErrNotFound
	}
	header := new(types.Header)
	if err := rlp.Decode(bytes.NewReader(data), header); err != nil {
		return nil, err
	}
	db.blockReads.Inc(1)
	return header, nil
}

// Stores the tip batch header hash into the database.
func (db *DB) writeTipBlockHash(w ethdb.KeyValueWriter, hash gethcommon.Hash) error {
	return w.Put(tipBlockHash, hash.Bytes())
}

// Retrieves the hash of the head block.
func (db *DB) readTipBlockHash() (*gethcommon.Hash, error) {
	value, err := db.kvStore.Get(tipBlockHash)
	if err != nil {
		return nil, err
	}
	h := gethcommon.BytesToHash(value)
	return &h, nil
}
