package main

import (
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/tools/contractdeployer"
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
