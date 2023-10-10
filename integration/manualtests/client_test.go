package manualtests

import (
	"math/big"
	"os"
	"testing"

	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/obsclient"
	"github.com/obscuronet/go-obscuro/go/rpc"
	"github.com/stretchr/testify/assert"
)

func TestClientGetRollup(t *testing.T) {
	// test is skipped by default to avoid breaking CI - set env flag in `Run Configurations` to run it in IDE
	if os.Getenv(_IDEFlag) == "" {
		t.Skipf("set flag %s to run this test in the IDE", _IDEFlag)
	}
	hostRPCAddress := "http://erpc.sepolia-testnet.obscu.ro:80"
	client, err := rpc.NewNetworkClient(hostRPCAddress)
	assert.Nil(t, err)

	obsClient := obsclient.NewObsClient(client)

	rollupHeader, err := obsClient.BatchHeaderByNumber(big.NewInt(4392))
	assert.Nil(t, err)

	var rollup *common.ExtRollup
	err = client.Call(&rollup, rpc.GetBatch, rollupHeader.Hash())
	assert.Nil(t, err)
}

func TestClientGetTransactionCount(t *testing.T) {
	// test is skipped by default to avoid breaking CI - set env flag in `Run Configurations` to run it in IDE
	if os.Getenv(_IDEFlag) == "" {
		t.Skipf("set flag %s to run this test in the IDE", _IDEFlag)
	}
	hostRPCAddress := "http://127.0.0.1:37801"
	client, err := rpc.NewNetworkClient(hostRPCAddress)
	assert.Nil(t, err)

	obsClient := obsclient.NewObsClient(client)

	count, err := obsClient.GetTotalContractCount()
	assert.Nil(t, err)
	assert.NotEqual(t, count, 0)

	var contractCount int
	err = client.Call(&contractCount, rpc.GetTotalContractCount)
	assert.Nil(t, err)
}
