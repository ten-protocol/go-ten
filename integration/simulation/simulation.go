package simulation

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/obscuronet/go-obscuro/integration/common/testlog"

	"github.com/ethereum/go-ethereum/eth/filters"

	"github.com/ethereum/go-ethereum"

	"github.com/obscuronet/go-obscuro/integration/erc20contract"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/go-obscuro/go/enclave/bridge"
	"github.com/obscuronet/go-obscuro/integration/simulation/network"

	"github.com/obscuronet/go-obscuro/go/common"

	"github.com/obscuronet/go-obscuro/go/ethadapter"

	testcommon "github.com/obscuronet/go-obscuro/integration/common"
	"github.com/obscuronet/go-obscuro/integration/simulation/params"

	"github.com/obscuronet/go-obscuro/integration/simulation/stats"
)

const (
	allocObsWallets = 750_000_000_000_000 // The amount the faucet allocates to each Obscuro wallet.
)

var initialBalance = common.ValueInWei(big.NewInt(5000))

// Simulation represents all the data required to inject transactions on a network
type Simulation struct {
	RPCHandles       *network.RPCHandles
	AvgBlockDuration uint64
	TxInjector       *TransactionInjector
	SimulationTime   time.Duration
	Stats            *stats.Stats
	Params           *params.SimParams
	LogChannels      map[string][]chan common.IDAndLog // Maps an owner to the channels on which they receive logs for each client.
	Subscriptions    []ethereum.Subscription           // A slice of all created event subscriptions.
	ctx              context.Context
}

// Start executes the simulation given all the Params. Injects transactions.
func (s *Simulation) Start() {
	testlog.Logger().Info(fmt.Sprintf("Genesis block: b_%d.", common.ShortHash(common.GenesisBlock.Hash())))
	s.ctx = context.Background() // use injected context for graceful shutdowns

	s.waitForObscuroGenesisOnL1()

	// Arbitrary sleep to wait for RPC clients to get up and running
	time.Sleep(1 * time.Second)

	s.trackLogs()              // Create log subscriptions, to validate that they're working correctly later.
	s.prefundObscuroAccounts() // Prefund every L2 wallet
	s.deployObscuroERC20s()    // Deploy the Obscuro HOC and POC ERC20 contracts
	s.prefundL1Accounts()      // Prefund every L1 wallet

	timer := time.Now()
	fmt.Printf("Starting injection\n")
	testlog.Logger().Info("Starting injection")
	go s.TxInjector.Start()

	stoppingDelay := s.Params.AvgBlockDuration * 7

	// Wait for the simulation time
	time.Sleep(s.SimulationTime - stoppingDelay)
	fmt.Printf("Stopping injection\n")
	testlog.Logger().Info("Stopping injection")

	s.TxInjector.Stop()

	// Allow for some time after tx injection was stopped so that the network can process all transactions
	time.Sleep(stoppingDelay)

	fmt.Printf("Ran simulation for %f secs, configured to run for: %s ... \n", time.Since(timer).Seconds(), s.SimulationTime)
	testlog.Logger().Info(fmt.Sprintf("Ran simulation for %f secs, configured to run for: %s ... \n", time.Since(timer).Seconds(), s.SimulationTime))
}

func (s *Simulation) Stop() {
	// nothing to do for now
}

func (s *Simulation) waitForObscuroGenesisOnL1() {
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
		testlog.Logger().Trace("Waiting for the Obscuro genesis rollup...")
	}
}

// We subscribe to logs on every client for every wallet.
func (s *Simulation) trackLogs() {
	// In-memory clients cannot handle subscriptions for now.
	if s.Params.IsInMem {
		return
	}

	for owner, clients := range s.RPCHandles.AuthObsClients {
		// There is a subscription, and corresponding log channel, per owner per client.
		s.LogChannels[owner] = []chan common.IDAndLog{}

		for _, client := range clients {
			channel := make(chan common.IDAndLog)

			// To exercise the filtering mechanism, we subscribe for HOC events only, ignoring POC events.
			hocFilter := filters.FilterCriteria{
				Addresses: []gethcommon.Address{gethcommon.HexToAddress("0x" + bridge.HOCAddr)},
			}
			sub, err := client.SubscribeFilterLogs(context.Background(), hocFilter, channel)
			if err != nil {
				panic(fmt.Errorf("subscription failed. Cause: %w", err))
			}
			s.Subscriptions = append(s.Subscriptions, sub)

			s.LogChannels[owner] = append(s.LogChannels[owner], channel)
		}
	}
}

// Prefunds the L2 wallets with `allocObsWallets` each.
func (s *Simulation) prefundObscuroAccounts() {
	faucetWallet := s.Params.Wallets.L2FaucetWallet
	faucetClient := s.RPCHandles.ObscuroWalletRndClient(faucetWallet)
	nonce := NextNonce(s.ctx, s.RPCHandles, faucetWallet)
	testcommon.PrefundWallets(s.ctx, faucetWallet, faucetClient, nonce, s.Params.Wallets.AllObsWallets(), big.NewInt(allocObsWallets))
}

// This deploys an ERC20 contract on Obscuro, which is used for token arithmetic.
func (s *Simulation) deployObscuroERC20s() {
	tokens := []bridge.ERC20{bridge.HOC, bridge.POC}

	wg := sync.WaitGroup{}
	for _, token := range tokens {
		wg.Add(1)
		go func(token bridge.ERC20) {
			defer wg.Done()
			owner := s.Params.Wallets.Tokens[token].L2Owner
			/*	if token == "HOC" { //todo:: remove when not using HOC owner key for synthetic transactions
				owner.SetNonce(1)
			}*/

			contractBytes := erc20contract.L2BytecodeWithDefaultSupply(string(token))

			deployContractTx := types.DynamicFeeTx{
				Nonce:     NextNonce(s.ctx, s.RPCHandles, owner),
				Gas:       1025_000_000,
				GasFeeCap: gethcommon.Big1, // This field is used to derive the gas price for dynamic fee transactions.
				Data:      contractBytes,
			}
			signedTx, err := owner.SignTransaction(&deployContractTx)
			if err != nil {
				panic(err)
			}

			err = s.RPCHandles.ObscuroWalletRndClient(owner).SendTransaction(s.ctx, signedTx)
			if err != nil {
				panic(err)
			}

			err = testcommon.AwaitReceipt(s.ctx, s.RPCHandles.ObscuroWalletRndClient(owner), signedTx.Hash())
			if err != nil {
				panic(fmt.Sprintf("ERC20 deployment transaction unsuccessful. Cause: %s", err))
			}
		}(token)
	}
	wg.Wait()
}

// Sends an amount from the faucet to each L1 account, to pay for transactions.
func (s *Simulation) prefundL1Accounts() {
	for _, w := range s.Params.Wallets.SimEthWallets {
		receiver := w.Address()
		tokenOwner := s.Params.Wallets.Tokens[bridge.HOC].L1Owner
		ownerAddr := tokenOwner.Address()
		txData := &ethadapter.L1DepositTx{
			Amount:        initialBalance,
			To:            &receiver,
			TokenContract: s.Params.Wallets.Tokens[bridge.HOC].L1ContractAddress,
			Sender:        &ownerAddr,
		}
		tx := s.Params.ERC20ContractLib.CreateDepositTx(txData, tokenOwner.GetNonceAndIncrement())
		signedTx, err := tokenOwner.SignTransaction(tx)
		if err != nil {
			panic(err)
		}
		err = s.RPCHandles.RndEthClient().SendTransaction(signedTx)
		if err != nil {
			panic(err)
		}

		txClone := signedTx
		go func() {
			time.Sleep(3 * time.Second)
			receipt, err := s.RPCHandles.RndEthClient().TransactionReceipt(txClone.Hash())
			if err != nil {
				panic(err)
			}

			res, err := s.RPCHandles.RndEthClient().CallContract(ethereum.CallMsg{
				From:       ownerAddr,
				To:         txClone.To(),
				Gas:        txClone.Gas(),
				GasPrice:   big.NewInt(20000000000),
				GasTipCap:  big.NewInt(0),
				Value:      txClone.Value(),
				Data:       txClone.Data(),
				AccessList: txClone.AccessList(),
			})
			if err != nil {
				fmt.Printf(fmt.Sprintf("Deposit %s ERROR - %+v", txClone.Hash(), err))
			} else {
				fmt.Printf("Signed Tx - %s bn - %d\n", signedTx.Hash().Hex(), receipt.BlockNumber.Uint64())
				fmt.Printf(fmt.Sprintf("Deposit %s bn Deposit res - %+v", txClone.Hash(), res))
			}
		}()

		s.Stats.Deposit(initialBalance)
		go s.TxInjector.TxTracker.trackL1Tx(txData)
	}
}
