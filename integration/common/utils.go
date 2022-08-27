package common

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"sync"
	"time"

	"github.com/obscuronet/go-obscuro/go/obsclient"

	"github.com/obscuronet/go-obscuro/go/wallet"

	"github.com/obscuronet/go-obscuro/go/common"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/go-obscuro/go/rpcclientlib"
)

const (
	receiptTimeoutMillis = 30000 // The timeout in millis to wait for a receipt for a transaction.
)

func RndBtw(min uint64, max uint64) uint64 {
	if min >= max {
		panic(fmt.Sprintf("RndBtw requires min (%d) to be greater than max (%d)", min, max))
	}
	return uint64(rand.Int63n(int64(max-min))) + min //nolint:gosec
}

func RndBtwTime(min time.Duration, max time.Duration) time.Duration {
	if min <= 0 || max <= 0 {
		panic("invalid durations")
	}
	return time.Duration(RndBtw(uint64(min.Nanoseconds()), uint64(max.Nanoseconds()))) * time.Nanosecond
}

// AwaitReceipt blocks until the receipt for the transaction with the given hash has been received. Errors if the
// transaction is unsuccessful or times out.
func AwaitReceipt(ctx context.Context, client *obsclient.AuthObsClient, txHash gethcommon.Hash) error {
	counter := 0
	for {
		receipt, err := client.TransactionReceipt(ctx, txHash)
		if err != nil {
			if !errors.Is(err, rpcclientlib.ErrNilResponse) {
				return err
			}

			counter += 100
			if counter > receiptTimeoutMillis {
				return fmt.Errorf("could not retrieve transaction after timeout")
			}
			time.Sleep(100 * time.Millisecond)
			continue
		}

		if receipt.Status == types.ReceiptStatusFailed {
			return fmt.Errorf("receipt had status failed")
		}

		return nil
	}
}

// EncodeTx formats a transaction for sending to the enclave
func EncodeTx(tx *common.L2Tx) string {
	txBinary, err := tx.MarshalBinary()
	if err != nil {
		panic(err)
	}

	// We convert the transaction binary to the form expected for sending transactions via RPC.
	txBinaryHex := gethcommon.Bytes2Hex(txBinary)

	return "0x" + txBinaryHex
}

// PrefundWallets sends an amount `alloc` from the faucet wallet to each listed wallet.
// The transactions are sent with sequential nonces, starting with `startingNonce`.
func PrefundWallets(ctx context.Context, faucetWallet wallet.Wallet, faucetClient *obsclient.AuthObsClient, startingNonce uint64, wallets []wallet.Wallet, alloc *big.Int) {
	// We send the transactions serially, so that we can precompute the nonces.
	txHashes := make([]gethcommon.Hash, len(wallets))
	for idx, w := range wallets {
		destAddr := w.Address()
		tx := &types.LegacyTx{
			Nonce:    startingNonce + uint64(idx),
			Value:    alloc,
			Gas:      uint64(1_000_000),
			GasPrice: gethcommon.Big1,
			To:       &destAddr,
		}
		signedTx, err := faucetWallet.SignTransaction(tx)
		if err != nil {
			panic(err)
		}

		err = faucetClient.SendTransaction(ctx, signedTx)
		if err != nil {
			panic(fmt.Sprintf("could not transfer from faucet. Cause: %s", err))
		}

		txHashes[idx] = signedTx.Hash()
	}

	// Then we await the receipts in parallel.
	wg := sync.WaitGroup{}
	for _, txHash := range txHashes {
		wg.Add(1)
		go func(txHash gethcommon.Hash) {
			defer wg.Done()
			err := AwaitReceipt(ctx, faucetClient, txHash)
			if err != nil {
				panic(fmt.Sprintf("faucet transfer transaction failed. Cause: %s", err))
			}
		}(txHash)
	}
	wg.Wait()
}
