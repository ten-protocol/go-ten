package enclaverunner

// Flag names, defaults and usages.
const (
	configName  = "config"
	configUsage = "The path to the node's config file. Overrides all other flags"

	HostIDName  = "hostID"
	hostIDUsage = "The 20 bytes of the address of the Obscuro host this enclave serves"

	AddressName  = "address"
	addressUsage = "The address on which to serve the Obscuro enclave service"

	l1ChainIDName  = "l1ChainID"
	l1ChainIDUsage = "An integer representing the unique chain id of the Ethereum chain used as an L1(default 1337)"

	obscuroChainIDName  = "obscuroChainID"
	obscuroChainIDUsage = "An integer representing the unique chain id of the Obscuro chain (default 777)"

	willAttestName  = "willAttest"
	willAttestUsage = "Whether the enclave will produce a verified attestation report"

	validateL1BlocksName  = "validateL1Blocks"
	validateL1BlocksUsage = "Whether to validate incoming blocks using the hardcoded L1 genesis.json config"

	speculativeExecutionName  = "speculativeExecution"
	speculativeExecutionUsage = "Whether to enable speculative execution"

	ManagementContractAddressName  = "managementContractAddress"
	managementContractAddressUsage = "The management contract address on the L1"

	Erc20ContractAddrsName  = "erc20ContractAddresses"
	erc20ContractAddrsUsage = "The ERC20 contract addresses to monitor on the L1"

	writeToLogsName  = "writeToLogs"
	writeToLogsUsage = "Whether to redirect the output to the log file."

	logPathName  = "logPath"
	logPathUsage = "The path to use for the enclave service's log file"

	useInMemoryDBName  = "useInMemoryDB"
	useInMemoryDBUsage = "Whether the enclave will use an in-memory DB rather than persist data"
)
