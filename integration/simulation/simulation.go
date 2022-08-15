package simulation

import (
	"fmt"
	"math/big"
	"sync"
	"time"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/go-obscuro/go/enclave/bridge"
	"github.com/obscuronet/go-obscuro/go/rpcclientlib"
	"github.com/obscuronet/go-obscuro/go/wallet"
	"github.com/obscuronet/go-obscuro/integration/erc20contract"

	"github.com/obscuronet/go-obscuro/integration/simulation/network"

	"github.com/obscuronet/go-obscuro/go/common/log"

	"github.com/obscuronet/go-obscuro/go/common"

	"github.com/obscuronet/go-obscuro/go/ethadapter"

	"github.com/obscuronet/go-obscuro/integration/simulation/params"

	"github.com/obscuronet/go-obscuro/integration/simulation/stats"
)

const initialBalance = 5000

// Simulation represents all the data required to inject transactions on a network
type Simulation struct {
	RPCHandles       *network.RPCHandles
	AvgBlockDuration uint64
	TxInjector       *TransactionInjector
	SimulationTime   time.Duration
	Stats            *stats.Stats
	Params           *params.SimParams
}

// Start executes the simulation given all the Params. Injects transactions.
func (s *Simulation) Start() {
	log.Info(fmt.Sprintf("Genesis block: b_%d.", common.ShortHash(common.GenesisBlock.Hash())))

	s.waitForObscuroGenesis()

	// arbitrary sleep to wait for RPC clients to get up and running
	time.Sleep(1 * time.Second)

	// deposit some initial amount into every L2 wallet
	s.prefundObscuroAccounts()

	// deploy the Obscuro ERC20 contracts
	s.deployObscuroERC20(s.TxInjector.wallets.Tokens[bridge.OBX].L2Owner)
	s.deployObscuroERC20(s.TxInjector.wallets.Tokens[bridge.ETH].L2Owner)

	// deposit some initial amount into every L1 wallet
	s.prefundL1Accounts()

	// enough time to process everywhere
	time.Sleep(s.Params.AvgBlockDuration * 6)

	timer := time.Now()
	log.Info("Starting injection")
	go s.TxInjector.Start()

	stoppingDelay := s.Params.AvgBlockDuration * 7

	// Wait for the simulation time
	time.Sleep(s.SimulationTime - stoppingDelay)
	log.Info("Stopping injection")

	s.TxInjector.Stop()

	// allow for some time after tx injection was stopped so that the network can process all transactions
	time.Sleep(stoppingDelay)

	fmt.Printf("Ran simulation for %f secs, configured to run for: %s ... \n", time.Since(timer).Seconds(), s.SimulationTime)
}

func (s *Simulation) Stop() {
	// nothing to do for now
}

func (s *Simulation) waitForObscuroGenesis() {
	// grab an L1 client
	client := s.RPCHandles.EthClients[0]

	for {
		// spin through the L1 blocks periodically to see if the genesis rollup has arrived
		head := client.FetchHeadBlock()
		for _, b := range client.BlocksBetween(common.GenesisBlock, head) {
			for _, tx := range b.Transactions() {
				t := s.Params.MgmtContractLib.DecodeTx(tx)
				if t == nil {
					continue
				}
				if _, ok := t.(*ethadapter.L1RollupTx); ok {
					// exit at the first obscuro rollup we see
					return
				}
			}
		}
		time.Sleep(s.Params.AvgBlockDuration)
		log.Trace("Waiting for the Obscuro genesis rollup...")
	}
}

// Sends an amount from the faucet to each Obscuro account, to pay for transactions.
func (s *Simulation) prefundObscuroAccounts() {
	faucetWallet := s.TxInjector.wallets.L2FaucetWallet

	wg := sync.WaitGroup{}
	for _, w := range s.Params.Wallets.AllObsWallets() {
		wg.Add(1)
		go func(wallet wallet.Wallet) {
			defer wg.Done()
			destAddr := wallet.Address()
			tx := &types.LegacyTx{
				Nonce:    NextNonce(s.TxInjector.rpcHandles, faucetWallet),
				Value:    big.NewInt(allocObsWallets),
				Gas:      uint64(1_000_000),
				GasPrice: gethcommon.Big1,
				Data:     nil,
				To:       &destAddr,
			}
			signedTx, err := faucetWallet.SignTransaction(tx)
			if err != nil {
				panic(err)
			}

			err = s.TxInjector.rpcHandles.ObscuroWalletRndClient(faucetWallet).Call(nil, rpcclientlib.RPCSendRawTransaction, encodeTx(signedTx))
			if err != nil {
				panic(fmt.Sprintf("could not transfer from faucet. Cause: %s", err))
			}

			err = s.TxInjector.awaitReceipt(faucetWallet, signedTx.Hash())
			if err != nil {
				panic(fmt.Sprintf("did not get receipt for faucet transfer transaction. Cause: %s", err))
			}
		}(w)
	}

	wg.Wait()
}

// This deploys an ERC20 contract on Obscuro, which is used for token arithmetic.
func (s *Simulation) deployObscuroERC20(owner wallet.Wallet) {
	contractBytes := erc20contract.L2BytecodeWithDefaultSupply(string(bridge.OBX))

	deployContractTx := types.DynamicFeeTx{
		Nonce: NextNonce(s.TxInjector.rpcHandles, owner),
		Gas:   1025_000_000,
		Data:  contractBytes,
	}
	signedTx, err := owner.SignTransaction(&deployContractTx)
	if err != nil {
		panic(err)
	}

	err = s.TxInjector.rpcHandles.ObscuroWalletRndClient(owner).Call(nil, rpcclientlib.RPCSendRawTransaction, encodeTx(signedTx))
	if err != nil {
		panic(err)
	}
}

// Sends an amount from the faucet to each L1 account, to pay for transactions.
func (s *Simulation) prefundL1Accounts() {
	for _, w := range s.TxInjector.wallets.SimEthWallets {
		addr := w.Address()
		txData := &ethadapter.L1DepositTx{
			Amount:        initialBalance,
			To:            s.Params.MgmtContractAddr,
			TokenContract: s.Params.Wallets.Tokens[bridge.OBX].L1ContractAddress,
			Sender:        &addr,
		}
		tx := s.Params.ERC20ContractLib.CreateDepositTx(txData, w.GetNonceAndIncrement())
		signedTx, err := w.SignTransaction(tx)
		if err != nil {
			panic(err)
		}
		err = s.TxInjector.rpcHandles.RndEthClient().SendTransaction(signedTx)
		if err != nil {
			panic(err)
		}

		s.Stats.Deposit(initialBalance)
		go s.TxInjector.TxTracker.trackL1Tx(txData)
	}
}
