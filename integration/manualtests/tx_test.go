package manualtests

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/ethadapter"
	"github.com/obscuronet/go-obscuro/go/obsclient"
	"github.com/obscuronet/go-obscuro/go/rpc"
	"github.com/obscuronet/go-obscuro/go/wallet"
	"github.com/obscuronet/go-obscuro/integration/datagenerator"
	"github.com/stretchr/testify/assert"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
)

const (
	_IDEFlag = "IDE"
)

// Testnet values used for quick debugging:
// l1 host: testnet-eth2network.uksouth.azurecontainer.io
// l1 port: 9000
// l2 host: testnet.obscu.ro
// l2 port: 13001
// l2wallet: 8dfb8083da6275ae3e4f41e3e8a8c19d028d32c9247e24530933782f2a05035b

var (
	l1Wallet = wallet.NewInMemoryWalletFromConfig(
		"5d1cffab85ddad285de2485ff09339e66e1e0acbfb9960c0df8231a1deb4994a",
		1337,
		gethlog.New())
	l1Host = "testnet-eth2network.uksouth.azurecontainer.io"
	l1Port = 9000

	l2Wallet = wallet.NewInMemoryWalletFromConfig(
		"4bfe14725e685901c062ccd4e220c61cf9c189897b6c78bd18d7f51291b2b8f8",
		777,
		gethlog.New())
	l2Host = "localhost"
	l2Port = 37900
)

func TestL1IssueContractInteractWaitReceipt(t *testing.T) {
	// test is skipped by default to avoid breaking CI - set env flag in `Run Configurations` to run it in IDE
	if os.Getenv(_IDEFlag) == "" {
		t.Skipf("set flag %s to run this test in the IDE", _IDEFlag)
	}
	var err error
	ethClient, err := ethadapter.NewEthClient(l1Host, uint(l1Port), 30*time.Second, common.L2Address{}, gethlog.New())
	assert.Nil(t, err)

	storeContractBytecode := "0x608060405234801561001057600080fd5b5061020b806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c80632e64cec11461003b57806370ef6e0b14610059575b600080fd5b610043610075565b60405161005091906100ab565b60405180910390f35b610073600480360381019061006e9190610161565b61007e565b005b60008054905090565b836000819055508260008190555050505050565b6000819050919050565b6100a581610092565b82525050565b60006020820190506100c0600083018461009c565b92915050565b600080fd5b600080fd5b6100d981610092565b81146100e457600080fd5b50565b6000813590506100f6816100d0565b92915050565b600080fd5b600080fd5b600080fd5b60008083601f840112610121576101206100fc565b5b8235905067ffffffffffffffff81111561013e5761013d610101565b5b60208301915083600182028301111561015a57610159610106565b5b9250929050565b6000806000806060858703121561017b5761017a6100c6565b5b6000610189878288016100e7565b945050602061019a878288016100e7565b935050604085013567ffffffffffffffff8111156101bb576101ba6100cb565b5b6101c78782880161010b565b92509250509295919450925056fea2646970667358221220eda68578fb741c32f26000b6c0273945f8322dd35f536c918e3d5a6193aaf62564736f6c63430008120033"

	nonce, err := ethClient.Nonce(l1Wallet.Address())
	assert.Nil(t, err)

	l1Wallet.SetNonce(nonce)
	estimatedTx, err := ethClient.EstimateGasAndGasPrice(&types.LegacyTx{
		Nonce: l1Wallet.GetNonceAndIncrement(),
		Data:  gethcommon.FromHex(storeContractBytecode),
	}, l1Wallet.Address())
	assert.Nil(t, err)

	signedTx, err := l1Wallet.SignTransaction(estimatedTx)
	assert.Nil(t, err)

	err = ethClient.SendTransaction(signedTx)
	assert.Nil(t, err)

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
	// test is skipped by default to avoid breaking CI - set env flag in `Run Configurations` to run it in IDE
	if os.Getenv(_IDEFlag) == "" {
		t.Skipf("set flag %s to run this test in the IDE", _IDEFlag)
	}

	var err error
	ethClient, err := ethadapter.NewEthClient(l1Host, uint(l1Port), 30*time.Second, common.L2Address{}, gethlog.New())
	assert.Nil(t, err)

	toAddr := datagenerator.RandomAddress()
	nonce, err := ethClient.Nonce(l1Wallet.Address())
	assert.Nil(t, err)

	l1Wallet.SetNonce(nonce)
	estimatedTx, err := ethClient.EstimateGasAndGasPrice(&types.LegacyTx{
		Nonce: l1Wallet.GetNonceAndIncrement(),
		To:    &toAddr,
		Value: big.NewInt(100),
	}, l1Wallet.Address())
	assert.Nil(t, err)

	signedTx, err := l1Wallet.SignTransaction(estimatedTx)
	assert.Nil(t, err)

	err = ethClient.SendTransaction(signedTx)
	assert.Nil(t, err)

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
	// test is skipped by default to avoid breaking CI - set env flag in `Run Configurations` to run it in IDE
	if os.Getenv(_IDEFlag) == "" {
		t.Skipf("set flag %s to run this test in the IDE", _IDEFlag)
	}

	ctx := context.Background()

	vk, err := rpc.GenerateAndSignViewingKey(l2Wallet)
	assert.Nil(t, err)
	client, err := rpc.NewEncNetworkClient(fmt.Sprintf("ws://%s:%d", l2Host, l2Port), vk, gethlog.New())
	assert.Nil(t, err)
	authClient := obsclient.NewAuthObsClient(client)

	balance, err := authClient.BalanceAt(context.Background(), nil)
	assert.Nil(t, err)

	if balance.Cmp(big.NewInt(0)) <= 0 {
		t.Errorf("not enough balance: has %s has %s obx", l2Wallet.Address().Hex(), balance.String())
	}

	storeContractBytecode := "0x608060405234801561001057600080fd5b5061020b806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c80632e64cec11461003b57806370ef6e0b14610059575b600080fd5b610043610075565b60405161005091906100ab565b60405180910390f35b610073600480360381019061006e9190610161565b61007e565b005b60008054905090565b836000819055508260008190555050505050565b6000819050919050565b6100a581610092565b82525050565b60006020820190506100c0600083018461009c565b92915050565b600080fd5b600080fd5b6100d981610092565b81146100e457600080fd5b50565b6000813590506100f6816100d0565b92915050565b600080fd5b600080fd5b600080fd5b60008083601f840112610121576101206100fc565b5b8235905067ffffffffffffffff81111561013e5761013d610101565b5b60208301915083600182028301111561015a57610159610106565b5b9250929050565b6000806000806060858703121561017b5761017a6100c6565b5b6000610189878288016100e7565b945050602061019a878288016100e7565b935050604085013567ffffffffffffffff8111156101bb576101ba6100cb565b5b6101c78782880161010b565b92509250509295919450925056fea2646970667358221220eda68578fb741c32f26000b6c0273945f8322dd35f536c918e3d5a6193aaf62564736f6c63430008120033"
	nonce, err := authClient.NonceAt(ctx, nil)
	assert.Nil(t, err)

	l2Wallet.SetNonce(nonce)
	estimatedTx := authClient.EstimateGasAndGasPrice(&types.LegacyTx{
		Nonce:    l2Wallet.GetNonceAndIncrement(),
		Value:    big.NewInt(100),
		Gas:      uint64(1_000_000),
		GasPrice: gethcommon.Big1,
		Data:     gethcommon.FromHex(storeContractBytecode),
	})
	assert.Nil(t, err)

	signedTx, err := l2Wallet.SignTransaction(estimatedTx)
	assert.Nil(t, err)

	_, err = authClient.SendTransaction(ctx, signedTx)
	assert.Nil(t, err)

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

func TestL2IssueTxWaitReceipt(t *testing.T) {
	// test is skipped by default to avoid breaking CI - set env flag in `Run Configurations` to run it in IDE
	if os.Getenv(_IDEFlag) == "" {
		t.Skipf("set flag %s to run this test in the IDE", _IDEFlag)
	}

	ctx := context.Background()

	vk, err := rpc.GenerateAndSignViewingKey(l2Wallet)
	assert.Nil(t, err)
	client, err := rpc.NewEncNetworkClient(fmt.Sprintf("ws://%s:%d", l2Host, l2Port), vk, gethlog.New())
	assert.Nil(t, err)
	authClient := obsclient.NewAuthObsClient(client)

	balance, err := authClient.BalanceAt(context.Background(), nil)
	assert.Nil(t, err)

	if balance.Cmp(big.NewInt(0)) <= 0 {
		t.Errorf("not enough balance: has %s has %s obx", l2Wallet.Address().Hex(), balance.String())
	}

	toAddr := datagenerator.RandomAddress()
	nonce, err := authClient.NonceAt(ctx, nil)
	assert.Nil(t, err)

	l2Wallet.SetNonce(nonce)
	estimatedTx := authClient.EstimateGasAndGasPrice(&types.LegacyTx{
		Nonce:    l2Wallet.GetNonceAndIncrement(),
		To:       &toAddr,
		Value:    big.NewInt(100),
		Gas:      uint64(1_000_000),
		GasPrice: gethcommon.Big1,
	})
	assert.Nil(t, err)

	signedTx, err := l2Wallet.SignTransaction(estimatedTx)
	assert.Nil(t, err)

	err = authClient.SendTransaction(ctx, signedTx)
	assert.Nil(t, err)

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
