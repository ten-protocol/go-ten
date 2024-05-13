package manualtests

import (
	"math/big"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/obsclient"
	"github.com/ten-protocol/go-ten/go/rpc"
)

func TestClientGetRollup(t *testing.T) {
	// test is skipped by default to avoid breaking CI - set env flag in `Run Configurations` to run it in IDE
	if os.Getenv(_IDEFlag) == "" {
		t.Skipf("set flag %s to run this test in the IDE", _IDEFlag)
	}
	hostRPCAddress := "http://erpc.sepolia-testnet.ten.xyz:80"
	client, err := rpc.NewNetworkClient(hostRPCAddress)
	assert.Nil(t, err)

	obsClient := obsclient.NewObsClient(client)

	batchHeader, err := obsClient.GetBatchHeaderByNumber(big.NewInt(4392))
	assert.Nil(t, err)

	var rollup *common.ExtRollup
	err = client.Call(&rollup, rpc.GetBatch, batchHeader.Hash())
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
