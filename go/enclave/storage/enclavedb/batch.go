package enclavedb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/params"

	"github.com/ethereum/go-ethereum/core/types"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/errutil"
	"github.com/ten-protocol/go-ten/go/enclave/core"
)

const (
	selectBatch = "select b.header, bb.content from batch b join batch_body bb on b.body=bb.id"

	queryReceipts = "select exec_tx.receipt, tx.content, batch.hash, batch.height from exec_tx join tx on tx.id=exec_tx.tx join batch on batch.sequence=exec_tx.batch "
)

// WriteBatchAndTransactions - persists the batch and the transactions
func WriteBatchAndTransactions(ctx context.Context, dbtx *sql.Tx, batch *core.Batch, convertedHash gethcommon.Hash, blockId int64) error {
	// todo - optimize for reorgs
	batchBodyID := batch.SeqNo().Uint64()

	body, err := rlp.EncodeToBytes(batch.Transactions)
	if err != nil {
		return fmt.Errorf("could not encode L2 transactions. Cause: %w", err)
	}
	header, err := rlp.EncodeToBytes(batch.Header)
	if err != nil {
		return fmt.Errorf("could not encode batch header. Cause: %w", err)
	}

	_, err = dbtx.ExecContext(ctx, "replace into batch_body values (?,?)", batchBodyID, body)
	if err != nil {
		return err
	}

	var isCanon bool
	err = dbtx.QueryRowContext(ctx,
		"select is_canonical from block where hash=? ",
		batch.Header.L1Proof.Bytes(),
	).Scan(&isCanon)
	if err != nil {
		// if the block is not found, we assume it is non-canonical
		// fmt.Printf("IsCanon %s err: %s\n", batch.Header.L1Proof, err)
		isCanon = false
	}

	args := []any{
		batch.Header.SequencerOrderNo.Uint64(), // sequence
		convertedHash,                          // converted_hash
		batch.Hash(),                           // hash
		batch.Header.Number.Uint64(),           // height
		isCanon,                                // is_canonical
		header,                                 // header blob
		batchBodyID,                            // reference to the batch body
		batch.Header.L1Proof.Bytes(),           // l1 proof hash
	}
	if blockId == 0 {
		args = append(args, nil) // l1_proof block id
	} else {
		args = append(args, blockId)
	}
	args = append(args, false) // executed
	_, err = dbtx.ExecContext(ctx, "insert into batch values (?,?,?,?,?,?,?,?,?,?)", args...)
	if err != nil {
		return err
	}

	// creates a big insert statement for all transactions
	if len(batch.Transactions) > 0 {
		insert := "replace into tx (hash, content, sender_address, nonce, idx, body) values " + repeat("(?,?,?,?,?,?)", ",", len(batch.Transactions))

		args := make([]any, 0)
		for i, transaction := range batch.Transactions {
			txBytes, err := rlp.EncodeToBytes(transaction)
			if err != nil {
				return fmt.Errorf("failed to encode block receipts. Cause: %w", err)
			}

			from, err := types.Sender(types.LatestSignerForChainID(transaction.ChainId()), transaction)
			if err != nil {
				return fmt.Errorf("unable to convert tx to message - %w", err)
			}

			args = append(args, transaction.Hash())  // tx_hash
			args = append(args, txBytes)             // content
			args = append(args, from.Bytes())        // sender_address
			args = append(args, transaction.Nonce()) // nonce
			args = append(args, i)                   // idx
			args = append(args, batchBodyID)         // the batch body which contained it
		}
		_, err = dbtx.ExecContext(ctx, insert, args...)
		if err != nil {
			return err
		}
	}

	return nil
}

// WriteBatchExecution - save receipts
func WriteBatchExecution(ctx context.Context, dbtx *sql.Tx, seqNo *big.Int, receipts []*types.Receipt) error {
	_, err := dbtx.ExecContext(ctx, "update batch set is_executed=true where sequence=?", seqNo.Uint64())
	if err != nil {
		return err
	}

	args := make([]any, 0)
	for _, receipt := range receipts {
		// Convert the receipt into their storage form and serialize them
		storageReceipt := (*types.ReceiptForStorage)(receipt)
		receiptBytes, err := rlp.EncodeToBytes(storageReceipt)
		if err != nil {
			return fmt.Errorf("failed to encode block receipts. Cause: %w", err)
		}

		// ignore the error because synthetic transactions will not be inserted
		txId, _ := GetTxId(ctx, dbtx, storageReceipt.TxHash)
		args = append(args, receipt.ContractAddress.Bytes()) // created_contract_address
		args = append(args, receiptBytes)                    // the serialised receipt
		if txId == 0 {
			args = append(args, nil) // tx id
		} else {
			args = append(args, txId) // tx id
		}
		args = append(args, seqNo.Uint64()) // batch_seq
	}
	if len(args) > 0 {
		insert := "insert into exec_tx (created_contract_address, receipt, tx, batch) values " + repeat("(?,?,?,?)", ",", len(receipts))
		_, err = dbtx.ExecContext(ctx, insert, args...)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetTxId(ctx context.Context, dbtx *sql.Tx, txHash gethcommon.Hash) (int64, error) {
	var txId int64
	err := dbtx.QueryRowContext(ctx, "select id from tx where hash=? ", txHash.Bytes()).Scan(&txId)
	return txId, err
}

func ReadBatchBySeqNo(ctx context.Context, db *sql.DB, seqNo uint64) (*core.Batch, error) {
	return fetchBatch(ctx, db, " where sequence=?", seqNo)
}

func ReadBatchByHash(ctx context.Context, db *sql.DB, hash common.L2BatchHash) (*core.Batch, error) {
	return fetchBatch(ctx, db, " where b.hash=? ", hash.Bytes())
}

func ReadCanonicalBatchByHeight(ctx context.Context, db *sql.DB, height uint64) (*core.Batch, error) {
	return fetchBatch(ctx, db, " where b.height=? and is_canonical=true", height)
}

func ReadNonCanonicalBatches(ctx context.Context, db *sql.DB, startAtSeq uint64, endSeq uint64) ([]*core.Batch, error) {
	return fetchBatches(ctx, db, " where b.sequence>=? and b.sequence <=? and b.is_canonical=false order by b.sequence", startAtSeq, endSeq)
}

// todo - is there a better way to write this query?
func ReadCurrentHeadBatch(ctx context.Context, db *sql.DB) (*core.Batch, error) {
	return fetchBatch(ctx, db, " where b.is_canonical=true and b.is_executed=true and b.height=(select max(b1.height) from batch b1 where b1.is_canonical=true and b1.is_executed=true)")
}

func ReadBatchesByBlock(ctx context.Context, db *sql.DB, hash common.L1BlockHash) ([]*core.Batch, error) {
	return fetchBatches(ctx, db, " join block l1b on b.l1_proof=l1b.id where l1b.hash=?  order by b.sequence", hash.Bytes())
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

func fetchBatch(ctx context.Context, db *sql.DB, whereQuery string, args ...any) (*core.Batch, error) {
	var header string
	var body []byte
	query := selectBatch + " " + whereQuery
	var err error
	if len(args) > 0 {
		err = db.QueryRowContext(ctx, query, args...).Scan(&header, &body)
	} else {
		err = db.QueryRowContext(ctx, query).Scan(&header, &body)
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
	txs := new([]*common.L2Tx)
	if err := rlp.DecodeBytes(body, txs); err != nil {
		return nil, fmt.Errorf("could not decode L2 transactions %v. Cause: %w", body, err)
	}

	b := core.Batch{
		Header:       h,
		Transactions: *txs,
	}

	return &b, nil
}

func fetchBatches(ctx context.Context, db *sql.DB, whereQuery string, args ...any) ([]*core.Batch, error) {
	result := make([]*core.Batch, 0)

	rows, err := db.QueryContext(ctx, selectBatch+" "+whereQuery, args...)
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
		var body []byte
		err := rows.Scan(&header, &body)
		if err != nil {
			return nil, err
		}
		h := new(common.BatchHeader)
		if err := rlp.DecodeBytes([]byte(header), h); err != nil {
			return nil, fmt.Errorf("could not decode batch header. Cause: %w", err)
		}
		txs := new([]*common.L2Tx)
		if err := rlp.DecodeBytes(body, txs); err != nil {
			return nil, fmt.Errorf("could not decode L2 transactions %v. Cause: %w", body, err)
		}

		result = append(result,
			&core.Batch{
				Header:       h,
				Transactions: *txs,
			})
	}
	return result, nil
}

func selectReceipts(ctx context.Context, db *sql.DB, config *params.ChainConfig, query string, args ...any) (types.Receipts, error) {
	var allReceipts types.Receipts

	// where batch=?
	rows, err := db.QueryContext(ctx, queryReceipts+" "+query, args...)
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
		var receiptData []byte
		var txData []byte
		var batchHash []byte
		var height uint64
		err := rows.Scan(&receiptData, &txData, &batchHash, &height)
		if err != nil {
			return nil, err
		}
		tx := new(common.L2Tx)
		if err := rlp.DecodeBytes(txData, tx); err != nil {
			return nil, fmt.Errorf("could not decode L2 transaction. Cause: %w", err)
		}
		transactions := []*common.L2Tx{tx}

		storageReceipt := new(types.ReceiptForStorage)
		if err := rlp.DecodeBytes(receiptData, storageReceipt); err != nil {
			return nil, fmt.Errorf("unable to decode receipt. Cause : %w", err)
		}
		receipts := (types.Receipts)([]*types.Receipt{(*types.Receipt)(storageReceipt)})

		hash := common.L2BatchHash{}
		hash.SetBytes(batchHash)
		if err = receipts.DeriveFields(config, hash, height, 0, big.NewInt(0), big.NewInt(0), transactions); err != nil {
			return nil, fmt.Errorf("failed to derive block receipts fields. hash = %s; number = %d; err = %w", hash, height, err)
		}
		allReceipts = append(allReceipts, receipts[0])
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return allReceipts, nil
}

// ReadReceiptsByBatchHash retrieves all the transaction receipts belonging to a block, including
// its corresponding metadata fields. If it is unable to populate these metadata
// fields then nil is returned.
//
// The current implementation populates these metadata fields by reading the receipts'
// corresponding block body, so if the block body is not found it will return nil even
// if the receipt itself is stored.
func ReadReceiptsByBatchHash(ctx context.Context, db *sql.DB, hash common.L2BatchHash, config *params.ChainConfig) (types.Receipts, error) {
	return selectReceipts(ctx, db, config, "where batch.hash=? ", hash.Bytes())
}

func ReadReceipt(ctx context.Context, db *sql.DB, txHash common.L2TxHash, config *params.ChainConfig) (*types.Receipt, error) {
	// todo - canonical?
	row := db.QueryRowContext(ctx, queryReceipts+" where tx.hash=? ", txHash.Bytes())
	// receipt, tx, batch, height
	var receiptData []byte
	var txData []byte
	var batchHash []byte
	var height uint64
	err := row.Scan(&receiptData, &txData, &batchHash, &height)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// make sure the error is converted to obscuro-wide not found error
			return nil, errutil.ErrNotFound
		}
		return nil, err
	}
	tx := new(common.L2Tx)
	if err := rlp.DecodeBytes(txData, tx); err != nil {
		return nil, fmt.Errorf("could not decode L2 transaction. Cause: %w", err)
	}
	transactions := []*common.L2Tx{tx}

	storageReceipt := new(types.ReceiptForStorage)
	if err := rlp.DecodeBytes(receiptData, storageReceipt); err != nil {
		return nil, fmt.Errorf("unable to decode receipt. Cause : %w", err)
	}
	receipts := (types.Receipts)([]*types.Receipt{(*types.Receipt)(storageReceipt)})

	batchhash := common.L2BatchHash{}
	batchhash.SetBytes(batchHash)
	// todo base fee
	if err = receipts.DeriveFields(config, batchhash, height, 0, big.NewInt(1), big.NewInt(0), transactions); err != nil {
		return nil, fmt.Errorf("failed to derive block receipts fields. txHash = %s; number = %d; err = %w", txHash, height, err)
	}
	return receipts[0], nil
}

func ReadTransaction(ctx context.Context, db *sql.DB, txHash gethcommon.Hash) (*types.Transaction, common.L2BatchHash, uint64, uint64, error) {
	row := db.QueryRowContext(ctx,
		"select tx.content, batch.hash, batch.height, tx.idx from exec_tx join tx on tx.id=exec_tx.tx join batch on batch.sequence=exec_tx.batch where batch.is_canonical=true and tx.hash=?",
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

func GetContractCreationTx(ctx context.Context, db *sql.DB, address gethcommon.Address) (*gethcommon.Hash, error) {
	row := db.QueryRowContext(ctx, "select tx.hash from exec_tx join tx on tx.id=exec_tx.tx where created_contract_address=? ", address.Bytes())

	var txHashBytes []byte
	err := row.Scan(&txHashBytes)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// make sure the error is converted to obscuro-wide not found error
			return nil, errutil.ErrNotFound
		}
		return nil, err
	}
	txHash := gethcommon.Hash{}
	txHash.SetBytes(txHashBytes)
	return &txHash, nil
}

func ReadContractCreationCount(ctx context.Context, db *sql.DB) (*big.Int, error) {
	row := db.QueryRowContext(ctx, "select count( distinct created_contract_address) from exec_tx ")

	var count int64
	err := row.Scan(&count)
	if err != nil {
		return nil, err
	}

	return big.NewInt(count), nil
}

func ReadUnexecutedBatches(ctx context.Context, db *sql.DB, from *big.Int) ([]*core.Batch, error) {
	return fetchBatches(ctx, db, "where is_executed=false and is_canonical=true and sequence >= ? order by b.sequence", from.Uint64())
}

func BatchWasExecuted(ctx context.Context, db *sql.DB, hash common.L2BatchHash) (bool, error) {
	row := db.QueryRowContext(ctx, "select is_executed from batch where is_canonical=true and hash=? ", hash.Bytes())

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

func GetTransactionsPerAddress(ctx context.Context, db *sql.DB, config *params.ChainConfig, address *gethcommon.Address, pagination *common.QueryPagination) (types.Receipts, error) {
	return selectReceipts(ctx, db, config, "where tx.sender_address = ? ORDER BY height DESC LIMIT ? OFFSET ? ", address.Bytes(), pagination.Size, pagination.Offset)
}

func CountTransactionsPerAddress(ctx context.Context, db *sql.DB, address *gethcommon.Address) (uint64, error) {
	row := db.QueryRowContext(ctx, "select count(1) from exec_tx join tx on tx.id=exec_tx.tx join batch on batch.sequence=exec_tx.batch "+" where tx.sender_address = ?", address.Bytes())

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
