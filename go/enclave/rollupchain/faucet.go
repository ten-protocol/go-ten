package rollupchain

import (
	"math/big"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/obscuronet/go-obscuro/go/enclave/db"
)

const (
	faucetAddressHex = "0xA58C60cc047592DE97BF1E8d2f225Fc5D959De77" // The faucet address.
	faucetPrealloc   = 7500000000000000000                          // The balance preallocated to the faucet address.
)

// Faucet handles the preallocation of funds in the network.
type Faucet struct {
	storage db.Storage
}

func NewFaucet(storage db.Storage) Faucet {
	return Faucet{
		storage: storage,
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
	s.SetBalance(gethcommon.HexToAddress(faucetAddressHex), big.NewInt(faucetPrealloc))
	return s
}
