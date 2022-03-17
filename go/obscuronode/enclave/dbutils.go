package enclave

import (
	"fmt"

	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
)

func ParentRollup(db DB, r *Rollup) *Rollup {
	parent, found := db.FetchRollup(r.Header.ParentHash)
	if !found {
		panic(fmt.Sprintf("Could not find rollup: r_%s", r.Hash()))
	}
	return parent
}

func HeightRollup(db DB, r *Rollup) int {
	if height := r.Height.Load(); height != nil {
		return height.(int)
	}
	if r.Hash() == GenesisRollup.Hash() {
		r.Height.Store(obscurocommon.L2GenesisHeight)
		return obscurocommon.L2GenesisHeight
	}
	v := HeightRollup(db, ParentRollup(db, r)) + 1
	r.Height.Store(v)
	return v
}
