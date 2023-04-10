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

func TestL1IssueTxWaitReceipt(t *testing.T) {
	// test is skipped by default to avoid breaking CI - set env flag in `Run Configurations` to run it in IDE
	if os.Getenv(_IDEFlag) == "" {
		t.Skipf("set flag %s to run this test in the IDE", _IDEFlag)
	}

	w := wallet.NewInMemoryWalletFromConfig(
		"5d1cffab85ddad285de2485ff09339e66e1e0acbfb9960c0df8231a1deb4994a",
		1337,
		gethlog.New())
	host := "localhost"
	port := 38000

	var err error
	ethClient, err := ethadapter.NewEthClient(host, uint(port), 30*time.Second, common.L2Address{}, gethlog.New())
	assert.Nil(t, err)

	toAddr := datagenerator.RandomAddress()
	nonce, err := ethClient.Nonce(w.Address())
	assert.Nil(t, err)

	w.SetNonce(nonce)
	estimatedTx, err := ethClient.EstimateGasAndGasPrice(&types.LegacyTx{
		Nonce: w.GetNonceAndIncrement(),
		To:    &toAddr,
		Value: big.NewInt(100),
	}, w.Address())
	assert.Nil(t, err)

	signedTx, err := w.SignTransaction(estimatedTx)
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

func TestL2IssueTxWaitReceipt(t *testing.T) {
	// test is skipped by default to avoid breaking CI - set env flag in `Run Configurations` to run it in IDE
	if os.Getenv(_IDEFlag) == "" {
		t.Skipf("set flag %s to run this test in the IDE", _IDEFlag)
	}

	ctx := context.Background()
	w := wallet.NewInMemoryWalletFromConfig(
		"4bfe14725e685901c062ccd4e220c61cf9c189897b6c78bd18d7f51291b2b8f8",
		777,
		gethlog.New())
	host := "localhost"
	port := 13011

	vk, err := rpc.GenerateAndSignViewingKey(w)
	assert.Nil(t, err)
	client, err := rpc.NewEncNetworkClient(fmt.Sprintf("ws://%s:%d", host, port), vk, gethlog.New())
	assert.Nil(t, err)
	authClient := obsclient.NewAuthObsClient(client)

	balance, err := authClient.BalanceAt(context.Background(), nil)
	assert.Nil(t, err)

	if balance.Cmp(big.NewInt(0)) <= 0 {
		t.Errorf("not enough balance: has %s has %s obx", w.Address().Hex(), balance.String())
	}

	toAddr := datagenerator.RandomAddress()
	nonce, err := authClient.NonceAt(ctx, nil)
	assert.Nil(t, err)

	w.SetNonce(nonce)
	estimatedTx := authClient.EstimateGasAndGasPrice(&types.LegacyTx{
		Nonce:    w.GetNonceAndIncrement(),
		To:       &toAddr,
		Value:    big.NewInt(100),
		Gas:      uint64(1_000_000),
		GasPrice: gethcommon.Big1,
	})
	assert.Nil(t, err)

	signedTx, err := w.SignTransaction(estimatedTx)
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
}
