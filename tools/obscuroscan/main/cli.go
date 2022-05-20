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

	addressName  = "address"
	addressUsage = "The address to serve Obscuroscan on"
)

type obscuroscanConfig struct {
	nodeID           string
	clientServerAddr string
	address          string
}

func defaultObscuroClientConfig() obscuroscanConfig {
	return obscuroscanConfig{
		nodeID:           "",
		clientServerAddr: "20.68.160.65:13000",
		address:          "localhost:3000",
	}
}

func parseCLIArgs() obscuroscanConfig {
	defaultConfig := defaultObscuroClientConfig()

	nodeID := flag.String(nodeIDName, defaultConfig.nodeID, nodeIDUsage)
	clientServerAddr := flag.String(clientServerAddrName, defaultConfig.clientServerAddr, clientServerAddrUsage)
	address := flag.String(addressName, defaultConfig.address, addressUsage)

	flag.Parse()

	return obscuroscanConfig{
		nodeID:           *nodeID,
		clientServerAddr: *clientServerAddr,
		address:          *address,
	}
}
