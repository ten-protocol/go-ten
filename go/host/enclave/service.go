package enclave

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync/atomic"
	"time"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/errutil"
	"github.com/ten-protocol/go-ten/go/common/host"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/ethadapter"
	hostconfig "github.com/ten-protocol/go-ten/go/host/config"
	"github.com/ten-protocol/go-ten/go/responses"
	"github.com/ten-protocol/go-ten/lib/gethfork/rpc"
)

const (
	_promoteSeqRetryInterval = 1 * time.Second
)

var _noActiveSequencer = &common.EnclaveID{}

// This private interface enforces the services that the enclaves service depends on
type enclaveServiceLocator interface {
	L1Data() host.L1DataService
	L1Publisher() host.L1Publisher
	L2Repo() host.L2BatchRepository
	P2P() host.P2P
}

// Service is a host service that provides access to the enclave(s) - it handles failover, load balancing, circuit breaking when a host has multiple enclaves
type Service struct {
	hostData host.Identity
	sl       enclaveServiceLocator

	// The service goes via the Guardians to talk to the enclave (because guardian knows if the enclave is healthy etc.)
	enclaveGuardians  []*Guardian
	activeSequencerID atomic.Pointer[common.EnclaveID] // atomic pointer for thread safety

	// rollup config
	rollupInterval time.Duration
	blockTime      time.Duration
	maxRollupSize  uint64

	running atomic.Bool
	logger  gethlog.Logger
}

func NewService(config *hostconfig.HostConfig, hostData host.Identity, serviceLocator enclaveServiceLocator, enclaveGuardians []*Guardian, logger gethlog.Logger) *Service {
	return &Service{
		hostData:         hostData,
		sl:               serviceLocator,
		enclaveGuardians: enclaveGuardians,
		rollupInterval:   config.RollupInterval,
		blockTime:        config.L1BlockTime,
		maxRollupSize:    config.MaxRollupSize,
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
		go e.managePeriodicRollups()
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

const batchCompressionFactor = 0.85

// managePeriodicRollups is a background goroutine that periodically produces a rollup
// where possible it will prefer to use a non-active sequencer enclave to avoid disrupting the production of batches
// note: this function runs in a separate goroutine for the lifetime of the service
func (e *Service) managePeriodicRollups() {
	e.logger.Info("Starting periodic rollups.")
	lastSuccessfulRollup := time.Now()

	time.Sleep(e.blockTime)

	for e.running.Load() {
		// block time seems a reasonable scaling cadence to check if rollup required, no need to check after every batch
		time.Sleep(e.blockTime)

		var rollupToPublish *common.CreateRollupResult
		var err error

		rollupRequired, fromBatch := e.isRollupRequired(lastSuccessfulRollup)
		if !rollupRequired {
			// the rollup required check contains appropriate logging, so no need to log here
			continue
		}

		// find a client to produce rollup. Skip active sequencer at first, then try active sequencer if needed.
		for _, guardian := range e.enclaveGuardians {
			if guardian.state.IsEnclaveActiveSequencer() {
				continue // skip active sequencer for now
			}

			rollupToPublish, err = e.prepareRollup(guardian, fromBatch)
			if err != nil {
				e.logger.Error("Enclave failed to prepare rollup.", log.ErrKey, err)
				continue // try next guardian
			}
		}

		if rollupToPublish == nil {
			// if we didn't find a non-active sequencer to produce the rollup, try the active sequencer
			guardian, err := e.getActiveSequencerGuardian()
			if err != nil {
				e.logger.Error("no active sequencer guardian found, cannot prepare rollup", log.ErrKey, err)
				continue // try again later
			}
			rollupToPublish, err = e.prepareRollup(guardian, fromBatch)
			if err != nil {
				e.logger.Error("Seq failed to prepare rollup.", log.ErrKey, err)
				continue // try again later
			}
		}

		// this method waits until the receipt is received
		err = e.sl.L1Publisher().PublishBlob(*rollupToPublish)
		if err != nil {
			e.logger.Error("Failed to publish rollup ", log.ErrKey, err)
			continue // try again later
		}
		lastSuccessfulRollup = time.Now()
	}

	e.logger.Info("Stopping periodic rollups.")
}

// returns true if a rollup is required, and the batch number to start from
func (e *Service) isRollupRequired(lastSuccessfulRollup time.Time) (bool, uint64) {
	if e.activeSequencerID.Load() == _noActiveSequencer {
		e.logger.Debug("No active sequencer, skipping periodic rollup.")
		return false, 0
	}

	fromBatch, err := e.getLatestBatchNo()
	if err != nil {
		e.logger.Error("Encountered error while trying to retrieve latest batch", log.ErrKey, err)
		return false, 0
	}

	// estimate the size of a compressed rollup
	availBatchesSumSize, err := e.calculateNonRolledupBatchesSize(fromBatch)
	if err != nil {
		e.logger.Error("Unable to estimate the size of the current rollup", log.ErrKey, err, "from_batch", fromBatch)
		// Note: this should not happen. If it does, we will assume the size is 0, meaning only time will trigger a rollup
		availBatchesSumSize = 0
	}

	// adjust the availBatchesSumSize
	estimatedRunningRollupSize := uint64(float64(availBatchesSumSize) * batchCompressionFactor)

	// produce and issue rollup when either:
	// it has passed g.rollupInterval from last lastSuccessfulRollup
	// or the size of accumulated batches is > g.maxRollupSize
	timeExpired := time.Since(lastSuccessfulRollup) > e.rollupInterval
	sizeExceeded := estimatedRunningRollupSize >= e.maxRollupSize

	return timeExpired || sizeExceeded, fromBatch
}

func (e *Service) prepareRollup(guardian *Guardian, fromBatch uint64) (*common.CreateRollupResult, error) {
	enclID := guardian.GetEnclaveID()
	e.logger.Info("Attempting to produce rollup.", log.EnclaveIDKey, enclID)
	result, err := guardian.GetEnclaveClient().CreateRollup(context.Background(), fromBatch)
	if err != nil {
		e.logger.Info("Unable to produce rollup", log.EnclaveIDKey, enclID, log.ErrKey, err)
		return nil, err
	}
	rollup, err := ethadapter.ReconstructRollup(result.Blobs)
	if err != nil {
		e.logger.Error("Failed to reconstruct rollup", log.ErrKey, err)
		return nil, err
	}

	canonBlock, err := e.sl.L1Data().FetchBlockByHeight(rollup.Header.CompressionL1Number)
	if err != nil {
		e.logger.Error("Failed to fetch canonical block for rollup", log.ErrKey, err)
		return nil, err
	}

	// only publish if the block used for compression is canonical
	if canonBlock.Hash() != rollup.Header.CompressionL1Head {
		e.logger.Info("Skipping rollup publication because compression block is not canonical.", "block", canonBlock.Hash())
		return nil, fmt.Errorf("compression block is not canonical, block=%s", canonBlock.Hash())
	}
	return result, nil
}

func (e *Service) getLatestBatchNo() (uint64, error) {
	lastBatchNo, err := e.sl.L1Publisher().FetchLatestSeqNo()
	if err != nil {
		return 0, err
	}
	fromBatch := lastBatchNo.Uint64()
	if lastBatchNo.Uint64() > common.L2GenesisSeqNo {
		fromBatch++
	}
	return fromBatch, nil
}

func (e *Service) calculateNonRolledupBatchesSize(seqNo uint64) (uint64, error) {
	if seqNo == 0 { // don't calculate for seqNo 0 batches
		return 0, nil
	}

	return e.sl.L2Repo().EstimateRollupSize(context.Background(), big.NewInt(int64(seqNo)))
}

func (e *Service) getActiveSequencerGuardian() (*Guardian, error) {
	activeSequencerID := e.activeSequencerID.Load()
	if activeSequencerID == _noActiveSequencer {
		return nil, errors.New("no active sequencer found")
	}

	for _, guardian := range e.enclaveGuardians {
		if *(guardian.GetEnclaveID()) == *activeSequencerID {
			return guardian, nil
		}
	}
	return nil, errors.New("active sequencer not found in guardians")
}
