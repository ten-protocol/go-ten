package enclavedb

import (
	"bytes"
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/params"

	"github.com/ethereum/go-ethereum/core/types"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/errutil"
	"github.com/obscuronet/go-obscuro/go/enclave/core"
)

const (
	bodyInsert    = "replace into batch_body values (?,?)"
	txInsert      = "replace into tx values "
	txInsertValue = "(?,?,?,?,?,?,?)"

	bInsert             = "insert into batch values (?,?,?,?,?,?,?,?,?,?)"
	updateBatchExecuted = "update batch set is_executed=true where sequence=?"

	selectBatch  = "select b.header, bb.content from batch b join batch_body bb on b.body=bb.id"
	selectHeader = "select b.header from batch b"

	txExecInsert       = "insert into exec_tx values "
	txExecInsertValue  = "(?,?,?,?,?)"
	queryReceipts      = "select exec_tx.receipt, tx.content, batch.full_hash, batch.height from exec_tx join tx on tx.hash=exec_tx.tx join batch on batch.sequence=exec_tx.batch "
	queryReceiptsCount = "select count(1) from exec_tx join tx on tx.hash=exec_tx.tx join batch on batch.sequence=exec_tx.batch "

	selectTxQuery = "select tx.content, batch.full_hash, batch.height, tx.idx from exec_tx join tx on tx.hash=exec_tx.tx join batch on batch.sequence=exec_tx.batch where batch.is_canonical=true and tx.hash=?"

	selectContractCreationTx    = "select tx.full_hash from exec_tx join tx on tx.hash=exec_tx.tx where created_contract_address=?"
	selectTotalCreatedContracts = "select count( distinct created_contract_address) from exec_tx "
	queryBatchWasExecuted       = "select is_executed from batch where is_canonical=true and hash=?"

	isCanonQuery = "select is_canonical from block where hash=?"

	queryTxList      = "select tx.full_hash, batch.height from exec_tx join batch on batch.sequence=exec_tx.batch join tx on tx.hash=exec_tx.tx"
	queryTxCountList = "select count(1) from exec_tx join batch on batch.sequence=exec_tx.batch"
)

// WriteBatchAndTransactions - persists the batch and the transactions
func WriteBatchAndTransactions(dbtx DBTransaction, batch *core.Batch) error {
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

	dbtx.ExecuteSQL(bodyInsert, batchBodyID, body)

	var parentBytes []byte
	if batch.Number().Uint64() > 0 {
		parentBytes = truncTo16(batch.Header.ParentHash)
	}

	var isCanon bool
	err = dbtx.GetDB().QueryRow(isCanonQuery, truncTo16(batch.Header.L1Proof)).Scan(&isCanon)
	if err != nil {
		// if the block is not found, we assume it is non-canonical
		// fmt.Printf("IsCanon %s err: %s\n", batch.Header.L1Proof, err)
		isCanon = false
	}

	dbtx.ExecuteSQL(bInsert,
		batch.Header.SequencerOrderNo.Uint64(), // sequence
		batch.Hash(),                           // full hash
		truncTo16(batch.Hash()),                // index hash
		parentBytes,                            // parent
		batch.Header.Number.Uint64(),           // height
		isCanon,                                // is_canonical
		header,                                 // header blob
		batchBodyID,                            // reference to the batch body
		truncTo16(batch.Header.L1Proof),        // l1_proof
		false,                                  // executed
	)

	// creates a big insert statement for all transactions
	if len(batch.Transactions) > 0 {
		insert := txInsert + strings.Repeat(txInsertValue+",", len(batch.Transactions))
		insert = insert[0 : len(insert)-1] // remove trailing comma

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

			args = append(args, truncTo16(transaction.Hash())) // indexed tx_hash
			args = append(args, transaction.Hash())            // full tx_hash
			args = append(args, txBytes)                       // content
			args = append(args, from.Bytes())                  // sender_address
			args = append(args, transaction.Nonce())           // nonce
			args = append(args, i)                             // idx
			args = append(args, batchBodyID)                   // the batch body which contained it
		}
		dbtx.ExecuteSQL(insert, args...)
	}

	return nil
}

// WriteBatchExecution - insert all receipts to the db
func WriteBatchExecution(dbtx DBTransaction, seqNo *big.Int, receipts []*types.Receipt) error {
	dbtx.ExecuteSQL(updateBatchExecuted, seqNo.Uint64())

	args := make([]any, 0)
	for _, receipt := range receipts {
		// Convert the receipt into their storage form and serialize them
		storageReceipt := (*types.ReceiptForStorage)(receipt)
		receiptBytes, err := rlp.EncodeToBytes(storageReceipt)
		if err != nil {
			return fmt.Errorf("failed to encode block receipts. Cause: %w", err)
		}

		args = append(args, executedTransactionID(&receipt.BlockHash, &receipt.TxHash)) // PK
		args = append(args, receipt.ContractAddress.Bytes())                            // created_contract_address
		args = append(args, receiptBytes)                                               // the serialised receipt
		args = append(args, truncTo16(receipt.TxHash))                                  // tx_hash
		args = append(args, seqNo.Uint64())                                             // batch_seq
	}
	if len(args) > 0 {
		insert := txExecInsert + strings.Repeat(txExecInsertValue+",", len(receipts))
		insert = insert[0 : len(insert)-1] // remove trailing comma
		dbtx.ExecuteSQL(insert, args...)
	}
	return nil
}

// concatenates the batch_hash with the tx_hash to create a PK for the executed transaction
func executedTransactionID(batchHash *common.L2BatchHash, txHash *common.L2TxHash) []byte {
	execTxID := make([]byte, 0)
	execTxID = append(execTxID, batchHash.Bytes()...)
	execTxID = append(execTxID, txHash.Bytes()...)
	return truncTo16(sha256.Sum256(execTxID))
}

func ReadBatchBySeqNo(db *sql.DB, seqNo uint64) (*core.Batch, error) {
	return fetchBatch(db, " where sequence=?", seqNo)
}

func ReadBatchByHash(db *sql.DB, hash common.L2BatchHash) (*core.Batch, error) {
	return fetchBatch(db, " where b.hash=?", truncTo16(hash))
}

func ReadCanonicalBatchByHeight(db *sql.DB, height uint64) (*core.Batch, error) {
	return fetchBatch(db, " where b.height=? and is_canonical=true", height)
}

func ReadBatchHeader(db *sql.DB, hash gethcommon.Hash) (*common.BatchHeader, error) {
	return fetchBatchHeader(db, " where hash=?", truncTo16(hash))
}

// todo - is there a better way to write this query?
func ReadCurrentHeadBatch(db *sql.DB) (*core.Batch, error) {
	return fetchBatch(db, " where b.is_canonical=true and b.is_executed=true and b.height=(select max(b1.height) from batch b1 where b1.is_canonical=true and b1.is_executed=true)")
}

func ReadBatchesByBlock(db *sql.DB, hash common.L1BlockHash) ([]*core.Batch, error) {
	return fetchBatches(db, " where b.l1_proof=? order by b.sequence", truncTo16(hash))
}

func ReadCurrentSequencerNo(db *sql.DB) (*big.Int, error) {
	var seq sql.NullInt64
	query := "select max(sequence) from batch"
	err := db.QueryRow(query).Scan(&seq)
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

func ReadHeadBatchForBlock(db *sql.DB, l1Hash common.L1BlockHash) (*core.Batch, error) {
	query := " where b.is_canonical=true and b.is_executed=true and b.height=(select max(b1.height) from batch b1 where b1.is_canonical=true and b1.is_executed=true and b1.l1_proof=?)"
	return fetchBatch(db, query, truncTo16(l1Hash))
}

func fetchBatch(db *sql.DB, whereQuery string, args ...any) (*core.Batch, error) {
	var header string
	var body []byte
	query := selectBatch + " " + whereQuery
	var err error
	if len(args) > 0 {
		err = db.QueryRow(query, args...).Scan(&header, &body)
	} else {
		err = db.QueryRow(query).Scan(&header, &body)
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

func fetchBatches(db *sql.DB, whereQuery string, args ...any) ([]*core.Batch, error) {
	result := make([]*core.Batch, 0)

	rows, err := db.Query(selectBatch+" "+whereQuery, args...)
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

func fetchBatchHeader(db *sql.DB, whereQuery string, args ...any) (*common.BatchHeader, error) {
	var header string
	query := selectHeader + " " + whereQuery
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
	h := new(common.BatchHeader)
	if err := rlp.Decode(bytes.NewReader([]byte(header)), h); err != nil {
		return nil, fmt.Errorf("could not decode batch header. Cause: %w", err)
	}

	return h, nil
}

func selectReceipts(db *sql.DB, config *params.ChainConfig, query string, args ...any) (types.Receipts, error) {
	var allReceipts types.Receipts

	// where batch=?
	rows, err := db.Query(queryReceipts+" "+query, args...)
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
func ReadReceiptsByBatchHash(db *sql.DB, hash common.L2BatchHash, config *params.ChainConfig) (types.Receipts, error) {
	return selectReceipts(db, config, "where batch.hash = ?", truncTo16(hash))
}

func ReadReceipt(db *sql.DB, hash common.L2TxHash, config *params.ChainConfig) (*types.Receipt, error) {
	row := db.QueryRow(queryReceipts+" where tx=?", truncTo16(hash))
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
	if err = receipts.DeriveFields(config, batchhash, height, 0, big.NewInt(0), big.NewInt(0), transactions); err != nil {
		return nil, fmt.Errorf("failed to derive block receipts fields. hash = %s; number = %d; err = %w", hash, height, err)
	}
	return receipts[0], nil
}

func ReadTransaction(db *sql.DB, txHash gethcommon.Hash) (*types.Transaction, gethcommon.Hash, uint64, uint64, error) {
	row := db.QueryRow(selectTxQuery, truncTo16(txHash))

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

func GetContractCreationTx(db *sql.DB, address gethcommon.Address) (*gethcommon.Hash, error) {
	row := db.QueryRow(selectContractCreationTx, address.Bytes())

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

func ReadContractCreationCount(db *sql.DB) (*big.Int, error) {
	row := db.QueryRow(selectTotalCreatedContracts)

	var count int64
	err := row.Scan(&count)
	if err != nil {
		return nil, err
	}

	return big.NewInt(count), nil
}

func ReadUnexecutedBatches(db *sql.DB, from *big.Int) ([]*core.Batch, error) {
	return fetchBatches(db, "where is_executed=false and is_canonical=true and sequence >= ? order by b.sequence", from.Uint64())
}

func BatchWasExecuted(db *sql.DB, hash common.L2BatchHash) (bool, error) {
	row := db.QueryRow(queryBatchWasExecuted, truncTo16(hash))

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

func GetReceiptsPerAddress(db *sql.DB, config *params.ChainConfig, address *gethcommon.Address, pagination *common.QueryPagination) (types.Receipts, error) {
	return selectReceipts(db, config, "where tx.sender_address = ? ORDER BY height DESC LIMIT ? OFFSET ? ", address.Bytes(), pagination.Size, pagination.Offset)
}

func GetReceiptsPerAddressCount(db *sql.DB, address *gethcommon.Address) (uint64, error) {
	row := db.QueryRow(queryReceiptsCount+" where tx.sender_address = ?", address.Bytes())

	var count uint64
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func GetPublicTransactionData(db *sql.DB, pagination *common.QueryPagination) ([]common.PublicTransaction, error) {
	return selectPublicTxsBySender(db, " ORDER BY height DESC LIMIT ? OFFSET ? ", pagination.Size, pagination.Offset)
}

func selectPublicTxsBySender(db *sql.DB, query string, args ...any) ([]common.PublicTransaction, error) {
	var publicTxs []common.PublicTransaction

	rows, err := db.Query(queryTxList+" "+query, args...)
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
		err := rows.Scan(&txHash, &batchHeight)
		if err != nil {
			return nil, err
		}

		publicTxs = append(publicTxs, common.PublicTransaction{
			TransactionHash: gethcommon.BytesToHash(txHash),
			BatchHeight:     big.NewInt(0).SetUint64(batchHeight),
			Finality:        common.BatchFinal,
		})
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return publicTxs, nil
}

func GetPublicTransactionCount(db *sql.DB) (uint64, error) {
	row := db.QueryRow(queryTxCountList)

	var count uint64
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
