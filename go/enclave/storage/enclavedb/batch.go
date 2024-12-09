package enclavedb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/errutil"
	"github.com/ten-protocol/go-ten/go/enclave/core"
)

func WriteBatchHeader(ctx context.Context, dbtx *sql.Tx, batch *core.Batch, convertedHash gethcommon.Hash, blockId int64, isCanonical bool) error {
	header, err := rlp.EncodeToBytes(batch.Header)
	if err != nil {
		return fmt.Errorf("could not encode batch header. Cause: %w", err)
	}
	args := []any{
		batch.Header.SequencerOrderNo.Uint64(), // sequence
		convertedHash,                          // converted_hash
		batch.Hash(),                           // hash
		batch.Header.Number.Uint64(),           // height
		isCanonical,                            // is_canonical
		header,                                 // header blob
		batch.Header.L1Proof.Bytes(),           // l1 proof hash
	}
	if blockId == 0 {
		args = append(args, nil) // l1_proof block id
	} else {
		args = append(args, blockId)
	}
	args = append(args, false) // executed
	_, err = dbtx.ExecContext(ctx, "insert into batch values (?,?,?,?,?,?,?,?,?)", args...)
	return err
}

func UpdateCanonicalBatch(ctx context.Context, dbtx *sql.Tx, isCanonical bool, blocks []common.L1BlockHash) error {
	if len(blocks) == 0 {
		return nil
	}
	args := make([]any, 0)
	args = append(args, isCanonical)
	for _, blockHash := range blocks {
		args = append(args, blockHash.Bytes())
	}

	updateBatches := "update batch set is_canonical=? where " + repeat(" l1_proof_hash=? ", "OR", len(blocks))
	_, err := dbtx.ExecContext(ctx, updateBatches, args...)
	return err
}

func ExistsBatchAtHeight(ctx context.Context, dbTx *sql.Tx, height *big.Int) (bool, error) {
	var exists bool
	err := dbTx.QueryRowContext(ctx, "select exists(select 1 from batch where height=?)", height.Uint64()).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// WriteTransactions - persists the batch and the transactions
func WriteTransactions(ctx context.Context, dbtx *sql.Tx, transactions []*core.TxWithSender, height uint64, isSynthetic bool, senderIds []uint64, toContractIds []*uint64, fromIdx int) error {
	if len(transactions) == 0 {
		return nil
	}
	// creates a batch insert statement for all entries
	insert := "insert into tx (hash, content, to_address, type, sender_address, idx, batch_height, is_synthetic) values " + repeat("(?,?,?,?,?,?,?,?)", ",", len(transactions))

	args := make([]any, 0)
	for i, transaction := range transactions {
		txBytes, err := rlp.EncodeToBytes(transaction.Tx)
		if err != nil {
			return fmt.Errorf("failed to encode block receipts. Cause: %w", err)
		}

		args = append(args, transaction.Tx.Hash().Bytes()) // tx_hash
		args = append(args, txBytes)                       // content
		args = append(args, toContractIds[i])              // To
		args = append(args, transaction.Tx.Type())         // Type
		args = append(args, senderIds[i])                  // sender_address
		args = append(args, fromIdx+i)                     // idx
		args = append(args, height)                        // the batch height which contained it
		args = append(args, isSynthetic)                   // is_synthetic if the transaction is a synthetic (internally derived transaction)
	}
	_, err := dbtx.ExecContext(ctx, insert, args...)
	if err != nil {
		return err
	}

	return nil
}

func IsCanonicalBatchHash(ctx context.Context, dbtx *sql.Tx, hash *gethcommon.Hash) (bool, error) {
	var isCanon bool
	err := dbtx.QueryRowContext(ctx, "select is_canonical from batch where hash=? ", hash.Bytes()).Scan(&isCanon)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return isCanon, err
}

func IsCanonicalBatchSeq(ctx context.Context, db *sql.DB, seqNo uint64) (bool, error) {
	var isCanon bool
	err := db.QueryRowContext(ctx, "select is_canonical from batch where sequence=? ", seqNo).Scan(&isCanon)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return isCanon, err
}

func MarkBatchExecuted(ctx context.Context, dbtx *sql.Tx, seqNo *big.Int) error {
	_, err := dbtx.ExecContext(ctx, "update batch set is_executed=true where sequence=?", seqNo.Uint64())
	return err
}

func WriteReceipt(ctx context.Context, dbtx *sql.Tx, batchSeqNo uint64, txId *uint64, receipt *types.Receipt) (uint64, error) {
	insert := "insert into receipt (post_state, status, cumulative_gas_used, effective_gas_price, created_contract_address, tx, batch) values " + "(?,?,?,?,?,?,?)"
	addr := &receipt.ContractAddress
	if *addr == (gethcommon.Address{}) {
		addr = nil
	}
	var effPrice uint64
	if receipt.EffectiveGasPrice != nil {
		effPrice = receipt.EffectiveGasPrice.Uint64()
	}
	res, err := dbtx.ExecContext(ctx, insert, receipt.PostState, receipt.Status, receipt.CumulativeGasUsed, effPrice, addr, txId, batchSeqNo)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return uint64(id), nil
}

func ReadTransactionIdAndSender(ctx context.Context, dbtx *sql.Tx, txHash gethcommon.Hash) (*uint64, *uint64, error) {
	var txId uint64
	var senderId uint64
	err := dbtx.QueryRowContext(ctx, "select id,sender_address from tx where hash=? ", txHash.Bytes()).Scan(&txId, &senderId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// make sure the error is converted to obscuro-wide not found error
			return nil, nil, errutil.ErrNotFound
		}
		return nil, nil, err
	}
	return &txId, &senderId, err
}

func ReadBatchHeaderBySeqNo(ctx context.Context, db *sql.DB, seqNo uint64) (*common.BatchHeader, error) {
	return fetchBatchHeader(ctx, db, " where sequence=?", seqNo)
}

func ReadBatchHeaderByHash(ctx context.Context, db *sql.DB, hash common.L2BatchHash) (*common.BatchHeader, error) {
	return fetchBatchHeader(ctx, db, " where b.hash=? ", hash.Bytes())
}

func ReadCanonicalBatchHeaderByHeight(ctx context.Context, db *sql.DB, height uint64) (*common.BatchHeader, error) {
	return fetchBatchHeader(ctx, db, " where b.height=? and is_canonical=true", height)
}

func ReadNonCanonicalBatches(ctx context.Context, db *sql.DB, startAtSeq uint64, endSeq uint64) ([]*common.BatchHeader, error) {
	return fetchBatches(ctx, db, " where b.sequence>=? and b.sequence <=? and b.is_canonical=false order by b.sequence", startAtSeq, endSeq)
}

func ReadCanonicalBatches(ctx context.Context, db *sql.DB, startAtSeq uint64, endSeq uint64) ([]*common.BatchHeader, error) {
	return fetchBatches(ctx, db, " where b.sequence>=? and b.sequence <=? and b.is_canonical=true order by b.sequence", startAtSeq, endSeq)
}

// todo - is there a better way to write this query?
func ReadCurrentHeadBatchHeader(ctx context.Context, db *sql.DB) (*common.BatchHeader, error) {
	return fetchBatchHeader(ctx, db, " where b.is_canonical=true and b.is_executed=true and b.height=(select max(b1.height) from batch b1 where b1.is_canonical=true and b1.is_executed=true)")
}

func ReadBatchesByBlock(ctx context.Context, db *sql.DB, hash common.L1BlockHash) ([]*common.BatchHeader, error) {
	return fetchBatches(ctx, db, " where l1_proof_hash=?  order by b.sequence", hash.Bytes())
}

func ReadCurrentSequencerNo(ctx context.Context, db *sql.DB) (*big.Int, error) {
	var seq sql.NullInt64
	query := "select max(sequence) from batch"
	err := db.QueryRowContext(ctx, query).Scan(&seq)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// make sure the error is converted to obscuro-wide not found error
			return nil, errutil.ErrNotFound
		}
		return nil, err
	}
	if !seq.Valid {
		return nil, errutil.ErrNotFound
	}
	return big.NewInt(seq.Int64), nil
}

func fetchBatchHeader(ctx context.Context, db *sql.DB, whereQuery string, args ...any) (*common.BatchHeader, error) {
	var header string
	query := "select b.header from batch b " + whereQuery
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
	h := new(common.BatchHeader)
	if err := rlp.DecodeBytes([]byte(header), h); err != nil {
		return nil, fmt.Errorf("could not decode batch header. Cause: %w", err)
	}

	return h, nil
}

func fetchBatches(ctx context.Context, db *sql.DB, whereQuery string, args ...any) ([]*common.BatchHeader, error) {
	result := make([]*common.BatchHeader, 0)

	rows, err := db.QueryContext(ctx, "select b.header from batch b "+whereQuery, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// make sure the error is converted to obscuro-wide not found error
			return nil, errutil.ErrNotFound
		}
		return nil, err
	}
	defer rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}
	for rows.Next() {
		var header string
		err := rows.Scan(&header)
		if err != nil {
			return nil, err
		}
		h := new(common.BatchHeader)
		if err := rlp.DecodeBytes([]byte(header), h); err != nil {
			return nil, fmt.Errorf("could not decode batch header. Cause: %w", err)
		}

		result = append(result, h)
	}
	return result, nil
}

func ReadReceipt(ctx context.Context, db *sql.DB, txHash common.L2TxHash, requester *gethcommon.Address) (*core.InternalReceipt, error) {
	rec, _, err := loadReceiptsAndEventLogs(ctx, db, requester, " AND curr_tx.hash=?", []any{txHash.Bytes()}, true)
	if err != nil {
		return nil, err
	}
	// there should be only a single receipt
	return rec[0], nil
}

func ExistsReceipt(ctx context.Context, db *sql.DB, txHash common.L2TxHash) (bool, error) {
	query := "select count(1) from receipt rec join tx curr_tx on rec.tx=curr_tx.id where curr_tx.hash=?"
	row := db.QueryRowContext(ctx, query, txHash.Bytes())
	var cnt uint
	err := row.Scan(&cnt)
	if err != nil {
		return false, err
	}
	return cnt > 0, nil
}

func ReadTransaction(ctx context.Context, db *sql.DB, txHash gethcommon.Hash) (*types.Transaction, common.L2BatchHash, uint64, uint64, error) {
	row := db.QueryRowContext(ctx,
		"select tx.content, batch.hash, batch.height, tx.idx from receipt join tx on tx.id=receipt.tx join batch on batch.sequence=receipt.batch where batch.is_canonical=true and tx.hash=?",
		txHash.Bytes())

	// tx, batch, height, idx
	var txData []byte
	var batchHash []byte
	var height uint64
	var idx uint64
	err := row.Scan(&txData, &batchHash, &height, &idx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// make sure the error is converted to obscuro-wide not found error
			return nil, gethcommon.Hash{}, 0, 0, errutil.ErrNotFound
		}
		return nil, gethcommon.Hash{}, 0, 0, err
	}
	tx := new(common.L2Tx)
	if err := rlp.DecodeBytes(txData, tx); err != nil {
		return nil, gethcommon.Hash{}, 0, 0, fmt.Errorf("could not decode L2 transaction. Cause: %w", err)
	}
	batch := gethcommon.Hash{}
	batch.SetBytes(batchHash)
	return tx, batch, height, idx, nil
}

func ReadBatchTransactions(ctx context.Context, db *sql.DB, height uint64) ([]*common.L2Tx, error) {
	var txs []*common.L2Tx

	rows, err := db.QueryContext(ctx, "select content from tx where batch_height=? and is_synthetic=? order by idx", height, false)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// make sure the error is converted to obscuro-wide not found error
			return nil, errutil.ErrNotFound
		}
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		// receipt, tx, batch, height
		var txContent []byte
		err := rows.Scan(&txContent)
		if err != nil {
			return nil, err
		}
		tx := new(common.L2Tx)
		if err := rlp.DecodeBytes(txContent, tx); err != nil {
			return nil, fmt.Errorf("could not decode L2 transaction. Cause: %w", err)
		}
		txs = append(txs, tx)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return txs, nil
}

func ReadContractCreationCount(ctx context.Context, db *sql.DB) (*big.Int, error) {
	row := db.QueryRowContext(ctx, "select count(id) from contract")

	var count int64
	err := row.Scan(&count)
	if err != nil {
		return nil, err
	}

	return big.NewInt(count), nil
}

func ReadUnexecutedBatches(ctx context.Context, db *sql.DB, from *big.Int) ([]*common.BatchHeader, error) {
	return fetchBatches(ctx, db, "where is_executed=false and is_canonical=true and sequence >= ? order by b.sequence", from.Uint64())
}

func BatchWasExecuted(ctx context.Context, db *sql.DB, hash common.L2BatchHash) (bool, error) {
	row := db.QueryRowContext(ctx, "select is_executed from batch where hash=? ", hash.Bytes())

	var result bool
	err := row.Scan(&result)
	if err != nil {
		// When there are no rows returned it means there is no canonical batch with that hash.
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return result, nil
}

func GetTransactionsPerAddress(ctx context.Context, db *sql.DB, address *gethcommon.Address, pagination *common.QueryPagination) ([]*core.InternalReceipt, error) {
	receipts, err := loadReceiptList(ctx, db, address, " AND tx_sender.address = ? ", []any{address.Bytes()}, " ORDER BY b.height DESC LIMIT ? OFFSET ?", []any{pagination.Size, pagination.Offset})
	if err != nil {
		return nil, err
	}
	return receipts, nil
}

func CountTransactionsPerAddress(ctx context.Context, db *sql.DB, address *gethcommon.Address) (uint64, error) {
	row := db.QueryRowContext(ctx, "select count(1) from receipt "+
		"join tx on tx.id=receipt.tx "+
		"join externally_owned_account eoa on eoa.id = tx.sender_address "+
		"where eoa.address = ?", address.Bytes())

	var count uint64
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func FetchConvertedBatchHash(ctx context.Context, db *sql.DB, seqNo uint64) (gethcommon.Hash, error) {
	var hash []byte

	query := "select converted_hash from batch where sequence=?"
	err := db.QueryRowContext(ctx, query, seqNo).Scan(&hash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// make sure the error is converted to obscuro-wide not found error
			return gethcommon.Hash{}, errutil.ErrNotFound
		}
		return gethcommon.Hash{}, err
	}
	return gethcommon.BytesToHash(hash), nil
}
