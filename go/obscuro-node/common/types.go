package common

import (
	common2 "github.com/obscuronet/obscuro-playground/go/common"
	"io"
	"sync"
	"sync/atomic"

	c "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"golang.org/x/crypto/sha3"
)

var GenesisHash = c.HexToHash("1000000000000000000000000000000000000000000000000000000000000000")

// Todo - this has to be a trie root eventually
type StateRoot = string
type EncryptedTx []byte

type EncryptedTransactions []EncryptedTx

// Header is in plaintext
type Header struct {
	ParentHash  common2.L2RootHash
	Agg         common2.NodeId
	Nonce       common2.Nonce
	L1Proof     common2.L1RootHash // the L1 block where the Parent was published
	State       StateRoot
	Withdrawals []Withdrawal
}

type Withdrawal struct {
	Amount  uint64
	Address common2.Address
}

type Rollup struct {
	Header *Header

	hash   atomic.Value
	Height atomic.Value
	size   atomic.Value

	Transactions EncryptedTransactions
}

// ExtRollup Data structure that is used to communicate between the enclave and the outside world
type ExtRollup struct {
	Header *Header
	Txs    EncryptedTransactions
}

func (eb ExtRollup) ToRollup() *Rollup {
	return &Rollup{
		Header:       eb.Header,
		Transactions: eb.Txs,
	}
}

func (b Rollup) ToExtRollup() ExtRollup {
	return ExtRollup{
		Header: b.Header,
		Txs:    b.Transactions,
	}
}

func (r Rollup) Proof(l1BlockResolver common2.BlockResolver) *common2.Block {
	v, f := l1BlockResolver.Resolve(r.Header.L1Proof)
	if !f {
		panic("Could not find proof for this rollup")
	}
	return v
}

// ProofHeight - return the height of the L1 proof, or -1 - if the block is not known
// todo - find a better way. This is a workaround to handle rollups created with proofs that haven't propagated yet
func (r Rollup) ProofHeight(l1BlockResolver common2.BlockResolver) int {
	v, f := l1BlockResolver.Resolve(r.Header.L1Proof)
	if !f {
		return -1
	}
	return v.Height(l1BlockResolver)
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

// Hash returns the block hash of the header, which is simply the keccak256 hash of its
// RLP encoding.
func (h *Header) Hash() common2.L2RootHash {
	return rlpHash(h)
}

// rlpHash encodes x and hashes the encoded bytes.
func rlpHash(x interface{}) (h common2.L2RootHash) {
	sha := hasherPool.Get().(crypto.KeccakState)
	defer hasherPool.Put(sha)
	sha.Reset()
	rlp.Encode(sha, x)
	sha.Read(h[:])
	return h
}

var hasherPool = sync.Pool{
	New: func() interface{} { return sha3.NewLegacyKeccak256() },
}

// DecodeRLP decodes the Ethereum
func (b *Rollup) DecodeRLP(s *rlp.Stream) error {
	var eb ExtRollup
	_, size, _ := s.Kind()
	if err := s.Decode(&eb); err != nil {
		return err
	}
	b.Header, b.Transactions = eb.Header, eb.Txs
	b.size.Store(c.StorageSize(rlp.ListSize(size)))
	return nil
}

// EncodeRLP serializes b into the Ethereum RLP block format.
func (b Rollup) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, b.ToExtRollup())
}
