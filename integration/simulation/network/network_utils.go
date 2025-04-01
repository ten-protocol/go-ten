package network

import (
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/ten-protocol/go-ten/go/ethadapter/contractlib"

	gethlog "github.com/ethereum/go-ethereum/log"

	"github.com/ethereum/go-ethereum/params"
	enclaveconfig "github.com/ten-protocol/go-ten/go/enclave/config"
	"github.com/ten-protocol/go-ten/go/host"
	hostconfig "github.com/ten-protocol/go-ten/go/host/config"
	hostcontainer "github.com/ten-protocol/go-ten/go/host/container"
	"github.com/ten-protocol/go-ten/go/host/l1"

	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/common/metrics"
	"github.com/ten-protocol/go-ten/go/enclave"
	"github.com/ten-protocol/go-ten/go/ethadapter"
	"github.com/ten-protocol/go-ten/go/wallet"
	"github.com/ten-protocol/go-ten/integration"
	"github.com/ten-protocol/go-ten/integration/common/testlog"
	"github.com/ten-protocol/go-ten/integration/ethereummock"
	"github.com/ten-protocol/go-ten/integration/simulation/stats"

	gethcommon "github.com/ethereum/go-ethereum/common"
	hostcommon "github.com/ten-protocol/go-ten/go/common/host"
	testcommon "github.com/ten-protocol/go-ten/integration/common"
)

const (
	Localhost               = "127.0.0.1"
	EnclaveClientRPCTimeout = 5 * time.Minute
	DefaultL1RPCTimeout     = 15 * time.Second
)

func createMockEthNode(id int, nrNodes int, avgBlockDuration time.Duration, avgNetworkLatency time.Duration, stats *stats.Stats, blobResolver l1.BlobResolver) *ethereummock.Node {
	mockEthNetwork := ethereummock.NewMockEthNetwork(avgBlockDuration, avgNetworkLatency, stats)
	ethereumMockCfg := defaultMockEthNodeCfg(nrNodes, avgBlockDuration)
	logger := log.New(log.EthereumL1Cmp, int(gethlog.LvlInfo), ethereumMockCfg.LogFile, log.NodeIDKey, id)
	// create an in memory mock ethereum node responsible with notifying the layer 2 node about blocks
	miner := ethereummock.NewMiner(gethcommon.BigToAddress(big.NewInt(int64(id))), ethereumMockCfg, mockEthNetwork, stats, blobResolver, logger)
	mockEthNetwork.CurrentNode = miner
	return miner
}

func createInMemTenNode(
	id int64,
	isGenesis bool,
	nodeType common.NodeType,
	contractRegistryLib contractlib.ContractRegistryLib,
	ethWallet wallet.Wallet,
	ethClient ethadapter.EthClient,
	mockP2P hostcommon.P2PHostService,
	l1BusAddress gethcommon.Address,
	l1BridgeAddress gethcommon.Address,
	l1StartBlk gethcommon.Hash,
	batchInterval time.Duration,
	incomingP2PDisabled bool,
	l1BlockTime time.Duration,
	blobResolver l1.BlobResolver,
) *hostcontainer.HostContainer {
	networkConfigAddr := contractRegistryLib.NetworkConfigLib().GetContractAddr()
	hostConfig := &hostconfig.HostConfig{
		ID:                   fmt.Sprintf("%d", id),
		IsGenesis:            isGenesis,
		NodeType:             nodeType,
		HasClientRPCHTTP:     false,
		P2PPublicAddress:     fmt.Sprintf("%d", id),
		L1StartHash:          l1StartBlk,
		NetworkConfigAddress: *networkConfigAddr,
		BatchInterval:        batchInterval,
		CrossChainInterval:   11 * time.Second, // todo @matt fix where this default comes from
		IsInboundP2PDisabled: incomingP2PDisabled,
		L1BlockTime:          l1BlockTime,
		UseInMemoryDB:        true,
	}

	contracts := contractRegistryLib.GetContractAddresses()

	enclaveConfig := &enclaveconfig.EnclaveConfig{
		NodeID:                          hostConfig.ID,
		L1ChainID:                       integration.EthereumChainID,
		TenChainID:                      integration.TenChainID,
		WillAttest:                      false,
		UseInMemoryDB:                   true,
		MinGasPrice:                     gethcommon.Big1,
		MessageBusAddress:               l1BusAddress,
		BridgeAddress:                   l1BridgeAddress,
		DataAvailabilityRegistryAddress: contracts.DataAvailabilityRegistry,
		EnclaveRegistryAddress:          contracts.EnclaveRegistry,
		SystemContractOwner:             gethcommon.BigToAddress(big.NewInt(1)), // Irrelevant for in-mem nodes
		MaxBatchSize:                    1024 * 55,
		MaxRollupSize:                   1024 * 128,
		BaseFee:                         big.NewInt(1), // todo @siliev:: fix test transaction builders so this can be different
		GasLocalExecutionCapFlag:        params.MaxGasLimit / 2,
		GasBatchExecutionLimit:          30_000_000,
		RPCTimeout:                      5 * time.Second,
		StoreExecutedTransactions:       true,
		DecompressionLimit:              1024 * 1024 * 2,
	}

	enclaveLogger := testlog.Logger().New(log.NodeIDKey, id, log.CmpKey, log.EnclaveCmp)
	enclaveClients := []common.Enclave{enclave.NewEnclave(enclaveConfig, &TestnetGenesis, contractRegistryLib, enclaveLogger)}

	hostLogger := testlog.Logger().New(log.NodeIDKey, id, log.CmpKey, log.HostCmp)
	// create an in memory TEN node
	metricsService := metrics.New(hostConfig.MetricsEnabled, hostConfig.MetricsHTTPPort, hostLogger)
	l1Data := l1.NewL1DataService(ethClient, hostLogger, contractRegistryLib, blobResolver)
	currentContainer := hostcontainer.NewHostContainer(hostConfig, host.NewServicesRegistry(hostLogger), mockP2P, ethClient, l1Data, enclaveClients, ethWallet, nil, hostLogger, metricsService, blobResolver, contractRegistryLib)

	return currentContainer
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
