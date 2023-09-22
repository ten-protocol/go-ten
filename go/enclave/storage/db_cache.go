package storage

import (
	"context"

	"github.com/eko/gocache/lib/v4/cache"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/go-obscuro/go/common/log"
)

func getCachedValue[V any](cache *cache.Cache[[]byte], logger gethlog.Logger, key any, onFailed func(any) (V, error)) (V, error) {
	value, err := cache.Get(context.Background(), key)
	if err != nil {
		// todo metrics for cache misses
		b, err := onFailed(key)
		if err != nil {
			return b, err
		}
		cacheValue(cache, logger, key, b)
		return b, err
	}

	v := new(V)
	err = rlp.DecodeBytes(value, v)
	return *v, err
}

func cacheValue(cache *cache.Cache[[]byte], logger gethlog.Logger, key any, v any) {
	encoded, err := rlp.EncodeToBytes(v)
	if err != nil {
		logger.Error("Could not encode value to store in cache", log.ErrKey, err)
		return
	}
	err = cache.Set(context.Background(), key, encoded)
	if err != nil {
		logger.Error("Could not store value in cache", log.ErrKey, err)
	}
}
