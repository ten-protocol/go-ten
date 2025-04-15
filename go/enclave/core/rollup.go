package core

import (
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ten-protocol/go-ten/go/common"
)

// Rollup - is an internal data structure useful during creation
type Rollup struct {
	Header  *common.RollupHeader
	Batches []*Batch
	Blocks  map[common.L1BlockHash]*types.Header // these are the blocks required during compression. The key is the hash
	// hash    atomic.Pointer[common.L2RollupHash]
}

// Hash returns the keccak256 hash of b's header.
// The hash is computed on the first call and cached thereafter - disabled because we mutate the header
func (r *Rollup) Hash() common.L2BatchHash {
	//if hash := r.hash.Load(); hash != nil {
	//	return *hash
	//}
	h := r.Header.Hash()
	// r.hash.Store(&h)
	return h
}
