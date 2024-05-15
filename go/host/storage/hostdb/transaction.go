package hostdb

import (
	"fmt"
	"math/big"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ten-protocol/go-ten/go/common"
)

const (
	selectTxCount = "SELECT total FROM transaction_count WHERE id = 1"
	selectTx      = "SELECT full_hash, b_sequence FROM transaction_host WHERE hash = "
	selectTxs     = "SELECT t.full_hash, b.ext_batch FROM transaction_host t JOIN batch_host b ON t.b_sequence = b.sequence ORDER BY b.height DESC "
)

// GetTransactionListing returns a paginated list of transactions in descending order
func GetTransactionListing(db HostDB, pagination *common.QueryPagination) (*common.TransactionListingResponse, error) {
	query := selectTxs + db.GetSQLStatement().Pagination
	rows, err := db.GetSQLDB().Query(query, pagination.Size, pagination.Offset)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query %s - %w", query, err)
	}
	defer rows.Close()
	var txs []common.PublicTransaction

	for rows.Next() {
		var fullHash, extBatch []byte

		err = rows.Scan(&fullHash, &extBatch)
		if err != nil {
			return nil, fmt.Errorf("failed to scan query %s - %w", query, err)
		}

		b := new(common.ExtBatch)
		if err := rlp.DecodeBytes(extBatch, b); err != nil {
			return nil, fmt.Errorf("could not decode rollup hash. Cause: %w", err)
		}
		tx := common.PublicTransaction{
			TransactionHash: gethcommon.HexToHash(bytesToHexString(fullHash)),
			BatchHeight:     b.Header.Number,
			BatchTimestamp:  b.Header.Time,
			// TODO @will this will be implemented under #3336
			Finality: common.BatchFinal,
		}
		txs = append(txs, tx)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	println("GetTransactionListing offset: ", pagination.Offset)
	println("GetTransactionListing size: ", pagination.Size)
	println("------")
	return &common.TransactionListingResponse{
		TransactionsData: txs,
		Total:            uint64(len(txs)),
	}, nil
}

// GetTransaction returns a transaction given its hash
func GetTransaction(db HostDB, hash gethcommon.Hash) (*common.PublicTransaction, error) {
	query := selectTx + db.GetSQLStatement().Placeholder

	var fullHash []byte
	var seq int
	err := db.GetSQLDB().QueryRow(query, truncTo16(hash)).Scan(&fullHash, &seq)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve transaction sequence number: %w", err)
	}

	batch, err := GetBatchBySequenceNumber(db, uint64(seq))
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve batch by sequence number: %w", err)
	}

	tx := &common.PublicTransaction{
		TransactionHash: gethcommon.BytesToHash(fullHash),
		BatchHeight:     batch.Header.Number,
		BatchTimestamp:  batch.Header.Time,
		Finality:        common.BatchFinal,
	}

	return tx, nil
}

// GetTotalTxCount returns the total number of batched transactions
func GetTotalTxCount(db HostDB) (*big.Int, error) {
	var totalCount int
	err := db.GetSQLDB().QueryRow(selectTxCount).Scan(&totalCount)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve total transaction count: %w", err)
	}
	return big.NewInt(int64(totalCount)), nil
}
