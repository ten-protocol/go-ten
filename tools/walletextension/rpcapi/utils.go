package rpcapi

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ten-protocol/go-ten/go/common/viewingkey"

	"github.com/ten-protocol/go-ten/lib/gethfork/rpc"

	"github.com/status-im/keycard-go/hexutils"

	"github.com/ten-protocol/go-ten/tools/walletextension/cache"

	"github.com/ethereum/go-ethereum/common"
)

const (
	ethCallPaddedArgLen = 64
	ethCallAddrPadding  = "000000000000000000000000"

	notAuthorised = "not authorised"

	longCacheTTL  = 5 * time.Hour
	shortCacheTTL = 1 * time.Second
)

var rpcNotImplemented = fmt.Errorf("rpc endpoint not implemented")

type ExecCfg struct {
	account             *common.Address
	computeFromCallback func(user *GWUser) *common.Address
	tryAll              bool
	tryUntilAuthorised  bool
	adjustArgs          func(acct *GWAccount) []any
	cacheCfg            *CacheCfg
}

type CacheCfg struct {
	// ResetWhenNewBlock bool todo
	TTL time.Duration
	// logic based on block
	// todo - handle block in the future
	TTLCallback func() time.Duration
}

func UnauthenticatedTenRPCCall[R any](ctx context.Context, w *Services, cfg *CacheCfg, method string, args ...any) (*R, error) {
	audit(w, "RPC start method=%s args=%v", method, args)
	requestStartTime := time.Now()
	cacheArgs := []any{method}
	cacheArgs = append(cacheArgs, args...)

	res, err := withCache(w.Cache, cfg, generateCacheKey(cacheArgs), func() (*R, error) {
		var resp *R
		unauthedRPC, err := w.UnauthenticatedClient()
		if err != nil {
			return nil, err
		}
		if ctx == nil {
			err = unauthedRPC.Call(&resp, method, args...)
		} else {
			err = unauthedRPC.CallContext(ctx, &resp, method, args...)
		}
		return resp, err
	})
	audit(w, "RPC call. method=%s args=%v result=%s error=%s time=%d", method, args, res, err, time.Since(requestStartTime).Milliseconds())
	return res, err
}

func ExecAuthRPC[R any](ctx context.Context, w *Services, cfg *ExecCfg, method string, args ...any) (*R, error) {
	audit(w, "RPC start method=%s args=%v", method, args)
	requestStartTime := time.Now()
	userID, err := extractUserID(ctx, w)
	if err != nil {
		return nil, err
	}

	user, err := getUser(userID, w)
	if err != nil {
		return nil, err
	}

	cacheArgs := []any{userID, method}
	cacheArgs = append(cacheArgs, args...)

	res, err := withCache(w.Cache, cfg.cacheCfg, generateCacheKey(cacheArgs), func() (*R, error) {
		// determine candidate "from"
		candidateAccts, err := getCandidateAccounts(user, w, cfg)
		if err != nil {
			return nil, err
		}
		if len(candidateAccts) == 0 {
			return nil, fmt.Errorf("illegal access")
		}

		var rpcErr error
		for _, acct := range candidateAccts {
			var result *R
			rpcClient, err := acct.connect(w.HostAddrHTTP, w.Logger())
			if err != nil {
				rpcErr = err
				continue
			}
			adjustedArgs := args
			if cfg.adjustArgs != nil {
				adjustedArgs = cfg.adjustArgs(acct)
			}
			err = rpcClient.CallContext(ctx, &result, method, adjustedArgs...)
			if err != nil {
				// for calls where we know the expected error we can return early
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
	audit(w, "RPC call. uid=%s, method=%s args=%v result=%s error=%s time=%d", hexutils.BytesToHex(userID), method, args, res, err, time.Since(requestStartTime).Milliseconds())
	return res, err
}

func getCandidateAccounts(user *GWUser, _ *Services, cfg *ExecCfg) ([]*GWAccount, error) {
	candidateAccts := make([]*GWAccount, 0)
	// for users with multiple accounts try to determine a candidate account based on the available information
	switch {
	case cfg.account != nil:
		acc := user.accounts[*cfg.account]
		if acc != nil {
			candidateAccts = append(candidateAccts, acc)
			return candidateAccts, nil
		}

	case cfg.computeFromCallback != nil:
		suggestedAddress := cfg.computeFromCallback(user)
		if suggestedAddress != nil {
			acc := user.accounts[*suggestedAddress]
			if acc != nil {
				candidateAccts = append(candidateAccts, acc)
				return candidateAccts, nil
			}
		}
	}

	if cfg.tryAll || cfg.tryUntilAuthorised {
		for _, acc := range user.accounts {
			candidateAccts = append(candidateAccts, acc)
		}
	}

	return candidateAccts, nil
}

func extractUserID(ctx context.Context, _ *Services) ([]byte, error) {
	token, ok := ctx.Value(rpc.GWTokenKey{}).(string)
	if !ok {
		return nil, fmt.Errorf("invalid userid: %s", ctx.Value(rpc.GWTokenKey{}))
	}
	userID := common.FromHex(token)
	if len(userID) != viewingkey.UserIDLength {
		return nil, fmt.Errorf("invalid userid: %s", token)
	}
	return userID, nil
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

func withCache[R any](cache cache.Cache, cfg *CacheCfg, cacheKey []byte, onCacheMiss func() (*R, error)) (*R, error) {
	if cfg == nil {
		return onCacheMiss()
	}

	cacheTTL := cfg.TTL
	if cfg.TTLCallback != nil {
		cacheTTL = cfg.TTLCallback()
	}
	isCacheable := cacheTTL > 0

	if isCacheable {
		if cachedValue, ok := cache.Get(cacheKey); ok {
			// cloning?
			returnValue, ok := cachedValue.(*R)
			if !ok {
				return nil, fmt.Errorf("unexpected error. Invalid format cached. %v", cachedValue)
			}
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

func audit(services *Services, msg string, params ...any) {
	println(fmt.Sprintf(msg, params...))
	if services.Config.VerboseFlag {
		services.FileLogger.Info(fmt.Sprintf(msg, params...))
	}
}

func cacheTTLBlockNumberOrHash(blockNrOrHash rpc.BlockNumberOrHash) time.Duration {
	if blockNrOrHash.BlockNumber != nil && blockNrOrHash.BlockNumber.Int64() <= 0 {
		return shortCacheTTL
	}
	return longCacheTTL
}

func cacheTTLBlockNumber(lastBlock rpc.BlockNumber) time.Duration {
	if lastBlock > 0 {
		return longCacheTTL
	}
	return shortCacheTTL
}
