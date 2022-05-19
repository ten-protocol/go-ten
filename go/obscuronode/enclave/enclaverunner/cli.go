package enclaverunner

import (
	"flag"
	"strings"

	"github.com/ethereum/go-ethereum/common"
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

	erc20contractAddrsName  = "erc20contractAddresses"
	erc20contractAddrsUsage = "The erc20 contract addresses to monitor on the L1"
)

type EnclaveConfig struct {
	NodeID             int64
	Address            string
	ContractAddress    string
	WriteToLogs        bool
	LogPath            string
	ERC20ContractAddrs []*common.Address
}

func DefaultEnclaveConfig() EnclaveConfig {
	return EnclaveConfig{
		NodeID:             1,
		Address:            "localhost:11000",
		ContractAddress:    "",
		WriteToLogs:        false,
		LogPath:            "enclave_logs.txt",
		ERC20ContractAddrs: []*common.Address{},
	}
}

func ParseCLIArgs() EnclaveConfig {
	defaultConfig := DefaultEnclaveConfig()

	nodeID := flag.Int64(nodeIDName, defaultConfig.NodeID, nodeIDUsage)
	port := flag.String(addressName, defaultConfig.Address, addressUsage)
	writeToLogs := flag.Bool(writeToLogsName, defaultConfig.WriteToLogs, writeToLogsUsage)
	contractAddress := flag.String(contractAddrName, defaultConfig.ContractAddress, contractAddrUsage)
	logPath := flag.String(logPathName, defaultConfig.LogPath, logPathUsage)
	erc20ContractAddrs := flag.String(erc20contractAddrsName, "", erc20contractAddrsUsage)

	flag.Parse()

	parsedERC20ContractAddrs := strings.Split(*erc20ContractAddrs, ",")
	erc20contractAddresses := make([]*common.Address, len(parsedERC20ContractAddrs))
	if *erc20ContractAddrs != "" {
		for i, addr := range parsedERC20ContractAddrs {
			hexAddr := common.HexToAddress(addr)
			erc20contractAddresses[i] = &hexAddr
		}
	}

	return EnclaveConfig{
		NodeID:          *nodeID,
		Address:         *port,
		ContractAddress: *contractAddress,
		WriteToLogs:     *writeToLogs,
		LogPath:         *logPath,
	}
}
