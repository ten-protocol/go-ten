package mempool

import (
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
)

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
	nt.accountNonces[address] += 1
}
