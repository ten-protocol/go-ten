package simulation

import (
	"fmt"
	"math/big"
	"time"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/host/p2p"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave"

	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/host"

	"github.com/obscuronet/obscuro-playground/go/log"

	ethereum_mock "github.com/obscuronet/obscuro-playground/integration/ethereummock"
)

const (
	INITIAL_BALANCE         = 5000        // nolint:revive,stylecheck
	P2P_START_PORT          = 10000       // nolint:revive,stylecheck
	ENCLAVE_CONN_START_PORT = 11000       // nolint:revive,stylecheck
	LOCALHOST               = "localhost" // nolint:revive,stylecheck
)

// Simulation represents the data which to set up and run a simulated network
type Simulation struct {
	l1NodeConfig     *ethereum_mock.MiningConfig
	l1Network        *L1NetworkCfg
	l2NodeConfig     *host.AggregatorCfg
	l2Network        *L2NetworkCfg
	avgBlockDuration uint64
}

// NewSimulation defines a new simulation network
func NewSimulation(
	nrNodes int,
	l1NetworkCfg *L1NetworkCfg,
	l2NetworkCfg *L2NetworkCfg,
	avgBlockDuration uint64,
	gossipPeriod uint64,
	localEnclave bool,
	stats *Stats,
) *Simulation {
	l1NodeCfg := ethereum_mock.MiningConfig{
		PowTime: func() uint64 {
			// This formula might feel counter-intuitive, but it is a good approximation for Proof of Work.
			// It creates a uniform distribution up to nrMiners*avgDuration
			// Which means on average, every round, the winner (miner who gets the lowest nonce) will pick a number around "avgDuration"
			// while everyone else will have higher values.
			// Over a large number of rounds, the actual average block duration will be around the desired value, while the number of miners who get very close numbers will be limited.
			return obscurocommon.RndBtw(avgBlockDuration/uint64(nrNodes), uint64(nrNodes)*avgBlockDuration)
		},
	}

	l2NodeCfg := host.AggregatorCfg{ClientRPCTimeoutSecs: host.ClientRPCTimeoutSecs, GossipRoundDuration: gossipPeriod}

	// We generate the P2P addresses for each node on the network.
	for i := 1; i <= nrNodes; i++ {
		l2NetworkCfg.nodeAddresses = append(l2NetworkCfg.nodeAddresses, fmt.Sprintf("%s:%d", LOCALHOST, P2P_START_PORT+i))
	}

	for i := 1; i <= nrNodes; i++ {
		genesis := false
		if i == 1 {
			genesis = true
		}

		// create an enclave server
		nodeID := common.BigToAddress(big.NewInt(int64(i)))
		var enclaveClient nodecommon.Enclave
		if localEnclave {
			enclaveClient = enclave.NewEnclave(nodeID, true, stats)
		} else {
			port := uint64(ENCLAVE_CONN_START_PORT + i)
			timeout := time.Duration(l2NodeCfg.ClientRPCTimeoutSecs) * time.Second
			err := enclave.StartServer(port, nodeID, stats)
			if err != nil {
				panic(fmt.Sprintf("failed to create enclave server: %v", err))
			}
			enclaveClient = host.NewEnclaveRPCClient(fmt.Sprintf("%s:%d", LOCALHOST, port), timeout)
		}

		// create a layer 2 node
		aggP2P := p2p.NewP2P(l2NetworkCfg.nodeAddresses[i-1], l2NetworkCfg.nodeAddresses)
		agg := host.NewAgg(nodeID, l2NodeCfg, nil, stats, genesis, enclaveClient, aggP2P)
		l2NetworkCfg.nodes = append(l2NetworkCfg.nodes, &agg)

		// create a layer 1 node responsible with notifying the layer 2 node about blocks
		miner := ethereum_mock.NewMiner(common.BigToAddress(big.NewInt(int64(i))), l1NodeCfg, &agg, l1NetworkCfg, stats)
		l1NetworkCfg.nodes = append(l1NetworkCfg.nodes, &miner)
		agg.L1Node = &miner
	}

	log.Log(fmt.Sprintf("Genesis block: b_%d.", obscurocommon.ShortHash(obscurocommon.GenesisBlock.Hash())))

	return &Simulation{
		l1NodeConfig:     &l1NodeCfg,
		l1Network:        l1NetworkCfg,
		l2NodeConfig:     &l2NodeCfg,
		l2Network:        l2NetworkCfg,
		avgBlockDuration: avgBlockDuration,
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

	log.Log(fmt.Sprintf("Genesis block: b_%d.", obscurocommon.ShortHash(obscurocommon.GenesisBlock.Hash())))

	// The sequence of starting the nodes is important to catch various edge cases.
	// Here we first start the mock layer 1 nodes, with a pause between them of a fraction of a block duration.
	// The reason is to make sure that they catch up correctly.
	// Then we pause for a while, to give the L1 network enough time to create a number of blocks, which will have to be ingested by the Obscuro nodes
	// Then, we begin the starting sequence of the Obscuro nodes, again with a delay between them, to test that they are able to cach up correctly.
	// Note: Other simulations might test variations of this pattern.
	s.l1Network.Start(time.Duration(s.avgBlockDuration / 8))
	time.Sleep(time.Duration(s.avgBlockDuration * 20))
	s.l2Network.Start(time.Duration(s.avgBlockDuration / 3))

	// time in micro seconds to run the simulation
	timeInUs := simulationTime * 1000 * 1000

	timer := time.Now()
	go txManager.Start(timeInUs)

	// Wait for the simulation time
	time.Sleep(obscurocommon.Duration(uint64(timeInUs)))

	fmt.Printf("Ran simulation for %f secs, configured to run for: %s ... \n", time.Since(timer).Seconds(), obscurocommon.Duration(uint64(timeInUs)))
	time.Sleep(time.Second)
}

// Stop closes down the L2 and L1 networks.
func (s *Simulation) Stop() {
	// stop L2 first and then L1
	go s.l2Network.Stop()
	go s.l1Network.Stop()
}
