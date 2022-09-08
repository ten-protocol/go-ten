package clientutil

import (
	"context"
	"fmt"
	"github.com/obscuronet/go-obscuro/go/common/retry"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/go-obscuro/go/obsclient"
)

var defaultTimeoutInterval = 1 * time.Second

func AwaitTransactionReceipt(ctx context.Context, client *obsclient.AuthObsClient, txHash common.Hash, timeout time.Duration) (*types.Receipt, error) {
	timeoutStrategy := retry.NewTimeoutStrategy(timeout, defaultTimeoutInterval)
	return AwaitTransactionReceiptWithRetryStrategy(ctx, client, txHash, timeoutStrategy)
}

func AwaitTransactionReceiptWithRetryStrategy(ctx context.Context, client *obsclient.AuthObsClient, txHash common.Hash, retryStrategy retry.RetryStrategy) (*types.Receipt, error) {
	retryStrategy.Reset()
	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("context cancelled before receipt was received")
		case <-time.After(retryStrategy.NextRetryInterval()):
			receipt, err := client.TransactionReceipt(ctx, txHash)
			if err == nil {
				return receipt, nil
			}
			if retryStrategy.Done() {
				return nil, fmt.Errorf("receipt not found - %s - %w", retryStrategy.Summary(), err)
			}
		}
	}
}
