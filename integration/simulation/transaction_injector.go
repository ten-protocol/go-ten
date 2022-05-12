package simulation

import (
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum"

	"github.com/obscuronet/obscuro-playground/go/ethclient/wallet"

	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/obscuro-playground/contracts"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/host"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/obscuroclient"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/core"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/go/ethclient"
	"github.com/obscuronet/obscuro-playground/go/log"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
	"golang.org/x/sync/errgroup"

	stats2 "github.com/obscuronet/obscuro-playground/integration/simulation/stats"
	wallet_mock "github.com/obscuronet/obscuro-playground/integration/walletmock"
)

// TransactionInjector is a structure that generates, issues and tracks transactions
type TransactionInjector struct {
	// settings
	avgBlockDuration time.Duration
	stats            *stats2.Stats
	wallets          []wallet_mock.Wallet

	l1Nodes       []ethclient.EthClient
	l2NodeClients []*obscuroclient.Client

	l1TransactionsLock sync.RWMutex
	l1Transactions     []obscurocommon.L1Transaction

	l2TransactionsLock sync.RWMutex
	l2Transactions     core.L2Txs

	interruptRun      *int32
	fullyStoppedChan  chan bool
	ethWallets        []wallet.Wallet
	erc20ContractAddr common.Address
}

// NewTransactionInjector returns a transaction manager with a given number of wallets
// todo Add methods that generate deterministic scenarios
func NewTransactionInjector(
	numberWallets int,
	avgBlockDuration time.Duration,
	stats *stats2.Stats,
	l1Nodes []ethclient.EthClient,
	l2Nodes []*host.Node,
	ethWallets []wallet.Wallet,
	addr common.Address,
	l2NodeClients []*obscuroclient.Client,
) *TransactionInjector {
	// create a bunch of wallets
	wallets := make([]wallet_mock.Wallet, numberWallets)
	for i := 0; i < numberWallets; i++ {
		wallets[i] = wallet_mock.New()
	}
	interrupt := int32(0)

	return &TransactionInjector{
		wallets:           wallets,
		avgBlockDuration:  avgBlockDuration,
		stats:             stats,
		l1Nodes:           l1Nodes,
		l2NodeClients:     l2NodeClients,
		interruptRun:      &interrupt,
		fullyStoppedChan:  make(chan bool),
		ethWallets:        ethWallets,
		erc20ContractAddr: addr,
	}
}

// Start begins the execution on the TransactionInjector
// Deposits an initial balance in to each wallet
// Generates and issues L1 and L2 transactions to the network
func (m *TransactionInjector) Start() {
	// deposit some initial amount into every user
	for _, u := range m.wallets {
		txData := &obscurocommon.L1DepositTx{
			Amount: initialBalance,
			To:     u.Address,
		}
		m.rndL1Node().BroadcastTx(txData)
		m.stats.Deposit(initialBalance)
		go m.trackL1Tx(txData)
		time.Sleep(m.avgBlockDuration / 3)
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
		m.issueInvalidWithdrawals()
		return nil
	})

	//wg.Go(func() error {
	//	m.issueRandomERC20Deposits()
	//	return nil
	//})

	_ = wg.Wait() // future proofing to return errors
	m.fullyStoppedChan <- true
}

func (m *TransactionInjector) Stop() {
	atomic.StoreInt32(m.interruptRun, 1)
	for range m.fullyStoppedChan {
		log.Log("TransactionInjector stopped successfully")
		return
	}
}

// trackL1Tx adds an L1Tx to the internal list
func (m *TransactionInjector) trackL1Tx(tx obscurocommon.L1Transaction) {
	m.l1TransactionsLock.Lock()
	defer m.l1TransactionsLock.Unlock()
	m.l1Transactions = append(m.l1Transactions, tx)
}

// trackL2Tx adds an L2Tx to the internal list
func (m *TransactionInjector) trackL2Tx(tx nodecommon.L2Tx) {
	m.l2TransactionsLock.Lock()
	defer m.l2TransactionsLock.Unlock()
	m.l2Transactions = append(m.l2Transactions, tx)
}

// GetL1Transactions returns all generated L1 L2Txs
func (m *TransactionInjector) GetL1Transactions() []obscurocommon.L1Transaction {
	return m.l1Transactions
}

// GetL2Transactions returns all generated non-WithdrawalTx transactions
func (m *TransactionInjector) GetL2Transactions() (core.L2Txs, core.L2Txs) {
	var transfers, withdrawals core.L2Txs
	for _, req := range m.l2Transactions {
		r := req
		switch core.TxData(&r).Type {
		case core.TransferTx:
			transfers = append(transfers, req)
		case core.WithdrawalTx:
			withdrawals = append(withdrawals, req)
		case core.DepositTx:
		}
	}
	return transfers, withdrawals
}

// GetL2WithdrawalRequests returns generated stored WithdrawalTx transactions
func (m *TransactionInjector) GetL2WithdrawalRequests() []nodecommon.Withdrawal {
	var withdrawals []nodecommon.Withdrawal
	for _, req := range m.l2Transactions {
		tx := core.TxData(&req) //nolint:gosec
		if tx.Type == core.WithdrawalTx {
			withdrawals = append(withdrawals, nodecommon.Withdrawal{Amount: tx.Amount, Address: tx.To})
		}
	}
	return withdrawals
}

// issueRandomTransfers creates and issues a number of L2 transfer transactions proportional to the simulation time, such that they can be processed
func (m *TransactionInjector) issueRandomTransfers() {
	for ; atomic.LoadInt32(m.interruptRun) == 0; time.Sleep(obscurocommon.RndBtwTime(m.avgBlockDuration/4, m.avgBlockDuration)) {
		fromWallet := rndWallet(m.wallets)
		to := rndWallet(m.wallets).Address
		for fromWallet.Address == to {
			to = rndWallet(m.wallets).Address
		}
		tx := wallet_mock.NewL2Transfer(fromWallet.Address, to, obscurocommon.RndBtw(1, 500))
		signedTx := wallet_mock.SignTx(tx, fromWallet.Key.PrivateKey)
		encryptedTx := core.EncryptTx(signedTx)
		m.stats.Transfer()

		err := (*m.rndL2NodeClient()).Call(nil, obscuroclient.RPCSendTransactionEncrypted, encryptedTx)
		if err != nil {
			log.Log("Failed to issue transfer via RPC.")
			continue
		}

		go m.trackL2Tx(*signedTx)
	}
}

// issueRandomDeposits creates and issues a number of transactions proportional to the simulation time, such that they can be processed
// Generates L1 common.L1DepositTx transactions
func (m *TransactionInjector) issueRandomDeposits() {
	for ; atomic.LoadInt32(m.interruptRun) == 0; time.Sleep(obscurocommon.RndBtwTime(m.avgBlockDuration, m.avgBlockDuration*2)) {
		v := obscurocommon.RndBtw(1, 100)
		txData := &obscurocommon.L1DepositTx{
			Amount: v,
			To:     rndWallet(m.wallets).Address,
		}
		m.rndL1Node().BroadcastTx(txData)
		m.stats.Deposit(v)
		go m.trackL1Tx(txData)
	}
}

// issueRandomERC20Deposits creates and issues a number of transactions proportional to the simulation time, such that they can be processed
// Generates L1 common.L1DepositTx transactions
func (m *TransactionInjector) issueRandomERC20Deposits() {
	defaultGasPrice := big.NewInt(20000000000)
	defaultGas := uint64(1024_000_000)
	timeout := 30 * time.Second
	w := m.ethWallets[0]
	for ; atomic.LoadInt32(m.interruptRun) == 0; time.Sleep(obscurocommon.RndBtwTime(m.avgBlockDuration, m.avgBlockDuration*2)) {
		tx := &types.LegacyTx{
			Nonce:    w.GetNonceAndIncrement(),
			GasPrice: defaultGasPrice,
			Gas:      defaultGas,
			To:       &m.erc20ContractAddr,
		}

		v := obscurocommon.RndBtw(1, 100)
		data, err := contracts.PedroERC20ContractABIJSON.Pack("transfer", w.Address(), big.NewInt(int64(v)))
		if err != nil {
			panic(err)
		}
		tx.Data = data

		node := m.rndL1Node()
		issuedTx, err := node.SubmitTransaction(tx)
		if err != nil {
			panic(err)
		}

		var receipt *types.Receipt
		for start := time.Now(); time.Since(start) < timeout; time.Sleep(time.Second) {
			receipt, err = node.FetchTxReceipt(issuedTx.Hash())
			if err == nil && receipt != nil {
				break
			}
			if !errors.Is(err, ethereum.NotFound) {
				panic(err)
			}
			log.Log(fmt.Sprintf("Tx has not been mined into a block after %s...", time.Since(start)))
		}

		if receipt == nil || receipt.Status != types.ReceiptStatusSuccessful {
			panic(fmt.Errorf("transaction not minted into a block after %s", timeout))
		}
		m.stats.Deposit(v)
		go m.trackL1Tx(&obscurocommon.L1DepositTx{
			Amount:        v,
			To:            w.Address(),
			TokenContract: m.erc20ContractAddr,
		})
	}
}

// issueRandomWithdrawals creates and issues a number of transactions proportional to the simulation time, such that they can be processed
// Generates L2 enclave2.WithdrawalTx transactions
func (m *TransactionInjector) issueRandomWithdrawals() {
	for ; atomic.LoadInt32(m.interruptRun) == 0; time.Sleep(obscurocommon.RndBtwTime(m.avgBlockDuration, m.avgBlockDuration*2)) {
		v := obscurocommon.RndBtw(1, 100)
		wallet := rndWallet(m.wallets)
		tx := wallet_mock.NewL2Withdrawal(wallet.Address, v)
		signedTx := wallet_mock.SignTx(tx, wallet.Key.PrivateKey)
		encryptedTx := core.EncryptTx(signedTx)

		err := (*m.rndL2NodeClient()).Call(nil, obscuroclient.RPCSendTransactionEncrypted, encryptedTx)
		if err != nil {
			log.Log("Failed to issue withdrawal via RPC.")
			continue
		}

		m.stats.Withdrawal(v)
		go m.trackL2Tx(*signedTx)
	}
}

// issueInvalidWithdrawals creates and issues a number of invalidly-signed L2 withdrawal transactions proportional to the simulation time.
// These transactions should be rejected by the nodes, and thus we expect them not to show up in the simulation withdrawal checks.
func (m *TransactionInjector) issueInvalidWithdrawals() {
	for ; atomic.LoadInt32(m.interruptRun) == 0; time.Sleep(obscurocommon.RndBtwTime(m.avgBlockDuration/4, m.avgBlockDuration)) {
		fromWallet := rndWallet(m.wallets)
		tx := wallet_mock.NewL2Withdrawal(fromWallet.Address, obscurocommon.RndBtw(1, 100))
		signedTx := createInvalidSignature(tx, &fromWallet)
		encryptedTx := core.EncryptTx(signedTx)

		err := (*m.rndL2NodeClient()).Call(nil, obscuroclient.RPCSendTransactionEncrypted, encryptedTx)
		if err != nil {
			log.Log("Failed to issue withdrawal via RPC.")
			continue
		}
	}
}

// Uses one of three approaches to create an invalidly-signed transaction.
func createInvalidSignature(tx *nodecommon.L2Tx, fromWallet *wallet_mock.Wallet) *nodecommon.L2Tx {
	i := rand.Intn(3) //nolint:gosec
	switch i {
	case 0: // We sign the transaction with a bad signer.
		incorrectChainID := int64(enclave.ChainID + 1)
		signer := types.NewLondonSigner(big.NewInt(incorrectChainID))
		signedTx, _ := types.SignTx(tx, signer, fromWallet.Key.PrivateKey)
		return signedTx

	case 1: // We do not sign the transaction.
		return tx

	case 2: // We modify the transaction after signing.
		// We create a new transaction, as we need access to the transaction's encapsulated transaction data.
		txData := core.L2TxData{Type: core.WithdrawalTx, From: fromWallet.Address, Amount: obscurocommon.RndBtw(1, 100)}
		newTx := wallet_mock.NewL2Tx(txData)
		wallet_mock.SignTx(newTx, fromWallet.Key.PrivateKey)
		// After signing the transaction, we create a new transaction based on the transaction data, breaking the signature.
		return wallet_mock.NewL2Tx(txData)
	}
	panic("Expected i to be in the range [0,2).")
}

func rndWallet(wallets []wallet_mock.Wallet) wallet_mock.Wallet {
	return wallets[rand.Intn(len(wallets))] //nolint:gosec
}

func (m *TransactionInjector) rndL1Node() ethclient.EthClient {
	return m.l1Nodes[rand.Intn(len(m.l1Nodes))] //nolint:gosec
}

func (m *TransactionInjector) rndL2NodeClient() *obscuroclient.Client {
	return m.l2NodeClients[rand.Intn(len(m.l2NodeClients))] //nolint:gosec
}
