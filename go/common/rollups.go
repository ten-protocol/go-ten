package common

import (
	"sync/atomic"
)

// ExtRollup is an encrypted form of rollup used when passing the rollup around outside of an enclave.
// TODO - #718 - This structure can now be deleted, since there is no private information in the "vanilla" rollup (as rollups no longer contain transactions).
type ExtRollup struct {
	Header      *RollupHeader
	BatchHashes []L2RootHash // The hashes of the batches included in the rollup
	hash        atomic.Value
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
