package mempool

import (
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
)

// NonceTracker - a struct that helps us maintain the nonces for each account.
// If it gets asked for an account it does not know the nonce for, it will pull it
// from stateDB. Used when selecting transactions in order to ensure transactions get
// applied at correct nonces and correct order without any gaps.
type NonceTracker struct {
	accountNonces map[gethcommon.Address]uint64
	stateDB       *state.StateDB
}

func NewNonceTracker(stateDB *state.StateDB) *NonceTracker {
	return &NonceTracker{
		stateDB:       stateDB,
		accountNonces: make(map[gethcommon.Address]uint64),
	}
}

func (nt *NonceTracker) GetNonce(address gethcommon.Address) uint64 {
	if nonce, ok := nt.accountNonces[address]; ok {
		return nonce
	}

	nonce := nt.nonceFromState(address)
	nt.accountNonces[address] = nonce
	return nonce
}

func (nt *NonceTracker) nonceFromState(address gethcommon.Address) uint64 {
	return nt.stateDB.GetNonce(address)
}

func (nt *NonceTracker) IncrementNonce(address gethcommon.Address) {
	nt.accountNonces[address]++
}
