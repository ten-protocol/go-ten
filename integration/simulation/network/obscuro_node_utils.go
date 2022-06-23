package network

import (
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/config"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave"
	"github.com/obscuronet/obscuro-playground/integration"

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
			params.Wallets.NodeWallets[i],
			l1Clients[i],
			params.ViewingKeysEnabled,
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
	}

	// Create a handle to each node
	obscuroClients := make([]obscuroclient.Client, params.NumberOfNodes)
	for i, node := range obscuroNodes {
		obscuroClients[i] = host.NewInMemObscuroClient(node)
	}
	time.Sleep(100 * time.Millisecond)
	return obscuroClients
}

func startStandaloneObscuroNodes(params *params.SimParams, stats *stats.Stats, gethClients []ethclient.EthClient, enclaveAddresses []string) []obscuroclient.Client {
	// handle to the obscuro clients
	obscuroClients := make([]obscuroclient.Client, params.NumberOfNodes)
	obscuroNodes := make([]*host.Node, params.NumberOfNodes)

	for i := 0; i < params.NumberOfNodes; i++ {
		isGenesis := i == 0

		// We use the convention to determine the rpc ports of the node
		nodeRPCPortHTTP := params.StartPort + DefaultHostRPCHTTPOffset + i
		nodeRPCPortWS := params.StartPort + DefaultHostRPCWSOffset + i

		// create a remote enclave server
		obscuroNodes[i] = createSocketObscuroNode(
			int64(i),
			isGenesis,
			params.AvgGossipPeriod,
			stats,
			fmt.Sprintf("%s:%d", Localhost, params.StartPort+DefaultHostP2pOffset+i),
			enclaveAddresses[i],
			Localhost,
			uint64(nodeRPCPortHTTP),
			uint64(nodeRPCPortWS),
			params.Wallets.NodeWallets[i],
			params.MgmtContractLib,
			gethClients[i],
		)

		nodeRPCAddress := fmt.Sprintf("%s:%d", Localhost, nodeRPCPortHTTP)
		obscuroClients[i] = obscuroclient.NewClient(nodeRPCAddress)
	}

	// start each obscuro node
	for _, m := range obscuroNodes {
		t := m
		go t.Start()
		time.Sleep(params.AvgBlockDuration / 3)
	}

	// wait for the clients to be connected
	for i, client := range obscuroClients {
		started := false
		for !started {
			err := client.Call(nil, obscuroclient.RPCGetID)
			started = err == nil
			if !started {
				fmt.Printf("Could not connect to client %d. Err %s. Retrying..\n", i, err)
			}
			time.Sleep(50 * time.Millisecond)
		}
	}

	return obscuroClients
}

func startRemoteEnclaveServers(startAt int, params *params.SimParams, stats *stats.Stats) {
	for i := startAt; i < params.NumberOfNodes; i++ {
		// create a remote enclave server
		enclaveAddr := fmt.Sprintf("%s:%d", Localhost, params.StartPort+DefaultEnclaveOffset+i)
		hostAddr := fmt.Sprintf("%s:%d", Localhost, params.StartPort+DefaultHostP2pOffset+i)
		enclaveConfig := config.EnclaveConfig{
			HostID:             common.BigToAddress(big.NewInt(int64(i))),
			HostAddress:        hostAddr,
			Address:            enclaveAddr,
			L1ChainID:          integration.EthereumChainID,
			ObscuroChainID:     integration.ObscuroChainID,
			ValidateL1Blocks:   false,
			WillAttest:         false,
			GenesisJSON:        nil,
			UseInMemoryDB:      false,
			ViewingKeysEnabled: params.ViewingKeysEnabled,
		}
		_, err := enclave.StartServer(enclaveConfig, params.MgmtContractLib, params.ERC20ContractLib, stats)
		if err != nil {
			panic(fmt.Sprintf("failed to create enclave server: %v", err))
		}
	}
}

func StopObscuroNodes(clients []obscuroclient.Client) {
	var wg sync.WaitGroup
	for _, client := range clients {
		wg.Add(1)
		go func(c obscuroclient.Client) {
			defer wg.Done()
			err := c.Call(nil, obscuroclient.RPCStopHost)
			if err != nil {
				log.Error("Failed to stop client %s", err)
			}
			c.Stop()
		}(client)
	}
	if waitTimeout(&wg, 2*time.Second) {
		log.Error("Timed out waiting for the obscuro nodes to stop")
	} else {
		log.Info("Obscuro nodes stopped")
	}
	// Wait a bit for the nodes to shut down.
	time.Sleep(2 * time.Second)
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
