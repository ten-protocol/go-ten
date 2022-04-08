package ethclient

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
)

// Client defines the interface for RPC communications with the ethereum nodes
// Some of these methods are composed calls that should be decoupled in the future (ie: BlocksBetween or IsBlockAncestor)
type Client interface {
	FetchBlock(id common.Hash) (*types.Block, bool) // retrieves a block
	FetchHeadBlock() (*types.Block, uint64)         // retrieves the block at head height

	IssueTx(tx obscurocommon.EncodedL1Tx) // requests the node to broadcast a transaction

	Info() Info // retrieves the node Info

	// BlocksBetween returns the blocks between two blocks
	BlocksBetween(block *types.Block, head *types.Block) []*types.Block
	// IsBlockAncestor checks if the node recognizes a block like the ancestor
	IsBlockAncestor(block *types.Block, proof obscurocommon.L1RootHash) bool
}

// Info forces the RPC Client returns the data in the same format (independently of it's implementation)
type Info struct {
	ID common.Address
}
