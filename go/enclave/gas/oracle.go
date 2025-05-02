package gas

import (
	"math/big"

	"github.com/ethereum/go-ethereum/consensus/misc/eip4844"

	"github.com/ethereum/go-ethereum/params"

	"github.com/ten-protocol/go-ten/go/common"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ten-protocol/go-ten/go/common/gethapi"
)

const MovingAverageWindow = 50

// L1TxGas - a crude estimation of the cost of publishing an L1 tx
const L1TxGas = 150_000

// TxsPerRollup - the number of transactions in a rollup. A conservative estimation.
const TxsPerRollup = 200

// Oracle - the interface for the future precompiled gas oracle contract
// which will expose necessary l1 information.
type Oracle interface {
	EstimateL1StorageGasCost(tx *types.Transaction, header *common.BatchHeader) (*big.Int, error)
	EstimateL1CostForMsg(args *gethapi.TransactionArgs, header *common.BatchHeader) (*big.Int, error)
	SubmitL1Block(block *types.Header) error
}

type oracle struct {
	l1ChainCfg *params.ChainConfig
	blobFeeMA  *big.Int
	baseFeeMA  *big.Int
	baseFees   [MovingAverageWindow]*big.Int
	blobFees   [MovingAverageWindow]*big.Int
	maIndex    int
}

func NewGasOracle(l1ChainCfg *params.ChainConfig) Oracle {
	return &oracle{
		l1ChainCfg: l1ChainCfg,
		baseFeeMA:  big.NewInt(0),
		blobFeeMA:  big.NewInt(0),
		baseFees:   [MovingAverageWindow]*big.Int{},
		blobFees:   [MovingAverageWindow]*big.Int{},
		maIndex:    0,
	}
}

// EstimateL1StorageGasCost - Returns the expected l1 gas cost for a transaction at a given l1 block.
func (o *oracle) EstimateL1StorageGasCost(tx *types.Transaction, header *common.BatchHeader) (*big.Int, error) {
	encodedTx, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return nil, err
	}

	return o.calculateL1Cost(header, encodedTx)
}

func (o *oracle) EstimateL1CostForMsg(args *gethapi.TransactionArgs, header *common.BatchHeader) (*big.Int, error) {
	encoded, err := rlp.EncodeToBytes(args)
	if err != nil {
		return nil, err
	}

	return o.calculateL1Cost(header, encoded)
}

func (o *oracle) SubmitL1Block(block *types.Header) error {
	blobFeePerByte := big.NewInt(0)
	if block.ExcessBlobGas != nil {
		blobFeePerByte = eip4844.CalcBlobFee(o.l1ChainCfg, block)
	}
	o.updateMA(block.BaseFee, blobFeePerByte)
	return nil
}

func (o *oracle) updateMA(baseFee *big.Int, blobFeePerByte *big.Int) {
	o.baseFees[o.maIndex] = baseFee
	o.blobFees[o.maIndex] = blobFeePerByte

	o.maIndex = (o.maIndex + 1) % MovingAverageWindow

	baseFeeSum := big.NewInt(0)
	count := 0
	for _, fee := range o.baseFees {
		if fee != nil {
			baseFeeSum.Add(baseFeeSum, fee)
			count++
		}
	}

	blobFeeSum := big.NewInt(0)
	for _, fee := range o.blobFees {
		if fee != nil {
			blobFeeSum.Add(blobFeeSum, fee)
			count++
		}
	}

	if count > 0 {
		o.baseFeeMA = new(big.Int).Div(baseFeeSum, big.NewInt(int64(count)))
		o.blobFeeMA = new(big.Int).Div(blobFeeSum, big.NewInt(int64(count)))
	}
}

// calculateL1Cost - Calculates the L1 cost as a multiple of the L2 base fee.
// it takes into account the share of the blob cost and the share of the L1 TX cost - which submits and stores the rollup header.
func (o *oracle) calculateL1Cost(l2Batch *common.BatchHeader, encodedTx []byte) (*big.Int, error) {
	totalCost := big.NewInt(0)

	// 1. Calculate the cost of including the tx in a blob
	// price in Wei for a single unit of blob
	txL1Size := CalculateL1Size(encodedTx)
	shareOfBlobCost := big.NewInt(0).Mul(txL1Size, o.blobFeeMA)

	// 2. Estimate how much this tx should absorb from the L1 tx cost that submits the rollup
	shareOfL1TxGas := big.NewInt(L1TxGas / TxsPerRollup)
	// todo - use a moving average for the L1 base fee
	shareOfL1TxCost := big.NewInt(0).Mul(shareOfL1TxGas, o.baseFeeMA)

	// 3. The total cost is the sum of the share of the blob cost and the share of the L1 tx cost
	totalCost.Add(shareOfBlobCost, shareOfL1TxCost)

	// 4. round the shareOfBlobCost up to the nearest multiple of l2Batch.BaseFee
	remainder := new(big.Int).Mod(totalCost, l2Batch.BaseFee)
	if remainder.Sign() > 0 {
		totalCost.Add(totalCost, new(big.Int).Sub(l2Batch.BaseFee, remainder))
	}

	return totalCost, nil
}
