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

type simInfo struct {
	ethClients  []ethclient.EthClient
	wallets     []wallet.Wallet
	gethNetwork *gethnetwork.GethNetwork
}

func runNetwork(t *testing.T) *simInfo {
	// make sure the geth network binaries exist
	path, err := gethnetwork.EnsureBinariesExist(gethnetwork.LatestVersion)
	if err != nil {
		t.Fatal(err)
	}

	// prefund one wallet as the worker wallet
	workerWallet := datagenerator.RandomWallet(integration.EthereumChainID)

	gethNetwork := gethnetwork.NewGethNetwork(
		integration.StartPortSmartContractTests,
		integration.StartPortSmartContractTests+100,
		path,
		3,
		1,
		[]string{workerWallet.Address().String()},
	)

	client, err := ethclient.NewEthClient(config.HostConfig{
		ID:                  common.Address{1},
		L1NodeHost:          "127.0.0.1",
		L1NodeWebsocketPort: gethNetwork.WebSocketPorts[0],
		L1ConnectionTimeout: 30 * time.Second,
	})
	if err != nil {
		return nil
	}

	return &simInfo{
		ethClients:  []ethclient.EthClient{client},
		wallets:     []wallet.Wallet{workerWallet},
		gethNetwork: gethNetwork,
	}
}

func TestManagementContract(t *testing.T) {
	sim := runNetwork(t)
	defer sim.gethNetwork.StopNodes()

	for name, test := range map[string]func(*testing.T, *debugMgmtContractLib, *debugWallet, ethclient.EthClient){
		"nonAttestedNodesCannotCreateRollup": nonAttestedNodesCannotCreateRollup,
		"attestedNodesCreateRollup":          attestedNodesCreateRollup,
	} {
		t.Run(name, func(t *testing.T) {
			client := sim.ethClients[0]
			w := sim.wallets[0]

			contractAddr, err := network.DeployContract(client, w, common.Hex2Bytes(mgmtcontractlib.MgmtContractByteCode))
			if err != nil {
				t.Error(err)
			}

			test(t,
				newDebugMgmtContractLib(*contractAddr, client.EthClient(), mgmtcontractlib.NewMgmtContractLib(contractAddr)),
				newDebugWallet(w),
				client,
			)
		})
	}
}

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

func attestedNodesCreateRollup(t *testing.T, mgmtContractLib *debugMgmtContractLib, w *debugWallet, client ethclient.EthClient) {
	rollup := datagenerator.RandomRollup()
	requesterID := rollup.Header.Agg
	requesterPubKey := common.HexToHash("0x1337")
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

	txData = mgmtContractLib.CreateRollup(&obscurocommon.L1RollupTx{Rollup: nodecommon.EncodeRollup(&rollup)}, w.GetNonceAndIncrement())

	_, receipt, err = w.AwaitedSignAndSendTransaction(client, txData)
	if err != nil {
		t.Error(err)
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		t.Errorf("transaction should have sucessed, expected %d got %d", 1, receipt.Status)
	}

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
