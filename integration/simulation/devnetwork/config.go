package devnetwork

import (
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/obscuronet/go-obscuro/go/enclave/genesis"
	"github.com/obscuronet/go-obscuro/go/wallet"
	testcommon "github.com/obscuronet/go-obscuro/integration/common"
	"math/big"
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
	BatchInterval     time.Duration
	RollupInterval    time.Duration
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
		networkWallets: networkWallets,
		l1Network:      l1Network,
		obscuroConfig: ObscuroConfig{
			PortStart:         integration.StartPortSimulationFullNetwork,
			InitNumValidators: 3,
			BatchInterval:     1 * time.Second,
			RollupInterval:    10 * time.Second,
		},
		faucetLock: sync.Mutex{},
	}
}

func LiveL1DevNetwork(seqWallet wallet.Wallet, validatorWallets []wallet.Wallet, rpcAddress string) *InMemDevNetwork {
	// setup the host and deployer wallets to be the prefunded wallets

	// create the L2 faucet wallet
	l2FaucetPrivKey, err := crypto.HexToECDSA(genesis.TestnetPrefundedPK)
	if err != nil {
		panic("could not initialise L2 faucet private key")
	}
	l2FaucetWallet := wallet.NewInMemoryWalletFromPK(big.NewInt(integration.ObscuroChainID), l2FaucetPrivKey, testlog.Logger())
	networkWallets := &params.SimWallets{
		MCOwnerWallet:  seqWallet,
		NodeWallets:    append([]wallet.Wallet{seqWallet}, validatorWallets...),
		L2FaucetWallet: l2FaucetWallet,
		Tokens:         map[testcommon.ERC20]*params.SimToken{},
	}

	l1Network := &liveL1Network{
		deployWallet:     seqWallet, // use the same wallet for deploying the contracts
		seqWallet:        seqWallet,
		validatorWallets: validatorWallets,
		rpcAddress:       rpcAddress,
	}

	return &InMemDevNetwork{
		logger:         testlog.Logger(),
		networkWallets: networkWallets,
		l1Network:      l1Network,
		obscuroConfig: ObscuroConfig{
			PortStart:         integration.StartPortSimulationFullNetwork,
			InitNumValidators: len(validatorWallets),
			BatchInterval:     5 * time.Second,
			RollupInterval:    1 * time.Minute,
		},
	}
}
