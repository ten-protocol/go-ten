package enclave

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ten-protocol/go-ten/go/ethadapter"
	hostconfig "github.com/ten-protocol/go-ten/go/host/config"

	"github.com/ten-protocol/go-ten/go/host/storage"

	"github.com/ten-protocol/go-ten/go/common/stopcontrol"

	gethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/ten-protocol/go-ten/go/common/gethutil"

	"github.com/ethereum/go-ethereum/core/types"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/pkg/errors"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/errutil"
	"github.com/ten-protocol/go-ten/go/common/host"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/common/retry"
	"github.com/ten-protocol/go-ten/go/host/l1"
)

const (
	// time between loops on mainLoop, will be retry time if things are failing
	_retryInterval = 100 * time.Millisecond

	// when enclave is healthy this is the time before we call its status (can be slow, is just a sanity check)
	_monitoringInterval = 1 * time.Second

	// when we have submitted request to L1 for the secret, how long do we wait for an answer before we retry
	_maxWaitForSecretResponse = 2 * time.Minute
)

// This private interface enforces the services that the guardian depends on
type guardianServiceLocator interface {
	P2P() host.P2P
	L1Publisher() host.L1Publisher
	L1Data() host.L1DataService
	L2Repo() host.L2BatchRepository
	LogSubs() host.LogSubscriptionManager
	Enclaves() host.EnclaveService
}

// Guardian is a host service which monitors an enclave, it's responsibilities include:
// - monitor the enclave state and feed it the data it needs
// - if it is an active sequencer then the guardian will trigger batch/rollup creation
// - guardian provides access to the enclave data and reports the enclave status for other services - acting as a gatekeeper
type Guardian struct {
	hostData      host.Identity
	state         *StateTracker // state machine that tracks our view of the enclave's state
	enclaveClient common.Enclave

	sl      guardianServiceLocator
	storage storage.Storage

	submitDataLock sync.Mutex // we only submit one block, batch or transaction to enclave at a time

	batchInterval      time.Duration
	rollupInterval     time.Duration
	blockTime          time.Duration
	crossChainInterval time.Duration
	l1StartHash        gethcommon.Hash
	maxRollupSize      uint64

	hostInterrupter *stopcontrol.StopControl // host hostInterrupter so we can stop quickly
	running         atomic.Bool

	logger    gethlog.Logger
	enclaveID *common.EnclaveID

	cleanupFuncs []func()
}

func NewGuardian(cfg *hostconfig.HostConfig, hostData host.Identity, serviceLocator guardianServiceLocator, enclaveClient common.Enclave, storage storage.Storage, interrupter *stopcontrol.StopControl, logger gethlog.Logger) *Guardian {
	return &Guardian{
		hostData:           hostData,
		state:              NewStateTracker(logger),
		enclaveClient:      enclaveClient,
		sl:                 serviceLocator,
		batchInterval:      cfg.BatchInterval,
		rollupInterval:     cfg.RollupInterval,
		l1StartHash:        cfg.L1StartHash,
		maxRollupSize:      cfg.MaxRollupSize,
		blockTime:          cfg.L1BlockTime,
		crossChainInterval: cfg.CrossChainInterval,
		storage:            storage,
		hostInterrupter:    interrupter,
		logger:             logger,
	}
}

func (g *Guardian) Start() error {
	// sanity check, Start() spawns new go-routines and should only be called once
	if g.running.Load() {
		return errors.New("guardian already started")
	}

	g.running.Store(true)
	g.hostInterrupter.OnStop(func() {
		g.logger.Info("Guardian stopping (because host is stopping)")
		g.running.Store(false)
	})

	// Identify the enclave before starting (the enclave generates its ID immediately at startup)
	// (retry until we get the enclave ID or the host is stopping)
	enclaveConnectRetries := 0
	for g.enclaveID == nil && g.running.Load() {
		enclID, err := g.enclaveClient.EnclaveID(context.Background())
		if err != nil {
			if enclaveConnectRetries%10 == 0 { // avoid spamming the log while waiting for the enclave to start
				g.logger.Warn("could not get enclave ID, retrying", log.ErrKey, err, "retries", enclaveConnectRetries)
			} else {
				g.logger.Trace("could not get enclave ID, retrying", log.ErrKey, err, "retries", enclaveConnectRetries)
			}
			enclaveConnectRetries++
			time.Sleep(_retryInterval)
			continue
		}
		g.enclaveID = &enclID
		// include the enclave ID in guardian log messages (for multi-enclave nodes)
		g.logger = g.logger.New(log.EnclaveIDKey, g.enclaveID)
		// recreate status with new logger
		g.state = NewStateTracker(g.logger)
		g.state.OnReceivedBatch(g.sl.L2Repo().FetchLatestBatchSeqNo())
		g.logger.Info("Starting guardian process.")
	}

	go g.mainLoop()

	// subscribe for L1 and P2P data
	txUnsub := g.sl.P2P().SubscribeForTx(g)

	// note: not keeping the unsubscribe functions because the lifespan of the guardian is the same as the host
	l1Unsub := g.sl.L1Data().Subscribe(g)
	batchesUnsub := g.sl.L2Repo().SubscribeNewBatches(g)

	g.cleanupFuncs = []func(){txUnsub, l1Unsub, batchesUnsub}

	// start streaming data from the enclave
	go g.streamEnclaveData()

	return nil
}

func (g *Guardian) Stop() error {
	g.running.Store(false)
	err := g.enclaveClient.Stop()
	if err != nil {
		g.logger.Error("error stopping enclave", log.ErrKey, err)
	}

	err = g.enclaveClient.StopClient()
	if err != nil {
		g.logger.Error("error stopping enclave client", log.ErrKey, err)
	}

	// unsubscribe
	for _, cleanup := range g.cleanupFuncs {
		cleanup()
	}

	return nil
}

func (g *Guardian) HealthStatus(context.Context) host.HealthStatus {
	// todo (@matt) do proper health status based on enclave state
	errMsg := ""
	if !g.running.Load() {
		errMsg = "not running"
	}
	return &host.BasicErrHealthStatus{ErrMsg: errMsg}
}

func (g *Guardian) IsLive() bool {
	return g.state.IsLive()
}

func (g *Guardian) InSyncWithL1() bool {
	return g.state.InSyncWithL1()
}

func (g *Guardian) IsEnclaveL2AheadOfHost() bool {
	return g.state.IsEnclaveAheadOfHost()
}

func (g *Guardian) GetEnclaveState() *StateTracker {
	return g.state
}

// GetEnclaveClient returns the enclave client for use by other services
// todo (@matt) avoid exposing client directly and return errors if enclave is not ready for requests
func (g *Guardian) GetEnclaveClient() common.Enclave {
	return g.enclaveClient
}

func (g *Guardian) GetEnclaveID() *common.EnclaveID {
	return g.enclaveID
}

func (g *Guardian) PromoteToActiveSequencer() error {
	if g.state.IsEnclaveActiveSequencer() {
		// this shouldn't happen and shouldn't be an issue if it does, but good to have visibility on it
		g.logger.Error("Unable to promote to active sequencer, already active")
		return nil
	}
	l2Head := g.state.GetEnclaveL2Head()
	if l2Head != nil && l2Head.Cmp(big.NewInt(0)) > 0 && !g.state.IsLive() {
		// enclave has an L2 head so it's not just starting up, it can't be promoted to active sequencer until it is
		// up-to-date with the L2 head according to the host's database.
		return errors.New("cannot promote to active sequencer while behind the L2 head, it must finish syncing first")
	}
	err := g.enclaveClient.MakeActive()
	if err != nil {
		return errors.Wrap(err, "could not promote enclave to active sequencer")
	}
	g.state.SetActiveSequencer(true)
	return nil
}

// DemoteFromActiveSequencer stops the guardian from being the active sequencer,
// stopping batch and rollup production. The enclave can be promoted again later if it catches up and failover is needed.
func (g *Guardian) DemoteFromActiveSequencer() {
	if !g.state.IsEnclaveActiveSequencer() {
		g.logger.Warn("Guardian demoted when not currently active - this is unexpected")
	}
	g.logger.Info("Guardian demoted from active sequencer")
	g.state.SetActiveSequencer(false)

	// if this has been called the enclave is probably already dead, but we should try to stop it to ensure the enclave
	// is not left running as an active sequencer
	err := g.enclaveClient.Stop()
	if err != nil {
		// log the error at info, this will be common as the enclave is probably already dead
		g.logger.Info("could not stop enclave after demotion", log.ErrKey, err)
	}
}

// HandleBlock is called by the L1 repository when new blocks arrive.
// Note: The L1 processing behaviour has two modes based on the state, either
// - enclave is behind: lookup blocks to feed it 1-by-1 (see `catchupWithL1()`), ignore new live blocks that arrive here
// - enclave is up-to-date: feed it these live blocks as they arrive, no need to lookup blocks
func (g *Guardian) HandleBlock(block *types.Header) {
	if !g.running.Load() {
		return
	}

	g.logger.Debug("Received L1 block", log.BlockHashKey, block.Hash(), log.BlockHeightKey, block.Number)
	// record the newest block we've seen
	g.state.OnReceivedBlock(block.Hash())
	if !g.state.InSyncWithL1() {
		// the enclave is still catching up with the L1 chain, it won't be able to process this new head block yet so return
		return
	}
	_, err := g.submitL1Block(block, true)
	if err != nil {
		g.logger.Warn("failure processing L1 block", log.ErrKey, err)
	}
}

// HandleBatch is called by the L2 repository when a new batch arrives
// Note: this should only be called for validators, sequencers produce their own batches
func (g *Guardian) HandleBatch(batch *common.ExtBatch) {
	if !g.running.Load() {
		return
	}

	g.logger.Debug("Host received L2 batch", log.BatchHashKey, batch.Hash(), log.BatchSeqNoKey, batch.Header.SequencerOrderNo)
	// record the newest batch we've seen
	g.state.OnReceivedBatch(batch.Header.SequencerOrderNo)
	// Sequencer enclaves produce batches, they cannot receive them. Also, enclave will reject new batches if it is not up-to-date
	if g.state.IsEnclaveActiveSequencer() {
		g.logger.Debug("Active sequencer cannot receive batches")
		return
	}
	// expect one after the enclave L2 head, if it's further ahead than that then we don't process it yet
	nextExpectedSeqNo := new(big.Int).Add(g.state.GetEnclaveL2Head(), big.NewInt(1))
	if g.state.GetEnclaveL2Head() != nil && nextExpectedSeqNo.Cmp(batch.Header.SequencerOrderNo) < 0 {
		g.logger.Debug("Enclave is behind, ignoring batch", log.BatchSeqNoKey, batch.Header.SequencerOrderNo)
		return // ignore batches until we're up-to-date
	}
	err := g.submitL2Batch(batch)
	if err != nil {
		g.logger.Error("Error submitting batch to enclave", log.ErrKey, err)
	}
}

func (g *Guardian) HandleTransaction(tx common.EncryptedTx) {
	if !g.running.Load() {
		return
	}

	if g.GetEnclaveState().status == Disconnected ||
		g.GetEnclaveState().status == Unavailable ||
		g.GetEnclaveState().status == AwaitingSecret {
		g.logger.Info("Enclave is not ready yet, dropping transaction.")
		return // ignore transactions when enclave unavailable
	}
	resp, sysError := g.enclaveClient.EncryptedRPC(context.Background(), common.EncryptedRequest(tx))
	if sysError != nil {
		g.logger.Warn("could not submit transaction due to sysError", log.ErrKey, sysError)
		return
	}
	if resp.Error() != nil {
		g.logger.Trace("could not submit transaction", log.ErrKey, resp.Error())
	}
}

// mainLoop runs until the enclave guardian is stopped. It checks the state of the enclave and takes action as
// required to improve the state (e.g. provide a secret, catch up with L1, etc.)
func (g *Guardian) mainLoop() {
	g.logger.Debug("starting guardian main loop")
	unavailableCounter := 0
	for g.running.Load() {
		// check enclave status on every loop (this will happen whenever we hit an error while trying to resolve a state,
		// or after the monitoring interval if we are healthy)
		g.checkEnclaveStatus()
		g.logger.Trace("mainLoop - enclave status", "status", g.state.GetStatus())
		status := g.state.GetStatus()
		if status != Disconnected && status != Unavailable {
			// if we are not disconnected or unavailable, we reset the counter
			unavailableCounter = 0
		}
		switch status {
		case Disconnected, Unavailable:
			// nothing to do, we are waiting for the enclave to be available
			time.Sleep(_retryInterval)
			unavailableCounter++
			if unavailableCounter%50 == 0 { // log every 5 seconds on 100ms retry interval
				g.logger.Error("Enclave is unavailable, retrying connection continuously.")
			}
		case AwaitingSecret:
			err := g.provideSecret()
			if err != nil {
				g.logger.Warn("could not provide secret to enclave", log.ErrKey, err)
				time.Sleep(_retryInterval)
			}
		case L1Catchup:
			// catchUpWithL1 will feed blocks 1-by-1 to the enclave until we are up-to-date, we hit an error or the guardian is stopped
			err := g.catchupWithL1()
			if err != nil {
				g.logger.Warn("could not catch up with L1", log.ErrKey, err)
				time.Sleep(_retryInterval)
			}
		case L2Catchup:
			// catchUpWithL2 will feed batches 1-by-1 to the enclave until we are up-to-date, we hit an error or the guardian is stopped
			err := g.catchupWithL2()
			if err != nil {
				g.logger.Warn("could not catch up with L2", log.ErrKey, err)
				time.Sleep(_retryInterval)
			}
		case Live:
			// we're healthy: loop back to enclave status again after long monitoring interval
			select {
			case <-time.After(_monitoringInterval):
				// loop back to check status
			case <-g.hostInterrupter.Done():
				// stop sleeping, we've been interrupted by the host stopping
			}
		}
	}
	g.logger.Debug("stopping guardian main loop")
}

func (g *Guardian) checkEnclaveStatus() {
	s, err := g.enclaveClient.Status(context.Background())
	if err != nil {
		g.logger.Trace("Could not get enclave status", log.ErrKey, err)
		// we record this as a disconnection, we can't get any more info from the enclave about status currently
		g.state.OnDisconnected()
		return
	}
	g.state.OnEnclaveStatus(s)
}

// This method implements the procedure by which a node obtains the secret
func (g *Guardian) provideSecret() error {
	if g.hostData.IsGenesis {
		// instead of requesting a secret, we generate one and broadcast it
		return g.generateAndBroadcastSecret()
	}
	att, err := g.enclaveClient.Attestation(context.Background())
	if err != nil {
		return fmt.Errorf("could not retrieve attestation from enclave. Cause: %w", err)
	}
	if att.EnclaveID != *g.enclaveID {
		return fmt.Errorf("enclave has ID %s, but it has produced an attestation using ID %s", g.enclaveID.Hex(), att.EnclaveID.Hex())
	}

	g.logger.Info("Requesting secret.")
	// returns the L1 block when the request was published, any response will be after that block
	awaitFromBlock, err := g.sl.L1Publisher().RequestSecret(att)
	if err != nil {
		return errors.Wrap(err, "could not request secret from L1")
	}

	// keep checking L1 blocks until we find a secret response for our request or timeout
	err = retry.Do(func() error {
		nextBlock, _, err := g.sl.L1Data().FetchNextBlock(awaitFromBlock)
		if err != nil {
			return fmt.Errorf("next block after block=%s not found - %w", awaitFromBlock, err)
		}

		processedData, err := g.sl.L1Data().GetTenRelevantTransactions(nextBlock)
		if err != nil {
			return fmt.Errorf("failed to extract Ten transactions from block=%s", nextBlock.Hash())
		}
		secretResponseEvents := processedData.GetEvents(common.SecretResponseTx)
		secretRespTxs := g.sl.L1Publisher().FindSecretResponseTx(secretResponseEvents)
		for _, scrt := range secretRespTxs {
			if scrt.RequesterID.Hex() == g.enclaveID.Hex() {
				err = g.enclaveClient.InitEnclave(context.Background(), scrt.Secret)
				if err != nil {
					g.logger.Error("Could not initialize enclave with received secret response", log.ErrKey, err)
					continue // try the next secret response in the block if there are more
				}
				return nil // successfully initialized enclave with secret, break out of retry loop function
			}
		}
		awaitFromBlock = nextBlock.Hash()
		return errors.New("no valid secret received in block")
	}, retry.NewTimeoutStrategy(_maxWaitForSecretResponse, 500*time.Millisecond))
	if err != nil {
		// something went wrong, check the enclave status in case it is an enclave problem and let the main loop try again when appropriate
		return errors.Wrap(err, "no valid secret received for enclave")
	}

	g.logger.Info("Secret received")
	g.state.OnSecretProvided()

	return nil
}

func (g *Guardian) generateAndBroadcastSecret() error {
	g.logger.Info("Node is genesis node. Publishing secret to L1 enclave registry contract.")
	// Create the shared secret and submit it to the management contract for storage
	attestation, err := g.enclaveClient.Attestation(context.Background())
	if err != nil {
		return fmt.Errorf("could not retrieve attestation from enclave. Cause: %w", err)
	}
	if attestation.EnclaveID != *g.enclaveID {
		return fmt.Errorf("genesis enclave has ID %s, but its enclave produced an attestation using ID %s", g.enclaveID.Hex(), attestation.EnclaveID.Hex())
	}

	secret, err := g.enclaveClient.GenerateSecret(context.Background())
	if err != nil {
		return fmt.Errorf("could not generate secret. Cause: %w", err)
	}

	err = g.sl.L1Publisher().InitializeSecret(attestation, secret)
	if err != nil {
		return errors.Wrap(err, "failed to publish generated enclave secret")
	}
	g.logger.Info("Node is genesis node. Secret generation was published to L1.")
	g.state.OnSecretProvided()
	return nil
}

func (g *Guardian) catchupWithL1() error {
	// while we are behind the L1 head and still running, fetch and submit L1 blocks
	for g.running.Load() && g.state.GetStatus() == L1Catchup {
		// generally we will be feeding the block after the enclave's current head
		enclaveHead := g.state.GetEnclaveL1Head()
		if enclaveHead == gethutil.EmptyHash {
			// but if enclave has no current head, then we use the configured hash to find the first block to feed
			enclaveHead = g.l1StartHash
		}

		l1Block, isLatest, err := g.sl.L1Data().FetchNextBlock(enclaveHead)
		if err != nil {
			if errors.Is(err, gethutil.ErrAncestorNotFound) {
				g.logger.Error("should not happen. Chain fork cannot be calculated because there are missing blocks")
			}
			if errors.Is(err, l1.ErrNoNextBlock) {
				if g.state.GetHostL1Head() == gethutil.EmptyHash {
					// this is usually temporary after a restart until new heads stream starts receiving blocks
					return fmt.Errorf("host not received any new L1 blocks since startup, cannot catch up with L1")
				}
				return nil // we are up-to-date
			}
			return errors.Wrap(err, "could not fetch next L1 block")
		}
		_, err = g.submitL1Block(l1Block, isLatest)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *Guardian) catchupWithL2() error {
	// while we are behind the L2 head and still running:
	for g.running.Load() && g.state.GetStatus() == L2Catchup {
		if g.hostData.IsSequencer && g.state.IsEnclaveActiveSequencer() {
			return errors.New("l2 catchup is not supported for active sequencer")
		}
		// request the next batch by sequence number (based on what the enclave has been fed so far)
		prevHead := g.state.GetEnclaveL2Head()
		nextHead := prevHead.Add(prevHead, big.NewInt(1))

		g.logger.Trace("fetching next batch", log.BatchSeqNoKey, nextHead)
		batch, err := g.sl.L2Repo().FetchBatchBySeqNo(context.Background(), nextHead)
		if err != nil {
			return errors.Wrap(err, "could not fetch next L2 batch")
		}

		err = g.submitL2Batch(batch)
		if err != nil {
			return err
		}
	}
	return nil
}

// returns false if the block was not processed
// todo - @matt - think about removing the TryLock
func (g *Guardian) submitL1Block(block *types.Header, isLatest bool) (bool, error) {
	g.logger.Trace("submitting L1 block", log.BlockHashKey, block.Hash(), log.BlockHeightKey, block.Number)
	if !g.submitDataLock.TryLock() {
		g.logger.Debug("Unable to submit block, enclave is busy processing data")
		return false, nil
	}
	processedData, err := g.sl.L1Data().GetTenRelevantTransactions(block)
	if err != nil {
		g.submitDataLock.Unlock() // lock must be released before returning
		return false, fmt.Errorf("could not extract ten transaction for block=%s - %w", block.Hash(), err)
	}

	rollupTxs := g.getRollupTxs(*processedData)

	resp, err := g.enclaveClient.SubmitL1Block(context.Background(), processedData)

	g.submitDataLock.Unlock() // lock is only guarding the enclave call, so we can release it now
	if resp.RejectError != nil {
		if strings.Contains(resp.RejectError.Error(), errutil.ErrBlockAlreadyProcessed.Error()) {
			// we have already processed this block, let's try the next canonical block
			// this is most common when we are returning to a previous fork and the enclave has already seen some of the blocks on it
			// note: logging this because we don't expect it to happen often and would like visibility on that.
			g.logger.Info("L1 block already processed by enclave, trying the next block", "block", block.Hash())
			nextHeight := big.NewInt(0).Add(block.Number, big.NewInt(1))
			nextCanonicalBlock, err := g.sl.L1Data().FetchBlockByHeight(nextHeight)
			if err != nil {
				return false, fmt.Errorf("failed to fetch next block after forking block=%s: %w", block.Hash(), err)
			}
			return g.submitL1Block(nextCanonicalBlock, isLatest)
		}
		// something went wrong, return error and let the main loop check status and try again when appropriate
		return false, errors.Wrap(err, "could not submit L1 block to enclave")
	}
	// successfully processed block, update the state
	g.state.OnProcessedBlock(block.Hash())
	g.processL1BlockTransactions(block, resp.RollupMetadata, rollupTxs, g.shouldSyncContracts(*processedData), g.shouldSyncAdditionalContracts(*processedData))

	// todo: make sure this doesn't respond to old requests (once we have a proper protocol for that)
	err = g.publishSharedSecretResponses(resp.ProducedSecretResponses)
	if err != nil {
		g.logger.Error("Failed to publish response to secret request", log.ErrKey, err)
	}
	return true, nil
}

func (g *Guardian) processL1BlockTransactions(block *types.Header, metadatas []common.ExtRollupMetadata, rollupTxs []*common.L1RollupTx, syncContracts bool, syncAdditionalContracts bool) {
	for idx, rollup := range rollupTxs {
		r, err := common.DecodeRollup(rollup.Rollup)
		if err != nil {
			g.logger.Error("Could not decode rollup.", log.ErrKey, err)
		}

		metaData, err := g.enclaveClient.GetRollupData(context.Background(), r.Header.Hash())
		if err != nil {
			g.logger.Error("Could not fetch rollup metadata from enclave.", log.RollupHashKey, r.Header.Hash(), log.ErrKey, err)
		} else {
			// TODO - This is a temporary fix, arrays should always match in practice...
			extMetadata := common.ExtRollupMetadata{}
			if len(metadatas) > idx {
				extMetadata = metadatas[idx]
			}
			err = g.storage.AddRollup(r, &extMetadata, metaData, block)
		}
		if err != nil {
			if errors.Is(err, errutil.ErrAlreadyExists) {
				g.logger.Info("Rollup already stored", log.RollupHashKey, r.Hash())
			} else {
				g.logger.Error("Could not store rollup.", log.ErrKey, err)
			}
		}
	}

	if syncContracts || syncAdditionalContracts {
		go func() {
			err := g.sl.L1Publisher().ResyncImportantContracts()
			if err != nil {
				g.logger.Error("Could not resync important contracts", log.ErrKey, err)
			}
		}()
	}
}

func (g *Guardian) publishSharedSecretResponses(scrtResponses []*common.ProducedSecretResponse) error {
	for _, scrtResponse := range scrtResponses {
		// todo (#1624) - implement proper protocol so only one host responds to this secret requests initially
		// 	for now we just have the genesis host respond until protocol implemented
		if !g.hostData.IsSequencer {
			g.logger.Trace("Not genesis node, not publishing response to secret request.",
				"requester", scrtResponse.RequesterID)
			return nil
		}

		err := g.sl.L1Publisher().PublishSecretResponse(scrtResponse)
		if err != nil {
			return errors.Wrap(err, "could not publish secret response")
		}
	}
	return nil
}

func (g *Guardian) submitL2Batch(batch *common.ExtBatch) error {
	g.submitDataLock.Lock()
	err := g.enclaveClient.SubmitBatch(context.Background(), batch)
	g.submitDataLock.Unlock()
	if err != nil {
		// something went wrong, return error and let the main loop check status and try again when appropriate
		return errors.Wrap(err, "could not submit L2 batch to enclave")
	}
	// successfully processed batch, update the state
	g.state.OnProcessedBatch(batch.Header.SequencerOrderNo)
	return nil
}

func (g *Guardian) ProduceBatch() error {
	// todo @matt remove skipIfEmpty flag and get rid of relevant config
	return g.enclaveClient.CreateBatch(context.Background(), false)
}

func (g *Guardian) streamEnclaveData() {
	defer g.logger.Info("Stopping enclave data stream")
	g.logger.Info("Starting L2 update stream from enclave")

	streamChan, stopStream := g.enclaveClient.StreamL2Updates()
	var lastBatch *common.ExtBatch
	for {
		select {
		case resp, ok := <-streamChan:
			if !ok {
				stopStream()
				g.logger.Warn("Batch streaming failed. Reconnecting after 3 seconds")
				time.Sleep(3 * time.Second)
				streamChan, stopStream = g.enclaveClient.StreamL2Updates()

				continue
			}

			if resp.Batch != nil { //nolint:nestif
				lastBatch = resp.Batch
				g.logger.Trace("Received batch from stream", log.BatchHashKey, lastBatch.Hash())
				err := g.sl.L2Repo().AddBatch(resp.Batch)
				if err != nil && !errors.Is(err, errutil.ErrAlreadyExists) {
					g.logger.Crit("failed to add batch to L2 repo", log.BatchHashKey, resp.Batch.Hash(), log.ErrKey, err)
				}

				if g.state.IsEnclaveActiveSequencer() { // active sequencer enclave should broadcast the batch to peers
					g.logger.Info("Batch produced. Sending to peers..", log.BatchHeightKey, resp.Batch.Header.Number, log.BatchHashKey, resp.Batch.Hash())

					err = g.sl.P2P().BroadcastBatches([]*common.ExtBatch{resp.Batch})
					if err != nil {
						g.logger.Error("Failed to broadcast batch", log.BatchHashKey, resp.Batch.Hash(), log.ErrKey, err)
					}
				} else {
					g.logger.Debug("Received validated batch from enclave", log.BatchSeqNoKey, resp.Batch.Header.SequencerOrderNo, log.BatchHashKey, resp.Batch.Hash())
				}
				// Notify the L2 repo that an enclave has validated a batch, so it can update its validated head and notify subscribers
				g.sl.L2Repo().NotifyNewValidatedHead(resp.Batch)
				g.state.OnProcessedBatch(resp.Batch.Header.SequencerOrderNo)
			}

			if resp.Logs != nil {
				g.sl.LogSubs().SendLogsToSubscribers(&resp.Logs)
			}

		case <-g.hostInterrupter.Done():
			// interrupted by host stopping
			return
		}
	}
}

func (g *Guardian) getRollupTxs(processed common.ProcessedL1Data) []*common.L1RollupTx {
	rollupTxs := make([]*common.L1RollupTx, 0)

	for _, txData := range processed.GetEvents(common.RollupTx) {
		encodedRlp, err := ethadapter.DecodeBlobs(txData.BlobsWithSignature.ToBlobs())
		if err != nil {
			g.logger.Crit("could not decode blobs.", log.ErrKey, err)
			continue
		}

		rlp := &common.L1RollupTx{
			Rollup: encodedRlp,
		}
		rollupTxs = append(rollupTxs, rlp)
	}

	return rollupTxs
}

func (g *Guardian) shouldSyncContracts(processed common.ProcessedL1Data) bool {
	return len(processed.GetEvents(common.NetworkContractAddressAddedTx)) > 0
}

func (g *Guardian) shouldSyncAdditionalContracts(processed common.ProcessedL1Data) bool {
	return len(processed.GetEvents(common.AdditionalContractAddressAddedTx)) > 0
}
