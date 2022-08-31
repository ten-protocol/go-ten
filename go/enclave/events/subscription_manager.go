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
	subscriptions map[uuid.UUID]logSubscription
}

func NewSubscriptionManager() *SubscriptionManager {
	return &SubscriptionManager{}
}

// todo - describe
func (s *SubscriptionManager) AddSubscription(id uuid.UUID, subscription common.EncryptedLogSubscription) error {
	// todo - decrypt and deserialize the subscription
	subcription := logSubscription{}

	// todo
	// check that each account is signed with a valid viewing key which in turn is signed with the account key
	//for _, account := range logSubscription.Accounts {
	//}

	s.subscriptions[id] = subcription
	return nil
}

// todo - describe
func (s *SubscriptionManager) DeleteSubscription(id uuid.UUID) {
	delete(s.subscriptions, id)
}

// todo - describe
func (s *SubscriptionManager) FilterRelevantLogs(logs []*types.Log, stateDB *state.StateDB) map[uuid.UUID][]*types.Log {
	relevantLogs := map[uuid.UUID][]*types.Log{}
	for _, log := range logs {
		for subID, sub := range s.subscriptions {
			isRelevant, _ := sub.isRelevant(log, stateDB)
			// todo return the account somehow, maybe as a tuple with the log
			if isRelevant {
				subResult, found := relevantLogs[subID]
				if !found {
					subResult = make([]*types.Log, 0)
				}
				subResult = append(subResult, log)
				relevantLogs[subID] = subResult
			}
		}
	}
	return relevantLogs
}

// todo - describe
func (s *SubscriptionManager) EncryptLogs(logsBySubID map[uuid.UUID][]*types.Log) (map[uuid.UUID]common.EncryptedLogs, error) {
	result := map[uuid.UUID]common.EncryptedLogs{}
	for subID, logs := range logsBySubID {
		enc := make([]byte, 0)
		for _, log := range logs {
			txReceiptBytes, err := log.MarshalJSON()
			if err != nil {
				return nil, fmt.Errorf("could not marshal log to JSON. Cause: %w", err)
			}

			// todo - add encryption
			// todo - separator
			enc = append(enc, txReceiptBytes...)
		}
		result[subID] = enc
	}
	return result, nil
}

func (s logSubscription) isRelevant(log *types.Log, db *state.StateDB) (bool, *subscriptionAccount) {
	// todo
	// transform the log into a useful data structure by extracting addresses from the log (according to the design)
	// identify what type of log it is ( account specific or lifecycle log)
	// if account-specific go through each SubscriptionAccount and check whether the log is relevant
	// note - the above logic has to be reused to filter out the logs when someone requests a transaction receipt
	// for logs that pass the above the step apply the filters
	// return the first account for which the log matches, so it can be used for encryption
	return true, nil
}

// From the design - call must take a list of signed owning accounts.
// Each account must be signed with the latest viewing key (to prevent someone from asking for random logs, just to leak info).
// The call will fail if there are no viewing keys for all those accounts.
type logSubscription struct {
	accounts []*subscriptionAccount
	// todo Filters - the geth log filters
}

// SubscriptionAccount is an authenticated account used for subscribing to logs.
type subscriptionAccount struct {
	// The account the events relate to.
	account gethcommon.Address
	// A signature over the subscription ID using the private viewing key. Prevents attackers from subscribing to
	// (encrypted) logs for other accounts to see the pattern of logs.
	signature []byte
}
