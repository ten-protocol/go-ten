package simulation

import (
	"sync"

	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/core"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
)

type txInjectorCounter struct {
	l1TransactionsLock sync.RWMutex
	l1Transactions     []obscurocommon.L1Transaction
	l2TransactionsLock sync.RWMutex
	l2Transactions     core.L2Txs
}

func newCounter() *txInjectorCounter {
	return &txInjectorCounter{
		l1TransactionsLock: sync.RWMutex{},
		l1Transactions:     []obscurocommon.L1Transaction{},
		l2TransactionsLock: sync.RWMutex{},
		l2Transactions:     []nodecommon.L2Tx{},
	}
}

// trackL1Tx adds an L1Tx to the internal list
func (m *txInjectorCounter) trackL1Tx(tx obscurocommon.L1Transaction) {
	m.l1TransactionsLock.Lock()
	defer m.l1TransactionsLock.Unlock()
	m.l1Transactions = append(m.l1Transactions, tx)
}

// trackL2Tx adds an L2Tx to the internal list
func (m *txInjectorCounter) trackL2Tx(tx nodecommon.L2Tx) {
	m.l2TransactionsLock.Lock()
	defer m.l2TransactionsLock.Unlock()
	m.l2Transactions = append(m.l2Transactions, tx)
}

// GetL1Transactions returns all generated L1 L2Txs
func (m *txInjectorCounter) GetL1Transactions() []obscurocommon.L1Transaction {
	return m.l1Transactions
}

// GetL2Transactions returns all generated non-WithdrawalTx transactions
func (m *txInjectorCounter) GetL2Transactions() (core.L2Txs, core.L2Txs) {
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
func (m *txInjectorCounter) GetL2WithdrawalRequests() []nodecommon.Withdrawal {
	var withdrawals []nodecommon.Withdrawal
	for _, req := range m.l2Transactions {
		tx := core.TxData(&req) //nolint:gosec
		if tx.Type == core.WithdrawalTx {
			withdrawals = append(withdrawals, nodecommon.Withdrawal{Amount: tx.Amount, Address: tx.To})
		}
	}
	return withdrawals
}
