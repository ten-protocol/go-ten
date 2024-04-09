package config

import "github.com/ten-protocol/go-ten/go/common/flag"

// Flag names.
const (
	configFlag                   = "config"
	nodeIDFlag                   = "id"
	isGenesisFlag                = "isGenesis"
	nodeTypeFlag                 = "nodeType"
	clientRPCPortHTTPFlag        = "clientRPCPortHttp"
	clientRPCPortWSFlag          = "clientRPCPortWs"
	clientRPCHostFlag            = "clientRPCHost"
	enclaveRPCAddressFlag        = "enclaveRPCAddress"
	p2pBindAddressFlag           = "p2pBindAddress"
	p2pPublicAddressFlag         = "p2pPublicAddress"
	l1WebsocketURLFlag           = "l1WSURL"
	enclaveRPCTimeoutSecsFlag    = "enclaveRPCTimeoutSecs"
	l1RPCTimeoutSecsFlag         = "l1RPCTimeoutSecs"
	p2pConnectionTimeoutSecsFlag = "p2pConnectionTimeoutSecs"
	managementContractAddrFlag   = "managementContractAddress"
	messageBusContractAddrFlag   = "messageBusContractAddress"
	logLevelFlag                 = "logLevel"
	logPathFlag                  = "logPath"
	privateKeyFlag               = "privateKey"
	l1ChainIDFlag                = "l1ChainID"
	tenChainIDFlag               = "tenChainID"
	profilerEnabledFlag          = "profilerEnabled"
	l1StartHashFlag              = "l1Start"
	sequencerIDFlag              = "sequencerID"
	metricsEnabledFlag           = "metricsEnabled"
	metricsHTTPPortFlag          = "metricsHTTPPort"
	useInMemoryDBFlag            = "useInMemoryDB"
	levelDBPathFlag              = "levelDBPath"
	debugNamespaceEnabledFlag    = "debugNamespaceEnabled"
	batchIntervalFlag            = "batchInterval"
	maxBatchIntervalFlag         = "maxBatchInterval"
	rollupIntervalFlag           = "rollupInterval"
	isInboundP2PDisabledFlag     = "isInboundP2PDisabled"
	maxRollupSizeFlagFlag        = "maxRollupSize"
)

// HostFlags are the flags that the host can receive
var HostFlags map[string]*flag.TenFlag

func init() {
	hostFlagDefaults := DefaultHostConfig() // for setting flag default values
	flagUsageMap := getFlagUsageMap()

	HostFlags = map[string]*flag.TenFlag{
		nodeIDFlag:                   flag.NewStringFlag(nodeIDFlag, hostFlagDefaults.ID.String(), flagUsageMap[nodeIDFlag]),
		isGenesisFlag:                flag.NewBoolFlag(isGenesisFlag, hostFlagDefaults.IsGenesis, flagUsageMap[isGenesisFlag]),
		nodeTypeFlag:                 flag.NewStringFlag(nodeTypeFlag, hostFlagDefaults.NodeType.String(), flagUsageMap[nodeTypeFlag]),
		clientRPCPortHTTPFlag:        flag.NewUint64Flag(clientRPCPortHTTPFlag, hostFlagDefaults.ClientRPCPortHTTP, flagUsageMap[clientRPCPortHTTPFlag]),
		clientRPCPortWSFlag:          flag.NewUint64Flag(clientRPCPortWSFlag, hostFlagDefaults.ClientRPCPortWS, flagUsageMap[clientRPCPortWSFlag]),
		clientRPCHostFlag:            flag.NewStringFlag(clientRPCHostFlag, hostFlagDefaults.ClientRPCHost, flagUsageMap[clientRPCHostFlag]),
		enclaveRPCAddressFlag:        flag.NewStringFlag(enclaveRPCAddressFlag, hostFlagDefaults.EnclaveRPCAddress, flagUsageMap[enclaveRPCAddressFlag]),
		p2pBindAddressFlag:           flag.NewStringFlag(p2pBindAddressFlag, hostFlagDefaults.P2PBindAddress, flagUsageMap[p2pBindAddressFlag]),
		p2pPublicAddressFlag:         flag.NewStringFlag(p2pPublicAddressFlag, hostFlagDefaults.P2PPublicAddress, flagUsageMap[p2pPublicAddressFlag]),
		l1WebsocketURLFlag:           flag.NewStringFlag(l1WebsocketURLFlag, hostFlagDefaults.L1WebsocketURL, flagUsageMap[l1WebsocketURLFlag]),
		enclaveRPCTimeoutSecsFlag:    flag.NewDurationFlag(enclaveRPCTimeoutSecsFlag, hostFlagDefaults.EnclaveRPCTimeout, flagUsageMap[enclaveRPCTimeoutSecsFlag]),
		l1RPCTimeoutSecsFlag:         flag.NewDurationFlag(l1RPCTimeoutSecsFlag, hostFlagDefaults.L1RPCTimeout, flagUsageMap[l1RPCTimeoutSecsFlag]),
		p2pConnectionTimeoutSecsFlag: flag.NewDurationFlag(p2pConnectionTimeoutSecsFlag, hostFlagDefaults.P2PConnectionTimeout, flagUsageMap[p2pConnectionTimeoutSecsFlag]),
		managementContractAddrFlag:   flag.NewStringFlag(managementContractAddrFlag, hostFlagDefaults.ManagementContractAddress.String(), flagUsageMap[managementContractAddrFlag]),
		messageBusContractAddrFlag:   flag.NewStringFlag(messageBusContractAddrFlag, hostFlagDefaults.MessageBusAddress.String(), flagUsageMap[messageBusContractAddrFlag]),
		logLevelFlag:                 flag.NewIntFlag(logLevelFlag, hostFlagDefaults.LogLevel, flagUsageMap[logPathFlag]),
		logPathFlag:                  flag.NewStringFlag(logPathFlag, hostFlagDefaults.LogPath, flagUsageMap[logPathFlag]),
		privateKeyFlag:               flag.NewStringFlag(privateKeyFlag, hostFlagDefaults.PrivateKeyString, flagUsageMap[privateKeyFlag]),
		l1ChainIDFlag:                flag.NewInt64Flag(l1ChainIDFlag, hostFlagDefaults.L1ChainID, flagUsageMap[l1ChainIDFlag]),
		tenChainIDFlag:               flag.NewInt64Flag(tenChainIDFlag, hostFlagDefaults.TenChainID, flagUsageMap[tenChainIDFlag]),
		profilerEnabledFlag:          flag.NewBoolFlag(profilerEnabledFlag, hostFlagDefaults.ProfilerEnabled, flagUsageMap[profilerEnabledFlag]),
		l1StartHashFlag:              flag.NewStringFlag(l1StartHashFlag, hostFlagDefaults.L1StartHash.Hex(), flagUsageMap[l1StartHashFlag]),
		sequencerIDFlag:              flag.NewStringFlag(sequencerIDFlag, hostFlagDefaults.SequencerID.String(), flagUsageMap[sequencerIDFlag]),
		metricsEnabledFlag:           flag.NewBoolFlag(metricsEnabledFlag, hostFlagDefaults.MetricsEnabled, flagUsageMap[metricsEnabledFlag]),
		metricsHTTPPortFlag:          flag.NewUIntFlag(metricsHTTPPortFlag, hostFlagDefaults.MetricsHTTPPort, flagUsageMap[metricsHTTPPortFlag]),
		useInMemoryDBFlag:            flag.NewBoolFlag(useInMemoryDBFlag, hostFlagDefaults.UseInMemoryDB, flagUsageMap[useInMemoryDBFlag]),
		levelDBPathFlag:              flag.NewStringFlag(levelDBPathFlag, hostFlagDefaults.LevelDBPath, flagUsageMap[levelDBPathFlag]),
		debugNamespaceEnabledFlag:    flag.NewBoolFlag(debugNamespaceEnabledFlag, hostFlagDefaults.DebugNamespaceEnabled, flagUsageMap[debugNamespaceEnabledFlag]),
		batchIntervalFlag:            flag.NewDurationFlag(batchIntervalFlag, hostFlagDefaults.BatchInterval, flagUsageMap[batchIntervalFlag]),
		maxBatchIntervalFlag:         flag.NewDurationFlag(maxBatchIntervalFlag, hostFlagDefaults.MaxBatchInterval, flagUsageMap[maxBatchIntervalFlag]),
		rollupIntervalFlag:           flag.NewDurationFlag(rollupIntervalFlag, hostFlagDefaults.RollupInterval, flagUsageMap[rollupIntervalFlag]),
		isInboundP2PDisabledFlag:     flag.NewBoolFlag(isInboundP2PDisabledFlag, hostFlagDefaults.IsInboundP2PDisabled, flagUsageMap[isInboundP2PDisabledFlag]),
		maxRollupSizeFlagFlag:        flag.NewUint64Flag(maxRollupSizeFlagFlag, hostFlagDefaults.MaxRollupSize, flagUsageMap[maxRollupSizeFlagFlag]),
	}
}

// Returns a map of the flag usages.
// While we could just use constants instead of a map, this approach allows us to test that all the expected flags are defined.
func getFlagUsageMap() map[string]string {
	return map[string]string{
		configFlag:                   "The path to the host's config file. Overrides all other flags",
		nodeIDFlag:                   "The 20 bytes of the host's address",
		isGenesisFlag:                "Whether the host is the first host to join the network",
		nodeTypeFlag:                 "The node's type (e.g. aggregator, validator)",
		clientRPCPortHTTPFlag:        "The port on which to listen for client application RPC requests over HTTP",
		clientRPCPortWSFlag:          "The port on which to listen for client application RPC requests over websockets",
		clientRPCHostFlag:            "The host on which to handle client application RPC requests",
		enclaveRPCAddressFlag:        "The address to use to connect to the TEN enclave service",
		p2pBindAddressFlag:           "The address where the p2p server is bound to. Defaults to 0.0.0.0:10000",
		p2pPublicAddressFlag:         "The P2P address where the other servers should connect to. Defaults to 127.0.0.1:10000",
		l1WebsocketURLFlag:           "The websocket RPC address the host can use for L1 requests",
		enclaveRPCTimeoutSecsFlag:    "The timeout for host <-> enclave RPC communication",
		l1RPCTimeoutSecsFlag:         "The timeout for connecting to, and communicating with, the Ethereum client",
		p2pConnectionTimeoutSecsFlag: "The timeout for host <-> host P2P messaging",
		managementContractAddrFlag:   "The management contract address on the L1",
		messageBusContractAddrFlag:   "The message bus contract address on the L1",
		logLevelFlag:                 "The verbosity level of logs. (Defaults to Info)",
		logPathFlag:                  "The path to use for the host's log file",
		privateKeyFlag:               "The private key for the L1 host account",
		l1ChainIDFlag:                "An integer representing the unique chain id of the Ethereum chain used as an L1 (default 1337)",
		tenChainIDFlag:               "An integer representing the unique chain id of the TEN chain (default 443)",
		profilerEnabledFlag:          "Runs a profiler instance (Defaults to false)",
		l1StartHashFlag:              "The L1 block hash where the management contract was deployed",
		sequencerIDFlag:              "The ID of the sequencer",
		metricsEnabledFlag:           "Whether the metrics are enabled (Defaults to true)",
		metricsHTTPPortFlag:          "The port on which the metrics are served (Defaults to 0.0.0.0:14000)",
		useInMemoryDBFlag:            "Whether the host will use an in-memory DB rather than persist data",
		levelDBPathFlag:              "Filepath for the levelDB persistence dir (can be empty if a throwaway file in /tmp/ is acceptable or if using InMemory DB)",
		debugNamespaceEnabledFlag:    "Whether the debug names is enabled",
		batchIntervalFlag:            "Duration between each batch. Can be put down as 1.0s",
		maxBatchIntervalFlag:         "Max interval between each batch, if greater than batchInterval then some empty batches will be skipped. Can be put down as 1.0s",
		rollupIntervalFlag:           "Duration between each rollup. Can be put down as 1.0s",
		isInboundP2PDisabledFlag:     "Whether inbound p2p is enabled",
		maxRollupSizeFlagFlag:        "Max size of a rollup",
	}
}
