package l1

import (
	"context"
	"fmt"
	"strings"
	"time"

	gethlog "github.com/ethereum/go-ethereum/log"

	"github.com/ten-protocol/go-ten/go/common/retry"

	gethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"

	"github.com/ten-protocol/go-ten/go/ethadapter"
)

// RetryStrategy represents the type of retry strategy to use
type RetryStrategy string

const (
	RetryStrategyNone      RetryStrategy = "none"
	RetryStrategyTransient RetryStrategy = "transient"
	RetryStrategyRateLimit RetryStrategy = "rate_limit"
	RetryStrategyStandard  RetryStrategy = "standard"
)

var (
	_maxWaitForBlobs = 60 * time.Second
	_maxRetries      = 5
)

// BlobResolver is an interface for fetching blobs
type BlobResolver interface {
	// FetchBlobs Fetches the blob data using beacon chain APIs
	FetchBlobs(ctx context.Context, b *types.Header, hashes []gethcommon.Hash) ([]*kzg4844.Blob, error)
	// StoreBlobs is used to store blobs for the in-memory testing nodes
	StoreBlobs(slot uint64, blobs []*kzg4844.Blob) error
}

type beaconBlobResolver struct {
	beaconClient *ethadapter.L1BeaconClient
	logger       gethlog.Logger
}

func NewBlobResolver(beaconClient *ethadapter.L1BeaconClient, logger gethlog.Logger) BlobResolver {
	return &beaconBlobResolver{beaconClient: beaconClient, logger: logger}
}

func (r *beaconBlobResolver) FetchBlobs(ctx context.Context, b *types.Header, hashes []gethcommon.Hash) ([]*kzg4844.Blob, error) {
	// try fetching once to get the initial error and set appropriate retry strategy
	blobs, initialErr := r.beaconClient.FetchBlobs(ctx, b, hashes)
	if initialErr == nil {
		return blobs, nil
	}
	err := retry.DoWithCount(func(retryNum int) error {
		var fetchErr error
		blobs, fetchErr = r.beaconClient.FetchBlobs(ctx, b, hashes)

		if fetchErr == nil {
			return nil
		}

		retryStrategy := classifyError(fetchErr)

		switch retryStrategy {
		case RetryStrategyTransient:
			r.logger.Warn("Transient error while fetching blobs, will retry with backoff",
				"error", fetchErr, "retryNum", retryNum, "blockHash", b.Hash().Hex())
		case RetryStrategyRateLimit:
			r.logger.Warn("Rate limit or method unavailable error, will retry with longer intervals",
				"error", fetchErr, "retryNum", retryNum, "blockHash", b.Hash().Hex())
		default:
			r.logger.Warn("Error while fetching blobs, will retry",
				"error", fetchErr, "retryNum", retryNum, "blockHash", b.Hash().Hex())
		}

		return fetchErr
	}, r.getRetryStrategy(initialErr))

	if err != nil {
		return nil, fmt.Errorf("failed to fetch blobs after retries: %w", err)
	}

	return blobs, nil
}

// classifyError determines the retry strategy based on error type
func classifyError(err error) RetryStrategy {
	if err == nil {
		return RetryStrategyNone
	}

	errStr := strings.ToLower(err.Error())

	// transient error - use exponential backoff
	transientPatterns := []string{
		"timeout",
		"connection refused",
		"connection reset",
		"network is unreachable",
		"no route to host",
		"temporary failure",
		"server overloaded",
		"cloudflare",
		"524", // cloudflare timeout error
		"503",
		"502",
		"504",
		"origin web server timed out",
		"overloaded background task",
		"overloaded application",
		"stressing the resources",
		"free up resources",
	}

	for _, pattern := range transientPatterns {
		if strings.Contains(errStr, pattern) {
			return RetryStrategyTransient
		}
	}

	// rate limit and method unavailable errors - use longer intervals
	rateLimitPatterns := []string{
		"rate limit",
		"too many requests",
		"method is not available",
		"freetier",
	}

	for _, pattern := range rateLimitPatterns {
		if strings.Contains(errStr, pattern) {
			return RetryStrategyRateLimit
		}
	}

	return RetryStrategyStandard
}

func (r *beaconBlobResolver) getRetryStrategy(err error) retry.Strategy {
	retryStrategy := classifyError(err)

	switch retryStrategy {
	case RetryStrategyTransient:
		// exponential back off for transient errors
		return retry.NewDoublingBackoffStrategy(2*time.Second, uint64(_maxRetries))
	case RetryStrategyRateLimit:
		// longer intervals for rate limits and method unavailable errors
		return retry.NewTimeoutStrategy(_maxWaitForBlobs, 10*time.Second)
	default:
		// standard timeout strategy with fixed intervals
		return retry.NewTimeoutStrategy(_maxWaitForBlobs, 2*time.Second)
	}
}

func (r *beaconBlobResolver) StoreBlobs(_ uint64, _ []*kzg4844.Blob) error {
	panic("provided by the ethereum consensus layer")
}
