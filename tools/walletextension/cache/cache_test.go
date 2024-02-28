package cache

import (
	"reflect"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/log"

	"github.com/ten-protocol/go-ten/tools/walletextension/common"
)

var tests = map[string]func(t *testing.T){
	"testCacheableMethods":                    testCacheableMethods,
	"testNonCacheableMethods":                 testNonCacheableMethods,
	"testMethodsWithLatestOrPendingParameter": testMethodsWithLatestOrPendingParameter,
}

var cacheTests = map[string]func(cache Cache, t *testing.T){
	"testResultsAreCached":               testResultsAreCached,
	"testCacheTTL":                       testCacheTTL,
	"testCachingAuthenticatedMethods":    testCachingAuthenticatedMethods,
	"testCachingNonAuthenticatedMethods": testCachingNonAuthenticatedMethods,
}

var (
	nonCacheableMethods = []string{"eth_sendrawtransaction", "eth_sendtransaction", "join", "authenticate"}
	encryptionToken     = "test"
	encryptionToken2    = "not-test"
)

func TestGatewayCaching(t *testing.T) {
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			test(t)
		})
	}

	// cache tests
	for name, test := range cacheTests {
		t.Run(name, func(t *testing.T) {
			logger := log.New()
			cache, err := NewCache(logger)
			if err != nil {
				t.Errorf("failed to create cache: %v", err)
			}
			test(cache, t)
		})
	}
}

// testCacheableMethods tests if the cacheable methods are cacheable
func testCacheableMethods(t *testing.T) {
	for method := range cacheableRPCMethods {
		key := &common.RPCRequest{Method: method}
		isCacheable, _, _ := IsCacheable(key, encryptionToken)
		if isCacheable != true {
			t.Errorf("method %s should be cacheable", method)
		}
	}
}

// testNonCacheableMethods tests if the non-cacheable methods are not cacheable
func testNonCacheableMethods(t *testing.T) {
	for _, method := range nonCacheableMethods {
		key := &common.RPCRequest{Method: method}
		isCacheable, _, _ := IsCacheable(key, encryptionToken)
		if isCacheable == true {
			t.Errorf("method %s should not be cacheable", method)
		}
	}
}

// testMethodsWithLatestOrPendingParameter tests if the methods with latest or pending parameter are cacheable
func testMethodsWithLatestOrPendingParameter(t *testing.T) {
	methods := []string{"eth_getCode", "eth_estimateGas", "eth_call"}
	for _, method := range methods {
		key := &common.RPCRequest{Method: method, Params: []interface{}{"0x123", "latest"}}
		_, _, ttl := IsCacheable(key, encryptionToken)
		if ttl != shortCacheTTL {
			t.Errorf("method %s with latest parameter should have TTL of %s, but %s received", method, shortCacheTTL, ttl)
		}

		key = &common.RPCRequest{Method: method, Params: []interface{}{"0x123", "pending"}}
		_, _, ttl = IsCacheable(key, encryptionToken)
		if ttl != shortCacheTTL {
			t.Errorf("method %s with pending parameter should have TTL of %s, but %s received", method, shortCacheTTL, ttl)
		}
	}
}

// testResultsAreCached tests if the results are cached as expected
func testResultsAreCached(cache Cache, t *testing.T) {
	// prepare a cacheable request and imaginary response
	req := &common.RPCRequest{Method: "eth_getBlockByNumber", Params: []interface{}{"0x123"}}
	res := map[string]interface{}{"result": "block"}
	isCacheable, key, ttl := IsCacheable(req, encryptionToken)
	if !isCacheable {
		t.Errorf("method %s should be cacheable", req.Method)
	}
	// set the response in the cache with a TTL
	if !cache.Set(key, res, ttl) {
		t.Errorf("failed to set value in cache for %s", req)
	}

	time.Sleep(50 * time.Millisecond) // wait for the cache to be set
	value, ok := cache.Get(key)
	if !ok {
		t.Errorf("failed to get cached value for %s", req)
	}

	if !reflect.DeepEqual(value, res) {
		t.Errorf("expected %v, got %v", res, value)
	}
}

// testCacheTTL tests if the cache TTL is working as expected
func testCacheTTL(cache Cache, t *testing.T) {
	req := &common.RPCRequest{Method: "eth_blockNumber", Params: []interface{}{"0x123"}}
	res := map[string]interface{}{"result": "100"}
	isCacheable, key, ttl := IsCacheable(req, encryptionToken)

	if !isCacheable {
		t.Errorf("method %s should be cacheable", req.Method)
	}

	if ttl != shortCacheTTL {
		t.Errorf("method %s should have TTL of %s, but %s received", req.Method, shortCacheTTL, ttl)
	}

	// set the response in the cache with a TTL
	if !cache.Set(key, res, ttl) {
		t.Errorf("failed to set value in cache for %s", req)
	}
	time.Sleep(50 * time.Millisecond) // wait for the cache to be set

	// check if the value is in the cache
	value, ok := cache.Get(key)
	if !ok {
		t.Errorf("failed to get cached value for %s", req)
	}

	if !reflect.DeepEqual(value, res) {
		t.Errorf("expected %v, got %v", res, value)
	}

	// sleep for the TTL to expire
	time.Sleep(shortCacheTTL + 100*time.Millisecond)
	_, ok = cache.Get(key)
	if ok {
		t.Errorf("value should not be in the cache after TTL")
	}
}

func testCachingAuthenticatedMethods(cache Cache, t *testing.T) {
	// eth_getTransactionByHash
	authMethods := []string{
		"eth_getTransactionByHash",
		"eth_getCode",
		"eth_getTransactionReceipt",
		"eth_call",
		"eth_estimateGas",
	}
	for _, method := range authMethods {
		req := &common.RPCRequest{Method: method, Params: []interface{}{"0x123"}}
		res := map[string]interface{}{"result": "transaction"}

		// store the response in cache for the first user using encryptionToken
		isCacheable, key, ttl := IsCacheable(req, encryptionToken)

		if !isCacheable {
			t.Errorf("method %s should be cacheable", req.Method)
		}

		// set the response in the cache with a TTL
		if !cache.Set(key, res, ttl) {
			t.Errorf("failed to set value in cache for %s", req)
		}
		time.Sleep(50 * time.Millisecond) // wait for the cache to be set

		// check if the value is in the cache
		value, ok := cache.Get(key)
		if !ok {
			t.Errorf("failed to get cached value for %s", req)
		}

		// for the first error we should have the value in cache
		if !reflect.DeepEqual(value, res) {
			t.Errorf("expected %v, got %v", res, value)
		}

		// now check with the second user asking for the same request, but with a different encryptionToken
		_, key2, _ := IsCacheable(req, encryptionToken2)

		_, okSecondUser := cache.Get(key2)
		if okSecondUser {
			t.Errorf("another user should not see a value the first user cached %s", req)
		}
	}
}

func testCachingNonAuthenticatedMethods(cache Cache, t *testing.T) {
	// eth_getTransactionByHash
	nonAuthMethods := []string{
		"eth_getBlockByNumber",
		"eth_getBlockByHash",
		"eth_chainId",
		"eth_blockNumber",
		"eth_gasPrice",
		"eth_feeHistory",
	}

	for _, method := range nonAuthMethods {
		req := &common.RPCRequest{Method: method, Params: []interface{}{"0x123"}}
		res := map[string]interface{}{"result": "transaction"}

		// store the response in cache for the first user using encryptionToken
		isCacheable, key, ttl := IsCacheable(req, encryptionToken)

		if !isCacheable {
			t.Errorf("method %s should be cacheable", req.Method)
		}

		// set the response in the cache with a TTL
		if !cache.Set(key, res, ttl) {
			t.Errorf("failed to set value in cache for %s", req)
		}
		time.Sleep(50 * time.Millisecond) // wait for the cache to be set

		// check if the value is in the cache
		value, ok := cache.Get(key)
		if !ok {
			t.Errorf("failed to get cached value for %s", req)
		}

		// for the first error we should have the value in cache
		if !reflect.DeepEqual(value, res) {
			t.Errorf("expected %v, got %v", res, value)
		}

		// now check with the second user asking for the same request, but with a different encryptionToken
		_, key2, _ := IsCacheable(req, encryptionToken2)

		_, okSecondUser := cache.Get(key2)
		if !okSecondUser {
			t.Errorf("another user should see a value the first user cached %s", req)
		}
	}
}
