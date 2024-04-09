package hostdb

import (
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ten-protocol/go-ten/go/common"
)

// AddBlock stores a block header with the given rollupHash it contains in the host DB
func AddBlock(dbtx *dbTransaction, b *types.Header, rollupHash common.L2RollupHash) error {
	header, err := rlp.EncodeToBytes(b)
	if err != nil {
		return fmt.Errorf("could not encode block header. Cause: %w", err)
	}

	r, err := rlp.EncodeToBytes(rollupHash)
	if err != nil {
		return fmt.Errorf("could not encode rollup hash transactions: %w", err)
	}

	_, err = dbtx.GetDB().Exec(dbtx.GetSQLStatements().InsertBlock,
		b.Hash(), // hash
		header,   // l1 block header
		r,        // rollup hash
	)
	if err != nil {
		return fmt.Errorf("could not insert block. Cause: %w", err)
	}

	return nil
}

// GetBlockListing returns a paginated list of blocks in descending order against the order they were added
func GetBlockListing(dbtx *dbTransaction, pagination *common.QueryPagination) (*common.BlockListingResponse, error) {
	rows, err := dbtx.GetDB().Query(dbtx.GetSQLStatements().SelectBlocks, pagination.Size, pagination.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var blocks []common.PublicBlock

	for rows.Next() {
		var id int
		var hash, header, rollupHash []byte

		err = rows.Scan(&id, &hash, &header, &rollupHash)
		if err != nil {
			return nil, err
		}

		blockHeader := new(types.Header)
		if err := rlp.DecodeBytes(header, blockHeader); err != nil {
			return nil, fmt.Errorf("could not decode block header. Cause: %w", err)
		}
		r := new(common.L2RollupHash)
		if err := rlp.DecodeBytes(rollupHash, r); err != nil {
			return nil, fmt.Errorf("could not decode rollup hash. Cause: %w", err)
		}
		block := common.PublicBlock{
			BlockHeader: *blockHeader,
			RollupHash:  *r,
		}
		blocks = append(blocks, block)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &common.BlockListingResponse{
		BlocksData: blocks,
		Total:      uint64(len(blocks)),
	}, nil
}
