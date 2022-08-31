package events

import (
	"fmt"

	gethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/google/uuid"
	"github.com/obscuronet/go-obscuro/go/common"
)

type SubscriptionManager struct {
	subscriptions map[uuid.UUID]EventSubscription
}

func NewSubscriptionManager() *SubscriptionManager {
	return &SubscriptionManager{}
}

// todo - describe
func (s *SubscriptionManager) AddSubscription(id uuid.UUID, subscription common.EncryptedEventSubscription) error {
	// todo - decrypt and deserialize the subscription
	eventSubscription := EventSubscription{}

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
			matches, _ := sub.matches(event, stateDB)
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

func (s EventSubscription) matches(log *types.Log, db *state.StateDB) (bool, *SubscriptionAccount) {
	// todo
	// transform the log into a useful data structure by extracting addresses from the log (according to the design)
	// identify what type of log it is ( account specific or lifecycle log)
	// if account-specific go through each SubscriptionAccount and check whether the log is relevant
	// note - the above logic has to be reused to filter out the logs when someone requests a transaction receipt
	// for logs that pass the above the step apply the filters
	// return the first account for which the log matches, so it can be used for encryption
	return true, nil
}

// EventSubscription
// From the design - call must take a list of signed owning accounts.
// Each account must be signed with the latest viewing key (to prevent someone from asking random events, just to leak info).
// The call will fail if there are no viewing keys for all those accounts.
type EventSubscription struct {
	Accounts []*SubscriptionAccount
	// todo Filters - the geth log filters
}

// SubscriptionAccount is an authenticated account used for subscribing to events.
type SubscriptionAccount struct {
	// The account the events relate to.
	Account gethcommon.Address
	// A signature over the subscription ID using the private viewing key. Prevents attackers from subscribing to
	// events for other accounts to see the pattern of events.
	Signature []byte
}
