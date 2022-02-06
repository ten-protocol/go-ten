package common

import (
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/google/uuid"
	"sync"
	"time"
)

// L1TxType - Just two types of relevant L1 transactions: Deposits and Rollups
// this does not actually exist in the real implementation
type L1TxType uint8

type EncodedL1Tx []byte
type EncodedBlock []byte

const (
	DepositTx L1TxType = iota
	RollupTx
)

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

func (tx L1Tx) Hash() TxHash {
	return tx.Id
}

type Block struct {
	H            uint32
	RootHash     RootHash
	Nonce        Nonce
	Miner        NodeId
	ParentHash   RootHash
	CreationTime time.Time
	Transactions []L1Tx
}

var blockCache = make(map[RootHash]Block)
var rbm = sync.RWMutex{}

func (b Block) ParentBlock() *Block {
	rbm.RLock()
	defer rbm.RUnlock()
	block, found := blockCache[b.ParentHash]
	if !found {
		panic("could not find block")
	}
	return &block
}
func (b Block) Parent() ChainNode {
	return b.ParentBlock()
}
func (b Block) Height() uint32 {
	return b.H
}
func (r Block) Root() RootHash {
	return r.RootHash
}

func (b Block) Txs() []Tx {
	txs := make([]Tx, len(b.Transactions))
	// todo - inefficient
	for i, tx := range b.Transactions {
		txs[i] = Tx(tx)
	}
	return txs
}
func (b Block) L1Txs() []L1Tx {
	return b.Transactions
}

func (b Block) Encode() (EncodedBlock, error) {
	return rlp.EncodeToBytes(b)
}

func (b EncodedBlock) Decode() (Block, error) {
	bl := Block{}
	err := rlp.DecodeBytes(b, &bl)
	return bl, err
}

func (tx L1Tx) Encode() (EncodedL1Tx, error) {
	return rlp.EncodeToBytes(tx)
}

func (tx EncodedL1Tx) Decode() (L1Tx, error) {
	tx1 := L1Tx{}
	err := rlp.DecodeBytes(tx, &tx1)
	return tx1, err
}

func NewBlock(parent *Block, nonce uint64, m NodeId, txs []L1Tx) Block {
	rootHash := uuid.New()
	parentHash := rootHash
	height := GenesisHeight
	if parent != nil {
		parentHash = parent.RootHash
		height = parent.H + 1
	}
	b := Block{
		H:            height,
		RootHash:     rootHash,
		Nonce:        nonce,
		Miner:        m,
		ParentHash:   parentHash,
		CreationTime: time.Now(),
		Transactions: txs,
	}
	rbm.Lock()
	blockCache[rootHash] = b
	rbm.Unlock()
	return b
}

var GenesisBlock = NewBlock(nil, 0, 0, []L1Tx{})
