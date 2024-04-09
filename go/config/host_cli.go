package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ten-protocol/go-ten/go/common"

	"github.com/naoina/toml"
	"gopkg.in/yaml.v3"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

// HostFileConfig is the structure that a host's .toml config is parsed into.
type HostFileConfig struct {
	IsGenesis                 bool
	NodeType                  string
	HasClientRPCHTTP          bool
	ClientRPCPortHTTP         uint
	HasClientRPCWebsockets    bool
	ClientRPCPortWS           uint
	ClientRPCHost             string
	EnclaveRPCAddress         string
	P2PBindAddress            string
	P2PPublicAddress          string
	L1WebsocketURL            string
	EnclaveRPCTimeout         int
	L1RPCTimeout              int
	P2PConnectionTimeout      int
	ManagementContractAddress string
	MessageBusAddress         string
	LogLevel                  int
	LogPath                   string
	ID                        string
	PrivateKeyString          string
	L1ChainID                 int64
	TenChainID                int64
	ProfilerEnabled           bool
	L1StartHash               string
	SequencerID               string
	MetricsEnabled            bool
	MetricsHTTPPort           uint
	UseInMemoryDB             bool
	LevelDBPath               string
	DebugNamespaceEnabled     bool
	BatchInterval             string
	MaxBatchInterval          string
	RollupInterval            string
	IsInboundP2PDisabled      bool
	L1BlockTime               int
	MaxRollupSize             int
}

// ParseHostConfig returns a config.HostConfig based on either the file identified by the `config` flag, or with
// specific defaults (if the `config` flag isn't specified). Flags may override values. Order of override, flag > (file || default).
func ParseHostConfig() (*HostConfig, error) {
	flagUsageMap := getFlagUsageMap()
	cfg := &HostConfig{}

	// sets the base template to either the file or default
	configPath := flag.String(configName, "", flagUsageMap[configName])
	if *configPath != "" {
		fileCfg, err := fileBasedConfig(*configPath)
		if err != nil {
			return nil, err
		}
		cfg = fileCfg
	} else {
		cfg = DefaultHostConfig()
	}

	// process flag overrides
	isGenesis := flag.Bool(isGenesisName, cfg.IsGenesis, flagUsageMap[isGenesisName])
	nodeTypeStr := flag.String(nodeTypeName, cfg.NodeType.String(), flagUsageMap[nodeTypeName])
	clientRPCPortHTTP := flag.Uint64(clientRPCPortHTTPName, cfg.ClientRPCPortHTTP, flagUsageMap[clientRPCPortHTTPName])
	clientRPCPortWS := flag.Uint64(clientRPCPortWSName, cfg.ClientRPCPortWS, flagUsageMap[clientRPCPortWSName])
	clientRPCHost := flag.String(clientRPCHostName, cfg.ClientRPCHost, flagUsageMap[clientRPCHostName])
	enclaveRPCAddress := flag.String(enclaveRPCAddressName, cfg.EnclaveRPCAddress, flagUsageMap[enclaveRPCAddressName])
	p2pBindAddress := flag.String(p2pBindAddressName, cfg.P2PBindAddress, flagUsageMap[p2pBindAddressName])
	p2pPublicAddress := flag.String(p2pPublicAddressName, cfg.P2PPublicAddress, flagUsageMap[p2pPublicAddressName])
	l1WSURL := flag.String(l1WebsocketURLName, cfg.L1WebsocketURL, flagUsageMap[l1WebsocketURLName])
	enclaveRPCTimeoutSecs := flag.Uint64(enclaveRPCTimeoutSecsName, uint64(cfg.EnclaveRPCTimeout.Seconds()), flagUsageMap[enclaveRPCTimeoutSecsName])
	l1RPCTimeoutSecs := flag.Uint64(l1RPCTimeoutSecsName, uint64(cfg.L1RPCTimeout.Seconds()), flagUsageMap[l1RPCTimeoutSecsName])
	p2pConnectionTimeoutSecs := flag.Uint64(p2pConnectionTimeoutSecsName, uint64(cfg.P2PConnectionTimeout.Seconds()), flagUsageMap[p2pConnectionTimeoutSecsName])
	managementContractAddress := flag.String(managementContractAddrName, cfg.ManagementContractAddress.Hex(), flagUsageMap[managementContractAddrName])
	messageBusContractAddress := flag.String(messageBusContractAddrName, cfg.MessageBusAddress.Hex(), flagUsageMap[messageBusContractAddrName])
	logLevel := flag.Int(logLevelName, cfg.LogLevel, flagUsageMap[logLevelName])
	logPath := flag.String(logPathName, cfg.LogPath, flagUsageMap[logPathName])
	l1ChainID := flag.Int64(l1ChainIDName, cfg.L1ChainID, flagUsageMap[l1ChainIDName])
	obscuroChainID := flag.Int64(obscuroChainIDName, cfg.ObscuroChainID, flagUsageMap[obscuroChainIDName])
	privateKeyStr := flag.String(privateKeyName, cfg.PrivateKeyString, flagUsageMap[privateKeyName])
	profilerEnabled := flag.Bool(profilerEnabledName, cfg.ProfilerEnabled, flagUsageMap[profilerEnabledName])
	l1StartHash := flag.String(l1StartHashName, cfg.L1StartHash.Hex(), flagUsageMap[l1StartHashName])
	sequencerID := flag.String(sequencerIDName, cfg.SequencerID.Hex(), flagUsageMap[sequencerIDName])
	metricsEnabled := flag.Bool(metricsEnabledName, cfg.MetricsEnabled, flagUsageMap[metricsEnabledName])
	metricsHTPPPort := flag.Uint(metricsHTTPPortName, cfg.MetricsHTTPPort, flagUsageMap[metricsHTTPPortName])
	useInMemoryDB := flag.Bool(useInMemoryDBName, cfg.UseInMemoryDB, flagUsageMap[useInMemoryDBName])
	levelDBPath := flag.String(levelDBPathName, cfg.LevelDBPath, flagUsageMap[levelDBPathName])
	debugNamespaceEnabled := flag.Bool(debugNamespaceEnabledName, cfg.DebugNamespaceEnabled, flagUsageMap[debugNamespaceEnabledName])
	batchInterval := flag.String(batchIntervalName, cfg.BatchInterval.String(), flagUsageMap[batchIntervalName])
	maxBatchInterval := flag.String(maxBatchIntervalName, cfg.MaxBatchInterval.String(), flagUsageMap[maxBatchIntervalName])
	rollupInterval := flag.String(rollupIntervalName, cfg.RollupInterval.String(), flagUsageMap[rollupIntervalName])
	isInboundP2PDisabled := flag.Bool(isInboundP2PDisabledName, cfg.IsInboundP2PDisabled, flagUsageMap[isInboundP2PDisabledName])
	maxRollupSize := flag.Uint64(maxRollupSizeFlagName, cfg.MaxRollupSize, flagUsageMap[maxRollupSizeFlagName])

	flag.Parse()

	nodeType, err := common.ToNodeType(*nodeTypeStr)
	if err != nil {
		return &HostConfig{}, fmt.Errorf("unrecognised node type '%s'", *nodeTypeStr)
	}

	cfg.IsGenesis = *isGenesis
	cfg.NodeType = nodeType
	cfg.HasClientRPCHTTP = true
	cfg.ClientRPCPortHTTP = *clientRPCPortHTTP
	cfg.HasClientRPCWebsockets = true
	cfg.ClientRPCPortWS = *clientRPCPortWS
	cfg.ClientRPCHost = *clientRPCHost
	cfg.EnclaveRPCAddress = *enclaveRPCAddress
	cfg.P2PBindAddress = *p2pBindAddress
	cfg.P2PPublicAddress = *p2pPublicAddress
	cfg.L1WebsocketURL = *l1WSURL
	cfg.EnclaveRPCTimeout = time.Duration(*enclaveRPCTimeoutSecs) * time.Second
	cfg.L1RPCTimeout = time.Duration(*l1RPCTimeoutSecs) * time.Second
	cfg.P2PConnectionTimeout = time.Duration(*p2pConnectionTimeoutSecs) * time.Second
	cfg.ManagementContractAddress = gethcommon.HexToAddress(*managementContractAddress)
	cfg.MessageBusAddress = gethcommon.HexToAddress(*messageBusContractAddress)
	cfg.PrivateKeyString = *privateKeyStr
	cfg.LogLevel = *logLevel
	cfg.LogPath = *logPath
	cfg.L1ChainID = *l1ChainID
	cfg.ObscuroChainID = *obscuroChainID
	cfg.ProfilerEnabled = *profilerEnabled
	cfg.L1StartHash = gethcommon.HexToHash(*l1StartHash)
	cfg.SequencerID = gethcommon.HexToAddress(*sequencerID)
	cfg.MetricsEnabled = *metricsEnabled
	cfg.MetricsHTTPPort = *metricsHTPPPort
	cfg.UseInMemoryDB = *useInMemoryDB
	cfg.LevelDBPath = *levelDBPath
	cfg.DebugNamespaceEnabled = *debugNamespaceEnabled
	cfg.BatchInterval, err = time.ParseDuration(*batchInterval)
	if err != nil {
		return nil, err
	}
	cfg.MaxBatchInterval, err = time.ParseDuration(*maxBatchInterval)
	if err != nil {
		return nil, err
	}
	cfg.RollupInterval, err = time.ParseDuration(*rollupInterval)
	if err != nil {
		return nil, err
	}
	cfg.IsInboundP2PDisabled = *isInboundP2PDisabled
	cfg.MaxRollupSize = *maxRollupSize

	return cfg, nil
}

// Parses the config from the .yaml or .toml file at configPath.
func fileBasedConfig(configPath string) (*HostConfig, error) {
	bytes, err := os.ReadFile(configPath)
	if err != nil {
		panic(fmt.Sprintf("could not read config file at %s. Cause: %s", configPath, err))
	}

	var fileConfig HostFileConfig
	err = yaml.Unmarshal(bytes, &fileConfig)
	if err != nil {
		// If yaml unmarshall failed, use toml
		err = toml.Unmarshal(bytes, &fileConfig)
		if err != nil {
			panic(fmt.Sprintf("could not read config file (.yaml or .toml) at %s. Cause: %s", configPath, err))
		}
	}

	nodeType, err := common.ToNodeType(fileConfig.NodeType)
	if err != nil {
		return &HostConfig{}, fmt.Errorf("unrecognised node type '%s'", fileConfig.NodeType)
	}

	batchInterval, maxBatchInterval, rollupInterval := 1*time.Second, 1*time.Second, 5*time.Second
	if interval, err := time.ParseDuration(fileConfig.BatchInterval); err == nil {
		batchInterval = interval
	}
	if interval, err := time.ParseDuration(fileConfig.RollupInterval); err == nil {
		rollupInterval = interval
	}
	if interval, err := time.ParseDuration(fileConfig.MaxBatchInterval); err == nil {
		maxBatchInterval = interval
	}

	return &HostConfig{
		IsGenesis:                 fileConfig.IsGenesis,
		NodeType:                  nodeType,
		HasClientRPCHTTP:          fileConfig.HasClientRPCHTTP,
		ClientRPCPortHTTP:         uint64(fileConfig.ClientRPCPortHTTP),
		HasClientRPCWebsockets:    fileConfig.HasClientRPCWebsockets,
		ClientRPCPortWS:           uint64(fileConfig.ClientRPCPortWS),
		ClientRPCHost:             fileConfig.ClientRPCHost,
		EnclaveRPCAddress:         fileConfig.EnclaveRPCAddress,
		P2PBindAddress:            fileConfig.P2PBindAddress,
		P2PPublicAddress:          fileConfig.P2PPublicAddress,
		L1WebsocketURL:            fileConfig.L1WebsocketURL,
		EnclaveRPCTimeout:         time.Duration(fileConfig.EnclaveRPCTimeout) * time.Second,
		L1RPCTimeout:              time.Duration(fileConfig.L1RPCTimeout) * time.Second,
		P2PConnectionTimeout:      time.Duration(fileConfig.P2PConnectionTimeout) * time.Second,
		ManagementContractAddress: gethcommon.HexToAddress(fileConfig.ManagementContractAddress),
		MessageBusAddress:         gethcommon.HexToAddress(fileConfig.MessageBusAddress),
		LogLevel:                  fileConfig.LogLevel,
		LogPath:                   fileConfig.LogPath,
		PrivateKeyString:          fileConfig.PrivateKeyString,
		L1ChainID:                 fileConfig.L1ChainID,
		ObscuroChainID:            fileConfig.ObscuroChainID,
		ProfilerEnabled:           fileConfig.ProfilerEnabled,
		L1StartHash:               gethcommon.HexToHash(fileConfig.L1StartHash),
		SequencerID:               gethcommon.HexToAddress(fileConfig.SequencerID),
		MetricsEnabled:            fileConfig.MetricsEnabled,
		MetricsHTTPPort:           fileConfig.MetricsHTTPPort,
		UseInMemoryDB:             fileConfig.UseInMemoryDB,
		LevelDBPath:               fileConfig.LevelDBPath,
		BatchInterval:             batchInterval,
		MaxBatchInterval:          maxBatchInterval,
		RollupInterval:            rollupInterval,
		IsInboundP2PDisabled:      fileConfig.IsInboundP2PDisabled,
		L1BlockTime:               time.Duration(fileConfig.L1BlockTime) * time.Second,
	}, nil
}
