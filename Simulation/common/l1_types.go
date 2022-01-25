package common

import (
	"github.com/google/uuid"
	wallet_mock "simulation/wallet-mock"
	"time"
)

// Just two types of relevant L1 transactions: Deposits and Rollups
type L1TxType int64

const (
	DepositTx L1TxType = iota
	RollupTx
)

type L1Tx struct {
	Id     TxHash
	TxType L1TxType

	// if the type is rollup
	Rollup Rollup //interface{} // todo - make it a serialised data structure

	// if the type is deposit
	Amount int
	Dest   wallet_mock.Address
}

func (tx L1Tx) Hash() TxHash {
	return tx.Id
}

type Block struct {
	h            int
	root         RootHash
	Nonce        Nonce
	Miner        NodeId
	p            *Block
	CreationTime time.Time
	txs          []L1Tx
}

func (b Block) ParentBlock() *Block {
	return b.p
}
func (b Block) Parent() ChainNode {
	return *b.p
}
func (b Block) Height() int {
	return b.h
}
func (b Block) RootHash() RootHash {
	return b.root
}
func (b Block) Txs() []Tx {
	txs := make([]Tx, len(b.txs))
	// todo - inefficient
	for i, tx := range b.txs {
		txs[i] = Tx(tx)
	}
	return txs
}
func (b Block) L1Txs() []L1Tx {
	return b.txs
}

var GenesisBlock = Block{h: -1, root: uuid.New(), Nonce: 0, CreationTime: time.Now(), txs: []L1Tx{{Id: uuid.New(), TxType: RollupTx, Rollup: GenesisRollup}}}

func NewBlock(cb *Block, nonce int, m NodeId, txsCopy []L1Tx) *Block {
	return &Block{h: cb.Height() + 1, root: uuid.New(), Nonce: nonce, Miner: m, p: cb, CreationTime: time.Now(), txs: txsCopy}
}
