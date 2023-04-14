package manualtests

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/ethadapter"
	"github.com/obscuronet/go-obscuro/go/obsclient"
	"github.com/obscuronet/go-obscuro/go/wallet"
	"github.com/obscuronet/go-obscuro/integration/common/testlog"
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
	err = awaitL1Tx(ethClient, signedTx)
	assert.Nil(t, err)
	fmt.Printf("Successfully minted Tx: %s \n", signedTx.Hash().Hex())
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

	err = awaitL1Tx(ethClient, signedTx)
	assert.Nil(t, err)
	fmt.Printf("Successfully minted Tx: %s \n", signedTx.Hash().Hex())
}

func TestL2IssueContractInteractWaitReceipt(t *testing.T) {
	// test is skipped by default to avoid breaking CI - set env flag in `Run Configurations` to run it in IDE
	if os.Getenv(_IDEFlag) == "" {
		t.Skipf("set flag %s to run this test in the IDE", _IDEFlag)
	}

	authClient, err := obsclient.DialWithAuth(fmt.Sprintf("ws://%s:%d", l2Host, l2Port), l2Wallet, testlog.Logger())
	assert.Nil(t, err)

	// check if account has balance of the ops
	balance, err := authClient.BalanceAt(context.Background(), nil)
	assert.Nil(t, err)
	if balance.Cmp(big.NewInt(0)) <= 0 {
		t.Errorf("not enough balance: has %s has %s obx", l2Wallet.Address().Hex(), balance.String())
	}

	// update the nonce
	nonce, err := authClient.NonceAt(context.Background(), nil)
	assert.Nil(t, err)
	l2Wallet.SetNonce(nonce)

	// estimate and deploy contract
	storeContractBytecode := "0x608060405234801561001057600080fd5b5061020b806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c80632e64cec11461003b57806370ef6e0b14610059575b600080fd5b610043610075565b60405161005091906100ab565b60405180910390f35b610073600480360381019061006e9190610161565b61007e565b005b60008054905090565b836000819055508260008190555050505050565b6000819050919050565b6100a581610092565b82525050565b60006020820190506100c0600083018461009c565b92915050565b600080fd5b600080fd5b6100d981610092565b81146100e457600080fd5b50565b6000813590506100f6816100d0565b92915050565b600080fd5b600080fd5b600080fd5b60008083601f840112610121576101206100fc565b5b8235905067ffffffffffffffff81111561013e5761013d610101565b5b60208301915083600182028301111561015a57610159610106565b5b9250929050565b6000806000806060858703121561017b5761017a6100c6565b5b6000610189878288016100e7565b945050602061019a878288016100e7565b935050604085013567ffffffffffffffff8111156101bb576101ba6100cb565b5b6101c78782880161010b565b92509250509295919450925056fea2646970667358221220eda68578fb741c32f26000b6c0273945f8322dd35f536c918e3d5a6193aaf62564736f6c63430008120033"
	estimatedTx := authClient.EstimateGasAndGasPrice(&types.LegacyTx{
		Nonce:    l2Wallet.GetNonceAndIncrement(),
		Gas:      uint64(1_000_000),
		GasPrice: gethcommon.Big1,
		Data:     gethcommon.FromHex(storeContractBytecode),
	})
	assert.Nil(t, err)

	signedTx, err := l2Wallet.SignTransaction(estimatedTx)
	assert.Nil(t, err)

	_, err = authClient.SendTransaction(context.Background(), signedTx)
	assert.Nil(t, err)

	err = awaitL2Tx(authClient, signedTx)
	assert.Nil(t, err)
	fmt.Printf("Successfully minted Tx: %s \n", signedTx.Hash().Hex())
}

func TestL2IssueTxWaitReceipt(t *testing.T) {
	// test is skipped by default to avoid breaking CI - set env flag in `Run Configurations` to run it in IDE
	if os.Getenv(_IDEFlag) == "" {
		t.Skipf("set flag %s to run this test in the IDE", _IDEFlag)
	}

	authClient, err := obsclient.DialWithAuth(fmt.Sprintf("ws://%s:%d", l2Host, l2Port), l2Wallet, testlog.Logger())
	assert.Nil(t, err)

	// check if account has balance of the ops
	balance, err := authClient.BalanceAt(context.Background(), nil)
	assert.Nil(t, err)
	if balance.Cmp(big.NewInt(0)) <= 0 {
		t.Errorf("not enough balance: has %s has %s obx", l2Wallet.Address().Hex(), balance.String())
	}

	toAddr := datagenerator.RandomAddress()
	nonce, err := authClient.NonceAt(context.Background(), nil)
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

	//_, err = authClient.SendTransaction(context.Background(), signedTx)
	//assert.Nil(t, err)

	err = awaitL2Tx(authClient, signedTx)
	assert.Nil(t, err)
	fmt.Printf("Successfully minted Tx: %s \n", signedTx.Hash().Hex())
}
