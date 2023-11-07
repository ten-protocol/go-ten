package subscriptions

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/go-kit/kit/transport/http/jsonrpc"

	gethlog "github.com/ethereum/go-ethereum/log"
	gethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/rpc"
	wecommon "github.com/obscuronet/go-obscuro/tools/walletextension/common"
	"github.com/obscuronet/go-obscuro/tools/walletextension/userconn"
)

type SubscriptionManager struct {
	subscriptionMappings map[string][]*gethrpc.ClientSubscription
	logger               gethlog.Logger
	mu                   sync.Mutex
}

func New(logger gethlog.Logger) *SubscriptionManager {
	return &SubscriptionManager{
		subscriptionMappings: make(map[string][]*gethrpc.ClientSubscription),
		logger:               logger,
	}
}

// HandleNewSubscriptions subscribes to an event with all the clients provided.
// Doing this is necessary because we have relevancy rule, and we want to subscribe sometimes with all clients to get all the events
func (sm *SubscriptionManager) HandleNewSubscriptions(clients []rpc.Client, req *wecommon.RPCRequest, resp *interface{}, userConn userconn.UserConn) error {
	if len(req.Params) == 0 {
		return fmt.Errorf("could not subscribe as no subscription namespace was provided")
	}

	sm.logger.Info(fmt.Sprintf("Subscribing to event %s with %d clients", req.Params, len(clients)))

	// create subscriptionID which will enable user to unsubscribe from all subscriptions
	userSubscriptionID := gethrpc.NewID()

	// create a common channel for subscriptions from all accounts
	funnelMultipleAccountsChan := make(chan common.IDAndLog)

	// read from a multiple accounts channel and write results to userConn
	go readFromChannelAndWriteToUserConn(funnelMultipleAccountsChan, userConn, userSubscriptionID, sm.logger)

	// iterate over all clients and subscribe for each of them
	for _, client := range clients {
		subscription, err := client.Subscribe(context.Background(), resp, rpc.SubscribeNamespace, funnelMultipleAccountsChan, req.Params...)
		if err != nil {
			return fmt.Errorf("could not call %s with params %v. Cause: %w", req.Method, req.Params, err)
		}
		sm.UpdateSubscriptionMapping(string(userSubscriptionID), subscription)

		// We periodically check if the websocket is closed, and terminate the subscription.
		go checkIfUserConnIsClosedAndUnsubscribe(userConn, subscription, &sm.mu)
	}

	// We return subscriptionID with resp interface. We want to use userSubscriptionID to allow unsubscribing
	*resp = userSubscriptionID
	return nil
}

func readFromChannelAndWriteToUserConn(channel chan common.IDAndLog, userConn userconn.UserConn, userSubscriptionID gethrpc.ID, logger gethlog.Logger) {
	buffer := NewCircularBuffer(wecommon.DeduplicationBufferSize)
	for data := range channel {
		// create unique identifier for current log
		uniqueLogKey := LogKey{
			BlockHash: data.Log.BlockHash,
			TxHash:    data.Log.TxHash,
			Index:     data.Log.Index,
		}

		// check if the current event is a duplicate (and skip it if it is)
		if buffer.Contains(uniqueLogKey) {
			continue
		}

		jsonResponse, err := prepareLogResponse(data, userSubscriptionID)
		if err != nil {
			logger.Error("could not marshal log response to JSON on subscription.", log.SubIDKey, data.SubID, log.ErrKey, err)
			continue
		}

		// the current log is unique, and we want to add it to our buffer and proceed with forwarding to the user
		buffer.Push(uniqueLogKey)

		logger.Trace(fmt.Sprintf("Forwarding log from Obscuro node: %s", jsonResponse), log.SubIDKey, data.SubID)
		err = userConn.WriteResponse(jsonResponse)
		if err != nil {
			logger.Error("could not write the JSON log to the websocket on subscription %", log.SubIDKey, data.SubID, log.ErrKey, err)
			continue
		}
	}
}

func checkIfUserConnIsClosedAndUnsubscribe(userConn userconn.UserConn, subscription *gethrpc.ClientSubscription, mu *sync.Mutex) {
	for {
		mu.Lock()
		if userConn.IsClosed() {
			subscription.Unsubscribe()
			return
		}
		mu.Unlock()
		time.Sleep(100 * time.Millisecond)
	}
}

func (sm *SubscriptionManager) UpdateSubscriptionMapping(userSubscriptionID string, subscription *gethrpc.ClientSubscription) {
	// Ensure there is no concurrent map writes
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// Check if the userSubscriptionID already exists in the map
	subscriptions, exists := sm.subscriptionMappings[userSubscriptionID]

	// If it doesn't exist, create a new slice for it
	if !exists {
		subscriptions = []*gethrpc.ClientSubscription{}
	}

	// Check if the subscription is already in the slice, if not, add it
	subscriptionExists := false
	for _, sub := range subscriptions {
		if sub == subscription {
			subscriptionExists = true
			break
		}
	}

	if !subscriptionExists {
		sm.subscriptionMappings[userSubscriptionID] = append(subscriptions, subscription)
	}
}

// Formats the log to be sent as an Eth JSON-RPC response.
func prepareLogResponse(idAndLog common.IDAndLog, userSubscriptionID gethrpc.ID) ([]byte, error) {
	paramsMap := make(map[string]interface{})
	paramsMap[wecommon.JSONKeySubscription] = userSubscriptionID
	paramsMap[wecommon.JSONKeyResult] = idAndLog.Log

	respMap := make(map[string]interface{})
	respMap[wecommon.JSONKeyRPCVersion] = jsonrpc.Version
	respMap[wecommon.JSONKeyMethod] = wecommon.MethodEthSubscription
	respMap[wecommon.JSONKeyParams] = paramsMap

	jsonResponse, err := json.Marshal(respMap)
	if err != nil {
		return nil, fmt.Errorf("could not marshal log response to JSON. Cause: %w", err)
	}
	return jsonResponse, nil
}

func (sm *SubscriptionManager) HandleUnsubscribe(userSubscriptionID string, rpcResp *interface{}) {
	subscriptions, exists := sm.subscriptionMappings[userSubscriptionID]

	if !exists {
		*rpcResp = false
		return
	}

	sm.mu.Lock()
	defer sm.mu.Unlock()
	for _, sub := range subscriptions {
		sub.Unsubscribe()
	}
	delete(sm.subscriptionMappings, userSubscriptionID)
	*rpcResp = true
}
