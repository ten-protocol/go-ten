package smartcontract

import (
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/stretchr/testify/assert"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/constants"
	"github.com/ten-protocol/go-ten/go/common/signature"
	"github.com/ten-protocol/go-ten/go/ethadapter"
	"github.com/ten-protocol/go-ten/go/ethadapter/mgmtcontractlib"
	"github.com/ten-protocol/go-ten/go/wallet"
	"github.com/ten-protocol/go-ten/integration"
	"github.com/ten-protocol/go-ten/integration/common/testlog"
	"github.com/ten-protocol/go-ten/integration/datagenerator"
	"github.com/ten-protocol/go-ten/integration/eth2network"
	"github.com/ten-protocol/go-ten/integration/simulation/network"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

// netInfo is a bag holder struct for output data from the execution/run of a network
type netInfo struct {
	ethClients  []ethadapter.EthClient
	wallets     []wallet.Wallet
	eth2Network eth2network.PosEth2Network
}

var testLogs = "../.build/noderunner/"

func init() { //nolint:gochecknoinits
	testlog.Setup(&testlog.Cfg{
		LogDir:      testLogs,
		TestType:    "smartcontracts",
		TestSubtype: "test",
		LogLevel:    log.LvlInfo,
	})
}

// runGethNetwork runs a geth network with one prefunded wallet
func runGethNetwork(t *testing.T) *netInfo {
	// make sure the geth network binaries exist
	binDir, err := eth2network.EnsureBinariesExist()
	if err != nil {
		t.Fatal(err)
	}

	// prefund one wallet as the worker wallet
	workerWallet := datagenerator.RandomWallet(integration.EthereumChainID)

	startPort := integration.TestPorts.TestManagementContractPort
	eth2Network := eth2network.NewPosEth2Network(
		binDir,
		startPort+integration.DefaultGethNetworkPortOffset,
		startPort+integration.DefaultPrysmP2PPortOffset,
		startPort+integration.DefaultGethAUTHPortOffset, // RPC
		startPort+integration.DefaultGethWSPortOffset,
		startPort+integration.DefaultGethHTTPPortOffset,
		startPort+integration.DefaultPrysmRPCPortOffset,
		startPort+integration.DefaultPrysmGatewayPortOffset,
		integration.EthereumChainID,
		3*time.Minute,
		workerWallet.Address().String(),
	)

	if err = eth2Network.Start(); err != nil {
		t.Fatal(err)
	}

	// create a client that is connected to node 0 of the network
	client, err := ethadapter.NewEthClient("127.0.0.1", uint(startPort+100), 60*time.Second, gethcommon.HexToAddress("0x0"), testlog.Logger())
	if err != nil {
		t.Fatal(err)
	}

	return &netInfo{
		ethClients:  []ethadapter.EthClient{client},
		wallets:     []wallet.Wallet{workerWallet},
		eth2Network: eth2Network,
	}
}

func TestManagementContract(t *testing.T) {
	// run tests on one network
	sim := runGethNetwork(t)
	defer sim.eth2Network.Stop() //nolint: errcheck

	// set up the client and the (debug) wallet
	client := sim.ethClients[0]
	w := newDebugWallet(sim.wallets[0])

	for name, test := range map[string]func(*testing.T, *debugMgmtContractLib, *debugWallet, ethadapter.EthClient){
		"secretCannotBeInitializedTwice":     secretCannotBeInitializedTwice,
		"nonAttestedNodesCannotCreateRollup": nonAttestedNodesCannotCreateRollup,
		"attestedNodesCreateRollup":          attestedNodesCreateRollup,
		"nonAttestedNodesCannotAttest":       nonAttestedNodesCannotAttest,
		"newlyAttestedNodesCanAttest":        newlyAttestedNodesCanAttest,
	} {
		t.Run(name, func(t *testing.T) {
			bytecode, err := constants.Bytecode()
			if err != nil {
				panic(err)
			}

			nonce, err := client.Nonce(w.Address())
			if err != nil {
				t.Error(err)
			}

			w.SetNonce(nonce)
			// deploy the same contract to a new address
			receipt, err := network.DeployContract(client, w, bytecode)
			if err != nil {
				t.Error(err)
			}

			_, err = network.InitializeContract(client, w, receipt.ContractAddress)
			if err != nil {
				t.Error(err)
			}

			nonce, err = client.Nonce(w.Address())
			if err != nil {
				t.Error(err)
			}

			w.SetNonce(nonce)

			// run the test using the new contract, but same wallet
			test(t,
				newDebugMgmtContractLib(receipt.ContractAddress, client.EthClient(),
					mgmtcontractlib.NewMgmtContractLib(&receipt.ContractAddress, testlog.Logger())),
				w,
				client,
			)
		})
	}
}

// nonAttestedNodesCannotCreateRollup issues a rollup from a node that did not receive the secret network key
func nonAttestedNodesCannotCreateRollup(t *testing.T, mgmtContractLib *debugMgmtContractLib, w *debugWallet, client ethadapter.EthClient) {
	block, err := client.FetchHeadBlock()
	if err != nil {
		t.Error(err)
	}

	rollup := datagenerator.RandomRollup(block)

	encodedRollup, err := common.EncodeRollup(&rollup)
	if err != nil {
		t.Error(err)
	}
	txData := mgmtContractLib.CreateRollup(&ethadapter.L1RollupTx{Rollup: encodedRollup})

	_, _, err = w.AwaitedSignAndSendTransaction(client, txData)
	if err == nil || !assert.Contains(t, err.Error(), "execution reverted") {
		t.Error(err)
	}
}

// secretCannotBeInitializedTwice issues the InitializeNetworkSecret twice, failing the second time
func secretCannotBeInitializedTwice(t *testing.T, mgmtContractLib *debugMgmtContractLib, w *debugWallet, client ethadapter.EthClient) {
	aggregatorID := datagenerator.RandomAddress()
	txData := mgmtContractLib.CreateInitializeSecret(
		&ethadapter.L1InitializeSecretTx{
			EnclaveID: &aggregatorID,
		},
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
			EnclaveID: &aggregatorID,
		},
	)

	_, _, err = w.AwaitedSignAndSendTransaction(client, txData)
	if err == nil || !assert.Contains(t, err.Error(), "execution reverted") {
		t.Error(err)
	}
}

// attestedNodesCreateRollup attests a node by issuing a InitializeNetworkSecret, issues a rollups from the same node and verifies the rollup was stored
func attestedNodesCreateRollup(t *testing.T, mgmtContractLib *debugMgmtContractLib, w *debugWallet, client ethadapter.EthClient) {
	block, err := client.FetchHeadBlock()
	if err != nil {
		t.Error(err)
	}

	pk := datagenerator.RandomPrivateKey()
	enclaveID := crypto.PubkeyToAddress(pk.PublicKey)

	rollup := datagenerator.RandomRollup(block)
	rollup.Header.Signature, err = signature.Sign(rollup.Hash().Bytes(), pk)
	if err != nil {
		t.Error(err)
	}

	// the aggregator starts the network
	txData := mgmtContractLib.CreateInitializeSecret(
		&ethadapter.L1InitializeSecretTx{
			EnclaveID: &enclaveID,
		},
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
			EnclaveID: &aggAID,
		},
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
		true,
	)

	_, _, err = w.AwaitedSignAndSendTransaction(client, txData)
	if err == nil || !assert.Contains(t, err.Error(), "execution reverted") {
		t.Error(err)
	}

	// agg c responds to the secret AGAIN, but trying to mimick aggregator A
	txData = mgmtContractLib.CreateRespondSecret(
		(&ethadapter.L1RespondSecretTx{
			Secret:      fakeSecret,
			RequesterID: aggBID,
			AttesterID:  aggAID,
		}).Sign(aggCPrivateKey),
		true,
	)

	_, _, err = w.AwaitedSignAndSendTransaction(client, txData)
	if err == nil || !assert.Contains(t, err.Error(), "execution reverted") {
		t.Error(err)
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
			EnclaveID:     &aggAID,
			InitialSecret: secretBytes,
		},
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
