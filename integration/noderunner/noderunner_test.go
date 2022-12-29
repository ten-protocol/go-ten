package noderunner

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"testing"
	"time"

	enclavecontainer "github.com/obscuronet/go-obscuro/go/enclave/container"
	hostcontainer "github.com/obscuronet/go-obscuro/go/host/container"

	"github.com/obscuronet/go-obscuro/go/common/profiler"

	"github.com/ethereum/go-ethereum/common/hexutil"

	gethlog "github.com/ethereum/go-ethereum/log"

	"github.com/obscuronet/go-obscuro/integration/common/testlog"

	"github.com/obscuronet/go-obscuro/go/config"

	"github.com/ethereum/go-ethereum/common"

	"github.com/obscuronet/go-obscuro/integration"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/obscuronet/go-obscuro/go/rpc"
	"github.com/obscuronet/go-obscuro/integration/gethnetwork"
)

// TODO - Use the NewHostContainerFromConfig/NewEnclaveContainerFromConfig methods in the socket-based integration tests, and retire this smoketest.

const (
	testLogs             = "../.build/noderunner/"
	gethPort             = integration.StartPortNodeRunnerTest + 2
	defaultWsPortOffset  = 100
	gethWebsocketPort    = gethPort + defaultWsPortOffset
	localhost            = "127.0.0.1"
	obscuroWebsocketPort = integration.StartPortNodeRunnerTest + 1
	p2pBindAddress       = integration.StartPortNodeRunnerTest + 3
)

// A smoke test to check that we can stand up a standalone Obscuro host and enclave.
func TestCanStartStandaloneObscuroHostAndEnclave(t *testing.T) {
	testlog.Setup(&testlog.Cfg{
		LogDir:      testLogs,
		TestType:    "noderunner",
		TestSubtype: "test",
		LogLevel:    gethlog.LvlInfo,
	})

	enclaveAddr := fmt.Sprintf("%s:%d", localhost, integration.StartPortNodeRunnerTest)
	rpcAddress := fmt.Sprintf("ws://%s:%d", localhost, obscuroWebsocketPort)

	privateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	hostAddress := crypto.PubkeyToAddress(privateKey.PublicKey)

	hostConfig := config.DefaultHostParsedConfig()
	hostConfig.PrivateKeyString = hex.EncodeToString(crypto.FromECDSA(privateKey))
	hostConfig.EnclaveRPCAddress = enclaveAddr
	hostConfig.ClientRPCPortWS = obscuroWebsocketPort
	hostConfig.L1NodeWebsocketPort = uint(gethWebsocketPort)
	hostConfig.ProfilerEnabled = true
	hostConfig.P2PBindAddress = fmt.Sprintf("0.0.0.0:%d", p2pBindAddress)
	hostConfig.LogPath = testlog.LogFile()

	enclaveConfig := config.DefaultEnclaveConfig()
	enclaveConfig.HostID = hostAddress
	enclaveConfig.Address = enclaveAddr
	dummyContractAddress := common.BytesToAddress([]byte("AA"))
	enclaveConfig.ERC20ContractAddresses = []*common.Address{&dummyContractAddress, &dummyContractAddress}
	enclaveConfig.ProfilerEnabled = true
	enclaveConfig.LogPath = testlog.LogFile()

	gethBinaryPath, err := gethnetwork.EnsureBinariesExist(gethnetwork.LatestVersion)
	if err != nil {
		t.Fatal(err)
	}
	network := gethnetwork.NewGethNetwork(int(gethPort), int(gethWebsocketPort), gethBinaryPath, 1, 1, []string{hostAddress.String()}, "", int(gethlog.LvlDebug))
	defer network.StopNodes()

	go func() {
		err := enclavecontainer.NewEnclaveContainerFromConfig(enclaveConfig).Start()
		if err != nil {
			panic(err)
		}
	}()

	var hostContainer *hostcontainer.HostContainer
	go func() {
		hostContainer = hostcontainer.NewHostContainerFromConfig(hostConfig)
		err = hostContainer.Start()
		if err != nil {
			panic(err)
		}
	}()

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
	defer func() {
		// the container stops the enclave
		if err = hostContainer.Stop(); err != nil {
			t.Fatal("unable to properly stop the host container")
		}
	}()

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
	// We retry 20 times to check if the network has produced any rollups, sleeping half a second between each attempt.
	for counter < 20 {
		counter++
		time.Sleep(500 * time.Millisecond)

		var rollupNumber hexutil.Uint64
		err = obscuroClient.Call(&rollupNumber, rpc.RollupNumber)
		if err == nil && rollupNumber > 0 {
			return
		}
	}

	t.Fatalf("Zero rollups have been produced after ten seconds. Something is wrong. Latest error was: %s", err)
}
