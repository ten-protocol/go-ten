package hostdb

import (
	"database/sql"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/errutil"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

const (
	selectBatch             = "SELECT sequence, hash, height, ext_batch FROM batch_host b"
	selectExtBatch          = "SELECT ext_batch FROM batch_host"
	selectLatestBatch       = "SELECT sequence, hash, height, ext_batch FROM batch_host ORDER BY sequence DESC LIMIT 1"
	selectTxsAndBatch       = "SELECT t.hash FROM transaction_host t JOIN batch_host b ON t.b_sequence = b.sequence WHERE b.hash = ?"
	selectBatchSeqByTx      = "SELECT b_sequence FROM transaction_host WHERE hash = ?"
	selectTxBySeq           = "SELECT hash FROM transaction_host WHERE b_sequence = ?"
	selectBatchTxs          = "SELECT t.hash, b.sequence, b.height, b.ext_batch FROM transaction_host t JOIN batch_host b ON t.b_sequence = b.sequence"
	selectSumBatchSizes     = "SELECT SUM(txs_size) FROM batch_host WHERE sequence >= ?"
	insertBatch             = "INSERT INTO batch_host (sequence, hash, height, ext_batch, txs_size) VALUES (?, ?, ?, ?, ?)"
	insertTransactions      = "INSERT INTO transaction_host (hash, b_sequence) VALUES "
	updateTxCount           = "UPDATE transaction_count SET total=? WHERE id=1"
	updateHistoricalTxCount = "UPDATE historical_transaction_count SET total=? WHERE id=1"

	whereHash     = " WHERE hash = ?"
	whereBHash    = " WHERE b.hash = ?"
	whereSequence = " WHERE sequence = ?"
	whereHeight   = " WHERE height = ?"
)

// AddBatch adds a batch and its header to the DB
func AddBatch(dbtx *dbTransaction, db HostDB, batch *common.ExtBatch) error {
	extBatch, err := rlp.EncodeToBytes(batch)
	if err != nil {
		return fmt.Errorf("could not encode L2 transactions: %w", err)
	}

	// this value is used on host side rollup size estimation
	// empty batches still contribute to rollup size through:
	// - RLP encoded empty transaction array (~3 bytes)
	// - Time delta (~1-2 bytes when compressed)
	// - L1 height delta (~1-2 bytes when compressed)
	// actual data from mainnet: 12,769 empty batches = 26KB rollup = ~2 bytes/batch final
	// batchCompressionFactor=0.1: 26KB ÷ 0.1 ÷ 12,769 = ~20 bytes raw estimate
	compressedSize := 20
	if len(batch.TxHashes) > 0 {
		// non-empty batches add transaction data on top of the base overhead
		compressedSize += len(batch.EncryptedTxBlob)
	}

	reboundInsertBatch := db.GetSQLDB().Rebind(insertBatch)
	_, err = dbtx.Tx.Exec(reboundInsertBatch,
		batch.SeqNo().Uint64(),       // sequence
		batch.Hash(),                 // full hash
		batch.Header.Number.Uint64(), // height
		extBatch,                     // ext_batch
		compressedSize,               // txs_size
	)
	if err != nil {
		if IsRowExistsError(err) {
			return errutil.ErrAlreadyExists
		}
		return fmt.Errorf("host failed to insert batch: %w", err)
	}

	if len(batch.TxHashes) > 0 {
		valuePlaceholders := make([]string, len(batch.TxHashes))
		args := make([]any, 0)
		for i, txHash := range batch.TxHashes {
			valuePlaceholders[i] = "(?, ?)"
			args = append(args, txHash.Bytes(), batch.SeqNo().Uint64())
		}

		insert := insertTransactions + strings.Join(valuePlaceholders, ", ")
		insert = db.GetSQLDB().Rebind(insert)
		_, err = dbtx.Tx.Exec(insert, args...)
		if err != nil {
			return fmt.Errorf("failed to insert transactions. cause: %w", err)
		}
	}

	// update current transaction count
	var currentTotal int
	reboundSelectTxCount := db.GetSQLDB().Rebind("SELECT total FROM transaction_count WHERE id = 1")
	err = dbtx.Tx.QueryRow(reboundSelectTxCount).Scan(&currentTotal)
	if err != nil {
		return fmt.Errorf("failed to query transaction count: %w", err)
	}

	newTotal := currentTotal + len(batch.TxHashes)
	reboundUpdateTxCount := db.GetSQLDB().Rebind(updateTxCount)
	_, err = dbtx.Tx.Exec(reboundUpdateTxCount, newTotal)
	if err != nil {
		return fmt.Errorf("failed to update transaction count: %w", err)
	}

	// update historical transaction count
	var currentHistoricalTotal int
	reboundSelectHistoricalTxCount := db.GetSQLDB().Rebind("SELECT total FROM historical_transaction_count WHERE id = 1")
	err = dbtx.Tx.QueryRow(reboundSelectHistoricalTxCount).Scan(&currentHistoricalTotal)
	if err != nil {
		return fmt.Errorf("failed to query historical transaction count: %w", err)
	}

	newHistoricalTotal := currentHistoricalTotal + len(batch.TxHashes)
	reoundHistoricalQuery := db.GetSQLDB().Rebind(updateHistoricalTxCount)
	_, err = dbtx.Tx.Exec(reoundHistoricalQuery, newHistoricalTotal)
	if err != nil {
		return fmt.Errorf("failed to update historical transaction count: %w", err)
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
	reboundWhereQuery := db.GetSQLDB().Rebind(whereSequence)
	return fetchPublicBatch(db, reboundWhereQuery, seqNo)
}

// GetTxsBySequenceNumber returns the transaction hashes with sequence number.
func GetTxsBySequenceNumber(db HostDB, seqNo uint64) ([]common.TxHash, error) {
	return fetchTx(db, seqNo)
}

// GetBatchBySequenceNumber returns the ext batch for a given sequence number.
func GetBatchBySequenceNumber(db HostDB, seqNo uint64) (*common.ExtBatch, error) {
	reboundWhereQuery := db.GetSQLDB().Rebind(whereSequence)
	return fetchFullBatch(db, reboundWhereQuery, seqNo)
}

// GetCurrentHeadBatch retrieves the current head batch with the largest sequence number (or height).
func GetCurrentHeadBatch(db HostDB) (*common.PublicBatch, error) {
	return fetchHeadBatch(db)
}

// GetBatchHeader returns the batch header given the hash.
func GetBatchHeader(db HostDB, hash gethcommon.Hash) (*common.BatchHeader, error) {
	reboundWhereQuery := db.GetSQLDB().Rebind(whereHash)
	return fetchBatchHeader(db, reboundWhereQuery, hash.Bytes())
}

// GetBatchHashByNumber returns the hash of a batch given its number.
func GetBatchHashByNumber(db HostDB, number *big.Int) (*gethcommon.Hash, error) {
	reboundWhereQuery := db.GetSQLDB().Rebind(whereHeight)
	batch, err := fetchBatchHeader(db, reboundWhereQuery, number.Uint64())
	if err != nil {
		return nil, fmt.Errorf("failed to fetch batch header - %w", err)
	}
	l2BatchHash := batch.Hash()
	return &l2BatchHash, nil
}

// GetHeadBatchHeader returns the latest batch header.
func GetHeadBatchHeader(db HostDB) (*common.BatchHeader, error) {
	batch, err := fetchHeadBatch(db)
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
	reboundQuery := db.GetSQLDB().Rebind(selectTxsAndBatch)
	rows, err := db.GetSQLDB().Query(reboundQuery, batchHash.Bytes())
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
	return fetchPublicBatch(db, whereBHash, hash.Bytes())
}

// GetBatchByTx returns the batch with the given hash.
func GetBatchByTx(db HostDB, txHash gethcommon.Hash) (*common.PublicBatch, error) {
	var seqNo uint64
	reboundQuery := db.GetSQLDB().Rebind(selectBatchSeqByTx)
	err := db.GetSQLDB().QueryRow(reboundQuery, txHash.Bytes()).Scan(&seqNo)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errutil.ErrNotFound
		}
		return nil, fmt.Errorf("failed to execute query %s - %w", reboundQuery, err)
	}
	extBatch, err := GetBatchBySequenceNumber(db, seqNo)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch ext batch - %w", err)
	}
	return toPublicBatch(extBatch), nil
}

// GetBatchByHash returns the batch with the given hash.
func GetBatchByHash(db HostDB, hash common.L2BatchHash) (*common.PublicBatch, error) {
	reboundWhereQuery := db.GetSQLDB().Rebind(whereHash)
	return fetchPublicBatch(db, reboundWhereQuery, hash.Bytes())
}

// GetBatchHeaderByHeight returns the batch header given the height
func GetBatchHeaderByHeight(db HostDB, height *big.Int) (*common.BatchHeader, error) {
	reboundWhereQuery := db.GetSQLDB().Rebind(whereHeight)
	return fetchBatchHeader(db, reboundWhereQuery, height.Uint64())
}

// GetBatchByHeight returns the batch header given the height
func GetBatchByHeight(db HostDB, height *big.Int) (*common.PublicBatch, error) {
	reboundWhereQuery := db.GetSQLDB().Rebind(whereHeight)
	return fetchPublicBatch(db, reboundWhereQuery, height.Uint64())
}

// GetBatchTransactions returns the TransactionListingResponse for a given batch hash
func GetBatchTransactions(db HostDB, batchHash gethcommon.Hash, pagination *common.QueryPagination) (*common.TransactionListingResponse, error) {
	orderQuery := " ORDER BY t.id DESC "
	query := selectBatchTxs + whereBHash + orderQuery + paginationQuery
	countQuery := "SELECT COUNT(*) FROM transaction_host t JOIN batch_host b ON t.b_sequence = b.sequence" + whereBHash
	reboundQuery := db.GetSQLDB().Rebind(query)
	reboundCountQuery := db.GetSQLDB().Rebind(countQuery)
	var total uint64
	err := db.GetSQLDB().QueryRow(reboundCountQuery, batchHash.Bytes()).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %w", err)
	}

	rows, err := db.GetSQLDB().Query(reboundQuery, batchHash.Bytes(), int64(pagination.Size), int64(pagination.Offset))
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
	reboundQuery := db.GetSQLDB().Rebind(selectSumBatchSizes)
	err := db.GetSQLDB().QueryRow(reboundQuery, fromSeqNo.Uint64()).Scan(&totalTx)
	if err != nil {
		return 0, fmt.Errorf("failed to query sum of rollup batches: %w", err)
	}
	return totalTx, nil
}

func fetchBatchHeader(db HostDB, whereQuery string, args ...any) (*common.BatchHeader, error) {
	var extBatch []byte
	query := selectExtBatch + whereQuery
	var err error
	reboundQuery := db.GetSQLDB().Rebind(query)
	if len(args) > 0 {
		err = db.GetSQLDB().QueryRow(reboundQuery, args...).Scan(&extBatch)
	} else {
		err = db.GetSQLDB().QueryRow(reboundQuery).Scan(&extBatch)
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
	reboundQuery := db.GetSQLDB().Rebind(selectBatchSeqByTx)
	var err error
	if len(args) > 0 {
		err = db.GetSQLDB().QueryRow(reboundQuery, args...).Scan(&seqNo)
	} else {
		err = db.GetSQLDB().QueryRow(reboundQuery).Scan(&seqNo)
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errutil.ErrNotFound
		}
		return nil, fmt.Errorf("failed to scan with query %s - %w", reboundQuery, err)
	}
	batch, err := GetPublicBatchBySequenceNumber(db, seqNo)
	if err != nil {
		return nil, fmt.Errorf("could not fetch batch by seq no. Cause: %w", err)
	}
	return batch.Height, nil
}

func fetchPublicBatch(db HostDB, whereQuery string, args ...any) (*common.PublicBatch, error) {
	var sequenceInt64 uint64
	var fullHash common.TxHash
	var heightInt64 int
	var extBatch []byte

	query := selectBatch + whereQuery
	reboundQuery := db.GetSQLDB().Rebind(query)
	var err error
	if len(args) > 0 {
		err = db.GetSQLDB().QueryRow(reboundQuery, args...).Scan(&sequenceInt64, &fullHash, &heightInt64, &extBatch)
	} else {
		err = db.GetSQLDB().QueryRow(reboundQuery).Scan(&sequenceInt64, &fullHash, &heightInt64, &extBatch)
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
		Header:           b.Header,
		EncryptedTxBlob:  b.EncryptedTxBlob,
		TxHashes:         b.TxHashes,
	}

	return batch, nil
}

func fetchFullBatch(db HostDB, whereQuery string, args ...any) (*common.ExtBatch, error) {
	var sequenceInt64 uint64
	var fullHash common.TxHash
	var heightInt64 int
	var extBatch []byte

	query := selectBatch + whereQuery
	reboundQuery := db.GetSQLDB().Rebind(query)
	var err error
	if len(args) > 0 {
		err = db.GetSQLDB().QueryRow(reboundQuery, args...).Scan(&sequenceInt64, &fullHash, &heightInt64, &extBatch)
	} else {
		err = db.GetSQLDB().QueryRow(reboundQuery).Scan(&sequenceInt64, &fullHash, &heightInt64, &extBatch)
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

func fetchHeadBatch(db HostDB) (*common.PublicBatch, error) {
	var sequenceInt64 int
	var fullHash gethcommon.Hash // common.Hash
	var heightInt64 int
	var extBatch []byte

	reboundQuery := db.GetSQLDB().Rebind(selectLatestBatch)
	err := db.GetSQLDB().QueryRow(reboundQuery).Scan(&sequenceInt64, &fullHash, &heightInt64, &extBatch)
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
		Header:           b.Header,
		EncryptedTxBlob:  b.EncryptedTxBlob,
		TxHashes:         b.TxHashes,
	}

	return batch, nil
}

func fetchTx(db HostDB, seqNo uint64) ([]common.TxHash, error) {
	reboundQuery := db.GetSQLDB().Rebind(selectTxBySeq)
	rows, err := db.GetSQLDB().Query(reboundQuery, seqNo)
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

func IsRowExistsError(err error) bool {
	return strings.Contains(strings.ToLower(err.Error()), "unique") || strings.Contains(strings.ToLower(err.Error()), "duplicate key")
}

func toPublicBatch(b *common.ExtBatch) *common.PublicBatch {
	return &common.PublicBatch{
		SequencerOrderNo: b.SeqNo(),
		FullHash:         b.Hash(),
		Height:           b.Header.Number,
		TxHashes:         b.TxHashes,
		Header:           b.Header,
		EncryptedTxBlob:  b.EncryptedTxBlob,
	}
}
