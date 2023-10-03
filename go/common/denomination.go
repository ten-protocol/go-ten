package common

import "math/big"

// Analogous to denomination.go constants in geth

// These are the multipliers for ERC20 and native ETH value denominations
// Example: To get the wei value of an amount in whole tokens, use
//
//	new(big.Int).Mul(value, big.NewInt(common.Token))
var (
	Wei   = big.NewInt(1)
	GWei  = big.NewInt(1e9)
	Token = big.NewInt(1e18)
)

// ValueInWei converts a quantity of tokens (e.g. 20.5 HOC) to the integer value in wei (e.g. 2.05 x 10^19 HOC-wei)
func ValueInWei(tokenAmount *big.Int) *big.Int {
	return new(big.Int).Mul(tokenAmount, Token)
}
