package events

import (
	"sync"

	"github.com/obscuronet/go-obscuro/go/common/host"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/pkg/errors"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/obscuronet/go-obscuro/go/common"
)

// LogEventManager manages the routing of logs back to their subscribers.
// todo (@matt) currently, this operates as a service but maybe it would make more sense to be owned by enclave service?
type LogEventManager struct {
	enclaveService    common.Enclave
	subscriptions     map[rpc.ID]*subscription // The channels that logs are sent to, one per subscription
	subscriptionMutex *sync.RWMutex
	logger            gethlog.Logger
}

func NewLogEventManager(enclaveService common.Enclave, logger gethlog.Logger) *LogEventManager {
	return &LogEventManager{
		enclaveService:    enclaveService,
		subscriptions:     map[rpc.ID]*subscription{},
		subscriptionMutex: &sync.RWMutex{},
		logger:            logger,
	}
}

func (l *LogEventManager) HealthStatus() host.HealthStatus {
	// always healthy for now
	return &host.BasicErrHealthStatus{ErrMsg: ""}
}

func (l *LogEventManager) Subscribe(id rpc.ID, encryptedLogSubscription common.EncryptedParamsLogSubscription, matchedLogsCh chan []byte) error {
	err := l.enclaveService.Subscribe(id, encryptedLogSubscription)
	if err != nil {
		return errors.Wrap(err, "could not create subscription with enclave")
	}
	l.subscriptionMutex.Lock()
	defer l.subscriptionMutex.Unlock()

	l.subscriptions[id] = &subscription{ch: matchedLogsCh}
	return nil
}

func (l *LogEventManager) Unsubscribe(id rpc.ID) {
	err := l.enclaveService.Unsubscribe(id)
	if err != nil {
		l.logger.Warn("could not terminate enclave subscription", log.ErrKey, err)
	}
	l.subscriptionMutex.Lock()
	defer l.subscriptionMutex.Unlock()

	logSubscription, found := l.subscriptions[id]
	if found {
		close(logSubscription.ch)
		delete(l.subscriptions, id)
	}
}

// SendLogsToSubscribers distributes logs to subscribed clients.
func (l *LogEventManager) SendLogsToSubscribers(result *common.EncryptedSubscriptionLogs) {
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
