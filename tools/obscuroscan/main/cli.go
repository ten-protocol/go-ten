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

	logPathName  = "logPath"
	logPathUsage = "The path to use for Obscuroscan's log file"
)

type obscuroscanConfig struct {
	nodeID        string
	rpcServerAddr string
	address       string
	logPath       string
}

func defaultObscuroClientConfig() obscuroscanConfig {
	return obscuroscanConfig{
		nodeID:        "",
		rpcServerAddr: "http://localhost:37400",
		address:       "127.0.0.1:3000",
		logPath:       "obscuroscan_logs.txt",
	}
}

func parseCLIArgs() obscuroscanConfig {
	defaultConfig := defaultObscuroClientConfig()

	nodeID := flag.String(nodeIDName, defaultConfig.nodeID, nodeIDUsage)
	rpcServerAddr := flag.String(rpcServerAddrName, defaultConfig.rpcServerAddr, rpcServerAddrUsage)
	address := flag.String(addressName, defaultConfig.address, addressUsage)
	logPath := flag.String(logPathName, defaultConfig.logPath, logPathUsage)

	flag.Parse()

	return obscuroscanConfig{
		nodeID:        *nodeID,
		rpcServerAddr: *rpcServerAddr,
		address:       *address,
		logPath:       *logPath,
	}
}
