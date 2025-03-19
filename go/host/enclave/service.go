package enclave

import (
	"context"
	"fmt"
	"math/big"
	"sync/atomic"
	"time"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/errutil"
	"github.com/ten-protocol/go-ten/go/common/host"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/responses"
	"github.com/ten-protocol/go-ten/lib/gethfork/rpc"
)

const (
	_promoteSeqRetryInterval = 1 * time.Second
)

var _noActiveSequencer = &common.EnclaveID{}

// This private interface enforces the services that the enclaves service depends on
type enclaveServiceLocator interface {
	P2P() host.P2P
}

// Service is a host service that provides access to the enclave(s) - it handles failover, load balancing, circuit breaking when a host has multiple enclaves
type Service struct {
	hostData host.Identity
	sl       enclaveServiceLocator

	// The service goes via the Guardians to talk to the enclave (because guardian knows if the enclave is healthy etc.)
	enclaveGuardians  []*Guardian
	activeSequencerID atomic.Pointer[common.EnclaveID] // atomic pointer for thread safety

	running atomic.Bool
	logger  gethlog.Logger
}

func NewService(hostData host.Identity, serviceLocator enclaveServiceLocator, enclaveGuardians []*Guardian, logger gethlog.Logger) *Service {
	return &Service{
		hostData:         hostData,
		sl:               serviceLocator,
		enclaveGuardians: enclaveGuardians,
		logger:           logger,
	}
}

func (e *Service) Start() error {
	e.running.Store(true)
	for _, guardian := range e.enclaveGuardians {
		if err := guardian.Start(); err != nil {
			// abandon starting the rest of the guardians if one fails
			return err
		}
	}
	e.activeSequencerID.Store(_noActiveSequencer)
	if e.hostData.IsSequencer {
		go e.promoteNewActiveSequencer()
	}
	return nil
}

func (e *Service) Stop() error {
	e.running.Store(false)
	var errors []error
	for i, guardian := range e.enclaveGuardians {
		if err := guardian.Stop(); err != nil {
			errors = append(errors, fmt.Errorf("error stopping enclave guardian [%d]: %w", i, err))
		}
	}
	if len(errors) > 0 {
		return fmt.Errorf("errors stopping enclave guardians: %v", errors)
	}
	return nil
}

func (e *Service) HealthStatus(ctx context.Context) host.HealthStatus {
	if !e.running.Load() {
		return &host.BasicErrHealthStatus{ErrMsg: "not running"}
	}

	errors := make([]error, 0, len(e.enclaveGuardians))

	for i, guardian := range e.enclaveGuardians {
		// check the enclave health, which in turn checks the DB health
		enclaveHealthy, err := guardian.enclaveClient.HealthCheck(ctx)
		if err != nil {
			errors = append(errors, fmt.Errorf("unable to HealthCheck enclave[%d] - %w", i, err))
		} else if !enclaveHealthy {
			errors = append(errors, fmt.Errorf("enclave[%d] reported itself not healthy", i))
		}

		if !guardian.GetEnclaveState().InSyncWithL1() {
			errors = append(errors, fmt.Errorf("enclave[%d - %s] not in sync with L1 - %s", i, guardian.GetEnclaveID(), guardian.GetEnclaveState()))
		}
	}

	// empty error msg means healthy
	return &host.GroupErrsHealthStatus{Errors: errors}
}

// LookupBatchBySeqNo is used to fetch batch data from the enclave - it is only used as a fallback for the sequencer
// host if it's missing a batch (other host services should use L2Repo to fetch batch data)
func (e *Service) LookupBatchBySeqNo(ctx context.Context, seqNo *big.Int) (*common.ExtBatch, error) {
	// todo (@matt) revisit this flow to make sure it handles HA scenarios properly
	hg := e.enclaveGuardians[0]
	state := hg.GetEnclaveState()
	if state.GetEnclaveL2Head().Cmp(seqNo) < 0 {
		return nil, errutil.ErrNotFound
	}
	client := hg.GetEnclaveClient()
	return client.GetBatchBySeqNo(ctx, seqNo.Uint64())
}

func (e *Service) GetEnclaveClient() common.Enclave {
	// for now we always return first guardian's enclave client
	// in future be good to load balance and failover but need to improve subscribe/unsubscribe (unsubscribe from same enclave)
	return e.enclaveGuardians[0].GetEnclaveClient()
}

func (e *Service) GetEnclaveClients() []common.Enclave {
	clients := make([]common.Enclave, len(e.enclaveGuardians))
	for i, guardian := range e.enclaveGuardians {
		clients[i] = guardian.enclaveClient
	}
	return clients
}

// NotifyUnavailable is called by enclave guardians when they detect that the enclave is unavailable.
// If this is a sequencer host then this function will start a search for a live standby enclave to promote to active sequencer.
func (e *Service) NotifyUnavailable(enclaveID *common.EnclaveID) {
	if e.activeSequencerID.Load() != enclaveID {
		e.logger.Debug("Failed enclave is not an active sequencer, no action required.", log.EnclaveIDKey, enclaveID)
		return
	}
	failedEnclaveIdx := -1
	for i, guardian := range e.enclaveGuardians {
		if *(guardian.GetEnclaveID()) == *enclaveID {
			failedEnclaveIdx = i
			break
		}
	}
	if failedEnclaveIdx == -1 {
		e.logger.Warn("Could not find failed enclave to evict.", log.EnclaveIDKey, enclaveID)
		return
	}
	e.enclaveGuardians[failedEnclaveIdx].DemoteFromActiveSequencer()
	e.activeSequencerID.Store(_noActiveSequencer)

	go e.promoteNewActiveSequencer()
}

func (e *Service) SubmitAndBroadcastTx(ctx context.Context, encryptedParams common.EncryptedRequest) (*responses.RawTx, error) {
	encryptedTx := common.EncryptedTx(encryptedParams)

	enclaveResponse, sysError := e.GetEnclaveClient().EncryptedRPC(ctx, encryptedParams)
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
	return e.GetEnclaveClient().Subscribe(context.Background(), id, encryptedParams)
}

func (e *Service) Unsubscribe(id rpc.ID) error {
	return e.GetEnclaveClient().Unsubscribe(id)
}

// promoteNewActiveSequencer is a background goroutine that promotes a new active sequencer at startup or when the current one fails.
// It will never give up, it just cycles through current enclaves until one can be successfully promoted.
func (e *Service) promoteNewActiveSequencer() {
	for e.activeSequencerID.Load() == _noActiveSequencer && e.running.Load() {
		for _, guardian := range e.enclaveGuardians {
			enclID := guardian.GetEnclaveID()
			e.logger.Info("Attempting to promote new sequencer.", log.EnclaveIDKey, enclID)
			err := guardian.PromoteToActiveSequencer()
			if err != nil {
				e.logger.Info("Unable to promote new sequencer.", log.EnclaveIDKey, enclID, log.ErrKey, err)
				continue
			}
			e.activeSequencerID.Store(enclID)
			e.logger.Warn("Successfully promoted new sequencer.", log.EnclaveIDKey, enclID)
			return
		}
		// wait for retry interval before trying again, enclaves may not be ready yet or may be awaiting permissioning
		time.Sleep(_promoteSeqRetryInterval)
	}
}
