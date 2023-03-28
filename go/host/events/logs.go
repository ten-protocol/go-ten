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
	l.subscriptionMutex.Lock()
	defer l.subscriptionMutex.Unlock()

	logSubscription, found := l.subscriptions[id]
	if found {
		close(logSubscription.ch)
		delete(l.subscriptions, id)
	}
}

// SendLogsToSubscribers distributes logs to subscribed clients.
func (l *LogEventManager) SendLogsToSubscribers(result *common.BlockSubmissionResponse) {
	l.subscriptionMutex.RLock()
	defer l.subscriptionMutex.RUnlock()

	for id, encryptedLogs := range result.SubscribedLogs {
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
