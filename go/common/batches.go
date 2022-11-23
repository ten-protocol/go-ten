package common

import "sync/atomic"

// ExtBatch is an encrypted form of batch used when passing the batch around outside of an enclave.
// TODO - #718 - Expand this structure to contain the required fields.
type ExtBatch struct {
	Header          *Header
	TxHashes        []TxHash // The hashes of the transactions included in the batch.
	EncryptedTxBlob EncryptedTransactions
	hash            atomic.Value
}

// Hash returns the keccak256 hash of the batch's header.
// The hash is computed on the first call and cached thereafter.
func (r *ExtBatch) Hash() L2RootHash {
	if hash := r.hash.Load(); hash != nil {
		return hash.(L2RootHash)
	}
	v := r.Header.Hash()
	r.hash.Store(v)
	return v
}
