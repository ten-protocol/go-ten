package obsclient

import (
	"context"
	"github.com/stretchr/testify/mock"
)

// this mock lets us simulate responses from the RPC client, so we can verify the obsclient usage and response handling
type rpcClientMock struct {
	mock.Mock
}

func (m *rpcClientMock) Call(result interface{}, method string, args ...interface{}) error {
	arguments := m.Called(result, method, args)
	return arguments.Error(0)
}
func (m *rpcClientMock) CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	arguments := m.Called(ctx, result, method, args)
	return arguments.Error(0)
}
func (m *rpcClientMock) Stop() {
	m.Called()
}
