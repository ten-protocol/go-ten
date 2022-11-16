package db

import (
	"bytes"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/go-obscuro/go/common/log"
)

// DB methods relating to batches.

// headerKey = blockHeaderPrefix  + hash
func blockHeaderKey(hash gethcommon.Hash) []byte {
	return append(blockHeaderPrefix, hash.Bytes()...)
}

// GetHeadBlockHeader returns the block header of the current head block, or (nil, false) if no such header is found.
func (db *DB) GetHeadBlockHeader() (*types.Header, bool) {
	head := db.readHeadBlock(db.kvStore)
	if head == nil {
		return nil, false
	}
	return db.readBlockHeader(db.kvStore, *head)
}

// GetBlockHeader returns the block header given the hash, or (nil, false) if no such header is found.
func (db *DB) GetBlockHeader(hash gethcommon.Hash) (*types.Header, bool) {
	return db.readBlockHeader(db.kvStore, hash)
}

// AddBlockHeader adds a types.Header to the known headers
func (db *DB) AddBlockHeader(header *types.Header) {
	b := db.kvStore.NewBatch()
	db.writeBlockHeader(b, header)

	// update the head if the new height is greater than the existing one
	headBlockHeader, found := db.GetHeadBlockHeader()
	if !found || headBlockHeader.Number.Int64() <= header.Number.Int64() {
		db.writeHeadBlockHash(b, header.Hash())
	}

	if err := b.Write(); err != nil {
		db.logger.Crit("Could not write rollup.", log.ErrKey, err)
	}
}

// Stores a block header into the database
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

func (db *DB) writeHeadBlockHash(w ethdb.KeyValueWriter, val gethcommon.Hash) {
	err := w.Put(headBlock, val.Bytes())
	if err != nil {
		db.logger.Crit("could not write head block.", log.ErrKey, err)
	}
}

// Retrieves the block header corresponding to the hash, or (nil, false) if no such header is found.
func (db *DB) readBlockHeader(r ethdb.KeyValueReader, hash gethcommon.Hash) (*types.Header, bool) {
	f, err := r.Has(blockHeaderKey(hash))
	if err != nil {
		db.logger.Crit("could not retrieve block header.", log.ErrKey, err)
	}
	if !f {
		return nil, false
	}
	data, err := r.Get(blockHeaderKey(hash))
	if err != nil {
		db.logger.Crit("could not retrieve block header.", log.ErrKey, err)
	}
	if len(data) == 0 {
		return nil, false
	}
	header := new(types.Header)
	if err := rlp.Decode(bytes.NewReader(data), header); err != nil {
		db.logger.Crit("could not decode block header.", log.ErrKey, err)
	}
	return header, true
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
