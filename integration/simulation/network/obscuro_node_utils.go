package network

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/obscuronet/go-obscuro/go/common/host"
	"github.com/obscuronet/go-obscuro/go/ethadapter"
	hostcontainer "github.com/obscuronet/go-obscuro/go/host/container"
	"github.com/obscuronet/go-obscuro/integration/simulation/p2p"

	"github.com/obscuronet/go-obscuro/go/common/viewingkey"

	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/obsclient"
	"github.com/obscuronet/go-obscuro/go/rpc"
	"github.com/obscuronet/go-obscuro/go/wallet"
	"github.com/obscuronet/go-obscuro/integration/common/testlog"
	"github.com/obscuronet/go-obscuro/integration/simulation/params"
	"golang.org/x/sync/errgroup"
)

const (
	protocolSeparator = "://"
	networkTCP        = "tcp"
)

func startInMemoryObscuroNodes(params *params.SimParams, genesisJSON []byte, l1Clients []ethadapter.EthClient) []rpc.Client {
	// Create the in memory obscuro nodes, each connect each to a geth node
	obscuroNodes := make([]*hostcontainer.HostContainer, params.NumberOfNodes)
	obscuroHosts := make([]host.Host, params.NumberOfNodes)
	mockP2PNetw := p2p.NewMockP2PNetwork(params.AvgBlockDuration, params.AvgNetworkLatency, params.NodeWithInboundP2PDisabled)

	for i := 0; i < params.NumberOfNodes; i++ {
		isGenesis := i == 0

		obscuroNodes[i] = createInMemObscuroNode(
			int64(i),
			isGenesis,
			GetNodeType(i),
			params.MgmtContractLib,
			true,
			genesisJSON,
			params.Wallets.NodeWallets[i],
			l1Clients[i],
			mockP2PNetw.NewNode(i),
			params.L1SetupData.MessageBusAddr,
			params.L1SetupData.ObscuroStartBlock,
			params.AvgBlockDuration/3,
			true,
			params.AvgBlockDuration,
		)
		obscuroHosts[i] = obscuroNodes[i].Host()
	}

	// start each obscuro node
	for _, m := range obscuroNodes {
		t := m
		go func() {
			err := t.Start()
			if err != nil {
				panic(err)
			}
		}()
	}

	// Create a handle to each node
	obscuroClients := make([]rpc.Client, params.NumberOfNodes)
	for i, node := range obscuroNodes {
		obscuroClients[i] = p2p.NewInMemObscuroClient(node)
	}
	time.Sleep(100 * time.Millisecond)

	return obscuroClients
}

func createAuthClientsPerWallet(clients []rpc.Client, wallets *params.SimWallets) map[string][]*obsclient.AuthObsClient {
	walletClients := make(map[string][]*obsclient.AuthObsClient)
	// loop through all the L2 wallets we're using and round-robin allocate them the rpc clients we have for each host
	for _, w := range append(wallets.SimObsWallets, wallets.L2FaucetWallet) {
		walletClients[w.Address().String()] = CreateAuthClients(clients, w)
	}
	for _, t := range wallets.Tokens {
		w := t.L2Owner
		walletClients[w.Address().String()] = CreateAuthClients(clients, w)
	}
	return walletClients
}

func CreateAuthClients(clients []rpc.Client, wal wallet.Wallet) []*obsclient.AuthObsClient {
	authClients := make([]*obsclient.AuthObsClient, len(clients))
	for i, client := range clients {
		vk, err := viewingkey.GenerateViewingKeyForWallet(wal)
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

// StopObscuroNodes stops the Obscuro nodes and their RPC clients.
func StopObscuroNodes(clients []rpc.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	eg, _ := errgroup.WithContext(ctx)

	for _, client := range clients {
		c := client
		eg.Go(func() error {
			err := c.Call(nil, rpc.StopHost)
			if err != nil {
				testlog.Logger().Error("Could not stop Obscuro node.", log.ErrKey, err)
				return err
			}
			c.Stop()
			return nil
		})
	}

	err := eg.Wait()
	if err != nil {
		testlog.Logger().Error(fmt.Sprintf("Error waiting for the Obscuro nodes to stop - %s", err))
	}

	testlog.Logger().Info("Obscuro nodes stopped")
}

// CheckHostRPCServersStopped checks whether the hosts' RPC server addresses have been freed up.
func CheckHostRPCServersStopped(hostWSURLS []string) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	eg, _ := errgroup.WithContext(ctx)

	for _, hostWSURL := range hostWSURLS {
		url := hostWSURL
		// We cannot stop the RPC server synchronously. This is because the host itself is being stopped by an RPC
		// call, so there is a deadlock. The RPC server is waiting for all connections to close, but a single
		// connection remains open, waiting for the RPC server to close. Instead, we check whether the RPC port
		// becomes free.
		eg.Go(func() error {
			for !isAddressAvailable(url) {
				time.Sleep(100 * time.Millisecond)
			}
			return nil
		})
	}

	err := eg.Wait()
	if err != nil {
		panic(fmt.Sprintf("Timed out waiting for the Obscuro host RPC addresses to become available - %s", err))
	}

	testlog.Logger().Info("Obscuro host RPC addresses freed")
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
