package enclave

import (
	common2 "github.com/ethereum/go-ethereum/common"
	"simulation/common"
	oc "simulation/obscuro/common"
	"sync/atomic"
)

type Transactions []oc.L2Tx

// todo - this should become an elaborate data structure
type EnclaveSecret []byte

// Data structure only for the internal use of the enclave
// since transactions are in clear
type enclaveRollup struct {
	Header *oc.Header

	hash   atomic.Value
	Height atomic.Value
	size   atomic.Value

	Transactions Transactions
}

func NewRollup(b *common.Block, parent *oc.Rollup, a common.NodeId, txs []oc.L2Tx, withdrawals []oc.Withdrawal, nonce common.Nonce, state oc.StateRoot) oc.Rollup {
	parentHash := common2.HexToHash(oc.GenesisHash)
	if parent != nil {
		parentHash = parent.Hash()
	}
	h := oc.Header{
		Agg:         a,
		ParentHash:  parentHash,
		L1Proof:     b.Hash(),
		Nonce:       nonce,
		State:       state,
		Withdrawals: withdrawals,
	}
	r := oc.Rollup{
		Header:       &h,
		Transactions: txs,
	}
	return r
}

var GenesisRollup = NewRollup(&common.GenesisBlock, nil, 0, []oc.L2Tx{}, []oc.Withdrawal{}, common.GenerateNonce(), "")
