package host

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/naoina/toml"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/errutil"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/common/profiler"
	"github.com/obscuronet/go-obscuro/go/common/retry"
	"github.com/obscuronet/go-obscuro/go/config"
	"github.com/obscuronet/go-obscuro/go/ethadapter"
	"github.com/obscuronet/go-obscuro/go/ethadapter/mgmtcontractlib"
	"github.com/obscuronet/go-obscuro/go/host/batchmanager"
	"github.com/obscuronet/go-obscuro/go/host/db"
	"github.com/obscuronet/go-obscuro/go/host/events"
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
	exitHostCh            chan bool
	stopHostInterrupt     *int32
	bootstrappingComplete *int32 // Marks when the host is done bootstrapping

	l1BlockProvider hostcommon.ReconnectingBlockProvider
	txP2PCh         chan common.EncryptedTx         // The channel that new transactions from peers are sent to
	batchP2PCh      chan common.EncodedBatchMsg     // The channel that new batches from peers are sent to
	batchRequestCh  chan common.EncodedBatchRequest // The channel that batch requests from peers are sent to

	db *db.DB // Stores the host's publicly-available data

	mgmtContractLib mgmtcontractlib.MgmtContractLib // Library to handle Management Contract lib operations
	ethWallet       wallet.Wallet                   // Wallet used to issue ethereum transactions
	logEventManager events.LogEventManager
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
	database := db.NewInMemoryDB(regMetrics) // todo - make this config driven
	host := &host{
		// config
		config:  config,
		shortID: common.ShortAddress(config.ID),

		// Communication layers.
		p2p:           p2p,
		ethClient:     ethClient,
		enclaveClient: enclaveClient,

		// lifecycle channels
		exitHostCh:            make(chan bool),
		stopHostInterrupt:     new(int32),
		bootstrappingComplete: new(int32),

		// incoming data
		l1BlockProvider: ethadapter.NewEthBlockProvider(ethClient, logger),
		txP2PCh:         make(chan common.EncryptedTx),
		batchP2PCh:      make(chan common.EncodedBatchMsg),
		batchRequestCh:  make(chan common.EncodedBatchRequest),

		// Initialize the host DB
		db: database,

		mgmtContractLib: mgmtContractLib, // library that provides a handler for Management Contract
		ethWallet:       ethWallet,       // the host's ethereum wallet
		logEventManager: events.NewLogEventManager(logger),
		batchManager:    batchmanager.NewBatchManager(database, config.P2PPublicAddress),

		logger:         logger,
		metricRegistry: regMetrics,
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
func (h *host) Start() {
	h.validateConfig()

	tomlConfig, err := toml.Marshal(h.config)
	if err != nil {
		h.logger.Crit("could not print host config")
	}
	h.logger.Info("Host started with following config", log.CfgKey, string(tomlConfig))

	go func() {
		// wait for the Enclave to be available
		h.waitForEnclave()

		// TODO the host should only connect to enclaves with the same ID as the host.ID
		// TODO Issue: https://github.com/obscuronet/obscuro-internal/issues/1265

		// todo: we should try to recover the key from a previous run of the node here? Before generating or requesting the key.
		if h.config.IsGenesis {
			err = h.broadcastSecret()
			if err != nil {
				h.logger.Crit("Could not broadcast secret", log.ErrKey, err.Error())
			}
		} else {
			err = h.requestSecret()
			if err != nil {
				h.logger.Crit("Could not request secret", log.ErrKey, err.Error())
			}
		}

		// start the host's main processing loop
		h.startProcessing()
	}()
}

func (h *host) broadcastSecret() error {
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

func (h *host) SubmitAndBroadcastTx(encryptedParams common.EncryptedParamsSendRawTx) (common.EncryptedResponseSendRawTx, error) {
	encryptedTx := common.EncryptedTx(encryptedParams)

	// TODO - #718 - We only need to submit to the enclave as the sequencer. But we still need to return the encrypted
	//  transaction hash, so some round-trip to the enclave is required.
	encryptedResponse, err := h.enclaveClient.SubmitTx(encryptedTx)
	if err != nil {
		return nil, fmt.Errorf("could not submit transaction. Cause: %w", err)
	}

	if h.config.NodeType != common.Sequencer {
		err = h.p2p.SendTxToSequencer(encryptedTx)
		if err != nil {
			return nil, fmt.Errorf("could not broadcast transaction to sequencer. Cause: %w", err)
		}
	}

	return encryptedResponse, nil
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
	err := h.EnclaveClient().Subscribe(id, encryptedLogSubscription)
	if err != nil {
		return fmt.Errorf("could not create subscription with enclave. Cause: %w", err)
	}
	h.logEventManager.AddSubscription(id, matchedLogsCh)
	return nil
}

func (h *host) Unsubscribe(id rpc.ID) {
	err := h.EnclaveClient().Unsubscribe(id)
	if err != nil {
		h.logger.Error("could not terminate subscription", log.SubIDKey, id, log.ErrKey, err)
	}
	h.logEventManager.RemoveSubscription(id)
}

func (h *host) Stop() {
	// block all requests
	atomic.StoreInt32(h.stopHostInterrupt, 1)

	if err := h.p2p.StopListening(); err != nil {
		h.logger.Error("failed to close transaction P2P listener cleanly", log.ErrKey, err)
	}
	if err := h.enclaveClient.Stop(); err != nil {
		h.logger.Error("could not stop enclave server", log.ErrKey, err)
	}
	if err := h.enclaveClient.StopClient(); err != nil {
		h.logger.Error("failed to stop enclave RPC client", log.ErrKey, err)
	}

	// Leave some time for all processing to finish before exiting the main loop.
	time.Sleep(time.Second)
	h.exitHostCh <- true

	h.logger.Info("Host shut down successfully.")
}

// HealthCheck returns whether the host, enclave and DB are healthy
func (h *host) HealthCheck() (*hostcommon.HealthCheck, error) {
	// check the enclave health, which in turn checks the DB health
	enclaveHealthy, err := h.enclaveClient.HealthCheck()
	if err != nil {
		// simplest iteration, log the error and just return that it's not healthy
		h.logger.Error("unable to HealthCheck enclave", "err", err)
	}

	// Overall health is achieved when all parts are healthy
	obscuroNodeHealth := h.p2p.HealthCheck() && enclaveHealthy

	return &hostcommon.HealthCheck{
		HealthCheckHost: &hostcommon.HealthCheckHost{
			P2PStatus: h.p2p.Status(),
		},
		HealthCheckEnclave: &hostcommon.HealthCheckEnclave{
			EnclaveHealthy: enclaveHealthy,
		},
		OverallHealth: obscuroNodeHealth,
	}, nil
}

// Waits for enclave to be available, printing a wait message every two seconds.
func (h *host) waitForEnclave() {
	counter := 0
	for _, err := h.enclaveClient.Status(); err != nil; {
		if counter >= 20 {
			h.logger.Info(fmt.Sprintf("Waiting for enclave on %s. Latest connection attempt failed", h.config.EnclaveRPCAddress), log.ErrKey, err)
			counter = 0
		}

		time.Sleep(100 * time.Millisecond)
		counter++
	}
	h.logger.Info("Connected to enclave service.")
}

// starts the host main processing loop
func (h *host) startProcessing() {
	h.p2p.StartListening(h)

	// The blockStream channel is a stream of consecutive, canonical blocks. BlockStream may be replaced with a new
	// stream ch during the main loop if enclave gets out-of-sync, and we need to stream from an earlier block
	blockStream, err := h.l1BlockProvider.StartStreamingFromHash(h.config.L1StartHash)
	if err != nil {
		// maybe start hash wasn't provided or couldn't be found, instead we stream from L1 genesis
		// note: in production this could be expensive, hence the WARN log message, todo: review whether we should fail here
		h.logger.Warn("unable to stream from L1StartHash", log.ErrKey, err, "l1StartHash", h.config.L1StartHash)
		blockStream, err = h.l1BlockProvider.StartStreamingFromHeight(big.NewInt(1))
		if err != nil {
			h.logger.Crit("unable to stream l1 blocks for enclave", log.ErrKey, err)
		}
	}

	// use the roundInterrupt as a signaling mechanism for interrupting block processing
	// stops processing the current round if a new block arrives
	i := int32(0)
	roundInterrupt := &i

	// Main Processing Loop -
	// - Process new blocks from the L1 node
	// - Process new Transactions gossiped from L2 Peers
	for {
		select {
		case b := <-blockStream.Stream:
			roundInterrupt = triggerInterrupt(roundInterrupt)
			isLive := h.l1BlockProvider.IsLive(b.Hash()) // checks whether the block is the current head of the L1 (false if there is a newer block available)
			err := h.processL1Block(b, isLive)
			if err != nil {
				// handle the error, replace the blockStream if necessary (e.g. if stream needs resetting based on enclave's reported L1 head)
				blockStream = h.handleProcessBlockErr(b, blockStream, err)
			}

		case tx := <-h.txP2PCh:
			// todo: discard p2p messages if enclave won't be able to make use of them (e.g. we're way behind L1 head)
			if _, err := h.enclaveClient.SubmitTx(tx); err != nil {
				h.logger.Warn("Could not submit transaction. ", log.ErrKey, err)
			}

		case batchMsg := <-h.batchP2PCh:
			// todo: discard p2p messages if enclave won't be able to make use of them (e.g. we're way behind L1 head)
			if err := h.handleBatches(&batchMsg); err != nil {
				h.logger.Error("Could not handle batches. ", log.ErrKey, err)
			}

		case batchRequest := <-h.batchRequestCh:
			if err := h.handleBatchRequest(&batchRequest); err != nil {
				h.logger.Error("Could not handle batch request. ", log.ErrKey, err)
			}

		case <-h.exitHostCh:
			return
		}
	}
}

func (h *host) handleProcessBlockErr(processedBlock *types.Block, stream *hostcommon.BlockStream, err error) *hostcommon.BlockStream {
	var rejErr *common.BlockRejectError
	if !errors.As(err, &rejErr) {
		// received unexpected error (no useful information from the enclave)
		// we log it out and ignore it until the enclave tells us more information
		h.logger.Warn("Error processing block.", log.ErrKey, err)
		return stream
	}
	h.logger.Info("Block rejected by enclave.", log.ErrKey, rejErr, "blk", processedBlock.Hash(), "blkHeight", processedBlock.Number())
	if errors.Is(rejErr, common.ErrBlockAlreadyProcessed) {
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
	h.logger.Info("Resetting block provider stream to enclave latest head.", "streamFrom", rejErr.L1Head)
	// streaming from the latest canonical ancestor of the enclave's L1 head (we may end up re-streaming some things it's
	//	already processed, but we tolerate those failures)
	replacementStream, err := h.l1BlockProvider.StartStreamingFromHash(rejErr.L1Head)
	if err != nil {
		h.logger.Warn("Could not reset block provider, continuing with previous stream", log.ErrKey, err)
		return stream
	}
	stream.Stop() // cancel the previous stream and return the replacement
	return replacementStream
}

// activates the given interrupt (atomically) and returns a new interrupt
func triggerInterrupt(interrupt *int32) *int32 {
	// Notify the previous round to stop work
	atomic.StoreInt32(interrupt, 1)
	i := int32(0)
	return &i
}

func (h *host) processL1Block(block *types.Block, isLatestBlock bool) error {
	// For the genesis block the parent is nil
	if block == nil {
		return nil
	}

	h.processL1BlockTransactions(block)

	// submit each block to the enclave for ingestion plus validation
	result, err := h.enclaveClient.SubmitL1Block(*block, h.extractReceipts(block), isLatestBlock)
	if err != nil {
		return fmt.Errorf("did not ingest block b_%d. Cause: %w", common.ShortHash(block.Hash()), err)
	}
	err = h.db.AddBlockHeader(block.Header())
	if err != nil {
		return fmt.Errorf("submitted block to enclave but could not store the block processing result. Cause: %w", err)
	}

	h.logEventManager.SendLogsToSubscribers(result)

	err = h.publishSharedSecretResponses(result.ProducedSecretResponses)
	if err != nil {
		h.logger.Error("failed to publish response to secret request", log.ErrKey, err)
	}

	// If we're not the sequencer, we do not need to produce the genesis or publish and distribute rollups.
	if h.config.NodeType != common.Sequencer {
		return nil
	}

	if result.ProducedBatch != nil && result.ProducedBatch.Header != nil {
		// TODO - #718 - Unlink rollup production from L1 cadence.
		h.publishRollup(result.ProducedBatch.ToExtRollup())
		// TODO - #718 - Unlink batch production from L1 cadence.
		h.storeAndDistributeBatch(result.ProducedBatch)
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

		// node received a secret response, we should add them to our p2p connection pool
		if scrtRespTx, ok := t.(*ethadapter.L1RespondSecretTx); ok {
			err := h.processSharedSecretResponse(scrtRespTx)
			if err != nil {
				h.logger.Error("Failed to process shared secret response", log.ErrKey, err)
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

	h.logger.Trace("Sending transaction to publish rollup", "rollup_header",
		gethlog.Lazy{Fn: func() string {
			header, err := json.MarshalIndent(producedRollup.Header, "", "   ")
			if err != nil {
				return err.Error()
			}

			return string(header[:])
		}})

	rollupTx := h.mgmtContractLib.CreateRollup(tx, h.ethWallet.GetNonceAndIncrement())

	// fire-and-forget (track the receipt asynchronously)
	err = h.signAndBroadcastL1Tx(rollupTx, l1TxTriesRollup, false)
	if err != nil {
		h.logger.Error("could not broadcast rollup", log.ErrKey, err)
	}
}

// Creates a batch based on the rollup and distributes it to all other nodes.
func (h *host) storeAndDistributeBatch(producedBatch *common.ExtBatch) {
	err := h.db.AddBatchHeader(producedBatch)
	if err != nil {
		h.logger.Error("could not store batch", log.ErrKey, err)
	}

	batchMsg := hostcommon.BatchMsg{
		Batches:   []*common.ExtBatch{producedBatch},
		IsCatchUp: false,
	}
	err = h.p2p.BroadcastBatch(&batchMsg)
	if err != nil {
		h.logger.Error("could not broadcast batch", log.ErrKey, err)
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

	err = retry.Do(func() error {
		return h.ethClient.SendTransaction(signedTx)
	}, retry.NewDoublingBackoffStrategy(time.Second, tries)) // doubling retry wait (3 tries = 7sec, 7 tries = 63sec)
	if err != nil {
		return fmt.Errorf("broadcasting L1 transaction failed after %d tries. Cause: %w", tries, err)
	}
	h.logger.Trace("L1 transaction sent successfully, watching for receipt.")

	if awaitReceipt {
		// block until receipt is found and then return
		return h.waitForReceipt(signedTx.Hash())
	}

	// else just watch for receipt asynchronously and log if it fails
	go func() {
		// todo: consider how to handle the various ways that L1 transactions could fail to improve node operator QoL
		err := h.waitForReceipt(signedTx.Hash())
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
	h.logger.Trace("Successful L1 transaction receipt found.", "blk", receipt.BlockNumber, "blkHash", receipt.BlockHash)
	return nil
}

// This method implements the procedure by which a node obtains the secret
func (h *host) requestSecret() error {
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
	for _, scrtResponse := range scrtResponses {
		// todo: implement proper protocol so only one host responds to this secret requests initially
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
		// TODO review: l1tx.Sign(a.attestationPubKey) doesn't matter as the waitSecret will process a tx that was reverted
		respondSecretTx := h.mgmtContractLib.CreateRespondSecret(l1tx, h.ethWallet.GetNonceAndIncrement(), false)
		h.logger.Trace("Broadcasting secret response L1 tx.", "requester", scrtResponse.RequesterID)
		// fire-and-forget (track the receipt asynchronously)
		err := h.signAndBroadcastL1Tx(respondSecretTx, l1TxTriesSecret, false)
		if err != nil {
			return fmt.Errorf("could not broadcast secret response. Cause %w", err)
		}
	}
	return nil
}

// Whenever we receive a new shared secret response transaction, we update our list of P2P peers, as another aggregator
// may have joined the network.
func (h *host) processSharedSecretResponse(_ *ethadapter.L1RespondSecretTx) error {
	// We make a call to the L1 node to retrieve the new list of aggregators. An alternative would be to check that the
	// transaction succeeded, and if so, extract the additional host address from the transaction arguments. But we
	// believe this would be more brittle than just asking the L1 contract for its view of the current aggregators.
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

	// We filter down the list of retrieved addresses.
	var filteredHostAddresses []string //nolint:prealloc
	for _, hostAddress := range hostAddresses {
		// We exclude our own address.
		if hostAddress == h.config.P2PPublicAddress {
			continue
		}

		// We exclude any duplicate host addresses.
		isDup := false
		for _, existingHostAddress := range filteredHostAddresses {
			if hostAddress == existingHostAddress {
				isDup = true
				break
			}
		}
		if isDup {
			continue
		}

		filteredHostAddresses = append(filteredHostAddresses, hostAddress)
	}

	h.p2p.UpdatePeerList(filteredHostAddresses)
	return nil
}

// TODO: Perhaps extract only relevant logs. There were missing ones when requesting
// the logs filtered from geth.
func (h *host) extractReceipts(block *types.Block) types.Receipts {
	receipts := make(types.Receipts, 0)

	for _, transaction := range block.Transactions() {
		receipt, err := h.ethClient.TransactionReceipt(transaction.Hash())

		if err != nil || receipt == nil {
			h.logger.Error("Problem with retrieving the receipt on the host!", log.ErrKey, err, log.CmpKey, log.CrossChainCmp)
			continue
		}

		h.logger.Trace(fmt.Sprintf("Adding receipt[%d] for block %d, TX: %d",
			receipt.Status,
			common.ShortHash(block.Hash()),
			common.ShortHash(transaction.Hash())),
			log.CmpKey, log.CrossChainCmp)

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
			h.logger.Trace("checking block for secret resp", "height", blk.Number())
			if h.checkBlockForSecretResponse(blk) {
				return nil
			}

		case <-time.After(blockStreamWarningTimeout):
			// This will provide useful feedback if things are stuck (and in tests if any goroutines got stranded on this select)
			h.logger.Warn(fmt.Sprintf(" Waiting for secret from the L1. No blocks received for over %s", blockStreamWarningTimeout))

		case <-h.exitHostCh:
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
// a batch without its parent), we request the sequencer to resend the batches we've just received, plus any missing
// historical batches, then discard the received batches. We will store all of these at once when we receive them
func (h *host) handleBatches(encodedBatchMsg *common.EncodedBatchMsg) error {
	var batchMsg *hostcommon.BatchMsg
	err := rlp.DecodeBytes(*encodedBatchMsg, &batchMsg)
	if err != nil {
		return fmt.Errorf("could not decode batches using RLP. Cause: %w", err)
	}

	for _, batch := range batchMsg.Batches {
		// TODO - #718 - Consider moving to a model where the enclave manages the entire state, to avoid inconsistency.

		// If we do not have the block the rollup is tied to, we skip processing the batches for now. We'll catch them
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
		if err = h.enclaveClient.SubmitBatch(batch); err != nil {
			return fmt.Errorf("could not submit batch. Cause: %w", err)
		}
		if err = h.db.AddBatchHeader(batch); err != nil {
			return fmt.Errorf("could not store batch header. Cause: %w", err)
		}
	}

	return nil
}

func (h *host) handleBatchRequest(encodedBatchRequest *common.EncodedBatchRequest) error {
	var batchRequest *common.BatchRequest
	err := rlp.DecodeBytes(*encodedBatchRequest, &batchRequest)
	if err != nil {
		return fmt.Errorf("could not decode batch request using RLP. Cause: %w", err)
	}

	batches, err := h.batchManager.GetBatches(batchRequest)
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
