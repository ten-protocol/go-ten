package rollupchain

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/obscuronet/go-obscuro/go/enclave/db"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

const (
	FaucetPrivateKeyHex = "8dfb8083da6275ae3e4f41e3e8a8c19d028d32c9247e24530933782f2a05035b" // The faucet's private key.
	faucetPrealloc      = "7500000000000000000000000000000"                                  // The balance preallocated to the faucet address.
)

// Faucet handles the preallocation of funds in the network.
type Faucet struct {
	faucetAddress gethcommon.Address
}

func NewFaucet() Faucet {
	faucetPrivateKey, err := crypto.HexToECDSA(FaucetPrivateKeyHex)
	if err != nil {
		panic("could not convert faucet private key from hex to ECDSA")
	}

	return Faucet{
		faucetAddress: crypto.PubkeyToAddress(faucetPrivateKey.PublicKey),
	}
}

// GetGenesisRoot applies the faucet preallocation on top of an empty state DB, and returns the corresponding trie
// root.
func (f *Faucet) GetGenesisRoot(storage db.Storage) (*gethcommon.Hash, error) {
	stateDB, err := f.applyFaucetPrealloc(storage)
	if err != nil {
		return nil, err
	}
	stateHash := stateDB.IntermediateRoot(true)
	return &stateHash, nil
}

// CommitGenesisState applies the faucet preallocation on top of an empty state DB and commits the result.
func (f *Faucet) CommitGenesisState(storage db.Storage) (*state.StateDB, error) {
	stateDB, err := f.applyFaucetPrealloc(storage)
	if err != nil {
		return nil, err
	}
	_, err = stateDB.Commit(true)
	if err != nil {
		return nil, err
	}
	return stateDB, nil
}

// Applies the faucet preallocation on top of an empty state DB.
func (f *Faucet) applyFaucetPrealloc(storage db.Storage) (*state.StateDB, error) {
	s, err := storage.EmptyStateDB()
	if err != nil {
		return nil, fmt.Errorf("could not initialise empty state DB. Cause: %w", err)
	}

	faucetPreallocBig, success := big.NewInt(0).SetString(faucetPrealloc, 10)
	if !success {
		return nil, fmt.Errorf("could not initialise faucet prealloc Big from string %s", faucetPrealloc)
	}

	s.SetBalance(f.faucetAddress, faucetPreallocBig)
	return s, nil
}
