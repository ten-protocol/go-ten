package ethclient

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	ethereum_mock "github.com/obscuronet/obscuro-playground/integration/ethereummock"
)

// MockClient implements the eth Client interface for mock clients
type MockClient struct {
	node *ethereum_mock.Node
}

func NewMockClient(node *ethereum_mock.Node) Client {
	return &MockClient{node: node}
}

func (m *MockClient) FetchBlock(hash common.Hash) (*types.Block, bool) {
	return m.node.Resolver.FetchBlock(hash)
}

func (m *MockClient) BroadcastTx(t obscurocommon.EncodedL1Tx) {
	m.node.BroadcastTx(t)
}
