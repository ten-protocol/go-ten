package main

import (
	"flag"
)

const (
	// Flag names, defaults and usages.
	localNetworkName    = "localNetwork"
	localNetworkDefault = false
	localNetworkUsage   = "Whether to connect the wallet extension to a new local Ethereum network"
)

type walletExtensionConfig struct {
	localNetwork *bool
}

func parseCLIArgs() walletExtensionConfig {
	localNetwork := flag.Bool(localNetworkName, localNetworkDefault, localNetworkUsage)
	flag.Parse()

	return walletExtensionConfig{localNetwork: localNetwork}
}
