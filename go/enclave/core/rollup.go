package core

import "C"
import (
	"sync/atomic"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/obscuronet/go-obscuro/go/common"
)

// Rollup - is an internal data structure useful during creation
type Rollup struct {
	Header  *common.RollupHeader
	Batches []*Batch
	Blocks  map[common.L1BlockHash]*types.Block // these are the blocks required during compression. The key is the hash
	hash    atomic.Value
}

// Hash returns the keccak256 hash of b's header.
// The hash is computed on the first call and cached thereafter.
func (r *Rollup) Hash() common.L2BatchHash {
	if hash := r.hash.Load(); hash != nil {
		return hash.(common.L2BatchHash)
	}
	v := r.Header.Hash()
	r.hash.Store(v)
	return v
}
