package obscurocommon

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const (
	L2GenesisHeight = uint64(0)
	L1GenesisHeight = uint64(0)
)

// Number of blocks deep a transaction must be before being considered safe from reorganisations.
const HeightCommittedBlocks = 20

type (
	L1RootHash = common.Hash
	L2RootHash = common.Hash
	TxHash     = common.Hash
)

type Nonce = uint64

type EncodedRollup []byte

type NotifyNewBlock interface {
	RPCNewHead(b EncodedBlock, p EncodedBlock)
	RPCNewFork(b []EncodedBlock)
}

type L1Node interface {
	RPCBlockchainFeed() []*types.Block
	BroadcastTx(t EncodedL1Tx)
}
