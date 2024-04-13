package config

// Flag names.
const (
	isGenesisFlag                 = "isGenesis"
	nodeTypeFlag                  = "nodeType"
	clientRPCPortHTTPFlag         = "clientRPCPortHTTP"
	clientRPCPortWSFlag           = "clientRPCPortWS"
	clientRPCHostFlag             = "clientRPCHost"
	enclaveRPCAddressesFlag       = "enclaveRPCAddresses"
	p2pBindAddressFlag            = "p2pBindAddress"
	p2pPublicAddressFlag          = "p2pPublicAddress"
	l1WebsocketURLFlag            = "l1WebsocketURL"
	enclaveRPCTimeoutFlag         = "enclaveRPCTimeout"
	l1RPCTimeoutFlag              = "l1RPCTimeout"
	p2pConnectionTimeoutFlag      = "p2pConnectionTimeout"
	managementContractAddressFlag = "managementContractAddress"
	messageBusAddressFlag         = "messageBusAddress"
	logLevelFlag                  = "logLevel"
	logPathFlag                   = "logPath"
	privateKeyFlag                = "privateKey"
	l1ChainIDFlag                 = "l1ChainID"
	tenChainIDFlag                = "tenChainID"
	profilerEnabledFlag           = "profilerEnabled"
	l1StartHashFlag               = "l1StartHash"
	sequencerIDFlag               = "sequencerID"
	metricsEnabledFlag            = "metricsEnabled"
	metricsHTTPPortFlag           = "metricsHTTPPort"
	useInMemoryDBFlag             = "useInMemoryDB"
	postgresDBHostFlag            = "postgresDBHost"
	sqliteDBPathFlag              = "sqliteDBPath"
	debugNamespaceEnabledFlag     = "debugNamespaceEnabled"
	batchIntervalFlag             = "batchInterval"
	maxBatchIntervalFlag          = "maxBatchInterval"
	rollupIntervalFlag            = "rollupInterval"
	l1BlockTimeFlag               = "l1BlockTime"
	isInboundP2PDisabledFlag      = "isInboundP2PDisabled"
	maxRollupSizeFlag             = "maxRollupSize"
	LevelDBPathFlag               = "levelDBPath"
)

// UsageMap Returns a map of the flag usages.
// While we could just use constants instead of a map, this approach allows us to test that all the expected flags are defined.
func (*HostInputConfig) UsageMap() map[string]string {
	return map[string]string{
		isGenesisFlag:                 "Whether the host is the first host to join the network",
		nodeTypeFlag:                  "The node's type (e.g. aggregator, validator)",
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
		managementContractAddressFlag: "The management contract address on the L1",
		messageBusAddressFlag:         "The message bus contract address on the L1",
		logLevelFlag:                  "The verbosity level of logs. (Defaults to Info)",
		logPathFlag:                   "The path to use for the host's log file",
		privateKeyFlag:                "The private key for the L1 host account",
		l1ChainIDFlag:                 "An integer representing the unique chain id of the Ethereum chain used as an L1 (default 1337)",
		tenChainIDFlag:                "An integer representing the unique chain id of the Obscuro chain (default 443)",
		profilerEnabledFlag:           "Runs a profiler instance (Defaults to false)",
		l1StartHashFlag:               "The L1 block hash where the management contract was deployed",
		sequencerIDFlag:               "The ID of the sequencer",
		metricsEnabledFlag:            "Whether the metrics are enabled (Defaults to true)",
		metricsHTTPPortFlag:           "The port on which the metrics are served (Defaults to 0.0.0.0:14000)",
		useInMemoryDBFlag:             "Whether the host will use an in-memory DB rather than persist data",
		postgresDBHostFlag:            "The host for the Postgres DB instance",
		sqliteDBPathFlag:              "Filepath for the sqlite DB persistence file (can be empty if a throwaway file in /tmp/ is acceptable or if using InMemory DB",
		debugNamespaceEnabledFlag:     "Whether the debug names is enabled",
		batchIntervalFlag:             "Duration between each batch. Can be put down as 1.0s",
		maxBatchIntervalFlag:          "Max interval between each batch, if greater than batchInterval then some empty batches will be skipped. Can be put down as 1.0s",
		rollupIntervalFlag:            "Duration between each rollup. Can be put down as 1.0s",
		l1BlockTimeFlag:               "Time of 1l Blocks",
		isInboundP2PDisabledFlag:      "Whether inbound p2p is enabled",
		maxRollupSizeFlag:             "Max size of a rollup",
		LevelDBPathFlag:               "LevelDBPath",
	}
}
