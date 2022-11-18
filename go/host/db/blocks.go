package db

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/obscuronet/go-obscuro/go/common/errutil"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
)

// DB methods relating to batches.

// GetHeadBlockHeader returns the block header of the current head block.
func (db *DB) GetHeadBlockHeader() (*types.Header, error) {
	head, err := db.readHeadBlock(db.kvStore)
	if err != nil {
		return nil, err
	}
	return db.readBlockHeader(db.kvStore, *head)
}

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

	// update the head if the new height is greater than the existing one
	headBlockHeader, err := db.GetHeadBlockHeader()
	if err != nil && !errors.Is(err, errutil.ErrNotFound) {
		return fmt.Errorf("could not retrieve head block header. Cause: %w", err)
	}
	if errors.Is(err, errutil.ErrNotFound) || headBlockHeader.Number.Int64() <= header.Number.Int64() {
		err = db.writeHeadBlockHash(header.Hash())
		if err != nil {
			return fmt.Errorf("could not write new head block hash. Cause: %w", err)
		}
	}

	if err = b.Write(); err != nil {
		return fmt.Errorf("could not write batch to DB. Cause: %w", err)
	}

	return nil
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
	return nil
}

func (db *DB) writeHeadBlockHash(val gethcommon.Hash) error {
	err := db.kvStore.Put(headBlock, val.Bytes())
	if err != nil {
		return err
	}
	return nil
}

// Retrieves the block header corresponding to the hash.
func (db *DB) readBlockHeader(r ethdb.KeyValueReader, hash gethcommon.Hash) (*types.Header, error) {
	f, err := r.Has(blockHeaderKey(hash))
	if err != nil {
		return nil, err
	}
	if !f {
		return nil, errutil.ErrNotFound
	}
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
	return header, nil
}

func (db *DB) readHeadBlock(r ethdb.KeyValueReader) (*gethcommon.Hash, error) {
	f, err := r.Has(headBlock)
	if err != nil {
		return nil, err
	}
	if !f {
		return nil, errutil.ErrNotFound
	}
	value, err := r.Get(headBlock)
	if err != nil {
		return nil, err
	}
	h := gethcommon.BytesToHash(value)
	return &h, nil
}
