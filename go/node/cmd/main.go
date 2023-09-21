package main

import (
	"github.com/obscuronet/go-obscuro/go/node"
)

func main() {
	cliConfig := ParseConfigCLI()
	// todo (#1618) - allow for multiple operation (start, stop, status)

	nodeCfg := node.NewNodeConfig(
		node.WithNodeName(cliConfig.nodeName),
		node.WithNodeType(cliConfig.nodeType),
		node.WithGenesis(cliConfig.isGenesis),
		node.WithSGXEnabled(cliConfig.isSGXEnabled),
		node.WithEnclaveImage(cliConfig.enclaveDockerImage),                  // "local_enclave"
		node.WithHostImage(cliConfig.hostDockerImage),                        // "local_host"
		node.WithL1WebsocketURL(cliConfig.l1WebsocketURL),                    // "ws://eth2network:9000"
		node.WithHostP2PPort(cliConfig.hostP2PPort),                          // 14000
		node.WithHostP2PHost(cliConfig.hostP2PHost),                          // 0.0.0.0
		node.WithHostPublicP2PAddr(cliConfig.hostP2PPublicAddr),              // node public facing ip and port
		node.WithHostHTTPPort(cliConfig.hostHTTPPort),                        // 12000
		node.WithHostWSPort(cliConfig.hostWSPort),                            // 12001
		node.WithEnclaveWSPort(cliConfig.enclaveWSPort),                      // 81
		node.WithPrivateKey(cliConfig.privateKey),                            // "8ead642ca80dadb0f346a66cd6aa13e08a8ac7b5c6f7578d4bac96f5db01ac99"
		node.WithHostID(cliConfig.hostID),                                    // "0x0654D8B60033144D567f25bF41baC1FB0D60F23B"),
		node.WithSequencerID(cliConfig.sequencerID),                          // "0x0654D8B60033144D567f25bF41baC1FB0D60F23B"),
		node.WithManagementContractAddress(cliConfig.managementContractAddr), // "0xeDa66Cc53bd2f26896f6Ba6b736B1Ca325DE04eF"),
		node.WithMessageBusContractAddress(cliConfig.messageBusContractAddr), // "0xFD03804faCA2538F4633B3EBdfEfc38adafa259B"),
		node.WithL1Start(cliConfig.l1Start),
		node.WithPCCSAddr(cliConfig.pccsAddr),
		node.WithEdgelessDBImage(cliConfig.edgelessDBImage),
		node.WithDebugNamespaceEnabled(cliConfig.isDebugNamespaceEnabled), // false
		node.WithLogLevel(cliConfig.logLevel),
		node.WithInboundP2PDisabled(cliConfig.isInboundP2PDisabled),
		node.WithBatchInterval(cliConfig.batchInterval),
		node.WithRollupInterval(cliConfig.rollupInterval),
		node.WithL1ChainID(cliConfig.l1ChainID),
	)

	dockerNode := node.NewDockerNode(nodeCfg)
	var err error
	switch cliConfig.nodeAction {
	case startAction:
		// write the network-level config to disk for future restarts
		err = node.WriteNetworkConfigToDisk(nodeCfg)
		if err != nil {
			panic(err)
		}
		err = dockerNode.Start()
	case upgradeAction:
		// load network-specific details from the initial node setup from disk
		var ntwCfg *node.NetworkConfig
		ntwCfg, err = node.ReadNetworkConfigFromDisk()
		if err != nil {
			panic(err)
		}

		err = dockerNode.Upgrade(ntwCfg)
	default:
		panic("unrecognized node action: " + cliConfig.nodeAction)
	}
	if err != nil {
		panic(err)
	}
}
