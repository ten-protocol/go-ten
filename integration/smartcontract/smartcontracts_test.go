package smartcontract

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/obscuronet/go-obscuro/contracts/managementcontract"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/ethadapter"
	"github.com/obscuronet/go-obscuro/go/ethadapter/mgmtcontractlib"
	"github.com/obscuronet/go-obscuro/go/wallet"
	"github.com/obscuronet/go-obscuro/integration"
	"github.com/obscuronet/go-obscuro/integration/datagenerator"
	"github.com/obscuronet/go-obscuro/integration/gethnetwork"
	"github.com/obscuronet/go-obscuro/integration/simulation/network"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

// netInfo is a bag holder struct for output data from the execution/run of a network
type netInfo struct {
	ethClients  []ethadapter.EthClient
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
	address := fmt.Sprintf("ws://127.0.0.1:%d", gethNetwork.WebSocketPorts[0])
	client, err := ethadapter.NewEthClient(address, 30*time.Second, gethcommon.HexToAddress("0x0"))
	if err != nil {
		return nil
	}

	return &netInfo{
		ethClients:  []ethadapter.EthClient{client},
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

	for name, test := range map[string]func(*testing.T, *debugMgmtContractLib, *debugWallet, ethadapter.EthClient){
		"secretCannotBeInitializedTwice":     secretCannotBeInitializedTwice,
		"nonAttestedNodesCannotCreateRollup": nonAttestedNodesCannotCreateRollup,
		"attestedNodesCreateRollup":          attestedNodesCreateRollup,
		"nonAttestedNodesCannotAttest":       nonAttestedNodesCannotAttest,
		"newlyAttestedNodesCanAttest":        newlyAttestedNodesCanAttest,
		"attestedNodeHostAddressesAreStored": attestedNodeHostAddressesAreStored,
		"detectSimpleFork":                   detectSimpleFork,
	} {
		t.Run(name, func(t *testing.T) {
			bytecode, err := managementcontract.Bytecode()
			if err != nil {
				panic(err)
			}
			// deploy the same contract to a new address
			contractAddr, err := network.DeployContract(client, w, bytecode)
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
func nonAttestedNodesCannotCreateRollup(t *testing.T, mgmtContractLib *debugMgmtContractLib, w *debugWallet, client ethadapter.EthClient) {
	rollup := datagenerator.RandomRollup()
	encodedRollup, err := common.EncodeRollup(&rollup)
	if err != nil {
		t.Error(err)
	}
	txData := mgmtContractLib.CreateRollup(
		&ethadapter.L1RollupTx{Rollup: encodedRollup},
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
func secretCannotBeInitializedTwice(t *testing.T, mgmtContractLib *debugMgmtContractLib, w *debugWallet, client ethadapter.EthClient) {
	aggregatorID := datagenerator.RandomAddress()
	txData := mgmtContractLib.CreateInitializeSecret(
		&ethadapter.L1InitializeSecretTx{
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
		&ethadapter.L1InitializeSecretTx{
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
func attestedNodesCreateRollup(t *testing.T, mgmtContractLib *debugMgmtContractLib, w *debugWallet, client ethadapter.EthClient) {
	rollup := datagenerator.RandomRollup()
	requesterID := &rollup.Header.Agg

	// the aggregator starts the network
	txData := mgmtContractLib.CreateInitializeSecret(
		&ethadapter.L1InitializeSecretTx{
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
	err = mgmtContractLib.AwaitedIssueRollup(rollup, client, w)
	if err != nil {
		t.Error(err)
	}
}

// nonAttestedNodesCannotAttest agg A initializes the network, agg B requests the secret, agg C issues response, but it's reverted
func nonAttestedNodesCannotAttest(t *testing.T, mgmtContractLib *debugMgmtContractLib, w *debugWallet, client ethadapter.EthClient) {
	aggAPrivateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Error(err)
	}
	aggAID := crypto.PubkeyToAddress(aggAPrivateKey.PublicKey)

	// aggregator A starts the network secret
	txData := mgmtContractLib.CreateInitializeSecret(
		&ethadapter.L1InitializeSecretTx{
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
		&ethadapter.L1RequestSecretTx{
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
		(&ethadapter.L1RespondSecretTx{
			AttesterID:  aggCID,
			RequesterID: aggBID,
			Secret:      fakeSecret,
		}).Sign(aggCPrivateKey),
		w.GetNonceAndIncrement(),
		true,
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
		(&ethadapter.L1RespondSecretTx{
			Secret:      fakeSecret,
			RequesterID: aggBID,
			AttesterID:  aggAID,
		}).Sign(aggCPrivateKey),
		w.GetNonceAndIncrement(),
		true,
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
func newlyAttestedNodesCanAttest(t *testing.T, mgmtContractLib *debugMgmtContractLib, w *debugWallet, client ethadapter.EthClient) {
	secretBytes := []byte("This is super random")
	// crypto.GenerateKey will generate a PK that does not play along this test
	aggAPrivateKey, err := crypto.ToECDSA(hexutil.MustDecode("0xc0083389f7a5925b662f8982080ced523bcc5e5dc33c6b1eaf11e288183e3c95"))
	if err != nil {
		t.Fatal(err)
	}
	aggAID := crypto.PubkeyToAddress(aggAPrivateKey.PublicKey)

	// the aggregator starts the network
	txData := mgmtContractLib.CreateInitializeSecret(
		&ethadapter.L1InitializeSecretTx{
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
		&ethadapter.L1RequestSecretTx{
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
		&ethadapter.L1RequestSecretTx{
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
		(&ethadapter.L1RespondSecretTx{
			Secret:      secretBytes,
			RequesterID: aggCID,
			AttesterID:  aggAID,
		}).Sign(aggAPrivateKey),
		w.GetNonceAndIncrement(),
		true,
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
		(&ethadapter.L1RespondSecretTx{
			Secret:      secretBytes,
			RequesterID: aggBID,
			AttesterID:  aggCID,
		}).Sign(aggCPrivateKey),
		w.GetNonceAndIncrement(),
		true,
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

// attestedNodeHostAddressesAreStored agg A initializes the network, agg B becomes attested, agg C is rejected. Only A and B's host addresses are stored in the management contract
func attestedNodeHostAddressesAreStored(t *testing.T, mgmtContractLib *debugMgmtContractLib, w *debugWallet, client ethadapter.EthClient) {
	aggAHostAddr := "aggAHostAddr"
	aggBHostAddr := "aggBHostAddr"

	secretBytes := []byte("This is super random")
	// crypto.GenerateKey will generate a PK that does not play along this test
	aggAPrivateKey, err := crypto.ToECDSA(hexutil.MustDecode("0xc0083389f7a5925b662f8982080ced523bcc5e5dc33c6b1eaf11e288183e3c95"))
	if err != nil {
		t.Fatal(err)
	}
	aggAID := crypto.PubkeyToAddress(aggAPrivateKey.PublicKey)

	// the aggregator starts the network
	txData := mgmtContractLib.CreateInitializeSecret(
		&ethadapter.L1InitializeSecretTx{
			AggregatorID:  &aggAID,
			InitialSecret: secretBytes,
			HostAddress:   aggAHostAddr,
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
	aggBPrivateKey, err := crypto.ToECDSA(hexutil.MustDecode("0x0d3de78eb7f26239a7ee32895a0b0898699ad3c4e5a910d0ffd65f707d2e63c4"))
	if err != nil {
		t.Fatal(err)
	}
	aggBID := crypto.PubkeyToAddress(aggBPrivateKey.PublicKey)

	txData = mgmtContractLib.CreateRequestSecret(
		&ethadapter.L1RequestSecretTx{
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
	txData = mgmtContractLib.CreateRequestSecret(
		&ethadapter.L1RequestSecretTx{
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

	// Agg A only responds to Agg B request
	txData = mgmtContractLib.CreateRespondSecret(
		(&ethadapter.L1RespondSecretTx{
			Secret:      secretBytes,
			RequesterID: aggBID,
			AttesterID:  aggAID,
			HostAddress: aggBHostAddr,
		}).Sign(aggAPrivateKey),
		w.GetNonceAndIncrement(),
		true,
	)
	_, receipt, err = w.AwaitedSignAndSendTransaction(client, txData)
	if err != nil {
		t.Error(err)
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		t.Errorf("transaction should have succeeded, expected %d got %d", 1, receipt.Status)
	}

	hostAddresses, err := mgmtContractLib.GenContract.GetHostAddresses(nil)
	if err != nil {
		t.Error(err)
	}
	expectedHostAddresses := []string{aggAHostAddr, aggBHostAddr}
	if !reflect.DeepEqual(hostAddresses, expectedHostAddresses) {
		t.Errorf("expected to find host addresses %s, found %s", expectedHostAddresses, hostAddresses)
	}
}

// detectSimpleFork agg A initializes the network, agg A creates 3 correct rollups, then makes a depth 2 fork and expects the contract to detect
//
// -> 4'-> 5'
// 0 -> 1 -> 2 -> 3 -> 4 -> 5  -> 6 (contract marked with invalid withdrawals)
func detectSimpleFork(t *testing.T, mgmtContractLib *debugMgmtContractLib, w *debugWallet, client ethadapter.EthClient) {
	secretBytes := []byte("This is super random")
	// crypto.GenerateKey will generate a PK that does not play along this test
	aggAPrivateKey, err := crypto.ToECDSA(hexutil.MustDecode("0xc0083389f7a5925b662f8982080ced523bcc5e5dc33c6b1eaf11e288183e3c95"))
	if err != nil {
		t.Fatal(err)
	}
	aggAID := crypto.PubkeyToAddress(aggAPrivateKey.PublicKey)

	// the aggregator starts the network
	txData := mgmtContractLib.CreateInitializeSecret(
		&ethadapter.L1InitializeSecretTx{
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

	// Issue a genesis rollup
	rollup := datagenerator.RandomRollup()
	rollup.Header.Agg = aggAID

	err = mgmtContractLib.AwaitedIssueRollup(rollup, client, w)
	if err != nil {
		t.Error(err)
	}
	t.Logf("Issued Rollup: %s parent: %s", rollup.Hash(), rollup.Header.ParentHash)

	// Issues 3 rollups
	parentRollup := rollup
	for i := 0; i < 3; i++ {
		// issue rollup - make sure it comes from the attested aggregator
		r := datagenerator.RandomRollup()
		r.Header.Agg = aggAID
		r.Header.ParentHash = parentRollup.Header.Hash()

		// each rollup is child of the previous rollup
		parentRollup = r

		// issue the rollup
		err = mgmtContractLib.AwaitedIssueRollup(r, client, w)
		if err != nil {
			t.Error(err)
		}
		t.Logf("Issued Rollup: %s parent: %s", r.Hash(), r.Header.ParentHash)
	}

	// inserts a fork ( two rollups at same height / same parent )
	splitPoint := make([]common.EncryptedRollup, 2)
	for i := 0; i < 2; i++ {
		r := datagenerator.RandomRollup()
		r.Header.Agg = aggAID

		// same parent
		r.Header.ParentHash = parentRollup.Header.Hash()

		// store these on the side as fork branches
		splitPoint[i] = r

		// issue the rollup
		err = mgmtContractLib.AwaitedIssueRollup(r, client, w)
		if err != nil {
			t.Error(err)
		}
		t.Logf("Issued Rollup: %s parent: %s", r.Hash(), r.Header.ParentHash)
	}

	// create the fork
	forks := make([]common.EncryptedRollup, 2)
	for i, parentRollup := range splitPoint {
		r := datagenerator.RandomRollup()
		r.Header.Agg = aggAID
		r.Header.ParentHash = parentRollup.Header.Hash()

		forks[i] = r

		// issue the rollup
		err = mgmtContractLib.AwaitedIssueRollup(r, client, w)
		if err != nil {
			t.Error(err)
		}
		t.Logf("Issued Rollup: %s parent: %s", r.Hash(), r.Header.ParentHash)
	}

	available, err := mgmtContractLib.GenContract.IsWithdrawalAvailable(nil)
	if err != nil {
		t.Error(err)
	}

	if !available {
		t.Error("Withdrawals should be available at this stage")
	}

	// lock the contract
	parentRollup = forks[1]

	r := datagenerator.RandomRollup()
	r.Header.Agg = aggAID
	r.Header.ParentHash = parentRollup.Header.Hash()

	t.Logf("LAST Issued Rollup: %s parent: %s", r.Hash(), r.Header.ParentHash)

	encodedRollup, err := common.EncodeRollup(&r)
	if err != nil {
		t.Error(err)
	}
	txData = mgmtContractLib.CreateRollup(
		&ethadapter.L1RollupTx{Rollup: encodedRollup},
		w.GetNonceAndIncrement(),
	)

	_, receipt, err = w.AwaitedSignAndSendTransaction(client, txData)
	if err != nil {
		t.Error(err)
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		t.Errorf("transaction should have succeeded, expected %d got %d ", types.ReceiptStatusSuccessful, receipt.Status)
	}

	available, err = mgmtContractLib.GenContract.IsWithdrawalAvailable(nil)
	if err != nil {
		t.Error(err)
	}

	if available {
		t.Error("Withdrawals should NOT be available at this stage")
	}
}
