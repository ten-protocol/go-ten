package simulation

import (
	"fmt"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/obscuroclient"

	"github.com/obscuronet/obscuro-playground/go/log"
)

const testLogs = "../.build/simulations/"

func setupTestLog() *os.File {
	// create a folder specific for the test
	err := os.MkdirAll(testLogs, 0o700)
	if err != nil {
		panic(err)
	}
	f, err := os.CreateTemp(testLogs, fmt.Sprintf("simulation-result-%d-*.txt", time.Now().Unix()))
	if err != nil {
		panic(err)
	}
	log.SetLog(f)
	return f
}

func minMax(arr []uint64) (min uint64, max uint64) {
	min = ^uint64(0)
	for _, no := range arr {
		if no < min {
			min = no
		}
		if no > max {
			max = no
		}
	}
	return
}

// Uses the client to retrieve the height of the current block head.
func getCurrentBlockHeadHeight(client *obscuroclient.Client) int64 {
	var l1Height int64
	err := (*client).Call(&l1Height, obscuroclient.RPCGetCurrentBlockHeadHeight)
	if err != nil {
		panic("Could not retrieve current block head.")
	}
	return l1Height
}

// Uses the client to retrieve the current rollup head.
func getCurrentRollupHead(client *obscuroclient.Client) *nodecommon.Header {
	var result *nodecommon.Header
	err := (*client).Call(&result, obscuroclient.RPCGetCurrentRollupHead)
	if err != nil {
		panic(fmt.Errorf("simulation failed due to failed RPC call. Cause: %w", err))
	}
	return result
}

// Uses the client to retrieve the rollup header with the matching hash.
func getRollupHeader(client *obscuroclient.Client, hash common.Hash) *nodecommon.Header {
	var result *nodecommon.Header
	err := (*client).Call(&result, obscuroclient.RPCGetRollupHeader, hash)
	if err != nil {
		panic(fmt.Errorf("simulation failed due to failed RPC call. Cause: %w", err))
	}
	return result
}
