package host

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"sync/atomic"
	"time"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/go-obscuro/go/host/batchmanager"

	"github.com/ethereum/go-ethereum/rlp"

	hostcommon "github.com/obscuronet/go-obscuro/go/common/host"

	"github.com/obscuronet/go-obscuro/go/common/retry"

	gethlog "github.com/ethereum/go-ethereum/log"

	"github.com/obscuronet/go-obscuro/go/host/events"

	"github.com/obscuronet/go-obscuro/go/host/rpc/clientapi"

	"github.com/obscuronet/go-obscuro/go/host/db"
	"github.com/obscuronet/go-obscuro/go/host/rpc/clientrpc"

	"github.com/ethereum/go-ethereum/rpc"

	"github.com/obscuronet/go-obscuro/go/common/profiler"

	"github.com/obscuronet/go-obscuro/go/common/log"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/naoina/toml"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/config"
	"github.com/obscuronet/go-obscuro/go/ethadapter"
	"github.com/obscuronet/go-obscuro/go/ethadapter/mgmtcontractlib"
	"github.com/obscuronet/go-obscuro/go/wallet"
)

const (
	APIVersion1             = "1.0"
	APINamespaceObscuro     = "obscuro"
	APINamespaceEth         = "eth"
	apiNamespaceObscuroScan = "obscuroscan"
	apiNamespaceNetwork     = "net"
	apiNamespaceTest        = "test"

	// Attempts to broadcast the rollup transaction to the L1. Worst-case, equates to 7 seconds, plus time per request.
	l1TxTriesRollup = 3
	// Attempts to send secret initialisation, request or response transactions to the L1. Worst-case, equates to 63 seconds, plus time per request.
	l1TxTriesSecret = 7

	maxWaitForL1Receipt       = 100 * time.Second
	retryIntervalForL1Receipt = 10 * time.Second
	blockStreamWarningTimeout = 30 * time.Second
)

var (
	latestRollup *common.ExtRollup = nil
)

// Implementation of host.Host.
type host struct {
	config          config.HostConfig
	shortID         uint64
	isSequencer     bool
	genesisRequired bool

	p2p           hostcommon.P2P       // For communication with other Obscuro nodes
	ethClient     ethadapter.EthClient // For communication with the L1 node
	enclaveClient common.Enclave       // For communication with the enclave
	rpcServer     clientrpc.Server     // For communication with Obscuro client applications

	// control the host lifecycle
	exitHostCh            chan bool
	stopHostInterrupt     *int32
	bootstrappingComplete *int32 // Marks when the host is done bootstrapping

	l1BlockProvider hostcommon.ReconnectingBlockProvider
	txP2PCh         chan common.EncryptedTx         // The channel that new transactions from peers are sent to
	batchP2PCh      chan common.EncodedBatches      // The channel that new batches from peers are sent to
	batchRequestCh  chan common.EncodedBatchRequest // The channel that batch requests from peers are sent to

	db *db.DB // Stores the host's publicly-available data

	mgmtContractLib mgmtcontractlib.MgmtContractLib // Library to handle Management Contract lib operations
	ethWallet       wallet.Wallet                   // Wallet used to issue ethereum transactions
	logEventManager events.LogEventManager
	batchManager    *batchmanager.BatchManager

	logger gethlog.Logger
}

func NewHost(
	config config.HostConfig,
	p2p hostcommon.P2P,
	ethClient ethadapter.EthClient,
	enclaveClient common.Enclave,
	ethWallet wallet.Wallet,
	mgmtContractLib mgmtcontractlib.MgmtContractLib,
	logger gethlog.Logger,
) hostcommon.Host {
	database := db.NewInMemoryDB() // todo - make this config driven
	host := &host{
		// config
		config:      config,
		shortID:     common.ShortAddress(config.ID),
		isSequencer: config.IsGenesis && config.NodeType == common.Aggregator,

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
		batchP2PCh:      make(chan common.EncodedBatches),
		batchRequestCh:  make(chan common.EncodedBatchRequest),

		// Initialize the host DB
		db: database,

		mgmtContractLib: mgmtContractLib, // library that provides a handler for Management Contract
		ethWallet:       ethWallet,       // the host's ethereum wallet
		logEventManager: events.NewLogEventManager(logger),
		batchManager:    batchmanager.NewBatchManager(database),

		logger: logger,
	}

	if config.HasClientRPCHTTP || config.HasClientRPCWebsockets {
		rpcAPIs := []rpc.API{
			{
				Namespace: APINamespaceObscuro,
				Version:   APIVersion1,
				Service:   clientapi.NewObscuroAPI(host),
				Public:    true,
			},
			{
				Namespace: APINamespaceEth,
				Version:   APIVersion1,
				Service:   clientapi.NewEthereumAPI(host, logger),
				Public:    true,
			},
			{
				Namespace: apiNamespaceObscuroScan,
				Version:   APIVersion1,
				Service:   clientapi.NewObscuroScanAPI(host),
				Public:    true,
			},
			{
				Namespace: apiNamespaceNetwork,
				Version:   APIVersion1,
				Service:   clientapi.NewNetworkAPI(host),
				Public:    true,
			},
			{
				Namespace: apiNamespaceTest,
				Version:   APIVersion1,
				Service:   clientapi.NewTestAPI(host),
				Public:    true,
			},
			{
				Namespace: APINamespaceEth,
				Version:   APIVersion1,
				Service:   clientapi.NewFilterAPI(host, logger),
				Public:    true,
			},
		}
		host.rpcServer = clientrpc.NewServer(config, rpcAPIs, logger)
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

func (h *host) Start() {
	h.validateConfig()

	tomlConfig, err := toml.Marshal(h.config)
	if err != nil {
		h.logger.Crit("could not print host config")
	}
	h.logger.Info("Host started with following config", log.CfgKey, string(tomlConfig))

	// wait for the Enclave to be available
	h.waitForEnclave()

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

	if h.config.IsGenesis {
		// the genesis node gets a flag set on it until it has published the genesis block
		// todo: handle genesis separately as an initial step for the genesis node
		h.genesisRequired = true
	}
	// start the obscuro RPC endpoints
	if h.rpcServer != nil {
		h.rpcServer.Start()
		h.logger.Info("Started client server.")
	}

	// start the host's main processing loop
	h.startProcessing()
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
	err = h.signAndBroadcastL1Tx(initialiseSecretTx, l1TxTriesSecret)
	if err != nil {
		return fmt.Errorf("failed to initialise enclave secret. Cause: %w", err)
	}
	h.logger.Info("Node is genesis node. Secret was broadcast.")
	return nil
}

func (h *host) Config() *config.HostConfig {
	return &h.config
}

func (h *host) DB() *db.DB {
	return h.db
}

func (h *host) EnclaveClient() common.Enclave {
	return h.enclaveClient
}

func (h *host) SubmitAndBroadcastTx(encryptedParams common.EncryptedParamsSendRawTx) (common.EncryptedResponseSendRawTx, error) {
	encryptedTx := common.EncryptedTx(encryptedParams)
	encryptedResponse, err := h.enclaveClient.SubmitTx(encryptedTx)
	if err != nil {
		return nil, fmt.Errorf("could not submit transaction. Cause: %w", err)
	}

	err = h.p2p.BroadcastTx(encryptedTx)
	if err != nil {
		return nil, fmt.Errorf("could not broadcast transaction. Cause: %w", err)
	}

	return encryptedResponse, nil
}

func (h *host) ReceiveTx(tx common.EncryptedTx) {
	h.txP2PCh <- tx
}

func (h *host) ReceiveBatches(batches common.EncodedBatches) {
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
	if h.rpcServer != nil {
		// We cannot stop the RPC server synchronously. This is because the host itself is being stopped by an RPC
		// call, so there is a deadlock. The RPC server is waiting for all connections to close, but a single
		// connection remains open, waiting for the RPC server to close.
		go h.rpcServer.Stop()
	}

	// Leave some time for all processing to finish before exiting the main loop.
	time.Sleep(time.Second)
	h.exitHostCh <- true

	h.logger.Info("Host shut down successfully.")
}

// HealthCheck returns whether the host, enclave and DB are healthy
func (h *host) HealthCheck() (bool, error) {
	// check the enclave health, which in turn checks the DB health
	enclaveHealthy, err := h.enclaveClient.HealthCheck()
	if err != nil {
		// simplest iteration, log the error and just return that it's not healthy
		h.logger.Error("unable to HealthCheck enclave", "err", err)
		return false, nil
	}
	// TODO host healthcheck operations
	hostHealthy := true
	return enclaveHealthy && hostHealthy, nil
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

		case batches := <-h.batchP2PCh:
			// todo: discard p2p messages if enclave won't be able to make use of them (e.g. we're way behind L1 head)
			if err := h.handleBatches(&batches); err != nil {
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
	result, err := h.enclaveClient.SubmitL1Block(*block, isLatestBlock, h.isSequencer)
	if err != nil {
		return fmt.Errorf("did not ingest block b_%d. Cause: %w", common.ShortHash(block.Hash()), err)
	}
	if h.shortID == 0 {
		println(fmt.Sprintf("jjj back at node. Is updated? %t; new rollup? %t; isLatestBlock? %t",
			result.IngestedRollupHeader != nil, result.ProducedRollup != nil, isLatestBlock))
	}
	err = h.storeBlockProcessingResult(result, block.Header())
	if err != nil {
		return fmt.Errorf("submitted block to enclave but could not store the block processing result. Cause: %w", err)
	}

	h.logEventManager.SendLogsToSubscribers(result)

	err = h.publishSharedSecretResponses(result.ProducedSecretResponses)
	if err != nil {
		h.logger.Error("failed to publish response to secret request", log.ErrKey, err)
	}

	// If we're the sequencer, and we're processing the latest block, we produce, publish and distribute a new rollup.
	if !h.isSequencer {
		return nil
	}

	if h.genesisRequired {
		if err = h.initialiseProtocol(block); err != nil {
			h.logger.Crit("Could not initialise protocol.", log.ErrKey, err)
		}
		h.genesisRequired = false
		return nil // nothing further to process since network had no genesis
	}

	if result.ProducedRollup.Header != nil {
		// TODO - #718 - Unlink rollup production from L1 cadence.
		h.publishRollup(result.ProducedRollup)
		// TODO - #718 - Unlink batch production from L1 cadence.
		h.storeAndDistributeBatch(result.ProducedRollup)
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

	rollupTx := h.mgmtContractLib.CreateRollup(tx, h.ethWallet.GetNonceAndIncrement())
	err = h.signAndBroadcastL1Tx(rollupTx, l1TxTriesRollup)
	if err != nil {
		h.logger.Error("could not broadcast rollup", log.ErrKey, err)
	}
}

// Creates a batch based on the rollup and distributes it to all other nodes.
func (h *host) storeAndDistributeBatch(producedRollup *common.ExtRollup) {
	batch := common.ExtBatch{
		Header:          producedRollup.Header,
		TxHashes:        producedRollup.TxHashes,
		EncryptedTxBlob: producedRollup.EncryptedTxBlob,
	}

	println(fmt.Sprintf("jjj sequencer distributing batch %d (hash: %s, parent hash: %s, block: %s)", batch.Header.Number, batch.Hash(), batch.Header.ParentHash, batch.Header.L1Proof))

	err := h.db.AddBatchHeader(&batch)
	if err != nil {
		h.logger.Error("could not store batch", log.ErrKey, err)
	}

	err = h.p2p.BroadcastBatch(&batch)
	if err != nil {
		h.logger.Error("could not broadcast batch", log.ErrKey, err)
	}
}

func (h *host) storeBlockProcessingResult(result *common.BlockSubmissionResponse, blockHeader *types.Header) error {
	// only update the host rollup headers if the enclave has found a new rollup head
	if result.IngestedRollupHeader != nil {
		// adding a header will update the head if it has a higher height
		err := h.db.AddRollupHeader(result.IngestedRollupHeader)
		if err != nil {
			return err
		}
	}

	// adding a header will update the head if it has a higher height
	return h.db.AddBlockHeader(blockHeader)
}

// Called only by the first enclave to bootstrap the network

func (h *host) initialiseProtocol(block *types.Block) error {
	// Create the genesis rollup
	genesisRollup, err := h.enclaveClient.ProduceGenesis(block.Hash())
	if err != nil {
		return fmt.Errorf("could not produce genesis. Cause: %w", err)
	}
	h.logger.Info(
		fmt.Sprintf("Initialising network. Genesis rollup r_%d.",
			common.ShortHash(genesisRollup.Header.Hash()),
		))

	// Distribute the corresponding batch.
	h.storeAndDistributeBatch(genesisRollup)

	// Submit the rollup to the management contract.
	encodedRollup, err := common.EncodeRollup(genesisRollup)
	if err != nil {
		return fmt.Errorf("could not encode rollup. Cause: %w", err)
	}
	l1tx := &ethadapter.L1RollupTx{
		Rollup: encodedRollup,
	}
	rollupTx := h.mgmtContractLib.CreateRollup(l1tx, h.ethWallet.GetNonceAndIncrement())
	err = h.signAndBroadcastL1Tx(rollupTx, l1TxTriesRollup)
	if err != nil {
		return fmt.Errorf("could not initialise protocol. Cause: %w", err)
	}

	return nil
}

// `tries` is the number of times to attempt broadcasting the transaction.
func (h *host) signAndBroadcastL1Tx(tx types.TxData, tries uint64) error {
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

	// asynchronously watch for a successful receipt
	// todo: consider how to handle the various ways that L1 transactions could fail to improve node operator QoL
	go h.watchForReceipt(signedTx.Hash())

	return nil
}

func (h *host) watchForReceipt(txHash common.TxHash) {
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
		h.logger.Error("receipt for L1 transaction never found despite 'successful' broadcast",
			"err", err, "signer", h.ethWallet.Address().Hex(),
		)
		return
	}

	if err == nil && receipt.Status != types.ReceiptStatusSuccessful {
		h.logger.Error("unsuccessful receipt found for published L1 transaction",
			"status", receipt.Status,
			"signer", h.ethWallet.Address().Hex())
	}
	h.logger.Trace("Successful L1 transaction receipt found.", "blk", receipt.BlockNumber, "blkHash", receipt.BlockHash)
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
	err = h.signAndBroadcastL1Tx(requestSecretTx, l1TxTriesSecret)
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
		err := h.signAndBroadcastL1Tx(respondSecretTx, l1TxTriesSecret)
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
func (h *host) handleBatches(encodedBatches *common.EncodedBatches) error {
	var batches []*common.ExtBatch
	err := rlp.DecodeBytes(*encodedBatches, &batches)
	if err != nil {
		return fmt.Errorf("could not decode batches using RLP. Cause: %w", err)
	}
	if len(batches) == 0 {
		return nil
	}

	// We store the batches.
	err = h.batchManager.StoreBatches(batches, h.shortID)
	if err != nil {
		if !errors.Is(err, batchmanager.ErrBatchesMissing) {
			return fmt.Errorf("could not store batches. Cause: %w", err)
		}

		// We have encountered missing batches. We abort the storage operation and request the missing batches.
		batchRequest, err := h.batchManager.CreateBatchRequest(h.config.P2PPublicAddress)
		if err != nil {
			return fmt.Errorf("could not create batch request. Cause: %w", err)
		}
		if err = h.p2p.RequestBatches(batchRequest); err != nil {
			return fmt.Errorf("could not request historical batches. Cause: %w", err)
		}

		// If we requested any batches, we return early and wait for the missing batches to arrive.
		return nil
	}

	// TODO - #718 - We should probably submit each batch after storing it, and not submitting each one only if *every*
	//  batch was stored correctly.
	for _, batch := range batches {
		if err = h.enclaveClient.SubmitBatch(batch); err != nil {
			return fmt.Errorf("could not submit batch. Cause: %w", err)
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

	return h.p2p.SendBatches(batches, batchRequest.Requester)
}

// Checks the host config is valid.
func (h *host) validateConfig() {
	if h.config.IsGenesis && h.config.NodeType != common.Aggregator {
		h.logger.Crit("genesis node must be an aggregator")
	}
	if !h.config.IsGenesis && h.config.NodeType == common.Aggregator {
		h.logger.Crit("only the genesis node can be an aggregator")
	}

	if h.config.P2PPublicAddress == "" {
		h.logger.Crit("the host must specify a public P2P address")
	}
}
