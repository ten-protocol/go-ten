package hostdb

import (
	"database/sql"
	"fmt"
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ten-protocol/go-ten/go/common"
)

// AddBatch adds a batch and its header to the DB
func AddBatch(db *sql.DB, batch *common.ExtBatch) error {
	return BeginTx(db, func(tx *sql.Tx) error {

		batchBodyStmt, err := db.Prepare("INSERT INTO batch_body (id, body) VALUES (?, ?)")
		if err != nil {
			return fmt.Errorf("failed to prepare body insert statement: %w", err)
		}
		defer batchBodyStmt.Close()

		// BATCH INSERT
		batchStmt, err := db.Prepare("INSERT INTO batches (sequence, full_hash, hash, height, tx_count, header_blob, body_id) VALUES (?, ?, ?, ?, ?, ?)")
		if err != nil {
			return fmt.Errorf("failed to prepare batch insert statement: %w", err)
		}
		defer batchStmt.Close()

		//TX INSERT
		txStmt, err := db.Prepare("INSERT INTO transactions (tx_hash_indexed, tx_hash_full, content, sender_address, nonce, idx, body_id) VALUES (?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			return fmt.Errorf("failed to prepare transaction insert statement: %w", err)
		}
		defer txStmt.Close()

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

// GetHeadBatchHeader returns the header of the node's current head batch.
func GetHeadBatchHeader(db *sql.DB) (*common.BatchHeader, error) {
	panic("implement me")
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
	panic("implement me")
}

// GetBatch returns the batch with the given hash.
func GetBatch(db *sql.DB, batchHash gethcommon.Hash) (*common.ExtBatch, error) {
	panic("implement me")
}

// GetBatchBySequenceNumber returns the batch with the given sequence number.
func GetBatchBySequenceNumber(db *sql.DB, sequenceNumber *big.Int) (*common.ExtBatch, error) {
	panic("implement me")
}

// GetBatchListing returns latest batches given a pagination.
// For example, page 0, size 10 will return the latest 10 batches.
func GetBatchListing(db *sql.DB, pagination *common.QueryPagination) (*common.BatchListingResponse, error) {
	panic("implement me")
}

// Retrieves the batch header corresponding to the hash.
func readBatchHeader(db *sql.DB, hash gethcommon.Hash) (*common.BatchHeader, error) {
	panic("implement me")
}

// Retrieves the hash of the head batch.
func readHeadBatchHash(db *sql.DB) (*gethcommon.Hash, error) {
	panic("implement me")
}

// Stores a batch header into the database.
func writeBatchHeader(db *sql.DB, w ethdb.KeyValueWriter, header *common.BatchHeader) error {
	panic("implement me")
}

// Stores the head batch header hash into the database.
func writeHeadBatchHash(db *sql.DB, w ethdb.KeyValueWriter, val gethcommon.Hash) error {
	panic("implement me")
}

// Stores a batch's hash in the database, keyed by the batch's number.
func writeBatchHash(db *sql.DB, w ethdb.KeyValueWriter, header *common.BatchHeader) error {
	panic("implement me")
}

// Stores a batch's hash in the database, keyed by the batch's sequencer number.
func writeBatchSeqNo(db *sql.DB, w ethdb.KeyValueWriter, header *common.BatchHeader) error {
	panic("implement me")
}

// Retrieves the hash for the batch with the given number..
func readBatchHash(db *sql.DB, number *big.Int) (*gethcommon.Hash, error) {
	panic("implement me")
}

// Returns the transaction hashes in the batch with the given hash.
func readBatchTxHashes(db *sql.DB, batchHash common.L2BatchHash) ([]gethcommon.Hash, error) {
	panic("implement me")
}

// Stores a batch's number in the database, keyed by the hash of a transaction in that rollup.
func writeBatchNumber(db *sql.DB, w ethdb.KeyValueWriter, header *common.BatchHeader, txHash gethcommon.Hash) error {
	panic("implement me")
}

// Writes the transaction hashes against the batch containing them.
func writeBatchTxHashes(db *sql.DB, w ethdb.KeyValueWriter, batchHash common.L2BatchHash, txHashes []gethcommon.Hash) error {
	panic("implement me")
}

// Retrieves the number of the batch containing the transaction with the given hash.
func readBatchNumber(db *sql.DB, txHash gethcommon.Hash) (*big.Int, error) {
	panic("implement me")
}

func readBatchHashBySequenceNumber(db *sql.DB, seqNum *big.Int) (*gethcommon.Hash, error) {
	panic("implement me")
}

// Retrieves the total number of rolled-up transactions - returns 0 if no tx count is found
func readTotalTransactions(db *sql.DB) (*big.Int, error) {
	panic("implement me")
}

// Stores the total number of transactions in the database.
func writeTotalTransactions(db *sql.DB, w ethdb.KeyValueWriter, newTotal *big.Int) error {
	panic("implement me")
}

// Stores a batch into the database.
func writeBatch(db *sql.DB, w ethdb.KeyValueWriter, batch *common.ExtBatch) error {
	panic("implement me")
}

// Retrieves the batch corresponding to the hash.
func readBatch(db *sql.DB, hash gethcommon.Hash) (*common.ExtBatch, error) {
	panic("implement me")
}
