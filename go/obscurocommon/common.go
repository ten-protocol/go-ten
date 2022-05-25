package obscurocommon

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// todo - this should be configured  - move away from the constant
var ChainID = big.NewInt(777) // The unique ID for the Obscuro chain.

const (
	L2GenesisHeight = uint64(0)
	L1GenesisHeight = uint64(0)
)

// HeightCommittedBlocks is the number of blocks deep a transaction must be to be considered safe from reorganisations.
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
