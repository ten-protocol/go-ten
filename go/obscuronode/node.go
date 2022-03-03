package obscuronode

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/google/uuid"

	"github.com/obscuronet/obscuro-playground/go/common"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave"

	obscuroCommon "github.com/obscuronet/obscuro-playground/go/obscuronode/common"
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
	L2Recalc(id common.NodeID)
	NewBlock(block *common.Block)
	NewRollup(rollup *obscuroCommon.Rollup)
	RollupWithMoreRecentProof()
}

// Node this will become the Obscuro "Node" type
type Node struct {
	ID common.NodeID

	l2Network L2Network
	L1Node    common.L1Node

	mining  bool // true -if this is an aggregator, false if it is a validator
	genesis bool // true - if this is the first Obscuro node which has to initialize the network
	cfg     AggregatorCfg

	stats StatsCollector

	// control the lifecycle
	exitNodeCh chan bool
	interrupt  *int32

	// where the connected L1Node node drops new blocks
	blockRPCCh chan blockAndParent
	forkRPCCh  chan []common.EncodedBlock

	rollupsP2PCh chan common.EncodedRollup

	// Interface to the logic running inside the TEE
	Enclave enclave.Enclave
}

func (a *Node) Start() {
	if a.genesis {
		// Create the shared secret and submit it to the management contract for storage
		secret := a.Enclave.GenerateSecret()
		a.broadcastTx(common.L1Tx{
			ID:          uuid.New(),
			TxType:      common.StoreSecretTx,
			Secret:      secret,
			Attestation: a.Enclave.Attestation(),
		})
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
	extblocks := make([]common.ExtBlock, len(allblocks))

	for i, b := range allblocks {
		extblocks[i] = b.ToExtBlock()
	}

	a.Enclave.IngestBlocks(extblocks)
	// todo - what happens with the blocks received while processing ?
	go a.Enclave.Start(extblocks[len(extblocks)-1])

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

func sendInterrupt(interrupt *int32) *int32 {
	// Notify the previous round to stop work
	atomic.StoreInt32(interrupt, 1)
	i := int32(0)
	return &i
}

func (a *Node) processBlocks(blocks []common.EncodedBlock, interrupt *int32) {
	var result enclave.SubmitBlockResponse
	for _, block := range blocks {
		// For the genesis block the parent is nil
		if block != nil {
			a.checkForSharedSecretRequests(block)
			result = a.Enclave.SubmitBlock(block.DecodeBlock().ToExtBlock())
		}
	}

	if !result.Processed {
		b := blocks[len(blocks)-1].DecodeBlock()
		common.Log(fmt.Sprintf(">   Agg%d: Could not process block b_%s", a.ID, common.Str(b.Hash())))
		return
	}

	a.l2Network.BroadcastRollup(obscuroCommon.EncodeRollup(result.Rollup.ToRollup()))

	common.ScheduleInterrupt(a.cfg.GossipRoundDuration, interrupt, func() {
		if atomic.LoadInt32(a.interrupt) == 1 {
			return
		}
		// Request the round winner for the current head
		winnerRollup, submit := a.Enclave.RoundWinner(result.Hash)
		if submit {
			tx := common.L1Tx{ID: uuid.New(), TxType: common.RollupTx, Rollup: obscuroCommon.EncodeRollup(winnerRollup.ToRollup())}
			t, err := tx.Encode()
			if err != nil {
				panic(err)
			}
			a.L1Node.BroadcastTx(t)
			// collect Stats
			// a.stats.NewRollup(DecodeRollup(winnerRollup))
		}
	})
}

type blockAndParent struct {
	b common.EncodedBlock
	p common.EncodedBlock
}

// RPCNewHead Receive notifications From the L1Node Node when there's a new block
func (a *Node) RPCNewHead(b common.EncodedBlock, p common.EncodedBlock) {
	if atomic.LoadInt32(a.interrupt) == 1 {
		return
	}
	a.blockRPCCh <- blockAndParent{b, p}
}

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

func (a *Node) P2PReceiveTx(tx obscuroCommon.EncryptedTx) {
	if atomic.LoadInt32(a.interrupt) == 1 {
		return
	}
	// Ignore gossiped transactions while the node is still initialisng
	if a.Enclave.IsInitialised() {
		go a.Enclave.SubmitTx(tx)
	}
}

func (a *Node) RPCBalance(address common.Address) uint64 {
	return a.Enclave.Balance(address)
}

func (a *Node) Stop() {
	// block all requests
	atomic.StoreInt32(a.interrupt, 1)
	a.Enclave.Stop()
	time.Sleep(time.Millisecond * 10)
	a.exitNodeCh <- true
}

// Called only by the first enclave to bootstrap the network
func (a *Node) initialiseProtocol() common.L2RootHash {
	// Create the genesis rollup and submit it to the MC
	genesis := a.Enclave.ProduceGenesis()
	a.broadcastTx(common.L1Tx{ID: uuid.New(), TxType: common.RollupTx, Rollup: obscuroCommon.EncodeRollup(genesis.Rollup.ToRollup())})

	return genesis.Hash
}

func (a *Node) broadcastTx(tx common.L1Tx) {
	t, err := tx.Encode()
	if err != nil {
		panic(err)
	}
	a.L1Node.BroadcastTx(t)
}

// This method implements the procedure by which a node obtains the secret
func (a *Node) requestSecret() {
	a.broadcastTx(common.L1Tx{
		ID:          uuid.New(),
		TxType:      common.RequestSecretTx,
		Attestation: a.Enclave.Attestation(),
	})

	// start listening for l1 blocks that contain the response to the request
	for {
		select {
		case b := <-a.blockRPCCh:
			txs := b.b.DecodeBlock().Transactions
			for _, tx := range txs {
				if tx.TxType == common.StoreSecretTx && tx.Attestation.Owner == a.ID {
					// someone has replied
					a.Enclave.Init(tx.Secret)
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
	for _, tx := range b.Transactions {
		if tx.TxType == common.RequestSecretTx {
			a.broadcastTx(common.L1Tx{
				ID:          uuid.New(),
				TxType:      common.StoreSecretTx,
				Secret:      a.Enclave.FetchSecret(tx.Attestation),
				Attestation: tx.Attestation,
			})
		}
	}
}

func NewAgg(
	id common.NodeID,
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
	}
}
