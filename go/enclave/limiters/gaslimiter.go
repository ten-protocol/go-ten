package limiters

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
)

type GasLimiter struct {
	currentLimit uint64
}

func NewGasLimiter() *GasLimiter {
	return &GasLimiter{
		currentLimit: 1, // Avoid divide by zero
	}
}

func (gl *GasLimiter) ProcessBlock(block *types.Header) {
	gl.currentLimit = block.GasLimit
}

// GetCalldataLimit - returns the byte limit for all non zero bytes considering the latest block limit.
// Note that this is the worst case scenario limit. Consider that we can't just estimate based on batches
// as they get compressed and then encrypted.
func (gl *GasLimiter) GetCalldataLimit() uint64 {
	gasPerByte := params.TxDataNonZeroGasEIP2028

	// the number of bytes that fit within the current limit
	maxCalldataBytes := gl.currentLimit / gasPerByte

	return maxCalldataBytes
}
