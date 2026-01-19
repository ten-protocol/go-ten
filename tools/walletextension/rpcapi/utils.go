package rpcapi

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ten-protocol/go-ten/go/common/log"
	tencommonrpc "github.com/ten-protocol/go-ten/go/common/rpc"

	"github.com/ten-protocol/go-ten/tools/walletextension/common"

	"github.com/ten-protocol/go-ten/tools/walletextension/cache"

	"github.com/ten-protocol/go-ten/tools/walletextension/services"

	"github.com/status-im/keycard-go/hexutils"

	"github.com/ten-protocol/go-ten/go/common/viewingkey"
	tenrpc "github.com/ten-protocol/go-ten/go/rpc"

	"github.com/ten-protocol/go-ten/lib/gethfork/rpc"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

const (
	ethCallPaddedArgLen = 64
	ethCallAddrPadding  = "000000000000000000000000"

	notAuthorised = "not authorised"
	serverBusy    = "server busy. please retry later"

	// hardcoding the maximum time for an RPC request
	// this value will be propagated to the node and enclave and all the operations
	maximumRPCCallDuration  = 5 * time.Second
	sendTransactionDuration = 20 * time.Second
)

var (
	ErrRPCNotImplemented          = errors.New("rpc endpoint not implemented")
	ErrAuthenticationTokenMissing = errors.New("authentication token missing")
)

type AuthExecCfg struct {
	// these 4 fields specify the account(s) that should make the backend call
	account             *gethcommon.Address
	computeFromCallback func(user *common.GWUser) *gethcommon.Address
	tryAll              bool
	tryUntilAuthorised  bool

	adjustArgs func(acct *common.GWAccount) []any
	cacheCfg   *cache.Cfg
	timeout    time.Duration
}

func SendRawTx(ctx context.Context, w *services.Services, input hexutil.Bytes) (gethcommon.Hash, error) {
	txRec, err := ExecAuthRPC[gethcommon.Hash](ctx, w, &AuthExecCfg{tryAll: true, timeout: sendTransactionDuration}, tencommonrpc.ERPCSendRawTransaction, input)
	if err != nil {
		return gethcommon.Hash{}, err
	}
	return *txRec, err
}

func UnauthenticatedTenRPCCall[R any](ctx context.Context, w *services.Services, cfg *cache.Cfg, method string, args ...any) (*R, error) {
	if ctx == nil {
		return nil, errors.New("invalid call. nil Context")
	}
	w.Logger().Debug("RPC start", "method", method, "args", SafeArgsForLogging(args))
	requestStartTime := time.Now()
	cacheArgs := []any{method}
	cacheArgs = append(cacheArgs, args...)

	res, err := cache.WithCache(w.RPCResponsesCache, cfg, generateCacheKey(cacheArgs), func() (*R, error) {
		return services.WithPlainRPCConnection(ctx, w.BackendRPC, func(client *rpc.Client) (*R, error) {
			resp := new(R)
			var err error

			// wrap the context with a timeout to prevent long executions
			timeoutContext, cancelCtx := context.WithTimeout(ctx, maximumRPCCallDuration)
			defer cancelCtx()

			err = client.CallContext(timeoutContext, &resp, method, args...)
			return resp, err
		})
	})
	if err != nil {
		w.Logger().Error("RPC call failed", "method", method, "args", SafeArgsForLogging(args), "err", err, "time", time.Since(requestStartTime).Milliseconds())
		return nil, err
	}

	w.Logger().Info("RPC call succeeded", "method", method, "args", SafeArgsForLogging(args), "result", SafeValueForLogging(res), "time", time.Since(requestStartTime).Milliseconds())
	return res, err
}

func ExecAuthRPC[R any](ctx context.Context, w *services.Services, cfg *AuthExecCfg, method string, args ...any) (*R, error) {
	w.Logger().Debug("RPC start", "method", method, "args", SafeArgsForLogging(args))
	requestStartTime := time.Now()

	// get the user from the request
	user, err := extractUserForRequest(ctx, w)

	switch err {
	case nil:
		// proced with the user from the request
	case ErrAuthenticationTokenMissing:
		// use the default user for public access & return error if not found
		user = w.DefaultUser
		if user == nil {
			w.Logger().Warn("Default user not found")
			return nil, errors.New("default user not found")
		}
	default:
		// return the error
		return nil, err
	}

	w.MetricsTracker.RecordUserActivity(user.ID)

	// use rate limiting for non-default users
	if user != w.DefaultUser {
		rateLimitAllowed, requestUUID := w.RateLimiter.Allow(gethcommon.Address(user.ID))
		if !rateLimitAllowed {
			w.Logger().Warn("Rate limit exceeded for user", "userID", hexutils.BytesToHex(user.ID))
			return nil, errors.New("rate limit exceeded")
		}
		defer w.RateLimiter.SetRequestEnd(gethcommon.Address(user.ID), requestUUID)
	}

	cacheArgs := []any{user.ID, method}
	cacheArgs = append(cacheArgs, args...)

	res, err := cache.WithCache(w.RPCResponsesCache, cfg.cacheCfg, generateCacheKey(cacheArgs), func() (*R, error) {
		// determine candidate "from"
		candidateAccts, err := getCandidateAccounts(user, w, cfg)
		if err != nil {
			return nil, err
		}
		if len(candidateAccts) == 0 {
			return nil, errors.New("illegal access")
		}

		var rpcErr error
		for i := range candidateAccts {
			acct := candidateAccts[i]
			result, err := services.WithEncRPCConnection(ctx, w.BackendRPC, acct, func(rpcClient *tenrpc.EncRPCClient) (*R, error) {
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
					return nil, errors.New(serverBusy)
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
	w.Logger().Info("RPC call", "uid", hexutils.BytesToHex(user.ID), "method", method, "args", SafeArgsForLogging(args), "result", SafeValueForLogging(res), "err", err, "time", time.Since(requestStartTime).Milliseconds())
	return res, err
}

func getCandidateAccounts(user *common.GWUser, we *services.Services, cfg *AuthExecCfg) ([]*common.GWAccount, error) {
	candidateAccts := make([]*common.GWAccount, 0)
	// for users with multiple accounts try to determine a candidate account based on the available information
	switch {
	case cfg.account != nil:
		acc := user.AllAccounts()[*cfg.account]
		if acc != nil {
			candidateAccts = append(candidateAccts, acc)
			return candidateAccts, nil
		}

	case cfg.computeFromCallback != nil:
		suggestedAddress := cfg.computeFromCallback(user)
		if suggestedAddress != nil {
			acc := user.AllAccounts()[*suggestedAddress]
			if acc != nil {
				candidateAccts = append(candidateAccts, acc)
				return candidateAccts, nil
			} else {
				// this can only happen when the "from" is not one of the registered accounts.
				return nil, fmt.Errorf("account: %s not registered to current user. Please register first", suggestedAddress.Hex())
			}
		}
	}

	if cfg.tryAll || cfg.tryUntilAuthorised {
		for _, acc := range user.AllAccounts() {
			candidateAccts = append(candidateAccts, acc)
		}
	}

	return candidateAccts, nil
}

func extractUserID(ctx context.Context, _ *services.Services) ([]byte, error) {
	token, ok := ctx.Value(rpc.GWTokenKey{}).(string)
	if !ok {
		return nil, fmt.Errorf("invalid authentication token: %s", ctx.Value(rpc.GWTokenKey{}))
	}
	if len(strings.TrimSpace(token)) == 0 {
		return nil, ErrAuthenticationTokenMissing
	}
	userID := gethcommon.FromHex(token)
	if len(userID) != viewingkey.UserIDLength {
		return nil, fmt.Errorf("invalid authentication token: %s", token)
	}
	return userID, nil
}

func extractUserForRequest(ctx context.Context, w *services.Services) (*common.GWUser, error) {
	userID, err := extractUserID(ctx, w)
	if err != nil {
		return nil, err
	}
	user, err := w.Storage.GetUser(userID)
	if err != nil {
		w.Logger().Info("authentication failed: user not found", "userID", hex.EncodeToString(userID), log.ErrKey, err)
		return nil, fmt.Errorf("authentication failed (userID=%s): %w", hex.EncodeToString(userID), err)
	}
	return user, nil
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

func cacheBlockNumberOrHash(blockNrOrHash rpc.BlockNumberOrHash) cache.Strategy {
	if blockNrOrHash.BlockNumber != nil && blockNrOrHash.BlockNumber.Int64() <= 0 {
		return cache.LatestBatch
	}
	return cache.LongLiving
}

func cacheBlockNumber(lastBlock rpc.BlockNumber) cache.Strategy {
	if lastBlock > 0 {
		return cache.LongLiving
	}
	return cache.LatestBatch
}

// SafeArgsForLogging replaces nil fmt.Stringer pointers in args with "<nil>" to prevent segfaults
func SafeArgsForLogging(args []any) string {
	if len(args) == 0 {
		return "[]"
	}
	var b strings.Builder
	b.WriteByte('[')
	for i, arg := range args {
		if i > 0 {
			b.WriteString(", ")
		}
		safeStringify(&b, reflect.ValueOf(arg), 0)
	}
	b.WriteByte(']')
	return b.String()
}

// SafeValueForLogging safely converts a value to a string for logging, handling nil pointers
func SafeValueForLogging(v any) string {
	var b strings.Builder
	safeStringify(&b, reflect.ValueOf(v), 0)
	return b.String()
}

const (
	// maxDepth limits the recursion depth when stringifying nested structures to prevent infinite loops
	maxDepth = 5
	// maxSliceElements limits the number of slice elements to log before truncating with "...+X more"
	maxSliceElements = 10
	// maxMapEntries limits the number of map key-value pairs to log before truncating
	maxMapEntries = 5
	// maxStructFields limits the number of struct fields to log before truncating
	maxStructFields = 10
)

func safeStringify(b *strings.Builder, rv reflect.Value, depth int) {
	if depth > maxDepth {
		b.WriteString("...")
		return
	}

	if !rv.IsValid() {
		b.WriteString("<nil>")
		return
	}

	switch rv.Kind() {
	case reflect.Interface:
		if rv.IsZero() {
			fmt.Fprintf(b, "<%s nil>", rv.Type())
			return
		}
		safeStringify(b, rv.Elem(), depth+1)
		return

	case reflect.Ptr:
		if rv.IsNil() {
			fmt.Fprintf(b, "<%s nil>", rv.Type())
			return
		}
		safeStringify(b, rv.Elem(), depth+1)
		return

	case reflect.Slice:
		if rv.IsNil() {
			fmt.Fprintf(b, "<%s nil>", rv.Type())
			return
		}
		b.WriteByte('[')
		for i := 0; i < rv.Len() && i < maxSliceElements; i++ {
			if i > 0 {
				b.WriteString(", ")
			}
			safeStringify(b, rv.Index(i), depth+1)
		}
		if rv.Len() > maxSliceElements {
			fmt.Fprintf(b, "...+%d more", rv.Len()-maxSliceElements)
		}
		b.WriteByte(']')
		return

	case reflect.Map:
		if rv.IsNil() {
			fmt.Fprintf(b, "<%s nil>", rv.Type())
			return
		}
		b.WriteString("map[")
		iter := rv.MapRange()
		count := 0
		for iter.Next() && count < maxMapEntries {
			if count > 0 {
				b.WriteString(", ")
			}
			safeStringify(b, iter.Key(), depth+1)
			b.WriteByte(':')
			safeStringify(b, iter.Value(), depth+1)
			count++
		}
		b.WriteByte(']')
		return

	case reflect.Struct:
		b.WriteByte('{')
		t := rv.Type()
		for i := 0; i < rv.NumField() && i < maxStructFields; i++ {
			if i > 0 {
				b.WriteString(", ")
			}
			b.WriteString(t.Field(i).Name)
			b.WriteByte(':')
			field := rv.Field(i)
			if field.CanInterface() {
				safeStringify(b, field, depth+1)
			} else {
				b.WriteString("<unexported>")
			}
		}
		b.WriteByte('}')
		return

	case reflect.Chan, reflect.Func:
		if rv.IsNil() {
			fmt.Fprintf(b, "<%s nil>", rv.Type())
		} else {
			fmt.Fprintf(b, "<%s>", rv.Type())
		}
		return

	case reflect.String:
		fmt.Fprintf(b, "%q", rv.String())
	case reflect.Bool:
		fmt.Fprintf(b, "%t", rv.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fmt.Fprintf(b, "%d", rv.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		fmt.Fprintf(b, "%d", rv.Uint())
	case reflect.Float32, reflect.Float64:
		fmt.Fprintf(b, "%g", rv.Float())
	default:
		fmt.Fprintf(b, "<%s>", rv.Type())
	}
}
