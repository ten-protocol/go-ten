package hostdb

import (
	"database/sql"
	"fmt"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ten-protocol/go-ten/go/common"
)

const (
	blockInsert  = "insert into block_host values (?,?,?,?)"
	selectBlocks = "SELECT id, hash, header, rollup_hash FROM block_host ORDER BY id DESC LIMIT ? OFFSET ?"
)

func AddBlock(db *sql.DB, b *types.Header, rollupHash common.L2RollupHash) error {
	header, err := rlp.EncodeToBytes(b)
	if err != nil {
		return fmt.Errorf("could not encode block header. Cause: %w", err)
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(blockInsert,
		truncTo16(b.Hash()), // hash
		header,              // l1 block header
		rollupHash,          // rollup hash
	)
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("could not store block in db: %w", err)
	}
	return nil
}

func GetBlockListing(db *sql.DB, pagination *common.QueryPagination) (*common.BlockListingResponse, error) {
	rows, err := db.Query(selectBlocks, pagination.Size, pagination.Offset)
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
		r := new(gethcommon.Hash)
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
