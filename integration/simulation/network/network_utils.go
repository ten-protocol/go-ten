package network

import (
	"fmt"
	"math"
	"math/big"
	"time"

	gethlog "github.com/ethereum/go-ethereum/log"
	hostcommon "github.com/obscuronet/go-obscuro/go/common/host"
	hostenclave "github.com/obscuronet/go-obscuro/go/host/enclave"
	"github.com/obscuronet/go-obscuro/go/host/events"
	"github.com/obscuronet/go-obscuro/go/host/l1"
	"github.com/obscuronet/go-obscuro/go/host/l2"
	"github.com/obscuronet/go-obscuro/go/rpc"

	"github.com/obscuronet/go-obscuro/go/host"

	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/common/metrics"
	"github.com/obscuronet/go-obscuro/go/config"
	"github.com/obscuronet/go-obscuro/go/enclave"
	"github.com/obscuronet/go-obscuro/go/enclave/genesis"
	"github.com/obscuronet/go-obscuro/go/ethadapter"
	"github.com/obscuronet/go-obscuro/go/ethadapter/mgmtcontractlib"
	"github.com/obscuronet/go-obscuro/go/host/container"
	"github.com/obscuronet/go-obscuro/go/wallet"
	"github.com/obscuronet/go-obscuro/integration"
	"github.com/obscuronet/go-obscuro/integration/common/testlog"
	"github.com/obscuronet/go-obscuro/integration/ethereummock"
	"github.com/obscuronet/go-obscuro/integration/simulation/stats"

	gethcommon "github.com/ethereum/go-ethereum/common"
	testcommon "github.com/obscuronet/go-obscuro/integration/common"
)

const (
	Localhost               = "127.0.0.1"
	EnclaveClientRPCTimeout = 5 * time.Minute
	DefaultL1RPCTimeout     = 15 * time.Second
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
	validateBlocks bool,
	genesisJSON []byte,
	ethWallet wallet.Wallet,
	ethClient ethadapter.EthClient,
	mockP2PFactory host.ServiceFactory[host.P2PService],
	l1BusAddress *gethcommon.Address,
	l1StartBlk gethcommon.Hash,
	batchInterval time.Duration,
) (*container.HostContainer, rpc.Client) {
	mgtContractAddress := mgmtContractLib.GetContractAddr()

	hostConfig := &config.HostConfig{
		ID:                        gethcommon.BigToAddress(big.NewInt(id)),
		IsGenesis:                 isGenesis,
		NodeType:                  nodeType,
		HasClientRPCHTTP:          false,
		P2PPublicAddress:          fmt.Sprintf("%d", id),
		L1StartHash:               l1StartBlk,
		ManagementContractAddress: *mgtContractAddress,
		BatchInterval:             batchInterval,
	}

	enclaveConfig := &config.EnclaveConfig{
		SequencerID:               gethcommon.BigToAddress(big.NewInt(0)),
		HostID:                    hostConfig.ID,
		NodeType:                  nodeType,
		L1ChainID:                 integration.EthereumChainID,
		ObscuroChainID:            integration.ObscuroChainID,
		WillAttest:                false,
		ValidateL1Blocks:          validateBlocks,
		GenesisJSON:               genesisJSON,
		UseInMemoryDB:             true,
		MinGasPrice:               big.NewInt(1),
		MessageBusAddress:         *l1BusAddress,
		ManagementContractAddress: *mgtContractAddress,
		MaxBatchSize:              1024 * 25,
		MaxRollupSize:             1024 * 64,
	}

	enclaveLogger := testlog.Logger().New(log.NodeIDKey, id, log.CmpKey, log.EnclaveCmp)
	enclaveClient := enclave.NewEnclave(enclaveConfig, &genesis.TestnetGenesis, mgmtContractLib, enclaveLogger)

	// create an in memory obscuro node
	hostLogger := testlog.Logger().New(log.NodeIDKey, id, log.CmpKey, log.HostCmp)

	currentContainer := &container.HostContainer{Logger: hostLogger}

	mockRPCServerFactory, mockRPCClient := SetupMockServer(currentContainer.Stop)

	h := host.NewHost(hostConfig, hostLogger,
		mockP2PFactory, // we replace the real p2p service with a service that passes msgs directly between hosts
		mockEthClientL1RepoFactory(ethClient),
		mockInjectedL1PublisherFactory(ethClient, ethWallet, mgmtContractLib),
		l2.RepoFactory,
		injectedEnclaveFactory(enclaveClient), // need a different factory for enclave service because no RPC client to create
		events.LogEventFactory,
		metrics.ServiceFactory,
		// the rpc server currently needs a way to stop the host, so we wire through this callback (aim to remove this soon)
		mockRPCServerFactory,
		host.DBServiceFactory,
	)
	currentContainer.Host = h

	return currentContainer, mockRPCClient
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

// The following factories are used to initialise the host services with the mocked dependencies

func injectedEnclaveFactory(enclaveClient common.Enclave) host.ServiceFactory[host.EnclaveHostService] {
	return func(config *config.HostConfig, serviceLocator host.ServiceLocator, logger gethlog.Logger) (host.EnclaveHostService, error) {
		hostData := hostcommon.NewIdentity(config)

		guardian := hostenclave.NewGuardian(hostData, config.BatchInterval, config.RollupInterval, config.L1StartHash, serviceLocator,
			enclaveClient, logger.New(log.CmpKey, "guardian-0"))

		return hostenclave.NewService(hostData, serviceLocator, guardian, logger), nil
	}
}

func mockEthClientL1RepoFactory(ethClient ethadapter.EthClient) host.ServiceFactory[host.L1BlockRepositoryService] {
	return func(_ *config.HostConfig, _ host.ServiceLocator, logger gethlog.Logger) (host.L1BlockRepositoryService, error) {
		return l1.NewL1Repository(ethClient, logger), nil
	}
}

func mockInjectedL1PublisherFactory(ethClient ethadapter.EthClient, ethWallet wallet.Wallet, mgmtContractLib mgmtcontractlib.MgmtContractLib) host.ServiceFactory[host.L1PublisherService] {
	return func(config *config.HostConfig, serviceLocator host.ServiceLocator, logger gethlog.Logger) (host.L1PublisherService, error) {
		return l1.NewL1Publisher(hostcommon.NewIdentity(config), ethWallet, ethClient, mgmtContractLib, serviceLocator, logger), nil
	}
}
