package network

import (
	"fmt"
	"math/big"
	"time"

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
	Localhost        = "Localhost"
	p2pStartPort     = 10000
	EnclaveStartPort = 11000
)

func createMockEthNode(id int64, nrNodes int, avgBlockDurationUSecs time.Duration, avgNetworkLatency time.Duration, stats *stats.Stats) *ethereum_mock.Node {
	mockEthNetwork := ethereum_mock.NewMockEthNetwork(avgBlockDurationUSecs, avgNetworkLatency, stats)
	ethereumMockCfg := defaultMockEthNodeCfg(nrNodes, avgBlockDurationUSecs)
	// create an in memory mock ethereum node responsible with notifying the layer 2 node about blocks
	miner := ethereum_mock.NewMiner(common.BigToAddress(big.NewInt(id)), ethereumMockCfg, mockEthNetwork, stats)
	mockEthNetwork.CurrentNode = &miner
	return &miner
}

func createInMemObscuroNode(id int64, genesis bool, avgGossipPeriod time.Duration, avgBlockDurationUSecs time.Duration, avgNetworkLatency time.Duration, stats *stats.Stats) *host.Node {
	obscuroInMemNetwork := p2p2.NewMockP2P(avgBlockDurationUSecs, avgNetworkLatency)

	obscuroNodeCfg := defaultObscuroNodeCfg(avgGossipPeriod)

	nodeID := common.BigToAddress(big.NewInt(id))
	enclaveClient := enclave.NewEnclave(nodeID, true, stats)

	// create an in memory obscuro node
	node := host.NewObscuroAggregator(nodeID, obscuroNodeCfg, nil, stats, genesis, enclaveClient, obscuroInMemNetwork)
	obscuroInMemNetwork.CurrentNode = &node
	return &node
}

func createSocketObscuroNode(id int64, genesis bool, avgGossipPeriod time.Duration, stats *stats.Stats, p2pAddr string, peerAddrs []string, enclavePort uint64) *host.Node {
	nodeID := common.BigToAddress(big.NewInt(id))

	// create an enclave client
	enclaveAddr := fmt.Sprintf("%s:%d", Localhost, enclavePort)
	enclaveClient := host.NewEnclaveRPCClient(enclaveAddr, host.ClientRPCTimeoutSecs*time.Second, nodeID)

	// create a socket obscuro node
	nodeP2p := p2p.NewSocketP2PLayer(p2pAddr, peerAddrs)
	obscuroNodeCfg := defaultObscuroNodeCfg(avgGossipPeriod)
	node := host.NewObscuroAggregator(nodeID, obscuroNodeCfg, nil, stats, genesis, enclaveClient, nodeP2p)

	return &node
}

func defaultObscuroNodeCfg(gossipPeriod time.Duration) host.AggregatorCfg {
	return host.AggregatorCfg{ClientRPCTimeoutSecs: host.ClientRPCTimeoutSecs, GossipRoundDuration: gossipPeriod}
}

func defaultMockEthNodeCfg(nrNodes int, avgBlockDuration time.Duration) ethereum_mock.MiningConfig {
	return ethereum_mock.MiningConfig{
		PowTime: func() uint64 {
			// This formula might feel counter-intuitive, but it is a good approximation for Proof of Work.
			// It creates a uniform distribution up to nrMiners*avgDuration
			// Which means on average, every round, the winner (miner who gets the lowest nonce) will pick a number around "avgDuration"
			// while everyone else will have higher values.
			// Over a large number of rounds, the actual average block duration will be around the desired value, while the number of miners who get very close numbers will be limited.
			return obscurocommon.RndBtw(uint64(avgBlockDuration.Nanoseconds()/int64(nrNodes)), uint64(int64(nrNodes)*avgBlockDuration.Nanoseconds()))
		},
	}
}
