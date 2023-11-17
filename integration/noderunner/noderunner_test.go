package noderunner

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ten-protocol/go-ten/go/common/profiler"
	"github.com/ten-protocol/go-ten/go/node"
	"github.com/ten-protocol/go-ten/go/rpc"
	"github.com/ten-protocol/go-ten/integration"
	"github.com/ten-protocol/go-ten/integration/common/testlog"
	"github.com/ten-protocol/go-ten/integration/eth2network"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
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

	// todo run the noderunner test with different obscuro node instances
	newNode, hostAddr := createInMemoryNode(t)

	binariesPath, err := eth2network.EnsureBinariesExist()
	if err != nil {
		panic(err)
	}

	network := eth2network.NewEth2Network(
		binariesPath,
		true,
		_startPort,
		_startPort+integration.DefaultGethWSPortOffset,
		_startPort+integration.DefaultGethAUTHPortOffset,
		_startPort+integration.DefaultGethNetworkPortOffset,
		_startPort+integration.DefaultPrysmHTTPPortOffset,
		_startPort+integration.DefaultPrysmP2PPortOffset,
		1337,
		1,
		1,
		2,
		2,
		[]string{hostAddr.String()},
		time.Minute,
	)
	defer network.Stop() //nolint: errcheck
	err = network.Start()
	if err != nil {
		panic(err)
	}

	err = newNode.Start()
	if err != nil {
		t.Fatal(err)
	}

	// we create the node RPC client
	wsURL := fmt.Sprintf("ws://127.0.0.1:%d", _startPort+integration.DefaultGethWSPortOffset)
	var obscuroClient rpc.Client
	wait := 30 // max wait in seconds
	for {
		obscuroClient, err = rpc.NewNetworkClient(wsURL)
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
		if err = newNode.Stop(); err != nil {
			t.Fatalf("unable to properly stop the host container - %s", err)
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
		err = obscuroClient.Call(&rollupNumber, rpc.BatchNumber)
		if err == nil && rollupNumber > 0 {
			return
		}
	}

	t.Fatalf("Zero rollups have been produced after ten seconds. Something is wrong. Latest error was: %s", err)
}

func createInMemoryNode(t *testing.T) (node.Node, gethcommon.Address) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	hostAddress := crypto.PubkeyToAddress(privateKey.PublicKey)

	nodeCfg := node.NewNodeConfig(
		node.WithPrivateKey(hex.EncodeToString(crypto.FromECDSA(privateKey))),
		node.WithHostID(hostAddress.String()),
		node.WithEnclaveWSPort(_startPort+integration.DefaultEnclaveOffset),
		node.WithHostHTTPPort(_startPort+integration.DefaultHostRPCHTTPOffset),
		node.WithHostWSPort(_startPort+integration.DefaultHostRPCWSOffset),
		node.WithL1WebsocketURL(fmt.Sprintf("ws://%s:%d", _localhost, _startPort+integration.DefaultGethWSPortOffset)),
		node.WithGenesis(true),
		node.WithProfiler(true),
		node.WithL1BlockTime(1*time.Second),
	)

	return NewInMemNode(nodeCfg), hostAddress
}
