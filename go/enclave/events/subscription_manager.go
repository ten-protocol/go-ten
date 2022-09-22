package events

import (
	"encoding/json"
	"fmt"

	"github.com/obscuronet/go-obscuro/go/enclave/db"

	"github.com/ethereum/go-ethereum/core/state"

	gethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/obscuronet/go-obscuro/go/enclave/rpc"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/google/uuid"
	"github.com/obscuronet/go-obscuro/go/common"
)

const (
	// The leading zero bytes in a hash indicating that it is possibly an address, since it only has 20 bytes of data.
	zeroBytesHex = "000000000000000000000000"
)

// SubscriptionManager manages the creation/deletion of subscriptions, and the filtering and encryption of logs for
// active subscriptions.
type SubscriptionManager struct {
	rpcEncryptionManager *rpc.EncryptionManager
	storage              db.Storage
	subscriptions        map[uuid.UUID]*common.LogSubscription
}

func NewSubscriptionManager(rpcEncryptionManager *rpc.EncryptionManager, storage db.Storage) *SubscriptionManager {
	return &SubscriptionManager{
		rpcEncryptionManager: rpcEncryptionManager,
		storage:              storage,
		subscriptions:        map[uuid.UUID]*common.LogSubscription{},
	}
}

// AddSubscription adds a log subscription to the enclave under the given ID, provided the request is authenticated
// correctly. If there is an existing subscription with the given ID, it is overwritten.
func (s *SubscriptionManager) AddSubscription(id uuid.UUID, encryptedSubscription common.EncryptedParamsLogSubscription) error {
	jsonSubscription, err := s.rpcEncryptionManager.DecryptBytes(encryptedSubscription)
	if err != nil {
		return fmt.Errorf("could not decrypt params in eth_subscribe logs request. Cause: %w", err)
	}

	var subscriptions []common.LogSubscription
	if err := json.Unmarshal(jsonSubscription, &subscriptions); err != nil {
		return fmt.Errorf("could not unmarshal log subscription from JSON. Cause: %w", err)
	}

	if len(subscriptions) != 1 {
		return fmt.Errorf("expected a single log subscription, received %d", len(subscriptions))
	}
	subscription := subscriptions[0]

	err = s.rpcEncryptionManager.AuthenticateSubscriptionRequest(subscription)
	if err != nil {
		return err
	}

	s.subscriptions[id] = &subscription
	return nil
}

// RemoveSubscription removes the log subscription with the given ID from the enclave. If there is no subscription with
// the given ID, nothing is deleted.
func (s *SubscriptionManager) RemoveSubscription(id uuid.UUID) {
	delete(s.subscriptions, id)
}

// FilteredLogs filters out irrelevant logs.
func (s *SubscriptionManager) FilteredLogs(logs []*types.Log, rollupHash common.L2RootHash, account *gethcommon.Address) []*types.Log {
	allLogs := []*types.Log{}
	stateDB := s.storage.CreateStateDB(rollupHash)

	for _, log := range logs {
		userAddrs := getUserAddrs(log, stateDB)
		if isRelevant(userAddrs, account) {
			allLogs = append(allLogs, log)
		}
	}

	return allLogs
}

// FilteredSubscribedLogs filters out irrelevant logs and those that are not subscribed to, and organises them by their subscribing ID.
func (s *SubscriptionManager) FilteredSubscribedLogs(logs []*types.Log, rollupHash common.L2RootHash) map[uuid.UUID][]*types.Log {
	relevantLogs := map[uuid.UUID][]*types.Log{}

	// If there are no subscriptions, we do not need to do any processing.
	if len(s.subscriptions) == 0 {
		return relevantLogs
	}

	for subscriptionID := range s.subscriptions {
		relevantLogs[subscriptionID] = []*types.Log{}
	}

	stateDB := s.storage.CreateStateDB(rollupHash)

	for _, log := range logs {
		userAddrs := getUserAddrs(log, stateDB)

		// We check whether the log is relevant to each subscription.
		for subscriptionID, subscription := range s.subscriptions {
			if isRelevant(userAddrs, subscription.Account) {
				relevantLogs[subscriptionID] = append(relevantLogs[subscriptionID], log)
			}
		}
	}

	return relevantLogs
}

// EncryptLogs encrypts each log with the appropriate viewing key.
func (s *SubscriptionManager) EncryptLogs(logsBySubID map[uuid.UUID][]*types.Log) (map[uuid.UUID]common.EncryptedLogs, error) {
	result := map[uuid.UUID]common.EncryptedLogs{}
	for subID, logs := range logsBySubID {
		subscription, found := s.subscriptions[subID]
		if !found {
			return nil, fmt.Errorf("could not find subscription with ID %s", subID.String())
		}

		jsonLogs, err := json.Marshal(logs)
		if err != nil {
			return nil, fmt.Errorf("could not marshal logs to JSON. Cause: %w", err)
		}

		encryptedLogs, err := s.rpcEncryptionManager.EncryptWithViewingKey(*subscription.Account, jsonLogs)
		if err != nil {
			return nil, err
		}

		result[subID] = encryptedLogs
	}
	return result, nil
}

// Extracts the (potential) user addresses from the topics. If there is no code associated with an address, it's a user
// address.
func getUserAddrs(log *types.Log, db *state.StateDB) []string {
	var userAddrs []string //nolint:prealloc

	for _, topic := range log.Topics {
		// Since addresses are 20 bytes long, while hashes are 32, only topics with 12 leading zero bytes can
		// (potentially) be user addresses.
		topicHex := topic.Hex()
		if topicHex[2:len(zeroBytesHex)+2] != zeroBytesHex {
			continue
		}

		// If there is code associated with an address, it is not a user address.
		addr := gethcommon.HexToAddress(topicHex)
		if db.GetCode(addr) != nil {
			continue
		}

		userAddrs = append(userAddrs, addr.Hex())
	}

	return userAddrs
}

// Indicates whether the log is relevant for the subscription.
func isRelevant(userAddrs []string, account *gethcommon.Address) bool {
	// If there are no potential user addresses, this is a lifecycle event, and is therefore relevant to everyone.
	if len(userAddrs) == 0 {
		return true
	}

	for _, addr := range userAddrs {
		if addr == account.Hex() {
			return true
		}
	}

	return false
}
