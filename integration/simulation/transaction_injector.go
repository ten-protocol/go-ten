package simulation

import (
	"math"
	"math/big"
	"math/rand"
	"sync/atomic"
	"time"

	"github.com/obscuronet/obscuro-playground/integration"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/obscuro-playground/go/ethclient"
	"github.com/obscuronet/obscuro-playground/go/ethclient/erc20contractlib"
	"github.com/obscuronet/obscuro-playground/go/ethclient/mgmtcontractlib"
	"github.com/obscuronet/obscuro-playground/go/log"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/core"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/obscuroclient"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/wallet"
	"golang.org/x/sync/errgroup"

	stats2 "github.com/obscuronet/obscuro-playground/integration/simulation/stats"
)

// TransactionInjector is a structure that generates, issues and tracks transactions
type TransactionInjector struct {
	// counters
	counter *txInjectorCounter
	stats   *stats2.Stats

	// settings
	avgBlockDuration time.Duration

	// connections
	wallets   []wallet.Wallet
	l1Clients []ethclient.EthClient
	l2Clients []*obscuroclient.Client

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
		wallets:           ethWallets,
		counter:           newCounter(),
	}
}

// Start begins the execution on the TransactionInjector
// Deposits an initial balance in to each wallet
// Generates and issues L1 and L2 transactions to the network
func (m *TransactionInjector) Start() {
	// deposit some initial amount into every simulation wallet
	for _, w := range m.wallets {
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
		fromWallet := m.rndWallet()
		toWallet := m.rndWallet()
		for fromWallet.Address().Hex() == toWallet.Address().Hex() {
			toWallet = m.rndWallet()
		}
		tx := NewL2Transfer(fromWallet.Address(), toWallet.Address(), obscurocommon.RndBtw(1, 500))
		signedTx, err := fromWallet.SignTransaction(tx)
		if err != nil {
			panic(err)
		}

		encryptedTx := core.EncryptTx(signedTx)
		m.stats.Transfer()

		err = (*m.rndL2NodeClient()).Call(nil, obscuroclient.RPCSendTransactionEncrypted, encryptedTx)
		if err != nil {
			log.Info("Failed to issue transfer via RPC.")
			continue
		}

		go m.counter.trackL2Tx(*signedTx)
	}
}

// issueRandomDeposits creates and issues a number of transactions proportional to the simulation time, such that they can be processed
func (m *TransactionInjector) issueRandomDeposits() {
	for ; atomic.LoadInt32(m.interruptRun) == 0; time.Sleep(obscurocommon.RndBtwTime(m.avgBlockDuration, m.avgBlockDuration*2)) {
		v := obscurocommon.RndBtw(1, 100)
		ethWallet := m.rndWallet()
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
		obsWallet := m.rndWallet()
		tx := NewL2Withdrawal(obsWallet.Address(), v)
		signedTx, err := obsWallet.SignTransaction(tx)
		if err != nil {
			panic(err)
		}
		encryptedTx := core.EncryptTx(signedTx)

		err = (*m.rndL2NodeClient()).Call(nil, obscuroclient.RPCSendTransactionEncrypted, encryptedTx)
		if err != nil {
			log.Info("Failed to issue withdrawal via RPC.")
			continue
		}

		m.stats.Withdrawal(v)
		go m.counter.trackL2Tx(*signedTx)
	}
}

// issueInvalidL2Txs creates and issues invalidly-signed L2 transactions proportional to the simulation time.
// These transactions should be rejected by the nodes, and thus we expect them to not affect the simulation
func (m *TransactionInjector) issueInvalidL2Txs() {
	for ; atomic.LoadInt32(m.interruptRun) == 0; time.Sleep(obscurocommon.RndBtwTime(m.avgBlockDuration/4, m.avgBlockDuration)) {
		fromWallet := m.rndWallet()
		toWallet := m.rndWallet()
		for fromWallet.Address().Hex() == toWallet.Address().Hex() {
			toWallet = m.rndWallet()
		}
		var tx types.TxData
		switch rand.Intn(1) {
		case 0:
			tx = NewL2Withdrawal(fromWallet.Address(), obscurocommon.RndBtw(1, 100))
		case 1:
			tx = NewL2Transfer(fromWallet.Address(), toWallet.Address(), obscurocommon.RndBtw(1, 500))
		}

		signedTx := m.createInvalidSignage(tx, fromWallet)
		encryptedTx := core.EncryptTx(signedTx)

		err := (*m.rndL2NodeClient()).Call(nil, obscuroclient.RPCSendTransactionEncrypted, encryptedTx)
		if err != nil {
			log.Info("Failed to issue withdrawal via RPC.")
			continue
		}
	}
}

// Uses one of the approaches to create an invalidly-signed transaction.
func (m *TransactionInjector) createInvalidSignage(tx types.TxData, w wallet.Wallet) *types.Transaction {
	switch rand.Intn(1) {
	case 0: // We sign the transaction with a bad signer.
		incorrectChainID := int64(integration.ChainID + 1)
		signer := types.NewLondonSigner(big.NewInt(incorrectChainID))
		signedTx, _ := types.SignNewTx(w.PrivateKey(), signer, tx)
		return signedTx

	case 1: // We do not sign the transaction.
		return types.NewTx(tx)
	}
	return nil
}

func (m *TransactionInjector) rndWallet() wallet.Wallet {
	return m.wallets[rand.Intn(len(m.wallets)-1)] //nolint:gosec
}

func (m *TransactionInjector) rndL1NodeClient() ethclient.EthClient {
	return m.l1Clients[rand.Intn(len(m.l1Clients))] //nolint:gosec
}

func (m *TransactionInjector) rndL2NodeClient() *obscuroclient.Client {
	return m.l2Clients[rand.Intn(len(m.l2Clients))] //nolint:gosec
}

// NewL2Transfer creates an enclave.L2Tx of type enclave.TransferTx
func NewL2Transfer(from common.Address, dest common.Address, amount uint64) types.TxData {
	txData := core.L2TxData{Type: core.TransferTx, From: from, To: dest, Amount: amount}
	return NewL2Tx(txData)
}

// NewL2Withdrawal creates an enclave.L2Tx of type enclave.WithdrawalTx
func NewL2Withdrawal(from common.Address, amount uint64) types.TxData {
	txData := core.L2TxData{Type: core.WithdrawalTx, From: from, Amount: amount}
	return NewL2Tx(txData)
}

// NewL2Tx creates an enclave.L2Tx.
//
// A random nonce is used to avoid hash collisions. The enclave.L2TxData is encoded and stored in the transaction's
// data field.
func NewL2Tx(data core.L2TxData) types.TxData {
	// We should probably use a deterministic nonce instead, as in the L1.
	nonce := rand.Intn(math.MaxInt) //nolint:gosec

	enc, err := rlp.EncodeToBytes(data)
	if err != nil {
		// TODO - Surface this error properly.
		panic(err)
	}

	return &types.LegacyTx{
		Nonce:    uint64(nonce),
		Value:    big.NewInt(1),
		Gas:      1,
		GasPrice: big.NewInt(1),
		Data:     enc,
	}
}
