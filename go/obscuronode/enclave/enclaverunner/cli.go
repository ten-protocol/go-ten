package enclaverunner

import (
	"flag"
)

const (
	// Flag names, defaults and usages.
	nodeIDName  = "nodeID"
	nodeIDUsage = "A integer representing the 20 bytes of the node's address (default 1)"

	addressName  = "address"
	addressUsage = "The address on which to serve the Obscuro enclave service"

	writeToLogsName  = "writeToLogs"
	writeToLogsUsage = "Whether to redirect the output to the log file."

	contractAddrName  = "contractAddress"
	contractAddrUsage = "The management contract address on the L1"

	logPathName  = "logPath"
	logPathUsage = "The path to use for the enclave service's log file"

	verifyL1BlocksName  = "verifyL1Blocks"
	verifyL1BlocksUsage = "Whether to verify incoming blocks using the hardcoded L1 genesis.json config"
)

type EnclaveConfig struct {
	NodeID          int64
	Address         string
	ContractAddress string
	WriteToLogs     bool
	LogPath         string
	VerifyL1Blocks  bool
}

func DefaultEnclaveConfig() EnclaveConfig {
	return EnclaveConfig{
		NodeID:          1,
		Address:         "localhost:11000",
		ContractAddress: "",
		WriteToLogs:     false,
		LogPath:         "enclave_logs.txt",
		VerifyL1Blocks:  false,
	}
}

func ParseCLIArgs() EnclaveConfig {
	defaultConfig := DefaultEnclaveConfig()

	nodeID := flag.Int64(nodeIDName, defaultConfig.NodeID, nodeIDUsage)
	port := flag.String(addressName, defaultConfig.Address, addressUsage)
	writeToLogs := flag.Bool(writeToLogsName, defaultConfig.WriteToLogs, writeToLogsUsage)
	contractAddress := flag.String(contractAddrName, defaultConfig.ContractAddress, contractAddrUsage)
	logPath := flag.String(logPathName, defaultConfig.LogPath, logPathUsage)
	verifyL1Blocks := flag.Bool(verifyL1BlocksName, defaultConfig.VerifyL1Blocks, verifyL1BlocksUsage)

	flag.Parse()

	return EnclaveConfig{
		NodeID:          *nodeID,
		Address:         *port,
		ContractAddress: *contractAddress,
		WriteToLogs:     *writeToLogs,
		LogPath:         *logPath,
		VerifyL1Blocks:  *verifyL1Blocks,
	}
}
