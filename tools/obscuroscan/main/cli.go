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

	startPortName  = "startPort"
	startPortUsage = "The first port to allocate. Ports will be allocated incrementally from this port as needed"
)

type obscuroscanConfig struct {
	nodeID           string
	clientServerAddr string
	startPort        int
}

func defaultObscuroClientConfig() obscuroscanConfig {
	return obscuroscanConfig{
		nodeID:           "",
		clientServerAddr: "20.68.160.65:13000",
		startPort:        3000,
	}
}

func parseCLIArgs() obscuroscanConfig {
	defaultConfig := defaultObscuroClientConfig()

	nodeID := flag.String(nodeIDName, defaultConfig.nodeID, nodeIDUsage)
	clientServerAddr := flag.String(clientServerAddrName, defaultConfig.clientServerAddr, clientServerAddrUsage)
	startPort := flag.Int(startPortName, defaultConfig.startPort, startPortUsage)

	flag.Parse()

	return obscuroscanConfig{
		nodeID:           *nodeID,
		clientServerAddr: *clientServerAddr,
		startPort:        *startPort,
	}
}
