package network

import (
	"math"
	"math/big"
	"time"

	"github.com/obscuronet/obscuro-playground/go/common"

	"github.com/obscuronet/obscuro-playground/integration/simulation/params"

	"github.com/obscuronet/obscuro-playground/go/ethclient"

	"github.com/obscuronet/obscuro-playground/go/config"

	"github.com/obscuronet/obscuro-playground/go/enclave"
	"github.com/obscuronet/obscuro-playground/go/ethclient/erc20contractlib"
	"github.com/obscuronet/obscuro-playground/integration"
	simp2p "github.com/obscuronet/obscuro-playground/integration/simulation/p2p"

	"github.com/obscuronet/obscuro-playground/go/wallet"

	"github.com/obscuronet/obscuro-playground/go/ethclient/mgmtcontractlib"

	"github.com/obscuronet/obscuro-playground/integration/simulation/stats"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/obscuro-playground/go/host"
	"github.com/obscuronet/obscuro-playground/go/host/p2p"
	ethereum_mock "github.com/obscuronet/obscuro-playground/integration/ethereummock"
)

const (
	Localhost                  = "127.0.0.1"
	DefaultWsPortOffset        = 100 // The default offset between a Geth node's port and websocket ports.
	DefaultHostP2pOffset       = 200 // The default offset for the host P2p
	DefaultEnclaveOffset       = 300 // The default offset between a Geth nodes port and the enclave ports. Used in Socket Simulations.
	DefaultHostRPCHTTPOffset   = 400 // The default offset for the host's RPC HTTP port
	DefaultHostRPCWSOffset     = 500 // The default offset for the host's RPC websocket port
	ClientRPCTimeoutSecs       = 5
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
	ethClient ethclient.EthClient,
	viewingKeysEnabled bool,
	wallets *params.SimWallets,
) *host.Node {
	obscuroInMemNetwork := simp2p.NewMockP2P(avgBlockDuration, avgNetworkLatency)

	hostConfig := config.HostConfig{
		ID:                  gethcommon.BigToAddress(big.NewInt(id)),
		IsGenesis:           isGenesis,
		GossipRoundDuration: avgGossipPeriod,
		HasClientRPCHTTP:    false,
	}

	enclaveConfig := config.EnclaveConfig{
		HostID:                 hostConfig.ID,
		L1ChainID:              integration.EthereumChainID,
		ObscuroChainID:         integration.ObscuroChainID,
		WillAttest:             false,
		ValidateL1Blocks:       validateBlocks,
		GenesisJSON:            genesisJSON,
		UseInMemoryDB:          true,
		ViewingKeysEnabled:     viewingKeysEnabled,
		ERC20ContractAddresses: wallets.AllEthAddresses(),
	}
	enclaveClient := enclave.NewEnclave(enclaveConfig, mgmtContractLib, stableTokenContractLib, stats)

	// create an in memory obscuro node
	node := host.NewHost(
		hostConfig,
		stats,
		obscuroInMemNetwork,
		nil,
		enclaveClient,
		ethWallet,
		mgmtContractLib,
	)
	obscuroInMemNetwork.CurrentNode = node
	node.ConnectToEthNode(ethClient)
	return node
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
	ethClient ethclient.EthClient,
) *host.Node {
	hostConfig := config.HostConfig{
		ID:                     gethcommon.BigToAddress(big.NewInt(id)),
		IsGenesis:              isGenesis,
		GossipRoundDuration:    avgGossipPeriod,
		HasClientRPCHTTP:       true,
		ClientRPCPortHTTP:      clientRPCPortHTTP,
		HasClientRPCWebsockets: true,
		ClientRPCPortWS:        clientRPCPortWS,
		ClientRPCHost:          clientRPCHost,
		ClientRPCTimeout:       ClientRPCTimeoutSecs * time.Second,
		EnclaveRPCTimeout:      ClientRPCTimeoutSecs * time.Second,
		EnclaveRPCAddress:      enclaveAddr,
		P2PAddress:             p2pAddr,
		ChainID:                config.DefaultHostConfig().ChainID,
	}

	// create an enclave client
	enclaveClient := host.NewEnclaveRPCClient(hostConfig)

	// create a socket obscuro node
	nodeP2p := p2p.NewSocketP2PLayer(hostConfig)

	node := host.NewHost(
		hostConfig,
		stats,
		nodeP2p,
		nil,
		enclaveClient,
		ethWallet,
		mgmtContractLib,
	)

	node.ConnectToEthNode(ethClient)
	return node
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
			return common.RndBtwTime(avgBlockDuration/time.Duration(span), avgBlockDuration*time.Duration(span))
		},
	}
}
