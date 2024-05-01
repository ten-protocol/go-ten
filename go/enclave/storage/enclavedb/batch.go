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

	queryReceipts = "select exec_tx.receipt, tx.content, batch.full_hash, batch.height from exec_tx join tx on tx.id=exec_tx.tx join batch on batch.sequence=exec_tx.batch "
)

// WriteBatchAndTransactions - persists the batch and the transactions
func WriteBatchAndTransactions(ctx context.Context, dbtx DBTransaction, batch *core.Batch, convertedHash gethcommon.Hash, blockId *uint64) error {
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

	dbtx.ExecuteSQL("replace into batch_body values (?,?)", batchBodyID, body)

	var isCanon bool
	err = dbtx.GetDB().QueryRowContext(ctx,
		"select is_canonical from block where hash=? and full_hash=?",
		truncTo4(batch.Header.L1Proof), batch.Header.L1Proof.Bytes(),
	).Scan(&isCanon)
	if err != nil {
		// if the block is not found, we assume it is non-canonical
		// fmt.Printf("IsCanon %s err: %s\n", batch.Header.L1Proof, err)
		isCanon = false
	}

	dbtx.ExecuteSQL("insert into batch values (?,?,?,?,?,?,?,?,?,?)",
		batch.Header.SequencerOrderNo.Uint64(), // sequence
		batch.Hash(),                           // full hash
		convertedHash,                          // converted_hash
		truncTo4(batch.Hash()),                 // index hash
		batch.Header.Number.Uint64(),           // height
		isCanon,                                // is_canonical
		header,                                 // header blob
		batchBodyID,                            // reference to the batch body
		blockId,                                // indexed l1_proof
		false,                                  // executed
	)

	// creates a big insert statement for all transactions
	if len(batch.Transactions) > 0 {
		insert := "replace into tx (hash, full_hash, content, sender_address, nonce, idx, body) values " + repeat("(?,?,?,?,?,?,?)", ",", len(batch.Transactions))

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

			args = append(args, truncTo4(transaction.Hash())) // truncated tx_hash
			args = append(args, transaction.Hash())           // full tx_hash
			args = append(args, txBytes)                      // content
			args = append(args, from.Bytes())                 // sender_address
			args = append(args, transaction.Nonce())          // nonce
			args = append(args, i)                            // idx
			args = append(args, batchBodyID)                  // the batch body which contained it
		}
		dbtx.ExecuteSQL(insert, args...)
	}

	return nil
}

// WriteBatchExecution - insert all receipts to the db
func WriteBatchExecution(ctx context.Context, dbtx DBTransaction, seqNo *big.Int, receipts []*types.Receipt) error {
	dbtx.ExecuteSQL("update batch set is_executed=true where sequence=?", seqNo.Uint64())

	args := make([]any, 0)
	for _, receipt := range receipts {
		// Convert the receipt into their storage form and serialize them
		storageReceipt := (*types.ReceiptForStorage)(receipt)
		receiptBytes, err := rlp.EncodeToBytes(storageReceipt)
		if err != nil {
			return fmt.Errorf("failed to encode block receipts. Cause: %w", err)
		}

		// ignore the error because synthetic transactions will not be inserted
		txId, _ := ReadTxId(ctx, dbtx, storageReceipt.TxHash)
		args = append(args, truncBTo4(receipt.ContractAddress.Bytes())) // created_contract_address
		args = append(args, receipt.ContractAddress.Bytes())            // created_contract_address
		args = append(args, receiptBytes)                               // the serialised receipt
		args = append(args, txId)                                       // tx id
		args = append(args, seqNo.Uint64())                             // batch_seq
	}
	if len(args) > 0 {
		insert := "insert into exec_tx (created_contract_address,created_contract_address_full, receipt, tx, batch) values " + repeat("(?,?,?,?,?)", ",", len(receipts))
		dbtx.ExecuteSQL(insert, args...)
	}
	return nil
}

func ReadTxId(ctx context.Context, dbtx DBTransaction, txHash gethcommon.Hash) (*uint64, error) {
	var txId uint64
	err := dbtx.GetDB().QueryRowContext(ctx, "select id from tx where hash=? and full_hash=?", truncTo4(txHash), txHash.Bytes()).Scan(&txId)
	if err != nil {
		return nil, err
	}
	return &txId, err
}

func ReadBatchBySeqNo(ctx context.Context, db *sql.DB, seqNo uint64) (*core.Batch, error) {
	return fetchBatch(ctx, db, " where sequence=?", seqNo)
}

func ReadBatchByHash(ctx context.Context, db *sql.DB, hash common.L2BatchHash) (*core.Batch, error) {
	return fetchBatch(ctx, db, " where b.hash=? and b.full_hash=?", truncTo4(hash), hash.Bytes())
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
	return fetchBatches(ctx, db, " join block l1b on b.l1_proof=l1b.id where l1b.hash=? and l1b.full_l1_proof=? order by b.sequence", truncTo4(hash), hash.Bytes())
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

func ReadHeadBatchForBlock(ctx context.Context, db *sql.DB, l1Hash common.L1BlockHash) (*core.Batch, error) {
	query := " where b.is_canonical=true and b.is_executed=true and b.height=(select max(b1.height) from batch b1 where b1.is_canonical=true and b1.is_executed=true and b1.l1_proof=? and b1.full_l1_proof=?)"
	return fetchBatch(ctx, db, query, truncTo4(l1Hash), l1Hash.Bytes())
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
	return selectReceipts(ctx, db, config, "where batch.hash=? and batch.full_hash=?", truncTo4(hash), hash.Bytes())
}

func ReadReceipt(ctx context.Context, db *sql.DB, txHash common.L2TxHash, config *params.ChainConfig) (*types.Receipt, error) {
	// todo - canonical?
	row := db.QueryRowContext(ctx, queryReceipts+" where tx.hash=? and tx.full_hash=?", truncTo4(txHash), txHash.Bytes())
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
		"select tx.content, batch.full_hash, batch.height, tx.idx from exec_tx join tx on tx.id=exec_tx.tx join batch on batch.sequence=exec_tx.batch where batch.is_canonical=true and tx.hash=? and tx.full_hash=?",
		truncTo4(txHash), txHash.Bytes())

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
	row := db.QueryRowContext(ctx, "select tx.full_hash from exec_tx join tx on tx.id=exec_tx.tx where created_contract_address=? and created_contract_address_full=?", truncBTo4(address.Bytes()), address.Bytes())

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
	row := db.QueryRowContext(ctx, "select is_executed from batch where is_canonical=true and hash=? and full_hash=?", truncTo4(hash), hash.Bytes())

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

func GetReceiptsPerAddress(ctx context.Context, db *sql.DB, config *params.ChainConfig, address *gethcommon.Address, pagination *common.QueryPagination) (types.Receipts, error) {
	// todo - not indexed
	return selectReceipts(ctx, db, config, "where tx.sender_address = ? ORDER BY height DESC LIMIT ? OFFSET ? ", address.Bytes(), pagination.Size, pagination.Offset)
}

func GetReceiptsPerAddressCount(ctx context.Context, db *sql.DB, address *gethcommon.Address) (uint64, error) {
	// todo - this is not indexed and will do a full table scan!
	row := db.QueryRowContext(ctx, "select count(1) from exec_tx join tx on tx.id=exec_tx.tx join batch on batch.sequence=exec_tx.batch "+" where tx.sender_address = ?", address.Bytes())

	var count uint64
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func GetPublicTransactionData(ctx context.Context, db *sql.DB, pagination *common.QueryPagination) ([]common.PublicTransaction, error) {
	return selectPublicTxsBySender(ctx, db, " ORDER BY height DESC LIMIT ? OFFSET ? ", pagination.Size, pagination.Offset)
}

func selectPublicTxsBySender(ctx context.Context, db *sql.DB, query string, args ...any) ([]common.PublicTransaction, error) {
	var publicTxs []common.PublicTransaction

	q := "select tx.full_hash, batch.height, batch.header from exec_tx join batch on batch.sequence=exec_tx.batch join tx on tx.id=exec_tx.tx where batch.is_canonical=true " + query
	rows, err := db.QueryContext(ctx, q, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// make sure the error is converted to obscuro-wide not found error
			return nil, errutil.ErrNotFound
		}
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var txHash []byte
		var batchHeight uint64
		var batchHeader string
		err := rows.Scan(&txHash, &batchHeight, &batchHeader)
		if err != nil {
			return nil, err
		}

		h := new(common.BatchHeader)
		if err := rlp.DecodeBytes([]byte(batchHeader), h); err != nil {
			return nil, fmt.Errorf("could not decode batch header. Cause: %w", err)
		}

		publicTxs = append(publicTxs, common.PublicTransaction{
			TransactionHash: gethcommon.BytesToHash(txHash),
			BatchHeight:     big.NewInt(0).SetUint64(batchHeight),
			BatchTimestamp:  h.Time,
			Finality:        common.BatchFinal,
		})
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return publicTxs, nil
}

func GetPublicTransactionCount(ctx context.Context, db *sql.DB) (uint64, error) {
	row := db.QueryRowContext(ctx, "select count(1) from exec_tx join batch on batch.sequence=exec_tx.batch where batch.is_canonical=true")

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
