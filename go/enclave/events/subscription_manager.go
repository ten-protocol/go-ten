package events

import (
	"encoding/json"
	"fmt"
	"math/big"
	"sync"

	gethlog "github.com/ethereum/go-ethereum/log"

	"github.com/ethereum/go-ethereum/rlp"

	"github.com/ethereum/go-ethereum/eth/filters"

	"github.com/obscuronet/go-obscuro/go/enclave/db"

	"github.com/ethereum/go-ethereum/core/state"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethrpc "github.com/ethereum/go-ethereum/rpc"

	"github.com/obscuronet/go-obscuro/go/enclave/rpc"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/go-obscuro/go/common"
)

const (
	// The leading zero bytes in a hash indicating that it is possibly an address, since it only has 20 bytes of data.
	zeroBytesHex = "000000000000000000000000"
)

// TODO - Ensure chain reorgs are handled gracefully.

// SubscriptionManager manages the creation/deletion of subscriptions, and the filtering and encryption of logs for
// active subscriptions.
type SubscriptionManager struct {
	rpcEncryptionManager *rpc.EncryptionManager
	storage              db.Storage

	subscriptions     map[gethrpc.ID]*common.LogSubscription
	lastHead          map[gethrpc.ID]*big.Int // This is the batch height up to which events were returned to the user
	subscriptionMutex *sync.RWMutex
	logger            gethlog.Logger
}

func NewSubscriptionManager(rpcEncryptionManager *rpc.EncryptionManager, storage db.Storage, logger gethlog.Logger) *SubscriptionManager {
	return &SubscriptionManager{
		rpcEncryptionManager: rpcEncryptionManager,
		storage:              storage,

		subscriptions:     map[gethrpc.ID]*common.LogSubscription{},
		lastHead:          map[gethrpc.ID]*big.Int{},
		subscriptionMutex: &sync.RWMutex{},
		logger:            logger,
	}
}

func (s *SubscriptionManager) ForEachSubscription(f func(gethrpc.ID, *common.LogSubscription, *big.Int) error) error {
	s.subscriptionMutex.Lock()
	defer s.subscriptionMutex.Unlock()

	for id, subscription := range s.subscriptions {
		err := f(id, subscription, s.lastHead[id])
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *SubscriptionManager) SetLastHead(id gethrpc.ID, nr *big.Int) {
	s.lastHead[id] = nr
}

// AddSubscription adds a log subscription to the enclave under the given ID, provided the request is authenticated
// correctly. If there is an existing subscription with the given ID, it is overwritten.
func (s *SubscriptionManager) AddSubscription(id gethrpc.ID, encryptedSubscription common.EncryptedParamsLogSubscription) error {
	encodedSubscription, err := s.rpcEncryptionManager.DecryptBytes(encryptedSubscription)
	if err != nil {
		return fmt.Errorf("could not decrypt params in eth_subscribe logs request. Cause: %w", err)
	}

	var subscription common.LogSubscription
	if err = rlp.DecodeBytes(encodedSubscription, &subscription); err != nil {
		return fmt.Errorf("could not decocde log subscription from RLP. Cause: %w", err)
	}

	err = s.rpcEncryptionManager.AuthenticateSubscriptionRequest(subscription)
	if err != nil {
		return err
	}

	// For subscriptions, only the Topics and Addresses fields of the filter are applied.
	subscription.Filter.BlockHash = nil
	subscription.Filter.ToBlock = nil

	// We set the FromBlock to the current rollup height, so that historical logs aren't returned.
	// Todo - this is probably not the right behaviour
	height := big.NewInt(0)
	rollup, err := s.storage.FetchHeadBatch()
	if err == nil {
		height = rollup.Number()
	}
	subscription.Filter.FromBlock = big.NewInt(0).Add(height, big.NewInt(1))

	// Always start from the FromBlock
	s.lastHead[id] = subscription.Filter.FromBlock
	s.subscriptionMutex.Lock()
	defer s.subscriptionMutex.Unlock()
	s.subscriptions[id] = &subscription
	return nil
}

// RemoveSubscription removes the log subscription with the given ID from the enclave. If there is no subscription with
// the given ID, nothing is deleted.
func (s *SubscriptionManager) RemoveSubscription(id gethrpc.ID) {
	s.subscriptionMutex.Lock()
	defer s.subscriptionMutex.Unlock()
	delete(s.subscriptions, id)
}

// FilterLogs takes a list of logs and the hash of the rollup to use to create the state DB. It returns the logs
// filtered based on the provided account and filter.
func (s *SubscriptionManager) FilterLogs(logs []*types.Log, rollupHash common.L2RootHash, account *gethcommon.Address, filter *filters.FilterCriteria) ([]*types.Log, error) {
	filteredLogs := []*types.Log{}
	stateDB, err := s.storage.CreateStateDB(rollupHash)
	if err != nil {
		return nil, fmt.Errorf("could not create state DB to filter logs. Cause: %w", err)
	}

	for _, logItem := range logs {
		userAddrs := getUserAddrsFromLogTopics(logItem, stateDB)
		if isRelevant(logItem, userAddrs, account, filter, s.logger) {
			filteredLogs = append(filteredLogs, logItem)
		}
	}

	return filteredLogs, nil
}

// EncryptLogs Encrypts each log with the appropriate viewing key.
func (s *SubscriptionManager) EncryptLogs(logsByID map[gethrpc.ID][]*types.Log) (map[gethrpc.ID][]byte, error) {
	encryptedLogsByID := map[gethrpc.ID][]byte{}

	for subID, logs := range logsByID {
		subscription, found := s.getSubscriptionThreadsafe(subID)
		if !found {
			continue // The subscription has been removed, so there's no need to return anything.
		}

		jsonLogs, err := json.Marshal(logs)
		if err != nil {
			return nil, fmt.Errorf("could not marshal logs to JSON. Cause: %w", err)
		}

		encryptedLogs, err := s.rpcEncryptionManager.EncryptWithViewingKey(*subscription.Account, jsonLogs)
		if err != nil {
			return nil, err
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
func getUserAddrsFromLogTopics(log *types.Log, db *state.StateDB) []string {
	var userAddrs []string

	// We skip over the first topic, which is always the hash of the event.
	for _, topic := range log.Topics[1:len(log.Topics)] {
		potentialAddr := gethcommon.HexToAddress(topic.Hex())

		if topic.Hex()[2:len(zeroBytesHex)+2] != zeroBytesHex {
			continue
		}

		// A user address must have a non-zero nonce. This prevents accidental or malicious sending of funds to an
		// address matching a topic, forcing its events to become permanently private.
		if db.GetNonce(potentialAddr) != 0 {
			// If the address has code, it's a smart contract address instead.
			if db.GetCode(potentialAddr) == nil {
				userAddrs = append(userAddrs, potentialAddr.Hex())
			}
		}
	}

	return userAddrs
}

// Locks the subscription map and retrieves the subscription with subID, or (nil, false) if so such subscription is found.
func (s *SubscriptionManager) getSubscriptionThreadsafe(subID gethrpc.ID) (*common.LogSubscription, bool) {
	s.subscriptionMutex.RLock()
	defer s.subscriptionMutex.RUnlock()

	subscription, found := s.subscriptions[subID]
	return subscription, found
}

// Indicates whether BOTH of the following apply:
//   - One of the log's user addresses matches the subscription's account
//   - The log matches the filter
func isRelevant(logItem *types.Log, userAddrs []string, account *gethcommon.Address, filter *filters.FilterCriteria, logger gethlog.Logger) bool {
	logger.Info(fmt.Sprintf("Checking if log = %v is relevant for account - %s. Addresses extracted from topics =  %v", logItem, account.String(), userAddrs))

	filteredLogs := filterLogs([]*types.Log{logItem}, filter.FromBlock, filter.ToBlock, filter.Addresses, filter.Topics, logger)
	if len(filteredLogs) == 0 {
		return false
	}

	// If there are no user addresses, this is a lifecycle event, and is therefore relevant to everyone.
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

// Lifted from eth/filters/filter.go in the go-ethereum repository.
// filterLogs creates a slice of logs matching the given criteria.
func filterLogs(logs []*types.Log, fromBlock, toBlock *big.Int, addresses []gethcommon.Address, topics [][]gethcommon.Hash, logger gethlog.Logger) []*types.Log { //nolint:gocognit
	var ret []*types.Log
Logs:
	for _, logItem := range logs {
		if fromBlock != nil && fromBlock.Int64() >= 0 && fromBlock.Uint64() > logItem.BlockNumber {
			logger.Info(fmt.Sprintf("Skipping log = %v", logItem), "reason", "In the past. The starting block num for filter is bigger than log")
			continue
		}
		if toBlock != nil && toBlock.Int64() >= 0 && toBlock.Uint64() < logItem.BlockNumber {
			logger.Info(fmt.Sprintf("Skipping log = %v", logItem), "reason", "In the future. The ending block num for filter is smaller than log")
			continue
		}

		if len(addresses) > 0 && !includes(addresses, logItem.Address) {
			logger.Info(fmt.Sprintf("Skipping log = %v", logItem), "reason", "The contract address of the log is not an address of interest")
			continue
		}
		// If the to filtered topics is greater than the amount of topics in logs, skip.
		if len(topics) > len(logItem.Topics) {
			logger.Info(fmt.Sprintf("Skipping log = %v", logItem), "reason", "Insufficient topics. The log has less topics than the required one to satisfy the query")
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
				logger.Info(fmt.Sprintf("Skipping log = %v", logItem), "reason", "Topics do not match.")
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
