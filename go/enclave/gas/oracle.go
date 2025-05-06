package gas

import (
	"context"
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/consensus/misc/eip4844"
	"github.com/ten-protocol/go-ten/go/enclave/storage"

	"github.com/ethereum/go-ethereum/params"

	"github.com/ten-protocol/go-ten/go/common"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ten-protocol/go-ten/go/common/gethapi"
)

// MovingAverageWindow - the more traffic on the network, the lower this number can get. Should be roughly the number of blocks between rollups.
const MovingAverageWindow = 1 // `3600 / 12` - last 1 hour

// L1TxGas - a crude estimation of the cost of publishing an L1 tx
const L1TxGas = 300_000

// TxsPerRollup - the number of transactions in a rollup. A conservative estimation.
const TxsPerRollup = 250

// Oracle - the interface for the future precompiled gas oracle contract
// which will expose necessary l1 information.
type Oracle interface {
	SubmitL1Block(ctx context.Context, headBlock *types.Header) error
	EstimateL1StorageGasCost(tx *types.Transaction, block *types.Header, header *common.BatchHeader) (*big.Int, error)
	EstimateL1CostForMsg(args *gethapi.TransactionArgs, header *common.BatchHeader) (*big.Int, error)
}

type oracle struct {
	l1ChainCfg *params.ChainConfig
	storage    storage.BlockResolver
	headMutex  sync.RWMutex
	headBlock  *types.Header
	blobFeeMA  map[uint64]*big.Int
	baseFeeMA  map[uint64]*big.Int
}

func NewGasOracle(l1ChainCfg *params.ChainConfig, storage storage.BlockResolver) Oracle {
	return &oracle{
		l1ChainCfg: l1ChainCfg,
		storage:    storage,
		headMutex:  sync.RWMutex{},
		baseFeeMA:  make(map[uint64]*big.Int),
		blobFeeMA:  make(map[uint64]*big.Int),
	}
}

// EstimateL1StorageGasCost - Returns the expected l1 gas cost for a transaction at a given l1 block.
func (o *oracle) EstimateL1StorageGasCost(tx *types.Transaction, block *types.Header, header *common.BatchHeader) (*big.Int, error) {
	encodedTx, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return nil, err
	}

	return o.calculateL1Cost(block, header, encodedTx)
}

func (o *oracle) EstimateL1CostForMsg(args *gethapi.TransactionArgs, header *common.BatchHeader) (*big.Int, error) {
	encoded, err := rlp.EncodeToBytes(args)
	if err != nil {
		return nil, err
	}

	return o.calculateL1Cost(nil, header, encoded)
}

func (o *oracle) SubmitL1Block(ctx context.Context, headBlock *types.Header) error {
	o.headMutex.Lock()
	defer o.headMutex.Unlock()
	o.headBlock = headBlock
	calculateBlobs := headBlock.ExcessBlobGas == nil

	baseFeeSum := big.NewInt(0)
	blobFeeSum := big.NewInt(0)
	count := 0

	b := headBlock
	var err error
	for ; count < MovingAverageWindow; count++ {
		baseFeeSum = baseFeeSum.Add(baseFeeSum, b.BaseFee)
		if calculateBlobs {
			blobFeeSum = blobFeeSum.Add(blobFeeSum, eip4844.CalcBlobFee(o.l1ChainCfg, b))
		}
		parent := b.ParentHash
		b, err = o.storage.FetchBlock(ctx, b.ParentHash)
		if err != nil {
			fmt.Printf("Block %s, Err %v\n", parent.Hex(), err)
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

	o.baseFeeMA[headBlock.Number.Uint64()] = baseFeeMA
	o.blobFeeMA[headBlock.Number.Uint64()] = blobFeeMA
	return nil
}

// calculateL1Cost - Calculates the L1 cost as a multiple of the L2 base fee.
// it takes into account the share of the blob cost and the share of the L1 TX cost - which submits and stores the rollup header.
func (o *oracle) calculateL1Cost(block *types.Header, l2Batch *common.BatchHeader, encodedTx []byte) (*big.Int, error) {
	o.headMutex.RLock()
	defer o.headMutex.RUnlock()

	if block == nil {
		block = o.headBlock
	}

	// 1. Calculate the cost of including the tx in a blob
	// price in Wei for a single unit of blob
	shareOfBlobCost := big.NewInt(0)
	txL1Size := CalculateL1Size(encodedTx)
	bl := o.blobFeeMA[block.Number.Uint64()]
	if bl != nil && bl.Sign() > 0 {
		shareOfBlobCost = big.NewInt(0).Mul(txL1Size, o.blobFeeMA[block.Number.Uint64()])
	}

	// 2. Estimate how much this tx should absorb from the L1 tx cost that submits the rollup
	shareOfL1TxGas := big.NewInt(L1TxGas / TxsPerRollup)
	shareOfL1TxCost := big.NewInt(0)
	bf := o.baseFeeMA[block.Number.Uint64()]
	if bf != nil && bf.Sign() > 0 {
		shareOfL1TxCost = big.NewInt(0).Mul(shareOfL1TxGas, bf)
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
