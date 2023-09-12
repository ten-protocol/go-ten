package core

import "C"
import (
	"sync/atomic"

	"github.com/obscuronet/go-obscuro/go/common"
)

// todo - This should be a synthetic datastructure
type Rollup struct {
	Header  *common.RollupHeader
	Batches []*Batch
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
