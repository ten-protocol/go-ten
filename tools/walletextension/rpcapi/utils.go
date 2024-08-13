package rpcapi

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/ten-protocol/go-ten/go/common/measure"
	"github.com/ten-protocol/go-ten/go/enclave/core"

	"github.com/ten-protocol/go-ten/go/common/log"

	gethlog "github.com/ethereum/go-ethereum/log"
	pool "github.com/jolestar/go-commons-pool/v2"
	tenrpc "github.com/ten-protocol/go-ten/go/rpc"
	wecommon "github.com/ten-protocol/go-ten/tools/walletextension/common"

	"github.com/ten-protocol/go-ten/go/common/viewingkey"

	"github.com/ten-protocol/go-ten/lib/gethfork/rpc"

	"github.com/status-im/keycard-go/hexutils"

	"github.com/ten-protocol/go-ten/tools/walletextension/cache"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

const (
	ethCallPaddedArgLen = 64
	ethCallAddrPadding  = "000000000000000000000000"

	notAuthorised = "not authorised"
	serverBusy    = "server busy. please retry later"

	longCacheTTL  = 5 * time.Hour
	shortCacheTTL = 1 * time.Minute

	// hardcoding the maximum time for an RPC request
	// this value will be propagated to the node and enclave and all the operations
	maximumRPCCallDuration  = 5 * time.Second
	sendTransactionDuration = 20 * time.Second
)

var rpcNotImplemented = fmt.Errorf("rpc endpoint not implemented")

type ExecCfg struct {
	account             *gethcommon.Address
	computeFromCallback func(user *GWUser) *gethcommon.Address
	tryAll              bool
	tryUntilAuthorised  bool
	adjustArgs          func(acct *GWAccount) []any
	cacheCfg            *CacheCfg
	timeout             time.Duration
}

type CacheStrategy uint8

const (
	NoCache     CacheStrategy = iota
	LatestBatch CacheStrategy = iota
	LongLiving  CacheStrategy = iota
)

type CacheCfg struct {
	CacheType        CacheStrategy
	CacheTypeDynamic func() CacheStrategy
}

func UnauthenticatedTenRPCCall[R any](ctx context.Context, w *Services, cfg *CacheCfg, method string, args ...any) (*R, error) {
	if ctx == nil {
		return nil, errors.New("invalid call. nil Context")
	}
	audit(w, "RPC start method=%s args=%v", method, args)
	requestStartTime := time.Now()
	cacheArgs := []any{method}
	cacheArgs = append(cacheArgs, args...)

	res, err := withCache(w.Cache, cfg, generateCacheKey(cacheArgs), func() (*R, error) {
		return withPlainRPCConnection(ctx, w, func(client *rpc.Client) (*R, error) {
			var resp *R
			var err error

			// wrap the context with a timeout to prevent long executions
			timeoutContext, cancelCtx := context.WithTimeout(ctx, maximumRPCCallDuration)
			defer cancelCtx()

			err = client.CallContext(timeoutContext, &resp, method, args...)
			return resp, err
		})
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

	rateLimitAllowed, requestUUID := w.RateLimiter.Allow(gethcommon.Address(userID))
	defer w.RateLimiter.SetRequestEnd(gethcommon.Address(userID), requestUUID)
	if !rateLimitAllowed {
		return nil, fmt.Errorf("rate limit exceeded")
	}

	cacheArgs := []any{userID, method}
	cacheArgs = append(cacheArgs, args...)

	res, err := withCache(w.Cache, cfg.cacheCfg, generateCacheKey(cacheArgs), func() (*R, error) {
		user, err := getUser(userID, w)
		if err != nil {
			return nil, err
		}

		// determine candidate "from"
		candidateAccts, err := getCandidateAccounts(user, w, cfg)
		if err != nil {
			return nil, err
		}
		if len(candidateAccts) == 0 {
			return nil, fmt.Errorf("illegal access")
		}

		var rpcErr error
		for i := range candidateAccts {
			acct := candidateAccts[i]
			result, err := withEncRPCConnection(ctx, w, acct, func(rpcClient *tenrpc.EncRPCClient) (*R, error) {
				var result *R
				adjustedArgs := args
				if cfg.adjustArgs != nil {
					adjustedArgs = cfg.adjustArgs(acct)
				}

				// wrap the context with a timeout to prevent long executions
				deadline := cfg.timeout
				// if not set, use default
				if deadline == 0 {
					deadline = maximumRPCCallDuration
				}
				timeoutContext, cancelCtx := context.WithTimeout(ctx, deadline)
				defer cancelCtx()

				err := rpcClient.CallContext(timeoutContext, &result, method, adjustedArgs...)
				// return a friendly error to the user
				if err != nil && errors.Is(err, context.DeadlineExceeded) {
					return nil, fmt.Errorf(serverBusy)
				}
				return result, err
			})
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
	userID := gethcommon.FromHex(token)
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

	cacheType := cfg.CacheType
	if cfg.CacheTypeDynamic != nil {
		cacheType = cfg.CacheTypeDynamic()
	}

	if cacheType == NoCache {
		return onCacheMiss()
	}

	// we implement a custom cache eviction logic for the cache strategy of type LatestBatch.
	// when a new batch is created, all entries with "LatestBatch" are considered evicted.
	// elements not cached for a specific batch are not evicted
	isEvicted := false
	ttl := longCacheTTL
	if cacheType == LatestBatch {
		ttl = shortCacheTTL
		isEvicted = cache.IsEvicted(cacheKey, ttl)
	}

	if !isEvicted {
		cachedValue, foundInCache := cache.Get(cacheKey)
		if foundInCache {
			returnValue, ok := cachedValue.(*R)
			if !ok {
				return nil, fmt.Errorf("unexpected error. Invalid format cached. %v", cachedValue)
			}
			return returnValue, nil
		}
	}

	result, err := onCacheMiss()

	// cache only non-nil values
	if err == nil && result != nil {
		cache.Set(cacheKey, result, ttl)
	}

	return result, err
}

func audit(services *Services, msg string, params ...any) {
	if services.Config.VerboseFlag {
		services.FileLogger.Info(fmt.Sprintf(msg, params...))
	}
}

func cacheBlockNumberOrHash(blockNrOrHash rpc.BlockNumberOrHash) CacheStrategy {
	if blockNrOrHash.BlockNumber != nil && blockNrOrHash.BlockNumber.Int64() <= 0 {
		return LatestBatch
	}
	return LongLiving
}

func cacheBlockNumber(lastBlock rpc.BlockNumber) CacheStrategy {
	if lastBlock > 0 {
		return LongLiving
	}
	return LatestBatch
}

func connectWS(ctx context.Context, account *GWAccount, logger gethlog.Logger) (*tenrpc.EncRPCClient, error) {
	return conn(ctx, account.user.services.rpcWSConnPool, account, logger)
}

func conn(ctx context.Context, p *pool.ObjectPool, account *GWAccount, logger gethlog.Logger) (*tenrpc.EncRPCClient, error) {
	defer core.LogMethodDuration(logger, measure.NewStopwatch(), "get rpc connection")
	connectionObj, err := p.BorrowObject(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot fetch rpc connection to backend node %w", err)
	}
	conn := connectionObj.(*rpc.Client)
	encClient, err := wecommon.CreateEncClient(conn, account.address.Bytes(), account.user.userKey, account.signature, account.signatureType, logger)
	if err != nil {
		_ = returnConn(p, conn, logger)
		return nil, fmt.Errorf("error creating new client, %w", err)
	}
	return encClient, nil
}

func returnConn(p *pool.ObjectPool, conn tenrpc.Client, logger gethlog.Logger) error {
	err := p.ReturnObject(context.Background(), conn)
	if err != nil {
		logger.Error("Error returning connection to pool", log.ErrKey, err)
	}
	return err
}

func withEncRPCConnection[R any](ctx context.Context, w *Services, acct *GWAccount, execute func(*tenrpc.EncRPCClient) (*R, error)) (*R, error) {
	rpcClient, err := conn(ctx, acct.user.services.rpcHTTPConnPool, acct, w.logger)
	if err != nil {
		return nil, fmt.Errorf("could not connect to backed. Cause: %w", err)
	}
	defer returnConn(w.rpcHTTPConnPool, rpcClient.BackingClient(), w.logger)
	return execute(rpcClient)
}

func withPlainRPCConnection[R any](ctx context.Context, w *Services, execute func(client *rpc.Client) (*R, error)) (*R, error) {
	connectionObj, err := w.rpcHTTPConnPool.BorrowObject(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot fetch rpc connection to backend node %w", err)
	}
	rpcClient := connectionObj.(*rpc.Client)
	defer returnConn(w.rpcHTTPConnPool, rpcClient, w.logger)
	return execute(rpcClient)
}
