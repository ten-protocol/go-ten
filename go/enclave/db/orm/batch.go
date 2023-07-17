package orm

import (
	"bytes"
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
	obscurosql "github.com/obscuronet/go-obscuro/go/enclave/db/sql"
)

const (
	bodyInsert = "replace into batch_body values (?,?)"
	bInsert    = "insert into batch values (?,?,?,?,?,?,?,?,?)"

	txInsert      = "insert into tx values "
	txInsertValue = "(?,?,?,?,?,?)"

	selectBatch  = "select b.header, bb.content from batch b join batch_body bb on b.body=bb.hash"
	selectHeader = "select b.header from batch b"

	txExecInsert      = "insert into exec_tx values "
	txExecInsertValue = "(?,?,?,?,?)"
	queryReceipts     = "select exec_tx.receipt, tx.content, exec_tx.batch, batch.height from exec_tx join tx on tx.hash=exec_tx.tx join batch on batch.hash=exec_tx.batch "

	selectTxQuery = "select tx.content, exec_tx.batch, batch.height, tx.idx from exec_tx join tx on tx.hash=exec_tx.tx join batch on batch.hash=exec_tx.batch where tx.hash=?"
)

func WriteBatch(dbtx *obscurosql.Batch, batch *core.Batch) error {
	bodyHash := batch.Header.TxHash.Bytes()

	body, err := rlp.EncodeToBytes(batch.Transactions)
	if err != nil {
		return fmt.Errorf("could not encode L2 transactions. Cause: %w", err)
	}
	header, err := rlp.EncodeToBytes(batch.Header)
	if err != nil {
		return fmt.Errorf("could not encode batch header. Cause: %w", err)
	}

	dbtx.ExecuteSQL(bodyInsert, bodyHash, body)

	var parentBytes []byte = nil
	if batch.Number().Uint64() > 0 {
		parentBytes = batch.Header.ParentHash.Bytes()
	}

	dbtx.ExecuteSQL(bInsert,
		batch.Hash().Bytes(),
		parentBytes,
		batch.Header.SequencerOrderNo.Uint64(),
		batch.Header.Number.Uint64(),
		true,
		string(header),
		bodyHash,
		batch.Header.L1Proof.Bytes(),
		"", // todo
	)

	if len(batch.Transactions) > 0 {
		insert := txInsert + strings.Repeat(txInsertValue+",", len(batch.Transactions))
		args := make([]any, 0)
		for i, transaction := range batch.Transactions {
			txBytes, err := rlp.EncodeToBytes(transaction)
			if err != nil {
				return fmt.Errorf("failed to encode block receipts. Cause: %w", err)
			}

			args = append(args, transaction.Hash().Bytes())
			args = append(args, txBytes)
			args = append(args, nil)
			args = append(args, transaction.Nonce())
			args = append(args, i)
			args = append(args, bodyHash)
		}
		dbtx.ExecuteSQL(insert[0:len(insert)-1], args...)
	}

	return nil
}

func FindBatchBySeqNo(db *sql.DB, seqNo uint64) (*core.Batch, error) {
	return fetchBatch(db, " where sequence=?", seqNo)
}

func FetchBatch(db *sql.DB, hash common.L2BatchHash) (*core.Batch, error) {
	return fetchBatch(db, " where b.hash=?", hash.Bytes())
}

func FetchCanonicalBatchByHeight(db *sql.DB, height uint64) (*core.Batch, error) {
	return fetchBatch(db, " where b.height=? and is_canonical", height)
}

func ReadBatchHeader(db *sql.DB, hash gethcommon.Hash) (*common.BatchHeader, error) {
	return fetchBatchHeader(db, " where hash=?", hash.Bytes())
}

func FetchHeadBatch(db *sql.DB) (*core.Batch, error) {
	return fetchBatch(db, " where b.height=(select max(b1.height) from batch b1) and is_canonical")
}

func FetchCurrentSequencerNo(db *sql.DB) (*big.Int, error) {
	var seq int64
	query := "select max(sequence) from batch"
	err := db.QueryRow(query).Scan(&seq)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// make sure the error is converted to obscuro-wide not found error
			return nil, errutil.ErrNotFound
		}
		return nil, err
	}
	return big.NewInt(seq), nil
}

func FetchHeadBatchForBlock(db *sql.DB, l1Hash common.L1BlockHash) (*core.Batch, error) {
	query := " where is_canonical and b.height=(select max(b1.height) from batch b1 where b1.is_canonical and b1.l1_proof=?)"

	return fetchBatch(db, query, l1Hash.Bytes())
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

	return &core.Batch{
		Header:       h,
		Transactions: *txs,
	}, nil
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

func WriteReceipts(dbtx *obscurosql.Batch, receipts []*types.Receipt) error {
	var args []any
	for _, receipt := range receipts {
		// Convert the receipt into their storage form and serialize them
		storageReceipt := (*types.ReceiptForStorage)(receipt)
		receiptBytes, err := rlp.EncodeToBytes(storageReceipt)
		if err != nil {
			return fmt.Errorf("failed to encode block receipts. Cause: %w", err)
		}

		execTxId := make([]byte, 0)
		execTxId = append(execTxId, receipt.BlockHash.Bytes()...)
		execTxId = append(execTxId, receipt.TxHash.Bytes()...)
		// println("rec: " + string(execTxId))

		args = append(args, execTxId)
		args = append(args, receipt.ContractAddress.Bytes())
		args = append(args, receiptBytes)
		args = append(args, receipt.TxHash.Bytes())
		args = append(args, receipt.BlockHash.Bytes())
	}
	if len(args) > 0 {
		insert := txExecInsert + strings.Repeat(txExecInsertValue+",", len(receipts))
		dbtx.ExecuteSQL(insert[0:len(insert)-1], args...)
	}
	return nil
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
		if err = receipts.DeriveFields(config, hash, height, transactions); err != nil {
			return nil, fmt.Errorf("failed to derive block receipts fields. hash = %s; number = %d; err = %w", hash, height, err)
		}
		allReceipts = append(allReceipts, receipts[0])
	}

	return allReceipts, nil
}

// ReadReceipts retrieves all the transaction receipts belonging to a block, including
// its corresponding metadata fields. If it is unable to populate these metadata
// fields then nil is returned.
//
// The current implementation populates these metadata fields by reading the receipts'
// corresponding block body, so if the block body is not found it will return nil even
// if the receipt itself is stored.
func ReadReceipts(db *sql.DB, hash common.L2BatchHash, config *params.ChainConfig) (types.Receipts, error) {
	return selectReceipts(db, config, "where batch.hash = ?", hash.Bytes())
}

func ReadReceipt(db *sql.DB, hash common.L2TxHash, config *params.ChainConfig) (*types.Receipt, error) {
	row := db.QueryRow(queryReceipts+" where tx=?", hash.Bytes())
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
	if err = receipts.DeriveFields(config, batchhash, height, transactions); err != nil {
		return nil, fmt.Errorf("failed to derive block receipts fields. hash = %s; number = %d; err = %w", hash, height, err)
	}
	return receipts[0], nil
}

func ReadTransaction(db *sql.DB, txHash gethcommon.Hash) (*types.Transaction, gethcommon.Hash, uint64, uint64, error) {
	row := db.QueryRow(selectTxQuery, txHash.Bytes())

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
