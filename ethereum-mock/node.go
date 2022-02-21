package ethereum_mock

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/otherview/obscuro-playground/common"
)

type L1Network interface {
	// BroadcastBlock - send the block and the parent to make sure there are no gaps
	BroadcastBlock(b common.EncodedBlock, p common.EncodedBlock)
	BroadcastTx(tx common.EncodedL1Tx)
}

type NotifyNewBlock interface {
	RPCNewHead(b common.EncodedBlock)
	RPCNewFork(b []common.EncodedBlock)
}

type MiningConfig struct {
	PowTime common.Latency
}

type TxDb interface {
	Txs(block *common.Block) (map[common.TxHash]*common.L1Tx, bool)
	AddTxs(*common.Block, map[common.TxHash]*common.L1Tx)
}

type StatsCollector interface {
	// Register when a miner has to process a reorg (a winning block from a fork)
	L1Reorg(id common.NodeId)
}

type Node struct {
	Id       common.NodeId
	cfg      MiningConfig
	clients  []NotifyNewBlock
	network  L1Network
	mining   bool
	stats    StatsCollector
	Resolver common.BlockResolver
	db       TxDb

	// Channels
	exitCh       chan bool // the Node stops
	exitMiningCh chan bool // the mining loop is notified to stop
	interrupt    *int32

	p2pCh       chan *common.Block // this is where blocks received from peers are dropped
	miningCh    chan *common.Block // this is where blocks created by the mining setup of the current node are dropped
	canonicalCh chan *common.Block // this is where the main processing routine drops blocks that are canonical
	mempoolCh   chan common.L1Tx   // where l1 transactions to be published in the next block are added

	// internal
	headInCh  chan bool
	headOutCh chan *common.Block
}

// Start runs an infinite loop that listens to the two block producing channels and processes them.
// it outputs the winning blocks to the roundWinnerCh channel
func (m *Node) Start() {
	if m.mining {
		// This starts the mining
		go m.startMining()
	}

	head := m.setHead(&common.GenesisBlock)
	m.Resolver.Store(&common.GenesisBlock)

	for {
		select {
		case p2pb := <-m.p2pCh: // Received from peers
			_, received := m.Resolver.Resolve(p2pb.Hash())
			// only process blocks if they haven't been processed before
			if !received {
				head = m.processBlock(p2pb, head)
			}

		case mb := <-m.miningCh: // Received from the local mining
			m.Resolver.Store(mb)
			if mb.Height(m.Resolver) > head.Height(m.Resolver) { // Ignore the locally produced block if someone else found one already
				head = m.setHead(mb)
				p, _ := mb.Parent(m.Resolver)
				m.network.BroadcastBlock(mb.EncodeBlock(), p.EncodeBlock())
			}
		case <-m.headInCh:
			m.headOutCh <- head
		case <-m.exitCh:
			return
		}
	}
}

func (m *Node) processBlock(b *common.Block, head *common.Block) *common.Block {
	m.Resolver.Store(b)
	_, f := m.Resolver.Resolve(b.Header.ParentHash)
	// only proceed if the parent is available
	if f {
		if b.Height(m.Resolver) > head.Height(m.Resolver) {
			// Check for Reorgs
			if !common.IsAncestor(head, b, m.Resolver) {
				m.stats.L1Reorg(m.Id)
				fork := LCA(head, b, m.Resolver)
				common.Log(fmt.Sprintf("> M%d: L1Reorg new=b_%s(%d), old=b_%s(%d), fork=b_%s(%d)", m.Id, common.Str(b.Hash()), b.Height(m.Resolver), common.Str(head.Hash()), head.Height(m.Resolver), common.Str(fork.Hash()), fork.Height(m.Resolver)))
				head = m.setFork(BlocksBetween(fork, b, m.Resolver))
			} else {
				head = m.setHead(b)
			}
		}
	} else {
		common.Log(fmt.Sprintf("> M%d: Not found=b_%s", m.Id, common.Str(b.Header.ParentHash)))
	}
	return head
}

// Notifies the Miner to start mining on the new block and the aggregtor to produce rollups
func (m *Node) setHead(b *common.Block) *common.Block {
	if atomic.LoadInt32(m.interrupt) == 1 {
		return b
	}

	// notify the clients
	for _, c := range m.clients {
		t := c
		go t.RPCNewHead(b.EncodeBlock())
	}
	m.canonicalCh <- b
	return b
}

func (m *Node) setFork(blocks []*common.Block) *common.Block {
	h := blocks[len(blocks)-1]
	if atomic.LoadInt32(m.interrupt) == 1 {
		return h
	}
	encoded := make([]common.EncodedBlock, len(blocks))
	for i, block := range blocks {
		encoded[i] = block.EncodeBlock()
	}
	// notify the clients
	for _, c := range m.clients {
		c.RPCNewFork(encoded)
	}
	m.canonicalCh <- h
	return h
}

// P2PReceiveBlock is called by counterparties when there is a block to broadcast
// All it does is drop the blocks in a channel for processing.
func (m *Node) P2PReceiveBlock(b common.EncodedBlock, p common.EncodedBlock) {
	if atomic.LoadInt32(m.interrupt) == 1 {
		return
	}
	m.p2pCh <- p.DecodeBlock()
	m.p2pCh <- b.DecodeBlock()
}

// startMining - listens on the canonicalCh and schedule a go routine that produces a block after a PowTime and drop it on the miningCh channel
func (m *Node) startMining() {
	// stores all transactions seen from the beginning of time.
	mempool := make([]*common.L1Tx, 0)
	z := int32(0)
	interrupt := &z

	for {
		select {
		case <-m.exitMiningCh:
			return
		case tx := <-m.mempoolCh:
			mempool = append(mempool, &tx)

		case cb := <-m.canonicalCh:
			// A new canonical block was found. Start a new round based on that block.

			// remove transactions that are already considered committed
			mempool = removeCommittedTransactions(cb, mempool, m.Resolver, m.db)

			// notify the existing mining go routine to stop mining
			atomic.StoreInt32(interrupt, 1)
			c := int32(0)
			interrupt = &c

			// Generate a random number, and wait for that number of ms. Equivalent to PoW
			// Include all rollups received during this period.
			nonce := m.cfg.PowTime()
			common.ScheduleInterrupt(nonce, interrupt, func() {
				toInclude := findNotIncludedTxs(cb, mempool, m.Resolver, m.db)
				// todo - iterate through the rollup transactions and include only the ones with the proof on the canonical chain
				if atomic.LoadInt32(m.interrupt) == 1 {
					return
				}
				b := common.NewBlock(cb, nonce, m.Id, toInclude)
				m.miningCh <- &b
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

func (m *Node) RPCBlockchainFeed() []*common.Block {
	m.headInCh <- true
	h := <-m.headOutCh
	return BlocksBetween(&common.GenesisBlock, h, m.Resolver)
}

func (m *Node) Stop() {
	// block all requests
	atomic.StoreInt32(m.interrupt, 1)
	time.Sleep(time.Millisecond * 100)

	m.exitMiningCh <- true
	m.exitCh <- true
}

func NewMiner(id common.NodeId, cfg MiningConfig, client NotifyNewBlock, network L1Network, statsCollector StatsCollector) Node {
	return Node{
		Id:           id,
		mining:       true,
		cfg:          cfg,
		stats:        statsCollector,
		Resolver:     NewResolver(),
		db:           NewTxDb(),
		clients:      []NotifyNewBlock{client},
		network:      network,
		exitCh:       make(chan bool),
		exitMiningCh: make(chan bool),
		interrupt:    new(int32),
		p2pCh:        make(chan *common.Block),
		miningCh:     make(chan *common.Block),
		canonicalCh:  make(chan *common.Block),
		mempoolCh:    make(chan common.L1Tx),
		headInCh:     make(chan bool),
		headOutCh:    make(chan *common.Block),
	}
}
