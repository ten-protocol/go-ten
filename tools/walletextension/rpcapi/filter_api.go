package rpcapi

import (
	"context"
	"fmt"
	"sort"
	"sync/atomic"
	"time"

	"github.com/status-im/keycard-go/hexutils"

	"github.com/ethereum/go-ethereum/log"

	subscriptioncommon "github.com/ten-protocol/go-ten/go/common/subscription"

	tenrpc "github.com/ten-protocol/go-ten/go/rpc"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ten-protocol/go-ten/go/common"

	wecommon "github.com/ten-protocol/go-ten/tools/walletextension/common"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ten-protocol/go-ten/lib/gethfork/rpc"
)

type FilterAPI struct {
	we     *Services
	logger log.Logger
}

func NewFilterAPI(we *Services) *FilterAPI {
	return &FilterAPI{
		we:     we,
		logger: we.logger,
	}
}

func (api *FilterAPI) NewPendingTransactionFilter(_ *bool) rpc.ID {
	return "not supported"
}

func (api *FilterAPI) NewPendingTransactions(ctx context.Context, fullTx *bool) (*rpc.Subscription, error) {
	return nil, fmt.Errorf("not supported")
}

func (api *FilterAPI) NewBlockFilter() rpc.ID {
	// not implemented
	return ""
}

func (api *FilterAPI) NewHeads(ctx context.Context) (*rpc.Subscription, error) {
	notifier, supported := rpc.NotifierFromContext(ctx)
	if !supported {
		return nil, fmt.Errorf("creation of subscriptions is not supported")
	}
	subscription := notifier.CreateSubscription()
	api.we.NewHeadsService.RegisterNotifier(notifier, subscription)
	return subscription, nil
}

func (api *FilterAPI) Logs(ctx context.Context, crit common.FilterCriteria) (*rpc.Subscription, error) {
	audit(api.we, "start Logs subscription %v", crit)
	subNotifier, user, err := getUserAndNotifier(ctx, api)
	if err != nil {
		return nil, err
	}

	// determine the accounts to use for the backend subscriptions
	candidateAddresses := user.GetAllAddresses()
	if len(candidateAddresses) > 1 {
		candidateAddresses = searchForAddressInFilterCriteria(crit, user.GetAllAddresses())
		// when we can't determine which addresses to use based on the criteria, use all of them
		if len(candidateAddresses) == 0 {
			candidateAddresses = user.GetAllAddresses()
		}
	}

	backendWSConnections := make([]*tenrpc.EncRPCClient, 0)
	inputChannels := make([]chan types.Log, 0)
	errorChannels := make([]<-chan error, 0)
	backendSubscriptions := make([]*rpc.ClientSubscription, 0)
	for _, address := range candidateAddresses {
		rpcWSClient, err := connectWS(ctx, user.accounts[*address], api.we.Logger())
		if err != nil {
			return nil, err
		}
		backendWSConnections = append(backendWSConnections, rpcWSClient)

		inCh := make(chan types.Log)
		backendSubscription, err := rpcWSClient.Subscribe(ctx, "eth", inCh, "logs", crit)
		if err != nil {
			fmt.Printf("could not connect to backend %s", err)
			return nil, err
		}

		inputChannels = append(inputChannels, inCh)
		errorChannels = append(errorChannels, backendSubscription.Err())
		backendSubscriptions = append(backendSubscriptions, backendSubscription)
	}

	dedupeBuffer := NewCircularBuffer(wecommon.DeduplicationBufferSize)
	subscription := subNotifier.CreateSubscription()

	unsubscribedByClient := atomic.Bool{}
	unsubscribedByBackend := atomic.Bool{}
	go subscriptioncommon.ForwardFromChannels(
		inputChannels,
		func(log types.Log) error {
			uniqueLogKey := LogKey{
				BlockHash: log.BlockHash,
				TxHash:    log.TxHash,
				Index:     log.Index,
			}

			if !dedupeBuffer.Contains(uniqueLogKey) {
				dedupeBuffer.Push(uniqueLogKey)
				return subNotifier.Notify(subscription.ID, log)
			}
			return nil
		},
		func() {
			// release resources
			api.closeConnections(backendSubscriptions, backendWSConnections)
		}, // todo - we can implement reconnect logic here
		&unsubscribedByBackend,
		&unsubscribedByClient,
		12*time.Hour,
		api.logger,
	)

	// handles any of the backend connections being closed
	go subscriptioncommon.HandleUnsubscribeErrChan(errorChannels, func() {
		unsubscribedByBackend.Store(true)
	})

	// handles "unsubscribe" from the user
	go subscriptioncommon.HandleUnsubscribe(subscription, func() {
		unsubscribedByClient.Store(true)
		api.closeConnections(backendSubscriptions, backendWSConnections)
	})

	return subscription, err
}

func (api *FilterAPI) closeConnections(backendSubscriptions []*rpc.ClientSubscription, backendWSConnections []*tenrpc.EncRPCClient) {
	for _, backendSub := range backendSubscriptions {
		backendSub.Unsubscribe()
	}
	for _, connection := range backendWSConnections {
		_ = returnConn(api.we.rpcWSConnPool, connection.BackingClient(), api.logger)
	}
}

func getUserAndNotifier(ctx context.Context, api *FilterAPI) (*rpc.Notifier, *GWUser, error) {
	subNotifier, supported := rpc.NotifierFromContext(ctx)
	if !supported {
		return nil, nil, fmt.Errorf("creation of subscriptions is not supported")
	}

	// todo - we might want to allow access to public logs
	if len(subNotifier.UserID) == 0 {
		return nil, nil, fmt.Errorf("illegal access")
	}

	user, err := getUser(subNotifier.UserID, api.we)
	if err != nil {
		return nil, nil, fmt.Errorf("illegal access: %s, %w", subNotifier.UserID, err)
	}
	return subNotifier, user, nil
}

func searchForAddressInFilterCriteria(filterCriteria common.FilterCriteria, possibleAddresses []*gethcommon.Address) []*gethcommon.Address {
	result := make([]*gethcommon.Address, 0)
	addrMap := toMap(possibleAddresses)
	for _, topicCondition := range filterCriteria.Topics {
		for _, topic := range topicCondition {
			potentialAddr := common.ExtractPotentialAddress(topic)
			if potentialAddr != nil && addrMap[*potentialAddr] != nil {
				result = append(result, potentialAddr)
			}
		}
	}
	return result
}

func (api *FilterAPI) NewFilter(crit common.FilterCriteria) (rpc.ID, error) {
	return rpc.NewID(), rpcNotImplemented
}

func (api *FilterAPI) GetLogs(ctx context.Context, crit common.FilterCriteria) ([]*types.Log, error) {
	method := "eth_getLogs"
	audit(api.we, "RPC start method=%s args=%v", method, ctx)
	requestStartTime := time.Now()
	userID, err := extractUserID(ctx, api.we)
	if err != nil {
		return nil, err
	}

	rateLimitAllowed, requestUUID := api.we.RateLimiter.Allow(gethcommon.Address(userID))
	if !rateLimitAllowed {
		return nil, fmt.Errorf("rate limit exceeded")
	}

	res, err := withCache(
		api.we.Cache,
		&CacheCfg{
			CacheTypeDynamic: func() CacheStrategy {
				if crit.ToBlock != nil && crit.ToBlock.Int64() > 0 {
					return LongLiving
				}
				if crit.BlockHash != nil {
					return LongLiving
				}
				// when the toBlock or the block Hash are not specified, the request is open-ended
				return LatestBatch
			},
		},
		generateCacheKey([]any{userID, method, common.SerializableFilterCriteria(crit)}),
		func() (*[]*types.Log, error) { // called when there is no entry in the cache
			user, err := getUser(userID, api.we)
			if err != nil {
				return nil, err
			}

			allEventLogsMap := make(map[LogKey]*types.Log)
			// for each account registered for the current user
			// execute the get_Logs function
			// dedupe and concatenate the results
			for _, acct := range user.accounts {
				eventLogs, err := withEncRPCConnection(ctx, api.we, acct, func(rpcClient *tenrpc.EncRPCClient) (*[]*types.Log, error) {
					var result []*types.Log

					// wrap the context with a timeout to prevent long executions
					timeoutContext, cancelCtx := context.WithTimeout(ctx, maximumRPCCallDuration)
					defer cancelCtx()

					err := rpcClient.CallContext(timeoutContext, &result, method, common.SerializableFilterCriteria(crit))
					return &result, err
				})
				if err != nil {
					return nil, fmt.Errorf("could not read logs. cause: %w", err)
				}
				// dedupe event logs
				for _, eventLog := range *eventLogs {
					allEventLogsMap[LogKey{
						BlockHash: eventLog.BlockHash,
						TxHash:    eventLog.TxHash,
						Index:     eventLog.Index,
					}] = eventLog
				}
			}

			result := make([]*types.Log, 0)
			for _, eventLog := range allEventLogsMap {
				result = append(result, eventLog)
			}
			sort.Slice(result, func(i, j int) bool {
				if result[i].BlockNumber == result[j].BlockNumber {
					return result[i].Index < result[j].Index
				}
				return result[i].BlockNumber < result[j].BlockNumber
			})
			return &result, nil
		})
	if err != nil {
		return nil, err
	}
	api.we.RateLimiter.SetRequestEnd(gethcommon.Address(userID), requestUUID)
	audit(api.we, "RPC call. uid=%s, method=%s args=%v result=%v error=%v time=%d", hexutils.BytesToHex(userID), method, crit, res, err, time.Since(requestStartTime).Milliseconds())
	return *res, err
}

func (api *FilterAPI) UninstallFilter(id rpc.ID) bool {
	// not implemented
	return false
}

func (api *FilterAPI) GetFilterLogs(ctx context.Context, id rpc.ID) ([]*types.Log, error) {
	//txRec, err := ExecAuthRPC[[]*types.Log](ctx, api.we, "GetFilterLogs", ExecCfg{account: args.From}, id)
	//if txRec != nil {
	//	return *txRec, err
	//}
	//return common.Hash{}, err

	// not implemented
	return nil, nil
}

func (api *FilterAPI) GetFilterChanges(id rpc.ID) (interface{}, error) {
	return nil, rpcNotImplemented
}
