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
	nodeHostDefault = "testnet.obscu.ro"
	nodeHostUsage   = "The host on which to connect to the Obscuro node. Default: `testnet.obscu.ro`."

	nodeHTTPPortName    = "nodePort"
	nodeHTTPPortDefault = 13000
	nodeHTTPPortUsage   = "The port on which to connect to the Obscuro node via RPC over HTTP. Default: 13000 ."

	faucetPKName    = "pk"
	faucetPKDefault = ""
	faucetPKUsage   = "The prefunded PK used to fund other accounts. No default, must be set."

	jwtSecretName    = "jwtSecret"
	jwtSecretDefault = ""
	jwtSecretUsage   = "The jwt request secret string. No default, must be set." //nolint: gosec
)

func parseCLIArgs() *faucet.Config {
	faucetPort := flag.Int(faucetPortName, faucetPortDefault, faucetPortUsage)
	nodeHost := flag.String(nodeHostName, nodeHostDefault, nodeHostUsage)
	nodeHTTPPort := flag.Int(nodeHTTPPortName, nodeHTTPPortDefault, nodeHTTPPortUsage)
	faucetPK := flag.String(faucetPKName, faucetPKDefault, faucetPKUsage)
	jwtSecret := flag.String(jwtSecretName, jwtSecretDefault, jwtSecretUsage)
	flag.Parse()

	return &faucet.Config{
		Port:      *faucetPort,
		Host:      *nodeHost,
		HTTPPort:  *nodeHTTPPort,
		PK:        *faucetPK,
		JWTSecret: *jwtSecret,
		ChainID:   big.NewInt(777), // TODO make this configurable
	}
}
