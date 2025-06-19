package storage

import (
	"context"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ten-protocol/go-ten/go/common/errutil"

	"github.com/ten-protocol/go-ten/go/enclave/storage/enclavedb"

	"github.com/ten-protocol/go-ten/go/enclave/core"

	"github.com/TwiN/gocache/v2"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/log"
)

type CacheService struct {
	// cache for the immutable blocks headers
	blockCache *gocache.Cache

	// stores batches using the sequence number as key
	batchCacheBySeqNo *gocache.Cache

	// mapping between the hash and the sequence number
	// note:  to fetch a batch by hash will require 2 cache hits
	seqCacheByHash *gocache.Cache

	// mapping between the height and the sequence number
	// note: to fetch a batch by height will require 2 cache hits
	seqCacheByHeight *gocache.Cache

	// store the converted ethereum header which is passed to the evm
	convertedGethHeaderCache *gocache.Cache

	// batch hash - geth converted hash
	convertedHashCache *gocache.Cache

	// from address ( either eoa or contract) to the id of the db entry
	eoaCache             *gocache.Cache
	contractAddressCache *gocache.Cache
	eventTopicCache      *gocache.Cache

	// store the last few batches together with the content
	lastBatchesCache *gocache.Cache

	// store all recent receipts in a cache
	// together with the sender - and for each log whether it is visible by the sender
	// only sender can view configured
	receiptCache *gocache.Cache

	// store the enclaves from the network
	attestedEnclavesCache *gocache.Cache

	// cache for sequencer enclave IDs
	sequencerIDsCache []common.EnclaveID

	logger gethlog.Logger
}

func NewCacheService(logger gethlog.Logger, testMode bool) *CacheService {
	nrL1Blocks := 500        // ~1M - note that the value needs to be more than the moving average window of the gas oracle
	nrBatches := 10_000      // ~25M
	nrConvertedEth := 10_000 // ~25M

	nrEventTypes := 10_000        // ~2M
	nrEOA := 100_000              // ~1M
	nrContractAddresses := 10_000 // ~1M

	nrBatchesWithContent := 50 // ~100M

	nrEnclaves := 20

	nrReceipts := 15_000 // ~100M
	receiptsTimeout := 4 * time.Minute
	if testMode {
		nrReceipts = 2500
	}

	return &CacheService{
		blockCache: newFifoCache(logger, nrL1Blocks, gocache.NoExpiration),

		batchCacheBySeqNo: newFifoCache(logger, nrBatches, gocache.NoExpiration),
		seqCacheByHash:    newFifoCache(logger, nrBatches, gocache.NoExpiration),
		seqCacheByHeight:  newFifoCache(logger, nrBatches, gocache.NoExpiration),

		convertedGethHeaderCache: newFifoCache(logger, nrConvertedEth, gocache.NoExpiration),
		convertedHashCache:       newFifoCache(logger, nrConvertedEth, gocache.NoExpiration),

		eoaCache:             newFifoCache(logger, nrEOA, gocache.NoExpiration),
		contractAddressCache: newFifoCache(logger, nrContractAddresses, gocache.NoExpiration),
		eventTopicCache:      newFifoCache(logger, nrEventTypes, gocache.NoExpiration),

		receiptCache:          newFifoCache(logger, nrReceipts, receiptsTimeout),
		attestedEnclavesCache: newFifoCache(logger, nrEnclaves, gocache.NoExpiration),

		// cache the latest received batches to avoid a lookup when streaming it back to the host after processing
		lastBatchesCache: newFifoCache(logger, nrBatchesWithContent, gocache.NoExpiration),

		sequencerIDsCache: make([]common.EnclaveID, 0),

		logger: logger,
	}
}

func (cs *CacheService) Stop() {
	cs.receiptCache.StopJanitor()
	cs.lastBatchesCache.StopJanitor()
	cs.blockCache.StopJanitor()
	cs.batchCacheBySeqNo.StopJanitor()
	cs.seqCacheByHash.StopJanitor()
	cs.seqCacheByHeight.StopJanitor()
	cs.convertedGethHeaderCache.StopJanitor()
	cs.convertedHashCache.StopJanitor()
	cs.eoaCache.StopJanitor()
	cs.contractAddressCache.StopJanitor()
	cs.eventTopicCache.StopJanitor()
	cs.attestedEnclavesCache.StopJanitor()
}

func newFifoCache(logger gethlog.Logger, nrElem int, ttl time.Duration) *gocache.Cache {
	cache := gocache.NewCache().WithMaxSize(nrElem).WithEvictionPolicy(gocache.FirstInFirstOut).WithDefaultTTL(ttl)
	err := cache.StartJanitor()
	if err != nil {
		logger.Crit("Could not initialise fifo cache", log.ErrKey, err)
	}
	return cache
}

func (cs *CacheService) CacheConvertedHash(ctx context.Context, batchHash, convertedHash gethcommon.Hash) {
	cacheValue(ctx, cs.convertedHashCache, cs.logger, batchHash.Bytes(), &convertedHash)
}

func (cs *CacheService) CacheBlock(ctx context.Context, b *types.Header) {
	cacheValue(ctx, cs.blockCache, cs.logger, b.Hash().Bytes(), b)
}

func (cs *CacheService) CacheBatch(ctx context.Context, batch *core.Batch) {
	cacheValue(ctx, cs.batchCacheBySeqNo, cs.logger, batch.SeqNo().Uint64(), batch.Header)
	cacheValue(ctx, cs.seqCacheByHash, cs.logger, batch.Hash().Bytes(), batch.SeqNo())
	// note: the key is (height+1), because for some reason it doesn't like a key of 0
	// should always contain the canonical batch because the cache is overwritten by each new batch after a reorg
	cacheValue(ctx, cs.seqCacheByHeight, cs.logger, batch.NumberU64()+1, batch.SeqNo())

	cacheValue(ctx, cs.lastBatchesCache, cs.logger, batch.SeqNo(), batch)
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
	b, found := cs.lastBatchesCache.Get(fmt.Sprintf("%d", seqNum))
	if !found {
		b1, err := onCacheMiss()
		if err != nil {
			return nil, err
		}
		cs.CacheBatch(ctx, b1)
		b = b1
	}
	cb, ok := b.(*core.Batch)
	if !ok {
		return nil, fmt.Errorf("should not happen. invalid cached batch")
	}
	return cb, nil
}

func (cs *CacheService) ReadEOA(ctx context.Context, addr gethcommon.Address, onCacheMiss func() (*uint64, error)) (*uint64, error) {
	return getCachedValue(ctx, cs.eoaCache, cs.logger, addr.Bytes(), onCacheMiss, true)
}

func (cs *CacheService) ReadContractAddr(ctx context.Context, addr gethcommon.Address, onCacheMiss func() (*enclavedb.Contract, error)) (*enclavedb.Contract, error) {
	return getCachedValue(ctx, cs.contractAddressCache, cs.logger, addr.Bytes(), onCacheMiss, true)
}

func (cs *CacheService) InvalidateContract(addr gethcommon.Address) {
	cs.contractAddressCache.Delete(toString(addr.Bytes()))
}

func (cs *CacheService) ReadEventTopic(ctx context.Context, topic []byte, eventTypeId uint64, onCacheMiss func() (*enclavedb.EventTopic, error)) (*enclavedb.EventTopic, error) {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, eventTypeId)
	key := append(topic, b...)
	return getCachedValue(ctx, cs.eventTopicCache, cs.logger, key, onCacheMiss, true)
}

// CachedReceipt - when all values are nil, it means there is no receipt
type CachedReceipt struct {
	Receipt *types.Receipt
	Tx      *types.Transaction
	From    *gethcommon.Address
	To      *gethcommon.Address
}

var ReceiptDoesNotExist = errors.New("receipt does not exist")

func (cs *CacheService) CacheReceipts(results core.TxExecResults) {
	receipts := make(map[string]any)
	for _, txExecResult := range results {
		receipts[txExecResult.TxWithSender.Tx.Hash().String()] = &CachedReceipt{
			Receipt: txExecResult.Receipt,
			Tx:      txExecResult.TxWithSender.Tx,
			From:    txExecResult.TxWithSender.Sender,
			To:      txExecResult.TxWithSender.Tx.To(),
		}
		cs.logger.Debug("Cache receipt", "tx", txExecResult.TxWithSender.Tx.Hash().String())
	}
	cs.receiptCache.SetAll(receipts)
}

func (cs *CacheService) ReceiptDoesNotExist(txHash gethcommon.Hash) {
	cacheValue(context.Background(), cs.receiptCache, cs.logger, txHash.Bytes(), &CachedReceipt{})
}

func (cs *CacheService) ReadReceipt(_ context.Context, txHash gethcommon.Hash) (*CachedReceipt, error) {
	r, found := cs.receiptCache.Get(txHash.String())
	if !found {
		return nil, errutil.ErrNotFound
	}
	cr, ok := r.(*CachedReceipt)
	if !ok {
		return nil, fmt.Errorf("should not happen. invalid cached receipt")
	}
	if cr.Receipt == nil {
		return nil, ReceiptDoesNotExist
	}
	return cr, nil
}

func (cs *CacheService) DelReceipt(_ context.Context, txHash gethcommon.Hash) error {
	cs.receiptCache.Delete(txHash.String())
	return nil
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
func getCachedValue[V any](ctx context.Context, cache *gocache.Cache, logger gethlog.Logger, key any, onCacheMiss func() (V, error), cacheIfMissing bool) (V, error) {
	var result V

	value, found := cache.Get(toString(key))
	if found {
		// Safely convert the cached value to type V
		if typedValue, ok := value.(V); ok {
			return typedValue, nil
		}
		// Log the type mismatch
		logger.Crit("Cached value type mismatch", "key", key, "expected", fmt.Sprintf("%T", result), "actual", fmt.Sprintf("%T", value))
		// Continue to fetch fresh value
	}

	if onCacheMiss == nil {
		return result, errutil.ErrNotFound
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

func cacheValue(_ context.Context, cache *gocache.Cache, _ gethlog.Logger, key any, v any) {
	cache.Set(toString(key), v)
}

// if anything else is presented, it will use MD5
func toString(key any) string {
	switch k := key.(type) {
	case string:
		return k
	case []byte:
		return base64.StdEncoding.EncodeToString(k)
	case uint64, int64, int, uint:
		return fmt.Sprint(k)
	case *big.Int:
		return fmt.Sprint(k)
	default:
		panic("should not happen. Invalid cache type")
	}
}
