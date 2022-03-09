package ethereummock

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"sync/atomic"
	"time"

	common2 "github.com/obscuronet/obscuro-playground/go/common"
)

type L1Network interface {
	// BroadcastBlock - send the block and the parent to make sure there are no gaps
	BroadcastBlock(b common2.EncodedBlock, p common2.EncodedBlock)
	BroadcastTx(tx common2.EncodedL1Tx)
}

type MiningConfig struct {
	PowTime common2.Latency
}

type TxDB interface {
	Txs(block *common2.Block) (map[common2.TxHash]*common2.L1Tx, bool)
	AddTxs(*common2.Block, map[common2.TxHash]*common2.L1Tx)
}

type StatsCollector interface {
	// Register when a miner has to process a reorg (a winning block from a fork)
	L1Reorg(id common2.NodeID)
}

type Node struct {
	ID       common2.NodeID
	cfg      MiningConfig
	clients  []common2.NotifyNewBlock
	network  L1Network
	mining   bool
	stats    StatsCollector
	Resolver common2.BlockResolver
	db       TxDB

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
	// common2.Log(fmt.Sprintf("Starting miner %d..", m.ID))
	if m.mining {
		// This starts the mining
		go m.startMining()
	}

	head := m.setHead(common2.GenesisBlock)
	m.Resolver.Store(common2.GenesisBlock)

	for {
		select {
		case p2pb := <-m.p2pCh: // Received from peers
			_, received := m.Resolver.Resolve(p2pb.Hash())
			// only process blocks if they haven't been processed before
			if !received {
				head = m.processBlock(p2pb, head)
			}

		case mb := <-m.miningCh: // Received from the local mining
			head = m.processBlock(mb, head)
			if head.Hash() == mb.Hash() { // Ignore the locally produced block if someone else found one already
				p, found := m.Resolver.Parent(mb)
				if !found {
					panic("noo")
				}
				m.network.BroadcastBlock(common2.EncodeBlock(mb), common2.EncodeBlock(p))
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
	_, f := m.Resolver.Resolve(b.Header().ParentHash)

	// only proceed if the parent is available
	if !f {
		common2.Log(fmt.Sprintf("> M%d: Parent block not found=b_%s", m.ID, common2.Str(b.Header().ParentHash)))
		return head
	}

	// Ignore superseeded blocks
	if m.Resolver.Height(b) <= m.Resolver.Height(head) {
		return head
	}

	// Check for Reorgs
	if !common2.IsAncestor(head, b, m.Resolver) {
		m.stats.L1Reorg(m.ID)
		fork := LCA(head, b, m.Resolver)
		common2.Log(fmt.Sprintf("> M%d: L1Reorg new=b_%s(%d), old=b_%s(%d), fork=b_%s(%d)", m.ID, common2.Str(b.Hash()), m.Resolver.Height(b), common2.Str(head.Hash()), m.Resolver.Height(head), common2.Str(fork.Hash()), m.Resolver.Height(fork)))
		return m.setFork(BlocksBetween(fork, b, m.Resolver))
	}

	if m.Resolver.Height(b) > (m.Resolver.Height(head) + 1) {
		panic(fmt.Sprintf("> M%d: Should not happen", m.ID))
	}

	return m.setHead(b)
}

// Notifies the Miner to start mining on the new block and the aggregtor to produce rollups
func (m *Node) setHead(b *common2.Block) *common2.Block {
	if atomic.LoadInt32(m.interrupt) == 1 {
		return b
	}

	// notify the clients
	for _, c := range m.clients {
		t := c
		if m.Resolver.Height(b) == 0 {
			go t.RPCNewHead(common2.EncodeBlock(b), nil)
		} else {
			p, f := m.Resolver.Parent(b)
			if !f {
				panic("This should not happen")
			}
			go t.RPCNewHead(common2.EncodeBlock(b), common2.EncodeBlock(p))
		}
	}
	m.canonicalCh <- b

	return b
}

func (m *Node) setFork(blocks []*common2.Block) *common2.Block {
	head := blocks[len(blocks)-1]
	if atomic.LoadInt32(m.interrupt) == 1 {
		return head
	}

	fork := make([]common2.EncodedBlock, len(blocks))
	for i, block := range blocks {
		fork[i] = common2.EncodeBlock(block)
	}

	// notify the clients
	for _, c := range m.clients {
		go c.RPCNewFork(fork)
	}
	m.canonicalCh <- head

	return head
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

// startMining - listens on the canonicalCh and schedule a go routine that produces a block after a PowTime and drop it
// on the miningCh channel
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

		case canonicalBlock := <-m.canonicalCh:
			// A new canonical block was found. Start a new round based on that block.

			// remove transactions that are already considered committed
			mempool = removeCommittedTransactions(canonicalBlock, mempool, m.Resolver, m.db)

			// notify the existing mining go routine to stop mining
			atomic.StoreInt32(interrupt, 1)
			c := int32(0)
			interrupt = &c

			// Generate a random number, and wait for that number of ms. Equivalent to PoW
			// Include all rollups received during this period.
			nonce := m.cfg.PowTime()
			common2.ScheduleInterrupt(nonce, interrupt, func() {
				toInclude := findNotIncludedTxs(canonicalBlock, mempool, m.Resolver, m.db)
				// todo - iterate through the rollup transactions and include only the ones with the proof on the canonical chain
				if atomic.LoadInt32(m.interrupt) == 1 {
					return
				}
				b := common2.NewBlock(canonicalBlock, nonce, common.Address(m.ID), toInclude)
				m.miningCh <- b
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
	return BlocksBetween(common2.GenesisBlock, h, m.Resolver)
}

func (m *Node) Stop() {
	// block all requests
	atomic.StoreInt32(m.interrupt, 1)
	time.Sleep(time.Millisecond * 100)

	m.exitMiningCh <- true
	m.exitCh <- true
}

func NewMiner(
	id common2.NodeID,
	cfg MiningConfig,
	client common2.NotifyNewBlock,
	network L1Network,
	statsCollector StatsCollector,
) Node {
	return Node{
		ID:           id,
		mining:       true,
		cfg:          cfg,
		stats:        statsCollector,
		Resolver:     NewResolver(),
		db:           NewTxDB(),
		clients:      []common2.NotifyNewBlock{client},
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
