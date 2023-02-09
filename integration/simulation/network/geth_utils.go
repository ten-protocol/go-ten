package network

import (
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/go-obscuro/contracts/generated/ManagementContract"
	"github.com/obscuronet/go-obscuro/go/common/constants"
	"github.com/obscuronet/go-obscuro/go/ethadapter"
	"github.com/obscuronet/go-obscuro/go/wallet"
	"github.com/obscuronet/go-obscuro/integration"
	"github.com/obscuronet/go-obscuro/integration/common/testlog"
	"github.com/obscuronet/go-obscuro/integration/erc20contract"
	"github.com/obscuronet/go-obscuro/integration/eth2network"
	"github.com/obscuronet/go-obscuro/integration/simulation/params"
)

const (
	// These are the addresses that the end-to-end tests expect to be prefunded when run locally. Corresponds to
	// private key hex "f52e5418e349dccdda29b6ac8b0abe6576bb7713886aa85abea6181ba731f9bb".
	e2eTestPrefundedL1Addr = "0x13E23Ca74DE0206C56ebaE8D51b5622EFF1E9944"
	// TODO - Also prefund the L1 HOC and POC addresses used for the end-to-end tests when run locally.
)

func SetUpGethNetwork(wallets *params.SimWallets, startPort int, nrNodes int, blockDurationSeconds int) (*params.L1SetupData, []ethadapter.EthClient, eth2network.Eth2Network) {
	// make sure the geth network binaries exist
	path, err := eth2network.EnsureBinariesExist()
	if err != nil {
		panic(err)
	}

	// get the node wallet addresses to prefund them with Eth, so they can submit rollups, deploy contracts, deposit to the bridge, etc
	walletAddresses := []string{e2eTestPrefundedL1Addr}
	for _, w := range wallets.AllEthWallets() {
		walletAddresses = append(walletAddresses, w.Address().String())
	}

	// kickoff the network with the prefunded wallet addresses
	eth2Network := eth2network.NewEth2Network(
		path,
		startPort,
		startPort+integration.DefaultGethWSPortOffset,
		startPort+integration.DefaultGethAUTHPortOffset,
		startPort+integration.DefaultGethNetworkPortOffset,
		startPort+integration.DefaultPrysmHTTPPortOffset,
		startPort+integration.DefaultPrysmP2PPortOffset,
		1337,
		1,
		blockDurationSeconds,
		walletAddresses,
	)

	err = eth2Network.Start()
	if err != nil {
		panic(err)
	}
	// connect to the first host to deploy
	tmpEthClient, err := ethadapter.NewEthClient(Localhost, uint(startPort+100), DefaultL1RPCTimeout, common.HexToAddress("0x0"), testlog.Logger())
	if err != nil {
		panic(err)
	}

	bytecode, err := constants.Bytecode()
	if err != nil {
		panic(err)
	}
	mgmtContractReceipt, err := DeployContract(tmpEthClient, wallets.MCOwnerWallet, bytecode)
	if err != nil {
		panic(fmt.Sprintf("failed to deploy management contract. Cause: %s", err))
	}

	managementContract, _ := ManagementContract.NewManagementContract(mgmtContractReceipt.ContractAddress, tmpEthClient.EthClient())
	l1BusAddress, _ := managementContract.MessageBus(&bind.CallOpts{})

	erc20ContractAddr := make([]common.Address, 0)
	for _, token := range wallets.Tokens {
		erc20receipt, err := DeployContract(tmpEthClient, token.L1Owner, erc20contract.L1BytecodeWithDefaultSupply(string(token.Name), mgmtContractReceipt.ContractAddress))
		if err != nil {
			panic(fmt.Sprintf("failed to deploy ERC20 contract. Cause: %s", err))
		}
		token.L1ContractAddress = &erc20receipt.ContractAddress
		erc20ContractAddr = append(erc20ContractAddr, erc20receipt.ContractAddress)
	}

	ethClients := make([]ethadapter.EthClient, nrNodes)
	for i := 0; i < nrNodes; i++ {
		ethClients[i] = createEthClientConnection(int64(i), uint(startPort+100))
	}

	l1Data := &params.L1SetupData{
		ObscuroStartBlock:   mgmtContractReceipt.BlockHash,
		MgmtContractAddress: mgmtContractReceipt.ContractAddress,
		ObxErc20Address:     erc20ContractAddr[0],
		EthErc20Address:     erc20ContractAddr[1],
		MessageBusAddr:      &l1BusAddress,
	}

	return l1Data, ethClients, eth2Network
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

// DeployContract returns receipt of deployment
// todo -this should live somewhere else
func DeployContract(workerClient ethadapter.EthClient, w wallet.Wallet, contractBytes []byte) (*types.Receipt, error) {
	deployContractTx, err := workerClient.EstimateGasAndGasPrice(&types.LegacyTx{
		Nonce: w.GetNonceAndIncrement(),
		Data:  contractBytes,
	}, w.Address())
	if err != nil {
		w.SetNonce(w.GetNonce() - 1)
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
	for start = time.Now(); time.Since(start) < 80*time.Second; time.Sleep(2 * time.Second) {
		receipt, err = workerClient.TransactionReceipt(signedTx.Hash())
		if err == nil && receipt != nil {
			if receipt.Status != types.ReceiptStatusSuccessful {
				return nil, errors.New("unable to deploy contract")
			}
			testlog.Logger().Info(fmt.Sprintf("Contract successfully deployed to %s", receipt.ContractAddress))
			return receipt, nil
		}

		testlog.Logger().Info(fmt.Sprintf("Contract deploy tx has not been mined into a block after %s...", time.Since(start)))
	}

	return nil, fmt.Errorf("failed to mine contract deploy tx into a block after %s. Aborting", time.Since(start))
}

func createEthClientConnection(id int64, port uint) ethadapter.EthClient {
	ethnode, err := ethadapter.NewEthClient(Localhost, port, DefaultL1RPCTimeout, common.BigToAddress(big.NewInt(id)), testlog.Logger())
	if err != nil {
		panic(err)
	}
	return ethnode
}
