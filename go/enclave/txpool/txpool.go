package txpool

import (
	"fmt"
	"math/big"
	"strings"
	"sync"
	_ "unsafe"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"

	// unsafe package imported in order to link to a private function in go-ethereum.
	// This allows us to validate transactions against the tx pool rules.

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common/log"

	gethtxpool "github.com/ethereum/go-ethereum/core/txpool"
	"github.com/ethereum/go-ethereum/core/txpool/legacypool"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/enclave/evm/ethchainadapter"
)

// TxPool is an obscuro wrapper around geths transaction pool
type TxPool struct {
	txPoolConfig legacypool.Config
	legacyPool   *legacypool.LegacyPool
	pool         *gethtxpool.TxPool
	Chain        *ethchainadapter.EthChainAdapter
	gasTip       *big.Int
	running      bool
	stateMutex   sync.Mutex
	logger       gethlog.Logger
}

// NewTxPool returns a new instance of the tx pool
func NewTxPool(blockchain *ethchainadapter.EthChainAdapter, gasTip *big.Int, logger gethlog.Logger) (*TxPool, error) {
	txPoolConfig := ethchainadapter.NewLegacyPoolConfig()
	legacyPool := legacypool.New(txPoolConfig, blockchain)

	return &TxPool{
		Chain:        blockchain,
		txPoolConfig: txPoolConfig,
		legacyPool:   legacyPool,
		gasTip:       gasTip,
		stateMutex:   sync.Mutex{},
		logger:       logger,
	}, nil
}

// Start starts the pool
// can only be started after t.blockchain has at least one block inside
func (t *TxPool) Start() error {
	if t.running {
		return fmt.Errorf("tx pool already started")
	}

	t.logger.Info("Starting tx pool")
	memp, err := gethtxpool.New(t.gasTip.Uint64(), t.Chain, []gethtxpool.SubPool{t.legacyPool})
	if err != nil {
		return fmt.Errorf("unable to init geth tx pool - %w", err)
	}
	t.logger.Info("Tx pool started")

	t.pool = memp
	t.running = true
	return nil
}

// PendingTransactions returns all pending transactions grouped per address and ordered per nonce
func (t *TxPool) PendingTransactions() map[gethcommon.Address][]*gethtxpool.LazyTransaction {
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

// Add adds a new transactions to the pool
func (t *TxPool) Add(transaction *common.L2Tx) error {
	if !t.running {
		return fmt.Errorf("tx pool not running")
	}
	var strErrors []string
	for _, err := range t.pool.Add([]*types.Transaction{transaction}, false, false) {
		if err != nil {
			strErrors = append(strErrors, err.Error())
		}
	}

	if len(strErrors) > 0 {
		return fmt.Errorf(strings.Join(strErrors, "; "))
	}
	return nil
}

//go:linkname validateTxBasics github.com/ethereum/go-ethereum/core/txpool/legacypool.(*LegacyPool).validateTxBasics
func validateTxBasics(_ *legacypool.LegacyPool, _ *types.Transaction, _ bool) error

//go:linkname validateTx github.com/ethereum/go-ethereum/core/txpool/legacypool.(*LegacyPool).validateTx
func validateTx(_ *legacypool.LegacyPool, _ *types.Transaction, _ bool) error

// Validate - run the underlying tx pool validation logic
func (t *TxPool) Validate(tx *common.L2Tx) error {
	// validate against the consensus rules
	err := validateTxBasics(t.legacyPool, tx, false)
	if err != nil {
		return err
	}

	t.stateMutex.Lock()
	defer t.stateMutex.Unlock()
	// validate against the state. Things like nonce, balance, etc
	return validateTx(t.legacyPool, tx, false)
}

func (t *TxPool) Running() bool {
	return t.running
}

func (t *TxPool) Close() error {
	defer func() {
		if err := recover(); err != nil {
			t.logger.Error("Could not close legacy pool", log.ErrKey, err)
		}
	}()
	return t.pool.Close()
}
