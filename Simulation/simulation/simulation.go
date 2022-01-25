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

type Stats struct {
	nrMiners         int
	simulationTime   int
	avgBlockDuration int
	avgLatency       int
	gossipPeriod     int

	l1Height int
	totalL1  int

	l2Height           int
	totalL2            int
	l2Head             *common.Rollup
	maxRollupsPerBlock int
	nrEmptyBlocks      int

	totalL2Txs int
	noL1Reorgs map[common.NodeId]int
	noL2Reorgs map[common.NodeId]int
	// todo - actual avg block Duration

	totalDepositedAmount   int
	nrTransferTransactions int
}

func (s *Stats) Reorg(id common.NodeId) {
	statsMu.Lock()
	s.noL1Reorgs[id]++
	statsMu.Unlock()
}

func RunSimulation(nrUsers int, nrMiners int, simulationTime int, avgBlockDuration int, avgLatency int, gossipPeriod int) NetworkCfg {

	var stats = Stats{
		nrMiners:         nrMiners,
		simulationTime:   simulationTime,
		avgBlockDuration: avgBlockDuration,
		avgLatency:       avgLatency,
		gossipPeriod:     gossipPeriod,
		noL1Reorgs:       map[common.NodeId]int{},
		noL2Reorgs:       map[common.NodeId]int{},
	}

	var network = NetworkCfg{delay: func() int {
		return common.RndBtw(avgLatency/10, 2*avgLatency)
	}, Stats: &stats}

	l1Config := ethereum_mock.L1MiningConfig{PowTime: func() int {
		return common.RndBtw(avgBlockDuration/nrMiners, nrMiners*avgBlockDuration)
	}}

	l2Cfg := obscuro.L2Cfg{GossipPeriod: gossipPeriod}

	for i := 1; i <= nrMiners; i++ {
		agg := obscuro.NewAgg(common.NodeId(i), l2Cfg, nil, network, network)
		miner := ethereum_mock.NewMiner(common.NodeId(i), l1Config, &agg, &network, &stats)
		stats.noL1Reorgs[common.NodeId(i)] = 0
		agg.L1 = &miner
		network.allAgg = append(network.allAgg, agg)
		network.allMiners = append(network.allMiners, miner)
	}

	common.Log(fmt.Sprintf("Genesis block: b_%d.", common.GenesisBlock.RootHash().ID()))
	common.Log(fmt.Sprintf("Genesis rollup: r_%d.", common.GenesisRollup.RootHash().ID()))

	for _, m := range network.allMiners {
		//fmt.Printf("Starting Miner M%d....\n", m.Id)
		t := m
		go t.Start()
		defer t.Stop()
		// don't start everything at once
		time.Sleep(common.Duration(avgBlockDuration / 4))
	}

	for _, a := range network.allAgg {
		//fmt.Printf("Starting Miner M%d....\n", m.Id)
		t := a
		go t.Start()
		defer t.Stop()
		// don't start everything at once
		time.Sleep(common.Duration(avgBlockDuration / 4))
	}

	var users = make([]wallet_mock.Wallet, 0)
	for i := 1; i <= nrUsers; i++ {
		users = append(users, wallet_mock.Wallet{Address: uuid.New()})
	}

	go injectUserTxs(users, &network, avgBlockDuration)

	time.Sleep(common.Duration(simulationTime * 1000 * 1000))

	return network
}

const INITIAL_BALANCE = 5000

//var total = 0
//var nrTransf = 0

func injectUserTxs(users []wallet_mock.Wallet, network *NetworkCfg, avgBlockDuration int) {
	// deposit some initial amount into every user
	for _, u := range users {
		tx := deposit(u, INITIAL_BALANCE)
		//total += INITIAL_BALANCE
		network.BroadcastL1Tx(tx)
		time.Sleep(common.Duration(avgBlockDuration / 3))
	}

	go injectDeposits(users, network, avgBlockDuration)

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
		network.BroadcastL2Tx(tx)
		//nrTransf++
		time.Sleep(common.Duration(common.RndBtw(avgBlockDuration/3, avgBlockDuration)))
	}
}

func injectDeposits(users []wallet_mock.Wallet, network *NetworkCfg, avgBlockDuration int) {
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

func findDups(list []uuid.UUID) map[uuid.UUID]int {

	duplicate_frequency := make(map[uuid.UUID]int)

	for _, item := range list {
		// check if the item/element exist in the duplicate_frequency map

		_, exist := duplicate_frequency[item]

		if exist {
			duplicate_frequency[item] += 1 // increase counter by 1 if already in the map
		} else {
			duplicate_frequency[item] = 1 // else start counting from 1
		}
	}
	dups := make(map[uuid.UUID]int)
	for u, i := range duplicate_frequency {
		if i > 1 {
			dups[u] = i
		}
	}

	return dups
}
