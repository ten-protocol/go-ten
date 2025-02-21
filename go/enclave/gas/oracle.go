package gas

import (
	"math/big"

	"github.com/ten-protocol/go-ten/go/common"

	"github.com/ethereum/go-ethereum/consensus/misc/eip4844"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ten-protocol/go-ten/go/common/gethapi"
)

// L1TxGas - a crude estimation of the cost of publishing an L1 tx
const L1TxGas = 300_000

// TxsPerBatch - the number of transactions in a batch. A conservative estimation.
const TxsPerBatch = 50

// Oracle - the interface for the future precompiled gas oracle contract
// which will expose necessary l1 information.
type Oracle interface {
	EstimateL1StorageGasCost(tx *types.Transaction, block *types.Header, header *common.BatchHeader) (*big.Int, error)
	EstimateL1CostForMsg(args *gethapi.TransactionArgs, block *types.Header, header *common.BatchHeader) (*big.Int, error)
}

type oracle struct{}

func NewGasOracle() Oracle {
	return &oracle{}
}

// EstimateL1StorageGasCost - Returns the expected l1 gas cost for a transaction at a given l1 block.
func (o *oracle) EstimateL1StorageGasCost(tx *types.Transaction, block *types.Header, header *common.BatchHeader) (*big.Int, error) {
	encodedTx, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return nil, err
	}

	return o.calculateL1Cost(block, header, encodedTx)
}

func (o *oracle) EstimateL1CostForMsg(args *gethapi.TransactionArgs, block *types.Header, header *common.BatchHeader) (*big.Int, error) {
	encoded := make([]byte, 0)
	if args.Data != nil {
		encoded = append(encoded, *args.Data...)
	}

	return o.calculateL1Cost(block, header, encoded)
}

// calculateL1Cost - Calculates the L1 cost as a multiple of the L2 base fee.
func (o *oracle) calculateL1Cost(l1Block *types.Header, l2Batch *common.BatchHeader, encodedTx []byte) (*big.Int, error) {
	totalCost := big.NewInt(0)

	// If the l1Block does not have excess blob gas, we can't estimate the cost
	if l1Block.ExcessBlobGas == nil {
		return totalCost, nil
	}

	// price in Wei for a single unit of blob
	blobFeePerByte := eip4844.CalcBlobFee(*l1Block.ExcessBlobGas)

	// Calculate teh cost of including the tx in a blob
	txL1Size := CalculateL1Size(encodedTx)
	shareOfBlobCost := big.NewInt(0).Mul(txL1Size, blobFeePerByte)

	// Add a value to the shareOfBlobCost to account for the cost of the transaction itself
	shareOfL1TxGas := big.NewInt(L1TxGas / TxsPerBatch)
	shareOfL1TxCost := big.NewInt(0).Mul(shareOfL1TxGas, l1Block.BaseFee)

	totalCost.Add(shareOfBlobCost, shareOfL1TxCost)

	// round the shareOfBlobCost up to the nearest multiple of l2Batch.BaseFee
	remainder := new(big.Int).Mod(totalCost, l2Batch.BaseFee)
	if remainder.Sign() > 0 {
		totalCost.Add(totalCost, new(big.Int).Sub(l2Batch.BaseFee, remainder))
	}

	return totalCost, nil
}
