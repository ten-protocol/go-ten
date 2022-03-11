package enclave

import (
	"sync/atomic"

	"github.com/ethereum/go-ethereum/rlp"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	common2 "github.com/obscuronet/obscuro-playground/go/common"
	oc "github.com/obscuronet/obscuro-playground/go/obscuronode/common"
)

// L2TxType indicates the type of L2 transaction - either a transfer or a withdrawal for now
type L2TxType uint8

const (
	TransferTx L2TxType = iota
	WithdrawalTx
)

// L2TxData is the Obscuro transaction data that will be stored encoded in the types.Transaction data field.
type L2TxData struct {
	Type   L2TxType
	From   common.Address
	To     common.Address
	Amount uint64
}

type L2Tx = types.Transaction

// TxData returns the decoded L2 data stored in the transaction's data field.
func TxData(tx *L2Tx) L2TxData {
	data := L2TxData{}

	err := rlp.DecodeBytes(tx.Data(), &data)
	if err != nil {
		// TODO - Surface this error properly.
		panic(err)
	}

	return data
}

var GenesisRollup = NewRollup(common2.GenesisBlock, nil, common.HexToAddress("0x0"), []L2Tx{}, []oc.Withdrawal{}, common2.GenerateNonce(), "")

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

func NewRollup(b *types.Block, parent *Rollup, a common.Address, txs []L2Tx, withdrawals []oc.Withdrawal, nonce common2.Nonce, state oc.StateRoot) Rollup {
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
	v, f := l1BlockResolver.ResolveBlock(r.Header.L1Proof)
	if !f {
		return -1
	}
	return l1BlockResolver.HeightBlock(v)
}

func (r *Rollup) ToExtRollup() oc.ExtRollup {
	return oc.ExtRollup{
		Header: r.Header,
		Txs:    encryptTransactions(r.Transactions),
	}
}

func (r *Rollup) Proof(l1BlockResolver common2.BlockResolver) *types.Block {
	v, f := l1BlockResolver.ResolveBlock(r.Header.L1Proof)
	if !f {
		panic("Could not find proof for this rollup")
	}
	return v
}
