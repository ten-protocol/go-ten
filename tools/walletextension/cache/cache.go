package cache

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/log"

	"github.com/ten-protocol/go-ten/tools/walletextension/common"
)

const (
	longCacheTTL  = 5 * time.Hour
	shortCacheTTL = 1 * time.Second
)

// CacheableRPCMethods is a map of Ethereum JSON-RPC methods that can be cached and their TTL
var cacheableRPCMethods = map[string]time.Duration{
	// Ethereum JSON-RPC methods that can be cached long time
	"eth_getBlockByNumber":     longCacheTTL,
	"eth_getBlockByHash":       longCacheTTL,
	"eth_getTransactionByHash": longCacheTTL,
	"eth_chainId":              longCacheTTL,

	// Ethereum JSON-RPC methods that can be cached short time
	"eth_blockNumber":           shortCacheTTL,
	"eth_getCode":               shortCacheTTL,
	"eth_getBalance":            shortCacheTTL,
	"eth_getTransactionReceipt": shortCacheTTL,
	"eth_call":                  shortCacheTTL,
	"eth_gasPrice":              shortCacheTTL,
	"eth_getTransactionCount":   shortCacheTTL,
	"eth_estimateGas":           shortCacheTTL,
	"eth_feeHistory":            shortCacheTTL,
}

type Cache interface {
	Set(key string, value map[string]interface{}, ttl time.Duration) bool
	Get(key string) (value map[string]interface{}, ok bool)
}

func NewCache(logger log.Logger) (Cache, error) {
	return NewRistrettoCache(logger)
}

// IsCacheable checks if the given RPC request is cacheable and returns the cache key and TTL
func IsCacheable(key *common.RPCRequest) (bool, string, time.Duration) {
	if key == nil || key.Method == "" {
		return false, "", 0
	}

	// Check if the method is cacheable
	ttl, isCacheable := cacheableRPCMethods[key.Method]

	if isCacheable {
		// method is cacheable - select cache key
		switch key.Method {
		case "eth_getCode", "eth_getBalance", "eth_getTransactionCount", "eth_estimateGas", "eth_call":
			if len(key.Params) == 1 || len(key.Params) == 2 && (key.Params[1] == "latest" || key.Params[1] == "pending") {
				return true, GenerateCacheKey(key.Method, key.Params...), ttl
			}
			// in this case, we have a fixed block number, and we can cache the result for a long time
			return true, GenerateCacheKey(key.Method, key.Params...), longCacheTTL
		case "eth_feeHistory":
			if len(key.Params) == 2 || len(key.Params) == 3 && (key.Params[2] == "latest" || key.Params[2] == "pending") {
				return true, GenerateCacheKey(key.Method, key.Params...), ttl
			}
			// in this case, we have a fixed block number, and we can cache the result for a long time
			return true, GenerateCacheKey(key.Method, key.Params...), longCacheTTL
		default:
			return true, GenerateCacheKey(key.Method, key.Params...), ttl
		}
	}

	// method is not cacheable
	return false, "", 0
}

// GenerateCacheKey generates a cache key for the given method and parameters
func GenerateCacheKey(method string, params ...interface{}) string {
	// Serialize parameters
	paramBytes, err := json.Marshal(params)
	if err != nil {
		return ""
	}

	// Concatenate method name and parameters
	rawKey := method + string(paramBytes)

	// Optional: Apply hashing
	hasher := sha256.New()
	hasher.Write([]byte(rawKey))
	hashedKey := fmt.Sprintf("%x", hasher.Sum(nil))

	return hashedKey
}
