package obscuro

import (
	"simulation/common"
	"simulation/ethereum-mock"
	wallet_mock "simulation/wallet-mock"
	"sync/atomic"
	"time"
)

type AggregatorCfg struct {
	GossipPeriod int
}

type L2Network interface {
	BroadcastRollup(r common.Rollup)
	BroadcastTx(tx common.L2Tx)
}

// Node this will become the Obscuro "Node" type
type Node struct {
	Id common.NodeId

	l2Network L2Network
	L1Node    *ethereum_mock.Node

	mining bool
	cfg    AggregatorCfg

	statsCollector common.StatsCollector

	// control the lifecycle
	exitNodeCh        chan bool
	exitAggregatingCh chan bool
	interrupt         *int32

	// where rollups and transactions are gossipped From peers
	p2pChRollup chan *common.Rollup
	p2pChTx     chan common.L2Tx

	// where the connected L1Node node drops new blocks
	blockRpcCh chan common.Block

	// responds to balance requests
	//balanceRpcInCh  chan wallet_mock.Address
	//balanceRpcOutCh chan int

	// used for internal communication between the gossi agent and the processing agent
	// todo - probably can use a single channel
	rollupInCh  chan int
	rollupOutCh chan []*common.Rollup

	// used for internal communication between the gossip agent and the processing agent
	txsInCh  chan bool
	txsOutCh chan []common.Tx

	// communicate the speculative work done during a pobi round
	speculativeWorkInCh  chan bool
	speculativeWorkOutCh chan currentWork

	// when a POBI round finishes and a winner is discovered. Notifies the gossip actor to start eagerly processing transactions on top of this state.
	roundWinnerCh chan winner

	// a database of work already executed
	Db Db
}

// internal structure to pass information
type currentWork struct {
	r   common.Rollup
	s   State
	txs []common.L2Tx
}

// internal structure to pass information
type winner struct {
	r common.Rollup
	s State
}

func NewAgg(id common.NodeId, cfg AggregatorCfg, l1 *ethereum_mock.Node, l2Network L2Network, collector common.StatsCollector) Node {
	return Node{
		Id:                id,
		cfg:               cfg,
		l2Network:         l2Network,
		statsCollector:    collector,
		L1Node:            l1,
		mining:            true,
		exitNodeCh:        make(chan bool),
		exitAggregatingCh: make(chan bool),
		interrupt:         new(int32),
		p2pChRollup:       make(chan *common.Rollup),
		p2pChTx:           make(chan common.L2Tx),
		blockRpcCh:        make(chan common.Block),
		//balanceRpcInCh:    make(chan wallet_mock.Address),
		//balanceRpcOutCh:   make(chan int),
		rollupInCh:           make(chan int),
		rollupOutCh:          make(chan []*common.Rollup),
		txsInCh:              make(chan bool),
		txsOutCh:             make(chan []common.Tx),
		roundWinnerCh:        make(chan winner),
		speculativeWorkInCh:  make(chan bool),
		speculativeWorkOutCh: make(chan currentWork),
		Db:                   NewInMemoryDb(),
	}
}

func (a Node) Start() {
	if a.mining {
		go a.startAggregating()
	}

	// used as a signaling mechanism to stop processing the old block if a new L1 block arrives earlier
	var doneCh *chan bool = nil

	// Main loop - Listen for notifications From the L1 node and process them
	// Note that during processing, more recent notifications can be received.
	for {
		select {
		case b := <-a.blockRpcCh:
			if a.mining {
				if doneCh != nil {
					*doneCh <- true
				}

				c := make(chan bool)
				doneCh = &c

				go a.newPobiRound(b, doneCh)
			} else {
				// Observer L2 nodes only calculate the state
				go a.updateState(b)
			}
		case <-a.exitNodeCh:
			if doneCh != nil {
				*doneCh <- true
			}
			return
		}
	}
}

// RPCNewHead Receive notifications From the L1Node Node when there's a new block
func (a Node) RPCNewHead(b common.Block) {
	if atomic.LoadInt32(a.interrupt) == 1 {
		return
	}
	a.blockRpcCh <- b
}

// P2PGossipRollup is called by counterparties when there is a Rollup to broadcast
// All it does is drop the Rollups in a channel for processing.
func (a Node) P2PGossipRollup(r *common.Rollup) {
	if atomic.LoadInt32(a.interrupt) == 1 {
		return
	}
	a.p2pChRollup <- r
}

func (a Node) P2PReceiveTx(tx common.L2Tx) {
	if atomic.LoadInt32(a.interrupt) == 1 {
		return
	}
	a.p2pChTx <- tx
}

func (a Node) RPCBalance(address wallet_mock.Address) int {
	return a.Db.Balance(address)
}

func (a Node) Stop() {
	// block all requests
	atomic.StoreInt32(a.interrupt, 1)
	time.Sleep(time.Millisecond * 10)
	a.exitAggregatingCh <- true
	a.exitNodeCh <- true
}
