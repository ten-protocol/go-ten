package network

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ten-protocol/go-ten/contracts/generated/CrossChainMessenger"
	"github.com/ten-protocol/go-ten/contracts/generated/MerkleTreeMessageBus"
	"github.com/ten-protocol/go-ten/contracts/generated/TenBridge"
	"github.com/ten-protocol/go-ten/contracts/generated/TransparentUpgradeableProxy"

	"github.com/ten-protocol/go-ten/contracts/generated/NetworkConfig"
	"github.com/ten-protocol/go-ten/contracts/generated/NetworkEnclaveRegistry"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
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
	_, merkleMessageBusReceipt, err := deploMerkleMessageBusContract(client, wallets.ContractOwnerWallet)
	if err != nil {
		return nil, err
	}
	crossChainContract, crossChainReceipt, err := deployCrossChainContract(client, wallets.ContractOwnerWallet, merkleMessageBusReceipt.ContractAddress)
	if err != nil {
		return nil, err
	}

	_, daRegistryReceipt, err := deployDataAvailabilityRegistry(client, wallets.ContractOwnerWallet, merkleMessageBusReceipt.ContractAddress, enclaveRegistryReceipt.ContractAddress)
	if err != nil {
		return nil, err
	}

	// Create a new instance of MerkleTreeMessageBus bound to the proxy address
	merkleTreeMessageBus, err := MerkleTreeMessageBus.NewMerkleTreeMessageBus(merkleMessageBusReceipt.ContractAddress, client.EthClient())
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate MerkleTreeMessageBus contract. Cause: %w", err)
	}

	opts, err := createTransactor(wallets.ContractOwnerWallet)
	if err != nil {
		return nil, err
	}

	// add rollup contract as stateroot manager
	tx, err := merkleTreeMessageBus.AddStateRootManager(opts, daRegistryReceipt.ContractAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to add state root manager to MerkleTreeMessageBus contract. Cause: %w", err)
	}

	_, err = integrationCommon.AwaitReceiptEth(context.Background(), client.EthClient(), tx.Hash(), 25*time.Second)
	if err != nil {
		return nil, fmt.Errorf("no receipt for MerkleTreeMessageBus contract state root manager addition")
	}

	// add cross chain contract as withdrawal manager
	opts.Nonce = big.NewInt(int64(wallets.ContractOwnerWallet.GetNonceAndIncrement()))
	tx1, err := merkleTreeMessageBus.AddWithdrawalManager(opts, crossChainReceipt.ContractAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to add withdrawal manager to MerkleTreeMessageBus contract. Cause: %w", err)
	}

	_, err = integrationCommon.AwaitReceiptEth(context.Background(), client.EthClient(), tx1.Hash(), 25*time.Second)
	if err != nil {
		return nil, fmt.Errorf("no receipt for MerkleTreeMessageBus contract withdrawal manager addition")
	}

	opts.Nonce = big.NewInt(int64(wallets.ContractOwnerWallet.GetNonceAndIncrement()))
	ccAddress, tx, ccCtr, err := CrossChainMessenger.DeployCrossChainMessenger(opts, client.EthClient())
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate CrossChainMessenger contract. Cause: %w", err)
	}

	_, err = integrationCommon.AwaitReceiptEth(context.Background(), client.EthClient(), tx.Hash(), 25*time.Second)
	if err != nil {
		return nil, fmt.Errorf("no receipt for CrossChainMessenger contract deployment")
	}

	opts.Nonce = big.NewInt(int64(wallets.ContractOwnerWallet.GetNonceAndIncrement()))
	tx, err = ccCtr.Initialize(opts, merkleMessageBusReceipt.ContractAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize CrossChainMessenger contract. Cause: %w", err)
	}

	_, err = integrationCommon.AwaitReceiptEth(context.Background(), client.EthClient(), tx.Hash(), 25*time.Second)
	if err != nil {
		return nil, fmt.Errorf("no receipt for CrossChainMessenger contract initialization")
	}

	opts.Nonce = big.NewInt(int64(wallets.ContractOwnerWallet.GetNonceAndIncrement()))
	addr, tx, tenBridge, err := TenBridge.DeployTenBridge(opts, client.EthClient())
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate TenBridge contract. Cause: %w", err)
	}

	_, err = integrationCommon.AwaitReceiptEth(context.Background(), client.EthClient(), tx.Hash(), 25*time.Second)
	if err != nil {
		return nil, fmt.Errorf("no receipt for TenBridge contract deployment")
	}

	opts.Nonce = big.NewInt(int64(wallets.ContractOwnerWallet.GetNonceAndIncrement()))
	tx, err = tenBridge.Initialize(opts, ccAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize TenBridge contract. Cause: %w", err)
	}

	_, err = integrationCommon.AwaitReceiptEth(context.Background(), client.EthClient(), tx.Hash(), 25*time.Second)
	if err != nil {
		return nil, fmt.Errorf("no receipt for TenBridge contract initialization")
	}

	// Get the MessageBus address from the CrossChain contract
	messageBusAddr, err := crossChainContract.MessageBus(&bind.CallOpts{})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch MessageBus address. Cause: %w", err)
	}

	// Create the Addresses struct to pass to initialize
	addresses := NetworkConfig.NetworkConfigFixedAddresses{
		CrossChain:               crossChainReceipt.ContractAddress,
		MessageBus:               messageBusAddr,
		NetworkEnclaveRegistry:   enclaveRegistryReceipt.ContractAddress,
		DataAvailabilityRegistry: daRegistryReceipt.ContractAddress,
	}

	networkConfigContract, networkConfigReceipt, err := deployNetworkConfigContract(client, wallets.ContractOwnerWallet, addresses)
	if err != nil {
		return nil, err
	}

	opts, err = createTransactor(wallets.ContractOwnerWallet)
	if err != nil {
		return nil, err
	}

	tx, err = networkConfigContract.SetL1CrossChainMessengerAddress(opts, ccAddress)
	if err != nil {
		return nil, err
	}
	_, err = integrationCommon.AwaitReceiptEth(context.Background(), client.EthClient(), tx.Hash(), 25*time.Second)
	if err != nil {
		return nil, fmt.Errorf("no receipt for NetworkConfig contract additional address addition")
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
			CrossChainMessengerAddress:      ccAddress,
			BridgeAddress:                   addr,
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
		CrossChainMessengerAddress:      ccAddress,
		BridgeAddress:                   addr,
	}, nil
}

func ConnectTenNetworkBridge(client ethadapter.EthClient, wallets *params.SimWallets, l1Data *params.L1TenData, l2BridgeAddress common.Address) error {
	bridgeAddr := l1Data.BridgeAddress
	bridgeCtr, err := TenBridge.NewTenBridge(bridgeAddr, client.EthClient())
	if err != nil {
		return fmt.Errorf("failed to instantiate TenBridge contract. Cause: %w", err)
	}

	opts, err := createTransactor(wallets.ContractOwnerWallet)
	if err != nil {
		return fmt.Errorf("failed to create transactor. Cause: %w", err)
	}

	tx, err := bridgeCtr.SetRemoteBridge(opts, l2BridgeAddress)
	if err != nil {
		return fmt.Errorf("failed to set remote bridge. Cause: %w", err)
	}

	_, err = integrationCommon.AwaitReceiptEth(context.Background(), client.EthClient(), tx.Hash(), 25*time.Second)
	if err != nil {
		return fmt.Errorf("no receipt for TenBridge contract remote bridge setting")
	}

	return nil
}

// deployProxyContract handles the common pattern of deploying a contract with a proxy
func deployProxyContract[T any](
	client ethadapter.EthClient,
	contractOwner wallet.Wallet,
	bytecode []byte,
	contractMetadata interface{ GetAbi() (*abi.ABI, error) },
	initArgs ...interface{},
) (*T, *types.Receipt, error) {
	// deploy implementation
	implReceipt, err := DeployContract(client, contractOwner, bytecode)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to deploy implementation from %s. Cause: %w", contractOwner.Address(), err)
	}

	opts, err := createTransactor(contractOwner)
	if err != nil {
		return nil, nil, err
	}

	// encode the initialization data
	parsed, err := contractMetadata.GetAbi()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get contract ABI. Cause: %w", err)
	}

	initData, err := parsed.Pack("initialize", initArgs...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to encode initialization data. Cause: %w", err)
	}

	// deploy proxy
	proxyAddr, tx, _, err := TransparentUpgradeableProxy.DeployTransparentUpgradeableProxy(
		opts,
		client.EthClient(),
		implReceipt.ContractAddress,
		contractOwner.Address(),
		initData,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to deploy proxy contract. Cause: %w", err)
	}

	receipt, err := integrationCommon.AwaitReceiptEth(context.Background(), client.EthClient(), tx.Hash(), 25*time.Second)
	if err != nil {
		return nil, nil, fmt.Errorf("no receipt for proxy contract deployment")
	}

	// bind contract instance to proxy
	contract, err := NewContract[T](proxyAddr, client.EthClient())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to instantiate contract at proxy address. Cause: %w", err)
	}

	err = waitForContractDeployment(context.Background(), client, proxyAddr)
	if err != nil {
		return nil, nil, fmt.Errorf("proxy contract not available after time. Cause: %w", err)
	}

	return contract, receipt, nil
}

// NewContract is a type constraint helper to create new contract instances
func NewContract[T any](address common.Address, client *ethclient.Client) (*T, error) {
	switch any(new(T)).(type) {
	case *NetworkEnclaveRegistry.NetworkEnclaveRegistry:
		contract, err := NetworkEnclaveRegistry.NewNetworkEnclaveRegistry(address, client)
		return any(contract).(*T), err
	case *CrossChain.CrossChain:
		contract, err := CrossChain.NewCrossChain(address, client)
		return any(contract).(*T), err
	case *DataAvailabilityRegistry.DataAvailabilityRegistry:
		contract, err := DataAvailabilityRegistry.NewDataAvailabilityRegistry(address, client)
		return any(contract).(*T), err
	case *MerkleTreeMessageBus.MerkleTreeMessageBus:
		contract, err := MerkleTreeMessageBus.NewMerkleTreeMessageBus(address, client)
		return any(contract).(*T), err
	case *NetworkConfig.NetworkConfig:
		contract, err := NetworkConfig.NewNetworkConfig(address, client)
		return any(contract).(*T), err
	default:
		return nil, fmt.Errorf("unsupported contract type")
	}
}

func deployEnclaveRegistryContract(client ethadapter.EthClient, contractOwner wallet.Wallet) (*NetworkEnclaveRegistry.NetworkEnclaveRegistry, *types.Receipt, error) {
	bytecode, err := constants.EnclaveRegistryBytecode()
	if err != nil {
		return nil, nil, err
	}
	return deployProxyContract[NetworkEnclaveRegistry.NetworkEnclaveRegistry](
		client,
		contractOwner,
		bytecode,
		NetworkEnclaveRegistry.NetworkEnclaveRegistryMetaData,
		contractOwner.Address(),
	)
}

func deployNetworkConfigContract(client ethadapter.EthClient, contractOwner wallet.Wallet, addresses NetworkConfig.NetworkConfigFixedAddresses) (*NetworkConfig.NetworkConfig, *types.Receipt, error) {
	bytecode, err := constants.NetworkConfigBytecode()
	if err != nil {
		return nil, nil, err
	}
	return deployProxyContract[NetworkConfig.NetworkConfig](
		client,
		contractOwner,
		bytecode,
		NetworkConfig.NetworkConfigMetaData,
		addresses,
		contractOwner.Address(),
	)
}

func deployDataAvailabilityRegistry(client ethadapter.EthClient, contractOwner wallet.Wallet, messageBus common.Address, enclaveRegistryAddress common.Address) (*DataAvailabilityRegistry.DataAvailabilityRegistry, *types.Receipt, error) {
	bytecode, err := constants.DataAvailabilityRegistryBytecode()
	if err != nil {
		return nil, nil, err
	}
	return deployProxyContract[DataAvailabilityRegistry.DataAvailabilityRegistry](
		client,
		contractOwner,
		bytecode,
		DataAvailabilityRegistry.DataAvailabilityRegistryMetaData,
		messageBus,
		enclaveRegistryAddress,
		contractOwner.Address(),
	)
}

func deploMerkleMessageBusContract(client ethadapter.EthClient, contractOwner wallet.Wallet) (*MerkleTreeMessageBus.MerkleTreeMessageBus, *types.Receipt, error) {
	bytecode, err := constants.MerkleTreeMessageBusBytecode()
	if err != nil {
		return nil, nil, err
	}
	return deployProxyContract[MerkleTreeMessageBus.MerkleTreeMessageBus](
		client,
		contractOwner,
		bytecode,
		CrossChain.CrossChainMetaData,
		contractOwner.Address(),
		contractOwner.Address(),
	)
}

func deployCrossChainContract(client ethadapter.EthClient, contractOwner wallet.Wallet, messageBus common.Address) (*CrossChain.CrossChain, *types.Receipt, error) {
	bytecode, err := constants.CrossChainBytecode()
	if err != nil {
		return nil, nil, err
	}
	return deployProxyContract[CrossChain.CrossChain](
		client,
		contractOwner,
		bytecode,
		CrossChain.CrossChainMetaData,
		contractOwner.Address(),
		messageBus,
	)
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
