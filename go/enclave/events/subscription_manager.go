package events

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"sync"

	"github.com/obscuronet/go-obscuro/go/common/errutil"

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
	subscriptionMutex *sync.RWMutex
	logger            gethlog.Logger
}

func NewSubscriptionManager(rpcEncryptionManager *rpc.EncryptionManager, storage db.Storage, logger gethlog.Logger) *SubscriptionManager {
	return &SubscriptionManager{
		rpcEncryptionManager: rpcEncryptionManager,
		storage:              storage,

		subscriptions:     map[gethrpc.ID]*common.LogSubscription{},
		subscriptionMutex: &sync.RWMutex{},
		logger:            logger,
	}
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
	rollup, err := s.storage.FetchHeadRollup()
	if err != nil {
		return fmt.Errorf("unable to fetch head rollup. Cause: %w", err)
	}
	if rollup == nil {
		return fmt.Errorf("no head rollup is stored")
	}
	subscription.Filter.FromBlock = big.NewInt(0).Add(rollup.Number(), big.NewInt(1))

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

// GetFilteredLogs returns the logs across the entire canonical chain that match the provided account and filter.
func (s *SubscriptionManager) GetFilteredLogs(account *gethcommon.Address, filter *filters.FilterCriteria) ([]*types.Log, error) {
	headBlock, err := s.storage.FetchHeadBlock()
	if err != nil {
		if errors.Is(err, errutil.ErrNotFound) {
			// There is no head block, and thus no logs to retrieve.
			return nil, nil
		}
		return nil, err
	}

	// We collect all the block hashes in the canonical chain.
	// TODO: Only collect blocks within the filter's range.
	blockHashes := []gethcommon.Hash{}
	currentBlock := headBlock
	for {
		blockHashes = append(blockHashes, currentBlock.Hash())

		if currentBlock.NumberU64() <= common.L1GenesisHeight {
			break // We have reached the end of the chain.
		}

		parentHash := currentBlock.ParentHash()
		var found bool
		currentBlock, found = s.storage.FetchBlock(parentHash)
		if !found {
			return nil, fmt.Errorf("could not retrieve block %s to extract its logs", parentHash)
		}
	}

	// We gather the logs across all the blocks in the canonical chain.
	logs := []*types.Log{}
	for _, hash := range blockHashes {
		blockLogs, err := s.storage.FetchLogs(hash)
		if err != nil {
			if errors.Is(err, errutil.ErrNotFound) {
				break // Blocks before the genesis rollup do not have associated logs (or block state).
			}
			return nil, fmt.Errorf("could not fetch logs for block hash. Cause: %w", err)
		}
		logs = append(logs, blockLogs...)
	}

	// We proceed in this way instead of calling `FetchHeadRollup` because we want to ensure the chain has not advanced
	// causing a head block/head rollup mismatch.
	headBlockState, err := s.storage.FetchBlockState(headBlock.Hash())
	if err != nil {
		return nil, fmt.Errorf("could not filter logs as block state for head block could not be retrieved. Cause: %w", err)
	}
	return s.FilterLogs(logs, headBlockState.HeadRollup, account, filter)
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
		if isRelevant(logItem, userAddrs, account, filter) {
			filteredLogs = append(filteredLogs, logItem)
		}
	}

	return filteredLogs, nil
}

// GetSubscribedLogsEncrypted returns, for each subscription, the logs filtered and encrypted with the appropriate
// viewing key.
func (s *SubscriptionManager) GetSubscribedLogsEncrypted(logs []*types.Log, rollupHash common.L2RootHash) (map[gethrpc.ID][]byte, error) {
	filteredLogs, err := s.getSubscribedLogs(logs, rollupHash)
	if err != nil {
		return nil, fmt.Errorf("could not get subscribed logs. Cause: %w", err)
	}
	return s.encryptLogs(filteredLogs)
}

// Filters out irrelevant logs, those that are not subscribed to, and those the subscription has seen before, and
// organises them by their subscribing ID.
func (s *SubscriptionManager) getSubscribedLogs(logs []*types.Log, rollupHash common.L2RootHash) (map[gethrpc.ID][]*types.Log, error) {
	relevantLogsByID := map[gethrpc.ID][]*types.Log{}

	// If there are no subscriptions, we return early, to avoid the overhead of creating the state DB.
	if s.getNumberOfSubsThreadsafe() == 0 {
		return map[gethrpc.ID][]*types.Log{}, nil
	}

	stateDB, err := s.storage.CreateStateDB(rollupHash)
	if err != nil {
		return nil, fmt.Errorf("could not create stateDB to extract user addresses. Cause: %w", err)
	}

	for _, logItem := range logs {
		userAddrs := getUserAddrsFromLogTopics(logItem, stateDB)
		s.updateRelevantLogs(logItem, userAddrs, relevantLogsByID)
	}

	return relevantLogsByID, nil
}

// Encrypts each log with the appropriate viewing key.
func (s *SubscriptionManager) encryptLogs(logsByID map[gethrpc.ID][]*types.Log) (map[gethrpc.ID][]byte, error) {
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

// Extracts the user addresses from the topics.
func getUserAddrsFromLogTopics(log *types.Log, db *state.StateDB) []string {
	var userAddrs []string

	for idx, topic := range log.Topics {
		// The first topic is always the hash of the event.
		if idx == 0 {
			continue
		}

		potentialAddr := gethcommon.HexToAddress(topic.Hex())

		// A user address must have (at least) 12 leading zero bytes, since addresses are 20 bytes long, while hashes
		// are 32.
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

// For each subscription, updates the relevant logs in the provided map.
func (s *SubscriptionManager) updateRelevantLogs(logItem *types.Log, userAddrs []string, relevantLogsByID map[gethrpc.ID][]*types.Log) {
	s.subscriptionMutex.RLock()
	defer s.subscriptionMutex.RUnlock()

	for subscriptionID, subscription := range s.subscriptions {
		// We ignore irrelevant logs.
		if !isRelevant(logItem, userAddrs, subscription.Account, subscription.Filter) {
			continue
		}

		// We update the relevant logs for the subscription.
		if relevantLogsByID[subscriptionID] == nil {
			relevantLogsByID[subscriptionID] = []*types.Log{}
		}
		relevantLogsByID[subscriptionID] = append(relevantLogsByID[subscriptionID], logItem)
	}
}

// Locks the subscription map and retrieves the subscription with subID, or (nil, false) if so such subscription is found.
func (s *SubscriptionManager) getSubscriptionThreadsafe(subID gethrpc.ID) (*common.LogSubscription, bool) {
	s.subscriptionMutex.RLock()
	defer s.subscriptionMutex.RUnlock()

	subscription, found := s.subscriptions[subID]
	return subscription, found
}

// Locks the subscription map and retrieves the number of subscriptions.
func (s *SubscriptionManager) getNumberOfSubsThreadsafe() int {
	s.subscriptionMutex.RLock()
	defer s.subscriptionMutex.RUnlock()

	return len(s.subscriptions)
}

// Indicates whether the log is relevant for the subscription.
func isRelevant(logItem *types.Log, userAddrs []string, account *gethcommon.Address, filter *filters.FilterCriteria) bool {
	return userAddrsContainAccount(account, userAddrs) && logMatchesFilter(logItem, filter)
}

// Indicates whether the account is contained in the user addresses.
func userAddrsContainAccount(account *gethcommon.Address, userAddrs []string) bool {
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

// Applies `filterLogs`, below, to determine whether the log should be filtered out based on the user's subscription
// criteria.
func logMatchesFilter(log *types.Log, filterCriteria *filters.FilterCriteria) bool {
	filteredLogs := filterLogs([]*types.Log{log}, filterCriteria.FromBlock, filterCriteria.ToBlock, filterCriteria.Addresses, filterCriteria.Topics)
	return len(filteredLogs) != 0
}

// Lifted from eth/filters/filter.go in the go-ethereum repository.
// filterLogs creates a slice of logs matching the given criteria.
func filterLogs(logs []*types.Log, fromBlock, toBlock *big.Int, addresses []gethcommon.Address, topics [][]gethcommon.Hash) []*types.Log { //nolint:gocognit
	var ret []*types.Log
Logs:
	for _, logItem := range logs {
		if fromBlock != nil && fromBlock.Int64() >= 0 && fromBlock.Uint64() > logItem.BlockNumber {
			continue
		}
		if toBlock != nil && toBlock.Int64() >= 0 && toBlock.Uint64() < logItem.BlockNumber {
			continue
		}

		if len(addresses) > 0 && !includes(addresses, logItem.Address) {
			continue
		}
		// If the to filtered topics is greater than the amount of topics in logs, skip.
		if len(topics) > len(logItem.Topics) {
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
