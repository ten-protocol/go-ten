package enclavedb

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/big"

	gethlog "github.com/ethereum/go-ethereum/log"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/errutil"
)

func WriteBlock(ctx context.Context, dbtx *sql.Tx, b *types.Header) error {
	header, err := rlp.EncodeToBytes(b)
	if err != nil {
		return fmt.Errorf("could not encode block header. Cause: %w", err)
	}

	_, err = dbtx.ExecContext(ctx, "insert into block (hash,is_canonical,header,height) values (?,?,?,?)",
		b.Hash().Bytes(),  // hash
		true,              // is_canonical
		header,            // header
		b.Number.Uint64(), // height
	)
	return err
}

func UpdateCanonicalBlocks(ctx context.Context, dbtx *sql.Tx, canonical []common.L1BlockHash, nonCanonical []common.L1BlockHash, logger gethlog.Logger) error {
	if len(nonCanonical) > 0 {
		err := updateCanonicalValue(ctx, dbtx, false, nonCanonical, logger)
		if err != nil {
			return err
		}
	}
	if len(canonical) > 0 {
		err := updateCanonicalValue(ctx, dbtx, true, canonical, logger)
		if err != nil {
			return err
		}
	}
	return nil
}

func updateCanonicalValue(ctx context.Context, dbtx *sql.Tx, isCanonical bool, blocks []common.L1BlockHash, _ gethlog.Logger) error {
	currentBlocks := repeat(" hash=? ", "OR", len(blocks))

	args := make([]any, 0)
	args = append(args, isCanonical)
	for _, blockHash := range blocks {
		args = append(args, blockHash.Bytes())
	}

	updateBlocks := "update block set is_canonical=? where " + currentBlocks
	_, err := dbtx.ExecContext(ctx, updateBlocks, args...)
	if err != nil {
		return err
	}

	updateBatches := "update batch set is_canonical=? where " + repeat(" l1_proof_hash=? ", "OR", len(blocks))
	_, err = dbtx.ExecContext(ctx, updateBatches, args...)
	if err != nil {
		return err
	}

	return nil
}

// HandleBlockArrivedAfterBatches- handle the corner case where the block wasn't available when the batch was received
func HandleBlockArrivedAfterBatches(ctx context.Context, dbtx *sql.Tx, blockId int64, blockHash common.L1BlockHash) error {
	_, err := dbtx.ExecContext(ctx, "update batch set l1_proof=?, is_canonical=true where l1_proof_hash=?", blockId, blockHash.Bytes())
	return err
}

// todo - remove this. For now creates a "block" but without a body.
func FetchBlock(ctx context.Context, db *sql.DB, hash common.L1BlockHash) (*types.Block, error) {
	return fetchBlock(ctx, db, " where hash=?", hash.Bytes())
}

func FetchHeadBlock(ctx context.Context, db *sql.DB) (*types.Block, error) {
	return fetchBlock(ctx, db, "order by id desc limit 1")
}

func FetchBlockHeaderByHeight(ctx context.Context, db *sql.DB, height *big.Int) (*types.Header, error) {
	return fetchBlockHeader(ctx, db, "where is_canonical=true and height=?", height.Int64())
}

func GetBlockId(ctx context.Context, db *sql.Tx, hash common.L1BlockHash) (int64, error) {
	var id int64
	err := db.QueryRowContext(ctx, "select id from block where hash=? ", hash.Bytes()).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, err
}

func WriteL1Messages[T any](ctx context.Context, db *sql.Tx, blockId int64, messages []T, isValueTransfer bool) error {
	insert := "insert into l1_msg (message, block, is_transfer) values " + repeat("(?,?,?)", ",", len(messages))

	args := make([]any, 0)

	for _, msg := range messages {
		data, err := rlp.EncodeToBytes(msg)
		if err != nil {
			return err
		}
		args = append(args, data)
		args = append(args, blockId)
		args = append(args, isValueTransfer)
	}
	if len(messages) > 0 {
		_, err := db.ExecContext(ctx, insert, args...)
		return err
	}
	return nil
}

func FetchL1Messages[T any](ctx context.Context, db *sql.DB, blockHash common.L1BlockHash, isTransfer bool) ([]T, error) {
	var result []T
	query := "select message from l1_msg m join block b on m.block=b.id where b.hash = ? and is_transfer = ?"
	rows, err := db.QueryContext(ctx, query, blockHash.Bytes(), isTransfer)
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
		ccm := new(T)
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

func WriteRollup(ctx context.Context, dbtx *sql.Tx, rollup *common.RollupHeader, blockId int64, internalHeader *common.CalldataRollupHeader) error {
	// Write the encoded header
	data, err := rlp.EncodeToBytes(rollup)
	if err != nil {
		return fmt.Errorf("could not encode batch header. Cause: %w", err)
	}
	_, err = dbtx.ExecContext(ctx, "replace into rollup (hash, start_seq, end_seq, time_stamp, header, compression_block) values (?,?,?,?,?,?)",
		rollup.Hash().Bytes(),
		internalHeader.FirstBatchSequence.Uint64(),
		rollup.LastBatchSeqNo,
		internalHeader.StartTime,
		data,
		blockId,
	)
	if err != nil {
		return err
	}

	return nil
}

func FetchReorgedRollup(ctx context.Context, db *sql.DB, reorgedBlocks []common.L1BlockHash) (*common.L2BatchHash, error) {
	whereClause := repeat(" b.hash=? ", "OR", len(reorgedBlocks))

	query := "select r.hash from rollup r join block b on r.compression_block=b.id where " + whereClause

	args := make([]any, 0)
	for _, blockHash := range reorgedBlocks {
		args = append(args, blockHash.Bytes())
	}
	rollup := new(common.L2BatchHash)
	err := db.QueryRowContext(ctx, query, args...).Scan(&rollup)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// make sure the error is converted to obscuro-wide not found error
			return nil, errutil.ErrNotFound
		}
		return nil, err
	}
	return rollup, nil
}

func FetchRollupMetadata(ctx context.Context, db *sql.DB, hash common.L2RollupHash) (*common.PublicRollupMetadata, error) {
	var startSeq int64
	var startTime uint64

	rollup := new(common.PublicRollupMetadata)
	err := db.QueryRowContext(ctx,
		"select start_seq, time_stamp from rollup where hash = ?", hash.Bytes(),
	).Scan(&startSeq, &startTime)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errutil.ErrNotFound
		}
		return nil, err
	}
	rollup.FirstBatchSequence = big.NewInt(startSeq)
	rollup.StartTime = startTime
	return rollup, nil
}

func fetchBlockHeader(ctx context.Context, db *sql.DB, whereQuery string, args ...any) (*types.Header, error) {
	var header string
	query := "select header from block " + whereQuery
	var err error
	if len(args) > 0 {
		err = db.QueryRowContext(ctx, query, args...).Scan(&header)
	} else {
		err = db.QueryRowContext(ctx, query).Scan(&header)
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

func fetchBlock(ctx context.Context, db *sql.DB, whereQuery string, args ...any) (*types.Block, error) {
	h, err := fetchBlockHeader(ctx, db, whereQuery, args...)
	if err != nil {
		return nil, err
	}
	return types.NewBlockWithHeader(h), nil
}
