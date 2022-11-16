package common

import (
	"sync/atomic"
)

// ExtRollup is an encrypted form of rollup used when passing the rollup around outside of an enclave.
type ExtRollup struct {
	Header          *Header
	TxHashes        []TxHash // The hashes of the transactions included in the rollup
	EncryptedTxBlob EncryptedTransactions
	hash            atomic.Value
}

func (r ExtRollup) ToExtRollup() *ExtRollup {
	return &ExtRollup{
		Header:          r.Header,
		TxHashes:        r.TxHashes,
		EncryptedTxBlob: r.EncryptedTxBlob,
	}
}

// Hash returns the keccak256 hash of the rollup's header.
// The hash is computed on the first call and cached thereafter.
func (r *ExtRollup) Hash() L2RootHash {
	if hash := r.hash.Load(); hash != nil {
		return hash.(L2RootHash)
	}
	v := r.Header.Hash()
	r.hash.Store(v)
	return v
}
