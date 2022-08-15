package rollupchain

import (
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/obscuronet/go-obscuro/go/enclave/db"
)

const (
	FaucetPrivateKeyHex = "8dfb8083da6275ae3e4f41e3e8a8c19d028d32c9247e24530933782f2a05035b" // The faucet's private key.
	faucetPrealloc      = 7500000000000000000                                                // The balance preallocated to the faucet address.
)

// Faucet handles the preallocation of funds in the network.
type Faucet struct {
	storage       db.Storage
	faucetAddress gethcommon.Address
}

func NewFaucet(storage db.Storage) Faucet {
	faucetPrivateKey, err := crypto.HexToECDSA(FaucetPrivateKeyHex)
	if err != nil {
		panic("could not convert faucet private key from hex to ECDSA")
	}
	faucetAddress := crypto.PubkeyToAddress(faucetPrivateKey.PublicKey)

	return Faucet{
		storage:       storage,
		faucetAddress: faucetAddress,
	}
}

// GetGenesisRoot applies the faucet preallocation on top of an empty state DB, and returns the corresponding trie
// root.
func (f *Faucet) GetGenesisRoot(storage db.Storage) gethcommon.Hash {
	stateDB := f.applyFaucetPrealloc(storage)
	return stateDB.IntermediateRoot(true)
}

// CalculateGenesisState applies the faucet preallocation on top of an empty state DB and commits the result.
func (f *Faucet) CalculateGenesisState(storage db.Storage) error {
	stateDB := f.applyFaucetPrealloc(storage)
	_, err := stateDB.Commit(true)
	if err != nil {
		return err
	}
	return nil
}

// Applies the faucet preallocation on top of an empty state DB.
func (f *Faucet) applyFaucetPrealloc(storage db.Storage) *state.StateDB {
	s := storage.EmptyStateDB()
	s.SetBalance(f.faucetAddress, big.NewInt(faucetPrealloc))
	return s
}
