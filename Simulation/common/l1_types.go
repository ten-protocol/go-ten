package common

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"golang.org/x/crypto/sha3"
	"io"
	"sync"
	"sync/atomic"
)

// Todo - replace these data structures with the actual Ethereum structures from geth

// L1TxType - Just two types of relevant L1 transactions: Deposits and Rollups
// this does not actually exist in the real implementation
type L1TxType uint8

const (
	DepositTx L1TxType = iota
	RollupTx
)

//todo - replace with the ethereum Transaction type
type L1Tx struct {
	Id     TxHash
	TxType L1TxType

	// if the type is rollup
	//todo -payload
	Rollup EncodedRollup

	// if the type is deposit
	Amount uint64
	Dest   Address
}
type EncodedL1Tx []byte
type Transactions []*L1Tx

type Header struct {
	ParentHash L1RootHash
	Miner      NodeId // this is actually coinbase
	Nonce      Nonce
}

// todo - split into header and payload and then replace with the ethereum Block
type Block struct {
	Header       *Header
	Transactions Transactions

	//ReceiveTime time.Time

	hash   atomic.Value
	height atomic.Value
	size   atomic.Value
}

// the encoded version of an ExtBlock
type EncodedBlock []byte

const GenesisHash = "0000000000000000000000000000000000000000000000000000000000000000"

func NewBlock(parent *Block, nonce uint64, m NodeId, txs []*L1Tx) Block {
	var parentHash = common.HexToHash(GenesisHash)
	if parent != nil {
		parentHash = parent.Hash()
	}
	h := Header{
		Nonce:      nonce,
		Miner:      m,
		ParentHash: parentHash,
	}
	b := Block{
		Header:       &h,
		Transactions: txs,
	}
	return b
}

var GenesisBlock = NewBlock(nil, 0, 0, []*L1Tx{})

// "external" block encoding. used for eth protocol, etc.
type ExtBlock struct {
	Header *Header
	Txs    []*L1Tx
}

// Hash returns the keccak256 hash of b's header.
// The hash is computed on the first call and cached thereafter.
func (b *Block) Hash() L1RootHash {
	if hash := b.hash.Load(); hash != nil {
		return hash.(L1RootHash)
	}
	v := b.Header.Hash()
	b.hash.Store(v)
	return v
}

// Hash returns the block hash of the header, which is simply the keccak256 hash of its
// RLP encoding.
func (h *Header) Hash() L1RootHash {
	return rlpHash(h)
}

//func (h *Header) Height() uint32 {
//	return rlpHash(h)
//}

// rlpHash encodes x and hashes the encoded bytes.
func rlpHash(x interface{}) (h L1RootHash) {
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

func (eb ExtBlock) ToBlock() *Block {
	return &Block{
		Header:       eb.Header,
		Transactions: eb.Txs,
	}
}
func (b Block) ToExtBlock() ExtBlock {
	return ExtBlock{
		Header: b.Header,
		Txs:    b.Transactions,
	}
}

// DecodeRLP decodes the Ethereum
func (b *Block) DecodeRLP(s *rlp.Stream) error {
	var eb ExtBlock
	_, size, _ := s.Kind()
	if err := s.Decode(&eb); err != nil {
		return err
	}
	b.Header, b.Transactions = eb.Header, eb.Txs
	b.size.Store(common.StorageSize(rlp.ListSize(size)))
	return nil
}

// EncodeRLP serializes b into the Ethereum RLP block format.
func (b Block) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, ExtBlock{
		Header: b.Header,
		Txs:    b.Transactions,
	})
}

func (b *Block) Height(resolver BlockResolver) int {
	if height := b.height.Load(); height != nil {
		return height.(int)
	}
	if b.Hash() == GenesisBlock.Hash() {
		b.height.Store(L1GenesisHeight)
		return L1GenesisHeight
	}

	p, f := b.Parent(resolver)
	if !f {
		panic("wtf")
	}
	v := p.Height(resolver) + 1
	b.height.Store(v)
	return v
}
