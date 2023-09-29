package faucet

import "math/big"

type Config struct {
	Port              int
	Host              string
	HTTPPort          int
	PK                string
	JWTSecret         string
	ChainID           *big.Int
	ServerPort        int
	DefaultFundAmount *big.Int // how much token to fund by default (in wei)
}
