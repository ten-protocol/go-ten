package common

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"sync"
	"time"

	"github.com/ten-protocol/go-ten/integration/common/testlog"

	"github.com/ethereum/go-ethereum"

	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ten-protocol/go-ten/go/common/errutil"
	"github.com/ten-protocol/go-ten/go/common/retry"
	"github.com/ten-protocol/go-ten/go/obsclient"

	"github.com/ten-protocol/go-ten/go/wallet"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

var _awaitReceiptPollingInterval = 200 * time.Millisecond

func RndBtw(min uint64, max uint64) uint64 {
	if min >= max {
		panic(fmt.Sprintf("RndBtw requires min (%d) to be greater than max (%d)", min, max))
	}
	return uint64(rand.Int63n(int64(max-min))) + min //nolint:gosec
}

func RndBtwTime(min time.Duration, max time.Duration) time.Duration {
	if min <= 0 || max <= 0 {
		panic(fmt.Sprintf("invalid durations min=%s max=%s", min, max))
	}
	return time.Duration(RndBtw(uint64(min.Nanoseconds()), uint64(max.Nanoseconds()))) * time.Nanosecond
}

// AwaitReceipt blocks until the receipt for the transaction with the given hash has been received. Errors if the
// transaction is unsuccessful or times out.
func AwaitReceipt(ctx context.Context, client *obsclient.AuthObsClient, txHash gethcommon.Hash, timeout time.Duration) error {
	var receipt *types.Receipt
	var err error
	err = retry.Do(func() error {
		receipt, err = client.TransactionReceipt(ctx, txHash)
		if err != nil && !errors.Is(err, ethereum.NotFound) {
			return retry.FailFast(err)
		}
		return err
	}, retry.NewTimeoutStrategy(timeout, _awaitReceiptPollingInterval))
	if err != nil {
		return fmt.Errorf("could not retrieve receipt for transaction %s - %w", txHash.Hex(), err)
	}

	if receipt.Status == types.ReceiptStatusFailed {
		return fmt.Errorf("receipt had status failed for transaction %s", txHash.Hex())
	}

	return nil
}

func AwaitReceiptEth(ctx context.Context, client *ethclient.Client, txHash gethcommon.Hash, timeout time.Duration) (*types.Receipt, error) {
	var receipt *types.Receipt
	var err error
	startTime := time.Now()

	testlog.Logger().Info("Fetching receipt for tx: ", "tx", txHash.Hex())
	err = retry.Do(func() error {
		receipt, err = client.TransactionReceipt(ctx, txHash)
		if err != nil && !errors.Is(err, errutil.ErrNotFound) {
			// we only retry for a nil "not found" response. This is a different error, so we bail out of the retry loop
			return retry.FailFast(err)
		}
		testlog.Logger().Info("No tx receipt after: ", "time", time.Since(startTime))
		return err
	}, retry.NewTimeoutStrategy(timeout, _awaitReceiptPollingInterval))
	if err != nil {
		return nil, fmt.Errorf("could not retrieve receipt for transaction %s - %w", txHash.Hex(), err)
	}

	if receipt.Status == types.ReceiptStatusFailed {
		return nil, fmt.Errorf("receipt had status failed for transaction %s", txHash.Hex())
	}

	return receipt, nil
}

// PrefundWallets sends an amount `alloc` from the faucet wallet to each listed wallet.
// The transactions are sent with sequential nonces, starting with `startingNonce`.
func PrefundWallets(ctx context.Context, faucetWallet wallet.Wallet, faucetClient *obsclient.AuthObsClient, startingNonce uint64, wallets []wallet.Wallet, alloc *big.Int, timeout time.Duration) {
	// We send the transactions serially, so that we can precompute the nonces.
	txHashes := make([]gethcommon.Hash, len(wallets))
	for idx, w := range wallets {
		destAddr := w.Address()
		fmt.Printf("L2 prefund: %s\n", destAddr.Hex())
		txData := &types.LegacyTx{
			Nonce:    startingNonce + uint64(idx),
			Value:    alloc,
			Gas:      uint64(100_000),
			GasPrice: gethcommon.Big1,
			To:       &destAddr,
		}

		tx := faucetClient.EstimateGasAndGasPrice(txData) //nolint: contextcheck
		signedTx, err := faucetWallet.SignTransaction(tx)
		if err != nil {
			panic(err)
		}

		err = faucetClient.SendTransaction(ctx, signedTx)
		if err != nil {
			var txJSON []byte
			txJSON, _ = signedTx.MarshalJSON()
			panic(fmt.Sprintf("could not transfer from faucet for tx %s. Cause: %s", string(txJSON[:]), err))
		}

		txHashes[idx] = signedTx.Hash()
	}

	// Then we await the receipts in parallel.
	wg := sync.WaitGroup{}
	for _, txHash := range txHashes {
		wg.Add(1)
		go func(txHash gethcommon.Hash) {
			defer wg.Done()
			err := AwaitReceipt(ctx, faucetClient, txHash, timeout)
			if err != nil {
				panic(fmt.Sprintf("faucet transfer transaction %s unsuccessful. Cause: %s", txHash, err))
			}
		}(txHash)
	}
	wg.Wait()
}

func InteractWithSmartContract(client *ethclient.Client, wallet wallet.Wallet, contractAbi abi.ABI, methodName string, methodParam string, contractAddress gethcommon.Address) (*types.Receipt, error) {
	contractInteractionData, err := contractAbi.Pack(methodName, methodParam)
	if err != nil {
		return nil, err
	}

	price, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	interactionTx := types.LegacyTx{
		Nonce:    wallet.GetNonceAndIncrement(),
		To:       &contractAddress,
		Gas:      uint64(10_000_000),
		GasPrice: price,
		Data:     contractInteractionData,
	}

	gas, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
		From: wallet.Address(),
		To:   &contractAddress,
		Data: contractInteractionData,
	})
	if err != nil {
		return nil, err
	}
	interactionTx.Gas = gas

	signedTx, err := wallet.SignTransaction(&interactionTx)
	if err != nil {
		return nil, err
	}
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return nil, err
	}

	txReceipt, err := AwaitReceiptEth(context.Background(), client, signedTx.Hash(), 10*time.Second)
	if err != nil {
		return nil, err
	}

	return txReceipt, nil
}
