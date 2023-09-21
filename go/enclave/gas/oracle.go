package gas

import (
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

// Oracle - the interface for the future precompiled gas oracle contract
// which will expose necessary l1 information.
type Oracle interface {
	ProcessL1Block(block *types.Block)
	EstimateL1StorageGasCost(tx *types.Transaction, block *types.Block) (*big.Int, error)
}

type oracle struct {
	baseFee *big.Int
}

func NewGasOracle() Oracle {
	return &oracle{
		baseFee: big.NewInt(1),
	}
}

// ProcessL1Block - should be used to update the gas oracle. Currently does not really
// fit into phase 1 gas mechanics as the information needs to be available per block.
// would be fixed when this becomes a smart contract using the stateDB
func (o *oracle) ProcessL1Block(block *types.Block) {
	blockBaseFee := block.BaseFee()
	if blockBaseFee != nil {
		o.baseFee = blockBaseFee
	}
}

// EstimateL1StorageGasCost - Returns the expected l1 gas cost for a transaction at a given l1 block.
func (o *oracle) EstimateL1StorageGasCost(tx *types.Transaction, block *types.Block) (*big.Int, error) {
	encodedTx, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return nil, err
	}

	blockBaseFee := block.BaseFee()
	if blockBaseFee == nil {
		return big.NewInt(0), nil
	}

	l1Gas := CalculateL1GasUsed(encodedTx, big.NewInt(0))
	return big.NewInt(0).Mul(l1Gas, block.BaseFee()), nil
}
