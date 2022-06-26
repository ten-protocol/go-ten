package main

import (
	"github.com/obscuronet/obscuro-playground/go/log"
	"github.com/obscuronet/obscuro-playground/tools/contractdeployer"
)

func main() {
	log.SetLogLevel(log.DisabledLevel)
	config := contractdeployer.ParseConfig()
	deployer := contractdeployer.NewContractDeployer(config)
	err := deployer.Run()
	if err != nil {
		log.SetLogLevel(log.TraceLevel)
		log.Panic("%w", err)
	}
}
