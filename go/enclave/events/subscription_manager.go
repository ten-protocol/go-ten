package events

import (
	"encoding/json"
	"fmt"
	"math/big"
	"sync"

	"github.com/obscuronet/go-obscuro/go/enclave/core"

	"github.com/obscuronet/go-obscuro/go/common/log"

	"github.com/obscuronet/go-obscuro/go/enclave/storage"

	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/filters"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/enclave/rpc"
	"github.com/obscuronet/go-obscuro/go/enclave/vkhandler"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
	gethrpc "github.com/ethereum/go-ethereum/rpc"
)

const (
	// The leading zero bytes in a hash indicating that it is possibly an address, since it only has 20 bytes of data.
	zeroBytesHex = "000000000000000000000000"
)

// SubscriptionManager manages the creation/deletion of subscriptions, and the filtering and encryption of logs for
// active subscriptions.
type SubscriptionManager struct {
	rpcEncryptionManager *rpc.EncryptionManager
	storage              storage.Storage

	subscriptions     map[gethrpc.ID]*common.LogSubscription
	subscriptionMutex *sync.RWMutex // the mutex guards the subscriptions/lastHead pair

	logger gethlog.Logger
}

func NewSubscriptionManager(rpcEncryptionManager *rpc.EncryptionManager, storage storage.Storage, logger gethlog.Logger) *SubscriptionManager {
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

	subscription := &common.LogSubscription{}
	if err = rlp.DecodeBytes(encodedSubscription, subscription); err != nil {
		return fmt.Errorf("could not decocde log subscription from RLP. Cause: %w", err)
	}

	// create viewing key encryption handler for pushing future logs
	encryptor, err := vkhandler.New(subscription.Account, subscription.PublicViewingKey, subscription.Signature)
	if err != nil {
		return fmt.Errorf("unable to create vk encryption for request - %w", err)
	}
	subscription.VkHandler = encryptor

	s.subscriptionMutex.Lock()
	s.subscriptions[id] = subscription
	s.subscriptionMutex.Unlock()

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
func (s *SubscriptionManager) FilterLogs(logs []*types.Log, batchHash common.L2BatchHash, account *gethcommon.Address, filter *filters.FilterCriteria) ([]*types.Log, error) {
	var filteredLogs []*types.Log
	stateDB, err := s.storage.CreateStateDB(batchHash)
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

// GetSubscribedLogsForBatch - Retrieves and encrypts the logs for the batch in live mode.
// The assumption is that this function is called synchronously after the batch is produced
func (s *SubscriptionManager) GetSubscribedLogsForBatch(batch *core.Batch, receipts types.Receipts) (common.EncryptedSubscriptionLogs, error) {
	result := map[gethrpc.ID][]*types.Log{}

	var allLogs []*types.Log
	for _, receipt := range receipts {
		allLogs = append(allLogs, receipt.Logs...)
	}

	s.subscriptionMutex.RLock()
	defer s.subscriptionMutex.RUnlock()

	var err error
	for id, sub := range s.subscriptions {
		result[id], err = s.FilterLogs(allLogs, batch.Hash(), sub.Account, sub.Filter)
		if err != nil {
			s.logger.Error("Could not retrieve subscription logs", log.ErrKey, err)
			return nil, err
		}
	}
	// Encrypt the results
	return s.encryptLogs(result)
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

		encryptedLogs, err := subscription.VkHandler.Encrypt(jsonLogs)
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
