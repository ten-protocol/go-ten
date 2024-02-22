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

// Define a struct to hold the cache TTL and auth requirement
type RPCMethodCacheConfig struct {
	CacheTTL     time.Duration
	RequiresAuth bool
}

// CacheableRPCMethods is a map of Ethereum JSON-RPC methods that can be cached and their TTL
var cacheableRPCMethods = map[string]RPCMethodCacheConfig{
	// Ethereum JSON-RPC methods that can be cached long time
	"eth_getBlockByNumber":     {longCacheTTL, false},
	"eth_getBlockByHash":       {longCacheTTL, false},
	"eth_getTransactionByHash": {longCacheTTL, true},
	"eth_chainId":              {longCacheTTL, false},

	// Ethereum JSON-RPC methods that can be cached short time
	"eth_blockNumber": {shortCacheTTL, false},
	"eth_getCode":     {shortCacheTTL, true},
	// "eth_getBalance":            {longCacheTTL, true},// excluded for test: gen_cor_059
	"eth_getTransactionReceipt": {shortCacheTTL, true},
	"eth_call":                  {shortCacheTTL, true},
	"eth_gasPrice":              {shortCacheTTL, false},
	// "eth_getTransactionCount": {longCacheTTL, true}, // excluded for test: gen_cor_009
	"eth_estimateGas": {shortCacheTTL, true},
	"eth_feeHistory":  {shortCacheTTL, false},
}

type Cache interface {
	Set(key string, value map[string]interface{}, ttl time.Duration) bool
	Get(key string) (value map[string]interface{}, ok bool)
}

func NewCache(logger log.Logger) (Cache, error) {
	return NewRistrettoCache(logger)
}

// IsCacheable checks if the given RPC request is cacheable and returns the cache key and TTL
func IsCacheable(key *common.RPCRequest, encryptionToken string) (bool, string, time.Duration) {
	if key == nil || key.Method == "" {
		return false, "", 0
	}

	// Check if the method is cacheable
	methodCacheConfig, isCacheable := cacheableRPCMethods[key.Method]

	// If method does not need to be authenticated, we can don't need to cache it per user
	if !methodCacheConfig.RequiresAuth {
		encryptionToken = ""
	}

	if isCacheable {
		// method is cacheable - select cache key and ttl
		switch key.Method {
		case "eth_getCode", "eth_getBalance", "eth_estimateGas", "eth_call":
			if len(key.Params) == 1 || len(key.Params) == 2 && (key.Params[1] == "latest" || key.Params[1] == "pending") {
				return true, GenerateCacheKey(key.Method, encryptionToken, key.Params...), methodCacheConfig.CacheTTL
			}
			// in this case, we have a fixed block number, and we can cache the result for a long time
			return true, GenerateCacheKey(key.Method, encryptionToken, key.Params...), longCacheTTL
		case "eth_feeHistory":
			if len(key.Params) == 2 || len(key.Params) == 3 && (key.Params[2] == "latest" || key.Params[2] == "pending") {
				return true, GenerateCacheKey(key.Method, encryptionToken, key.Params...), methodCacheConfig.CacheTTL
			}
			// in this case, we have a fixed block number, and we can cache the result for a long time
			return true, GenerateCacheKey(key.Method, encryptionToken, key.Params...), longCacheTTL

		default:
			return true, GenerateCacheKey(key.Method, encryptionToken, key.Params...), methodCacheConfig.CacheTTL
		}
	}

	// method is not cacheable
	return false, "", 0
}

// GenerateCacheKey generates a cache key for the given method, encryptionToken and parameters
// encryptionToken is used to generate a unique cache key for each user and empty string should be used for public data
func GenerateCacheKey(method string, encryptionToken string, params ...interface{}) string {
	// Serialize parameters
	paramBytes, err := json.Marshal(params)
	if err != nil {
		return ""
	}

	// Concatenate method name and parameters
	rawKey := method + encryptionToken + string(paramBytes)

	// Optional: Apply hashing
	hasher := sha256.New()
	hasher.Write([]byte(rawKey))
	hashedKey := fmt.Sprintf("%x", hasher.Sum(nil))

	return hashedKey
}
