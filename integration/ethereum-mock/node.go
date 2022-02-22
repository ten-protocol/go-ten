package ethereum_mock

import (
	"fmt"
	common2 "github.com/obscuronet/obscuro-playground/go/common"
	"sync/atomic"
	"time"
)

type L1Network interface {
	// BroadcastBlock - send the block and the parent to make sure there are no gaps
	BroadcastBlock(b common2.EncodedBlock, p common2.EncodedBlock)
	BroadcastTx(tx common2.EncodedL1Tx)
}

type NotifyNewBlock interface {
	RPCNewHead(b common2.EncodedBlock)
	RPCNewFork(b []common2.EncodedBlock)
}

type MiningConfig struct {
	PowTime common2.Latency
}

type TxDb interface {
	Txs(block *common2.Block) (map[common2.TxHash]*common2.L1Tx, bool)
	AddTxs(*common2.Block, map[common2.TxHash]*common2.L1Tx)
}

type StatsCollector interface {
	// Register when a miner has to process a reorg (a winning block from a fork)
	L1Reorg(id common2.NodeId)
}

type Node struct {
	Id       common2.NodeId
	cfg      MiningConfig
	clients  []NotifyNewBlock
	network  L1Network
	mining   bool
	stats    StatsCollector
	Resolver common2.BlockResolver
	db       TxDb

	// Channels
	exitCh       chan bool // the Node stops
	exitMiningCh chan bool // the mining loop is notified to stop
	interrupt    *int32

	p2pCh       chan *common2.Block // this is where blocks received from peers are dropped
	miningCh    chan *common2.Block // this is where blocks created by the mining setup of the current node are dropped
	canonicalCh chan *common2.Block // this is where the main processing routine drops blocks that are canonical
	mempoolCh   chan common2.L1Tx   // where l1 transactions to be published in the next block are added

	// internal
	headInCh  chan bool
	headOutCh chan *common2.Block
}

// Start runs an infinite loop that listens to the two block producing channels and processes them.
// it outputs the winning blocks to the roundWinnerCh channel
func (m *Node) Start() {
	if m.mining {
		// This starts the mining
		go m.startMining()
	}

	head := m.setHead(&common2.GenesisBlock)
	m.Resolver.Store(&common2.GenesisBlock)

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

func (m *Node) processBlock(b *common2.Block, head *common2.Block) *common2.Block {
	m.Resolver.Store(b)
	_, f := m.Resolver.Resolve(b.Header.ParentHash)
	// only proceed if the parent is available
	if f {
		if b.Height(m.Resolver) > head.Height(m.Resolver) {
			// Check for Reorgs
			if !common2.IsAncestor(head, b, m.Resolver) {
				m.stats.L1Reorg(m.Id)
				fork := LCA(head, b, m.Resolver)
				common2.Log(fmt.Sprintf("> M%d: L1Reorg new=b_%s(%d), old=b_%s(%d), fork=b_%s(%d)", m.Id, common2.Str(b.Hash()), b.Height(m.Resolver), common2.Str(head.Hash()), head.Height(m.Resolver), common2.Str(fork.Hash()), fork.Height(m.Resolver)))
				head = m.setFork(BlocksBetween(fork, b, m.Resolver))
			} else {
				head = m.setHead(b)
			}
		}
	} else {
		common2.Log(fmt.Sprintf("> M%d: Not found=b_%s", m.Id, common2.Str(b.Header.ParentHash)))
	}
	return head
}

// Notifies the Miner to start mining on the new block and the aggregtor to produce rollups
func (m *Node) setHead(b *common2.Block) *common2.Block {
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

func (m *Node) setFork(blocks []*common2.Block) *common2.Block {
	h := blocks[len(blocks)-1]
	if atomic.LoadInt32(m.interrupt) == 1 {
		return h
	}
	encoded := make([]common2.EncodedBlock, len(blocks))
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
func (m *Node) P2PReceiveBlock(b common2.EncodedBlock, p common2.EncodedBlock) {
	if atomic.LoadInt32(m.interrupt) == 1 {
		return
	}
	m.p2pCh <- p.DecodeBlock()
	m.p2pCh <- b.DecodeBlock()
}

// startMining - listens on the canonicalCh and schedule a go routine that produces a block after a PowTime and drop it on the miningCh channel
func (m *Node) startMining() {
	// stores all transactions seen from the beginning of time.
	mempool := make([]*common2.L1Tx, 0)
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
			common2.ScheduleInterrupt(nonce, interrupt, func() {
				toInclude := findNotIncludedTxs(cb, mempool, m.Resolver, m.db)
				// todo - iterate through the rollup transactions and include only the ones with the proof on the canonical chain
				if atomic.LoadInt32(m.interrupt) == 1 {
					return
				}
				b := common2.NewBlock(cb, nonce, m.Id, toInclude)
				m.miningCh <- &b
			})
		}
	}
}

// P2PGossipTx receive rollups to publish from the linked aggregators
func (m *Node) P2PGossipTx(tx common2.EncodedL1Tx) {
	if atomic.LoadInt32(m.interrupt) == 1 {
		return
	}
	t, err := tx.Decode()
	if err != nil {
		panic(err)
	}

	m.mempoolCh <- t
}

func (m *Node) BroadcastTx(tx common2.EncodedL1Tx) {
	m.network.BroadcastTx(tx)
}

func (m *Node) RPCBlockchainFeed() []*common2.Block {
	m.headInCh <- true
	h := <-m.headOutCh
	return BlocksBetween(&common2.GenesisBlock, h, m.Resolver)
}

func (m *Node) Stop() {
	// block all requests
	atomic.StoreInt32(m.interrupt, 1)
	time.Sleep(time.Millisecond * 100)

	m.exitMiningCh <- true
	m.exitCh <- true
}

func NewMiner(id common2.NodeId, cfg MiningConfig, client NotifyNewBlock, network L1Network, statsCollector StatsCollector) Node {
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
		p2pCh:        make(chan *common2.Block),
		miningCh:     make(chan *common2.Block),
		canonicalCh:  make(chan *common2.Block),
		mempoolCh:    make(chan common2.L1Tx),
		headInCh:     make(chan bool),
		headOutCh:    make(chan *common2.Block),
	}
}
