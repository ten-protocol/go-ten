package gas

import (
	"math/big"
)

// we choose a very conservative compression factor of 10%
// in practice, for most transactions, it will be much better
const compressionFactor = 90

// CalculateL1GasUsed - calculates the gas cost of having a transaction on the l1.
func CalculateL1GasUsed(data []byte) *big.Int {
	compressedSize := (uint64(len(data)) * compressionFactor) / 100
	return new(big.Int).SetUint64(compressedSize)
}
