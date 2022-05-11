package core

import (
	"sync/atomic"

	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
)

// Rollup Data structure only for the internal use of the enclave since transactions are in clear
// Making changes to this struct will require GRPC + GRPC Converters regen
type Rollup struct {
	Header *nodecommon.Header

	hash atomic.Value
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

func NewHeader(parent *common.Hash, height uint64, a common.Address) *nodecommon.Header {
	parentHash := obscurocommon.GenesisHash
	if parent != nil {
		parentHash = *parent
	}
	return &nodecommon.Header{
		Agg:        a,
		ParentHash: parentHash,
		Number:     height,
	}
}

func NewRollupFromHeader(header *nodecommon.Header, blkHash common.Hash, txs []nodecommon.L2Tx, nonce obscurocommon.Nonce, state nodecommon.StateRoot) Rollup {
	h := nodecommon.Header{
		Agg:        header.Agg,
		ParentHash: header.ParentHash,
		L1Proof:    blkHash,
		Nonce:      nonce,
		State:      state,
		Number:     header.Number,
	}
	r := Rollup{
		Header:       &h,
		Transactions: txs,
	}
	return r
}

func NewRollup(blkHash common.Hash, parent *Rollup, height uint64, a common.Address, txs []nodecommon.L2Tx, withdrawals []nodecommon.Withdrawal, nonce obscurocommon.Nonce, state nodecommon.StateRoot) Rollup {
	parentHash := obscurocommon.GenesisHash
	if parent != nil {
		parentHash = parent.Hash()
	}
	h := nodecommon.Header{
		Agg:         a,
		ParentHash:  parentHash,
		L1Proof:     blkHash,
		Nonce:       nonce,
		State:       state,
		Number:      height,
		Withdrawals: withdrawals,
	}
	r := Rollup{
		Header:       &h,
		Transactions: txs,
	}
	return r
}

func (r *Rollup) ToExtRollup() nodecommon.ExtRollup {
	return nodecommon.ExtRollup{
		Header: r.Header,
		Txs:    EncryptTransactions(r.Transactions),
	}
}

func EncryptTransactions(transactions L2Txs) nodecommon.EncryptedTransactions {
	result := make([]nodecommon.EncryptedTx, 0)
	for i := range transactions {
		result = append(result, EncryptTx(&transactions[i]))
	}
	return result
}
