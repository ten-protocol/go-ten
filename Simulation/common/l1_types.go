package common

import (
	"github.com/google/uuid"
	"time"
)

// L1TxType - Just two types of relevant L1 transactions: Deposits and Rollups
// this does not actually exist in the real implementation
type L1TxType uint8

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
type EncodedL1Tx []byte

// todo - split into header and payload
type Block struct {
	Height       uint32
	RootHash     RootHash
	Nonce        Nonce
	Miner        NodeId
	ParentHash   RootHash
	CreationTime time.Time
	Transactions []L1Tx
}
type EncodedBlock []byte

func NewBlock(parent *Block, nonce uint64, m NodeId, txs []L1Tx) Block {
	rootHash := uuid.New()
	parentHash := rootHash
	height := GenesisHeight
	if parent != nil {
		parentHash = parent.RootHash
		height = parent.Height + 1
	}
	b := Block{
		Height:       height,
		RootHash:     rootHash,
		Nonce:        nonce,
		Miner:        m,
		ParentHash:   parentHash,
		CreationTime: time.Now(),
		Transactions: txs,
	}
	return b
}

var GenesisBlock = NewBlock(nil, 0, 0, []L1Tx{})
