package main

import (
	"flag"
	"strings"
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

type walletExtensionConfig struct {
	localNetwork      *bool
	prefundedAccounts []string
}

func parseCLIArgs() walletExtensionConfig {
	localNetwork := flag.Bool(localNetworkName, localNetworkDefault, localNetworkUsage)
	prefundedAccounts := flag.String(prefundedAccountsName, prefundedAccountsDefault, prefundedAccountsUsage)
	flag.Parse()

	return walletExtensionConfig{
		localNetwork:      localNetwork,
		prefundedAccounts: strings.Split(*prefundedAccounts, ","),
	}
}
