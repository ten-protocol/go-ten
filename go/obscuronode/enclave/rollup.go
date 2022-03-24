package enclave

import (
	"sync/atomic"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
)

var GenesisRollup = NewRollup(obscurocommon.GenesisBlock, nil, common.HexToAddress("0x0"), []nodecommon.L2Tx{}, []nodecommon.Withdrawal{}, obscurocommon.GenerateNonce(), "")

// Rollup Data structure only for the internal use of the enclave since transactions are in clear
type Rollup struct {
	Header *nodecommon.Header

	hash   atomic.Value
	Height atomic.Value
	// size   atomic.Value

	Transactions L2Txs
}

// Hash returns the keccak256 hash of b's header.
// The hash is computed on the first call and cached thereafter.
func (r *Rollup) Hash() obscurocommon.L2RootHash {
	if hash := r.hash.Load(); hash != nil {
		return hash.(obscurocommon.L2RootHash)
	}
	v := r.Header.Hash()
	r.hash.Store(v)
	return v
}

func NewRollup(b *types.Block, parent *Rollup, a common.Address, txs []nodecommon.L2Tx, withdrawals []nodecommon.Withdrawal, nonce obscurocommon.Nonce, state nodecommon.StateRoot) *Rollup {
	parentHash := obscurocommon.GenesisHash
	if parent != nil {
		parentHash = parent.Hash()
	}
	h := nodecommon.Header{
		Agg:         a,
		ParentHash:  parentHash,
		L1Proof:     b.Hash(),
		Nonce:       nonce,
		State:       state,
		Withdrawals: withdrawals,
	}
	r := Rollup{
		Header:       &h,
		Transactions: txs,
	}
	return &r
}

// ProofHeight - return the height of the L1 proof, or GenesisHeight-1 - if the block is not known
// todo - find a better way. This is a workaround to handle rollups created with proofs that haven't propagated yet
func (r *Rollup) ProofHeight(l1BlockResolver BlockResolver) uint64 {
	v, f := l1BlockResolver.FetchBlock(r.Header.L1Proof)
	if !f {
		return obscurocommon.L1GenesisHeight - 1
	}
	return l1BlockResolver.HeightBlock(v)
}

func (r *Rollup) ToExtRollup() nodecommon.ExtRollup {
	return nodecommon.ExtRollup{
		Header: r.Header,
		Txs:    encryptTransactions(r.Transactions),
	}
}

func (r *Rollup) Proof(l1BlockResolver BlockResolver) *types.Block {
	v, f := l1BlockResolver.FetchBlock(r.Header.L1Proof)
	if !f {
		panic("Could not find proof for this rollup")
	}
	return v
}
