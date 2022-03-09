package simulation

import (
	"fmt"
	"math/rand"
	"time"

	obscuro_node "github.com/obscuronet/obscuro-playground/go/obscuronode"
	ethereum_mock "github.com/obscuronet/obscuro-playground/integration/ethereummock"

	"github.com/google/uuid"
	"github.com/obscuronet/obscuro-playground/go/common"
	"github.com/obscuronet/obscuro-playground/go/log"

	enclave2 "github.com/obscuronet/obscuro-playground/go/obscuronode/enclave"
	wallet_mock "github.com/obscuronet/obscuro-playground/integration/walletmock"
)

const (
	INITIAL_BALANCE      = 5000 // nolint:revive,stylecheck
	NODE_BOOTUP_DELAY_MS = 100  // nolint:revive,stylecheck
)

// Simulation represents the data which to set up and run a simulated network
type Simulation struct {
	l1NodeConfig *ethereum_mock.MiningConfig
	l1Network    *L1NetworkCfg
	l2NodeConfig *obscuro_node.AggregatorCfg
	l2Network    *L2NetworkCfg
}

// NewSimulation defines a new simulation network
func NewSimulation(nrNodes int, l1NetworkCfg *L1NetworkCfg, l2NetworkCfg *L2NetworkCfg, avgBlockDuration uint64, gossipPeriod uint64, stats *Stats) *Simulation {
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

	return &Simulation{
		l1NodeConfig: &l1NodeCfg,
		l1Network:    l1NetworkCfg,
		l2NodeConfig: &l2NodeCfg,
		l2Network:    l2NetworkCfg,
	}
}

// RunSimulation executes the simulation given all the params
// todo - introduce 2 parameters for nrNodes and random L1-L2 allocation
// todo - random add or remove l1 or l2 nodes - logic for catching up
func (s *Simulation) Start(
	txManager *TransactionManager,
	simulationTime int,
) {
	// todo - add observer nodes
	// todo read balance

	log.Log(fmt.Sprintf("Genesis block: b_%s.", common.Str(common.GenesisBlock.Hash())))

	s.l1Network.Start()
	s.l2Network.Start()

	timeInUs := simulationTime * 1000 * 1000

	go txManager.Start(timeInUs)

	// Wait for the simulation time
	time.Sleep(common.Duration(uint64(timeInUs)))

	fmt.Printf("Stopping simulation after running it for: %s ... \n", common.Duration(uint64(timeInUs)))

	// stop L2 first and then L1
	go s.l2Network.Stop()
	go s.l1Network.Stop()

	time.Sleep(time.Second)
}

func withdrawal(wallet wallet_mock.Wallet, amount uint64) enclave2.L2Tx {
	return enclave2.L2Tx{
		ID:     uuid.New(),
		TxType: enclave2.WithdrawalTx,
		Amount: amount,
		From:   wallet.Address,
	}
}

func rndWallet(wallets []wallet_mock.Wallet) wallet_mock.Wallet {
	return wallets[rand.Intn(len(wallets))] //nolint:gosec
}

func deposit(wallet wallet_mock.Wallet, amount uint64) common.L1Tx {
	return common.L1Tx{
		ID:     uuid.New(),
		TxType: common.DepositTx,
		Amount: amount,
		Dest:   wallet.Address,
	}
}
