package simulation

import (
	"context"
	"fmt"
	"math/big"
	"math/rand"
	"strings"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ten-protocol/go-ten/contracts/generated/ManagementContract"
	"github.com/ten-protocol/go-ten/contracts/generated/MessageBus"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/enclave/crosschain"
	"github.com/ten-protocol/go-ten/go/ethadapter/erc20contractlib"
	"github.com/ten-protocol/go-ten/go/ethadapter/mgmtcontractlib"
	"github.com/ten-protocol/go-ten/go/host/rpc/clientapi"
	"github.com/ten-protocol/go-ten/go/wallet"
	"github.com/ten-protocol/go-ten/integration"
	"github.com/ten-protocol/go-ten/integration/common/testlog"
	"github.com/ten-protocol/go-ten/integration/datagenerator"
	"github.com/ten-protocol/go-ten/integration/simulation/network"
	"github.com/ten-protocol/go-ten/integration/simulation/params"
	"golang.org/x/sync/errgroup"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
	testcommon "github.com/ten-protocol/go-ten/integration/common"
	simstats "github.com/ten-protocol/go-ten/integration/simulation/stats"
)

const (
	nonceTimeoutMillis = 30000 // The timeout in millis to wait for an updated nonce for a wallet.
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

	// The number of transactions of each type to issue, or 0 for unlimited transactions
	txsToIssue int

	// context for the transaction injector so in-flight requests can be cancelled gracefully
	ctx context.Context

	params *params.SimParams

	logger gethlog.Logger
}

// NewTransactionInjector returns a transaction manager with a given number of obsWallets
func NewTransactionInjector(
	avgBlockDuration time.Duration,
	stats *simstats.Stats,
	rpcHandles *network.RPCHandles,
	wallets *params.SimWallets,
	mgmtContractAddr *gethcommon.Address,
	mgmtContractLib mgmtcontractlib.MgmtContractLib,
	erc20ContractLib erc20contractlib.ERC20ContractLib,
	txsToIssue int,
	params *params.SimParams,
) *TransactionInjector {
	interrupt := int32(0)

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
		txsToIssue:       txsToIssue,
		params:           params,
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

	// in mem sim does not support the contract libraries required
	// to do complex bridge transactions
	if !ti.params.IsInMem {
		wg.Go(func() error {
			ti.bridgeRandomGasTransfers()
			return nil
		})
	}

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
		tenClient := ti.rpcHandles.TenWalletRndClient(fromWallet)
		// We avoid transfers to self, unless there is only a single L2 wallet.
		for len(ti.wallets.SimObsWallets) > 1 && fromWallet.Address().Hex() == toWallet.Address().Hex() {
			toWallet = ti.rndObsWallet()
		}
		toWalletAddr := toWallet.Address()
		txData := &types.LegacyTx{
			Nonce:    fromWallet.GetNonceAndIncrement(),
			Value:    big.NewInt(int64(testcommon.RndBtw(1, 100))),
			Gas:      uint64(50_000),
			GasPrice: gethcommon.Big1,
			To:       &toWalletAddr,
		}

		tx := tenClient.EstimateGasAndGasPrice(txData)
		signedTx, err := fromWallet.SignTransaction(tx)
		if err != nil {
			panic(err)
		}
		ti.logger.Info("Native value transfer transaction injected into L2.", log.TxKey, signedTx.Hash(), "fromAddress", fromWallet.Address(), "toAddress", toWallet.Address())

		ti.stats.NativeTransfer()

		err = tenClient.SendTransaction(ti.ctx, signedTx)
		if err != nil {
			ti.logger.Info("Failed to issue transfer via RPC.", log.ErrKey, err)
			continue
		}

		// todo (@pedro) - retrieve receipt

		go ti.TxTracker.trackNativeValueTransferL2Tx(signedTx)
		sleepRndBtw(ti.avgBlockDuration/10, ti.avgBlockDuration/4)
	}
}

// issueRandomTransfers creates and issues a number of L2 transfer transactions proportional to the simulation time, such that they can be processed
func (ti *TransactionInjector) issueRandomTransfers() {
	for txCounter := 0; ti.shouldKeepIssuing(txCounter); txCounter++ {
		fromWallet := ti.rndObsWallet()
		toWallet := ti.rndObsWallet()
		tenClient := ti.rpcHandles.TenWalletRndClient(fromWallet)
		// We avoid transfers to self, unless there is only a single L2 wallet.
		for len(ti.wallets.SimObsWallets) > 1 && fromWallet.Address().Hex() == toWallet.Address().Hex() {
			toWallet = ti.rndObsWallet()
		}
		tx := ti.newTenTransferTx(fromWallet, toWallet.Address(), testcommon.RndBtw(1, 500), testcommon.HOC)
		tx = tenClient.EstimateGasAndGasPrice(tx)
		signedTx, err := fromWallet.SignTransaction(tx)
		if err != nil {
			panic(err)
		}
		ti.logger.Info("Transfer transaction injected into L2.", log.TxKey, signedTx.Hash(), "fromAddress", fromWallet.Address(), "toAddress", toWallet.Address())

		ti.stats.Transfer()

		err = tenClient.SendTransaction(ti.ctx, signedTx)
		if err != nil {
			ti.logger.Info("Failed to issue transfer via RPC.", log.ErrKey, err)
		}

		// todo (@pedro) - retrieve receipt

		go ti.TxTracker.trackTransferL2Tx(signedTx)
		sleepRndBtw(ti.avgBlockDuration/100, ti.avgBlockDuration/20)
	}
}

func (ti *TransactionInjector) bridgeRandomGasTransfers() {
	gasWallet := ti.wallets.GasBridgeWallet

	ethClient := ti.rpcHandles.RndEthClient()

	mgmtCtr, err := ManagementContract.NewManagementContract(*ti.mgmtContractAddr, ethClient.EthClient())
	if err != nil {
		panic(err)
	}
	busAddr, err := mgmtCtr.MessageBus(&bind.CallOpts{})
	if err != nil {
		panic(err)
	}

	for txCounter := 0; ti.shouldKeepIssuing(txCounter); txCounter++ {
		ethClient = ti.rpcHandles.RndEthClient()

		busCtr, err := MessageBus.NewMessageBus(busAddr, ethClient.EthClient())
		if err != nil {
			panic(err)
		}

		opts, err := bind.NewKeyedTransactorWithChainID(gasWallet.PrivateKey(), gasWallet.ChainID())
		if err != nil {
			panic(err)
		}

		receiverWallet := datagenerator.RandomWallet(ti.rndObsWallet().ChainID().Int64())
		amount := big.NewInt(0).SetUint64(testcommon.RndBtw(500, 100_000))
		opts.Value = big.NewInt(0).Set(amount)

		tx, err := busCtr.SendValueToL2(opts, receiverWallet.Address(), amount)
		if err != nil {
			panic(err)
		}

		go ti.TxTracker.trackGasBridgingTx(tx, receiverWallet)

		sleepRndBtw(ti.avgBlockDuration/3, ti.avgBlockDuration)
	}
}

// issueRandomDeposits creates and issues a number of transactions proportional to the simulation time, such that they can be processed
func (ti *TransactionInjector) issueRandomDeposits() {
	// todo (@stefan) - this implementation transfers from the hoc and poc owner contracts
	// a better implementation should use the bridge
	for txCounter := 0; ti.shouldKeepIssuing(txCounter); txCounter++ {
		fromWalletToken := testcommon.HOC
		if txCounter%2 == 0 {
			fromWalletToken = testcommon.POC
		}
		fromWallet := ti.wallets.Tokens[fromWalletToken].L2Owner
		toWallet := ti.rndObsWallet()
		tenClient := ti.rpcHandles.TenWalletRndClient(fromWallet)
		v := testcommon.RndBtw(500, 2000)
		txData := ti.newTenTransferTx(fromWallet, toWallet.Address(), v, fromWalletToken)
		tx := tenClient.EstimateGasAndGasPrice(txData)
		signedTx, err := fromWallet.SignTransaction(tx)
		if err != nil {
			panic(err)
		}
		ti.logger.Info("Deposit transaction injected into L2.", log.TxKey, signedTx.Hash(), "fromAddress", fromWallet.Address(), "toAddress", toWallet.Address())

		ti.stats.Deposit(big.NewInt(int64(v)))

		err = tenClient.SendTransaction(ti.ctx, signedTx)
		if err != nil {
			ti.logger.Info("Failed to issue deposit via RPC.", log.ErrKey, err)
		} else {
			go ti.TxTracker.trackTransferL2Tx(signedTx)
		}
		// todo (@pedro) - retrieve receipt

		sleepRndBtw(ti.avgBlockDuration/3, ti.avgBlockDuration)
	}
	// todo (@stefan) - rework this when old contract deployer is phased out?
}

func (ti *TransactionInjector) awaitAndFinalizeWithdrawal(tx *types.Transaction, fromWallet wallet.Wallet) {
	err := testcommon.AwaitReceipt(ti.ctx, ti.rpcHandles.TenWalletRndClient(fromWallet), tx.Hash(), 45*time.Second)
	if err != nil {
		ti.logger.Error("Failed to await receipt for withdrawal transaction", log.ErrKey, err)
		return
	}

	receipt, err := ti.rpcHandles.TenWalletRndClient(fromWallet).TransactionReceipt(ti.ctx, tx.Hash())
	if err != nil {
		ti.logger.Error("Failed to retrieve receipt for withdrawal transaction", log.ErrKey, err)
		return
	}

	logs := make([]types.Log, len(receipt.Logs))
	for i, log := range receipt.Logs {
		logs[i] = *log
	}

	transfers, err := crosschain.ConvertLogsToValueTransfers(logs, crosschain.ValueTransferEventName, crosschain.MessageBusABI)
	if err != nil {
		panic(err)
	}

	vTransfers := crosschain.ValueTransfers(transfers)

	var proof clientapi.CrossChainProof
	for {
		proof, err = ti.rpcHandles.TenWalletRndClient(fromWallet).GetCrossChainProof(ti.ctx, "v", vTransfers.ForMerkleTree()[0][1].(gethcommon.Hash))
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				ti.logger.Info("Proof not found, retrying...", log.ErrKey, err)
				time.Sleep(1 * time.Second)
				continue
			}
			if strings.Contains(err.Error(), "database closed") {
				ti.logger.Info("Database closed, test over", log.ErrKey, err)
				return
			}
			panic(fmt.Errorf("unable to get proof for value transfer. cause: %w", err))
		}
		break
	}

	if len(proof.Proof) == 0 {
		return
	}

	proofBytes := [][]byte{}
	if err := rlp.DecodeBytes(proof.Proof, &proofBytes); err != nil {
		panic("unable to decode proof")
	}

	// In mem sim does not support the l1 interaction required for the rest of the function.
	if ti.mgmtContractLib.IsMock() {
		return
	}

	mCtr, err := ManagementContract.NewManagementContract(*ti.mgmtContractAddr, ti.rpcHandles.RndEthClient().EthClient())
	if err != nil {
		panic(err)
	}

	opts, err := bind.NewKeyedTransactorWithChainID(ti.wallets.GasWithdrawalWallet.PrivateKey(), ti.wallets.GasWithdrawalWallet.ChainID())
	if err != nil {
		panic(err)
	}

	proof32 := make([][32]byte, 0)
	for i := 0; i < len(proofBytes); i++ {
		proof32 = append(proof32, [32]byte(proofBytes[i][0:32]))
	}

	time.Sleep(20 * time.Second)

	oldBalance, err := ti.rpcHandles.RndEthClient().BalanceAt(vTransfers[0].Receiver, nil)
	if err != nil {
		ti.logger.Error("Failed to retrieve balance of receiver", log.ErrKey, err)
		return
	}

	withdrawalTx, err := mCtr.ExtractNativeValue(
		opts,
		ManagementContract.StructsValueTransferMessage(vTransfers[0]),
		proof32,
		proof.Root,
	)
	if err != nil {
		ti.logger.Error("Failed to extract value transfer from L2", log.ErrKey, err)
		return
	}

	receipt, err = testcommon.AwaitReceiptEth(ti.ctx, ti.rpcHandles.RndEthClient().EthClient(), withdrawalTx.Hash(), 30*time.Second)
	if err != nil {
		ti.logger.Error("Failed to await receipt for withdrawal transaction", log.ErrKey, err)
		return
	}

	if receipt.Status != 1 {
		ti.logger.Error("Withdrawal transaction failed", log.TxKey, withdrawalTx.Hash())
		return
	}

	newBalance, err := ti.rpcHandles.RndEthClient().BalanceAt(vTransfers[0].Receiver, nil)
	if err != nil {
		ti.logger.Error("Failed to retrieve balance of receiver", log.ErrKey, err)
		return
	}

	if newBalance.Sub(newBalance, oldBalance).Cmp(vTransfers[0].Amount) != 0 {
		ti.logger.Error("Balance of receiver did not increase by the expected amount", "expected", vTransfers[0].Amount, "actual", newBalance.Sub(newBalance, oldBalance))
		return
	}

	ti.logger.Info("Successful bridge withdrawal", log.TxKey, withdrawalTx.Hash())
}

// issueRandomWithdrawals creates and issues a number of transactions proportional to the simulation time, such that they can be processed
func (ti *TransactionInjector) issueRandomWithdrawals() {
	cfg, err := ti.rpcHandles.TenWalletRndClient(ti.wallets.L2FaucetWallet).GetConfig()
	if err != nil {
		panic(err)
	}
	msgBusAddr := cfg.L2MessageBusAddress

	for txCounter := 0; ti.shouldKeepIssuing(txCounter); txCounter++ {
		fromWallet := ti.rndObsWallet()
		client := ti.rpcHandles.TenWalletRndClient(fromWallet)
		price, err := client.GasPrice(ti.ctx)
		if err != nil {
			ti.logger.Error("unable to estimate gas price", log.ErrKey, err)
			continue
		}

		tx := &types.LegacyTx{
			Nonce:    fromWallet.GetNonceAndIncrement(),
			Value:    gethcommon.Big1,
			Gas:      uint64(1_000_000),
			GasPrice: price,
			Data:     nil,
			To:       &msgBusAddr,
		}
		signedTx, err := fromWallet.SignTransaction(tx)
		if err != nil {
			ti.logger.Error("[CrossChain] unable to sign withdrawal transaction", log.ErrKey, err)
			continue
		}

		err = client.SendTransaction(ti.ctx, signedTx)
		if err != nil {
			ti.logger.Error("[CrossChain] unable to send withdrawal transaction", log.ErrKey, err)
		}

		go ti.TxTracker.trackWithdrawalFromL2(signedTx)

		ti.logger.Info("[CrossChain] successful withdrawal tx", log.TxKey, signedTx.Hash())

		go ti.awaitAndFinalizeWithdrawal(signedTx, fromWallet)

		time.Sleep(testcommon.RndBtwTime(ti.avgBlockDuration/4, ti.avgBlockDuration))
	}
}

// issueInvalidL2Txs creates and issues invalidly-signed L2 transactions proportional to the simulation time.
// These transactions should be rejected by the nodes, and thus we expect them to not affect the simulation
func (ti *TransactionInjector) issueInvalidL2Txs() {
	// todo (@tudor) - also issue transactions with insufficient gas
	for txCounter := 0; ti.shouldKeepIssuing(txCounter); txCounter++ {
		fromWallet := ti.rndObsWallet()
		toWallet := ti.rndObsWallet()
		// We avoid transfers to self, unless there is only a single L2 wallet.
		for len(ti.wallets.SimObsWallets) > 1 && fromWallet.Address().Hex() == toWallet.Address().Hex() {
			toWallet = ti.rndObsWallet()
		}
		txData := ti.newCustomTenWithdrawalTx(testcommon.RndBtw(1, 100))

		tx := ti.rpcHandles.TenWalletRndClient(fromWallet).EstimateGasAndGasPrice(txData)
		signedTx := ti.createInvalidSignage(tx, fromWallet)

		err := ti.rpcHandles.TenWalletRndClient(fromWallet).SendTransaction(ti.ctx, signedTx)
		if err != nil {
			ti.logger.Info("Failed to issue withdrawal via RPC. ", log.ErrKey, err)
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

func (ti *TransactionInjector) newTenTransferTx(from wallet.Wallet, dest gethcommon.Address, amount uint64, ercType testcommon.ERC20) types.TxData {
	data := erc20contractlib.CreateTransferTxData(dest, common.ValueInWei(big.NewInt(int64(amount))))
	return ti.newTx(data, from.GetNonceAndIncrement(), ercType)
}

func (ti *TransactionInjector) newCustomTenWithdrawalTx(amount uint64) types.TxData {
	transferERC20data := erc20contractlib.CreateTransferTxData(testcommon.BridgeAddress, common.ValueInWei(big.NewInt(int64(amount))))
	return ti.newTx(transferERC20data, 1, testcommon.HOC)
}

func (ti *TransactionInjector) newTx(data []byte, nonce uint64, ercType testcommon.ERC20) types.TxData {
	return &types.LegacyTx{
		Nonce:    nonce,
		Value:    gethcommon.Big0,
		Gas:      uint64(1_000_000),
		GasPrice: gethcommon.Big1,
		Data:     data,
		To:       ti.wallets.Tokens[ercType].L2ContractAddress,
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
