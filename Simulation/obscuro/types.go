package obscuro

import (
	"github.com/google/uuid"
	"simulation/common"
	"time"
)

// Todo - this has to be a trie root eventually
type StateRoot = string
type EncodedL2Tx []byte

type Rollup struct {
	// header
	Height       uint32
	RootHash     common.L2RootHash
	Agg          common.NodeId
	ParentHash   common.L2RootHash
	CreationTime time.Time
	L1Proof      common.L1RootHash // the L1 block where the Parent was published
	Nonce        common.Nonce
	State        StateRoot
	Withdrawals  []Withdrawal

	// payload - todo move to body
	Transactions []L2Tx
}

// Transfers and Withdrawals for now
type L2TxType uint64

const (
	TransferTx L2TxType = iota
	WithdrawalTx
)

type Withdrawal struct {
	Amount  uint64
	Address common.Address
}

// todo - signing
type L2Tx struct {
	Id     common.TxHash
	TxType L2TxType
	Amount uint64
	From   common.Address
	To     common.Address
}

var GenesisRollup = NewRollup(&common.GenesisBlock, nil, 0, []L2Tx{}, []Withdrawal{}, common.GenerateNonce(), "")
var encodedGenesis, _ = GenesisRollup.Encode()
var GenesisTx = common.L1Tx{Id: uuid.New(), TxType: common.RollupTx, Rollup: encodedGenesis}

func (r Rollup) Proof(l1BlockResolver common.BlockResolver) common.Block {
	v, f := l1BlockResolver.Resolve(r.L1Proof)
	if !f {
		panic("Could not find proof for this rollup")
	}
	return v
}

// ProofHeight - return the height of the L1 proof, or 0 - if the block is not known
//todo - find a better way. This is a workaround to handle rollups created with proofs that haven't propagated yet
func (r Rollup) ProofHeight(l1BlockResolver common.BlockResolver) int {
	v, f := l1BlockResolver.Resolve(r.L1Proof)
	if !f {
		return 0
	}
	return v.Height(l1BlockResolver)
}

func NewRollup(b *common.Block, parent *Rollup, a common.NodeId, txs []L2Tx, withdrawals []Withdrawal, nonce common.Nonce, state StateRoot) Rollup {
	rootHash := uuid.New()
	parentHash := rootHash
	height := common.L2GenesisHeight
	if parent != nil {
		parentHash = parent.RootHash
		height = parent.Height + 1
	}
	r := Rollup{
		Height:       height,
		RootHash:     rootHash,
		Agg:          a,
		ParentHash:   parentHash,
		CreationTime: time.Now(),
		L1Proof:      b.Hash(),
		Nonce:        nonce,
		State:        state,
		Withdrawals:  withdrawals,
		Transactions: txs,
	}
	return r
}
