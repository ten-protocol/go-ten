package core

import (
	"math/big"
	"sync/atomic"

	"github.com/obscuronet/go-obscuro/go/common"
)

type Rollup struct {
	Header      *common.RollupHeader
	BatchHashes []common.L2RootHash // todo - joel - references?
	hash        atomic.Value
}

// Hash returns the keccak256 hash of b's header.
// The hash is computed on the first call and cached thereafter.
func (r *Rollup) Hash() *common.L2RootHash {
	// Temporarily disabling the caching of the hash because it's causing bugs.
	// Transforming a Rollup to an ExtRollup and then back to a Rollup will generate a different hash if caching is enabled.
	// Todo - re-enable
	//if hash := r.hash.Load(); hash != nil {
	//	return hash.(common.L2RootHash)
	//}
	v := r.Header.Hash()
	r.hash.Store(v)
	return &v
}

func (r *Rollup) NumberU64() uint64 { return r.Header.Number.Uint64() }
func (r *Rollup) Number() *big.Int  { return new(big.Int).Set(r.Header.Number) }

// IsGenesis indicates whether the rollup is the genesis rollup.
// TODO - #718 - Change this to a check against a hardcoded genesis hash.
func (r *Rollup) IsGenesis() bool {
	return r.Header.Number.Cmp(big.NewInt(int64(common.L2GenesisHeight))) == 0
}

func (r *Rollup) ToExtRollup() *common.ExtRollup {
	return &common.ExtRollup{
		Header:      r.Header,
		BatchHashes: r.BatchHashes,
	}
}

func ToRollup(encryptedRollup *common.ExtRollup) *Rollup {
	return &Rollup{
		Header: encryptedRollup.Header,
	}
}
