package events

import (
	"fmt"

	"github.com/ethereum/go-ethereum/rlp"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/google/uuid"
	"github.com/obscuronet/go-obscuro/go/common"
)

type SubscriptionManager struct {
	subscriptions map[uuid.UUID]*common.LogSubscription
}

func NewSubscriptionManager() *SubscriptionManager {
	return &SubscriptionManager{
		subscriptions: map[uuid.UUID]*common.LogSubscription{},
	}
}

// AddSubscription adds a log subscription to the enclave under the given ID, provided the request is authenticated
// correctly. If there is an existing subscription with the given ID, it is overwritten.
// TODO - #453 - Decrypt subscriptions (currently expects unencrypted serialised bytes).
// TODO - #453 - Check each account in the subscription request is authenticated with a signature.
func (s *SubscriptionManager) AddSubscription(id uuid.UUID, encryptedSubscription common.EncryptedLogSubscription) error {
	var subscription common.LogSubscription
	if err := rlp.DecodeBytes(encryptedSubscription, &subscription); err != nil {
		return fmt.Errorf("could not decode encoded log subscription. Cause: %w", err)
	}

	s.subscriptions[id] = &subscription
	println("jjj added subscription")
	return nil
}

// RemoveSubscription removes the log subscription with the given ID from the enclave. If there is no subscription with
// the given ID, nothing is deleted.
// TODO - #453 - Consider whether the deletion needs to be authenticated as well, to prevent attackers deleting subscriptions.
func (s *SubscriptionManager) RemoveSubscription(id uuid.UUID) {
	delete(s.subscriptions, id)
	println("jjj deleted subscription")
}

// FilterRelevantLogs returns only those logs for which one or more subscriptions exist.
// TODO - #453 - Return which account each log is relevant to.
func (s *SubscriptionManager) FilterRelevantLogs(logs []*types.Log) map[uuid.UUID][]*types.Log {
	relevantLogs := map[uuid.UUID][]*types.Log{}

	for _, log := range logs {
		for subID, sub := range s.subscriptions {
			logIsRelevant, _ := isRelevant(sub, log)
			if !logIsRelevant {
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
// TODO - #453 - Encrypt logs, rather than just serialising them.
func (s *SubscriptionManager) EncryptLogs(logsBySubID map[uuid.UUID][]*types.Log) (map[uuid.UUID]common.EncryptedLogs, error) {
	result := map[uuid.UUID]common.EncryptedLogs{}
	for subID, logs := range logsBySubID {
		enc := make([]byte, 0)
		for _, log := range logs {
			logBytes, err := log.MarshalJSON()
			if err != nil {
				return nil, fmt.Errorf("could not marshal log to JSON. Cause: %w", err)
			}

			// TODO - #453 - Add separator between each serialised log.
			enc = append(enc, logBytes...)
		}
		result[subID] = enc
	}
	return result, nil
}

// Indicates whether the log is relevant for any subscription, and returns the matching subscription account if so.
// A lifecycle log is considered relevant to everyone.
// TODO - #453 - Filter logs, instead of considering all logs relevant to everyone.
func isRelevant(sub *common.LogSubscription, log *types.Log) (bool, *common.SubscriptionAccount) {
	// Extract addresses from the logs
	// Work out if this is an account or lifecycle log
	// If the former, establish whether it is relevant to any subscription
	// Return the first account for which the log matches, so it can be used for encryption
	return true, nil
}
