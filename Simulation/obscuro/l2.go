package obscuro

import (
	"github.com/google/uuid"
	"simulation/common"
	"simulation/ethereum-mock"
	wallet_mock "simulation/wallet-mock"
	"sync/atomic"
	"time"
)

type AggregatorCfg struct {
	GossipPeriod uint64
}

type L2Network interface {
	BroadcastRollup(r common.EncodedRollup)
	BroadcastTx(tx common.EncodedL2Tx)
}

// Node this will become the Obscuro "Node" type
type Node struct {
	Id common.NodeId

	l2Network L2Network
	L1Node    *ethereum_mock.Node

	mining bool // true -if this is an aggregator, false if it is a validator
	cfg    AggregatorCfg

	statsCollector common.StatsCollector

	// control the lifecycle
	exitNodeCh chan bool
	interrupt  *int32

	// where the connected L1Node node drops new blocks
	blockRpcCh chan common.EncodedBlock

	Enclave Enclave
}

func NewAgg(id common.NodeId, cfg AggregatorCfg, l1 *ethereum_mock.Node, l2Network L2Network, collector common.StatsCollector) Node {
	return Node{
		// config
		Id:        id,
		cfg:       cfg,
		mining:    true,
		L1Node:    l1,
		l2Network: l2Network,

		statsCollector: collector,

		// lifecycle channels
		exitNodeCh: make(chan bool),
		interrupt:  new(int32),

		// incoming data
		blockRpcCh: make(chan common.EncodedBlock),

		// State processing
		Enclave: NewEnclave(id, true, collector),
	}
}

func (a *Node) Start() {
	// used as a signaling mechanism to stop processing the old block if a new L1 block arrives earlier
	var stopProcessingOldBlock *chan bool = nil

	go a.Enclave.Start()

	var currentHead common.RootHash
	// Main loop - Listen for notifications From the L1 node and process them
	// Note that during processing, more recent notifications can be received.
	for {
		select {
		case b := <-a.blockRpcCh:
			if a.mining {
				if stopProcessingOldBlock != nil {
					*stopProcessingOldBlock <- true
				}

				c := make(chan bool)
				stopProcessingOldBlock = &c
			}

			//common.Log(fmt.Sprintf(">   Agg%d: receive b_%d .", a.Id, common.DecodeBlock(b).RootHash.ID()))
			var rollup common.EncodedRollup
			rollup, currentHead = a.Enclave.SubmitBlock(b)
			//common.Log(fmt.Sprintf(">   Agg%d: produce r_%d .", a.Id, common.DecodeRollup(rollup).RootHash.ID()))

			if a.mining {
				a.l2Network.BroadcastRollup(rollup)

				common.ScheduleInterrupt(a.cfg.GossipPeriod, stopProcessingOldBlock, func() {
					if atomic.LoadInt32(a.interrupt) == 1 {
						return
					}
					// Request the round winner for the current head
					winnerRollup, submit := a.Enclave.RoundWinner(currentHead)
					if submit {
						tx := common.L1Tx{Id: uuid.New(), TxType: common.RollupTx, Rollup: winnerRollup}
						t, err := tx.Encode()
						if err != nil {
							panic(err)
						}
						a.L1Node.BroadcastTx(t)
					}
				})
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

// P2PGossipRollup is called by counterparties when there is a Rollup to broadcast
// All it does is forward the rollup for processing to the enclave
func (a *Node) P2PGossipRollup(r common.EncodedRollup) {
	if atomic.LoadInt32(a.interrupt) == 1 {
		return
	}
	go a.Enclave.SubmitRollup(r)
}

func (a *Node) P2PReceiveTx(tx common.EncodedL2Tx) {
	if atomic.LoadInt32(a.interrupt) == 1 {
		return
	}
	go a.Enclave.SubmitTx(tx)
}

func (a *Node) RPCBalance(address wallet_mock.Address) uint64 {
	return a.Enclave.Balance(address)
}

func (a *Node) Stop() {
	// block all requests
	atomic.StoreInt32(a.interrupt, 1)
	a.Enclave.Stop()
	time.Sleep(time.Millisecond * 10)
	a.exitNodeCh <- true
}
