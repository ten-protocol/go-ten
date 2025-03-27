package network

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ten-protocol/go-ten/contracts/generated/MerkleTreeMessageBus"

	"github.com/ten-protocol/go-ten/contracts/generated/NetworkConfig"
	"github.com/ten-protocol/go-ten/contracts/generated/NetworkEnclaveRegistry"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ten-protocol/go-ten/contracts/generated/CrossChain"
	"github.com/ten-protocol/go-ten/contracts/generated/DataAvailabilityRegistry"
	"github.com/ten-protocol/go-ten/go/common/constants"
	"github.com/ten-protocol/go-ten/go/ethadapter"
	"github.com/ten-protocol/go-ten/go/wallet"
	"github.com/ten-protocol/go-ten/integration"
	integrationCommon "github.com/ten-protocol/go-ten/integration/common"
	"github.com/ten-protocol/go-ten/integration/common/testlog"
	"github.com/ten-protocol/go-ten/integration/erc20contract"
	"github.com/ten-protocol/go-ten/integration/eth2network"
	"github.com/ten-protocol/go-ten/integration/simulation/params"
)

const _contractTimeout = time.Second * 30

func SetUpGethNetwork(wallets *params.SimWallets, startPort int, nrNodes int) (*params.L1TenData, []ethadapter.EthClient, eth2network.PosEth2Network) {
	eth2Network, err := StartGethNetwork(wallets, startPort)
	if err != nil {
		panic(fmt.Errorf("error starting geth network %w", err))
	}

	// connect to the first host to deploy
	tmpEthClient, err := ethadapter.NewEthClient(Localhost, uint(startPort+100), DefaultL1RPCTimeout, testlog.Logger())
	if err != nil {
		panic(fmt.Errorf("error connecting to te first host %w", err))
	}

	l1Data, err := DeployTenNetworkContracts(tmpEthClient, wallets, true)
	if err != nil {
		panic(fmt.Errorf("error deploying contracts. Cause: %w", err))
	}

	ethClients := make([]ethadapter.EthClient, nrNodes)
	for i := 0; i < nrNodes; i++ {
		ethClients[i] = CreateEthClientConnection(uint(startPort + 100))
	}

	return l1Data, ethClients, eth2Network
}

func StartGethNetwork(wallets *params.SimWallets, startPort int) (eth2network.PosEth2Network, error) {
	// make sure the geth network binaries exist
	binDir, err := eth2network.EnsureBinariesExist()
	if err != nil {
		return nil, err
	}

	// get the node wallet addresses to prefund them with Eth, so they can submit rollups, deploy contracts, deposit to the bridge, etc
	walletAddresses := []string{integration.GethNodeAddress}
	for _, w := range wallets.AllEthWallets() {
		walletAddresses = append(walletAddresses, w.Address().String())
	}

	network := eth2network.NewPosEth2Network(
		binDir,
		false,
		startPort+integration.DefaultGethNetworkPortOffset,
		startPort+integration.DefaultPrysmP2PPortOffset,
		startPort+integration.DefaultGethAUTHPortOffset,
		startPort+integration.DefaultGethWSPortOffset,
		startPort+integration.DefaultGethHTTPPortOffset,
		startPort+integration.DefaultPrysmRPCPortOffset,
		startPort+integration.DefaultPrysmGatewayPortOffset,
		integration.EthereumChainID,
		3*time.Minute,
		walletAddresses...,
	)

	fmt.Println("Starting Geth network with WS port", startPort+integration.DefaultGethWSPortOffset)
	err = network.Start()
	if err != nil {
		return nil, err
	}

	return network, nil
}

func DeployTenNetworkContracts(client ethadapter.EthClient, wallets *params.SimWallets, deployERC20s bool) (*params.L1TenData, error) {
	_, enclaveRegistryReceipt, err := deployEnclaveRegistryContract(client, wallets.ContractOwnerWallet)
	if err != nil {
		return nil, err
	}

	crossChainContract, crossChainReceipt, err := deployCrossChainContract(client, wallets.ContractOwnerWallet)
	if err != nil {
		return nil, err
	}

	// Get the MessageBus address from the CrossChain contract
	messageBusAddr, err := crossChainContract.MessageBus(&bind.CallOpts{})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch MessageBus address. Cause: %w", err)
	}

	// Wait for MessageBus contract to be deployed
	err = waitForContractDeployment(context.Background(), client, messageBusAddr)
	if err != nil {
		return nil, fmt.Errorf("MessageBus contract not available after deployment: %w", err)
	}

	_, daRegistryReceipt, err := deployDataAvailabilityRegistry(client, wallets.ContractOwnerWallet, messageBusAddr, enclaveRegistryReceipt.ContractAddress)
	if err != nil {
		return nil, err
	}

	merkleTreeMessageBus, err := MerkleTreeMessageBus.NewMerkleTreeMessageBus(messageBusAddr, client.EthClient())
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate MerkleTreeMessageBus contract. Cause: %w", err)
	}

	opts, err := createTransactor(wallets.ContractOwnerWallet)
	if err != nil {
		return nil, err
	}

	tx, err := merkleTreeMessageBus.AddStateRootManager(opts, daRegistryReceipt.ContractAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to add state root manager to MerkleTreeMessageBus contract. Cause: %w", err)
	}

	_, err = integrationCommon.AwaitReceiptEth(context.Background(), client.EthClient(), tx.Hash(), 25*time.Second)
	if err != nil {
		return nil, fmt.Errorf("no receipt for MerkleTreeMessageBus contract state root manager addition")
	}

	// Create the Addresses struct to pass to initialize
	addresses := NetworkConfig.NetworkConfigFixedAddresses{
		CrossChain:               crossChainReceipt.ContractAddress,
		MessageBus:               messageBusAddr,
		NetworkEnclaveRegistry:   enclaveRegistryReceipt.ContractAddress,
		DataAvailabilityRegistry: daRegistryReceipt.ContractAddress,
	}

	_, networkConfigReceipt, err := deployNetworkConfigContract(client, wallets.ContractOwnerWallet, addresses)
	if err != nil {
		return nil, err
	}

	fmt.Println("Deployed All Network Contracts successfully",
		"networkConfigContract: ", networkConfigReceipt.ContractAddress,
		"dataAvailabilityRegistryContract: ", daRegistryReceipt.ContractAddress, "enclaveRegistryContract: ", enclaveRegistryReceipt.ContractAddress,
		"crossChainContract: ", crossChainReceipt.ContractAddress, "messageBusAddr: ", messageBusAddr)

	if !deployERC20s {
		return &params.L1TenData{
			TenStartBlock:                   networkConfigReceipt.BlockHash,
			NetworkConfigAddress:            networkConfigReceipt.ContractAddress,
			EnclaveRegistryAddress:          enclaveRegistryReceipt.ContractAddress,
			DataAvailabilityRegistryAddress: daRegistryReceipt.ContractAddress,
			CrossChainContractAddress:       crossChainReceipt.ContractAddress,
			MessageBusAddr:                  messageBusAddr,
		}, nil
	}

	erc20ContractAddr := make([]common.Address, 0)
	for _, token := range wallets.Tokens {
		erc20receipt, err := DeployContract(client, token.L1Owner, erc20contract.L1BytecodeWithDefaultSupply(string(token.Name), crossChainReceipt.ContractAddress))
		if err != nil {
			return nil, fmt.Errorf("failed to deploy ERC20 contract. Cause: %w", err)
		}
		token.L1ContractAddress = &erc20receipt.ContractAddress
		erc20ContractAddr = append(erc20ContractAddr, erc20receipt.ContractAddress)
	}

	return &params.L1TenData{
		TenStartBlock:                   networkConfigReceipt.BlockHash,
		NetworkConfigAddress:            networkConfigReceipt.ContractAddress,
		EnclaveRegistryAddress:          enclaveRegistryReceipt.ContractAddress,
		DataAvailabilityRegistryAddress: daRegistryReceipt.ContractAddress,
		CrossChainContractAddress:       crossChainReceipt.ContractAddress,
		ObxErc20Address:                 erc20ContractAddr[0],
		EthErc20Address:                 erc20ContractAddr[1],
		MessageBusAddr:                  messageBusAddr,
	}, nil
}

func deployEnclaveRegistryContract(client ethadapter.EthClient, contractOwner wallet.Wallet) (*NetworkEnclaveRegistry.NetworkEnclaveRegistry, *types.Receipt, error) {
	bytecode, err := constants.EnclaveRegistryBytecode()
	if err != nil {
		return nil, nil, err
	}
	networkEnclaveRegistryReceipt, err := DeployContract(client, contractOwner, bytecode)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to deploy NetworkEnclaveRegistry contract from %s. Cause: %w", contractOwner.Address(), err)
	}
	networkEnclaveRegistryContract, err := NetworkEnclaveRegistry.NewNetworkEnclaveRegistry(networkEnclaveRegistryReceipt.ContractAddress, client.EthClient())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to instantiate NetworkEnclaveRegistry contract. Cause: %w", err)
	}

	// Create a fresh transactor for initialization
	opts, err := createTransactor(contractOwner)
	if err != nil {
		return nil, nil, err
	}

	tx, err := networkEnclaveRegistryContract.Initialize(opts, contractOwner.Address())
	if err != nil {
		return nil, nil, fmt.Errorf("unable to initialize NetworkEnclaveRegistry contract. Cause: %w", err)
	}

	_, err = integrationCommon.AwaitReceiptEth(context.Background(), client.EthClient(), tx.Hash(), 25*time.Second)
	if err != nil {
		return nil, nil, fmt.Errorf("no receipt for NetworkEnclaveRegistry contract initialization")
	}
	return networkEnclaveRegistryContract, networkEnclaveRegistryReceipt, nil
}

func deployNetworkConfigContract(client ethadapter.EthClient, contractOwner wallet.Wallet, addresses NetworkConfig.NetworkConfigFixedAddresses) (*NetworkConfig.NetworkConfig, *types.Receipt, error) {
	// Deploy NetworkConfig contract
	bytecode, err := constants.NetworkConfigBytecode()
	if err != nil {
		return nil, nil, err
	}
	networkConfigReceipt, err := DeployContract(client, contractOwner, bytecode)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to deploy NetworkConfig contract from %s. Cause: %w", contractOwner.Address(), err)
	}
	networkConfigContract, err := NetworkConfig.NewNetworkConfig(networkConfigReceipt.ContractAddress, client.EthClient())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to instantiate NetworkConfig contract. Cause: %w", err)
	}

	opts, err := createTransactor(contractOwner)
	if err != nil {
		return nil, nil, err
	}
	tx, err := networkConfigContract.Initialize(opts, addresses, contractOwner.Address())
	if err != nil {
		return nil, nil, fmt.Errorf("unable to initialize NetworkConfig contract. Cause: %w", err)
	}

	_, err = integrationCommon.AwaitReceiptEth(context.Background(), client.EthClient(), tx.Hash(), 25*time.Second)
	if err != nil {
		return nil, nil, fmt.Errorf("no receipt for NetworkConfig contract initialization")
	}
	return networkConfigContract, networkConfigReceipt, nil
}

func deployDataAvailabilityRegistry(client ethadapter.EthClient, contractOwner wallet.Wallet, messageBus common.Address, enclaveRegistryAddress common.Address) (*DataAvailabilityRegistry.DataAvailabilityRegistry, *types.Receipt, error) {
	bytecode, err := constants.DataAvailabilityRegistryBytecode()
	if err != nil {
		return nil, nil, err
	}
	daRegistryContractReceipt, err := DeployContract(client, contractOwner, bytecode)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to deploy DataAvailabilityRegistry from %s. Cause: %w", contractOwner.Address(), err)
	}
	daRegistryContract, err := DataAvailabilityRegistry.NewDataAvailabilityRegistry(daRegistryContractReceipt.ContractAddress, client.EthClient())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to instantiate DataAvailabilityRegistry. Cause: %w", err)
	}

	opts, err := createTransactor(contractOwner)
	if err != nil {
		return nil, nil, err
	}
	tx, err := daRegistryContract.Initialize(opts, messageBus, enclaveRegistryAddress, contractOwner.Address())
	if err != nil {
		return nil, nil, fmt.Errorf("unable to initialize DataAvailabilityRegistry. Cause: %w", err)
	}

	_, err = integrationCommon.AwaitReceiptEth(context.Background(), client.EthClient(), tx.Hash(), 25*time.Second)
	if err != nil {
		return nil, nil, fmt.Errorf("no receipt for DataAvailabilityRegistry initialization")
	}

	err = waitForContractDeployment(context.Background(), client, daRegistryContractReceipt.ContractAddress)
	if err != nil {
		return nil, nil, fmt.Errorf("DataAvailabilityRegistry contract not available after time. Cause: %w", err)
	}

	return daRegistryContract, daRegistryContractReceipt, nil
}

func deployCrossChainContract(client ethadapter.EthClient, contractOwner wallet.Wallet) (*CrossChain.CrossChain, *types.Receipt, error) {
	// Deploy CrossChain contract
	bytecode, err := constants.CrossChainBytecode()
	if err != nil {
		return nil, nil, err
	}
	crossChainReceipt, err := DeployContract(client, contractOwner, bytecode)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to deploy CrossChain contract from %s. Cause: %w", contractOwner.Address(), err)
	}
	crossChainContract, err := CrossChain.NewCrossChain(crossChainReceipt.ContractAddress, client.EthClient())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to instantiate CrossChain contract. Cause: %w", err)
	}

	opts, err := createTransactor(contractOwner)
	if err != nil {
		return nil, nil, err
	}

	tx, err := crossChainContract.Initialize(opts, contractOwner.Address())
	if err != nil {
		return nil, nil, fmt.Errorf("unable to initialize CrossChain contract. Cause: %w", err)
	}

	_, err = integrationCommon.AwaitReceiptEth(context.Background(), client.EthClient(), tx.Hash(), 25*time.Second)
	if err != nil {
		return nil, nil, fmt.Errorf("no receipt for CrossChain initialization")
	}
	return crossChainContract, crossChainReceipt, nil
}

func PermissionTenSequencerEnclave(contractOwner wallet.Wallet, client ethadapter.EthClient, enclaveRegistryAddr common.Address, seqEnclaveID common.Address) error {
	ctr, err := NetworkEnclaveRegistry.NewNetworkEnclaveRegistry(enclaveRegistryAddr, client.EthClient())
	if err != nil {
		return err
	}

	opts, err := createTransactor(contractOwner)
	if err != nil {
		return err
	}

	tx, err := ctr.GrantSequencerEnclave(opts, seqEnclaveID)
	if err != nil {
		return fmt.Errorf("unable to grant enclave sequencer permission: %w", err)
	}

	_, err = integrationCommon.AwaitReceiptEth(context.Background(), client.EthClient(), tx.Hash(), 25*time.Second)
	if err != nil {
		return fmt.Errorf("unable to fetch receipt for granting enclave sequencer: %w", err)
	}

	return nil
}

func PermissionDataAvailabilityRegistryStateRoot(contractOwner wallet.Wallet, client ethadapter.EthClient, crossChainAddress common.Address, daRegistryAddress common.Address) error {
	crossChainContract, err := CrossChain.NewCrossChain(crossChainAddress, client.EthClient())
	if err != nil {
		return err
	}

	// fetch merkle message bus address from cross chain contract
	merkleTreeMessageBus, err := crossChainContract.MerkleMessageBus(&bind.CallOpts{})
	if err != nil {
		return fmt.Errorf("unable to get merkletreeMessageBus contract address: %w", err)
	}
	ctr, err := MerkleTreeMessageBus.NewMerkleTreeMessageBus(merkleTreeMessageBus, client.EthClient())
	if err != nil {
		return fmt.Errorf("unable to get merkletreeMessageBus contract: %w", err)
	}

	opts, err := createTransactor(contractOwner)
	if err != nil {
		return err
	}

	tx, err := ctr.AddStateRootManager(opts, daRegistryAddress)
	if err != nil {
		return fmt.Errorf("unable to grant rollup contract state root manager acces: %w", err)
	}

	_, err = integrationCommon.AwaitReceiptEth(context.Background(), client.EthClient(), tx.Hash(), 25*time.Second)
	if err != nil {
		return fmt.Errorf("unable to fetch receipt for granting rollup contract state root permission: %w", err)
	}

	return nil
}

func waitForContractDeployment(ctx context.Context, client ethadapter.EthClient, address common.Address) error {
	start := time.Now()
	for {
		if time.Since(start) > _contractTimeout {
			return fmt.Errorf("timeout waiting for contract code at address %s", address.Hex())
		}

		code, err := client.EthClient().CodeAt(ctx, address, nil)
		if err != nil {
			return fmt.Errorf("error checking contract code: %w", err)
		}

		if len(code) > 0 {
			return nil
		}

		time.Sleep(100 * time.Millisecond)
	}
}

func StopEth2Network(clients []ethadapter.EthClient, network eth2network.PosEth2Network) {
	// Stop the clients first
	for _, c := range clients {
		if c != nil {
			c.Stop()
		}
	}
	// Stop the nodes second
	if network != nil { // If network creation failed, we may be attempting to tear down the Geth network before it even exists.
		err := network.Stop()
		if err != nil {
			fmt.Println(err)
		}
	}
}

// DeployContract returns receipt of deployment
// todo (@matt) - this should live somewhere else
func DeployContract(workerClient ethadapter.EthClient, w wallet.Wallet, contractBytes []byte) (*types.Receipt, error) {
	deployContractTx, err := ethadapter.SetTxGasPrice(context.Background(), workerClient, &types.LegacyTx{Data: contractBytes}, w.Address(), w.GetNonceAndIncrement(), 0, nil, testlog.Logger())
	if err != nil {
		return nil, err
	}

	signedTx, err := w.SignTransaction(deployContractTx)
	if err != nil {
		return nil, err
	}

	err = workerClient.SendTransaction(signedTx)
	if err != nil {
		return nil, err
	}

	var start time.Time
	var receipt *types.Receipt

	time.Sleep(1 * time.Millisecond)

	// todo (@matt) these timings should be driven by the L2 batch times and L1 block times
	for start = time.Now(); time.Since(start) < 50*time.Second; time.Sleep(1 * time.Second) {
		receipt, err = workerClient.TransactionReceipt(signedTx.Hash())
		if err == nil && receipt != nil {
			if receipt.Status != types.ReceiptStatusSuccessful {
				return nil, errors.New("unable to deploy contract")
			}
			testlog.Logger().Info(fmt.Sprintf("Contract successfully deployed to %s", receipt.ContractAddress))
			return receipt, nil
		}

		testlog.Logger().Info(fmt.Sprintf("Contract deploy tx (%s) has not been mined into a block after %s...", signedTx.Hash(), time.Since(start)))
	}

	return nil, fmt.Errorf("failed to mine contract deploy tx (%s) into a block after %s. Aborting", signedTx.Hash(), time.Since(start))
}

func CreateEthClientConnection(port uint) ethadapter.EthClient {
	ethnode, err := ethadapter.NewEthClient(Localhost, port, DefaultL1RPCTimeout, testlog.Logger())
	if err != nil {
		panic(err)
	}
	return ethnode
}

// createTransactor creates a new transactor with the current nonce from the wallet
func createTransactor(wallet wallet.Wallet) (*bind.TransactOpts, error) {
	opts, err := bind.NewKeyedTransactorWithChainID(wallet.PrivateKey(), wallet.ChainID())
	if err != nil {
		return nil, err
	}
	opts.Nonce = big.NewInt(int64(wallet.GetNonceAndIncrement()))
	return opts, nil
}
