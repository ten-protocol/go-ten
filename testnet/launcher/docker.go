package launcher

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/sanity-io/litter"
	"github.com/ten-protocol/go-ten/go/common/retry"
	"github.com/ten-protocol/go-ten/go/config"
	"github.com/ten-protocol/go-ten/go/node"
	"github.com/ten-protocol/go-ten/go/obsclient"
	"github.com/ten-protocol/go-ten/go/rpc"
	"github.com/ten-protocol/go-ten/testnet/launcher/eth2network"
	"github.com/ten-protocol/go-ten/testnet/launcher/faucet"
	"github.com/ten-protocol/go-ten/testnet/launcher/gateway"

	l1cd "github.com/ten-protocol/go-ten/testnet/launcher/l1contractdeployer"
	l1gs "github.com/ten-protocol/go-ten/testnet/launcher/l1grantsequencers"
	l2cd "github.com/ten-protocol/go-ten/testnet/launcher/l2contractdeployer"
)

type Testnet struct {
	cfg *Config
}

func NewTestnetLauncher(cfg *Config) *Testnet {
	// todo (@pedro) - bind testnet specific options like number of nodes, etc
	return &Testnet{cfg: cfg}
}

func (t *Testnet) Start() error {
	litter.Config.HidePrivateFields = true
	fmt.Printf("Starting Testnet with config: \n", litter.Sdump(*t.cfg))

	err := startEth2Network()
	if err != nil {
		return fmt.Errorf("unable to start eth2network - %w", err)
	}

	networkConfig, err := t.deployL1Contracts()
	if err != nil {
		return fmt.Errorf("unable to deploy l1 contracts - %w", err)
	}

	edgelessDBImage := "ghcr.io/edgelesssys/edgelessdb-sgx-4gb:v0.3.2"
	// todo: revisit how we should configure the image, this condition is not ideal
	if !t.cfg.isSGXEnabled {
		edgelessDBImage = "ghcr.io/edgelesssys/edgelessdb-sgx-1gb:v0.3.2"
	}

	sequencerCfg, err := config.LoadTenConfig(
		"defaults/testnet-launcher/1-testnet-launcher.yaml",
		"defaults/testnet-launcher/2-sequencer.yaml",
	)
	if err != nil {
		return fmt.Errorf("unable to load sequencer config - %w", err)
	}
	sequencerCfg.Network.L1.StartHash = common.HexToHash(networkConfig.L1StartHash)
	sequencerCfg.Network.L1.L1Contracts.EnclaveRegistryContract = common.HexToAddress(networkConfig.EnclaveRegistryAddress)
	sequencerCfg.Network.L1.L1Contracts.DataAvailabilityRegistry = common.HexToAddress(networkConfig.DataAvailabilityRegistryAddress)
	sequencerCfg.Network.L1.L1Contracts.CrossChainContract = common.HexToAddress(networkConfig.CrossChainAddress)
	sequencerCfg.Network.L1.L1Contracts.NetworkConfigContract = common.HexToAddress(networkConfig.NetworkConfigAddress)
	sequencerCfg.Network.L1.L1Contracts.MessageBusContract = common.HexToAddress(networkConfig.MessageBusAddress)

	sequencerNode := node.NewDockerNode(sequencerCfg, "testnetobscuronet.azurecr.io/obscuronet/host:latest", "testnetobscuronet.azurecr.io/obscuronet/enclave:latest", edgelessDBImage, false, "", 1)

	err = sequencerNode.Start()
	if err != nil {
		return fmt.Errorf("unable to start the TEN node - %w", err)
	}
	fmt.Println("TEN node was successfully started...")

	// wait until the node is healthy
	err = waitForHealthyNode(80)
	if err != nil {
		return fmt.Errorf("sequencer TEN node not healthy - %w", err)
	}

	validatorNodeCfg, err := config.LoadTenConfig(
		"defaults/testnet-launcher/1-testnet-launcher.yaml",
		"defaults/testnet-launcher/2-validator.yaml",
	)
	if err != nil {
		return fmt.Errorf("unable to load validator config - %w", err)
	}
	validatorNodeCfg.Network.L1.StartHash = common.HexToHash(networkConfig.L1StartHash)
	validatorNodeCfg.Network.L1.L1Contracts.EnclaveRegistryContract = common.HexToAddress(networkConfig.EnclaveRegistryAddress)
	validatorNodeCfg.Network.L1.L1Contracts.DataAvailabilityRegistry = common.HexToAddress(networkConfig.DataAvailabilityRegistryAddress)
	validatorNodeCfg.Network.L1.L1Contracts.CrossChainContract = common.HexToAddress(networkConfig.CrossChainAddress)
	validatorNodeCfg.Network.L1.L1Contracts.NetworkConfigContract = common.HexToAddress(networkConfig.NetworkConfigAddress)
	validatorNodeCfg.Network.L1.L1Contracts.MessageBusContract = common.HexToAddress(networkConfig.MessageBusAddress)

	validatorNode := node.NewDockerNode(validatorNodeCfg, "testnetobscuronet.azurecr.io/obscuronet/host:latest", "testnetobscuronet.azurecr.io/obscuronet/enclave:latest", edgelessDBImage, false, "", 1)
	err = validatorNode.Start()
	if err != nil {
		return fmt.Errorf("unable to start the obscuro node - %w", err)
	}
	fmt.Println("TEN validator node was successfully started...")

	// wait until the node it healthy
	err = waitForHealthyNode(13010)
	if err != nil {
		return fmt.Errorf("validator obscuro node not healthy - %w", err)
	}

	l2ContractDeployer, err := l2cd.NewDockerContractDeployer(
		l2cd.NewContractDeployerConfig(
			l2cd.WithL1HTTPURL("http://eth2network:8025"),
			l2cd.WithL2Host("sequencer-host"),
			l2cd.WithL2WSPort(81),
			l2cd.WithL1PrivateKey("f52e5418e349dccdda29b6ac8b0abe6576bb7713886aa85abea6181ba731f9bb"),
			l2cd.WithMessageBusContractAddress(networkConfig.MessageBusAddress),
			l2cd.WithNetworkConfigAddress(networkConfig.NetworkConfigAddress),
			l2cd.WithEnclaveRegistryAddress(networkConfig.EnclaveRegistryAddress),
			l2cd.WithDataAvailabilityRegistryAddress(networkConfig.DataAvailabilityRegistryAddress),
			l2cd.WithCrossChainAddress(networkConfig.CrossChainAddress),
			l2cd.WithL2PrivateKey("8dfb8083da6275ae3e4f41e3e8a8c19d028d32c9247e24530933782f2a05035b"),
			l2cd.WithDockerImage(t.cfg.contractDeployerDockerImage),
			l2cd.WithDebugEnabled(t.cfg.contractDeployerDebug),
			l2cd.WithFaucetFunds("10000"),
		),
	)
	if err != nil {
		return fmt.Errorf("unable to configure the l2 contract deployer - %w", err)
	}

	err = t.grantSequencerStatus(networkConfig.EnclaveRegistryAddress)
	if err != nil {
		return fmt.Errorf("failed to grant sequencer status: %w", err)
	}

	err = l2ContractDeployer.Start()
	if err != nil {
		return fmt.Errorf("unable to start the l2 contract deployer - %w", err)
	}

	err = l2ContractDeployer.WaitForFinish()
	if err != nil {
		return fmt.Errorf("unexpected error waiting for l2 contract deployer { ID = %s } to finish - %w", l2ContractDeployer.GetID(), err)
	}
	fmt.Println("L2 Contracts were successfully deployed...")

	faucetPort := 99
	faucetInst, err := faucet.NewDockerFaucet(
		faucet.NewFaucetConfig(
			faucet.WithFaucetPort(faucetPort),
			faucet.WithTenNodePort(13010),
			faucet.WithTenNodeHost("validator-host"),
			faucet.WithFaucetPrivKey("0x8dfb8083da6275ae3e4f41e3e8a8c19d028d32c9247e24530933782f2a05035b"),
			faucet.WithDockerImage("testnetobscuronet.azurecr.io/obscuronet/faucet:latest"),
		),
	)
	if err != nil {
		return fmt.Errorf("unable to instantiate faucet - %w", err)
	}

	if err = faucetInst.Start(); err != nil {
		return fmt.Errorf("unable to start faucet - %w", err)
	}

	if err = faucetInst.IsReady(); err != nil {
		return fmt.Errorf("unable to wait for faucet to be ready - %w", err)
	}

	fmt.Printf("Faucet ready to be accessed at http://127.0.0.1:%d/ ...\n", faucetPort)
	fmt.Printf("Fund your account with `curl --request POST 'http://127.0.0.1:%d/fund/eth' --header 'Content-Type: application/json' --data-raw '{ \"address\":\"0x0....\" } `\n", faucetPort)

	gatewayPort := 3000
	gatewayInst, err := gateway.NewDockerGateway(
		gateway.NewGatewayConfig(
			gateway.WithGatewayHTTPPort(gatewayPort),
			gateway.WithGatewayWSPort(3001),
			gateway.WithTenNodeHTTPPort(13010),
			gateway.WithTenNodeWSPort(13011),
			gateway.WithTenNodeHost("validator-host"),
			gateway.WithRateLimitUserComputeTime(0), // disable rate limiting for local network
			gateway.WithDockerImage("testnetobscuronet.azurecr.io/obscuronet/obscuro_gateway:latest"),
		),
	)
	if err != nil {
		return fmt.Errorf("unable to instantiate gateway - %w", err)
	}

	if err = gatewayInst.Start(); err != nil {
		return fmt.Errorf("unable to start gateway - %w", err)
	}

	if err = gatewayInst.IsReady(); err != nil {
		return fmt.Errorf("unable to wait for gateway to be ready - %w", err)
	}

	fmt.Printf("Gateway ready to be accessed at http://127.0.0.1:%d ...\n", gatewayPort)

	fmt.Println("Network successfully launched !")
	return nil
}

func startEth2Network() error {
	eth2Network, err := eth2network.NewDockerEth2Network(
		eth2network.NewEth2NetworkConfig(
			eth2network.WithGethHTTPStartPort(8025),
			eth2network.WithGethWSStartPort(9000),
			eth2network.WithGethPrefundedAddrs([]string{
				"0x13E23Ca74DE0206C56ebaE8D51b5622EFF1E9944", // contract deployment pk - f52e5418e349dccdda29b6ac8b0abe6576bb7713886aa85abea6181ba731f9bb
				"0x0654D8B60033144D567f25bF41baC1FB0D60F23B", // sequencer pk - 8ead642ca80dadb0f346a66cd6aa13e08a8ac7b5c6f7578d4bac96f5db01ac99
				"0x2f7fCaA34b38871560DaAD6Db4596860744e1e8A", // validator pk - ebca545772d6438bbbe1a16afbed455733eccf96157b52384f1722ea65ccfa89
				"0xE09a37ABc1A63441404007019E5BC7517bE2c43f", // bridge admin pk - 4bfe14725e685901c062ccd4e220c61cf9c189897b6c78bd18d7f51291b2b8f1
			}),
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

	return nil
}

func (t *Testnet) deployL1Contracts() (*node.NetworkConfig, error) {
	l1ContractDeployer, err := l1cd.NewDockerContractDeployer(
		l1cd.NewContractDeployerConfig(
			l1cd.WithL1HTTPURL("http://eth2network:8025"),
			l1cd.WithPrivateKey("f52e5418e349dccdda29b6ac8b0abe6576bb7713886aa85abea6181ba731f9bb"),
			l1cd.WithDockerImage(t.cfg.contractDeployerDockerImage),
			l1cd.WithDebugEnabled(t.cfg.contractDeployerDebug),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to configure l1 contract deployer - %w", err)
	}

	err = l1ContractDeployer.Start()
	if err != nil {
		return nil, fmt.Errorf("unable to start l1 contract deployer - %w", err)
	}

	networkConfig, err := l1ContractDeployer.RetrieveL1ContractAddresses()
	if err != nil {
		return nil, fmt.Errorf("unable to fetch l1 contract addresses - %w", err)
	}

	fmt.Printf("Network Config Values:\n"+
		"  NetworkConfig Address: %s\n"+
		"  CrossChain Address: %s\n"+
		"  MessageBus Address: %s\n"+
		"  EnclaveRegistry Address: %s\n"+
		"  DataAvailabilityRegistry Address: %s\n"+
		"  L1 Start Hash: %s\n",
		networkConfig.NetworkConfigAddress,
		networkConfig.CrossChainAddress,
		networkConfig.MessageBusAddress,
		networkConfig.EnclaveRegistryAddress,
		networkConfig.DataAvailabilityRegistryAddress,
		networkConfig.L1StartHash)
	fmt.Println("L1 Contracts were successfully deployed...")
	return networkConfig, nil
}

// waitForHealthyNode retries continuously for the node to respond to a healthcheck http request
func waitForHealthyNode(port int) error { // todo: hook the cfg
	timeStart := time.Now()

	hostURL := fmt.Sprintf("http://localhost:%d", port)
	fmt.Println("Waiting for Obscuro node to be healthy...")
	err := retry.Do(
		func() error {
			client, err := rpc.NewNetworkClient(hostURL)
			if err != nil {
				return err
			}
			defer client.Stop()

			obsClient := obsclient.NewObsClient(client)
			health, err := obsClient.Health()
			if err != nil {
				return err
			}
			if health.OverallHealth {
				fmt.Println("Obscuro node is ready")
				return nil
			}

			return fmt.Errorf("node OverallHealth is not good yet")
		}, retry.NewTimeoutStrategy(7*time.Minute, 1*time.Second),
	)
	if err != nil {
		return err
	}
	fmt.Printf("Node became healthy after %f seconds\n", time.Since(timeStart).Seconds())
	return nil
}

func (t *Testnet) grantSequencerStatus(enclaveRegistryAddr string) error {
	// fetch enclaveIDs
	hostURL := fmt.Sprintf("http://localhost:%d", 80)

	l1grantsequencers, err := l1gs.NewGrantSequencers(
		l1gs.NewGrantSequencerConfig(
			l1gs.WithL1HTTPURL("http://eth2network:8025"),
			l1gs.WithPrivateKey("f52e5418e349dccdda29b6ac8b0abe6576bb7713886aa85abea6181ba731f9bb"),
			l1gs.WithDockerImage(t.cfg.contractDeployerDockerImage),
			l1gs.WithEnclaveContractAddress(enclaveRegistryAddr),
			l1gs.WithSequencerURL(hostURL),
		),
	)
	if err != nil {
		return fmt.Errorf("unable to configure l1 grant sequencersr - %w", err)
	}

	err = l1grantsequencers.Start()
	if err != nil {
		return fmt.Errorf("unable to start l1 grant sequencers - %w", err)
	}

	err = l1grantsequencers.WaitForFinish()
	if err != nil {
		return fmt.Errorf("unable to wait for l1 grant sequencers to finish - %w", err)
	}

	fmt.Println("Enclaves were successfully granted sequencer roles...")

	return nil
}
