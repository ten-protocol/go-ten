package events

import (
	"context"
	"sync"

	"github.com/pkg/errors"
	"github.com/ten-protocol/go-ten/go/common/host"
	"github.com/ten-protocol/go-ten/go/common/log"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/lib/gethfork/rpc"
)

type logSubsServiceLocator interface {
	Enclaves() host.EnclaveService
}

// LogEventManager manages the routing of logs back to their subscribers.
// todo (@matt) currently, this operates as a service but maybe it would make more sense to be owned by enclave service?
type LogEventManager struct {
	sl                logSubsServiceLocator
	subscriptions     map[rpc.ID]*subscription // The channels that logs are sent to, one per subscription
	subscriptionMutex *sync.RWMutex
	logger            gethlog.Logger
}

func NewLogEventManager(serviceLocator logSubsServiceLocator, logger gethlog.Logger) *LogEventManager {
	return &LogEventManager{
		sl:                serviceLocator,
		subscriptions:     map[rpc.ID]*subscription{},
		subscriptionMutex: &sync.RWMutex{},
		logger:            logger,
	}
}

func (l *LogEventManager) Start() error {
	return nil
}

func (l *LogEventManager) Stop() error {
	return nil
}

func (l *LogEventManager) HealthStatus(context.Context) host.HealthStatus {
	// always healthy for now
	return &host.BasicErrHealthStatus{ErrMsg: ""}
}

func (l *LogEventManager) Subscribe(id rpc.ID, encryptedLogSubscription common.EncryptedParamsLogSubscription, matchedLogsCh chan []byte) error {
	err := l.sl.Enclaves().Subscribe(id, encryptedLogSubscription)
	if err != nil {
		return errors.Wrap(err, "could not create subscription with enclave")
	}
	l.subscriptionMutex.Lock()
	defer l.subscriptionMutex.Unlock()

	l.subscriptions[id] = &subscription{ch: matchedLogsCh}
	return nil
}

func (l *LogEventManager) Unsubscribe(id rpc.ID) {
	enclaveUnsubErr := l.sl.Enclaves().Unsubscribe(id)
	if enclaveUnsubErr != nil {
		// this can happen when the client passes an invalid subscription id
		l.logger.Debug("Could not terminate enclave subscription", log.SubIDKey, id, log.ErrKey, enclaveUnsubErr)
	}
	l.subscriptionMutex.RLock()
	logSubscription, found := l.subscriptions[id]
	ch := logSubscription.ch
	l.subscriptionMutex.RUnlock()

	if found {
		l.subscriptionMutex.Lock()
		delete(l.subscriptions, id)
		l.subscriptionMutex.Unlock()
		if enclaveUnsubErr != nil {
			l.logger.Error("The subscription management between the host and the enclave is out of sync", log.SubIDKey, id, log.ErrKey, enclaveUnsubErr)
		}
		close(ch)
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

// Simple wrapper over the channel that logs for this subscription are sent to.
type subscription struct {
	ch chan []byte
}
