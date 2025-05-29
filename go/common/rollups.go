package common

// ExtRollup is an encrypted form of rollup used when passing the rollup around outside an enclave.
type ExtRollup struct {
	Header               *RollupHeader // the fields required by the management contract
	CalldataRollupHeader []byte        // encrypted header useful for recreating the batches
	BatchPayloads        []byte        // The transactions included in the rollup, in external/encrypted form.
	// hash                 atomic.Pointer[L2RollupHash]
}

// ExtRollupMetadata metadata that should not be in the rollup, but rather is derived from one.
// This should all be public information as it is passed back to the host!
type ExtRollupMetadata struct {
	CrossChainTree []byte // All the elements of the cross chain tree when building the rollup; Host uses this for serving cross chain proofs;
	// todo: Move signature here maybe?
}

// Hash returns the keccak256 hash of the rollup's header.
// The hash is computed on the first call and cached thereafter - disabled because we mutate the header
func (r *ExtRollup) Hash() L2RollupHash {
	//if hash := r.hash.Load(); hash != nil {
	//	return *hash
	//}
	h := r.Header.Hash()
	// r.hash.Store(&h)
	return h
}
