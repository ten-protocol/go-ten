package obscuro

import (
	"fmt"
	"github.com/google/uuid"
	"simulation/common"
	"simulation/ethereum-mock"
	"sync/atomic"
	"time"
)

type AggregatorCfg struct {
	GossipPeriod uint64
}

type L2Network interface {
	BroadcastRollup(r common.EncodedRollup)
	BroadcastTx(tx EncodedL2Tx)
}

type StatsCollector interface {
	L2Recalc(id common.NodeId)
	NewBlock(block common.Block)
	NewRollup(rollup Rollup)
}

// Node this will become the Obscuro "Node" type
type Node struct {
	Id common.NodeId

	l2Network L2Network
	L1Node    *ethereum_mock.Node

	mining bool // true -if this is an aggregator, false if it is a validator
	cfg    AggregatorCfg

	stats StatsCollector

	// control the lifecycle
	exitNodeCh chan bool
	interrupt  *int32

	// where the connected L1Node node drops new blocks
	blockRpcCh chan common.EncodedBlock
	forkRpcCh  chan []common.EncodedBlock

	Enclave Enclave
}

func (a *Node) Start() {
	// used as a signaling mechanism to stop processing the old block if a new L1 block arrives earlier
	z := int32(0)
	var interrupt = &z

	go a.Enclave.Start()

	// Main loop - Listen for notifications From the L1 node and process them
	// Note that during processing, more recent notifications can be received.
	for {
		select {
		case b := <-a.blockRpcCh:
			if a.mining {
				atomic.StoreInt32(interrupt, 1)
				c := int32(0)
				interrupt = &c
			}

			result := a.Enclave.SubmitBlock(b)

			if result.processed {
				if a.mining {
					a.l2Network.BroadcastRollup(result.rollup)

					common.ScheduleInterrupt(a.cfg.GossipPeriod, interrupt, func() {
						if atomic.LoadInt32(a.interrupt) == 1 {
							return
						}
						// Request the round winner for the current head
						winnerRollup, submit := a.Enclave.RoundWinner(result.root)
						if submit {
							tx := common.L1Tx{Id: uuid.New(), TxType: common.RollupTx, Rollup: winnerRollup}
							t, err := tx.Encode()
							if err != nil {
								panic(err)
							}
							a.L1Node.BroadcastTx(t)
							// collect Stats
							a.stats.NewRollup(DecodeRollup(winnerRollup))
						}
					})
				}
			} else {
				common.Log(fmt.Sprintf(">   Agg%d: Could not process block b_%d", a.Id, b.DecodeBlock().RootHash.ID()))
			}

		case b := <-a.forkRpcCh:
			if a.mining {
				atomic.StoreInt32(interrupt, 1)
				c := int32(0)
				interrupt = &c
			}

			var result SubmitBlockResponse
			for _, block := range b {
				result = a.Enclave.SubmitBlock(block)
			}

			if result.processed {
				if a.mining {
					a.l2Network.BroadcastRollup(result.rollup)

					common.ScheduleInterrupt(a.cfg.GossipPeriod, interrupt, func() {
						if atomic.LoadInt32(a.interrupt) == 1 {
							return
						}
						// Request the round winner for the current head
						winnerRollup, submit := a.Enclave.RoundWinner(result.root)
						if submit {
							tx := common.L1Tx{Id: uuid.New(), TxType: common.RollupTx, Rollup: winnerRollup}
							t, err := tx.Encode()
							if err != nil {
								panic(err)
							}
							a.L1Node.BroadcastTx(t)
							// collect Stats
							a.stats.NewRollup(DecodeRollup(winnerRollup))
						}
					})
				}
			} else {
				common.Log(fmt.Sprintf(">   Agg%d: Could not process fork ", a.Id))
			}

		case <-a.exitNodeCh:
			a.Enclave.Stop()
			return
		}
	}
}

// RPCNewHead Receive notifications From the L1Node Node when there's a new block
func (a *Node) RPCNewHead(b common.EncodedBlock) {
	if atomic.LoadInt32(a.interrupt) == 1 {
		return
	}
	a.blockRpcCh <- b
}

func (a *Node) RPCNewFork(b []common.EncodedBlock) {
	if atomic.LoadInt32(a.interrupt) == 1 {
		return
	}
	a.forkRpcCh <- b
}

// P2PGossipRollup is called by counterparties when there is a Rollup to broadcast
// All it does is forward the rollup for processing to the enclave
func (a *Node) P2PGossipRollup(r common.EncodedRollup) {
	if atomic.LoadInt32(a.interrupt) == 1 {
		return
	}
	go a.Enclave.SubmitRollup(r)
}

func (a *Node) P2PReceiveTx(tx EncodedL2Tx) {
	if atomic.LoadInt32(a.interrupt) == 1 {
		return
	}
	go a.Enclave.SubmitTx(tx)
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

func NewAgg(id common.NodeId, cfg AggregatorCfg, l1 *ethereum_mock.Node, l2Network L2Network, collector StatsCollector) Node {
	return Node{
		// config
		Id:        id,
		cfg:       cfg,
		mining:    true,
		L1Node:    l1,
		l2Network: l2Network,

		stats: collector,

		// lifecycle channels
		exitNodeCh: make(chan bool),
		interrupt:  new(int32),

		// incoming data
		blockRpcCh: make(chan common.EncodedBlock),
		forkRpcCh:  make(chan []common.EncodedBlock),

		// State processing
		Enclave: NewEnclave(id, true, collector),
	}
}
