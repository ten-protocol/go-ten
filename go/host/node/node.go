package node

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/google/uuid"

	"github.com/obscuronet/go-obscuro/go/host/rpc/clientapi"

	"github.com/obscuronet/go-obscuro/go/host"

	"github.com/obscuronet/go-obscuro/go/host/db"
	"github.com/obscuronet/go-obscuro/go/host/rpc/clientrpc"

	"github.com/ethereum/go-ethereum/rpc"

	"github.com/obscuronet/go-obscuro/go/common/profiler"

	"github.com/obscuronet/go-obscuro/go/common/log"

	"github.com/ethereum/go-ethereum"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/naoina/toml"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/config"
	"github.com/obscuronet/go-obscuro/go/ethadapter"
	"github.com/obscuronet/go-obscuro/go/ethadapter/mgmtcontractlib"
	"github.com/obscuronet/go-obscuro/go/wallet"
)

const (
	apiVersion1             = "1.0"
	apiNamespaceObscuro     = "obscuro"
	apiNamespaceEthereum    = "eth"
	apiNamespaceObscuroScan = "obscuroscan"
	apiNamespaceNetwork     = "net"
	apiNamespaceTest        = "test"

	// Attempts to broadcast the rollup transaction to the L1. Worst-case, equates to 7 seconds, plus time per request.
	l1TxTriesRollup = 3
	// Attempts to send secret initialisation, request or response transactions to the L1. Worst-case, equates to 63 seconds, plus time per request.
	l1TxTriesSecret = 7
)

// Node is an implementation of host.Host.
type Node struct {
	config  config.HostConfig
	shortID uint64

	p2p           host.P2P             // For communication with other Obscuro nodes
	ethClient     ethadapter.EthClient // For communication with the L1 node
	enclaveClient common.Enclave       // For communication with the enclave
	rpcServer     clientrpc.Server     // For communication with Obscuro client applications

	stats host.StatsCollector

	// control the host lifecycle
	exitNodeCh            chan bool
	stopNodeInterrupt     *int32
	bootstrappingComplete *int32 // Marks when the node is done bootstrapping

	blockRPCCh   chan blockAndParent        // The channel that new blocks from the L1 node are sent to
	forkRPCCh    chan []common.EncodedBlock // The channel that new forks from the L1 node are sent to
	rollupsP2PCh chan common.EncodedRollup  // The channel that new rollups from peers are sent to
	txP2PCh      chan common.EncryptedTx    // The channel that new transactions from peers are sent to
	logsCh       chan []*types.Log          // The channel that logs are sent to

	nodeDB *db.DB // Stores the node's publicly-available data

	// library to handle Management Contract lib operations
	mgmtContractLib mgmtcontractlib.MgmtContractLib

	// Wallet used to issue ethereum transactions
	ethWallet wallet.Wallet
}

func NewHost(
	config config.HostConfig,
	stats host.StatsCollector,
	p2p host.P2P,
	ethClient ethadapter.EthClient,
	enclaveClient common.Enclave,
	ethWallet wallet.Wallet,
	mgmtContractLib mgmtcontractlib.MgmtContractLib,
) host.MockHost {
	node := &Node{
		// config
		config:  config,
		shortID: common.ShortAddress(config.ID),

		// Communication layers.
		p2p:           p2p,
		ethClient:     ethClient,
		enclaveClient: enclaveClient,

		// statistics and metrics
		stats: stats,

		// lifecycle channels
		exitNodeCh:            make(chan bool),
		stopNodeInterrupt:     new(int32),
		bootstrappingComplete: new(int32),

		// incoming data
		blockRPCCh:   make(chan blockAndParent),
		forkRPCCh:    make(chan []common.EncodedBlock),
		rollupsP2PCh: make(chan common.EncodedRollup),
		txP2PCh:      make(chan common.EncryptedTx),
		logsCh:       make(chan []*types.Log),

		// Initialize the node DB
		// nodeDB:       NewLevelDBBackedDB(), // todo - make this config driven
		nodeDB: db.NewInMemoryDB(),

		// library that provides a handler for Management Contract
		mgmtContractLib: mgmtContractLib,
		// the nodes ethereum wallet
		ethWallet: ethWallet,
	}

	if config.HasClientRPCHTTP || config.HasClientRPCWebsockets {
		rpcAPIs := []rpc.API{
			{
				Namespace: apiNamespaceObscuro,
				Version:   apiVersion1,
				Service:   clientapi.NewObscuroAPI(node),
				Public:    true,
			},
			{
				Namespace: apiNamespaceEthereum,
				Version:   apiVersion1,
				Service:   clientapi.NewEthereumAPI(node),
				Public:    true,
			},
			{
				Namespace: apiNamespaceObscuroScan,
				Version:   apiVersion1,
				Service:   clientapi.NewObscuroScanAPI(node),
				Public:    true,
			},
			{
				Namespace: apiNamespaceNetwork,
				Version:   apiVersion1,
				Service:   clientapi.NewNetworkAPI(node),
				Public:    true,
			},
			{
				Namespace: apiNamespaceTest,
				Version:   apiVersion1,
				Service:   clientapi.NewTestAPI(node),
				Public:    true,
			},
			{
				Namespace: "eth",
				Version:   apiVersion1,
				Service:   clientapi.NewFilterAPI(node, node.logsCh),
				Public:    true,
			},
		}
		node.rpcServer = clientrpc.NewServer(config, rpcAPIs)
	}

	var prof *profiler.Profiler
	if config.ProfilerEnabled {
		prof = profiler.NewProfiler(profiler.DefaultHostPort)
		err := prof.Start()
		if err != nil {
			log.Panic("unable to start the profiler: %s", err)
		}
	}

	jsonConfig, _ := json.MarshalIndent(config, "", "  ")
	log.Info("Host service created with following config:\n%s", string(jsonConfig))

	return node
}

func (a *Node) Start() {
	tomlConfig, err := toml.Marshal(a.config)
	if err != nil {
		log.Panic("could not print host config")
	}
	common.LogWithID(a.shortID, "Host started with following config:\n%s", tomlConfig)

	// wait for the Enclave to be available
	a.waitForEnclave()

	// todo: we should try to recover the key from a previous run of the node here? Before generating or requesting the key.
	if a.config.IsGenesis {
		err = a.broadcastSecret()
		if err != nil {
			log.Panic(err.Error())
		}
	} else {
		err = a.requestSecret()
		if err != nil {
			log.Panic(err.Error())
		}
	}

	// attach the l1 monitor
	go a.monitorBlocks()

	// bootstrap the node
	latestBlock := a.bootstrapNode()

	// start the enclave speculative work from last block
	a.enclaveClient.Start(latestBlock)

	if a.config.IsGenesis {
		_, err = a.initialiseProtocol(&latestBlock)
		if err != nil {
			log.Panic(err.Error())
		}
	}
	// start the obscuro RPC endpoints
	if a.rpcServer != nil {
		a.rpcServer.Start()
		common.LogWithID(a.shortID, "Started client server.")
	}

	// start the node main processing loop
	a.startProcessing()
}

func (a *Node) broadcastSecret() error {
	common.LogWithID(a.shortID, "Node is genesis node. Broadcasting secret.")
	// Create the shared secret and submit it to the management contract for storage
	attestation := a.enclaveClient.Attestation()
	if attestation.Owner != a.config.ID {
		return fmt.Errorf("genesis node has ID %s, but its enclave produced an attestation using ID %s", a.config.ID.Hex(), attestation.Owner.Hex())
	}

	encodedAttestation, err := common.EncodeAttestation(attestation)
	if err != nil {
		return fmt.Errorf("could not encode attestation Cause: %w", err)
	}
	l1tx := &ethadapter.L1InitializeSecretTx{
		AggregatorID:  &a.config.ID,
		Attestation:   encodedAttestation,
		InitialSecret: a.enclaveClient.GenerateSecret(),
		HostAddress:   a.config.P2PPublicAddress,
	}
	initialiseSecretTx := a.mgmtContractLib.CreateInitializeSecret(l1tx, a.ethWallet.GetNonceAndIncrement())
	err = a.signAndBroadcastTx(initialiseSecretTx, l1TxTriesSecret)
	if err != nil {
		return fmt.Errorf("failed to initialise enclave secret. Cause: %w", err)
	}
	common.LogWithID(a.shortID, "Node is genesis node. Secret was broadcast.")
	return nil
}

func (a *Node) Config() *config.HostConfig {
	return &a.config
}

func (a *Node) DB() *db.DB {
	return a.nodeDB
}

func (a *Node) EnclaveClient() common.Enclave {
	return a.enclaveClient
}

func (a *Node) MockedNewHead(b common.EncodedBlock, p common.EncodedBlock) {
	if atomic.LoadInt32(a.stopNodeInterrupt) == 1 {
		return
	}
	a.blockRPCCh <- blockAndParent{b, p}
}

func (a *Node) MockedNewFork(b []common.EncodedBlock) {
	if atomic.LoadInt32(a.stopNodeInterrupt) == 1 {
		return
	}
	a.forkRPCCh <- b
}

func (a *Node) SubmitAndBroadcastTx(encryptedParams common.EncryptedParamsSendRawTx) (common.EncryptedResponseSendRawTx, error) {
	encryptedTx := common.EncryptedTx(encryptedParams)
	encryptedResponse, err := a.enclaveClient.SubmitTx(encryptedTx)
	if err != nil {
		common.LogWithID(a.shortID, "Could not submit transaction: %s", err)
		return nil, err
	}

	err = a.p2p.BroadcastTx(encryptedTx)
	if err != nil {
		return nil, fmt.Errorf("could not broadcast transaction. Cause: %w", err)
	}

	return encryptedResponse, nil
}

// ReceiveRollup is called by counterparties when there is a Rollup to broadcast
// All it does is forward the rollup for processing to the enclave
func (a *Node) ReceiveRollup(r common.EncodedRollup) {
	if atomic.LoadInt32(a.stopNodeInterrupt) == 1 {
		return
	}
	a.rollupsP2PCh <- r
}

// ReceiveTx receives a new transaction
func (a *Node) ReceiveTx(tx common.EncryptedTx) {
	if atomic.LoadInt32(a.stopNodeInterrupt) == 1 {
		return
	}
	a.txP2PCh <- tx
}

// CreateSubscription sets up a subscription between the host and the enclave.
func (a *Node) CreateSubscription() error {
	id, err := uuid.NewUUID()
	if err != nil {
		return fmt.Errorf("could not generate new UUID for subscription")
	}

	// TODO - #453 - Link subscriptions to specific accounts.
	subscription := common.LogSubscription{
		Accounts: []*common.SubscriptionAccount{},
	}

	// TODO - #453 - Encrypt subscription requests, rather than just encoding them.
	encodedSubscription, err := rlp.EncodeToBytes(subscription)
	if err != nil {
		return fmt.Errorf("could not encrypt subscription request")
	}

	err = a.EnclaveClient().Subscribe(id, encodedSubscription)
	if err != nil {
		return fmt.Errorf("could not create subscription with enclave")
	}

	return nil
}

func (a *Node) Stop() {
	// block all requests
	atomic.StoreInt32(a.stopNodeInterrupt, 1)

	if err := a.p2p.StopListening(); err != nil {
		common.ErrorWithID(a.shortID, "failed to close transaction P2P listener cleanly: %s", err)
	}
	if err := a.enclaveClient.Stop(); err != nil {
		common.ErrorWithID(a.shortID, "could not stop enclave server. Cause: %s", err)
	}
	if err := a.enclaveClient.StopClient(); err != nil {
		common.ErrorWithID(a.shortID, "failed to stop enclave RPC client. Cause: %s", err)
	}
	if a.rpcServer != nil {
		// We cannot stop the RPC server synchronously. This is because the host itself is being stopped by an RPC
		// call, so there is a deadlock. The RPC server is waiting for all connections to close, but a single
		// connection remains open, waiting for the RPC server to close.
		go a.rpcServer.Stop()
	}

	// Leave some time for all processing to finish before exiting the main loop.
	time.Sleep(time.Second)
	a.exitNodeCh <- true

	common.LogWithID(a.shortID, "Node shut down successfully.")
}

// Waits for enclave to be available, printing a wait message every two seconds.
func (a *Node) waitForEnclave() {
	counter := 0
	for err := a.enclaveClient.IsReady(); err != nil; {
		if counter >= 20 {
			common.LogWithID(a.shortID, "Waiting for enclave on %s. Latest connection attempt failed with: %v", a.config.EnclaveRPCAddress, err)
			counter = 0
		}

		time.Sleep(100 * time.Millisecond)
		counter++
	}
	common.LogWithID(a.shortID, "Connected to enclave service.")
}

// starts the host main processing loop
func (a *Node) startProcessing() {
	//	time.Sleep(time.Second)
	// Only open the p2p connection when the node is fully initialised
	a.p2p.StartListening(a)

	// use the roundInterrupt as a signaling mechanism for interrupting block processing
	// stops processing the current round if a new block arrives
	i := int32(0)
	roundInterrupt := &i

	// Main Processing Loop -
	// - Process new blocks from the L1 node
	// - Process new Rollups gossiped from L2 Peers
	// - Process new Transactions gossiped from L2 Peers
	for {
		select {
		case b := <-a.blockRPCCh:
			roundInterrupt = triggerInterrupt(roundInterrupt)
			err := a.processBlocks([]common.EncodedBlock{b.p, b.b}, roundInterrupt)
			if err != nil {
				common.WarnWithID(a.shortID, "Could not process block received via RPC. Cause: %v", err)
			}

		case f := <-a.forkRPCCh:
			roundInterrupt = triggerInterrupt(roundInterrupt)
			err := a.processBlocks(f, roundInterrupt)
			if err != nil {
				common.WarnWithID(a.shortID, "Could not process fork received via RPC. Cause: %v", err)
			}

		case r := <-a.rollupsP2PCh:
			rol, err := common.DecodeRollup(r)
			common.TraceWithID(a.shortID, "Received rollup: r_%d(%d) parent: r_%d from A%d",
				common.ShortHash(rol.Hash()),
				rol.Header.Number,
				common.ShortHash(rol.Header.ParentHash),
				common.ShortAddress(rol.Header.Agg),
			)
			if err != nil {
				common.WarnWithID(a.shortID, "Could not check enclave initialisation. Cause: %v", err)
			}

			go a.enclaveClient.SubmitRollup(common.ExtRollup{
				Header:          rol.Header,
				TxHashes:        rol.TxHashes,
				EncryptedTxBlob: rol.Transactions,
			})

		case tx := <-a.txP2PCh:
			if _, err := a.enclaveClient.SubmitTx(tx); err != nil {
				common.WarnWithID(a.shortID, "Could not submit transaction. Cause: %s", err)
			}

		case <-a.exitNodeCh:
			return
		}
	}
}

// activates the given interrupt (atomically) and returns a new interrupt
func triggerInterrupt(interrupt *int32) *int32 {
	// Notify the previous round to stop work
	atomic.StoreInt32(interrupt, 1)
	i := int32(0)
	return &i
}

type blockAndParent struct {
	b common.EncodedBlock
	p common.EncodedBlock
}

func (a *Node) processBlocks(blocks []common.EncodedBlock, interrupt *int32) error {
	var result common.BlockSubmissionResponse
	for _, block := range blocks {
		// For the genesis block the parent is nil
		if block == nil {
			continue
		}

		decoded, err := block.DecodeBlock()
		if err != nil {
			return err
		}
		a.processBlock(decoded)

		// submit each block to the enclave for ingestion plus validation
		result, err = a.enclaveClient.SubmitBlock(*decoded)
		if err != nil {
			return err
		}

		a.storeBlockProcessingResult(result)
		a.sendLogsToSubscribers(result)
	}

	if !result.IngestedBlock {
		b, err := blocks[len(blocks)-1].DecodeBlock()
		if err != nil {
			return fmt.Errorf("did not ingest block. Cause: %s", result.BlockNotIngestedCause)
		}
		return fmt.Errorf("did not ingest block b_%d. Cause: %s", common.ShortHash(b.Hash()), result.BlockNotIngestedCause)
	}

	// Nodes can start before the genesis was published, and it makes no sense to enter the protocol.
	if result.ProducedRollup.Header == nil {
		return nil
	}

	encodedRollup, err := common.EncodeRollup(result.ProducedRollup.ToRollup())
	if err != nil {
		return fmt.Errorf("could not encode rollup. Cause: %w", err)
	}
	err = a.p2p.BroadcastRollup(encodedRollup)
	if err != nil {
		return fmt.Errorf("could not broadcast rollup. Cause: %w", err)
	}

	common.ScheduleInterrupt(a.config.GossipRoundDuration, interrupt, a.handleRoundWinner(result))
	return nil
}

// Looks at each transaction in the block, and kicks off special handling for the transaction if needed.
func (a *Node) processBlock(b *types.Block) {
	for _, tx := range b.Transactions() {
		t := a.mgmtContractLib.DecodeTx(tx)
		if t == nil {
			continue
		}

		if scrtReqTx, ok := t.(*ethadapter.L1RequestSecretTx); ok {
			common.LogWithID(a.shortID, "Process shared secret request. Block: %d. Tx: %d", b.NumberU64(), common.ShortHash(tx.Hash()))
			err := a.processSharedSecretRequest(scrtReqTx)
			if err != nil {
				common.ErrorWithID(a.shortID, "Failed to process shared secret request. Cause: %s", err)
				continue
			}
		}

		if scrtRespTx, ok := t.(*ethadapter.L1RespondSecretTx); ok {
			err := a.processSharedSecretResponse(scrtRespTx)
			if err != nil {
				common.ErrorWithID(a.shortID, "Failed to process shared secret response. Cause: %s", err)
				continue
			}
		}

		if initSecretTx, ok := t.(*ethadapter.L1InitializeSecretTx); ok {
			// TODO - Ensure that we don't accidentally skip over the real `L1InitializeSecretTx` message. Otherwise
			//  our node will never be able to speak to other nodes.
			// there must be a way to make sure that this transaction can only be sent once.
			att, err := common.DecodeAttestation(initSecretTx.Attestation)
			if err != nil {
				common.ErrorWithID(a.shortID, "Could not decode attestation report. Cause: %s", err)
			}
			err = a.enclaveClient.StoreAttestation(att)
			if err != nil {
				common.ErrorWithID(a.shortID, "Could not store the attestation report. Cause: %s", err)
			}
		}
	}
}

func (a *Node) handleRoundWinner(result common.BlockSubmissionResponse) func() {
	return func() {
		if atomic.LoadInt32(a.stopNodeInterrupt) == 1 {
			return
		}
		// Request the round winner for the current head
		winnerRollup, isWinner, err := a.enclaveClient.RoundWinner(result.ProducedRollup.Header.ParentHash)
		if err != nil {
			log.Panic("could not determine round winner. Cause: %s", err)
		}
		if isWinner {
			common.LogWithID(a.shortID, "Winner (b_%d) r_%d(%d).",
				common.ShortHash(result.BlockHeader.Hash()),
				common.ShortHash(winnerRollup.Header.Hash()),
				winnerRollup.Header.Number,
			)

			encodedRollup, err := common.EncodeRollup(winnerRollup.ToRollup())
			if err != nil {
				log.Panic("could not encode rollup. Cause: %s", err)
			}
			tx := &ethadapter.L1RollupTx{
				Rollup: encodedRollup,
			}

			// That handler can get called multiple times for the same height. And it will return the same winner rollup.
			// In case the winning rollup belongs to the current enclave it will be submitted again, which is inefficient.
			if !a.nodeDB.WasSubmitted(winnerRollup.Header.Hash()) {
				rollupTx := a.mgmtContractLib.CreateRollup(tx, a.ethWallet.GetNonceAndIncrement())
				err = a.signAndBroadcastTx(rollupTx, l1TxTriesRollup)
				if err != nil {
					log.Error("could not broadcast winning rollup. Cause: %s", err)
				}
				a.nodeDB.AddSubmittedRollup(winnerRollup.Header.Hash())
			}
		}
	}
}

func (a *Node) storeBlockProcessingResult(result common.BlockSubmissionResponse) {
	// only update the node rollup headers if the enclave has found a new rollup head
	if result.FoundNewHead {
		// adding a header will update the head if it has a higher height
		headerWithHashes := common.HeaderWithTxHashes{Header: result.RollupHead, TxHashes: result.ProducedRollup.TxHashes}
		a.nodeDB.AddRollupHeader(&headerWithHashes)
	}

	// adding a header will update the head if it has a higher height
	if result.IngestedBlock {
		a.nodeDB.AddBlockHeader(result.BlockHeader)
	}
}

// Distributes logs to subscribed clients.
// TODO - #453 - Encrypt logs, rather than just serialising them as JSON.
// TODO - #453 - Distribute logs specifically based on subscription IDs, rather than sending all logs to everyone.
func (a *Node) sendLogsToSubscribers(result common.BlockSubmissionResponse) {
	for _, jsonLogs := range result.SubscribedLogs {
		var logs []*types.Log
		err := json.Unmarshal(jsonLogs, &logs)
		if err != nil {
			log.Error("could not send logs to subscribers as could not unmarshal logs from JSON. Cause: %s", err)
		}
		a.logsCh <- logs
	}
}

// Called only by the first enclave to bootstrap the network
func (a *Node) initialiseProtocol(block *types.Block) (common.L2RootHash, error) {
	// Create the genesis rollup and submit it to the management contract
	genesisResponse := a.enclaveClient.ProduceGenesis(block.Hash())
	common.LogWithID(
		a.shortID,
		"Initialising network. Genesis rollup r_%d.",
		common.ShortHash(genesisResponse.ProducedRollup.Header.Hash()),
	)
	encodedRollup, err := common.EncodeRollup(genesisResponse.ProducedRollup.ToRollup())
	if err != nil {
		return common.L2RootHash{}, fmt.Errorf("could not encode rollup. Cause: %w", err)
	}
	l1tx := &ethadapter.L1RollupTx{
		Rollup: encodedRollup,
	}

	rollupTx := a.mgmtContractLib.CreateRollup(l1tx, a.ethWallet.GetNonceAndIncrement())
	err = a.signAndBroadcastTx(rollupTx, l1TxTriesRollup)
	if err != nil {
		return common.L2RootHash{}, fmt.Errorf("could not initialise protocol. Cause: %w", err)
	}

	return genesisResponse.ProducedRollup.Header.ParentHash, nil
}

// `tries` is the number of times to attempt broadcasting the transaction.
func (a *Node) signAndBroadcastTx(tx types.TxData, tries int) error {
	signedTx, err := a.ethWallet.SignTransaction(tx)
	if err != nil {
		return err
	}

	funcBroadcastTx := func() error { return a.ethClient.SendTransaction(signedTx) }
	err = retryWithBackoff(tries, funcBroadcastTx)
	if err == nil {
		return nil
	}
	return fmt.Errorf("broadcasting L1 transaction failed after %d tries. Cause: %w", tries, err)
}

// This method implements the procedure by which a node obtains the secret
func (a *Node) requestSecret() error {
	common.LogWithID(a.shortID, "Requesting secret.")
	att := a.enclaveClient.Attestation()
	if att.Owner != a.config.ID {
		return fmt.Errorf("node has ID %s, but its enclave produced an attestation using ID %s", a.config.ID.Hex(), att.Owner.Hex())
	}
	encodedAttestation, err := common.EncodeAttestation(att)
	if err != nil {
		return fmt.Errorf("could not encode attestation. Cause: %w", err)
	}
	l1tx := &ethadapter.L1RequestSecretTx{
		Attestation: encodedAttestation,
	}
	requestSecretTx := a.mgmtContractLib.CreateRequestSecret(l1tx, a.ethWallet.GetNonceAndIncrement())
	err = a.signAndBroadcastTx(requestSecretTx, l1TxTriesSecret)
	if err != nil {
		return err
	}

	err = a.awaitSecret()
	if err != nil {
		log.Panic(err.Error())
	}
	return nil
}

func (a *Node) handleStoreSecretTx(t *ethadapter.L1RespondSecretTx) bool {
	if t.RequesterID.Hex() != a.config.ID.Hex() {
		// this secret is for somebody else
		return false
	}

	// someone has replied for us
	err := a.enclaveClient.InitEnclave(t.Secret)
	if err != nil {
		common.LogWithID(a.shortID, "Failed to initialise enclave with received secret. Err: %s", err)
		return false
	}
	return true
}

func (a *Node) processSharedSecretRequest(scrtReqTx *ethadapter.L1RequestSecretTx) error {
	att, err := common.DecodeAttestation(scrtReqTx.Attestation)
	if err != nil {
		common.LogWithID(a.shortID, "Failed to decode attestation. %s", err)
		return nil
	}

	jsonAttestation, err := json.Marshal(att)
	if err == nil {
		common.LogWithID(a.shortID, "Received attestation request: %s", jsonAttestation)
	} else {
		common.LogWithID(a.shortID, "Received attestation request but it was unprintable.")
	}

	secret, err := a.enclaveClient.ShareSecret(att)
	if err != nil {
		common.LogWithID(a.shortID, "Secret request failed, no response will be published. %s", err)
		return nil
	}

	// Store the attested key only if the attestation process succeeded.
	err = a.enclaveClient.StoreAttestation(att)
	if err != nil {
		return fmt.Errorf("could not store attestation. Cause: %w", err)
	}

	// todo: implement proper protocol so only one host responds to this secret requests initially
	// 	for now we just have the genesis host respond until protocol implemented
	if !a.config.IsGenesis {
		return nil
	}

	l1tx := &ethadapter.L1RespondSecretTx{
		Secret:      secret,
		RequesterID: att.Owner,
		AttesterID:  a.config.ID,
		HostAddress: att.HostAddress,
	}
	// TODO review: l1tx.Sign(a.attestationPubKey) doesn't matter as the waitSecret will process a tx that was reverted
	respondSecretTx := a.mgmtContractLib.CreateRespondSecret(l1tx, a.ethWallet.GetNonceAndIncrement(), false)
	err = a.signAndBroadcastTx(respondSecretTx, l1TxTriesSecret)
	if err != nil {
		return fmt.Errorf("could not broadcast secret response. Cause %w", err)
	}
	return nil
}

// Whenever we receive a new shared secret response transaction, we update our list of P2P peers, as another aggregator
// may have joined the network.
func (a *Node) processSharedSecretResponse(_ *ethadapter.L1RespondSecretTx) error {
	// We make a call to the L1 node to retrieve the new list of aggregators. An alternative would be to check that the
	// transaction succeeded, and if so, extract the additional host address from the transaction arguments. But we
	// believe this would be more brittle than just asking the L1 contract for its view of the current aggregators.
	msg, err := a.mgmtContractLib.GetHostAddresses()
	if err != nil {
		return err
	}
	response, err := a.ethClient.CallContract(msg)
	if err != nil {
		return err
	}
	decodedResponse, err := a.mgmtContractLib.DecodeCallResponse(response)
	if err != nil {
		return err
	}
	hostAddresses := decodedResponse[0]

	// We filter down the list of retrieved addresses.
	var filteredHostAddresses []string //nolint:prealloc
	for _, hostAddress := range hostAddresses {
		// We exclude our own address.
		if hostAddress == a.config.P2PPublicAddress {
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

	a.p2p.UpdatePeerList(filteredHostAddresses)
	return nil
}

// monitors the L1 client for new blocks and injects them into the aggregator
func (a *Node) monitorBlocks() {
	var lastKnownBlkHash gethcommon.Hash
	listener, subs := a.ethClient.BlockListener()
	common.LogWithID(a.shortID, "Start monitoring Ethereum blocks..")

	// only process blocks if the node is running
	for atomic.LoadInt32(a.stopNodeInterrupt) == 0 {
		select {
		case err := <-subs.Err():
			log.Error("L1 block monitoring error: %s", err)
			log.Info("Restarting L1 block Monitoring...")
			// it's fine to immediately restart the listener, any incoming blocks will be on hold in the queue
			listener, subs = a.ethClient.BlockListener()

			err = a.catchupMissedBlocks(lastKnownBlkHash)
			if err != nil {
				log.Panic("could not catch up missed blocks. Cause: %s", err)
			}

		case blkHeader := <-listener:
			// don't process blocks if the node is stopping
			if atomic.LoadInt32(a.stopNodeInterrupt) == 1 {
				break
			}

			// ignore blocks if bootstrapping is happening
			if atomic.LoadInt32(a.bootstrappingComplete) == 0 {
				log.Trace("Node in bootstrap - ignoring block %s", blkHeader.Hash())
				continue
			}

			block, err := a.ethClient.BlockByHash(blkHeader.Hash())
			if err != nil {
				log.Panic("could not fetch block for hash %s. Cause: %s", blkHeader.Hash().String(), err)
			}
			blockParent, err := a.ethClient.BlockByHash(block.ParentHash())
			if err != nil {
				log.Panic("could not fetch block's parent with hash %s. Cause: %s", block.ParentHash().String(), err)
			}

			common.LogWithID(
				a.shortID,
				"Received a new block b_%d(%d)",
				common.ShortHash(blkHeader.Hash()),
				blkHeader.Number.Uint64(),
			)

			// issue the block to the ingestion channel
			err = a.encodeAndIngest(block, blockParent)
			if err != nil {
				log.Panic(err.Error())
			}
			lastKnownBlkHash = block.Hash()
		}
	}

	log.Info("Stopped monitoring for l1 blocks")
	// make sure it cleanly unsubscribes
	// todo this should be defered when the errors are upstreamed instead of panic'd
	subs.Unsubscribe()
}

func (a *Node) catchupMissedBlocks(lastKnownBlkHash gethcommon.Hash) error {
	var lastBlkNumber *big.Int
	var reingestBlocks []*types.Block

	// get the blockchain tip block
	blk, err := a.ethClient.BlockByNumber(lastBlkNumber)
	if err != nil {
		return fmt.Errorf("catching up on missed blocks, unable to fetch tip block - reason: %w", err)
	}

	if blk.Hash().Hex() == lastKnownBlkHash.Hex() {
		// if no new blocks have been issued then nothing to catchup
		return nil
	}
	reingestBlocks = append(reingestBlocks, blk)

	// get all blocks from the blockchain tip to the last block ingested by the node
	for blk.Hash().Hex() != lastKnownBlkHash.Hex() {
		blockParent, err := a.ethClient.BlockByHash(blk.ParentHash())
		if err != nil {
			return fmt.Errorf("catching up on missed blocks, could not fetch block's parent with hash %s. Cause: %w", blk.ParentHash(), err)
		}

		reingestBlocks = append(reingestBlocks, blockParent)
		blk = blockParent
	}

	// make sure to have the last ingested block available for ingestion (because we always ingest ( blk, blk_parent)
	lastKnownBlk, err := a.ethClient.BlockByHash(lastKnownBlkHash)
	if err != nil {
		return fmt.Errorf("catching up on missed blocks, unable to feth last known block - reason: %w", err)
	}
	reingestBlocks = append(reingestBlocks, lastKnownBlk)

	// issue the block to the ingestion channel in reverse, with the parent attached too
	for i := len(reingestBlocks) - 2; i >= 0; i-- {
		log.Debug("Ingesting %s and %s blocks of %v", reingestBlocks[i].Hash(), reingestBlocks[i+1].Hash(), reingestBlocks)
		err = a.encodeAndIngest(reingestBlocks[i], reingestBlocks[i+1])
		if err != nil {
			log.Panic(err.Error())
		}
	}

	return nil
}

func (a *Node) encodeAndIngest(block *types.Block, blockParent *types.Block) error {
	encodedBlock, err := common.EncodeBlock(block)
	if err != nil {
		return fmt.Errorf("could not encode block with hash %s. Cause: %w", block.Hash().String(), err)
	}

	encodedBlockParent, err := common.EncodeBlock(blockParent)
	if err != nil {
		return fmt.Errorf("could not encode block's parent with hash %s. Cause: %w", block.ParentHash().String(), err)
	}

	a.blockRPCCh <- blockAndParent{encodedBlock, encodedBlockParent}
	return nil
}

func (a *Node) bootstrapNode() types.Block {
	var err error
	var nextBlk *types.Block

	// build up from the genesis block
	// todo update to bootstrap from the last block in storage
	// todo the genesis block should be the block where the contract was deployed
	currentBlock, err := a.ethClient.BlockByNumber(big.NewInt(0))
	if err != nil {
		log.Panic(err.Error())
	}

	common.LogWithID(a.shortID, "Started node bootstrap with block %d", currentBlock.NumberU64())

	startTime, logTime := time.Now(), time.Now()
	for {
		cb := *currentBlock
		a.processBlock(&cb)
		// TODO ingest one block at a time or batch the blocks
		result := a.enclaveClient.IngestBlocks([]*types.Block{&cb})
		if !result[0].IngestedBlock && result[0].BlockNotIngestedCause != "" {
			common.LogWithID(
				a.shortID,
				"Failed to ingest block b_%d. Cause: %s",
				common.ShortHash(result[0].BlockHeader.Hash()),
				result[0].BlockNotIngestedCause,
			)
		}
		a.storeBlockProcessingResult(result[0])

		nextBlk, err = a.ethClient.BlockByNumber(big.NewInt(cb.Number().Int64() + 1))
		if err != nil {
			if errors.Is(err, ethereum.NotFound) {
				break
			}
			log.Panic(err.Error())
		}
		currentBlock = nextBlk

		if time.Since(logTime) > 30*time.Second {
			common.LogWithID(a.shortID, "Bootstrapping node at block... %d", cb.NumberU64())
			logTime = time.Now()
		}
	}
	atomic.StoreInt32(a.bootstrappingComplete, 1)
	common.LogWithID(
		a.shortID,
		"Finished bootstrap process with block %d after %s",
		currentBlock.NumberU64(),
		time.Since(startTime),
	)
	return *currentBlock
}

func (a *Node) awaitSecret() error {
	// start listening for l1 blocks that contain the response to the request
	listener, subs := a.ethClient.BlockListener()

	for {
		select {
		case err := <-subs.Err():
			log.Error("Restarting L1 block monitoring while awaiting for secret. Errored with: %s", err)
			// todo this is a very simple way of reconnecting the node, it might need catching up logic
			listener, subs = a.ethClient.BlockListener()

		// todo: find a way to get rid of this case and only listen for blocks on the expected channels
		case header := <-listener:
			block, err := a.ethClient.BlockByHash(header.Hash())
			if err != nil {
				return fmt.Errorf("failed to retrieve block. Cause: %w", err)
			}
			if a.checkBlockForSecretResponse(block) {
				// todo this should be defered when the errors are upstreamed instead of panic'd
				subs.Unsubscribe()
				return nil
			}

		case bAndParent := <-a.blockRPCCh:
			block, err := bAndParent.b.DecodeBlock()
			if err != nil {
				return fmt.Errorf("failed to decode block received via RPC. Cause: %w", err)
			}
			if a.checkBlockForSecretResponse(block) {
				subs.Unsubscribe()
				return nil
			}

		case <-time.After(time.Second * 10):
			// This will provide useful feedback if things are stuck (and in tests if any goroutines got stranded on this select
			common.LogWithID(a.shortID, "Still waiting for secret from the L1...")

		case <-a.exitNodeCh:
			subs.Unsubscribe()
			return nil
		}
	}
}

func (a *Node) checkBlockForSecretResponse(block *types.Block) bool {
	for _, tx := range block.Transactions() {
		t := a.mgmtContractLib.DecodeTx(tx)
		if t == nil {
			continue
		}
		if scrtTx, ok := t.(*ethadapter.L1RespondSecretTx); ok {
			ok := a.handleStoreSecretTx(scrtTx)
			if ok {
				common.LogWithID(a.shortID, "Stored enclave secret.")
				return true
			}
		}
	}
	// response not found
	return false
}

// We retry calling `funcToRetry`, with a pause in between that starts at one second and doubles on each retry.
// This is equal to 7 seconds for 3 tries, and 63 seconds for 7 tries.
func retryWithBackoff(tries int, funcToRetry func() error) error {
	pause := time.Second
	var err error

	for i := 0; ; i++ {
		if i >= tries {
			// We've performed all the tries, so we return the (possibly nil) error.
			return err
		}

		err = funcToRetry()
		if err == nil {
			// If the transaction sent successfully, we break out.
			return nil
		}

		time.Sleep(pause)
		pause = pause * 2
	}
}
