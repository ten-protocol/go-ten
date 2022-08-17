package network

import (
	"math"
	"math/big"
	"time"

	"github.com/obscuronet/go-obscuro/go/host/node"

	"github.com/obscuronet/go-obscuro/go/host/rpc/enclaverpc"

	testcommon "github.com/obscuronet/go-obscuro/integration/common"

	"github.com/obscuronet/go-obscuro/integration/simulation/params"

	"github.com/obscuronet/go-obscuro/go/ethadapter"

	"github.com/obscuronet/go-obscuro/go/config"

	"github.com/obscuronet/go-obscuro/go/enclave"
	"github.com/obscuronet/go-obscuro/go/ethadapter/erc20contractlib"
	"github.com/obscuronet/go-obscuro/integration"
	simp2p "github.com/obscuronet/go-obscuro/integration/simulation/p2p"

	"github.com/obscuronet/go-obscuro/go/wallet"

	"github.com/obscuronet/go-obscuro/go/ethadapter/mgmtcontractlib"

	"github.com/obscuronet/go-obscuro/integration/simulation/stats"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/go-obscuro/go/host"
	"github.com/obscuronet/go-obscuro/go/host/p2p"
	ethereum_mock "github.com/obscuronet/go-obscuro/integration/ethereummock"
)

const (
	Localhost                  = "127.0.0.1"
	DefaultWsPortOffset        = 100 // The default offset between a Geth node's port and websocket ports.
	DefaultHostP2pOffset       = 200 // The default offset for the host P2p
	DefaultEnclaveOffset       = 300 // The default offset between a Geth nodes port and the enclave ports. Used in Socket Simulations.
	DefaultHostRPCHTTPOffset   = 400 // The default offset for the host's RPC HTTP port
	DefaultHostRPCWSOffset     = 500 // The default offset for the host's RPC websocket port
	ClientRPCTimeout           = 5 * time.Second
	DefaultL1ConnectionTimeout = 15 * time.Second
)

func createMockEthNode(id int64, nrNodes int, avgBlockDuration time.Duration, avgNetworkLatency time.Duration, stats *stats.Stats) *ethereum_mock.Node {
	mockEthNetwork := ethereum_mock.NewMockEthNetwork(avgBlockDuration, avgNetworkLatency, stats)
	ethereumMockCfg := defaultMockEthNodeCfg(nrNodes, avgBlockDuration)
	// create an in memory mock ethereum node responsible with notifying the layer 2 node about blocks
	miner := ethereum_mock.NewMiner(gethcommon.BigToAddress(big.NewInt(id)), ethereumMockCfg, mockEthNetwork, stats)
	mockEthNetwork.CurrentNode = miner
	return miner
}

func createInMemObscuroNode(
	id int64,
	isGenesis bool,
	mgmtContractLib mgmtcontractlib.MgmtContractLib,
	stableTokenContractLib erc20contractlib.ERC20ContractLib,
	avgGossipPeriod time.Duration,
	avgBlockDuration time.Duration,
	avgNetworkLatency time.Duration,
	stats *stats.Stats,
	validateBlocks bool,
	genesisJSON []byte,
	ethWallet wallet.Wallet,
	ethClient ethadapter.EthClient,
	wallets *params.SimWallets,
) host.MockHost {
	obscuroInMemNetwork := simp2p.NewMockP2P(avgBlockDuration, avgNetworkLatency)

	hostConfig := config.HostConfig{
		ID:                  gethcommon.BigToAddress(big.NewInt(id)),
		IsGenesis:           isGenesis,
		GossipRoundDuration: avgGossipPeriod,
		HasClientRPCHTTP:    false,
	}

	enclaveConfig := config.DefaultEnclaveConfig()
	enclaveConfig.HostID = hostConfig.ID
	enclaveConfig.ValidateL1Blocks = validateBlocks
	enclaveConfig.GenesisJSON = genesisJSON
	enclaveConfig.ERC20ContractAddresses = wallets.AllEthAddresses()

	enclaveClient := enclave.NewEnclave(enclaveConfig, mgmtContractLib, stableTokenContractLib, stats)

	// create an in memory obscuro node
	inMemNode := node.NewHost(
		hostConfig,
		stats,
		obscuroInMemNetwork,
		nil,
		enclaveClient,
		ethWallet,
		mgmtContractLib,
	)
	obscuroInMemNetwork.CurrentNode = inMemNode
	inMemNode.ConnectToEthNode(ethClient)
	return inMemNode
}

func createSocketObscuroNode(
	id int64,
	isGenesis bool,
	avgGossipPeriod time.Duration,
	stats *stats.Stats,
	p2pAddr string,
	enclaveAddr string,
	clientRPCHost string,
	clientRPCPortHTTP uint64,
	clientRPCPortWS uint64,
	ethWallet wallet.Wallet,
	mgmtContractLib mgmtcontractlib.MgmtContractLib,
	ethClient ethadapter.EthClient,
) host.Host {
	hostConfig := config.HostConfig{
		ID:                     gethcommon.BigToAddress(big.NewInt(id)),
		IsGenesis:              isGenesis,
		GossipRoundDuration:    avgGossipPeriod,
		HasClientRPCHTTP:       true,
		ClientRPCPortHTTP:      clientRPCPortHTTP,
		HasClientRPCWebsockets: true,
		ClientRPCPortWS:        clientRPCPortWS,
		ClientRPCHost:          clientRPCHost,
		ClientRPCTimeout:       ClientRPCTimeout,
		EnclaveRPCTimeout:      ClientRPCTimeout,
		EnclaveRPCAddress:      enclaveAddr,
		P2PBindAddress:         p2pAddr,
		L1ChainID:              integration.EthereumChainID,
		ObscuroChainID:         integration.ObscuroChainID,
	}

	// create an enclave client
	enclaveClient := enclaverpc.NewClient(hostConfig)

	// create a socket obscuro node
	nodeP2p := p2p.NewSocketP2PLayer(hostConfig)

	socketNode := node.NewHost(
		hostConfig,
		stats,
		nodeP2p,
		nil,
		enclaveClient,
		ethWallet,
		mgmtContractLib,
	)

	socketNode.ConnectToEthNode(ethClient)
	return socketNode
}

func defaultMockEthNodeCfg(nrNodes int, avgBlockDuration time.Duration) ethereum_mock.MiningConfig {
	return ethereum_mock.MiningConfig{
		PowTime: func() time.Duration {
			// This formula might feel counter-intuitive, but it is a good approximation for Proof of Work.
			// It creates a uniform distribution up to nrMiners*avgDuration
			// Which means on average, every round, the winner (miner who gets the lowest nonce) will pick a number around "avgDuration"
			// while everyone else will have higher values.
			// Over a large number of rounds, the actual average block duration will be around the desired value, while the number of miners who get very close numbers will be limited.
			span := math.Max(2, float64(nrNodes)) // We handle the special cases of zero or one nodes.
			return testcommon.RndBtwTime(avgBlockDuration/time.Duration(span), avgBlockDuration*time.Duration(span))
		},
	}
}
