package gas

import (
	"math/big"

	"github.com/ethereum/go-ethereum/params"
)

// CalculateL1GasUsed - calculates the gas cost of having a transaction on the l1.
func CalculateL1GasUsed(data []byte, overhead *big.Int) *big.Int {
	reducedTxSize := uint64(len(data))
	reducedTxSize = (reducedTxSize * 90) / 100
	reducedTxSize = reducedTxSize * params.TxDataNonZeroGasEIP2028

	l1Gas := new(big.Int).SetUint64(reducedTxSize)
	return new(big.Int).Add(l1Gas, overhead)
}
