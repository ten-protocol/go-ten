package simulation

import (
	"fmt"
	"github.com/google/uuid"
	"simulation/common"
	"simulation/ethereum-mock"
	"simulation/obscuro"
	"simulation/wallet-mock"
	"time"
)

func RunSimulation(nrUsers int, nrNodes int, simulationTime int, avgBlockDuration int, avgLatency int, gossipPeriod int) (L1NetworkCfg, L2NetworkCfg) {

	stats := NewStats(nrNodes, simulationTime, avgBlockDuration, avgLatency, gossipPeriod)

	l1_network := L1NetworkCfg{delay: func() int {
		return common.RndBtw(avgLatency/10, 2*avgLatency)
	}, Stats: &stats}
	l1_cfg := ethereum_mock.L1MiningConfig{PowTime: func() int {
		return common.RndBtw(avgBlockDuration/nrNodes, nrNodes*avgBlockDuration)
	}}

	l2_network := L2NetworkCfg{delay: func() int {
		return common.RndBtw(avgLatency/10, 2*avgLatency)
	}, Stats: &stats}
	l2_cfg := obscuro.L2Cfg{GossipPeriod: gossipPeriod}

	for i := 1; i <= nrNodes; i++ {
		// create a layer 2 node
		agg := obscuro.NewAgg(common.NodeId(i), l2_cfg, nil, l1_network, l2_network)
		l2_network.nodes = append(l2_network.nodes, agg)

		// create a layer 1 node responsible to notify the layer 2 node about blocks
		miner := ethereum_mock.NewMiner(common.NodeId(i), l1_cfg, &agg, l1_network)
		l1_network.nodes = append(l1_network.nodes, miner)
		agg.L1 = &miner

		// initialize stats
		stats.noL1Reorgs[common.NodeId(i)] = 0
	}

	common.Log(fmt.Sprintf("Genesis block: b_%d.", common.GenesisBlock.RootHash().ID()))
	common.Log(fmt.Sprintf("Genesis rollup: r_%d.", common.GenesisRollup.RootHash().ID()))

	l1_network.Start(common.Duration(avgBlockDuration / 4))
	defer l1_network.Stop()
	l2_network.Start(common.Duration(avgBlockDuration / 4))
	defer l2_network.Stop()

	// Create a bunch of users and inject transactions
	var users = make([]wallet_mock.Wallet, 0)
	for i := 1; i <= nrUsers; i++ {
		users = append(users, wallet_mock.Wallet{Address: uuid.New()})
	}

	go injectUserTxs(users, &l1_network, &l2_network, avgBlockDuration)

	// Wait for the simulation time
	time.Sleep(common.Duration(simulationTime * 1000 * 1000))

	return l1_network, l2_network
}

const INITIAL_BALANCE = 5000

//var total = 0
//var nrTransf = 0

func injectUserTxs(users []wallet_mock.Wallet, l1Network common.L1Network, l2Network common.L2Network, avgBlockDuration int) {
	// deposit some initial amount into every user
	for _, u := range users {
		tx := deposit(u, INITIAL_BALANCE)
		//total += INITIAL_BALANCE
		l1Network.BroadcastL1Tx(tx)
		time.Sleep(common.Duration(avgBlockDuration / 3))
	}

	go injectDeposits(users, l1Network, avgBlockDuration)

	// generate random L2 transfers
	for {
		f := rndUser(users).Address
		t := rndUser(users).Address
		if f == t {
			continue
		}
		tx := common.L2Tx{
			Id:     uuid.New(),
			TxType: common.TransferTx,
			Amount: common.RndBtw(1, 500),
			From:   f,
			Dest:   t,
		}
		l2Network.BroadcastL2Tx(tx)
		//nrTransf++
		time.Sleep(common.Duration(common.RndBtw(avgBlockDuration/3, avgBlockDuration)))
	}
}

func injectDeposits(users []wallet_mock.Wallet, network common.L1Network, avgBlockDuration int) {
	i := 0
	for {
		if i == 1000 {
			break
		}
		v := common.RndBtw(1, 100)
		//v := INITIAL_BALANCE
		//total += v
		tx := deposit(rndUser(users), v)
		network.BroadcastL1Tx(tx)
		time.Sleep(common.Duration(common.RndBtw(avgBlockDuration, avgBlockDuration*2)))
		i++
	}
}

func rndUser(users []wallet_mock.Wallet) wallet_mock.Wallet {
	return users[common.RndBtw(0, len(users))]
}

func deposit(u wallet_mock.Wallet, amount int) common.L1Tx {
	return common.L1Tx{
		Id:     uuid.New(),
		TxType: common.DepositTx,
		Amount: amount,
		Dest:   u.Address,
	}
}
