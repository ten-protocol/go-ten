package cache

import (
	"time"

	"github.com/ethereum/go-ethereum/log"

	"github.com/dgraph-io/ristretto"
)

const (
	numCounters  = 1e7     // number of keys to track frequency of (10M).
	maxCost      = 1 << 30 // maximum cost of cache (1GB).
	bufferItems  = 64      // number of keys per Get buffer.
	defaultConst = 1       // default cost of cache.
)

type RistrettoCache struct {
	cache *ristretto.Cache
	quit  chan struct{}
}

// NewRistrettoCache returns a new RistrettoCache.
func NewRistrettoCache(logger log.Logger) (*RistrettoCache, error) {
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: numCounters,
		MaxCost:     maxCost,
		BufferItems: bufferItems,
		Metrics:     true,
	})
	if err != nil {
		return nil, err
	}

	c := &RistrettoCache{
		cache: cache,
		quit:  make(chan struct{}),
	}

	// Start the metrics logging
	go c.startMetricsLogging(logger)

	return c, nil
}

// Set adds the key and value to the cache.
func (c *RistrettoCache) Set(key string, value map[string]interface{}, ttl time.Duration) bool {
	return c.cache.SetWithTTL(key, value, defaultConst, ttl)
}

// Get returns the value for the given key if it exists.
func (c *RistrettoCache) Get(key string) (value map[string]interface{}, ok bool) {
	item, found := c.cache.Get(key)
	if !found {
		return nil, false
	}

	// Assuming the item is stored as a map[string]interface{}, otherwise you need to type assert to the correct type.
	value, ok = item.(map[string]interface{})
	if !ok {
		// The item isn't of type map[string]interface{}
		return nil, false
	}

	return value, true
}

// startMetricsLogging starts logging cache metrics every hour.
func (c *RistrettoCache) startMetricsLogging(logger log.Logger) {
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
