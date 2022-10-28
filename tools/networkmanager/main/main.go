package main

import (
	"os"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/obscuronet/go-obscuro/go/common/log"

	"github.com/obscuronet/go-obscuro/tools/networkmanager"
)

func main() {
	config, args := networkmanager.ParseCLIArgs()
	logger := log.New(log.NetwMngCmp, int(gethlog.LvlError), log.SysOut)

	switch config.Command {
	case networkmanager.DeployMgmtContract, networkmanager.DeployERC20Contract:
		networkmanager.DeployContract(config, logger)
	case networkmanager.InjectTxs:
		networkmanager.InjectTransactions(config, args, logger)
	default:
		panic("unrecognised command type")
	}

	os.Exit(0)
}
