package components

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/ethereum/go-ethereum/core/vm"

	"github.com/ten-protocol/go-ten/go/common/gethapi"
	enclaveconfig "github.com/ten-protocol/go-ten/go/enclave/config"
	tencore "github.com/ten-protocol/go-ten/go/enclave/core"
	gethrpc "github.com/ten-protocol/go-ten/lib/gethfork/rpc"

	"github.com/ten-protocol/go-ten/go/enclave/gas"
	"github.com/ten-protocol/go-ten/go/enclave/storage"

	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ten-protocol/go-ten/go/common/log"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"

	gethtxpool "github.com/ethereum/go-ethereum/core/txpool"
	"github.com/ethereum/go-ethereum/core/txpool/legacypool"
	"github.com/ethereum/go-ethereum/core/types"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common"
)

const (
	// txSlotSize is used to calculate how many data slots a single transaction
	// takes up based on its size. The slots are used as DoS protection, ensuring
	// that validating a new transaction remains a constant operation (in reality
	// O(maxslots), where max slots are 4 currently).
	txSlotSize = 32 * 1024

	// we assume that at the limit, a single "uncompressable" tx is in a batch which gets rolled-up, and must fit in a 128kb blob
	rollupOverhead = 5 * 1024

	// txMaxSize is the maximum size a single transaction can have. This field has
	// non-trivial consequences: larger transactions are significantly harder and
	// more expensive to propagate; larger transactions also take more resources
	// to validate whether they fit into the pool or not.
	txMaxSize = 4*txSlotSize - rollupOverhead // 128KB - overhead
)

// this is how long the node waits to receive the second batch (longer now as we have to wait for all the contracts to be deployed)
var startMempoolTimeout = 180 * time.Second

// TxPool is an obscuro wrapper around geths transaction pool
type TxPool struct {
	txPoolConfig     legacypool.Config
	chainconfig      *params.ChainConfig
	legacyPool       *legacypool.LegacyPool
	pool             *gethtxpool.TxPool
	Chain            *EthChainAdapter
	gasOracle        gas.Oracle
	batchRegistry    BatchRegistry
	storage          storage.Storage
	l1BlockProcessor L1BlockProcessor
	gasTip           *big.Int
	running          atomic.Bool
	stateMutex       sync.Mutex
	logger           gethlog.Logger
	validateOnly     atomic.Bool
	config           *enclaveconfig.EnclaveConfig
	tenChain         TENChain
}

// NewTxPool returns a new instance of the tx pool
func NewTxPool(blockchain *EthChainAdapter, config *enclaveconfig.EnclaveConfig, tenChain TENChain, storage storage.Storage, batchRegistry BatchRegistry, l1BlockProcessor L1BlockProcessor, gasOracle gas.Oracle, gasTip *big.Int, validateOnly bool, logger gethlog.Logger) (*TxPool, error) {
	txPoolConfig := legacypool.Config{
		Locals:       nil,
		NoLocals:     false,
		Journal:      "",
		Rejournal:    0,
		PriceLimit:   legacypool.DefaultConfig.PriceLimit,
		PriceBump:    legacypool.DefaultConfig.PriceBump,
		AccountSlots: 32,
		GlobalSlots:  (4096 + 1024) * 2,
		AccountQueue: 2048,
		GlobalQueue:  2048 * 4,
		Lifetime:     legacypool.DefaultConfig.Lifetime,
	}
	legacyPool := legacypool.New(txPoolConfig, blockchain)

	txp := &TxPool{
		Chain:            blockchain,
		config:           config,
		tenChain:         tenChain,
		chainconfig:      blockchain.Config(),
		txPoolConfig:     txPoolConfig,
		storage:          storage,
		gasOracle:        gasOracle,
		batchRegistry:    batchRegistry,
		l1BlockProcessor: l1BlockProcessor,
		legacyPool:       legacyPool,
		gasTip:           gasTip,
		stateMutex:       sync.Mutex{},
		validateOnly:     atomic.Bool{},
		logger:           logger,
	}
	txp.validateOnly.Store(validateOnly)
	go txp.start()
	return txp, nil
}

func (t *TxPool) SetValidateMode(validateOnly bool) {
	t.validateOnly.Store(validateOnly)
}

// can only be started after t.blockchain has at least one block inside
// note - blocking method that waits for the block.Call only as goroutine
func (t *TxPool) start() {
	if t.running.Load() {
		return
	}

	cb := t.Chain.CurrentBlock()
	if cb != nil && cb.Number.Uint64() > common.L2GenesisHeight+1 {
		err := t._startInternalPool()
		if err != nil {
			t.logger.Crit("Failed to start tx pool", log.ErrKey, err)
		}
		return
	}

	var (
		newHeadCh  = make(chan core.ChainHeadEvent)
		newHeadSub = t.Chain.SubscribeChainHeadEvent(newHeadCh)
	)
	defer close(newHeadCh)
	defer newHeadSub.Unsubscribe()
	for {
		select {
		case event := <-newHeadCh:
			newHead := event.Header
			if newHead.Number.Uint64() > common.L2GenesisHeight+1 {
				err := t._startInternalPool()
				if err != nil {
					t.logger.Crit("Failed to start tx pool", log.ErrKey, err)
				}
				return
			}
		case <-time.After(startMempoolTimeout):
			t.logger.Crit("Timeout waiting to start mempool.")
			return
		}
	}
}

func (t *TxPool) _startInternalPool() error {
	t.logger.Info("Starting tx pool")
	memp, err := gethtxpool.New(t.gasTip.Uint64(), t.Chain, []gethtxpool.SubPool{t.legacyPool})
	if err != nil {
		return fmt.Errorf("unable to init geth tx pool - %w", err)
	}
	t.logger.Info("Tx pool started")

	t.pool = memp
	t.running.Store(true)
	return nil
}

func (t *TxPool) SubmitTx(transaction *common.L2Tx) (error, error) {
	err := t.waitUntilPoolRunning()
	if err != nil {
		return nil, err
	}

	if t.validateOnly.Load() {
		return t.validate(transaction)
	}

	// this code runs only on the sequencer
	// it sets the time of entry into the mempool
	// so it can be checked in the smart contract
	transaction.SetTime(time.Now())

	return t.add(transaction), nil
}

func (t *TxPool) waitUntilPoolRunning() error {
	if t.running.Load() {
		return nil
	}

	timeout := time.After(startMempoolTimeout)
	tick := time.NewTicker(500 * time.Millisecond)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			if t.running.Load() {
				return nil
			}
		case <-timeout:
			return fmt.Errorf("timed out waiting for tx pool to start")
		}
	}
}

// PendingTransactions returns all pending transactions grouped per address and ordered per nonce
func (t *TxPool) PendingTransactions(batchTime uint64) map[gethcommon.Address][]*gethtxpool.LazyTransaction {
	if !t.running.Load() {
		t.logger.Error("tx pool not running")
		return nil
	}

	if t.validateOnly.Load() {
		t.logger.Error("Pending transactions requested while in validate only mode")
		return nil
	}

	// todo - for now using the base fee from the block
	currentBlock := t.Chain.CurrentBlock()
	if currentBlock == nil {
		return make(map[gethcommon.Address][]*gethtxpool.LazyTransaction)
	}
	baseFee := currentBlock.BaseFee
	txs := t.pool.Pending(gethtxpool.PendingFilter{
		BaseFee:      uint256.NewInt(baseFee.Uint64()),
		OnlyPlainTxs: true,
	})

	// Filter out transactions that have "Time" greater than batchTime + MaxNegativeTxTimeDeltaMs
	// this is required for serialising the transactiong together with their timestamp delta (which can't be negative).
	maxDeltaSec := (common.MaxNegativeTxTimeDeltaMs / 1000) - 1
	maxTime := time.Unix(int64(batchTime)+int64(maxDeltaSec), 0)

	filteredTxs := make(map[gethcommon.Address][]*gethtxpool.LazyTransaction)
	for addr, addrTxs := range txs {
		var validTxs []*gethtxpool.LazyTransaction
		for _, tx := range addrTxs {
			if tx.Time.Before(maxTime) {
				validTxs = append(validTxs, tx)
			} else {
				t.logger.Warn("Transaction excluded from PendingTransactions for being too recent. Should not happen", log.TxKey, tx.Hash, "tx_time", tx.Time, "block_time", batchTime)
			}
		}
		if len(validTxs) > 0 {
			filteredTxs[addr] = validTxs
		}
	}
	return filteredTxs
}

func (t *TxPool) Close() error {
	defer func() {
		if err := recover(); err != nil {
			t.logger.Error("Could not close legacy pool", log.ErrKey, err)
		}
	}()
	return t.pool.Close()
}

// Add adds a new transactions to the pool
func (t *TxPool) add(transaction *common.L2Tx) error {
	// validate against the consensus rules
	err := t.validateTxBasics(transaction, false)
	if err != nil {
		return err
	}

	var strErrors []string
	for _, err := range t.pool.Add([]*types.Transaction{transaction}, false) {
		if err != nil {
			strErrors = append(strErrors, err.Error())
		}
	}

	if len(strErrors) > 0 {
		return fmt.Errorf(strings.Join(strErrors, "; ")) // nolint
	}
	return nil
}

//go:linkname validateTx github.com/ethereum/go-ethereum/core/txpool/legacypool.(*LegacyPool).validateTx
func validateTx(_ *legacypool.LegacyPool, _ *types.Transaction, _ bool) error

// Validate - run the underlying tx pool validation logic
func (t *TxPool) validate(tx *common.L2Tx) (error, error) {
	// validate against the consensus rules
	err := t.validateTxBasics(tx, false)
	if err != nil {
		return err, nil
	}

	t.stateMutex.Lock()
	// validate against the state. Things like nonce, balance, etc
	err = validateTx(t.legacyPool, tx, false)
	if err != nil {
		t.stateMutex.Unlock()
		return err, nil
	}
	t.stateMutex.Unlock()
	return t.validateTotalGas(tx)
}

func (t *TxPool) Stats() (int, int) {
	return t.legacyPool.Stats()
}

// validateTxBasics checks whether a transaction is valid according to the consensus
// rules, but does not check state-dependent validation such as sufficient balance.
// This check is meant as an early check which only needs to be performed once,
// and does not require the pool mutex to be held.
func (t *TxPool) validateTxBasics(tx *types.Transaction, local bool) error {
	opts := &gethtxpool.ValidationOptions{
		Config: t.chainconfig,
		Accept: 0 |
			1<<types.LegacyTxType |
			1<<types.AccessListTxType |
			1<<types.DynamicFeeTxType,
		MaxSize: txMaxSize,
		MinTip:  t.gasTip,
	}

	// we need to access some private variables from the legacy pool to run validation with our own consensus options
	v := reflect.ValueOf(t.legacyPool).Elem()

	chField := v.FieldByName("currentHead")
	chFieldPtr := unsafe.Pointer(chField.UnsafeAddr())
	ch, ok := reflect.NewAt(chField.Type(), chFieldPtr).Elem().Interface().(atomic.Pointer[types.Header]) //nolint:govet
	if !ok {
		t.logger.Crit("invalid mempool. should not happen")
	}

	sigField := v.FieldByName("signer")
	sigFieldPtr := unsafe.Pointer(sigField.UnsafeAddr())
	sig, ok1 := reflect.NewAt(sigField.Type(), sigFieldPtr).Elem().Interface().(types.Signer)
	if !ok1 {
		t.logger.Crit("invalid mempool. should not happen")
	}

	if err := gethtxpool.ValidateTransaction(tx, ch.Load(), sig, opts); err != nil {
		return err
	}
	return nil
}

// check that the tx gas can pay for the l1
func (t *TxPool) validateTotalGas(tx *common.L2Tx) (error, error) {
	headBatchSeq := t.batchRegistry.HeadBatchSeq()

	// don't perform the check while the network is initialising
	if headBatchSeq == nil {
		return nil, nil
	}

	headBatch, err := t.storage.FetchBatchHeaderBySeqNo(context.Background(), headBatchSeq.Uint64())
	if err != nil {
		return nil, fmt.Errorf("could not retrieve head batch. Cause: %w", err)
	}

	serTx, err := tx.MarshalJSON()
	if err != nil {
		return fmt.Errorf("could not marshal tx. Cause: %w", err), nil
	}
	txArgs := gethapi.TransactionArgs{}
	err = json.Unmarshal(serTx, &txArgs)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal tx. Cause: %w", err)
	}
	from, err := tencore.GetExternalTxSigner(tx)
	if err != nil {
		return nil, fmt.Errorf("could not extract sender from transaction. Cause: %w", err)
	}
	txArgs.From = &from
	ge := NewGasEstimator(t.storage, t.tenChain, t.gasOracle, t.logger)
	latest := gethrpc.LatestBlockNumber
	leastGas, publishingGas, userErr, sysErr := ge.EstimateTotalGas(context.Background(), &txArgs, &latest, headBatch, t.config.GasLocalExecutionCapFlag)

	// if the transaction reverts we let it through
	if userErr != nil && errors.Is(userErr, vm.ErrExecutionReverted) {
		return nil, nil
	}

	if userErr != nil || sysErr != nil {
		return userErr, sysErr
	}

	// make sure the tx has enough gas to cover the execution and the tx won't be rejected by the sequencer
	leastGas = leastGas * 8 / 10 // reduce gas estimate by 20%
	if tx.Gas() < leastGas {
		return fmt.Errorf("insufficient gas. Want at least: %d have: %d", leastGas, tx.Gas()), nil
	}

	// make sure the tx has enough gas to cover the publishing cost even if by some chance the 20% deducted was too much
	if tx.Gas() < publishingGas+params.TxGas {
		return fmt.Errorf("insufficient gas to publish the transaction to the DA. Want at least: %d have: %d", publishingGas, tx.Gas()), nil
	}

	return nil, nil
}
