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
	selectTxCount           = "SELECT count FROM transaction_count WHERE id = 1"
	selectBatch             = "SELECT b.sequence_order, b.full_hash, b.hash, b.height, b.tx_count, b.header, b.body_id, bb.body FROM batch b JOIN batch_body bb ON b.body_id = bb.id"
	selectBatchBody         = "SELECT content FROM batch_body WHERE id = ?"
	selectDescendingBatches = `
		SELECT b.sequence_order, b.full_hash, b.hash, b.height, b.tx_count, b.header, b.body_id
		FROM batch b
		JOIN batch_body bb ON b.body_id = bb.id
		ORDER BY b.sequence_order DESC
		LIMIT 1
	`
	selectHeader = "select b.header from batch b"

	insertBatchBody = "INSERT INTO batch_body (id, content) VALUES (?, ?)"
	insertBatch     = "INSERT INTO batch (sequence_order, full_hash, hash, height, tx_count, header, body_id) VALUES (?, ?, ?, ?, ?, ?,?)"
	insertTxCount   = "INSERT INTO transaction_count (id, count) VALUES (?, ?) ON DUPLICATE KEY UPDATE count = ?"
)

// AddBatch adds a batch and its header to the DB
func AddBatch(db *sql.DB, batch *common.ExtBatch) error {

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
	_, err = db.Exec(insertBatchBody, batchBodyID, body)
	//_, err = batchBodyStmt.Exec(batchBodyID, body)
	if err != nil {
		return fmt.Errorf("failed to insert body: %w", err)
	}
	if len(batch.TxHashes) > 0 {
		var currentTotal int
		err := db.QueryRow(selectTxCount).Scan(&currentTotal)
		if err != nil {
			return fmt.Errorf("failed to retrieve current tx total value: %w", err)
		}
		newTotal := currentTotal + len(batch.TxHashes)
		_, err = db.Exec(insertTxCount, 1, newTotal, newTotal)
		if err != nil {
			return fmt.Errorf("failed to update transaction count: %w", err)
		}
	}

	_, err = db.Exec(insertBatch,
		batch.Header.SequencerOrderNo.Uint64(), // sequence
		batch.Hash(),                           // full hash
		truncTo16(batch.Hash()),                // shortened hash
		batch.Header.Number.Uint64(),           // height
		len(batch.TxHashes),                    // tx_count
		header,                                 // header blob
		batchBodyID,                            // reference to the batch body
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

// GetCurrentHeadBatch retrieves the current head batch with the largest sequence number (or height)
func GetCurrentHeadBatch(db *sql.DB) (*common.PublicBatch, error) {
	return fetchHeadBatch(db)
}

// GetBatchHeader returns the batch header given the hash.
func GetBatchHeader(db *sql.DB, hash gethcommon.Hash) (*common.BatchHeader, error) {
	return fetchBatchHeader(db, " where hash=?", truncTo16(hash))
}

// GetBatchHash returns the hash of a batch given its number.
func GetBatchHash(db *sql.DB, number *big.Int) (*gethcommon.Hash, error) {
	panic("implement me")
}

// GetBatchTxs returns the transaction hashes of the batch with the given hash.
func GetBatchTxs(db *sql.DB, batchHash gethcommon.Hash) ([]gethcommon.Hash, error) {
	panic("implement me")
}

func GetHeadBatchHeader(db *sql.DB) (*common.BatchHeader, error) {
	batch, err := fetchHeadBatch(db)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch head batch: %w", err)
	}
	return batch.Header, nil
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

func fetchHeadBatch(db *sql.DB) (*common.PublicBatch, error) {
	var sequenceInt64 int
	var fullHash gethcommon.Hash //common.Hash
	var hash []byte
	var heightInt64 int
	var txCountInt64 int
	var headerBlob []byte
	var body_id int

	err := db.QueryRow(selectDescendingBatches).Scan(&sequenceInt64, &fullHash, &hash, &heightInt64, &txCountInt64, &headerBlob, &body_id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no batches found")
		}
		return nil, fmt.Errorf("failed to fetch current head batch: %w", err)
	}
	//Select from batch_body table
	var content []byte
	err = db.QueryRow(selectBatchBody, &body_id).Scan(&content)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch batch content given the id: %w", err)
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
		EncryptedTxBlob:  content,
	}

	return batch, nil
}
