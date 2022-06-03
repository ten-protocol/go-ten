package smartcontract

import (
	"bytes"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/go/ethclient"
	"github.com/obscuronet/obscuro-playground/go/ethclient/mgmtcontractlib"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/config"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/wallet"
	"github.com/obscuronet/obscuro-playground/integration"
	"github.com/obscuronet/obscuro-playground/integration/datagenerator"
	"github.com/obscuronet/obscuro-playground/integration/gethnetwork"
	"github.com/obscuronet/obscuro-playground/integration/simulation/network"
)

// netInfo is a bag holder struct for output data from the execution/run of a network
type netInfo struct {
	ethClients  []ethclient.EthClient
	wallets     []wallet.Wallet
	gethNetwork *gethnetwork.GethNetwork
}

// runGethNetwork runs a geth network with one prefunded wallet
func runGethNetwork(t *testing.T) *netInfo {
	// make sure the geth network binaries exist
	path, err := gethnetwork.EnsureBinariesExist(gethnetwork.LatestVersion)
	if err != nil {
		t.Fatal(err)
	}

	// prefund one wallet as the worker wallet
	workerWallet := datagenerator.RandomWallet(integration.EthereumChainID)

	// define + run the network
	gethNetwork := gethnetwork.NewGethNetwork(
		integration.StartPortSmartContractTests,
		integration.StartPortSmartContractTests+100,
		path,
		3,
		1,
		[]string{workerWallet.Address().String()},
	)

	// create a client that is connected to node 0 of the network
	client, err := ethclient.NewEthClient(config.HostConfig{
		ID:                  common.Address{1},
		L1NodeHost:          "127.0.0.1",
		L1NodeWebsocketPort: gethNetwork.WebSocketPorts[0],
		L1ConnectionTimeout: 30 * time.Second,
	})
	if err != nil {
		return nil
	}

	return &netInfo{
		ethClients:  []ethclient.EthClient{client},
		wallets:     []wallet.Wallet{workerWallet},
		gethNetwork: gethNetwork,
	}
}

func TestManagementContract(t *testing.T) {
	// run tests on one network
	sim := runGethNetwork(t)
	defer sim.gethNetwork.StopNodes()

	// setup the client and the (debug) wallet
	client := sim.ethClients[0]
	w := newDebugWallet(sim.wallets[0])

	for name, test := range map[string]func(*testing.T, *debugMgmtContractLib, *debugWallet, ethclient.EthClient){
		"nonAttestedNodesCannotCreateRollup": nonAttestedNodesCannotCreateRollup,
		"attestedNodesCreateRollup":          attestedNodesCreateRollup,
	} {
		t.Run(name, func(t *testing.T) {
			// deploy the same contract to a new address
			contractAddr, err := network.DeployContract(client, w, common.Hex2Bytes(mgmtcontractlib.MgmtContractByteCode))
			if err != nil {
				t.Error(err)
			}

			// run the test using the new contract, but same wallet
			test(t,
				newDebugMgmtContractLib(*contractAddr, client.EthClient(), mgmtcontractlib.NewMgmtContractLib(contractAddr)),
				w,
				client,
			)
		})
	}
}

// nonAttestedNodesCannotCreateRollup issues a rollup from a node that did not receive the secret network key
func nonAttestedNodesCannotCreateRollup(t *testing.T, mgmtContractLib *debugMgmtContractLib, w *debugWallet, client ethclient.EthClient) {
	rollup := datagenerator.RandomRollup()
	txData := mgmtContractLib.CreateRollup(
		&obscurocommon.L1RollupTx{Rollup: nodecommon.EncodeRollup(&rollup)},
		w.GetNonceAndIncrement(),
	)

	_, receipt, err := w.AwaitedSignAndSendTransaction(client, txData)
	if err != nil {
		t.Error(err)
	}

	if receipt.Status != types.ReceiptStatusFailed {
		t.Errorf("transaction should have failed, expected %d got %d", 0, receipt.Status)
	}
}

// attestedNodesCreateRollup attests a node by issuing a CreateRespondSecret, issues a rollups from the same node and verifies the rollup was stored
func attestedNodesCreateRollup(t *testing.T, mgmtContractLib *debugMgmtContractLib, w *debugWallet, client ethclient.EthClient) {
	rollup := datagenerator.RandomRollup()
	requesterID := rollup.Header.Agg
	requesterPubKey := common.HexToHash("0x1337")

	// mark the aggregator as attested
	txData := mgmtContractLib.CreateRespondSecret(
		&obscurocommon.L1RespondSecretTx{
			RequesterPubKey: requesterPubKey.Bytes(),
			RequesterID:     requesterID,
		},
		w.GetNonceAndIncrement(),
	)

	_, receipt, err := w.AwaitedSignAndSendTransaction(client, txData)
	if err != nil {
		t.Error(err)
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		t.Errorf("transaction should have sucessed, expected %d got %d", 1, receipt.Status)
	}

	// issue a rollup from the attested node
	txData = mgmtContractLib.CreateRollup(&obscurocommon.L1RollupTx{Rollup: nodecommon.EncodeRollup(&rollup)}, w.GetNonceAndIncrement())
	_, receipt, err = w.AwaitedSignAndSendTransaction(client, txData)
	if err != nil {
		t.Error(err)
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		t.Errorf("transaction should have sucessed, expected %d got %d", 1, receipt.Status)
	}

	// make sure the rollup was stored in the contract
	storedRollup, err := mgmtContractLib.genContract.Rollups(nil, receipt.BlockNumber, big.NewInt(0))
	if err != nil {
		t.Error(err)
	}

	if storedRollup.Number.Int64() != int64(rollup.Header.Number) ||
		!bytes.Equal(storedRollup.ParentHash[:], rollup.Header.ParentHash.Bytes()) ||
		!bytes.Equal(storedRollup.AggregatorID[:], rollup.Header.Agg.Bytes()) ||
		!bytes.Equal(storedRollup.L1Block[:], rollup.Header.L1Proof.Bytes()) {
		t.Error("stored rollup does not match the generated rollup")
	}
}
