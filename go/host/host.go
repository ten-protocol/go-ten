package host

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/obscuronet/go-obscuro/go/host/l2"

	"github.com/obscuronet/go-obscuro/go/host/l1"
	"github.com/pkg/errors"

	"github.com/obscuronet/go-obscuro/go/host/enclave"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/kamilsk/breaker"
	"github.com/naoina/toml"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/errutil"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/common/measure"
	"github.com/obscuronet/go-obscuro/go/common/profiler"
	"github.com/obscuronet/go-obscuro/go/common/retry"
	"github.com/obscuronet/go-obscuro/go/common/stopcontrol"
	"github.com/obscuronet/go-obscuro/go/config"
	"github.com/obscuronet/go-obscuro/go/ethadapter"
	"github.com/obscuronet/go-obscuro/go/ethadapter/mgmtcontractlib"
	"github.com/obscuronet/go-obscuro/go/host/db"
	"github.com/obscuronet/go-obscuro/go/host/events"
	"github.com/obscuronet/go-obscuro/go/responses"
	"github.com/obscuronet/go-obscuro/go/wallet"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
	gethmetrics "github.com/ethereum/go-ethereum/metrics"
	hostcommon "github.com/obscuronet/go-obscuro/go/common/host"
)

const (
	maxWaitForSecretResponse = 120 * time.Second
)

type P2PHostService interface {
	hostcommon.Service
	hostcommon.P2P
}

// Implementation of host.Host.
type host struct {
	config   *config.HostConfig
	shortID  uint64
	services map[string]hostcommon.Service // host services - registered by name for health checks, start and stop

	enclaveClient     common.Enclave        // For communication with the enclave
	enclaveState      *enclave.StateTracker // StateTracker machine that maintains the current state of enclave (stepping stone to enclave-guardian)
	statusReqInFlight atomic.Bool           // Flag to ensure that only one enclave status request is sent at a time (the requests are triggered in separate threads when unexpected enclave errors occur)

	// control the host lifecycle
	interrupter   breaker.Interface
	shutdownGroup sync.WaitGroup

	// ignore incoming requests
	stopControl *stopcontrol.StopControl

	txP2PCh chan common.EncryptedTx // The channel that new transactions from peers are sent to

	submitDataToEnclaveLock sync.Mutex // host should only submit one block or batch to enclave at a time (probably not needed in enclave-guardian but here for safety for now)

	db *db.DB // Stores the host's publicly-available data

	logEventManager *events.LogEventManager

	logger gethlog.Logger

	metricRegistry gethmetrics.Registry
}

func NewHost(
	config *config.HostConfig,
	p2p P2PHostService,
	ethClient ethadapter.EthClient,
	enclaveClient common.Enclave,
	ethWallet wallet.Wallet,
	mgmtContractLib mgmtcontractlib.MgmtContractLib,
	logger gethlog.Logger,
	regMetrics gethmetrics.Registry,
) hostcommon.Host {
	database, err := db.CreateDBFromConfig(config, regMetrics, logger)
	if err != nil {
		logger.Crit("unable to create database for host", log.ErrKey, err)
	}
	l1Repo := l1.NewL1Repository(ethClient, logger)
	l1StartHash := config.L1StartHash
	if l1StartHash == (gethcommon.Hash{}) {
		startBlock, err := l1Repo.FetchBlockByHeight(0)
		if err != nil {
			logger.Crit("unable to fetch start block so stream from", log.ErrKey, err)
		}
		l1StartHash = startBlock.Hash()
	}
	enclStateTracker := enclave.NewStateTracker(logger)
	enclStateTracker.OnProcessedBlock(l1StartHash) // this makes sure we start streaming from the right block, will be less clunky in the enclave guardian
	hostIdentity := hostcommon.NewIdentity(config.ID, config.P2PPublicAddress)
	host := &host{
		// config
		config:  config,
		shortID: common.ShortAddress(config.ID),

		// services
		services: make(map[string]hostcommon.Service),

		// Communication layers.
		enclaveClient: enclaveClient,
		enclaveState:  enclStateTracker,

		// incoming data
		txP2PCh: make(chan common.EncryptedTx),

		submitDataToEnclaveLock: sync.Mutex{},

		// Initialize the host DB
		db: database,

		logEventManager: events.NewLogEventManager(logger),

		logger:         logger,
		metricRegistry: regMetrics,

		stopControl: stopcontrol.New(),
	}
	l2Repo := l2.NewBatchRepository(config, p2p, database, host.logger)

	host.RegisterService(hostcommon.P2PName, p2p)
	host.RegisterService(hostcommon.L1BlockRepositoryName, l1Repo)
	host.RegisterService(hostcommon.L1PublisherName, l1.NewL1Publisher(hostIdentity, ethWallet, ethClient, mgmtContractLib, l1Repo, logger))
	host.RegisterService(hostcommon.L2BatchRepositoryName, l2Repo)

	var prof *profiler.Profiler
	if config.ProfilerEnabled {
		prof = profiler.NewProfiler(profiler.DefaultHostPort, logger)
		err := prof.Start()
		if err != nil {
			logger.Crit("unable to start the profiler: %s", log.ErrKey, err)
		}
	}

	jsonConfig, _ := json.MarshalIndent(config, "", "  ")
	logger.Info("Host service created with following config:", log.CfgKey, string(jsonConfig))

	return host
}

func (h *host) RegisterService(name string, service hostcommon.Service) {
	if _, ok := h.services[name]; ok {
		h.logger.Crit("service already registered", "name", name)
	}
	h.services[name] = service
}

// Start validates the host config and starts the Host in a go routine - immediately returns after
func (h *host) Start() error {
	if h.stopControl.IsStopping() {
		return responses.ToInternalError(fmt.Errorf("requested Start with the host stopping"))
	}

	h.interrupter = breaker.Multiplex(
		breaker.BreakBySignal(
			os.Kill,
			os.Interrupt,
		),
	)

	// start all registered services
	for name, service := range h.services {
		err := service.Start()
		if err != nil {
			return errors.Wrapf(err, "could not start service=%s", name)
		}
	}

	h.validateConfig()

	tomlConfig, err := toml.Marshal(h.config)
	if err != nil {
		return fmt.Errorf("could not print host config - %w", err)
	}
	h.logger.Info("Host started with following config", log.CfgKey, string(tomlConfig))

	go func() {
		// wait for the Enclave to be available
		enclStatus := h.waitForEnclave()

		// todo (#1474) - the host should only connect to enclaves with the same ID as the host.ID
		if enclStatus.StatusCode == common.AwaitingSecret {
			err = h.requestSecret()
			if err != nil {
				h.logger.Crit("Could not request secret", log.ErrKey, err.Error())
			}
		}

		err := h.refreshP2PPeerList()
		if err != nil {
			h.logger.Warn("unable to sync current p2p peer list on startup", log.ErrKey, err)
		}

		// start the host's main processing loop
		h.startProcessing()
	}()

	return nil
}

// HandleBlock is called by the L1 repository. The host is subscribed to receive new blocks.
func (h *host) HandleBlock(block *types.Block) {
	h.logger.Debug("Received L1 block", log.BlockHashKey, block.Hash(), log.BlockHeightKey, block.Number())
	// record the newest block we've seen
	h.enclaveState.OnReceivedBlock(block.Hash())
	if !h.enclaveState.InSyncWithL1() {
		return // ignore blocks until we're up-to-date
	}
	h.submitDataToEnclaveLock.Lock()
	defer h.submitDataToEnclaveLock.Unlock()
	err := h.processL1Block(block, true)
	if err != nil {
		h.logger.Warn("error processing L1 block", log.ErrKey, err)
	}
}

// HandleBatch is called by the L2 repository when a new batch arrives
func (h *host) HandleBatch(batch *common.ExtBatch) {
	h.logger.Debug("Received L2 block", log.BatchHashKey, batch.Hash(), log.BatchSeqNoKey, batch.Header.SequencerOrderNo)
	// record the newest batch we've seen
	h.enclaveState.OnReceivedBatch(batch.Header.SequencerOrderNo)
	if !h.enclaveState.IsUpToDate() {
		return // ignore batches until we're up-to-date
	}
	h.submitDataToEnclaveLock.Lock()
	defer h.submitDataToEnclaveLock.Unlock()
	err := h.enclaveClient.SubmitBatch(batch)
	if err != nil {
		h.logger.Warn("error submitting batch to enclave", log.ErrKey, err)
	} else {
		h.enclaveState.OnProcessedBatch(batch.Header.SequencerOrderNo)
	}
}

func (h *host) generateAndBroadcastSecret() error {
	h.logger.Info("Node is genesis node. Broadcasting secret.")
	// Create the shared secret and submit it to the management contract for storage
	attestation, err := h.enclaveClient.Attestation()
	if err != nil {
		return fmt.Errorf("could not retrieve attestation from enclave. Cause: %w", err)
	}
	if attestation.Owner != h.config.ID {
		return fmt.Errorf("genesis node has ID %s, but its enclave produced an attestation using ID %s", h.config.ID.Hex(), attestation.Owner.Hex())
	}

	secret, err := h.enclaveClient.GenerateSecret()
	if err != nil {
		return fmt.Errorf("could not generate secret. Cause: %w", err)
	}

	err = h.l1Publisher().InitializeSecret(attestation, secret)
	if err != nil {
		return errors.Wrap(err, "failed to initialise enclave secret")
	}
	h.logger.Info("Node is genesis node. Secret was broadcast.")
	h.enclaveState.OnSecretProvided()
	return nil
}

func (h *host) Config() *config.HostConfig {
	return h.config
}

func (h *host) DB() *db.DB {
	return h.db
}

func (h *host) EnclaveClient() common.Enclave {
	return h.enclaveClient
}

func (h *host) SubmitAndBroadcastTx(encryptedParams common.EncryptedParamsSendRawTx) (*responses.RawTx, error) {
	if h.stopControl.IsStopping() {
		return nil, responses.ToInternalError(fmt.Errorf("requested SubmitAndBroadcastTx with the host stopping"))
	}
	encryptedTx := common.EncryptedTx(encryptedParams)

	enclaveResponse, sysError := h.enclaveClient.SubmitTx(encryptedTx)
	if sysError != nil {
		h.logger.Warn("Could not submit transaction due to sysError.", log.ErrKey, sysError)
		return nil, sysError
	}
	if enclaveResponse.Error() != nil {
		h.logger.Trace("Could not submit transaction.", log.ErrKey, enclaveResponse.Error())
		return enclaveResponse, nil //nolint: nilerr
	}

	if h.config.NodeType != common.Sequencer {
		err := h.p2p().SendTxToSequencer(encryptedTx)
		if err != nil {
			return nil, fmt.Errorf("could not broadcast transaction to sequencer. Cause: %w", err)
		}
	}

	return enclaveResponse, nil
}

func (h *host) HandleTransaction(tx common.EncryptedTx) {
	h.txP2PCh <- tx
}

func (h *host) Subscribe(id rpc.ID, encryptedLogSubscription common.EncryptedParamsLogSubscription, matchedLogsCh chan []byte) error {
	if h.stopControl.IsStopping() {
		return responses.ToInternalError(fmt.Errorf("requested Subscribe with the host stopping"))
	}
	err := h.EnclaveClient().Subscribe(id, encryptedLogSubscription)
	if err != nil {
		return fmt.Errorf("could not create subscription with enclave. Cause: %w", err)
	}
	h.logEventManager.AddSubscription(id, matchedLogsCh)
	return nil
}

func (h *host) Unsubscribe(id rpc.ID) {
	if h.stopControl.IsStopping() {
		h.logger.Error("requested Subscribe with the host stopping")
	}
	err := h.EnclaveClient().Unsubscribe(id)
	if err != nil {
		h.logger.Error("could not terminate subscription", log.SubIDKey, id, log.ErrKey, err)
	}
	h.logEventManager.RemoveSubscription(id)
}

func (h *host) Stop() error {
	// block all incoming requests
	h.stopControl.Stop()

	h.logger.Info("Host received a stop command. Attempting shutdown...")
	h.interrupter.Close()
	h.shutdownGroup.Wait()

	// stop all registered services
	for name, service := range h.services {
		if err := service.Stop(); err != nil {
			h.logger.Error("failed to stop service", "service", name, log.ErrKey, err)
		}
	}

	// Leave some time for all processing to finish before exiting the main loop.
	time.Sleep(time.Second)

	if err := h.enclaveClient.Stop(); err != nil {
		return fmt.Errorf("failed to stop enclave server - %w", err)
	}
	if err := h.enclaveClient.StopClient(); err != nil {
		return fmt.Errorf("failed to stop enclave RPC client - %w", err)
	}

	if err := h.db.Stop(); err != nil {
		return fmt.Errorf("failed to stop DB - %w", err)
	}

	h.logger.Info("Host shut down successfully.")
	return nil
}

// HealthCheck returns whether the host, enclave and DB are healthy
func (h *host) HealthCheck() (*hostcommon.HealthCheck, error) {
	if h.stopControl.IsStopping() {
		return nil, responses.ToInternalError(fmt.Errorf("requested HealthCheck with the host stopping"))
	}

	// check the enclave health, which in turn checks the DB health
	enclaveHealthy, err := h.enclaveClient.HealthCheck()
	if err != nil {
		// simplest iteration, log the error and just return that it's not healthy
		h.logger.Error("unable to HealthCheck enclave", log.ErrKey, err)
	}

	// todo (@matt) make the host health check more generic, it should just collate the health of all services
	//   (so at this point we can just loop through all services and call their health check, collate the responses)
	l1RepoService := h.l1Repo().(hostcommon.Service)
	l1RepoHealthy := l1RepoService.HealthStatus().OK()
	p2pService := h.p2p().(hostcommon.Service)
	p2pHealthy := p2pService.HealthStatus().OK()

	isL1Synced := h.enclaveState.InSyncWithL1()

	// Overall health is achieved when all parts are healthy
	obscuroNodeHealth := p2pHealthy && l1RepoHealthy && enclaveHealthy && isL1Synced

	return &hostcommon.HealthCheck{
		HealthCheckHost: &hostcommon.HealthCheckHost{
			// P2PStatus: h.p2p().Status(),
			L1Repo:   l1RepoHealthy,
			L1Synced: isL1Synced,
		},
		HealthCheckEnclave: &hostcommon.HealthCheckEnclave{
			EnclaveHealthy: enclaveHealthy,
		},
		OverallHealth: obscuroNodeHealth,
	}, nil
}

// Waits for enclave to be available, printing a wait message every two seconds.
func (h *host) waitForEnclave() common.Status {
	counter := 0
	var status common.Status
	var err error
	for status, err = h.enclaveClient.Status(); err != nil; {
		if counter >= 20 {
			h.logger.Info(fmt.Sprintf("Waiting for enclave on %s. Latest connection attempt failed", h.config.EnclaveRPCAddress), log.ErrKey, err)
			counter = 0
		}

		time.Sleep(100 * time.Millisecond)
		counter++
	}
	h.logger.Info("Connected to enclave service.", "enclaveStatus", status)
	return status
}

// starts the host main processing loop
func (h *host) startProcessing() {
	h.p2p().SubscribeForTx(h)
	if h.config.NodeType == common.Sequencer {
		go h.startBatchProduction()  // periodically request a new batch from enclave
		go h.startRollupProduction() // periodically request a new rollup from enclave
	}

	h.l1Repo().Subscribe(h)    // start receiving new L1 blocks as they arrive (they'll be ignored if we're still behind)
	h.l2Repo().Subscribe(h)    // start receiving new L2 batches as they arrive (they'll be ignored if we're still behind)
	go h.startBatchStreaming() // streams batches and events from the enclave.

	// Main Processing Loop -
	// - Process new blocks from the L1 node
	// - Process new Transactions gossiped from L2 Peers
	for {
		// if enclave is behind the L1 head then this will process the next block it needs. Just one block and then defer
		// to the select to see if there's anything waiting to process.
		// todo (@matt) this looping with sleep method is temporary while we still have queues for p2p, step towards enclave-guardian PR
		catchingUp := h.catchUpL1Block()
		loopTime := 10 * time.Millisecond
		if catchingUp {
			// we're catching up on L1 blocks, so we don't want to sleep we want to keep processing
			loopTime = 0
		} else { // we're caught up on L1 blocks, now make sure we're caught up on L2 batches (again, just one batch at a time)
			catchingUp = h.catchUpL2Batch()
			if catchingUp {
				// we're catching up on L2 batches, so we don't want to sleep we want to keep processing
				loopTime = 0
			}
		}
		// check if anything is waiting for the host to process before looping back to the catch-up
		select {
		case tx := <-h.txP2PCh:
			resp, sysError := h.enclaveClient.SubmitTx(tx)
			if sysError != nil {
				h.logger.Warn("Could not submit transaction due to sysError", log.ErrKey, sysError)
				continue
			}
			if resp.Error() != nil {
				h.logger.Trace("Could not submit transaction", log.ErrKey, resp.Error())
			}

		case <-h.interrupter.Done():
			return

		case <-time.After(loopTime):
			// todo (@matt) this is temporary, step towards enclave-guardian PR, ensures we look back to catching up
		}
	}
}

func (h *host) processL1Block(block *types.Block, isLatestBlock bool) error {
	// For the genesis block the parent is nil
	if block == nil {
		return nil
	}

	h.logger.Info("Processing L1 block", log.BlockHashKey, block.Hash(), log.BlockHeightKey, block.Number(), "isLatestBlock", isLatestBlock)
	h.processL1BlockTransactions(block)

	// submit each block to the enclave for ingestion plus validation
	blockSubmissionResponse, err := h.enclaveClient.SubmitL1Block(*block, h.l1Repo().FetchReceipts(block), isLatestBlock)
	if err != nil {
		go h.checkEnclaveStatus()
		return fmt.Errorf("did not ingest block %s. Cause: %w", block.Hash(), err)
	}
	h.enclaveState.OnProcessedBlock(block.Hash())
	if blockSubmissionResponse == nil {
		return fmt.Errorf("no block submission response given for a submitted l1 block")
	}

	err = h.db.AddBlockHeader(block.Header())
	if err != nil {
		return fmt.Errorf("submitted block to enclave but could not store the block processing result. Cause: %w", err)
	}

	err = h.publishSharedSecretResponses(blockSubmissionResponse.ProducedSecretResponses)
	if err != nil {
		h.logger.Error("failed to publish response to secret request", log.ErrKey, err)
	}

	return nil
}

// React to Obscuro L1 transactions if needed
func (h *host) processL1BlockTransactions(b *types.Block) {
	// if there are any secret responses in the block we should refresh our P2P list to re-sync with the network
	respTxs := h.l1Publisher().ExtractSecretResponses(b)
	if len(respTxs) > 0 {
		err := h.refreshP2PPeerList()
		if err != nil {
			h.logger.Error("Failed to update p2p peer list", log.ErrKey, err)
		}
	}
}

func (h *host) storeBatch(producedBatch *common.ExtBatch) {
	defer h.logger.Info("Batch stored", log.BatchHashKey, producedBatch.Hash(), log.DurationKey, measure.NewStopwatch())

	h.enclaveState.OnProcessedBatch(producedBatch.Header.SequencerOrderNo)
	err := h.l2Repo().AddBatch(producedBatch)
	// if batch already exists, we can ignore it (this happens, for example, when enclave finds batches in rollups that host has already seen)
	if err != nil && !errors.Is(err, errutil.ErrAlreadyExists) {
		// unexpected error storing batch
		h.logger.Error("could not store batch", log.BatchHashKey, producedBatch.Hash(), log.ErrKey, err)
	}
}

// Creates a batch based on the rollup and distributes it to all other nodes.
func (h *host) storeAndDistributeBatch(producedBatch *common.ExtBatch) {
	defer h.logger.Info("Batch stored and distributed", log.BatchHashKey, producedBatch.Hash(), log.DurationKey, measure.NewStopwatch())

	h.storeBatch(producedBatch)

	err := h.p2p().BroadcastBatches([]*common.ExtBatch{producedBatch})
	if err != nil {
		h.logger.Error("could not broadcast batch", log.BatchHashKey, producedBatch.Hash(), log.ErrKey, err)
	}
}

// This method implements the procedure by which a node obtains the secret
func (h *host) requestSecret() error {
	if h.config.IsGenesis {
		return h.generateAndBroadcastSecret()
	}
	h.logger.Info("Requesting secret.")
	att, err := h.enclaveClient.Attestation()
	if err != nil {
		return fmt.Errorf("could not retrieve attestation from enclave. Cause: %w", err)
	}
	if att.Owner != h.config.ID {
		return fmt.Errorf("host has ID %s, but its enclave produced an attestation using ID %s", h.config.ID.Hex(), att.Owner.Hex())
	}

	// returns the L1 block when the request was published, any response will be after that block
	awaitFromBlock, err := h.l1Publisher().RequestSecret(att)
	if err != nil {
		h.logger.Crit("could not receive the secret", log.ErrKey, err)
	}

	// keep checking L1 blocks until we find a secret response for our request or timeout
	err = retry.Do(func() error {
		nextBlock, _, err := h.l1Repo().FetchNextBlock(awaitFromBlock)
		if err != nil {
			return fmt.Errorf("next block after block=%s not found - %w", awaitFromBlock, err)
		}
		secretRespTxs := h.l1Publisher().ExtractSecretResponses(nextBlock)
		if err != nil {
			return fmt.Errorf("could not extract secret responses from block=%s - %w", nextBlock.Hash(), err)
		}
		for _, s := range secretRespTxs {
			if s.RequesterID.Hex() == h.config.ID.Hex() {
				err = h.enclaveClient.InitEnclave(s.Secret)
				if err != nil {
					h.logger.Warn("could not initialize enclave with received secret response", "err", err)
					continue // try the next secret response in the block if there are more
				}
				return nil // successfully initialized enclave with secret
			}
		}
		awaitFromBlock = nextBlock.Hash()
		return errors.New("no valid secret received in block")
	}, retry.NewTimeoutStrategy(maxWaitForSecretResponse, 500*time.Millisecond))
	if err != nil {
		// something went wrong, check the enclave status in case it is an enclave problem and let the main loop try again when appropriate
		return errors.Wrap(err, "no valid secret received for enclave")
	}

	h.logger.Info("Secret received")
	h.enclaveState.OnSecretProvided()
	return nil
}

func (h *host) publishSharedSecretResponses(scrtResponses []*common.ProducedSecretResponse) error {
	for _, scrtResponse := range scrtResponses {
		// todo (#1624) - implement proper protocol so only one host responds to this secret requests initially
		// 	for now we just have the genesis host respond until protocol implemented
		if !h.config.IsGenesis {
			h.logger.Trace("Not genesis node, not publishing response to secret request.",
				"requester", scrtResponse.RequesterID)
			return nil
		}

		err := h.l1Publisher().PublishSecretResponse(scrtResponse)
		if err != nil {
			return errors.Wrap(err, "could not publish secret response")
		}
	}
	return nil
}

// Whenever we receive a new shared secret response transaction or restart the host, we update our list of P2P peers
func (h *host) refreshP2PPeerList() error {
	// We make a call to the L1 node to retrieve the latest list of aggregators
	peers, err := h.l1Publisher().FetchLatestPeersList()
	if err != nil {
		return errors.Wrap(err, "could not fetch latest peers list from L1")
	}

	h.p2p().UpdatePeerList(peers)
	return nil
}

func (h *host) startBatchProduction() {
	defer h.logger.Info("Stopping batch production")

	h.shutdownGroup.Add(1)
	defer h.shutdownGroup.Done()

	interval := h.config.BatchInterval
	if interval == 0 {
		interval = 1 * time.Second
	}
	batchProdTicker := time.NewTicker(interval)
	for {
		select {
		case <-batchProdTicker.C:
			if !h.enclaveState.InSyncWithL1() {
				// if we're behind the L1, we don't want to produce batches
				h.logger.Debug("skipping batch production because L1 is not up to date")
				continue
			}
			h.logger.Debug("create batch")
			err := h.enclaveClient.CreateBatch()
			if err != nil {
				h.logger.Warn("unable to produce batch", log.ErrKey, err)
			}
		case <-h.interrupter.Done():
			return
		}
	}
}

func (h *host) startBatchStreaming() {
	defer h.logger.Info("Stopping batch streaming")

	h.shutdownGroup.Add(1)
	defer h.shutdownGroup.Done()

	var startingBatch *gethcommon.Hash
	header, err := h.db.GetHeadBatchHeader()
	if err != nil {
		h.logger.Warn("Could not retrieve head batch header for batch streaming", log.ErrKey, err)
	} else {
		batchHash := header.Hash()
		startingBatch = &batchHash
		h.enclaveState.OnReceivedBatch(header.SequencerOrderNo)
		h.logger.Info("Streaming from latest known head batch", log.BatchHashKey, startingBatch)
	}

	streamChan, stop := h.enclaveClient.StreamL2Updates(startingBatch)
	var lastBatch *common.ExtBatch
	for {
		select {
		case <-h.interrupter.Done():
			stop()
			return
		case resp, ok := <-streamChan:
			if !ok {
				stop()
				h.logger.Warn("Batch streaming failed. Reconnecting from latest received batch after 3 seconds")
				time.Sleep(3 * time.Second)

				if lastBatch != nil {
					bHash := lastBatch.Hash()
					streamChan, stop = h.enclaveClient.StreamL2Updates(&bHash)
				} else {
					streamChan, stop = h.enclaveClient.StreamL2Updates(nil)
				}
				continue
			}

			if resp.Batch != nil {
				lastBatch = resp.Batch
				h.logger.Trace("Received batch from stream", log.BatchHashKey, lastBatch.Hash())
				if h.config.NodeType == common.Sequencer {
					h.logger.Info("Batch produced", log.BatchHeightKey, resp.Batch.Header.Number, log.BatchHashKey, resp.Batch.Hash())
					h.enclaveState.OnReceivedBatch(resp.Batch.Header.SequencerOrderNo)
					h.storeAndDistributeBatch(resp.Batch)
				} else {
					h.logger.Info("Batch streamed", log.BatchHeightKey, resp.Batch.Header.Number, log.BatchHashKey, resp.Batch.Hash())
					h.storeBatch(resp.Batch)
				}
			}

			if resp.Logs != nil {
				h.logEventManager.SendLogsToSubscribers(&resp.Logs)
			}
		}
	}
}

func (h *host) startRollupProduction() {
	defer h.logger.Info("Stopping rollup production")

	h.shutdownGroup.Add(1)
	defer h.shutdownGroup.Done()

	interval := h.config.RollupInterval
	if interval == 0 {
		interval = 5 * time.Second
	}
	rollupTicker := time.NewTicker(interval)
	for {
		select {
		case <-rollupTicker.C:
			if !h.enclaveState.IsUpToDate() {
				// if we're behind the L1, we don't want to produce rollups
				h.logger.Debug("skipping rollup production because L1 is not up to date", "enclaveState", h.enclaveState)
				continue
			}
			producedRollup, err := h.enclaveClient.CreateRollup()
			if err != nil {
				h.logger.Error("unable to produce rollup", log.ErrKey, err)
			} else {
				// fire-and-forget (track the receipt asynchronously)
				h.l1Publisher().PublishRollup(producedRollup)
			}
		case <-h.interrupter.Done():
			return
		}
	}
}

// Checks the host config is valid.
func (h *host) validateConfig() {
	if h.config.IsGenesis && h.config.NodeType != common.Sequencer {
		h.logger.Crit("genesis node must be the sequencer")
	}
	if !h.config.IsGenesis && h.config.NodeType == common.Sequencer {
		h.logger.Crit("only the genesis node can be a sequencer")
	}

	if h.config.P2PPublicAddress == "" {
		h.logger.Crit("the host must specify a public P2P address")
	}
}

// this function should be fired off in a new goroutine whenever the status of the enclave needs to be verified
// (e.g. if we've seen unexpected errors from the enclave client)
func (h *host) checkEnclaveStatus() {
	// only allow one status request at a time, if flag is false, atomically swap it to true and continue
	if h.statusReqInFlight.CompareAndSwap(false, true) {
		defer h.statusReqInFlight.Store(false) // clear flag after request completed
		s, err := h.enclaveClient.Status()
		if err != nil {
			h.logger.Error("could not get enclave status", log.ErrKey, err)
			// we record this as a disconnection, we can't get any more info from the enclave about status currently
			h.enclaveState.OnDisconnected()
			return
		}
		h.enclaveState.OnEnclaveStatus(s)
	}
}

// returns true if processed a block
func (h *host) catchUpL1Block() bool {
	// nothing to do if host is stopping or L1 is up-to-date
	if h.stopControl.IsStopping() || h.enclaveState.InSyncWithL1() {
		return false
	}
	prevHead := h.enclaveState.GetEnclaveL1Head()
	h.logger.Trace("fetching next block", log.BlockHashKey, prevHead)
	block, isLatest, err := h.l1Repo().FetchNextBlock(prevHead)
	if err != nil {
		// ErrNoNext block occurs sometimes if we caught up with the L1 head, but other errors are unexpected
		if !errors.Is(err, l1.ErrNoNextBlock) {
			h.logger.Warn("unexpected error fetching next L1 block", log.ErrKey, err)
		}
		return false // nothing to do if we can't fetch the next block
	}
	h.submitDataToEnclaveLock.Lock()
	defer h.submitDataToEnclaveLock.Unlock()
	err = h.processL1Block(block, isLatest)
	if err != nil {
		h.logger.Warn("unable to process L1 block", log.ErrKey, err)
	}
	return true
}

// returns true if processed a batch to indicate to the host that we're still catching up (so, false => up-to-date)
func (h *host) catchUpL2Batch() bool {
	if h.stopControl.IsStopping() || h.enclaveState.IsUpToDate() {
		return false
	}
	prevHead := h.enclaveState.GetEnclaveL2Head()
	nextHead := big.NewInt(0)
	if prevHead != nil {
		nextHead = prevHead.Add(prevHead, big.NewInt(1))
	}
	h.logger.Trace("fetching next batch", log.BatchSeqNoKey, nextHead)
	batch, err := h.l2Repo().FetchBatchBySeqNo(nextHead)
	if err != nil {
		if !errors.Is(err, errutil.ErrNotFound) {
			h.logger.Warn("unexpected error fetching next L2 batch", log.ErrKey, err)
		}
		return false // nothing to do if we can't fetch the next batch, the repo will request it p2p if it's missing, so it should become available soon
	}
	h.submitDataToEnclaveLock.Lock()
	defer h.submitDataToEnclaveLock.Unlock()
	if err = h.enclaveClient.SubmitBatch(batch); err != nil {
		h.logger.Warn("unable to submit batch to enclave", log.ErrKey, err)
	} else {
		// successfully submitted, update the enclave state
		h.enclaveState.OnProcessedBatch(batch.Header.SequencerOrderNo)
	}
	return true
}

func (h *host) getService(name string) hostcommon.Service {
	service, ok := h.services[name]
	if !ok {
		h.logger.Crit("requested service not registered", "name", name)
	}
	return service
}

func (h *host) p2p() hostcommon.P2P {
	return h.getService(hostcommon.P2PName).(hostcommon.P2P)
}

func (h *host) l1Repo() hostcommon.L1BlockRepository {
	return h.getService(hostcommon.L1BlockRepositoryName).(hostcommon.L1BlockRepository)
}

func (h *host) l1Publisher() hostcommon.L1Publisher {
	return h.getService(hostcommon.L1PublisherName).(hostcommon.L1Publisher)
}

func (h *host) l2Repo() hostcommon.L2BatchRepository {
	return h.getService(hostcommon.L2BatchRepositoryName).(hostcommon.L2BatchRepository)
}
