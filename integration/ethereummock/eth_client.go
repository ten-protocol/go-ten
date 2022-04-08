package ethereummock

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/go/ethclient"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
)

// MockEthClient mocks the client connections wrapper to ethereum
type MockEthClient struct {
	node *Node
}

func NewEthClient(node *Node) ethclient.Client {
	return &MockEthClient{
		node: node,
	}
}

func (m *MockEthClient) IsBlockAncestor(block types.Block, proof obscurocommon.L1RootHash) bool {
	return m.node.Resolver.IsBlockAncestor(&block, proof)
}

func (m *MockEthClient) FetchHeadBlock() (*types.Block, uint64) {
	return m.node.Resolver.FetchHeadBlock()
}

func (m *MockEthClient) Info() ethclient.Info {
	return ethclient.Info{ID: m.node.ID}
}

func (m *MockEthClient) BlocksBetween(block *types.Block, head *types.Block) []*types.Block {
	return m.node.BlocksBetween(block, head)
}

func (m *MockEthClient) FetchBlock(id common.Hash) (*types.Block, bool) {
	return m.node.Resolver.FetchBlock(id)
}

func (m *MockEthClient) IssueTx(tx obscurocommon.EncodedL1Tx) {
	m.node.BroadcastTx(tx)
}
