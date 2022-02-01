package ethereum_mock

import (
	"fmt"
	"simulation/common"
	"sync/atomic"
	"time"
)

type L1Network interface {
	BroadcastBlock(b common.EncodedBlock)
	BroadcastTx(tx common.EncodedL1Tx)
}

type NotifyNewBlock interface {
	RPCNewHead(b common.EncodedBlock)
}

type MiningConfig struct {
	PowTime common.Latency
}

type Node struct {
	Id             common.NodeId
	cfg            MiningConfig
	clients        []NotifyNewBlock
	network        L1Network
	mining         bool
	statsCollector common.StatsCollector

	// Channels
	exitCh       chan bool // the Node stops
	exitMiningCh chan bool // the mining loop is notified to stop
	interrupt    *int32

	p2pCh       chan common.Block // this is where blocks received from peers are dropped
	miningCh    chan common.Block // this is where blocks created by the mining setup of the current node are dropped
	canonicalCh chan common.Block // this is where the main processing routine drops blocks that are canonical
	mempoolCh   chan common.L1Tx  // where l1 transactions to be published in the next block are added
}

func NewMiner(id common.NodeId, cfg MiningConfig, client NotifyNewBlock, network L1Network, statsCollector common.StatsCollector) Node {
	return Node{
		Id:             id,
		mining:         true,
		cfg:            cfg,
		statsCollector: statsCollector,
		clients:        []NotifyNewBlock{client},
		network:        network,
		exitCh:         make(chan bool),
		exitMiningCh:   make(chan bool),
		interrupt:      new(int32),
		p2pCh:          make(chan common.Block),
		miningCh:       make(chan common.Block),
		canonicalCh:    make(chan common.Block),
		mempoolCh:      make(chan common.L1Tx)}
}

// Start runs an infinite loop that listens to the two block producing channels and processes them.
// it outputs the winning blocks to the roundWinnerCh channel
func (m *Node) Start() {
	if m.mining {
		// This starts the mining
		go m.startMining()
	}

	var head = m.setHead(common.GenesisBlock)

	for {
		select {
		case p2pb := <-m.p2pCh: // Received from peers
			if p2pb.Height() > head.Height() {
				// Check for Reorgs
				if !common.IsAncestor(head, p2pb) {
					m.statsCollector.L1Reorg(m.Id)
					fork := common.LCA(head, p2pb)
					common.Log(fmt.Sprintf("> M%d: L1Reorg new=b_%d(%d), old=b_%d(%d), fork=b_%d(%d)", m.Id, p2pb.RootHash.ID(), p2pb.Height(), head.RootHash.ID(), head.Height(), fork.Root().ID(), fork.Height()))
				}
				head = m.setHead(p2pb)
			}
		case mb := <-m.miningCh: // Received from the local mining
			if mb.Height() > head.Height() { // Ignore the locally produced block if someone else found one already
				common.Log(m.printBlock(mb))
				head = m.setHead(mb)
				m.network.BroadcastBlock(encodeBlock(mb))
			}
		case <-m.exitCh:
			return
		}
	}
}

func (m *Node) printBlock(mb common.Block) string {
	// This is just for printing
	var txs []string
	for _, tx := range mb.L1Txs() {
		if tx.TxType == common.RollupTx {
			txs = append(txs, fmt.Sprintf("r_%d", tx.Rollup.Root().ID()))
		} else {
			txs = append(txs, fmt.Sprintf("deposit(%v=%d)", tx.Dest, tx.Amount))
		}
	}
	return fmt.Sprintf("> M%d: create b_%d(Height=%d, Nonce=%d)[p=b_%d]. Txs: %v", m.Id, mb.RootHash.ID(), mb.Height(), mb.Nonce, mb.Parent().Root().ID(), txs)
}

// Notifies the Miner to start mining on the new block and the aggregtor to produce rollups
func (m *Node) setHead(b common.Block) common.Block {
	if atomic.LoadInt32(m.interrupt) == 1 {
		return b
	}

	// notify the clients
	for _, c := range m.clients {
		ser := encodeBlock(b)
		c.RPCNewHead(ser)
	}
	m.canonicalCh <- b
	return b
}

// P2PReceiveBlock is called by counterparties when there is a block to broadcast
// All it does is drop the blocks in a channel for processing.
func (m *Node) P2PReceiveBlock(b common.EncodedBlock) {
	if atomic.LoadInt32(m.interrupt) == 1 {
		return
	}
	bl := decodeBlock(b)
	m.p2pCh <- bl
}

// startMining - listens on the canonicalCh and schedule a go routine that produces a block after a PowTime and drop it on the miningCh channel
func (m *Node) startMining() {

	// stores all transactions seen from the beginning of time.
	var mempool = make([]common.Tx, 0)
	var doneCh *chan bool = nil

	for {
		select {
		case <-m.exitMiningCh:
			return
		case tx := <-m.mempoolCh:
			mempool = append(mempool, tx)

		case cb := <-m.canonicalCh:
			// A new canonical block was found. Start a new round based on that block.

			// remove transactions that are already considered committed
			mempool = common.RemoveCommittedTransactions(cb, mempool)

			//notify the existing mining go routine to stop mining
			if doneCh != nil {
				*doneCh <- true
			}

			c := make(chan bool)
			doneCh = &c
			// Generate a random number, and wait for that number of ms. Equivalent to PoW
			// Include all rollups received during this period.
			nonce := m.cfg.PowTime()
			common.ScheduleInterrupt(nonce, doneCh, func() {
				toInclude := common.FindNotIncludedTxs(cb, mempool)
				txsCopy := make([]common.L1Tx, len(toInclude))
				for i, tx := range toInclude {
					txsCopy[i] = tx.(common.L1Tx)
				}
				// todo - iterate through the rollup transactions and include only the ones with the proof on the canonical chain
				if atomic.LoadInt32(m.interrupt) == 1 {
					return
				}
				m.miningCh <- common.NewBlock(&cb, nonce, m.Id, txsCopy)
			})
		}
	}
}

// P2PGossipTx receive rollups to publish from the linked aggregators
func (m *Node) P2PGossipTx(tx common.EncodedL1Tx) {
	if atomic.LoadInt32(m.interrupt) == 1 {
		return
	}
	t, err := tx.Decode()
	if err != nil {
		panic(err)
	}

	m.mempoolCh <- t
}

func (m *Node) BroadcastTx(tx common.EncodedL1Tx) {
	m.network.BroadcastTx(tx)
}

func (m *Node) Stop() {
	//block all requests
	atomic.StoreInt32(m.interrupt, 1)
	time.Sleep(time.Millisecond * 100)

	m.exitMiningCh <- true
	m.exitCh <- true
}

func decodeBlock(b common.EncodedBlock) common.Block {
	bl, err := b.Decode()
	if err != nil {
		panic(err)
	}
	return bl
}
func encodeBlock(b common.Block) common.EncodedBlock {
	ser, err := b.Encode()
	if err != nil {
		panic(err)
	}
	return ser
}
