package cache

import (
	"time"

	"github.com/ethereum/go-ethereum/log"

	"github.com/dgraph-io/ristretto"
)

const (
	numCounters = 1e7       // number of keys to track frequency of (10M).
	maxCost     = 1_000_000 // 1 million entries
	bufferItems = 64        // number of keys per Get buffer.
	defaultCost = 1         // default cost of cache.
)

type ristrettoCache struct {
	cache *ristretto.Cache
	quit  chan struct{}
}

// NewRistrettoCache returns a new ristrettoCache.
func NewRistrettoCache(logger log.Logger) (Cache, error) {
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: numCounters,
		MaxCost:     maxCost,
		BufferItems: bufferItems,
		Metrics:     true,
	})
	if err != nil {
		return nil, err
	}

	c := &ristrettoCache{
		cache: cache,
		quit:  make(chan struct{}),
	}

	// Start the metrics logging
	go c.startMetricsLogging(logger)

	return c, nil
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
			logger.Info("Cache metrics: Hits: %d, Misses: %d, Cost Added: %d\n",
				metrics.Hits(), metrics.Misses(), metrics.CostAdded())
		case <-c.quit:
			ticker.Stop()
			return
		}
	}
}
