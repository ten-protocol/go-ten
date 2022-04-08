package simulation

import (
	"math/big"
	"math/rand"
	"sync"
	"time"

	stats2 "github.com/obscuronet/obscuro-playground/integration/simulation/stats"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/host"
	ethereum_mock "github.com/obscuronet/obscuro-playground/integration/ethereummock"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave"
	"golang.org/x/sync/errgroup"

	wallet_mock "github.com/obscuronet/obscuro-playground/integration/walletmock"
)

// TransactionInjector is a structure that generates, issues and tracks transactions
type TransactionInjector struct {
	// settings
	avgBlockDuration uint64
	injectionTimeUs  time.Duration
	stats            *stats2.Stats
	wallets          []wallet_mock.Wallet

	l1Nodes []*ethereum_mock.Node
	l2Nodes []*host.Node

	l1TransactionsLock sync.RWMutex
	l1Transactions     []obscurocommon.L1TxData

	l2TransactionsLock sync.RWMutex
	l2Transactions     enclave.L2Txs
}

// NewTransactionInjector returns a transaction manager with a given number of wallets
// todo Add methods that generate deterministic scenarios
func NewTransactionInjector(numberWallets int, avgBlockDuration uint64, stats *stats2.Stats, injectionTime time.Duration, l1Nodes []*ethereum_mock.Node, l2Nodes []*host.Node) *TransactionInjector {
	// create a bunch of wallets
	wallets := make([]wallet_mock.Wallet, numberWallets)
	for i := 0; i < numberWallets; i++ {
		wallets[i] = wallet_mock.New()
	}

	return &TransactionInjector{
		wallets:          wallets,
		avgBlockDuration: avgBlockDuration,
		stats:            stats,
		injectionTimeUs:  injectionTime,
		l1Nodes:          l1Nodes,
		l2Nodes:          l2Nodes,
	}
}

// Start begins the execution on the TransactionInjector
// Deposits an initial balance in to each wallet
// Generates and issues L1 and L2 transactions to the network
func (m *TransactionInjector) Start() {
	// deposit some initial amount into every user
	for _, u := range m.wallets {
		txData := obscurocommon.L1TxData{
			TxType: obscurocommon.DepositTx,
			Amount: INITIAL_BALANCE,
			Dest:   u.Address,
		}
		tx := obscurocommon.NewL1Tx(txData)
		t, _ := obscurocommon.EncodeTx(tx)
		m.rndL1Node().Network.BroadcastTx(t)
		m.stats.Deposit(INITIAL_BALANCE)
		time.Sleep(obscurocommon.Duration(m.avgBlockDuration / 3))
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

	_ = wg.Wait() // future proofing to return errors
}

// trackL1Tx adds an L1Tx to the internal list
func (m *TransactionInjector) trackL1Tx(tx obscurocommon.L1TxData) {
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
func (m *TransactionInjector) GetL1Transactions() []obscurocommon.L1TxData {
	return m.l1Transactions
}

// GetL2Transactions returns all generated non-WithdrawalTx transactions
func (m *TransactionInjector) GetL2Transactions() enclave.L2Txs {
	var transactions enclave.L2Txs
	for _, req := range m.l2Transactions {
		if enclave.TxData(&req).Type != enclave.WithdrawalTx { //nolint:gosec
			transactions = append(transactions, req)
		}
	}
	return transactions
}

// GetL2WithdrawalRequests returns generated stored WithdrawalTx transactions
func (m *TransactionInjector) GetL2WithdrawalRequests() []nodecommon.Withdrawal {
	var withdrawals []nodecommon.Withdrawal
	for _, req := range m.l2Transactions {
		tx := enclave.TxData(&req) //nolint:gosec
		if tx.Type == enclave.WithdrawalTx {
			withdrawals = append(withdrawals, nodecommon.Withdrawal{Amount: tx.Amount, Address: tx.To})
		}
	}
	return withdrawals
}

// issueRandomTransfers creates and issues a number of L2 transfer transactions proportional to the simulation time, such that they can be processed
func (m *TransactionInjector) issueRandomTransfers() {
	n := uint64(m.injectionTimeUs.Microseconds()) / m.avgBlockDuration
	i := uint64(0)
	for {
		if i == n {
			break
		}
		fromWallet := rndWallet(m.wallets)
		to := rndWallet(m.wallets).Address
		for fromWallet.Address == to {
			to = rndWallet(m.wallets).Address
		}
		tx := wallet_mock.NewL2Transfer(fromWallet.Address, to, obscurocommon.RndBtw(1, 500))
		signedTx := wallet_mock.SignTx(tx, fromWallet.Key.PrivateKey)
		encryptedTx := enclave.EncryptTx(signedTx)
		m.stats.Transfer()
		m.rndL2Node().P2p.BroadcastTx(encryptedTx)
		go m.trackL2Tx(*signedTx)
		time.Sleep(obscurocommon.Duration(obscurocommon.RndBtw(m.avgBlockDuration/4, m.avgBlockDuration)))
		i++
	}
}

// issueRandomDeposits creates and issues a number of transactions proportional to the simulation time, such that they can be processed
// Generates L1 common.DepositTx transactions
func (m *TransactionInjector) issueRandomDeposits() {
	n := uint64(m.injectionTimeUs.Microseconds()) / (m.avgBlockDuration * 3)
	i := uint64(0)
	for {
		if i == n {
			break
		}
		v := obscurocommon.RndBtw(1, 100)
		txData := obscurocommon.L1TxData{
			TxType: obscurocommon.DepositTx,
			Amount: v,
			Dest:   rndWallet(m.wallets).Address,
		}
		tx := obscurocommon.NewL1Tx(txData)
		t, _ := obscurocommon.EncodeTx(tx)
		m.rndL1Node().BroadcastTx(t)
		m.stats.Deposit(v)
		go m.trackL1Tx(txData)
		time.Sleep(obscurocommon.Duration(obscurocommon.RndBtw(m.avgBlockDuration, m.avgBlockDuration*2)))
		i++
	}
}

// issueRandomWithdrawals creates and issues a number of transactions proportional to the simulation time, such that they can be processed
// Generates L2 enclave2.WithdrawalTx transactions
func (m *TransactionInjector) issueRandomWithdrawals() {
	n := uint64(m.injectionTimeUs.Microseconds()) / (m.avgBlockDuration * 3)
	i := uint64(0)
	for {
		if i == n {
			break
		}
		v := obscurocommon.RndBtw(1, 100)
		wallet := rndWallet(m.wallets)
		tx := wallet_mock.NewL2Withdrawal(wallet.Address, v)
		signedTx := wallet_mock.SignTx(tx, wallet.Key.PrivateKey)
		encryptedTx := enclave.EncryptTx(signedTx)
		m.rndL2Node().P2p.BroadcastTx(encryptedTx)
		m.stats.Withdrawal(v)
		go m.trackL2Tx(*signedTx)
		time.Sleep(obscurocommon.Duration(obscurocommon.RndBtw(m.avgBlockDuration, m.avgBlockDuration*2)))
		i++
	}
}

// issueInvalidWithdrawals creates and issues a number of invalidly-signed L2 withdrawal transactions proportional to the simulation time.
// These transactions should be rejected by the nodes, and thus we expect them not to show up in the simulation withdrawal checks.
func (m *TransactionInjector) issueInvalidWithdrawals() {
	n := uint64(m.injectionTimeUs.Microseconds()) / (m.avgBlockDuration * 3)
	i := uint64(0)
	for {
		if i == n {
			break
		}
		fromWallet := rndWallet(m.wallets)
		tx := wallet_mock.NewL2Withdrawal(fromWallet.Address, obscurocommon.RndBtw(1, 100))
		signedTx := createInvalidSignature(tx, &fromWallet)
		encryptedTx := enclave.EncryptTx(signedTx)
		m.rndL2Node().P2p.BroadcastTx(encryptedTx)
		time.Sleep(obscurocommon.Duration(obscurocommon.RndBtw(m.avgBlockDuration/4, m.avgBlockDuration)))
		i++
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
		txData := enclave.L2TxData{Type: enclave.WithdrawalTx, From: fromWallet.Address, Amount: obscurocommon.RndBtw(1, 100)}
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

func (m *TransactionInjector) rndL1Node() *ethereum_mock.Node {
	return m.l1Nodes[rand.Intn(len(m.l1Nodes))] //nolint:gosec
}

func (m *TransactionInjector) rndL2Node() *host.Node {
	return m.l2Nodes[rand.Intn(len(m.l2Nodes))] //nolint:gosec
}
