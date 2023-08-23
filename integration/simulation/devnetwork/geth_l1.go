package devnetwork

import (
	"fmt"
	"github.com/obscuronet/go-obscuro/integration"

	"github.com/obscuronet/go-obscuro/go/ethadapter"
	"github.com/obscuronet/go-obscuro/integration/eth2network"
	"github.com/obscuronet/go-obscuro/integration/simulation/network"
	"github.com/obscuronet/go-obscuro/integration/simulation/params"
)

type gethDockerNetwork struct {
	networkWallets *params.SimWallets
	l1Config       *L1Config
	l1Clients      []ethadapter.EthClient
	ethNetwork     eth2network.Eth2Network
}

func NewGethNetwork(networkWallets *params.SimWallets, l1Config *L1Config) L1Network {
	return &gethDockerNetwork{
		networkWallets: networkWallets,
		l1Config:       l1Config,
	}
}

func (g *gethDockerNetwork) Prepare() {
	gethNetwork, err := network.StartGethNetwork(g.networkWallets, g.l1Config.PortStart, int(g.l1Config.AvgBlockDuration.Seconds()))
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
