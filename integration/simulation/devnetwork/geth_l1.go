package devnetwork

import (
	"fmt"
	"time"

	"github.com/ten-protocol/go-ten/integration"

	"github.com/ten-protocol/go-ten/go/ethadapter"
	"github.com/ten-protocol/go-ten/integration/eth2network"
	"github.com/ten-protocol/go-ten/integration/simulation/network"
	"github.com/ten-protocol/go-ten/integration/simulation/params"
)

type gethDockerNetwork struct {
	networkWallets *params.SimWallets
	l1Config       *L1Config
	l1Clients      []ethadapter.EthClient
	ethNetwork     eth2network.PosEth2Network
}

func NewGethNetwork(networkWallets *params.SimWallets, l1Config *L1Config) L1Network {
	return &gethDockerNetwork{
		networkWallets: networkWallets,
		l1Config:       l1Config,
	}
}

func (g *gethDockerNetwork) Prepare() {
	gethNetwork, err := network.StartGethNetwork(g.networkWallets, g.l1Config.PortStart)
	if err != nil {
		panic(err)
	}
	g.l1Clients = make([]ethadapter.EthClient, g.l1Config.NumNodes)
	for i := 0; i < g.l1Config.NumNodes; i++ {
		g.l1Clients[i] = network.CreateEthClientConnection(int64(i), uint(g.l1Config.PortStart+integration.DefaultGethWSPortOffset))
	}
	g.ethNetwork = gethNetwork
}

func (g *gethDockerNetwork) CleanUp() {
	err := g.ethNetwork.Stop()
	if err != nil {
		fmt.Println("eth network failed to stop", err)
	}
}

func (g *gethDockerNetwork) NumNodes() int {
	return len(g.l1Clients)
}

func (g *gethDockerNetwork) GetClient(_ int) ethadapter.EthClient {
	return g.l1Clients[0]
}

func (g *gethDockerNetwork) GetBlockTime() time.Duration {
	return g.l1Config.AvgBlockDuration
}
