package gas

import (
	"math/big"
)

// we choose a very conservative compression factor
// in practice, for most transactions, it will be much better
const compressionFactor = 60

// CalculateL1Size - calculates the size of the published transaction.
func CalculateL1Size(data []byte) *big.Int {
	compressedSize := (uint64(len(data)) * compressionFactor) / 100
	return new(big.Int).SetUint64(compressedSize)
}
