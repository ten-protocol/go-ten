package storage

import (
	"context"
	"encoding/binary"
	"math/big"
	"time"

	"github.com/ten-protocol/go-ten/go/common/errutil"

	"github.com/ten-protocol/go-ten/go/enclave/storage/enclavedb"

	"github.com/ten-protocol/go-ten/go/enclave/core"

	"github.com/dgraph-io/ristretto/v2"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/log"
)

const (
	cacheCost = 1
)

type CacheService struct {
	// cache for the immutable blocks headers
	blockCache *ristretto.Cache[[]byte, *types.Header]

	// stores batches using the sequence number as key
	batchCacheBySeqNo *ristretto.Cache[uint64, *common.BatchHeader]

	// mapping between the hash and the sequence number
	// note:  to fetch a batch by hash will require 2 cache hits
	seqCacheByHash *ristretto.Cache[[]byte, *big.Int]

	// mapping between the height and the sequence number
	// note: to fetch a batch by height will require 2 cache hits
	seqCacheByHeight *ristretto.Cache[uint64, *big.Int]

	// store the converted ethereum header which is passed to the evm
	convertedGethHeaderCache *ristretto.Cache[[]byte, *types.Header]

	// batch hash - geth converted hash
	convertedHashCache *ristretto.Cache[[]byte, *gethcommon.Hash]

	// from address ( either eoa or contract) to the id of the db entry
	eoaCache             *ristretto.Cache[[]byte, *uint64]
	contractAddressCache *ristretto.Cache[[]byte, *enclavedb.Contract]
	eventTopicCache      *ristretto.Cache[[]byte, *enclavedb.EventTopic]

	// from contract_address||event_sig to the event_type object
	eventTypeCache *ristretto.Cache[[]byte, *enclavedb.EventType]

	// store the last few batches together with the content
	lastBatchesCache *ristretto.Cache[uint64, *core.Batch]

	// store all recent receipts in a cache
	// together with the sender - and for each log whether it is visible by the sender
	// only sender can view configured
	receiptCache *ristretto.Cache[[]byte, *CachedReceipt]

	// store the enclaves from the network
	attestedEnclavesCache *ristretto.Cache[[]byte, *AttestedEnclave]

	// cache for sequencer enclave IDs
	sequencerIDsCache []common.EnclaveID

	logger gethlog.Logger
}

func NewCacheService(logger gethlog.Logger, testMode bool) *CacheService {
	nrL1Blocks := 100        // ~200k
	nrBatches := 10_000      // ~25M
	nrConvertedEth := 10_000 // ~25M

	nrEventTypes := 10_000        // ~2M
	nrEOA := 100_000              // ~1M
	nrContractAddresses := 10_000 // ~1M

	nrBatchesWithContent := 50 // ~100M

	nrEnclaves := 20

	nrReceipts := 15_000 // ~100M
	if testMode {
		nrReceipts = 500
	}

	return &CacheService{
		blockCache: newCache[[]byte, *types.Header](logger, nrL1Blocks),

		batchCacheBySeqNo: newCache[uint64, *common.BatchHeader](logger, nrBatches),
		seqCacheByHash:    newCache[[]byte, *big.Int](logger, nrBatches),
		seqCacheByHeight:  newCache[uint64, *big.Int](logger, nrBatches),

		convertedGethHeaderCache: newCache[[]byte, *types.Header](logger, nrConvertedEth),
		convertedHashCache:       newCache[[]byte, *gethcommon.Hash](logger, nrConvertedEth),

		eoaCache:             newCache[[]byte, *uint64](logger, nrEOA),
		contractAddressCache: newCache[[]byte, *enclavedb.Contract](logger, nrContractAddresses),
		eventTypeCache:       newCache[[]byte, *enclavedb.EventType](logger, nrEventTypes),
		eventTopicCache:      newCache[[]byte, *enclavedb.EventTopic](logger, nrEventTypes),

		receiptCache:          newCache[[]byte, *CachedReceipt](logger, nrReceipts),
		attestedEnclavesCache: newCache[[]byte, *AttestedEnclave](logger, nrEnclaves),

		// cache the latest received batches to avoid a lookup when streaming it back to the host after processing
		lastBatchesCache: newCache[uint64, *core.Batch](logger, nrBatchesWithContent),

		sequencerIDsCache: make([]common.EnclaveID, 0),

		logger: logger,
	}
}

func newCache[K ristretto.Key, V any](logger gethlog.Logger, nrElem int) *ristretto.Cache[K, V] {
	ristrettoCache, err := ristretto.NewCache(&ristretto.Config[K, V]{
		NumCounters: int64(10 * nrElem), // 10 times the expected elements
		MaxCost:     int64(nrElem),      // calculate the max cost
		BufferItems: 64,                 // number of keys per Get buffer.
		Metrics:     true,
	})
	if err != nil {
		logger.Crit("Could not initialise ristretto cache", log.ErrKey, err)
	}
	return ristrettoCache
}

func (cs *CacheService) CacheConvertedHash(ctx context.Context, batchHash, convertedHash gethcommon.Hash) {
	cacheValue(ctx, cs.convertedHashCache, cs.logger, batchHash.Bytes(), &convertedHash)
}

func (cs *CacheService) CacheBlock(ctx context.Context, b *types.Header) {
	cacheValue(ctx, cs.blockCache, cs.logger, b.Hash().Bytes(), b)
}

func (cs *CacheService) CacheReceipt(_ context.Context, r *CachedReceipt) {
	// keep the receipts in cache for 100 seconds - 100 batches
	cs.receiptCache.SetWithTTL(r.Receipt.TxHash.Bytes(), r, cacheCost, 100*time.Second)
}

func (cs *CacheService) CacheBatch(ctx context.Context, batch *core.Batch) {
	cacheValue(ctx, cs.batchCacheBySeqNo, cs.logger, batch.SeqNo().Uint64(), batch.Header)
	cacheValue(ctx, cs.seqCacheByHash, cs.logger, batch.Hash().Bytes(), batch.SeqNo())
	// note: the key is (height+1), because for some reason it doesn't like a key of 0
	// should always contain the canonical batch because the cache is overwritten by each new batch after a reorg
	cacheValue(ctx, cs.seqCacheByHeight, cs.logger, batch.NumberU64()+1, batch.SeqNo())

	cacheValue(ctx, cs.lastBatchesCache, cs.logger, batch.SeqNo().Uint64(), batch)
}

func (cs *CacheService) ReadBlock(ctx context.Context, key gethcommon.Hash, onCacheMiss func() (*types.Header, error)) (*types.Header, error) {
	return getCachedValue(ctx, cs.blockCache, cs.logger, key.Bytes(), onCacheMiss, true)
}

func (cs *CacheService) ReadBatchSeqByHash(ctx context.Context, hash common.L2BatchHash, onCacheMiss func() (*big.Int, error)) (*big.Int, error) {
	return getCachedValue(ctx, cs.seqCacheByHash, cs.logger, hash.Bytes(), onCacheMiss, true)
}

func (cs *CacheService) ReadBatchSeqByHeight(ctx context.Context, height uint64, onCacheMiss func() (*big.Int, error)) (*big.Int, error) {
	// the key is (height+1), because for some reason it doesn't like a key of 0
	return getCachedValue(ctx, cs.seqCacheByHeight, cs.logger, height+1, onCacheMiss, true)
}

func (cs *CacheService) ReadConvertedHash(ctx context.Context, hash common.L2BatchHash, onCacheMiss func() (*gethcommon.Hash, error)) (*gethcommon.Hash, error) {
	return getCachedValue(ctx, cs.convertedHashCache, cs.logger, hash.Bytes(), onCacheMiss, true)
}

func (cs *CacheService) ReadBatchHeader(ctx context.Context, seqNum uint64, onCacheMiss func() (*common.BatchHeader, error)) (*common.BatchHeader, error) {
	return getCachedValue(ctx, cs.batchCacheBySeqNo, cs.logger, seqNum, onCacheMiss, true)
}

func (cs *CacheService) ReadBatch(ctx context.Context, seqNum uint64, onCacheMiss func() (*core.Batch, error)) (*core.Batch, error) {
	return getCachedValue(ctx, cs.lastBatchesCache, cs.logger, seqNum, onCacheMiss, true)
}

func (cs *CacheService) ReadEOA(ctx context.Context, addr gethcommon.Address, onCacheMiss func() (*uint64, error)) (*uint64, error) {
	return getCachedValue(ctx, cs.eoaCache, cs.logger, addr.Bytes(), onCacheMiss, true)
}

func (cs *CacheService) ReadContractAddr(ctx context.Context, addr gethcommon.Address, onCacheMiss func() (*enclavedb.Contract, error)) (*enclavedb.Contract, error) {
	return getCachedValue(ctx, cs.contractAddressCache, cs.logger, addr.Bytes(), onCacheMiss, true)
}

func (cs *CacheService) ReadEventTopic(ctx context.Context, topic []byte, eventTypeId uint64, onCacheMiss func() (*enclavedb.EventTopic, error)) (*enclavedb.EventTopic, error) {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, eventTypeId)
	key := append(topic, b...)
	return getCachedValue(ctx, cs.eventTopicCache, cs.logger, key, onCacheMiss, true)
}

type CachedReceipt struct {
	Receipt *types.Receipt
	From    *gethcommon.Address
	To      *gethcommon.Address
}

func (cs *CacheService) ReadReceipt(ctx context.Context, txHash gethcommon.Hash) (*CachedReceipt, error) {
	return getCachedValue(ctx, cs.receiptCache, cs.logger, txHash.Bytes(), nil, false)
}

func (cs *CacheService) DelReceipt(_ context.Context, txHash gethcommon.Hash) error {
	cs.receiptCache.Del(txHash.Bytes())
	return nil
}

func (cs *CacheService) ReadEventType(ctx context.Context, contractAddress gethcommon.Address, eventSignature gethcommon.Hash, onCacheMiss func() (*enclavedb.EventType, error)) (*enclavedb.EventType, error) {
	key := make([]byte, 0)
	key = append(key, contractAddress.Bytes()...)
	key = append(key, eventSignature.Bytes()...)
	return getCachedValue(ctx, cs.eventTypeCache, cs.logger, key, onCacheMiss, true)
}

func (cs *CacheService) ReadConvertedHeader(ctx context.Context, batchHash common.L2BatchHash, onCacheMiss func() (*types.Header, error)) (*types.Header, error) {
	return getCachedValue(ctx, cs.convertedGethHeaderCache, cs.logger, batchHash.Bytes(), onCacheMiss, true)
}

func (cs *CacheService) ReadEnclavePubKey(ctx context.Context, enclaveId common.EnclaveID, onCacheMiss func() (*AttestedEnclave, error)) (*AttestedEnclave, error) {
	return getCachedValue(ctx, cs.attestedEnclavesCache, cs.logger, enclaveId.Bytes(), onCacheMiss, true)
}

func (cs *CacheService) UpdateEnclaveNodeType(ctx context.Context, enclaveId common.EnclaveID, nodeType common.NodeType) {
	enclave, err := cs.ReadEnclavePubKey(ctx, enclaveId, nil)
	if err != nil {
		cs.logger.Debug("No cache entry found to update", log.ErrKey, err)
		return
	}
	enclave.Type = nodeType
	cacheValue(ctx, cs.attestedEnclavesCache, cs.logger, enclaveId.Bytes(), enclave)
}

func (cs *CacheService) CacheSequencerIDs(_ context.Context, sequencerIDs []common.EnclaveID) {
	cs.sequencerIDsCache = sequencerIDs
}

func (cs *CacheService) ReadSequencerIDs(_ context.Context, onCacheMiss func() ([]common.EnclaveID, error)) ([]common.EnclaveID, error) {
	if len(cs.sequencerIDsCache) == 0 {
		var err error
		cs.sequencerIDsCache, err = onCacheMiss()
		if err != nil {
			return nil, err
		}
	}
	return cs.sequencerIDsCache, nil
}

// getCachedValue - returns the cached value for the provided key. If the key is not found, then invoke the 'onCacheMiss' function
// which returns the value, and cache it
func getCachedValue[K ristretto.Key, V any](ctx context.Context, cache *ristretto.Cache[K, V], logger gethlog.Logger, key K, onCacheMiss func() (V, error), cacheIfMissing bool) (V, error) {
	value, found := cache.Get(key)
	if found {
		return value, nil
	}
	if onCacheMiss == nil {
		return value, errutil.ErrNotFound
	}

	v, err := onCacheMiss()
	if err != nil {
		return v, err
	}
	if cacheIfMissing {
		cacheValue(ctx, cache, logger, key, v)
	}
	return v, nil
}

func cacheValue[K ristretto.Key, V any](_ context.Context, cache *ristretto.Cache[K, V], logger gethlog.Logger, key K, v V) {
	added := cache.Set(key, v, cacheCost)
	if !added {
		logger.Debug("Did not store value in cache")
	}
}

//func logCacheMetrics(c *ristretto.Cache, name string, logger gethlog.Logger) {
//	metrics := c.Metrics
//	logger.Info(fmt.Sprintf("Cache %s metrics: Hits: %d, Misses: %d, Cost Added: %d",
//		name, metrics.Hits(), metrics.Misses(), metrics.CostAdded()))
//}
