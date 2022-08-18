package network

import (
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/obscuronet/go-obscuro/contracts/managementcontract"
	"github.com/obscuronet/go-obscuro/integration/erc20contract"

	"github.com/obscuronet/go-obscuro/go/common/log"

	"github.com/obscuronet/go-obscuro/integration/simulation/params"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/go-obscuro/go/ethadapter"
	"github.com/obscuronet/go-obscuro/go/wallet"
	"github.com/obscuronet/go-obscuro/integration/gethnetwork"
)

func SetUpGethNetwork(wallets *params.SimWallets, StartPort int, nrNodes int, blockDurationSeconds int) (*common.Address, *common.Address, *common.Address, []ethadapter.EthClient, *gethnetwork.GethNetwork) {
	// make sure the geth network binaries exist
	path, err := gethnetwork.EnsureBinariesExist(gethnetwork.LatestVersion)
	if err != nil {
		panic(err)
	}

	// get the node wallet addresses to prefund them with Eth, so they can submit rollups, deploy contracts, deposit to the bridge, etc
	walletAddresses := make([]string, len(wallets.AllEthWallets()))
	for i, w := range wallets.AllEthWallets() {
		walletAddresses[i] = w.Address().String()
	}

	// kickoff the network with the prefunded wallet addresses
	gethNetwork := gethnetwork.NewGethNetwork(
		StartPort,
		StartPort+DefaultWsPortOffset,
		path,
		nrNodes,
		blockDurationSeconds,
		walletAddresses,
	)

	// connect to the first host to deploy
	tmpEthClient, err := ethadapter.NewEthClient(Localhost, gethNetwork.WebSocketPorts[0], DefaultL1ConnectionTimeout, common.HexToAddress("0x0"))
	if err != nil {
		panic(err)
	}

	bytecode, err := managementcontract.Bytecode()
	if err != nil {
		panic(err)
	}
	mgmtContractAddr, err := DeployContract(tmpEthClient, wallets.MCOwnerWallet, bytecode)
	if err != nil {
		panic(fmt.Sprintf("failed to deploy management contract. Cause: %s", err))
	}

	erc20ContractAddr := make([]*common.Address, 0)
	for _, token := range wallets.Tokens {
		address, err := DeployContract(tmpEthClient, token.L1Owner, erc20contract.L1BytecodeWithDefaultSupply(string(token.Name)))
		if err != nil {
			panic(fmt.Sprintf("failed to deploy ERC20 contract. Cause: %s", err))
		}
		token.L1ContractAddress = address
		erc20ContractAddr = append(erc20ContractAddr, address)
	}

	ethClients := make([]ethadapter.EthClient, nrNodes)
	for i := 0; i < nrNodes; i++ {
		ethClients[i] = createEthClientConnection(int64(i), gethNetwork.WebSocketPorts[i])
	}

	return mgmtContractAddr, erc20ContractAddr[0], erc20ContractAddr[1], ethClients, gethNetwork
}

func StopGethNetwork(clients []ethadapter.EthClient, netw *gethnetwork.GethNetwork) {
	// Stop the clients first
	for _, c := range clients {
		if c != nil {
			c.Stop()
		}
	}
	// Stop the nodes second
	if netw != nil { // If network creation failed, we may be attempting to tear down the Geth network before it even exists.
		netw.StopNodes()
	}
}

// DeployContract todo -this should live somewhere else
func DeployContract(workerClient ethadapter.EthClient, w wallet.Wallet, contractBytes []byte) (*common.Address, error) {
	deployContractTx := types.LegacyTx{
		Nonce:    w.GetNonceAndIncrement(),
		GasPrice: big.NewInt(2000000000),
		Gas:      1025_000_000,
		Data:     contractBytes,
	}

	signedTx, err := w.SignTransaction(&deployContractTx)
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
			log.Info("Contract successfully deployed to %s", receipt.ContractAddress)
			return &receipt.ContractAddress, nil
		}

		log.Info("Contract deploy tx has not been mined into a block after %s...", time.Since(start))
	}

	return nil, fmt.Errorf("failed to mine contract deploy tx into a block after %s. Aborting", time.Since(start))
}

func createEthClientConnection(id int64, port uint) ethadapter.EthClient {
	ethnode, err := ethadapter.NewEthClient(Localhost, port, DefaultL1ConnectionTimeout, common.BigToAddress(big.NewInt(id)))
	if err != nil {
		panic(err)
	}
	return ethnode
}
