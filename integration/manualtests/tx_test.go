package manualtests

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ten-protocol/go-ten/integration/common/testlog"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/ten-protocol/go-ten/go/common/viewingkey"
	"github.com/ten-protocol/go-ten/go/ethadapter"
	"github.com/ten-protocol/go-ten/go/obsclient"
	"github.com/ten-protocol/go-ten/go/rpc"
	"github.com/ten-protocol/go-ten/go/wallet"
	"github.com/ten-protocol/go-ten/integration/datagenerator"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
)

const (
	_IDEFlag = "IDE"
)

// Testnet values used for quick debugging:
// l1 host: testnet-eth2network.uksouth.cloudapp.azure.com
// l1 port: 9000
// l2 host: erpc.testnet.ten.xyz
// l2 port: 81
// l2wallet: 8dfb8083da6275ae3e4f41e3e8a8c19d028d32c9247e24530933782f2a05035b -> 0xA58C60cc047592DE97BF1E8d2f225Fc5D959De77

var (
	l1Wallet = wallet.NewInMemoryWalletFromConfig(
		"5d1cffab85ddad285de2485ff09339e66e1e0acbfb9960c0df8231a1deb4994a",
		1337,
		gethlog.New())
	l1Host = "testnet-eth2network.uksouth.cloudapp.azure.com"
	l1Port = 9000

	l2Wallet = wallet.NewInMemoryWalletFromConfig(
		"4bfe14725e685901c062ccd4e220c61cf9c189897b6c78bd18d7f51291b2b8f8",
		777,
		gethlog.New())
	l2Host = "erpc.uat-testnet.ten.xyz"
	l2Port = 81
)

func TestL1IssueContractInteractWaitReceipt(t *testing.T) {
	t.Skipf("skip manual tests")

	var err error
	ethClient, err := ethadapter.NewEthClient(l1Host, uint(l1Port), 30*time.Second, gethlog.New())
	require.NoError(t, err)

	storeContractBytecode := "0x608060405234801561001057600080fd5b5061020b806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c80632e64cec11461003b57806370ef6e0b14610059575b600080fd5b610043610075565b60405161005091906100ab565b60405180910390f35b610073600480360381019061006e9190610161565b61007e565b005b60008054905090565b836000819055508260008190555050505050565b6000819050919050565b6100a581610092565b82525050565b60006020820190506100c0600083018461009c565b92915050565b600080fd5b600080fd5b6100d981610092565b81146100e457600080fd5b50565b6000813590506100f6816100d0565b92915050565b600080fd5b600080fd5b600080fd5b60008083601f840112610121576101206100fc565b5b8235905067ffffffffffffffff81111561013e5761013d610101565b5b60208301915083600182028301111561015a57610159610106565b5b9250929050565b6000806000806060858703121561017b5761017a6100c6565b5b6000610189878288016100e7565b945050602061019a878288016100e7565b935050604085013567ffffffffffffffff8111156101bb576101ba6100cb565b5b6101c78782880161010b565b92509250509295919450925056fea2646970667358221220eda68578fb741c32f26000b6c0273945f8322dd35f536c918e3d5a6193aaf62564736f6c63430008120033"

	estimatedTx, err := ethadapter.SetTxGasPrice(context.Background(), ethClient, &types.LegacyTx{
		Data: gethcommon.FromHex(storeContractBytecode),
	}, l1Wallet.Address(), l1Wallet.GetNonceAndIncrement(), 0, testlog.Logger())
	require.NoError(t, err)

	signedTx, err := l1Wallet.SignTransaction(estimatedTx)
	require.NoError(t, err)

	err = ethClient.SendTransaction(signedTx)
	require.NoError(t, err)

	fmt.Printf("Deployed Store Contract")
	fmt.Printf("Created Tx: %s \n", signedTx.Hash().Hex())
	fmt.Printf("Checking for tx receipt for %s \n", signedTx.Hash())
	var receipt *types.Receipt
	for start := time.Now(); time.Since(start) < time.Minute; time.Sleep(time.Second) {
		receipt, err = ethClient.TransactionReceipt(signedTx.Hash())
		if err == nil {
			break
		}
		if !errors.Is(err, ethereum.NotFound) {
			t.Fatal(err)
		}
	}

	if receipt == nil {
		t.Fatalf("Did not mine the transaction after %s seconds  - receipt: %+v", 30*time.Second, receipt)
	}
}

func TestL1IssueTxWaitReceipt(t *testing.T) {
	t.Skipf("skip manual tests")

	var err error
	ethClient, err := ethadapter.NewEthClient(l1Host, uint(l1Port), 30*time.Second, gethlog.New())
	require.NoError(t, err)

	toAddr := datagenerator.RandomAddress()
	nonce, err := ethClient.Nonce(l1Wallet.Address())
	require.NoError(t, err)

	l1Wallet.SetNonce(nonce)
	estimatedTx, err := ethadapter.SetTxGasPrice(context.Background(), ethClient, &types.LegacyTx{
		To:    &toAddr,
		Value: big.NewInt(100),
	}, l1Wallet.Address(), l1Wallet.GetNonceAndIncrement(), 0, testlog.Logger())
	require.NoError(t, err)

	signedTx, err := l1Wallet.SignTransaction(estimatedTx)
	require.NoError(t, err)

	err = ethClient.SendTransaction(signedTx)
	require.NoError(t, err)

	fmt.Printf("Created Tx: %s \n", signedTx.Hash().Hex())
	fmt.Printf("Checking for tx receipt for %s \n", signedTx.Hash())
	var receipt *types.Receipt
	for start := time.Now(); time.Since(start) < time.Minute; time.Sleep(time.Second) {
		receipt, err = ethClient.TransactionReceipt(signedTx.Hash())
		if err == nil {
			break
		}
		if !errors.Is(err, ethereum.NotFound) {
			t.Fatal(err)
		}
	}

	if receipt == nil {
		t.Fatalf("Did not mine the transaction after %s seconds  - receipt: %+v", 30*time.Second, receipt)
	}
}

func TestL2IssueContractInteractWaitReceipt(t *testing.T) {
	t.Skipf("skip manual tests")

	ctx := context.Background()

	vk, err := viewingkey.GenerateViewingKeyForWallet(l2Wallet)
	require.NoError(t, err)
	client, err := rpc.NewEncNetworkClient(fmt.Sprintf("ws://%s:%d", l2Host, l2Port), vk, gethlog.New())
	require.NoError(t, err)
	authClient := obsclient.NewAuthObsClient(client)

	balance, err := authClient.BalanceAt(context.Background(), nil)
	require.NoError(t, err)

	if balance.Cmp(big.NewInt(0)) <= 0 {
		t.Errorf("not enough balance: has %s has %s ten", l2Wallet.Address().Hex(), balance.String())
	}

	storeContractBytecode := "0x608060405234801561001057600080fd5b506101f3806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c80632e64cec11461003b5780636057361d14610059575b600080fd5b610043610075565b60405161005091906100da565b60405180910390f35b610073600480360381019061006e9190610126565b61007e565b005b60008054905090565b806000819055507febfcf7c0a1b09f6499e519a8d8bb85ce33cd539ec6cbd964e116cd74943ead1a33826040516100b6929190610194565b60405180910390a150565b6000819050919050565b6100d4816100c1565b82525050565b60006020820190506100ef60008301846100cb565b92915050565b600080fd5b610103816100c1565b811461010e57600080fd5b50565b600081359050610120816100fa565b92915050565b60006020828403121561013c5761013b6100f5565b5b600061014a84828501610111565b91505092915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b600061017e82610153565b9050919050565b61018e81610173565b82525050565b60006040820190506101a96000830185610185565b6101b660208301846100cb565b939250505056fea264697066735822122071ae4262a5da0c5c7b417d9c6cd13b57b8fcfe79c9c526b96f482ee67ff3136c64736f6c63430008120033"
	nonce, err := authClient.NonceAt(ctx, nil)
	require.NoError(t, err)

	l2Wallet.SetNonce(nonce)
	estimatedTx := authClient.EstimateGasAndGasPrice(&types.LegacyTx{
		Nonce:    l2Wallet.GetNonceAndIncrement(),
		Gas:      uint64(1_000_000),
		GasPrice: gethcommon.Big1,
		Data:     gethcommon.FromHex(storeContractBytecode),
	})
	require.NoError(t, err)

	contract, err := authClient.CallContract(ctx, ethereum.CallMsg{
		From: l2Wallet.Address(),
		Data: gethcommon.FromHex(storeContractBytecode),
	}, nil)
	if err != nil {
		fmt.Println(contract)
		return
	}

	signedTx, err := l2Wallet.SignTransaction(estimatedTx)
	require.NoError(t, err)

	err = authClient.SendTransaction(ctx, signedTx)
	require.NoError(t, err)

	fmt.Printf("Created Tx: %s \n", signedTx.Hash().Hex())
	fmt.Printf("Checking for tx receipt for %s \n", signedTx.Hash())
	var receipt *types.Receipt
	for start := time.Now(); time.Since(start) < time.Minute; time.Sleep(time.Second) {
		receipt, err = authClient.TransactionReceipt(ctx, signedTx.Hash())
		if err == nil {
			break
		}
		//
		// Currently when a receipt is not available the obscuro node is returning nil instead of err ethereum.NotFound
		// once that's fixed this commented block should be removed
		//if !errors.Is(err, ethereum.NotFound) {
		//	t.Fatal(err)
		//}
		if receipt != nil && receipt.Status == 1 {
			break
		}
		fmt.Printf("no tx receipt after %s - %s\n", time.Since(start), err)
	}

	if receipt == nil {
		t.Fatalf("Did not mine the transaction after %s seconds  - receipt: %+v", 30*time.Second, receipt)
	}
	if receipt.Status == 0 {
		t.Fatalf("Tx Failed")
	}

	abiJSONString := `
[
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": false,
				"internalType": "address",
				"name": "from",
				"type": "address"
			},
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "numb",
				"type": "uint256"
			}
		],
		"name": "Stored",
		"type": "event"
	},
	{
		"inputs": [
			{
				"internalType": "uint256",
				"name": "num",
				"type": "uint256"
			}
		],
		"name": "store",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "retrieve",
		"outputs": [
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"stateMutability": "view",
		"type": "function"
	}
]`

	myAbi, err := abi.JSON(strings.NewReader(abiJSONString))
	require.NoError(t, err)

	pack, err := myAbi.Pack("store", big.NewInt(1))
	if err != nil {
		return
	}

	estimatedTx = authClient.EstimateGasAndGasPrice(&types.LegacyTx{
		Nonce:    l2Wallet.GetNonceAndIncrement(),
		To:       &receipt.ContractAddress,
		Gas:      uint64(1_000_000),
		GasPrice: gethcommon.Big1,
		Data:     pack,
	})
	require.NoError(t, err)

	signedTx, err = l2Wallet.SignTransaction(estimatedTx)
	require.NoError(t, err)

	err = authClient.SendTransaction(ctx, signedTx)
	require.NoError(t, err)

	fmt.Printf("Created Tx: %s \n", signedTx.Hash().Hex())
	fmt.Printf("Checking for tx receipt for %s \n", signedTx.Hash())
	receipt = nil
	for start := time.Now(); time.Since(start) < time.Minute; time.Sleep(time.Second) {
		receipt, err = authClient.TransactionReceipt(ctx, signedTx.Hash())
		if err == nil {
			break
		}
		//
		// Currently when a receipt is not available the obscuro node is returning nil instead of err ethereum.NotFound
		// once that's fixed this commented block should be removed
		//if !errors.Is(err, ethereum.NotFound) {
		//	t.Fatal(err)
		//}
		if receipt != nil && receipt.Status == 1 {
			break
		}
		fmt.Printf("no tx receipt after %s - %s\n", time.Since(start), err)
	}

	if receipt == nil {
		t.Fatalf("Did not mine the transaction after %s seconds  - receipt: %+v", 30*time.Second, receipt)
	}
	if receipt.Status == 0 {
		t.Fatalf("Tx Failed")
	}
}

func TestL2IssueTxWaitReceipt(t *testing.T) {
	t.Skipf("skip manual tests")

	ctx := context.Background()

	vk, err := viewingkey.GenerateViewingKeyForWallet(l2Wallet)
	require.NoError(t, err)
	client, err := rpc.NewEncNetworkClient(fmt.Sprintf("ws://%s:%d", l2Host, l2Port), vk, gethlog.New())
	require.NoError(t, err)
	authClient := obsclient.NewAuthObsClient(client)

	balance, err := authClient.BalanceAt(context.Background(), nil)
	require.NoError(t, err)

	if balance.Cmp(big.NewInt(0)) <= 0 {
		t.Fatalf("not enough balance: has %s has %s ten", l2Wallet.Address().Hex(), balance.String())
	}
	fmt.Println("balance: ", balance.String())

	toAddr := datagenerator.RandomAddress()
	nonce, err := authClient.NonceAt(ctx, nil)
	require.NoError(t, err)

	l2Wallet.SetNonce(nonce)
	estimatedTx := authClient.EstimateGasAndGasPrice(&types.LegacyTx{
		Nonce:    l2Wallet.GetNonceAndIncrement(),
		To:       &toAddr,
		Value:    big.NewInt(100),
		Gas:      uint64(1_000_000),
		GasPrice: gethcommon.Big1,
	})
	require.NoError(t, err)

	signedTx, err := l2Wallet.SignTransaction(estimatedTx)
	require.NoError(t, err)

	err = authClient.SendTransaction(ctx, signedTx)
	require.NoError(t, err)

	fmt.Printf("Created Tx: %s \n", signedTx.Hash().Hex())
	fmt.Printf("Checking for tx receipt for %s \n", signedTx.Hash())
	var receipt *types.Receipt
	for start := time.Now(); time.Since(start) < time.Minute; time.Sleep(time.Second) {
		receipt, err = authClient.TransactionReceipt(ctx, signedTx.Hash())
		if err == nil {
			break
		}
		//
		// Currently when a receipt is not available the obscuro node is returning nil instead of err ethereum.NotFound
		// once that's fixed this commented block should be removed
		//if !errors.Is(err, ethereum.NotFound) {
		//	t.Fatal(err)
		//}
		if receipt != nil && receipt.Status == 1 {
			break
		}
		fmt.Printf("no tx receipt after %s - %s\n", time.Since(start), err)
	}

	if receipt == nil {
		t.Fatalf("Did not mine the transaction after %s seconds  - receipt: %+v", 30*time.Second, receipt)
	}
	if receipt.Status == 0 {
		t.Fatalf("Tx Failed")
	}
}

func TestL2IssueContractSubscribeChange(t *testing.T) {
	t.Skipf("skip manual tests")

	ctx := context.Background()

	vk, err := viewingkey.GenerateViewingKeyForWallet(l2Wallet)
	require.NoError(t, err)
	client, err := rpc.NewEncNetworkClient(fmt.Sprintf("ws://%s:%d", l2Host, l2Port), vk, gethlog.New())
	require.NoError(t, err)
	authClient := obsclient.NewAuthObsClient(client)

	balance, err := authClient.BalanceAt(context.Background(), nil)
	require.NoError(t, err)

	if balance.Cmp(big.NewInt(0)) <= 0 {
		t.Errorf("not enough balance: has %s has %s ten", l2Wallet.Address().Hex(), balance.String())
	}

	storeContractBytecode := "608060405234801561001057600080fd5b506101d7806100206000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c80630f2723ea14610030575b600080fd5b61004a60048036038101906100459190610137565b61004c565b005b8173ffffffffffffffffffffffffffffffffffffffff167fe86bd59ccd77aa1a9fbc46604e341e1dcc72f2a6e6637d5422736d645a71625e826040516100929190610186565b60405180910390a25050565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006100ce826100a3565b9050919050565b6100de816100c3565b81146100e957600080fd5b50565b6000813590506100fb816100d5565b92915050565b6000819050919050565b61011481610101565b811461011f57600080fd5b50565b6000813590506101318161010b565b92915050565b6000806040838503121561014e5761014d61009e565b5b600061015c858286016100ec565b925050602061016d85828601610122565b9150509250929050565b61018081610101565b82525050565b600060208201905061019b6000830184610177565b9291505056fea2646970667358221220d45d57b217a07cfd4ceecc8f5d2a9194a098e5fe1de98bc56c24cc400e21c96064736f6c63430008120033"
	nonce, err := authClient.NonceAt(ctx, nil)
	require.NoError(t, err)

	l2Wallet.SetNonce(nonce)
	estimatedTx := authClient.EstimateGasAndGasPrice(&types.LegacyTx{
		Nonce:    l2Wallet.GetNonceAndIncrement(),
		Gas:      uint64(1_000_000),
		GasPrice: gethcommon.Big1,
		Data:     gethcommon.FromHex(storeContractBytecode),
	})
	require.NoError(t, err)

	signedTx, err := l2Wallet.SignTransaction(estimatedTx)
	require.NoError(t, err)

	err = authClient.SendTransaction(ctx, signedTx)
	require.NoError(t, err)

	fmt.Printf("Created Tx: %s \n", signedTx.Hash().Hex())
	fmt.Printf("Checking for tx receipt for %s \n", signedTx.Hash())
	var receipt *types.Receipt
	for start := time.Now(); time.Since(start) < time.Minute; time.Sleep(time.Second) {
		receipt, err = authClient.TransactionReceipt(ctx, signedTx.Hash())
		if err == nil {
			break
		}
		//
		// Currently when a receipt is not available the obscuro node is returning nil instead of err ethereum.NotFound
		// once that's fixed this commented block should be removed
		//if !errors.Is(err, ethereum.NotFound) {
		//	t.Fatal(err)
		//}
		if receipt != nil && receipt.Status == 1 {
			break
		}
		fmt.Printf("no tx receipt after %s - %s\n", time.Since(start), err)
	}

	if receipt == nil {
		t.Fatalf("Did not mine the transaction after %s seconds  - receipt: %+v", 30*time.Second, receipt)
	}
	if receipt.Status == 0 {
		t.Fatalf("Tx Failed")
	}

	abiJSONString := `[
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": true,
				"internalType": "address",
				"name": "addr",
				"type": "address"
			},
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "value",
				"type": "uint256"
			}
		],
		"name": "IndexedAddressAndNumber",
		"type": "event"
	},
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "_addr",
				"type": "address"
			},
			{
				"internalType": "uint256",
				"name": "_value",
				"type": "uint256"
			}
		],
		"name": "echo",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	}
]`

	myAbi, err := abi.JSON(strings.NewReader(abiJSONString))
	require.NoError(t, err)

	pack, err := myAbi.Pack("echo", l2Wallet.Address(), big.NewInt(42))
	require.NoError(t, err)

	estimatedTx = authClient.EstimateGasAndGasPrice(&types.LegacyTx{
		Nonce:    l2Wallet.GetNonceAndIncrement(),
		To:       &receipt.ContractAddress,
		Gas:      uint64(1_000_000),
		GasPrice: gethcommon.Big1,
		Data:     pack,
	})
	require.NoError(t, err)

	signedTx, err = l2Wallet.SignTransaction(estimatedTx)
	require.NoError(t, err)

	err = authClient.SendTransaction(ctx, signedTx)
	require.NoError(t, err)

	fmt.Printf("Created Tx: %s \n", signedTx.Hash().Hex())
	fmt.Printf("Checking for tx receipt for %s \n", signedTx.Hash())
	receipt = nil
	for start := time.Now(); time.Since(start) < time.Minute; time.Sleep(time.Second) {
		receipt, err = authClient.TransactionReceipt(ctx, signedTx.Hash())
		if err == nil {
			break
		}
		//
		// Currently when a receipt is not available the obscuro node is returning nil instead of err ethereum.NotFound
		// once that's fixed this commented block should be removed
		//if !errors.Is(err, ethereum.NotFound) {
		//	t.Fatal(err)
		//}
		if receipt != nil && receipt.Status == 1 {
			break
		}
		fmt.Printf("no tx receipt after %s - %s\n", time.Since(start), err)
	}

	if receipt == nil {
		t.Fatalf("Did not mine the transaction after %s seconds  - receipt: %+v", 30*time.Second, receipt)
	}
	if receipt.Status == 0 {
		t.Fatalf("Tx Failed")
	}

	result := json.RawMessage{}
	err = client.Call(&result, "debug_eventLogRelevancy", signedTx.Hash())
	assert.NoError(t, err)
	fmt.Println(err)
}
