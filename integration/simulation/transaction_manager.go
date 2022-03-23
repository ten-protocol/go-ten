package simulation

import (
	"math/rand"
	"sync"
	"time"

	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave"
	"golang.org/x/sync/errgroup"

	wallet_mock "github.com/obscuronet/obscuro-playground/integration/walletmock"
)

// TransactionManager is a structure that generates, issues and tracks transactions
type TransactionManager struct {
	l1NetworkConfig    *L1NetworkCfg
	l2NetworkConfig    *L2NetworkCfg
	avgBlockDuration   uint64
	simulationTimeInUS int
	stats              *Stats
	wallets            []wallet_mock.Wallet
	l1TransactionsLock sync.RWMutex
	l1Transactions     []obscurocommon.L1TxData
	l2TransactionsLock sync.RWMutex
	l2Transactions     enclave.L2Txs
}

// NewTransactionManager returns a transaction manager with a given number of wallets
// todo Add methods that generate deterministic scenarios
func NewTransactionManager(numberWallets uint, l1 *L1NetworkCfg, l2 *L2NetworkCfg, avgBlockDuration uint64, stats *Stats) *TransactionManager {
	// create a bunch of wallets
	wallets := make([]wallet_mock.Wallet, numberWallets)
	for i := uint(0); i < numberWallets; i++ {
		wallets[i] = wallet_mock.New()
	}

	return &TransactionManager{
		wallets:          wallets,
		l1NetworkConfig:  l1,
		l2NetworkConfig:  l2,
		avgBlockDuration: avgBlockDuration,
		stats:            stats,
	}
}

// Start begins the execution on the TransactionManager
// Deposits an initial balance in to each wallet
// Generates and issues L1 and L2 transactions to the network
func (m *TransactionManager) Start(us int) {
	m.simulationTimeInUS = us

	// deposit some initial amount into every user
	for _, u := range m.wallets {
		txData := obscurocommon.L1TxData{
			TxType: obscurocommon.DepositTx,
			Amount: INITIAL_BALANCE,
			Dest:   u.Address,
		}
		tx := obscurocommon.NewL1Tx(txData)
		t, _ := obscurocommon.EncodeTx(tx)
		m.l1NetworkConfig.BroadcastTx(t)
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
		m.issueInvalidTransfers()
		return nil
	})

	_ = wg.Wait() // future proofing to return errors
}

// trackL1Tx adds a common.L1Tx to the internal list
func (m *TransactionManager) trackL1Tx(tx obscurocommon.L1TxData) {
	m.l1TransactionsLock.Lock()
	defer m.l1TransactionsLock.Unlock()
	m.l1Transactions = append(m.l1Transactions, tx)
}

// trackL2Tx adds an enclave.L2Tx to the internal list
func (m *TransactionManager) trackL2Tx(tx enclave.L2Tx) {
	m.l2TransactionsLock.Lock()
	defer m.l2TransactionsLock.Unlock()
	m.l2Transactions = append(m.l2Transactions, tx)
}

// GetL1Transactions returns all generated L1 L2Txs
func (m *TransactionManager) GetL1Transactions() []obscurocommon.L1TxData {
	return m.l1Transactions
}

// GetL2Transactions returns all generated non-WithdrawalTx transactions
func (m *TransactionManager) GetL2Transactions() enclave.L2Txs {
	var transactions enclave.L2Txs
	for _, req := range m.l2Transactions {
		if enclave.TxData(&req).Type != enclave.WithdrawalTx { //nolint:gosec
			transactions = append(transactions, req)
		}
	}
	return transactions
}

// GetL2WithdrawalRequests returns generated stored WithdrawalTx transactions
func (m *TransactionManager) GetL2WithdrawalRequests() []nodecommon.Withdrawal {
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
func (m *TransactionManager) issueRandomTransfers() {
	n := uint64(m.simulationTimeInUS) / m.avgBlockDuration
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
		m.l2NetworkConfig.BroadcastTx(encryptedTx)
		go m.trackL2Tx(*signedTx)
		time.Sleep(obscurocommon.Duration(obscurocommon.RndBtw(m.avgBlockDuration/4, m.avgBlockDuration)))
		i++
	}
}

// issueRandomDeposits creates and issues a number of transactions proportional to the simulation time, such that they can be processed
// Generates L1 common.DepositTx transactions
func (m *TransactionManager) issueRandomDeposits() {
	n := uint64(m.simulationTimeInUS) / (m.avgBlockDuration * 3)
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
		m.l1NetworkConfig.BroadcastTx(t)
		m.stats.Deposit(v)
		go m.trackL1Tx(txData)
		time.Sleep(obscurocommon.Duration(obscurocommon.RndBtw(m.avgBlockDuration, m.avgBlockDuration*2)))
		i++
	}
}

// issueRandomWithdrawals creates and issues a number of transactions proportional to the simulation time, such that they can be processed
// Generates L2 enclave2.WithdrawalTx transactions
func (m *TransactionManager) issueRandomWithdrawals() {
	n := uint64(m.simulationTimeInUS) / (m.avgBlockDuration * 3)
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
		m.l2NetworkConfig.BroadcastTx(encryptedTx)
		m.stats.Withdrawal(v)
		go m.trackL2Tx(*signedTx)
		time.Sleep(obscurocommon.Duration(obscurocommon.RndBtw(m.avgBlockDuration, m.avgBlockDuration*2)))
		i++
	}
}

// issueInvalidTransfers creates and issues a number of invalidly-signed L2 transfer transactions proportional to the simulation time.
// These transactions should be rejected by the nodes, and thus we expect them not to show up in the simulation stats.
func (m *TransactionManager) issueInvalidTransfers() {
	n := uint64(m.simulationTimeInUS) / (m.avgBlockDuration * 3)
	i := uint64(0)
	for {
		if i == n {
			break
		}
		fromWallet := rndWallet(m.wallets)
		toWallet := rndWallet(m.wallets)
		for fromWallet.Address == toWallet.Address {
			toWallet = rndWallet(m.wallets)
		}
		tx := wallet_mock.NewL2Transfer(fromWallet.Address, toWallet.Address, obscurocommon.RndBtw(1, 500))
		signedTx := createInvalidSignature(tx, &fromWallet, &toWallet)
		encryptedTx := enclave.EncryptTx(signedTx)
		m.l2NetworkConfig.BroadcastTx(encryptedTx)
		time.Sleep(obscurocommon.Duration(obscurocommon.RndBtw(m.avgBlockDuration/4, m.avgBlockDuration)))
		i++
	}
}

// Uses one of three approaches to create an invalidly-signed transaction.
func createInvalidSignature(tx *enclave.L2Tx, fromWallet *wallet_mock.Wallet, toWallet *wallet_mock.Wallet) *enclave.L2Tx {
	i := rand.Intn(3)
	switch i {
	case 0: // We sign the transaction with the wrong key.
		return wallet_mock.SignTx(tx, toWallet.Key.PrivateKey)

	case 1: // We do not sign the transaction.
		return tx

	case 2: // We modify the transaction after signing.
		// We create a new transaction, as we need access to the transaction's encapsulated transaction data.
		txData := enclave.L2TxData{Type: enclave.WithdrawalTx, From: fromWallet.Address, To: toWallet.Address, Amount: obscurocommon.RndBtw(1, 500)}
		newTx := wallet_mock.NewL2Tx(txData)
		wallet_mock.SignTx(newTx, fromWallet.Key.PrivateKey)
		// After signing the transaction, we modify its transaction data, breaking the signature.
		txData.Type = enclave.TransferTx
		modifiedTx := wallet_mock.NewL2Tx(txData)
		return modifiedTx
	}
	panic("Expected i to be in the range [0,2).")
}

func rndWallet(wallets []wallet_mock.Wallet) wallet_mock.Wallet {
	return wallets[rand.Intn(len(wallets))] //nolint:gosec
}
