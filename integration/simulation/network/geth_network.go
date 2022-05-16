package network

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/contracts"
	"github.com/obscuronet/obscuro-playground/go/ethclient"
	"github.com/obscuronet/obscuro-playground/go/ethclient/mgmtcontractlib"
	"github.com/obscuronet/obscuro-playground/go/ethclient/stabletokencontractlib"
	"github.com/obscuronet/obscuro-playground/go/ethclient/txhandler"
	"github.com/obscuronet/obscuro-playground/go/ethclient/wallet"
	"github.com/obscuronet/obscuro-playground/go/log"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/host"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/obscuroclient"
	"github.com/obscuronet/obscuro-playground/integration/gethnetwork"
	"github.com/obscuronet/obscuro-playground/integration/simulation/p2p"
	"github.com/obscuronet/obscuro-playground/integration/simulation/params"
	"github.com/obscuronet/obscuro-playground/integration/simulation/stats"
)

type networkInMemGeth struct {
	obscuroClients []*obscuroclient.Client
	gethNetwork    *gethnetwork.GethNetwork
	wallets        []wallet.Wallet
	contracts      []string
	workerWallet   wallet.Wallet
}

func NewNetworkInMemoryGeth(wallets []wallet.Wallet, workerWallet wallet.Wallet, contracts []string) Network {
	return &networkInMemGeth{
		wallets:      wallets,
		workerWallet: workerWallet,
		contracts:    contracts,
	}
}

// Create inits and starts the nodes, wires them up, and populates the network objects
func (n *networkInMemGeth) Create(params *params.SimParams, stats *stats.Stats) ([]ethclient.EthClient, []*obscuroclient.Client, []string) {
	// make sure the geth network binaries exist
	path, err := gethnetwork.EnsureBinariesExist(gethnetwork.LatestVersion)
	if err != nil {
		panic(err)
	}

	// get wallet addresses to prefund them
	walletAddresses := make([]string, len(n.wallets))
	for i, w := range n.wallets {
		walletAddresses[i] = w.Address().String()
	}

	// kickoff the network with the prefunded wallet addresses
	n.gethNetwork = gethnetwork.NewGethNetwork(
		40000,
		path,
		params.NumberOfNodes,
		int(params.AvgBlockDuration.Seconds()),
		walletAddresses,
	)

	tmpEthClient, err := ethclient.NewEthClient(common.Address{}, "127.0.0.1", n.gethNetwork.WebSocketPorts[0], n.workerWallet, nil)
	if err != nil {
		panic(err)
	}

	contractAddr := deployContract(tmpEthClient, n.workerWallet, common.Hex2Bytes(contracts.MgmtContractByteCode))
	erc20ContractAddr := deployContract(tmpEthClient, n.workerWallet, common.Hex2Bytes(contracts.StableTokenERC20ContractByteCode))

	params.MgmtContractAddr = contractAddr
	params.ERC20ContractAddr = erc20ContractAddr
	params.TxHandler = txhandler.NewTransactionHandler(
		mgmtcontractlib.NewHandler(contractAddr),
		stabletokencontractlib.NewHandler(erc20ContractAddr),
	)
	// Create the obscuro node, each connected to a geth node
	l1Clients := make([]ethclient.EthClient, params.NumberOfNodes)
	obscuroNodes := make([]*host.Node, params.NumberOfNodes)
	n.obscuroClients = make([]*obscuroclient.Client, params.NumberOfNodes)

	for i := 0; i < params.NumberOfNodes; i++ {
		isGenesis := i == 0

		// create the in memory l1 and l2 node
		miner := createEthClientConnection(
			int64(i),
			n.gethNetwork.WebSocketPorts[i],
			params.EthWallets[i],
			params.MgmtContractAddr,
		)
		agg := createInMemObscuroNode(
			int64(i),
			isGenesis,
			params.TxHandler,
			params.AvgGossipPeriod,
			params.AvgBlockDuration,
			params.AvgNetworkLatency,
			stats,
			true,
			n.gethNetwork.GenesisJSON,
		)
		obscuroClient := host.NewInMemObscuroClient(int64(i), agg)

		// and connect them to each other
		agg.ConnectToEthNode(miner)

		obscuroNodes[i] = agg
		n.obscuroClients[i] = &obscuroClient
		l1Clients[i] = miner
	}

	// make sure the aggregators can talk to each other
	for i := 0; i < params.NumberOfNodes; i++ {
		obscuroNodes[i].P2p.(*p2p.MockP2P).Nodes = obscuroNodes
	}

	// start each obscuro node
	for _, m := range obscuroNodes {
		t := m
		go t.Start()
		time.Sleep(params.AvgBlockDuration / 10)
	}

	return l1Clients, n.obscuroClients, nil
}

func (n *networkInMemGeth) TearDown() {
	for _, client := range n.obscuroClients {
		temp := client
		go (*temp).Call(nil, obscuroclient.RPCStopHost) //nolint:errcheck
		go (*temp).Stop()
	}

	n.gethNetwork.StopNodes()
}

func createEthClientConnection(id int64, port uint, wallet wallet.Wallet, contractAddr *common.Address) ethclient.EthClient {
	ethnode, err := ethclient.NewEthClient(common.BigToAddress(big.NewInt(id)), "127.0.0.1", port, wallet, contractAddr)
	if err != nil {
		panic(err)
	}
	return ethnode
}

func deployContract(tmpClient ethclient.EthClient, w wallet.Wallet, contractBytes []byte) *common.Address {
	deployContractTx := types.LegacyTx{
		Nonce:    w.GetNonceAndIncrement(),
		GasPrice: big.NewInt(2000000000),
		Gas:      1025_000_000,
		Data:     contractBytes,
	}

	signedTx, err := tmpClient.SubmitTransaction(&deployContractTx)
	if err != nil {
		panic(err)
	}

	var receipt *types.Receipt
	for start := time.Now(); time.Since(start) < 80*time.Second; time.Sleep(2 * time.Second) {
		receipt, err = tmpClient.FetchTxReceipt(signedTx.Hash())
		if err == nil && receipt != nil {
			if receipt.Status != types.ReceiptStatusSuccessful {
				panic("unable to deploy contract")
			}
			break
		}

		log.Info("Contract deploy tx has not been mined into a block after %s...", time.Since(start))
	}

	log.Info("Contract successfully deployed to %s", receipt.ContractAddress)
	return &receipt.ContractAddress
}
