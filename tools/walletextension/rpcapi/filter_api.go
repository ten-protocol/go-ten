package rpcapi

import (
	"context"
	"fmt"
	"reflect"
	"sync/atomic"
	"time"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ten-protocol/go-ten/go/common"

	wecommon "github.com/ten-protocol/go-ten/tools/walletextension/common"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/filters"
	"github.com/ten-protocol/go-ten/lib/gethfork/rpc"
)

type FilterAPI struct {
	we *Services
}

func NewFilterAPI(we *Services) *FilterAPI {
	return &FilterAPI{we: we}
}

func (api *FilterAPI) NewPendingTransactionFilter(_ *bool) rpc.ID {
	return "not supported"
}

/*func (api *FilterAPI) NewPendingTransactions(ctx context.Context, fullTx *bool) (*rpc.Subscription, error) {
	// not supported
	return nil, nil
}
*/

func (api *FilterAPI) NewBlockFilter() rpc.ID {
	// not implemented
	return ""
}

/*
	func (api *FilterAPI) NewHeads(ctx context.Context) (*rpc.Subscription, error) {
		// not implemented
		return nil, nil
	}
*/
func (api *FilterAPI) Logs(ctx context.Context, crit filters.FilterCriteria) (*rpc.Subscription, error) {
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

	inputChannels := make([]chan common.IDAndLog, 0)
	backendSubscriptions := make([]*rpc.ClientSubscription, 0)
	for _, address := range candidateAddresses {
		rpcWSClient, err := user.accounts[*address].connect(api.we.HostAddrWS, api.we.Logger())
		if err != nil {
			return nil, err
		}

		inCh := make(chan common.IDAndLog)
		backendSubscription, err := rpcWSClient.Subscribe(ctx, nil, "eth", inCh, "logs", crit)
		if err != nil {
			return nil, err
		}

		inputChannels = append(inputChannels, inCh)
		backendSubscriptions = append(backendSubscriptions, backendSubscription)
	}

	dedupeBuffer := NewCircularBuffer(wecommon.DeduplicationBufferSize)
	subscription := subNotifier.CreateSubscription()

	unsubscribed := atomic.Bool{}
	go forwardAndDedupe(inputChannels, backendSubscriptions, subscription, subNotifier, &unsubscribed, func(data common.IDAndLog) *types.Log {
		uniqueLogKey := LogKey{
			BlockHash: data.Log.BlockHash,
			TxHash:    data.Log.TxHash,
			Index:     data.Log.Index,
		}

		if !dedupeBuffer.Contains(uniqueLogKey) {
			dedupeBuffer.Push(uniqueLogKey)
			return data.Log
		}
		return nil
	})

	go handleUnsubscribe(subscription, backendSubscriptions, &unsubscribed)

	return subscription, err
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

	user, err := getUser(subNotifier.UserID, api.we.Storage)
	if err != nil {
		return nil, nil, fmt.Errorf("illegal access: %s, %w", subNotifier.UserID, err)
	}
	return subNotifier, user, nil
}

func searchForAddressInFilterCriteria(filterCriteria filters.FilterCriteria, possibleAddresses []*gethcommon.Address) []*gethcommon.Address {
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

// forwardAndDedupe - reads messages from the input channel, and forwards them to the notifier only if they are new
func forwardAndDedupe[R any, T any](inputChannels []chan R, _ []*rpc.ClientSubscription, outSub *rpc.Subscription, notifier *rpc.Notifier, unsubscribed *atomic.Bool, toForward func(elem R) *T) {
	inputCases := make([]reflect.SelectCase, len(inputChannels)+1)

	// create a ticker to handle cleanup
	inputCases[0] = reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(time.NewTicker(10 * time.Second).C),
	}

	// create a select "case" for each input channel
	for i, ch := range inputChannels {
		inputCases[i+1] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)}
	}

	unclosedInputChannels := len(inputCases)
	for unclosedInputChannels > 0 {
		chosen, value, ok := reflect.Select(inputCases)
		if !ok {
			// The chosen channel has been closed, so zero out the channel to disable the case
			inputCases[chosen].Chan = reflect.ValueOf(nil)
			unclosedInputChannels--
			continue
		}

		switch v := value.Interface().(type) {
		case time.Time:
			// exit the loop
			if unsubscribed.Load() {
				return
			}
		case R:
			valueToSubmit := toForward(v)
			if valueToSubmit != nil {
				err := notifier.Notify(outSub.ID, *valueToSubmit)
				if err != nil {
					return
				}
			}
		default:
			// unexpected element received
			continue
		}
	}
}

func handleUnsubscribe(connectionSub *rpc.Subscription, backendSubscriptions []*rpc.ClientSubscription, unsubscribed *atomic.Bool) {
	<-connectionSub.Err()
	for _, backendSub := range backendSubscriptions {
		backendSub.Unsubscribe()
		unsubscribed.Store(true)
	}
}

/*
	func (api *FilterAPI) NewFilter(crit filters.FilterCriteria) (rpc.ID, error) {
		// not implemented
		return "", nil
	}
*/
func (api *FilterAPI) GetLogs(ctx context.Context, crit filters.FilterCriteria) ([]*types.Log, error) {
	logs, err := ExecAuthRPC[[]*types.Log](
		ctx,
		api.we,
		&ExecCfg{
			cacheCfg: &CacheCfg{
				TTLCallback: func() time.Duration {
					// when the toBlock is not specified, the request is open-ended
					if crit.ToBlock != nil {
						return longCacheTTL
					}
					return shortCacheTTL
				},
			},
			tryUntilAuthorised: true,
			adjustArgs: func(acct *GWAccount) []any {
				// convert to something serializable
				return []any{common.FromCriteria(crit)}
			},
		},
		"eth_getLogs",
		crit,
	)
	if logs != nil {
		return *logs, err
	}
	return nil, err
}

/*func (api *FilterAPI) UninstallFilter(id rpc.ID) bool {
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
	// not implemented
	return nil, nil
}
*/
