package common

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
)

type NodeId uint64

// todo - use proper crypto
//type Address = uuid.UUID
type Address = uint32

const L2GenesisHeight uint32 = 0
const L1GenesisHeight = 0

// Number of blocks deep a transaction must be before being considered safe from reorganisations.
const HeightCommittedBlocks = 20

type L1RootHash = common.Hash

type L2RootHash = uuid.UUID
type TxHash = uuid.UUID

type Nonce = uint64

type EncodedRollup []byte
