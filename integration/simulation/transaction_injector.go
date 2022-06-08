package simulation

import (
	"math/big"
	"math/rand"
	"sync/atomic"
	"time"

	"github.com/obscuronet/obscuro-playground/integration/simulation/params"

	"github.com/obscuronet/obscuro-playground/integration/erc20contract"

	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/evm"
	"golang.org/x/sync/errgroup"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/obscuroclient"
	"github.com/obscuronet/obscuro-playground/integration"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/go/ethclient"
	"github.com/obscuronet/obscuro-playground/go/ethclient/erc20contractlib"
	"github.com/obscuronet/obscuro-playground/go/ethclient/mgmtcontractlib"
	"github.com/obscuronet/obscuro-playground/go/log"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/core"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/wallet"
	stats2 "github.com/obscuronet/obscuro-playground/integration/simulation/stats"
)

const timeoutMillis = 30000 // The timeout in millis to wait for an updated nonce for a wallet.

// TransactionInjector is a structure that generates, issues and tracks transactions
type TransactionInjector struct {
	// counters
	counter *txInjectorCounter
	stats   *stats2.Stats

	// settings
	avgBlockDuration time.Duration

	// wallets
	wallets *params.SimWallets

	// connections
	l1Clients []ethclient.EthClient
	l2Clients []obscuroclient.Client

	// addrs and libs
	erc20ContractAddr *common.Address
	mgmtContractAddr  *common.Address
	mgmtContractLib   mgmtcontractlib.MgmtContractLib
	erc20ContractLib  erc20contractlib.ERC20ContractLib

	// controls
	interruptRun     *int32
	fullyStoppedChan chan bool
}

// NewTransactionInjector returns a transaction manager with a given number of obsWallets
// todo Add methods that generate deterministic scenarios
func NewTransactionInjector(
	avgBlockDuration time.Duration,
	stats *stats2.Stats,
	l1Nodes []ethclient.EthClient,
	wallets *params.SimWallets,
	mgmtContractAddr *common.Address,
	erc20ContractAddr *common.Address,
	l2NodeClients []obscuroclient.Client,
	mgmtContractLib mgmtcontractlib.MgmtContractLib,
	erc20ContractLib erc20contractlib.ERC20ContractLib,
) *TransactionInjector {
	interrupt := int32(0)

	return &TransactionInjector{
		avgBlockDuration:  avgBlockDuration,
		stats:             stats,
		l1Clients:         l1Nodes,
		l2Clients:         l2NodeClients,
		interruptRun:      &interrupt,
		fullyStoppedChan:  make(chan bool),
		erc20ContractAddr: erc20ContractAddr,
		mgmtContractAddr:  mgmtContractAddr,
		mgmtContractLib:   mgmtContractLib,
		erc20ContractLib:  erc20ContractLib,
		wallets:           wallets,
		counter:           newCounter(),
	}
}

// Start begins the execution on the TransactionInjector
// Deposits an initial balance in to each wallet
// Generates and issues L1 and L2 transactions to the network
func (ti *TransactionInjector) Start() {
	// always deploy it from the first wallet
	// since it has a hardcoded key
	ti.deploySingleObscuroERC20(ti.wallets.Erc20ObsOwnerWallets[0])
	// enough time to process everywhere
	time.Sleep(ti.avgBlockDuration * 6)

	// deposit some initial amount into every simulation wallet
	for _, w := range ti.wallets.SimEthWallets {
		addr := w.Address()
		txData := &obscurocommon.L1DepositTx{
			Amount:        initialBalance,
			To:            ti.mgmtContractAddr,
			TokenContract: ti.erc20ContractAddr,
			Sender:        &addr,
		}
		tx := ti.erc20ContractLib.CreateDepositTx(txData, w.GetNonceAndIncrement())
		signedTx, err := w.SignTransaction(tx)
		if err != nil {
			panic(err)
		}
		err = ti.rndL1NodeClient().SendTransaction(signedTx)
		if err != nil {
			panic(err)
		}

		ti.stats.Deposit(initialBalance)
		go ti.counter.trackL1Tx(txData)
	}

	// start transactions issuance
	var wg errgroup.Group
	wg.Go(func() error {
		ti.issueRandomDeposits()
		return nil
	})

	wg.Go(func() error {
		ti.issueRandomWithdrawals()
		return nil
	})

	wg.Go(func() error {
		ti.issueRandomTransfers()
		return nil
	})

	wg.Go(func() error {
		ti.issueInvalidL2Txs()
		return nil
	})

	_ = wg.Wait() // future proofing to return errors
	ti.fullyStoppedChan <- true
}

// This deploys an ERC20 contract on Obscuro, which is used for token arithmetic.
func (ti *TransactionInjector) deploySingleObscuroERC20(w wallet.Wallet) {
	// deploy the ERC20
	contractBytes := common.Hex2Bytes(erc20contract.ContractByteCode)
	deployContractTx := types.LegacyTx{
		Nonce:    NextNonce(ti.l2Clients[0], w),
		Gas:      1025_000_000,
		GasPrice: common.Big0,
		Data:     contractBytes,
	}
	signedTx, err := w.SignTransaction(&deployContractTx)
	if err != nil {
		panic(err)
	}
	encryptedTx := core.EncryptTx(signedTx)
	err = ti.rndL2NodeClient().Call(nil, obscuroclient.RPCSendTransactionEncrypted, encryptedTx)
	if err != nil {
		panic(err)
	}
}

func (ti *TransactionInjector) Stop() {
	atomic.StoreInt32(ti.interruptRun, 1)
	for range ti.fullyStoppedChan {
		log.Info("TransactionInjector stopped successfully")
		return
	}
}

// issueRandomTransfers creates and issues a number of L2 transfer transactions proportional to the simulation time, such that they can be processed
func (ti *TransactionInjector) issueRandomTransfers() {
	for ; atomic.LoadInt32(ti.interruptRun) == 0; time.Sleep(obscurocommon.RndBtwTime(ti.avgBlockDuration/4, ti.avgBlockDuration)) {
		fromWallet := ti.rndObsWallet()
		toWallet := ti.rndObsWallet()
		for fromWallet.Address().Hex() == toWallet.Address().Hex() {
			toWallet = ti.rndObsWallet()
		}
		tx := newObscuroTransferTx(fromWallet, toWallet.Address(), obscurocommon.RndBtw(1, 500), ti.l2Clients[0])
		signedTx, err := fromWallet.SignTransaction(tx)
		if err != nil {
			panic(err)
		}

		encryptedTx := core.EncryptTx(signedTx)
		ti.stats.Transfer()

		err = ti.rndL2NodeClient().Call(nil, obscuroclient.RPCSendTransactionEncrypted, encryptedTx)
		if err != nil {
			log.Info("Failed to issue transfer via RPC. Cause: %s", err)
			continue
		}

		go ti.counter.trackTransferL2Tx(*signedTx)
	}
}

// issueRandomDeposits creates and issues a number of transactions proportional to the simulation time, such that they can be processed
func (ti *TransactionInjector) issueRandomDeposits() {
	for ; atomic.LoadInt32(ti.interruptRun) == 0; time.Sleep(obscurocommon.RndBtwTime(ti.avgBlockDuration, ti.avgBlockDuration*2)) {
		v := obscurocommon.RndBtw(1, 100)
		ethWallet := ti.rndEthWallet()
		addr := ethWallet.Address()
		txData := &obscurocommon.L1DepositTx{
			Amount:        v,
			To:            ti.mgmtContractAddr,
			TokenContract: ti.erc20ContractAddr,
			Sender:        &addr,
		}
		tx := ti.erc20ContractLib.CreateDepositTx(txData, ethWallet.GetNonceAndIncrement())
		signedTx, err := ethWallet.SignTransaction(tx)
		if err != nil {
			panic(err)
		}
		err = ti.rndL1NodeClient().SendTransaction(signedTx)
		if err != nil {
			panic(err)
		}

		ti.stats.Deposit(v)
		go ti.counter.trackL1Tx(txData)
	}
}

// issueRandomWithdrawals creates and issues a number of transactions proportional to the simulation time, such that they can be processed
func (ti *TransactionInjector) issueRandomWithdrawals() {
	for ; atomic.LoadInt32(ti.interruptRun) == 0; time.Sleep(obscurocommon.RndBtwTime(ti.avgBlockDuration, ti.avgBlockDuration*2)) {
		v := obscurocommon.RndBtw(1, 100)
		obsWallet := ti.rndObsWallet()
		// todo - random client
		tx := newObscuroWithdrawalTx(obsWallet, v, ti.l2Clients[0])
		signedTx, err := obsWallet.SignTransaction(tx)
		if err != nil {
			panic(err)
		}
		encryptedTx := core.EncryptTx(signedTx)

		err = ti.rndL2NodeClient().Call(nil, obscuroclient.RPCSendTransactionEncrypted, encryptedTx)
		if err != nil {
			log.Info("Failed to issue withdrawal via RPC. Cause: %s", err)
			continue
		}

		ti.stats.Withdrawal(v)
		go ti.counter.trackWithdrawalL2Tx(*signedTx)
	}
}

// issueInvalidL2Txs creates and issues invalidly-signed L2 transactions proportional to the simulation time.
// These transactions should be rejected by the nodes, and thus we expect them to not affect the simulation
func (ti *TransactionInjector) issueInvalidL2Txs() {
	for ; atomic.LoadInt32(ti.interruptRun) == 0; time.Sleep(obscurocommon.RndBtwTime(ti.avgBlockDuration/4, ti.avgBlockDuration)) {
		fromWallet := ti.rndObsWallet()
		toWallet := ti.rndObsWallet()
		for fromWallet.Address().Hex() == toWallet.Address().Hex() {
			toWallet = ti.rndObsWallet()
		}
		tx := newCustomObscuroWithdrawalTx(obscurocommon.RndBtw(1, 100))

		signedTx := ti.createInvalidSignage(tx, fromWallet)
		encryptedTx := core.EncryptTx(signedTx)

		err := ti.rndL2NodeClient().Call(nil, obscuroclient.RPCSendTransactionEncrypted, encryptedTx)
		if err != nil {
			log.Info("Failed to issue withdrawal via RPC. Cause: %s", err)
			continue
		}
	}
}

// Uses one of the approaches to create an invalidly-signed transaction.
func (ti *TransactionInjector) createInvalidSignage(tx types.TxData, w wallet.Wallet) *types.Transaction {
	switch rand.Intn(2) { //nolint:gosec
	case 0: // We sign the transaction with a bad signer.
		incorrectChainID := int64(integration.EthereumChainID + 1)
		signer := types.NewLondonSigner(big.NewInt(incorrectChainID))
		signedTx, _ := types.SignNewTx(w.PrivateKey(), signer, tx)
		return signedTx

	case 1: // We do not sign the transaction.
		return types.NewTx(tx)
	}
	return nil
}

func (ti *TransactionInjector) rndObsWallet() wallet.Wallet {
	return ti.wallets.SimObsWallets[rand.Intn(len(ti.wallets.SimObsWallets))] //nolint:gosec
}

func (ti *TransactionInjector) rndEthWallet() wallet.Wallet {
	return ti.wallets.SimEthWallets[rand.Intn(len(ti.wallets.SimEthWallets))] //nolint:gosec
}

func (ti *TransactionInjector) rndL1NodeClient() ethclient.EthClient {
	return ti.l1Clients[rand.Intn(len(ti.l1Clients))] //nolint:gosec
}

func (ti *TransactionInjector) rndL2NodeClient() obscuroclient.Client {
	return ti.l2Clients[rand.Intn(len(ti.l2Clients))] //nolint:gosec
}

func newObscuroTransferTx(from wallet.Wallet, dest common.Address, amount uint64, client obscuroclient.Client) types.TxData {
	data := erc20contractlib.CreateTransferTxData(dest, amount)
	return newTx(data, NextNonce(client, from))
}

func newObscuroWithdrawalTx(from wallet.Wallet, amount uint64, client obscuroclient.Client) types.TxData {
	transferERC20data := erc20contractlib.CreateTransferTxData(evm.WithdrawalAddress, amount)
	return newTx(transferERC20data, NextNonce(client, from))
}

func newCustomObscuroWithdrawalTx(amount uint64) types.TxData {
	transferERC20data := erc20contractlib.CreateTransferTxData(evm.WithdrawalAddress, amount)
	return newTx(transferERC20data, 1)
}

func newTx(data []byte, nonce uint64) types.TxData {
	return &types.LegacyTx{
		Nonce:    nonce,
		Value:    common.Big0,
		Gas:      1_000_000,
		GasPrice: common.Big0,
		Data:     data,
		To:       &evm.Erc20ContractAddress,
	}
}

func readNonce(cl obscuroclient.Client, a common.Address) uint64 {
	var result uint64
	err := cl.Call(&result, obscuroclient.RPCNonce, a)
	if err != nil {
		panic(err)
	}
	return result
}

func NextNonce(cl obscuroclient.Client, w wallet.Wallet) uint64 {
	counter := 0

	// only returns the nonce when the previous transaction was recorded
	for {
		remoteNonce := readNonce(cl, w.Address())
		localNonce := w.GetNonce()
		if remoteNonce == localNonce {
			return w.GetNonceAndIncrement()
		}
		if remoteNonce > localNonce {
			panic("remote nonce exceeds local nonce")
		}

		counter++
		if counter > timeoutMillis {
			panic("transaction injector failed to retrieve nonce after thirty seconds")
		}
		time.Sleep(time.Millisecond)
	}
}
