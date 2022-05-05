package main

import (
	"flag"
	"strings"

	"github.com/obscuronet/obscuro-playground/tools/walletextension"
)

const (
	// Flag names, defaults and usages.
	localNetworkName    = "localNetwork"
	localNetworkDefault = false
	localNetworkUsage   = "Whether to connect the wallet extension to a new local Ethereum network"

	prefundedAccountsName    = "prefundedAccounts"
	prefundedAccountsDefault = ""
	prefundedAccountsUsage   = "The account addresses to prefund if using a local network, as a comma-separated list"
)

func parseCLIArgs() walletextension.RunConfig {
	localNetwork := flag.Bool(localNetworkName, localNetworkDefault, localNetworkUsage)
	prefundedAccounts := flag.String(prefundedAccountsName, prefundedAccountsDefault, prefundedAccountsUsage)
	flag.Parse()

	return walletextension.RunConfig{
		LocalNetwork:      *localNetwork,
		PrefundedAccounts: strings.Split(*prefundedAccounts, ","),
	}
}
