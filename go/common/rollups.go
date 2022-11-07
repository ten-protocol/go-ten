package common

import (
	"sync/atomic"
)

// ExtRollup is an encrypted form of rollup used when passing the rollup around outside of an enclave.
type ExtRollup struct {
	Header          *Header
	TxHashes        []TxHash // The hashes of the transactions included in the rollup
	EncryptedTxBlob EncryptedTransactions
}

// ExtRollupWithHash extends ExtRollup with additional fields.
// This parallels the Block/extblock split in Geth.
type ExtRollupWithHash struct {
	ExtRollup
	hash atomic.Value
}

func (er ExtRollup) ToExtRollupWithHash() *ExtRollupWithHash {
	extRollup := ExtRollup{
		Header:          er.Header,
		TxHashes:        er.TxHashes,
		EncryptedTxBlob: er.EncryptedTxBlob,
	}
	return &ExtRollupWithHash{
		ExtRollup: extRollup,
	}
}

// Hash returns the keccak256 hash of b's header.
// The hash is computed on the first call and cached thereafter.
func (r *ExtRollupWithHash) Hash() L2RootHash {
	if hash := r.hash.Load(); hash != nil {
		return hash.(L2RootHash)
	}
	v := r.Header.Hash()
	r.hash.Store(v)

	return v
}
