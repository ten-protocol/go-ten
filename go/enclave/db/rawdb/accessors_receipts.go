package rawdb

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/obscuronet/go-obscuro/go/common/errutil"
)

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
