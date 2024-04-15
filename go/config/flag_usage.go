package config

// SHARED Enclave and Host flags
const (
	OverrideFlag                  = "override"
	ConfigFlag                    = "config"
	NodeNameFlag                  = "nodeName"
	IsSGXEnabledFlag              = "isSGXEnabled"
	PccsAddrFlag                  = "pccsAddr"
	HostImage                     = "hostImage"
	EnclaveImage                  = "enclaveImage"
	EdgelessDBImage               = "edgelessDBImage"
	NodeTypeFlag                  = "nodeType"
	ManagementContractAddressFlag = "managementContractAddress"
	MessageBusAddressFlag         = "messageBusAddress"
	LogLevelFlag                  = "logLevel"
	LogPathFlag                   = "logPath"
	L1ChainIDFlag                 = "l1ChainID"
	TenChainIDFlag                = "tenChainID"
	ProfilerEnabledFlag           = "profilerEnabled"
	MaxRollupSizeFlag             = "maxRollupSize"
	SQLiteDBPathFlag              = "sqliteDBPath"
	DebugNamespaceEnabledFlag     = "debugNamespaceEnabled"
	UseInMemoryDBFlag             = "useInMemoryDB"
	SequencerIDFlag               = "sequencerID"
	IsGenesisFlag                 = "isGenesis"
	ClientRPCPortHTTPFlag         = "clientRPCPortHTTP"
	ClientRPCPortWSFlag           = "clientRPCPortWS"
	ClientRPCHostFlag             = "clientRPCHost"
	EnclaveRPCAddressesFlag       = "enclaveRPCAddresses"
	P2pBindAddressFlag            = "p2pBindAddress"
	P2pPublicAddressFlag          = "p2pPublicAddress"
	L1WebsocketURLFlag            = "l1WebsocketURL"
	EnclaveRPCTimeoutFlag         = "enclaveRPCTimeout"
	L1RPCTimeoutFlag              = "l1RPCTimeout"
	P2pConnectionTimeoutFlag      = "p2pConnectionTimeout"
	PrivateKeyFlag                = "privateKey"
	L1StartHashFlag               = "l1StartHash"
	MetricsEnabledFlag            = "metricsEnabled"
	MetricsHTTPPortFlag           = "metricsHTTPPort"
	PostgresDBHostFlag            = "postgresDBHost"
	BatchIntervalFlag             = "batchInterval"
	MaxBatchIntervalFlag          = "maxBatchInterval"
	RollupIntervalFlag            = "rollupInterval"
	L1BlockTimeFlag               = "l1BlockTime"
	IsInboundP2PDisabledFlag      = "isInboundP2PDisabled"
	LevelDBPathFlag               = "levelDBPath"
	HostIDFlag                    = "hostID"
	HostAddressFlag               = "hostAddress"
	AddressFlag                   = "address"
	WillAttestFlag                = "willAttest"
	ValidateL1BlocksFlag          = "validateL1Blocks"
	EdgelessDBHostFlag            = "edgelessDBHost"
	MinGasPriceFlag               = "minGasPrice"
	GenesisJSONFlag               = "genesisJSON"
	TenGenesisFlag                = "tenGenesis"
	MaxBatchSizeFlag              = "maxBatchSize"
	L2BaseFeeFlag                 = "l2BaseFee"
	GasPaymentAddress             = "gasPaymentAddress"
	GasBatchExecutionLimit        = "gasBatchExecutionLimit"
	GasLocalExecutionCapFlag      = "gasLocalExecutionCap"
)

var FlagsByService = map[TypeConfig]map[string]bool{
	Enclave: {
		OverrideFlag: true,
		ConfigFlag:   true,
		//
		NodeTypeFlag:                  true,
		ManagementContractAddressFlag: true,
		MessageBusAddressFlag:         true,
		LogLevelFlag:                  true,
		LogPathFlag:                   true,
		L1ChainIDFlag:                 true,
		TenChainIDFlag:                true,
		ProfilerEnabledFlag:           true,
		MaxRollupSizeFlag:             true,
		SQLiteDBPathFlag:              true,
		DebugNamespaceEnabledFlag:     true,
		UseInMemoryDBFlag:             true,
		SequencerIDFlag:               true,
		//
		HostIDFlag:               true,
		HostAddressFlag:          true,
		AddressFlag:              true,
		WillAttestFlag:           true,
		ValidateL1BlocksFlag:     true,
		EdgelessDBHostFlag:       true,
		MinGasPriceFlag:          true,
		GenesisJSONFlag:          true,
		TenGenesisFlag:           true,
		MaxBatchSizeFlag:         true,
		L2BaseFeeFlag:            true,
		GasPaymentAddress:        true,
		GasBatchExecutionLimit:   true,
		GasLocalExecutionCapFlag: true,
	},
	Host: {
		OverrideFlag: true,
		ConfigFlag:   true,
		//
		NodeTypeFlag:                  true,
		ManagementContractAddressFlag: true,
		MessageBusAddressFlag:         true,
		LogLevelFlag:                  true,
		LogPathFlag:                   true,
		L1ChainIDFlag:                 true,
		TenChainIDFlag:                true,
		ProfilerEnabledFlag:           true,
		MaxRollupSizeFlag:             true,
		SQLiteDBPathFlag:              true,
		DebugNamespaceEnabledFlag:     true,
		UseInMemoryDBFlag:             true,
		SequencerIDFlag:               true,
		//
		IsGenesisFlag:            true,
		ClientRPCPortHTTPFlag:    true,
		ClientRPCPortWSFlag:      true,
		ClientRPCHostFlag:        true,
		EnclaveRPCAddressesFlag:  true,
		P2pBindAddressFlag:       true,
		P2pPublicAddressFlag:     true,
		L1WebsocketURLFlag:       true,
		EnclaveRPCTimeoutFlag:    true,
		L1RPCTimeoutFlag:         true,
		P2pConnectionTimeoutFlag: true,
		PrivateKeyFlag:           true,
		L1StartHashFlag:          true,
		MetricsEnabledFlag:       true,
		MetricsHTTPPortFlag:      true,
		PostgresDBHostFlag:       true,
		BatchIntervalFlag:        true,
		MaxBatchIntervalFlag:     true,
		RollupIntervalFlag:       true,
		L1BlockTimeFlag:          true,
		IsInboundP2PDisabledFlag: true,
		LevelDBPathFlag:          true,
	},
	Node: {
		OverrideFlag: true,
		ConfigFlag:   true,
		// NodeInputDetails
		NodeNameFlag:          true,
		HostIDFlag:            true,
		PrivateKeyFlag:        true,
		L1WebsocketURLFlag:    true,
		P2pPublicAddressFlag:  true,
		ClientRPCPortHTTPFlag: true,
		ClientRPCPortWSFlag:   true,
		// NodeInputSettings
		NodeTypeFlag:              true,
		IsSGXEnabledFlag:          true,
		PccsAddrFlag:              true,
		DebugNamespaceEnabledFlag: true,
		LogLevelFlag:              true,
		ProfilerEnabledFlag:       true,
		UseInMemoryDBFlag:         true,
		PostgresDBHostFlag:        true,
		// NodeInputImages
		HostImage:       true,
		EnclaveImage:    true,
		EdgelessDBImage: true,
	},
}

// EnclaveRestrictedFlags are the flags that the enclave can receive ONLY over (a) the Ego signed enclave.json (./go/enclave/main/enclave.json)
// or (b) if passed via EDG_<flag> as environment variable (see https://docs.edgeless.systems/ego/reference/config#environment-variables).
// In the case of running enclave as standalone process (./go/enclave/main/main w/o ego sign enclave.json stage) these flags will be checked
// to be set via EDG_<flag> env vars.
var EnclaveRestrictedFlags = map[string]string{
	L1ChainIDFlag:             "int64",
	TenChainIDFlag:            "int64",
	TenGenesisFlag:            "string",
	UseInMemoryDBFlag:         "bool",
	ProfilerEnabledFlag:       "bool",
	DebugNamespaceEnabledFlag: "bool",
}

// FlagUsageMap is a full indexing of available flags across all service configurations
func FlagUsageMap() map[string]string {
	return map[string]string{
		OverrideFlag:                  "Additive config file to apply on top of default or -config",
		ConfigFlag:                    "The path to the host's config file. Overrides all other flags",
		NodeNameFlag:                  "Common name for containers and reference",
		IsSGXEnabledFlag:              "Use SGX or simulation",
		PccsAddrFlag:                  "SGX attestation address",
		HostImage:                     "Docker image for host service",
		EnclaveImage:                  "Docker image for enclave service",
		EdgelessDBImage:               "Docker image for edgeless DB (enclave persistence)",
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
		GasPaymentAddress:             "Account used for gas payments of L1 transactions",
		GasBatchExecutionLimit:        "Max gas that can be executed in a single batch",
		GenesisJSONFlag:               "// When validating incoming blocks, the genesis config for the L1 chain",
		TenGenesisFlag:                "The json string with the TEN genesis",
		L1ChainIDFlag:                 "An integer representing the unique chain id of the Ethereum chain used as an L1 (default 1337)",
		TenChainIDFlag:                "An integer representing the unique chain id of the TEN chain (default 443)",
		UseInMemoryDBFlag:             "Whether the enclave will use an in-memory DB rather than persist data",
		ProfilerEnabledFlag:           "Runs a profiler instance (Defaults to false)",
		DebugNamespaceEnabledFlag:     "Whether the debug namespace is enabled",
		GasLocalExecutionCapFlag:      "Max gas usage when executing local transactions",
		IsGenesisFlag:                 "Whether the host is the first host to join the network",
		ClientRPCPortHTTPFlag:         "The port on which to listen for client application RPC requests over HTTP",
		ClientRPCPortWSFlag:           "The port on which to listen for client application RPC requests over websockets",
		ClientRPCHostFlag:             "The host on which to handle client application RPC requests",
		EnclaveRPCAddressesFlag:       "The comma-separated addresses to use to connect to the Ten enclaves",
		P2pBindAddressFlag:            "The address where the p2p server is bound to. Defaults to 0.0.0.0:10000",
		P2pPublicAddressFlag:          "The P2P address where the other servers should connect to. Defaults to 127.0.0.1:10000",
		L1WebsocketURLFlag:            "The websocket RPC address the host can use for L1 requests",
		EnclaveRPCTimeoutFlag:         "The timeout for host <-> enclave RPC communication",
		L1RPCTimeoutFlag:              "The timeout for connecting to, and communicating with, the Ethereum client",
		P2pConnectionTimeoutFlag:      "The timeout for host <-> host P2P messaging",
		PrivateKeyFlag:                "The private key for the L1 host account",
		L1StartHashFlag:               "The L1 block hash where the management contract was deployed",
		MetricsEnabledFlag:            "Whether the metrics are enabled (Defaults to true)",
		MetricsHTTPPortFlag:           "The port on which the metrics are served (Defaults to 0.0.0.0:14000)",
		PostgresDBHostFlag:            "The host for the Postgres DB instance",
		BatchIntervalFlag:             "Duration between each batch. Can be put down as 1.0s",
		MaxBatchIntervalFlag:          "Max interval between each batch, if greater than batchInterval then some empty batches will be skipped. Can be put down as 1.0s",
		RollupIntervalFlag:            "Duration between each rollup. Can be put down as 1.0s",
		L1BlockTimeFlag:               "Time of 1l Blocks",
		IsInboundP2PDisabledFlag:      "Whether inbound p2p is enabled",
		LevelDBPathFlag:               "LevelDBPath",
	}
}
