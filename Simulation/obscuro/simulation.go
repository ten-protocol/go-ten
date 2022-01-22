package obscuro

import (
	"fmt"
	"github.com/google/uuid"
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
	l2Head             *Rollup
	maxRollupsPerBlock int
	nrEmptyBlocks      int

	totalL2Txs int
	noL1Reorgs map[NodeId]int
	noL2Reorgs map[NodeId]int
	// todo - actual avg block Duration

	totalDepositedAmount   int
	nrTransferTransactions int
}

func RunSimulation(nrUsers int, nrMiners int, simulationTime int, avgBlockDuration int, avgLatency int, gossipPeriod int) NetworkCfg {

	var stats = Stats{
		nrMiners:         nrMiners,
		simulationTime:   simulationTime,
		avgBlockDuration: avgBlockDuration,
		avgLatency:       avgLatency,
		gossipPeriod:     gossipPeriod,
		noL1Reorgs:       map[NodeId]int{},
		noL2Reorgs:       map[NodeId]int{},
	}

	var network = NetworkCfg{delay: func() int {
		return RndBtw(avgLatency/10, 2*avgLatency)
	}, Stats: &stats}

	l1Config := L1MiningConfig{powTime: func() int {
		return RndBtw(avgBlockDuration/nrMiners, nrMiners*avgBlockDuration)
	}}

	l2Cfg := L2Cfg{gossipPeriodMs: gossipPeriod}

	for i := 1; i <= nrMiners; i++ {
		agg := NewAgg(NodeId(i), l2Cfg, nil, &network)
		miner := NewMiner(NodeId(i), l1Config, &agg, &network)
		stats.noL1Reorgs[NodeId(i)] = 0
		agg.l1 = &miner
		network.allAgg = append(network.allAgg, agg)
		network.allMiners = append(network.allMiners, miner)
	}

	log(fmt.Sprintf("Genesis block: b_%d.", GenesisBlock.rootHash.ID()))
	log(fmt.Sprintf("Genesis rollup: r_%d.", GenesisRollup.rootHash.ID()))

	for _, m := range network.allMiners {
		//fmt.Printf("Starting Miner M%d....\n", m.id)
		t := m
		go t.Start()
		defer t.Stop()
		go t.aggregator.Start()
		defer t.aggregator.Stop()
		// don't start everything at once
		time.Sleep(Duration(avgBlockDuration / 2))
	}

	var users = make([]Wallet, 0)
	for i := 1; i <= nrUsers; i++ {
		users = append(users, Wallet{address: uuid.New()})
	}

	go injectUserTxs(users, &network, avgBlockDuration)

	time.Sleep(Duration(simulationTime * 1000 * 1000))

	return network
}

const INITIAL_BALANCE = 5000

//var total = 0
//var nrTransf = 0

func injectUserTxs(users []Wallet, network *NetworkCfg, avgBlockDuration int) {
	// deposit some initial amount into every user
	for _, u := range users {
		tx := deposit(u, INITIAL_BALANCE)
		//total += INITIAL_BALANCE
		network.broadcastL1Tx(&tx)
		time.Sleep(Duration(avgBlockDuration / 3))
	}

	go injectDeposits(users, network, avgBlockDuration)

	// generate random L2 transfers
	for {
		f := rndUser(users).address
		t := rndUser(users).address
		if f == t {
			continue
		}
		tx := L2Tx{
			id:     uuid.New(),
			txType: TransferTx,
			amount: RndBtw(1, 500),
			from:   f,
			dest:   t,
		}
		network.broadcastL2Tx(&tx)
		//nrTransf++
		time.Sleep(Duration(RndBtw(avgBlockDuration/3, avgBlockDuration)))
	}
}

func injectDeposits(users []Wallet, network *NetworkCfg, avgBlockDuration int) {
	i := 0
	for {
		if i == 1000 {
			break
		}
		v := RndBtw(1, 100)
		//v := INITIAL_BALANCE
		//total += v
		tx := deposit(rndUser(users), v)
		network.broadcastL1Tx(&tx)
		time.Sleep(Duration(RndBtw(avgBlockDuration, avgBlockDuration*2)))
		i++
	}
}

func rndUser(users []Wallet) Wallet {
	return users[RndBtw(0, len(users))]
}

func deposit(u Wallet, amount int) L1Tx {
	return L1Tx{
		id:     uuid.New(),
		txType: DepositTx,
		amount: amount,
		dest:   u.address,
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
