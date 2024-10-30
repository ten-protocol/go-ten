package noderunner

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ten-protocol/go-ten/go/common/profiler"
	"github.com/ten-protocol/go-ten/go/config"
	"github.com/ten-protocol/go-ten/go/node"
	"github.com/ten-protocol/go-ten/go/rpc"
	"github.com/ten-protocol/go-ten/integration"
	"github.com/ten-protocol/go-ten/integration/common/testlog"
	"github.com/ten-protocol/go-ten/integration/eth2network"

	gethlog "github.com/ethereum/go-ethereum/log"
)

const (
	_testLogs  = "../.build/noderunner/"
	_localhost = "127.0.0.1"
)

// A smoke test to check that we can stand up a standalone TEN host and enclave.
func TestCanStartStandaloneTenHostAndEnclave(t *testing.T) {
	testlog.Setup(&testlog.Cfg{
		LogDir:      _testLogs,
		TestType:    "noderunner",
		TestSubtype: "test",
		LogLevel:    gethlog.LvlInfo,
	})
	startPort := integration.TestPorts.TestCanStartStandaloneTenHostAndEnclavePort
	// todo run the noderunner test with different TEN node instances
	newNode := createInMemoryNode(startPort)

	binDir, err := eth2network.EnsureBinariesExist()
	if err != nil {
		panic(err)
	}

	network := eth2network.NewPosEth2Network(
		binDir,
		startPort+integration.DefaultGethNetworkPortOffset,
		startPort+integration.DefaultPrysmP2PPortOffset,
		startPort+integration.DefaultGethAUTHPortOffset,
		startPort+integration.DefaultGethWSPortOffset,
		startPort+integration.DefaultGethHTTPPortOffset,
		startPort+integration.DefaultPrysmRPCPortOffset,
		startPort+integration.DefaultPrysmGatewayPortOffset,
		integration.EthereumChainID,
		3*time.Minute,
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
	wsURL := fmt.Sprintf("ws://127.0.0.1:%d", startPort+integration.DefaultGethWSPortOffset)
	var tenClient rpc.Client
	wait := 30 // max wait in seconds
	for {
		tenClient, err = rpc.NewNetworkClient(wsURL)
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
		err = tenClient.Call(&rollupNumber, rpc.BatchNumber)
		if err == nil && rollupNumber > 0 {
			return
		}
	}

	t.Fatalf("Zero rollups have been produced after ten seconds. Something is wrong. Latest error was: %s", err)
}

func createInMemoryNode(startPort int) node.Node {
	tenCfg, err := config.LoadTenConfig()
	if err != nil {
		panic(err)
	}

	tenCfg.Node.PrivateKeyString = integration.GethNodePK
	tenCfg.Node.ID = common.HexToAddress(integration.GethNodeAddress)
	tenCfg.Enclave.RPC.BindAddress = fmt.Sprintf("0.0.0.0:%d", startPort+integration.DefaultEnclaveOffset)
	tenCfg.Host.RPC.HTTPPort = uint64(startPort) + integration.DefaultHostRPCHTTPOffset
	tenCfg.Host.RPC.WSPort = uint64(startPort) + integration.DefaultHostRPCWSOffset
	tenCfg.Host.L1.WebsocketURL = fmt.Sprintf("ws://%s:%d", _localhost, startPort+integration.DefaultGethWSPortOffset)
	tenCfg.Node.IsGenesis = true

	return NewInMemNode(tenCfg)
}
