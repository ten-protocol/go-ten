package hostdb

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ten-protocol/go-ten/go/common/errutil"
	"math/big"

	gethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ten-protocol/go-ten/go/common"
)

const (
	selectBlocks     = "SELECT b.id, b.hash, b.header, r.hash FROM block_host b LEFT JOIN rollup_host r on r.compression_block=b.id ORDER BY b.id DESC "
	selectBlockId    = "SELECT id FROM block_host WHERE hash = "
	selectBlock      = "SELECT header FROM block_host WHERE hash = "
	selectBlockCount = "SELECT total FROM block_count WHERE id = 1"
)

// AddBlock stores a block header with the given rollupHash it contains in the host DB
func AddBlock(dbtx *sql.Tx, statements *SQLStatements, b *types.Header) error {
	header, err := json.Marshal(b)
	if err != nil {
		return fmt.Errorf("could not encode block header. Cause: %w", err)
	}

	println("ADDBLOCK with number: ", b.Number.Int64())
	_, err = dbtx.Exec(statements.InsertBlock,
		b.Hash().Bytes(), // hash
		header,           // l1 block header
	)
	if err != nil {
		return fmt.Errorf("could not insert block. Cause: %w", err)
	}

	var currentTotal int
	err = dbtx.QueryRow(selectBlockCount).Scan(&currentTotal)
	if err != nil {
		return fmt.Errorf("failed to query block count: %w", err)
	}

	newTotal := currentTotal + 1
	_, err = dbtx.Exec(statements.UpdateBlockCount, newTotal)
	if err != nil {
		return fmt.Errorf("failed to update block count: %w", err)
	}

	return nil
}

// GetBlockId returns the block ID given the hash.
func GetBlockId(db *sql.Tx, statements *SQLStatements, hash gethcommon.Hash) (*int64, error) {
	query := selectBlockId + statements.Placeholder
	var blockId int64
	err := db.QueryRow(query, hash.Bytes()).Scan(&blockId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errutil.ErrNotFound
		}
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
		var hash, header []byte
		var rollupHash []byte

		err = rows.Scan(&id, &hash, &header, &rollupHash)
		if err != nil {
			return nil, err
		}

		blockHeader := new(types.Header)
		if err := json.Unmarshal(header, blockHeader); err != nil {
			return nil, fmt.Errorf("could not decode block header. Cause: %w", err)
		}

		var rollupHashValue common.L2RollupHash
		if rollupHash != nil {
			rollupHashValue.SetBytes(rollupHash)
		} else {
			// set to zero hash if no rollup associated
			rollupHashValue = common.L2RollupHash{}
		}

		block := common.PublicBlock{
			BlockHeader: *blockHeader,
			RollupHash:  rollupHashValue,
		}
		blocks = append(blocks, block)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	totalBlocks, err := GetTotalBlockCount(db)
	if err != nil {
		return nil, fmt.Errorf("could not fetch the block count. Cause: %w", err)
	}

	return &common.BlockListingResponse{
		BlocksData: blocks,
		Total:      totalBlocks.Uint64(),
	}, nil
}

// GetTotalBlockCount returns value from the block count table
func GetTotalBlockCount(db HostDB) (*big.Int, error) {
	var totalCount int
	err := db.GetSQLDB().QueryRow(selectBlockCount).Scan(&totalCount)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve total block count: %w", err)
	}
	return big.NewInt(int64(totalCount)), nil
}
