package config

// SHARED Enclave and Host flags
const (
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
)

// Host Flag names.
const (
	isGenesisFlag            = "isGenesis"
	clientRPCPortHTTPFlag    = "clientRPCPortHTTP"
	clientRPCPortWSFlag      = "clientRPCPortWS"
	clientRPCHostFlag        = "clientRPCHost"
	enclaveRPCAddressesFlag  = "enclaveRPCAddresses"
	p2pBindAddressFlag       = "p2pBindAddress"
	p2pPublicAddressFlag     = "p2pPublicAddress"
	l1WebsocketURLFlag       = "l1WebsocketURL"
	enclaveRPCTimeoutFlag    = "enclaveRPCTimeout"
	l1RPCTimeoutFlag         = "l1RPCTimeout"
	p2pConnectionTimeoutFlag = "p2pConnectionTimeout"
	privateKeyFlag           = "privateKey"
	l1StartHashFlag          = "l1StartHash"
	metricsEnabledFlag       = "metricsEnabled"
	metricsHTTPPortFlag      = "metricsHTTPPort"
	postgresDBHostFlag       = "postgresDBHost"
	batchIntervalFlag        = "batchInterval"
	maxBatchIntervalFlag     = "maxBatchInterval"
	rollupIntervalFlag       = "rollupInterval"
	l1BlockTimeFlag          = "l1BlockTime"
	isInboundP2PDisabledFlag = "isInboundP2PDisabled"
	LevelDBPathFlag          = "levelDBPath"
)

// Enclave Flag names.
const (
	HostIDFlag               = "hostID"
	HostAddressFlag          = "hostAddress"
	AddressFlag              = "address"
	WillAttestFlag           = "willAttest"
	ValidateL1BlocksFlag     = "validateL1Blocks"
	EdgelessDBHostFlag       = "edgelessDBHost"
	MinGasPriceFlag          = "minGasPrice"
	GenesisJSONFlag          = "genesisJSON"
	TenGenesisFlag           = "tenGenesis"
	MaxBatchSizeFlag         = "maxBatchSize"
	L2BaseFeeFlag            = "l2BaseFee"
	GasPaymentAddress        = "gasPaymentAddress"
	GasBatchExecutionLimit   = "gasBatchExecutionLimit"
	GasLocalExecutionCapFlag = "gasLocalExecutionCap"
)

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
		isGenesisFlag:                 "Whether the host is the first host to join the network",
		clientRPCPortHTTPFlag:         "The port on which to listen for client application RPC requests over HTTP",
		clientRPCPortWSFlag:           "The port on which to listen for client application RPC requests over websockets",
		clientRPCHostFlag:             "The host on which to handle client application RPC requests",
		enclaveRPCAddressesFlag:       "The comma-separated addresses to use to connect to the Ten enclaves",
		p2pBindAddressFlag:            "The address where the p2p server is bound to. Defaults to 0.0.0.0:10000",
		p2pPublicAddressFlag:          "The P2P address where the other servers should connect to. Defaults to 127.0.0.1:10000",
		l1WebsocketURLFlag:            "The websocket RPC address the host can use for L1 requests",
		enclaveRPCTimeoutFlag:         "The timeout for host <-> enclave RPC communication",
		l1RPCTimeoutFlag:              "The timeout for connecting to, and communicating with, the Ethereum client",
		p2pConnectionTimeoutFlag:      "The timeout for host <-> host P2P messaging",
		privateKeyFlag:                "The private key for the L1 host account",
		l1StartHashFlag:               "The L1 block hash where the management contract was deployed",
		metricsEnabledFlag:            "Whether the metrics are enabled (Defaults to true)",
		metricsHTTPPortFlag:           "The port on which the metrics are served (Defaults to 0.0.0.0:14000)",
		postgresDBHostFlag:            "The host for the Postgres DB instance",
		batchIntervalFlag:             "Duration between each batch. Can be put down as 1.0s",
		maxBatchIntervalFlag:          "Max interval between each batch, if greater than batchInterval then some empty batches will be skipped. Can be put down as 1.0s",
		rollupIntervalFlag:            "Duration between each rollup. Can be put down as 1.0s",
		l1BlockTimeFlag:               "Time of 1l Blocks",
		isInboundP2PDisabledFlag:      "Whether inbound p2p is enabled",
		LevelDBPathFlag:               "LevelDBPath",
	}
}
