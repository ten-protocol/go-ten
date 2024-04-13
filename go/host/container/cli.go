package container

import (
	"flag"
	"fmt"
	"github.com/ten-protocol/go-ten/go/config"
	"strings"
)

// ParseConfig returns a config.HostInputConfig based on either the file identified by the `config` flag, or the flags with
// specific defaults (if the `config` flag isn't specified).
func ParseConfig() (*config.HostConfig, error) {
	inputCfg, err := config.LoadDefaultInputConfig(config.Host)
	if err != nil {
		panic(fmt.Errorf("issues loading default and override config from file. Cause: %w", err))
	}
	cfg := inputCfg.(*config.HostInputConfig) // assert
	flagUsageMap := getFlagUsageMap()

	isGenesis := flag.Bool(isGenesisFlag, config.GetEnvBool(isGenesisFlag, cfg.IsGenesis), flagUsageMap[isGenesisFlag])
	nodeType := flag.String(nodeTypeFlag, config.GetEnvString(nodeTypeFlag, cfg.NodeType), flagUsageMap[nodeTypeFlag])
	clientRPCPortHTTP := flag.Uint64(clientRPCPortHTTPFlag, config.GetEnvUint64(clientRPCPortHTTPFlag, cfg.ClientRPCPortHTTP), flagUsageMap[clientRPCPortHTTPFlag])
	clientRPCPortWS := flag.Uint64(clientRPCPortWSFlag, config.GetEnvUint64(clientRPCPortWSFlag, cfg.ClientRPCPortWS), flagUsageMap[clientRPCPortWSFlag])
	clientRPCHost := flag.String(clientRPCHostFlag, config.GetEnvString(clientRPCHostFlag, cfg.ClientRPCHost), flagUsageMap[clientRPCHostFlag])
	enclaveRPCAddressesStr := flag.String(enclaveRPCAddressesFlag, config.GetEnvString(enclaveRPCAddressesFlag, strings.Join(cfg.EnclaveRPCAddresses, ",")), flagUsageMap[enclaveRPCAddressesFlag])
	p2pBindAddress := flag.String(p2pBindAddressFlag, config.GetEnvString(p2pBindAddressFlag, cfg.P2PBindAddress), flagUsageMap[p2pBindAddressFlag])
	p2pPublicAddress := flag.String(p2pPublicAddressFlag, config.GetEnvString(p2pPublicAddressFlag, cfg.P2PPublicAddress), flagUsageMap[p2pPublicAddressFlag])
	l1WebsocketURL := flag.String(l1WebsocketURLFlag, config.GetEnvString(l1WebsocketURLFlag, cfg.L1WebsocketURL), flagUsageMap[l1WebsocketURLFlag])
	enclaveRPCTimeout := flag.Int(enclaveRPCTimeoutSecsFlag, config.GetEnvInt(enclaveRPCTimeoutSecsFlag, cfg.EnclaveRPCTimeout), flagUsageMap[enclaveRPCTimeoutSecsFlag])
	l1RPCTimeout := flag.Int(l1RPCTimeoutSecsFlag, config.GetEnvInt(l1RPCTimeoutSecsFlag, cfg.L1RPCTimeout), flagUsageMap[l1RPCTimeoutSecsFlag])
	p2pConnectionTimeout := flag.Int(p2pConnectionTimeoutSecsFlag, config.GetEnvInt(p2pConnectionTimeoutSecsFlag, cfg.P2PConnectionTimeout), flagUsageMap[p2pConnectionTimeoutSecsFlag])
	managementContractAddress := flag.String(managementContractAddrFlag, config.GetEnvString(managementContractAddrFlag, cfg.ManagementContractAddress), flagUsageMap[managementContractAddrFlag])
	messageBusContractAddress := flag.String(messageBusContractAddrFlag, config.GetEnvString(messageBusContractAddrFlag, cfg.MessageBusAddress), flagUsageMap[messageBusContractAddrFlag])
	logLevel := flag.Int(logLevelFlag, config.GetEnvInt(logLevelFlag, cfg.LogLevel), flagUsageMap[logLevelFlag])
	logPath := flag.String(logPathFlag, config.GetEnvString(logPathFlag, cfg.LogPath), flagUsageMap[logPathFlag])
	l1ChainID := flag.Int64(l1ChainIDFlag, config.GetEnvInt64(l1ChainIDFlag, cfg.L1ChainID), flagUsageMap[l1ChainIDFlag])
	tenChainID := flag.Int64(tenChainIDFlag, config.GetEnvInt64(tenChainIDFlag, cfg.TenChainID), flagUsageMap[tenChainIDFlag])
	privateKey := flag.String(privateKeyFlag, config.GetEnvString(privateKeyFlag, cfg.PrivateKey), flagUsageMap[privateKeyFlag])
	profilerEnabled := flag.Bool(profilerEnabledFlag, config.GetEnvBool(profilerEnabledFlag, cfg.ProfilerEnabled), flagUsageMap[profilerEnabledFlag])
	l1StartHash := flag.String(l1StartHashFlag, config.GetEnvString(l1StartHashFlag, cfg.L1StartHash), flagUsageMap[l1StartHashFlag])
	sequencerID := flag.String(sequencerIDFlag, config.GetEnvString(sequencerIDFlag, cfg.SequencerID), flagUsageMap[sequencerIDFlag])
	metricsEnabled := flag.Bool(metricsEnabledFlag, config.GetEnvBool(metricsEnabledFlag, cfg.MetricsEnabled), flagUsageMap[metricsEnabledFlag])
	metricsHTTPPort := flag.Uint(metricsHTTPPortFlag, config.GetEnvUint(metricsHTTPPortFlag, cfg.MetricsHTTPPort), flagUsageMap[metricsHTTPPortFlag])
	useInMemoryDB := flag.Bool(useInMemoryDBFlag, config.GetEnvBool(useInMemoryDBFlag, cfg.UseInMemoryDB), flagUsageMap[useInMemoryDBFlag])
	postgresDBHost := flag.String(postgresDBHostFlag, config.GetEnvString(postgresDBHostFlag, cfg.PostgresDBHost), flagUsageMap[postgresDBHostFlag])
	debugNamespaceEnabled := flag.Bool(debugNamespaceEnabledFlag, config.GetEnvBool(debugNamespaceEnabledFlag, cfg.DebugNamespaceEnabled), flagUsageMap[debugNamespaceEnabledFlag])
	batchInterval := flag.Int(batchIntervalFlag, config.GetEnvInt(batchIntervalFlag, cfg.BatchInterval), flagUsageMap[batchIntervalFlag])
	maxBatchInterval := flag.Int(maxBatchIntervalFlag, config.GetEnvInt(maxBatchIntervalFlag, cfg.MaxBatchInterval), flagUsageMap[maxBatchIntervalFlag])
	rollupInterval := flag.Int(rollupIntervalFlag, config.GetEnvInt(rollupIntervalFlag, cfg.RollupInterval), flagUsageMap[rollupIntervalFlag])
	isInboundP2PDisabled := flag.Bool(isInboundP2PDisabledFlag, config.GetEnvBool(isInboundP2PDisabledFlag, cfg.IsInboundP2PDisabled), flagUsageMap[isInboundP2PDisabledFlag])
	maxRollupSize := flag.Uint64(maxRollupSizeFlag, config.GetEnvUint64(maxRollupSizeFlag, cfg.MaxRollupSize), flagUsageMap[maxRollupSizeFlag])
	flag.Parse()

	cfg.IsGenesis = *isGenesis
	cfg.NodeType = *nodeType
	cfg.ClientRPCPortHTTP = *clientRPCPortHTTP
	cfg.ClientRPCPortWS = *clientRPCPortWS
	cfg.ClientRPCHost = *clientRPCHost
	cfg.EnclaveRPCAddresses = strings.Split(*enclaveRPCAddressesStr, ",")
	cfg.P2PBindAddress = *p2pBindAddress
	cfg.P2PPublicAddress = *p2pPublicAddress
	cfg.L1WebsocketURL = *l1WebsocketURL
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
	cfg.MetricsHTTPPort = *metricsHTTPPort
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
