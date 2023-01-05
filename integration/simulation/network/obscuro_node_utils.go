package network

import (
	"fmt"
	"math/big"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/obscuronet/go-obscuro/go/common/container"
	"github.com/obscuronet/go-obscuro/go/common/host"
	"github.com/obscuronet/go-obscuro/go/enclave"

	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/integration/common/testlog"

	"github.com/obscuronet/go-obscuro/go/obsclient"

	"github.com/obscuronet/go-obscuro/go/wallet"

	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/config"
	"github.com/obscuronet/go-obscuro/integration"

	"github.com/obscuronet/go-obscuro/go/ethadapter"
	"github.com/obscuronet/go-obscuro/go/rpc"
	"github.com/obscuronet/go-obscuro/integration/simulation/p2p"
	"github.com/obscuronet/go-obscuro/integration/simulation/params"

	gethcommon "github.com/ethereum/go-ethereum/common"
	enclavecontainer "github.com/obscuronet/go-obscuro/go/enclave/container"
)

const (
	protocolSeparator = "://"
	networkTCP        = "tcp"
)

func startInMemoryObscuroNodes(params *params.SimParams, genesisJSON []byte, l1Clients []ethadapter.EthClient) []rpc.Client {
	// Create the in memory obscuro nodes, each connect each to a geth node
	obscuroNodes := make([]host.Host, params.NumberOfNodes)
	p2pLayers := make([]*p2p.MockP2P, params.NumberOfNodes)
	for i := 0; i < params.NumberOfNodes; i++ {
		isGenesis := i == 0
		p2pLayers[i] = p2p.NewMockP2P(params.AvgBlockDuration, params.AvgNetworkLatency)

		obscuroNodes[i] = createInMemObscuroNode(
			int64(i),
			isGenesis,
			GetNodeType(i),
			params.MgmtContractLib,
			params.ERC20ContractLib,
			true,
			genesisJSON,
			params.Wallets.NodeWallets[i],
			l1Clients[i],
			params.Wallets,
			p2pLayers[i],
			params.L1SetupData.MessageBusAddr,
			params.L1SetupData.ObscuroStartBlock,
		)
	}
	// make sure the aggregators can talk to each other
	for i := 0; i < params.NumberOfNodes; i++ {
		p2pLayers[i].Nodes = obscuroNodes
	}

	// start each obscuro node
	for _, m := range obscuroNodes {
		t := m
		go t.Start()
	}

	// Create a handle to each node
	obscuroClients := make([]rpc.Client, params.NumberOfNodes)
	for i, node := range obscuroNodes {
		obscuroClients[i] = p2p.NewInMemObscuroClient(node)
	}
	time.Sleep(100 * time.Millisecond)

	return obscuroClients
}

func StartStandaloneObscuroNodes(params *params.SimParams, gethClients []ethadapter.EthClient, enclaveAddresses []string) ([]rpc.Client, []string) {
	// handle to the obscuro clients
	nodeRPCAddresses := make([]string, params.NumberOfNodes)
	obscuroClients := make([]rpc.Client, params.NumberOfNodes)
	obscuroNodes := make([]container.Container, params.NumberOfNodes)

	for i := 0; i < params.NumberOfNodes; i++ {
		isGenesis := i == 0

		// We use the convention to determine the rpc ports of the node
		nodeRPCPortHTTP := params.StartPort + DefaultHostRPCHTTPOffset + i
		nodeRPCPortWS := params.StartPort + DefaultHostRPCWSOffset + i

		// create an Obscuro node
		obscuroNodes[i] = createSocketObscuroHostContainer(
			int64(i),
			isGenesis,
			GetNodeType(i),
			fmt.Sprintf("%s:%d", Localhost, params.StartPort+DefaultHostP2pOffset+i),
			enclaveAddresses[i],
			Localhost,
			uint64(nodeRPCPortHTTP),
			uint64(nodeRPCPortWS),
			params.Wallets.NodeWallets[i],
			params.MgmtContractLib,
			gethClients[i],
			params.L1SetupData.ObscuroStartBlock,
		)

		nodeRPCAddresses[i] = fmt.Sprintf("ws://%s:%d", Localhost, nodeRPCPortWS)
	}

	// start each obscuro node
	for _, m := range obscuroNodes {
		// TODO bubble up errors from the containers/host/enclave Start https://github.com/obscuronet/obscuro-internal/issues/1315
		err := m.Start()
		if err != nil {
			testlog.Logger().Crit("unable to start obscuro node ", log.ErrKey, err)
		}
		time.Sleep(params.AvgBlockDuration / 3)
	}

	// create the RPC clients
	for i, rpcAddress := range nodeRPCAddresses {
		var client rpc.Client
		var err error

		// create a connection to the newly created nodes - panic if no connection is made after some time
		startTime := time.Now()
		for connected := false; !connected; time.Sleep(500 * time.Millisecond) {
			client, err = rpc.NewNetworkClient(rpcAddress)
			connected = err == nil // The client cannot be created until the node has started.
			if time.Now().After(startTime.Add(2 * time.Minute)) {
				testlog.Logger().Crit("failed to create a connect to node after 2 minute")
			}

			testlog.Logger().Info(fmt.Sprintf("Could not create client %d. Retrying...", i), log.ErrKey, err)
		}

		obscuroClients[i] = client
	}

	return obscuroClients, nodeRPCAddresses
}

func createAuthClientsPerWallet(clients []rpc.Client, wallets *params.SimWallets) map[string][]*obsclient.AuthObsClient {
	walletClients := make(map[string][]*obsclient.AuthObsClient)
	// loop through all the L2 wallets we're using and round-robin allocate them the rpc clients we have for each host
	for _, w := range append(wallets.SimObsWallets, wallets.L2FaucetWallet) {
		walletClients[w.Address().String()] = createAuthClients(clients, w)
	}
	for _, t := range wallets.Tokens {
		w := t.L2Owner
		walletClients[w.Address().String()] = createAuthClients(clients, w)
	}
	return walletClients
}

func createAuthClients(clients []rpc.Client, wal wallet.Wallet) []*obsclient.AuthObsClient {
	authClients := make([]*obsclient.AuthObsClient, len(clients))
	for i, client := range clients {
		vk, err := rpc.GenerateAndSignViewingKey(wal)
		if err != nil {
			panic(err)
		}
		// todo - use a child logger
		encClient, err := rpc.NewEncRPCClient(client, vk, testlog.Logger())
		if err != nil {
			panic(err)
		}
		authClients[i] = obsclient.NewAuthObsClient(encClient)
	}
	return authClients
}

func StartRemoteEnclaveServers(params *params.SimParams) {
	for i := 0; i < params.NumberOfNodes; i++ {
		// create a remote enclave server
		enclaveAddr := fmt.Sprintf("%s:%d", Localhost, params.StartPort+DefaultEnclaveOffset+i)
		hostAddr := fmt.Sprintf("%s:%d", Localhost, params.StartPort+DefaultHostP2pOffset+i)

		enclaveConfig := config.EnclaveConfig{
			HostID:                 gethcommon.BigToAddress(big.NewInt(int64(i))),
			HostAddress:            hostAddr,
			Address:                enclaveAddr,
			NodeType:               GetNodeType(i),
			L1ChainID:              integration.EthereumChainID,
			ObscuroChainID:         integration.ObscuroChainID,
			ValidateL1Blocks:       false,
			WillAttest:             false,
			GenesisJSON:            nil,
			UseInMemoryDB:          false,
			ERC20ContractAddresses: params.Wallets.AllEthAddresses(),
			MinGasPrice:            big.NewInt(1),
			MessageBusAddress:      *params.L1SetupData.MessageBusAddr,
		}
		enclaveLogger := testlog.Logger().New(log.NodeIDKey, i, log.CmpKey, log.EnclaveCmp)
		encl := enclave.NewEnclave(enclaveConfig, params.MgmtContractLib, params.ERC20ContractLib, enclaveLogger)
		enclaveContainer := enclavecontainer.EnclaveContainer{
			Enclave:   encl,
			RPCServer: enclave.NewEnclaveRPCServer(enclaveConfig.Address, encl, enclaveLogger),
			Logger:    enclaveLogger,
		}
		err := enclaveContainer.Start()
		if err != nil {
			panic(fmt.Sprintf("could not start enclave server: %s", err))
		}
	}
}

// StopObscuroNodes stops the Obscuro nodes and their RPC clients.
func StopObscuroNodes(clients []rpc.Client) {
	var wg sync.WaitGroup
	for _, client := range clients {
		wg.Add(1)
		go func(c rpc.Client) {
			defer wg.Done()
			err := c.Call(nil, rpc.StopHost)
			if err != nil {
				testlog.Logger().Error("Could not stop Obscuro node.", log.ErrKey, err)
			}
			c.Stop()
		}(client)
	}

	if waitTimeout(&wg, 10*time.Second) {
		panic("Timed out waiting for the Obscuro nodes to stop")
	} else {
		testlog.Logger().Info("Obscuro nodes stopped")
	}
}

// CheckHostRPCServersStopped checks whether the hosts' RPC server addresses have been freed up.
func CheckHostRPCServersStopped(hostRPCAddresses []string) {
	var wg sync.WaitGroup
	for _, hostRPCAddress := range hostRPCAddresses {
		wg.Add(1)

		// We cannot stop the RPC server synchronously. This is because the host itself is being stopped by an RPC
		// call, so there is a deadlock. The RPC server is waiting for all connections to close, but a single
		// connection remains open, waiting for the RPC server to close. Instead, we check whether the RPC port
		// becomes free.
		go func(rpcAddress string) {
			defer wg.Done()
			for !isAddressAvailable(rpcAddress) {
				time.Sleep(100 * time.Millisecond)
			}
		}(hostRPCAddress)
	}

	if waitTimeout(&wg, 10*time.Second) {
		panic("Timed out waiting for the Obscuro host RPC addresses to become available")
	} else {
		testlog.Logger().Info("Obscuro host RPC addresses freed")
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

func isAddressAvailable(address string) bool {
	// `net.Listen` requires us to strip the protocol, if it exists.
	addressNoProtocol := address
	splitAddress := strings.Split(address, protocolSeparator)
	if len(splitAddress) == 2 {
		addressNoProtocol = splitAddress[1]
	}

	ln, err := net.Listen(networkTCP, addressNoProtocol)
	if ln != nil {
		err = ln.Close()
		if err != nil {
			testlog.Logger().Error(fmt.Sprintf("could not close listener when checking if address %s was available", address))
		}
	}
	if err != nil {
		return false
	}

	return true
}

// GetNodeType returns the type of the node based on its ID.
func GetNodeType(i int) common.NodeType {
	// Only the genesis node is assigned the role of sequencer.
	if i == 0 {
		return common.Sequencer
	}
	return common.Validator
}
