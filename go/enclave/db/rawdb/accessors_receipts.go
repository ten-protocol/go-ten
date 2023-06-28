package rawdb

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/go-obscuro/go/common/errutil"

	common2 "github.com/obscuronet/go-obscuro/go/common"
)

// make sure that reading and incrementing is done atomically
var incrementSync sync.Mutex

// ReadReceiptsRLP retrieves all the transaction receipts belonging to a batch in RLP encoding.
func ReadReceiptsRLP(db ethdb.Reader, hash common.Hash) (rlp.RawValue, error) {
	data, err := db.Get(batchReceiptsKey(hash))
	if err != nil {
		return nil, fmt.Errorf("could not read receipts. Cause: %w", err)
	}
	return data, nil
}

// ReadRawReceipts retrieves all the transaction receipts belonging to a block.
// The receipt metadata fields are not guaranteed to be populated, so they
// should not be used. Use ReadReceipts instead if the metadata is needed.
func ReadRawReceipts(db ethdb.Reader, hash common.Hash) (types.Receipts, error) {
	// Retrieve the flattened receipt slice
	data, err := ReadReceiptsRLP(db, hash)
	if err != nil {
		return nil, err
	}
	// Convert the receipts from their storage form to their internal representation
	storageReceipts := []*types.ReceiptForStorage{}
	if err := rlp.DecodeBytes(data, &storageReceipts); err != nil {
		return nil, fmt.Errorf("invalid receipt array RLP. hash = %s; err = %w", hash, err)
	}
	receipts := make(types.Receipts, len(storageReceipts))
	for i, storageReceipt := range storageReceipts {
		receipts[i] = (*types.Receipt)(storageReceipt)
	}
	return receipts, nil
}

// ReadReceipts retrieves all the transaction receipts belonging to a block, including
// its corresponding metadata fields. If it is unable to populate these metadata
// fields then nil is returned.
//
// The current implementation populates these metadata fields by reading the receipts'
// corresponding block body, so if the block body is not found it will return nil even
// if the receipt itself is stored.
func ReadReceipts(db ethdb.Reader, hash common.Hash, number uint64, config *params.ChainConfig) (types.Receipts, error) {
	// We're deriving many fields from the block body, retrieve beside the receipt
	receipts, err := ReadRawReceipts(db, hash)
	if err != nil {
		return nil, fmt.Errorf("could not read receipt. Cause: %w", err)
	}
	body, err := readBatchBody(db, hash)
	if err != nil {
		return nil, fmt.Errorf("missing body but have receipt. Cause: %w", err)
	}

	if err = receipts.DeriveFields(config, hash, number, body); err != nil {
		return nil, fmt.Errorf("failed to derive block receipts fields. hash = %s; number = %d; err = %w", hash, number, err)
	}
	return receipts, nil
}

// WriteReceipts stores all the transaction receipts belonging to a batch.
func WriteReceipts(db ethdb.KeyValueWriter, hash common2.L2BatchHash, receipts types.Receipts) error {
	// Convert the receipts into their storage form and serialize them
	storageReceipts := make([]*types.ReceiptForStorage, len(receipts))
	for i, receipt := range receipts {
		storageReceipts[i] = (*types.ReceiptForStorage)(receipt)
	}
	bytes, err := rlp.EncodeToBytes(storageReceipts)
	if err != nil {
		return fmt.Errorf("failed to encode block receipts. Cause: %w", err)
	}
	// Store the flattened receipt slice
	if err = db.Put(batchReceiptsKey(hash), bytes); err != nil {
		return fmt.Errorf("failed to store block receipts. Cause: %w", err)
	}
	return nil
}

// WriteContractCreationTxs stores a mapping between each contract and the tx that created it
func WriteContractCreationTxs(db ethdb.KeyValueWriter, receipts types.Receipts) error {
	for _, receipt := range receipts {
		// determine receipts which create accounts and store the txHash
		if !bytes.Equal(receipt.ContractAddress.Bytes(), (common.Address{}).Bytes()) {
			if err := db.Put(contractReceiptKey(receipt.ContractAddress), receipt.TxHash.Bytes()); err != nil {
				return fmt.Errorf("failed to store contract receipt. Cause: %w", err)
			}
		}
	}
	return nil
}

func WriteCurrentBatchSequenceNumber(db ethdb.KeyValueWriter, sequenceNo *big.Int) error {
	return db.Put(sequenceNumber, sequenceNo.Bytes())
}

func ReadBatchSequenceNumber(db ethdb.KeyValueReader) (*big.Int, error) {
	value, err := db.Get(sequenceNumber)
	if err != nil {
		return nil, err
	}
	return big.NewInt(0).SetBytes(value), nil
}

// ReadContractTransaction - returns the tx that created a contract
func ReadContractTransaction(db ethdb.Reader, address common.Address) (*common.Hash, error) {
	value, err := db.Get(contractReceiptKey(address))
	if err != nil {
		return nil, err
	}
	hash := common.BytesToHash(value)
	return &hash, nil
}

func IncrementContractCreationCount(db ethdb.Reader, batch ethdb.KeyValueWriter, receipts []*types.Receipt) error {
	incrementSync.Lock()
	defer incrementSync.Unlock()

	contractCreationCounter := 0
	for _, receipt := range receipts {
		// receipts only have Contract Address when a new contract is created
		if !bytes.Equal(receipt.ContractAddress.Bytes(), (common.Address{}).Bytes()) {
			contractCreationCounter++
		}
	}
	if contractCreationCounter == 0 {
		return nil
	}

	current, err := ReadContractCreationCount(db)
	if err != nil {
		return err
	}
	total := big.NewInt(0).Add(current, big.NewInt(int64(contractCreationCounter)))

	if err = batch.Put(contractCreationCountKey(), total.Bytes()); err != nil {
		return fmt.Errorf("failed to store contract creation count - %w", err)
	}

	return nil
}

func ReadContractCreationCount(db ethdb.Reader) (*big.Int, error) {
	value, err := db.Get(contractCreationCountKey())
	if err != nil {
		if errors.Is(err, errutil.ErrNotFound) {
			return common.Big0, nil
		}
		return nil, fmt.Errorf("unable to read stored contract creation count - %w", err)
	}
	return big.NewInt(0).SetBytes(value), nil
}
