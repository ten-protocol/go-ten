package network

import (
	"math/big"
	"time"

	"github.com/obscuronet/obscuro-playground/go/ethclient/mgmtcontractlib"

	p2p2 "github.com/obscuronet/obscuro-playground/integration/simulation/p2p"

	"github.com/obscuronet/obscuro-playground/integration/simulation/stats"

	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/host"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/host/p2p"
	ethereum_mock "github.com/obscuronet/obscuro-playground/integration/ethereummock"
)

const (
	Localhost            = "127.0.0.1"
	DefaultWsPortOffset  = 100 // The default offset between a Geth node's HTTP and websocket ports.
	ClientRPCTimeoutSecs = 5
)

func createMockEthNode(id int64, nrNodes int, avgBlockDuration time.Duration, avgNetworkLatency time.Duration, stats *stats.Stats) *ethereum_mock.Node {
	mockEthNetwork := ethereum_mock.NewMockEthNetwork(avgBlockDuration, avgNetworkLatency, stats)
	ethereumMockCfg := defaultMockEthNodeCfg(nrNodes, avgBlockDuration)
	// create an in memory mock ethereum node responsible with notifying the layer 2 node about blocks
	miner := ethereum_mock.NewMiner(common.BigToAddress(big.NewInt(id)), ethereumMockCfg, mockEthNetwork, stats)
	mockEthNetwork.CurrentNode = miner
	return miner
}

func createInMemObscuroNode(
	id int64,
	isGenesis bool,
	txHandler mgmtcontractlib.TxHandler,
	avgGossipPeriod time.Duration,
	avgBlockDuration time.Duration,
	avgNetworkLatency time.Duration,
	stats *stats.Stats,
	validateBlocks bool,
	genesisJSON []byte,
) *host.Node {
	hostConfig := host.Config{
		ID:                  common.BigToAddress(big.NewInt(id)),
		IsGenesis:           isGenesis,
		GossipRoundDuration: avgGossipPeriod,
		HasClientRPC:        false,
	}

	obscuroInMemNetwork := p2p2.NewMockP2P(avgBlockDuration, avgNetworkLatency)
	enclaveClient := enclave.NewEnclave(hostConfig.ID, true, txHandler, validateBlocks, genesisJSON, stats)
	node := host.NewHost(hostConfig, stats, obscuroInMemNetwork, nil, enclaveClient, txHandler)

	obscuroInMemNetwork.CurrentNode = &node
	return &node
}

func createSocketObscuroNode(id int64, isGenesis bool, avgGossipPeriod time.Duration, stats *stats.Stats, p2pAddr string, peerAddrs []string, enclaveAddr string, clientServerAddr string, txHandler mgmtcontractlib.TxHandler) *host.Node {
	hostConfig := host.Config{
		ID:                   common.BigToAddress(big.NewInt(id)),
		IsGenesis:            isGenesis,
		GossipRoundDuration:  avgGossipPeriod,
		HasClientRPC:         true,
		ClientRPCAddress:     &clientServerAddr,
		ClientRPCTimeoutSecs: ClientRPCTimeoutSecs,
		EnclaveRPCTimeout:    ClientRPCTimeoutSecs * time.Second,
		EnclaveRPCAddress:    &enclaveAddr,
	}

	enclaveClient := host.NewEnclaveRPCClient(hostConfig)
	nodeP2p := p2p.NewSocketP2PLayer(p2pAddr, peerAddrs, hostConfig.ID)
	node := host.NewHost(hostConfig, stats, nodeP2p, nil, enclaveClient, txHandler)

	return &node
}

func defaultMockEthNodeCfg(nrNodes int, avgBlockDuration time.Duration) ethereum_mock.MiningConfig {
	return ethereum_mock.MiningConfig{
		PowTime: func() time.Duration {
			// This formula might feel counter-intuitive, but it is a good approximation for Proof of Work.
			// It creates a uniform distribution up to nrMiners*avgDuration
			// Which means on average, every round, the winner (miner who gets the lowest nonce) will pick a number around "avgDuration"
			// while everyone else will have higher values.
			// Over a large number of rounds, the actual average block duration will be around the desired value, while the number of miners who get very close numbers will be limited.
			return obscurocommon.RndBtwTime(avgBlockDuration/time.Duration(nrNodes), avgBlockDuration*time.Duration(nrNodes))
		},
	}
}
