package common

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
)

type NodeID common.Address

const (
	L2GenesisHeight = 0
	L1GenesisHeight = 0
)

// Number of blocks deep a transaction must be before being considered safe from reorganisations.
const HeightCommittedBlocks = 20

type L1RootHash = common.Hash

type (
	L2RootHash = common.Hash
	TxHash     = common.Hash
	L2TxHash   = uuid.UUID
)

type Nonce = uint64

type EncodedRollup []byte

type NotifyNewBlock interface {
	RPCNewHead(b EncodedBlock, p EncodedBlock)
	RPCNewFork(b []EncodedBlock)
}

type L1Node interface {
	RPCBlockchainFeed() []*Block
	BroadcastTx(t EncodedL1Tx)
}
