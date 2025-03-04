package enclavedb

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	gethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/errutil"
)

func WriteBlock(ctx context.Context, dbtx *sql.Tx, b *types.Header) error {
	header, err := encodeHeader(b)
	if err != nil {
		return fmt.Errorf("could not encode block header. Cause: %w", err)
	}

	_, err = dbtx.ExecContext(ctx, "insert into block (hash,is_canonical,header,height, processed) values (?,?,?,?,?)",
		b.Hash().Bytes(),  // hash
		true,              // is_canonical
		header,            // header
		b.Number.Uint64(), // height
		false,             // processed
	)
	return err
}

func UpdateCanonicalBlock(ctx context.Context, dbtx *sql.Tx, isCanonical bool, blocks []common.L1BlockHash) error {
	if len(blocks) == 0 {
		return nil
	}
	args := make([]any, 0)
	args = append(args, isCanonical)
	for _, blockHash := range blocks {
		args = append(args, blockHash.Bytes())
	}

	updateBlocks := "update block set is_canonical=? where " + repeat(" hash=? ", "OR", len(blocks))
	_, err := dbtx.ExecContext(ctx, updateBlocks, args...)
	return err
}

func IsCanonicalBlock(ctx context.Context, dbtx *sql.Tx, hash *gethcommon.Hash) (bool, error) {
	var isCanon bool
	err := dbtx.QueryRowContext(ctx, "select is_canonical from block where hash=? ", hash.Bytes()).Scan(&isCanon)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return isCanon, err
}

/*// CheckCanonicalValidity - expensive but useful for debugging races
func CheckCanonicalValidity(ctx context.Context, dbtx *sql.Tx, batchId int64) error {
	rows, err := dbtx.QueryContext(ctx, "select count(*), height from batch where height >=? AND is_canonical=true group by height having count(*) >1", batchId)
	if err != nil {
		return err
	}
	defer rows.Close()
	if err := rows.Err(); err != nil {
		return err
	}
	if rows.Next() {
		var cnt uint64
		var height uint64
		err := rows.Scan(&cnt, &height)
		if err != nil {
			return err
		}
		return fmt.Errorf("found multiple (%d) canonical batches for height %d", cnt, height)
	}
	return nil
}
*/

// HandleBlockArrivedAfterBatches - handle the corner case where the block wasn't available when the batch was received
func HandleBlockArrivedAfterBatches(ctx context.Context, dbtx *sql.Tx, _ int64, blockHash common.L1BlockHash) error {
	_, err := dbtx.ExecContext(ctx, "update batch set is_canonical=true where l1_proof_hash=?", blockHash.Bytes())
	return err
}

func FetchBlockHeader(ctx context.Context, db *sql.DB, hash common.L1BlockHash) (*types.Header, error) {
	return fetchBlock(ctx, db, " where hash=?", hash.Bytes())
}

func FetchHeadBlock(ctx context.Context, db *sql.DB) (*types.Header, error) {
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
		rollup.FirstBatchSeqNo,
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
	return decodeHeader([]byte(header))
}

func fetchBlock(ctx context.Context, db *sql.DB, whereQuery string, args ...any) (*types.Header, error) {
	return fetchBlockHeader(ctx, db, whereQuery, args...)
}

func encodeHeader(h *types.Header) ([]byte, error) {
	return json.Marshal(h)
}

func decodeHeader(b []byte) (*types.Header, error) {
	h := new(types.Header)
	err := json.Unmarshal(b, h)
	if err != nil {
		return nil, fmt.Errorf("could not decode l1 block header. Cause: %w", err)
	}
	return h, nil
}

func UpdateBlockProcessed(ctx context.Context, dbtx *sql.Tx, block common.L1BlockHash) error {
	_, err := dbtx.ExecContext(ctx, "update block set processed=true where hash=?", block.Bytes())
	return err
}

func SelectUnprocessedBlocks(ctx context.Context, dbtx *sql.Tx) ([]uint64, error) {
	query := "select id from block where processed=false"
	rows, err := dbtx.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ids := make([]uint64, 0)
	for rows.Next() {
		var id uint64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, rows.Err()
}

func DeleteUnprocessedRollups(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, "delete from rollup where compression_block in (select id from block where processed=false)")
	return err
}

func DeleteUnprocessedL1Messages(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, "delete from l1_msg where block in (select id from block where processed=false)")
	return err
}

func DeleteUnProcessedBlock(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, "delete from block where processed=false")
	return err
}
