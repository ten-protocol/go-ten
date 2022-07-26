package main

import (
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/tools/contractdeployer"
)

func main() {
	log.SetLogLevel(log.DisabledLevel)
	config := contractdeployer.ParseConfig()
	deployer, err := contractdeployer.NewContractDeployer(config)
	if err != nil {
		// todo: why is this log level stuff setup in this way (why not print here or use logs everywhere)
		log.SetLogLevel(log.TraceLevel)
		log.Panic("%s", err)
	}
	if err := deployer.Run(); err != nil {
		log.SetLogLevel(log.TraceLevel)
		log.Panic("%s", err)
	}
}
