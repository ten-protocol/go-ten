package main

import (
	"flag"
)

const (
	// Flag names, defaults and usages.
	nodeIDName  = "nodeID"
	nodeIDUsage = "The 20 bytes of the node's address"

	privateKeyName  = "privateKey"
	privateKeyUsage = "The private key for the L1 node account"

	contractAddrName  = "contractAddress"
	contractAddrUsage = "The management contract address on the L1"

	ethClientHostName  = "ethClientHost"
	ethClientHostUsage = "The host on which to connect to the Ethereum client"

	ethClientPortName  = "ethClientPort"
	ethClientPortUsage = "The port on which to connect to the Ethereum client"

	clientServerAddrName  = "clientServerAddress"
	clientServerAddrUsage = "The address on which to send RPC requests"
)

type obscuroscanConfig struct {
	nodeID           string
	clientServerAddr string
	privateKeyString string
	contractAddress  string
	ethClientHost    string
	ethClientPort    uint64
}

func defaultObscuroClientConfig() obscuroscanConfig {
	return obscuroscanConfig{
		nodeID:           "",
		clientServerAddr: "127.0.0.1:13000",
		privateKeyString: "0000000000000000000000000000000000000000000000000000000000000001",
		contractAddress:  "",
		ethClientHost:    "127.0.0.1",
		ethClientPort:    8546,
	}
}

func parseCLIArgs() obscuroscanConfig {
	defaultConfig := defaultObscuroClientConfig()

	nodeID := flag.String(nodeIDName, defaultConfig.nodeID, nodeIDUsage)
	clientServerAddr := flag.String(clientServerAddrName, defaultConfig.clientServerAddr, clientServerAddrUsage)
	privateKeyStr := flag.String(privateKeyName, defaultConfig.privateKeyString, privateKeyUsage)
	contractAddress := flag.String(contractAddrName, defaultConfig.contractAddress, contractAddrUsage)
	ethClientHost := flag.String(ethClientHostName, defaultConfig.ethClientHost, ethClientHostUsage)
	ethClientPort := flag.Uint64(ethClientPortName, defaultConfig.ethClientPort, ethClientPortUsage)

	flag.Parse()

	return obscuroscanConfig{
		nodeID:           *nodeID,
		clientServerAddr: *clientServerAddr,
		privateKeyString: *privateKeyStr,
		contractAddress:  *contractAddress,
		ethClientHost:    *ethClientHost,
		ethClientPort:    *ethClientPort,
	}
}
