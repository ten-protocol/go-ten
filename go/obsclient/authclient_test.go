package obsclient

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/obscuronet/go-obscuro/go/rpcclientlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

var (
	// since the RPC client is mocked out in these tests we are never testing the context and we just use this
	testCtx = context.TODO()
	testAcc = common.BytesToAddress(common.Hex2Bytes("0000000000000000000000000000000000000abc"))
)

func TestNonceAt_ConvertsNilBlockNumberToLatest(t *testing.T) {
	mockRpc, authClient := createAuthClientWithMockRPCClient()

	// expect mock to be called once with the nonce request, it should have translated nil blockNumber to "latest" string
	mockRpc.On(
		"CallContext",
		testCtx, mock.AnythingOfType("*hexutil.Uint64"), rpcclientlib.RPCNonce, []interface{}{testAcc, "latest"},
	).Return(nil).Run(func(args mock.Arguments) {
		res := args.Get(1).(*hexutil.Uint64)
		// set the result pointer in the RPC client
		*res = 2
	})

	nonce, err := authClient.NonceAt(testCtx, testAcc, nil)

	// assert mock called as expected
	mockRpc.AssertExpectations(t)
	// assert no error
	assert.Nil(t, err)
	// assert nonce returned correctly
	assert.Equal(t, uint64(2), nonce)
}

func createAuthClientWithMockRPCClient() (*rpcClientMock, *AuthObsClient) {
	mockRpc := new(rpcClientMock)
	authClient := &AuthObsClient{
		ObsClient: ObsClient{c: mockRpc},
		c:         mockRpc,
	}
	return mockRpc, authClient
}
