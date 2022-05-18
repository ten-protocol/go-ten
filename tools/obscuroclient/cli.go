package main

import (
	"flag"
)

const (
	// Flag names, defaults and usages.
	nodeIDName  = "nodeID"
	nodeIDUsage = "The 20 bytes of the node's address"

	clientServerAddrName  = "clientServerAddress"
	clientServerAddrUsage = "The address on which to send RPC requests"
)

type obscuroClientConfig struct {
	nodeID           string
	clientServerAddr string
}

func defaultObscuroClientConfig() obscuroClientConfig {
	return obscuroClientConfig{
		nodeID:           "",
		clientServerAddr: "20.68.160.65:13000",
	}
}

func parseCLIArgs() obscuroClientConfig {
	defaultConfig := defaultObscuroClientConfig()

	nodeID := flag.String(nodeIDName, defaultConfig.nodeID, nodeIDUsage)
	clientServerAddr := flag.String(clientServerAddrName, defaultConfig.clientServerAddr, clientServerAddrUsage)

	flag.Parse()

	return obscuroClientConfig{
		nodeID:           *nodeID,
		clientServerAddr: *clientServerAddr,
	}
}
