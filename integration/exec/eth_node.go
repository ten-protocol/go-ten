package exec

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/go/ethclient"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
)

// EthNode is the interface for handling ethereum nodes lifecycle and properties
// It's a wrapper for mocks and real nodes - allows for simulations to seamlessly use both types of nodes
type EthNode interface {
	Start()
	Stop()
	Client() ethclient.Client
	KnownPeers(nodes []EthNode)
	IsBlockAncestor(block *types.Block, proof obscurocommon.L1RootHash) bool
	FetchHeadBlock() (*types.Block, uint64)
	Info() EthNodeInfo
	BlocksBetween(block *types.Block, head *types.Block) []*types.Block //  ethereum_mock.BlocksBetween(obscurocommon.GenesisBlock, head, node.Resolver)
}

// EthNodeInfo wraps eth node info in consistent manner
type EthNodeInfo struct {
	ID common.Address
}
