package enclave

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
	"math/rand"
	"sync/atomic"

	common2 "github.com/obscuronet/obscuro-playground/go/common"
	oc "github.com/obscuronet/obscuro-playground/go/obscuronode/common"
)

// Transfers and Withdrawals for now
type L2TxType uint64

const (
	TransferTx L2TxType = iota
	WithdrawalTx
)

// L2Tx Only in clear inside the enclave
type L2Tx struct {
	Tx *types.Transaction
	// TODO - Joel - Explain why these wrapper fields are needed.
	From common.Address
	Type L2TxType
}

// TODO - Joel - Describe.
func L2TxNew(to common.Address, value uint64, from common.Address, txType L2TxType) L2Tx {
	// TODO - Joel - Create the tx data.
	tx := types.NewTx(&types.LegacyTx{
		// TODO - Joel - Review this taking of a reference.
		To: &to,
		// TODO - Joel - Review this conversion.
		Value: big.NewInt(int64(value)),
		// TODO - Joel - Choose a better nonce. Just using this to avoid duplicates.
		Nonce: rand.Uint64(),
	})
	return L2Tx{tx, from, txType}
}

var GenesisRollup = NewRollup(&common2.GenesisBlock, nil, 0, []L2Tx{}, []oc.Withdrawal{}, common2.GenerateNonce(), "")

type Transactions []L2Tx

// todo - this should become an elaborate data structure
type SharedEnclaveSecret []byte

// Rollup Data structure only for the internal use of the enclave since transactions are in clear
type Rollup struct {
	Header *oc.Header

	hash   atomic.Value
	Height atomic.Value
	// size   atomic.Value

	Transactions Transactions
}

// Hash returns the keccak256 hash of b's header.
// The hash is computed on the first call and cached thereafter.
func (r *Rollup) Hash() common2.L2RootHash {
	if hash := r.hash.Load(); hash != nil {
		return hash.(common2.L2RootHash)
	}
	v := r.Header.Hash()
	r.hash.Store(v)
	return v
}

func NewRollup(b *common2.Block, parent *Rollup, a common2.NodeID, txs []L2Tx, withdrawals []oc.Withdrawal, nonce common2.Nonce, state oc.StateRoot) Rollup {
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
	r := Rollup{
		Header:       &h,
		Transactions: txs,
	}
	return r
}

// ProofHeight - return the height of the L1 proof, or -1 - if the block is not known
// todo - find a better way. This is a workaround to handle rollups created with proofs that haven't propagated yet
func (r *Rollup) ProofHeight(l1BlockResolver common2.BlockResolver) int {
	v, f := l1BlockResolver.Resolve(r.Header.L1Proof)
	if !f {
		return -1
	}
	return v.Height(l1BlockResolver)
}

func (r *Rollup) ToExtRollup() oc.ExtRollup {
	return oc.ExtRollup{
		Header: r.Header,
		Txs:    encryptTransactions(r.Transactions),
	}
}

func (r *Rollup) Proof(l1BlockResolver common2.BlockResolver) *common2.Block {
	v, f := l1BlockResolver.Resolve(r.Header.L1Proof)
	if !f {
		panic("Could not find proof for this rollup")
	}
	return v
}
