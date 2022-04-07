package exec

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/go/ethclient"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	ethereum_mock "github.com/obscuronet/obscuro-playground/integration/ethereummock"
)

// MockEthNode implements the eth node interface for simulated eth nodes
type MockEthNode struct {
	node   *ethereum_mock.Node // Accesses the node lifecycle
	client ethclient.Client    // Accesses the nodes RPC endpoints
}

func NewMockEthNode(node *ethereum_mock.Node) EthNode {
	return &MockEthNode{node: node, client: ethclient.NewMockClient(node)}
}

func (m *MockEthNode) FetchHeadBlock() (*types.Block, uint64) {
	return m.node.Resolver.FetchHeadBlock()
}

func (m *MockEthNode) Info() EthNodeInfo {
	return EthNodeInfo{ID: m.node.ID}
}

func (m *MockEthNode) BlocksBetween(block *types.Block, head *types.Block) []*types.Block {
	return m.node.BlocksBetween(block, head)
}

func (m *MockEthNode) IsBlockAncestor(block *types.Block, proof obscurocommon.L1RootHash) bool {
	return m.node.Resolver.IsBlockAncestor(block, proof)
}

func (m *MockEthNode) KnownPeers(nodes []EthNode) {
	var enodes []*ethereum_mock.Node
	for _, n := range nodes {
		enodes = append(enodes, n.(*MockEthNode).node)
	}
	m.node.Network.(*ethereum_mock.MockEthNetwork).AllNodes = enodes
}

func (m *MockEthNode) Start() {
	m.node.Start()
}

func (m *MockEthNode) Stop() {
	m.node.Stop()
}

func (m *MockEthNode) Client() ethclient.Client {
	return m.client
}
