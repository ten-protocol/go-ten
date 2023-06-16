package faucet

import (
	"fmt"
	"github.com/ethereum/go-ethereum/log"
	"github.com/obscuronet/go-obscuro/integration"
	"github.com/obscuronet/go-obscuro/integration/common/testlog"
	"github.com/obscuronet/go-obscuro/integration/ethereummock"
	"github.com/obscuronet/go-obscuro/integration/simulation/network"
	"github.com/obscuronet/go-obscuro/integration/simulation/params"
	"github.com/obscuronet/go-obscuro/tools/faucet/container"
	"github.com/obscuronet/go-obscuro/tools/faucet/faucet"
	"github.com/stretchr/testify/assert"
	"io"
	"math/big"
	"net/http"
	"strings"
	"testing"
	"time"
)

func init() { //nolint:gochecknoinits
	testlog.Setup(&testlog.Cfg{
		LogDir:      testLogs,
		TestType:    "faucet",
		TestSubtype: "test",
		LogLevel:    log.LvlInfo,
	})
}

const (
	contractDeployerPrivateKeyHex = "4bfe14725e685901c062ccd4e220c61cf9c189897b6c78bd18d7f51291b2b8f8"
	latestBlock                   = "latest"
	emptyCode                     = "0x"
	erc20ParamOne                 = "Hocus"
	erc20ParamTwo                 = "Hoc"
	erc20ParamThree               = "1000000000000000000"
	testLogs                      = "../.build/noderunner/"
	receiptTimeout                = 30 * time.Second // The time to wait for a receipt for a transaction.
	_portOffset                   = 1000
)

func TestFaucet(t *testing.T) {
	startPort := integration.StartPortFaucetUnitTest
	createObscuroNetwork(t, startPort)
	// This sleep is required to ensure the initial rollup exists, and thus contract deployer can check its balance.
	time.Sleep(2 * time.Second)

	faucetContainer, err := container.NewFaucetContainerFromConfig(&faucet.Config{
		Port:       startPort,
		Host:       "localhost",
		HTTPPort:   13000,
		PK:         "0x" + contractDeployerPrivateKeyHex,
		JWTSecret:  "This_is_secret",
		ChainID:    big.NewInt(integration.ObscuroChainID),
		ServerPort: 80,
	})
	assert.NoError(t, err)

	err = faucetContainer.Start()
	assert.NoError(t, err)

	err = issueRequest()
	assert.NoError(t, err)

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

func issueRequest() error {

	url := "http://localhost/fund/obx"
	method := "POST"

	payload := strings.NewReader(`{
    "address":"0x731ed18A8B84e83C79Da742052763272C4D802ee"
}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body))
	return nil
}
