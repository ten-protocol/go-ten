package common

import (
	"sort"
	"sync/atomic"
)

// ExtRollup is an encrypted form of rollup used when passing the rollup around outside an enclave.
type ExtRollup struct {
	Header *RollupHeader
	// todo (#718) - consider compressing these batches before submitting to the L1.
	Batches []*ExtBatch // The batches included in the rollup, in external/encrypted form.
	hash    atomic.Value
}

// Hash returns the keccak256 hash of the rollup's header.
// The hash is computed on the first call and cached thereafter.
func (r *ExtRollup) Hash() L2BatchHash {
	if hash := r.hash.Load(); hash != nil {
		return hash.(L2BatchHash)
	}
	v := r.Header.Hash()
	r.hash.Store(v)
	return v
}

// ExtRollupFromExtBatches converts a set of ExtBatch into an ExtRollup. The hash of the head batch (the one with the
// highest number) is stored in the ExtRollup header's `HeadBatchHash` field.
func ExtRollupFromExtBatches(batches []*ExtBatch) *ExtRollup {
	sort.Slice(batches, func(i, j int) bool {
		// Ascending order sort.
		return batches[i].Header.Number.Cmp(batches[j].Header.Number) < 0
	})

	return &ExtRollup{
		Header:  batches[len(batches)-1].Header.ToRollupHeader(),
		Batches: batches,
	}
}
