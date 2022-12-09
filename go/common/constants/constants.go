package constants

import "math/big"

// TODO this package is a temporary structure to ensure code wide constants are gathered before being properly decoupled

var (
	DefaultGasPrice = big.NewInt(2000000000)
	DefaultGasLimit = uint64(1025_000_000)
)
