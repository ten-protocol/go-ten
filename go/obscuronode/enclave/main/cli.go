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

	writeToLogsName    = "writeToLogs"
	writeToLogsDefault = false
	writeToLogsUsage   = "Whether to redirect the output to the log file."
)

type enclaveConfig struct {
	nodeID      *string
	address     *string
	writeToLogs *bool
}

func parseCLIArgs() enclaveConfig {
	nodeID := flag.String(nodeIDName, nodeIDDefault, nodeIDUsage)
	port := flag.String(addressName, addressDefault, addressUsage)
	writeToLogs := flag.Bool(writeToLogsName, writeToLogsDefault, writeToLogsUsage)
	flag.Parse()

	return enclaveConfig{nodeID, port, writeToLogs}
}
