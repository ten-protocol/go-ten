package core

import (
	"math/big"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/trie"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/obscuro-playground/go/common"
)

// Rollup Data structure only for the internal use of the enclave since transactions are in clear
// Making changes to this struct will require GRPC + GRPC Converters regen
type Rollup struct {
	Header *common.Header

	hash atomic.Value
	// size   atomic.Value

	Transactions L2Txs
}

// Hash returns the keccak256 hash of b's header.
// The hash is computed on the first call and cached thereafter.
func (r *Rollup) Hash() common.L2RootHash {
	if hash := r.hash.Load(); hash != nil {
		return hash.(common.L2RootHash)
	}
	v := r.Header.Hash()
	r.hash.Store(v)
	return v
}

func (r *Rollup) NumberU64() uint64 { return r.Header.Number.Uint64() }
func (r *Rollup) Number() *big.Int  { return new(big.Int).Set(r.Header.Number) }

func NewHeader(parent *gethcommon.Hash, height uint64, a gethcommon.Address) *common.Header {
	parentHash := common.GenesisHash
	if parent != nil {
		parentHash = *parent
	}
	return &common.Header{
		Agg:        a,
		ParentHash: parentHash,
		Number:     big.NewInt(int64(height)),
	}
}

func NewRollupFromHeader(header *common.Header, blkHash gethcommon.Hash, txs []*common.L2Tx, nonce common.Nonce, state common.StateRoot) Rollup {
	h := common.Header{
		Agg:        header.Agg,
		ParentHash: header.ParentHash,
		L1Proof:    blkHash,
		Nonce:      nonce,
		Root:       state,
		Number:     header.Number,
	}
	transactions := make([]*common.L2Tx, len(txs))
	copy(transactions, txs)
	r := Rollup{
		Header:       &h,
		Transactions: transactions,
	}
	if len(txs) == 0 {
		h.TxHash = types.EmptyRootHash
	} else {
		h.TxHash = types.DeriveSha(types.Transactions(txs), trie.NewStackTrie(nil))
	}

	return r
}

// NewRollup - produces a new rollup. only used for genesis. todo - review
func NewRollup(blkHash gethcommon.Hash, parent *Rollup, height uint64, a gethcommon.Address, txs []*common.L2Tx, withdrawals []common.Withdrawal, nonce common.Nonce, state common.StateRoot) *Rollup {
	parentHash := common.GenesisHash
	if parent != nil {
		parentHash = parent.Hash()
	}
	h := common.Header{
		Agg:         a,
		ParentHash:  parentHash,
		L1Proof:     blkHash,
		Nonce:       nonce,
		Root:        state,
		TxHash:      types.EmptyRootHash,
		Number:      big.NewInt(int64(height)),
		Withdrawals: withdrawals,
		ReceiptHash: types.EmptyRootHash,
	}
	r := Rollup{
		Header:       &h,
		Transactions: txs,
	}
	return &r
}
