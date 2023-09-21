package simulation

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/filters"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/errutil"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/ethadapter"
	"github.com/obscuronet/go-obscuro/go/wallet"
	"github.com/obscuronet/go-obscuro/integration/common/testlog"
	"github.com/obscuronet/go-obscuro/integration/erc20contract"
	"github.com/obscuronet/go-obscuro/integration/ethereummock"
	"github.com/obscuronet/go-obscuro/integration/simulation/network"
	"github.com/obscuronet/go-obscuro/integration/simulation/params"
	"github.com/obscuronet/go-obscuro/integration/simulation/stats"

	gethcommon "github.com/ethereum/go-ethereum/common"
	testcommon "github.com/obscuronet/go-obscuro/integration/common"
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
	testlog.Logger().Info(fmt.Sprintf("Genesis block: b_%d.", common.ShortHash(ethereummock.MockGenesisBlock.Hash())))
	s.ctx = context.Background() // use injected context for graceful shutdowns

	s.waitForObscuroGenesisOnL1()

	// Arbitrary sleep to wait for RPC clients to get up and running
	time.Sleep(1 * time.Second)

	s.trackLogs()              // Create log subscriptions, to validate that they're working correctly later.
	s.prefundObscuroAccounts() // Prefund every L2 wallet
	s.deployObscuroERC20s()    // Deploy the Obscuro HOC and POC ERC20 contracts
	s.prefundL1Accounts()      // Prefund every L1 wallet
	s.checkHealthStatus()      // Checks the nodes health status

	timer := time.Now()
	fmt.Printf("Starting injection\n")
	testlog.Logger().Info("Starting injection")
	go s.TxInjector.Start()

	// Allow for some time after tx injection was stopped so that the network can process all transactions, catch up
	// on missed batches, etc.

	// Wait for the simulation time
	time.Sleep(s.SimulationTime - s.Params.StoppingDelay)
	fmt.Printf("Stopping injection\n")
	testlog.Logger().Info("Stopping injection")

	s.TxInjector.Stop()

	time.Sleep(s.Params.StoppingDelay)

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
		head, err := client.FetchHeadBlock()
		if err != nil && !errors.Is(err, errutil.ErrNotFound) {
			panic(fmt.Errorf("could not fetch head block. Cause: %w", err))
		}
		if err == nil {
			for _, b := range client.BlocksBetween(ethereummock.MockGenesisBlock, head) {
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
			channel := make(chan common.IDAndLog, 1000)

			// To exercise the filtering mechanism, we subscribe for HOC events only, ignoring POC events.
			hocFilter := filters.FilterCriteria{
				Addresses: []gethcommon.Address{gethcommon.HexToAddress("0x" + testcommon.HOCAddr)},
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
	testcommon.PrefundWallets(s.ctx, faucetWallet, faucetClient, nonce, s.Params.Wallets.AllObsWallets(), big.NewInt(allocObsWallets), s.Params.ReceiptTimeout)
}

// This deploys an ERC20 contract on Obscuro, which is used for token arithmetic.
func (s *Simulation) deployObscuroERC20s() {
	tokens := []testcommon.ERC20{testcommon.HOC, testcommon.POC}

	wg := sync.WaitGroup{}
	for _, token := range tokens {
		wg.Add(1)
		go func(token testcommon.ERC20) {
			defer wg.Done()
			owner := s.Params.Wallets.Tokens[token].L2Owner
			// 0x526c84529b2b8c11f57d93d3f5537aca3aecef9b - this is the address of the L2 contract which is currently hardcoded.
			contractBytes := erc20contract.L2BytecodeWithDefaultSupply(string(token), gethcommon.HexToAddress("0x526c84529b2b8c11f57d93d3f5537aca3aecef9b"))

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

			err = testcommon.AwaitReceipt(s.ctx, s.RPCHandles.ObscuroWalletRndClient(owner), signedTx.Hash(), s.Params.ReceiptTimeout)
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
		ethClient := s.RPCHandles.RndEthClient()
		receiver := w.Address()
		tokenOwner := s.Params.Wallets.Tokens[testcommon.HOC].L1Owner
		ownerAddr := tokenOwner.Address()
		txData := &ethadapter.L1DepositTx{
			Amount:        initialBalance,
			To:            &receiver,
			TokenContract: s.Params.Wallets.Tokens[testcommon.HOC].L1ContractAddress,
			Sender:        &ownerAddr,
		}
		tx := s.Params.ERC20ContractLib.CreateDepositTx(txData)
		estimatedTx, err := ethClient.PrepareTransactionToSend(tx, tokenOwner.Address(), tokenOwner.GetNonceAndIncrement())
		if err != nil {
			// ignore txs that are not able to be estimated/execute
			testlog.Logger().Error("unable to estimate tx", log.ErrKey, err)
			tokenOwner.SetNonce(tokenOwner.GetNonce() - 1)
			continue
		}
		signedTx, err := tokenOwner.SignTransaction(estimatedTx)
		if err != nil {
			panic(err)
		}
		err = s.RPCHandles.RndEthClient().SendTransaction(signedTx)
		if err != nil {
			panic(err)
		}

		go s.TxInjector.TxTracker.trackL1Tx(txData)
	}
}

func (s *Simulation) checkHealthStatus() {
	for _, client := range s.RPCHandles.ObscuroClients {
		if healthy, err := client.Health(); !healthy || err != nil {
			panic("Client is not healthy")
		}
	}
}

func NextNonce(ctx context.Context, clients *network.RPCHandles, w wallet.Wallet) uint64 {
	counter := 0

	// only returns the nonce when the previous transaction was recorded
	for {
		remoteNonce, err := clients.ObscuroWalletRndClient(w).NonceAt(ctx, nil)
		if err != nil {
			panic(err)
		}
		localNonce := w.GetNonce()
		if remoteNonce == localNonce {
			return w.GetNonceAndIncrement()
		}
		if remoteNonce > localNonce {
			panic("remote nonce exceeds local nonce")
		}

		counter++
		if counter > nonceTimeoutMillis {
			panic(fmt.Sprintf("transaction injector could not retrieve nonce after thirty seconds for address %s. "+
				"Local nonce was %d, remote nonce was %d", w.Address().Hex(), localNonce, remoteNonce))
		}
		time.Sleep(time.Millisecond)
	}
}
