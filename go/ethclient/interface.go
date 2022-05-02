package ethclient

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
)

// EthClient defines the interface for RPC communications with the ethereum nodes
// TODO Some of these methods are composed calls that should be decoupled in the future (ie: BlocksBetween or IsBlockAncestor)
type EthClient interface {
	FetchBlock(id common.Hash) (*types.Block, error)     // retrieves a block
	FetchBlockByNumber(n *big.Int) (*types.Block, error) // retrieves a block given a number - returns head block if n is nil
	FetchHeadBlock() (*types.Block, uint64)              // retrieves the block at head height

	Info() Info // retrieves the node Info

	// BlocksBetween returns the blocks between two blocks
	BlocksBetween(block *types.Block, head *types.Block) []*types.Block
	// IsBlockAncestor checks if the node recognizes a block like the ancestor
	IsBlockAncestor(block *types.Block, proof obscurocommon.L1RootHash) bool

	RPCBlockchainFeed() []*types.Block                             // returns all blocks from genesis to head
	BroadcastTx(t *obscurocommon.L1TxData)                         // issues an obscurocommon.L1TxData to the L1 network
	BlockListener() chan *types.Header                             // subscribes to new blocks and returns a listener with the blocks heads
	SubmitTransaction(tx types.TxData) (*types.Transaction, error) // submits an ethereum transaction
	FetchTxReceipt(hash common.Hash) (*types.Receipt, error)       // fetches the ethereum transaction receipt

	Stop() // tries to cleanly stop the client and release any resources
}

// Info forces the RPC EthClient to return the data in the same format (independently of it's implementation)
type Info struct {
	ID common.Address
}
