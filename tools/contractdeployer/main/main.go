package main

import (
	"github.com/obscuronet/obscuro-playground/tools/contractdeployer"
	"os"
)

func main() {
	config := contractdeployer.ParseCLIArgs()

	switch config.Command {
	case contractdeployer.DeployMgmtContract, contractdeployer.DeployERC20Contract:
		contractdeployer.DeployContract(config)
	default:
		panic("unrecognised command type")
	}

	os.Exit(0)
}
