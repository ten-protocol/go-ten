package obsclient

import (
	"context"
	"testing"

	tenrpc "github.com/ten-protocol/go-ten/go/common/rpc"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ten-protocol/go-ten/go/responses"
)

// These tests use a mocked RPC client, they test any transformations of the Go objects -> RPC params, as well as any
//	transformations of the RPC resp -> return value

var (
	// since the RPC client is mocked out in these tests we are never testing the context (no async call to cancel)
	testCtx = context.Background()
	testAcc = common.BytesToAddress(common.Hex2Bytes("0000000000000000000000000000000000000abc"))
)

func TestNonceAt_ConvertsNilBlockNumberToLatest(t *testing.T) {
	mockRPC, authClient := createAuthClientWithMockRPCClient()

	// expect mock to be called once with the nonce request, it should have translated nil blockNumber to "latest" string
	mockRPC.On(
		"CallContext",
		testCtx, mock.AnythingOfType("*string"), tenrpc.ERPCGetTransactionCount, []interface{}{testAcc, "latest"},
	).Return(nil).Run(func(args mock.Arguments) {
		res := args.Get(1).(*responses.NonceType)
		// set the result pointer in the RPC client
		num := "0x2"
		*res = num
	})

	nonce, err := authClient.NonceAt(testCtx, nil)

	// assert mock called as expected
	mockRPC.AssertExpectations(t)
	// assert no error
	assert.Nil(t, err)
	// assert nonce returned correctly
	assert.Equal(t, uint64(2), nonce)
}

func createAuthClientWithMockRPCClient() (*rpcClientMock, *AuthObsClient) {
	mockRPC := new(rpcClientMock)
	authClient := &AuthObsClient{
		ObsClient: *NewObsClient(mockRPC),
		account:   testAcc,
	}
	return mockRPC, authClient
}
