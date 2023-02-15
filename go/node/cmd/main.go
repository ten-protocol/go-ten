package main

import (
	node2 "github.com/obscuronet/go-obscuro/go/node"
)

func main() {
	// todo: hook up config / CLI flags

	nodeCfg := node2.NewNodeConfig(
		node2.WithNodeType("sequencer"),
		node2.WithGenesis(true),
		node2.WithSGXEnabled(false),
		node2.WithEnclaveImage("local_enclave"),
		node2.WithHostImage("local_host"),
		node2.WithL1Host("eth2network"),
		node2.WithL1WSPort(9000),
		node2.WithHostP2PPort(14000),
		node2.WithEnclaveHTTPPort(13000),
		node2.WithEnclaveWSPort(13001),
		node2.WithPrivateKey("8ead642ca80dadb0f346a66cd6aa13e08a8ac7b5c6f7578d4bac96f5db01ac99"),
		node2.WithHostID("0x0654D8B60033144D567f25bF41baC1FB0D60F23B"),
		node2.WithSequencerID("0x0654D8B60033144D567f25bF41baC1FB0D60F23B"),
		node2.WithManagementContractAddress("0xeDa66Cc53bd2f26896f6Ba6b736B1Ca325DE04eF"),
		node2.WithMessageBusContractAddress("0xFD03804faCA2538F4633B3EBdfEfc38adafa259B"),
	)

	dockerNode, err := node2.NewDockerNode(nodeCfg)
	if err != nil {
		panic(err)
	}

	err = dockerNode.Start()
	if err != nil {
		panic(err)
	}
}
