package simulation

import (
	"sync"
	"time"

	common2 "github.com/obscuronet/obscuro-playground/go/obscuronode/common"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave"

	"github.com/google/uuid"
	"github.com/obscuronet/obscuro-playground/go/common"
	wallet_mock "github.com/obscuronet/obscuro-playground/integration/walletmock"
)

type TransactionManager struct {
	l1NetworkConfig    *L1NetworkCfg
	l2NetworkConfig    *L2NetworkCfg
	blockDuration      uint64
	simulationTimeInUS int
	stats              *Stats
	wallets            []wallet_mock.Wallet
	l1TransactionsLock sync.RWMutex
	l1Transactions     common.Transactions
	l2TransactionsLock sync.RWMutex
	l2Transactions     enclave.Transactions
}

func (m *TransactionManager) Start(l1 *L1NetworkCfg, l2 *L2NetworkCfg, duration uint64, us int, stats *Stats) {
	m.l1NetworkConfig = l1
	m.l2NetworkConfig = l2
	m.blockDuration = duration
	m.simulationTimeInUS = us
	m.stats = stats

	// deposit some initial amount into every user
	for _, u := range m.wallets {
		tx := deposit(u, INITIAL_BALANCE)
		t, _ := tx.Encode()
		m.l1NetworkConfig.BroadcastTx(t)
		m.stats.Deposit(INITIAL_BALANCE)
		time.Sleep(common.Duration(m.blockDuration / 3))
	}

	// inject numbers of transactions proportional to the simulation time, such that they can be processed
	go injectRandomDeposits(m.wallets, m.l1NetworkConfig, m.blockDuration, m.simulationTimeInUS, m.stats, m.TrackL1Tx)
	go injectRandomWithdrawals(m.wallets, m.l2NetworkConfig, m.blockDuration, m.simulationTimeInUS, m.stats, m.TrackL2Tx)
	injectRandomTransfers(m.wallets, m.l2NetworkConfig, m.blockDuration, m.simulationTimeInUS, m.stats, m.TrackL2Tx)
}

func (m *TransactionManager) TrackL1Tx(tx common.L1Tx) {
	m.l1TransactionsLock.Lock()
	defer m.l1TransactionsLock.Unlock()
	m.l1Transactions = append(m.l1Transactions, &tx)
}

func (m *TransactionManager) TrackL2Tx(tx enclave.L2Tx) {
	m.l2TransactionsLock.Lock()
	defer m.l2TransactionsLock.Unlock()
	m.l2Transactions = append(m.l2Transactions, tx)
}

func (m *TransactionManager) GetL1Transactions() common.Transactions {
	return m.l1Transactions
}

func (m *TransactionManager) GetL2Transactions() enclave.Transactions {
	var transactions enclave.Transactions
	for _, req := range m.l2Transactions {
		if req.TxType == enclave.TransferTx {
			transactions = append(transactions, req)
		}
	}
	return transactions
}

func (m *TransactionManager) GetL2WithdrawalRequests() []common2.Withdrawal {
	var withdrawals []common2.Withdrawal
	for _, req := range m.l2Transactions {
		if req.TxType == enclave.WithdrawalTx {
			withdrawals = append(withdrawals, common2.Withdrawal{Amount: req.Amount, Address: req.To})
		}
	}
	return withdrawals
}

func NewTransactionGenerator(numberWallets uint) *TransactionManager {
	// create a bunch of wallets
	wallets := make([]wallet_mock.Wallet, numberWallets)
	for i := uint(0); i < numberWallets; i++ {
		wallets[i] = wallet_mock.Wallet{Address: uuid.New().ID()}
	}

	return &TransactionManager{
		wallets: wallets,
	}
}
