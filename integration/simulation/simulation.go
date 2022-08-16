package simulation

import (
	"errors"
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

	s.waitForObscuroGenesisOnL1()

	// Arbitrary sleep to wait for RPC clients to get up and running
	time.Sleep(1 * time.Second)

	s.prefundObscuroAccounts() // Prefund every L2 wallet
	s.deployObscuroERC20s()    // Deploy the Obscuro OBX and ETH ERC20 contracts
	s.prefundL1Accounts()      // Prefund every L1 wallet

	timer := time.Now()
	log.Info("Starting injection")
	go s.TxInjector.Start()

	stoppingDelay := s.Params.AvgBlockDuration * 7

	// Wait for the simulation time
	time.Sleep(s.SimulationTime - stoppingDelay)
	log.Info("Stopping injection")

	s.TxInjector.Stop()

	// Allow for some time after tx injection was stopped so that the network can process all transactions
	time.Sleep(stoppingDelay)

	fmt.Printf("Ran simulation for %f secs, configured to run for: %s ... \n", time.Since(timer).Seconds(), s.SimulationTime)
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
		log.Trace("Waiting for the Obscuro genesis rollup...")
	}
}

// Sends an amount from the faucet to each Obscuro account, to pay for transactions.
func (s *Simulation) prefundObscuroAccounts() {
	faucetWallet := s.Params.Wallets.L2FaucetWallet
	nonce := NextNonce(s.RPCHandles, faucetWallet)

	// We send the transactions serially, so that we can precompute the nonces.
	txHashes := make([]gethcommon.Hash, len(s.Params.Wallets.AllObsWallets()))
	for idx, w := range s.Params.Wallets.AllObsWallets() {
		destAddr := w.Address()
		tx := &types.LegacyTx{
			Nonce:    nonce + uint64(idx),
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

		err = s.RPCHandles.ObscuroWalletRndClient(faucetWallet).Call(nil, rpcclientlib.RPCSendRawTransaction, encodeTx(signedTx))
		if err != nil {
			panic(fmt.Sprintf("could not transfer from faucet. Cause: %s", err))
		}

		txHashes[idx] = signedTx.Hash()
	}

	// Then we await the receipts in parallel.
	wg := sync.WaitGroup{}
	for _, txHash := range txHashes {
		wg.Add(1)
		go func(txHash gethcommon.Hash) {
			defer wg.Done()
			err := s.awaitReceipt(faucetWallet, txHash)
			if err != nil {
				panic(fmt.Sprintf("faucet transfer transaction failed. Cause: %s", err))
			}
		}(txHash)
	}
	wg.Wait()
}

// This deploys an ERC20 contract on Obscuro, which is used for token arithmetic.
func (s *Simulation) deployObscuroERC20s() {
	tokens := []bridge.ERC20{bridge.OBX, bridge.ETH}

	wg := sync.WaitGroup{}
	for _, token := range tokens {
		wg.Add(1)
		go func(token bridge.ERC20) {
			defer wg.Done()
			owner := s.Params.Wallets.Tokens[token].L2Owner
			contractBytes := erc20contract.L2BytecodeWithDefaultSupply(string(token))

			deployContractTx := types.DynamicFeeTx{
				Nonce: NextNonce(s.RPCHandles, owner),
				Gas:   1025_000_000,
				Data:  contractBytes,
			}
			signedTx, err := owner.SignTransaction(&deployContractTx)
			if err != nil {
				panic(err)
			}

			err = s.RPCHandles.ObscuroWalletRndClient(owner).Call(nil, rpcclientlib.RPCSendRawTransaction, encodeTx(signedTx))
			if err != nil {
				panic(err)
			}

			err = s.awaitReceipt(owner, signedTx.Hash())
			if err != nil {
				panic(fmt.Sprintf("ERC20 deployment transaction failed. Cause: %s", err))
			}
		}(token)
	}
	wg.Wait()
}

// Sends an amount from the faucet to each L1 account, to pay for transactions.
func (s *Simulation) prefundL1Accounts() {
	for _, w := range s.Params.Wallets.SimEthWallets {
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
		err = s.RPCHandles.RndEthClient().SendTransaction(signedTx)
		if err != nil {
			panic(err)
		}

		s.Stats.Deposit(initialBalance)
		go s.TxInjector.TxTracker.trackL1Tx(txData)
	}
}

// Blocks until the receipt for the transaction has been received. Errors if the transaction is unsuccessful or we time
// out.
func (s *Simulation) awaitReceipt(wallet wallet.Wallet, signedTxHash gethcommon.Hash) error {
	client := s.RPCHandles.ObscuroWalletRndClient(wallet)

	var receipt types.Receipt
	counter := 0
	for {
		err := client.Call(&receipt, rpcclientlib.RPCGetTxReceipt, signedTxHash)
		if err != nil {
			if !errors.Is(err, rpcclientlib.ErrNilResponse) {
				return err
			}

			counter++
			if counter > receiptTimeoutMillis {
				return fmt.Errorf("could not retrieve transaction after timeout")
			}
			time.Sleep(time.Millisecond)
			continue
		}

		if receipt.Status == types.ReceiptStatusFailed {
			return fmt.Errorf("receipt had status failed")
		}

		return nil
	}
}
