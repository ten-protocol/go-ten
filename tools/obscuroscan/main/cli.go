package main

import (
	"flag"
)

const (
	// Flag names and usages.
	nodeIDName  = "nodeID"
	nodeIDUsage = "The 20 bytes of the node's address"

	rpcServerAddrName  = "rpcServerAddress"
	rpcServerAddrUsage = "The address on which to send RPC requests"

	addressName  = "address"
	addressUsage = "The address to serve Obscuroscan on"
)

type obscuroscanConfig struct {
	nodeID        string
	rpcServerAddr string
	address       string
}

func defaultObscuroClientConfig() obscuroscanConfig {
	return obscuroscanConfig{
		nodeID:        "",
		rpcServerAddr: "127.0.0.1:37400",
		address:       "127.0.0.1:3000",
	}
}

func parseCLIArgs() obscuroscanConfig {
	defaultConfig := defaultObscuroClientConfig()

	nodeID := flag.String(nodeIDName, defaultConfig.nodeID, nodeIDUsage)
	rpcServerAddr := flag.String(rpcServerAddrName, defaultConfig.rpcServerAddr, rpcServerAddrUsage)
	address := flag.String(addressName, defaultConfig.address, addressUsage)

	flag.Parse()

	return obscuroscanConfig{
		nodeID:        *nodeID,
		rpcServerAddr: *rpcServerAddr,
		address:       *address,
	}
}
