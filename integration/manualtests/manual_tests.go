package manualtests

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/go-obscuro/go/ethadapter"
	"github.com/obscuronet/go-obscuro/go/obsclient"
)

func awaitL1Tx(ethClient ethadapter.EthClient, signedTx *types.Transaction) error {
	fmt.Printf("Checking for tx receipt for %s \n", signedTx.Hash())

	var err error
	var receipt *types.Receipt
	for start := time.Now(); time.Since(start) < time.Minute; time.Sleep(time.Second) {
		receipt, err = ethClient.TransactionReceipt(signedTx.Hash())
		if err == nil {
			break
		}
		if !errors.Is(err, ethereum.NotFound) {
			return err
		}
	}

	if receipt == nil {
		return fmt.Errorf("did not mine the transaction after %s seconds  - receipt: %+v", 30*time.Second, receipt)
	}
	if receipt.Status == 0 {
		return fmt.Errorf("tx Failed")
	}
	return nil
}

func awaitL2Tx(authClient *obsclient.AuthObsClient, signedTx *types.Transaction) error {
	fmt.Printf("Checking for tx receipt for %s \n", signedTx.Hash())

	var err error
	var receipt *types.Receipt
	for start := time.Now(); time.Since(start) < time.Minute; time.Sleep(time.Second) {
		receipt, err = authClient.TransactionReceipt(context.Background(), signedTx.Hash())
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
		return fmt.Errorf("did not mine the transaction after %s seconds  - receipt: %+v", 30*time.Second, receipt)
	}
	if receipt.Status == 0 {
		return fmt.Errorf("tx Failed")
	}
	return nil
}
