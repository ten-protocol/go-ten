package events

import (
	"fmt"

	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/google/uuid"
	"github.com/obscuronet/go-obscuro/go/common"
)

type SubscriptionManager struct {
	subscriptions map[uuid.UUID]common.EventSubscription
}

func NewSubscriptionManager() *SubscriptionManager {
	return &SubscriptionManager{}
}

// todo - describe
func (s *SubscriptionManager) AddSubscription(id uuid.UUID, subscription common.EncryptedEventSubscription) error {
	// todo - decrypt and deserialize the subscription
	eventSubscription := common.EventSubscription{}

	// todo
	// check that each account is signed with a valid viewing key which in turn is signed with the account key
	//for _, account := range eventSubscription.Accounts {
	//}

	s.subscriptions[id] = eventSubscription
	return nil
}

// todo - describe
func (s *SubscriptionManager) DeleteSubscription(id uuid.UUID) {
	delete(s.subscriptions, id)
}

// todo - describe
func (s *SubscriptionManager) ExtractEvents(events []*types.Log, stateDB *state.StateDB) map[uuid.UUID][]*types.Log {
	result := map[uuid.UUID][]*types.Log{}
	for _, event := range events {
		for subID, sub := range s.subscriptions {
			matches, _ := sub.Matches(event, stateDB)
			// todo return the account somehow, maybe as a tuple with the log
			if matches {
				subResult, found := result[subID]
				if !found {
					subResult = make([]*types.Log, 0)
				}
				subResult = append(subResult, event)
				result[subID] = subResult
			}
		}
	}
	return result
}

// todo - describe
func (s *SubscriptionManager) EncryptEvents(eventsPerSubscription map[uuid.UUID][]*types.Log) (map[uuid.UUID]common.EncryptedEvents, error) {
	result := map[uuid.UUID]common.EncryptedEvents{}
	for u, events := range eventsPerSubscription {
		enc := make([]byte, 0)
		for _, event := range events {
			txReceiptBytes, err := event.MarshalJSON()
			if err != nil {
				return nil, fmt.Errorf("could not marshal event log to JSON. Cause: %w", err)
			}

			// todo - add encryption
			// todo - separator
			enc = append(enc, txReceiptBytes...)
		}
		result[u] = enc
	}
	return result, nil
}
