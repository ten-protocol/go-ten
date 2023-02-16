package manualtests

import (
	"context"
	"errors"
	"fmt"
	"github.com/obscuronet/go-obscuro/go/obsclient"
	"github.com/obscuronet/go-obscuro/go/rpc"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/ethadapter"
	"github.com/obscuronet/go-obscuro/go/wallet"
	"github.com/obscuronet/go-obscuro/integration/datagenerator"
	"github.com/stretchr/testify/assert"

	gethlog "github.com/ethereum/go-ethereum/log"
)

func TestL1IssueTxWaitReceipt(t *testing.T) {
	t.Skip("manual tests should not be used for unit testing")

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
	//t.Skip("manual tests should not be used for unit testing")

	ctx := context.Background()
	w := wallet.NewInMemoryWalletFromConfig(
		"4bfe14725e685901c062ccd4e220c61cf9c189897b6c78bd18d7f51291b2b8f8",
		777,
		gethlog.New())
	host := "51.132.32.212"
	port := 13001

	vk, err := rpc.GenerateAndSignViewingKey(w)
	assert.Nil(t, err)
	client, err := rpc.NewEncNetworkClient(fmt.Sprintf("ws://%s:%d", host, port), vk, gethlog.New())
	assert.Nil(t, err)
	authClient := obsclient.NewAuthObsClient(client)

	toAddr := datagenerator.RandomAddress()
	nonce, err := authClient.NonceAt(ctx, nil)
	assert.Nil(t, err)

	w.SetNonce(nonce)
	estimatedTx := authClient.EstimateGasAndGasPrice(&types.LegacyTx{
		Nonce: w.GetNonceAndIncrement(),
		To:    &toAddr,
		Value: big.NewInt(100),
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
		//if !errors.Is(err, ethereum.NotFound) {
		//	t.Fatal(err)
		//}
		fmt.Printf("no tx receipt after %s - %s\n", time.Since(start), err)
	}

	if receipt == nil {
		t.Fatalf("Did not mine the transaction after %s seconds  - receipt: %+v", 30*time.Second, receipt)
	}
}
