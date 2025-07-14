package hostdb

import (
	"database/sql"
	"encoding/json"
	"fmt"

	gethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ten-protocol/go-ten/go/common"
)

const (
	selectBlocks  = "SELECT b.id, b.hash, b.header, r.hash FROM block_host b join rollup_host r on r.compression_block=b.id ORDER BY b.id DESC "
	selectBlockId = "SELECT id FROM block_host WHERE hash = "
	selectBlock   = "SELECT header FROM block_host WHERE hash = "
)

// AddBlock stores a block header with the given rollupHash it contains in the host DB
func AddBlock(dbtx *sql.Tx, statements *SQLStatements, b *types.Header) error {
	header, err := json.Marshal(b)
	if err != nil {
		return fmt.Errorf("could not encode block header. Cause: %w", err)
	}

	_, err = dbtx.Exec(statements.InsertBlock,
		b.Hash().Bytes(), // hash
		header,           // l1 block header
	)
	if err != nil {
		return fmt.Errorf("could not insert block. Cause: %w", err)
	}

	return nil
}

// GetBlockId returns the block ID given the hash.
func GetBlockId(db *sql.Tx, statements *SQLStatements, hash gethcommon.Hash) (*int64, error) {
	query := selectBlockId + statements.Placeholder
	var blockId int64
	err := db.QueryRow(query, hash.Bytes()).Scan(&blockId)
	if err != nil {
		return nil, fmt.Errorf("query execution for select block failed: %w", err)
	}

	return &blockId, nil
}

// GetBlock returns the block ID given the hash.
func GetBlock(db HostDB, statements *SQLStatements, hash *gethcommon.Hash) (*types.Header, error) {
	query := selectBlock + statements.Placeholder
	var header []byte
	err := db.GetSQLDB().QueryRow(query, hash.Bytes()).Scan(&header)
	if err != nil {
		return nil, fmt.Errorf("query execution for select block failed: %w", err)
	}
	h := new(types.Header)
	if err := json.Unmarshal(header, h); err != nil {
		return nil, fmt.Errorf("could not decode block header. Cause: %w", err)
	}
	return h, nil
}

// GetBlockListing returns a paginated list of blocks in descending order against the order they were added
func GetBlockListing(db HostDB, pagination *common.QueryPagination) (*common.BlockListingResponse, error) {
	query := selectBlocks + db.GetSQLStatement().Pagination
	rows, err := db.GetSQLDB().Query(query, int64(pagination.Size), int64(pagination.Offset))
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
		if err := json.Unmarshal(header, blockHeader); err != nil {
			return nil, fmt.Errorf("could not decode block header. Cause: %w", err)
		}
		r := new(common.L2RollupHash)
		r.SetBytes(rollupHash)
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
