package gas

import (
	"math/big"

	"github.com/ten-protocol/go-ten/go/common"

	"github.com/ethereum/go-ethereum/consensus/misc/eip4844"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ten-protocol/go-ten/go/common/gethapi"
)

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

	return o.calculateL1Cost(block, encodedTx, header)
}

func (o *oracle) EstimateL1CostForMsg(args *gethapi.TransactionArgs, block *types.Header, header *common.BatchHeader) (*big.Int, error) {
	encoded := make([]byte, 0)
	if args.Data != nil {
		encoded = append(encoded, *args.Data...)
	}

	return o.calculateL1Cost(block, encoded, header)
}

func (o *oracle) calculateL1Cost(block *types.Header, encodedTx []byte, header *common.BatchHeader) (*big.Int, error) {
	// If the block does not have excess blob gas, we can't estimate the cost
	if block.ExcessBlobGas == nil {
		return big.NewInt(0), nil
	}

	blobFee := eip4844.CalcBlobFee(*block.ExcessBlobGas)

	l1Gas := CalculateL1GasUsed(encodedTx)
	gasCost := big.NewInt(0).Mul(l1Gas, blobFee)

	remainder := new(big.Int).Mod(gasCost, header.BaseFee)
	if remainder.Sign() > 0 {
		gasCost.Add(gasCost, new(big.Int).Sub(header.BaseFee, remainder))
	}

	return gasCost, nil
}
