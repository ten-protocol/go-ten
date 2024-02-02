package txpool

import (
	"fmt"
	"math/big"
	"strings"

	gethcommon "github.com/ethereum/go-ethereum/common"
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
	blockchain   *ethchainadapter.EthChainAdapter
	gasTip       *big.Int
	running      bool
}

// NewTxPool returns a new instance of the tx pool
func NewTxPool(blockchain *ethchainadapter.EthChainAdapter, gasTip *big.Int) (*TxPool, error) {
	txPoolConfig := ethchainadapter.NewLegacyPoolConfig()
	legacyPool := legacypool.New(txPoolConfig, blockchain)

	return &TxPool{
		blockchain:   blockchain,
		txPoolConfig: txPoolConfig,
		legacyPool:   legacyPool,
		gasTip:       gasTip,
	}, nil
}

// Start starts the pool
// can only be started after t.blockchain has at least one block inside
func (t *TxPool) Start() error {
	if t.pool != nil {
		return fmt.Errorf("tx pool already started")
	}

	memp, err := gethtxpool.New(t.gasTip, t.blockchain, []gethtxpool.SubPool{t.legacyPool})
	if err != nil {
		return fmt.Errorf("unable to init geth tx pool - %w", err)
	}

	t.pool = memp
	t.running = true
	return nil
}

// PendingTransactions returns all pending transactions grouped per address and ordered per nonce
func (t *TxPool) PendingTransactions() map[gethcommon.Address][]*gethtxpool.LazyTransaction {
	return t.pool.Pending(false)
}

// Add adds a new transactions to the pool
func (t *TxPool) Add(transaction *common.L2Tx) error {
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

func (t *TxPool) Running() bool {
	return t.running
}
