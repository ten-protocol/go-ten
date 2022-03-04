package simulation

import (
	"github.com/obscuronet/obscuro-playground/go/common"
	obscuro_node "github.com/obscuronet/obscuro-playground/go/obscuronode"
	ethereum_mock "github.com/obscuronet/obscuro-playground/integration/ethereummock"
)

// Network is a simulation network
type Network struct {
	l1NodeConfig *ethereum_mock.MiningConfig
	l1Network    *L1NetworkCfg
	l2NodeConfig *obscuro_node.AggregatorCfg
	l2Network    *L2NetworkCfg
}

// NewSimulationNetwork creates a simulation network
func NewSimulationNetwork(nrNodes int, l1NetworkCfg *L1NetworkCfg, l2NetworkCfg *L2NetworkCfg, avgBlockDuration uint64, gossipPeriod uint64, stats *Stats) *Network {
	l1NodeCfg := ethereum_mock.MiningConfig{
		PowTime: func() uint64 {
			// This formula might feel counter-intuitive, but it is a good approximation for Proof of Work.
			// It creates a uniform distribution up to nrMiners*avgDuration
			// Which means on average, every round, the winner (miner who gets the lowest nonce) will pick a number around "avgDuration"
			// while everyone else will have higher values.
			// Over a large number of rounds, the actual average block duration will be around the desired value, while the number of miners who get very close numbers will be limited.
			return common.RndBtw(avgBlockDuration/uint64(nrNodes), uint64(nrNodes)*avgBlockDuration)
		},
	}

	l2NodeCfg := obscuro_node.AggregatorCfg{GossipRoundDuration: gossipPeriod}

	for i := 1; i <= nrNodes; i++ {
		genesis := false
		if i == 1 {
			genesis = true
		}
		// create a layer 2 node
		agg := obscuro_node.NewAgg(common.NodeID(i), l2NodeCfg, nil, l2NetworkCfg, stats, genesis)
		l2NetworkCfg.nodes = append(l2NetworkCfg.nodes, &agg)

		// create a layer 1 node responsible with notifying the layer 2 node about blocks
		miner := ethereum_mock.NewMiner(common.NodeID(i), l1NodeCfg, &agg, l1NetworkCfg, stats)
		l1NetworkCfg.nodes = append(l1NetworkCfg.nodes, &miner)
		agg.L1Node = &miner
	}

	return &Network{
		l1NodeConfig: &l1NodeCfg,
		l1Network:    l1NetworkCfg,
		l2NodeConfig: &l2NodeCfg,
		l2Network:    l2NetworkCfg,
	}
}
