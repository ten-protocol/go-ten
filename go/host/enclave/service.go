package enclave

import (
	"fmt"
	"math/big"
	"sync/atomic"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/errutil"
	"github.com/obscuronet/go-obscuro/go/common/host"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/responses"
)

// This private interface enforces the services that the enclaves service depends on
type enclaveServiceLocator interface {
	P2P() host.P2P
}

// Service is a host service that provides access to the enclave(s) - it handles failover, load balancing, circuit breaking when a host has multiple enclaves
type Service struct {
	hostData host.Identity
	sl       enclaveServiceLocator
	// eventually this service will support multiple enclaves for HA, but currently there's only one
	// The service goes via the Guardian to talk to the enclave (because guardian knows if the enclave is healthy etc.)
	enclaveGuardian *Guardian

	running atomic.Bool
	logger  gethlog.Logger
}

func NewService(hostData host.Identity, serviceLocator enclaveServiceLocator, enclaveGuardian *Guardian, logger gethlog.Logger) *Service {
	return &Service{
		hostData:        hostData,
		sl:              serviceLocator,
		enclaveGuardian: enclaveGuardian,
		logger:          logger,
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

// LookupBatchBySeqNo is used to fetch batch data from the enclave - it is only used as a fallback for the sequencer
// host if it's missing a batch (other host services should use L2Repo to fetch batch data)
func (e *Service) LookupBatchBySeqNo(seqNo *big.Int) (*common.ExtBatch, error) {
	state := e.enclaveGuardian.GetEnclaveState()
	if state.GetEnclaveL2Head().Cmp(seqNo) < 0 {
		return nil, errutil.ErrNotFound
	}
	client := e.enclaveGuardian.GetEnclaveClient()
	return client.GetBatchBySeqNo(seqNo.Uint64())
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

	if !e.hostData.IsSequencer {
		err := e.sl.P2P().SendTxToSequencer(encryptedTx)
		if err != nil {
			return nil, fmt.Errorf("could not broadcast transaction to sequencer. Cause: %w", err)
		}
	}

	return enclaveResponse, nil
}

func (e *Service) Subscribe(id rpc.ID, encryptedParams common.EncryptedParamsLogSubscription) error {
	return e.enclaveGuardian.GetEnclaveClient().Subscribe(id, encryptedParams)
}

func (e *Service) Unsubscribe(id rpc.ID) error {
	return e.enclaveGuardian.GetEnclaveClient().Unsubscribe(id)
}
