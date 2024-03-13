package hostdb

import (
	"database/sql"
	"errors"
	"fmt"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/errutil"
	"math/big"
)

const (
	selectTxCount           = "SELECT total FROM transaction_count WHERE id = 1"
	selectBatch             = "SELECT sequence, full_hash, hash, height, ext_batch FROM batch_host"
	selectExtBatch          = "SELECT ext_batch FROM batch_host"
	selectDescendingBatches = `
		SELECT sequence, full_hash, hash, height, ext_batch
		FROM batch_host
		ORDER BY sequence DESC
		LIMIT 1
	`
	selectTxsAndBatch                 = "SELECT t.full_hash FROM transactions_host t JOIN batch_host b ON t.b_sequence = b.sequence WHERE b.full_hash = ?"
	selectBatchNumberFromTransactions = "SELECT b_sequence FROM transactions_host WHERE full_hash = ?"
	selectTxByHash                    = "SELECT b_sequence FROM transaction_host t WHERE full_hash = ?"

	insertBatch        = "INSERT INTO batch_host (sequence, full_hash, hash, height, ext_batch) VALUES (?, ?, ?, ?, ?)"
	insertTransactions = "INSERT INTO transactions_host (hash, full_hash, b_sequence) VALUES (?, ?, ?)"
	insertTxCount      = "INSERT INTO transaction_count (id, total) VALUES (?, ?) ON CONFLICT(id) DO UPDATE SET total = excluded.total;"
)

// AddBatch adds a batch and its header to the DB
func AddBatch(db *sql.DB, batch *common.ExtBatch) error {
	// Check if the Batch is already stored
	_, err := GetBatchHeader(db, batch.Hash())
	if err == nil {
		return errutil.ErrAlreadyExists
	}

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	extBatch, err := rlp.EncodeToBytes(batch)
	if err != nil {
		return fmt.Errorf("could not encode L2 transactions: %w", err)
	}

	_, err = tx.Exec(insertBatch,
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
			shortHash := truncTo16(transaction)
			fullHash := transaction.Bytes()
			_, err := tx.Exec(insertTransactions, shortHash, fullHash, batch.SeqNo().Uint64())
			if err != nil {
				return fmt.Errorf("failed to insert transaction with hash: %d", err)
			}
		}
		//Increment total count
		var currentTotal int
		err := tx.QueryRow(selectTxCount).Scan(&currentTotal)
		newTotal := currentTotal + len(batch.TxHashes)
		// Increase the TX count
		_, err = tx.Exec(insertTxCount, 1, newTotal, newTotal)
		if err != nil {
			return fmt.Errorf("failed to update transaction count: %w", err)
		}
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("could not commit batch tx: %w", err)
	}

	return nil
}

// GetBatchListing returns latest batches given a pagination.
// For example, page 0, size 10 will return the latest 10 batches.
func GetBatchListing(db *sql.DB, pagination *common.QueryPagination) (*common.BatchListingResponse, error) {
	headBatch, err := GetCurrentHeadBatch(db)
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
		batch, err := GetBatchBySequenceNumber(db, i)
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

// GetBatchBySequenceNumber returns the batch with the given sequence number.
func GetBatchBySequenceNumber(db *sql.DB, seqNo uint64) (*common.PublicBatch, error) {
	return fetchPublicBatch(db, " WHERE sequence=?", seqNo)
}

func GetFullBatchBySequenceNumber(db *sql.DB, seqNo uint64) (*common.ExtBatch, error) {
	return fetchFullBatch(db, " WHERE sequence=?", seqNo)
}

// GetCurrentHeadBatch retrieves the current head batch with the largest sequence number (or height)
func GetCurrentHeadBatch(db *sql.DB) (*common.PublicBatch, error) {
	return fetchHeadBatch(db)
}

// GetBatchHeader returns the batch header given the hash.
func GetBatchHeader(db *sql.DB, hash gethcommon.Hash) (*common.BatchHeader, error) {
	return fetchBatchHeader(db, " where hash=?", truncTo16(hash))
}

// GetBatchHashByNumber returns the hash of a batch given its number.
func GetBatchHashByNumber(db *sql.DB, number *big.Int) (*gethcommon.Hash, error) {
	batch, err := fetchBatchHeader(db, " where height=?", number.Uint64())
	if err != nil {
		return nil, err
	}
	l2BatchHash := batch.Hash()
	return &l2BatchHash, nil
}

func GetHeadBatchHeader(db *sql.DB) (*common.BatchHeader, error) {
	batch, err := fetchHeadBatch(db)
	if err != nil {
		return nil, err
	}
	return batch.Header, nil
}

// GetBatchNumber returns the number of the batch containing the given transaction hash.
func GetBatchNumber(db *sql.DB, txHash gethcommon.Hash) (*big.Int, error) {
	batchNumber, err := fetchBatchNumber(db, txHash)
	if err != nil {
		return nil, err
	}
	return batchNumber, nil
}

// GetBatchTxs returns the transaction hashes of the batch with the given hash.
func GetBatchTxs(db *sql.DB, batchHash gethcommon.Hash) ([]gethcommon.Hash, error) {
	rows, err := db.Query(selectTxsAndBatch, batchHash)
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

// GetTotalTransactions returns the total number of batched transactions.
func GetTotalTransactions(db *sql.DB) (*big.Int, error) {
	var totalCount int
	err := db.QueryRow(selectTxCount).Scan(&totalCount)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve total transaction count: %w", err)
	}
	return big.NewInt(int64(totalCount)), nil
}

// GetPublicBatch returns the batch with the given hash.
func GetPublicBatch(db *sql.DB, hash common.L2BatchHash) (*common.PublicBatch, error) {
	return fetchPublicBatch(db, " where b.hash=?", truncTo16(hash))
}

// GetFullBatchByTx returns the batch with the given hash.
func GetFullBatchByTx(db *sql.DB, txHash gethcommon.Hash) (*common.ExtBatch, error) {
	var seqNo uint64
	err := db.QueryRow(selectTxByHash, txHash).Scan(&seqNo)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// make sure the error is converted to obscuro-wide not found error
			return nil, errutil.ErrNotFound
		}
		return nil, err
	}
	return GetFullBatchBySequenceNumber(db, seqNo)
}

// GetFullBatch returns the batch with the given hash.
func GetFullBatch(db *sql.DB, hash common.L2BatchHash) (*common.ExtBatch, error) {
	return fetchFullBatch(db, " where hash=?", truncTo16(hash))
}

func fetchBatchHeader(db *sql.DB, whereQuery string, args ...any) (*common.BatchHeader, error) {
	var extBatch []byte
	query := selectExtBatch + " " + whereQuery
	var err error
	if len(args) > 0 {
		err = db.QueryRow(query, args...).Scan(&extBatch)
	} else {
		err = db.QueryRow(query).Scan(&extBatch)
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// make sure the error is converted to obscuro-wide not found error
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

func fetchBatchNumber(db *sql.DB, args ...any) (*big.Int, error) {
	var batchNumber uint64
	query := selectBatchNumberFromTransactions
	var err error
	if len(args) > 0 {
		err = db.QueryRow(query, args...).Scan(&batchNumber)
	} else {
		err = db.QueryRow(query).Scan(&batchNumber)
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errutil.ErrNotFound
		}
		return nil, err
	}
	var bigIntBatchNumber = new(big.Int).SetInt64(int64(batchNumber))
	return bigIntBatchNumber, nil
}

func fetchPublicBatch(db *sql.DB, whereQuery string, args ...any) (*common.PublicBatch, error) {
	var sequenceInt64 uint64
	var fullHash common.TxHash
	var hash []byte
	var heightInt64 int
	var extBatch []byte

	query := selectBatch + " " + whereQuery

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
	// Decode batch
	var b common.ExtBatch
	err = rlp.DecodeBytes(extBatch, &b)

	// Construct the batch
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
	// Decode batch
	var b common.ExtBatch
	err = rlp.DecodeBytes(extBatch, &b)

	return &b, nil
}

func fetchHeadBatch(db *sql.DB) (*common.PublicBatch, error) {
	var sequenceInt64 int
	var fullHash gethcommon.Hash //common.Hash
	var hash []byte
	var heightInt64 int
	var extBatch []byte

	err := db.QueryRow(selectDescendingBatches).Scan(&sequenceInt64, &fullHash, &hash, &heightInt64, &extBatch)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errutil.ErrNotFound
		}
		return nil, fmt.Errorf("failed to fetch current head batch: %w", err)
	}

	var b common.ExtBatch
	err = rlp.DecodeBytes(extBatch, &b)

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
