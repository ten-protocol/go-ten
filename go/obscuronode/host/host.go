package host

import (
	"fmt"
	"sync/atomic"
	"time"

	p2p2 "github.com/obscuronet/obscuro-playground/go/obscuronode/host/p2p"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/go/log"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"

	"github.com/ethereum/go-ethereum/common"
)

const ClientRPCTimeoutSecs = 5

type AggregatorCfg struct {
	// duration of the gossip round
	GossipRoundDuration uint64
	// timeout duration in seconds for RPC requests to the enclave service
	ClientRPCTimeoutSecs uint64
}

type StatsCollector interface {
	// Register when a node has to discard the speculative work built on top of the winner of the gossip round.
	L2Recalc(id common.Address)
	NewBlock(block *types.Block)
	NewRollup(rollup *nodecommon.Rollup)
	RollupWithMoreRecentProof()
}

// Node this will become the Obscuro "Node" type
type Node struct {
	ID common.Address

	L1Node obscurocommon.L1Node

	mining  bool // true -if this is an aggregator, false if it is a validator
	genesis bool // true - if this is the first Obscuro node which has to initialize the network
	cfg     AggregatorCfg

	stats StatsCollector

	// control the lifecycle
	exitNodeCh chan bool
	interrupt  *int32

	blockRPCCh   chan blockAndParent               // The channel that new blocks from the L1 node are sent to
	forkRPCCh    chan []obscurocommon.EncodedBlock // The channel that new forks from the L1 node are sent to
	rollupsP2PCh chan obscurocommon.EncodedRollup  // The channel that new rollups from peers are sent to
	txP2PCh      chan nodecommon.EncryptedTx       // The channel that new transactions from peers are sent to

	// Interface to the logic running inside the TEE
	Enclave nodecommon.Enclave

	// Node nodeDB - stores the node public available data
	nodeDB *DB

	p2p p2p2.P2P
}

func NewAgg(
	id common.Address,
	cfg AggregatorCfg,
	l1 obscurocommon.L1Node,
	collector StatsCollector,
	genesis bool,
	enclaveClient nodecommon.Enclave,
	p2p p2p2.P2P,
) Node {
	return Node{
		// config
		ID:      id,
		cfg:     cfg,
		mining:  true,
		genesis: genesis,
		L1Node:  l1,

		stats: collector,

		// lifecycle channels
		exitNodeCh: make(chan bool),
		interrupt:  new(int32),

		// incoming data
		blockRPCCh:   make(chan blockAndParent),
		forkRPCCh:    make(chan []obscurocommon.EncodedBlock),
		rollupsP2PCh: make(chan obscurocommon.EncodedRollup),
		txP2PCh:      make(chan nodecommon.EncryptedTx),

		// State processing
		Enclave: enclaveClient,

		// Initialized the node nodeDB
		nodeDB: NewDB(),

		p2p: p2p,
	}
}

// Start initializes the main loop of the node
func (a *Node) Start() {
	a.p2p.Listen(a.txP2PCh, a.rollupsP2PCh)
	defer a.p2p.StopListening()

	if a.genesis {
		// Create the shared secret and submit it to the management contract for storage
		txData := obscurocommon.L1TxData{
			TxType:      obscurocommon.StoreSecretTx,
			Secret:      a.Enclave.GenerateSecret(),
			Attestation: a.Enclave.Attestation(),
		}
		a.broadcastTx(*obscurocommon.NewL1Tx(txData))
	}

	if !a.Enclave.IsInitialised() {
		a.requestSecret()
	}

	allBlocks := a.waitForL1Blocks()

	// todo create a channel between request secret and start processing
	a.startProcessing(allBlocks)
}

// Waits for blocks from the L1 node, printing a wait message every two seconds.
func (a *Node) waitForL1Blocks() []*types.Block {
	// It feeds the entire L1 blockchain into the enclave when it starts
	// todo - what happens with the blocks received while processing ?
	allBlocks := a.L1Node.RPCBlockchainFeed()
	counter := 0

	for len(allBlocks) == 0 {
		if counter >= 20 {
			log.Log(fmt.Sprintf(">   Agg%d: Waiting for blocks from L1 node...", obscurocommon.ShortAddress(a.ID)))
			counter = 0
		}

		time.Sleep(100 * time.Millisecond)
		allBlocks = a.L1Node.RPCBlockchainFeed()
		counter++
	}

	return allBlocks
}

func (a *Node) startProcessing(allblocks []*types.Block) {
	// Todo: This is a naive implementation.
	results := a.Enclave.IngestBlocks(allblocks)
	for _, result := range results {
		a.storeBlockProcessingResult(result)
	}

	a.Enclave.Start(*allblocks[len(allblocks)-1])

	if a.genesis {
		a.initialiseProtocol()
	}

	// used as a signaling mechanism to stop processing the old block if a new L1 block arrives earlier
	i := int32(0)
	interrupt := &i

	// Main loop - Listen for notifications From the L1 node and process them
	// Note that during processing, more recent notifications can be received.
	for {
		select {
		case b := <-a.blockRPCCh:
			interrupt = sendInterrupt(interrupt)
			a.processBlocks([]obscurocommon.EncodedBlock{b.p, b.b}, interrupt)

		case f := <-a.forkRPCCh:
			interrupt = sendInterrupt(interrupt)
			a.processBlocks(f, interrupt)

		case r := <-a.rollupsP2PCh:
			rol, err := nodecommon.DecodeRollup(r)
			if err != nil {
				log.Log(fmt.Sprintf(">   Agg%d: Could not check enclave initialisation: %v", obscurocommon.ShortAddress(a.ID), err))
			}

			go a.Enclave.SubmitRollup(nodecommon.ExtRollup{
				Header: rol.Header,
				Txs:    rol.Transactions,
			})

		case tx := <-a.txP2PCh:
			// Ignore gossiped transactions while the node is still initialising
			if a.Enclave.IsInitialised() {
				if err := a.Enclave.SubmitTx(tx); err != nil {
					log.Log(fmt.Sprintf(">   Agg%d: Could not submit transaction: %s", obscurocommon.ShortAddress(a.ID), err))
				}
			}

		case <-a.exitNodeCh:
			return
		}
	}
}

// RPCNewHead receives the notification of new blocks from the L1Node Node
func (a *Node) RPCNewHead(b obscurocommon.EncodedBlock, p obscurocommon.EncodedBlock) {
	if atomic.LoadInt32(a.interrupt) == 1 {
		return
	}
	a.blockRPCCh <- blockAndParent{b, p}
}

// RPCNewFork receives the notification of a new fork from the L1Node
func (a *Node) RPCNewFork(b []obscurocommon.EncodedBlock) {
	if atomic.LoadInt32(a.interrupt) == 1 {
		return
	}
	a.forkRPCCh <- b
}

// RPCBalance allows to fetch the balance of one address
func (a *Node) RPCBalance(address common.Address) uint64 {
	return a.Enclave.Balance(address)
}

// RPCCurrentBlockHead returns the current head of the blocks (l1)
func (a *Node) RPCCurrentBlockHead() *BlockHeader {
	return a.nodeDB.GetCurrentBlockHead()
}

// RPCCurrentRollupHead returns the current head of the rollups (l2)
func (a *Node) RPCCurrentRollupHead() *RollupHeader {
	return a.nodeDB.GetCurrentRollupHead()
}

// DB returns the DB of the node
func (a *Node) DB() *DB {
	return a.nodeDB
}

// Stop gracefully stops the node execution
func (a *Node) Stop() {
	// block all requests
	atomic.StoreInt32(a.interrupt, 1)
	a.Enclave.Stop()

	time.Sleep(time.Millisecond * 1000)
	a.exitNodeCh <- true

	a.Enclave.StopClient()
}

func sendInterrupt(interrupt *int32) *int32 {
	// Notify the previous round to stop work
	atomic.StoreInt32(interrupt, 1)
	i := int32(0)
	return &i
}

type blockAndParent struct {
	b obscurocommon.EncodedBlock
	p obscurocommon.EncodedBlock
}

func (a *Node) processBlocks(blocks []obscurocommon.EncodedBlock, interrupt *int32) {
	var result nodecommon.BlockSubmissionResponse
	for _, block := range blocks {
		// For the genesis block the parent is nil
		if block != nil {
			a.checkForSharedSecretRequests(block)

			// submit each block to the enclave for ingestion plus validation
			result = a.Enclave.SubmitBlock(*block.DecodeBlock())
			a.storeBlockProcessingResult(result)
		}
	}

	if !result.IngestedBlock {
		b := blocks[len(blocks)-1].DecodeBlock()
		log.Log(fmt.Sprintf(">   Agg%d: Could not process block b_%d", obscurocommon.ShortAddress(a.ID), obscurocommon.ShortHash(b.Hash())))
		return
	}

	// todo -make this a better check
	if result.ProducedRollup.Header != nil {
		a.p2p.BroadcastRollup(nodecommon.EncodeRollup(result.ProducedRollup.ToRollup()))

		obscurocommon.ScheduleInterrupt(a.cfg.GossipRoundDuration, interrupt, func() {
			if atomic.LoadInt32(a.interrupt) == 1 {
				return
			}
			// Request the round winner for the current head
			winnerRollup, submit := a.Enclave.RoundWinner(result.L2Hash)
			if submit {
				txData := obscurocommon.L1TxData{TxType: obscurocommon.RollupTx, Rollup: nodecommon.EncodeRollup(winnerRollup.ToRollup())}
				tx := obscurocommon.NewL1Tx(txData)
				t, err := obscurocommon.EncodeTx(tx)
				if err != nil {
					panic(err)
				}
				a.L1Node.BroadcastTx(t)
				// collect Stats
				// a.stats.NewRollup(DecodeRollupOrPanic(winnerRollup))
			}
		})
	}
}

func (a *Node) storeBlockProcessingResult(result nodecommon.BlockSubmissionResponse) {
	// only update the node rollup headers if the enclave has ingested it
	if result.IngestedNewRollup {
		// adding a header will update the head if it has a higher height
		a.DB().AddRollupHeader(
			&RollupHeader{
				ID:          result.L2Hash,
				Parent:      result.L2Parent,
				Withdrawals: result.Withdrawals,
				Height:      result.L2Height,
			},
		)
	}

	// adding a header will update the head if it has a higher height
	a.DB().AddBlockHeader(
		&BlockHeader{
			ID:     result.L1Hash,
			Parent: result.L1Parent,
			Height: result.L1Height,
		},
	)
}

// Called only by the first enclave to bootstrap the network
func (a *Node) initialiseProtocol() obscurocommon.L2RootHash {
	// Create the genesis rollup and submit it to the MC
	genesis := a.Enclave.ProduceGenesis()
	txData := obscurocommon.L1TxData{TxType: obscurocommon.RollupTx, Rollup: nodecommon.EncodeRollup(genesis.ProducedRollup.ToRollup())}
	a.broadcastTx(*obscurocommon.NewL1Tx(txData))

	return genesis.L2Hash
}

func (a *Node) broadcastTx(tx obscurocommon.L1Tx) {
	t, err := obscurocommon.EncodeTx(&tx)
	if err != nil {
		panic(err)
	}
	a.L1Node.BroadcastTx(t)
}

// This method implements the procedure by which a node obtains the secret
func (a *Node) requestSecret() {
	attestation := a.Enclave.Attestation()
	txData := obscurocommon.L1TxData{
		TxType:      obscurocommon.RequestSecretTx,
		Attestation: attestation,
	}
	a.broadcastTx(*obscurocommon.NewL1Tx(txData))

	// start listening for l1 blocks that contain the response to the request
	for {
		select {
		case b := <-a.blockRPCCh:
			txs := b.b.DecodeBlock().Transactions()
			for _, tx := range txs {
				t := obscurocommon.TxData(tx)
				if t.TxType == obscurocommon.StoreSecretTx && t.Attestation.Owner == a.ID {
					// someone has replied
					a.Enclave.InitEnclave(t.Secret)
					return
				}
			}

		case <-a.forkRPCCh:
			// todo

		case <-a.rollupsP2PCh:
			// ignore rolllups from peers as we're not part of the network just yet

		case <-a.exitNodeCh:
			return
		}
	}
}

func (a *Node) checkForSharedSecretRequests(block obscurocommon.EncodedBlock) {
	b := block.DecodeBlock()
	for _, tx := range b.Transactions() {
		t := obscurocommon.TxData(tx)
		if t.TxType == obscurocommon.RequestSecretTx {
			txData := obscurocommon.L1TxData{
				TxType:      obscurocommon.StoreSecretTx,
				Secret:      a.Enclave.FetchSecret(t.Attestation),
				Attestation: t.Attestation,
			}
			a.broadcastTx(*obscurocommon.NewL1Tx(txData))
		}
	}
}
