package events

import (
	"sync"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/obscuronet/go-obscuro/go/common"
)

// LogEventManager manages the routing of logs back to their subscribers.
type LogEventManager struct {
	subscriptions     map[rpc.ID]*subscription // The channels that logs are sent to, one per subscription
	subscriptionMutex *sync.RWMutex
	logger            gethlog.Logger
}

func NewLogEventManager(logger gethlog.Logger) LogEventManager {
	return LogEventManager{
		subscriptions:     map[rpc.ID]*subscription{},
		subscriptionMutex: &sync.RWMutex{},
		logger:            logger,
	}
}

// AddSubscription adds a subscription to the set of managed subscriptions.
func (l *LogEventManager) AddSubscription(id rpc.ID, matchedLogsCh chan []byte) {
	l.subscriptionMutex.Lock()
	defer l.subscriptionMutex.Unlock()

	l.subscriptions[id] = &subscription{ch: matchedLogsCh}
}

// RemoveSubscription removes a subscription from the set of managed subscriptions.
func (l *LogEventManager) RemoveSubscription(id rpc.ID) {
	logSubscription, found := l.getSubscriptionThreadsafe(id)
	if found {
		close(logSubscription.ch)

		l.subscriptionMutex.RLock()
		defer l.subscriptionMutex.RUnlock()
		delete(l.subscriptions, id)
	}
}

// SendLogsToSubscribers distributes logs to subscribed clients.
func (l *LogEventManager) SendLogsToSubscribers(result *common.BlockSubmissionResponse) {
	for id, encryptedLogs := range result.SubscribedLogs {
		logSub, found := l.getSubscriptionThreadsafe(id)
		if !found {
			continue
		}
		logSub.ch <- encryptedLogs
	}
}

// Locks the subscription map and retrieves the subscription with subID, or (nil, false) if so such subscription is found.
func (l *LogEventManager) getSubscriptionThreadsafe(subID rpc.ID) (*subscription, bool) {
	l.subscriptionMutex.RLock()
	defer l.subscriptionMutex.RUnlock()

	sub, found := l.subscriptions[subID]
	return sub, found
}

// Pairs the latest seen rollup for a log subscription with the channel on which new logs should be sent.
type subscription struct {
	ch chan []byte // The channel that logs for this subscription are sent to.
}
