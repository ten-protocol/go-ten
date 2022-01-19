package obscuro

import (
	"fmt"
	"github.com/google/uuid"
	"sync"
	"time"
)

type L1MiningConfig struct {
	powTime Latency
}

type L1Miner struct {
	id         int
	cfg        L1MiningConfig
	aggregator *L2Agg
	network    *NetworkCfg

	// Channels
	runCh1      chan bool   // add false when the miner has to stop
	runCh2      chan bool   // the mining loop is notified to stop
	p2pCh       chan *Block // this is where blocks received from peers are dropped
	miningCh    chan *Block // this is where blocks created by the mining setup of the current node are dropped
	canonicalCh chan *Block // this is where the main processing routine drops blocks that are canonical
	mempoolCh   chan *L1Tx  // where l1 transactions to be published in the next block are added
}

// Just two types of relevant L1 transactions: Deposits and Rollups
type L1TxType int64

const (
	DepositTx L1TxType = iota
	RollupTx
)

type L1TxId = uuid.UUID
type L1Tx struct {
	id     L1TxId
	txType L1TxType
	rollup *Rollup
	amount int
	dest   Address
}

func NewMiner(id int, cfg L1MiningConfig, agg *L2Agg, network *NetworkCfg) L1Miner {
	return L1Miner{id: id, cfg: cfg, aggregator: agg, network: network, runCh1: make(chan bool), p2pCh: make(chan *Block), miningCh: make(chan *Block), canonicalCh: make(chan *Block), runCh2: make(chan bool), mempoolCh: make(chan *L1Tx)}
}

type L1RootHash = uuid.UUID
type Nonce = int

type Block struct {
	height       int
	rootHash     L1RootHash
	nonce        Nonce
	miner        *L1Miner
	parent       *Block
	creationTime time.Time
	txs          []*L1Tx
}

var GenesisBlock = Block{height: -1, rootHash: uuid.New(), nonce: 0, creationTime: time.Now(), txs: []*L1Tx{{id: uuid.New(), txType: RollupTx, rollup: &GenesisRollup}}}

// Start runs an infinite loop that listens to the two block producing channels and processes them.
// it outputs the winning blocks to the canonicalCh channel
func (m L1Miner) Start() {
	// This starts the mining
	go m.startMining()

	var head = m.setHead(&GenesisBlock)

	for {
		select {
		case p2pb := <-m.p2pCh: // Received from peers
			if p2pb.height > head.height {
				head = m.setHead(p2pb)
			}
		case mb := <-m.miningCh: // Received from the local mining
			if mb.height > head.height { // Ignore the locally produced block if someone else found one already
				head = m.setHead(mb)
				m.network.broadcastBlockL1(mb)
				m.network.f.WriteString(m.printBlock(mb))
			}
		case _ = <-m.runCh1:
			return
		}
	}
}

func (m L1Miner) printBlock(mb *Block) string {
	// This is just for printing
	var txs []string
	for _, tx := range mb.txs {
		if tx.txType == RollupTx {
			txs = append(txs, fmt.Sprintf("r%d", tx.rollup.rootHash.ID()))
		} else {
			txs = append(txs, fmt.Sprintf("d%d=%d", tx.dest.ID(), tx.amount))
		}
	}
	return fmt.Sprintf("> M%d create b%d(height=%d, nonce=%d)[p=b%d]. Txs: %v\n", m.id, mb.rootHash.ID(), mb.height, mb.nonce, mb.parent.rootHash.ID(), txs)
}

// Notifies the miner to start mining on the new block and the aggregtor to produce rollups
func (m L1Miner) setHead(b *Block) *Block {
	m.aggregator.RPCNewHead(b)
	m.canonicalCh <- b
	return b
}

func (m L1Miner) Stop() {
	m.runCh1 <- false
	m.runCh2 <- false
}

// P2PReceiveBlock is called by counterparties when there is a block to broadcast
// All it does is drop the blocks in a channel for processing.
func (m L1Miner) P2PReceiveBlock(b *Block) {
	//fmt.Printf("%d receive %s\n", m.id, b.id)
	m.p2pCh <- b
}

// startMining - listens on the canonicalCh and schedule a go routine that produces a block after a powTime and drop it on the miningCh channel
func (m L1Miner) startMining() {

	// stores all rollups seen from the beginning of time.
	// store rollups grouped by height, to optimize the algorithm that determines what needs to be included in a block
	// todo - move this out into some state processing
	var mempool = make(map[int][]*L1Tx)

	var deposits = make([]*L1Tx, 0)
	var mut = &sync.RWMutex{}

	var currentHeight = -1

	for {
		select {
		case _ = <-m.runCh2:
			return
		case tx := <-m.mempoolCh:
			mut.Lock()
			if tx.txType == RollupTx {
				r := tx.rollup
				currentHeight = Max(currentHeight, r.height)
				val, found := mempool[r.height]
				if found {
					mempool[r.height] = append(val, tx)
				} else {
					mempool[r.height] = []*L1Tx{tx}
				}
			} else if tx.txType == DepositTx {
				deposits = append(deposits, tx)
			}
			mut.Unlock()

		case cb := <-m.canonicalCh:
			// A new canonical block was found. Start a new round based on that block.

			// Generate a random number, and wait for that number of ms. Equivalent to PoW
			// Include all rollups received during this period.
			nonce := m.cfg.powTime()
			Schedule(nonce, func() {
				var rollups = make([]*L1Tx, 0)
				mut.RLock()

				// only look back 10 rollups - this is an ugly hack for performance
				// todo - move this out
				for i := 0; i < 10; i++ {
					v, f := mempool[currentHeight-i]
					if f {
						rollups = append(rollups, v...)
					}
				}
				mut.RUnlock()
				var toInclude = FindNotIncludedL1Txs(cb, append(deposits, rollups...))
				m.miningCh <- &Block{height: cb.height + 1, rootHash: uuid.New(), nonce: nonce, miner: &m, parent: cb, creationTime: time.Now(), txs: toInclude}
			})
		}
	}
}

// L1P2PGossipTx receive rollups to publish from the linked aggregators
func (m L1Miner) L1P2PGossipTx(tx *L1Tx) {
	m.mempoolCh <- tx
}
