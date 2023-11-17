package db

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/errutil"
)

// DB methods relating to block headers.

// GetBlockByHash returns the block header given the hash.
func (db *DB) GetBlockByHash(hash gethcommon.Hash) (*types.Header, error) {
	return db.readBlock(db.kvStore, blockHashKey(hash))
}

// GetBlockByHeight returns the block header given the height
func (db *DB) GetBlockByHeight(height *big.Int) (*types.Header, error) {
	return db.readBlock(db.kvStore, blockNumberKey(height))
}

// AddBlock adds a types.Header to the known headers
func (db *DB) AddBlock(header *types.Header) error {
	b := db.kvStore.NewBatch()
	err := db.writeBlockByHash(header)
	if err != nil {
		return fmt.Errorf("could not write block header. Cause: %w", err)
	}

	err = db.writeBlockByHeight(header)
	if err != nil {
		return fmt.Errorf("could not write block header. Cause: %w", err)
	}

	// Update the tip if the new height is greater than the existing one.
	tipBlockHeader, err := db.GetBlockAtTip()
	if err != nil && !errors.Is(err, errutil.ErrNotFound) {
		return fmt.Errorf("could not retrieve block header at tip. Cause: %w", err)
	}
	if tipBlockHeader == nil || tipBlockHeader.Number.Cmp(header.Number) == -1 {
		err = db.writeBlockAtTip(b, header.Hash())
		if err != nil {
			return fmt.Errorf("could not write new block hash at tip. Cause: %w", err)
		}
	}

	if err = b.Write(); err != nil {
		return fmt.Errorf("could not write batch to DB. Cause: %w", err)
	}

	return nil
}

// GetBlockListing returns latest L1 blocks given the pagination.
// For example, page 0, size 10 will return the latest 10 blocks.
func (db *DB) GetBlockListing(pagination *common.QueryPagination) (*common.BlockListingResponse, error) {
	// fetch the total blocks so we can paginate
	tipHeader, err := db.GetBlockAtTip()
	if err != nil {
		return nil, err
	}

	blocksFrom := tipHeader.Number.Uint64() - pagination.Offset
	blocksToInclusive := int(blocksFrom) - int(pagination.Size) + 1
	// if blocksToInclusive would be negative, set it to 0
	if blocksToInclusive < 0 {
		blocksToInclusive = 0
	}

	// fetch requested batches
	var blocks []common.PublicBlock
	for i := blocksFrom; i > uint64(blocksToInclusive); i-- {
		header, err := db.GetBlockByHeight(big.NewInt(int64(i)))
		if err != nil {
			return nil, err
		}

		// check if the block has a rollup
		rollup, err := db.GetRollupHeaderByBlock(header.Hash())
		if err != nil && !errors.Is(err, errutil.ErrNotFound) {
			return nil, err
		}

		listedBlock := common.PublicBlock{BlockHeader: *header}
		if rollup != nil {
			listedBlock.RollupHash = rollup.Hash()
			fmt.Println("added at block: ", header.Number.Int64(), " - ", listedBlock.RollupHash)
		}
		blocks = append(blocks, listedBlock)
	}

	return &common.BlockListingResponse{
		BlocksData: blocks,
		Total:      tipHeader.Number.Uint64(),
	}, nil
}

// GetBlockAtTip returns the block at current Head or Tip
func (db *DB) GetBlockAtTip() (*types.Header, error) {
	value, err := db.kvStore.Get(blockHeadedAtTip)
	if err != nil {
		return nil, err
	}
	h := gethcommon.BytesToHash(value)

	return db.GetBlockByHash(h)
}

// Stores the hash of the block at tip
func (db *DB) writeBlockAtTip(w ethdb.KeyValueWriter, hash gethcommon.Hash) error {
	err := w.Put(blockHeadedAtTip, hash.Bytes())
	if err != nil {
		return err
	}
	return nil
}

// Stores a block header into the database using the hash as key
func (db *DB) writeBlockByHash(header *types.Header) error {
	// Write the encoded header
	data, err := rlp.EncodeToBytes(header)
	if err != nil {
		return err
	}
	key := blockHashKey(header.Hash())
	if err := db.kvStore.Put(key, data); err != nil {
		return err
	}
	db.blockWrites.Inc(1)
	return nil
}

// Stores a block header into the database using the height as key
func (db *DB) writeBlockByHeight(header *types.Header) error {
	// Write the encoded header
	data, err := rlp.EncodeToBytes(header)
	if err != nil {
		return err
	}
	key := blockNumberKey(header.Number)
	return db.kvStore.Put(key, data)
}

// Retrieves the block header corresponding to the key.
func (db *DB) readBlock(r ethdb.KeyValueReader, key []byte) (*types.Header, error) {
	data, err := r.Get(key)
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

// headerKey = blockNumberHeaderPrefix  + hash
func blockNumberKey(height *big.Int) []byte {
	return append(blockNumberHeaderPrefix, height.Bytes()...)
}

// headerKey = blockHeaderPrefix  + hash
func blockHashKey(hash gethcommon.Hash) []byte {
	return append(blockHeaderPrefix, hash.Bytes()...)
}
