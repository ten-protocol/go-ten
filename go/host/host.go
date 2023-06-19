package host

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
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
	"github.com/obscuronet/go-obscuro/go/host/batchmanager"
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
	// Attempts to broadcast the rollup transaction to the L1. Worst-case, equates to 7 seconds, plus time per request.
	l1TxTriesRollup = 3
	// Attempts to send secret initialisation, request or response transactions to the L1. Worst-case, equates to 63 seconds, plus time per request.
	l1TxTriesSecret = 7

	// todo - these values have to be configurable
	maxWaitForL1Receipt       = 100 * time.Second
	retryIntervalForL1Receipt = 10 * time.Second
	blockStreamWarningTimeout = 30 * time.Second
)

// Implementation of host.Host.
type host struct {
	config  *config.HostConfig
	shortID uint64

	p2p           hostcommon.P2P       // For communication with other Obscuro nodes
	ethClient     ethadapter.EthClient // For communication with the L1 node
	enclaveClient common.Enclave       // For communication with the enclave

	// control the host lifecycle
	interrupter   breaker.Interface
	shutdownGroup sync.WaitGroup

	// ignore incoming requests
	stopControl *stopcontrol.StopControl

	l1BlockProvider hostcommon.ReconnectingBlockProvider
	l1UpToDate      atomic.Bool                     // Whether the last submitted L1 block was the head L1 block
	txP2PCh         chan common.EncryptedTx         // The channel that new transactions from peers are sent to
	batchP2PCh      chan common.EncodedBatchMsg     // The channel that new batches from peers are sent to
	batchRequestCh  chan common.EncodedBatchRequest // The channel that batch requests from peers are sent to

	db *db.DB // Stores the host's publicly-available data

	mgmtContractLib mgmtcontractlib.MgmtContractLib // Library to handle Management Contract lib operations
	ethWallet       wallet.Wallet                   // Wallet used to issue ethereum transactions
	logEventManager *events.LogEventManager
	batchManager    *batchmanager.BatchManager

	logger gethlog.Logger

	metricRegistry gethmetrics.Registry
}

func NewHost(
	config *config.HostConfig,
	p2p hostcommon.P2P,
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
	host := &host{
		// config
		config:  config,
		shortID: common.ShortAddress(config.ID),

		// Communication layers.
		p2p:           p2p,
		ethClient:     ethClient,
		enclaveClient: enclaveClient,

		// incoming data
		l1BlockProvider: ethadapter.NewEthBlockProvider(ethClient, logger),
		l1UpToDate:      atomic.Bool{},
		txP2PCh:         make(chan common.EncryptedTx),
		batchP2PCh:      make(chan common.EncodedBatchMsg),
		batchRequestCh:  make(chan common.EncodedBatchRequest),

		// Initialize the host DB
		db: database,

		mgmtContractLib: mgmtContractLib, // library that provides a handler for Management Contract
		ethWallet:       ethWallet,       // the host's ethereum wallet
		logEventManager: events.NewLogEventManager(logger),
		batchManager:    batchmanager.NewBatchManager(database, config.P2PPublicAddress, logger),

		logger:         logger,
		metricRegistry: regMetrics,

		stopControl: stopcontrol.New(),
	}

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

	encodedAttestation, err := common.EncodeAttestation(attestation)
	if err != nil {
		return fmt.Errorf("could not encode attestation. Cause: %w", err)
	}

	secret, err := h.enclaveClient.GenerateSecret()
	if err != nil {
		return fmt.Errorf("could not generate secret. Cause: %w", err)
	}

	l1tx := &ethadapter.L1InitializeSecretTx{
		AggregatorID:  &h.config.ID,
		Attestation:   encodedAttestation,
		InitialSecret: secret,
		HostAddress:   h.config.P2PPublicAddress,
	}
	initialiseSecretTx := h.mgmtContractLib.CreateInitializeSecret(l1tx, h.ethWallet.GetNonceAndIncrement())
	initialiseSecretTx, err = h.ethClient.EstimateGasAndGasPrice(initialiseSecretTx, h.ethWallet.Address())
	if err != nil {
		h.ethWallet.SetNonce(h.ethWallet.GetNonce() - 1)
		return err
	}
	// we block here until we confirm a successful receipt. It is important this is published before the initial rollup.
	err = h.signAndBroadcastL1Tx(initialiseSecretTx, l1TxTriesSecret, true)
	if err != nil {
		return fmt.Errorf("failed to initialise enclave secret. Cause: %w", err)
	}
	h.logger.Info("Node is genesis node. Secret was broadcast.")
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
		err := h.p2p.SendTxToSequencer(encryptedTx)
		if err != nil {
			return nil, fmt.Errorf("could not broadcast transaction to sequencer. Cause: %w", err)
		}
	}

	return enclaveResponse, nil
}

func (h *host) ReceiveTx(tx common.EncryptedTx) {
	h.txP2PCh <- tx
}

func (h *host) ReceiveBatches(batches common.EncodedBatchMsg) {
	h.batchP2PCh <- batches
}

func (h *host) ReceiveBatchRequest(batchRequest common.EncodedBatchRequest) {
	h.batchRequestCh <- batchRequest
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

	if err := h.p2p.StopListening(); err != nil {
		return fmt.Errorf("failed to close transaction P2P listener cleanly - %w", err)
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

	l1BlockProviderStatus := h.l1BlockProvider.HealthStatus()
	isL1Synced := h.l1UpToDate.Load()

	// Overall health is achieved when all parts are healthy
	obscuroNodeHealth := h.p2p.HealthCheck() && l1BlockProviderStatus.Healthy && enclaveHealthy && isL1Synced

	return &hostcommon.HealthCheck{
		HealthCheckHost: &hostcommon.HealthCheckHost{
			P2PStatus:       h.p2p.Status(),
			L1BlockProvider: &l1BlockProviderStatus,
			L1Synced:        isL1Synced,
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
	h.p2p.StartListening(h)
	if h.config.NodeType == common.Sequencer {
		go h.startBatchProduction()  // periodically request a new batch from enclave
		go h.startRollupProduction() // periodically request a new rollup from enclave
	}

	go h.startBatchStreaming() // streams batches and events from the enclave.

	// The blockStream channel is a stream of consecutive, canonical blocks. BlockStream may be replaced with a new
	// stream ch during the main loop if enclave gets out-of-sync, and we need to stream from an earlier block
	blockStream, err := h.l1BlockProvider.StartStreamingFromHash(h.config.L1StartHash)
	if err != nil {
		// maybe start hash wasn't provided or couldn't be found, instead we stream from L1 genesis
		// note: in production this could be expensive, hence the WARN log message
		// todo (@matt) - review whether we should fail here
		h.logger.Warn("unable to stream from L1StartHash", log.ErrKey, err, "l1StartHash", h.config.L1StartHash)
		blockStream, err = h.l1BlockProvider.StartStreamingFromHeight(big.NewInt(1))
		if err != nil {
			h.logger.Crit("unable to stream l1 blocks for enclave", log.ErrKey, err)
		}
	}

	// Main Processing Loop -
	// - Process new blocks from the L1 node
	// - Process new Transactions gossiped from L2 Peers
	for {
		select {
		case b := <-blockStream.Stream:
			isLive := h.l1BlockProvider.IsLatest(b) // checks whether the block is the current head of the L1 (false if there is a newer block available)
			err = h.processL1Block(b, isLive)
			if err != nil {
				// handle the error, replace the blockStream if necessary (e.g. if stream needs resetting based on enclave's reported L1 head)
				blockStream = h.handleProcessBlockErr(b, blockStream, err)
				// failed to update the L1 head, so assume we're behind
				h.l1UpToDate.Store(false)
				continue
			}
			h.l1UpToDate.Store(isLive)

		case tx := <-h.txP2PCh:
			resp, sysError := h.enclaveClient.SubmitTx(tx)
			if sysError != nil {
				h.logger.Warn("Could not submit transaction due to sysError", log.ErrKey, sysError)
				continue
			}
			if resp.Error() != nil {
				h.logger.Trace("Could not submit transaction", log.ErrKey, resp.Error())
			}

		// todo (#718) - adopt a similar approach to blockStream, where we have a BatchProvider that streams new batches.
		case batchMsg := <-h.batchP2PCh:
			// todo (#1623) - discard p2p messages if enclave won't be able to make use of them (e.g. we're way behind L1 head)
			if err := h.handleBatches(&batchMsg); err != nil {
				h.logger.Error("Could not handle batches. ", log.ErrKey, err)
			}

		case batchRequest := <-h.batchRequestCh:
			if err := h.handleBatchRequest(&batchRequest); err != nil {
				h.logger.Error("Could not handle batch request. ", log.ErrKey, err)
			}

		case <-h.interrupter.Done():
			blockStream.Stop()
			return
		}
	}
}

func (h *host) handleProcessBlockErr(processedBlock *types.Block, stream *hostcommon.BlockStream, err error) *hostcommon.BlockStream {
	var rejErr *errutil.BlockRejectError
	var resetFrom gethcommon.Hash
	if errors.As(err, &rejErr) {
		h.logger.Info("Block rejected by enclave.", log.ErrKey, rejErr, log.BlockHashKey, processedBlock.Hash(), log.BlockHeightKey, processedBlock.Number())
		if errors.Is(rejErr, errutil.ErrBlockAlreadyProcessed) {
			// resetting stream after rejection for duplicate is a possible optimisation in future but it's rarely an expensive case and
			// it's a risky optimisation (need to ensure it can't get stuck in a loop)
			// Instead we assume that only one or two blocks are being repeated (probably from revisiting a fork that was
			// abandoned) and then the enclave will be progressing again
			return stream
		}
		if rejErr.L1Head == (gethcommon.Hash{}) {
			h.logger.Warn("No L1 head information provided by enclave, continuing with existing stream")
			return stream
		}
		// prepare to reset the stream from the L1 head provided by the enclave
		resetFrom = rejErr.L1Head
	} else {
		// received unexpected error (no useful information from the enclave)
		// we log it out and retry the stream from the same block
		h.logger.Warn("Error processing block, resetting block provider to retry", log.ErrKey, err)
		resetFrom = processedBlock.Hash()
	}
	h.logger.Info("Resetting block provider stream", "streamFrom", resetFrom)
	// streaming from the latest canonical ancestor of the enclave's L1 head (we may end up re-streaming some blocks it's
	//	already processed, but we tolerate those inefficiencies for simplicity for now)
	replacementStream, err := h.l1BlockProvider.StartStreamingFromHash(resetFrom)
	if err != nil {
		h.logger.Warn("Could not reset block provider, continuing with previous stream", log.ErrKey, err)
		return stream
	}
	stream.Stop() // cancel the previous stream and return the replacement
	return replacementStream
}

func (h *host) processL1Block(block *types.Block, isLatestBlock bool) error {
	// For the genesis block the parent is nil
	if block == nil {
		return nil
	}

	h.processL1BlockTransactions(block)

	// submit each block to the enclave for ingestion plus validation
	blockSubmissionResponse, err := h.enclaveClient.SubmitL1Block(*block, h.extractReceipts(block), isLatestBlock)
	if err != nil {
		return fmt.Errorf("did not ingest block %s. Cause: %w", block.Hash(), err)
	}
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

// Looks at each transaction in the block, and kicks off special handling for the transaction if needed.
func (h *host) processL1BlockTransactions(b *types.Block) {
	for _, tx := range b.Transactions() {
		t := h.mgmtContractLib.DecodeTx(tx)
		if t == nil {
			continue
		}

		// node received a secret response, we should make sure our p2p addresses are up-to-date
		if _, ok := t.(*ethadapter.L1RespondSecretTx); ok {
			err := h.refreshP2PPeerList()
			if err != nil {
				h.logger.Error("Failed to update p2p peer list", log.ErrKey, err)
				continue
			}
		}
	}
}

// Publishes a rollup to the L1.
func (h *host) publishRollup(producedRollup *common.ExtRollup) {
	encodedRollup, err := common.EncodeRollup(producedRollup)
	if err != nil {
		h.logger.Crit("could not encode rollup.", log.ErrKey, err)
	}
	tx := &ethadapter.L1RollupTx{
		Rollup: encodedRollup,
	}
	h.logger.Info("Publishing rollup", log.RollupHeightKey, producedRollup.Header.Number, "size", len(encodedRollup)/1024, log.RollupHashKey, producedRollup.Hash())

	h.logger.Trace("Sending transaction to publish rollup", "rollup_header",
		gethlog.Lazy{Fn: func() string {
			header, err := json.MarshalIndent(producedRollup.Header, "", "   ")
			if err != nil {
				return err.Error()
			}

			return string(header)
		}}, "rollup_hash", producedRollup.Header.Hash().Hex(), "batches_len", len(producedRollup.BatchPayloads))

	rollupTx := h.mgmtContractLib.CreateRollup(tx, h.ethWallet.GetNonceAndIncrement())
	rollupTx, err = h.ethClient.EstimateGasAndGasPrice(rollupTx, h.ethWallet.Address())
	if err != nil {
		// todo (#1624) - make rollup submission a separate workflow (design and implement the flow etc)
		h.ethWallet.SetNonce(h.ethWallet.GetNonce() - 1)
		h.logger.Error("could not estimate rollup tx", log.ErrKey, err)
		return
	}

	// fire-and-forget (track the receipt asynchronously)
	// todo (#1624) - With a single sequencer, it is problematic if rollup publication fails; handle this case better
	err = h.signAndBroadcastL1Tx(rollupTx, l1TxTriesRollup, true)
	if err != nil {
		h.logger.Error("could not issue rollup tx", log.ErrKey, err)
	} else {
		h.logger.Info("Rollup included in L1", "height", producedRollup.Header.Number, "hash", producedRollup.Hash())
	}
}

func (h *host) storeBatch(producedBatch *common.ExtBatch) {
	defer h.logger.Info("Batch stored", log.BatchHashKey, producedBatch.Hash(), log.DurationKey, measure.NewStopwatch())

	err := h.db.AddBatch(producedBatch)
	if err != nil {
		h.logger.Error("could not store batch", log.BatchHashKey, producedBatch.Hash(), log.ErrKey, err)
	}
}

// Creates a batch based on the rollup and distributes it to all other nodes.
func (h *host) storeAndDistributeBatch(producedBatch *common.ExtBatch) {
	defer h.logger.Info("Batch stored and distributed", log.BatchHashKey, producedBatch.Hash(), log.DurationKey, measure.NewStopwatch())

	h.storeBatch(producedBatch)

	batchMsg := hostcommon.BatchMsg{
		Batches:   []*common.ExtBatch{producedBatch},
		IsCatchUp: false,
	}
	err := h.p2p.BroadcastBatch(&batchMsg)
	if err != nil {
		h.logger.Error("could not broadcast batch", log.BatchHashKey, producedBatch.Hash(), log.ErrKey, err)
	}
}

// `tries` is the number of times to attempt broadcasting the transaction.
// if awaitReceipt is true then this method will block and synchronously wait to check the receipt, otherwise it is fire
// and forget and the receipt tracking will happen in a separate go-routine
func (h *host) signAndBroadcastL1Tx(tx types.TxData, tries uint64, awaitReceipt bool) error {
	var err error
	tx, err = h.ethClient.EstimateGasAndGasPrice(tx, h.ethWallet.Address())
	if err != nil {
		return fmt.Errorf("unable to estimate gas limit and gas price - %w", err)
	}

	signedTx, err := h.ethWallet.SignTransaction(tx)
	if err != nil {
		return err
	}

	h.logger.Info("Host issuing l1 tx", log.TxKey, signedTx.Hash(), "size", signedTx.Size()/1024)

	err = retry.Do(func() error {
		return h.ethClient.SendTransaction(signedTx)
	}, retry.NewDoublingBackoffStrategy(time.Second, tries)) // doubling retry wait (3 tries = 7sec, 7 tries = 63sec)
	if err != nil {
		return fmt.Errorf("broadcasting L1 transaction failed after %d tries. Cause: %w", tries, err)
	}
	h.logger.Info("Successfully submitted tx to L1", "txHash", signedTx.Hash())

	if awaitReceipt {
		// block until receipt is found and then return
		return h.waitForReceipt(signedTx.Hash())
	}

	// else just watch for receipt asynchronously and log if it fails
	go func() {
		// todo (#1624) - consider how to handle the various ways that L1 transactions could fail to improve node operator QoL
		err = h.waitForReceipt(signedTx.Hash())
		if err != nil {
			h.logger.Error("L1 transaction failed", log.ErrKey, err)
		}
	}()

	return nil
}

func (h *host) waitForReceipt(txHash common.TxHash) error {
	var receipt *types.Receipt
	var err error
	err = retry.Do(
		func() error {
			receipt, err = h.ethClient.TransactionReceipt(txHash)
			if err != nil {
				// adds more info on the error
				return fmt.Errorf("unable to get receipt for tx: %s - %w", txHash.Hex(), err)
			}
			return err
		},
		retry.NewTimeoutStrategy(maxWaitForL1Receipt, retryIntervalForL1Receipt),
	)
	if err != nil {
		return fmt.Errorf("receipt for L1 transaction never found despite 'successful' broadcast - %w", err)
	}

	if err == nil && receipt.Status != types.ReceiptStatusSuccessful {
		return fmt.Errorf("unsuccessful receipt found for published L1 transaction, status=%d", receipt.Status)
	}
	h.logger.Debug("L1 transaction receipt found.", log.TxKey, txHash, log.BlockHeightKey, receipt.BlockNumber, log.BlockHashKey, receipt.BlockHash)
	return nil
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
	encodedAttestation, err := common.EncodeAttestation(att)
	if err != nil {
		return fmt.Errorf("could not encode attestation. Cause: %w", err)
	}
	l1tx := &ethadapter.L1RequestSecretTx{
		Attestation: encodedAttestation,
	}
	// record the L1 head height before we submit the secret request so we know which block to watch from
	l1Head, err := h.ethClient.FetchHeadBlock()
	if err != nil {
		panic(fmt.Errorf("could not fetch head L1 block. Cause: %w", err))
	}
	requestSecretTx := h.mgmtContractLib.CreateRequestSecret(l1tx, h.ethWallet.GetNonceAndIncrement())
	requestSecretTx, err = h.ethClient.EstimateGasAndGasPrice(requestSecretTx, h.ethWallet.Address())
	if err != nil {
		h.ethWallet.SetNonce(h.ethWallet.GetNonce() - 1)
		return err
	}
	// we wait until the secret req transaction has succeeded before we start polling for the secret
	err = h.signAndBroadcastL1Tx(requestSecretTx, l1TxTriesSecret, true)
	if err != nil {
		return err
	}

	err = h.awaitSecret(l1Head.Number())
	if err != nil {
		h.logger.Crit("could not receive the secret", log.ErrKey, err)
	}
	return nil
}

func (h *host) handleStoreSecretTx(t *ethadapter.L1RespondSecretTx) bool {
	if t.RequesterID.Hex() != h.config.ID.Hex() {
		// this secret is for somebody else
		return false
	}

	// someone has replied for us
	err := h.enclaveClient.InitEnclave(t.Secret)
	if err != nil {
		h.logger.Info("Failed to initialise enclave with received secret.", log.ErrKey, err.Error())
		return false
	}
	return true
}

func (h *host) publishSharedSecretResponses(scrtResponses []*common.ProducedSecretResponse) error {
	var err error

	for _, scrtResponse := range scrtResponses {
		// todo (#1624) - implement proper protocol so only one host responds to this secret requests initially
		// 	for now we just have the genesis host respond until protocol implemented
		if !h.config.IsGenesis {
			h.logger.Trace("Not genesis node, not publishing response to secret request.",
				"requester", scrtResponse.RequesterID)
			return nil
		}

		l1tx := &ethadapter.L1RespondSecretTx{
			Secret:      scrtResponse.Secret,
			RequesterID: scrtResponse.RequesterID,
			AttesterID:  h.config.ID,
			HostAddress: scrtResponse.HostAddress,
		}
		// todo (#1624) - l1tx.Sign(a.attestationPubKey) doesn't matter as the waitSecret will process a tx that was reverted
		respondSecretTx := h.mgmtContractLib.CreateRespondSecret(l1tx, h.ethWallet.GetNonceAndIncrement(), false)
		respondSecretTx, err = h.ethClient.EstimateGasAndGasPrice(respondSecretTx, h.ethWallet.Address())
		if err != nil {
			h.ethWallet.SetNonce(h.ethWallet.GetNonce() - 1)
			return err
		}
		h.logger.Trace("Broadcasting secret response L1 tx.", "requester", scrtResponse.RequesterID)
		// fire-and-forget (track the receipt asynchronously)
		err = h.signAndBroadcastL1Tx(respondSecretTx, l1TxTriesSecret, false)
		if err != nil {
			return fmt.Errorf("could not broadcast secret response. Cause %w", err)
		}
	}
	return nil
}

// Whenever we receive a new shared secret response transaction or restart the host, we update our list of P2P peers
func (h *host) refreshP2PPeerList() error {
	// We make a call to the L1 node to retrieve the latest list of aggregators
	msg, err := h.mgmtContractLib.GetHostAddresses()
	if err != nil {
		return err
	}
	response, err := h.ethClient.CallContract(msg)
	if err != nil {
		return err
	}
	decodedResponse, err := h.mgmtContractLib.DecodeCallResponse(response)
	if err != nil {
		return err
	}
	hostAddresses := decodedResponse[0]

	// We remove any duplicate addresses and our own address from the retrieved peer list
	var filteredHostAddresses []string
	uniqueHostKeys := make(map[string]bool) // map to track addresses we've seen already
	for _, hostAddress := range hostAddresses {
		// We exclude our own address.
		if hostAddress == h.config.P2PPublicAddress {
			continue
		}
		if _, found := uniqueHostKeys[hostAddress]; !found {
			uniqueHostKeys[hostAddress] = true
			filteredHostAddresses = append(filteredHostAddresses, hostAddress)
		}
	}

	h.p2p.UpdatePeerList(filteredHostAddresses)
	return nil
}

// todo (@stefan) - perhaps extract only relevant logs. There were missing ones when requesting
// the logs filtered from geth.
func (h *host) extractReceipts(block *types.Block) types.Receipts {
	receipts := make(types.Receipts, 0)

	for _, transaction := range block.Transactions() {
		receipt, err := h.ethClient.TransactionReceipt(transaction.Hash())

		if err != nil || receipt == nil {
			h.logger.Error("Problem with retrieving the receipt on the host!", log.ErrKey, err, log.CmpKey, log.CrossChainCmp)
			continue
		}

		h.logger.Trace("Adding receipt", "status", receipt.Status, log.TxKey, transaction.Hash(),
			log.BlockHashKey, block.Hash(), log.CmpKey, log.CrossChainCmp)

		receipts = append(receipts, receipt)
	}

	return receipts
}

func (h *host) awaitSecret(fromHeight *big.Int) error {
	blkStream, err := h.l1BlockProvider.StartStreamingFromHeight(fromHeight)
	if err != nil {
		return err
	}
	defer blkStream.Stop()

	for {
		select {
		case blk := <-blkStream.Stream:
			h.logger.Trace("checking block for secret resp", log.BlockHeightKey, blk.Number())
			if h.checkBlockForSecretResponse(blk) {
				return nil
			}

		case <-time.After(blockStreamWarningTimeout):
			// This will provide useful feedback if things are stuck (and in tests if any goroutines got stranded on this select)
			h.logger.Warn(fmt.Sprintf(" Waiting for secret from the L1. No blocks received for over %s", blockStreamWarningTimeout))

		case <-h.interrupter.Done():
			return nil
		}
	}
}

func (h *host) checkBlockForSecretResponse(block *types.Block) bool {
	for _, tx := range block.Transactions() {
		t := h.mgmtContractLib.DecodeTx(tx)
		if t == nil {
			continue
		}
		if scrtTx, ok := t.(*ethadapter.L1RespondSecretTx); ok {
			ok := h.handleStoreSecretTx(scrtTx)
			if ok {
				h.logger.Info("Stored enclave secret.")
				return true
			}
		}
	}
	// response not found
	return false
}

// Handles an incoming set of batches. There are two possibilities:
// (1) There are no gaps in the historical chain of batches. The new batches can be added immediately
// (2) There are gaps in the historical chain of batches. To avoid an inconsistent state (i.e. one where we have stored
//
//	a batch without its parent), we request the sequencer to resend the batches we've just received, plus any missing
//	historical batches, then discard the received batches. We will store all of these at once when we receive them
func (h *host) handleBatches(encodedBatchMsg *common.EncodedBatchMsg) error {
	var batchMsg *hostcommon.BatchMsg
	err := rlp.DecodeBytes(*encodedBatchMsg, &batchMsg)
	if err != nil {
		return fmt.Errorf("could not decode batches using RLP. Cause: %w", err)
	}

	for _, batch := range batchMsg.Batches {
		h.logger.Info("Received batch from peer", log.BatchHeightKey, batch.Header.Number, log.BatchHashKey, batch.Hash())
		// todo (@stefan) - consider moving to a model where the enclave manages the entire state, to avoid inconsistency.

		// If we do not have the block the batch is tied to, we skip processing the batches for now. We'll catch them
		// up later, once we've received the L1 block.
		_, err = h.db.GetBlockHeader(batch.Header.L1Proof)
		if err != nil {
			if errors.Is(err, errutil.ErrNotFound) {
				return nil
			}
			return fmt.Errorf("could not retrieve block header. Cause: %w", err)
		}

		isParentStored, batchRequest, err := h.batchManager.IsParentStored(batch)
		if err != nil {
			return fmt.Errorf("could not determine whether batch parent was missing. Cause: %w", err)
		}

		// We have encountered a missing parent batch. We abort the storage operation and request the missing batches.
		if !isParentStored {
			h.logger.Info("Parent batch not found. Requesting missing batches.", "fromBatch", batchRequest.CurrentHeadBatch, "isCatchUp", batchMsg.IsCatchUp)
			// We only request the missing batches if the batches did not themselves arrive as part of catch-up, to
			// avoid excessive P2P pressure.
			if !batchMsg.IsCatchUp {
				if err = h.p2p.RequestBatchesFromSequencer(batchRequest); err != nil {
					return fmt.Errorf("could not request historical batches. Cause: %w", err)
				}
				return nil
			}
			return nil
		}

		// We only store the batch locally if it stores successfully on the enclave.
		// todo (@stefan) - edge case when the enclave is restarted and loses some state; move to having enclave as source
		//  of truth re: stored batches
		if err = h.enclaveClient.SubmitBatch(batch); err != nil {
			return fmt.Errorf("could not submit batch. Cause: %w", err)
		}
		if err = h.db.AddBatch(batch); err != nil {
			return fmt.Errorf("could not store batch header. Cause: %w", err)
		}
	}

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
			if !h.l1UpToDate.Load() {
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
					h.logger.Info("Batch produced", log.RollupHeightKey, resp.Batch.Header.Number, log.RollupHashKey, resp.Batch.Hash())
					h.storeAndDistributeBatch(resp.Batch)
				} else {
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
			if !h.l1UpToDate.Load() {
				// if we're behind the L1, we don't want to produce rollups
				h.logger.Debug("skipping rollup production because L1 is not up to date")
				continue
			}
			producedRollup, err := h.enclaveClient.CreateRollup()
			if err != nil {
				h.logger.Error("unable to produce rollup", log.ErrKey, err)
			} else {
				h.publishRollup(producedRollup)
			}
		case <-h.interrupter.Done():
			return
		}
	}
}

// todo (#1625) - only allow requests for batches since last rollup, to avoid DoS attacks.
func (h *host) handleBatchRequest(encodedBatchRequest *common.EncodedBatchRequest) error {
	var batchRequest *common.BatchRequest
	err := rlp.DecodeBytes(*encodedBatchRequest, &batchRequest)
	if err != nil {
		return fmt.Errorf("could not decode batch request using RLP. Cause: %w", err)
	}

	batches, err := h.batchManager.GetBatches(batchRequest, h.enclaveClient)
	if err != nil {
		return fmt.Errorf("could not retrieve batches based on request. Cause: %w", err)
	}

	batchMsg := hostcommon.BatchMsg{
		Batches:   batches,
		IsCatchUp: true,
	}
	return h.p2p.SendBatches(&batchMsg, batchRequest.Requester)
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
