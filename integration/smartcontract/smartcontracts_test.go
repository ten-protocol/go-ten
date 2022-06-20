package smartcontract

import (
	"bytes"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"

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
		//"secretCannotBeInitializedTwice":     secretCannotBeInitializedTwice,
		//"nonAttestedNodesCannotCreateRollup": nonAttestedNodesCannotCreateRollup,
		//"attestedNodesCreateRollup":          attestedNodesCreateRollup,
		//"nonAttestedNodesCannotAttest":       nonAttestedNodesCannotAttest,
		//"newlyAttestedNodesCanAttest":        newlyAttestedNodesCanAttest,
		"detectSimpleFork": detectSimpleFork,
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

// secretCannotBeInitializedTwice issues the InitializeNetworkSecret twice, failing the second time
func secretCannotBeInitializedTwice(t *testing.T, mgmtContractLib *debugMgmtContractLib, w *debugWallet, client ethclient.EthClient) {
	aggregatorID := datagenerator.RandomAddress()
	txData := mgmtContractLib.CreateInitializeSecret(
		&obscurocommon.L1InitializeSecretTx{
			AggregatorID: &aggregatorID,
		},
		w.GetNonceAndIncrement(),
	)

	_, receipt, err := w.AwaitedSignAndSendTransaction(client, txData)
	if err != nil {
		t.Error(err)
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		t.Errorf("transaction should have succeeded, expected %d got %d", 1, receipt.Status)
	}

	// was the pubkey stored ?
	attested, err := mgmtContractLib.GenContract.Attested(nil, aggregatorID)
	if err != nil {
		t.Error(err)
	}
	if !attested {
		t.Error("expected agg to be attested")
	}

	// do the same again
	aggregatorID = datagenerator.RandomAddress()
	txData = mgmtContractLib.CreateInitializeSecret(
		&obscurocommon.L1InitializeSecretTx{
			AggregatorID: &aggregatorID,
		},
		w.GetNonceAndIncrement(),
	)

	_, receipt, err = w.AwaitedSignAndSendTransaction(client, txData)
	if err != nil {
		t.Error(err)
	}

	if receipt.Status != types.ReceiptStatusFailed {
		t.Errorf("transaction should have failed, expected %d got %d", 0, receipt.Status)
	}
}

// attestedNodesCreateRollup attests a node by issuing a InitializeNetworkSecret, issues a rollups from the same node and verifies the rollup was stored
func attestedNodesCreateRollup(t *testing.T, mgmtContractLib *debugMgmtContractLib, w *debugWallet, client ethclient.EthClient) {
	rollup := datagenerator.RandomRollup()
	requesterID := &rollup.Header.Agg

	// the aggregator starts the network
	txData := mgmtContractLib.CreateInitializeSecret(
		&obscurocommon.L1InitializeSecretTx{
			AggregatorID: requesterID,
		},
		w.GetNonceAndIncrement(),
	)

	_, receipt, err := w.AwaitedSignAndSendTransaction(client, txData)
	if err != nil {
		t.Error(err)
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		t.Errorf("transaction should have succeeded, expected %d got %d", types.ReceiptStatusSuccessful, receipt.Status)
	}

	// issue a rollup from the attested node
	txData = mgmtContractLib.CreateRollup(&obscurocommon.L1RollupTx{Rollup: nodecommon.EncodeRollup(&rollup)}, w.GetNonceAndIncrement())
	_, receipt, err = w.AwaitedSignAndSendTransaction(client, txData)
	if err != nil {
		t.Error(err)
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		t.Errorf("transaction should have succeeded, expected %d got %d", types.ReceiptStatusSuccessful, receipt.Status)
	}

	// make sure the rollup was stored in the contract
	storedRollup, err := mgmtContractLib.GenContract.Rollups(nil, receipt.BlockNumber, big.NewInt(0))
	if err != nil {
		t.Error(err)
	}

	if storedRollup.Number.Int64() != rollup.Header.Number.Int64() ||
		!bytes.Equal(storedRollup.ParentHash[:], rollup.Header.ParentHash.Bytes()) ||
		!bytes.Equal(storedRollup.AggregatorID[:], rollup.Header.Agg.Bytes()) ||
		!bytes.Equal(storedRollup.L1Block[:], rollup.Header.L1Proof.Bytes()) {
		t.Error("stored rollup does not match the generated rollup")
	}
}

// nonAttestedNodesCannotAttest agg A initializes the network, agg B requests the secret, agg C issues response, but it's reverted
func nonAttestedNodesCannotAttest(t *testing.T, mgmtContractLib *debugMgmtContractLib, w *debugWallet, client ethclient.EthClient) {
	aggAPrivateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Error(err)
	}
	aggAID := crypto.PubkeyToAddress(aggAPrivateKey.PublicKey)

	// aggregator A starts the network secret
	txData := mgmtContractLib.CreateInitializeSecret(
		&obscurocommon.L1InitializeSecretTx{
			AggregatorID: &aggAID,
		},
		w.GetNonceAndIncrement(),
	)

	_, receipt, err := w.AwaitedSignAndSendTransaction(client, txData)
	if err != nil {
		t.Error(err)
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		t.Errorf("transaction should have succeeded, expected %d got %d", 1, receipt.Status)
	}

	// agg b requests the secret
	aggBPrivateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Error(err)
	}
	aggBID := crypto.PubkeyToAddress(aggBPrivateKey.PublicKey)

	txData = mgmtContractLib.CreateRequestSecret(
		&obscurocommon.L1RequestSecretTx{
			Attestation: datagenerator.RandomBytes(10),
		},
		w.GetNonceAndIncrement(),
	)

	_, receipt, err = w.AwaitedSignAndSendTransaction(client, txData)
	if err != nil {
		t.Error(err)
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		t.Errorf("transaction should have succeeded, expected %d got %d", 1, receipt.Status)
	}

	// agg c responds to the secret
	aggCPrivateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Error(err)
	}
	aggCID := crypto.PubkeyToAddress(aggCPrivateKey.PublicKey)

	fakeSecret := []byte{123}

	txData = mgmtContractLib.CreateRespondSecret(
		(&obscurocommon.L1RespondSecretTx{
			AttesterID:  aggCID,
			RequesterID: aggBID,
			Secret:      fakeSecret,
		}).Sign(aggCPrivateKey),
		w.GetNonceAndIncrement(),
	)

	_, receipt, err = w.AwaitedSignAndSendTransaction(client, txData)
	if err != nil {
		t.Error(err)
	}
	if receipt.Status != types.ReceiptStatusFailed {
		t.Errorf("transaction should have failed, expected %d got %d", 1, receipt.Status)
	}

	// agg c responds to the secret AGAIN, but trying to mimick aggregator A
	txData = mgmtContractLib.CreateRespondSecret(
		(&obscurocommon.L1RespondSecretTx{
			Secret:      fakeSecret,
			RequesterID: aggBID,
			AttesterID:  aggAID,
		}).Sign(aggCPrivateKey),
		w.GetNonceAndIncrement(),
	)

	_, receipt, err = w.AwaitedSignAndSendTransaction(client, txData)
	if err != nil {
		t.Error(err)
	}
	if receipt.Status != types.ReceiptStatusFailed {
		t.Errorf("transaction should have failed, expected %d got %d", 0, receipt.Status)
	}
}

// newlyAttestedNodesCanAttest agg A initializes the network, agg B requests the secret, agg C requests the secret, agg C is attested by agg A and agg B is attested by agg C
func newlyAttestedNodesCanAttest(t *testing.T, mgmtContractLib *debugMgmtContractLib, w *debugWallet, client ethclient.EthClient) {
	secretBytes := []byte("This is super random")
	// crypto.GenerateKey will generate a PK that does not play along this test
	aggAPrivateKey, err := crypto.ToECDSA(hexutil.MustDecode("0xc0083389f7a5925b662f8982080ced523bcc5e5dc33c6b1eaf11e288183e3c95"))
	if err != nil {
		t.Fatal(err)
	}
	aggAID := crypto.PubkeyToAddress(aggAPrivateKey.PublicKey)

	// the aggregator starts the network
	txData := mgmtContractLib.CreateInitializeSecret(
		&obscurocommon.L1InitializeSecretTx{
			AggregatorID:  &aggAID,
			InitialSecret: secretBytes,
		},
		w.GetNonceAndIncrement(),
	)

	_, receipt, err := w.AwaitedSignAndSendTransaction(client, txData)
	if err != nil {
		t.Error(err)
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		t.Errorf("transaction should have succeeded, expected %d got %d", 1, receipt.Status)
	}
	attested, err := mgmtContractLib.GenContract.Attested(nil, aggAID)
	if err != nil {
		t.Error(err)
	}
	if !attested {
		t.Error("expected agg to be attested")
	}

	// agg b requests the secret
	aggBPrivateKey, err := crypto.ToECDSA(hexutil.MustDecode("0x0d3de78eb7f26239a7ee32895a0b0898699ad3c4e5a910d0ffd65f707d2e63c4"))
	if err != nil {
		t.Fatal(err)
	}
	aggBID := crypto.PubkeyToAddress(aggBPrivateKey.PublicKey)

	txData = mgmtContractLib.CreateRequestSecret(
		&obscurocommon.L1RequestSecretTx{
			Attestation: datagenerator.RandomBytes(10),
		},
		w.GetNonceAndIncrement(),
	)
	_, receipt, err = w.AwaitedSignAndSendTransaction(client, txData)
	if err != nil {
		t.Error(err)
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		t.Errorf("transaction should have succeeded, expected %d got %d", 1, receipt.Status)
	}

	// agg C requests the secret
	aggCPrivateKey, err := crypto.ToECDSA(hexutil.MustDecode("0x2464a793cee0ea7103121fb1dfb6d021d80f43f3b5af39c7944b52db19a7ef30"))
	if err != nil {
		t.Fatal(err)
	}
	aggCID := crypto.PubkeyToAddress(aggCPrivateKey.PublicKey)

	txData = mgmtContractLib.CreateRequestSecret(
		&obscurocommon.L1RequestSecretTx{
			Attestation: datagenerator.RandomBytes(10),
		},
		w.GetNonceAndIncrement(),
	)

	_, receipt, err = w.AwaitedSignAndSendTransaction(client, txData)
	if err != nil {
		t.Error(err)
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		t.Errorf("transaction should have succeeded, expected %d got %d", 1, receipt.Status)
	}

	// Agg A responds to Agg C request
	txData = mgmtContractLib.CreateRespondSecret(
		(&obscurocommon.L1RespondSecretTx{
			Secret:      secretBytes,
			RequesterID: aggCID,
			AttesterID:  aggAID,
		}).Sign(aggAPrivateKey),
		w.GetNonceAndIncrement(),
	)
	_, receipt, err = w.AwaitedSignAndSendTransaction(client, txData)
	if err != nil {
		t.Error(err)
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		t.Errorf("transaction should have succeeded, expected %d got %d", 1, receipt.Status)
	}

	// test if aggregator is attested
	attested, err = mgmtContractLib.GenContract.Attested(nil, aggCID)
	if err != nil {
		t.Error(err)
	}
	if !attested {
		t.Error("expected agg to be attested")
	}

	// agg C attests agg B
	txData = mgmtContractLib.CreateRespondSecret(
		(&obscurocommon.L1RespondSecretTx{
			Secret:      secretBytes,
			RequesterID: aggBID,
			AttesterID:  aggCID,
		}).Sign(aggCPrivateKey),
		w.GetNonceAndIncrement(),
	)
	_, receipt, err = w.AwaitedSignAndSendTransaction(client, txData)
	if err != nil {
		t.Error(err)
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		t.Errorf("transaction should have succeeded, expected %d got %d", 1, receipt.Status)
	}

	// test if aggregator is attested
	attested, err = mgmtContractLib.GenContract.Attested(nil, aggBID)
	if err != nil {
		t.Error(err)
	}
	if !attested {
		t.Error("expected agg to be attested")
	}
}

// detectSimpleFork agg A initializes the network, agg A creates 3 correct rollups, then makes a depth 2 fork and expects the contract to detect
//
//               -> 3' -> 4' (contract marked with invalid withdrawals)
//   0 -> 1 -> 2 -> 3  -> 4
//

func detectSimpleFork(t *testing.T, mgmtContractLib *debugMgmtContractLib, w *debugWallet, client ethclient.EthClient) {
	secretBytes := []byte("This is super random")
	// crypto.GenerateKey will generate a PK that does not play along this test
	aggAPrivateKey, err := crypto.ToECDSA(hexutil.MustDecode("0xc0083389f7a5925b662f8982080ced523bcc5e5dc33c6b1eaf11e288183e3c95"))
	if err != nil {
		t.Fatal(err)
	}
	aggAID := crypto.PubkeyToAddress(aggAPrivateKey.PublicKey)

	// the aggregator starts the network
	txData := mgmtContractLib.CreateInitializeSecret(
		&obscurocommon.L1InitializeSecretTx{
			AggregatorID:  &aggAID,
			InitialSecret: secretBytes,
		},
		w.GetNonceAndIncrement(),
	)

	_, receipt, err := w.AwaitedSignAndSendTransaction(client, txData)
	if err != nil {
		t.Error(err)
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		t.Errorf("transaction should have succeeded, expected %d got %d", 1, receipt.Status)
	}
	attested, err := mgmtContractLib.GenContract.Attested(nil, aggAID)
	if err != nil {
		t.Error(err)
	}
	if !attested {
		t.Error("expected agg to be attested")
	}

	// Issue the genesis rollup
	rollup := datagenerator.RandomRollup()
	rollup.Header.Agg = aggAID

	txData = mgmtContractLib.CreateRollup(
		&obscurocommon.L1RollupTx{Rollup: nodecommon.EncodeRollup(&rollup)},
		w.GetNonceAndIncrement(),
	)

	issuedTx, receipt, err := w.AwaitedSignAndSendTransaction(client, txData)
	if err != nil {
		t.Error(err)
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		_, err := w.debugTransaction(client, issuedTx)
		if err != nil {
			t.Errorf("transaction should have suceeded, expected %d got %d - reason: %s", types.ReceiptStatusSuccessful, receipt.Status, err)
		}
	}

	// rollup meta data is actually stored
	found, rollupElement, err := mgmtContractLib.GenContract.GetRollupByHash(nil, rollup.Hash())
	if err != nil {
		t.Error(err)
	}

	if !found {
		t.Error("rollup not stored in tree")
	}

	if rollupElement.Rollup.Number.Int64() != rollup.Header.Number.Int64() ||
		!bytes.Equal(rollupElement.Rollup.ParentHash[:], rollup.Header.ParentHash.Bytes()) ||
		!bytes.Equal(rollupElement.Rollup.AggregatorID[:], rollup.Header.Agg.Bytes()) ||
		!bytes.Equal(rollupElement.Rollup.L1Block[:], rollup.Header.L1Proof.Bytes()) {
		t.Error("stored rollup does not match the generated rollup")
	}

	// Issues 3 rollups
	parentRollup := rollup
	for i := 0; i < 3; i++ {
		// issue rollup - make sure it comes from the attested aggregator
		r := datagenerator.RandomRollup()
		r.Header.Agg = aggAID
		r.Header.ParentHash = parentRollup.Header.Hash()

		// each rollup is child of the previous rollup
		parentRollup = r
		txData = mgmtContractLib.CreateRollup(
			&obscurocommon.L1RollupTx{Rollup: nodecommon.EncodeRollup(&r)},
			w.GetNonceAndIncrement(),
		)

		issuedTx, receipt, err := w.AwaitedSignAndSendTransaction(client, txData)
		if err != nil {
			t.Error(err)
		}

		if receipt.Status != types.ReceiptStatusSuccessful {
			_, err := w.debugTransaction(client, issuedTx)
			if err != nil {
				t.Errorf("transaction should have suceeded, expected %d got %d - reason: %s", types.ReceiptStatusSuccessful, receipt.Status, err)
			}
		}

		// rollup meta data is actually stored
		found, rollupElement, err := mgmtContractLib.GenContract.GetRollupByHash(nil, r.Header.Hash())
		if err != nil {
			t.Error(err)
		}

		if !found {
			t.Error("rollup not stored in tree")
		}

		if rollupElement.Rollup.Number.Int64() != r.Header.Number.Int64() ||
			!bytes.Equal(rollupElement.Rollup.ParentHash[:], r.Header.ParentHash.Bytes()) ||
			!bytes.Equal(rollupElement.Rollup.AggregatorID[:], r.Header.Agg.Bytes()) ||
			!bytes.Equal(rollupElement.Rollup.L1Block[:], r.Header.L1Proof.Bytes()) {
			t.Error("stored rollup does not match the generated rollup")
		}
	}

	// inserts a fork ( two rollups at same height / same parent )
	forks := make([]nodecommon.Rollup, 2)
	for i := 0; i < 2; i++ {
		r := datagenerator.RandomRollup()
		r.Header.Agg = aggAID

		// same parent
		r.Header.ParentHash = parentRollup.Header.Hash()

		// store these on the side as fork branches
		forks[i] = r

		fmt.Printf("insertion %d: new rollup %s - parent %s \n", i, r.Header.Hash(), parentRollup.Header.Hash())

		txData = mgmtContractLib.CreateRollup(
			&obscurocommon.L1RollupTx{Rollup: nodecommon.EncodeRollup(&r)},
			w.GetNonceAndIncrement(),
		)

		issuedTx, receipt, err = w.AwaitedSignAndSendTransaction(client, txData)
		if err != nil {
			t.Error(err)
		}

		if receipt.Status != types.ReceiptStatusSuccessful {
			_, err := w.debugTransaction(client, issuedTx)
			if err != nil {
				t.Errorf("transaction should have suceeded, expected %d got %d - reason: %s", types.ReceiptStatusSuccessful, receipt.Status, err)
			}
		}

		// rollup meta data is actually stored
		found, rollupElement, err = mgmtContractLib.GenContract.GetRollupByHash(nil, r.Hash())
		if err != nil {
			t.Error(err)
		}

		if !found {
			t.Error("rollup not stored in tree")
		}

		if rollupElement.Rollup.Number.Int64() != r.Header.Number.Int64() ||
			!bytes.Equal(rollupElement.Rollup.ParentHash[:], r.Header.ParentHash.Bytes()) ||
			!bytes.Equal(rollupElement.Rollup.AggregatorID[:], r.Header.Agg.Bytes()) ||
			!bytes.Equal(rollupElement.Rollup.L1Block[:], r.Header.L1Proof.Bytes()) {
			t.Error("stored rollup does not match the generated rollup")
		}
	}

	// create the fork
	for i, parentRollup := range forks {
		r := datagenerator.RandomRollup()
		r.Header.Agg = aggAID
		r.Header.ParentHash = parentRollup.Header.Hash()

		forks = append(forks, r)

		fmt.Printf("insertion %d: new rollup %s - parent %s \n", i, r.Header.Hash(), parentRollup.Header.Hash())

		txData = mgmtContractLib.CreateRollup(
			&obscurocommon.L1RollupTx{Rollup: nodecommon.EncodeRollup(&r)},
			w.GetNonceAndIncrement(),
		)

		issuedTx, receipt, err := w.AwaitedSignAndSendTransaction(client, txData)
		if err != nil {
			t.Error(err)
		}

		if receipt.Status != types.ReceiptStatusSuccessful {
			_, err := w.debugTransaction(client, issuedTx)
			if err != nil {
				t.Errorf("transaction should have suceeded, expected %d got %d - reason: %s", types.ReceiptStatusSuccessful, receipt.Status, err)
			}
		}

		// ensure it's retrievable
		found, ele, err := mgmtContractLib.GenContract.GetRollupByHash(nil, r.Header.Hash())
		if err != nil {
			t.Error(err)
		}
		if !found {
			fmt.Println(ele)
			t.Error("not found hash")
		}
	}

	available, err := mgmtContractLib.GenContract.IsWithdrawalAvailable(nil)
	if err != nil {
		t.Error(err)
	}

	if !available {
		t.Error("Withdrawals should be available at this stage")
	}

	// lock the contract
	parentRollup = forks[3]

	r := datagenerator.RandomRollup()
	r.Header.Agg = aggAID
	r.Header.ParentHash = parentRollup.Header.Hash()

	fmt.Printf("LAST Insertion : new rollup %s - parent %s \n", r.Header.Hash(), parentRollup.Header.Hash())

	txData = mgmtContractLib.CreateRollup(
		&obscurocommon.L1RollupTx{Rollup: nodecommon.EncodeRollup(&r)},
		w.GetNonceAndIncrement(),
	)

	issuedTx, receipt, err = w.AwaitedSignAndSendTransaction(client, txData)
	if err != nil {
		t.Error(err)
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		t.Errorf("transaction should have suceeded, expected %d got %d ", types.ReceiptStatusSuccessful, receipt.Status)
	}

	available, err = mgmtContractLib.GenContract.IsWithdrawalAvailable(nil)
	if err != nil {
		t.Error(err)
	}

	if available {
		t.Error("Withdrawals should NOT be available at this stage")
	}
}
