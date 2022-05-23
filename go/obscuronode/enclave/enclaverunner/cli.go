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

	logPathName  = "logPath"
	logPathUsage = "The path to use for the enclave service's log file"

	contractAddrName  = "contractAddress"
	contractAddrUsage = "The management contract address on the L1"

	verifyL1BlocksName  = "verifyL1Blocks"
	verifyL1BlocksUsage = "Whether to verify incoming blocks using the hardcoded L1 genesis.json config"

	disableAttestationName  = "DisableAttestation"
	disableAttestationUsage = "Whether to disable the attestation process (use a mock attestation)."
)

type EnclaveConfig struct {
	NodeID             int64
	Address            string
	WriteToLogs        bool
	LogPath            string
	VerifyL1Blocks     bool
	ContractAddress    string
	DisableAttestation bool
}

func DefaultEnclaveConfig() EnclaveConfig {
	return EnclaveConfig{
		NodeID:             1,
		Address:            "127.0.0.1:11000",
		ContractAddress:    "",
		WriteToLogs:        false,
		LogPath:            "enclave_logs.txt",
		VerifyL1Blocks:     false,
		DisableAttestation: true,
	}
}

func ParseCLIArgs() EnclaveConfig {
	defaultConfig := DefaultEnclaveConfig()

	nodeID := flag.Int64(nodeIDName, defaultConfig.NodeID, nodeIDUsage)
	port := flag.String(addressName, defaultConfig.Address, addressUsage)
	contractAddress := flag.String(contractAddrName, defaultConfig.ContractAddress, contractAddrUsage)
	writeToLogs := flag.Bool(writeToLogsName, defaultConfig.WriteToLogs, writeToLogsUsage)
	logPath := flag.String(logPathName, defaultConfig.LogPath, logPathUsage)
	verifyL1Blocks := flag.Bool(verifyL1BlocksName, defaultConfig.VerifyL1Blocks, verifyL1BlocksUsage)
	disableAttestation := flag.Bool(disableAttestationName, defaultConfig.DisableAttestation, disableAttestationUsage)

	flag.Parse()

	return EnclaveConfig{
		NodeID:             *nodeID,
		Address:            *port,
		ContractAddress:    *contractAddress,
		WriteToLogs:        *writeToLogs,
		LogPath:            *logPath,
		VerifyL1Blocks:     *verifyL1Blocks,
		DisableAttestation: *disableAttestation,
	}
}
