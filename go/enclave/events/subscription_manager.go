package events

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/ten-protocol/go-ten/go/enclave/vkhandler"
	gethrpc "github.com/ten-protocol/go-ten/lib/gethfork/rpc"

	"github.com/ten-protocol/go-ten/go/common/log"

	"github.com/ten-protocol/go-ten/go/enclave/core"

	"github.com/ten-protocol/go-ten/go/enclave/storage"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common"
)

const (
	// The leading zero bytes in a hash indicating that it is possibly an address, since it only has 20 bytes of data.
	zeroBytesHex = "000000000000000000000000"
)

type logSubscription struct {
	Subscription *common.LogSubscription
	// Handles the viewing key encryption
	ViewingKeyEncryptor *vkhandler.AuthenticatedViewingKey
}

// SubscriptionManager manages the creation/deletion of subscriptions, and the filtering and encryption of logs for
// active subscriptions.
type SubscriptionManager struct {
	storage storage.Storage

	subscriptions     map[gethrpc.ID]*logSubscription
	chainID           int64
	subscriptionMutex *sync.RWMutex // the mutex guards the subscriptions/lastHead pair

	logger gethlog.Logger
}

func NewSubscriptionManager(storage storage.Storage, chainID int64, logger gethlog.Logger) *SubscriptionManager {
	return &SubscriptionManager{
		storage: storage,

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

// FilterLogsForReceipt removes the logs that the sender of a transaction is not allowed to view
func FilterLogsForReceipt(ctx context.Context, receipt *types.Receipt, account *gethcommon.Address, storage storage.Storage) ([]*types.Log, error) {
	filteredLogs := []*types.Log{}
	stateDB, err := storage.CreateStateDB(ctx, receipt.BlockHash)
	if err != nil {
		return nil, fmt.Errorf("could not create state DB to filter logs. Cause: %w", err)
	}

	for _, logItem := range receipt.Logs {
		userAddrs := getUserAddrsFromLogTopics(logItem, stateDB)
		if isRelevant(account, userAddrs) {
			filteredLogs = append(filteredLogs, logItem)
		}
	}

	return filteredLogs, nil
}

// GetSubscribedLogsForBatch - Retrieves and encrypts the logs for the batch in live mode.
// The assumption is that this function is called synchronously after the batch is produced
func (s *SubscriptionManager) GetSubscribedLogsForBatch(ctx context.Context, batch *core.Batch, receipts types.Receipts) (common.EncryptedSubscriptionLogs, error) {
	s.subscriptionMutex.RLock()
	defer s.subscriptionMutex.RUnlock()

	// exit early if there are no subscriptions
	if len(s.subscriptions) == 0 {
		return nil, nil
	}

	relevantLogsPerSubscription := map[gethrpc.ID][]*types.Log{}

	// extract the logs from all receipts
	var allLogs []*types.Log
	for _, receipt := range receipts {
		allLogs = append(allLogs, receipt.Logs...)
	}

	if len(allLogs) == 0 {
		return nil, nil
	}

	// the stateDb is needed to extract the user addresses from the topics
	stateDB, err := s.storage.CreateStateDB(ctx, batch.Hash())
	if err != nil {
		return nil, fmt.Errorf("could not create state DB to filter logs. Cause: %w", err)
	}

	// cache for the user addresses extracted from the individual logs
	// this is an expensive operation so we are doing it lazy, and caching the result
	userAddrsForLog := map[*types.Log][]*gethcommon.Address{}

	for id, sub := range s.subscriptions {
		// first filter the logs
		filteredLogs := filterLogs(allLogs, sub.Subscription.Filter.FromBlock, sub.Subscription.Filter.ToBlock, sub.Subscription.Filter.Addresses, sub.Subscription.Filter.Topics, s.logger)

		// the account requesting the logs is retrieved from the Viewing Key
		requestingAccount := sub.ViewingKeyEncryptor.AccountAddress
		relevantLogsForSub := []*types.Log{}
		for _, logItem := range filteredLogs {
			userAddrs, f := userAddrsForLog[logItem]
			if !f {
				userAddrs = getUserAddrsFromLogTopics(logItem, stateDB)
				userAddrsForLog[logItem] = userAddrs
			}
			relevant := isRelevant(requestingAccount, userAddrs)
			if relevant {
				relevantLogsForSub = append(relevantLogsForSub, logItem)
			}
			s.logger.Debug("Subscription", log.SubIDKey, id, "acc", requestingAccount, "log", logItem, "extr_addr", userAddrs, "relev", relevant)
		}
		if len(relevantLogsForSub) > 0 {
			relevantLogsPerSubscription[id] = relevantLogsForSub
		}
	}

	// Encrypt the results
	return s.encryptLogs(relevantLogsPerSubscription)
}

func isRelevant(sub *gethcommon.Address, userAddrs []*gethcommon.Address) bool {
	// If there are no user addresses, this is a lifecycle event, and is therefore relevant to everyone.
	if len(userAddrs) == 0 {
		return true
	}
	for _, addr := range userAddrs {
		if *addr == *sub {
			return true
		}
	}
	return false
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

// Of the log's topics, returns those that are (potentially) user addresses. A topic is considered a user address if:
//   - It has 12 leading zero bytes (since addresses are 20 bytes long, while hashes are 32)
//   - It has a non-zero nonce (to prevent accidental or malicious creation of the address matching a given topic,
//     forcing its events to become permanently private
//   - It does not have associated code (meaning it's a smart-contract address)
func getUserAddrsFromLogTopics(log *types.Log, db *state.StateDB) []*gethcommon.Address {
	var userAddrs []*gethcommon.Address

	// We skip over the first topic, which is always the hash of the event.
	for _, topic := range log.Topics[1:len(log.Topics)] {
		if topic.Hex()[2:len(zeroBytesHex)+2] != zeroBytesHex {
			continue
		}

		potentialAddr := gethcommon.BytesToAddress(topic.Bytes())

		// A user address must have a non-zero nonce. This prevents accidental or malicious sending of funds to an
		// address matching a topic, forcing its events to become permanently private.
		if db.GetNonce(potentialAddr) != 0 {
			// If the address has code, it's a smart contract address instead.
			if db.GetCode(potentialAddr) == nil {
				userAddrs = append(userAddrs, &potentialAddr)
			}
		}
	}

	return userAddrs
}

// Lifted from eth/filters/filter.go in the go-ethereum repository.
// filterLogs creates a slice of logs matching the given criteria.
func filterLogs(logs []*types.Log, fromBlock, toBlock *gethrpc.BlockNumber, addresses []gethcommon.Address, topics [][]gethcommon.Hash, logger gethlog.Logger) []*types.Log { //nolint:gocognit
	var ret []*types.Log
Logs:
	for _, logItem := range logs {
		if fromBlock != nil && fromBlock.Int64() >= 0 && fromBlock.Int64() > int64(logItem.BlockNumber) {
			logger.Debug("Skipping log ", "log", logItem, "reason", "In the past. The starting block num for filter is bigger than log")
			continue
		}
		if toBlock != nil && toBlock.Int64() > 0 && toBlock.Int64() < int64(logItem.BlockNumber) {
			logger.Debug("Skipping log ", "log", logItem, "reason", "In the future. The ending block num for filter is smaller than log")
			continue
		}

		if len(addresses) > 0 && !includes(addresses, logItem.Address) {
			logger.Debug("Skipping log ", "log", logItem, "reason", "The contract address of the log is not an address of interest")
			continue
		}
		// If the to filtered topics is greater than the amount of topics in logs, skip.
		if len(topics) > len(logItem.Topics) {
			logger.Debug("Skipping log ", "log", logItem, "reason", "Insufficient topics. The log has less topics than the required one to satisfy the query")
			continue
		}
		for i, sub := range topics {
			match := len(sub) == 0 // empty rule set == wildcard
			for _, topic := range sub {
				if logItem.Topics[i] == topic {
					match = true
					break
				}
			}
			if !match {
				logger.Debug("Skipping log ", "log", logItem, "reason", "Topics do not match.")
				continue Logs
			}
		}
		ret = append(ret, logItem)
	}
	return ret
}

// Lifted from eth/filters/filter.go in the go-ethereum repository.
func includes(addresses []gethcommon.Address, a gethcommon.Address) bool {
	for _, addr := range addresses {
		if addr == a {
			return true
		}
	}

	return false
}
