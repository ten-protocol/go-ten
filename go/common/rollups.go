package common

import (
	"sync/atomic"
)

// ExtRollup is an encrypted form of rollup used when passing the rollup around outside an enclave.
type ExtRollup struct {
	Header *RollupHeader
	// todo (#718) - consider compressing these batches before submitting to the L1.
	BatchPayloads []byte // The batches included in the rollup, in external/encrypted form.
	BatchHeaders  []byte // compressed blob of a serialised list of batch headers
	hash          atomic.Value
}

// Hash returns the keccak256 hash of the rollup's header.
// The hash is computed on the first call and cached thereafter.
func (r *ExtRollup) Hash() L2BatchHash {
	//if hash := r.hash.Load(); hash != nil {
	//	return hash.(L2BatchHash)
	//}
	v := r.Header.Hash()
	r.hash.Store(v)
	return v
}
