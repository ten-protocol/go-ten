package noderunner

import (
	"encoding/hex"
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/obscuronet/go-obscuro/go/common/profiler"
	"github.com/obscuronet/go-obscuro/integration/common/testlog"

	"github.com/obscuronet/go-obscuro/go/config"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/obscuronet/go-obscuro/integration"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/obscuronet/go-obscuro/go/enclave/enclaverunner"
	"github.com/obscuronet/go-obscuro/go/host/hostrunner"
	"github.com/obscuronet/go-obscuro/go/rpc"
	"github.com/obscuronet/go-obscuro/integration/gethnetwork"
)

// TODO - Use the HostRunner/EnclaveRunner methods in the socket-based integration tests, and retire this smoketest.

const (
	testLogs             = "../.build/noderunner/"
	gethPort             = integration.StartPortNodeRunnerTest + 2
	defaultWsPortOffset  = 100
	gethWebsocketPort    = gethPort + defaultWsPortOffset
	localhost            = "127.0.0.1"
	obscuroWebsocketPort = integration.StartPortNodeRunnerTest + 1
)

// A smoke test to check that we can stand up a standalone Obscuro host and enclave.
func TestCanStartStandaloneObscuroHostAndEnclave(t *testing.T) {
	testlog.Setup(&testlog.Cfg{
		LogDir:      testLogs,
		TestType:    "noderunner",
		TestSubtype: "test",
	})

	enclaveAddr := fmt.Sprintf("%s:%d", localhost, integration.StartPortNodeRunnerTest)
	rpcAddress := fmt.Sprintf("ws://%s:%d", localhost, obscuroWebsocketPort)

	privateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	address := crypto.PubkeyToAddress(privateKey.PublicKey)

	hostConfig := config.DefaultHostConfig()
	hostConfig.PrivateKeyString = hex.EncodeToString(crypto.FromECDSA(privateKey))
	hostConfig.EnclaveRPCAddress = enclaveAddr
	hostConfig.ClientRPCPortWS = uint64(obscuroWebsocketPort)
	hostConfig.L1NodeWebsocketPort = uint(gethWebsocketPort)
	hostConfig.ProfilerEnabled = true

	enclaveConfig := config.DefaultEnclaveConfig()
	enclaveConfig.Address = enclaveAddr
	dummyContractAddress := common.BytesToAddress([]byte("AA"))
	enclaveConfig.ERC20ContractAddresses = []*common.Address{&dummyContractAddress, &dummyContractAddress}
	enclaveConfig.ProfilerEnabled = true

	gethBinaryPath, err := gethnetwork.EnsureBinariesExist(gethnetwork.LatestVersion)
	if err != nil {
		t.Fatal(err)
	}
	network := gethnetwork.NewGethNetwork(gethPort, gethWebsocketPort, gethBinaryPath, 1, 1, []string{address.String()})
	defer network.StopNodes()

	go enclaverunner.RunEnclave(enclaveConfig)
	go hostrunner.RunHost(hostConfig)

	// we create the node RPC client
	var obscuroClient rpc.Client
	wait := 30 // max wait in seconds
	for {
		obscuroClient, err = rpc.NewNetworkClient(rpcAddress)
		if err == nil {
			break
		}
		if wait <= 0 {
			t.Fatal("RPC client server never became available")
		}
		time.Sleep(time.Second)
		wait--
	}
	defer teardown(obscuroClient, rpcAddress)

	// Check if the host profiler is up
	_, err = http.Get(fmt.Sprintf("http://%s:%d", localhost, profiler.DefaultHostPort)) //nolint
	if err != nil {
		t.Errorf("host profiler is not reachable: %s", err)
	}
	// Check if the enclave profiler is up
	_, err = http.Get(fmt.Sprintf("http://%s:%d", localhost, profiler.DefaultEnclavePort)) //nolint
	if err != nil {
		t.Errorf("enclave profiler is not reachable: %s", err)
	}

	counter := 0
	// We retry 20 times to check if the network has produced any blocks, sleeping half a second between each attempt.
	for counter < 20 {
		counter++
		time.Sleep(500 * time.Millisecond)

		var result types.Header
		err = obscuroClient.Call(&result, rpc.RPCGetCurrentBlockHead)
		if err == nil && result.Number.Uint64() > 0 {
			return
		}
	}

	t.Fatal("Zero blocks have been produced after ten seconds. Something is wrong.")
}

func teardown(obscuroClient rpc.Client, rpcServerAddr string) {
	if obscuroClient == nil {
		return
	}

	obscuroClient.Call(nil, rpc.RPCStopHost) //nolint:errcheck

	// We wait for the client server port to be closed.
	wait := 0
	for tcpConnectionAvailable(rpcServerAddr) {
		if wait == 20 { // max wait in seconds
			panic(fmt.Sprintf("RPC client server had not shut down after %d seconds", wait))
		}
		time.Sleep(time.Second)
		wait++
	}
}

func tcpConnectionAvailable(addr string) bool {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return false
	}
	_ = conn.Close()
	// we don't worry about failure while closing, it connected successfully so let test proceed
	return true
}
