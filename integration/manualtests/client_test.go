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
	hostRPCAddress := l2Host
	client, err := rpc.NewNetworkClient(hostRPCAddress)
	assert.Nil(t, err)

	obsClient := obsclient.NewObsClient(client)

	rollupHeader, err := obsClient.RollupHeaderByNumber(big.NewInt(4392))
	assert.Nil(t, err)

	var rollup *common.ExtRollup
	err = client.Call(&rollup, rpc.GetBatch, rollupHeader.Hash())
	assert.Nil(t, err)
}
