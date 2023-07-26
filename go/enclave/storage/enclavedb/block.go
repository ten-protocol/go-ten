package enclavedb

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/errutil"
)

const (
	blockInsert       = "insert into block values (?,?,?,?,?)"
	selectBlockHeader = "select header from block"

	l1msgInsert = "insert into l1_msg (message, block) values "
	l1msgValue  = "(?,?)"
	selectL1Msg = "select message from l1_msg "

	rollupInsert = "insert into rollup values (?,?,?,?)"

	updateCanonicalBlock   = "update block set is_canonical=? where hash in "
	updateCanonicalBatches = "update batch set is_canonical=? where l1_proof in "
)

func WriteBlock(dbtx DBTransaction, b *types.Header) error {
	header, err := rlp.EncodeToBytes(b)
	if err != nil {
		return fmt.Errorf("could not encode block header. Cause: %w", err)
	}

	var parentBytes []byte
	if b.Number.Uint64() > 1 {
		parentBytes = b.ParentHash.Bytes()
	}
	dbtx.ExecuteSQL(blockInsert,
		b.Hash().Bytes(),
		parentBytes,
		true,
		header,
		b.Number.Uint64(),
	)
	return nil
}

func UpdateCanonicalBlocks(dbtx DBTransaction, canonical []common.L1BlockHash, nonCanonical []common.L1BlockHash) {
	if len(nonCanonical) > 0 {
		updateCanonicalValue(dbtx, false, nonCanonical)
	}
	if len(canonical) > 0 {
		updateCanonicalValue(dbtx, true, canonical)
	}
}

func updateCanonicalValue(dbtx DBTransaction, isCanonical bool, values []common.L1BlockHash) {
	vals := strings.Repeat("?,", len(values))
	updateBlocks := updateCanonicalBlock + "(" + vals[0:len(vals)-1] + ")"
	updateBatches := updateCanonicalBatches + "(" + vals[0:len(vals)-1] + ")"

	args := make([]any, 0)
	args = append(args, isCanonical)
	for _, value := range values {
		args = append(args, value.Bytes())
	}
	dbtx.ExecuteSQL(updateBlocks, args...)
	dbtx.ExecuteSQL(updateBatches, args...)
}

func FetchBlockHeader(db *sql.DB, hash common.L2BatchHash) (*types.Header, error) {
	return fetchBlockHeader(db, " where hash=?", hash.Bytes())
}

// todo - remove this. For now creates a "block" but without a body.
func FetchBlock(db *sql.DB, hash common.L2BatchHash) (*types.Block, error) {
	return fetchBlock(db, " where hash=?", hash.Bytes())
}

func FetchHeadBlock(db *sql.DB) (*types.Block, error) {
	return fetchBlock(db, "where is_canonical and height=(select max(b.height) from block b where is_canonical)")
}

func WriteL1Messages(db *sql.DB, blockHash common.L1BlockHash, messages common.CrossChainMessages) error {
	insert := l1msgInsert + strings.Repeat(l1msgValue+",", len(messages))
	args := make([]any, 0)

	for _, msg := range messages {
		data, err := rlp.EncodeToBytes(msg)
		if err != nil {
			return err
		}
		args = append(args, data)
		args = append(args, blockHash.Bytes())
	}
	if len(messages) > 0 {
		_, err := db.Exec(insert[0:len(insert)-1], args...)
		return err
	}
	return nil
}

func FetchL1Messages(db *sql.DB, blockHash common.L1BlockHash) (common.CrossChainMessages, error) {
	var result common.CrossChainMessages
	query := selectL1Msg + " where block = ?"
	rows, err := db.Query(query, blockHash.Bytes())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// make sure the error is converted to obscuro-wide not found error
			return nil, errutil.ErrNotFound
		}
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var msg []byte
		err := rows.Scan(&msg)
		if err != nil {
			return nil, err
		}
		ccm := new(common.CrossChainMessage)
		if err := rlp.Decode(bytes.NewReader(msg), ccm); err != nil {
			return nil, fmt.Errorf("could not decode cross chain message. Cause: %w", err)
		}
		result = append(result, *ccm)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return result, nil
}

func WriteRollup(dbtx DBTransaction, rollup *common.RollupHeader) error {
	// Write the encoded header
	data, err := rlp.EncodeToBytes(rollup)
	if err != nil {
		return fmt.Errorf("could not encode batch header. Cause: %w", err)
	}
	dbtx.ExecuteSQL(rollupInsert,
		0,
		0,
		data,
		rollup.L1Proof.Bytes(),
	)
	return nil
}

func fetchBlockHeader(db *sql.DB, whereQuery string, args ...any) (*types.Header, error) {
	var header string
	query := selectBlockHeader + " " + whereQuery
	var err error
	if len(args) > 0 {
		err = db.QueryRow(query, args...).Scan(&header)
	} else {
		err = db.QueryRow(query).Scan(&header)
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// make sure the error is converted to obscuro-wide not found error
			return nil, errutil.ErrNotFound
		}
		return nil, err
	}
	h := new(types.Header)
	if err := rlp.Decode(bytes.NewReader([]byte(header)), h); err != nil {
		return nil, fmt.Errorf("could not decode l1 block header. Cause: %w", err)
	}

	return h, nil
}

func fetchBlock(db *sql.DB, whereQuery string, args ...any) (*types.Block, error) {
	h, err := fetchBlockHeader(db, whereQuery, args...)
	if err != nil {
		return nil, err
	}
	return types.NewBlockWithHeader(h), nil
}
