package enclaverunner

import (
	"flag"
	"strings"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/config"

	"github.com/ethereum/go-ethereum/common"
)

const (
	// Flag names, defaults and usages.
	HostIDName  = "hostID"
	hostIDUsage = "The 20 bytes of the address of the Obscuro host this enclave serves"

	AddressName  = "address"
	addressUsage = "The address on which to serve the Obscuro enclave service"

	chainIDName  = "chainID"
	chainIDUsage = "A integer representing the unique chain id the enclave will connect to (default 1337)"

	validateL1BlocksName  = "validateL1Blocks"
	validateL1BlocksUsage = "Whether to validate incoming blocks using the hardcoded L1 genesis.json config"

	speculativeExecutionName  = "speculativeExecution"
	speculativeExecutionUsage = "Whether to enable speculative execution"

	managementContractAddressName  = "managementContractAddress"
	managementContractAddressUsage = "The management contract address on the L1"

	erc20contractAddrsName  = "erc20ContractAddresses"
	erc20contractAddrsUsage = "The ERC20 contract addresses to monitor on the L1"

	writeToLogsName  = "writeToLogs"
	writeToLogsUsage = "Whether to redirect the output to the log file."

	logPathName  = "logPath"
	logPathUsage = "The path to use for the enclave service's log file"
)

func ParseCLIArgs() config.EnclaveConfig {
	defaultConfig := config.DefaultEnclaveConfig()

	hostID := flag.String(HostIDName, defaultConfig.HostID.Hex(), hostIDUsage)
	address := flag.String(AddressName, defaultConfig.Address, addressUsage)
	chainID := flag.Int64(chainIDName, defaultConfig.ChainID, chainIDUsage)
	validateL1Blocks := flag.Bool(validateL1BlocksName, defaultConfig.ValidateL1Blocks, validateL1BlocksUsage)
	speculativeExecution := flag.Bool(speculativeExecutionName, defaultConfig.SpeculativeExecution, speculativeExecutionUsage)
	managementContractAddress := flag.String(managementContractAddressName, defaultConfig.ManagementContractAddress.Hex(), managementContractAddressUsage)
	erc20ContractAddrs := flag.String(erc20contractAddrsName, "", erc20contractAddrsUsage)
	writeToLogs := flag.Bool(writeToLogsName, defaultConfig.WriteToLogs, writeToLogsUsage)
	logPath := flag.String(logPathName, defaultConfig.LogPath, logPathUsage)

	flag.Parse()

	parsedERC20ContractAddrs := strings.Split(*erc20ContractAddrs, ",")
	erc20contractAddresses := make([]*common.Address, len(parsedERC20ContractAddrs))
	if *erc20ContractAddrs != "" {
		for i, addr := range parsedERC20ContractAddrs {
			hexAddr := common.HexToAddress(addr)
			erc20contractAddresses[i] = &hexAddr
		}
	}

	// todo - joel - add speculative execution flag

	defaultConfig.HostID = common.HexToAddress(*hostID)
	defaultConfig.Address = *address
	defaultConfig.ChainID = *chainID
	defaultConfig.ValidateL1Blocks = *validateL1Blocks
	defaultConfig.SpeculativeExecution = *speculativeExecution
	defaultConfig.ManagementContractAddress = common.HexToAddress(*managementContractAddress)
	defaultConfig.ERC20ContractAddresses = erc20contractAddresses
	defaultConfig.WriteToLogs = *writeToLogs
	defaultConfig.LogPath = *logPath

	return defaultConfig
}
