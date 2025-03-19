package components

import (
	"context"
	"fmt"
	"math/big"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

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
var startMempoolTimeout = 150 * time.Second

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
}

// NewTxPool returns a new instance of the tx pool
func NewTxPool(blockchain *EthChainAdapter, storage storage.Storage, batchRegistry BatchRegistry, l1BlockProcessor L1BlockProcessor, gasOracle gas.Oracle, gasTip *big.Int, validateOnly bool, logger gethlog.Logger) (*TxPool, error) {
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

func (t *TxPool) SubmitTx(transaction *common.L2Tx) error {
	err := t.waitUntilPoolRunning()
	if err != nil {
		return err
	}

	if t.validateOnly.Load() {
		return t.validate(transaction)
	}
	return t.add(transaction)
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
func (t *TxPool) PendingTransactions() map[gethcommon.Address][]*gethtxpool.LazyTransaction {
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
	return t.pool.Pending(gethtxpool.PendingFilter{
		BaseFee:      uint256.NewInt(baseFee.Uint64()),
		OnlyPlainTxs: true,
	})
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
func (t *TxPool) validate(tx *common.L2Tx) error {
	// validate against the consensus rules
	err := t.validateTxBasics(tx, false)
	if err != nil {
		return err
	}

	t.stateMutex.Lock()
	// validate against the state. Things like nonce, balance, etc
	err = validateTx(t.legacyPool, tx, false)
	if err != nil {
		t.stateMutex.Unlock()
		return err
	}
	t.stateMutex.Unlock()
	return t.validateL1Gas(tx)
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
func (t *TxPool) validateL1Gas(tx *common.L2Tx) error {
	headBatchSeq := t.batchRegistry.HeadBatchSeq()

	// don't perform the check while the network is initialising
	if headBatchSeq == nil {
		return nil
	}

	headBatch, err := t.storage.FetchBatchHeaderBySeqNo(context.Background(), headBatchSeq.Uint64())
	if err != nil {
		return fmt.Errorf("could not retrieve head batch. Cause: %w", err)
	}
	headBlock, err := t.l1BlockProcessor.GetHead(context.Background())
	if err != nil {
		return fmt.Errorf("could not retrieve head block. Cause: %w", err)
	}

	l1Cost, err := t.gasOracle.EstimateL1StorageGasCost(tx, headBlock, headBatch)
	if err != nil {
		t.logger.Error("Unable to get gas cost for tx. Should not happen at this point.", log.TxKey, tx.Hash(), log.ErrKey, err)
		return fmt.Errorf("unable to get gas cost for tx. Cause: %w", err)
	}

	// calculate the cost in l2 gas
	l2Gas := big.NewInt(0).Div(l1Cost, headBatch.BaseFee)

	intrGas, err := core.IntrinsicGas(tx.Data(), tx.AccessList(), tx.SetCodeAuthorizations(), tx.To() == nil, true, true, true)
	if err != nil {
		return err
	}

	leastGas := l2Gas.Uint64() + intrGas
	// The gas limit of the transaction (evm message) should always be higher than the gas overhead
	// used to cover the l1 cost
	if tx.Gas() < leastGas {
		return fmt.Errorf("insufficient gas. Want at least: %d have: %d", leastGas, tx.Gas())
	}

	return nil
}
