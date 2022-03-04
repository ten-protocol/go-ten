package simulation

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/obscuronet/obscuro-playground/go/common"
	"github.com/obscuronet/obscuro-playground/go/log"

	enclave2 "github.com/obscuronet/obscuro-playground/go/obscuronode/enclave"
	wallet_mock "github.com/obscuronet/obscuro-playground/integration/walletmock"
)

const INITIAL_BALANCE = 5000 // nolint:revive,stylecheck

// RunSimulation executes the simulation given all the params
// todo - introduce 2 parameters for nrNodes and random L1-L2 allocation
// todo - random add or remove l1 or l2 nodes - logic for catching up
func RunSimulation(
	txManager *TransactionManager,
	network *Network,
	simulationTime int,
) {
	// todo - add observer nodes
	// todo read balance

	log.Log(fmt.Sprintf("Genesis block: b_%s.", common.Str(common.GenesisBlock.Hash())))

	network.l1Network.Start()
	network.l2Network.Start()

	timeInUs := simulationTime * 1000 * 1000

	go txManager.Start(timeInUs)

	// Wait for the simulation time
	time.Sleep(common.Duration(uint64(timeInUs)))

	fmt.Printf("Stopping simulation after running it for: %s ... \n", common.Duration(uint64(timeInUs)))

	// stop L2 first and then L1
	go network.l2Network.Stop()
	go network.l1Network.Stop()

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
