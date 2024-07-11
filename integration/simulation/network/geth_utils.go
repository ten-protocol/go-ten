package network

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ten-protocol/go-ten/contracts/generated/ManagementContract"
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

const (
	// These are the addresses that the end-to-end tests expect to be prefunded when run locally. Corresponds to
	// private key hex "f52e5418e349dccdda29b6ac8b0abe6576bb7713886aa85abea6181ba731f9bb".
	e2eTestPrefundedL1Addr = "0x13E23Ca74DE0206C56ebaE8D51b5622EFF1E9944"
)

func SetUpGethNetwork(wallets *params.SimWallets, startPort int, nrNodes int, blockDurationSeconds int) (*params.L1TenData, []ethadapter.EthClient, eth2network.Eth2Network) {
	eth2Network, err := StartGethNetwork(wallets, startPort, blockDurationSeconds)
	if err != nil {
		panic(fmt.Errorf("error starting geth network %w", err))
	}

	// connect to the first host to deploy
	tmpEthClient, err := ethadapter.NewEthClient(Localhost, uint(startPort+100), DefaultL1RPCTimeout, common.HexToAddress("0x0"), testlog.Logger())
	if err != nil {
		panic(fmt.Errorf("error connecting to te first host %w", err))
	}

	l1Data, err := DeployObscuroNetworkContracts(tmpEthClient, wallets, true)
	if err != nil {
		panic(fmt.Errorf("error deploying obscuro contract %w", err))
	}

	ethClients := make([]ethadapter.EthClient, nrNodes)
	for i := 0; i < nrNodes; i++ {
		ethClients[i] = CreateEthClientConnection(int64(i), uint(startPort+100))
	}

	return l1Data, ethClients, eth2Network
}

func StartGethNetwork(wallets *params.SimWallets, startPort int) (eth2network.PosEth2Network, error) {
	// make sure the geth network binaries exist
	binDir, err := eth2network.EnsureBinariesExist()
	if err != nil {
		return nil, err
	}

	network := eth2network.NewPosEth2Network(
		binDir,
		startPort+integration.DefaultGethAUTHPortOffset, // RPC
		startPort+integration.DefaultGethWSPortOffset,
		startPort+integration.DefaultGethNetworkPortOffset, // HTTP
		startPort+integration.DefaultPrysmP2PPortOffset,    // RPC
		3*time.Minute,
	)

	err = network.Start()
	if err != nil {
		return nil, err
	}

	return network, nil
}

func DeployObscuroNetworkContracts(client ethadapter.EthClient, wallets *params.SimWallets, deployERC20s bool) (*params.L1TenData, error) {
	bytecode, err := constants.Bytecode()
	if err != nil {
		return nil, err
	}
	mgmtContractReceipt, err := DeployContract(client, wallets.MCOwnerWallet, bytecode)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy management contract from %s. Cause: %w", wallets.MCOwnerWallet.Address(), err)
	}

	managementContract, err := ManagementContract.NewManagementContract(mgmtContractReceipt.ContractAddress, client.EthClient())
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate management contract. Cause: %w", err)
	}

	opts, err := bind.NewKeyedTransactorWithChainID(wallets.MCOwnerWallet.PrivateKey(), wallets.MCOwnerWallet.ChainID())
	if err != nil {
		return nil, fmt.Errorf("unable to create a keyed transactor for initializing the management contract. Cause: %w", err)
	}

	tx, err := managementContract.Initialize(opts)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize management contract. Cause: %w", err)
	}

	_, err = integrationCommon.AwaitReceiptEth(context.Background(), client.EthClient(), tx.Hash(), 25*time.Second)
	if err != nil {
		return nil, fmt.Errorf("no receipt for management contract initialization")
	}

	l1BusAddress, err := managementContract.MessageBus(&bind.CallOpts{})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch MessageBus address. Cause: %w", err)
	}

	fmt.Println("Deployed Management Contract successfully",
		"address: ", mgmtContractReceipt.ContractAddress, "txHash: ", mgmtContractReceipt.TxHash,
		"blockHash: ", mgmtContractReceipt.BlockHash, "l1BusAddress: ", l1BusAddress)

	if !deployERC20s {
		return &params.L1TenData{
			TenStartBlock:       mgmtContractReceipt.BlockHash,
			MgmtContractAddress: mgmtContractReceipt.ContractAddress,
			MessageBusAddr:      l1BusAddress,
		}, nil
	}

	erc20ContractAddr := make([]common.Address, 0)
	for _, token := range wallets.Tokens {
		erc20receipt, err := DeployContract(client, token.L1Owner, erc20contract.L1BytecodeWithDefaultSupply(string(token.Name), mgmtContractReceipt.ContractAddress))
		if err != nil {
			return nil, fmt.Errorf("failed to deploy ERC20 contract. Cause: %w", err)
		}
		token.L1ContractAddress = &erc20receipt.ContractAddress
		erc20ContractAddr = append(erc20ContractAddr, erc20receipt.ContractAddress)
	}

	return &params.L1TenData{
		TenStartBlock:       mgmtContractReceipt.BlockHash,
		MgmtContractAddress: mgmtContractReceipt.ContractAddress,
		ObxErc20Address:     erc20ContractAddr[0],
		EthErc20Address:     erc20ContractAddr[1],
		MessageBusAddr:      l1BusAddress,
	}, nil
}

func StopEth2Network(clients []ethadapter.EthClient, netw eth2network.Eth2Network) {
	// Stop the clients first
	for _, c := range clients {
		if c != nil {
			c.Stop()
		}
	}
	// Stop the nodes second
	if netw != nil { // If network creation failed, we may be attempting to tear down the Geth network before it even exists.
		err := netw.Stop()
		if err != nil {
			fmt.Println(err)
		}
	}
}

func InitializeContract(workerClient ethadapter.EthClient, w wallet.Wallet, contractAddress common.Address) (*types.Receipt, error) {
	ctr, err := ManagementContract.NewManagementContract(contractAddress, workerClient.EthClient())
	if err != nil {
		return nil, err
	}

	opts, err := bind.NewKeyedTransactorWithChainID(w.PrivateKey(), w.ChainID())
	if err != nil {
		return nil, err
	}

	tx, err := ctr.Initialize(opts)
	if err != nil {
		return nil, err
	}
	w.SetNonce(w.GetNonce())

	var start time.Time
	var receipt *types.Receipt
	// todo (@matt) these timings should be driven by the L2 batch times and L1 block times
	for start = time.Now(); time.Since(start) < 80*time.Second; time.Sleep(2 * time.Second) {
		receipt, err = workerClient.TransactionReceipt(tx.Hash())
		if err == nil && receipt != nil {
			if receipt.Status != types.ReceiptStatusSuccessful {
				return nil, errors.New("unable to initialize contract")
			}
			testlog.Logger().Info("Contract initialized")
			return receipt, nil
		}
	}

	return receipt, nil
}

// DeployContract returns receipt of deployment
// todo (@matt) - this should live somewhere else
func DeployContract(workerClient ethadapter.EthClient, w wallet.Wallet, contractBytes []byte) (*types.Receipt, error) {
	deployContractTx, err := workerClient.PrepareTransactionToSend(
		context.Background(),
		&types.LegacyTx{Data: contractBytes},
		w.Address(),
	)
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

func CreateEthClientConnection(id int64, port uint) ethadapter.EthClient {
	ethnode, err := ethadapter.NewEthClient(Localhost, port, DefaultL1RPCTimeout, common.BigToAddress(big.NewInt(id)), testlog.Logger())
	if err != nil {
		panic(err)
	}
	return ethnode
}
