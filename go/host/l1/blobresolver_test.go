package l1

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/stretchr/testify/require"
	"github.com/ten-protocol/go-ten/go/ethadapter"
)

const (
	vHash1 = "0x01c509cb5b108f8edf2333fa35a93acc5e4b24808179b0604f1945a8c7e98a3a"
	vHash2 = "0x013cc291176322a963ca75c09d27e31ab7690d88afdefb87c436815baeeb1078"
)

func TestBlobResolver(t *testing.T) {
	t.Skipf("TODO need to work out new params that will work with a new fallback provider")
	logger := log.New()
    beaconClient := ethadapter.NewBeaconHTTPClient(new(http.Client), logger, "https://docs-demo.quiknode.pro/")
    fallback := ethadapter.NewArchivalHTTPClient(new(http.Client), logger, "https://api.ethernow.xyz")
    blobResolver := NewBlobResolver(ethadapter.NewL1BeaconClient(beaconClient, logger, fallback), logger)

	// this will convert to slot 5 which will return 404 from the quicknode api, causing the fallback to be used
	b := &types.Header{
		Time: 1742476343,
	}

	blobs, err := blobResolver.FetchBlobs(context.Background(), b, []gethcommon.Hash{gethcommon.HexToHash(vHash1), gethcommon.HexToHash(vHash2)})
	require.NoError(t, err)
	require.Len(t, blobs, 2)
}

// TestSepoliaBlobResolver checks the public node sepolia beacon APIs work as expected
func TestSepoliaBlobResolver(t *testing.T) {
	t.Skipf("Test will occasionally not pass due to the time window landing on a block with no blobs")
	logger := log.New()
	// l1_beacon_url for sepolia
    beaconClient := ethadapter.NewBeaconHTTPClient(new(http.Client), logger, "https://ethereum-sepolia-beacon-api.publicnode.com")
	// l1_blob_archive_url for sepolia
    fallback := ethadapter.NewBeaconHTTPClient(new(http.Client), logger, "https://eth-beacon-chain-sepolia.drpc.org/rest/")
    blobResolver := NewBlobResolver(ethadapter.NewL1BeaconClient(beaconClient, logger, fallback), logger)

	// this is a moving point in time so we can't compare hashes or be certain there will be blobs in the block
	// create block with timestamp 30 days ago relative to current time
	historicalBlock := &types.Header{
		Time: uint64(time.Now().Unix()) - (30 * 24 * 60 * 60), // 30 days in seconds
	}

	_, err := blobResolver.FetchBlobs(context.Background(), historicalBlock, []gethcommon.Hash{})
	require.NoError(t, err)
}

// TestBlobResolverErrorHandling tests the retry strategy with different error conditions
func TestBlobResolverErrorHandling(t *testing.T) {
	t.Skipf("Test takes a long time to complete and is non-deterministic.")

    beaconClient := ethadapter.NewBeaconHTTPClient(new(http.Client), log.New(), "https://ethereum-sepolia-beacon-api.publicnode.com")

    blobResolver := NewBlobResolver(ethadapter.NewL1BeaconClient(beaconClient, log.New()), log.New())

	// test different slot numbers to potentially trigger different error states
	testSlots := []uint64{7861391, 7861392, 7861393, 7861394, 7861395}

	blocks := make([]*types.Header, len(testSlots))
	for i, slot := range testSlots {
		// approximate slot to timestamp ssuming 12 seconds per slot and genesis time around 1606824023
		genesisTime := uint64(1606824023)
		secondsPerSlot := uint64(12)
		timestamp := genesisTime + (slot * secondsPerSlot)

		blocks[i] = &types.Header{
			Time: timestamp,
		}
	}

	// Test 1: Concurrent requests to same slot (potential rate limiting)
	t.Run("ConcurrentSameSlot", func(t *testing.T) {
		block := blocks[0] // Use first slot
		var wg sync.WaitGroup
		results := make([]error, 10)

		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()
				_, err := blobResolver.FetchBlobs(context.Background(), block, []gethcommon.Hash{})
				results[index] = err
			}(i)
		}

		wg.Wait()

		// Log results to see what errors we get
		for i, err := range results {
			if err != nil {
				t.Logf("Concurrent request %d failed: %v", i, err)
			}
		}
	})

	// Test 2: Sequential requests to different slots (potential overload)
	t.Run("SequentialDifferentSlots", func(t *testing.T) {
		for i, block := range blocks {
			t.Logf("Testing slot %d (block time: %d)", testSlots[i], block.Time)
			_, err := blobResolver.FetchBlobs(context.Background(), block, []gethcommon.Hash{})
			if err != nil {
				t.Logf("Slot %d failed: %v", testSlots[i], err)
			} else {
				t.Logf("Slot %d succeeded", testSlots[i])
			}
			// Small delay between requests
			time.Sleep(100 * time.Millisecond)
		}
	})

	// Test 3: Rapid fire requests (potential timeout/overload)
	t.Run("RapidFireRequests", func(t *testing.T) {
		block := blocks[0]
		for i := 0; i < 5; i++ {
			_, err := blobResolver.FetchBlobs(context.Background(), block, []gethcommon.Hash{})
			if err != nil {
				t.Logf("Rapid fire request %d failed: %v", i, err)
			}
			// Very small delay
			time.Sleep(10 * time.Millisecond)
		}
	})

	// Test 4: Invalid slot numbers (potential 404/not found errors)
	t.Run("InvalidSlots", func(t *testing.T) {
		invalidSlots := []uint64{999999999, 1000000000, 1000000001}

		for _, slot := range invalidSlots {
			// Convert to timestamp
			genesisTime := uint64(1606824023)
			secondsPerSlot := uint64(12)
			timestamp := genesisTime + (slot * secondsPerSlot)

			block := &types.Header{Time: timestamp}
			_, err := blobResolver.FetchBlobs(context.Background(), block, []gethcommon.Hash{})
			if err != nil {
				t.Logf("Invalid slot %d failed as expected: %v", slot, err)
			} else {
				t.Logf("Invalid slot %d unexpectedly succeeded", slot)
			}
		}
	})

	// Test 5: Mixed concurrent requests to different slots
	t.Run("MixedConcurrentRequests", func(t *testing.T) {
		var wg sync.WaitGroup
		results := make([]error, len(blocks)*3) // 3 requests per block

		for i, block := range blocks {
			for j := 0; j < 3; j++ {
				wg.Add(1)
				go func(blockIndex, requestIndex int, currentBlock *types.Header) {
					defer wg.Done()
					_, err := blobResolver.FetchBlobs(context.Background(), currentBlock, []gethcommon.Hash{})
					results[blockIndex*3+requestIndex] = err
				}(i, j, block)
			}
		}

		wg.Wait()

		// Log results
		for i, err := range results {
			blockIndex := i / 3
			requestIndex := i % 3
			if err != nil {
				t.Logf("Mixed concurrent request block %d request %d failed: %v", blockIndex, requestIndex, err)
			}
		}
	})
}

// TestBlobResolverRetryStrategy tests the error classification and retry strategy logic
func TestBlobResolverRetryStrategy(t *testing.T) {
	// Test error classification
	tests := []struct {
		name             string
		errorMessage     string
		expectedStrategy RetryStrategy
	}{
		{
			name:             "Cloudflare 524 timeout",
			errorMessage:     "failed to get blob sidecars for Block Header 0xea18449cb5d47fa7006f9c020a8074174e4383a490ab54311e0e5601d0963665: failed to fetch blob sidecars for slot 7859641: The connection to the origin web server was made, but the origin web server timed out before responding. To resolve, please work with your hosting provider or web development team to free up resources for your database or overloaded application. Additional troubleshooting information here. Cloudflare Ray ID: 9507ac702cad9d0f",
			expectedStrategy: RetryStrategyTransient,
		},
		{
			name:             "Method not available on freetier",
			errorMessage:     "failed request with status 400: {\"message\":\"method is not available on freetier\",\"code\":35}",
			expectedStrategy: RetryStrategyRateLimit,
		},
		{
			name:             "Rate limit error",
			errorMessage:     "rate limit exceeded: too many requests",
			expectedStrategy: RetryStrategyRateLimit,
		},
		{
			name:             "Connection timeout",
			errorMessage:     "connection timeout after 30s",
			expectedStrategy: RetryStrategyTransient,
		},
		{
			name:             "Generic error",
			errorMessage:     "some other error occurred",
			expectedStrategy: RetryStrategyStandard,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := errors.New(tt.errorMessage)
			strategy := classifyError(err)
			require.Equal(t, tt.expectedStrategy, strategy, fmt.Sprintf("Error classification mismatch for: %s", tt.name))
		})
	}
}
