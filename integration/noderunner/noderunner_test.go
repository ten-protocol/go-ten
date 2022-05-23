package noderunner

import (
	"encoding/hex"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/obscuronet/obscuro-playground/integration"

	"github.com/obscuronet/obscuro-playground/go/log"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/enclaverunner"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/host/hostrunner"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/obscuroclient"
	"github.com/obscuronet/obscuro-playground/integration/gethnetwork"
)

// TODO - Use the HostRunner/EnclaveRunner methods in the socket-based integration tests, and retire this smoketest.

const testLogs = "../.build/noderunner/"

// A smoke test to check that we can stand up a standalone Obscuro host and enclave.
func TestCanStartStandaloneObscuroHostAndEnclave(t *testing.T) {
	setupTestLog()

	startPort := integration.StartPortNodeRunnerTest
	enclaveAddr := fmt.Sprintf("127.0.0.1:%d", startPort)
	clientServerAddr := fmt.Sprintf("127.0.0.1:%d", startPort+1)
	ethClientPort := startPort + 2

	privateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	address := crypto.PubkeyToAddress(privateKey.PublicKey)

	hostConfig := hostrunner.DefaultHostConfig()
	hostConfig.PrivateKeyString = hex.EncodeToString(crypto.FromECDSA(privateKey))
	hostConfig.EnclaveAddr = enclaveAddr
	hostConfig.ClientServerAddr = clientServerAddr
	hostConfig.EthClientPort = ethClientPort

	enclaveConfig := enclaverunner.DefaultEnclaveConfig()
	enclaveConfig.Address = enclaveAddr

	gethBinaryPath, err := gethnetwork.EnsureBinariesExist(gethnetwork.LatestVersion)
	if err != nil {
		t.Fatal(err)
	}
	network := gethnetwork.NewGethNetwork(int(ethClientPort)-100, gethBinaryPath, 1, 1, []string{address.String()})
	defer network.StopNodes()
	go enclaverunner.RunEnclave(enclaveConfig)
	go hostrunner.RunHost(hostConfig)

	// We sleep to give the network time to produce some blocks.
	time.Sleep(3 * time.Second)

	obscuroClient := obscuroclient.NewClient(common.BytesToAddress([]byte(hostConfig.NodeID)), clientServerAddr)
	var result types.Header
	err = obscuroClient.Call(&result, obscuroclient.RPCGetCurrentBlockHead)
	if err != nil {
		t.Fatal(err)
	}

	if result.Number.Uint64() == 0 {
		t.Fatal("Zero blocks have been produced. Something is wrong.")
	}
}

func setupTestLog() *os.File {
	// create a folder specific for the test
	err := os.MkdirAll(testLogs, 0o700)
	if err != nil {
		panic(err)
	}
	timeFormatted := time.Now().Format("2006-01-02_15-04-05")
	f, err := os.CreateTemp(testLogs, fmt.Sprintf("noderunner-%s-*.txt", timeFormatted))
	if err != nil {
		panic(err)
	}
	log.OutputToFile(f)
	return f
}
