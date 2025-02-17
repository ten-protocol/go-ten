package events

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/ten-protocol/go-ten/go/enclave/components"

	"github.com/ten-protocol/go-ten/go/enclave/vkhandler"
	gethrpc "github.com/ten-protocol/go-ten/lib/gethfork/rpc"

	"github.com/ten-protocol/go-ten/go/enclave/core"

	"github.com/ten-protocol/go-ten/go/enclave/storage"

	"github.com/ethereum/go-ethereum/core/types"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common"
)

type logSubscription struct {
	Subscription *common.LogSubscription
	// Handles the viewing key encryption
	ViewingKeyEncryptor *vkhandler.AuthenticatedViewingKey
}

// SubscriptionManager manages the creation/deletion of subscriptions, and the filtering and encryption of logs for
// active subscriptions.
type SubscriptionManager struct {
	storage  storage.Storage
	registry components.BatchRegistry

	subscriptions     map[gethrpc.ID]*logSubscription
	chainID           int64
	subscriptionMutex *sync.RWMutex // the mutex guards the subscriptions/lastHead pair

	logger gethlog.Logger
}

func NewSubscriptionManager(storage storage.Storage, registry components.BatchRegistry, chainID int64, logger gethlog.Logger) *SubscriptionManager {
	return &SubscriptionManager{
		storage:  storage,
		registry: registry,

		subscriptions:     map[gethrpc.ID]*logSubscription{},
		chainID:           chainID,
		subscriptionMutex: &sync.RWMutex{},
		logger:            logger,
	}
}

// AddSubscription adds a log subscription to the enclave under the given ID, provided the request is authenticated
// correctly. If there is an existing subscription with the given ID, it is overwritten.
func (s *SubscriptionManager) AddSubscription(id gethrpc.ID, encodedSubscription []byte) error {
	subscription := &common.LogSubscription{}
	if err := json.Unmarshal(encodedSubscription, subscription); err != nil {
		return fmt.Errorf("could not decode log subscription. Cause: %w", err)
	}

	// verify the viewing key
	authenticateViewingKey, err := vkhandler.VerifyViewingKey(subscription.ViewingKey, s.chainID)
	if err != nil {
		return fmt.Errorf("unable to authenticate the viewing key for subscription  - %w", err)
	}

	s.subscriptionMutex.Lock()
	defer s.subscriptionMutex.Unlock()
	s.subscriptions[id] = &logSubscription{
		Subscription:        subscription,
		ViewingKeyEncryptor: authenticateViewingKey,
	}

	return nil
}

// RemoveSubscription removes the log subscription with the given ID from the enclave. If there is no subscription with
// the given ID, nothing is deleted.
func (s *SubscriptionManager) RemoveSubscription(id gethrpc.ID) {
	s.subscriptionMutex.Lock()
	defer s.subscriptionMutex.Unlock()
	delete(s.subscriptions, id)
}

// GetSubscribedLogsForBatch - Retrieves and encrypts the logs for the batch in live mode.
// The assumption is that this function is called synchronously after the batch is produced
func (s *SubscriptionManager) GetSubscribedLogsForBatch(ctx context.Context, batch *core.Batch, receipts types.Receipts) (common.EncryptedSubscriptionLogs, error) {
	s.subscriptionMutex.RLock()
	subs := make(map[gethrpc.ID]*logSubscription)
	for key, value := range s.subscriptions {
		subs[key] = value
	}
	s.subscriptionMutex.RUnlock()

	// exit early if there are no subscriptions
	if len(s.subscriptions) == 0 {
		return nil, nil
	}

	h := batch.Hash()
	relevantLogsPerSubscription := map[gethrpc.ID][]*types.Log{}

	if len(receipts) == 0 {
		return nil, nil
	}

	for id, sub := range subs {
		relevantLogsForSub, err := s.storage.FilterLogs(ctx, sub.ViewingKeyEncryptor.AccountAddress, nil, nil, &h, sub.Subscription.Filter.Addresses, sub.Subscription.Filter.Topics)
		if err != nil {
			return nil, err
		}
		if len(relevantLogsForSub) > 0 {
			relevantLogsPerSubscription[id] = relevantLogsForSub
		}
	}

	// Encrypt the results
	return s.encryptLogs(relevantLogsPerSubscription)
}

// Encrypts each log with the appropriate viewing key.
func (s *SubscriptionManager) encryptLogs(logsByID map[gethrpc.ID][]*types.Log) (map[gethrpc.ID][]byte, error) {
	encryptedLogsByID := map[gethrpc.ID][]byte{}

	for subID, logs := range logsByID {
		subscription, found := s.subscriptions[subID]
		if !found {
			continue // The subscription has been removed, so there's no need to return anything.
		}

		jsonLogs, err := json.Marshal(logs)
		if err != nil {
			return nil, fmt.Errorf("could not marshal logs to JSON. Cause: %w", err)
		}

		encryptedLogs, err := subscription.ViewingKeyEncryptor.Encrypt(jsonLogs)
		if err != nil {
			return nil, fmt.Errorf("unable to encrypt logs - %w", err)
		}

		encryptedLogsByID[subID] = encryptedLogs
	}

	return encryptedLogsByID, nil
}
