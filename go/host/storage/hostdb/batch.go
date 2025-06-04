package hostdb

import (
	"database/sql"
	"errors"
	"fmt"
	"math/big"
	"strings"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/errutil"
)

const (
	selectBatch         = "SELECT sequence, hash, height, ext_batch FROM batch_host"
	selectExtBatch      = "SELECT ext_batch FROM batch_host"
	selectLatestBatch   = "SELECT sequence, hash, height, ext_batch FROM batch_host ORDER BY sequence DESC LIMIT 1"
	selectTxsAndBatch   = "SELECT t.hash FROM transaction_host t JOIN batch_host b ON t.b_sequence = b.sequence WHERE b.hash = "
	selectBatchSeqByTx  = "SELECT b_sequence FROM transaction_host WHERE hash = "
	selectTxBySeq       = "SELECT hash FROM transaction_host WHERE b_sequence = "
	selectBatchTxs      = "SELECT t.hash, b.sequence, b.height, b.ext_batch FROM transaction_host t JOIN batch_host b ON t.b_sequence = b.sequence"
	selectSumBatchSizes = "SELECT SUM(txs_size) FROM batch_host WHERE sequence >= "
)

// AddBatch adds a batch and its header to the DB
func AddBatch(dbtx *dbTransaction, statements *SQLStatements, batch *common.ExtBatch) error {
	extBatch, err := rlp.EncodeToBytes(batch)
	if err != nil {
		return fmt.Errorf("could not encode L2 transactions: %w", err)
	}

	_, err = dbtx.Tx.Exec(statements.InsertBatch,
		batch.SeqNo().Uint64(),       // sequence
		batch.Hash(),                 // full hash
		batch.Header.Number.Uint64(), // height
		extBatch,                     // ext_batch
		len(batch.EncryptedTxBlob),   // txs_size
	)
	if err != nil {
		if IsRowExistsError(err) {
			return errutil.ErrAlreadyExists
		}
		return fmt.Errorf("host failed to insert batch: %w", err)
	}

	if len(batch.TxHashes) > 0 {
		insert := statements.InsertTransactions
		args := make([]any, 0)
		for i, txHash := range batch.TxHashes {
			insert += fmt.Sprintf(" (%s, %s),", statements.GetPlaceHolder(i*2+1), statements.GetPlaceHolder(i*2+2))
			args = append(args, txHash.Bytes(), batch.SeqNo().Uint64())
		}
		insert = strings.TrimRight(insert, ",")
		_, err = dbtx.Tx.Exec(insert, args...)
		if err != nil {
			return fmt.Errorf("failed to insert transactions. cause: %w", err)
		}
	}

	var currentTotal int
	err = dbtx.Tx.QueryRow(selectTxCount).Scan(&currentTotal)
	if err != nil {
		return fmt.Errorf("failed to query transaction count: %w", err)
	}

	newTotal := currentTotal + len(batch.TxHashes)
	_, err = dbtx.Tx.Exec(statements.UpdateTxCount, newTotal)
	if err != nil {
		return fmt.Errorf("failed to update transaction count: %w", err)
	}

	return nil
}

// GetBatchListing returns latest batches given a pagination.
// For example, page 0, size 10 will return the latest 10 batches.
func GetBatchListing(db HostDB, pagination *common.QueryPagination) (*common.BatchListingResponse, error) {
	headBatch, err := GetCurrentHeadBatch(db)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch current head batch - %w", err)
	}
	batchesFrom := headBatch.SequencerOrderNo.Uint64() - pagination.Offset
	batchesTo := int(batchesFrom) - int(pagination.Size) + 1

	if batchesTo <= 0 {
		batchesTo = 1
	}

	var batches []common.PublicBatch
	for i := batchesFrom; i >= uint64(batchesTo); i-- {
		batch, err := GetPublicBatchBySequenceNumber(db, i)
		if err != nil && !errors.Is(err, errutil.ErrNotFound) {
			return nil, fmt.Errorf("failed to fetch batch by sequence number - %w", err)
		}
		if batch != nil {
			batches = append(batches, *batch)
		}
	}

	return &common.BatchListingResponse{
		BatchesData: batches,
		Total:       headBatch.SequencerOrderNo.Uint64(),
	}, nil
}

// GetPublicBatchBySequenceNumber returns the batch with the given sequence number.
func GetPublicBatchBySequenceNumber(db HostDB, seqNo uint64) (*common.PublicBatch, error) {
	whereQuery := " WHERE sequence=" + db.GetSQLStatement().Placeholder
	return fetchPublicBatch(db.GetSQLDB(), whereQuery, seqNo)
}

// GetTxsBySequenceNumber returns the transaction hashes with sequence number.
func GetTxsBySequenceNumber(db HostDB, seqNo uint64) ([]common.TxHash, error) {
	return fetchTx(db, seqNo)
}

// GetBatchBySequenceNumber returns the ext batch for a given sequence number.
func GetBatchBySequenceNumber(db HostDB, seqNo uint64) (*common.ExtBatch, error) {
	whereQuery := " WHERE sequence=" + db.GetSQLStatement().Placeholder
	return fetchFullBatch(db.GetSQLDB(), whereQuery, seqNo)
}

// GetCurrentHeadBatch retrieves the current head batch with the largest sequence number (or height).
func GetCurrentHeadBatch(db HostDB) (*common.PublicBatch, error) {
	return fetchHeadBatch(db.GetSQLDB())
}

// GetBatchHeader returns the batch header given the hash.
func GetBatchHeader(db HostDB, hash gethcommon.Hash) (*common.BatchHeader, error) {
	whereQuery := " WHERE hash=" + db.GetSQLStatement().Placeholder
	return fetchBatchHeader(db.GetSQLDB(), whereQuery, hash.Bytes())
}

// GetBatchHashByNumber returns the hash of a batch given its number.
func GetBatchHashByNumber(db HostDB, number *big.Int) (*gethcommon.Hash, error) {
	whereQuery := " WHERE height=" + db.GetSQLStatement().Placeholder
	batch, err := fetchBatchHeader(db.GetSQLDB(), whereQuery, number.Uint64())
	if err != nil {
		return nil, fmt.Errorf("failed to fetch batch header - %w", err)
	}
	l2BatchHash := batch.Hash()
	return &l2BatchHash, nil
}

// GetHeadBatchHeader returns the latest batch header.
func GetHeadBatchHeader(db HostDB) (*common.BatchHeader, error) {
	batch, err := fetchHeadBatch(db.GetSQLDB())
	if err != nil {
		return nil, fmt.Errorf("failed to fetch head batch header - %w", err)
	}
	return batch.Header, nil
}

// GetBatchNumber returns the height of the batch containing the given transaction hash.
func GetBatchNumber(db HostDB, txHash gethcommon.Hash) (*big.Int, error) {
	batchHeight, err := fetchBatchNumber(db, txHash.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to fetch batch height - %w", err)
	}
	return batchHeight, nil
}

// GetBatchTxHashes returns the transaction hashes of the batch with the given hash.
func GetBatchTxHashes(db HostDB, batchHash gethcommon.Hash) ([]gethcommon.Hash, error) {
	query := selectTxsAndBatch + db.GetSQLStatement().Placeholder
	rows, err := db.GetSQLDB().Query(query, batchHash.Bytes())
	if err != nil {
		return nil, fmt.Errorf("query execution failed: %w", err)
	}
	defer rows.Close()

	var transactions []gethcommon.Hash
	for rows.Next() {
		var txHashBytes []byte
		if err := rows.Scan(&txHashBytes); err != nil {
			return nil, fmt.Errorf("failed to scan transaction hash: %w", err)
		}
		txHash := gethcommon.BytesToHash(txHashBytes)
		transactions = append(transactions, txHash)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error looping through transacion rows: %w", err)
	}

	return transactions, nil
}

// GetPublicBatch returns the batch with the given hash.
func GetPublicBatch(db HostDB, hash common.L2BatchHash) (*common.PublicBatch, error) {
	whereQuery := " WHERE b.hash=" + db.GetSQLStatement().Placeholder
	return fetchPublicBatch(db.GetSQLDB(), whereQuery, hash.Bytes())
}

// GetBatchByTx returns the batch with the given hash.
func GetBatchByTx(db HostDB, txHash gethcommon.Hash) (*common.ExtBatch, error) {
	var seqNo uint64
	query := selectBatchSeqByTx + db.GetSQLStatement().Placeholder
	err := db.GetSQLDB().QueryRow(query, txHash.Bytes()).Scan(&seqNo)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errutil.ErrNotFound
		}
		return nil, fmt.Errorf("failed to execute query %s - %w", query, err)
	}
	return GetBatchBySequenceNumber(db, seqNo)
}

// GetBatchByHash returns the batch with the given hash.
func GetBatchByHash(db HostDB, hash common.L2BatchHash) (*common.ExtBatch, error) {
	whereQuery := " WHERE hash=" + db.GetSQLStatement().Placeholder
	return fetchFullBatch(db.GetSQLDB(), whereQuery, hash.Bytes())
}

// GetBatchHeaderByHeight returns the batch header given the height
func GetBatchHeaderByHeight(db HostDB, height *big.Int) (*common.BatchHeader, error) {
	whereQuery := " WHERE height=" + db.GetSQLStatement().Placeholder
	return fetchBatchHeader(db.GetSQLDB(), whereQuery, height.Uint64())
}

// GetBatchByHeight returns the batch header given the height
func GetBatchByHeight(db HostDB, height *big.Int) (*common.PublicBatch, error) {
	whereQuery := " WHERE height=" + db.GetSQLStatement().Placeholder
	return fetchPublicBatch(db.GetSQLDB(), whereQuery, height.Uint64())
}

// GetBatchTransactions returns the TransactionListingResponse for a given batch hash
func GetBatchTransactions(db HostDB, batchHash gethcommon.Hash, pagination *common.QueryPagination) (*common.TransactionListingResponse, error) {
	whereQuery := " WHERE b.hash=" + db.GetSQLStatement().Placeholder
	orderQuery := " ORDER BY t.timestamp DESC"
	limitQuery := fmt.Sprintf(" LIMIT %d OFFSET %d", pagination.Size, pagination.Offset)
	query := selectBatchTxs + whereQuery + orderQuery + limitQuery

	// First get total count
	countQuery := "SELECT COUNT(*) FROM transactions_host t JOIN batches b ON t.batch_hash = b.hash" + whereQuery
	var total uint64
	err := db.GetSQLDB().QueryRow(countQuery, batchHash.Bytes()).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %w", err)
	}

	rows, err := db.GetSQLDB().Query(query, batchHash.Bytes())
	if err != nil {
		return nil, fmt.Errorf("query execution for select batch transactions failed: %w", err)
	}
	defer rows.Close()

	var transactions []common.PublicTransaction
	for rows.Next() {
		var (
			txHash         gethcommon.Hash
			batchHeight    int64
			batchTimestamp uint64
			finality       string
		)
		err := rows.Scan(&txHash, &batchHeight, &batchTimestamp, &finality)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, errutil.ErrNotFound
			}
			return nil, fmt.Errorf("failed to fetch batch transactions: %w", err)
		}

		tx := common.PublicTransaction{
			TransactionHash: txHash,
			BatchHeight:     new(big.Int).SetInt64(batchHeight),
			BatchTimestamp:  batchTimestamp,
			Finality:        common.FinalityType(finality),
		}
		transactions = append(transactions, tx)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &common.TransactionListingResponse{
		TransactionsData: transactions,
		Total:            total,
	}, nil
}

func EstimateRollupSize(db HostDB, fromSeqNo *big.Int) (uint64, error) {
	var totalTx uint64
	query := selectSumBatchSizes + db.GetSQLStatement().Placeholder
	err := db.GetSQLDB().QueryRow(query, fromSeqNo.Uint64()).Scan(&totalTx)
	if err != nil {
		return 0, fmt.Errorf("failed to query sum of rollup batches: %w", err)
	}
	return totalTx, nil
}

func fetchBatchHeader(db *sql.DB, whereQuery string, args ...any) (*common.BatchHeader, error) {
	var extBatch []byte
	query := selectExtBatch + whereQuery
	var err error
	if len(args) > 0 {
		err = db.QueryRow(query, args...).Scan(&extBatch)
	} else {
		err = db.QueryRow(query).Scan(&extBatch)
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errutil.ErrNotFound
		}
		return nil, fmt.Errorf("failed to scan with query %s - %w", query, err)
	}
	// Decode batch
	var b common.ExtBatch
	err = rlp.DecodeBytes(extBatch, &b)
	if err != nil {
		return nil, fmt.Errorf("could not decode batch header. Cause: %w", err)
	}
	return b.Header, nil
}

func fetchBatchNumber(db HostDB, args ...any) (*big.Int, error) {
	var seqNo uint64
	query := selectBatchSeqByTx + db.GetSQLStatement().Placeholder
	var err error
	if len(args) > 0 {
		err = db.GetSQLDB().QueryRow(query, args...).Scan(&seqNo)
	} else {
		err = db.GetSQLDB().QueryRow(query).Scan(&seqNo)
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errutil.ErrNotFound
		}
		return nil, fmt.Errorf("failed to scan with query %s - %w", query, err)
	}
	batch, err := GetPublicBatchBySequenceNumber(db, seqNo)
	if err != nil {
		return nil, fmt.Errorf("could not fetch batch by seq no. Cause: %w", err)
	}
	return batch.Height, nil
}

func fetchPublicBatch(db *sql.DB, whereQuery string, args ...any) (*common.PublicBatch, error) {
	var sequenceInt64 uint64
	var fullHash common.TxHash
	var heightInt64 int
	var extBatch []byte

	query := selectBatch + whereQuery

	var err error
	if len(args) > 0 {
		err = db.QueryRow(query, args...).Scan(&sequenceInt64, &fullHash, &heightInt64, &extBatch)
	} else {
		err = db.QueryRow(query).Scan(&sequenceInt64, &fullHash, &heightInt64, &extBatch)
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errutil.ErrNotFound
		}
		return nil, fmt.Errorf("failed to scan with query %s - %w", query, err)
	}
	var b common.ExtBatch
	err = rlp.DecodeBytes(extBatch, &b)
	if err != nil {
		return nil, fmt.Errorf("could not decode ext batch. Cause: %w", err)
	}

	batch := &common.PublicBatch{
		SequencerOrderNo: new(big.Int).SetInt64(int64(sequenceInt64)),
		FullHash:         fullHash,
		Height:           new(big.Int).SetInt64(int64(heightInt64)),
		TxCount:          new(big.Int).SetInt64(int64(len(b.TxHashes))),
		Header:           b.Header,
		EncryptedTxBlob:  b.EncryptedTxBlob,
	}

	return batch, nil
}

func fetchFullBatch(db *sql.DB, whereQuery string, args ...any) (*common.ExtBatch, error) {
	var sequenceInt64 uint64
	var fullHash common.TxHash
	var heightInt64 int
	var extBatch []byte

	query := selectBatch + whereQuery

	var err error
	if len(args) > 0 {
		err = db.QueryRow(query, args...).Scan(&sequenceInt64, &fullHash, &heightInt64, &extBatch)
	} else {
		err = db.QueryRow(query).Scan(&sequenceInt64, &fullHash, &heightInt64, &extBatch)
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errutil.ErrNotFound
		}
		return nil, fmt.Errorf("failed to scan with query %s - %w", query, err)
	}
	var b common.ExtBatch
	err = rlp.DecodeBytes(extBatch, &b)
	if err != nil {
		return nil, fmt.Errorf("could not decode ext batch. Cause: %w", err)
	}

	return &b, nil
}

func fetchHeadBatch(db *sql.DB) (*common.PublicBatch, error) {
	var sequenceInt64 int
	var fullHash gethcommon.Hash // common.Hash
	var heightInt64 int
	var extBatch []byte

	err := db.QueryRow(selectLatestBatch).Scan(&sequenceInt64, &fullHash, &heightInt64, &extBatch)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errutil.ErrNotFound
		}
		return nil, fmt.Errorf("failed to fetch current head batch: %w", err)
	}

	var b common.ExtBatch
	err = rlp.DecodeBytes(extBatch, &b)
	if err != nil {
		return nil, fmt.Errorf("could not decode ext batch. Cause: %w", err)
	}

	batch := &common.PublicBatch{
		SequencerOrderNo: new(big.Int).SetInt64(int64(sequenceInt64)),
		FullHash:         fullHash,
		Height:           new(big.Int).SetInt64(int64(heightInt64)),
		TxCount:          new(big.Int).SetInt64(int64(len(b.TxHashes))),
		Header:           b.Header,
		EncryptedTxBlob:  b.EncryptedTxBlob,
	}

	return batch, nil
}

func fetchTx(db HostDB, seqNo uint64) ([]common.TxHash, error) {
	query := selectTxBySeq + db.GetSQLStatement().Placeholder
	rows, err := db.GetSQLDB().Query(query, seqNo)
	if err != nil {
		return nil, fmt.Errorf("query execution for select txs failed: %w", err)
	}
	defer rows.Close()

	var transactions []gethcommon.Hash
	for rows.Next() {
		var txHashBytes []byte
		if err := rows.Scan(&txHashBytes); err != nil {
			return nil, fmt.Errorf("failed to scan transaction hash: %w", err)
		}
		txHash := gethcommon.BytesToHash(txHashBytes)
		transactions = append(transactions, txHash)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error looping through transacion rows: %w", err)
	}

	return transactions, nil
}

func fetchBatchTxs(db *sql.DB, whereQuery string, batchHash []byte) (*common.TransactionListingResponse, error) {
	query := selectBatchTxs + whereQuery
	rows, err := db.Query(query, batchHash)
	if err != nil {
		return nil, fmt.Errorf("query execution for select batch txs failed: %w", err)
	}
	defer rows.Close()

	var transactions []common.PublicTransaction
	for rows.Next() {
		var (
			fullHash []byte
			sequence int
			height   int
			extBatch []byte
		)
		err := rows.Scan(&fullHash, &sequence, &height, &extBatch)
		if err != nil {
			return nil, fmt.Errorf("failed to scan with query %s - %w", query, err)
		}
		extBatchDecoded := new(common.ExtBatch)
		if err := rlp.DecodeBytes(extBatch, extBatchDecoded); err != nil {
			return nil, fmt.Errorf("could not decode batch. Cause: %w", err)
		}
		transaction := common.PublicTransaction{
			TransactionHash: gethcommon.BytesToHash(fullHash),
			BatchHeight:     big.NewInt(int64(height)),
			BatchTimestamp:  extBatchDecoded.Header.Time,
			Finality:        common.BatchFinal,
		}
		transactions = append(transactions, transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error looping through transacion rows: %w", err)
	}

	return &common.TransactionListingResponse{
		TransactionsData: transactions,
		Total:            uint64(len(transactions)),
	}, nil
}

func IsRowExistsError(err error) bool {
	return strings.Contains(strings.ToLower(err.Error()), "unique") || strings.Contains(strings.ToLower(err.Error()), "duplicate key")
}
