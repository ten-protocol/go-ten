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

	chainIDName  = "chainID"
	chainIDUsage = "A integer representing the unique chain id the enclave will connect to (default 1337)"

	addressName  = "address"
	addressUsage = "The address on which to serve the Obscuro enclave service"

	writeToLogsName  = "writeToLogs"
	writeToLogsUsage = "Whether to redirect the output to the log file."

	logPathName  = "logPath"
	logPathUsage = "The path to use for the enclave service's log file"

	erc20contractAddrsName  = "erc20contractAddresses"
	erc20contractAddrsUsage = "The erc20 contract addresses to monitor on the L1"

	contractAddrName  = "contractAddress"
	contractAddrUsage = "The management contract address on the L1"

	verifyL1BlocksName  = "verifyL1Blocks"
	verifyL1BlocksUsage = "Whether to verify incoming blocks using the hardcoded L1 genesis.json config"

	disableAttestationName  = "DisableAttestation"
	disableAttestationUsage = "Whether to disable the attestation process (use a mock attestation)."
)

type EnclaveConfig struct {
	NodeID             int64
	ChainID            int64
	Address            string
	ContractAddress    string
	WriteToLogs        bool
	LogPath            string
	ERC20ContractAddrs []*common.Address
	VerifyL1Blocks     bool
	DisableAttestation bool
}

func DefaultEnclaveConfig() EnclaveConfig {
	return EnclaveConfig{
		NodeID:             1,
		ChainID:            1337,
		Address:            "127.0.0.1:11000",
		ContractAddress:    "",
		WriteToLogs:        false,
		LogPath:            "enclave_logs.txt",
		ERC20ContractAddrs: []*common.Address{},
		VerifyL1Blocks:     false,
		DisableAttestation: true,
	}
}

func ParseCLIArgs() EnclaveConfig {
	defaultConfig := DefaultEnclaveConfig()

	nodeID := flag.Int64(nodeIDName, defaultConfig.NodeID, nodeIDUsage)
	chainID := flag.Int64(chainIDName, defaultConfig.ChainID, chainIDUsage)
	port := flag.String(addressName, defaultConfig.Address, addressUsage)
	contractAddress := flag.String(contractAddrName, defaultConfig.ContractAddress, contractAddrUsage)
	writeToLogs := flag.Bool(writeToLogsName, defaultConfig.WriteToLogs, writeToLogsUsage)
	logPath := flag.String(logPathName, defaultConfig.LogPath, logPathUsage)
	erc20ContractAddrs := flag.String(erc20contractAddrsName, "", erc20contractAddrsUsage)
	verifyL1Blocks := flag.Bool(verifyL1BlocksName, defaultConfig.VerifyL1Blocks, verifyL1BlocksUsage)
	disableAttestation := flag.Bool(disableAttestationName, defaultConfig.DisableAttestation, disableAttestationUsage)

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
		ChainID:         *chainID,
		VerifyL1Blocks:  *verifyL1Blocks,
		DisableAttestation: *disableAttestation,
	}
}
