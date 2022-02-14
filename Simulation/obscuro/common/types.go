package common

import (
	c "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"golang.org/x/crypto/sha3"
	"io"
	"simulation/common"
	"sync"
	"sync/atomic"
)

// Todo - this has to be a trie root eventually
type StateRoot = string
type EncryptedTx []byte

type EncryptedTransactionBlob []byte

// The header is in plaintext
type Header struct {
	ParentHash  common.L2RootHash
	Agg         common.NodeId
	Nonce       common.Nonce
	L1Proof     common.L1RootHash // the L1 block where the Parent was published
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
	size   atomic.Value

	//Transactions EncryptedTransactionBlob
	Transactions []L2Tx
}

//func (r *Rollup) Txs() Transactions {
//	if txs := r.txs.Load(); txs != nil {
//		return txs.(Transactions)
//	}
//	v := decrypt(r.Transactions, key)
//	r.txs.Store(v)
//	return v
//
//}

// Data structure that is used to communicate between the enclave and the outside world
type ExtRollup struct {
	Header *Header
	//Txs    EncryptedTransactionBlob
	Txs []L2Tx
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

// Transfers and Withdrawals for now
type L2TxType uint64

const (
	TransferTx L2TxType = iota
	WithdrawalTx
)

// todo - signing
type L2Tx struct {
	Id     common.TxHash
	TxType L2TxType
	Amount uint64
	From   common.Address
	To     common.Address
}

const GenesisHash = "1000000000000000000000000000000000000000000000000000000000000000"

func (r Rollup) Proof(l1BlockResolver common.BlockResolver) *common.Block {
	v, f := l1BlockResolver.Resolve(r.Header.L1Proof)
	if !f {
		panic("Could not find proof for this rollup")
	}
	return v
}

// ProofHeight - return the height of the L1 proof, or -1 - if the block is not known
//todo - find a better way. This is a workaround to handle rollups created with proofs that haven't propagated yet
func (r Rollup) ProofHeight(l1BlockResolver common.BlockResolver) int {
	v, f := l1BlockResolver.Resolve(r.Header.L1Proof)
	if !f {
		return -1
	}
	return v.Height(l1BlockResolver)
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

// Hash returns the block hash of the header, which is simply the keccak256 hash of its
// RLP encoding.
func (h *Header) Hash() common.L2RootHash {
	return rlpHash(h)
}

// rlpHash encodes x and hashes the encoded bytes.
func rlpHash(x interface{}) (h common.L2RootHash) {
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
