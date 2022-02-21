package enclave

import (
	"sync/atomic"

	"github.com/otherview/obscuro-playground/common"
	oc "github.com/otherview/obscuro-playground/obscuro/common"
)

// Transfers and Withdrawals for now
type L2TxType uint64

const (
	TransferTx L2TxType = iota
	WithdrawalTx
)

// L2Tx Only in clear inside the enclave
type L2Tx struct {
	Id     common.TxHash
	TxType L2TxType
	Amount uint64
	From   common.Address
	To     common.Address
}

var GenesisRollup = NewRollup(&common.GenesisBlock, nil, 0, []L2Tx{}, []oc.Withdrawal{}, common.GenerateNonce(), "")

type Transactions []L2Tx

// todo - this should become an elaborate data structure
type EnclaveSecret []byte

// EnclaveRollup Data structure only for the internal use of the enclave since transactions are in clear
type EnclaveRollup struct {
	Header *oc.Header

	hash   atomic.Value
	Height atomic.Value
	size   atomic.Value

	Transactions Transactions
}

// Hash returns the keccak256 hash of b's header.
// The hash is computed on the first call and cached thereafter.
func (r *EnclaveRollup) Hash() common.L2RootHash {
	if hash := r.hash.Load(); hash != nil {
		return hash.(common.L2RootHash)
	}
	v := r.Header.Hash()
	r.hash.Store(v)
	return v
}

func NewRollup(b *common.Block, parent *EnclaveRollup, a common.NodeId, txs []L2Tx, withdrawals []oc.Withdrawal, nonce common.Nonce, state oc.StateRoot) EnclaveRollup {
	parentHash := oc.GenesisHash
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
	r := EnclaveRollup{
		Header:       &h,
		Transactions: txs,
	}
	return r
}

// ProofHeight - return the height of the L1 proof, or -1 - if the block is not known
// todo - find a better way. This is a workaround to handle rollups created with proofs that haven't propagated yet
func (r *EnclaveRollup) ProofHeight(l1BlockResolver common.BlockResolver) int {
	v, f := l1BlockResolver.Resolve(r.Header.L1Proof)
	if !f {
		return -1
	}
	return v.Height(l1BlockResolver)
}

func (r *EnclaveRollup) ToExtRollup() oc.ExtRollup {
	return oc.ExtRollup{
		Header: r.Header,
		Txs:    encryptTransactions(r.Transactions),
	}
}

func (r *EnclaveRollup) Proof(l1BlockResolver common.BlockResolver) *common.Block {
	v, f := l1BlockResolver.Resolve(r.Header.L1Proof)
	if !f {
		panic("Could not find proof for this rollup")
	}
	return v
}
