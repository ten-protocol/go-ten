package ethereummock

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave"

	"github.com/obscuronet/obscuro-playground/go/obscurocommon"

	"github.com/obscuronet/obscuro-playground/go/log"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum/common"
)

type L1Network interface {
	// BroadcastBlock - send the block and the parent to make sure there are no gaps
	BroadcastBlock(b obscurocommon.EncodedBlock, p obscurocommon.EncodedBlock)
	BroadcastTx(tx obscurocommon.EncodedL1Tx)
}

type MiningConfig struct {
	PowTime obscurocommon.Latency
}

type TxDB interface {
	Txs(block *types.Block) (map[obscurocommon.TxHash]*obscurocommon.L1Tx, bool)
	AddTxs(*types.Block, map[obscurocommon.TxHash]*obscurocommon.L1Tx)
}

type StatsCollector interface {
	// Register when a miner has to process a reorg (a winning block from a fork)
	L1Reorg(id common.Address)
}

type Node struct {
	ID       common.Address
	cfg      MiningConfig
	clients  []obscurocommon.NotifyNewBlock
	network  L1Network
	mining   bool
	stats    StatsCollector
	Resolver enclave.BlockResolver
	db       TxDB

	// Channels
	exitCh       chan bool // the Node stops
	exitMiningCh chan bool // the mining loop is notified to stop
	interrupt    *int32

	p2pCh       chan *types.Block       // this is where blocks received from peers are dropped
	miningCh    chan *types.Block       // this is where blocks created by the mining setup of the current node are dropped
	canonicalCh chan *types.Block       // this is where the main processing routine drops blocks that are canonical
	mempoolCh   chan obscurocommon.L1Tx // where l1 transactions to be published in the next block are added

	// internal
	headInCh  chan bool
	headOutCh chan *types.Block
}

// Start runs an infinite loop that listens to the two block producing channels and processes them.
// it outputs the winning blocks to the roundWinnerCh channel
func (m *Node) Start() {
	// common.Log(fmt.Sprintf("Starting miner %d..", m.ID))
	if m.mining {
		// This starts the mining
		go m.startMining()
	}

	m.Resolver.StoreBlock(obscurocommon.GenesisBlock)
	head := m.setHead(obscurocommon.GenesisBlock)

	for {
		select {
		case p2pb := <-m.p2pCh: // Received from peers
			_, received := m.Resolver.FetchBlock(p2pb.Hash())
			// only process blocks if they haven't been processed before
			if !received {
				head = m.processBlock(p2pb, head)
			}

		case mb := <-m.miningCh: // Received from the local mining
			head = m.processBlock(mb, head)
			if head.Hash() == mb.Hash() { // Ignore the locally produced block if someone else found one already
				p, found := m.Resolver.ParentBlock(mb)
				if !found {
					panic("noo")
				}
				m.network.BroadcastBlock(obscurocommon.EncodeBlock(mb), obscurocommon.EncodeBlock(p))
			}
		case <-m.headInCh:
			m.headOutCh <- head
		case <-m.exitCh:
			return
		}
	}
}

func (m *Node) processBlock(b *types.Block, head *types.Block) *types.Block {
	m.Resolver.StoreBlock(b)
	_, f := m.Resolver.FetchBlock(b.Header().ParentHash)

	// only proceed if the parent is available
	if !f {
		log.Log(fmt.Sprintf("> M%d: Parent block not found=b_%d", obscurocommon.ShortAddress(m.ID), obscurocommon.ShortHash(b.Header().ParentHash)))
		return head
	}

	// Ignore superseded blocks
	if m.Resolver.HeightBlock(b) <= m.Resolver.HeightBlock(head) {
		return head
	}

	// Check for Reorgs
	if !m.Resolver.IsAncestor(b, head) {
		m.stats.L1Reorg(m.ID)
		fork := LCA(head, b, m.Resolver)
		log.Log(fmt.Sprintf("> M%d: L1Reorg new=b_%d(%d), old=b_%d(%d), fork=b_%d(%d)", obscurocommon.ShortAddress(m.ID), obscurocommon.ShortHash(b.Hash()), m.Resolver.HeightBlock(b), obscurocommon.ShortHash(head.Hash()), m.Resolver.HeightBlock(head), obscurocommon.ShortHash(fork.Hash()), m.Resolver.HeightBlock(fork)))
		return m.setFork(BlocksBetween(fork, b, m.Resolver))
	}

	if m.Resolver.HeightBlock(b) > (m.Resolver.HeightBlock(head) + 1) {
		panic(fmt.Sprintf("> M%d: Should not happen", obscurocommon.ShortAddress(m.ID)))
	}

	return m.setHead(b)
}

// Notifies the Miner to start mining on the new block and the aggregtor to produce rollups
func (m *Node) setHead(b *types.Block) *types.Block {
	if atomic.LoadInt32(m.interrupt) == 1 {
		return b
	}

	// notify the clients
	for _, c := range m.clients {
		t := c
		if m.Resolver.HeightBlock(b) == obscurocommon.L1GenesisHeight {
			go t.RPCNewHead(obscurocommon.EncodeBlock(b), nil)
		} else {
			p, f := m.Resolver.ParentBlock(b)
			if !f {
				panic("This should not happen")
			}
			go t.RPCNewHead(obscurocommon.EncodeBlock(b), obscurocommon.EncodeBlock(p))
		}
	}
	m.canonicalCh <- b

	return b
}

func (m *Node) setFork(blocks []*types.Block) *types.Block {
	head := blocks[len(blocks)-1]
	if atomic.LoadInt32(m.interrupt) == 1 {
		return head
	}

	fork := make([]obscurocommon.EncodedBlock, len(blocks))
	for i, block := range blocks {
		fork[i] = obscurocommon.EncodeBlock(block)
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
func (m *Node) P2PReceiveBlock(b obscurocommon.EncodedBlock, p obscurocommon.EncodedBlock) {
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
	mempool := make([]*obscurocommon.L1Tx, 0)
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
			obscurocommon.ScheduleInterrupt(nonce, interrupt, func() {
				toInclude := findNotIncludedTxs(canonicalBlock, mempool, m.Resolver, m.db)
				// todo - iterate through the rollup transactions and include only the ones with the proof on the canonical chain
				if atomic.LoadInt32(m.interrupt) == 1 {
					return
				}
				b := obscurocommon.NewBlock(canonicalBlock, nonce, m.ID, toInclude)
				m.miningCh <- b
			})
		}
	}
}

// P2PGossipTx receive rollups to publish from the linked aggregators
func (m *Node) P2PGossipTx(tx obscurocommon.EncodedL1Tx) {
	if atomic.LoadInt32(m.interrupt) == 1 {
		return
	}
	t, err := tx.Decode()
	if err != nil {
		panic(err)
	}

	m.mempoolCh <- t
}

func (m *Node) BroadcastTx(tx obscurocommon.EncodedL1Tx) {
	m.network.BroadcastTx(tx)
}

func (m *Node) RPCBlockchainFeed() []*types.Block {
	m.headInCh <- true
	h := <-m.headOutCh
	return BlocksBetween(obscurocommon.GenesisBlock, h, m.Resolver)
}

func (m *Node) Stop() {
	// block all requests
	atomic.StoreInt32(m.interrupt, 1)
	time.Sleep(time.Millisecond * 100)

	m.exitMiningCh <- true
	m.exitCh <- true
}

func NewMiner(
	id common.Address,
	cfg MiningConfig,
	client obscurocommon.NotifyNewBlock,
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
		clients:      []obscurocommon.NotifyNewBlock{client},
		network:      network,
		exitCh:       make(chan bool),
		exitMiningCh: make(chan bool),
		interrupt:    new(int32),
		p2pCh:        make(chan *types.Block),
		miningCh:     make(chan *types.Block),
		canonicalCh:  make(chan *types.Block),
		mempoolCh:    make(chan obscurocommon.L1Tx),
		headInCh:     make(chan bool),
		headOutCh:    make(chan *types.Block),
	}
}
