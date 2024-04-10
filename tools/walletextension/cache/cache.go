package cache

import (
	"time"

	"github.com/ethereum/go-ethereum/log"
)

type Cache interface {
	EvictShortLiving()
	IsEvicted(key any, originalTTL time.Duration) bool
	Set(key []byte, value any, ttl time.Duration) bool
	Get(key []byte) (value any, ok bool)
	Remove(key []byte)
}

func NewCache(logger log.Logger) (Cache, error) {
	return NewRistrettoCacheWithEviction(logger)
}
