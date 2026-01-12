package simulation

import (
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ten-protocol/go-ten/go/wallet"

	"github.com/ten-protocol/go-ten/go/common"
)

type txInjectorTracker struct {
	gasTransactionsLock               sync.RWMutex
	l1TransactionsLock                sync.RWMutex
	L1Transactions                    []common.L1TenTransaction
	l2TransactionsLock                sync.RWMutex
	TransferL2Transactions            []*common.L2Tx
	NativeValueTransferL2Transactions []*common.L2Tx
	WithdrawalL2Transactions          []*common.L2Tx
	GasBridgeTransactions             []GasBridgingRecord
	wethBridgeLock                    sync.RWMutex
	WETHBridgeTransactions            []WETHBridgingRecord
}

type GasBridgingRecord struct {
	L1BridgeTx     *types.Transaction
	ReceiverWallet wallet.Wallet
}

// WETHBridgingRecord tracks WETH bridge transactions for verification
type WETHBridgingRecord struct {
	BridgeTx       *types.Transaction // The bridge transaction (L1 or L2)
	Amount         *big.Int           // Amount of WETH bridged
	ReceiverWallet wallet.Wallet      // Wallet that should receive funds
	IsL1ToL2       bool               // true = L1->L2, false = L2->L1
}

func newCounter() *txInjectorTracker {
	return &txInjectorTracker{
		l1TransactionsLock:                sync.RWMutex{},
		L1Transactions:                    []common.L1TenTransaction{},
		l2TransactionsLock:                sync.RWMutex{},
		TransferL2Transactions:            []*common.L2Tx{},
		WithdrawalL2Transactions:          []*common.L2Tx{},
		NativeValueTransferL2Transactions: []*common.L2Tx{},
		GasBridgeTransactions:             []GasBridgingRecord{},
		WETHBridgeTransactions:            []WETHBridgingRecord{},
	}
}

func (m *txInjectorTracker) trackGasBridgingTx(tx *types.Transaction, receiverWallet wallet.Wallet) {
	m.gasTransactionsLock.Lock()
	defer m.gasTransactionsLock.Unlock()
	m.GasBridgeTransactions = append(m.GasBridgeTransactions, GasBridgingRecord{
		L1BridgeTx:     tx,
		ReceiverWallet: receiverWallet,
	})
}

// trackL1Tx adds an L1Tx to the internal list
func (m *txInjectorTracker) trackL1Tx(tx common.L1TenTransaction) {
	m.l1TransactionsLock.Lock()
	defer m.l1TransactionsLock.Unlock()
	m.L1Transactions = append(m.L1Transactions, tx)
}

func (m *txInjectorTracker) trackTransferL2Tx(tx *common.L2Tx) {
	m.l2TransactionsLock.Lock()
	defer m.l2TransactionsLock.Unlock()
	m.TransferL2Transactions = append(m.TransferL2Transactions, tx)
}

func (m *txInjectorTracker) trackWithdrawalFromL2(tx *common.L2Tx) {
	m.l2TransactionsLock.Lock()
	defer m.l2TransactionsLock.Unlock()
	m.WithdrawalL2Transactions = append(m.WithdrawalL2Transactions, tx)
}

func (m *txInjectorTracker) trackNativeValueTransferL2Tx(tx *common.L2Tx) {
	m.l2TransactionsLock.Lock()
	defer m.l2TransactionsLock.Unlock()
	m.NativeValueTransferL2Transactions = append(m.NativeValueTransferL2Transactions, tx)
}

// GetL1Transactions returns all generated L1 L2Txs
func (m *txInjectorTracker) GetL1Transactions() []common.L1TenTransaction {
	return m.L1Transactions
}

// GetL2Transactions returns all generated non-WithdrawalTx transactions, excluding prefund and ERC20 deploy transactions.
func (m *txInjectorTracker) GetL2Transactions() ([]*common.L2Tx, []*common.L2Tx, []*common.L2Tx) {
	return m.TransferL2Transactions, m.WithdrawalL2Transactions, m.NativeValueTransferL2Transactions
}

// trackWETHBridgingL1ToL2 tracks a WETH bridge transaction from L1 to L2
func (m *txInjectorTracker) trackWETHBridgingL1ToL2(tx *types.Transaction, amount *big.Int, receiverWallet wallet.Wallet) {
	m.wethBridgeLock.Lock()
	defer m.wethBridgeLock.Unlock()
	m.WETHBridgeTransactions = append(m.WETHBridgeTransactions, WETHBridgingRecord{
		BridgeTx:       tx,
		Amount:         amount,
		ReceiverWallet: receiverWallet,
		IsL1ToL2:       true,
	})
}

// trackWETHBridgingL2ToL1 tracks a WETH bridge transaction from L2 to L1
func (m *txInjectorTracker) trackWETHBridgingL2ToL1(tx *types.Transaction, amount *big.Int, receiverWallet wallet.Wallet) {
	m.wethBridgeLock.Lock()
	defer m.wethBridgeLock.Unlock()
	m.WETHBridgeTransactions = append(m.WETHBridgeTransactions, WETHBridgingRecord{
		BridgeTx:       tx,
		Amount:         amount,
		ReceiverWallet: receiverWallet,
		IsL1ToL2:       false,
	})
}
