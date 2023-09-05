package gas

import (
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

type Oracle interface {
	ProcessL1Block(block *types.Block)
	GetGasCostForTx(tx *types.Transaction, block *types.Block) (*big.Int, error)
}

type oracle struct {
	baseFee *big.Int
}

func NewGasOracle() Oracle {
	return &oracle{
		baseFee: big.NewInt(1),
	}
}

func (o *oracle) ProcessL1Block(block *types.Block) {
	blockBaseFee := block.BaseFee()
	if blockBaseFee != nil {
		o.baseFee = blockBaseFee
	}
}

func (o *oracle) GetGasCostForTx(tx *types.Transaction, block *types.Block) (*big.Int, error) {
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
