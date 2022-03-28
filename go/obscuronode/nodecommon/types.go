package nodecommon

import (
	"fmt"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/obscuronet/obscuro-playground/go/hashing"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"

	"github.com/ethereum/go-ethereum/common"
)

type (
	StateRoot             = string // Todo - this has to be a trie root eventually
	EncryptedTx           []byte
	L2Tx                  = types.Transaction
	EncryptedTransactions []EncryptedTx
)

// Header is a public / plaintext struct that holds common properties of the Rollup
// Making changes to this struct will require GRPC + GRPC Converters regen
type Header struct {
	ParentHash  obscurocommon.L2RootHash
	Agg         common.Address
	Nonce       obscurocommon.Nonce
	L1Proof     obscurocommon.L1RootHash // the L1 block where the Parent was published
	State       StateRoot
	Height      uint64 // the rollup height
	Withdrawals []Withdrawal
}

type Withdrawal struct {
	Amount  uint64
	Address common.Address
}

// ExtRollup is used for communication between the enclave and the outside world.
type ExtRollup struct {
	Header *Header
	Txs    EncryptedTransactions
}

// Rollup extends ExtRollup with additional fields.
// This parallels the Block/extblock split in Go Ethereum.
type Rollup struct {
	Header *Header

	hash   atomic.Value
	Height atomic.Value
	size   atomic.Value

	Transactions EncryptedTransactions
}

func (er ExtRollup) ToRollup() *Rollup {
	return &Rollup{
		Header:       er.Header,
		Transactions: er.Txs,
	}
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

// Hash returns the block hash of the header, which is simply the keccak256 hash of its
// RLP encoding.
func (h *Header) Hash() obscurocommon.L2RootHash {
	hash, err := hashing.RLPHash(h)
	if err != nil {
		// todo - log / surface these
		fmt.Printf("err hashing the l2roothash")
	}
	return hash
}
