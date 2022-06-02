package main

import (
	"flag"
	"strings"

	"github.com/obscuronet/go-obscuro/tools/walletextension"
)

const (
	// Flag names, defaults and usages.
	localNetworkName    = "localNetwork"
	localNetworkDefault = true
	localNetworkUsage   = "Whether to connect the wallet extension to a new local Ethereum network"

	prefundedAccountsName    = "prefundedAccounts"
	prefundedAccountsDefault = ""
	prefundedAccountsUsage   = "The account addresses to prefund if using a local network, as a comma-separated list"

	startPortName    = "startPort"
	startPortDefault = 3000
	startPortUsage   = "The first port to allocate. Ports will be allocated incrementally from this port as needed"
)

func parseCLIArgs() walletextension.RunConfig {
	localNetwork := flag.Bool(localNetworkName, localNetworkDefault, localNetworkUsage)
	prefundedAccounts := flag.String(prefundedAccountsName, prefundedAccountsDefault, prefundedAccountsUsage)
	startPort := flag.Int(startPortName, startPortDefault, startPortUsage)
	flag.Parse()

	var parsedAccounts []string
	if len(*prefundedAccounts) == 0 {
		parsedAccounts = []string{}
	} else {
		parsedAccounts = strings.Split(*prefundedAccounts, ",")
	}

	return walletextension.RunConfig{
		LocalNetwork:      *localNetwork,
		PrefundedAccounts: parsedAccounts,
		StartPort:         *startPort,
	}
}
