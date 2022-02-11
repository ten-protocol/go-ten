package enclave

import (
	"simulation/common"
)

// RollupResolver -database of rollups indexed by the root hash
type RollupResolver interface {
	FetchRollup(hash common.L2RootHash) *Rollup
}

func (r Rollup) Parent(db RollupResolver) *Rollup {
	return db.FetchRollup(r.Header.ParentHash)
}
