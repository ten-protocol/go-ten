package simulation

import (
	"fmt"
	"math/big"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/obscuronet/obscuro-playground/contracts"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/evm"
	"golang.org/x/sync/errgroup"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/obscuroclient"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/core"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/go/ethclient"
	"github.com/obscuronet/obscuro-playground/go/log"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
	stats2 "github.com/obscuronet/obscuro-playground/integration/simulation/stats"
	wallet_mock "github.com/obscuronet/obscuro-playground/integration/walletmock"
)

// TransactionInjector is a structure that generates, issues and tracks transactions
type TransactionInjector struct {
	// settings
	avgBlockDuration time.Duration
	stats            *stats2.Stats
	issuingWallet    *wallet_mock.Wallet // the wallet which deploys the erc20
	wallets          []*wallet_mock.Wallet

	l1Nodes       []ethclient.EthClient
	l2NodeClients []*obscuroclient.Client

	l1TransactionsLock sync.RWMutex
	l1Transactions     []obscurocommon.L1TxData

	l2TransactionsLock       sync.RWMutex
	transferL2Transactions   core.L2Txs
	withdrawalL2Transactions core.L2Txs

	interruptRun     *int32
	fullyStoppedChan chan bool
}

// NewTransactionInjector returns a transaction manager with a given number of wallets
// todo Add methods that generate deterministic scenarios
func NewTransactionInjector(
	numberWallets int,
	avgBlockDuration time.Duration,
	stats *stats2.Stats,
	l1Nodes []ethclient.EthClient,
	l2NodeClients []*obscuroclient.Client,
) *TransactionInjector {
	issuingWallet := wallet_mock.New(evm.Erc20OwnerKey)

	// create a bunch of wallets
	wallets := make([]*wallet_mock.Wallet, numberWallets)
	for i := 0; i < numberWallets; i++ {
		key, err := crypto.GenerateKey()
		if err != nil {
			panic(fmt.Sprintf("Could not generate keypair for wallet: %v", err))
		}
		wallets[i] = wallet_mock.New(key)
	}
	interrupt := int32(0)

	return &TransactionInjector{
		issuingWallet:    issuingWallet,
		wallets:          wallets,
		avgBlockDuration: avgBlockDuration,
		stats:            stats,
		l1Nodes:          l1Nodes,
		l2NodeClients:    l2NodeClients,
		interruptRun:     &interrupt,
		fullyStoppedChan: make(chan bool),
	}
}

// Start begins the execution on the TransactionInjector
// Deposits an initial balance in to each wallet
// Generates and issues L1 and L2 transactions to the network
func (m *TransactionInjector) Start() {
	// always deploy it from the first wallet
	// since it has a hardcoded key
	m.deployERC20(m.issuingWallet)
	// enough time to process everywhere
	time.Sleep(m.avgBlockDuration * 6)

	// deposit some initial amount into every user
	for _, u := range m.wallets {
		txData := &obscurocommon.L1TxData{
			TxType: obscurocommon.DepositTx,
			Amount: initialBalance,
			Dest:   u.Address,
		}
		// fmt.Printf("Injected l1 deposit tx: %v\n", txData)
		m.rndL1Node().BroadcastTx(txData)
		m.stats.Deposit(initialBalance)
		go m.trackL1Tx(*txData)
		time.Sleep(m.avgBlockDuration / 3)
	}

	// start transactions issuance
	var wg errgroup.Group
	wg.Go(func() error {
		m.issueRandomL1Deposits()
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

	// todo
	//wg.Go(func() error {
	//	m.issueInvalidWithdrawals()
	//	return nil
	//})

	_ = wg.Wait() // future proofing to return errors
	m.fullyStoppedChan <- true
}

func (m *TransactionInjector) deployERC20(w *wallet_mock.Wallet) {
	// deploy the ERC20
	contractBytes := common.Hex2Bytes(contracts.PedroERC20ContractByteCode)
	deployContractTx := types.LegacyTx{
		Nonce:    w.NextNonce(m.l2NodeClients[0]),
		Gas:      1025_000_000,
		GasPrice: common.Big0,
		Data:     contractBytes,
	}
	signedTx, err := types.SignTx(types.NewTx(&deployContractTx), types.NewLondonSigner(big.NewInt(int64(evm.ChainID))), w.Key.PrivateKey)
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

// trackL1Tx adds an L1Tx to the internal list
func (m *TransactionInjector) trackL1Tx(tx obscurocommon.L1TxData) {
	m.l1TransactionsLock.Lock()
	defer m.l1TransactionsLock.Unlock()
	m.l1Transactions = append(m.l1Transactions, tx)
}

func (m *TransactionInjector) trackWithdrawalL2Tx(tx nodecommon.L2Tx) {
	m.l2TransactionsLock.Lock()
	defer m.l2TransactionsLock.Unlock()
	m.withdrawalL2Transactions = append(m.withdrawalL2Transactions, tx)
}

func (m *TransactionInjector) trackTransferL2Tx(tx nodecommon.L2Tx) {
	m.l2TransactionsLock.Lock()
	defer m.l2TransactionsLock.Unlock()
	m.transferL2Transactions = append(m.transferL2Transactions, tx)
}

// GetL1Transactions returns all generated L1 L2Txs
func (m *TransactionInjector) GetL1Transactions() []obscurocommon.L1TxData {
	return m.l1Transactions
}

// GetL2Transactions returns all generated non-WithdrawalTx transactions
func (m *TransactionInjector) GetL2Transactions() (core.L2Txs, core.L2Txs) {
	return m.transferL2Transactions, m.withdrawalL2Transactions
}

// GetL2WithdrawalRequests returns generated stored WithdrawalTx transactions
func (m *TransactionInjector) GetL2WithdrawalRequests() []nodecommon.Withdrawal {
	withdrawals := make([]nodecommon.Withdrawal, 0)
	for _, req := range m.withdrawalL2Transactions {
		// todo - helper
		method, err := contracts.PedroERC20ContractABIJSON.MethodById(req.Data()[:4])
		if err != nil || method.Name != "transfer" {
			panic(err)
		}
		args := map[string]interface{}{}
		if err := method.Inputs.UnpackIntoMap(args, req.Data()[4:]); err != nil {
			panic(err)
		}
		withdrawals = append(withdrawals, nodecommon.Withdrawal{Amount: args["amount"].(*big.Int).Uint64(), Address: args["to"].(common.Address)})
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
		tx := wallet_mock.NewObscuroTransferTx(fromWallet, to, obscurocommon.RndBtw(1, 500), m.l2NodeClients[0])
		// fmt.Printf("Injected transfer tx: %d\n", obscurocommon.ShortHash(tx.Hash()))
		signedTx := wallet_mock.SignTx(tx, fromWallet.Key.PrivateKey)
		encryptedTx := core.EncryptTx(signedTx)
		m.stats.Transfer()

		err := (*m.rndL2NodeClient()).Call(nil, obscuroclient.RPCSendTransactionEncrypted, encryptedTx)
		if err != nil {
			log.Info("Failed to issue transfer via RPC.")
			continue
		}

		go m.trackTransferL2Tx(*signedTx)
	}
}

// issueRandomDeposits creates and issues a number of transactions proportional to the simulation time, such that they can be processed
// Generates L1 common.DepositTx transactions
func (m *TransactionInjector) issueRandomL1Deposits() {
	for ; atomic.LoadInt32(m.interruptRun) == 0; time.Sleep(obscurocommon.RndBtwTime(m.avgBlockDuration, m.avgBlockDuration*2)) {
		v := obscurocommon.RndBtw(1, 100)
		txData := obscurocommon.L1TxData{
			TxType: obscurocommon.DepositTx,
			Amount: v,
			Dest:   rndWallet(m.wallets).Address,
		}
		// fmt.Printf("Injected l1 deposit tx: %v\n", txData)
		m.rndL1Node().BroadcastTx(&txData)
		m.stats.Deposit(v)
		go m.trackL1Tx(txData)
	}
}

// issueRandomWithdrawals creates and issues a number of transactions proportional to the simulation time, such that they can be processed
// Generates L2 enclave2.WithdrawalTx transactions
func (m *TransactionInjector) issueRandomWithdrawals() {
	for ; atomic.LoadInt32(m.interruptRun) == 0; time.Sleep(obscurocommon.RndBtwTime(m.avgBlockDuration, m.avgBlockDuration*2)) {
		v := obscurocommon.RndBtw(1, 100)
		wallet := rndWallet(m.wallets)
		tx := wallet_mock.NewObscuroWithdrawalTx(v, wallet, m.l2NodeClients[0])
		// fmt.Printf("Injected withdrawal tx: %d\n", obscurocommon.ShortHash(tx.Hash()))
		signedTx := wallet_mock.SignTx(tx, wallet.Key.PrivateKey)
		encryptedTx := core.EncryptTx(signedTx)

		err := (*m.rndL2NodeClient()).Call(nil, obscuroclient.RPCSendTransactionEncrypted, encryptedTx)
		if err != nil {
			log.Info("Failed to issue withdrawal via RPC.")
			continue
		}

		m.stats.Withdrawal(v)
		go m.trackWithdrawalL2Tx(*signedTx)
	}
}

// issueInvalidWithdrawals creates and issues a number of invalidly-signed L2 withdrawal transactions proportional to the simulation time.
// These transactions should be rejected by the nodes, and thus we expect them not to show up in the simulation withdrawal checks.
func (m *TransactionInjector) issueInvalidWithdrawals() {
	for ; atomic.LoadInt32(m.interruptRun) == 0; time.Sleep(obscurocommon.RndBtwTime(m.avgBlockDuration/4, m.avgBlockDuration)) {
		fromWallet := rndWallet(m.wallets)
		tx := wallet_mock.NewObscuroWithdrawalTx(obscurocommon.RndBtw(1, 100), fromWallet, m.l2NodeClients[0])
		signedTx := createInvalidSignature(tx, fromWallet, m.l2NodeClients[0])
		encryptedTx := core.EncryptTx(signedTx)

		err := (*m.rndL2NodeClient()).Call(nil, obscuroclient.RPCSendTransactionEncrypted, encryptedTx)
		if err != nil {
			log.Info("Failed to issue withdrawal via RPC.")
			continue
		}
	}
}

// Uses one of three approaches to create an invalidly-signed transaction.
func createInvalidSignature(tx *nodecommon.L2Tx, fromWallet *wallet_mock.Wallet, client *obscuroclient.Client) *nodecommon.L2Tx {
	i := rand.Intn(3) //nolint:gosec
	switch i {
	case 0: // We sign the transaction with a bad signer.
		incorrectChainID := int64(evm.ChainID + 1)
		signer := types.NewLondonSigner(big.NewInt(incorrectChainID))
		signedTx, _ := types.SignTx(tx, signer, fromWallet.Key.PrivateKey)
		return signedTx

	case 1: // We do not sign the transaction.
		return tx

	case 2: // We modify the transaction after signing.
		// We create a new transaction, as we need access to the transaction's encapsulated transaction data.
		newTx := wallet_mock.NewObscuroWithdrawalTx(obscurocommon.RndBtw(1, 100), fromWallet, client)
		wallet_mock.SignTx(newTx, fromWallet.Key.PrivateKey)
		// After signing the transaction, we create a new transaction based on the transaction data, breaking the signature.
		return wallet_mock.NewObscuroWithdrawalTx(obscurocommon.RndBtw(1, 100), fromWallet, client)
	}
	panic("Expected i to be in the range [0,2).")
}

func rndWallet(wallets []*wallet_mock.Wallet) *wallet_mock.Wallet {
	return wallets[rand.Intn(len(wallets))] //nolint:gosec
}

func (m *TransactionInjector) rndL1Node() ethclient.EthClient {
	return m.l1Nodes[rand.Intn(len(m.l1Nodes))] //nolint:gosec
}

func (m *TransactionInjector) rndL2NodeClient() *obscuroclient.Client {
	return m.l2NodeClients[rand.Intn(len(m.l2NodeClients))] //nolint:gosec
}
