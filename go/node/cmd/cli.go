package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	startAction      = "start"
	upgradeAction    = "upgrade"
	validNodeActions = []string{startAction, upgradeAction}
)

// NodeConfigCLI represents the configurations passed into the node over CLI
type NodeConfigCLI struct {
	nodeAction              string
	nodeType                string
	isGenesis               bool
	isSGXEnabled            bool
	enclaveDockerImage      string
	hostDockerImage         string
	l1WebsocketURL          string
	hostP2PPort             int
	hostP2PHost             string
	hostP2PPublicAddr       string
	enclaveHTTPPort         int
	enclaveWSPort           int
	privateKey              string
	hostID                  string
	sequencerID             string
	managementContractAddr  string
	messageBusContractAddr  string
	l1Start                 string
	pccsAddr                string
	edgelessDBImage         string
	hostHTTPPort            int
	hostWSPort              int
	nodeName                string
	isDebugNamespaceEnabled bool
	logLevel                int
	isInboundP2PDisabled    bool
	batchInterval           string // format like 500ms or 2s (any time parsable by time.ParseDuration())
	rollupInterval          string // format like 500ms or 2s (any time parsable by time.ParseDuration())
	l1ChainID               int
}

// ParseConfigCLI returns a NodeConfigCLI based the cli params and defaults.
func ParseConfigCLI() *NodeConfigCLI {
	cfg := &NodeConfigCLI{}
	flagUsageMap := getFlagUsageMap()

	nodeName := flag.String(nodeNameFlag, "obscuronode", flagUsageMap[nodeNameFlag])
	nodeType := flag.String(nodeTypeFlag, "", flagUsageMap[nodeTypeFlag])
	isGenesis := flag.Bool(isGenesisFlag, false, flagUsageMap[isGenesisFlag])
	isSGXEnabled := flag.Bool(isSGXEnabledFlag, false, flagUsageMap[isSGXEnabledFlag])
	enclaveDockerImage := flag.String(enclaveDockerImageFlag, "", flagUsageMap[enclaveDockerImageFlag])
	hostDockerImage := flag.String(hostDockerImageFlag, "", flagUsageMap[hostDockerImageFlag])
	l1WebsocketURL := flag.String(l1WebsocketURLFlag, "ws://eth2network:9000", flagUsageMap[l1WebsocketURLFlag])
	hostP2PPort := flag.Int(hostP2PPortFlag, 14000, flagUsageMap[hostP2PPortFlag])
	hostP2PHost := flag.String(hostP2PHostFlag, "0.0.0.0", flagUsageMap[hostP2PHostFlag])
	hostP2PPublicAddr := flag.String(hostP2PPublicAddrFlag, "", flagUsageMap[hostP2PPublicAddrFlag])
	hostHTTPPort := flag.Int(hostHTTPPortFlag, 80, flagUsageMap[hostHTTPPortFlag])
	hostWSPort := flag.Int(hostWSPortFlag, 81, flagUsageMap[hostWSPortFlag])
	enclaveHTTPPort := flag.Int(enclaveHTTPPortFlag, 11000, flagUsageMap[enclaveHTTPPortFlag])
	enclaveWSPort := flag.Int(enclaveWSPortFlag, 11001, flagUsageMap[enclaveWSPortFlag])
	privateKey := flag.String(privateKeyFlag, "", flagUsageMap[privateKeyFlag])
	hostID := flag.String(hostIDFlag, "", flagUsageMap[hostIDFlag])
	sequencerID := flag.String(sequencerIDFlag, "", flagUsageMap[sequencerIDFlag])
	managementContractAddr := flag.String(managementContractAddrFlag, "", flagUsageMap[managementContractAddrFlag])
	messageBusContractAddr := flag.String(messageBusContractAddrFlag, "", flagUsageMap[messageBusContractAddrFlag])
	l1Start := flag.String(l1StartBlockFlag, "", flagUsageMap[l1StartBlockFlag])
	pccsAddr := flag.String(pccsAddrFlag, "", flagUsageMap[pccsAddrFlag])
	edgelessDBImage := flag.String(edgelessDBImageFlag, "ghcr.io/edgelesssys/edgelessdb-sgx-4gb:v0.3.2", flagUsageMap[edgelessDBImageFlag])
	isDebugNamespaceEnabled := flag.Bool(isDebugNamespaceEnabledFlag, false, flagUsageMap[isDebugNamespaceEnabledFlag])
	logLevel := flag.Int(logLevelFlag, 3, flagUsageMap[logLevelFlag])
	isInboundP2PDisabled := flag.Bool(isInboundP2PDisabledFlag, false, flagUsageMap[isInboundP2PDisabledFlag])
	batchInterval := flag.String(batchIntervalFlag, "1s", flagUsageMap[batchIntervalFlag])
	rollupInterval := flag.String(rollupIntervalFlag, "3s", flagUsageMap[rollupIntervalFlag])
	l1ChainID := flag.Int(l1ChainIDFlag, 1337, flagUsageMap[l1ChainIDFlag])

	flag.Parse()
	cfg.nodeName = *nodeName
	cfg.nodeType = *nodeType
	cfg.isGenesis = *isGenesis
	cfg.isSGXEnabled = *isSGXEnabled
	cfg.enclaveDockerImage = *enclaveDockerImage
	cfg.hostDockerImage = *hostDockerImage
	cfg.l1WebsocketURL = *l1WebsocketURL
	cfg.hostP2PPort = *hostP2PPort
	cfg.hostP2PHost = *hostP2PHost
	cfg.hostP2PPublicAddr = *hostP2PPublicAddr
	cfg.enclaveHTTPPort = *enclaveHTTPPort
	cfg.enclaveWSPort = *enclaveWSPort
	cfg.privateKey = *privateKey
	cfg.hostID = *hostID
	cfg.sequencerID = *sequencerID
	cfg.managementContractAddr = *managementContractAddr
	cfg.messageBusContractAddr = *messageBusContractAddr
	cfg.l1Start = *l1Start
	cfg.pccsAddr = *pccsAddr
	cfg.edgelessDBImage = *edgelessDBImage
	cfg.hostHTTPPort = *hostHTTPPort
	cfg.hostWSPort = *hostWSPort
	cfg.isDebugNamespaceEnabled = *isDebugNamespaceEnabled
	cfg.logLevel = *logLevel
	cfg.isInboundP2PDisabled = *isInboundP2PDisabled
	cfg.batchInterval = *batchInterval
	cfg.rollupInterval = *rollupInterval
	cfg.l1ChainID = *l1ChainID

	cfg.nodeAction = flag.Arg(0)
	if !validateNodeAction(cfg.nodeAction) {
		if cfg.nodeAction == "" {
			fmt.Printf("expected a node action string (%s) as the only argument after the flags but no argument provided\n",
				strings.Join(validNodeActions, ", "))
		} else {
			fmt.Printf("expected a node action string (%s) as the only argument after the flags but got %s\n",
				strings.Join(validNodeActions, ", "), cfg.nodeAction)
		}
		os.Exit(1)
	}

	return cfg
}

func validateNodeAction(action string) bool {
	for _, a := range validNodeActions {
		if a == action {
			return true
		}
	}
	return false
}
