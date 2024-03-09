package rpcapi

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/ten-protocol/go-ten/tools/walletextension/cache"

	"github.com/ethereum/go-ethereum/common"
	common2 "github.com/ten-protocol/go-ten/tools/walletextension/common"
)

const (
	exposedParams       = "exposedParams"
	ethCallPaddedArgLen = 64
	ethCallAddrPadding  = "000000000000000000000000"

	notAuthorised = "not authorised"

	longCacheTTL  = 5 * time.Hour
	shortCacheTTL = 1 * time.Second
)

type ExecCfg struct {
	account             *common.Address
	computeFromCallback func(user *GWUser) *common.Address
	tryAll              bool
	tryUntilAuthorised  bool
	adjustArgs          func(acct *GWAccount) []any
	cacheCfg            *CacheCfg
}

type CacheCfg struct {
	TTL time.Duration
	// logic based on block
	TTLCallback func() time.Duration
}

func ExecAuthRPC[R any](ctx context.Context, w *Services, cfg *ExecCfg, method string, args ...any) (*R, error) {
	// requestStartTime := time.Now()
	userId, err := extractUserId(ctx)
	if err != nil {
		return nil, err
	}

	user, err := getUser(userId, w.Storage)
	if err != nil {
		return nil, err
	}

	// todo - add method to args
	return withCache(w.Cache, cfg.cacheCfg, args, func() (*R, error) {
		// determine candidate "from"
		candidateAccts := make([]*GWAccount, 0)

		// for users with multiple accounts determine a candidate account
		switch {
		case cfg.account != nil:
			acc := user.accounts[*cfg.account]
			if acc == nil {
				return nil, fmt.Errorf("invalid account %s", cfg.account)
			}
			candidateAccts = append(candidateAccts, acc)

		case cfg.computeFromCallback != nil:
			addr := cfg.computeFromCallback(user)
			if addr == nil {
				return nil, fmt.Errorf("invalid request")
			}
			acc := user.accounts[*addr]
			if acc == nil {
				// todo - return an account
				a := reflect.ValueOf(user.accounts).MapKeys()[0].Interface().(common.Address)
				candidateAccts = append(candidateAccts, user.accounts[a])
				break
			}
			candidateAccts = append(candidateAccts, acc)

		case cfg.tryAll, cfg.tryUntilAuthorised:
			for _, acc := range user.accounts {
				candidateAccts = append(candidateAccts, acc)
			}

		default:
			return nil, fmt.Errorf("programming error. invalid owner detection strategy")
		}

		var rpcErr error
		for _, acct := range candidateAccts {
			result := new(R)
			rpcClient, err := acct.connect(w.HostAddrHTTP, w.Logger())
			if err != nil {
				rpcErr = err
				continue
			}
			adjustedArgs := args
			if cfg.adjustArgs != nil {
				adjustedArgs = cfg.adjustArgs(acct)
			}
			err = rpcClient.CallContext(ctx, result, method, adjustedArgs...)
			if err != nil {
				// todo - is this correct?
				if cfg.tryUntilAuthorised && err.Error() != notAuthorised {
					return nil, err
				}
				rpcErr = err
				continue
			}
			return result, nil
		}
		return nil, rpcErr
	})
}

func UnauthenticatedTenRPCCall[R any](ctx context.Context, w *Services, cfg *CacheCfg, method string, args ...any) (*R, error) {
	return withCache(w.Cache, cfg, args, func() (*R, error) {
		resp := new(R)
		unauthedRPC, err := w.UnauthenticatedClient()
		if err != nil {
			return nil, err
		}
		err = unauthedRPC.CallContext(ctx, resp, method, args...)
		return resp, err
	})
}

func extractUserId(ctx context.Context) ([]byte, error) {
	params := ctx.Value(exposedParams).(map[string]string)
	// todo handle errors
	userId, ok := params[common2.EncryptedTokenQueryParameter]
	if !ok || len(userId) < 3 {
		// todo - remove
		return []byte(common2.DefaultUser), nil
		// return nil, fmt.Errorf("invalid encryption token %s", userId)
	}
	return common2.GetUserIDbyte(userId)
}

// generateCacheKey generates a cache key for the given method, encryptionToken and parameters
// encryptionToken is used to generate a unique cache key for each user and empty string should be used for public data
func generateCacheKey(params []any) []byte {
	// Serialize parameters
	rawKey, err := json.Marshal(params)
	if err != nil {
		return nil
	}

	// Optional: Apply hashing
	hasher := sha256.New()
	hasher.Write(rawKey)

	return hasher.Sum(nil)
}

func withCache[R any](cache cache.Cache, cfg *CacheCfg, args []any, onCacheMiss func() (*R, error)) (*R, error) {
	if cfg == nil {
		return onCacheMiss()
	}

	cacheTTL := cfg.TTL
	if cfg.TTLCallback != nil {
		cacheTTL = cfg.TTLCallback()
	}
	isCacheable := cfg.TTL > 0

	var cacheKey []byte
	if isCacheable {
		cacheKey = generateCacheKey(args)
		if cachedValue, ok := cache.Get(cacheKey); ok {
			// cloning?
			returnValue := cachedValue.(*R)
			return returnValue, nil
		}
	}

	result, err := onCacheMiss()

	// cache only non-nil values
	if isCacheable && err == nil && result != nil {
		cache.Set(cacheKey, result, cacheTTL)
	}

	return result, err
}
