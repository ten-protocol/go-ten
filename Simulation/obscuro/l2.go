package obscuro

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

type L2Cfg struct {
	gossipPeriodMs int
}

type L2Agg struct {
	id      NodeId
	l1      *L1Miner
	cfg     L2Cfg
	network *NetworkCfg

	// control the lifecycle
	runCh1 chan bool
	runCh2 chan bool

	// where rollups and transactions are gossipped from peers
	p2pChRollup chan *Rollup
	p2pChTx     chan *L2Tx

	// where the connected L1 node drops new blocks
	rpcCh chan *Block

	// used for internal communication between the gossi agent and the processing agent
	// todo - probably can use a single channel
	rollupIntCh chan int
	rollupOutCh chan []*Rollup

	// used for internal communication between the gossi agent and the processing agent
	txsIntCh chan bool
	txsOutCh chan []*L2Tx

	// todo
	progIntCh chan bool
	progOutCh chan currentWork

	// when a round finishes and a winner is discovered. Notifies the gossip actor to start processing new transactions.
	roundWinnerCh chan BlockState

	// a database of work already executed
	db Db
}

func NewAgg(id NodeId, cfg L2Cfg, l1 *L1Miner, network *NetworkCfg) L2Agg {
	return L2Agg{
		id:            id,
		cfg:           cfg,
		network:       network,
		l1:            l1,
		runCh1:        make(chan bool),
		runCh2:        make(chan bool),
		p2pChRollup:   make(chan *Rollup),
		p2pChTx:       make(chan *L2Tx),
		rpcCh:         make(chan *Block),
		rollupIntCh:   make(chan int),
		rollupOutCh:   make(chan []*Rollup),
		txsIntCh:      make(chan bool),
		txsOutCh:      make(chan []*L2Tx),
		roundWinnerCh: make(chan BlockState),
		progIntCh:     make(chan bool),
		progOutCh:     make(chan currentWork),
		db:            NewInMemoryDb(),
	}
}

type L2RootHash = uuid.UUID

// Todo - this has to be a trie root eventually
type StateRoot = string

type Rollup struct {
	height       int
	rootHash     L2RootHash
	agg          *L2Agg
	parent       *Rollup
	creationTime time.Time
	l1Proof      *Block // the L1 block where the parent was published
	nonce        Nonce
	txs          []*L2Tx
	state        StateRoot
}

// Transfers and Withdrawals for now
type L2TxType int64

const (
	TransferTx L2TxType = iota
	WithdrawalTx
)

type L2TxId = uuid.UUID

// no signing for now
type L2Tx struct {
	id     L2TxId // this is the hash
	txType L2TxType
	amount int
	from   Address
	dest   Address
}

var GenesisRollup = Rollup{-1,
	uuid.New(),
	nil,
	nil,
	time.Now(),
	nil,
	0,
	[]*L2Tx{},
	serialize(emptyState()),
}

func (a L2Agg) Start() {
	go a.startGossip()
	var doneCh *chan bool = nil

	for {
		select {
		// Main loop
		// Listen for notifications from the L1 node and process them
		// Note that during processing, more recent notifications can be received
		case b := <-a.rpcCh:
			if doneCh != nil {
				*doneCh <- true
			}

			c := make(chan bool)
			doneCh = &c

			go a.processBlock(b, doneCh)
		case _ = <-a.runCh1:
			return
		}
	}
}

type currentWork struct {
	r   *Rollup
	s   State
	txs []*L2Tx
}

// actor that participates in rollup and transaction gossip
// processes transactions
func (a L2Agg) startGossip() {

	// Rollups grouped by height
	var allRollups = make(map[int][]*Rollup)

	// transactions
	var allTxs = make([]*L2Tx, 0)

	// Process transactions on the fly
	var currentHead = &GenesisRollup
	var currentState = emptyState()
	var currentProcessedTxs = make([]*L2Tx, 0)

	for {
		select {
		case r := <-a.roundWinnerCh:
			currentHead = r.head
			currentState = r.state
			currentProcessedTxs = make([]*L2Tx, 0)

			// determine the transactions that were not included in the previous head
			// this is terribly inefficient
			txs := FindNotIncludedTxs(currentHead, allTxs)

			// calculate the state after executing them
			currentState = calculateState(txs, currentState)

		case tx := <-a.p2pChTx:
			allTxs = append(allTxs, tx)
			currentProcessedTxs = append(currentProcessedTxs, tx)
			executeTx(currentState, tx)

		case _ = <-a.progIntCh:
			b := make([]*L2Tx, len(currentProcessedTxs))
			copy(b, currentProcessedTxs)
			a.progOutCh <- currentWork{
				r:   currentHead,
				s:   copyState(currentState),
				txs: b,
			}

		case r := <-a.p2pChRollup:
			val, found := allRollups[r.height]
			if found {
				allRollups[r.height] = append(val, r)
			} else {
				allRollups[r.height] = []*Rollup{r}
			}

		case requestedHeight := <-a.rollupIntCh:
			a.rollupOutCh <- allRollups[requestedHeight]
		case _ = <-a.txsIntCh:
			a.txsOutCh <- allTxs

		case _ = <-a.runCh2:
			return
		}
	}
}

// RPCNewHead Receive notifications from the L1 Node when there's a new block
func (a L2Agg) RPCNewHead(b *Block) {
	a.rpcCh <- b
}

// L2P2PGossipRollup is called by counterparties when there is a Rollup to broadcast
// All it does is drop the Rollups in a channel for processing.
func (a L2Agg) L2P2PGossipRollup(r *Rollup) {
	a.p2pChRollup <- r
}

func (a L2Agg) L2P2PReceiveTx(tx *L2Tx) {
	a.p2pChTx <- tx
}

// main L1 block processing logic - the POBI protocol
func (a L2Agg) processBlock(b *Block, doneCh *chan bool) {
	// A Pobi round starts when a new canonical L1 block was produced

	// Find the new canonical L2 chain and calculate the state
	bs := a.calculateL2State(b)

	// retrieve the calculated state based on the previous winner and the incoming transactions
	//a.progIntCh <- true
	//current := <-a.progOutCh

	// the transactions and the state as calculated during the round
	//txs := current.txs
	//stateAfter := current.s
	//newL2Head := current.r

	/*** this here toremove */
	//avoid using the precalculated stuff
	newL2Head := bs.head
	stateAfter := processDeposits(b, copyState(bs.state))
	txs := a.currentTxs(newL2Head)
	stateAfter = calculateState(txs, stateAfter)

	/***/

	// the transactions were processed on a wrong rollup
	// we were working on the wrong winner
	/*	if newL2Head.rootHash != bs.head.rootHash {
			if !IsRlpAncestor(newL2Head, bs.head) && !IsRlpAncestor(bs.head, newL2Head) {
				log(fmt.Sprintf(">   Agg%d: Reorg. published=r_%d(%d), existing=r_%d(%d)", a.id, newL2Head.rootHash.ID(), newL2Head.height, bs.head.rootHash.ID(), bs.head.height))
				statsMu.Lock()
				a.network.Stats.noL2Reorgs[a.id]++
				statsMu.Unlock()
			}

			newL2Head = bs.head
			stateAfter = processDeposits(b, bs.state)

			// determine the transactions to be included on the actual winner
			// and calculate the state after executing them
			txs = a.currentTxs(newL2Head)
			stateAfter = calculateState(txs, stateAfter)
		}
	*/
	// Create a new rollup based on the proof of inclusion of the previous, including all new transactions
	r := Rollup{newL2Head.height + 1, uuid.New(), &a, newL2Head, time.Now(), b, generateNonce(), txs, serialize(stateAfter)}
	a.network.broadcastRollupL2(&r)

	// wait to receive rollups from peers
	// todo - make this smarter. e.g: if 90% of the peers have sent rollups, proceed. Or if a nonce is very low and probabilistically there is no chance, etc
	ScheduleInterrupt(a.cfg.gossipPeriodMs, doneCh, func() {

		// request the new generation rolllups
		a.rollupIntCh <- newL2Head.height + 1
		rollups := <-a.rollupOutCh

		// filter out rollups with a different parent
		var usefulRollups = []*Rollup{&r}
		for _, rol := range rollups {
			if rol.parent.rootHash == newL2Head.rootHash {
				usefulRollups = append(usefulRollups, rol)
			}
		}

		// determine the winner of the round
		winner := a.findRoundWinner(usefulRollups, newL2Head, bs.state, b)

		// we are the winner
		if winner.head.agg.id == a.id {
			var txsString []string
			for _, t := range winner.head.txs {
				txsString = append(txsString, fmt.Sprintf("%v->%v(%d)", t.from, t.dest, t.amount))
			}
			sum := 0
			for _, bal := range winner.state {
				sum += bal
			}
			log(fmt.Sprintf(">   Agg%d: (b_%d)create rollup=r_%d(%d)[r_%d]{poof=b_%d}. Txs: %v. State=%v. Total=%d", a.id, b.rootHash.ID(), winner.head.rootHash.ID(), winner.head.height, winner.head.parent.rootHash.ID(), winner.head.l1Proof.rootHash.ID(), txsString, winner.state, sum))
			// build a L1 tx with the rollup
			a.network.broadcastL1Tx(&L1Tx{id: uuid.New(), txType: RollupTx, rollup: winner.head})
		}

		a.roundWinnerCh <- winner
	})
}

// Calculate transactions to be included in the current rollup
func (a L2Agg) currentTxs(head *Rollup) []*L2Tx {
	// Requests all l2 transactions received over gossip
	a.txsIntCh <- true
	txs := <-a.txsOutCh
	// and return only the ones not included in any rollup so far
	return FindNotIncludedTxs(head, txs)
}

// Complex logic to determine the new canonical head and to calculate the state
// Uses cache-ing to map the head rollup and the state to each l1 block in case of rollbacks.
func (a L2Agg) calculateL2State(b *Block) BlockState {
	val, found := a.db.fetch(b.rootHash)
	if found {
		return val
	}

	//  1. The genesis rollup is part of the canonical chain and will be included in an L1 block by the first Aggregator.
	if b.rootHash == GenesisBlock.rootHash {
		bs := BlockState{
			head:  &GenesisRollup,
			state: emptyState(),
		}
		a.db.set(b.rootHash, bs)
		return bs
	}

	parentState, parentFound := a.db.fetch(b.parent.rootHash)
	if !parentFound {
		// go back and calculate the state of the parent
		parentState = a.calculateL2State(b.parent)
	}

	bs := calculateBlockState(b, parentState)

	a.db.set(b.rootHash, bs)

	return bs
}

// given a new L1 block, and the state as it was in the parent block, calculates the state after the current block
func calculateBlockState(b *Block, parentState BlockState) BlockState {
	var newHead *Rollup = nil

	s := copyState(parentState.state)

	// always process deposits of the parent block
	s = processDeposits(b.parent, s)

	// find the head rollup according to the rules
	for _, t := range b.txs {
		// go through all rollup transactions
		if t.txType == RollupTx {
			r := t.rollup
			// only consider rollups if they advance the chain
			if r.height > parentState.head.height {
				if newHead == nil || r.height > newHead.height || r.l1Proof.height > newHead.l1Proof.height || (r.l1Proof.height == newHead.l1Proof.height && r.nonce < newHead.nonce) {
					newHead = r
				}
			}
		}
	}

	// only apply transactions if there is a new l2 head
	if newHead != nil {
		s = calculateState(newHead.txs, s)
	} else {
		newHead = parentState.head
	}

	bs := BlockState{
		head:  newHead,
		state: s,
	}
	return bs
}

func processDeposits(b *Block, s State) State {
	if b == nil {
		return emptyState()
	}
	for _, tx := range b.txs {
		if tx.txType == DepositTx {
			v, f := s[tx.dest]
			if f {
				s[tx.dest] = v + tx.amount
			} else {
				s[tx.dest] = tx.amount
			}
		}
	}
	return s
}

func (a L2Agg) findRoundWinner(receivedRollups []*Rollup, head *Rollup, state State, b *Block) BlockState {
	var win *Rollup
	for _, r := range receivedRollups {
		if r.parent.rootHash != head.rootHash {
			continue
		}
		if win == nil || r.l1Proof.height > win.l1Proof.height || (r.l1Proof.height == win.l1Proof.height && r.nonce < win.nonce) {
			win = r
		}
	}

	s := copyState(state)
	s = processDeposits(b, s)
	s = calculateState(win.txs, s)
	// todo - check that s is valid against the state in the rollup, if not - call the function again with this tx excluded
	return BlockState{
		head:  win,
		state: s,
	}
}

func (a L2Agg) Stop() {
	a.runCh1 <- false
	a.runCh2 <- false
}
