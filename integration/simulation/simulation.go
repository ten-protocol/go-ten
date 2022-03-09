package simulation

import (
	"fmt"
	common2 "github.com/ethereum/go-ethereum/common"
	"math/big"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/obscuronet/obscuro-playground/go/common"
	obscuro_node "github.com/obscuronet/obscuro-playground/go/obscuronode"
	ethereum_mock "github.com/obscuronet/obscuro-playground/integration/ethereummock"
	wallet_mock "github.com/obscuronet/obscuro-playground/integration/walletmock"
)

// todo - introduce 2 parameters for nrNodes and random L1-L2 allocation
// todo - random add or remove l1 or l2 nodes - logic for catching up
func RunSimulation(
	nrWallets int,
	nrNodes int,
	simulationTime int,
	avgBlockDuration uint64,
	avgLatency uint64,
	gossipPeriod uint64,
) (L1NetworkCfg, L2NetworkCfg) {
	// todo - add observer nodes
	// todo read balance

	stats := NewStats(nrNodes, simulationTime, avgBlockDuration, avgLatency, gossipPeriod)

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
		agg := obscuro_node.NewAgg(common.NodeID(common2.BigToAddress(big.NewInt(int64(i)))), l2Cfg, nil, &l2Network, &stats, genesis)
		l2Network.nodes = append(l2Network.nodes, &agg)

		// create a layer 1 node responsible with notifying the layer 2 node about blocks
		miner := ethereum_mock.NewMiner(common.NodeID(common2.BigToAddress(big.NewInt(int64(i)))), l1Cfg, &agg, &l1Network, &stats)
		l1Network.nodes = append(l1Network.nodes, &miner)
		agg.L1Node = &miner
	}

	common.Log(fmt.Sprintf("Genesis block: b_%s.", common.Str(common.GenesisBlock.Hash())))

	l1Network.Start(common.Duration(avgBlockDuration / 4))
	l2Network.Start(common.Duration(avgBlockDuration / 4))

	// Create a bunch of users and inject transactions
	wallets := make([]wallet_mock.Wallet, nrWallets)
	for i := 0; i < nrWallets; i++ {
		wallets[i] = wallet_mock.New()
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

const INITIAL_BALANCE = 5000 // nolint:revive,stylecheck

func injectUserTxs(wallets []wallet_mock.Wallet, l1Network ethereum_mock.L1Network, l2Network obscuro_node.L2Network, avgBlockDuration uint64, simulationTime int, s *Stats) {
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
		t, _ := common.EncodeTx(tx)
		l1Network.BroadcastTx(t)
		s.Deposit(INITIAL_BALANCE)
		time.Sleep(common.Duration(avgBlockDuration / 3))
	}
}

func injectRandomTransfers(wallets []wallet_mock.Wallet, l2Network obscuro_node.L2Network, avgBlockDuration uint64, simulationTime int, s *Stats) {
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
		encryptedTx := wallet_mock.NewEncryptedL2Transfer(f, t, common.RndBtw(1, 500))
		s.Transfer()
		l2Network.BroadcastTx(encryptedTx)
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
		t, _ := common.EncodeTx(tx)
		network.BroadcastTx(t)
		s.Deposit(v)
		time.Sleep(common.Duration(common.RndBtw(avgBlockDuration, avgBlockDuration*2)))
		i++
	}
}

func injectRandomWithdrawals(wallets []wallet_mock.Wallet, network obscuro_node.L2Network, avgBlockDuration uint64, simulationTime int, s *Stats) {
	n := uint64(simulationTime) / (avgBlockDuration * 3)
	i := uint64(0)
	for {
		if i == n {
			break
		}
		v := common.RndBtw(1, 100)
		encryptedTx := wallet_mock.NewEncryptedL2Withdrawal(rndWallet(wallets).Address, v)
		network.BroadcastTx(encryptedTx)
		s.Withdrawal(v)
		time.Sleep(common.Duration(common.RndBtw(avgBlockDuration, avgBlockDuration*2)))
		i++
	}
}

func rndWallet(wallets []wallet_mock.Wallet) wallet_mock.Wallet {
	return wallets[rand.Intn(len(wallets))] //nolint:gosec
}

func deposit(wallet wallet_mock.Wallet, amount uint64) *common.L1Tx {
	data := common.L1TxData{
		TxType: common.DepositTx,
		Amount: amount,
		Dest:   wallet.Address,
	}
	return common.NewL1Tx(data)
}
