package cache

import (
	"fmt"
	"time"

	"golang.org/x/sync/singleflight"

	"github.com/ethereum/go-ethereum/log"
)

type Cache interface {
	// EvictShortLiving - notify the cache that all short living elements cached before the events should be considered as evicted.
	EvictShortLiving()

	// DisableShortLiving disables the caching of short-living elements.
	DisableShortLiving()

	// IsShortLivingEnabled returns whether short-living elements are currently being cached.
	IsShortLivingEnabled() bool

	// IsEvicted - based on the eviction event and the time of caching, calculates whether the key was evicted
	IsEvicted(key []byte, originalTTL time.Duration) bool

	Set(key []byte, value any, ttl time.Duration) bool
	Get(key []byte) (value any, ok bool)
	Remove(key []byte)
}

func NewCache(nrElems int, logger log.Logger) (Cache, error) {
	return NewRistrettoCacheWithEviction(nrElems, logger)
}

type Strategy uint8

const (
	NoCache     Strategy = iota
	LatestBatch Strategy = iota
	LongLiving  Strategy = iota

	longCacheTTL  = 5 * time.Hour
	shortCacheTTL = 1 * time.Second
)

type Cfg struct {
	Type        Strategy
	DynamicType func() Strategy
}

// Global singleflight group
var sfGroup singleflight.Group

func WithCache[R any](cache Cache, cfg *Cfg, cacheKey []byte, onCacheMiss func() (*R, error)) (*R, error) {
	if cfg == nil {
		return onCacheMiss()
	}

	sfKey := string(cacheKey)

	// serialises and optimizes access to the cache for the same key
	res, err, _ := sfGroup.Do(sfKey, func() (interface{}, error) {
		cacheType := cfg.Type
		if cfg.DynamicType != nil {
			cacheType = cfg.DynamicType()
		}

		if cacheType == NoCache {
			return onCacheMiss()
		}

		// we implement a custom cache eviction logic for the cache strategy of type LatestBatch.
		// when a new batch is created, all entries with "LatestBatch" are considered evicted.
		// elements not cached for a specific batch are not evicted
		isEvicted := false
		ttl := longCacheTTL
		if cacheType == LatestBatch {
			// If short living is disabled, bypass cache entirely to prevent serving stale data
			if !cache.IsShortLivingEnabled() {
				return onCacheMiss()
			}
			ttl = shortCacheTTL
			isEvicted = cache.IsEvicted(cacheKey, ttl)
		}

		if !isEvicted {
			cachedValue, foundInCache := cache.Get(cacheKey)
			if foundInCache {
				returnValue, ok := cachedValue.(*R)
				if !ok {
					return nil, fmt.Errorf("unexpected error. Invalid format cached. %v", cachedValue)
				}
				return returnValue, nil
			}
		}

		result, err := onCacheMiss()

		// cache only non-nil values
		if err == nil && result != nil {
			cache.Set(cacheKey, result, ttl)
		}

		return result, err
	})

	if err != nil {
		return nil, err
	}

	// Convert back to the correct type
	result, ok := res.(*R)
	if !ok {
		return nil, fmt.Errorf("singleflight returned unexpected type: %T", res)
	}

	return result, nil
}

type noOpCache struct{}

func NewNoOpCache() Cache {
	return &noOpCache{}
}

func (c *noOpCache) EvictShortLiving() {}

func (c *noOpCache) DisableShortLiving() {}

func (c *noOpCache) IsShortLivingEnabled() bool {
	return false
}

func (c *noOpCache) IsEvicted(key []byte, originalTTL time.Duration) bool {
	return false
}

func (c *noOpCache) Set(key []byte, value any, ttl time.Duration) bool {
	return false
}

func (c *noOpCache) Get(key []byte) (value any, ok bool) {
	return nil, false
}

func (c *noOpCache) Remove(key []byte) {}
