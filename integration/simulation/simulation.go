package simulation

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/obscuronet/obscuro-playground/go/common"
	obscuro_node "github.com/obscuronet/obscuro-playground/go/obscuronode"
	enclave2 "github.com/obscuronet/obscuro-playground/go/obscuronode/enclave"
	ethereum_mock "github.com/obscuronet/obscuro-playground/integration/ethereummock"
	wallet_mock "github.com/obscuronet/obscuro-playground/integration/walletmock"
)

// todo - introduce 2 parameters for nrNodes and random L1-L2 allocation
// todo - random add or remove l1 or l2 nodes - logic for catching up
func RunSimulation(
	txManager *TransactionManager,
	nrNodes int,
	simulationTime int,
	avgBlockDuration uint64,
	avgLatency uint64,
	gossipPeriod uint64,
	stats Stats,
) (L1NetworkCfg, L2NetworkCfg) {
	// todo - add observer nodes
	// todo read balance

	l1Network := L1NetworkCfg{delay: func() uint64 {
		return common.RndBtw(avgLatency/10, 2*avgLatency)
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
		return common.RndBtw(avgLatency/10, 2*avgLatency)
	}}
	l2Cfg := obscuro_node.AggregatorCfg{GossipRoundDuration: gossipPeriod}

	for i := 1; i <= nrNodes; i++ {
		genesis := false
		if i == 1 {
			genesis = true
		}
		// create a layer 2 node
		agg := obscuro_node.NewAgg(common.NodeID(i), l2Cfg, nil, &l2Network, &stats, genesis)
		l2Network.nodes = append(l2Network.nodes, &agg)

		// create a layer 1 node responsible with notifying the layer 2 node about blocks
		miner := ethereum_mock.NewMiner(common.NodeID(i), l1Cfg, &agg, &l1Network, &stats)
		l1Network.nodes = append(l1Network.nodes, &miner)
		agg.L1Node = &miner
	}

	common.Log(fmt.Sprintf("Genesis block: b_%s.", common.Str(common.GenesisBlock.Hash())))

	l1Network.Start(common.Duration(avgBlockDuration / 4))
	l2Network.Start(common.Duration(avgBlockDuration / 4))

	timeInUs := simulationTime * 1000 * 1000

	go txManager.Start(&l1Network, &l2Network, avgBlockDuration, timeInUs, &stats)

	// Wait for the simulation time
	time.Sleep(common.Duration(uint64(timeInUs)))

	fmt.Printf("Stopping simulation after running it for: %s ... \n", common.Duration(uint64(timeInUs)))

	// stop L2 first and then L1
	go l2Network.Stop()
	go l1Network.Stop()

	time.Sleep(time.Second)

	return l1Network, l2Network
}

const INITIAL_BALANCE = 5000 // nolint:revive,stylecheck

func injectRandomTransfers(wallets []wallet_mock.Wallet, l2Network obscuro_node.L2Network, avgBlockDuration uint64, simulationTime int, s *Stats, trackTx func(tx enclave2.L2Tx)) {
	n := uint64(simulationTime) / (avgBlockDuration * 3)
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
		tx := enclave2.L2Tx{
			ID:     uuid.New(),
			TxType: enclave2.TransferTx,
			Amount: common.RndBtw(1, 500),
			From:   f,
			To:     t,
		}
		s.Transfer()
		encoded := enclave2.EncryptTx(tx)
		l2Network.BroadcastTx(encoded)
		go trackTx(tx)
		time.Sleep(common.Duration(common.RndBtw(avgBlockDuration/4, avgBlockDuration)))
		i++
	}
}

func injectRandomDeposits(wallets []wallet_mock.Wallet, network ethereum_mock.L1Network, avgBlockDuration uint64, simulationTime int, s *Stats, trackTx func(tx common.L1Tx)) {
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
		go trackTx(tx)
		time.Sleep(common.Duration(common.RndBtw(avgBlockDuration, avgBlockDuration*2)))
		i++
	}
}

func injectRandomWithdrawals(wallets []wallet_mock.Wallet, network obscuro_node.L2Network, avgBlockDuration uint64, simulationTime int, s *Stats, trackTx func(tx enclave2.L2Tx)) {
	n := uint64(simulationTime) / (avgBlockDuration * 3)
	i := uint64(0)
	for {
		if i == n {
			break
		}
		v := common.RndBtw(1, 100)
		tx := withdrawal(rndWallet(wallets), v)
		t := enclave2.EncryptTx(tx)
		network.BroadcastTx(t)
		s.Withdrawal(v)
		go trackTx(tx)
		time.Sleep(common.Duration(common.RndBtw(avgBlockDuration, avgBlockDuration*2)))
		i++
	}
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
