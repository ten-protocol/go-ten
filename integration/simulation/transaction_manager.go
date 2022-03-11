package simulation

import (
	"math/rand"
	"sync"
	"time"

	"github.com/obscuronet/obscuro-playground/go/common"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave"
	"golang.org/x/sync/errgroup"

	common2 "github.com/obscuronet/obscuro-playground/go/obscuronode/common"
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
	l1Transactions     []common.L1TxData
	l2TransactionsLock sync.RWMutex
	l2Transactions     enclave.Transactions
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
		txData := common.L1TxData{
			TxType: common.DepositTx,
			Amount: INITIAL_BALANCE,
			Dest:   u.Address,
		}
		tx := common.NewL1Tx(txData)
		t, _ := common.EncodeTx(tx)
		m.l1NetworkConfig.BroadcastTx(t)
		m.stats.Deposit(INITIAL_BALANCE)
		time.Sleep(common.Duration(m.avgBlockDuration / 3))
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

	_ = wg.Wait() // future proofing to return errors
}

// trackL1Tx adds a common.L1Tx to the internal list
func (m *TransactionManager) trackL1Tx(tx common.L1TxData) {
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

// GetL1Transactions returns all generated L1 Transactions
func (m *TransactionManager) GetL1Transactions() []common.L1TxData {
	return m.l1Transactions
}

// GetL2Transactions returns all generated non-WithdrawalTx transactions
func (m *TransactionManager) GetL2Transactions() enclave.Transactions {
	var transactions enclave.Transactions
	for _, req := range m.l2Transactions {
		if enclave.TxData(&req).Type != enclave.WithdrawalTx { //nolint:gosec
			transactions = append(transactions, req)
		}
	}
	return transactions
}

// GetL2WithdrawalRequests returns generated stored WithdrawalTx transactions
func (m *TransactionManager) GetL2WithdrawalRequests() []common2.Withdrawal {
	var withdrawals []common2.Withdrawal
	for _, req := range m.l2Transactions {
		tx := enclave.TxData(&req) //nolint:gosec
		if tx.Type == enclave.WithdrawalTx {
			withdrawals = append(withdrawals, common2.Withdrawal{Amount: tx.Amount, Address: tx.To})
		}
	}
	return withdrawals
}

// issueRandomTransfers creates and issues a numbers of L2 transfer transactions proportional to the simulation time, such that they can be processed
func (m *TransactionManager) issueRandomTransfers() {
	n := uint64(m.simulationTimeInUS) / m.avgBlockDuration
	i := uint64(0)
	for {
		if i == n {
			break
		}
		fromWallet := rndWallet(m.wallets)
		from := fromWallet.Address
		to := rndWallet(m.wallets).Address
		if from == to {
			continue
		}
		tx := wallet_mock.NewL2Transfer(from, to, common.RndBtw(1, 500))
		signedTx := wallet_mock.SignTx(tx, fromWallet.Key.PrivateKey)
		encryptedTx := enclave.EncryptTx(signedTx)
		m.stats.Transfer()
		m.l2NetworkConfig.BroadcastTx(encryptedTx)
		go m.trackL2Tx(*signedTx)
		time.Sleep(common.Duration(common.RndBtw(m.avgBlockDuration/4, m.avgBlockDuration)))
		i++
	}
}

// issueRandomDeposits creates and issues a numbers transactions proportional to the simulation time, such that they can be processed
// Generates L1 common.DepositTx transactions
func (m *TransactionManager) issueRandomDeposits() {
	n := uint64(m.simulationTimeInUS) / (m.avgBlockDuration * 3)
	i := uint64(0)
	for {
		if i == n {
			break
		}
		v := common.RndBtw(1, 100)
		txData := common.L1TxData{
			TxType: common.DepositTx,
			Amount: v,
			Dest:   rndWallet(m.wallets).Address,
		}
		tx := common.NewL1Tx(txData)
		t, _ := common.EncodeTx(tx)
		m.l1NetworkConfig.BroadcastTx(t)
		m.stats.Deposit(v)
		go m.trackL1Tx(txData)
		time.Sleep(common.Duration(common.RndBtw(m.avgBlockDuration, m.avgBlockDuration*2)))
		i++
	}
}

// issueRandomWithdrawals creates and issues a numbers transactions proportional to the simulation time, such that they can be processed
// Generates L2 enclave2.WithdrawalTx transactions
func (m *TransactionManager) issueRandomWithdrawals() {
	n := uint64(m.simulationTimeInUS) / (m.avgBlockDuration * 3)
	i := uint64(0)
	for {
		if i == n {
			break
		}
		v := common.RndBtw(1, 100)
		wallet := rndWallet(m.wallets)
		tx := wallet_mock.NewL2Withdrawal(wallet.Address, v)
		signedTx := wallet_mock.SignTx(tx, wallet.Key.PrivateKey)
		encryptedTx := enclave.EncryptTx(signedTx)
		m.l2NetworkConfig.BroadcastTx(encryptedTx)
		m.stats.Withdrawal(v)
		go m.trackL2Tx(*signedTx)
		time.Sleep(common.Duration(common.RndBtw(m.avgBlockDuration, m.avgBlockDuration*2)))
		i++
	}
}

func rndWallet(wallets []wallet_mock.Wallet) wallet_mock.Wallet {
	return wallets[rand.Intn(len(wallets))] //nolint:gosec
}
