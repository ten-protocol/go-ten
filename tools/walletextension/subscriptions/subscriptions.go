package subscriptions

import (
	"context"
	"encoding/json"
	"fmt"
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
	subscriptionMappings map[string][]string
	logger               gethlog.Logger
}

func New(logger gethlog.Logger) *SubscriptionManager {
	return &SubscriptionManager{
		subscriptionMappings: make(map[string][]string),
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
	// TODO: currently we use only first client (enabling subscriptions for all of them will be part of future PR)
	for _, client := range clients {
		subscription, err := client.Subscribe(context.Background(), resp, rpc.SubscribeNamespace, funnelMultipleAccountsChan, req.Params...)
		if err != nil {
			return fmt.Errorf("could not call %s with params %v. Cause: %w", req.Method, req.Params, err)
		}

		// We periodically check if the websocket is closed, and terminate the subscription.
		// TODO: test this feature in integration test
		go checkIfUserConnIsClosedAndUnsubscribe(userConn, subscription)

		// Make a connection between subscriptionID returned from node for current request and subscriptionID returned to user
		if currentNodeSubscriptionID, ok := (*resp).(string); ok {
			sm.UpdateSubscriptionMapping(string(userSubscriptionID), currentNodeSubscriptionID)
		} else {
			sm.logger.Error("Unable to read subscriptionID")
		}
	}

	// We return subscriptionID with resp interface. We want to use userSubscriptionID to allow unsubscribing
	*resp = userSubscriptionID
	return nil
}

func readFromChannelAndWriteToUserConn(channel chan common.IDAndLog, userConn userconn.UserConn, userSubscriptionID gethrpc.ID, logger gethlog.Logger) {
	for data := range channel {
		jsonResponse, err := PrepareLogResponse(data, userSubscriptionID)
		if err != nil {
			logger.Error("could not marshal log response to JSON on subscription.", log.SubIDKey, data.SubID, log.ErrKey, err)
			continue
		}

		logger.Trace(fmt.Sprintf("Forwarding log from Obscuro node: %s", jsonResponse), log.SubIDKey, data.SubID)
		err = userConn.WriteResponse(jsonResponse)
		if err != nil {
			logger.Error("could not write the JSON log to the websocket on subscription %", log.SubIDKey, data.SubID, log.ErrKey, err)
			continue
		}
	}
}

func checkIfUserConnIsClosedAndUnsubscribe(userConn userconn.UserConn, subscription *gethrpc.ClientSubscription) {
	for {
		if userConn.IsClosed() {
			subscription.Unsubscribe()
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func (sm *SubscriptionManager) UpdateSubscriptionMapping(userSubscriptionID string, obscuroNodeSubscriptionID string) {
	existingUserIDs, exists := sm.subscriptionMappings[userSubscriptionID]

	if !exists {
		sm.subscriptionMappings[userSubscriptionID] = []string{obscuroNodeSubscriptionID}
		return
	}

	// Check if obscuroNodeSubscriptionID already exists to avoid duplication
	alreadyExists := false
	for _, existingID := range existingUserIDs {
		if obscuroNodeSubscriptionID == existingID {
			alreadyExists = true
			break
		}
	}

	if !alreadyExists {
		sm.subscriptionMappings[userSubscriptionID] = append(existingUserIDs, obscuroNodeSubscriptionID)
	}
}

// Formats the log to be sent as an Eth JSON-RPC response.
func PrepareLogResponse(idAndLog common.IDAndLog, userSubscriptionID gethrpc.ID) ([]byte, error) {
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
