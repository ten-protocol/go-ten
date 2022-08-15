package rollupchain

import (
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/obscuronet/go-obscuro/go/enclave/db"
	"math/big"
)

const (
	faucetPrealloc   = 7500000000000000000 // The balance preallocated to the faucet address.
	faucetAddressHex = "0xA58C60cc047592DE97BF1E8d2f225Fc5D959De77"
)

// todo - joel - make this into a class

// GetGenesisRoot applies the faucet preallocation to an empty state DB, and returns the corresponding trie
// root.
func GetGenesisRoot(storage db.Storage) gethcommon.Hash {
	stateDB := applyFaucetPrealloc(storage)
	return stateDB.IntermediateRoot(true)
}

// CommitGenesis applies the faucet preallocation to an empty state DB and commits the result.
func CommitGenesis(storage db.Storage) error {
	stateDB := applyFaucetPrealloc(storage)
	_, err := stateDB.Commit(true)
	if err != nil {
		return err
	}
	return nil
}

// Applies the faucet preallocation to an empty state DB.
func applyFaucetPrealloc(storage db.Storage) *state.StateDB {
	s := storage.EmptyStateDB()
	s.SetBalance(gethcommon.HexToAddress(faucetAddressHex), big.NewInt(faucetPrealloc))
	return s
}
