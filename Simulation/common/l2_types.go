package common

import (
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/google/uuid"
	"simulation/wallet-mock"
	"sync"
	"time"
)

// Todo - this has to be a trie root eventually
type StateRoot = string
type EncodedL2Tx []byte
type EncodedRollup []byte

type Rollup struct {
	// header
	H            uint32
	RootHash     RootHash
	Agg          NodeId
	ParentHash   RootHash
	CreationTime time.Time
	L1Proof      RootHash // the L1 block where the Parent was published
	Nonce        Nonce
	State        StateRoot
	Withdrawals  []Withdrawal

	// payload - move to body
	Transactions []L2Tx
}

var rollupCache = make(map[RootHash]Rollup)
var rcm = sync.RWMutex{}

func (r Rollup) ParentRollup() *Rollup {
	rcm.RLock()
	defer rcm.RUnlock()
	rollup, found := rollupCache[r.ParentHash]
	if !found {
		panic("could not find rollup")
	}
	return &rollup
}
func (r Rollup) Parent() ChainNode {
	return r.ParentRollup()
}
func (r Rollup) Height() uint32 {
	return r.H
}
func (r Rollup) Root() RootHash {
	return r.RootHash
}
func (r Rollup) Txs() []Tx {
	txs := make([]Tx, len(r.Transactions))
	// todo - inefficient
	for i, tx := range r.Transactions {
		txs[i] = Tx(tx)
	}
	return txs
}
func (r Rollup) L2Txs() []L2Tx {
	return r.Transactions
}

func (r Rollup) Proof() Block {
	rbm.RLock()
	defer rbm.RUnlock()
	block, f := blockCache[r.L1Proof]
	if !f {
		panic("Couldn't find block")
	}
	return block
}

func NewRollup(b *Block, parent *Rollup, a NodeId, txs []L2Tx, withdrawals []Withdrawal, state StateRoot) Rollup {
	rootHash := uuid.New()
	parentHash := rootHash
	height := GenesisHeight
	if parent != nil {
		parentHash = parent.RootHash
		height = parent.H + 1
	}
	r := Rollup{
		H:            height,
		RootHash:     rootHash,
		Agg:          a,
		ParentHash:   parentHash,
		CreationTime: time.Now(),
		L1Proof:      b.RootHash,
		Nonce:        GenerateNonce(),
		State:        state,
		Withdrawals:  withdrawals,
		Transactions: txs,
	}
	rcm.Lock()
	rollupCache[rootHash] = r
	rcm.Unlock()
	return r
}

// Transfers and Withdrawals for now
type L2TxType uint64

const (
	TransferTx L2TxType = iota
	WithdrawalTx
)

type Withdrawal struct {
	Amount  uint64
	Address wallet_mock.Address
}

// no signing for now
type L2Tx struct {
	Id     TxHash
	TxType L2TxType
	Amount uint64
	From   wallet_mock.Address
	Dest   wallet_mock.Address
}

func (tx L2Tx) Hash() TxHash {
	return tx.Id
}

var GenesisRollup = NewRollup(&GenesisBlock, nil, 0, []L2Tx{}, []Withdrawal{}, "")
var GenesisTx = L1Tx{Id: uuid.New(), TxType: RollupTx, Rollup: GenesisRollup}

func (r Rollup) Encode() (EncodedRollup, error) {
	return rlp.EncodeToBytes(r)
}

func (r EncodedRollup) Decode() (Rollup, error) {
	bl := Rollup{}
	err := rlp.DecodeBytes(r, &bl)
	return bl, err
}

func (tx L2Tx) EncodeBytes() (EncodedL2Tx, error) {
	return rlp.EncodeToBytes(tx)
}

func (tx EncodedL2Tx) DecodeBytes() (L2Tx, error) {
	tx1 := L2Tx{}
	err := rlp.DecodeBytes(tx, &tx1)
	return tx1, err
}
