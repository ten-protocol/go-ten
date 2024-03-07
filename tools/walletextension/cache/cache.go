package cache

import (
	"time"

	"github.com/ethereum/go-ethereum/log"
)

type Cache interface {
	Set(key []byte, value any, ttl time.Duration) bool
	Get(key []byte) (value any, ok bool)
}

func NewCache(logger log.Logger) (Cache, error) {
	return NewRistrettoCache(logger)
}
