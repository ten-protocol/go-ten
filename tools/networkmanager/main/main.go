package main

import (
	"os"

	"github.com/obscuronet/go-obscuro/tools/networkmanager"
)

func main() {
	config, args := networkmanager.ParseCLIArgs()

	switch config.Command {
	case networkmanager.DeployMgmtContract, networkmanager.DeployERC20Contract:
		networkmanager.DeployContract(config)
	case networkmanager.InjectTxs:
		networkmanager.InjectTransactions(config, args)
	default:
		panic("unrecognised command type")
	}

	os.Exit(0)
}
