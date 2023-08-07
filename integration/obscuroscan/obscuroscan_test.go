package faucet

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/obscuronet/go-obscuro/tools/obscuroscan_v2/backend/config"
	"github.com/obscuronet/go-obscuro/tools/obscuroscan_v2/backend/container"
	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"

	"github.com/ethereum/go-ethereum/log"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/integration"
	"github.com/obscuronet/go-obscuro/integration/common/testlog"
	"github.com/obscuronet/go-obscuro/integration/ethereummock"
	"github.com/obscuronet/go-obscuro/integration/simulation/network"
	"github.com/obscuronet/go-obscuro/integration/simulation/params"
	"github.com/stretchr/testify/assert"
)

func init() { //nolint:gochecknoinits
	testlog.Setup(&testlog.Cfg{
		LogDir:      testLogs,
		TestType:    "obscuroscan",
		TestSubtype: "test",
		LogLevel:    log.LvlInfo,
	})
}

const (
	testLogs = "../.build/obscuroscan/"
)

func TestObscuroscan(t *testing.T) {
	t.Skip("Commented it out until more testing is driven from this test")
	startPort := integration.StartPortObscuroscanUnitTest
	createObscuroNetwork(t, startPort)

	obsScanConfig := &config.Config{
		NodeHostAddress: fmt.Sprintf("http://127.0.0.1:%d", startPort+integration.DefaultHostRPCHTTPOffset),
		ServerAddress:   fmt.Sprintf("127.0.0.1:%d", startPort+integration.DefaultObscuroscanHTTPPortOffset),
		LogPath:         "sys_out",
	}
	serverAddress := fmt.Sprintf("http://%s", obsScanConfig.ServerAddress)

	obsScanContainer, err := container.NewObscuroScanContainer(obsScanConfig)
	require.NoError(t, err)

	err = obsScanContainer.Start()
	require.NoError(t, err)

	// wait for the msg bus contract to be deployed
	time.Sleep(10 * time.Second)

	// make sure the server is ready to receive requests
	err = waitServerIsReady(serverAddress)
	require.NoError(t, err)

	// Issue tests
	statusCode, body, err := fasthttp.Get(nil, fmt.Sprintf("%s/count/contracts/", serverAddress))
	assert.NoError(t, err)
	assert.Equal(t, 200, statusCode)
	assert.Equal(t, "{\"count\":1}", string(body))

	statusCode, body, err = fasthttp.Get(nil, fmt.Sprintf("%s/count/transactions/", serverAddress))
	assert.NoError(t, err)
	assert.Equal(t, 200, statusCode)
	assert.Equal(t, "{\"count\":1}", string(body))

	statusCode, body, err = fasthttp.Get(nil, fmt.Sprintf("%s/items/batch/latest/", serverAddress))
	assert.NoError(t, err)
	assert.Equal(t, 200, statusCode)

	type itemRes struct {
		Item common.BatchHeader `json:"item"`
	}

	itemObj := itemRes{}
	err = json.Unmarshal(body, &itemObj)
	assert.NoError(t, err)
	batchHead := itemObj.Item

	statusCode, _, err = fasthttp.Get(nil, fmt.Sprintf("%s/items/rollup/latest/", serverAddress))
	assert.NoError(t, err)
	assert.Equal(t, 200, statusCode)

	statusCode, _, err = fasthttp.Get(nil, fmt.Sprintf("%s/batch/%s", serverAddress, batchHead.Hash().String()))
	assert.NoError(t, err)
	assert.Equal(t, 200, statusCode)

	statusCode, body, err = fasthttp.Get(nil, fmt.Sprintf("%s/items/transactions/", serverAddress))
	assert.NoError(t, err)
	assert.Equal(t, 200, statusCode)

	type publicTxsRes struct {
		Result []common.PublicTxData `json:"result"`
	}

	publicTxsObj := publicTxsRes{}
	err = json.Unmarshal(body, &publicTxsObj)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(publicTxsObj.Result))

	// Gracefully shutdown
	err = obsScanContainer.Stop()
	assert.NoError(t, err)
}

func waitServerIsReady(serverAddr string) error {
	for now := time.Now(); time.Since(now) < 30*time.Second; time.Sleep(500 * time.Millisecond) {
		statusCode, _, err := fasthttp.Get(nil, fmt.Sprintf("%s/health/", serverAddr))
		if err != nil {
			// give it time to boot up
			if strings.Contains(err.Error(), "connection") {
				continue
			}
			return err
		}

		if statusCode == http.StatusOK {
			return nil
		}
	}
	return fmt.Errorf("timed out before server was ready")
}

// Creates a single-node Obscuro network for testing.
func createObscuroNetwork(t *testing.T, startPort int) {
	// Create the Obscuro network.
	numberOfNodes := 1
	wallets := params.NewSimWallets(1, numberOfNodes, integration.EthereumChainID, integration.ObscuroChainID)
	simParams := params.SimParams{
		NumberOfNodes:    numberOfNodes,
		AvgBlockDuration: 1 * time.Second,
		MgmtContractLib:  ethereummock.NewMgmtContractLibMock(),
		ERC20ContractLib: ethereummock.NewERC20ContractLibMock(),
		Wallets:          wallets,
		StartPort:        startPort,
	}

	obscuroNetwork := network.NewNetworkOfSocketNodes(wallets)
	t.Cleanup(obscuroNetwork.TearDown)
	_, err := obscuroNetwork.Create(&simParams, nil)
	if err != nil {
		panic(fmt.Sprintf("failed to create test Obscuro network. Cause: %s", err))
	}
}
