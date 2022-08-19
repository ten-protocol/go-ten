package main

import (
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/tools/contractdeployer"
)

func main() {
	log.SetLogLevel(log.DisabledLevel)
	config := contractdeployer.ParseConfig()
	_, err := contractdeployer.Deploy(config)
	if err != nil {
		// the contract deployer's output is to be consumed by other applications
		// in case of a failure bump the log level and panic
		log.SetLogLevel(log.TraceLevel)
		log.Panic("%s", err)
	}
}
