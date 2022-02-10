package simulation

import (
	"fmt"
	"github.com/google/uuid"
	"math/rand"
	"simulation/common"
	"simulation/ethereum-mock"
	"simulation/obscuro"
	"simulation/wallet-mock"
	"time"
)

// todo - introduce 2 parameters for nrNodes and random L1-L2 allocation
// todo - random add or remove l1 or l2 nodes - logic for catching up
func RunSimulation(nrWallets int, nrNodes int, simulationTime int, avgBlockDuration uint64, avgLatency uint64, gossipPeriod uint64) (L1NetworkCfg, L2NetworkCfg) {

	//todo - add observer nodes
	//todo read balance

	stats := NewStats(nrNodes, simulationTime, avgBlockDuration, avgLatency, gossipPeriod)

	l1Network := L1NetworkCfg{delay: func() uint64 {
		return common.RndBtw(uint64(avgLatency/10), uint64(2*avgLatency))
	}, Stats: &stats, interrupt: new(int32)}
	l1Cfg := ethereum_mock.MiningConfig{PowTime: func() uint64 {
		// This formula might feel counter-intuitive, but it is a good approximation for Proof of Work.
		// It creates a uniform distribution up to nrMiners*avgDuration
		// Which means on average, every round, the winner (miner who gets the lowest nonce) will pick a number around "avgDuration"
		// while everyone else will have higher values.
		// Over a large number of rounds, the actual average block duration will be around the desired value, while the number of miners who get very close numbers will be limited.
		return common.RndBtw(avgBlockDuration/uint64(nrNodes), uint64(nrNodes)*avgBlockDuration)
	}}

	l2Network := L2NetworkCfg{delay: func() uint64 {
		return common.RndBtw(uint64(avgLatency/10), uint64(2*avgLatency))
	}}
	l2Cfg := obscuro.AggregatorCfg{GossipRoundDuration: gossipPeriod}

	for i := 1; i <= nrNodes; i++ {
		// create a layer 2 node
		agg := obscuro.NewAgg(common.NodeId(i), l2Cfg, nil, &l2Network, &stats)
		l2Network.nodes = append(l2Network.nodes, &agg)

		// create a layer 1 node responsible with notifying the layer 2 node about blocks
		miner := ethereum_mock.NewMiner(common.NodeId(i), l1Cfg, &agg, &l1Network, &stats)
		l1Network.nodes = append(l1Network.nodes, &miner)
		agg.L1Node = &miner
	}

	common.Log(fmt.Sprintf("Genesis block: b_%s.", common.Str(common.GenesisBlock.Hash())))
	common.Log(fmt.Sprintf("Genesis rollup: r_%d.", obscuro.GenesisRollup.RootHash.ID()))

	l1Network.Start(common.Duration(avgBlockDuration / 4))

	// publish the genesis rollup before the l2 nodes are started
	tx, _ := obscuro.GenesisTx.Encode()
	l1Network.BroadcastTx(tx)

	l2Network.Start(common.Duration(avgBlockDuration / 4))

	// Create a bunch of users and inject transactions
	var wallets = make([]wallet_mock.Wallet, nrWallets)
	for i := 0; i < nrWallets; i++ {
		wallets[i] = wallet_mock.Wallet{Address: uuid.New().ID()}
	}

	timeInUs := simulationTime * 1000 * 1000
	go injectUserTxs(wallets, &l1Network, &l2Network, avgBlockDuration, timeInUs, &stats)

	// Wait for the simulation time
	time.Sleep(common.Duration(uint64(timeInUs)))

	fmt.Println("Stopping..")

	// stop L2 first and then L1
	go l2Network.Stop()
	go l1Network.Stop()

	time.Sleep(time.Second)

	return l1Network, l2Network
}

const INITIAL_BALANCE = 5000

func injectUserTxs(wallets []wallet_mock.Wallet, l1Network ethereum_mock.L1Network, l2Network obscuro.L2Network, avgBlockDuration uint64, simulationTime int, s *Stats) {
	// deposit some initial amount into every user
	initialiseWallets(wallets, l1Network, avgBlockDuration, s)

	// inject numbers of transactions proportional to the simulation time, such that they can be processed
	go injectRandomDeposits(wallets, l1Network, avgBlockDuration, simulationTime, s)
	go injectRandomWithdrawals(wallets, l2Network, avgBlockDuration, simulationTime, s)
	injectRandomTransfers(wallets, l2Network, avgBlockDuration, simulationTime, s)
}

func initialiseWallets(wallets []wallet_mock.Wallet, l1Network ethereum_mock.L1Network, avgBlockDuration uint64, s *Stats) {
	for _, u := range wallets {
		tx := deposit(u, INITIAL_BALANCE)
		t, _ := tx.Encode()
		l1Network.BroadcastTx(t)
		s.Deposit(INITIAL_BALANCE)
		time.Sleep(common.Duration(uint64(avgBlockDuration / 3)))
	}
}

func injectRandomTransfers(wallets []wallet_mock.Wallet, l2Network obscuro.L2Network, avgBlockDuration uint64, simulationTime int, s *Stats) {
	n := uint64(simulationTime) / avgBlockDuration
	i := uint64(0)
	for {
		if i == n {
			break
		}
		f := rndWallet(wallets).Address
		t := rndWallet(wallets).Address
		if f == t {
			continue
		}
		tx := obscuro.L2Tx{
			Id:     uuid.New(),
			TxType: obscuro.TransferTx,
			Amount: common.RndBtw(1, 500),
			From:   f,
			To:     t,
		}
		s.Transfer()
		encoded, err := tx.Encode()
		if err != nil {
			panic(err)
		}
		l2Network.BroadcastTx(encoded)
		time.Sleep(common.Duration(common.RndBtw(avgBlockDuration/4, avgBlockDuration)))
		i++
	}
}

func injectRandomDeposits(wallets []wallet_mock.Wallet, network ethereum_mock.L1Network, avgBlockDuration uint64, simulationTime int, s *Stats) {
	n := uint64(simulationTime) / (avgBlockDuration * 3)
	i := uint64(0)
	for {
		if i == n {
			break
		}
		v := common.RndBtw(1, 100)
		tx := deposit(rndWallet(wallets), v)
		t, _ := tx.Encode()
		network.BroadcastTx(t)
		s.Deposit(v)
		time.Sleep(common.Duration(common.RndBtw(avgBlockDuration, avgBlockDuration*2)))
		i++
	}
}

func injectRandomWithdrawals(wallets []wallet_mock.Wallet, network obscuro.L2Network, avgBlockDuration uint64, simulationTime int, s *Stats) {
	n := uint64(simulationTime) / (avgBlockDuration * 3)
	i := uint64(0)
	for {
		if i == n {
			break
		}
		v := common.RndBtw(1, 100)
		tx := withdrawal(rndWallet(wallets), v)
		t, _ := tx.Encode()
		network.BroadcastTx(t)
		s.Withdrawal(v)
		time.Sleep(common.Duration(common.RndBtw(avgBlockDuration, avgBlockDuration*2)))
		i++
	}
}

func withdrawal(wallet wallet_mock.Wallet, amount uint64) obscuro.L2Tx {
	return obscuro.L2Tx{
		Id:     uuid.New(),
		TxType: obscuro.WithdrawalTx,
		Amount: amount,
		From:   wallet.Address,
	}
}

func rndWallet(wallets []wallet_mock.Wallet) wallet_mock.Wallet {
	return wallets[rand.Intn(len(wallets))]
}

func deposit(wallet wallet_mock.Wallet, amount uint64) common.L1Tx {
	return common.L1Tx{
		Id:     uuid.New(),
		TxType: common.DepositTx,
		Amount: amount,
		Dest:   wallet.Address,
	}
}
