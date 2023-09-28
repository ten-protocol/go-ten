package main

import (
	"flag"
	"math/big"

	"github.com/obscuronet/go-obscuro/tools/faucet/faucet"
)

const (
	// Flag names, defaults and usages.
	faucetPortName    = "port"
	faucetPortDefault = 80
	faucetPortUsage   = "The port on which to serve the faucet endpoint. Default: 80."

	nodeHostName    = "nodeHost"
	nodeHostDefault = "erpc.testnet.obscu.ro"
	nodeHostUsage   = "The host on which to connect to the Obscuro node. Default: `erpc.testnet.obscu.ro`."

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
)

func parseCLIArgs() *faucet.Config {
	faucetPort := flag.Int(faucetPortName, faucetPortDefault, faucetPortUsage)
	nodeHost := flag.String(nodeHostName, nodeHostDefault, nodeHostUsage)
	nodeHTTPPort := flag.Int(nodeHTTPPortName, nodeHTTPPortDefault, nodeHTTPPortUsage)
	faucetPK := flag.String(faucetPKName, faucetPKDefault, faucetPKUsage)
	jwtSecret := flag.String(jwtSecretName, jwtSecretDefault, jwtSecretUsage)
	serverPort := flag.Int(serverPortName, serverPortDefault, serverPortUsage)
	defaultAmount := flag.Float64(defaultAmountName, defaultAmountDefault, defaultAmountUsage)
	flag.Parse()

	return &faucet.Config{
		Port:              *faucetPort,
		Host:              *nodeHost,
		HTTPPort:          *nodeHTTPPort,
		PK:                *faucetPK,
		JWTSecret:         *jwtSecret,
		ServerPort:        *serverPort,
		ChainID:           big.NewInt(443), // TODO make this configurable
		DefaultFundAmount: toWei(defaultAmount),
	}
}

func toWei(amount *float64) *big.Int {
	amtFloat := new(big.Float).SetFloat64(*amount)
	weiFloat := new(big.Float).Mul(amtFloat, big.NewFloat(1e18))
	// don't care about the accuracy here, float should have less than 18 decimal places
	wei, _ := weiFloat.Int(nil)
	return wei
}
