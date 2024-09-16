package storage

import (
	"context"
	"math/big"

	"github.com/ten-protocol/go-ten/go/enclave/storage/enclavedb"

	"github.com/eko/gocache/lib/v4/store"

	"github.com/ten-protocol/go-ten/go/enclave/core"

	"github.com/dgraph-io/ristretto"
	"github.com/eko/gocache/lib/v4/cache"
	ristretto_store "github.com/eko/gocache/store/ristretto/v4"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/log"
)

// approximate cost in bytes of the cached values
const (
	blockCost = 1024
	batchCost = 1024
	hashCost  = 32
	idCost    = 8
)

type CacheService struct {
	// cache for the immutable blocks and batches.
	// this avoids a trip to the database.
	blockCache *cache.Cache[*types.Header]

	// stores batches using the sequence number as key
	batchCacheBySeqNo *cache.Cache[*common.BatchHeader]

	// mapping between the hash and the sequence number
	// note:  to fetch a batch by hash will require 2 cache hits
	seqCacheByHash *cache.Cache[*big.Int]

	// mapping between the height and the sequence number
	// note: to fetch a batch by height will require 2 cache hits
	seqCacheByHeight *cache.Cache[*big.Int]

	// batch hash - geth converted hash
	convertedHashCache *cache.Cache[*gethcommon.Hash]

	// from address ( either eoa or contract) to the id of the db entry
	eoaCache             *cache.Cache[*uint64]
	contractAddressCache *cache.Cache[*enclavedb.Contract]

	// from contract_address||event_sig to the event_type (id, isLifecycle) object
	eventTypeCache *cache.Cache[*enclavedb.EventType]

	// store the converted ethereum header which is passed to the evm
	convertedGethHeaderCache *cache.Cache[*types.Header]

	logger gethlog.Logger
}

func NewCacheService(logger gethlog.Logger) *CacheService {
	// todo (tudor) figure out the config
	ristrettoCache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 100_000_000,        // 10 times the expected elements
		MaxCost:     1024 * 1024 * 1024, // allocate 1GB
		BufferItems: 64,                 // number of keys per Get buffer.
	})
	if err != nil {
		logger.Crit("Could not initialise ristretto cache", log.ErrKey, err)
	}
	ristrettoStore := ristretto_store.NewRistretto(ristrettoCache)
	return &CacheService{
		blockCache:               cache.New[*types.Header](ristrettoStore),
		batchCacheBySeqNo:        cache.New[*common.BatchHeader](ristrettoStore),
		seqCacheByHash:           cache.New[*big.Int](ristrettoStore),
		seqCacheByHeight:         cache.New[*big.Int](ristrettoStore),
		convertedHashCache:       cache.New[*gethcommon.Hash](ristrettoStore),
		eoaCache:                 cache.New[*uint64](ristrettoStore),
		contractAddressCache:     cache.New[*enclavedb.Contract](ristrettoStore),
		eventTypeCache:           cache.New[*enclavedb.EventType](ristrettoStore),
		convertedGethHeaderCache: cache.New[*types.Header](ristrettoStore),
		logger:                   logger,
	}
}

func (cs *CacheService) CacheBlock(ctx context.Context, b *types.Header) {
	cacheValue(ctx, cs.blockCache, cs.logger, b.Hash(), b, blockCost)
}

func (cs *CacheService) CacheBatch(ctx context.Context, batch *core.Batch) {
	cacheValue(ctx, cs.batchCacheBySeqNo, cs.logger, batch.SeqNo().Uint64(), batch.Header, batchCost)
	cacheValue(ctx, cs.seqCacheByHash, cs.logger, batch.Hash(), batch.SeqNo(), idCost)
	// note: the key is (height+1), because for some reason it doesn't like a key of 0
	// should always contain the canonical batch because the cache is overwritten by each new batch after a reorg
	cacheValue(ctx, cs.seqCacheByHeight, cs.logger, batch.NumberU64()+1, batch.SeqNo(), idCost)
}

func (cs *CacheService) ReadBlock(ctx context.Context, key gethcommon.Hash, onCacheMiss func(any) (*types.Header, error)) (*types.Header, error) {
	return getCachedValue(ctx, cs.blockCache, cs.logger, key, blockCost, onCacheMiss)
}

func (cs *CacheService) ReadBatchSeqByHash(ctx context.Context, hash common.L2BatchHash, onCacheMiss func(any) (*big.Int, error)) (*big.Int, error) {
	return getCachedValue(ctx, cs.seqCacheByHash, cs.logger, hash, idCost, onCacheMiss)
}

func (cs *CacheService) ReadBatchSeqByHeight(ctx context.Context, height uint64, onCacheMiss func(any) (*big.Int, error)) (*big.Int, error) {
	// the key is (height+1), because for some reason it doesn't like a key of 0
	return getCachedValue(ctx, cs.seqCacheByHeight, cs.logger, height+1, idCost, onCacheMiss)
}

func (cs *CacheService) ReadConvertedHash(ctx context.Context, hash common.L2BatchHash, onCacheMiss func(any) (*gethcommon.Hash, error)) (*gethcommon.Hash, error) {
	return getCachedValue(ctx, cs.convertedHashCache, cs.logger, hash, hashCost, onCacheMiss)
}

func (cs *CacheService) ReadBatch(ctx context.Context, seqNum uint64, onCacheMiss func(any) (*common.BatchHeader, error)) (*common.BatchHeader, error) {
	return getCachedValue(ctx, cs.batchCacheBySeqNo, cs.logger, seqNum, batchCost, onCacheMiss)
}

func (cs *CacheService) ReadEOA(ctx context.Context, addr gethcommon.Address, onCacheMiss func(any) (*uint64, error)) (*uint64, error) {
	return getCachedValue(ctx, cs.eoaCache, cs.logger, addr, idCost, onCacheMiss)
}

func (cs *CacheService) ReadContractAddr(ctx context.Context, addr gethcommon.Address, onCacheMiss func(any) (*enclavedb.Contract, error)) (*enclavedb.Contract, error) {
	return getCachedValue(ctx, cs.contractAddressCache, cs.logger, addr, idCost, onCacheMiss)
}

func (cs *CacheService) ReadEventType(ctx context.Context, contractAddress gethcommon.Address, eventSignature gethcommon.Hash, onCacheMiss func(any) (*enclavedb.EventType, error)) (*enclavedb.EventType, error) {
	key := make([]byte, 0)
	key = append(key, contractAddress.Bytes()...)
	key = append(key, eventSignature.Bytes()...)
	return getCachedValue(ctx, cs.eventTypeCache, cs.logger, key, idCost, onCacheMiss)
}

func (cs *CacheService) ReadConvertedHeader(ctx context.Context, batchHash common.L2BatchHash, onCacheMiss func(any) (*types.Header, error)) (*types.Header, error) {
	return getCachedValue(ctx, cs.convertedGethHeaderCache, cs.logger, batchHash, blockCost, onCacheMiss)
}

// getCachedValue - returns the cached value for the provided key. If the key is not found, then invoke the 'onCacheMiss' function
// which returns the value, and cache it
func getCachedValue[V any](ctx context.Context, cache *cache.Cache[*V], logger gethlog.Logger, key any, cost int64, onCacheMiss func(any) (*V, error)) (*V, error) {
	value, err := cache.Get(ctx, key)
	if err != nil || value == nil {
		// todo metrics for cache misses
		v, err := onCacheMiss(key)
		if err != nil {
			return v, err
		}
		if v == nil {
			logger.Crit("Returned a nil value from the onCacheMiss function. Should not happen.")
		}
		cacheValue(ctx, cache, logger, key, v, cost)
		return v, nil
	}

	return value, err
}

func cacheValue[V any](ctx context.Context, cache *cache.Cache[*V], logger gethlog.Logger, key any, v *V, cost int64) {
	if v == nil {
		return
	}
	err := cache.Set(ctx, key, v, store.WithCost(cost))
	if err != nil {
		logger.Error("Could not store value in cache", log.ErrKey, err)
	}
}
