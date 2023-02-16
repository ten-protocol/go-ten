package launcher

import (
	"fmt"

	"github.com/obscuronet/go-obscuro/go/node"
	"github.com/obscuronet/go-obscuro/testnet/launcher/eth2network"

	l1cd "github.com/obscuronet/go-obscuro/testnet/launcher/l1contractdeployer"
	l2cd "github.com/obscuronet/go-obscuro/testnet/launcher/l2contractdeployer"
)

type Testnet struct{}

func NewTestnetLauncher(*Config) *Testnet {
	// todo bind testnet specific options like number of nodes, etc
	return &Testnet{}
}

func (t *Testnet) Start() error {
	eth2Network, err := eth2network.NewDockerEth2Network(
		eth2network.NewEth2NetworkConfig(
			eth2network.WithGethHTTPStartPort(8025),
			eth2network.WithGethWSStartPort(9000),
			eth2network.WithGethPrefundedAddrs([]string{"0x13E23Ca74DE0206C56ebaE8D51b5622EFF1E9944", "0x0654D8B60033144D567f25bF41baC1FB0D60F23B"}),
		),
	)
	if err != nil {
		return fmt.Errorf("unable to configure eth2network - %w", err)
	}

	err = eth2Network.Start()
	if err != nil {
		return fmt.Errorf("unable to start eth2network - %w", err)
	}
	fmt.Println("Eth2 network started...")

	err = eth2Network.IsReady()
	if err != nil {
		return fmt.Errorf("eth2network not ready in time - %w", err)
	}
	fmt.Println("Eth2 network is ready...")

	l1ContractDeployer, err := l1cd.NewDockerContractDeployer(
		l1cd.NewContractDeployerConfig(
			l1cd.WithL1Host("eth2network"),
			l1cd.WithL1Port(8025),
			l1cd.WithPrivateKey("f52e5418e349dccdda29b6ac8b0abe6576bb7713886aa85abea6181ba731f9bb"),
			l1cd.WithDockerImage("testnetobscuronet.azurecr.io/obscuronet/hardhatdeployer:latest"),
		),
	)
	if err != nil {
		return fmt.Errorf("unable to configure l1 contract deployer - %w", err)
	}

	err = l1ContractDeployer.Start()
	if err != nil {
		return fmt.Errorf("unable to start l1 contract deployer - %w", err)
	}

	managementContractAddr, messageBusContractAddr, err := l1ContractDeployer.RetrieveL1ContractAddresses()
	if err != nil {
		return fmt.Errorf("unable to fetch l1 contract addresses - %w", err)
	}
	fmt.Println("L1 Contracts were successfully deployed...")

	nodeCfg := node.NewNodeConfig(
		node.WithNodeType("sequencer"),
		node.WithGenesis(true),
		node.WithSGXEnabled(false),
		node.WithEnclaveImage("testnetobscuronet.azurecr.io/obscuronet/enclave:latest"),
		node.WithHostImage("testnetobscuronet.azurecr.io/obscuronet/host:latest"),
		node.WithL1Host("eth2network"),
		node.WithL1WSPort(9000),
		node.WithHostP2PPort(14000),
		node.WithEnclaveHTTPPort(13000),
		node.WithEnclaveWSPort(13001),
		node.WithHostHTTPPort(12000),
		node.WithHostWSPort(12000),
		node.WithPrivateKey("8ead642ca80dadb0f346a66cd6aa13e08a8ac7b5c6f7578d4bac96f5db01ac99"),
		node.WithHostID("0x0654D8B60033144D567f25bF41baC1FB0D60F23B"),
		node.WithSequencerID("0x0654D8B60033144D567f25bF41baC1FB0D60F23B"),
		node.WithManagementContractAddress(managementContractAddr),
		node.WithMessageBusContractAddress(messageBusContractAddr),
	)

	dockerNode, err := node.NewDockerNode(nodeCfg)
	if err != nil {
		return fmt.Errorf("unable to configure the obscuro node - %w", err)
	}

	err = dockerNode.Start()
	if err != nil {
		return fmt.Errorf("unable to start the obscuro node - %w", err)
	}
	fmt.Println("Obscuro node was successfully started...")

	l2ContractDeployer, err := l2cd.NewDockerContractDeployer(
		l2cd.NewContractDeployerConfig(
			l2cd.WithL1Host("eth2network"),
			l2cd.WithL1Port(8025),
			l2cd.WithL2Host("host"),
			l2cd.WithL2WSPort(13001),
			l2cd.WithL1PrivateKey("f52e5418e349dccdda29b6ac8b0abe6576bb7713886aa85abea6181ba731f9bb"),
			l2cd.WithMessageBusContractAddress("0xFD03804faCA2538F4633B3EBdfEfc38adafa259B"),
			l2cd.WithL2PrivateKey("8dfb8083da6275ae3e4f41e3e8a8c19d028d32c9247e24530933782f2a05035b"),
			l2cd.WithHocPKString("6e384a07a01263518a09a5424c7b6bbfc3604ba7d93f47e3a455cbdd7f9f0682"),
			l2cd.WithPocPKString("4bfe14725e685901c062ccd4e220c61cf9c189897b6c78bd18d7f51291b2b8f8"),
			l2cd.WithDockerImage("testnetobscuronet.azurecr.io/obscuronet/hardhatdeployer:latest"),
		),
	)
	if err != nil {
		return fmt.Errorf("unable to configure the l2 contract deployer - %w", err)
	}

	err = l2ContractDeployer.Start()
	if err != nil {
		return fmt.Errorf("unable to start the l2 contract deployer - %w", err)
	}
	fmt.Println("L2 Contracts were successfully deployed...")

	return nil
}
