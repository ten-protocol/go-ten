package common

import "math/big"

var Wei = big.NewInt(10 ^ 18)

func ToWei(amt uint64) *big.Int {
	bigAmt := big.NewInt(int64(amt))
	return bigAmt.Mul(bigAmt, Wei)
}
