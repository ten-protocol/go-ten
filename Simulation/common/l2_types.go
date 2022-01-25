package common

import (
	"github.com/google/uuid"
	"simulation/wallet-mock"
	"time"
)

// Todo - this has to be a trie root eventually
type StateRoot = string

type Rollup struct {
	h            int
	root         RootHash
	Agg          NodeId
	p            *Rollup
	CreationTime time.Time
	L1Proof      *Block // the L1 block where the Parent was published
	Nonce        Nonce
	txs          []L2Tx
	State        StateRoot
}

func (r Rollup) ParentRollup() *Rollup {
	return r.p
}
func (r Rollup) Parent() ChainNode {
	return r.p
}
func (r Rollup) Height() int {
	return r.h
}
func (r Rollup) RootHash() RootHash {
	return r.root
}
func (r Rollup) Txs() []Tx {
	txs := make([]Tx, len(r.txs))
	// todo - inefficient
	for i, tx := range r.txs {
		txs[i] = Tx(tx)
	}
	return txs
}
func (r Rollup) L2Txs() []L2Tx {
	return r.txs
}

func NewRollup(b *Block, newL2Head *Rollup, a NodeId, txs []L2Tx, state StateRoot) Rollup {
	return Rollup{newL2Head.Height() + 1, uuid.New(), a, newL2Head, time.Now(), b, GenerateNonce(), txs, state}
}

// Transfers and Withdrawals for now
type L2TxType int64

const (
	TransferTx L2TxType = iota
	WithdrawalTx
)

// no signing for now
type L2Tx struct {
	Id     TxHash
	TxType L2TxType
	Amount int
	From   wallet_mock.Address
	Dest   wallet_mock.Address
}

func (tx L2Tx) Hash() TxHash {
	return tx.Id
}

var GenesisRollup = Rollup{-1,
	uuid.New(),
	-1,
	nil,
	time.Now(),
	nil,
	0,
	[]L2Tx{},
	"",
}
