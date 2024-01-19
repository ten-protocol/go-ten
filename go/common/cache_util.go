package common

import (
	"context"

	"github.com/eko/gocache/lib/v4/cache"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common/log"
)

func GetCachedValue[V any](cache *cache.Cache[V], logger gethlog.Logger, key any, onFailed func(any) (V, error)) (V, error) {
	value, err := cache.Get(context.Background(), key)
	if err != nil {
		// todo metrics for cache misses
		b, err := onFailed(key)
		if err != nil {
			return b, err
		}
		CacheValue(cache, logger, key, b)
		return b, err
	}

	return value, err
}

func CacheValue[V any](cache *cache.Cache[V], logger gethlog.Logger, key any, v V) {
	err := cache.Set(context.Background(), key, v)
	if err != nil {
		logger.Error("Could not store value in cache", log.ErrKey, err)
	}
}
