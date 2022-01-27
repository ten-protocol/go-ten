package common

import (
	"github.com/google/uuid"
	"simulation/wallet-mock"
	"time"
)

// Todo - this has to be a trie root eventually
type StateRoot = string

type Rollup struct {
	// header
	h            int
	root         RootHash
	Agg          NodeId
	p            *Rollup
	CreationTime time.Time
	L1Proof      *Block // the L1 block where the Parent was published
	Nonce        Nonce
	State        StateRoot
	Withdrawals  []Withdrawal
	// payload
	txs []L2Tx
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

func NewRollup(b *Block, newL2Head *Rollup, a NodeId, txs []L2Tx, withdrawals []Withdrawal, state StateRoot) Rollup {
	return Rollup{
		h:            newL2Head.Height() + 1,
		root:         uuid.New(),
		Agg:          a,
		p:            newL2Head,
		CreationTime: time.Now(),
		L1Proof:      b,
		Nonce:        GenerateNonce(),
		State:        state,
		Withdrawals:  withdrawals,
		txs:          txs,
	}
}

// Transfers and Withdrawals for now
type L2TxType int64

const (
	TransferTx L2TxType = iota
	WithdrawalTx
)

type Withdrawal struct {
	Amount  int
	Address wallet_mock.Address
}

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

var GenesisRollup = Rollup{
	h:            GenesisHeight,
	root:         uuid.New(),
	Agg:          -1,
	CreationTime: time.Now(),
	Withdrawals:  []Withdrawal{},
	txs:          []L2Tx{},
}
