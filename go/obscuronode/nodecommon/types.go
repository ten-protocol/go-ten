package nodecommon

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	obscuroCommon "github.com/obscuronet/obscuro-playground/go/common"
	"golang.org/x/crypto/sha3"
)

// Todo - this has to be a trie root eventually
type StateRoot = string
type EncryptedTx []byte

type EncryptedTransactions []EncryptedTx

// Header is in plaintext
type Header struct {
	ParentHash  obscuroCommon.L2RootHash
	Agg         common.Address
	Nonce       obscuroCommon.Nonce
	L1Proof     obscuroCommon.L1RootHash // the L1 block where the Parent was published
	State       StateRoot
	Withdrawals []Withdrawal
}

type Withdrawal struct {
	Amount  uint64
	Address common.Address
}

type Rollup struct {
	Header *Header

	hash   atomic.Value
	Height atomic.Value
	size   atomic.Value //nolint

	Transactions EncryptedTransactions
}

// ExtRollup Data structure that is used to communicate between the enclave and the outside world
type ExtRollup struct {
	Header *Header
	Txs    EncryptedTransactions
}

func (er ExtRollup) ToRollup() *Rollup {
	return &Rollup{
		Header:       er.Header,
		Transactions: er.Txs,
	}
}

// Hash returns the keccak256 hash of b's header.
// The hash is computed on the first call and cached thereafter.
func (r *Rollup) Hash() obscuroCommon.L2RootHash {
	if hash := r.hash.Load(); hash != nil {
		return hash.(obscuroCommon.L2RootHash)
	}
	v := r.Header.Hash()
	r.hash.Store(v)

	return v
}

// Hash returns the block hash of the header, which is simply the keccak256 hash of its
// RLP encoding.
func (h *Header) Hash() obscuroCommon.L2RootHash {
	return rlpHash(h)
}

/// rlpHash encodes x and hashes the encoded bytes.
func rlpHash(x interface{}) (h common.Hash) {
	sha := hasherPool.Get().(crypto.KeccakState)
	defer hasherPool.Put(sha)
	sha.Reset()
	if err := rlp.Encode(sha, x); err != nil {
		// TODO hook this up with a logger in the future - shouldn't need errors to go upstream
		// Supplying the encoder so we shouldn't see any errors
		fmt.Printf("unexpected error on the rpl encoding %v\n", err)
	}

	if _, err := sha.Read(h[:]); err != nil {
		fmt.Printf("unexpected error on the KeccakState byte read  %v\n", err)
	}

	return h
}

var hasherPool = sync.Pool{
	New: func() interface{} { return sha3.NewLegacyKeccak256() },
}
