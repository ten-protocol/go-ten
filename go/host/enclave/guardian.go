package enclave

import (
	"fmt"
	"math/big"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/obscuronet/go-obscuro/go/host/events"

	"github.com/obscuronet/go-obscuro/go/common/gethutil"

	"github.com/kamilsk/breaker"

	"github.com/ethereum/go-ethereum/core/types"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/errutil"
	"github.com/obscuronet/go-obscuro/go/common/host"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/common/retry"
	"github.com/obscuronet/go-obscuro/go/config"
	"github.com/obscuronet/go-obscuro/go/host/db"
	"github.com/obscuronet/go-obscuro/go/host/l1"
	"github.com/pkg/errors"
)

const (
	// time between loops on mainLoop, will be retry time if things are failing
	_retryInterval = 100 * time.Millisecond

	// when enclave is healthy this is the time before we call its status (can be slow, is just a sanity check)
	_monitoringInterval = 1 * time.Second

	// when we have submitted request to L1 for the secret, how long do we wait for an answer before we retry
	_maxWaitForSecretResponse = 2 * time.Minute
)

// Guardian is a host service which monitors an enclave, it's responsibilities include:
// - monitor the enclave state and feed it the data it needs
// - if it is an active sequencer then the guardian will trigger batch/rollup creation
// - guardian provides access to the enclave data and reports the enclave status for other services - acting as a gatekeeper
type Guardian struct {
	hostData      host.Identity
	state         *StateTracker // state machine that tracks our view of the enclave's state
	enclaveClient common.Enclave
	db            *db.DB

	submitDataLock sync.Mutex // we only submit one block, batch or transaction to enclave at a time

	batchInterval  time.Duration
	rollupInterval time.Duration

	running         atomic.Bool
	hostInterrupter breaker.Interface // host hostInterrupter so we can stop quickly

	logger          gethlog.Logger
	p2p             host.P2P
	l1Publisher     host.L1Publisher
	l2Repo          host.L2BatchRepository
	logEventManager host.LogSubscriptionManager
	l1Repo          host.L1BlockRepository
	isSequencer     bool
}

func NewGuardian(
	cfg *config.HostConfig,
	hostData host.Identity,
	p2p host.P2P,
	l1Publisher host.L1Publisher,
	l1repo host.L1BlockRepository,
	l2Repo host.L2BatchRepository,
	enclaveClient common.Enclave,
	db *db.DB,
	interrupter breaker.Interface,
	logger gethlog.Logger,
) *Guardian {
	return &Guardian{
		hostData:        hostData,
		state:           NewStateTracker(logger),
		enclaveClient:   enclaveClient,
		batchInterval:   cfg.BatchInterval,
		rollupInterval:  cfg.RollupInterval,
		db:              db,
		hostInterrupter: interrupter,
		logger:          logger,
		p2p:             p2p,
		l1Publisher:     l1Publisher,
		l1Repo:          l1repo,
		l2Repo:          l2Repo,
		logEventManager: events.NewSubscriptionManager(),
		isSequencer:     cfg.NodeType == common.Sequencer,
	}
}

func (g *Guardian) Start() error {
	g.running.Store(true)
	go g.mainLoop()
	if g.hostData.IsSequencer {
		// if we are a sequencer then we need to start the periodic batch/rollup production
		// Note: after HA work this will need additional check that we are the **active** sequencer enclave
		go g.periodicBatchProduction()
		go g.periodicRollupProduction()
	}

	// subscribe for L1 and P2P data
	g.p2p.SubscribeForTx(g)
	g.l1Repo.Subscribe(g)
	g.l2Repo.Subscribe(g)

	// start streaming data from the enclave
	go g.streamEnclaveData()

	return nil
}

func (g *Guardian) Stop() error {
	g.running.Store(false)

	err := g.enclaveClient.Stop()
	if err != nil {
		g.logger.Warn("error stopping enclave", log.ErrKey, err)
	}

	err = g.enclaveClient.StopClient()
	if err != nil {
		g.logger.Warn("error stopping enclave client", log.ErrKey, err)
	}

	return nil
}

func (g *Guardian) HealthStatus() host.HealthStatus {
	// todo (@matt) do proper health status based on enclave state
	errMsg := ""
	if !g.running.Load() {
		errMsg = "not running"
	}
	return &host.BasicErrHealthStatus{ErrMsg: errMsg}
}

func (g *Guardian) GetEnclaveState() *StateTracker {
	return g.state
}

// GetEnclaveClient returns the enclave client for use by other services
// todo (@matt) avoid exposing client directly and return errors if enclave is not ready for requests
func (g *Guardian) GetEnclaveClient() common.Enclave {
	return g.enclaveClient
}

// HandleBlock is called by the L1 repository when new blocks arrive.
// Note: The L1 processing behaviour has two modes based on the state, either
// - enclave is behind: lookup blocks to feed it 1-by-1 (see `catchupWithL1()`), ignore new live blocks that arrive here
// - enclave is up-to-date: feed it these live blocks as they arrive, no need to lookup blocks
func (g *Guardian) HandleBlock(block *types.Block) {
	g.logger.Debug("Received L1 block", log.BlockHashKey, block.Hash(), log.BlockHeightKey, block.Number())
	// record the newest block we've seen
	g.state.OnReceivedBlock(block.Hash())
	if !g.state.InSyncWithL1() {
		// the enclave is still catching up with the L1 chain, it won't be able to process this new head block yet so return
		return
	}
	err := g.submitL1Block(block, true)
	if err != nil {
		g.logger.Warn("failure processing L1 block", log.ErrKey, err)
	}
}

// HandleBatch is called by the L2 repository when a new batch arrives
// Note: this should only be called for validators, sequencers produce their own batches
func (g *Guardian) HandleBatch(batch *common.ExtBatch) {
	if g.hostData.IsSequencer {
		g.logger.Error("repo received batch but we are a sequencer, ignoring")
		return
	}
	g.logger.Debug("Received L2 block", log.BatchHashKey, batch.Hash(), log.BatchSeqNoKey, batch.Header.SequencerOrderNo)
	// record the newest batch we've seen
	g.state.OnReceivedBatch(batch.Header.SequencerOrderNo)
	if !g.state.IsUpToDate() {
		return // ignore batches until we're up-to-date
	}
	err := g.submitL2Batch(batch)
	if err != nil {
		g.logger.Warn("error submitting batch to enclave", log.ErrKey, err)
	}
}

func (g *Guardian) HandleTransaction(tx common.EncryptedTx) {
	resp, sysError := g.enclaveClient.SubmitTx(tx)
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
	for g.running.Load() {
		// check enclave status on every loop (this will happen whenever we hit an error while trying to resolve a state,
		// or after the monitoring interval if we are healthy)
		g.checkEnclaveStatus()
		g.logger.Trace("mainLoop - enclave status", "status", g.state.GetStatus())
		switch g.state.GetStatus() {
		case Disconnected, Unavailable:
			// nothing to do, we are waiting for the enclave to be available
			time.Sleep(_retryInterval)
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
	s, err := g.enclaveClient.Status()
	if err != nil {
		g.logger.Error("could not get enclave status", log.ErrKey, err)
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
	g.logger.Info("Requesting secret.")
	att, err := g.enclaveClient.Attestation()
	if err != nil {
		return fmt.Errorf("could not retrieve attestation from enclave. Cause: %w", err)
	}
	if att.Owner != g.hostData.ID {
		return fmt.Errorf("host has ID %s, but its enclave produced an attestation using ID %s", g.hostData.ID.Hex(), att.Owner.Hex())
	}

	g.logger.Info("Requesting secret.")
	// returns the L1 block when the request was published, any response will be after that block
	awaitFromBlock, err := g.l1Publisher.RequestSecret(att)
	if err != nil {
		return errors.Wrap(err, "could not request secret from L1")
	}

	// keep checking L1 blocks until we find a secret response for our request or timeout
	err = retry.Do(func() error {
		nextBlock, _, err := g.l1Repo.FetchNextBlock(awaitFromBlock)
		if err != nil {
			return fmt.Errorf("next block after block=%s not found - %w", awaitFromBlock, err)
		}
		secretRespTxs := g.l1Publisher.ExtractSecretResponses(nextBlock)
		if err != nil {
			return fmt.Errorf("could not extract secret responses from block=%s - %w", nextBlock.Hash(), err)
		}
		for _, scrt := range secretRespTxs {
			if scrt.RequesterID.Hex() == g.hostData.ID.Hex() {
				err = g.enclaveClient.InitEnclave(scrt.Secret)
				if err != nil {
					g.logger.Warn("could not initialize enclave with received secret response", log.ErrKey, err)
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

	// we're now ready to catch up with network, sync peer list
	go g.p2p.RefreshPeerList()
	return nil
}

func (g *Guardian) generateAndBroadcastSecret() error {
	g.logger.Info("Node is genesis node. Broadcasting secret.")
	// Create the shared secret and submit it to the management contract for storage
	attestation, err := g.enclaveClient.Attestation()
	if err != nil {
		return fmt.Errorf("could not retrieve attestation from enclave. Cause: %w", err)
	}
	if attestation.Owner != g.hostData.ID {
		return fmt.Errorf("genesis node has ID %s, but its enclave produced an attestation using ID %s", g.hostData.ID.Hex(), attestation.Owner.Hex())
	}

	secret, err := g.enclaveClient.GenerateSecret()
	if err != nil {
		return fmt.Errorf("could not generate secret. Cause: %w", err)
	}

	err = g.l1Publisher.InitializeSecret(attestation, secret)
	if err != nil {
		return errors.Wrap(err, "failed to initialise enclave secret")
	}
	g.logger.Info("Node is genesis node. Secret was broadcast.")
	g.state.OnSecretProvided()
	return nil
}

func (g *Guardian) catchupWithL1() error {
	// while we are behind the L1 head and still running, fetch and submit L1 blocks
	for g.running.Load() && g.state.GetStatus() == L1Catchup {
		l1Block, isLatest, err := g.l1Repo.FetchNextBlock(g.state.GetEnclaveL1Head())
		if err != nil {
			if errors.Is(err, l1.ErrNoNextBlock) {
				if g.state.hostL1Head == gethutil.EmptyHash {
					return fmt.Errorf("no L1 blocks found in repository")
				}
				return nil // we are up-to-date
			}
			return errors.Wrap(err, "could not fetch next L1 block")
		}
		err = g.submitL1Block(l1Block, isLatest)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *Guardian) catchupWithL2() error {
	// while we are behind the L2 head and still running:
	for g.running.Load() && g.state.GetStatus() == L2Catchup {
		if g.hostData.IsSequencer {
			return errors.New("l2 catchup is not supported for sequencer")
		}
		// request the next batch by sequence number (based on what the enclave has been fed so far)
		prevHead := g.state.GetEnclaveL2Head()
		nextHead := prevHead.Add(prevHead, big.NewInt(1))

		g.logger.Trace("fetching next batch", log.BatchSeqNoKey, nextHead)
		batch, err := g.l2Repo.FetchBatchBySeqNo(nextHead)
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

// LookupBatchBySeqNo is used to fetch batch data from the enclave - it is only used as a fallback for the sequencer
// host if it's missing a batch (other host services should use L2Repo to fetch batch data)
func (g *Guardian) LookupBatchBySeqNo(seqNo *big.Int) (*common.ExtBatch, error) {
	state := g.GetEnclaveState()
	if state.GetEnclaveL2Head().Cmp(seqNo) < 0 {
		return nil, errutil.ErrNotFound
	}
	client := g.GetEnclaveClient()
	return client.GetBatchBySeqNo(seqNo.Uint64())
}

func (g *Guardian) submitL1Block(block *common.L1Block, isLatest bool) error {
	g.logger.Trace("submitting L1 block", log.BlockHashKey, block.Hash(), log.BlockHeightKey, block.Number())
	receipts := g.l1Repo.FetchReceipts(block)
	if !g.submitDataLock.TryLock() {
		// we are already submitting a block, and we don't want to leak goroutines, we wil catch up with the block later
		return errors.New("unable to submit block, already submitting another block")
	}
	resp, err := g.enclaveClient.SubmitL1Block(*block, receipts, isLatest)
	g.submitDataLock.Unlock()
	if err != nil {
		if strings.Contains(err.Error(), errutil.ErrBlockAlreadyProcessed.Error()) {
			// we have already processed this block, let's try the next canonical block
			// this is most common when we are returning to a previous fork and the enclave has already seen some of the blocks on it
			// note: logging this because we don't expect it to happen often and would like visibility on that.
			g.logger.Info("L1 block already processed by enclave, trying the next block", "block", block.Hash())
			nextHeight := big.NewInt(0).Add(block.Number(), big.NewInt(1))
			nextCanonicalBlock, err := g.l1Repo.FetchBlockByHeight(nextHeight)
			if err != nil {
				return fmt.Errorf("failed to fetch next block after forking block=%s: %w", block.Hash(), err)
			}
			return g.submitL1Block(nextCanonicalBlock, isLatest)
		}
		// something went wrong, return error and let the main loop check status and try again when appropriate
		return errors.Wrap(err, "could not submit L1 block to enclave")
	}
	// successfully processed block, update the state
	g.state.OnProcessedBlock(block.Hash())
	g.processL1BlockTransactions(block)

	// todo (@matt) this should not be here, it is only used by the RPC API server for batch data which will eventually just use L1 repo
	err = g.db.AddBlockHeader(block.Header())
	if err != nil {
		return fmt.Errorf("submitted block to enclave but could not store the block processing result. Cause: %w", err)
	}

	// todo: make sure this doesn't respond to old requests (once we have a proper protocol for that)
	err = g.publishSharedSecretResponses(resp.ProducedSecretResponses)
	if err != nil {
		g.logger.Error("failed to publish response to secret request", log.ErrKey, err)
	}
	return nil
}

func (g *Guardian) processL1BlockTransactions(block *common.L1Block) {
	// if there are any secret responses in the block we should refresh our P2P list to re-sync with the network
	respTxs := g.l1Publisher.ExtractSecretResponses(block)
	if len(respTxs) > 0 {
		// new peers may have been granted access to the network, notify p2p service to refresh its peer list
		go g.p2p.RefreshPeerList()
	}

	rollupTxs := g.l1Publisher.ExtractRollupTxs(block)
	for _, rollup := range rollupTxs {
		r, err := common.DecodeRollup(rollup.Rollup)
		if err != nil {
			g.logger.Error("could not decode rollup.", log.ErrKey, err)
		}
		err = g.db.AddRollupHeader(r)
		if err != nil {
			g.logger.Error("could not store rollup.", log.ErrKey, err)
		}
	}
}

func (g *Guardian) publishSharedSecretResponses(scrtResponses []*common.ProducedSecretResponse) error {
	for _, scrtResponse := range scrtResponses {
		// todo (#1624) - implement proper protocol so only one host responds to this secret requests initially
		// 	for now we just have the genesis host respond until protocol implemented
		if !g.hostData.IsGenesis {
			g.logger.Trace("Not genesis node, not publishing response to secret request.",
				"requester", scrtResponse.RequesterID)
			return nil
		}

		err := g.l1Publisher.PublishSecretResponse(scrtResponse)
		if err != nil {
			return errors.Wrap(err, "could not publish secret response")
		}
	}
	return nil
}

func (g *Guardian) submitL2Batch(batch *common.ExtBatch) error {
	g.submitDataLock.Lock()
	err := g.enclaveClient.SubmitBatch(batch)
	g.submitDataLock.Unlock()
	if err != nil {
		// something went wrong, return error and let the main loop check status and try again when appropriate
		return errors.Wrap(err, "could not submit L2 batch to enclave")
	}
	// successfully processed batch, update the state
	g.state.OnProcessedBatch(batch.Header.SequencerOrderNo)
	return nil
}

func (g *Guardian) periodicBatchProduction() {
	defer g.logger.Info("Stopping batch production")

	interval := g.batchInterval
	if interval == 0 {
		interval = 1 * time.Second
	}
	batchProdTicker := time.NewTicker(interval)
	// attempt to produce rollup every time the timer ticks until we are stopped/interrupted
	for {
		if !g.running.Load() {
			batchProdTicker.Stop()
			return // stop periodic rollup production
		}
		select {
		case <-batchProdTicker.C:
			if !g.state.InSyncWithL1() {
				// if we're behind the L1, we don't want to produce batches
				g.logger.Debug("skipping batch production because L1 is not up to date")
				continue
			}
			g.logger.Debug("create batch")
			err := g.enclaveClient.CreateBatch()
			if err != nil {
				g.logger.Error("unable to produce batch", log.ErrKey, err)
			}
		case <-g.hostInterrupter.Done():
			// interrupted - end periodic process
			batchProdTicker.Stop()
			return
		}
	}
}

func (g *Guardian) periodicRollupProduction() {
	defer g.logger.Info("Stopping rollup production")

	interval := g.rollupInterval
	if interval == 0 {
		interval = 3 * time.Second
	}
	rollupTicker := time.NewTicker(interval)
	// attempt to produce rollup every time the timer ticks until we are stopped/interrupted
	for {
		if !g.running.Load() {
			rollupTicker.Stop()
			return // stop periodic rollup production
		}
		select {
		case <-rollupTicker.C:
			if !g.state.IsUpToDate() {
				// if we're behind the L1, we don't want to produce rollups
				g.logger.Debug("skipping rollup production because L1 is not up to date", "state", g.state)
				continue
			}
			lastBatchNo, err := g.l1Publisher.FetchLatestSeqNo()
			if err != nil {
				g.logger.Warn("encountered error while trying to retrieve latest sequence number", log.ErrKey, err)
				continue
			}
			producedRollup, err := g.enclaveClient.CreateRollup(lastBatchNo.Uint64())
			if err != nil {
				g.logger.Error("unable to produce rollup", log.ErrKey, err)
			} else {
				g.l1Publisher.PublishRollup(producedRollup)
			}
		case <-g.hostInterrupter.Done():
			// interrupted - end periodic process
			rollupTicker.Stop()
			return
		}
	}
}

func (g *Guardian) streamEnclaveData() {
	defer g.logger.Info("Stopping enclave data stream")
	g.logger.Info("Starting L2 update stream from enclave")

	streamChan, stop := g.enclaveClient.StreamL2Updates()
	var lastBatch *common.ExtBatch
	for {
		select {
		case resp, ok := <-streamChan:
			if !ok {
				stop()
				g.logger.Warn("Batch streaming failed. Reconnecting after 3 seconds")
				time.Sleep(3 * time.Second)
				streamChan, stop = g.enclaveClient.StreamL2Updates()

				continue
			}

			if resp.Batch != nil {
				lastBatch = resp.Batch
				g.logger.Trace("Received batch from stream", log.BatchHashKey, lastBatch.Hash())
				err := g.l2Repo.AddBatch(resp.Batch)
				if err != nil && !errors.Is(err, errutil.ErrAlreadyExists) {
					// todo (@matt) this is a catastrophic scenario, the host may never get that batch - handle this
					g.logger.Crit("failed to add batch to L2 repo", log.BatchHashKey, resp.Batch.Hash(), log.ErrKey, err)
				}

				if g.hostData.IsSequencer { // if we are the sequencer we need to broadcast this new batch to the network
					g.logger.Info("Batch produced", log.BatchHeightKey, resp.Batch.Header.Number, log.BatchHashKey, resp.Batch.Hash())

					err = g.p2p.BroadcastBatches([]*common.ExtBatch{resp.Batch})
					if err != nil {
						g.logger.Error("failed to broadcast batch", log.BatchHashKey, resp.Batch.Hash(), log.ErrKey, err)
					}
				}
				g.logger.Info("Batch streamed", log.BatchHeightKey, resp.Batch.Header.Number, log.BatchHashKey, resp.Batch.Hash())
				g.state.OnProcessedBatch(resp.Batch.Header.SequencerOrderNo)
			}

			if resp.Logs != nil {
				g.logEventManager.SendLogsToSubscribers(&resp.Logs)
			}

		case <-time.After(1 * time.Second):
			if !g.running.Load() {
				// guardian service is stopped
				return
			}

		case <-g.hostInterrupter.Done():
			// interrupted - end periodic process
			return
		}
	}
}
