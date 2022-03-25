package main

import (
	"flag"
)

const (
	// Flag names, defaults and usages.
	nodeIDName    = "nodeAddress"
	nodeIDDefault = ""
	nodeIDUsage   = "The 20 bytes of the node's address (default \"\")"

	portName    = "port"
	portDefault = 11000
	portUsage   = "The port on which to serve the Obscuro enclave service"
)

type enclaveConfig struct {
	nodeID *string
	port   *uint64
}

func parseCLIArgs() enclaveConfig {
	nodeID := flag.String(nodeIDName, nodeIDDefault, nodeIDUsage)
	port := flag.Uint64(portName, portDefault, portUsage)
	flag.Parse()

	return enclaveConfig{nodeID, port}
}
