package noderunner

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/obscuronet/go-obscuro/go/common/profiler"
	"github.com/obscuronet/go-obscuro/go/config"
	"github.com/obscuronet/go-obscuro/go/rpc"
	"github.com/obscuronet/go-obscuro/integration"
	"github.com/obscuronet/go-obscuro/integration/common/testlog"
	"github.com/obscuronet/go-obscuro/integration/eth2network"

	gethlog "github.com/ethereum/go-ethereum/log"
	enclavecontainer "github.com/obscuronet/go-obscuro/go/enclave/container"
	hostcontainer "github.com/obscuronet/go-obscuro/go/host/container"
)

const (
	_testLogs  = "../.build/noderunner/"
	_localhost = "127.0.0.1"
	_startPort = integration.StartPortNodeRunnerTest
)

// A smoke test to check that we can stand up a standalone Obscuro host and enclave.
func TestCanStartStandaloneObscuroHostAndEnclave(t *testing.T) {
	testlog.Setup(&testlog.Cfg{
		LogDir:      _testLogs,
		TestType:    "noderunner",
		TestSubtype: "test",
		LogLevel:    gethlog.LvlInfo,
	})

	enclaveAddr := fmt.Sprintf("%s:%d", _localhost, _startPort+integration.DefaultEnclaveOffset)
	rpcAddress := fmt.Sprintf("ws://%s:%d", _localhost, _startPort+integration.DefaultGethWSPortOffset)

	privateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	hostAddress := crypto.PubkeyToAddress(privateKey.PublicKey)

	hostConfig := config.DefaultHostParsedConfig()
	hostConfig.PrivateKeyString = hex.EncodeToString(crypto.FromECDSA(privateKey))
	hostConfig.EnclaveRPCAddress = enclaveAddr
	hostConfig.HasClientRPCHTTP = false
	hostConfig.ClientRPCPortWS = _startPort + integration.DefaultHostRPCWSOffset
	hostConfig.L1NodeWebsocketPort = uint(_startPort + integration.DefaultGethWSPortOffset)
	hostConfig.ProfilerEnabled = true
	hostConfig.P2PBindAddress = fmt.Sprintf("0.0.0.0:%d", _startPort+integration.DefaultHostP2pOffset)
	hostConfig.LogPath = testlog.LogFile()

	enclaveConfig := config.DefaultEnclaveConfig()
	enclaveConfig.HostID = hostAddress
	enclaveConfig.Address = enclaveAddr
	enclaveConfig.ProfilerEnabled = true
	enclaveConfig.LogPath = testlog.LogFile()
	enclaveConfig.WillAttest = false

	binariesPath, err := eth2network.EnsureBinariesExist()
	if err != nil {
		panic(err)
	}

	network := eth2network.NewEth2Network(
		binariesPath,
		_startPort,
		_startPort+integration.DefaultGethWSPortOffset,
		_startPort+integration.DefaultGethAUTHPortOffset,
		_startPort+integration.DefaultGethNetworkPortOffset,
		_startPort+integration.DefaultPrysmHTTPPortOffset,
		_startPort+integration.DefaultPrysmP2PPortOffset,
		1337,
		1,
		1,
		[]string{hostAddress.String()},
	)
	defer network.Stop() //nolint: errcheck
	err = network.Start()
	if err != nil {
		panic(err)
	}

	go func() {
		err := enclavecontainer.NewEnclaveContainerFromConfig(enclaveConfig).Start()
		if err != nil {
			panic(err)
		}
	}()

	hostContainer := hostcontainer.NewHostContainerFromConfig(hostConfig)
	err = hostContainer.Start()
	if err != nil {
		panic(err)
	}

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
	_, err = http.Get(fmt.Sprintf("http://%s:%d", _localhost, profiler.DefaultHostPort)) //nolint
	if err != nil {
		t.Errorf("host profiler is not reachable: %s", err)
	}

	counter := 0
	// We retry 20 times to check if the network has produced any rollups, sleeping one second between each attempt.
	for counter < 20 {
		counter++
		time.Sleep(time.Second)

		var rollupNumber hexutil.Uint64
		err = obscuroClient.Call(&rollupNumber, rpc.RollupNumber)
		if err == nil && rollupNumber > 0 {
			return
		}
	}

	t.Fatalf("Zero rollups have been produced after ten seconds. Something is wrong. Latest error was: %s", err)
}
