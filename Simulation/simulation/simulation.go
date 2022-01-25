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

	l1Network := L1NetworkCfg{delay: func() int {
		return common.RndBtw(avgLatency/10, 2*avgLatency)
	}, Stats: &stats}
	l1Cfg := ethereum_mock.MiningConfig{PowTime: func() int {
		return common.RndBtw(avgBlockDuration/nrNodes, nrNodes*avgBlockDuration)
	}}

	l2Network := L2NetworkCfg{delay: func() int {
		return common.RndBtw(avgLatency/10, 2*avgLatency)
	}}
	l2Cfg := obscuro.AggregatorCfg{GossipPeriod: gossipPeriod}

	for i := 1; i <= nrNodes; i++ {
		// create a layer 2 node
		agg := obscuro.NewAgg(common.NodeId(i), l2Cfg, nil, &l1Network, &l2Network, &stats)
		l2Network.nodes = append(l2Network.nodes, &agg)

		// create a layer 1 node responsible with notifying the layer 2 node about blocks
		miner := ethereum_mock.NewMiner(common.NodeId(i), l1Cfg, &agg, &l1Network, &stats)
		l1Network.nodes = append(l1Network.nodes, &miner)
		agg.L1Node = &miner

		// initialize stats
		stats.noL1Reorgs[common.NodeId(i)] = 0
		stats.noL2Reorgs[common.NodeId(i)] = 0
	}

	common.Log(fmt.Sprintf("Genesis block: b_%d.", common.GenesisBlock.RootHash().ID()))
	common.Log(fmt.Sprintf("Genesis rollup: r_%d.", common.GenesisRollup.RootHash().ID()))

	l1Network.Start(common.Duration(avgBlockDuration / 4))
	defer l1Network.Stop()
	l2Network.Start(common.Duration(avgBlockDuration / 4))
	defer l2Network.Stop()

	// Create a bunch of users and inject transactions
	var users = make([]wallet_mock.Wallet, 0)
	for i := 1; i <= nrUsers; i++ {
		users = append(users, wallet_mock.Wallet{Address: uuid.New()})
	}

	go injectUserTxs(users, &l1Network, &l2Network, avgBlockDuration, simulationTime)

	// Wait for the simulation time
	time.Sleep(common.Duration(simulationTime * 1000 * 1000))

	return l1Network, l2Network
}

const INITIAL_BALANCE = 5000

func injectUserTxs(users []wallet_mock.Wallet, l1Network common.L1Network, l2Network common.L2Network, avgBlockDuration int, simulationTime int) {
	// deposit some initial amount into every user
	initialiseWallets(users, l1Network, avgBlockDuration)

	go injectRandomDeposits(users, l1Network, avgBlockDuration, simulationTime)

	// generate random L2 transfers
	injectRandomTransfers(users, l2Network, avgBlockDuration, simulationTime)
}

func initialiseWallets(users []wallet_mock.Wallet, l1Network common.L1Network, avgBlockDuration int) {
	for _, u := range users {
		tx := deposit(u, INITIAL_BALANCE)
		l1Network.BroadcastL1Tx(tx)
		time.Sleep(common.Duration(avgBlockDuration / 3))
	}
}

func injectRandomTransfers(users []wallet_mock.Wallet, l2Network common.L2Network, avgBlockDuration int, simulationTime int) {
	n := simulationTime * 2 / avgBlockDuration
	i := 0
	for {
		if i == n {
			break
		}
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
		time.Sleep(common.Duration(common.RndBtw(avgBlockDuration/4, avgBlockDuration)))
		i++
	}
}

func injectRandomDeposits(users []wallet_mock.Wallet, network common.L1Network, avgBlockDuration int, simulationTime int) {
	n := simulationTime / (avgBlockDuration * 2)
	i := 0
	for {
		if i == n {
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
