package main

import (
	"flag"
	"fmt"
	"github.com/ten-protocol/go-ten/go/node"
	"os"
	"strings"
)

var (
	startAction      = "start"
	upgradeAction    = "upgrade"
	validNodeActions = []string{startAction, upgradeAction}
)

// LoadDefaultConfig parses optional or default configuration file and returns struct.
func LoadDefaultConfig() *node.Config {
	flagUsageMap := getFlagUsageMap()
	configPath := flag.String("config", "./go/config/default_node_config.yaml", flagUsageMap[configFlag])
	overridePath := flag.String("override", "", flagUsageMap[overrideFlag])

	// Parse only once capturing all necessary flags
	flag.Parse()

	defaults, err := node.LoadConfig(*configPath)
	if err != nil {
		panic(err)
	}

	// Apply overrides if the override path is provided
	if *overridePath != "" {
		overrides, err := node.LoadConfig(*overridePath)
		if err != nil {
			panic(err)
		}
		defaults.ApplyOverrides(overrides)
	}

	return defaults
}

// ParseConfigCLI returns a NodeConfig based on the CLI params and defaults from YAML.
func ParseConfigCLI(defaults *node.Config) *node.Config {
	flagUsageMap := getFlagUsageMap()

	nodeAction := flag.String(nodeActionFlag, defaults.NodeAction, flagUsageMap[nodeActionFlag])
	nodeName := flag.String(nodeNameFlag, defaults.NodeName, flagUsageMap[nodeNameFlag])
	isGenesis := flag.Bool(isGenesisFlag, defaults.IsGenesis, flagUsageMap[isGenesisFlag])
	isSGXEnabled := flag.Bool(isSGXEnabledFlag, defaults.IsSGXEnabled, flagUsageMap[isSGXEnabledFlag])
	enclaveDockerImage := flag.String(enclaveDockerImageFlag, defaults.EnclaveDockerImage, flagUsageMap[enclaveDockerImageFlag])
	hostDockerImage := flag.String(hostDockerImageFlag, defaults.HostDockerImage, flagUsageMap[hostDockerImageFlag])
	l1WebsocketURL := flag.String(l1WebsocketURLFlag, defaults.L1WebsocketURL, flagUsageMap[l1WebsocketURLFlag])
	hostP2PPort := flag.Int(hostP2PPortFlag, defaults.HostP2PPort, flagUsageMap[hostP2PPortFlag])
	hostP2PPublicAddr := flag.String(hostP2PPublicAddrFlag, defaults.HostP2PPublicAddr, flagUsageMap[hostP2PPublicAddrFlag])
	hostHTTPPort := flag.Int(hostHTTPPortFlag, defaults.HostHTTPPort, flagUsageMap[hostHTTPPortFlag])
	hostWSPort := flag.Int(hostWSPortFlag, defaults.HostWSPort, flagUsageMap[hostWSPortFlag])
	enclaveWSPort := flag.Int(enclaveWSPortFlag, defaults.EnclaveWSPort, flagUsageMap[enclaveWSPortFlag])
	privateKey := flag.String(privateKeyFlag, defaults.PrivateKey, flagUsageMap[privateKeyFlag])
	hostID := flag.String(hostIDFlag, defaults.HostID, flagUsageMap[hostIDFlag])
	sequencerID := flag.String(sequencerIDFlag, defaults.SequencerID, flagUsageMap[sequencerIDFlag])
	managementContractAddr := flag.String(managementContractAddrFlag, defaults.ManagementContractAddr, flagUsageMap[managementContractAddrFlag])
	messageBusContractAddr := flag.String(messageBusContractAddrFlag, defaults.MessageBusContractAddr, flagUsageMap[messageBusContractAddrFlag])
	l1Start := flag.String(l1StartBlockFlag, defaults.L1Start, flagUsageMap[l1StartBlockFlag])
	pccsAddr := flag.String(pccsAddrFlag, defaults.PccsAddr, flagUsageMap[pccsAddrFlag])
	edgelessDBImage := flag.String(edgelessDBImageFlag, defaults.EdgelessDBImage, flagUsageMap[edgelessDBImageFlag])
	isDebugNamespaceEnabled := flag.Bool(isDebugNamespaceEnabledFlag, defaults.IsDebugNamespaceEnabled, flagUsageMap[isDebugNamespaceEnabledFlag])
	logLevel := flag.Int(logLevelFlag, defaults.LogLevel, flagUsageMap[logLevelFlag])
	isInboundP2PDisabled := flag.Bool(isInboundP2PDisabledFlag, defaults.IsInboundP2PDisabled, flagUsageMap[isInboundP2PDisabledFlag])
	batchInterval := flag.String(batchIntervalFlag, defaults.BatchInterval, flagUsageMap[batchIntervalFlag])
	maxBatchInterval := flag.String(maxBatchIntervalFlag, defaults.MaxBatchInterval, flagUsageMap[maxBatchIntervalFlag])
	rollupInterval := flag.String(rollupIntervalFlag, defaults.RollupInterval, flagUsageMap[rollupIntervalFlag])
	l1ChainID := flag.Int(l1ChainIDFlag, defaults.L1ChainID, flagUsageMap[l1ChainIDFlag])
	profilerEnabled := flag.Bool(profilerEnabledFlag, defaults.ProfilerEnabled, flagUsageMap[profilerEnabledFlag])
	metricsEnabled := flag.Bool(metricsEnabledFlag, defaults.MetricsEnabled, flagUsageMap[metricsEnabledFlag])
	coinbaseAddress := flag.String(coinbaseAddressFlag, defaults.CoinbaseAddress, flagUsageMap[coinbaseAddressFlag])
	l1BlockTime := flag.Int(l1BlockTimeFlag, defaults.L1BlockTime, flagUsageMap[l1BlockTimeFlag])
	tenGenesis := flag.String(tenGenesisFlag, defaults.TenGenesis, flagUsageMap[tenGenesisFlag])
	enclaveDebug := flag.Bool(enclaveDebugFlag, defaults.EnclaveDebug, flagUsageMap[enclaveDebugFlag])
	hostInMemDB := flag.Bool(hostInMemDBFlag, defaults.HostInMemDB, flagUsageMap[hostInMemDBFlag])
	hostExternalDBHost := flag.String(hostExternalDBHostFlag, defaults.HostExternalDBHost, flagUsageMap[hostExternalDBHostFlag])
	hostExternalDBHostUser := flag.String(hostExternalDBHostUserFlag, defaults.HostExternalDBUser, flagUsageMap[hostExternalDBHostUserFlag])
	hostExternalDBHostPass := flag.String(hostExternalDBHostPassFlag, defaults.HostExternalDBPass, flagUsageMap[hostExternalDBHostPassFlag])

	flag.Parse()

	defaults.NodeAction = *nodeAction
	defaults.NodeName = *nodeName
	defaults.IsGenesis = *isGenesis
	defaults.IsSGXEnabled = *isSGXEnabled
	defaults.EnclaveDockerImage = *enclaveDockerImage
	defaults.HostDockerImage = *hostDockerImage
	defaults.L1WebsocketURL = *l1WebsocketURL
	defaults.HostP2PPort = *hostP2PPort
	defaults.HostP2PPublicAddr = *hostP2PPublicAddr
	defaults.HostHTTPPort = *hostHTTPPort
	defaults.HostWSPort = *hostWSPort
	defaults.EnclaveWSPort = *enclaveWSPort
	defaults.PrivateKey = *privateKey
	defaults.HostID = *hostID
	defaults.SequencerID = *sequencerID
	defaults.ManagementContractAddr = *managementContractAddr
	defaults.MessageBusContractAddr = *messageBusContractAddr
	defaults.L1Start = *l1Start
	defaults.PccsAddr = *pccsAddr
	defaults.EdgelessDBImage = *edgelessDBImage
	defaults.IsDebugNamespaceEnabled = *isDebugNamespaceEnabled
	defaults.LogLevel = *logLevel
	defaults.IsInboundP2PDisabled = *isInboundP2PDisabled
	defaults.BatchInterval = *batchInterval
	defaults.MaxBatchInterval = *maxBatchInterval
	defaults.RollupInterval = *rollupInterval
	defaults.L1ChainID = *l1ChainID
	defaults.ProfilerEnabled = *profilerEnabled
	defaults.MetricsEnabled = *metricsEnabled
	defaults.CoinbaseAddress = *coinbaseAddress
	defaults.L1BlockTime = *l1BlockTime
	defaults.TenGenesis = *tenGenesis
	defaults.EnclaveDebug = *enclaveDebug
	defaults.HostInMemDB = *hostInMemDB
	defaults.HostExternalDBHost = *hostExternalDBHost
	defaults.HostExternalDBUser = *hostExternalDBHostUser
	defaults.HostExternalDBPass = *hostExternalDBHostPass

	if !validateNodeAction(defaults.NodeAction) {
		if defaults.NodeAction == "" {
			fmt.Printf("expected a node action string (%s) as the only argument after the flags but no argument provided\n",
				strings.Join(validNodeActions, ", "))
		} else {
			fmt.Printf("expected a node action string (%s) as the only argument after the flags but got %s\n",
				strings.Join(validNodeActions, ", "), defaults.NodeAction)
		}
		os.Exit(1)
	}

	return defaults
}

func validateNodeAction(action string) bool {
	for _, a := range validNodeActions {
		if a == action {
			return true
		}
	}
	return false
}
