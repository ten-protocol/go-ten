package obscuronode

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/obscuronet/obscuro-playground/go/common"
	"github.com/obscuronet/obscuro-playground/go/log"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave"

	obscuroCommon "github.com/obscuronet/obscuro-playground/go/obscuronode/common"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

type AggregatorCfg struct {
	// duration of the gossip round
	GossipRoundDuration uint64
}

type L2Network interface {
	BroadcastRollup(r common.EncodedRollup)
	BroadcastTx(tx obscuroCommon.EncryptedTx)
}

type StatsCollector interface {
	// Register when a node has to discard the speculative work built on top of the winner of the gossip round.
	L2Recalc(id gethcommon.Address)
	NewBlock(block *types.Block)
	NewRollup(rollup *obscuroCommon.Rollup)
	RollupWithMoreRecentProof()
}

// Node this will become the Obscuro "Node" type
type Node struct {
	ID gethcommon.Address

	l2Network L2Network
	L1Node    common.L1Node

	mining  bool // true -if this is an aggregator, false if it is a validator
	genesis bool // true - if this is the first Obscuro node which has to initialize the network
	cfg     AggregatorCfg

	stats StatsCollector

	// control the lifecycle
	exitNodeCh chan bool
	interrupt  *int32

	// blockRPCCh is where the connected L1Node node drops new blocks
	blockRPCCh chan blockAndParent

	// forkRPCCh is where new forks from the L1 notify the obscuro node
	forkRPCCh chan []common.EncodedBlock

	// rollupsP2PCh is the channel where new rollups are gossiped to
	rollupsP2PCh chan common.EncodedRollup

	// Interface to the logic running inside the TEE
	Enclave enclave.Enclave

	// Node nodeDB - stores the node public available data
	nodeDB *DB
}

func NewAgg(
	id gethcommon.Address,
	cfg AggregatorCfg,
	l1 common.L1Node,
	l2Network L2Network,
	collector StatsCollector,
	genesis bool,
) Node {
	return Node{
		// config
		ID:        id,
		cfg:       cfg,
		mining:    true,
		genesis:   genesis,
		L1Node:    l1,
		l2Network: l2Network,

		stats: collector,

		// lifecycle channels
		exitNodeCh: make(chan bool),
		interrupt:  new(int32),

		// incoming data
		blockRPCCh:   make(chan blockAndParent),
		forkRPCCh:    make(chan []common.EncodedBlock),
		rollupsP2PCh: make(chan common.EncodedRollup),

		// State processing
		Enclave: enclave.NewEnclave(id, true, collector),

		// Initialized the node nodeDB
		nodeDB: NewDB(),
	}
}

// Start initializes the main loop of the node
func (a *Node) Start() {
	if a.genesis {
		// Create the shared secret and submit it to the management contract for storage
		secret := a.Enclave.GenerateSecret()
		txData := common.L1TxData{
			TxType:      common.StoreSecretTx,
			Secret:      secret,
			Attestation: a.Enclave.Attestation(),
		}
		a.broadcastTx(*common.NewL1Tx(txData))
	}

	if !a.Enclave.IsInitialised() {
		a.requestSecret()
	}

	// todo create a channel between request secret and start processing
	a.startProcessing()
}

func (a *Node) startProcessing() {
	// Todo: This is a naive implementation.
	// It feeds the entire L1 blockchain into the enclave when it starts
	allblocks := a.L1Node.RPCBlockchainFeed()
	a.Enclave.IngestBlocks(allblocks)
	// todo - what happens with the blocks received while processing ?
	go a.Enclave.Start(*allblocks[len(allblocks)-1])

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
			a.processBlocks([]common.EncodedBlock{b.p, b.b}, interrupt)

		case f := <-a.forkRPCCh:
			interrupt = sendInterrupt(interrupt)
			a.processBlocks(f, interrupt)

		case r := <-a.rollupsP2PCh:
			rol, _ := obscuroCommon.Decode(r)
			go a.Enclave.SubmitRollup(obscuroCommon.ExtRollup{
				Header: rol.Header,
				Txs:    rol.Transactions,
			})

		case <-a.exitNodeCh:
			a.Enclave.Stop()
			return
		}
	}
}

// RPCNewHead receives the notification of new blocks from the L1Node Node
func (a *Node) RPCNewHead(b common.EncodedBlock, p common.EncodedBlock) {
	if atomic.LoadInt32(a.interrupt) == 1 {
		return
	}
	a.blockRPCCh <- blockAndParent{b, p}
}

// RPCNewFork receives the notification of a new fork from the L1Node
func (a *Node) RPCNewFork(b []common.EncodedBlock) {
	if atomic.LoadInt32(a.interrupt) == 1 {
		return
	}
	a.forkRPCCh <- b
}

// P2PGossipRollup is called by counterparties when there is a Rollup to broadcast
// All it does is forward the rollup for processing to the enclave
func (a *Node) P2PGossipRollup(r common.EncodedRollup) {
	if atomic.LoadInt32(a.interrupt) == 1 {
		return
	}
	a.rollupsP2PCh <- r
}

// P2PReceiveTx receives a new transactions from the P2P network
func (a *Node) P2PReceiveTx(tx obscuroCommon.EncryptedTx) {
	if atomic.LoadInt32(a.interrupt) == 1 {
		return
	}
	// Ignore gossiped transactions while the node is still initialising
	if a.Enclave.IsInitialised() {
		go func() {
			err := a.Enclave.SubmitTx(tx)
			if err != nil {
				log.Log(fmt.Sprintf(">   Agg%d: Could not submit transaction: %s", common.ShortAddress(a.ID), err))
			}
		}()
	}
}

// RPCBalance allows to fetch the balance of one address
func (a *Node) RPCBalance(address gethcommon.Address) uint64 {
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
	time.Sleep(time.Millisecond * 10)
	a.exitNodeCh <- true
}

func sendInterrupt(interrupt *int32) *int32 {
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
	var result enclave.SubmitBlockResponse
	for _, block := range blocks {
		// For the genesis block the parent is nil
		if block != nil {
			a.checkForSharedSecretRequests(block)

			// submit each block to the enclave for ingestion plus validation
			result = a.Enclave.SubmitBlock(*block.DecodeBlock())

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
	}

	if !result.IngestedBlock {
		b := blocks[len(blocks)-1].DecodeBlock()
		log.Log(fmt.Sprintf(">   Agg%d: Could not process block b_%d", common.ShortAddress(a.ID), common.ShortHash(b.Hash())))
		return
	}

	a.l2Network.BroadcastRollup(obscuroCommon.EncodeRollup(result.ProducedRollup.ToRollup()))

	common.ScheduleInterrupt(a.cfg.GossipRoundDuration, interrupt, func() {
		if atomic.LoadInt32(a.interrupt) == 1 {
			return
		}
		// Request the round winner for the current head
		winnerRollup, submit := a.Enclave.RoundWinner(result.L2Hash)
		if submit {
			txData := common.L1TxData{TxType: common.RollupTx, Rollup: obscuroCommon.EncodeRollup(winnerRollup.ToRollup())}
			tx := common.NewL1Tx(txData)
			t, err := common.EncodeTx(tx)
			if err != nil {
				panic(err)
			}
			a.L1Node.BroadcastTx(t)
			// collect Stats
			// a.stats.NewRollup(DecodeRollup(winnerRollup))
		}
	})
}

// Called only by the first enclave to bootstrap the network
func (a *Node) initialiseProtocol() common.L2RootHash {
	// Create the genesis rollup and submit it to the MC
	genesis := a.Enclave.ProduceGenesis()
	txData := common.L1TxData{TxType: common.RollupTx, Rollup: obscuroCommon.EncodeRollup(genesis.ProducedRollup.ToRollup())}
	a.broadcastTx(*common.NewL1Tx(txData))

	return genesis.L2Hash
}

func (a *Node) broadcastTx(tx common.L1Tx) {
	t, err := common.EncodeTx(&tx)
	if err != nil {
		panic(err)
	}
	a.L1Node.BroadcastTx(t)
}

// This method implements the procedure by which a node obtains the secret
func (a *Node) requestSecret() {
	txData := common.L1TxData{
		TxType:      common.RequestSecretTx,
		Attestation: a.Enclave.Attestation(),
	}
	a.broadcastTx(*common.NewL1Tx(txData))

	// start listening for l1 blocks that contain the response to the request
	for {
		select {
		case b := <-a.blockRPCCh:
			txs := b.b.DecodeBlock().Transactions()
			for _, tx := range txs {
				t := common.TxData(tx)
				if t.TxType == common.StoreSecretTx && t.Attestation.Owner == a.ID {
					// someone has replied
					a.Enclave.Init(t.Secret)
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

func (a *Node) checkForSharedSecretRequests(block common.EncodedBlock) {
	b := block.DecodeBlock()
	for _, tx := range b.Transactions() {
		t := common.TxData(tx)
		if t.TxType == common.RequestSecretTx {
			txData := common.L1TxData{
				TxType:      common.StoreSecretTx,
				Secret:      a.Enclave.FetchSecret(t.Attestation),
				Attestation: t.Attestation,
			}
			a.broadcastTx(*common.NewL1Tx(txData))
		}
	}
}
