package node

import (
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	ethereum_mock "github.com/obscuronet/obscuro-playground/integration/ethereummock"
)

func defaultMockEthNodeCfg(nrNodes int, avgBlockDuration uint64) ethereum_mock.MiningConfig {
	return ethereum_mock.MiningConfig{
		PowTime: func() uint64 {
			// This formula might feel counter-intuitive, but it is a good approximation for Proof of Work.
			// It creates a uniform distribution up to nrMiners*avgDuration
			// Which means on average, every round, the winner (miner who gets the lowest nonce) will pick a number around "avgDuration"
			// while everyone else will have higher values.
			// Over a large number of rounds, the actual average block duration will be around the desired value, while the number of miners who get very close numbers will be limited.
			return obscurocommon.RndBtw(avgBlockDuration/uint64(nrNodes), uint64(nrNodes)*avgBlockDuration)
		},
	}
}
