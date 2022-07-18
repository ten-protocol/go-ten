package simulation

import (
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/obscuronet/go-obscuro/integration/common/testlog"

	testcommon "github.com/obscuronet/go-obscuro/integration/common"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/enclave/rollupchain"
	"github.com/obscuronet/go-obscuro/go/ethadapter/erc20contractlib"
	"github.com/obscuronet/go-obscuro/go/rpcclientlib"
)

const (
	testLogs             = "../.build/simulations/"
	jsonKeyRoot          = "root"
	jsonKeyStatus        = "status"
	receiptStatusFailure = "0x0"
)

func setupSimTestLog(simType string) {
	testlog.Setup(&testlog.Cfg{
		LogDir:      testLogs,
		TestType:    "sim-log",
		TestSubtype: simType,
	})
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
func getCurrentBlockHeadHeight(client rpcclientlib.Client) int64 {
	method := rpcclientlib.RPCGetCurrentBlockHead

	var blockHead *types.Header
	err := client.Call(&blockHead, method)
	if err != nil {
		panic(fmt.Errorf("simulation failed due to failed %s RPC call. Cause: %w", method, err))
	}

	if blockHead == nil || blockHead.Number == nil {
		panic(fmt.Errorf("simulation failed - no current block head found in RPC response from host"))
	}

	return blockHead.Number.Int64()
}

// Uses the client to retrieve the current rollup head.
func getCurrentRollupHead(client rpcclientlib.Client) *common.Header {
	method := rpcclientlib.RPCGetCurrentRollupHead

	var result *common.Header
	err := client.Call(&result, method)
	if err != nil {
		panic(fmt.Errorf("simulation failed due to failed %s RPC call. Cause: %w", method, err))
	}

	return result
}

// Uses the client to retrieve the rollup header with the matching hash.
func getRollupHeader(client rpcclientlib.Client, hash gethcommon.Hash) *common.Header {
	method := rpcclientlib.RPCGetRollupHeader

	var result *common.Header
	err := client.Call(&result, method, hash.Hex())
	if err != nil {
		panic(fmt.Errorf("simulation failed due to failed %s RPC call. Cause: %w", method, err))
	}

	return result
}

// Uses the client to retrieve the transaction with the matching hash.
func getTransaction(client rpcclientlib.Client, txHash gethcommon.Hash) map[string]interface{} {
	var encryptedResponse string
	err := client.Call(&encryptedResponse, rpcclientlib.RPCGetTransactionByHash, txHash.Hex())
	if err != nil {
		panic(fmt.Errorf("simulation failed due to failed %s RPC call. Cause: %w", rpcclientlib.RPCGetTransactionByHash, err))
	}

	var responseJSONMap map[string]interface{}
	err = json.Unmarshal(gethcommon.Hex2Bytes(encryptedResponse), &responseJSONMap)
	if err != nil {
		panic(fmt.Errorf("simulation failed because could not unmarshal JSON response to %s RPC call. Cause: %w", rpcclientlib.RPCGetTransactionByHash, err))
	}

	return responseJSONMap
}

// Returns the transaction receipt for the given transaction hash.
func getTransactionReceipt(client rpcclientlib.Client, txHash gethcommon.Hash) map[string]interface{} {
	var encryptedResponse string
	err := client.Call(&encryptedResponse, rpcclientlib.RPCGetTxReceipt, txHash.Hex())
	if err != nil {
		panic(fmt.Errorf("simulation failed due to failed %s RPC call. Cause: %w", rpcclientlib.RPCGetTxReceipt, err))
	}

	var responseJSONMap map[string]interface{}
	err = json.Unmarshal(gethcommon.Hex2Bytes(encryptedResponse), &responseJSONMap)
	if err != nil {
		panic(fmt.Errorf("simulation failed because could not unmarshall JSON response to %s RPC call. Cause: %w", rpcclientlib.RPCGetTxReceipt, err))
	}

	return responseJSONMap
}

// Uses the client to retrieve the balance of the wallet with the given address.
func balance(client rpcclientlib.Client, address gethcommon.Address, l2ContractAddress *gethcommon.Address) uint64 {
	method := rpcclientlib.RPCCall
	balanceData := erc20contractlib.CreateBalanceOfData(address)
	convertedData := (hexutil.Bytes)(balanceData)

	params := map[string]interface{}{
		rollupchain.CallFieldFrom: address.Hex(),
		rollupchain.CallFieldTo:   l2ContractAddress.Hex(),
		rollupchain.CallFieldData: convertedData,
	}

	var encryptedResponse string
	err := client.Call(&encryptedResponse, method, params)
	if err != nil {
		panic(fmt.Errorf("simulation failed due to failed %s RPC call. Cause: %w", method, err))
	}
	bytes, err := hexutil.Decode(string(gethcommon.Hex2Bytes(encryptedResponse)))
	if err != nil {
		panic(fmt.Errorf("could not decode ERC20 balance response for node. Response was %s. Cause: %w", encryptedResponse, err))
	}
	r := new(big.Int)
	r = r.SetBytes(bytes)
	return r.Uint64()
}

// FindHashDups - returns a map of all hashes that appear multiple times, and how many times
func findHashDups(list []gethcommon.Hash) map[gethcommon.Hash]int {
	elementCount := make(map[gethcommon.Hash]int)

	for _, item := range list {
		// check if the item/element exist in the duplicate_frequency map
		_, exist := elementCount[item]
		if exist {
			elementCount[item]++ // increase counter by 1 if already in the map
		} else {
			elementCount[item] = 1 // else start counting from 1
		}
	}
	dups := make(map[gethcommon.Hash]int)
	for u, i := range elementCount {
		if i > 1 {
			dups[u] = i
			fmt.Printf("Dup: %s\n", u)
		}
	}
	return dups
}

// FindRollupDups - returns a map of all L2 root hashes that appear multiple times, and how many times
func findRollupDups(list []common.L2RootHash) map[common.L2RootHash]int {
	elementCount := make(map[common.L2RootHash]int)

	for _, item := range list {
		// check if the item/element exist in the duplicate_frequency map
		_, exist := elementCount[item]
		if exist {
			elementCount[item]++ // increase counter by 1 if already in the map
		} else {
			elementCount[item] = 1 // else start counting from 1
		}
	}
	dups := make(map[common.L2RootHash]int)
	for u, i := range elementCount {
		if i > 1 {
			dups[u] = i
			fmt.Printf("Dup: r_%d\n", common.ShortHash(u))
		}
	}
	return dups
}

func SleepRndBtw(min time.Duration, max time.Duration) {
	time.Sleep(testcommon.RndBtwTime(min, max))
}
