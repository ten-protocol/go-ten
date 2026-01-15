package simulation

import (
	"context"
	"fmt"
	"math/big"
	"math/rand"
	"strings"
	"sync/atomic"
	"time"

	"github.com/ten-protocol/go-ten/contracts/generated/CrossChainMessenger"
	"github.com/ten-protocol/go-ten/contracts/generated/EthereumBridge"
	"github.com/ten-protocol/go-ten/contracts/generated/TenBridge"
	"github.com/ten-protocol/go-ten/go/ethadapter"
	"github.com/ten-protocol/go-ten/go/ethadapter/contractlib"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/enclave/crosschain"
	"github.com/ten-protocol/go-ten/go/ethadapter/erc20contractlib"
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
	"github.com/ten-protocol/go-ten/go/obsclient"
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
	contractRegistryLib contractlib.ContractRegistryLib
	erc20ContractLib    erc20contractlib.ERC20ContractLib

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
	contractRegistryLib contractlib.ContractRegistryLib,
	erc20ContractLib erc20contractlib.ERC20ContractLib,
	txsToIssue int,
	params *params.SimParams,
) *TransactionInjector {
	interrupt := int32(0)

	return &TransactionInjector{
		avgBlockDuration:    avgBlockDuration,
		stats:               stats,
		rpcHandles:          rpcHandles,
		interruptRun:        &interrupt,
		fullyStoppedChan:    make(chan bool, 1),
		contractRegistryLib: contractRegistryLib,
		erc20ContractLib:    erc20ContractLib,
		wallets:             wallets,
		TxTracker:           newCounter(),
		txsToIssue:          txsToIssue,
		params:              params,
		ctx:                 context.Background(), // for now we create a new context here, should allow it to be passed in
		logger:              testlog.Logger().New(log.CmpKey, log.TxInjectCmp),
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

		// WETH bridging L1 -> L2
		wg.Go(func() error {
			ti.bridgeRandomWETHTransfers()
			return nil
		})

		// WETH bridging L2 -> L1
		wg.Go(func() error {
			ti.issueRandomWETHWithdrawals()
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
		price, err := tenClient.GasPrice(ti.ctx)
		if err != nil {
			panic(err)
		}
		price = new(big.Int).Mul(price, big.NewInt(2))
		// We avoid transfers to self, unless there is only a single L2 wallet.
		for len(ti.wallets.SimObsWallets) > 1 && fromWallet.Address().Hex() == toWallet.Address().Hex() {
			toWallet = ti.rndObsWallet()
		}
		toWalletAddr := toWallet.Address()
		txData := &types.LegacyTx{
			Nonce:    fromWallet.GetNonceAndIncrement(),
			Value:    big.NewInt(int64(testcommon.RndBtw(1, 100))),
			Gas:      uint64(10_000_000),
			GasPrice: price,
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

	for txCounter := 0; ti.shouldKeepIssuing(txCounter); txCounter++ {
		bridgeAddr := ti.params.L1TenData.BridgeAddress
		bridgeCtr, err := TenBridge.NewTenBridge(bridgeAddr, ti.rpcHandles.RndEthClient().EthClient())
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

		tx, err := bridgeCtr.SendNative(opts, receiverWallet.Address())
		if err != nil {
			panic(err)
		}

		go ti.TxTracker.trackGasBridgingTx(tx, receiverWallet)

		sleepRndBtw(ti.avgBlockDuration/3, ti.avgBlockDuration)
	}
}

// bridgeRandomWETHTransfers bridges WETH from L1 to L2 using TenBridge.SendERC20
// The flow is: wrap ETH to WETH on L1 -> approve TenBridge -> sendERC20(weth) -> receiver gets WETH on L2
func (ti *TransactionInjector) bridgeRandomWETHTransfers() {
	gasWallet := ti.wallets.GasBridgeWallet
	wethAddress := gethcommon.HexToAddress("0x1000000000000000000000000000000000000042")

	for txCounter := 0; ti.shouldKeepIssuing(txCounter); txCounter++ {
		ethClient := ti.rpcHandles.RndEthClient().EthClient()
		bridgeAddr := ti.params.L1TenData.BridgeAddress

		// Create random receiver and amount
		receiverWallet := datagenerator.RandomWallet(ti.rndObsWallet().ChainID().Int64())
		amount := big.NewInt(0).SetUint64(testcommon.RndBtw(500, 50_000))

		// Get gas price
		gasPrice, err := ethClient.SuggestGasPrice(ti.ctx)
		if err != nil {
			ti.logger.Error("[WETH Bridge] failed to get gas price", log.ErrKey, err)
			continue
		}
		gasPrice = new(big.Int).Mul(gasPrice, big.NewInt(2))

		// Fetch actual nonce from chain to avoid desync issues
		nonce, err := ethClient.PendingNonceAt(ti.ctx, gasWallet.Address())
		if err != nil {
			ti.logger.Error("[WETH Bridge] failed to get nonce", log.ErrKey, err)
			continue
		}

		// Step 1: Wrap ETH to WETH by sending ETH to WETH contract (triggers deposit())
		wrapTxData := &types.LegacyTx{
			Nonce:    nonce,
			To:       &wethAddress,
			Value:    amount,
			Gas:      uint64(100_000),
			GasPrice: gasPrice,
		}
		wrapTx, err := gasWallet.SignTransaction(wrapTxData)
		if err != nil {
			ti.logger.Error("[WETH Bridge] failed to sign wrap tx", log.ErrKey, err)
			continue
		}

		err = ethClient.SendTransaction(ti.ctx, wrapTx)
		if err != nil {
			ti.logger.Error("[WETH Bridge] failed to send wrap tx", log.ErrKey, err)
			continue
		}

		_, err = testcommon.AwaitReceiptEth(ti.ctx, ethClient, wrapTx.Hash(), 30*time.Second)
		if err != nil {
			ti.logger.Error("[WETH Bridge] wrap tx failed", log.ErrKey, err, log.TxKey, wrapTx.Hash())
			continue
		}
		ti.logger.Info("[WETH Bridge] wrapped ETH to WETH", log.TxKey, wrapTx.Hash(), "amount", amount)

		// Re-fetch nonce after wrap tx - other concurrent transactions may have used nonces
		nonce, err = ethClient.PendingNonceAt(ti.ctx, gasWallet.Address())
		if err != nil {
			ti.logger.Error("[WETH Bridge] failed to get nonce after wrap", log.ErrKey, err)
			continue
		}

		// Step 2: Approve TenBridge to spend WETH
		// approve(address spender, uint256 amount) selector: 0x095ea7b3
		approveData := make([]byte, 68)
		copy(approveData[0:4], []byte{0x09, 0x5e, 0xa7, 0xb3})
		copy(approveData[16:36], bridgeAddr.Bytes())
		amount.FillBytes(approveData[36:68])

		approveTxData := &types.LegacyTx{
			Nonce:    nonce,
			To:       &wethAddress,
			Value:    big.NewInt(0),
			Gas:      uint64(100_000),
			GasPrice: gasPrice,
			Data:     approveData,
		}
		approveTx, err := gasWallet.SignTransaction(approveTxData)
		if err != nil {
			ti.logger.Error("[WETH Bridge] failed to sign approve tx", log.ErrKey, err)
			continue
		}

		err = ethClient.SendTransaction(ti.ctx, approveTx)
		if err != nil {
			ti.logger.Error("[WETH Bridge] failed to send approve tx", log.ErrKey, err)
			continue
		}

		_, err = testcommon.AwaitReceiptEth(ti.ctx, ethClient, approveTx.Hash(), 30*time.Second)
		if err != nil {
			ti.logger.Error("[WETH Bridge] approve tx failed", log.ErrKey, err, log.TxKey, approveTx.Hash())
			continue
		}
		ti.logger.Info("[WETH Bridge] approved TenBridge to spend WETH", log.TxKey, approveTx.Hash())

		// Re-fetch nonce after approve tx
		nonce, err = ethClient.PendingNonceAt(ti.ctx, gasWallet.Address())
		if err != nil {
			ti.logger.Error("[WETH Bridge] failed to get nonce after approve", log.ErrKey, err)
			continue
		}

		// Step 3: Call TenBridge.SendERC20(wethAddress, amount, receiver)
		bridgeCtr, err := TenBridge.NewTenBridge(bridgeAddr, ethClient)
		if err != nil {
			ti.logger.Error("[WETH Bridge] failed to create TenBridge contract", log.ErrKey, err)
			continue
		}

		sendOpts, err := bind.NewKeyedTransactorWithChainID(gasWallet.PrivateKey(), gasWallet.ChainID())
		if err != nil {
			ti.logger.Error("[WETH Bridge] failed to create transactor for sendERC20", log.ErrKey, err)
			continue
		}
		sendOpts.GasLimit = uint64(500_000)
		sendOpts.Nonce = big.NewInt(int64(nonce))

		bridgeTx, err := bridgeCtr.SendERC20(sendOpts, wethAddress, amount, receiverWallet.Address())
		if err != nil {
			ti.logger.Error("[WETH Bridge] failed to call SendERC20", log.ErrKey, err)
			continue
		}

		ti.logger.Info("[WETH Bridge] L1->L2 WETH bridge tx sent", log.TxKey, bridgeTx.Hash(), "amount", amount, "receiver", receiverWallet.Address())
		go ti.awaitAndRelayL1ToL2Message(bridgeTx, amount, receiverWallet)

		sleepRndBtw(ti.avgBlockDuration/3, ti.avgBlockDuration)
	}
}

// awaitAndRelayL1ToL2Message waits for an L1→L2 bridge transaction to be processed,
// extracts the cross-chain message, and relays it on L2 to complete the bridge.
func (ti *TransactionInjector) awaitAndRelayL1ToL2Message(tx *types.Transaction, amount *big.Int, receiverWallet wallet.Wallet) {
	ethClient := ti.rpcHandles.RndEthClient().EthClient()

	// Step 1: Wait for L1 transaction receipt
	receipt, err := testcommon.AwaitReceiptEth(ti.ctx, ethClient, tx.Hash(), 45*time.Second)
	if err != nil {
		ti.logger.Error("[WETH Bridge L1->L2] failed to get L1 receipt", log.ErrKey, err, log.TxKey, tx.Hash())
		return
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		ti.logger.Error("[WETH Bridge L1->L2] L1 transaction failed", log.TxKey, tx.Hash())
		return
	}

	// Step 2: Extract LogMessagePublished events from L1 MessageBus
	l1MessageBusAddr := ti.params.L1TenData.MessageBusAddr
	filteredLogs, err := crosschain.FilterLogsFromReceipt(receipt, &l1MessageBusAddr, &ethadapter.CrossChainEventID)
	if err != nil {
		ti.logger.Error("[WETH Bridge L1->L2] failed to filter L1 logs", log.ErrKey, err)
		return
	}

	if len(filteredLogs) == 0 {
		ti.logger.Error("[WETH Bridge L1->L2] no cross-chain messages found in L1 tx")
		return
	}

	// Convert logs to CrossChainMessages
	messages, err := crosschain.ConvertLogsToMessages(filteredLogs, ethadapter.CrossChainEventName, ethadapter.MessageBusABI)
	if err != nil {
		ti.logger.Error("[WETH Bridge L1->L2] failed to convert logs to messages", log.ErrKey, err)
		return
	}

	// Find the receiveNativeWrapped message (Topic = TRANSFER = 0)
	// The sendNative message (Topic = VALUE = 2) is handled automatically by value transfer
	var targetMessage *common.CrossChainMessage
	for i := range messages {
		msg := &messages[i]
		// Topic 0 = TRANSFER, which is used for receiveNativeWrapped
		if msg.Topic == 0 {
			targetMessage = msg
			break
		}
	}

	if targetMessage == nil {
		ti.logger.Error("[WETH Bridge L1->L2] no TRANSFER message found for relay")
		return
	}

	ti.logger.Info("[WETH Bridge L1->L2] found cross-chain message to relay",
		"sender", targetMessage.Sender.Hex(),
		"sequence", targetMessage.Sequence,
		"topic", targetMessage.Topic)

	// Step 3: Wait for the message to be stored on L2 (synthetic tx processing)
	// Give time for the L1 block to be processed and synthetic tx to be created
	time.Sleep(30 * time.Second)

	// Step 4: Get L2 config and relay the message
	// Use L2FaucetWallet's client since receiverWallet is a newly generated wallet without registered clients
	l2Client := ti.rpcHandles.TenWalletRndClient(ti.wallets.L2FaucetWallet)
	cfg, err := l2Client.GetConfig()
	if err != nil {
		ti.logger.Error("[WETH Bridge L1->L2] failed to get L2 config", log.ErrKey, err)
		return
	}

	// Create L2 CrossChainMessenger contract instance
	crossChainMessenger, err := CrossChainMessenger.NewCrossChainMessenger(cfg.L2CrossChainMessenger, l2Client)
	if err != nil {
		ti.logger.Error("[WETH Bridge L1->L2] failed to create L2 CrossChainMessenger", log.ErrKey, err)
		return
	}

	// Use a wallet with L2 funds for the relay transaction
	relayWallet := ti.wallets.L2FaucetWallet
	opts, err := bind.NewKeyedTransactorWithChainID(relayWallet.PrivateKey(), relayWallet.ChainID())
	if err != nil {
		ti.logger.Error("[WETH Bridge L1->L2] failed to create transactor", log.ErrKey, err)
		return
	}
	opts.Nonce = big.NewInt(int64(relayWallet.GetNonceAndIncrement()))
	opts.GasLimit = uint64(500_000)

	// Convert to the struct expected by the contract
	msgStruct := CrossChainMessenger.StructsCrossChainMessage{
		Sender:           targetMessage.Sender,
		Sequence:         targetMessage.Sequence,
		Nonce:            targetMessage.Nonce,
		Topic:            targetMessage.Topic,
		Payload:          targetMessage.Payload,
		ConsistencyLevel: targetMessage.ConsistencyLevel,
	}

	// Step 5: Call relayMessage on L2
	relayTx, err := crossChainMessenger.RelayMessage(opts, msgStruct)
	if err != nil {
		ti.logger.Error("[WETH Bridge L1->L2] failed to call relayMessage", log.ErrKey, err)
		return
	}

	ti.logger.Info("[WETH Bridge L1->L2] relay transaction sent", log.TxKey, relayTx.Hash())

	// Wait for relay transaction receipt
	err = testcommon.AwaitReceipt(ti.ctx, l2Client, relayTx.Hash(), 30*time.Second)
	if err != nil {
		ti.logger.Error("[WETH Bridge L1->L2] relay transaction failed", log.ErrKey, err, log.TxKey, relayTx.Hash())
		return
	}

	ti.logger.Info("[WETH Bridge L1->L2] successfully relayed message", log.TxKey, relayTx.Hash(), "amount", amount, "receiver", receiverWallet.Address())

	// Track the successful bridge for verification
	ti.TxTracker.trackWETHBridgingL1ToL2(tx, amount, receiverWallet)
}

// issueRandomWETHWithdrawals bridges WETH from L2 to L1 using EthereumBridge.SendERC20
// The flow is: wrap ETH to WETH on L2 -> approve EthereumBridge -> sendERC20(weth) -> receiver gets native ETH on L1
func (ti *TransactionInjector) issueRandomWETHWithdrawals() {
	wethAddress := gethcommon.HexToAddress("0x1000000000000000000000000000000000000042")

	cfg, err := ti.rpcHandles.TenWalletRndClient(ti.wallets.L2FaucetWallet).GetConfig()
	if err != nil {
		ti.logger.Error("[WETH Withdrawal] failed to get config", log.ErrKey, err)
		return
	}
	l2BridgeAddr := cfg.L2Bridge

	for txCounter := 0; ti.shouldKeepIssuing(txCounter); txCounter++ {
		fromWallet := ti.rndObsWallet()
		client := ti.rpcHandles.TenWalletRndClient(fromWallet)

		amount := big.NewInt(0).SetUint64(testcommon.RndBtw(500, 50_000))

		// Step 1: Wrap native ETH to WETH on L2 by sending ETH to WETH contract
		price, err := client.GasPrice(ti.ctx)
		if err != nil {
			ti.logger.Error("[WETH Withdrawal] failed to get gas price", log.ErrKey, err)
			continue
		}
		price = new(big.Int).Mul(price, big.NewInt(2))

		wrapTxData := &types.LegacyTx{
			Nonce:    fromWallet.GetNonceAndIncrement(),
			To:       &wethAddress,
			Value:    amount,
			Gas:      uint64(100_000),
			GasPrice: price,
			Data:     nil, // sending ETH triggers deposit()
		}

		signedWrapTx, err := fromWallet.SignTransaction(wrapTxData)
		if err != nil {
			ti.logger.Error("[WETH Withdrawal] failed to sign wrap tx", log.ErrKey, err)
			continue
		}

		err = client.SendTransaction(ti.ctx, signedWrapTx)
		if err != nil {
			ti.logger.Error("[WETH Withdrawal] failed to send wrap tx", log.ErrKey, err)
			continue
		}

		err = testcommon.AwaitReceipt(ti.ctx, client, signedWrapTx.Hash(), 30*time.Second)
		if err != nil {
			ti.logger.Error("[WETH Withdrawal] wrap tx failed", log.ErrKey, err, log.TxKey, signedWrapTx.Hash())
			continue
		}
		ti.logger.Info("[WETH Withdrawal] wrapped ETH to WETH on L2", log.TxKey, signedWrapTx.Hash(), "amount", amount)

		// Step 2: Approve EthereumBridge to spend WETH
		approveData := make([]byte, 68)
		copy(approveData[0:4], []byte{0x09, 0x5e, 0xa7, 0xb3})
		copy(approveData[16:36], l2BridgeAddr.Bytes())
		amount.FillBytes(approveData[36:68])

		approveTxData := &types.LegacyTx{
			Nonce:    fromWallet.GetNonceAndIncrement(),
			To:       &wethAddress,
			Value:    big.NewInt(0),
			Gas:      uint64(100_000),
			GasPrice: price,
			Data:     approveData,
		}

		signedApproveTx, err := fromWallet.SignTransaction(approveTxData)
		if err != nil {
			ti.logger.Error("[WETH Withdrawal] failed to sign approve tx", log.ErrKey, err)
			continue
		}

		err = client.SendTransaction(ti.ctx, signedApproveTx)
		if err != nil {
			ti.logger.Error("[WETH Withdrawal] failed to send approve tx", log.ErrKey, err)
			continue
		}

		err = testcommon.AwaitReceipt(ti.ctx, client, signedApproveTx.Hash(), 30*time.Second)
		if err != nil {
			ti.logger.Error("[WETH Withdrawal] approve tx failed", log.ErrKey, err, log.TxKey, signedApproveTx.Hash())
			continue
		}
		ti.logger.Info("[WETH Withdrawal] approved EthereumBridge to spend WETH", log.TxKey, signedApproveTx.Hash())

		// Step 3: Call EthereumBridge.SendERC20(wethAddress, amount, receiver)
		ethereumBridge, err := EthereumBridge.NewEthereumBridge(l2BridgeAddr, client)
		if err != nil {
			ti.logger.Error("[WETH Withdrawal] failed to create EthereumBridge contract", log.ErrKey, err)
			continue
		}

		// Get the publish fee for the bridge message
		fee, err := ethereumBridge.Erc20Fee(&bind.CallOpts{From: fromWallet.Address()})
		if err != nil {
			ti.logger.Error("[WETH Withdrawal] failed to get bridge fee", log.ErrKey, err)
			continue
		}

		opts, err := bind.NewKeyedTransactorWithChainID(fromWallet.PrivateKey(), fromWallet.ChainID())
		if err != nil {
			ti.logger.Error("[WETH Withdrawal] failed to create transactor", log.ErrKey, err)
			continue
		}
		opts.Value = fee
		opts.GasLimit = uint64(500_000)
		opts.GasPrice = price
		opts.Nonce = big.NewInt(int64(fromWallet.GetNonceAndIncrement()))

		// Receiver gets native ETH on L1
		bridgeTx, err := ethereumBridge.SendERC20(opts, wethAddress, amount, fromWallet.Address())
		if err != nil {
			ti.logger.Error("[WETH Withdrawal] failed to call SendERC20", log.ErrKey, err)
			continue
		}

		ti.logger.Info("[WETH Withdrawal] L2->L1 WETH bridge tx sent", log.TxKey, bridgeTx.Hash(), "amount", amount, "receiver", fromWallet.Address())
		go ti.TxTracker.trackWETHBridgingL2ToL1(bridgeTx, amount, fromWallet)
		go ti.awaitAndFinalizeWithdrawal(bridgeTx, fromWallet)

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

	cfg, err := ti.rpcHandles.TenWalletRndClient(fromWallet).GetConfig()
	if err != nil {
		ti.logger.Error("Failed to retrieve config for withdrawal transaction", log.ErrKey, err)
		return
	}

	receipt, err := ti.rpcHandles.TenWalletRndClient(fromWallet).TransactionReceipt(ti.ctx, tx.Hash())
	if err != nil {
		ti.logger.Error("Failed to retrieve receipt for withdrawal transaction", log.ErrKey, err)
		return
	}

	if receipt.Status != 1 {
		ti.logger.Error("Withdrawal transaction failed", log.TxKey, tx.Hash())
		return
	}

	// Filter logs to only include LogMessagePublished events from the MessageBus
	messageBusAddr := cfg.L2MessageBus
	filteredLogs, err := crosschain.FilterLogsFromReceipt(receipt, &messageBusAddr, &ethadapter.CrossChainEventID)
	if err != nil {
		ti.logger.Error("Failed to filter logs for withdrawal transaction", log.ErrKey, err)
		return
	}

	transfers, err := crosschain.ConvertLogsToMessages(filteredLogs, ethadapter.CrossChainEventName, ethadapter.MessageBusABI)
	if err != nil {
		panic(err)
	}

	vTransfers := crosschain.MessageStructs(transfers)

	var proof clientapi.CrossChainProof
	for {
		mtree, err := vTransfers.ForMerkleTree()
		if err != nil {
			panic(err)
		}
		proof, err = ti.rpcHandles.TenWalletRndClient(fromWallet).GetCrossChainProof(ti.ctx, "m", mtree[0][1].(gethcommon.Hash))
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				ti.logger.Info("Proof not found, retrying...", log.ErrKey, err)
				time.Sleep(1 * time.Second)
				continue
			}
			if strings.Contains(err.Error(), "database closed") || strings.Contains(err.Error(), "database is closed") {
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
	if ti.contractRegistryLib.IsMock() {
		return
	}

	opts, err := bind.NewKeyedTransactorWithChainID(ti.wallets.GasWithdrawalWallet.PrivateKey(), ti.wallets.GasWithdrawalWallet.ChainID())
	if err != nil {
		panic(err)
	}

	proof32 := make([][32]byte, 0)
	for i := 0; i < len(proofBytes); i++ {
		proof32 = append(proof32, [32]byte(proofBytes[i][0:32]))
	}

	crossChainMessenger, err := CrossChainMessenger.NewCrossChainMessenger(cfg.L1CrossChainMessenger, ti.rpcHandles.RndEthClient().EthClient())
	if err != nil {
		panic(err)
	}

	oldBalance, err := ti.rpcHandles.RndEthClient().BalanceAt(fromWallet.Address(), nil)
	if err != nil {
		ti.logger.Error("Failed to retrieve balance of receiver", log.ErrKey, err)
		return
	}

	withdrawalTx, err := crossChainMessenger.RelayMessageWithProof(opts, CrossChainMessenger.StructsCrossChainMessage(vTransfers[0]), proof32, proof.Root)
	if err != nil {
		ti.logger.Error("Failed to relay message with proof", log.ErrKey, err)
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

	time.Sleep(15 * time.Second)

	newBalance, err := ti.rpcHandles.RndEthClient().BalanceAt(fromWallet.Address(), nil)
	if err != nil {
		ti.logger.Error("Failed to retrieve balance of receiver", log.ErrKey, err)
		return
	}

	if newBalance.Cmp(oldBalance) == 0 {
		ti.logger.Error("Balance of receiver did not change.")
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
	l2BridgeAddr := cfg.L2Bridge

	for txCounter := 0; ti.shouldKeepIssuing(txCounter); txCounter++ {
		fromWallet := ti.rndObsWallet()
		client := ti.rpcHandles.TenWalletRndClient(fromWallet)
		price, err := client.GasPrice(ti.ctx)
		if err != nil {
			ti.logger.Error("unable to estimate gas price", log.ErrKey, err)
			continue
		}

		price = new(big.Int).Mul(price, big.NewInt(2))

		// Create EthereumBridge contract binding
		ethereumBridge, err := EthereumBridge.NewEthereumBridge(l2BridgeAddr, client)
		if err != nil {
			ti.logger.Error("[CrossChain] unable to create EthereumBridge contract", log.ErrKey, err)
			continue
		}

		// Create transaction options
		opts, err := bind.NewKeyedTransactorWithChainID(fromWallet.PrivateKey(), fromWallet.ChainID())
		if err != nil {
			ti.logger.Error("[CrossChain] unable to create transaction options", log.ErrKey, err)
			continue
		}
		opts.Value = big.NewInt(100) // Send 1 wei
		opts.GasLimit = uint64(10_000_000)
		opts.GasPrice = price
		opts.Nonce = big.NewInt(int64(fromWallet.GetNonceAndIncrement()))

		// Call sendNative on the bridge
		signedTx, err := ethereumBridge.SendNative(opts, fromWallet.Address())
		if err != nil {
			ti.logger.Error("[CrossChain] unable to send withdrawal transaction", log.ErrKey, err)
			continue
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
	return ti.newTx(data, from.GetNonceAndIncrement(), ercType, ti.rpcHandles.TenWalletRndClient(from))
}

func (ti *TransactionInjector) newCustomTenWithdrawalTx(amount uint64) types.TxData {
	transferERC20data := erc20contractlib.CreateTransferTxData(testcommon.BridgeAddress, common.ValueInWei(big.NewInt(int64(amount))))
	return ti.newTx(transferERC20data, 1, testcommon.HOC, ti.rpcHandles.TenWalletRndClient(ti.wallets.L2FaucetWallet))
}

func (ti *TransactionInjector) newTx(data []byte, nonce uint64, ercType testcommon.ERC20, client *obsclient.AuthObsClient) types.TxData {
	price, err := client.GasPrice(ti.ctx)
	if err != nil {
		// Fallback to a reasonable gas price if we can't get current price
		price = big.NewInt(2000000000) // 2 gwei
	} else {
		price = new(big.Int).Mul(price, big.NewInt(2))
	}

	return &types.LegacyTx{
		Nonce:    nonce,
		Value:    gethcommon.Big0,
		Gas:      uint64(10_000_000),
		GasPrice: price,
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
