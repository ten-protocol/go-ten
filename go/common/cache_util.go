package common

import (
	"context"

	"github.com/eko/gocache/lib/v4/cache"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common/log"
)

// GetCachedValue - returns the cached value for the provided key. If the key is not found, then invoke the 'onFailed' function
// which returns the value, and cache it
func GetCachedValue[V any](cache *cache.Cache[*V], logger gethlog.Logger, key any, onCacheMiss func(any) (*V, error)) (*V, error) {
	value, err := cache.Get(context.Background(), key)
	if err != nil || value == nil {
		// todo metrics for cache misses
		v, err := onCacheMiss(key)
		if err != nil {
			return v, err
		}
		if v == nil {
			logger.Crit("Returned a nil value from the onCacheMiss function. Should not happen.")
		}
		CacheValue(cache, logger, key, v)
		return v, nil
	}

	return value, err
}

func CacheValue[V any](cache *cache.Cache[*V], logger gethlog.Logger, key any, v *V) {
	if v == nil {
		return
	}
	err := cache.Set(context.Background(), key, v)
	if err != nil {
		logger.Error("Could not store value in cache", log.ErrKey, err)
	}
}
