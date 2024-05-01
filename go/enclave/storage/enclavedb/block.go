package enclavedb

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/errutil"
)

func WriteBlock(_ context.Context, dbtx DBTransaction, b *types.Header) error {
	header, err := rlp.EncodeToBytes(b)
	if err != nil {
		return fmt.Errorf("could not encode block header. Cause: %w", err)
	}

	dbtx.ExecuteSQL("insert into block (hash,full_hash,is_canonical,header,height) values (?,?,?,?,?)",
		truncTo4(b.Hash()), // hash
		b.Hash().Bytes(),   // full_hash
		true,               // is_canonical
		header,             // header
		b.Number.Uint64(),  // height
	)
	return nil
}

func UpdateCanonicalBlocks(ctx context.Context, dbtx DBTransaction, canonical []common.L1BlockHash, nonCanonical []common.L1BlockHash) {
	if len(nonCanonical) > 0 {
		updateCanonicalValue(ctx, dbtx, false, nonCanonical)
	}
	if len(canonical) > 0 {
		updateCanonicalValue(ctx, dbtx, true, canonical)
	}
}

func updateCanonicalValue(_ context.Context, dbtx DBTransaction, isCanonical bool, blocks []common.L1BlockHash) {
	if len(blocks) > 1 {
		println("!!!FFFFFOOOOORRRR")
	}
	canonicalBlocks := repeat("(hash=? and full_hash=?)", "OR", len(blocks))

	args := make([]any, 0)
	args = append(args, isCanonical)
	for _, blockHash := range blocks {
		args = append(args, truncTo4(blockHash), blockHash.Bytes())
	}

	rows, err := dbtx.GetDB().Query("select id from block where "+canonicalBlocks, args[1:]...)
	defer rows.Close()
	if err != nil {
		panic(err)
		return
	}
	for rows.Next() {
		var id uint64
		rows.Scan(&id)
		fmt.Printf("Update canonical=%t block id: %v, hash: %s\n", isCanonical, id, blocks[0].Hex())
	}

	updateBlocks := "update block set is_canonical=? where " + canonicalBlocks
	dbtx.ExecuteSQL(updateBlocks, args...)

	updateBatches := "update batch set is_canonical=? where l1_proof in (select id from block where " + canonicalBlocks + ")"
	dbtx.ExecuteSQL(updateBatches, args...)
}

// todo - remove this. For now creates a "block" but without a body.
func FetchBlock(ctx context.Context, db *sql.DB, hash common.L1BlockHash) (*types.Block, error) {
	return fetchBlock(ctx, db, " where hash=? and full_hash=?", truncTo4(hash), hash.Bytes())
}

func FetchHeadBlock(ctx context.Context, db *sql.DB) (*types.Block, error) {
	// todo - just read the one with the max id
	return fetchBlock(ctx, db, "where is_canonical=true and height=(select max(b.height) from block b where is_canonical=true)")
}

func FetchBlockHeaderByHeight(ctx context.Context, db *sql.DB, height *big.Int) (*types.Header, error) {
	return fetchBlockHeader(ctx, db, "where is_canonical=true and height=?", height.Int64())
}

func GetBlockId(ctx context.Context, db *sql.DB, hash common.L1BlockHash) (*uint64, error) {
	var id uint64
	err := db.QueryRowContext(ctx, "select id from block where hash=? and full_hash=?", truncTo4(hash), hash).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &id, err
}

func WriteL1Messages[T any](ctx context.Context, db *sql.DB, blockId *uint64, messages []T, isValueTransfer bool) error {
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
	query := "select message from l1_msg m join block b on m.block=b.id where b.hash = ? and b.full_hash = ? and is_transfer = ?"
	rows, err := db.QueryContext(ctx, query, truncTo4(blockHash), blockHash.Bytes(), isTransfer)
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

func WriteRollup(_ context.Context, dbtx DBTransaction, rollup *common.RollupHeader, blockId *uint64, internalHeader *common.CalldataRollupHeader) error {
	// Write the encoded header
	data, err := rlp.EncodeToBytes(rollup)
	if err != nil {
		return fmt.Errorf("could not encode batch header. Cause: %w", err)
	}
	dbtx.ExecuteSQL("replace into rollup (hash, full_hash, start_seq, end_seq, time_stamp, header, compression_block) values (?,?,?,?,?,?,?)",
		truncTo4(rollup.Hash()),
		rollup.Hash().Bytes(),
		internalHeader.FirstBatchSequence.Uint64(),
		rollup.LastBatchSeqNo,
		internalHeader.StartTime,
		data,
		blockId,
	)
	return nil
}

func FetchReorgedRollup(ctx context.Context, db *sql.DB, reorgedBlocks []common.L1BlockHash) (*common.L2BatchHash, error) {
	whereClause := repeat("(b.hash=? and b.full_hash=?)", "OR", len(reorgedBlocks))

	query := "select full_hash from rollup r join block b on r.compression_block=b.id where " + whereClause

	args := make([]any, 0)
	for _, blockHash := range reorgedBlocks {
		args = append(args, truncTo4(blockHash), blockHash.Bytes())
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
		"select start_seq, time_stamp from rollup where hash = ? and full_hash=?", truncTo4(hash), hash.Bytes(),
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
