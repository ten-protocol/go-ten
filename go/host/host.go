package host

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"sync/atomic"
	"time"

	"github.com/obscuronet/obscuro-playground/go/common/log"

	"github.com/ethereum/go-ethereum"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/naoina/toml"
	"github.com/obscuronet/obscuro-playground/go/common"
	"github.com/obscuronet/obscuro-playground/go/config"
	"github.com/obscuronet/obscuro-playground/go/ethadapter"
	"github.com/obscuronet/obscuro-playground/go/ethadapter/mgmtcontractlib"
	"github.com/obscuronet/obscuro-playground/go/wallet"
)

// Node this will become the Obscuro "Node" type
type Node struct {
	config  config.HostConfig
	ID      gethcommon.Address
	shortID uint64

	P2p           P2P                  // For communication with other Obscuro nodes
	ethClient     ethadapter.EthClient // For communication with the L1 node
	EnclaveClient common.Enclave       // For communication with the enclave
	rpcServer     RPCServer            // For communication with Obscuro client applications

	stats StatsCollector

	// control the host lifecycle
	exitNodeCh            chan bool
	stopNodeInterrupt     *int32
	bootstrappingComplete *int32 // Marks when the node is done bootstrapping

	blockRPCCh   chan blockAndParent        // The channel that new blocks from the L1 node are sent to
	forkRPCCh    chan []common.EncodedBlock // The channel that new forks from the L1 node are sent to
	rollupsP2PCh chan common.EncodedRollup  // The channel that new rollups from peers are sent to
	txP2PCh      chan common.EncryptedTx    // The channel that new transactions from peers are sent to

	nodeDB       *DB    // Stores the node's publicly-available data
	readyForWork *int32 // Whether the node has bootstrapped the existing blocks and has the enclave secret

	// library to handle Management Contract lib operations
	mgmtContractLib mgmtcontractlib.MgmtContractLib

	// Wallet used to issue ethereum transactions
	ethWallet wallet.Wallet
}

func NewHost(
	config config.HostConfig,
	collector StatsCollector,
	p2p P2P,
	ethClient ethadapter.EthClient,
	enclaveClient common.Enclave,
	ethWallet wallet.Wallet,
	mgmtContractLib mgmtcontractlib.MgmtContractLib,
) *Node {
	host := &Node{
		// config
		config:  config,
		ID:      config.ID,
		shortID: common.ShortAddress(config.ID),

		// Communication layers.
		P2p:           p2p,
		ethClient:     ethClient,
		EnclaveClient: enclaveClient,

		// statistics and metrics
		stats: collector,

		// lifecycle channels
		exitNodeCh:            make(chan bool),
		stopNodeInterrupt:     new(int32),
		bootstrappingComplete: new(int32),

		// incoming data
		blockRPCCh:   make(chan blockAndParent),
		forkRPCCh:    make(chan []common.EncodedBlock),
		rollupsP2PCh: make(chan common.EncodedRollup),
		txP2PCh:      make(chan common.EncryptedTx),

		// Initialize the node DB
		// nodeDB:       NewLevelDBBackedDB(), // todo - make this config driven
		nodeDB:       NewInMemoryDB(),
		readyForWork: new(int32),

		// library that provides a handler for Management Contract
		mgmtContractLib: mgmtContractLib,
		// the nodes ethereum wallet
		ethWallet: ethWallet,
	}

	if config.HasClientRPCHTTP || config.HasClientRPCWebsockets {
		host.rpcServer = NewRPCServer(config, host)
	}

	return host
}

// Start initializes the main loop of the node
func (a *Node) Start() {
	tomlConfig, err := toml.Marshal(a.config)
	if err != nil {
		panic("could not print host config")
	}
	common.LogWithID(a.shortID, "Host started with following config:\n%s", tomlConfig)

	// wait for the Enclave to be available
	a.waitForEnclave()

	// todo: we should try to recover the key from a previous run of the node here? Before generating or requesting the key.
	if a.config.IsGenesis {
		common.LogWithID(a.shortID, "Node is genesis node. Broadcasting secret.")
		// Create the shared secret and submit it to the management contract for storage
		attestation := a.EnclaveClient.Attestation()
		if attestation.Owner != a.ID {
			log.Panic(">   Agg%d: genesis node has ID %s, but its enclave produced an attestation using ID %s", a.shortID, a.ID.Hex(), attestation.Owner.Hex())
		}

		l1tx := &ethadapter.L1InitializeSecretTx{
			AggregatorID:  &a.ID,
			InitialSecret: a.EnclaveClient.GenerateSecret(),
			HostAddress:   a.config.P2PAddress,
		}
		a.broadcastTx(a.mgmtContractLib.CreateInitializeSecret(l1tx, a.ethWallet.GetNonceAndIncrement()))
		common.LogWithID(a.shortID, "Node is genesis node. Secret was broadcasted.")
	} else {
		a.requestSecret()
	}

	// attach the l1 monitor
	go a.monitorBlocks()

	// bootstrap the node
	latestBlock := a.bootstrapNode()

	// start the enclave speculative work from last block
	a.EnclaveClient.Start(latestBlock)

	if a.config.IsGenesis {
		a.initialiseProtocol(&latestBlock)
	}
	// start the obscuro RPC endpoints
	if a.rpcServer != nil {
		a.rpcServer.Start()
		common.LogWithID(a.shortID, "Started client server.")
	}

	// start the node main processing loop
	a.startProcessing()
}

// MockedNewHead receives the notification of new blocks
// This endpoint is specific to the ethereum mock node
func (a *Node) MockedNewHead(b common.EncodedBlock, p common.EncodedBlock) {
	if atomic.LoadInt32(a.stopNodeInterrupt) == 1 {
		return
	}
	a.blockRPCCh <- blockAndParent{b, p}
}

// MockedNewFork receives the notification of a new fork
// This endpoint is specific to the ethereum mock node
func (a *Node) MockedNewFork(b []common.EncodedBlock) {
	if atomic.LoadInt32(a.stopNodeInterrupt) == 1 {
		return
	}
	a.forkRPCCh <- b
}

func (a *Node) SubmitAndBroadcastTx(encryptedParams common.EncryptedParamsSendRawTx) (common.EncryptedResponseSendRawTx, error) {
	encryptedTx := common.EncryptedTx(encryptedParams)
	encryptedResponse, err := a.EnclaveClient.SubmitTx(encryptedTx)
	if err != nil {
		log.Info(fmt.Sprintf(">   Agg%d: Could not submit transaction: %s", a.shortID, err))
	}
	a.P2p.BroadcastTx(encryptedTx)

	return encryptedResponse, err
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

// RPCExecuteOffChainTransaction allows execution of off chain transactions
func (a *Node) RPCExecuteOffChainTransaction(encryptedParams common.EncryptedParamsCall) (common.EncryptedResponseCall, error) {
	return a.EnclaveClient.ExecuteOffChainTransaction(encryptedParams)
}

// RPCCurrentBlockHead returns the current head of the blocks (l1)
func (a *Node) RPCCurrentBlockHead() *types.Header {
	return a.nodeDB.GetCurrentBlockHead()
}

// RPCCurrentRollupHead returns the current head of the rollups (l2)
func (a *Node) RPCCurrentRollupHead() *common.Header {
	return a.nodeDB.GetCurrentRollupHead()
}

// DB returns the DB of the node
func (a *Node) DB() *DB {
	return a.nodeDB
}

// Stop gracefully stops the node execution
func (a *Node) Stop() {
	// block all requests
	atomic.StoreInt32(a.stopNodeInterrupt, 1)

	if err := a.P2p.StopListening(); err != nil {
		common.ErrorWithID(a.shortID, "failed to close transaction P2P listener cleanly: %s", err)
	}
	if err := a.EnclaveClient.Stop(); err != nil {
		common.ErrorWithID(a.shortID, "could not stop enclave server. Cause: %s", err)
	}

	if err := a.EnclaveClient.StopClient(); err != nil {
		common.ErrorWithID(a.shortID, "failed to stop enclave RPC client. Cause: %s", err)
	}
	if a.rpcServer != nil {
		if err := a.rpcServer.Stop(); err != nil {
			common.ErrorWithID(a.shortID, "could not stop client RPC server. Cause: %s", err)
		}
	}

	// Leave some time for all processing to finish before exiting the main loop.
	time.Sleep(time.Second)
	a.exitNodeCh <- true
}

// ConnectToEthNode connects the Aggregator to the ethereum node
func (a *Node) ConnectToEthNode(node ethadapter.EthClient) {
	a.ethClient = node
}

// IsReady returns if the Aggregator is ready to work (process blocks, respond to RPC requests, etc..)
func (a *Node) IsReady() bool {
	return atomic.LoadInt32(a.readyForWork) == 1
}

// Waits for enclave to be available, printing a wait message every two seconds.
func (a *Node) waitForEnclave() {
	counter := 0
	for err := a.EnclaveClient.IsReady(); err != nil; {
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
	// Only open the p2p connection when the node is fully initialised
	a.P2p.StartListening(a)

	// use the roundInterrupt as a signaling mechanism for interrupting block processing
	// stops processing the current round if a new block arrives
	i := int32(0)
	roundInterrupt := &i

	// marks the node as ready to do work ( process blocks, respond to RPC requests, etc... )
	atomic.StoreInt32(a.readyForWork, 1)
	common.LogWithID(a.shortID, "Node is ready for work...")

	// Main Processing Loop -
	// - Process new blocks from the L1 node
	// - Process new Rollups gossiped from L2 Peers
	// - Process new Transactions gossiped from L2 Peers
	for {
		select {
		case b := <-a.blockRPCCh:
			roundInterrupt = triggerInterrupt(roundInterrupt)
			a.processBlocks([]common.EncodedBlock{b.p, b.b}, roundInterrupt)

		case f := <-a.forkRPCCh:
			roundInterrupt = triggerInterrupt(roundInterrupt)
			a.processBlocks(f, roundInterrupt)

		case r := <-a.rollupsP2PCh:
			rol, err := common.DecodeRollup(r)
			log.Trace(fmt.Sprintf(">   Agg%d: Received rollup: r_%d from A%d",
				a.shortID,
				common.ShortHash(rol.Hash()),
				common.ShortAddress(rol.Header.Agg),
			))
			if err != nil {
				common.LogWithID(a.shortID, "Could not check enclave initialisation. Cause: %v", err)
			}

			go a.EnclaveClient.SubmitRollup(common.ExtRollup{
				Header:          rol.Header,
				EncryptedTxBlob: rol.Transactions,
			})

		case tx := <-a.txP2PCh:
			if _, err := a.EnclaveClient.SubmitTx(tx); err != nil {
				log.Info(fmt.Sprintf(">   Agg%d: Could not submit transaction: %s", a.shortID, err))
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

func (a *Node) processBlocks(blocks []common.EncodedBlock, interrupt *int32) {
	var result common.BlockSubmissionResponse
	for _, block := range blocks {
		// For the genesis block the parent is nil
		if block != nil {
			a.processBlock(block)

			// submit each block to the enclave for ingestion plus validation
			result = a.EnclaveClient.SubmitBlock(*block.DecodeBlock())
			a.storeBlockProcessingResult(result)
		}
	}

	if !result.IngestedBlock {
		b := blocks[len(blocks)-1].DecodeBlock()
		common.LogWithID(a.shortID, "Did not ingest block b_%d. Cause: %s", common.ShortHash(b.Hash()), result.BlockNotIngestedCause)
		return
	}

	// Nodes can start before the genesis was published, and it makes no sense to enter the protocol.
	if result.ProducedRollup.Header != nil {
		a.P2p.BroadcastRollup(common.EncodeRollup(result.ProducedRollup.ToRollup()))

		common.ScheduleInterrupt(a.config.GossipRoundDuration, interrupt, a.handleRoundWinner(result))
	}
}

// Looks at each transaction in the block, and kicks off special handling for the transaction if needed.
func (a *Node) processBlock(block common.EncodedBlock) {
	b := block.DecodeBlock()
	for _, tx := range b.Transactions() {
		t := a.mgmtContractLib.DecodeTx(tx)
		if t == nil {
			continue
		}

		if scrtReqTx, ok := t.(*ethadapter.L1RequestSecretTx); ok {
			a.processSharedSecretRequest(scrtReqTx)
		}

		if scrtRespTx, ok := t.(*ethadapter.L1RespondSecretTx); ok {
			err := a.processSharedSecretResponse(scrtRespTx)
			if err != nil {
				common.LogWithID(a.shortID, "Failed to process shared secret response. Cause: %s", err)
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
		winnerRollup, isWinner, err := a.EnclaveClient.RoundWinner(result.ProducedRollup.Header.ParentHash)
		if err != nil {
			log.Panic("could not determine round winner. Cause: %s", err)
		}
		if isWinner {
			common.LogWithID(a.shortID, "Winner (b_%d) r_%d(%d).",
				common.ShortHash(result.BlockHeader.Hash()),
				common.ShortHash(winnerRollup.Header.Hash()),
				winnerRollup.Header.Number,
			)

			tx := &ethadapter.L1RollupTx{
				Rollup: common.EncodeRollup(winnerRollup.ToRollup()),
			}

			// That handler can get called multiple times for the same height. And it will return the same winner rollup.
			// In case the winning rollup belongs to the current enclave it will be submitted again, which is inefficient.
			if !a.DB().WasSubmitted(winnerRollup.Header.Hash()) {
				a.broadcastTx(a.mgmtContractLib.CreateRollup(tx, a.ethWallet.GetNonceAndIncrement()))
				a.DB().AddSubmittedRollup(winnerRollup.Header.Hash())
			}
		}
	}
}

func (a *Node) storeBlockProcessingResult(result common.BlockSubmissionResponse) {
	// only update the node rollup headers if the enclave has found a new rollup head
	if result.FoundNewHead {
		// adding a header will update the head if it has a higher height
		a.DB().AddRollupHeader(result.RollupHead)
	}

	// adding a header will update the head if it has a higher height
	if result.IngestedBlock {
		a.DB().AddBlockHeader(result.BlockHeader)
	}
}

// Called only by the first enclave to bootstrap the network
func (a *Node) initialiseProtocol(block *types.Block) common.L2RootHash {
	// Create the genesis rollup and submit it to the MC
	genesisResponse := a.EnclaveClient.ProduceGenesis(block.Hash())
	common.LogWithID(
		a.shortID,
		"Initialising network. Genesis rollup r_%d.",
		common.ShortHash(genesisResponse.ProducedRollup.Header.Hash()),
	)
	l1tx := &ethadapter.L1RollupTx{
		Rollup: common.EncodeRollup(genesisResponse.ProducedRollup.ToRollup()),
	}

	a.broadcastTx(a.mgmtContractLib.CreateRollup(l1tx, a.ethWallet.GetNonceAndIncrement()))

	return genesisResponse.ProducedRollup.Header.ParentHash
}

func (a *Node) broadcastTx(tx types.TxData) {
	// TODO add retry and deal with failures
	signedTx, err := a.ethWallet.SignTransaction(tx)
	if err != nil {
		panic(err)
	}

	err = a.ethClient.SendTransaction(signedTx)
	if err != nil {
		panic(err)
	}
}

// This method implements the procedure by which a node obtains the secret
func (a *Node) requestSecret() {
	common.LogWithID(a.shortID, "Requesting secret.")
	att := a.EnclaveClient.Attestation()
	if att.Owner != a.ID {
		log.Panic(">   Agg%d: node has ID %s, but its enclave produced an attestation using ID %s", a.shortID, a.ID.Hex(), att.Owner.Hex())
	}
	encodedAttestation := common.EncodeAttestation(att)
	l1tx := &ethadapter.L1RequestSecretTx{
		Attestation: encodedAttestation,
	}
	a.broadcastTx(a.mgmtContractLib.CreateRequestSecret(l1tx, a.ethWallet.GetNonceAndIncrement()))

	a.awaitSecret()
}

func (a *Node) handleStoreSecretTx(t *ethadapter.L1RespondSecretTx) bool {
	if t.RequesterID.Hex() != a.ID.Hex() {
		// this secret is for somebody else
		return false
	}

	// someone has replied for us
	err := a.EnclaveClient.InitEnclave(t.Secret)
	if err != nil {
		common.LogWithID(a.shortID, "Failed to initialise enclave with received secret. Err: %s", err)
		return false
	}
	return true
}

func (a *Node) processSharedSecretRequest(scrtReqTx *ethadapter.L1RequestSecretTx) {
	// todo: implement proper protocol so only one host responds to this secret requests initially
	// 	for now we just have the genesis host respond until protocol implemented
	if !a.config.IsGenesis {
		return
	}

	att, err := common.DecodeAttestation(scrtReqTx.Attestation)
	if err != nil {
		common.LogWithID(a.shortID, "Failed to decode attestation. %s", err)
		return
	}

	jsonAttestation, err := json.Marshal(att)
	if err == nil {
		common.LogWithID(a.shortID, "Received attestation request: %s", jsonAttestation)
	} else {
		common.LogWithID(a.shortID, "Received attestation request but it was unprintable.")
	}

	secret, err := a.EnclaveClient.ShareSecret(att)
	if err != nil {
		common.LogWithID(a.shortID, "Secret request failed, no response will be published. %s", err)
		return
	}

	l1tx := &ethadapter.L1RespondSecretTx{
		Secret:      secret,
		RequesterID: att.Owner,
		AttesterID:  a.ID,
		HostAddress: att.HostAddress,
	}
	// TODO review: l1tx.Sign(a.attestationPubKey) doesn't matter as the waitSecret will process a tx that was reverted
	a.broadcastTx(a.mgmtContractLib.CreateRespondSecret(l1tx, a.ethWallet.GetNonceAndIncrement(), false))
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
		if hostAddress == a.config.P2PAddress {
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

	a.P2p.UpdatePeerList(filteredHostAddresses)
	common.LogWithID(a.shortID, "Updated peer list to %s", strings.Join(filteredHostAddresses, ", "))
	return nil
}

// monitors the L1 client for new blocks and injects them into the aggregator
func (a *Node) monitorBlocks() {
	listener := a.ethClient.BlockListener()
	common.LogWithID(a.shortID, "Start monitoring Ethereum blocks..")

	// only process blocks if the node is running
	for atomic.LoadInt32(a.stopNodeInterrupt) == 0 {
		blkHeader := <-listener

		// don't process blocks if the node is stopping
		if atomic.LoadInt32(a.stopNodeInterrupt) == 1 {
			return
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
		a.blockRPCCh <- blockAndParent{common.EncodeBlock(block), common.EncodeBlock(blockParent)}
	}
}

func (a *Node) bootstrapNode() types.Block {
	var err error
	var nextBlk *types.Block

	// build up from the genesis block
	// todo update to bootstrap from the last block in storage
	// todo the genesis block should be the block where the contract was deployed
	currentBlock, err := a.ethClient.BlockByNumber(big.NewInt(0))
	if err != nil {
		panic(err)
	}

	common.LogWithID(a.shortID, "Started node bootstrap with block %d", currentBlock.NumberU64())

	startTime, logTime := time.Now(), time.Now()
	for {
		// TODO ingest one block at a time or batch the blocks
		result := a.EnclaveClient.IngestBlocks([]*types.Block{currentBlock})
		if !result[0].IngestedBlock && result[0].BlockNotIngestedCause != "" {
			common.LogWithID(
				a.shortID,
				"Failed to ingest block b_%d. Cause: %s",
				common.ShortHash(result[0].BlockHeader.Hash()),
				result[0].BlockNotIngestedCause,
			)
		}
		a.storeBlockProcessingResult(result[0])

		nextBlk, err = a.ethClient.BlockByNumber(big.NewInt(currentBlock.Number().Int64() + 1))
		if err != nil {
			if errors.Is(err, ethereum.NotFound) {
				break
			}
			panic(err)
		}
		currentBlock = nextBlk

		if time.Since(logTime) > 30*time.Second {
			common.LogWithID(a.shortID, "Bootstrapping node at block... %d", currentBlock.NumberU64())
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

func (a *Node) awaitSecret() {
	// start listening for l1 blocks that contain the response to the request
	for {
		select {
		// todo: find a way to get rid of this case and only listen for blocks on the expected channels
		case header := <-a.ethClient.BlockListener():
			block, err := a.ethClient.BlockByHash(header.Hash())
			if err != nil {
				log.Panic("failed to retrieve block. Cause: %s:", err)
			}
			if a.checkBlockForSecretResponse(block) {
				return
			}

		case b := <-a.blockRPCCh:
			if a.checkBlockForSecretResponse(b.b.DecodeBlock()) {
				return
			}

		case <-time.After(time.Second * 10):
			// This will provide useful feedback if things are stuck (and in tests if any goroutines got stranded on this select
			common.LogWithID(a.shortID, "Still waiting for secret from the L1...")

		case <-a.exitNodeCh:
			return
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
