package hostdb

import (
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ten-protocol/go-ten/go/common"
)

const (
	selectBlocks = "SELECT b.id, b.hash, b.header, r.hash FROM block_host b join rollup_host r on r.compression_block=b.id ORDER BY id DESC "
)

// AddBlock stores a block header with the given rollupHash it contains in the host DB
func AddBlock(dbtx *dbTransaction, statements *SQLStatements, b *types.Header) error {
	header, err := rlp.EncodeToBytes(b)
	if err != nil {
		return fmt.Errorf("could not encode block header. Cause: %w", err)
	}

	_, err = dbtx.tx.Exec(statements.InsertBlock,
		b.Hash().Bytes(), // hash
		header,           // l1 block header
	)
	if err != nil {
		return fmt.Errorf("could not insert block. Cause: %w", err)
	}

	return nil
}

// GetBlockListing returns a paginated list of blocks in descending order against the order they were added
func GetBlockListing(db HostDB, pagination *common.QueryPagination) (*common.BlockListingResponse, error) {
	query := selectBlocks + db.GetSQLStatement().Pagination
	rows, err := db.GetSQLDB().Query(query, pagination.Size, pagination.Offset)
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
