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
	"github.com/ten-protocol/go-ten/contracts/generated/RollupContract"
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

	_, rollupReceipt, err := deployRollupContract(client, wallets.ContractOwnerWallet, messageBusAddr, enclaveRegistryReceipt.ContractAddress)
	if err != nil {
		return nil, err
	}

	// Create the Addresses struct to pass to initialize
	addresses := NetworkConfig.NetworkConfigAddresses{
		CrossChain:             crossChainReceipt.ContractAddress,
		MessageBus:             messageBusAddr,
		NetworkEnclaveRegistry: enclaveRegistryReceipt.ContractAddress,
		RollupContract:         rollupReceipt.ContractAddress,
	}

	_, networkConfigReceipt, err := deployNetworkConfigContract(client, wallets.ContractOwnerWallet, addresses)
	if err != nil {
		return nil, err
	}

	fmt.Println("Deployed All Network Contracts successfully",
		"networkConfigContract: ", networkConfigReceipt.ContractAddress,
		"rollupContract: ", rollupReceipt.ContractAddress, "enclaveRegistryContract: ", enclaveRegistryReceipt.ContractAddress,
		"crossChainContract: ", crossChainReceipt.ContractAddress, "messageBusAddr: ", messageBusAddr)

	if !deployERC20s {
		return &params.L1TenData{
			TenStartBlock:             networkConfigReceipt.BlockHash,
			NetworkConfigAddress:      networkConfigReceipt.ContractAddress,
			EnclaveRegistryAddress:    enclaveRegistryReceipt.ContractAddress,
			RollupContractAddress:     rollupReceipt.ContractAddress,
			CrossChainContractAddress: crossChainReceipt.ContractAddress,
			MessageBusAddr:            messageBusAddr,
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
		TenStartBlock:             networkConfigReceipt.BlockHash,
		NetworkConfigAddress:      networkConfigReceipt.ContractAddress,
		EnclaveRegistryAddress:    enclaveRegistryReceipt.ContractAddress,
		RollupContractAddress:     rollupReceipt.ContractAddress,
		CrossChainContractAddress: crossChainReceipt.ContractAddress,
		ObxErc20Address:           erc20ContractAddr[0],
		EthErc20Address:           erc20ContractAddr[1],
		MessageBusAddr:            messageBusAddr,
	}, nil
}

func deployEnclaveRegistryContract(client ethadapter.EthClient, ownerKey wallet.Wallet) (*NetworkEnclaveRegistry.NetworkEnclaveRegistry, *types.Receipt, error) {
	bytecode, err := constants.EnclaveRegistryBytecode()
	if err != nil {
		return nil, nil, err
	}
	networkEnclaveRegistryReceipt, err := DeployContract(client, ownerKey, bytecode)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to deploy NetworkEnclaveRegistry contract from %s. Cause: %w", ownerKey.Address(), err)
	}
	networkEnclaveRegistryContract, err := NetworkEnclaveRegistry.NewNetworkEnclaveRegistry(networkEnclaveRegistryReceipt.ContractAddress, client.EthClient())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to instantiate NetworkEnclaveRegistry contract. Cause: %w", err)
	}

	// Create a fresh transactor for initialization
	opts, err := createTransactor(ownerKey)
	if err != nil {
		return nil, nil, err
	}

	tx, err := networkEnclaveRegistryContract.Initialize(opts)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to initialize NetworkEnclaveRegistry contract. Cause: %w", err)
	}

	_, err = integrationCommon.AwaitReceiptEth(context.Background(), client.EthClient(), tx.Hash(), 25*time.Second)
	if err != nil {
		return nil, nil, fmt.Errorf("no receipt for NetworkEnclaveRegistry contract initialization")
	}
	return networkEnclaveRegistryContract, networkEnclaveRegistryReceipt, nil
}

func deployNetworkConfigContract(client ethadapter.EthClient, ownerKey wallet.Wallet, addresses NetworkConfig.NetworkConfigAddresses) (*NetworkConfig.NetworkConfig, *types.Receipt, error) {
	// Deploy NetworkConfig contract
	bytecode, err := constants.NetworkConfigBytecode()
	if err != nil {
		return nil, nil, err
	}
	networkConfigReceipt, err := DeployContract(client, ownerKey, bytecode)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to deploy NetworkConfig contract from %s. Cause: %w", ownerKey.Address(), err)
	}
	networkConfigContract, err := NetworkConfig.NewNetworkConfig(networkConfigReceipt.ContractAddress, client.EthClient())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to instantiate NetworkConfig contract. Cause: %w", err)
	}

	opts, err := createTransactor(ownerKey)
	if err != nil {
		return nil, nil, err
	}
	tx, err := networkConfigContract.Initialize(opts, addresses)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to initialize NetworkConfig contract. Cause: %w", err)
	}

	_, err = integrationCommon.AwaitReceiptEth(context.Background(), client.EthClient(), tx.Hash(), 25*time.Second)
	if err != nil {
		return nil, nil, fmt.Errorf("no receipt for NetworkConfig contract initialization")
	}
	return networkConfigContract, networkConfigReceipt, nil
}

func deployRollupContract(client ethadapter.EthClient, ownerKey wallet.Wallet, messageBus common.Address, enclaveRegistryAddress common.Address) (*RollupContract.RollupContract, *types.Receipt, error) {
	bytecode, err := constants.RollupContractBytecode()
	if err != nil {
		return nil, nil, err
	}
	rollupContractReceipt, err := DeployContract(client, ownerKey, bytecode)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to deploy RollupContract from %s. Cause: %w", ownerKey.Address(), err)
	}
	rollupContract, err := RollupContract.NewRollupContract(rollupContractReceipt.ContractAddress, client.EthClient())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to instantiate RollupContract. Cause: %w", err)
	}

	opts, err := createTransactor(ownerKey)
	if err != nil {
		return nil, nil, err
	}
	tx, err := rollupContract.Initialize(opts, messageBus, enclaveRegistryAddress)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to initialize RollupContract. Cause: %w", err)
	}

	_, err = integrationCommon.AwaitReceiptEth(context.Background(), client.EthClient(), tx.Hash(), 25*time.Second)
	if err != nil {
		return nil, nil, fmt.Errorf("no receipt for RollupContract initialization")
	}
	return rollupContract, rollupContractReceipt, nil
}

func deployCrossChainContract(client ethadapter.EthClient, ownerKey wallet.Wallet) (*CrossChain.CrossChain, *types.Receipt, error) {
	// Deploy CrossChain contract
	bytecode, err := constants.CrossChainBytecode()
	if err != nil {
		return nil, nil, err
	}
	crossChainReceipt, err := DeployContract(client, ownerKey, bytecode)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to deploy CrossChain contract from %s. Cause: %w", ownerKey.Address(), err)
	}
	crossChainContract, err := CrossChain.NewCrossChain(crossChainReceipt.ContractAddress, client.EthClient())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to instantiate CrossChain contract. Cause: %w", err)
	}

	opts, err := createTransactor(ownerKey)
	if err != nil {
		return nil, nil, err
	}

	tx, err := crossChainContract.Initialize(opts)
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

func PermissionRollupContractStateRoot(contractOwner wallet.Wallet, client ethadapter.EthClient, crossChainAddress common.Address, rollupContractAddress common.Address) error {
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

	tx, err := ctr.AddStateRootManager(opts, rollupContractAddress)
	if err != nil {
		return fmt.Errorf("unable to grant rollup contract state root manager acces: %w", err)
	}

	_, err = integrationCommon.AwaitReceiptEth(context.Background(), client.EthClient(), tx.Hash(), 25*time.Second)
	if err != nil {
		return fmt.Errorf("unable to fetch receipt for granting rollup contract state root permission: %w", err)
	}

	return nil
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

func InitializeContract(workerClient ethadapter.EthClient, w wallet.Wallet, contractAddress common.Address) (*types.Receipt, error) {
	//ctr, err := ManagementContract.NewManagementContract(contractAddress, workerClient.EthClient())
	//if err != nil {
	//	return nil, err
	//}
	//
	//opts, err := bind.NewKeyedTransactorWithChainID(w.PrivateKey(), w.ChainID())
	//if err != nil {
	//	return nil, err
	//}
	//
	//tx, err := ctr.Initialize(opts)
	//if err != nil {
	//	return nil, err
	//}
	//w.SetNonce(w.GetNonce())
	//
	//var start time.Time
	//var receipt *types.Receipt
	//// todo (@matt) these timings should be driven by the L2 batch times and L1 block times
	//for start = time.Now(); time.Since(start) < 80*time.Second; time.Sleep(2 * time.Second) {
	//	receipt, err = workerClient.TransactionReceipt(tx.Hash())
	//	if err == nil && receipt != nil {
	//		if receipt.Status != types.ReceiptStatusSuccessful {
	//			return nil, errors.New("unable to initialize contract")
	//		}
	//		testlog.Logger().Info("Contract initialized")
	//		return receipt, nil
	//	}
	//}
	//
	return nil, nil
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
