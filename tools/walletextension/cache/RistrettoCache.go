package cache

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/log"

	"github.com/dgraph-io/ristretto/v2"
)

const (
	bufferItems = 64 // number of keys per Get buffer.
	defaultCost = 1  // default cost of cache.
)

type ristrettoCache struct {
	cache              *ristretto.Cache[[]byte, any]
	quit               chan struct{}
	lastEviction       time.Time
	shortLivingEnabled *atomic.Bool
}

// NewRistrettoCacheWithEviction returns a new ristrettoCache.
func NewRistrettoCacheWithEviction(nrElems int, logger log.Logger) (Cache, error) {
	cache, err := ristretto.NewCache[[]byte, any](&ristretto.Config[[]byte, any]{
		NumCounters: int64(nrElems * 10),
		MaxCost:     int64(nrElems),
		BufferItems: bufferItems,
		Metrics:     true,
	})
	if err != nil {
		return nil, err
	}

	c := &ristrettoCache{
		cache:              cache,
		quit:               make(chan struct{}),
		lastEviction:       time.Now(),
		shortLivingEnabled: &atomic.Bool{},
	}
	c.shortLivingEnabled.Store(true)

	// Start the metrics logging
	go c.startMetricsLogging(logger)

	return c, nil
}

func (c *ristrettoCache) EvictShortLiving() {
	// this event happens when a new batch is received, so the cache can be enabled
	c.shortLivingEnabled.Store(true)
	c.lastEviction = time.Now()
}

func (c *ristrettoCache) DisableShortLiving() {
	c.shortLivingEnabled.Store(false)
}

func (c *ristrettoCache) IsEvicted(key []byte, originalTTL time.Duration) bool {
	if !c.shortLivingEnabled.Load() {
		return true
	}
	remainingTTL, notExpired := c.cache.GetTTL(key)
	if !notExpired {
		return true
	}
	cachedTime := time.Now().Add(remainingTTL).Add(-originalTTL)
	// ... LE...Cached...Now - Valid
	// ... Cached...LE...Now - Evicted
	return c.lastEviction.After(cachedTime)
}

// Set adds the key and value to the cache.
func (c *ristrettoCache) Set(key []byte, value any, ttl time.Duration) bool {
	return c.cache.SetWithTTL(key, value, defaultCost, ttl)
}

// Get returns the value for the given key if it exists.
func (c *ristrettoCache) Get(key []byte) (value any, ok bool) {
	return c.cache.Get(key)
}

func (c *ristrettoCache) Remove(key []byte) {
	c.cache.Del(key)
}

// startMetricsLogging starts logging cache metrics every hour.
func (c *ristrettoCache) startMetricsLogging(logger log.Logger) {
	ticker := time.NewTicker(1 * time.Hour)
	for {
		select {
		case <-ticker.C:
			metrics := c.cache.Metrics
			logger.Info(fmt.Sprintf("Cache metrics: Hits: %d, Misses: %d, Cost Added: %d",
				metrics.Hits(), metrics.Misses(), metrics.CostAdded()))
		case <-c.quit:
			ticker.Stop()
			return
		}
	}
}
