package simulation

import (
	"context"
	"fmt"
	"math/big"
	"math/rand"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/ethadapter/erc20contractlib"
	"github.com/obscuronet/go-obscuro/go/ethadapter/mgmtcontractlib"
	"github.com/obscuronet/go-obscuro/go/wallet"
	"github.com/obscuronet/go-obscuro/integration"
	"github.com/obscuronet/go-obscuro/integration/common/testlog"
	"github.com/obscuronet/go-obscuro/integration/simulation/network"
	"github.com/obscuronet/go-obscuro/integration/simulation/params"
	"golang.org/x/sync/errgroup"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
	testcommon "github.com/obscuronet/go-obscuro/integration/common"
	simstats "github.com/obscuronet/go-obscuro/integration/simulation/stats"
)

const (
	nonceTimeoutMillis = 30000 // The timeout in millis to wait for an updated nonce for a wallet.

	// EnclavePublicKeyHex is the public key of the enclave.
	// TODO - Retrieve this key from the management contract instead.
	EnclavePublicKeyHex = "034d3b7e63a8bcd532ee3d1d6ecad9d67fca7821981a044551f0f0cbec74d0bc5e"
)

// TransactionInjector is a structure that generates, issues and tracks transactions
type TransactionInjector struct {
	// counters
	TxTracker *txInjectorTracker
	stats     *simstats.Stats

	// settings
	avgBlockDuration time.Duration

	// wallets
	wallets *params.SimWallets

	// connections
	rpcHandles *network.RPCHandles

	// addrs and libs
	mgmtContractAddr *gethcommon.Address
	mgmtContractLib  mgmtcontractlib.MgmtContractLib
	erc20ContractLib erc20contractlib.ERC20ContractLib

	// controls
	interruptRun     *int32
	fullyStoppedChan chan bool

	enclavePublicKey *ecies.PublicKey

	// The number of transactions of each type to issue, or 0 for unlimited transactions
	txsToIssue int

	// context for the transaction injector so in-flight requests can be cancelled gracefully
	ctx context.Context

	logger gethlog.Logger
}

// NewTransactionInjector returns a transaction manager with a given number of obsWallets
// todo Add methods that generate deterministic scenarios
func NewTransactionInjector(
	avgBlockDuration time.Duration,
	stats *simstats.Stats,
	rpcHandles *network.RPCHandles,
	wallets *params.SimWallets,
	mgmtContractAddr *gethcommon.Address,
	mgmtContractLib mgmtcontractlib.MgmtContractLib,
	erc20ContractLib erc20contractlib.ERC20ContractLib,
	txsToIssue int,
) *TransactionInjector {
	interrupt := int32(0)

	// We retrieve the enclave public key to encrypt transactions.
	enclavePublicKey, err := crypto.DecompressPubkey(gethcommon.Hex2Bytes(EnclavePublicKeyHex))
	if err != nil {
		panic(fmt.Errorf("could not decompress enclave public key from hex. Cause: %w", err))
	}
	enclavePublicKeyEcies := ecies.ImportECDSAPublic(enclavePublicKey)

	return &TransactionInjector{
		avgBlockDuration: avgBlockDuration,
		stats:            stats,
		rpcHandles:       rpcHandles,
		interruptRun:     &interrupt,
		fullyStoppedChan: make(chan bool, 1),
		mgmtContractAddr: mgmtContractAddr,
		mgmtContractLib:  mgmtContractLib,
		erc20ContractLib: erc20ContractLib,
		wallets:          wallets,
		TxTracker:        newCounter(),
		enclavePublicKey: enclavePublicKeyEcies,
		txsToIssue:       txsToIssue,
		ctx:              context.Background(), // for now we create a new context here, should allow it to be passed in
		logger:           testlog.Logger().New(log.CmpKey, log.TxInjectCmp),
	}
}

// Start begins the execution on the TransactionInjector
// Deposits an initial balance in to each wallet
// Generates and issues L1 and L2 transactions to the network
func (ti *TransactionInjector) Start() {
	var wg errgroup.Group
	wg.Go(func() error {
		ti.issueRandomDeposits()
		return nil
	})

	wg.Go(func() error {
		ti.issueRandomWithdrawals()
		return nil
	})

	wg.Go(func() error {
		ti.issueRandomTransfers()
		return nil
	})

	wg.Go(func() error {
		ti.issueRandomValueTransfers()
		return nil
	})

	wg.Go(func() error {
		ti.issueInvalidL2Txs()
		return nil
	})

	_ = wg.Wait() // future proofing to return errors
	ti.fullyStoppedChan <- true
}

func (ti *TransactionInjector) Stop() {
	atomic.StoreInt32(ti.interruptRun, 1)
	for range ti.fullyStoppedChan {
		ti.logger.Info("TransactionInjector stopped successfully")
		return
	}
}

// issueRandomValueTransfers creates and issues a number of L2 value transfer transactions proportional to the simulation time, such that they can be processed
func (ti *TransactionInjector) issueRandomValueTransfers() {
	for txCounter := 0; ti.shouldKeepIssuing(txCounter); txCounter++ {
		fromWallet := ti.rndObsWallet()
		toWallet := ti.rndObsWallet()
		obscuroClient := ti.rpcHandles.ObscuroWalletRndClient(fromWallet)
		// We avoid transfers to self, unless there is only a single L2 wallet.
		for len(ti.wallets.SimObsWallets) > 1 && fromWallet.Address().Hex() == toWallet.Address().Hex() {
			toWallet = ti.rndObsWallet()
		}
		toWalletAddr := toWallet.Address()
		txData := &types.LegacyTx{
			Nonce:    fromWallet.GetNonceAndIncrement(),
			Value:    big.NewInt(int64(testcommon.RndBtw(1, 500))),
			Gas:      uint64(1_000_000),
			GasPrice: gethcommon.Big1,
			To:       &toWalletAddr,
		}

		tx := obscuroClient.EstimateGasAndGasPrice(txData)
		signedTx, err := fromWallet.SignTransaction(tx)
		if err != nil {
			panic(err)
		}
		ti.logger.Info(fmt.Sprintf(
			"Transfer transaction injected into L2. Hash: %d. From address: %d. To address: %d",
			common.ShortHash(signedTx.Hash()),
			common.ShortAddress(fromWallet.Address()),
			common.ShortAddress(toWallet.Address()),
		))

		ti.stats.Transfer()

		txHash, err := obscuroClient.SendTransaction(ti.ctx, signedTx)
		if err != nil {
			ti.logger.Info("Failed to issue transfer via RPC.", log.ErrKey, err)
			continue
		}

		if *txHash != signedTx.Hash() {
			ti.logger.Error("The hash of the submitted transaction does not match the hash coming back!")
			continue
		}

		// todo - retrieve receipt

		go ti.TxTracker.trackNativeValueTransferL2Tx(signedTx)
		sleepRndBtw(ti.avgBlockDuration/4, ti.avgBlockDuration)
	}
}

// issueRandomTransfers creates and issues a number of L2 transfer transactions proportional to the simulation time, such that they can be processed
func (ti *TransactionInjector) issueRandomTransfers() {
	for txCounter := 0; ti.shouldKeepIssuing(txCounter); txCounter++ {
		fromWallet := ti.rndObsWallet()
		toWallet := ti.rndObsWallet()
		obscuroClient := ti.rpcHandles.ObscuroWalletRndClient(fromWallet)
		// We avoid transfers to self, unless there is only a single L2 wallet.
		for len(ti.wallets.SimObsWallets) > 1 && fromWallet.Address().Hex() == toWallet.Address().Hex() {
			toWallet = ti.rndObsWallet()
		}
		tx := ti.newObscuroTransferTx(fromWallet, toWallet.Address(), testcommon.RndBtw(1, 500))
		tx = obscuroClient.EstimateGasAndGasPrice(tx)
		signedTx, err := fromWallet.SignTransaction(tx)
		if err != nil {
			panic(err)
		}
		ti.logger.Info(fmt.Sprintf(
			"Transfer transaction injected into L2. Hash: %d. From address: %d. To address: %d",
			common.ShortHash(signedTx.Hash()),
			common.ShortAddress(fromWallet.Address()),
			common.ShortAddress(toWallet.Address()),
		))

		ti.stats.Transfer()

		_, err = obscuroClient.SendTransaction(ti.ctx, signedTx)
		if err != nil {
			ti.logger.Info("Failed to issue transfer via RPC.", log.ErrKey, err)
			continue
		}

		// todo - retrieve receipt

		go ti.TxTracker.trackTransferL2Tx(signedTx)
		sleepRndBtw(ti.avgBlockDuration/10, ti.avgBlockDuration/4)
	}
}

// issueRandomDeposits creates and issues a number of transactions proportional to the simulation time, such that they can be processed
func (ti *TransactionInjector) issueRandomDeposits() {
	// TODO - this implementation transfers from the hoc and poc owner contracts
	// a better implementation should use the bridge
	fromWalletHoc := ti.wallets.Tokens[testcommon.HOC].L2Owner
	fromWalletPoc := ti.wallets.Tokens[testcommon.POC].L2Owner

	for txCounter := 0; ti.shouldKeepIssuing(txCounter); txCounter++ {
		fromWallet := fromWalletHoc
		if txCounter%2 == 0 {
			fromWallet = fromWalletPoc
		}
		toWallet := ti.rndObsWallet()
		obscuroClient := ti.rpcHandles.ObscuroWalletRndClient(toWallet)
		v := testcommon.RndBtw(500, 2000)
		tx := ti.newObscuroTransferTx(fromWallet, toWallet.Address(), v)
		tx = obscuroClient.EstimateGasAndGasPrice(tx)
		signedTx, err := fromWallet.SignTransaction(tx)
		if err != nil {
			panic(err)
		}
		ti.logger.Info(fmt.Sprintf(
			"Deposit  transaction injected into L2. Hash: %d. From address: %d. To address: %d",
			common.ShortHash(signedTx.Hash()),
			common.ShortAddress(fromWallet.Address()),
			common.ShortAddress(toWallet.Address()),
		))

		ti.stats.Deposit(big.NewInt(int64(v)))

		err = obscuroClient.SendTransaction(ti.ctx, signedTx)
		if err != nil {
			ti.logger.Info("Failed to issue deposit via RPC.", log.ErrKey, err)
			continue
		}

		// todo - retrieve receipt

		go ti.TxTracker.trackTransferL2Tx(signedTx)
		sleepRndBtw(ti.avgBlockDuration/3, ti.avgBlockDuration)
	}
	// TODO: Rework this when old contract deployer is phased out?
}

// issueRandomWithdrawals creates and issues a number of transactions proportional to the simulation time, such that they can be processed
func (ti *TransactionInjector) issueRandomWithdrawals() {
	// TODO: Rework this when old contract deployer is phased out?
}

// issueInvalidL2Txs creates and issues invalidly-signed L2 transactions proportional to the simulation time.
// These transactions should be rejected by the nodes, and thus we expect them to not affect the simulation
func (ti *TransactionInjector) issueInvalidL2Txs() {
	// todo - also issue transactions with insufficient gas
	for txCounter := 0; ti.shouldKeepIssuing(txCounter); txCounter++ {
		fromWallet := ti.rndObsWallet()
		toWallet := ti.rndObsWallet()
		// We avoid transfers to self, unless there is only a single L2 wallet.
		for len(ti.wallets.SimObsWallets) > 1 && fromWallet.Address().Hex() == toWallet.Address().Hex() {
			toWallet = ti.rndObsWallet()
		}
		tx := ti.newCustomObscuroWithdrawalTx(testcommon.RndBtw(1, 100))

		signedTx := ti.createInvalidSignage(tx, fromWallet)

		_, err := ti.rpcHandles.ObscuroWalletRndClient(fromWallet).SendTransaction(ti.ctx, signedTx)
		if err != nil {
			ti.logger.Warn("Failed to issue withdrawal via RPC. ", log.ErrKey, err)
		}
		time.Sleep(testcommon.RndBtwTime(ti.avgBlockDuration/4, ti.avgBlockDuration))
	}
}

// Uses one of the approaches to create an invalidly-signed transaction.
func (ti *TransactionInjector) createInvalidSignage(tx types.TxData, w wallet.Wallet) *types.Transaction {
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

func (ti *TransactionInjector) rndObsWallet() wallet.Wallet {
	return ti.wallets.SimObsWallets[rand.Intn(len(ti.wallets.SimObsWallets))] //nolint:gosec
}

func (ti *TransactionInjector) newObscuroTransferTx(from wallet.Wallet, dest gethcommon.Address, amount uint64) types.TxData {
	data := erc20contractlib.CreateTransferTxData(dest, common.ValueInWei(big.NewInt(int64(amount))))
	return ti.newTx(data, from.GetNonceAndIncrement())
}

func (ti *TransactionInjector) newCustomObscuroWithdrawalTx(amount uint64) types.TxData {
	transferERC20data := erc20contractlib.CreateTransferTxData(testcommon.BridgeAddress, common.ValueInWei(big.NewInt(int64(amount))))
	return ti.newTx(transferERC20data, 1)
}

func (ti *TransactionInjector) newTx(data []byte, nonce uint64) types.TxData {
	// todo - reenable this logic when the nonce logic has been replaced by receipt confirmation
	//max := big.NewInt(1_000_000_000_000_000_000)
	//if nonce%3 == 0 {
	//	value = max
	//}

	return &types.LegacyTx{
		Nonce:    nonce,
		Value:    gethcommon.Big0,
		Gas:      uint64(1_000_000),
		GasPrice: gethcommon.Big1,
		Data:     data,
		To:       ti.wallets.Tokens[testcommon.HOC].L2ContractAddress,
	}
}

// Indicates whether to keep issuing transactions, or halt.
func (ti *TransactionInjector) shouldKeepIssuing(txCounter int) bool {
	isInterrupted := atomic.LoadInt32(ti.interruptRun) != 0

	// 0 is a special value indicating we should only stop issuing transactions when interrupted.
	if ti.txsToIssue == 0 {
		return !isInterrupted
	}

	return !isInterrupted && txCounter < ti.txsToIssue
}
