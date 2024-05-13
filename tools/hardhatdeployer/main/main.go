package main

import (
	"fmt"
	"time"

	contractdeployer "github.com/ten-protocol/go-ten/tools/hardhatdeployer"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common/log"
)

func main() {
	config := contractdeployer.ParseConfig()
	logger := log.New(log.DeployerCmp, int(gethlog.LevelError), log.SysOut)

	contractAddr, err := contractdeployer.Deploy(config, logger)
	if err != nil {
		panic(err)
	}
	// print the contract address, to be read if necessary by the caller (important: this must be the last message output by the script)
	fmt.Print(contractAddr)

	// this is a safety sleep to make sure the output is printed from the docker container
	time.Sleep(5 * time.Second)
}
