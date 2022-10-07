package events

import (
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/log"
)

// LogEventManager manages the routing of logs back to their subscribers.
type LogEventManager struct {
	subscriptions map[rpc.ID]*subscription // The channels that logs are sent to, one per subscription
}

func NewLogEventManager() LogEventManager {
	return LogEventManager{
		subscriptions: map[rpc.ID]*subscription{},
	}
}

// AddSubscription adds a subscription to the set of managed subscriptions.
func (l *LogEventManager) AddSubscription(id rpc.ID, matchedLogsCh chan []byte) {
	l.subscriptions[id] = &subscription{ch: matchedLogsCh}
}

// RemoveSubscription removes a subscription from the set of managed subscriptions.
func (l *LogEventManager) RemoveSubscription(id rpc.ID) {
	logSubscription, found := l.subscriptions[id]
	if found {
		close(logSubscription.ch)
		delete(l.subscriptions, id)
	}
}

// SendLogsToSubscribers distributes logs to subscribed clients. We only send logs for rollups the subscription hasn't seen before.
func (l *LogEventManager) SendLogsToSubscribers(result common.BlockSubmissionResponse) {
	latestSeenRollupByID := map[rpc.ID]uint64{}

	for subscriptionID, encLogsByRollup := range result.SubscribedLogs {
		logSub, found := l.subscriptions[subscriptionID]
		if !found {
			log.Error("received a log for subscription with ID %s, but no such subscription exists", subscriptionID)
			continue
		}

		for rollupNumber, encryptedLogs := range encLogsByRollup {
			// We have received a log from a rollup this subscription hasn't seen before.
			if rollupNumber > logSub.latestSeenRollup {
				l.subscriptions[subscriptionID].ch <- encryptedLogs
			}

			// We update the latest rollup number if this is the highest one we've seen so far.
			currentLatestSeenRollup, exists := latestSeenRollupByID[subscriptionID]
			if !exists || rollupNumber > currentLatestSeenRollup {
				latestSeenRollupByID[subscriptionID] = rollupNumber
			}
		}
	}

	// We update the latest seen rollup for each subscription. We must do this in a separate loop, as if we update it
	// as we go, we may miss a set of logs if we process the logs for rollup N before we process those for rollup N-1.
	for subscriptionID, logSub := range l.subscriptions {
		newLatestSeenRollup, exists := latestSeenRollupByID[subscriptionID]
		if exists && newLatestSeenRollup > logSub.latestSeenRollup {
			logSub.latestSeenRollup = newLatestSeenRollup
		}
	}
}

// Pairs the latest seen rollup for a log subscription with the channel on which new logs should be sent.
type subscription struct {
	latestSeenRollup uint64      // The latest rollup for which logs have been distributed to this subscription.
	ch               chan []byte // The channel that logs for this subscription are sent to.
}
