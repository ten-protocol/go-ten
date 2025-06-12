package gas

import (
	"context"
	"fmt"
	"math/big"
	"sync"

	gethlog "github.com/ethereum/go-ethereum/log"

	"github.com/ethereum/go-ethereum/consensus/misc/eip4844"
	"github.com/ten-protocol/go-ten/go/enclave/storage"

	"github.com/ethereum/go-ethereum/params"

	"github.com/ten-protocol/go-ten/go/common"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ten-protocol/go-ten/go/common/gethapi"
)

// MovingAverageWindow - the more traffic on the network, the lower this number can get. Should be roughly the number of blocks between rollups.
// - note that increasing this value will have to increase the number of cached blocks in the cache_service
const MovingAverageWindow = 300 // `3600 / 12` - last 1 hour

// MaxHistoricMA - the maximum number of historic blocks
const MaxHistoricMA = 50

// L1TxGas - a crude estimation of the cost of publishing an L1 tx
const L1TxGas = 300_000

// TxsPerRollup - the number of transactions in a rollup. A conservative estimation.
const TxsPerRollup = 250

// Oracle - the interface for the future precompiled gas oracle contract
// which will expose necessary l1 information.
type Oracle interface {
	SubmitL1Block(ctx context.Context, headBlock *types.Header) error
	EstimateL1StorageGasCost(ctx context.Context, tx *types.Transaction, block *types.Header, header *common.BatchHeader) (*big.Int, error)
	EstimateL1CostForMsg(ctx context.Context, args *gethapi.TransactionArgs, header *common.BatchHeader) (*big.Int, error)
}

type oracle struct {
	l1ChainCfg *params.ChainConfig
	storage    storage.BlockResolver
	headMutex  sync.RWMutex
	headBlock  *types.Header
	// hold the moving average of the base fee and blob fee per block
	blobFeeMA map[uint64]*big.Int
	baseFeeMA map[uint64]*big.Int
	logger    gethlog.Logger
}

func NewGasOracle(l1ChainCfg *params.ChainConfig, storage storage.BlockResolver, logger gethlog.Logger) Oracle {
	return &oracle{
		l1ChainCfg: l1ChainCfg,
		storage:    storage,
		logger:     logger,
		headMutex:  sync.RWMutex{},
		baseFeeMA:  make(map[uint64]*big.Int),
		blobFeeMA:  make(map[uint64]*big.Int),
	}
}

// EstimateL1StorageGasCost - Returns the expected l1 gas cost for a transaction at a given l1 block.
func (o *oracle) EstimateL1StorageGasCost(ctx context.Context, tx *types.Transaction, block *types.Header, header *common.BatchHeader) (*big.Int, error) {
	encodedTx, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return nil, err
	}

	return o.calculateL1Cost(ctx, block, header, encodedTx)
}

func (o *oracle) EstimateL1CostForMsg(ctx context.Context, args *gethapi.TransactionArgs, header *common.BatchHeader) (*big.Int, error) {
	encoded, err := rlp.EncodeToBytes(args)
	if err != nil {
		return nil, err
	}

	return o.calculateL1Cost(ctx, nil, header, encoded)
}

func (o *oracle) SubmitL1Block(ctx context.Context, headBlock *types.Header) error {
	o.headMutex.Lock()
	defer o.headMutex.Unlock()
	o.headBlock = headBlock
	blockNum := headBlock.Number.Uint64()

	baseFeeMA, blobFeeMA, err := o.calculateMA(ctx, blockNum)
	if err != nil {
		return err
	}

	o.baseFeeMA[blockNum] = baseFeeMA
	o.blobFeeMA[blockNum] = blobFeeMA

	if blockNum > MaxHistoricMA {
		// cleanup entries older than MaxHistoricMA
		delete(o.baseFeeMA, blockNum-MaxHistoricMA)
		delete(o.blobFeeMA, blockNum-MaxHistoricMA)
	}
	return nil
}

// calculateMA - calculates the baseFee and blobFee MA for the specified block
// by walking back the window and
func (o *oracle) calculateMA(ctx context.Context, blockHeight uint64) (*big.Int, *big.Int, error) {
	b, err := o.storage.FetchCanonicaBlockByHeight(ctx, big.NewInt(int64(blockHeight)))
	if err != nil {
		return nil, nil, fmt.Errorf("failed getting block %d: %w", blockHeight, err)
	}

	calculateBlobs := b.ExcessBlobGas == nil

	baseFeeSum := big.NewInt(0)
	blobFeeSum := big.NewInt(0)
	count := 0

	for ; count < MovingAverageWindow; count++ {
		baseFeeSum = baseFeeSum.Add(baseFeeSum, b.BaseFee)
		if calculateBlobs {
			blobFeeSum = blobFeeSum.Add(blobFeeSum, eip4844.CalcBlobFee(o.l1ChainCfg, b))
		}
		b, err = o.storage.FetchBlock(ctx, b.ParentHash)
		if err != nil {
			break
		}
	}

	baseFeeMA, blobFeeMA := big.NewInt(0), big.NewInt(0)
	if count > 0 {
		baseFeeMA = new(big.Int).Div(baseFeeSum, big.NewInt(int64(count)))
		if calculateBlobs {
			blobFeeMA = new(big.Int).Div(blobFeeSum, big.NewInt(int64(count)))
		}
	}
	return baseFeeMA, blobFeeMA, nil
}

// calculateL1Cost - Calculates the L1 cost as a multiple of the L2 base fee.
// it takes into account the share of the blob cost and the share of the L1 TX cost - which submits and stores the rollup header.
func (o *oracle) calculateL1Cost(ctx context.Context, block *types.Header, l2Batch *common.BatchHeader, encodedTx []byte) (*big.Int, error) {
	o.headMutex.RLock()
	defer o.headMutex.RUnlock()

	if block == nil {
		block = o.headBlock
	}

	shareOfBlobCost := big.NewInt(0)
	txL1Size := CalculateL1Size(encodedTx)

	var err error
	baseFee := o.baseFeeMA[block.Number.Uint64()]
	blobFee, f := o.blobFeeMA[block.Number.Uint64()]
	// the block MA is not in the cache, so we're calculating
	if !f {
		o.logger.Info("L1 Gas MA not found in cache, calculating it", "block_number", block.Number.Uint64())
		baseFee, blobFee, err = o.calculateMA(ctx, block.Number.Uint64())
		if err != nil {
			return nil, err
		}
	}

	// 1. Calculate the cost of including the tx in a blob
	// price in Wei for a single unit of blob
	if isNonZero(blobFee) {
		shareOfBlobCost = big.NewInt(0).Mul(txL1Size, o.blobFeeMA[block.Number.Uint64()])
	}

	// 2. Estimate how much this tx should absorb from the L1 tx cost that submits the rollup
	shareOfL1TxGas := big.NewInt(L1TxGas / TxsPerRollup)
	shareOfL1TxCost := big.NewInt(0)
	if isNonZero(baseFee) {
		shareOfL1TxCost = big.NewInt(0).Mul(shareOfL1TxGas, baseFee)
	}

	// 3. The total cost is the sum of the share of the blob cost and the share of the L1 tx cost
	totalCost := big.NewInt(0).Add(shareOfBlobCost, shareOfL1TxCost)

	// 4. round the shareOfBlobCost up to the nearest multiple of l2Batch.BaseFee
	remainder := new(big.Int).Mod(totalCost, l2Batch.BaseFee)
	if remainder.Sign() > 0 {
		totalCost = totalCost.Add(totalCost, new(big.Int).Sub(l2Batch.BaseFee, remainder))
	}

	return totalCost, nil
}

func isNonZero(nr *big.Int) bool {
	return nr != nil && nr.Sign() > 0
}
