package container

// Flag names.
const (
	overrideFlag                  = "override"
	configFlag                    = "config"
	HostIDFlag                    = "hostID"
	HostAddressFlag               = "hostAddress"
	AddressFlag                   = "address"
	NodeTypeFlag                  = "nodeType"
	L1ChainIDFlag                 = "l1ChainID"
	TenChainIDFlag                = "tenChainID"
	WillAttestFlag                = "willAttest"
	ValidateL1BlocksFlag          = "validateL1Blocks"
	ManagementContractAddressFlag = "managementContractAddress"
	LogLevelFlag                  = "logLevel"
	LogPathFlag                   = "logPath"
	UseInMemoryDBFlag             = "useInMemoryDB"
	EdgelessDBHostFlag            = "edgelessDBHost"
	SQLiteDBPathFlag              = "sqliteDBPath"
	ProfilerEnabledFlag           = "profilerEnabled"
	MinGasPriceFlag               = "minGasPrice"
	MessageBusAddressFlag         = "messageBusAddress"
	SequencerIDFlag               = "sequencerID"
	TenGenesisFlag                = "tenGenesis"
	DebugNamespaceEnabledFlag     = "debugNamespaceEnabled"
	MaxBatchSizeFlag              = "maxBatchSize"
	MaxRollupSizeFlag             = "maxRollupSize"
	L2BaseFeeFlag                 = "l2BaseFee"
	L2CoinbaseFlag                = "l2Coinbase"
	GasBatchExecutionLimit        = "gasBatchExecutionLimit"
	GasLocalExecutionCapFlag      = "gasLocalExecutionCap"
)

// EnclaveRestrictedFlags are the flags that the enclave can receive ONLY over (a) the Ego signed enclave.json (./go/enclave/main/enclave.json)
// or (b) if passed via EDG_<flag> as environment variable (see https://docs.edgeless.systems/ego/reference/config#environment-variables).
// In the case of running enclave as standalone process (./go/enclave/main/main w/o ego sign enclave.json stage) these flags will be checked
// to be set via EDG_<flag> env vars.
var enclaveRestrictedFlags = map[string]string{
	L1ChainIDFlag:             "int64",
	TenChainIDFlag:            "int64",
	TenGenesisFlag:            "string",
	UseInMemoryDBFlag:         "bool",
	ProfilerEnabledFlag:       "bool",
	DebugNamespaceEnabledFlag: "bool",
}

func getFlagUsageMap() map[string]string {
	return map[string]string{
		overrideFlag:                  "Additive config file to apply on top of default or -config",
		configFlag:                    "The path to the host's config file. Overrides all other flags",
		HostIDFlag:                    "The 20 bytes of the address of the TEN host this enclave serves",
		HostAddressFlag:               "The peer-to-peer IP address of the TEN host this enclave serves",
		AddressFlag:                   "The address on which to serve the TEN enclave service",
		NodeTypeFlag:                  "The node's type (e.g. sequencer, validator)",
		WillAttestFlag:                "Whether the enclave will produce a verified attestation report",
		ValidateL1BlocksFlag:          "Whether to validate incoming blocks using the hardcoded L1 genesis.json config",
		ManagementContractAddressFlag: "The management contract address on the L1",
		LogLevelFlag:                  "The verbosity level of logs. (Defaults to Info)",
		LogPathFlag:                   "The path to use for the enclave service's log file",
		EdgelessDBHostFlag:            "Host address for the edgeless DB instance (can be empty if useInMemoryDB is true or if not using attestation",
		SQLiteDBPathFlag:              "Filepath for the sqlite DB persistence file (can be empty if a throwaway file in /tmp/ is acceptable or if using InMemory DB or if using attestation/EdgelessDB)",
		MinGasPriceFlag:               "The minimum gas price for mining a transaction",
		MessageBusAddressFlag:         "The address of the L1 message bus contract owned by the management contract.",
		SequencerIDFlag:               "The 20 bytes of the address of the sequencer for this network",
		MaxBatchSizeFlag:              "The maximum size a batch is allowed to reach uncompressed",
		MaxRollupSizeFlag:             "The maximum size a rollup is allowed to reach",
		L2BaseFeeFlag:                 "Base gas fee",
		L2CoinbaseFlag:                "Account used for gas payments of L1 transactions",
		GasBatchExecutionLimit:        "Max gas that can be executed in a single batch",
		TenGenesisFlag:                "The json string with the TEN genesis",
		L1ChainIDFlag:                 "An integer representing the unique chain id of the Ethereum chain used as an L1 (default 1337)",
		TenChainIDFlag:                "An integer representing the unique chain id of the TEN chain (default 443)",
		UseInMemoryDBFlag:             "Whether the enclave will use an in-memory DB rather than persist data",
		ProfilerEnabledFlag:           "Runs a profiler instance (Defaults to false)",
		DebugNamespaceEnabledFlag:     "Whether the debug namespace is enabled",
		GasLocalExecutionCapFlag:      "Max gas usage when executing local transactions",
	}
}
