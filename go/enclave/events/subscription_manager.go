package events

import (
	"fmt"

	"github.com/ethereum/go-ethereum/rlp"

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

// AddSubscription adds a log subscription to the enclave, provided the request is authenticated correctly.
func (s *SubscriptionManager) AddSubscription(id uuid.UUID, encryptedSubscription common.EncryptedLogSubscription) error {
	// todo - decrypt the subscription

	var sub logSubscription
	if err := rlp.DecodeBytes(encryptedSubscription, &sub); err != nil {
		return fmt.Errorf("could not decode encoded log subscription. Cause: %w", err)
	}

	for _, account := range sub.accounts {
		println(account) // todo - check that each account is signed correctly
	}

	s.subscriptions[id] = sub
	return nil
}

// DeleteSubscription removes a log subscription from the enclave.
// TODO - Consider whether the deletion needs to be authenticated as well, to prevent attackers deleting subscriptions.
func (s *SubscriptionManager) DeleteSubscription(id uuid.UUID) {
	delete(s.subscriptions, id)
}

// FilterRelevantLogs returns only those logs for which one or more subscriptions exist.
func (s *SubscriptionManager) FilterRelevantLogs(logs []*types.Log, stateDB *state.StateDB) map[uuid.UUID][]*types.Log {
	relevantLogs := map[uuid.UUID][]*types.Log{}

	for _, log := range logs {
		for subID, sub := range s.subscriptions {
			isRelevant, _ := sub.isRelevant(log, stateDB)
			// todo return the account somehow, maybe as a tuple with the log
			if !isRelevant {
				continue
			}

			logsForSubID, found := relevantLogs[subID]
			if !found {
				relevantLogs[subID] = []*types.Log{log}
			} else {
				relevantLogs[subID] = append(logsForSubID, log)
			}
		}
	}

	return relevantLogs
}

// EncryptLogs encrypts the logs with the corresponding viewing keys.
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

// Indicates whether the log is relevant for any subscription, and returns the matching subscription account if so.
// A lifecycle log is considered relevant to everyone.
func (s logSubscription) isRelevant(log *types.Log, db *state.StateDB) (bool, *subscriptionAccount) {
	// todo - extract addresses from the logs
	// todo - work out if this is an account or lifecycle log
	// todo - if the former, establish whether it is relevant to any subscription
	// todo - return the first account for which the log matches, so it can be used for encryption
	return true, nil
}

// Information about a log subscription, including a list of authenticated accounts for the subscription.
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
