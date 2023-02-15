package main

import (
	"flag"
)

// NodeConfigCLI represents the configurations passed into the node over CLI
type NodeConfigCLI struct {
	nodeType               string
	isGenesis              bool
	isSGXEnabled           bool
	enclaveDockerImage     string
	hostDockerImage        string
	l1Addr                 string
	l1WSPort               int
	hostP2PPort            int
	hostP2PAddr            string
	enclaveHTTPPort        int
	enclaveWSPort          int
	privateKey             string
	hostID                 string
	sequencerID            string
	managementContractAddr string
	messageBusContractAddr string
	pccsAddr               string
	edgelessDBImage        string
	hostHTTPPort           int
	hostWSPort             int
}

// ParseConfigCLI returns a NodeConfigCLI based the cli params and defaults.
func ParseConfigCLI() *NodeConfigCLI {
	cfg := &NodeConfigCLI{}
	flagUsageMap := getFlagUsageMap()

	nodeType := flag.String(nodeTypeFlag, "", flagUsageMap[nodeTypeFlag])
	isGenesis := flag.Bool(isGenesisFlag, true, flagUsageMap[isGenesisFlag])
	isSGXEnabled := flag.Bool(isSGXEnabledFlag, false, flagUsageMap[isSGXEnabledFlag])
	enclaveDockerImage := flag.String(enclaveDockerImageFlag, "", flagUsageMap[enclaveDockerImageFlag])
	hostDockerImage := flag.String(hostDockerImageFlag, "", flagUsageMap[hostDockerImageFlag])
	l1Addr := flag.String(l1AddrFlag, "eth2network", flagUsageMap[l1AddrFlag])
	l1WSPort := flag.Int(l1WSPortFlag, 9000, flagUsageMap[l1WSPortFlag])
	hostP2PPort := flag.Int(hostP2PPortFlag, 14000, flagUsageMap[hostP2PPortFlag])
	hostP2PAddr := flag.String(hostP2PAddrFlag, "0.0.0.0", flagUsageMap[hostP2PAddrFlag])
	hostHTTPPort := flag.Int(hostHTTPPortFlag, 12000, flagUsageMap[hostHTTPPortFlag])
	hostWSPort := flag.Int(hostWSPortFlag, 12001, flagUsageMap[hostWSPortFlag])
	enclaveHTTPPort := flag.Int(enclaveHTTPPortFlag, 13000, flagUsageMap[enclaveHTTPPortFlag])
	enclaveWSPort := flag.Int(enclaveWSPortFlag, 13001, flagUsageMap[enclaveWSPortFlag])
	privateKey := flag.String(privateKeyFlag, "", flagUsageMap[privateKeyFlag])
	hostID := flag.String(hostIDFlag, "", flagUsageMap[hostIDFlag])
	sequencerID := flag.String(sequencerIDFlag, "", flagUsageMap[sequencerIDFlag])
	managementContractAddr := flag.String(managementContractAddrFlag, "", flagUsageMap[managementContractAddrFlag])
	messageBusContractAddr := flag.String(messageBusContractAddrFlag, "", flagUsageMap[messageBusContractAddrFlag])
	pccsAddr := flag.String(pccsAddrFlag, "", flagUsageMap[pccsAddrFlag])
	edgelessDBImage := flag.String(edgelessDBImageFlag, "ghcr.io/edgelesssys/edgelessdb-sgx-4gb:v0.3.2", flagUsageMap[edgelessDBImageFlag])

	flag.Parse()
	cfg.nodeType = *nodeType
	cfg.isGenesis = *isGenesis
	cfg.isSGXEnabled = *isSGXEnabled
	cfg.enclaveDockerImage = *enclaveDockerImage
	cfg.hostDockerImage = *hostDockerImage
	cfg.l1Addr = *l1Addr
	cfg.l1WSPort = *l1WSPort
	cfg.hostP2PPort = *hostP2PPort
	cfg.hostP2PAddr = *hostP2PAddr
	cfg.enclaveHTTPPort = *enclaveHTTPPort
	cfg.enclaveWSPort = *enclaveWSPort
	cfg.privateKey = *privateKey
	cfg.hostID = *hostID
	cfg.sequencerID = *sequencerID
	cfg.managementContractAddr = *managementContractAddr
	cfg.messageBusContractAddr = *messageBusContractAddr
	cfg.pccsAddr = *pccsAddr
	cfg.edgelessDBImage = *edgelessDBImage
	cfg.hostHTTPPort = *hostHTTPPort
	cfg.hostWSPort = *hostWSPort

	return cfg
}
