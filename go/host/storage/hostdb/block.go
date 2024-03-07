package hostdb

import (
	"database/sql"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ten-protocol/go-ten/go/common"
)

const (
	blockInsert = "insert into block_host values (?,?,?,?)"
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
