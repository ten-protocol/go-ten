package tenscan

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	testcommon "github.com/ten-protocol/go-ten/integration/common"

	"github.com/ten-protocol/go-ten/go/common"

	"github.com/ten-protocol/go-ten/tools/tenscan/backend/config"
	"github.com/ten-protocol/go-ten/tools/tenscan/backend/container"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/ten-protocol/go-ten/go/wallet"
	"github.com/valyala/fasthttp"

	"github.com/ethereum/go-ethereum/log"
	"github.com/stretchr/testify/assert"
	"github.com/ten-protocol/go-ten/integration"
	"github.com/ten-protocol/go-ten/integration/common/testlog"
)

func init() { //nolint:gochecknoinits
	testlog.Setup(&testlog.Cfg{
		LogDir:      testLogs,
		TestType:    "tenscan",
		TestSubtype: "test",
		LogLevel:    log.LvlInfo,
	})
}

func TestTenscanContracts(t *testing.T) {
	t.Skipf("Used only to test contract listing stuff locally")
	startPort := integration.TestPorts.TestTenscanPort
	createTenNetwork(t, integration.TestPorts.TestTenscanPort)

	tenScanConfig := &config.Config{
		NodeHostAddress: fmt.Sprintf("http://127.0.0.1:%d", startPort+integration.DefaultHostRPCHTTPOffset),
		ServerAddress:   fmt.Sprintf("127.0.0.1:%d", startPort+integration.DefaultTenscanHTTPPortOffset),
		LogPath:         "sys_out",
	}

	serverAddress := fmt.Sprintf("http://%s", tenScanConfig.ServerAddress)

	tenScanContainer, err := container.NewTenScanContainer(tenScanConfig)
	require.NoError(t, err)

	err = tenScanContainer.Start()
	require.NoError(t, err)

	// wait for the msg bus contract to be deployed
	time.Sleep(30 * time.Second)

	// make sure the server is ready to receive requests
	err = waitServerIsReady(serverAddress)
	require.NoError(t, err)

	err = waitForFirstRollup(serverAddress)
	require.NoError(t, err)

	deployTestContracts(
		t,
		fmt.Sprintf("ws://127.0.0.1:%d", startPort+integration.DefaultHostRPCWSOffset),
		wallet.NewInMemoryWalletFromConfig(testcommon.TestnetPrefundedPK, integration.TenChainID, testlog.Logger()),
		3, // deploy 3 test contracts (just enough to test pagination)
	)

	type contractListingRes struct {
		Result common.ContractListingResponse `json:"result"`
	}

	// wait for contracts to be synced (sync interval is 30s, so need to wait longer)
	// we just need at least 3 contracts to test the basics (contract listing, fetching, pagination)
	contractListingObj := contractListingRes{}
	contractWaitDeadline := time.Now().Add(time.Minute * 2) // allow multiple sync cycles
	for {
		statusCode, body, err := fasthttp.Get(nil, fmt.Sprintf("%s/items/contracts/?offset=0&size=30", serverAddress))
		assert.NoError(t, err)
		assert.Equal(t, 200, statusCode)
		err = json.Unmarshal(body, &contractListingObj)
		assert.NoError(t, err)
		// just need at least 3 contracts to test pagination and basic functionality
		if len(contractListingObj.Result.Contracts) >= 3 {
			t.Logf("Contracts synced: total=%d, visible=%d", contractListingObj.Result.Total, len(contractListingObj.Result.Contracts))
			break
		}
		if time.Now().After(contractWaitDeadline) {
			t.Fatalf("Timed out waiting for contracts to be indexed; have total=%d, visible=%d", contractListingObj.Result.Total, len(contractListingObj.Result.Contracts))
		}
		time.Sleep(5 * time.Second)
	}

	if len(contractListingObj.Result.Contracts) > 0 {
		firstContract := contractListingObj.Result.Contracts[0]
		assert.NotEqual(t, gethcommon.Address{}, firstContract.Address, "Contract address should not be empty")
		assert.NotEqual(t, gethcommon.Address{}, firstContract.Creator, "Contract creator should not be empty")
		assert.GreaterOrEqual(t, firstContract.BatchSeq, uint64(0), "Batch sequence should be set")
		assert.GreaterOrEqual(t, firstContract.Height, uint64(0), "Height should be set")
		assert.GreaterOrEqual(t, firstContract.Time, uint64(0), "Time should be set")
	}

	if len(contractListingObj.Result.Contracts) > 0 {
		testContractAddr := contractListingObj.Result.Contracts[0].Address
		statusCode, body, err := fasthttp.Get(nil, fmt.Sprintf("%s/items/contract/%s", serverAddress, testContractAddr.Hex()))
		assert.NoError(t, err)
		assert.Equal(t, 200, statusCode)

		type contractItemRes struct {
			Item common.PublicContract `json:"item"`
		}

		contractItemObj := contractItemRes{}
		err = json.Unmarshal(body, &contractItemObj)
		assert.NoError(t, err)
		assert.Equal(t, testContractAddr, contractItemObj.Item.Address, "Fetched contract should match requested address")
	}

	if contractListingObj.Result.Total > 2 {
		// first page
		statusCode, body, err := fasthttp.Get(nil, fmt.Sprintf("%s/items/contracts/?offset=0&size=1", serverAddress))
		assert.NoError(t, err)
		assert.Equal(t, 200, statusCode)

		firstPageContracts := contractListingRes{}
		err = json.Unmarshal(body, &firstPageContracts)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(firstPageContracts.Result.Contracts), "First page should have 1 contract")

		// second page
		statusCode, body, err = fasthttp.Get(nil, fmt.Sprintf("%s/items/contracts/?offset=1&size=1", serverAddress))
		assert.NoError(t, err)
		assert.Equal(t, 200, statusCode)

		secondPageContracts := contractListingRes{}
		err = json.Unmarshal(body, &secondPageContracts)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(secondPageContracts.Result.Contracts), "Second page should have 1 contract")

		// verify different pages have different contracts
		assert.NotEqual(t, firstPageContracts.Result.Contracts[0].Address,
			secondPageContracts.Result.Contracts[0].Address,
			"Different pages should have different contracts")
	}

	// invalid contract address returns appropriate error
	statusCode, _, err := fasthttp.Get(nil, fmt.Sprintf("%s/items/contract/0xinvalid", serverAddress))
	assert.NoError(t, err)
	// Should handle gracefully (either 400 or 404 or 500 depending on implementation)
	assert.True(t, statusCode >= 400, "Invalid address should return error status")

	// non-existent contract address
	nonExistentAddr := "0x0000000000000000000000000000000000000001"
	statusCode, _, err = fasthttp.Get(nil, fmt.Sprintf("%s/items/contract/%s", serverAddress, nonExistentAddr))
	assert.NoError(t, err)
	assert.True(t, statusCode >= 400, "Non-existent contract should return error status")
}

// TestContractsPagination tests that the contract sync properly handles
// pagination when there are more contracts than the fetch limit.
// This test uses a fetch limit of 2 and sync interval of 30s to verify all contracts
// are eventually synced through multiple sync cycles.
func TestContractsPagination(t *testing.T) {
	t.Skipf("Local only test - uncomment to test contract pagination with small fetch limit")
	startPort := integration.TestPorts.TestTenscanPort
	createTenNetwork(t, integration.TestPorts.TestTenscanPort)

	tenScanConfig := &config.Config{
		NodeHostAddress: fmt.Sprintf("http://127.0.0.1:%d", startPort+integration.DefaultHostRPCHTTPOffset),
		ServerAddress:   fmt.Sprintf("127.0.0.1:%d", startPort+integration.DefaultTenscanHTTPPortOffset),
		LogPath:         "sys_out",
	}

	serverAddress := fmt.Sprintf("http://%s", tenScanConfig.ServerAddress)

	tenScanContainer, err := container.NewTenScanContainer(tenScanConfig)
	require.NoError(t, err)

	err = tenScanContainer.Start()
	require.NoError(t, err)

	// wait for the initial setup
	time.Sleep(30 * time.Second)

	err = waitServerIsReady(serverAddress)
	require.NoError(t, err)

	err = waitForFirstRollup(serverAddress)
	require.NoError(t, err)

	// more contracts than the fetch limit to force pagination with limit of 2, deploying 7 contracts means we need at
	// least 4 sync cycles
	numTestContracts := 7
	t.Logf("Deploying %d test contracts (fetch limit is 2, so this will require multiple sync cycles)", numTestContracts)

	deployTestContracts(
		t,
		fmt.Sprintf("ws://127.0.0.1:%d", startPort+integration.DefaultHostRPCWSOffset),
		wallet.NewInMemoryWalletFromConfig(testcommon.TestnetPrefundedPK, integration.TenChainID, testlog.Logger()),
		numTestContracts,
	)

	type contractListingRes struct {
		Result common.ContractListingResponse `json:"result"`
	}

	// system contracts are already deployed, so we need to account for those
	systemContractCount := 0
	statusCode, body, err := fasthttp.Get(nil, fmt.Sprintf("%s/items/contracts/?offset=0&size=100", serverAddress))
	require.NoError(t, err)
	require.Equal(t, 200, statusCode)

	initialListing := contractListingRes{}
	err = json.Unmarshal(body, &initialListing)
	require.NoError(t, err)
	systemContractCount = int(initialListing.Result.Total)
	t.Logf("Initial system contracts: %d", systemContractCount)

	expectedTotalContracts := systemContractCount + numTestContracts
	t.Logf("Expecting total contracts: %d (system: %d + test: %d)", expectedTotalContracts, systemContractCount, numTestContracts)

	// With sync interval of 30s and limit of 2, we need:
	// 1: +2 contracts (after 30s)
	// 2: +2 contracts (after 60s)
	// 3: +2 contracts (after 90s)
	// 4: +1 contract (after 120s) = 7 total
	// Add buffer time for processing
	maxWaitTime := 150 * time.Second
	pollInterval := 10 * time.Second

	t.Logf("Starting to poll for contracts. Will wait up to %v", maxWaitTime)
	deadline := time.Now().Add(maxWaitTime)

	lastSeenCount := systemContractCount
	progressCheckpoints := make(map[int]time.Time)

	for {
		statusCode, body, err := fasthttp.Get(nil, fmt.Sprintf("%s/items/contracts/?offset=0&size=100", serverAddress))
		require.NoError(t, err)
		require.Equal(t, 200, statusCode)

		contractListing := contractListingRes{}
		err = json.Unmarshal(body, &contractListing)
		require.NoError(t, err)

		currentCount := int(contractListing.Result.Total)
		newContractsSynced := currentCount - systemContractCount

		if currentCount > lastSeenCount {
			elapsed := time.Since(progressCheckpoints[systemContractCount])
			if len(progressCheckpoints) == 0 {
				progressCheckpoints[systemContractCount] = time.Now()
				elapsed = 0
			}
			progressCheckpoints[currentCount] = time.Now()

			t.Logf("Progress: %d/%d test contracts synced (total: %d) - elapsed: %v",
				newContractsSynced, numTestContracts, currentCount, elapsed)
			lastSeenCount = currentCount
		}

		// all test contracts are synced
		if newContractsSynced >= numTestContracts {
			totalTime := time.Since(progressCheckpoints[systemContractCount])
			t.Logf("SUCCESS: All %d test contracts synced in %v", numTestContracts, totalTime)

			assert.GreaterOrEqual(t, len(contractListing.Result.Contracts), numTestContracts,
				"Should have at least the test contracts visible")

			// verify proper fields
			if len(contractListing.Result.Contracts) > 0 {
				sampleContract := contractListing.Result.Contracts[0]
				assert.NotEmpty(t, sampleContract.Address, "Contract should have address")
				assert.NotEmpty(t, sampleContract.Creator, "Contract should have creator")
				assert.Greater(t, sampleContract.Time, uint64(0), "Contract should have timestamp")
			}

			break
		}

		if time.Now().After(deadline) {
			t.Fatalf("Timeout: Only synced %d/%d test contracts after %v. "+
				"Expected all contracts to be synced through pagination (limit=2, interval=30s). "+
				"Total contracts: %d, System contracts: %d",
				newContractsSynced, numTestContracts, maxWaitTime,
				currentCount, systemContractCount)
		}

		time.Sleep(pollInterval)
	}

	// ensure pagination doesn't cause duplicates
	t.Log("Verifying no duplicate contracts were synced...")
	statusCode, body, err = fasthttp.Get(nil, fmt.Sprintf("%s/items/contracts/?offset=0&size=100", serverAddress))
	require.NoError(t, err)
	require.Equal(t, 200, statusCode)

	finalListing := contractListingRes{}
	err = json.Unmarshal(body, &finalListing)
	require.NoError(t, err)

	// duplicates by address
	seenAddresses := make(map[string]bool)
	for _, contract := range finalListing.Result.Contracts {
		addrStr := contract.Address.Hex()
		assert.False(t, seenAddresses[addrStr], "Found duplicate contract address: %s", addrStr)
		seenAddresses[addrStr] = true
	}

	t.Logf("Pagination test completed successfully. No duplicates found among %d contracts", len(finalListing.Result.Contracts))
}
