package ethereum_mock

import (
	"fmt"
	"simulation/common"
)

type L1MiningConfig struct {
	PowTime common.Latency
}

type L1Node struct {
	Id      common.NodeId
	cfg     L1MiningConfig
	Agg     common.NotifyNewBlock
	network common.L1Network
	mining  bool

	// Channels
	runCh1      chan bool          // add false when the Miner has to stop - todo - rename to done
	runCh2      chan bool          // the mining loop is notified to stop
	p2pCh       chan *common.Block // this is where blocks received from peers are dropped
	miningCh    chan *common.Block // this is where blocks created by the mining setup of the current node are dropped
	canonicalCh chan *common.Block // this is where the main processing routine drops blocks that are canonical
	mempoolCh   chan *common.L1Tx  // where l1 transactions to be published in the next block are added
}

func NewMiner(id common.NodeId, cfg L1MiningConfig, agg common.NotifyNewBlock, network common.L1Network) L1Node {
	return L1Node{Id: id, mining: true, cfg: cfg, Agg: agg, network: network, runCh1: make(chan bool), p2pCh: make(chan *common.Block), miningCh: make(chan *common.Block), canonicalCh: make(chan *common.Block), runCh2: make(chan bool), mempoolCh: make(chan *common.L1Tx)}
}

// Start runs an infinite loop that listens to the two block producing channels and processes them.
// it outputs the winning blocks to the roundWinnerCh channel
func (m L1Node) Start() {
	// This starts the mining
	go m.startMining()

	var head = m.setHead(&common.GenesisBlock)

	for {
		select {
		case p2pb := <-m.p2pCh: // Received from peers
			if p2pb.Height() > head.Height() {
				if !common.IsAncestor(*head, *p2pb) {
					// todo - find a way
					//simulation.statsMu.Lock()
					//m.network.Stats.noL1Reorgs[m.Id]++
					//simulation.statsMu.Unlock()
					fork := common.LCA(head, p2pb)
					common.Log(fmt.Sprintf("> M%d: L1Reorg new=b_%d(%d), old=b_%d(%d), fork=b_%d(%d)", m.Id, p2pb.RootHash().ID(), p2pb.Height, head.RootHash().ID(), head.Height, fork.RootHash().ID(), fork.Height))
				}
				head = m.setHead(p2pb)
			}
		case mb := <-m.miningCh: // Received from the local mining
			if mb.Height() > head.Height() { // Ignore the locally produced block if someone else found one already
				common.Log(m.printBlock(mb))
				head = m.setHead(mb)
				m.network.BroadcastBlockL1(*mb)
			}
		case _ = <-m.runCh1:
			return
		}
	}
}

func (m L1Node) printBlock(mb *common.Block) string {
	// This is just for printing
	var txs []string
	//for _, tx := range mb.Txs {
	//if tx.TxType == RollupTx {
	//	txs = append(txs, fmt.Sprintf("r_%d", tx.Rollup.RootHash.ID()))
	//} else {
	//	txs = append(txs, fmt.Sprintf("deposit(%v=%d)", tx.Dest, tx.Amount))
	//}
	//}
	return fmt.Sprintf("> M%d: create b_%d(Height=%d, Nonce=%d)[p=b_%d]. Txs: %v", m.Id, mb.RootHash().ID(), mb.Height, mb.Nonce, mb.Parent().RootHash().ID(), txs)
}

// Notifies the Miner to start mining on the new block and the aggregtor to produce rollups
func (m L1Node) setHead(b *common.Block) *common.Block {
	m.Agg.RPCNewHead(*b)
	m.canonicalCh <- b
	return b
}

func (m L1Node) Stop() {
	m.runCh1 <- false
	m.runCh2 <- false
}

// P2PReceiveBlock is called by counterparties when there is a block to broadcast
// All it does is drop the blocks in a channel for processing.
func (m L1Node) P2PReceiveBlock(b *common.Block) {
	//fmt.Printf("%d receive %s\n", m.Id, b.Id)
	m.p2pCh <- b
}

// startMining - listens on the canonicalCh and schedule a go routine that produces a block after a PowTime and drop it on the miningCh channel
func (m L1Node) startMining() {

	// stores all rollups seen from the beginning of time.
	// store rollups grouped by Height, to optimize the algorithm that determines what needs to be included in a block
	// todo - move this out into some state processing
	//var mempool = make(map[int][]*L1Tx)
	//var deposits = make([]*L1Tx, 0)
	//var mut = &sync.RWMutex{}
	var txs = make([]common.Tx, 0)

	//var currentHeight = -1

	var doneCh *chan bool = nil

	for {
		select {
		case _ = <-m.runCh2:
			return
		case tx := <-m.mempoolCh:
			txs = append(txs, *tx)
			//mut.Lock()
			//if tx.TxType == RollupTx {
			//	r := tx.Rollup
			//	currentHeight = common.Max(currentHeight, r.Height())
			//	val, found := mempool[r.Height()]
			//	if found {
			//		mempool[r.Height()] = append(val, tx)
			//	} else {
			//		mempool[r.Height()] = []*L1Tx{tx}
			//	}
			//} else if tx.TxType == DepositTx {
			//	deposits = append(deposits, tx)
			//}
			//mut.Unlock()

		case cb := <-m.canonicalCh:
			// A new canonical block was found. Start a new round based on that block.

			//notify the existing mining thread to stop mining
			if doneCh != nil {
				*doneCh <- true
			}

			c := make(chan bool)
			doneCh = &c
			// Generate a random number, and wait for that number of ms. Equivalent to PoW
			// Include all rollups received during this period.
			nonce := m.cfg.PowTime()
			common.ScheduleInterrupt(nonce, doneCh, func() {
				//var rollups = make([]*L1Tx, 0)
				//mut.RLock()
				//
				//// only look back 10 rollups - this is an ugly hack for performance
				//// todo - move this out
				//for i := 0; i < 10; i++ {
				//	v, f := mempool[currentHeight-i]
				//	if f {
				//		rollups = append(rollups, v...)
				//	}
				//}
				//mut.RUnlock()
				//all := make([]*L1Tx, 0)
				//all = append(all, rollups...)
				//all = append(all, deposits...)
				toInclude := common.FindNotIncludedTxs(cb, txs)
				txsCopy := make([]common.L1Tx, len(toInclude))
				for i, tx := range toInclude {
					txsCopy[i] = tx.(common.L1Tx)
				}
				m.miningCh <- common.NewBlock(cb, nonce, m.Id, txsCopy)
			})
		}
	}
}

// L1P2PGossipTx receive rollups to publish from the linked aggregators
func (m L1Node) L1P2PGossipTx(tx *common.L1Tx) {
	m.mempoolCh <- tx
}
