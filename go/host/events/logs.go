package events

import (
	"sync"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/obscuronet/go-obscuro/go/common"
)

// SubscriptionManager manages the routing of logs back to their subscribers.
type SubscriptionManager struct {
	subscriptions     map[rpc.ID]*subscription // The channels that logs are sent to, one per subscription
	subscriptionMutex *sync.RWMutex
}

func NewSubscriptionManager() *SubscriptionManager {
	return &SubscriptionManager{
		subscriptions:     map[rpc.ID]*subscription{},
		subscriptionMutex: &sync.RWMutex{},
	}
}

func (l *SubscriptionManager) Subscribe(id rpc.ID, matchedLogsCh chan []byte) error {
	l.subscriptionMutex.Lock()
	defer l.subscriptionMutex.Unlock()

	l.subscriptions[id] = &subscription{ch: matchedLogsCh}
	return nil
}

func (l *SubscriptionManager) Unsubscribe(id rpc.ID) {
	l.subscriptionMutex.Lock()
	defer l.subscriptionMutex.Unlock()

	logSubscription, found := l.subscriptions[id]
	if found {
		close(logSubscription.ch)
		delete(l.subscriptions, id)
	}
}

// SendLogsToSubscribers distributes logs to subscribed clients.
func (l *SubscriptionManager) SendLogsToSubscribers(result *common.EncryptedSubscriptionLogs) {
	l.subscriptionMutex.RLock()
	defer l.subscriptionMutex.RUnlock()

	for id, encryptedLogs := range *result {
		logSub, found := l.subscriptions[id]
		if !found {
			continue
		}
		logSub.ch <- encryptedLogs
	}
}

// Pairs the latest seen rollup for a log subscription with the channel on which new logs should be sent.
type subscription struct {
	ch chan []byte // The channel that logs for this subscription are sent to.
}
