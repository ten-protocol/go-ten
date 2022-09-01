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
		NodeHost:       "testnet.obscu.ro",
		NodePort:       13000,
		IsL1Deployment: false,
		PrivateKey:     "6e384a07a01263518a09a5424c7b6bbfc3604ba7d93f47e3a455cbdd7f9f0682",
		ChainID:        big.NewInt(777),
		ContractName:   "L2ERC20",
		ConstructorParams: []string{
			"Mattus",
			"MATT",
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
