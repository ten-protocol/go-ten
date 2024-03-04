package hostdb

import (
	"bytes"
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
	selectBatch             = "SELECT b.sequence_order, b.full_hash, b.hash, b.height, b.tx_count, b.header, b.body_id, bb.body FROM batch b JOIN batch_body_host bb ON b.body_id = bb.id"
	selectBatchBody         = "SELECT content FROM batch_body_host WHERE id = ?"
	selectDescendingBatches = `
		SELECT b.sequence_order, b.full_hash, b.hash, b.height, b.tx_count, b.header, b.body_id
		FROM batch_host b
		JOIN batch_body_host bb ON b.body_id = bb.id
		ORDER BY b.sequence_order DESC
		LIMIT 1
	`
	selectHeader                      = "SELECt b.header FROM batch_host b"
	selectTxsAndBatch                 = "SELECT t.full_hash FROM transactions_host t JOIN batch b ON t.body_id = b.body_id WHERE b.full_hash = ?"
	selectBatchNumberFromTransactions = "SELECT t.body_id FROM transactions_host t WHERE t.full_hash = ?"
	selectTxsBySequence               = "SELECT t.full_hash FROM transactions_host t WHERE t.body_id = ?"
	selectTxByHash                    = "SELECT t.body_id FROM transaction_host t WHERE t.full_hash = ?"

	insertBatchBody    = "INSERT INTO batch_body_host (id, content) VALUES (?, ?)"
	insertBatch        = "INSERT INTO batch_host (sequence_order, full_hash, hash, height, tx_count, header, body_id) VALUES (?, ?, ?, ?, ?, ?, ?)"
	insertTransactions = "INSERT INTO transactions_host (hash, full_hash, body_id) VALUES (?, ?, ?)"
	insertTxCount      = "INSERT INTO transaction_count (id, total) VALUES (?, ?) ON CONFLICT(id) DO UPDATE SET total = excluded.total;"
)

// AddBatch adds a batch and its header to the DB
func AddBatch(db *sql.DB, batch *common.ExtBatch) error {
	// Encode batch data
	batchBodyID := batch.SeqNo().Uint64()
	body, err := rlp.EncodeToBytes(batch.EncryptedTxBlob)
	if err != nil {
		return fmt.Errorf("could not encode L2 transactions: %w", err)
	}
	header, err := rlp.EncodeToBytes(batch.Header)
	if err != nil {
		return fmt.Errorf("could not encode batch header: %w", err)
	}

	// Insert the batch body data
	_, err = db.Exec(insertBatchBody, batchBodyID, body)
	if err != nil {
		return fmt.Errorf("failed to insert body: %w", err)
	}
	if len(batch.TxHashes) > 0 {
		//Insert transactions
		for _, transaction := range batch.TxHashes {
			shortHash := truncTo16(transaction)
			fullHash := transaction.Bytes()
			_, err := db.Exec(insertTransactions, shortHash, fullHash, batchBodyID)
			if err != nil {
				return fmt.Errorf("failed to insert transaction with hash: %d", err)
			}
		}
		//Increment total count
		var currentTotal int
		err := db.QueryRow(selectTxCount).Scan(&currentTotal)
		newTotal := currentTotal + len(batch.TxHashes)
		// Increase the TX count
		_, err = db.Exec(insertTxCount, 1, newTotal, newTotal)
		if err != nil {
			return fmt.Errorf("failed to update transaction count: %w", err)
		}
	}

	// Insert the batch data
	_, err = db.Exec(insertBatch,
		batch.SeqNo().Uint64(),       // sequence
		batch.Hash(),                 // full hash
		truncTo16(batch.Hash()),      // shortened hash
		batch.Header.Number.Uint64(), // height
		len(batch.TxHashes),          // tx_count
		header,                       // header blob
		batchBodyID,                  // reference to the batch body
	)

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
	if batchesTo < 0 {
		batchesTo = 0
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
	batch, err := fetchBatchHeader(db, " where sequence_order=?", number.Uint64())
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
	err := db.QueryRow(selectTxByHash).Scan(&seqNo)
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
	return fetchFullBatch(db, " where b.hash=?", truncTo16(hash))
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
	var sequenceInt64 int
	var fullHash common.TxHash
	var hash []byte
	var heightInt64 int
	var txCountInt64 int
	var headerBlob []byte
	var bodyID uint64

	query := selectBatch + " " + whereQuery

	var err error
	if len(args) > 0 {
		err = db.QueryRow(query, args...).Scan(&sequenceInt64, &fullHash, &hash, &heightInt64, &txCountInt64, &headerBlob, &bodyID)
	} else {
		err = db.QueryRow(query).Scan(&sequenceInt64, &fullHash, &hash, &heightInt64, &txCountInt64, &headerBlob, &bodyID)
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errutil.ErrNotFound
		}
		return nil, err
	}
	// Decode the batch header
	var header common.BatchHeader
	err = rlp.DecodeBytes(headerBlob, &header)
	if err != nil {
		return nil, fmt.Errorf("failed to decode batch header: %w", err)
	}
	sequence := new(big.Int).SetInt64(int64(sequenceInt64))
	height := new(big.Int).SetInt64(int64(heightInt64))
	txCount := new(big.Int).SetInt64(int64(txCountInt64))

	// Fetch batch_body from the database
	var encryptedTxBlob common.EncryptedTransactions
	err = db.QueryRow(selectBatchBody, bodyID).Scan(&encryptedTxBlob)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve batch body: %w", err)
	}

	var batchBody []byte
	err = rlp.DecodeBytes(encryptedTxBlob, &batchBody)
	if err != nil {
		return nil, fmt.Errorf("failed to decode batch body: %w", err)
	}

	// Construct the batch
	batch := &common.PublicBatch{
		SequencerOrderNo: sequence,
		Hash:             hash,
		FullHash:         fullHash,
		Height:           height,
		TxCount:          txCount,
		Header:           &header,
		EncryptedTxBlob:  encryptedTxBlob,
	}

	return batch, nil
}

func fetchFullBatch(db *sql.DB, whereQuery string, args ...any) (*common.ExtBatch, error) {
	var sequenceInt64 int
	var fullHash common.TxHash
	var shorthash []byte
	var heightInt64 int
	var txCountInt64 int
	var headerBlob []byte
	var bodyID uint64

	query := selectBatch + " " + whereQuery

	var err error
	if len(args) > 0 {
		err = db.QueryRow(query, args...).Scan(&sequenceInt64, &fullHash, &shorthash, &heightInt64, &txCountInt64, &headerBlob, &bodyID)
	} else {
		err = db.QueryRow(query).Scan(&sequenceInt64, &fullHash, &shorthash, &heightInt64, &txCountInt64, &headerBlob, &bodyID)
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errutil.ErrNotFound
		}
		return nil, err
	}
	// Decode the batch header
	var header common.BatchHeader
	err = rlp.DecodeBytes(headerBlob, &header)
	if err != nil {
		return nil, fmt.Errorf("failed to decode batch header: %w", err)
	}

	// Fetch batch_body from the database
	var encryptedTxBlob common.EncryptedTransactions
	err = db.QueryRow(selectBatchBody, bodyID).Scan(&encryptedTxBlob)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve batch body: %w", err)
	}

	var batchBody []byte
	err = rlp.DecodeBytes(encryptedTxBlob, &batchBody)
	if err != nil {
		return nil, fmt.Errorf("failed to decode batch body: %w", err)
	}

	txHashes, err := fetchTxBySeq(db, header.SequencerOrderNo.Uint64())
	if err != nil {
		return nil, fmt.Errorf("failed to get tx hashes for batch ")
	}

	batch := &common.ExtBatch{
		Header:          &header,
		TxHashes:        txHashes,
		EncryptedTxBlob: encryptedTxBlob,
	}

	return batch, nil
}

func fetchTxBySeq(db *sql.DB, seqNo uint64) ([]gethcommon.Hash, error) {
	rows, err := db.Query(selectTxsBySequence, seqNo)
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

func fetchHeadBatch(db *sql.DB) (*common.PublicBatch, error) {
	var sequenceInt64 int
	var fullHash gethcommon.Hash //common.Hash
	var hash []byte
	var heightInt64 int
	var txCountInt64 int
	var headerBlob []byte
	var bodyId int

	err := db.QueryRow(selectDescendingBatches).Scan(&sequenceInt64, &fullHash, &hash, &heightInt64, &txCountInt64, &headerBlob, &bodyId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errutil.ErrNotFound
		}
		return nil, fmt.Errorf("failed to fetch current head batch: %w", err)
	}

	var content []byte
	err = db.QueryRow(selectBatchBody, &bodyId).Scan(&content)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch batch content given the id: %w", err)
	}

	var header common.BatchHeader
	err = rlp.DecodeBytes(headerBlob, &header)
	if err != nil {
		return nil, fmt.Errorf("failed to decode batch header: %w", err)
	}

	sequence := new(big.Int).SetInt64(int64(sequenceInt64))
	height := new(big.Int).SetInt64(int64(heightInt64))
	txCount := new(big.Int).SetInt64(int64(txCountInt64))

	// Construct the batch
	batch := &common.PublicBatch{
		SequencerOrderNo: sequence,
		Hash:             hash,
		FullHash:         fullHash,
		Height:           height,
		TxCount:          txCount,
		Header:           &header,
		EncryptedTxBlob:  content,
	}

	return batch, nil
}
