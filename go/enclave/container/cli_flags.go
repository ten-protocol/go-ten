package container

// Flag names.
const (
	configName                    = "config"
	hostIDName                    = "hostID"
	hostAddressName               = "hostAddress"
	addressName                   = "address"
	nodeTypeName                  = "nodeType"
	l1ChainIDName                 = "l1ChainID"
	obscuroChainIDName            = "obscuroChainID"
	willAttestName                = "willAttest"
	validateL1BlocksName          = "validateL1Blocks"
	ManagementContractAddressName = "managementContractAddress"
	logLevelName                  = "logLevel"
	logPathName                   = "logPath"
	useInMemoryDBName             = "useInMemoryDB"
	edgelessDBHostName            = "edgelessDBHost"
	sqliteDBPathName              = "sqliteDBPath"
	profilerEnabledName           = "profilerEnabled"
	minGasPriceName               = "minGasPrice"
	messageBusAddressName         = "messageBusAddress"
	sequencerIDName               = "sequencerID"
	obscuroGenesisName            = "obscuroGenesis"
)

// Returns a map of the flag usages.
// While we could just use constants instead of a map, this approach allows us to test that all the expected flags are defined.
func getFlagUsageMap() map[string]string {
	return map[string]string{
		configName:                    "The path to the node's config file. Overrides all other flags",
		hostIDName:                    "The 20 bytes of the address of the Obscuro host this enclave serves",
		hostAddressName:               "The peer-to-peer IP address of the Obscuro host this enclave serves",
		addressName:                   "The address on which to serve the Obscuro enclave service",
		nodeTypeName:                  "The node's type (e.g. sequencer, validator)",
		l1ChainIDName:                 "An integer representing the unique chain id of the Ethereum chain used as an L1 (default 1337)",
		obscuroChainIDName:            "An integer representing the unique chain id of the Obscuro chain (default 777)",
		willAttestName:                "Whether the enclave will produce a verified attestation report",
		validateL1BlocksName:          "Whether to validate incoming blocks using the hardcoded L1 genesis.json config",
		ManagementContractAddressName: "The management contract address on the L1",
		logLevelName:                  "The verbosity level of logs. (Defaults to Info)",
		logPathName:                   "The path to use for the enclave service's log file",
		useInMemoryDBName:             "Whether the enclave will use an in-memory DB rather than persist data",
		edgelessDBHostName:            "Host address for the edgeless DB instance (can be empty if useInMemoryDB is true or if not using attestation",
		sqliteDBPathName:              "Filepath for the sqlite DB persistence file (can be empty if a throwaway file in /tmp/ is acceptable or if using InMemory DB or if using attestation/EdgelessDB)",
		profilerEnabledName:           "Runs a profiler instance (Defaults to false)",
		minGasPriceName:               "The minimum gas price for mining a transaction",
		messageBusAddressName:         "The address of the L1 message bus contract owned by the management contract.",
		sequencerIDName:               "The 20 bytes of the address of the sequencer for this network",
		obscuroGenesisName:            "The json string with the obscuro genesis",
	}
}
