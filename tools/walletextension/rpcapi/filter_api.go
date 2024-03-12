package rpcapi

import (
	"context"
	"fmt"
	"reflect"
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

/*
	func (api *FilterAPI) NewPendingTransactions(ctx context.Context, fullTx *bool) (*rpc.Subscription, error) {
		// not supported
		return nil, nil
	}
*/
func (api *FilterAPI) NewBlockFilter() rpc.ID {
	// not implemented
	return ""
}

/*func (api *FilterAPI) NewHeads(ctx context.Context) (*rpc.Subscription, error) {
	// not implemented
	return nil, nil
}
*/
// todo - unsubscribe
func (api *FilterAPI) Logs(ctx context.Context, crit filters.FilterCriteria) (*rpc.Subscription, error) {
	not, ok := rpc.NotifierFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("invalid subscription")
	}

	uid, err := wecommon.GetUserIDbyte(not.UserID)
	if err != nil {
		return nil, fmt.Errorf("invald token: %s, %w", not.UserID, err)
	}

	user, err := getUser(uid, api.we.Storage)

	notifier, supported := rpc.NotifierFromContext(ctx)
	if !supported {
		return nil, fmt.Errorf("creation of subscriptions is not supported")
	}
	subscription := notifier.CreateSubscription()

	candidateAddresses := user.GetAllAddresses()
	if len(candidateAddresses) > 1 {
		candidateAddresses = searchForAddressInFilterCriteria(crit, user.GetAllAddresses())
		if len(candidateAddresses) == 0 {
			candidateAddresses = user.GetAllAddresses()
		}
	}
	inputChannels := make([]chan common.IDAndLog, 0)
	clientSubscriptions := make([]*rpc.ClientSubscription, 0)
	for _, address := range candidateAddresses {
		rpcWSClient, err := user.accounts[*address].connect(api.we.HostAddrWS, api.we.Logger())
		if err != nil {
			return nil, err
		}
		inCh := make(chan common.IDAndLog)
		inputChannels = append(inputChannels, inCh)
		clientSubscription, err := rpcWSClient.Subscribe(ctx, nil, "eth", inCh, "logs", crit)
		if err != nil {
			return nil, err
		}
		clientSubscriptions = append(clientSubscriptions, clientSubscription)
	}
	go forwardMsgs(inputChannels, clientSubscriptions, subscription, notifier)
	return subscription, err
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

func forwardMsgs(inputChannels []chan common.IDAndLog, _ []*rpc.ClientSubscription, outSub *rpc.Subscription, notifier *rpc.Notifier) {
	buffer := NewCircularBuffer(wecommon.DeduplicationBufferSize)
	cases := make([]reflect.SelectCase, len(inputChannels))
	for i, ch := range inputChannels {
		cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)}
	}
	remaining := len(cases)
	for remaining > 0 {
		chosen, value, ok := reflect.Select(cases)
		if !ok {
			// The chosen channel has been closed, so zero out the channel to disable the case
			cases[chosen].Chan = reflect.ValueOf(nil)
			remaining--
			continue
		}

		data := value.Interface().(common.IDAndLog)
		uniqueLogKey := LogKey{
			BlockHash: data.Log.BlockHash,
			TxHash:    data.Log.TxHash,
			Index:     data.Log.Index,
		}

		// check if the current event is a duplicate (and skip it if it is)
		if !buffer.Contains(uniqueLogKey) {
			buffer.Push(uniqueLogKey)
			err := notifier.Notify(outSub.ID, data.Log)
			if err != nil {
				println(err)
				return
			}
		}
	}
}

/*
	func (api *FilterAPI) NewFilter(crit filters.FilterCriteria) (rpc.ID, error) {
		// not implemented
		return "", nil
	}
*/
func (api *FilterAPI) GetLogs(ctx context.Context, crit filters.FilterCriteria) ([]*types.Log, error) {
	// todo
	logs, err := ExecAuthRPC[[]*types.Log](
		ctx,
		api.we,
		&ExecCfg{
			cacheCfg: &CacheCfg{
				TTLCallback: func() time.Duration {
					if crit.ToBlock != nil {
						return longCacheTTL
					}
					return shortCacheTTL
				},
			},
			tryUntilAuthorised: true,
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
