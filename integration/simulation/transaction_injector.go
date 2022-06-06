package simulation

import (
	"math/big"
	"math/rand"
	"sync/atomic"
	"time"

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

const maxRetries = 13 // The number of times to retry waiting for an updated nonce for a wallet.

// TransactionInjector is a structure that generates, issues and tracks transactions
type TransactionInjector struct {
	// counters
	counter *txInjectorCounter
	stats   *stats2.Stats

	// settings
	avgBlockDuration time.Duration

	// connections
	issuingWallet wallet.Wallet // the wallet which deploys the erc20
	ethWallets    []wallet.Wallet
	obsWallets    []wallet.Wallet
	l1Clients     []ethclient.EthClient
	l2Clients     []*obscuroclient.Client

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
	ethWallets []wallet.Wallet,
	mgmtContractAddr *common.Address,
	erc20ContractAddr *common.Address,
	l2NodeClients []*obscuroclient.Client,
	mgmtContractLib mgmtcontractlib.MgmtContractLib,
	erc20ContractLib erc20contractlib.ERC20ContractLib,
) *TransactionInjector {
	issuingWallet := wallet.NewInMemoryWalletFromPK(big.NewInt(integration.ObscuroChainID), evm.Erc20OwnerKey)

	interrupt := int32(0)

	obsWallets := make([]wallet.Wallet, len(ethWallets))
	for i, w := range ethWallets {
		obsWallets[i] = wallet.NewInMemoryWalletFromPK(big.NewInt(integration.ObscuroChainID), w.PrivateKey())
	}
	return &TransactionInjector{
		issuingWallet:     issuingWallet,
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
		ethWallets:        ethWallets,
		obsWallets:        obsWallets,
		counter:           newCounter(),
	}
}

// Start begins the execution on the TransactionInjector
// Deposits an initial balance in to each wallet
// Generates and issues L1 and L2 transactions to the network
func (m *TransactionInjector) Start() {
	// always deploy it from the first wallet
	// since it has a hardcoded key
	m.deploySingleObscuroERC20(m.issuingWallet)
	// enough time to process everywhere
	time.Sleep(m.avgBlockDuration * 6)

	// deposit some initial amount into every simulation wallet
	for _, w := range m.ethWallets {
		addr := w.Address()
		txData := &obscurocommon.L1DepositTx{
			Amount:        initialBalance,
			To:            m.mgmtContractAddr,
			TokenContract: m.erc20ContractAddr,
			Sender:        &addr,
		}
		tx := m.erc20ContractLib.CreateDepositTx(txData, w.GetNonceAndIncrement())
		signedTx, err := w.SignTransaction(tx)
		if err != nil {
			panic(err)
		}
		err = m.rndL1NodeClient().SendTransaction(signedTx)
		if err != nil {
			panic(err)
		}

		m.stats.Deposit(initialBalance)
		go m.counter.trackL1Tx(txData)
	}

	// start transactions issuance
	var wg errgroup.Group
	wg.Go(func() error {
		m.issueRandomDeposits()
		return nil
	})

	wg.Go(func() error {
		m.issueRandomWithdrawals()
		return nil
	})

	wg.Go(func() error {
		m.issueRandomTransfers()
		return nil
	})

	wg.Go(func() error {
		m.issueInvalidL2Txs()
		return nil
	})

	_ = wg.Wait() // future proofing to return errors
	m.fullyStoppedChan <- true
}

// This deploys an ERC20 contract on Obscuro, which is used for token arithmetic.
func (m *TransactionInjector) deploySingleObscuroERC20(w wallet.Wallet) {
	// deploy the ERC20
	contractBytes := common.Hex2Bytes(erc20contract.ContractByteCode)
	deployContractTx := types.LegacyTx{
		Nonce:    NextNonce(m.l2Clients[0], w),
		Gas:      1025_000_000,
		GasPrice: common.Big0,
		Data:     contractBytes,
	}
	signedTx, err := w.SignTransaction(&deployContractTx)
	if err != nil {
		panic(err)
	}
	encryptedTx := core.EncryptTx(signedTx)
	err = (*m.rndL2NodeClient()).Call(nil, obscuroclient.RPCSendTransactionEncrypted, encryptedTx)
	if err != nil {
		panic(err)
	}
}

func (m *TransactionInjector) Stop() {
	atomic.StoreInt32(m.interruptRun, 1)
	for range m.fullyStoppedChan {
		log.Info("TransactionInjector stopped successfully")
		return
	}
}

// issueRandomTransfers creates and issues a number of L2 transfer transactions proportional to the simulation time, such that they can be processed
func (m *TransactionInjector) issueRandomTransfers() {
	for ; atomic.LoadInt32(m.interruptRun) == 0; time.Sleep(obscurocommon.RndBtwTime(m.avgBlockDuration/4, m.avgBlockDuration)) {
		fromWallet := m.rndObsWallet()
		toWallet := m.rndObsWallet()
		for fromWallet.Address().Hex() == toWallet.Address().Hex() {
			toWallet = m.rndObsWallet()
		}
		tx := newObscuroTransferTx(fromWallet, toWallet.Address(), obscurocommon.RndBtw(1, 500), m.l2Clients[0])
		signedTx, err := fromWallet.SignTransaction(tx)
		if err != nil {
			panic(err)
		}

		encryptedTx := core.EncryptTx(signedTx)
		m.stats.Transfer()

		err = (*m.rndL2NodeClient()).Call(nil, obscuroclient.RPCSendTransactionEncrypted, encryptedTx)
		if err != nil {
			log.Info("Failed to issue transfer via RPC. Cause: %s", err)
			continue
		}

		go m.counter.trackTransferL2Tx(*signedTx)
	}
}

// issueRandomDeposits creates and issues a number of transactions proportional to the simulation time, such that they can be processed
func (m *TransactionInjector) issueRandomDeposits() {
	for ; atomic.LoadInt32(m.interruptRun) == 0; time.Sleep(obscurocommon.RndBtwTime(m.avgBlockDuration, m.avgBlockDuration*2)) {
		v := obscurocommon.RndBtw(1, 100)
		ethWallet := m.rndEthWallet()
		addr := ethWallet.Address()
		txData := &obscurocommon.L1DepositTx{
			Amount:        v,
			To:            m.mgmtContractAddr,
			TokenContract: m.erc20ContractAddr,
			Sender:        &addr,
		}
		tx := m.erc20ContractLib.CreateDepositTx(txData, ethWallet.GetNonceAndIncrement())
		signedTx, err := ethWallet.SignTransaction(tx)
		if err != nil {
			panic(err)
		}
		err = m.rndL1NodeClient().SendTransaction(signedTx)
		if err != nil {
			panic(err)
		}

		m.stats.Deposit(v)
		go m.counter.trackL1Tx(txData)
	}
}

// issueRandomWithdrawals creates and issues a number of transactions proportional to the simulation time, such that they can be processed
func (m *TransactionInjector) issueRandomWithdrawals() {
	for ; atomic.LoadInt32(m.interruptRun) == 0; time.Sleep(obscurocommon.RndBtwTime(m.avgBlockDuration, m.avgBlockDuration*2)) {
		v := obscurocommon.RndBtw(1, 100)
		obsWallet := m.rndObsWallet()
		// todo - random client
		tx := newObscuroWithdrawalTx(obsWallet, v, m.l2Clients[0])
		signedTx, err := obsWallet.SignTransaction(tx)
		if err != nil {
			panic(err)
		}
		encryptedTx := core.EncryptTx(signedTx)

		err = (*m.rndL2NodeClient()).Call(nil, obscuroclient.RPCSendTransactionEncrypted, encryptedTx)
		if err != nil {
			log.Info("Failed to issue withdrawal via RPC. Cause: %s", err)
			continue
		}

		m.stats.Withdrawal(v)
		go m.counter.trackWithdrawalL2Tx(*signedTx)
	}
}

// issueInvalidL2Txs creates and issues invalidly-signed L2 transactions proportional to the simulation time.
// These transactions should be rejected by the nodes, and thus we expect them to not affect the simulation
func (m *TransactionInjector) issueInvalidL2Txs() {
	for ; atomic.LoadInt32(m.interruptRun) == 0; time.Sleep(obscurocommon.RndBtwTime(m.avgBlockDuration/4, m.avgBlockDuration)) {
		fromWallet := m.rndObsWallet()
		toWallet := m.rndObsWallet()
		for fromWallet.Address().Hex() == toWallet.Address().Hex() {
			toWallet = m.rndObsWallet()
		}
		tx := newCustomObscuroWithdrawalTx(obscurocommon.RndBtw(1, 100))

		signedTx := m.createInvalidSignage(tx, fromWallet)
		encryptedTx := core.EncryptTx(signedTx)

		err := (*m.rndL2NodeClient()).Call(nil, obscuroclient.RPCSendTransactionEncrypted, encryptedTx)
		if err != nil {
			log.Info("Failed to issue withdrawal via RPC. Cause: %s", err)
			continue
		}
	}
}

// Uses one of the approaches to create an invalidly-signed transaction.
func (m *TransactionInjector) createInvalidSignage(tx types.TxData, w wallet.Wallet) *types.Transaction {
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

func (m *TransactionInjector) rndObsWallet() wallet.Wallet {
	return m.obsWallets[rand.Intn(len(m.obsWallets))] //nolint:gosec
}

func (m *TransactionInjector) rndEthWallet() wallet.Wallet {
	return m.ethWallets[rand.Intn(len(m.ethWallets))] //nolint:gosec
}

func (m *TransactionInjector) rndL1NodeClient() ethclient.EthClient {
	return m.l1Clients[rand.Intn(len(m.l1Clients))] //nolint:gosec
}

func (m *TransactionInjector) rndL2NodeClient() *obscuroclient.Client {
	return m.l2Clients[rand.Intn(len(m.l2Clients))] //nolint:gosec
}

func newObscuroTransferTx(from wallet.Wallet, dest common.Address, amount uint64, client *obscuroclient.Client) types.TxData {
	data := erc20contractlib.CreateTransferTxData(dest, amount)
	return newTx(data, NextNonce(client, from))
}

func newObscuroWithdrawalTx(from wallet.Wallet, amount uint64, client *obscuroclient.Client) types.TxData {
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

func readNonce(cl *obscuroclient.Client, a common.Address) uint64 {
	var result uint64
	err := (*cl).Call(&result, obscuroclient.RPCNonce, a)
	if err != nil {
		panic(err)
	}
	return result
}

func NextNonce(cl *obscuroclient.Client, w wallet.Wallet) uint64 {
	retries := 0

	// only returns the nonce when the previous transaction was recorded
	for {
		result := readNonce(cl, w.Address())
		if result == w.GetNonce() {
			return w.GetNonceAndIncrement()
		}
		time.Sleep(time.Duration(2^retries) * time.Millisecond) // For 13 retries, this is around 17 seconds.
		retries++

		if retries > maxRetries {
			panic("Transaction injector failed to retrieve nonce after ten seconds...")
		}
	}
}
