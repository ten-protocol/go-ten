package faucet

import "math/big"

type Config struct {
	Port      int
	Host      string
	HTTPPort  int
	PK        string
	JWTSecret string
	ChainID   *big.Int
}
