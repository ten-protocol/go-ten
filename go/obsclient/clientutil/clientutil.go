package clientutil

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/go-obscuro/go/obsclient"
)

var defaultTimeoutInterval = 1 * time.Second

func AwaitTransactionReceipt(ctx context.Context, client *obsclient.AuthObsClient, hash common.Hash, timeout time.Duration) (*types.Receipt, error) {
	timeoutStrategy := NewTimeoutStrategy(timeout, defaultTimeoutInterval)
	return AwaitTransactionReceiptWithRetryStrategy(ctx, client, hash, timeoutStrategy)
}

func AwaitTransactionReceiptWithRetryStrategy(ctx context.Context, client *obsclient.AuthObsClient, hash common.Hash, waitingStrat retryStrategy) (*types.Receipt, error) {
	waitingStrat.Reset()
	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("context cancelled before receipt was received")
		case <-time.After(waitingStrat.WaitInterval()):
			fmt.Println("attempting receipt request " + waitingStrat.Summary())
			receipt, err := client.TransactionReceipt(ctx, hash)
			if err == nil {
				return receipt, nil
			}
			if waitingStrat.Done() {
				return nil, fmt.Errorf("receipt not found - %s - %w", waitingStrat.Summary(), err)
			}
		}
	}
}
