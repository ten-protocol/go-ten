package rawdb

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/obscuronet/go-obscuro/go/common/errutil"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/params"
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

// WriteTxLookupEntries is identical to WriteTxLookupEntry, but it works on
// a list of hashes
func WriteTxLookupEntries(db ethdb.KeyValueWriter, number uint64, hashes []common.Hash) error {
	numberBytes := new(big.Int).SetUint64(number).Bytes()
	for _, hash := range hashes {
		err := writeTxLookupEntry(db, hash, numberBytes)
		if err != nil {
			return err
		}
	}
	return nil
}

// WriteTxLookupEntriesByBlock stores a positional metadata for every transaction from
// a block, enabling hash based transaction and receipt lookups.
func WriteTxLookupEntriesByBlock(db ethdb.KeyValueWriter, rollup *core.Rollup) error {
	numberBytes := rollup.Number().Bytes()
	for _, tx := range rollup.Transactions {
		err := writeTxLookupEntry(db, tx.Hash(), numberBytes)
		if err != nil {
			return err
		}
	}
	return nil
}

// DeleteTxLookupEntries removes all transaction lookups for a given block.
func DeleteTxLookupEntries(db ethdb.KeyValueWriter, hashes []common.Hash) error {
	for _, hash := range hashes {
		err := deleteTxLookupEntry(db, hash)
		if err != nil {
			return fmt.Errorf("could not delete transaction lookuo entry. Cause: %w", err)
		}
	}
	return nil
}

// Removes all transaction data associated with a hash.
func deleteTxLookupEntry(db ethdb.KeyValueWriter, hash common.Hash) error {
	if err := db.Delete(txLookupKey(hash)); err != nil {
		return fmt.Errorf("failed to delete transaction lookup entry. Cause: %w", err)
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

	blockHash, err := ReadCanonicalHash(db, *blockNumber)
	if err != nil {
		return nil, common.Hash{}, 0, 0, fmt.Errorf("could not retrieve canonical hash for block number. Cause: %w", err)
	}

	transactions, err := ReadBody(db, *blockHash, *blockNumber)
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

// ReadReceipt retrieves a specific transaction receipt from the database, along with
// its added positional metadata.
func ReadReceipt(db ethdb.Reader, hash common.Hash, config *params.ChainConfig) (*types.Receipt, common.Hash, uint64, uint64, error) {
	// Retrieve the context of the receipt based on the transaction hash
	blockNumber, err := ReadTxLookupEntry(db, hash)
	if err != nil {
		return nil, common.Hash{}, 0, 0, fmt.Errorf("could not retrieve transaction lookup entry. Cause: %w", err)
	}
	blockHash, err := ReadCanonicalHash(db, *blockNumber)
	if err != nil {
		return nil, common.Hash{}, 0, 0, fmt.Errorf("could not retrieve canonical hash for block number. Cause: %w", err)
	}

	// Read all the receipts from the block and return the one with the matching hash
	receipts, err := ReadReceipts(db, *blockHash, *blockNumber, config)
	if err != nil {
		return nil, common.Hash{}, 0, 0, fmt.Errorf("could not retrieve receipts for block number. Cause: %w", err)
	}
	for receiptIndex, receipt := range receipts {
		if receipt.TxHash == hash {
			return receipt, *blockHash, *blockNumber, uint64(receiptIndex), nil
		}
	}
	return nil, common.Hash{}, 0, 0, fmt.Errorf("receipt not found. Cause: %w", err)
}

// ReadBloomBits retrieves the compressed bloom bit vector belonging to the given
// section and bit index.
func ReadBloomBits(db ethdb.KeyValueReader, bit uint, section uint64, head common.Hash) ([]byte, error) {
	return db.Get(bloomBitsKey(bit, section, head))
}

// WriteBloomBits stores the compressed bloom bits vector belonging to the given
// section and bit index.
func WriteBloomBits(db ethdb.KeyValueWriter, bit uint, section uint64, head common.Hash, bits []byte) error {
	if err := db.Put(bloomBitsKey(bit, section, head), bits); err != nil {
		return fmt.Errorf("failed to store bloom bits. Cause: %w", err)
	}
	return nil
}

// DeleteBloombits removes all compressed bloom bits vector belonging to the
// given section range and bit index.
func DeleteBloombits(db ethdb.Database, bit uint, from uint64, to uint64) error {
	start, end := bloomBitsKey(bit, from, common.Hash{}), bloomBitsKey(bit, to, common.Hash{})
	it := db.NewIterator(nil, start)
	defer it.Release()

	for it.Next() {
		if bytes.Compare(it.Key(), end) >= 0 {
			break
		}
		if len(it.Key()) != len(bloomBitsPrefix)+2+8+32 {
			continue
		}
		err := db.Delete(it.Key())
		if err != nil {
			return fmt.Errorf("failed to delete bloom bits. Cause: %w", err)
		}
	}
	if it.Error() != nil {
		return fmt.Errorf("failed to delete bloom bits. Cause: %w", it.Error())
	}
	return nil
}
