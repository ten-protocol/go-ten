package main

import (
	"fmt"
	"math/big"
	"time"

	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/tools/contractdeployer"
)

func main() {
	log.SetLogLevel(log.DisabledLevel)
	// config := contractdeployer.ParseConfig()
	config := &contractdeployer.Config{
		NodeHost:       "obscuronode-0-testnet-17.uksouth.cloudapp.azure.com",
		NodePort:       13000,
		IsL1Deployment: false,
		PrivateKey:     "4bfe14725e685901c062ccd4e220c61cf9c189897b6c78bd18d7f51291b2b8f8",
		ChainID:        big.NewInt(777),
		ContractName:   "L2ERC20",
		ConstructorParams: []string{
			"Matt",
			"MAT",
			"1000000000000000000000000000000",
		},
	}
	contractAddr, err := contractdeployer.Deploy(config)
	if err != nil {
		// the contract deployer's output is to be consumed by other applications
		// in case of a failure bump the log level and panic
		log.SetLogLevel(log.TraceLevel)
		log.Panic("%s", err)
	}
	// print the contract address, to be read if necessary by the caller (important: this must be the last message output by the script)
	fmt.Print(contractAddr)

	// this is a safety sleep to make sure the output is printed from the docker container
	time.Sleep(5 * time.Second)
}
