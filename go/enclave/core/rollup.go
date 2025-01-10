package core

import (
	"sync/atomic"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ten-protocol/go-ten/go/common"
)

// Rollup - is an internal data structure useful during creation
type Rollup struct {
	Header  *common.RollupHeader
	Batches []*Batch
	Blocks  map[common.L1BlockHash]*types.Header // these are the blocks required during compression. The key is the hash
	hash    atomic.Value
}

type BlobRollup struct {
	Rollup
	Blob        []byte          // The actual blob data
	BlobHash    gethcommon.Hash // Hash from blobhash opcode
	MessageRoot gethcommon.Hash // Root of cross-chain messages
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
