package launcher

import (
	"testing"

	"github.com/obscuronet/go-obscuro/node"
)

func TestName(t *testing.T) {
	eth2Network, err := NewDockerEth2Network(
		NewEth2NetworkConfig(
			WithGethHTTPStartPort(8025),
			WithGethWSStartPort(9000),
		),
	)
	if err != nil {
		panic(err)
	}

	err = eth2Network.Start()
	if err != nil {
		panic(err)
	}

	err = eth2Network.IsReady()
	if err != nil {
		panic(err)
	}

	l1ContractDeployer, err := NewDockerContractDeployer(
		NewContractDeployerConfig(),
	)
	if err != nil {
		panic(err)
	}

	err = l1ContractDeployer.Start()
	if err != nil {
		panic(err)
	}

	nodeCfg := node.NewNodeConfig(
		node.WithNodeType("sequencer"),
		node.WithGenesis(true),
		node.WithSGXEnabled(false),
		node.WithEnclaveImage("local_enclave"),
		node.WithHostImage("local_host"),
		node.WithL1Host("eth2network"),
		node.WithL1WSPort(9000),
		node.WithHostP2PPort(14000),
		node.WithEnclaveHTTPPort(13000),
		node.WithEnclaveWSPort(13001),
		node.WithPrivateKey("8ead642ca80dadb0f346a66cd6aa13e08a8ac7b5c6f7578d4bac96f5db01ac99"),
		node.WithHostID("0x0654D8B60033144D567f25bF41baC1FB0D60F23B"),
		node.WithSequencerID("0x0654D8B60033144D567f25bF41baC1FB0D60F23B"),
		node.WithManagementContractAddress("0xeDa66Cc53bd2f26896f6Ba6b736B1Ca325DE04eF"),
		node.WithMessageBusContractAddress("0xFD03804faCA2538F4633B3EBdfEfc38adafa259B"),
	)

	dockerNode, err := node.NewDockerNode(nodeCfg)
	if err != nil {
		panic(err)
	}

	err = dockerNode.Start()
	if err != nil {
		panic(err)
	}
}
