package devnetwork

import (
	"fmt"

	"github.com/obscuronet/go-obscuro/go/ethadapter"
	"github.com/obscuronet/go-obscuro/integration/eth2network"
	"github.com/obscuronet/go-obscuro/integration/simulation/network"
	"github.com/obscuronet/go-obscuro/integration/simulation/params"
)

type gethDockerNetwork struct {
	networkWallets *params.SimWallets
	l1Config       *L1Config
	l1SetupData    *params.L1SetupData
	l1Clients      []ethadapter.EthClient
	ethNetwork     eth2network.Eth2Network
}

func NewGethNetwork(networkWallets *params.SimWallets, l1Config *L1Config) L1Network {
	return &gethDockerNetwork{
		networkWallets: networkWallets,
		l1Config:       l1Config,
	}
}

func (g *gethDockerNetwork) Start() {
	l1SetupData, l1Clients, gethNetwork := network.SetUpGethNetwork(g.networkWallets, g.l1Config.PortStart, g.l1Config.NumNodes, int(g.l1Config.AvgBlockDuration.Seconds()))
	g.l1SetupData = l1SetupData
	g.l1Clients = l1Clients
	g.ethNetwork = gethNetwork
}

func (g *gethDockerNetwork) Stop() {
	err := g.ethNetwork.Stop()
	if err != nil {
		fmt.Println("eth network failed to stop", err)
	}
}

func (g *gethDockerNetwork) NumNodes() int {
	return len(g.l1Clients)
}

func (g *gethDockerNetwork) GetClient(idx int) ethadapter.EthClient {
	return g.l1Clients[idx]
}

func (g *gethDockerNetwork) ObscuroSetupData() *params.L1SetupData {
	return g.l1SetupData
}
