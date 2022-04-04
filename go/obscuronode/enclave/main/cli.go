package main

import (
	"flag"
)

const (
	// Flag names, defaults and usages.
	nodeIDName    = "nodeID"
	nodeIDDefault = ""
	nodeIDUsage   = "The 20 bytes of the node's address (default \"\")"

	addressName    = "address"
	addressDefault = ":11000"
	addressUsage   = "The address on which to serve the Obscuro enclave service"
)

type enclaveConfig struct {
	nodeID  *string
	address *string
}

func parseCLIArgs() enclaveConfig {
	nodeID := flag.String(nodeIDName, nodeIDDefault, nodeIDUsage)
	port := flag.String(addressName, addressDefault, addressUsage)
	flag.Parse()

	return enclaveConfig{nodeID, port}
}
