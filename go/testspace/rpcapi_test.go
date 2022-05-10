package testspace

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/rpc"
	"math/big"
	"testing"
)

func TestRPCAPIs(t *testing.T) {
	p2pConfig := node.Config{
		// Other config options are available for websockets and IPC.
		HTTPHost: "localhost",
		HTTPPort: 3000,
	}
	p2pNode, err := node.New(&p2pConfig)
	if err != nil {
		panic(err)
	}

	// We define our own APIs for specific namespaces. See the `EthAPI` type below.
	rpcAPIs := []rpc.API{
		{
			Namespace: "eth",
			Version:   "1.0",
			Service:   NewEthAPI(),
			Public:    true,
		},
	}
	p2pNode.RegisterAPIs(rpcAPIs)

	err = p2pNode.Start()
	if err != nil {
		panic(err)
	}

	client, err := ethclient.Dial("http://127.0.0.1:3000")
	if err != nil {
		panic(err)
	}

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		panic(err)
	}

	balance, err := client.BalanceAt(context.Background(), common.HexToAddress("0x"), big.NewInt(1))
	if err != nil {
		panic(err)
	}

	if chainID.Uint64() != 777 {
		t.Fatal("Did not retrieve correct chain ID.")
	}
	if balance.Uint64() != 888 {
		t.Fatal("Did not retrieve correct balance.")
	}
}

type EthAPI struct {
}

func NewEthAPI() *EthAPI {
	return &EthAPI{}
}

func (api *EthAPI) GetBalance(context.Context, common.Address, rpc.BlockNumberOrHash) (*hexutil.Big, error) {
	return (*hexutil.Big)(big.NewInt(888)), nil
}

func (api *EthAPI) ChainId() (*hexutil.Big, error) {
	return (*hexutil.Big)(big.NewInt(777)), nil
}
