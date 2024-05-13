package rpcapi

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	subscriptioncommon "github.com/ten-protocol/go-ten/go/common/subscription"

	tenrpc "github.com/ten-protocol/go-ten/go/rpc"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ten-protocol/go-ten/go/common"

	wecommon "github.com/ten-protocol/go-ten/tools/walletextension/common"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ten-protocol/go-ten/lib/gethfork/rpc"
)

type FilterAPI struct {
	we *Services
}

func NewFilterAPI(we *Services) *FilterAPI {
	return &FilterAPI{
		we: we,
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
	backendSubscriptions := make([]*rpc.ClientSubscription, 0)
	for _, address := range candidateAddresses {
		rpcWSClient, err := connectWS(user.accounts[*address], api.we.Logger())
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
		backendSubscriptions = append(backendSubscriptions, backendSubscription)
	}

	dedupeBuffer := NewCircularBuffer(wecommon.DeduplicationBufferSize)
	subscription := subNotifier.CreateSubscription()

	unsubscribed := atomic.Bool{}
	go subscriptioncommon.ForwardFromChannels(
		inputChannels,
		&unsubscribed,
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
		nil,
		12*time.Hour,
	)

	go subscriptioncommon.HandleUnsubscribe(subscription, &unsubscribed, func() {
		for _, backendSub := range backendSubscriptions {
			backendSub.Unsubscribe()
		}
		for _, connection := range backendWSConnections {
			_ = returnConn(api.we.rpcWSConnPool, connection.BackingClient())
		}
		unsubscribed.Store(true)
	})

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
	logs, err := ExecAuthRPC[[]*types.Log](
		ctx,
		api.we,
		&ExecCfg{
			cacheCfg: &CacheCfg{
				CacheTypeDynamic: func() CacheStrategy {
					// when the toBlock is not specified, the request is open-ended
					if crit.ToBlock != nil && crit.ToBlock.Int64() > 0 {
						return LongLiving
					}
					return LatestBatch
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
