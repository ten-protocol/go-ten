package cache

import (
	"time"

	"github.com/ethereum/go-ethereum/log"
)

type Cache interface {
	// EvictShortLiving - notify the cache that all short living elements cached before the events should be considered as evicted.
	EvictShortLiving()

	// IsEvicted - based on the eviction event and the time of caching, calculates whether the key was evicted
	IsEvicted(key any, originalTTL time.Duration) bool

	Set(key []byte, value any, ttl time.Duration) bool
	Get(key []byte) (value any, ok bool)
	Remove(key []byte)
}

func NewCache(logger log.Logger) (Cache, error) {
	return NewRistrettoCacheWithEviction(logger)
}
