package main

import (
	"flag"
	"math/big"

	"github.com/ethereum/go-ethereum/params"

	"github.com/ten-protocol/go-ten/tools/faucet/faucet"
)

const (
	// Flag names, defaults and usages.
	nodeHostName    = "nodeHost"
	nodeHostDefault = "erpc.sepolia-testnet.ten.xyz"
	nodeHostUsage   = "The host on which to connect to the Obscuro node. Default: `erpc.sepolia-testnet.ten.xyz`."

	nodeHTTPPortName    = "nodePort"
	nodeHTTPPortDefault = 80
	nodeHTTPPortUsage   = "The port on which to connect to the Obscuro node via RPC over HTTP. Default: 80 ."

	faucetPKName    = "pk"
	faucetPKDefault = ""
	faucetPKUsage   = "The prefunded PK used to fund other accounts. No default, must be set."

	jwtSecretName    = "jwtSecret"
	jwtSecretDefault = ""
	jwtSecretUsage   = "The jwt request secret string. No default, must be set." //nolint: gosec

	serverPortName    = "serverPort"
	serverPortDefault = 80
	serverPortUsage   = "Port where the web server binds to"

	defaultAmountName    = "defaultAmount"
	defaultAmountDefault = 100.0
	defaultAmountUsage   = "Default amount of token to fund (in ETH)"

	chainIDName    = "chainID"
	chainIDDefault = 443
	chainIDUsage   = "The chain ID to use for transactions. Default: 443."
)

func parseCLIArgs() *faucet.Config {
	nodeHost := flag.String(nodeHostName, nodeHostDefault, nodeHostUsage)
	nodeHTTPPort := flag.Int(nodeHTTPPortName, nodeHTTPPortDefault, nodeHTTPPortUsage)
	faucetPK := flag.String(faucetPKName, faucetPKDefault, faucetPKUsage)
	jwtSecret := flag.String(jwtSecretName, jwtSecretDefault, jwtSecretUsage)
	serverPort := flag.Int(serverPortName, serverPortDefault, serverPortUsage)
	defaultAmount := flag.Float64(defaultAmountName, defaultAmountDefault, defaultAmountUsage)
	chainID := flag.Int64(chainIDName, chainIDDefault, chainIDUsage)
	flag.Parse()

	return &faucet.Config{
		Host:              *nodeHost,
		HTTPPort:          *nodeHTTPPort,
		PK:                *faucetPK,
		JWTSecret:         *jwtSecret,
		ServerPort:        *serverPort,
		ChainID:           big.NewInt(*chainID),
		DefaultFundAmount: toWei(defaultAmount),
	}
}

func toWei(amount *float64) *big.Int {
	amtFloat := new(big.Float).SetFloat64(*amount)
	weiFloat := new(big.Float).Mul(amtFloat, big.NewFloat(params.Ether))
	// don't care about the accuracy here, float should have less than 18 decimal places
	wei, _ := weiFloat.Int(nil)
	return wei
}
