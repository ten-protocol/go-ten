package rawdb

import (
	"fmt"
	"math/big"

	"github.com/obscuronet/go-obscuro/go/common/errutil"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/obscuronet/go-obscuro/go/enclave/core"
)

// ReadTxLookupEntry retrieves the positional metadata associated with a transaction
// hash to allow retrieving the transaction or receipt by hash.
func ReadTxLookupEntry(db ethdb.Reader, hash common.Hash) (*uint64, error) {
	data, err := db.Get(txLookupKey(hash))
	if err != nil {
		return nil, errutil.ErrNotFound
	}

	// Database v6 tx lookup just stores the block number
	if len(data) >= common.HashLength {
		return nil, fmt.Errorf("transaction positional metadata was too long. Cause: %w", err)
	}

	number := new(big.Int).SetBytes(data).Uint64()
	return &number, nil
}

// writeTxLookupEntry stores a positional metadata for a transaction,
// enabling hash based transaction and receipt lookups.
func writeTxLookupEntry(db ethdb.KeyValueWriter, hash common.Hash, numberBytes []byte) error {
	if err := db.Put(txLookupKey(hash), numberBytes); err != nil {
		return fmt.Errorf("failed to store transaction lookup entry. Cause: %w", err)
	}
	return nil
}

// WriteTxLookupEntriesByBatch stores a positional metadata for every transaction from a batch, enabling hash based
// transaction and receipt lookups.
func WriteTxLookupEntriesByBatch(db ethdb.KeyValueWriter, batch *core.Batch) error {
	for _, tx := range batch.Transactions {
		err := writeTxLookupEntry(db, tx.Hash(), batch.Number().Bytes())
		if err != nil {
			return err
		}
	}
	return nil
}

// ReadTransaction retrieves a specific transaction from the database, along with
// its added positional metadata.
func ReadTransaction(db ethdb.Reader, hash common.Hash) (*types.Transaction, common.Hash, uint64, uint64, error) {
	blockNumber, err := ReadTxLookupEntry(db, hash)
	if err != nil {
		return nil, common.Hash{}, 0, 0, fmt.Errorf("could not retrieve transaction lookup entry. Cause: %w", err)
	}

	blockHash, err := ReadCanonicalBatchHash(db, *blockNumber)
	if err != nil {
		return nil, common.Hash{}, 0, 0, fmt.Errorf("could not retrieve canonical hash for block number. Cause: %w", err)
	}

	transactions, err := ReadBody(db, *blockHash)
	if err != nil {
		return nil, common.Hash{}, 0, 0, fmt.Errorf("could not retrieve block body. Cause: %w", err)
	}
	for txIndex, tx := range transactions {
		if tx.Hash() == hash {
			return tx, *blockHash, *blockNumber, uint64(txIndex), nil
		}
	}
	return nil, common.Hash{}, 0, 0, fmt.Errorf("transaction not found")
}
