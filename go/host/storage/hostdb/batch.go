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
	selectTxCount           = "SELECT count FROM transaction_count WHERE id = 1"
	selectBatch             = "SELECT b.sequence, b.full_hash, b.hash, b.height, b.tx_count, b.header_blob, b.body_id, bb.body FROM batches b JOIN batch_body bb ON b.body_id = bb.id"
	selectBatchBody         = "SELECT body FROM batch_body WHERE id = ?"
	selectDescendingBatches = `
		SELECT b.sequence, b.full_hash, b.hash, b.height, b.tx_count, b.header_blob, bb.body
		FROM batches b
		JOIN batch_body bb ON b.body_id = bb.id
		ORDER BY b.sequence DESC
		LIMIT 1
	`

	insertBatchBody = "INSERT INTO batch_body (id, body) VALUES (?, ?)"
	insertBatch     = "INSERT INTO batches (sequence, full_hash, hash, height, tx_count, header_blob, body_id) VALUES (?, ?, ?, ?, ?, ?)"
	insertTxCount   = "INSERT INTO transaction_count (id, count) VALUES (?,?) ON DUPLICATE KEY UPDATE count = count + ?"
)

// AddBatch adds a batch and its header to the DB
func AddBatch(db *sql.DB, batch *common.ExtBatch) error {
	return BeginTx(db, func(tx *sql.Tx) error {

		// Batch body insert
		batchBodyStmt, err := db.Prepare(insertBatchBody)
		if err != nil {
			return fmt.Errorf("failed to prepare body insert statement: %w", err)
		}
		defer batchBodyStmt.Close()

		// Batch insert
		batchStmt, err := db.Prepare(insertBatch)
		if err != nil {
			return fmt.Errorf("failed to prepare batch insert statement: %w", err)
		}
		defer batchStmt.Close()

		// Tx count insert
		txStmt, err := db.Prepare(insertTxCount)
		if err != nil {
			return fmt.Errorf("failed to prepare tx count insert statement: %w", err)
		}
		defer batchStmt.Close()

		// Encode batch data
		batchBodyID := batch.Header.SequencerOrderNo.Uint64()
		body, err := rlp.EncodeToBytes(batch.EncryptedTxBlob)
		if err != nil {
			return fmt.Errorf("could not encode L2 transactions: %w", err)
		}
		header, err := rlp.EncodeToBytes(batch.Header)
		if err != nil {
			return fmt.Errorf("could not encode batch header: %w", err)
		}

		// Execute body insert
		_, err = batchBodyStmt.Exec(batchBodyID, body)
		if err != nil {
			return fmt.Errorf("failed to insert body: %w", err)
		}

		if len(batch.TxHashes) > 0 {
			_, err = txStmt.Exec(1, len(batch.TxHashes), len(batch.TxHashes))
			if err != nil {
				return fmt.Errorf("failed to update transaction count: %w", err)
			}
		}

		_, err = batchStmt.Exec(
			batch.Header.SequencerOrderNo.Uint64(), // sequence
			batch.Hash(),                           // full hash
			truncTo16(batch.Hash()),                // shortened hash
			batch.Header.Number.Uint64(),           // height
			len(batch.TxHashes),                    // tx_count
			header,                                 // header blob
			batchBodyID,                            // reference to the batch body
		)

		return nil
	})
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

	var placeHolderTxHashes []common.TxHash //FIXME remove from ExtBatch?

	batch := &common.ExtBatch{
		Header:          &header,
		TxHashes:        placeHolderTxHashes,
		EncryptedTxBlob: encryptedTxBlob,
	}

	return batch, nil
}

// GetCurrentHeadBatch retrieves the current head batch with the largest sequence number (or height)
func GetCurrentHeadBatch(db *sql.DB) (*common.PublicBatch, error) {
	var sequenceInt64 int
	var fullHash gethcommon.Hash //common.Hash
	var hash []byte
	var heightInt64 int
	var txCountInt64 int
	var headerBlob []byte
	var encryptedTxBlob common.EncryptedTransactions

	err := db.QueryRow(selectDescendingBatches).Scan(&sequenceInt64, &fullHash, &hash, &heightInt64, &txCountInt64, &headerBlob, &encryptedTxBlob)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no batches found")
		}
		return nil, fmt.Errorf("failed to fetch current head batch: %w", err)
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

// GetBatchHeader returns the batch header given the hash.
func GetBatchHeader(db *sql.DB, hash gethcommon.Hash) (*common.BatchHeader, error) {
	panic("implement me")
}

// GetBatchHash returns the hash of a batch given its number.
func GetBatchHash(db *sql.DB, number *big.Int) (*gethcommon.Hash, error) {
	panic("implement me")
}

// GetBatchTxs returns the transaction hashes of the batch with the given hash.
func GetBatchTxs(db *sql.DB, batchHash gethcommon.Hash) ([]gethcommon.Hash, error) {
	panic("implement me")
}

// GetBatchNumber returns the number of the batch containing the given transaction hash.
func GetBatchNumber(db *sql.DB, txHash gethcommon.Hash) (*big.Int, error) {
	panic("implement me")
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

// GetFullBatch returns the batch with the given hash.
func GetFullBatch(db *sql.DB, hash common.L2BatchHash) (*common.ExtBatch, error) {
	return fetchFullBatch(db, " where b.hash=?", truncTo16(hash))
}
