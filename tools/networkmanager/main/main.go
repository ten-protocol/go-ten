package main

import (
	"os"

	"github.com/obscuronet/obscuro-playground/tools/networkmanager"
)

func main() {
	config := networkmanager.ParseCLIArgs()

	switch config.Command {
	case networkmanager.DeployMgmtContract, networkmanager.DeployERC20Contract:
		networkmanager.DeployContract(config)
	default:
		panic("unrecognised command type")
	}

	os.Exit(0)
}
