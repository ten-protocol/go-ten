package gas

import (
	"math/big"

	gethcore "github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
)

type ObscuroGasPool struct {
	gp gethcore.GasPool
}

func NewObscuroGasPool(gp *gethcore.GasPool) *ObscuroGasPool {
	return &ObscuroGasPool{
		gp: gethcore.GasPool(gp.Gas()),
	}
}

func (gp *ObscuroGasPool) ForTransaction(tx *types.Transaction) (*gethcore.GasPool, error) {
	encodedTx, err := rlp.EncodeToBytes(*tx)
	if err != nil {
		return nil, err
	}

	l1Gas := CalculateL1GasUsed(encodedTx, big.NewInt(0))

	gPool := gethcore.GasPool(l1Gas.Uint64())
	return &gPool, nil
}

// CalculateL1GasUsed - calculates the gas cost of having a transaction on the l1.
func CalculateL1GasUsed(data []byte, overhead *big.Int) *big.Int {
	reducedTxSize := uint64(len(data))
	reducedTxSize = (reducedTxSize * 75) / 100
	reducedTxSize = reducedTxSize * params.TxDataNonZeroGasEIP2028

	l1Gas := new(big.Int).SetUint64(reducedTxSize)
	return new(big.Int).Add(l1Gas, overhead)
}
