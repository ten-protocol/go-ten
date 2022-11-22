package db

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"

	"github.com/obscuronet/go-obscuro/go/common/errutil"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/go-obscuro/go/common"
)

// TODO - #718 - Remove any rollup-based getters that are superseded by the batch-based getters.

// DB methods relating to rollups.

// GetHeadRollupHeader returns the header of the node's current head rollup, or (nil, false) if no such header is found.
func (db *DB) GetHeadRollupHeader() (*common.Header, error) {
	headRollupHash, err := db.readHeadRollupHash()
	if err != nil {
		return nil, err
	}
	return db.readRollupHeader(*headRollupHash)
}

// GetRollupHeader returns the rollup header given the hash, or (nil, false) if no such header is found.
func (db *DB) GetRollupHeader(hash gethcommon.Hash) (*common.Header, error) {
	return db.readRollupHeader(hash)
}

// AddRollupHeader adds a rollup's header to the known headers
func (db *DB) AddRollupHeader(header *common.Header, txHashes []common.TxHash) error {
	b := db.kvStore.NewBatch()

	if err := db.writeRollupHeader(b, header); err != nil {
		return fmt.Errorf("could not write rollup header. Cause: %w", err)
	}
	// Required by ObscuroScan, to display a list of recent transactions.
	if err := db.writeRollupTxHashes(b, header.Hash(), txHashes); err != nil {
		return fmt.Errorf("could not write rollup transaction hashes. Cause: %w", err)
	}
	if err := db.writeRollupHash(b, header); err != nil {
		return fmt.Errorf("could not write rollup hash. Cause: %w", err)
	}
	for _, txHash := range txHashes {
		if err := db.writeRollupNumber(b, header, txHash); err != nil {
			return fmt.Errorf("could not write rollup number. Cause: %w", err)
		}
	}

	// There's a potential race here, but absolute accuracy of the number of transactions is not required.
	currentTotal, err := db.readTotalTransactions()
	if err != nil {
		return fmt.Errorf("could not retrieve total transactions. Cause: %w", err)
	}
	newTotal := big.NewInt(0).Add(currentTotal, big.NewInt(int64(len(txHashes))))
	err = db.writeTotalTransactions(b, newTotal)
	if err != nil {
		return fmt.Errorf("could not write total transactions. Cause: %w", err)
	}

	// update the head if the new height is greater than the existing one
	headRollupHeader, err := db.GetHeadRollupHeader()
	if err != nil && !errors.Is(err, errutil.ErrNotFound) {
		return fmt.Errorf("could not retrieve head rollup header. Cause: %w", err)
	}
	if errors.Is(err, errutil.ErrNotFound) || headRollupHeader.Number.Int64() <= header.Number.Int64() {
		err = db.writeHeadRollupHash(b, header.Hash())
		if err != nil {
			return fmt.Errorf("could not write new head rollup hash. Cause: %w", err)
		}
	}

	if err = b.Write(); err != nil {
		return fmt.Errorf("could not write batch to DB. Cause: %w", err)
	}
	return nil
}

// GetRollupHash returns the hash of a rollup given its number, or (nil, false) if no such rollup is found.
func (db *DB) GetRollupHash(number *big.Int) (*gethcommon.Hash, error) {
	return db.readRollupHash(number)
}

// GetRollupNumber returns the number of the rollup containing the given transaction hash, or (nil, false) if no such rollup is found.
func (db *DB) GetRollupNumber(txHash gethcommon.Hash) (*big.Int, error) {
	return db.readRollupNumber(txHash)
}

// GetTotalTransactions returns the total number of rolled-up transactions.
// TODO - #718 - Return number of batched transactions, instead.
func (db *DB) GetTotalTransactions() (*big.Int, error) {
	return db.readTotalTransactions()
}

// headerKey = rollupHeaderPrefix  + hash
func rollupHeaderKey(hash gethcommon.Hash) []byte {
	return append(rollupHeaderPrefix, hash.Bytes()...)
}

// headerKey = rollupTxHashesPrefix + rollup hash
func rollupTxHashesKey(hash gethcommon.Hash) []byte {
	return append(rollupTxHashesPrefix, hash.Bytes()...)
}

// headerKey = rollupHashPrefix + number
func rollupHashKey(num *big.Int) []byte {
	return append(rollupHashPrefix, []byte(num.String())...)
}

// headerKey = rollupNumberPrefix + hash
func rollupNumberKey(txHash gethcommon.Hash) []byte {
	return append(rollupNumberPrefix, txHash.Bytes()...)
}

// Stores a rollup header into the database
func (db *DB) writeRollupHeader(w ethdb.KeyValueWriter, header *common.Header) error {
	// Write the encoded header
	data, err := rlp.EncodeToBytes(header)
	if err != nil {
		return err
	}
	key := rollupHeaderKey(header.Hash())
	if err = w.Put(key, data); err != nil {
		return err
	}
	return nil
}

// Retrieves the rollup header corresponding to the hash, or (nil, false) if no such header is found.
func (db *DB) readRollupHeader(hash gethcommon.Hash) (*common.Header, error) {
	f, err := db.kvStore.Has(rollupHeaderKey(hash))
	if err != nil {
		return nil, err
	}
	if !f {
		return nil, errutil.ErrNotFound
	}
	data, err := db.kvStore.Get(rollupHeaderKey(hash))
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errutil.ErrNotFound
	}
	header := new(common.Header)
	if err := rlp.Decode(bytes.NewReader(data), header); err != nil {
		return nil, err
	}
	return header, nil
}

// Writes the transaction hashes against the rollup containing them.
func (db *DB) writeRollupTxHashes(w ethdb.KeyValueWriter, rollupHash common.L2RootHash, txHashes []gethcommon.Hash) error {
	data, err := rlp.EncodeToBytes(txHashes)
	if err != nil {
		return err
	}
	key := rollupTxHashesKey(rollupHash)
	if err = w.Put(key, data); err != nil {
		return err
	}
	return nil
}

// Returns the head rollup's hash, or (nil, false) is no such hash is found.
func (db *DB) readHeadRollupHash() (*gethcommon.Hash, error) {
	f, err := db.kvStore.Has(headRollup)
	if err != nil {
		return nil, err
	}
	if !f {
		return nil, errutil.ErrNotFound
	}
	value, err := db.kvStore.Get(headRollup)
	if err != nil {
		return nil, err
	}
	h := gethcommon.BytesToHash(value)
	return &h, nil
}

// Stores the head rollup's hash in the database.
func (db *DB) writeHeadRollupHash(w ethdb.KeyValueWriter, val gethcommon.Hash) error {
	err := w.Put(headRollup, val.Bytes())
	if err != nil {
		return err
	}
	return nil
}

// Stores a rollup's hash in the database, keyed by the rollup's number.
func (db *DB) writeRollupHash(w ethdb.KeyValueWriter, header *common.Header) error {
	key := rollupHashKey(header.Number)
	if err := w.Put(key, header.Hash().Bytes()); err != nil {
		return err
	}
	return nil
}

// Stores a rollup's number in the database, keyed by the hash of a transaction in that rollup.
func (db *DB) writeRollupNumber(w ethdb.KeyValueWriter, header *common.Header, txHash gethcommon.Hash) error {
	key := rollupNumberKey(txHash)
	// TODO - Investigate this off-by-one issue. The tx hashes that are in the `BlockSubmissionResponse` for rollup #1
	//  are actually the transactions for rollup #2.
	number := big.NewInt(0).Add(header.Number, big.NewInt(1))
	if err := w.Put(key, number.Bytes()); err != nil {
		return err
	}
	return nil
}

// Stores the total number of transactions in the database.
func (db *DB) writeTotalTransactions(w ethdb.KeyValueWriter, newTotal *big.Int) error {
	err := w.Put(totalTransactionsKey, newTotal.Bytes())
	if err != nil {
		return err
	}
	return nil
}

// Retrieves the hash for the rollup with the given number, or (nil, false) if no such rollup is found.
func (db *DB) readRollupHash(number *big.Int) (*gethcommon.Hash, error) {
	f, err := db.kvStore.Has(rollupHashKey(number))
	if err != nil {
		return nil, err
	}
	if !f {
		return nil, errutil.ErrNotFound
	}
	data, err := db.kvStore.Get(rollupHashKey(number))
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errutil.ErrNotFound
	}
	hash := gethcommon.BytesToHash(data)
	return &hash, nil
}

// Retrieves the number of the rollup containing the transaction with the given hash, or (nil, false) if no such rollup is found.
func (db *DB) readRollupNumber(txHash gethcommon.Hash) (*big.Int, error) {
	f, err := db.kvStore.Has(rollupNumberKey(txHash))
	if err != nil {
		return nil, err
	}
	if !f {
		return nil, errutil.ErrNotFound
	}
	data, err := db.kvStore.Get(rollupNumberKey(txHash))
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errutil.ErrNotFound
	}
	return big.NewInt(0).SetBytes(data), nil
}

// Retrieves the total number of rolled-up transactions.
func (db *DB) readTotalTransactions() (*big.Int, error) {
	f, err := db.kvStore.Has(totalTransactionsKey)
	if err != nil {
		return nil, err
	}
	if !f {
		return big.NewInt(0), nil
	}
	data, err := db.kvStore.Get(totalTransactionsKey)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return big.NewInt(0), nil
	}
	return big.NewInt(0).SetBytes(data), nil
}
