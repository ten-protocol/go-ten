package common

import (
	"fmt"
	"io"
	"sync"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"golang.org/x/crypto/sha3"
)

// Todo - replace these data structures with the actual Ethereum structures from geth

// L1TxType - Just two types of relevant L1 transactions: Deposits and Rollups
// this does not actually exist in the real implementation
type L1TxType uint8

const (
	DepositTx L1TxType = iota
	RollupTx
	StoreSecretTx
	RequestSecretTx
)

// todo - replace with the ethereum Transaction type
// all the fields are placeholders for arguments sent to the management contract
type L1Tx struct {
	ID     TxHash
	TxType L1TxType

	// if the type is rollup
	// todo -payload
	Rollup EncodedRollup

	Secret      EncryptedSharedEnclaveSecret
	Attestation AttestationReport

	// if the type is deposit
	Amount uint64
	Dest   common.Address
}

type (
	EncodedL1Tx  []byte
	Transactions []*L1Tx
)

type Header struct {
	ParentHash L1RootHash
	Miner      NodeID // this is actually coinbase
	Nonce      Nonce
}

// todo - replace with the ethereum Block
type Block struct {
	Header       *Header
	Transactions Transactions

	// ReceiveTime time.Time

	hash   atomic.Value
	height atomic.Value
	size   atomic.Value
}

// the encoded version of an ExtBlock
type EncodedBlock []byte

var GenesisHash = common.HexToHash("0000000000000000000000000000000000000000000000000000000000000000")

func NewBlock(parent *Block, nonce uint64, nodeID NodeID, txs []*L1Tx) Block {
	parentHash := GenesisHash
	if parent != nil {
		parentHash = parent.Hash()
	}

	header := Header{
		Nonce:      nonce,
		Miner:      nodeID,
		ParentHash: parentHash,
	}

	return Block{
		Header:       &header,
		Transactions: txs,
	}
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

	v, _ := b.Header.Hash()
	b.hash.Store(v)

	return v
}

// Hash returns the block hash of the header, which is simply the keccak256 hash of its
// RLP encoding.
func (h *Header) Hash() (L1RootHash, error) {
	return rlpHash(h)
}

//func (h *Header) Height() uint32 {
//	return rlpHash(h)
//}

// rlpHash encodes x and hashes the encoded bytes.
func rlpHash(value interface{}) (L1RootHash, error) {
	var hash L1RootHash

	sha := hasherPool.Get().(crypto.KeccakState)
	defer hasherPool.Put(sha)
	sha.Reset()
	err := rlp.Encode(sha, value)
	if err != nil {
		return hash, fmt.Errorf("unable to encode Value. %w", err)
	}

	_, err = sha.Read(hash[:])
	if err != nil {
		return hash, fmt.Errorf("unable to read encoded value. %w", err)
	}

	return hash, nil
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
	var extBlock ExtBlock

	_, size, _ := s.Kind()

	if err := s.Decode(&extBlock); err != nil {
		return err
	}

	b.Header, b.Transactions = extBlock.Header, extBlock.Txs
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
		panic(fmt.Sprintf("Could not find parent of block b_%s", b.Hash()))
	}

	v := p.Height(resolver) + 1
	b.height.Store(v)

	return v
}

type EncryptedSharedEnclaveSecret []byte

type AttestationReport struct {
	Owner NodeID
	// todo public key
	// hash of code
	// other stuff
}
