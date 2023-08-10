package enclave

import (
	"fmt"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/kamilsk/breaker"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/host"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/config"
	"github.com/obscuronet/go-obscuro/go/host/db"
	"github.com/obscuronet/go-obscuro/go/responses"
	"github.com/pkg/errors"

	gethlog "github.com/ethereum/go-ethereum/log"
)

// Service is a host service that provides access to the enclave(s) - it handles failover, load balancing, circuit breaking when a host has multiple enclaves
type Service struct {
	hostIdentity host.Identity
	p2p          host.P2P
	// eventually this service will support multiple enclaves for HA, but currently there's only one
	// The service goes via the Guardian to talk to the enclave (because guardian knows if the enclave is healthy etc.)
	enclaveGuardian *Guardian

	running atomic.Bool
	logger  gethlog.Logger
}

func NewService(
	cfg *config.HostConfig,
	p2p host.P2P,
	l1Publisher host.L1Publisher,
	l1RepoService host.L1BlockRepository,
	l2RepoService host.L2BatchRepository,
	hostIdentity host.Identity,
	enclaveClient common.Enclave,
	database *db.DB,
	interrupter breaker.Interface,
	logger gethlog.Logger,
) *Service {
	return &Service{
		hostIdentity: hostIdentity,
		p2p:          p2p,
		enclaveGuardian: NewGuardian(
			cfg,
			hostIdentity,
			p2p,
			l1Publisher,
			l1RepoService,
			l2RepoService,
			enclaveClient,
			database,
			interrupter,
			logger,
		),
		logger: logger,
	}
}

func (e *Service) Start() error {
	e.running.Store(true)
	return e.enclaveGuardian.Start()
}

func (e *Service) Stop() error {
	e.running.Store(false)
	return e.enclaveGuardian.Stop()
}

func (e *Service) HealthStatus() host.HealthStatus {
	if !e.running.Load() {
		return &host.BasicErrHealthStatus{ErrMsg: "not running"}
	}

	// check the enclave health, which in turn checks the DB health
	enclaveHealthy, err := e.enclaveGuardian.enclaveClient.HealthCheck()
	if err != nil {
		return &host.BasicErrHealthStatus{ErrMsg: fmt.Sprintf("unable to HealthCheck enclave - %s", err.Error())}
	} else if !enclaveHealthy {
		return &host.BasicErrHealthStatus{ErrMsg: "enclave reported itself as not healthy"}
	}

	if !e.enclaveGuardian.GetEnclaveState().InSyncWithL1() {
		return &host.BasicErrHealthStatus{ErrMsg: "enclave not in sync with L1"}
	}

	// empty error msg means healthy
	return &host.BasicErrHealthStatus{ErrMsg: ""}
}

func (e *Service) GetEnclaveClient() common.Enclave {
	return e.enclaveGuardian.GetEnclaveClient()
}

func (e *Service) SubmitAndBroadcastTx(encryptedParams common.EncryptedParamsSendRawTx) (*responses.RawTx, error) {
	encryptedTx := common.EncryptedTx(encryptedParams)

	enclaveResponse, sysError := e.enclaveGuardian.GetEnclaveClient().SubmitTx(encryptedTx)
	if sysError != nil {
		e.logger.Warn("Could not submit transaction due to sysError.", log.ErrKey, sysError)
		return nil, sysError
	}
	if enclaveResponse.Error() != nil {
		e.logger.Trace("Could not submit transaction.", log.ErrKey, enclaveResponse.Error())
		return enclaveResponse, nil //nolint: nilerr
	}

	if !e.hostIdentity.IsSequencer {
		err := e.p2p.SendTxToSequencer(encryptedTx)
		if err != nil {
			return nil, fmt.Errorf("could not broadcast transaction to sequencer. Cause: %w", err)
		}
	}

	return enclaveResponse, nil
}

func (e *Service) Subscribe(id rpc.ID, encryptedParams common.EncryptedParamsLogSubscription, matchedLogsCh chan []byte) error {
	err := e.enclaveGuardian.GetEnclaveClient().Subscribe(id, encryptedParams)
	if err != nil {
		return errors.Wrap(err, "could not create subscription with enclave")
	}
	return e.enclaveGuardian.logEventManager.Subscribe(id, matchedLogsCh)
}

func (e *Service) Unsubscribe(id rpc.ID) error {
	e.enclaveGuardian.logEventManager.Unsubscribe(id)
	err := e.enclaveGuardian.GetEnclaveClient().Unsubscribe(id)
	if err != nil {
		return errors.Wrap(err, "could not terminate enclave subscription")
	}
	return nil
}
