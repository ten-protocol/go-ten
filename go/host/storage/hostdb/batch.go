package hostdb

import (
	"database/sql"
	"errors"
	"fmt"
	"math/big"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/errutil"
)

const (
	selectTxCount      = "SELECT total FROM transaction_count WHERE id = 1"
	selectBatch        = "SELECT sequence, full_hash, hash, height, ext_batch FROM batch_host"
	selectExtBatch     = "SELECT ext_batch FROM batch_host"
	selectLatestBatch  = "SELECT sequence, full_hash, hash, height, ext_batch FROM batch_host ORDER BY sequence DESC LIMIT 1"
	selectTxsAndBatch  = "SELECT t.hash FROM transactions_host t JOIN batch_host b ON t.b_sequence = b.sequence WHERE b.full_hash = "
	selectBatchSeqByTx = "SELECT b_sequence FROM transactions_host WHERE hash = "
	selectTxBySeq      = "SELECT hash FROM transactions_host WHERE b_sequence = "
	selectBatchTxs     = "SELECT t.hash, b.sequence, b.height, b.ext_batch FROM transactions_host t JOIN batch_host b ON t.b_sequence = b.sequence"
)

// AddBatch adds a batch and its header to the DB
func AddBatch(dbtx *dbTransaction, statements *SQLStatements, batch *common.ExtBatch) error {
	extBatch, err := rlp.EncodeToBytes(batch)
	if err != nil {
		return fmt.Errorf("could not encode L2 transactions: %w", err)
	}

	_, err = dbtx.tx.Exec(statements.InsertBatch,
		batch.SeqNo().Uint64(),       // sequence
		batch.Hash(),                 // full hash
		truncTo16(batch.Hash()),      // shortened hash
		batch.Header.Number.Uint64(), // height
		extBatch,                     // ext_batch
	)
	if err != nil {
		return fmt.Errorf("host failed to insert batch: %w", err)
	}

	if len(batch.TxHashes) > 0 {
		for _, transaction := range batch.TxHashes {
			_, err = dbtx.tx.Exec(statements.InsertTransactions, transaction.Bytes(), batch.SeqNo().Uint64())
			if err != nil {
				return fmt.Errorf("failed to insert transaction with hash: %d", err)
			}
		}
	}

	var currentTotal int
	err = dbtx.tx.QueryRow(selectTxCount).Scan(&currentTotal)
	if err != nil {
		return fmt.Errorf("failed to query transaction count: %w", err)
	}

	newTotal := currentTotal + len(batch.TxHashes)
	_, err = dbtx.tx.Exec(statements.InsertTxCount, 1, newTotal)
	if err != nil {
		return fmt.Errorf("failed to update transaction count: %w", err)
	}

	return nil
}

// GetBatchListing returns latest batches given a pagination.
// For example, page 0, size 10 will return the latest 10 batches.
func GetBatchListing(db HostDB, pagination *common.QueryPagination) (*common.BatchListingResponse, error) {
	headBatch, err := GetCurrentHeadBatch(db.GetSQLDB())
	if err != nil {
		return nil, err
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
			return nil, err
		}
		if batch != nil {
			batches = append(batches, *batch)
		}
	}

	return &common.BatchListingResponse{
		BatchesData: batches,
		Total:       uint64(len(batches)),
	}, nil
}

// GetBatchListingDeprecated returns latest batches given a pagination.
// For example, page 0, size 10 will return the latest 10 batches.
func GetBatchListingDeprecated(db HostDB, pagination *common.QueryPagination) (*common.BatchListingResponseDeprecated, error) {
	headBatch, err := GetCurrentHeadBatch(db.GetSQLDB())
	if err != nil {
		return nil, err
	}
	batchesFrom := headBatch.SequencerOrderNo.Uint64() - pagination.Offset
	batchesTo := int(batchesFrom) - int(pagination.Size) + 1

	if batchesTo <= 0 {
		batchesTo = 1
	}

	var batches []common.PublicBatchDeprecated
	var txHashes []common.TxHash
	for i := batchesFrom; i >= uint64(batchesTo); i-- {
		batch, err := GetPublicBatchBySequenceNumber(db, i)
		if batch == nil {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("failed to get batch by seq no: %w", err)
		}

		txHashes, err = GetTxsBySequenceNumber(db, batch.Header.SequencerOrderNo.Uint64())
		if err != nil {
			return nil, fmt.Errorf("failed to get tx hashes by seq no: %w", err)
		}
		if batch == nil || batch.Header == nil {
			return nil, fmt.Errorf("batch or batch header is nil")
		} else {
			publicBatchDeprecated := common.PublicBatchDeprecated{
				BatchHeader: *batch.Header,
				TxHashes:    txHashes,
			}
			batches = append(batches, publicBatchDeprecated)
		}
	}

	return &common.BatchListingResponseDeprecated{
		BatchesData: batches,
		Total:       uint64(len(batches)),
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
func GetCurrentHeadBatch(db *sql.DB) (*common.PublicBatch, error) {
	return fetchHeadBatch(db)
}

// GetBatchHeader returns the batch header given the hash.
func GetBatchHeader(db HostDB, hash gethcommon.Hash) (*common.BatchHeader, error) {
	whereQuery := " WHERE hash=" + db.GetSQLStatement().Placeholder
	return fetchBatchHeader(db.GetSQLDB(), whereQuery, truncTo16(hash))
}

// GetBatchHashByNumber returns the hash of a batch given its number.
func GetBatchHashByNumber(db HostDB, number *big.Int) (*gethcommon.Hash, error) {
	whereQuery := " WHERE height=" + db.GetSQLStatement().Placeholder
	batch, err := fetchBatchHeader(db.GetSQLDB(), whereQuery, number.Uint64())
	if err != nil {
		return nil, err
	}
	l2BatchHash := batch.Hash()
	return &l2BatchHash, nil
}

// GetHeadBatchHeader returns the latest batch header.
func GetHeadBatchHeader(db *sql.DB) (*common.BatchHeader, error) {
	batch, err := fetchHeadBatch(db)
	if err != nil {
		return nil, err
	}
	return batch.Header, nil
}

// GetBatchNumber returns the height of the batch containing the given transaction hash.
func GetBatchNumber(db HostDB, txHash gethcommon.Hash) (*big.Int, error) {
	txBytes := txHash.Bytes()
	batchHeight, err := fetchBatchNumber(db, txBytes)
	if err != nil {
		return nil, err
	}
	return batchHeight, nil
}

// GetBatchTxHashes returns the transaction hashes of the batch with the given hash.
func GetBatchTxHashes(db HostDB, batchHash gethcommon.Hash) ([]gethcommon.Hash, error) {
	query := selectTxsAndBatch + db.GetSQLStatement().Placeholder
	rows, err := db.GetSQLDB().Query(query, batchHash)
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

// GetTotalTxCount returns the total number of batched transactions.
func GetTotalTxCount(db *sql.DB) (*big.Int, error) {
	var totalCount int
	err := db.QueryRow(selectTxCount).Scan(&totalCount)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve total transaction count: %w", err)
	}
	return big.NewInt(int64(totalCount)), nil
}

// GetPublicBatch returns the batch with the given hash.
func GetPublicBatch(db HostDB, hash common.L2BatchHash) (*common.PublicBatch, error) {
	whereQuery := " WHERE b.hash=" + db.GetSQLStatement().Placeholder
	return fetchPublicBatch(db.GetSQLDB(), whereQuery, truncTo16(hash))
}

// GetBatchByTx returns the batch with the given hash.
func GetBatchByTx(db HostDB, txHash gethcommon.Hash) (*common.ExtBatch, error) {
	var seqNo uint64
	query := selectBatchSeqByTx + db.GetSQLStatement().Placeholder
	err := db.GetSQLDB().QueryRow(query, txHash).Scan(&seqNo)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errutil.ErrNotFound
		}
		return nil, err
	}
	return GetBatchBySequenceNumber(db, seqNo)
}

// GetBatchByHash returns the batch with the given hash.
func GetBatchByHash(db HostDB, hash common.L2BatchHash) (*common.ExtBatch, error) {
	whereQuery := " WHERE hash=" + db.GetSQLStatement().Placeholder
	return fetchFullBatch(db.GetSQLDB(), whereQuery, truncTo16(hash))
}

// GetLatestBatch returns the head batch header
func GetLatestBatch(db *sql.DB) (*common.BatchHeader, error) {
	headBatch, err := fetchHeadBatch(db)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch head batch: %w", err)
	}
	return headBatch.Header, nil
}

// GetBatchByHeight returns the batch header given the height
func GetBatchByHeight(db HostDB, height *big.Int) (*common.BatchHeader, error) {
	whereQuery := " WHERE height=" + db.GetSQLStatement().Placeholder
	return fetchBatchHeader(db.GetSQLDB(), whereQuery, height.Uint64())
}

// GetBatchTransactions returns the TransactionListingResponse for a given batch hash
func GetBatchTransactions(db HostDB, batchHash gethcommon.Hash) (*common.TransactionListingResponse, error) {
	whereQuery := " WHERE b.full_hash=" + db.GetSQLStatement().Placeholder
	return fetchBatchTxs(db.GetSQLDB(), whereQuery, batchHash)
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
		return nil, err
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
		return nil, err
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
	var hash []byte
	var heightInt64 int
	var extBatch []byte

	query := selectBatch + whereQuery

	var err error
	if len(args) > 0 {
		err = db.QueryRow(query, args...).Scan(&sequenceInt64, &fullHash, &hash, &heightInt64, &extBatch)
	} else {
		err = db.QueryRow(query).Scan(&sequenceInt64, &fullHash, &hash, &heightInt64, &extBatch)
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errutil.ErrNotFound
		}
		return nil, err
	}
	var b common.ExtBatch
	err = rlp.DecodeBytes(extBatch, &b)
	if err != nil {
		return nil, fmt.Errorf("could not decode ext batch. Cause: %w", err)
	}

	batch := &common.PublicBatch{
		SequencerOrderNo: new(big.Int).SetInt64(int64(sequenceInt64)),
		Hash:             hash,
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
	var hash []byte
	var heightInt64 int
	var extBatch []byte

	query := selectBatch + whereQuery

	var err error
	if len(args) > 0 {
		err = db.QueryRow(query, args...).Scan(&sequenceInt64, &fullHash, &hash, &heightInt64, &extBatch)
	} else {
		err = db.QueryRow(query).Scan(&sequenceInt64, &fullHash, &hash, &heightInt64, &extBatch)
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errutil.ErrNotFound
		}
		return nil, err
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
	var hash []byte
	var heightInt64 int
	var extBatch []byte

	err := db.QueryRow(selectLatestBatch).Scan(&sequenceInt64, &fullHash, &hash, &heightInt64, &extBatch)
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
		Hash:             hash,
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

func fetchBatchTxs(db *sql.DB, whereQuery string, batchHash gethcommon.Hash) (*common.TransactionListingResponse, error) {
	query := selectBatchTxs + whereQuery
	rows, err := db.Query(query, batchHash)
	if err != nil {
		return nil, fmt.Errorf("query execution for select batch txs failed: %w", err)
	}
	defer rows.Close()

	var transactions []common.PublicTransaction
	for rows.Next() {
		var (
			hash     []byte
			sequence int
			height   int
			extBatch []byte
		)
		err := rows.Scan(&hash, &sequence, &height, &extBatch)
		if err != nil {
			return nil, err
		}
		extBatchDecoded := new(common.ExtBatch)
		if err := rlp.DecodeBytes(extBatch, extBatchDecoded); err != nil {
			return nil, fmt.Errorf("could not decode batch. Cause: %w", err)
		}
		transaction := common.PublicTransaction{
			TransactionHash: common.TxHash(hash),
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
