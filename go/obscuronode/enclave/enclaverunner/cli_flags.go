package enclaverunner

// Flag names, defaults and usages.
const (
	configName  = "config"
	configUsage = "The path to the node's config file. Overrides all other flags"

	HostIDName  = "hostID"
	hostIDUsage = "The 20 bytes of the address of the Obscuro host this enclave serves"

	AddressName  = "address"
	addressUsage = "The address on which to serve the Obscuro enclave service"

	chainIDName  = "chainID"
	chainIDUsage = "A integer representing the unique chain id the enclave will connect to (default 1337)"

	willAttestName  = "willAttest"
	willAttestUsage = "Whether the enclave will produce a verified attestation report"

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
