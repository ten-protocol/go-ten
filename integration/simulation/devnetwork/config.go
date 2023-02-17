package devnetwork

import (
	"sync"
	"time"

	"github.com/obscuronet/go-obscuro/integration"
	"github.com/obscuronet/go-obscuro/integration/common/testlog"
	"github.com/obscuronet/go-obscuro/integration/simulation/params"
)

// L1Config tells network admin how to setup the L1 network
type L1Config struct {
	PortStart          int
	WebsocketPortStart int
	NumNodes           int
	AvgBlockDuration   time.Duration
}

// ObscuroConfig tells the L2 node operators how to configure the nodes
type ObscuroConfig struct {
	PortStart         int
	InitNumValidators int
}

// DefaultDevNetwork provides an off-the-shelf default config for a sim network
func DefaultDevNetwork() *InMemDevNetwork {
	numNodes := 4 // Default sim currently uses 4 L1 nodes. Obscuro nodes: 1 seq, 3 validators
	networkWallets := params.NewSimWallets(0, numNodes, integration.EthereumChainID, integration.ObscuroChainID)
	l1Config := &L1Config{
		PortStart:        integration.StartPortSimulationFullNetwork,
		NumNodes:         4,
		AvgBlockDuration: 1 * time.Second,
	}
	l1Network := NewGethNetwork(networkWallets, l1Config)

	return &InMemDevNetwork{
		logger:         testlog.Logger(),
		networkWallets: networkWallets,
		l1Network:      l1Network,
		obscuroConfig: ObscuroConfig{
			PortStart:         integration.StartPortSimulationFullNetwork,
			InitNumValidators: 3,
		},
		faucetLock: sync.Mutex{},
	}
}
