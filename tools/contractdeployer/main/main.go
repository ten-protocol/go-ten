package main

import (
	"github.com/obscuronet/obscuro-playground/go/log"
	"github.com/obscuronet/obscuro-playground/tools/contractdeployer"
)

func main() {
	log.SetLogLevel(log.DisabledLevel)
	config := contractdeployer.ParseConfig()
	deployer := contractdeployer.NewContractDeployer(config)
	if err := deployer.Run(); err != nil {
		log.SetLogLevel(log.TraceLevel)
		log.Panic("%s", err)
	}
}
