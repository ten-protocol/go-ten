package container

import (
	"flag"
	"fmt"
	"github.com/ten-protocol/go-ten/go/config"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"strings"
)

// ParseConfig returns a config.HostInputConfig based on either the file identified by the `config` flag, or the flags with
// specific defaults (if the `config` flag isn't specified).
func ParseConfig() (*config.HostConfig, error) {
	cfg, err := loadDefaultHostInputConfig()
	if err != nil {
		panic(fmt.Errorf("issues loading default and override config from file. Cause: %w", err))
	}
	flagUsageMap := getFlagUsageMap()

	isGenesis := flag.Bool(isGenesisFlag, cfg.IsGenesis, flagUsageMap[isGenesisFlag])
	nodeType := flag.String(nodeTypeFlag, cfg.NodeType, flagUsageMap[nodeTypeFlag])
	clientRPCPortHTTP := flag.Uint64(clientRPCPortHTTPFlag, cfg.ClientRPCPortHTTP, flagUsageMap[clientRPCPortHTTPFlag])
	clientRPCPortWS := flag.Uint64(clientRPCPortWSFlag, cfg.ClientRPCPortWS, flagUsageMap[clientRPCPortWSFlag])
	clientRPCHost := flag.String(clientRPCHostFlag, cfg.ClientRPCHost, flagUsageMap[clientRPCHostFlag])
	enclaveRPCAddressesStr := flag.String(enclaveRPCAddressesFlag, strings.Join(cfg.EnclaveRPCAddresses, ","), flagUsageMap[enclaveRPCAddressesFlag])
	p2pBindAddress := flag.String(p2pBindAddressFlag, cfg.P2PBindAddress, flagUsageMap[p2pBindAddressFlag])
	p2pPublicAddress := flag.String(p2pPublicAddressFlag, cfg.P2PPublicAddress, flagUsageMap[p2pPublicAddressFlag])
	l1WSURL := flag.String(l1WebsocketURLFlag, cfg.L1WebsocketURL, flagUsageMap[l1WebsocketURLFlag])
	enclaveRPCTimeout := flag.Int(enclaveRPCTimeoutSecsFlag, cfg.EnclaveRPCTimeout, flagUsageMap[enclaveRPCTimeoutSecsFlag])
	l1RPCTimeout := flag.Int(l1RPCTimeoutSecsFlag, cfg.L1RPCTimeout, flagUsageMap[l1RPCTimeoutSecsFlag])
	p2pConnectionTimeout := flag.Int(p2pConnectionTimeoutSecsFlag, cfg.P2PConnectionTimeout, flagUsageMap[p2pConnectionTimeoutSecsFlag])
	managementContractAddress := flag.String(managementContractAddrFlag, cfg.ManagementContractAddress, flagUsageMap[managementContractAddrFlag])
	messageBusContractAddress := flag.String(messageBusContractAddrFlag, cfg.MessageBusAddress, flagUsageMap[messageBusContractAddrFlag])
	logLevel := flag.Int(logLevelFlag, cfg.LogLevel, flagUsageMap[logLevelFlag])
	logPath := flag.String(logPathFlag, cfg.LogPath, flagUsageMap[logPathFlag])
	l1ChainID := flag.Int64(l1ChainIDFlag, cfg.L1ChainID, flagUsageMap[l1ChainIDFlag])
	tenChainID := flag.Int64(tenChainIDFlag, cfg.TenChainID, flagUsageMap[tenChainIDFlag])
	privateKey := flag.String(privateKeyFlag, cfg.PrivateKey, flagUsageMap[privateKeyFlag])
	profilerEnabled := flag.Bool(profilerEnabledFlag, cfg.ProfilerEnabled, flagUsageMap[profilerEnabledFlag])
	l1StartHash := flag.String(l1StartHashFlag, cfg.L1StartHash, flagUsageMap[l1StartHashFlag])
	sequencerID := flag.String(sequencerIDFlag, cfg.SequencerID, flagUsageMap[sequencerIDFlag])
	metricsEnabled := flag.Bool(metricsEnabledFlag, cfg.MetricsEnabled, flagUsageMap[metricsEnabledFlag])
	metricsHTPPPort := flag.Uint(metricsHTTPPortFlag, cfg.MetricsHTTPPort, flagUsageMap[metricsHTTPPortFlag])
	useInMemoryDB := flag.Bool(useInMemoryDBFlag, cfg.UseInMemoryDB, flagUsageMap[useInMemoryDBFlag])
	postgresDBHost := flag.String(postgresDBHostFlag, cfg.PostgresDBHost, flagUsageMap[postgresDBHostFlag])
	debugNamespaceEnabled := flag.Bool(debugNamespaceEnabledFlag, cfg.DebugNamespaceEnabled, flagUsageMap[debugNamespaceEnabledFlag])
	batchInterval := flag.Int(batchIntervalFlag, cfg.BatchInterval, flagUsageMap[batchIntervalFlag])
	maxBatchInterval := flag.Int(maxBatchIntervalFlag, cfg.MaxBatchInterval, flagUsageMap[maxBatchIntervalFlag])
	rollupInterval := flag.Int(rollupIntervalFlag, cfg.RollupInterval, flagUsageMap[rollupIntervalFlag])
	isInboundP2PDisabled := flag.Bool(isInboundP2PDisabledFlag, cfg.IsInboundP2PDisabled, flagUsageMap[isInboundP2PDisabledFlag])
	maxRollupSize := flag.Uint64(maxRollupSizeFlag, cfg.MaxRollupSize, flagUsageMap[maxRollupSizeFlag])

	flag.Parse()

	cfg.IsGenesis = *isGenesis
	cfg.NodeType = *nodeType
	cfg.ClientRPCPortHTTP = *clientRPCPortHTTP
	cfg.ClientRPCPortWS = *clientRPCPortWS
	cfg.ClientRPCHost = *clientRPCHost
	cfg.EnclaveRPCAddresses = strings.Split(*enclaveRPCAddressesStr, ",")
	cfg.P2PBindAddress = *p2pBindAddress
	cfg.P2PPublicAddress = *p2pPublicAddress
	cfg.L1WebsocketURL = *l1WSURL
	cfg.EnclaveRPCTimeout = *enclaveRPCTimeout
	cfg.L1RPCTimeout = *l1RPCTimeout
	cfg.P2PConnectionTimeout = *p2pConnectionTimeout
	cfg.ManagementContractAddress = *managementContractAddress
	cfg.MessageBusAddress = *messageBusContractAddress
	cfg.PrivateKey = *privateKey
	cfg.LogLevel = *logLevel
	cfg.LogPath = *logPath
	cfg.L1ChainID = *l1ChainID
	cfg.TenChainID = *tenChainID
	cfg.ProfilerEnabled = *profilerEnabled
	cfg.L1StartHash = *l1StartHash
	cfg.SequencerID = *sequencerID
	cfg.MetricsEnabled = *metricsEnabled
	cfg.MetricsHTTPPort = *metricsHTPPPort
	cfg.UseInMemoryDB = *useInMemoryDB
	cfg.PostgresDBHost = *postgresDBHost
	cfg.DebugNamespaceEnabled = *debugNamespaceEnabled
	cfg.BatchInterval = *batchInterval
	cfg.MaxBatchInterval = *maxBatchInterval
	cfg.RollupInterval = *rollupInterval
	cfg.IsInboundP2PDisabled = *isInboundP2PDisabled
	cfg.MaxRollupSize = *maxRollupSize

	hostConfig, err := cfg.ToHostConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to convert HostInputConfig to HostConfig")
	}
	return hostConfig, nil
}

// loadDefaultHostInputConfig parses optional or default configuration file and returns struct.
func loadDefaultHostInputConfig() (*config.HostInputConfig, error) {
	flagUsageMap := getFlagUsageMap()
	configPath := flag.String("config", "./go/config/templates/default_host_config.yaml", flagUsageMap["configFlag"])
	overridePath := flag.String("override", "", flagUsageMap["overrideFlag"])

	// Parse only once capturing all necessary flags
	flag.Parse()

	var err error
	conf, err := loadHostConfigFromFile(*configPath)
	if err != nil {
		panic(err)
	}

	// Apply overrides if the override path is provided
	if *overridePath != "" {
		overridesConf, err := loadHostConfigFromFile(*overridePath)
		if err != nil {
			panic(err)
		}

		config.ApplyOverrides(conf, overridesConf)
	}

	return conf, nil
}

// loadHostConfigFromFile reads configuration from a file and environment variables
func loadHostConfigFromFile(configPath string) (*config.HostInputConfig, error) {
	defaultConfig := &config.HostInputConfig{}
	// Read YAML configuration
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, defaultConfig)
	if err != nil {
		return nil, err
	}

	return defaultConfig, nil
}
