package network

import (
	"fmt"
	"math"
	"math/big"
	"time"

	commonhost "github.com/obscuronet/go-obscuro/go/common/host"
	"github.com/obscuronet/go-obscuro/go/host"

	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/common/metrics"
	"github.com/obscuronet/go-obscuro/integration/common/testlog"

	"github.com/obscuronet/go-obscuro/go/host/rpc/enclaverpc"

	testcommon "github.com/obscuronet/go-obscuro/integration/common"

	"github.com/obscuronet/go-obscuro/integration/simulation/params"

	"github.com/obscuronet/go-obscuro/go/ethadapter"
	"github.com/obscuronet/go-obscuro/go/host/rpc/clientrpc"

	"github.com/obscuronet/go-obscuro/go/config"

	"github.com/obscuronet/go-obscuro/go/enclave"
	"github.com/obscuronet/go-obscuro/go/ethadapter/erc20contractlib"
	"github.com/obscuronet/go-obscuro/integration"
	simp2p "github.com/obscuronet/go-obscuro/integration/simulation/p2p"

	"github.com/obscuronet/go-obscuro/go/wallet"

	"github.com/obscuronet/go-obscuro/go/ethadapter/mgmtcontractlib"
	"github.com/obscuronet/go-obscuro/integration/simulation/stats"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/go-obscuro/go/host/p2p"
	"github.com/obscuronet/go-obscuro/integration/ethereummock"
)

const (
	Localhost                = "127.0.0.1"
	DefaultWsPortOffset      = 100 // The default offset between a Geth node's port and websocket ports.
	DefaultHostP2pOffset     = 200 // The default offset for the host P2p
	DefaultEnclaveOffset     = 300 // The default offset between a Geth nodes port and the enclave ports. Used in Socket Simulations.
	DefaultHostRPCHTTPOffset = 400 // The default offset for the host's RPC HTTP port
	DefaultHostRPCWSOffset   = 500 // The default offset for the host's RPC websocket port
	EnclaveClientRPCTimeout  = 5 * time.Second
	DefaultL1RPCTimeout      = 15 * time.Second
)

func createMockEthNode(id int64, nrNodes int, avgBlockDuration time.Duration, avgNetworkLatency time.Duration, stats *stats.Stats) *ethereummock.Node {
	mockEthNetwork := ethereummock.NewMockEthNetwork(avgBlockDuration, avgNetworkLatency, stats)
	ethereumMockCfg := defaultMockEthNodeCfg(nrNodes, avgBlockDuration)
	// create an in memory mock ethereum node responsible with notifying the layer 2 node about blocks
	miner := ethereummock.NewMiner(gethcommon.BigToAddress(big.NewInt(id)), ethereumMockCfg, mockEthNetwork, stats)
	mockEthNetwork.CurrentNode = miner
	return miner
}

func createInMemObscuroNode(
	id int64,
	isGenesis bool,
	nodeType common.NodeType,
	mgmtContractLib mgmtcontractlib.MgmtContractLib,
	stableTokenContractLib erc20contractlib.ERC20ContractLib,
	validateBlocks bool,
	genesisJSON []byte,
	ethWallet wallet.Wallet,
	ethClient ethadapter.EthClient,
	wallets *params.SimWallets,
	mockP2P *simp2p.MockP2P,
	l1BusAddress *gethcommon.Address,
	l1StartBlk gethcommon.Hash,
) commonhost.Host {
	hostConfig := &config.HostConfig{
		ID:               gethcommon.BigToAddress(big.NewInt(id)),
		IsGenesis:        isGenesis,
		NodeType:         nodeType,
		HasClientRPCHTTP: false,
		P2PPublicAddress: fmt.Sprintf("%d", id),
		L1StartHash:      l1StartBlk,
	}

	enclaveConfig := config.EnclaveConfig{
		HostID:                 hostConfig.ID,
		NodeType:               nodeType,
		L1ChainID:              integration.EthereumChainID,
		ObscuroChainID:         integration.ObscuroChainID,
		WillAttest:             false,
		ValidateL1Blocks:       validateBlocks,
		GenesisJSON:            genesisJSON,
		UseInMemoryDB:          true,
		ERC20ContractAddresses: wallets.AllEthAddresses(),
		MinGasPrice:            big.NewInt(1),
		MessageBusAddress:      *l1BusAddress,
	}

	enclaveLogger := testlog.Logger().New(log.NodeIDKey, id, log.CmpKey, log.EnclaveCmp)
	enclaveClient := enclave.NewEnclave(enclaveConfig, mgmtContractLib, stableTokenContractLib, enclaveLogger)

	// create an in memory obscuro node
	hostLogger := testlog.Logger().New(log.NodeIDKey, id, log.CmpKey, log.HostCmp)
	metricsService := metrics.New(hostConfig.MetricsEnabled, hostConfig.MetricsHTTPPort, hostLogger)
	rpcServer := clientrpc.NewServer(hostConfig, hostLogger)
	inMemNode := host.NewHost(hostConfig, mockP2P, ethClient, enclaveClient, rpcServer, ethWallet, mgmtContractLib, hostLogger, metricsService.Registry())
	mockP2P.CurrentNode = inMemNode
	return inMemNode
}

func createSocketObscuroNode(
	id int64,
	isGenesis bool,
	nodeType common.NodeType,
	p2pAddr string,
	enclaveAddr string,
	clientRPCHost string,
	clientRPCPortHTTP uint64,
	clientRPCPortWS uint64,
	ethWallet wallet.Wallet,
	mgmtContractLib mgmtcontractlib.MgmtContractLib,
	ethClient ethadapter.EthClient,
	l1StartBlk gethcommon.Hash,
) commonhost.Host {
	hostConfig := &config.HostConfig{
		ID:                     gethcommon.BigToAddress(big.NewInt(id)),
		IsGenesis:              isGenesis,
		NodeType:               nodeType,
		HasClientRPCHTTP:       true,
		ClientRPCPortHTTP:      clientRPCPortHTTP,
		HasClientRPCWebsockets: true,
		ClientRPCPortWS:        clientRPCPortWS,
		ClientRPCHost:          clientRPCHost,
		EnclaveRPCTimeout:      EnclaveClientRPCTimeout,
		EnclaveRPCAddress:      enclaveAddr,
		P2PBindAddress:         p2pAddr,
		P2PPublicAddress:       p2pAddr,
		L1ChainID:              integration.EthereumChainID,
		ObscuroChainID:         integration.ObscuroChainID,
		L1StartHash:            l1StartBlk,
	}

	// create an enclave client
	enclaveClient := enclaverpc.NewClient(hostConfig, testlog.Logger().New(log.NodeIDKey, id))

	hostLogger := testlog.Logger().New(log.NodeIDKey, id, log.CmpKey, log.HostCmp)

	// create the metrics service
	metricsService := metrics.New(hostConfig.MetricsEnabled, hostConfig.MetricsHTTPPort, hostLogger)

	// create a socket P2P layer
	p2pLogger := hostLogger.New(log.CmpKey, log.P2PCmp)
	nodeP2p := p2p.NewSocketP2PLayer(hostConfig, p2pLogger, metricsService.Registry())
	rpcServer := clientrpc.NewServer(hostConfig, hostLogger)

	return host.NewHost(hostConfig, nodeP2p, ethClient, enclaveClient, rpcServer, ethWallet, mgmtContractLib, hostLogger, metricsService.Registry())
}

func defaultMockEthNodeCfg(nrNodes int, avgBlockDuration time.Duration) ethereummock.MiningConfig {
	return ethereummock.MiningConfig{
		PowTime: func() time.Duration {
			// This formula might feel counter-intuitive, but it is a good approximation for Proof of Work.
			// It creates a uniform distribution up to nrMiners*avgDuration
			// Which means on average, every round, the winner (miner who gets the lowest nonce) will pick a number around "avgDuration"
			// while everyone else will have higher values.
			// Over a large number of rounds, the actual average block duration will be around the desired value, while the number of miners who get very close numbers will be limited.
			span := math.Max(2, float64(nrNodes)) // We handle the special cases of zero or one nodes.
			return testcommon.RndBtwTime(avgBlockDuration/time.Duration(span), avgBlockDuration*time.Duration(span))
		},
		LogFile: testlog.LogFile(),
	}
}
