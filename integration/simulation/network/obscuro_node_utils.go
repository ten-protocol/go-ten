package network

import (
	"fmt"
	"sync"
	"time"

	"github.com/obscuronet/obscuro-playground/go/log"

	"github.com/obscuronet/obscuro-playground/go/ethclient"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/host"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/obscuroclient"
	"github.com/obscuronet/obscuro-playground/integration/simulation/p2p"
	"github.com/obscuronet/obscuro-playground/integration/simulation/params"
	"github.com/obscuronet/obscuro-playground/integration/simulation/stats"
)

func startInMemoryObscuroNodes(params *params.SimParams, stats *stats.Stats, genesisJSON []byte, l1Clients []ethclient.EthClient) []obscuroclient.Client {
	// Create the in memory obscuro nodes, each connect each to a geth node
	obscuroNodes := make([]*host.Node, params.NumberOfNodes)
	for i := 0; i < params.NumberOfNodes; i++ {
		isGenesis := i == 0
		obscuroNodes[i] = createInMemObscuroNode(
			int64(i),
			isGenesis,
			params.MgmtContractLib,
			params.ERC20ContractLib,
			params.AvgGossipPeriod,
			params.AvgBlockDuration,
			params.AvgNetworkLatency,
			stats,
			true,
			genesisJSON,
			params.NodeEthWallets[i],
			l1Clients[i],
		)
	}
	// make sure the aggregators can talk to each other
	for _, m := range obscuroNodes {
		mockP2P := m.P2p.(*p2p.MockP2P)
		mockP2P.Nodes = obscuroNodes
	}

	// start each obscuro node
	for _, m := range obscuroNodes {
		t := m
		go t.Start()
		time.Sleep(params.AvgBlockDuration / 10)
	}

	// Create a handle to each node
	obscuroClients := make([]obscuroclient.Client, params.NumberOfNodes)
	for i, node := range obscuroNodes {
		obscuroClients[i] = host.NewInMemObscuroClient(node)
	}
	return obscuroClients
}

func startStandaloneObscuroNodes(params *params.SimParams, stats *stats.Stats, gethClients []ethclient.EthClient) ([]obscuroclient.Client, []string) {
	// handle to the obscuro clients
	obscuroClients := make([]obscuroclient.Client, params.NumberOfNodes)

	obscuroNodes := make([]*host.Node, params.NumberOfNodes)
	nodeP2pAddrs := make([]string, params.NumberOfNodes)

	for i := 0; i < params.NumberOfNodes; i++ {
		// We assign a P2P address to each node on the network according to the convention.
		nodeP2pAddrs[i] = fmt.Sprintf("%s:%d", Localhost, params.StartPort+DefaultHostP2pOffset+i)
	}

	for i := 0; i < params.NumberOfNodes; i++ {
		isGenesis := i == 0

		// We use the convention to determine the rpc ports of the node and the enclave
		nodeRpcAddress := fmt.Sprintf("%s:%d", Localhost, params.StartPort+DefaultHostRPCOffset+i)
		enclaveAddress := fmt.Sprintf("%s:%d", Localhost, params.StartPort+DefaultEnclaveOffset+i)

		// create a remote enclave server
		obscuroNodes[i] = createSocketObscuroNode(
			int64(i+1),
			isGenesis,
			params.AvgGossipPeriod,
			stats,
			nodeP2pAddrs[i],
			nodeP2pAddrs,
			enclaveAddress,
			nodeRpcAddress,
			params.NodeEthWallets[i],
			params.MgmtContractLib,
			gethClients[i],
		)
		obscuroClients[i] = obscuroclient.NewClient(nodeRpcAddress)
	}

	// start each obscuro node
	for _, m := range obscuroNodes {
		t := m
		go t.Start()
		time.Sleep(params.AvgBlockDuration / 3)
	}

	return obscuroClients, nodeP2pAddrs
}

func StopObscuroNodes(clients []obscuroclient.Client) {
	var wg sync.WaitGroup
	for _, client := range clients {
		wg.Add(1)
		temp := client
		go func() {
			defer wg.Done()
			err := temp.Call(nil, obscuroclient.RPCStopHost)
			if err != nil {
				log.Error("Failed to stop client %s", err)
			}
			temp.Stop()
		}()
	}
	if waitTimeout(&wg, 2*time.Second) {
		log.Error("Timed out waiting for the obscuro nodes to stop")
	} else {
		log.Info("Obscuro nodes stopped")
	}
}

// waitTimeout waits for the waitgroup for the specified max timeout.
// Returns true if waiting timed out.
func waitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return false // completed normally
	case <-time.After(timeout):
		return true // timed out
	}
}
