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

func setupTestLog(simType string) *os.File {
	// create a folder specific for the test
	err := os.MkdirAll(testLogs, 0o700)
	if err != nil {
		panic(err)
	}
	timeFormatted := time.Now().Format("2006-01-02-15-04-05")
	f, err := os.CreateTemp(testLogs, fmt.Sprintf("simulation-result-%s-%s-*.txt", timeFormatted, simType))
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
	method := obscuroclient.RPCGetCurrentBlockHeadHeight

	var l1Height int64
	err := (*client).Call(&l1Height, method)
	if err != nil {
		panic(fmt.Errorf("simulation failed due to failed %s RPC call. Cause: %w", method, err))
	}

	return l1Height
}

// Uses the client to retrieve the current rollup head.
func getCurrentRollupHead(client *obscuroclient.Client) *nodecommon.Header {
	method := obscuroclient.RPCGetCurrentRollupHead

	var result *nodecommon.Header
	err := (*client).Call(&result, method)
	if err != nil {
		panic(fmt.Errorf("simulation failed due to failed %s RPC call. Cause: %w", method, err))
	}

	return result
}

// Uses the client to retrieve the rollup header with the matching hash.
func getRollupHeader(client *obscuroclient.Client, hash common.Hash) *nodecommon.Header {
	method := obscuroclient.RPCGetRollupHeader

	var result *nodecommon.Header
	err := (*client).Call(&result, method, hash)
	if err != nil {
		panic(fmt.Errorf("simulation failed due to failed %s RPC call. Cause: %w", method, err))
	}

	return result
}

// Uses the client to retrieve the transaction with the matching hash.
func getTransaction(client *obscuroclient.Client, hash common.Hash) *nodecommon.L2Tx {
	method := obscuroclient.RPCGetTransaction

	var result *nodecommon.L2Tx
	err := (*client).Call(&result, method, hash)
	if err != nil {
		panic(fmt.Errorf("simulation failed due to failed %s RPC call. Cause: %w", method, err))
	}

	return result
}

// Uses the client to retrieve the balance of the wallet with the given address.
func balance(client *obscuroclient.Client, address common.Address) uint64 {
	method := obscuroclient.RPCBalance

	var result uint64
	err := (*client).Call(&result, method, address)
	if err != nil {
		panic(fmt.Errorf("simulation failed due to failed %s RPC call. Cause: %w", method, err))
	}

	return result
}
