package obscuro

import (
	"fmt"
	"github.com/google/uuid"
	"simulation/common"
	"simulation/ethereum-mock"
)

type L2Cfg struct {
	GossipPeriod int
}

// this will become the Obscuro "Node" type
type L2Agg struct {
	Id        common.NodeId
	L1        *ethereum_mock.L1Miner
	cfg       L2Cfg
	l1Network common.L1Network
	l2Network common.L2Network

	// control the lifecycle
	runCh1 chan bool
	runCh2 chan bool

	// where rollups and transactions are gossipped From peers
	p2pChRollup chan *common.Rollup
	p2pChTx     chan common.L2Tx

	// where the connected L1 node drops new blocks
	rpcCh chan *common.Block

	// used for internal communication between the gossi agent and the processing agent
	// todo - probably can use a single channel
	rollupIntCh chan int
	rollupOutCh chan []*common.Rollup

	// used for internal communication between the gossi agent and the processing agent
	txsIntCh chan bool
	txsOutCh chan []common.L2Tx

	// todo
	progIntCh chan bool
	progOutCh chan currentWork

	// when a round finishes and a winner is discovered. Notifies the gossip actor to start processing new transactions.
	roundWinnerCh chan BlockState

	// a database of work already executed
	Db Db
}

func NewAgg(id common.NodeId, cfg L2Cfg, l1 *ethereum_mock.L1Miner, l1Network common.L1Network, l2Network common.L2Network) L2Agg {
	return L2Agg{
		Id:            id,
		cfg:           cfg,
		l1Network:     l1Network,
		l2Network:     l2Network,
		L1:            l1,
		runCh1:        make(chan bool),
		runCh2:        make(chan bool),
		p2pChRollup:   make(chan *common.Rollup),
		p2pChTx:       make(chan common.L2Tx),
		rpcCh:         make(chan *common.Block),
		rollupIntCh:   make(chan int),
		rollupOutCh:   make(chan []*common.Rollup),
		txsIntCh:      make(chan bool),
		txsOutCh:      make(chan []common.L2Tx),
		roundWinnerCh: make(chan BlockState),
		progIntCh:     make(chan bool),
		progOutCh:     make(chan currentWork),
		Db:            NewInMemoryDb(),
	}
}

func (a L2Agg) Start() {
	go a.startGossip()
	var doneCh *chan bool = nil

	for {
		select {
		// Main loop
		// Listen for notifications From the L1 node and process them
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
	r   *common.Rollup
	s   State
	txs []common.L2Tx
}

// actor that participates in rollup and transaction gossip
// processes transactions
func (a L2Agg) startGossip() {

	// Rollups grouped by Height
	var allRollups = make(map[int][]*common.Rollup)

	// transactions
	var allTxs = make([]common.L2Tx, 0)

	// Process transactions on the fly
	var currentHead = &common.GenesisRollup
	var currentState = emptyState()
	var currentProcessedTxs = make([]common.L2Tx, 0)

	for {
		select {
		case r := <-a.roundWinnerCh:
			currentHead = r.Head
			currentState = r.State
			currentProcessedTxs = make([]common.L2Tx, 0)

			// determine the transactions that were not included in the previous Head
			// this is terribly inefficient
			copyTxs := make([]common.Tx, len(allTxs))
			for i, tx := range allTxs {
				copyTxs[i] = common.Tx(tx)
			}

			txs := common.FindNotIncludedTxs(currentHead, copyTxs)

			copyResult := make([]common.L2Tx, len(txs))
			for i, tx := range txs {
				copyResult[i] = tx.(common.L2Tx)
			}

			// calculate the State after executing them
			currentState = calculateState(copyResult, currentState)

		case tx := <-a.p2pChTx:
			allTxs = append(allTxs, tx)
			currentProcessedTxs = append(currentProcessedTxs, tx)
			executeTx(currentState, tx)

		case _ = <-a.progIntCh:
			b := make([]common.L2Tx, len(currentProcessedTxs))
			copy(b, currentProcessedTxs)
			a.progOutCh <- currentWork{
				r:   currentHead,
				s:   copyState(currentState),
				txs: b,
			}

		case r := <-a.p2pChRollup:
			val, found := allRollups[r.Height()]
			if found {
				allRollups[r.Height()] = append(val, r)
			} else {
				allRollups[r.Height()] = []*common.Rollup{r}
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

// RPCNewHead Receive notifications From the L1 Node when there's a new block
func (a L2Agg) RPCNewHead(b common.Block) {
	a.rpcCh <- &b
}

// L2P2PGossipRollup is called by counterparties when there is a Rollup to broadcast
// All it does is drop the Rollups in a channel for processing.
func (a L2Agg) L2P2PGossipRollup(r *common.Rollup) {
	a.p2pChRollup <- r
}

func (a L2Agg) L2P2PReceiveTx(tx common.L2Tx) {
	a.p2pChTx <- tx
}

// main L1 block processing logic - the POBI protocol
func (a L2Agg) processBlock(b *common.Block, doneCh *chan bool) {
	// A Pobi round starts when a new canonical L1 block was produced

	// Find the new canonical L2 chain and calculate the State
	bs := a.calculateL2State(b)

	// retrieve the calculated State based on the previous winner and the incoming transactions
	//a.progIntCh <- true
	//current := <-a.progOutCh

	// the transactions and the State as calculated during the round
	//Txs := current.Txs
	//stateAfter := current.s
	//newL2Head := current.r

	/*** this here toremove */
	//avoid using the precalculated stuff
	newL2Head := bs.Head
	stateAfter := processDeposits(b, copyState(bs.State))
	txs := a.currentTxs(newL2Head)
	stateAfter = calculateState(txs, stateAfter)

	/***/

	// the transactions were processed on a wrong rollup
	// we were working on the wrong winner
	/*	if newL2Head.RootHash != bs.Head.RootHash {
			if !IsRlpAncestor(newL2Head, bs.Head) && !IsRlpAncestor(bs.Head, newL2Head) {
				log(fmt.Sprintf(">   Agg%d: Reorg. published=r_%d(%d), existing=r_%d(%d)", a.Id, newL2Head.RootHash.ID(), newL2Head.Height, bs.Head.RootHash.ID(), bs.Head.Height))
				statsMu.Lock()
				a.network.Stats.noL2Reorgs[a.Id]++
				statsMu.Unlock()
			}

			newL2Head = bs.Head
			stateAfter = processDeposits(b, bs.State)

			// determine the transactions to be included on the actual winner
			// and calculate the State after executing them
			Txs = a.currentTxs(newL2Head)
			stateAfter = calculateState(Txs, stateAfter)
		}
	*/
	// Create a new rollup based on the proof of inclusion of the previous, including all new transactions
	r := common.NewRollup(b, newL2Head, a.Id, txs, serialize(stateAfter))
	a.l2Network.BroadcastRollupL2(r)

	// wait to receive rollups From peers
	// todo - make this smarter. e.g: if 90% of the peers have sent rollups, proceed. Or if a Nonce is very low and probabilistically there is no chance, etc
	common.ScheduleInterrupt(a.cfg.GossipPeriod, doneCh, func() {

		// request the new generation rolllups
		a.rollupIntCh <- newL2Head.Height() + 1
		rollups := <-a.rollupOutCh

		// filter out rollups with a different Parent
		var usefulRollups = []*common.Rollup{&r}
		for _, rol := range rollups {
			if rol.Parent().RootHash() == newL2Head.RootHash() {
				usefulRollups = append(usefulRollups, rol)
			}
		}

		// determine the winner of the round
		winner := a.findRoundWinner(usefulRollups, newL2Head, bs.State, b)

		// we are the winner
		if winner.Head.Agg == a.Id {
			var txsString []string
			for _, t := range winner.Head.Txs() {
				t1 := t.(common.L2Tx)
				txsString = append(txsString, fmt.Sprintf("%v->%v(%d)", t1.From, t1.Dest, t1.Amount))
			}
			sum := 0
			for _, bal := range winner.State {
				sum += bal
			}
			common.Log(fmt.Sprintf(">   Agg%d: (b_%d)create rollup=r_%d(%d)[r_%d]{poof=b_%d}. Txs: %v. State=%v. Total=%d", a.Id, b.RootHash().ID(), winner.Head.RootHash().ID(), winner.Head.Height, winner.Head.Parent().RootHash().ID(), winner.Head.L1Proof.RootHash().ID(), txsString, winner.State, sum))
			// build a L1 tx with the rollup
			a.l1Network.BroadcastL1Tx(common.L1Tx{Id: uuid.New(), TxType: common.RollupTx, Rollup: *winner.Head})
		}

		a.roundWinnerCh <- winner
	})
}

// Calculate transactions to be included in the current rollup
func (a L2Agg) currentTxs(head *common.Rollup) []common.L2Tx {
	// Requests all l2 transactions received over gossip
	a.txsIntCh <- true
	txs := <-a.txsOutCh
	copyTxs := make([]common.Tx, len(txs))
	for i, tx := range txs {
		copyTxs[i] = common.Tx(tx)
	}

	resultTxs := common.FindNotIncludedTxs(head, copyTxs)
	txsCopy := make([]common.L2Tx, len(resultTxs))
	for i, tx := range resultTxs {
		txsCopy[i] = tx.(common.L2Tx)
	}

	// and return only the ones not included in any rollup so far
	return txsCopy
}

// Complex logic to determine the new canonical Head and to calculate the State
// Uses cache-ing to map the Head rollup and the State to each L1 block in case of rollbacks.
func (a L2Agg) calculateL2State(b *common.Block) BlockState {
	val, found := a.Db.Fetch(b.RootHash())
	if found {
		return val
	}

	//  1. The genesis rollup is part of the canonical chain and will be included in an L1 block by the first Aggregator.
	if b.RootHash() == common.GenesisBlock.RootHash() {
		bs := BlockState{
			Head:  &common.GenesisRollup,
			State: emptyState(),
		}
		a.Db.Set(b.RootHash(), bs)
		return bs
	}

	parentState, parentFound := a.Db.Fetch(b.Parent().RootHash())
	if !parentFound {
		// go back and calculate the State of the Parent
		parentState = a.calculateL2State(b.ParentBlock())
	}

	bs := calculateBlockState(b, parentState)

	a.Db.Set(b.RootHash(), bs)

	return bs
}

// given a new L1 block, and the State as it was in the Parent block, calculates the State after the current block
func calculateBlockState(b *common.Block, parentState BlockState) BlockState {
	var newHead *common.Rollup = nil

	s := copyState(parentState.State)

	// always process deposits of the Parent block
	s = processDeposits(b.ParentBlock(), s)

	// find the Head rollup according to the rules
	for _, t := range b.L1Txs() {
		// go through all rollup transactions
		if t.TxType == common.RollupTx {
			r := t.Rollup
			// only consider rollups if they advance the chain
			if r.Height() > parentState.Head.Height() {
				if newHead == nil || r.Height() > newHead.Height() || r.L1Proof.Height() > newHead.L1Proof.Height() || (r.L1Proof.Height() == newHead.L1Proof.Height() && r.Nonce < newHead.Nonce) {
					newHead = &r
				}
			}
		}
	}

	// only apply transactions if there is a new l2 Head
	if newHead != nil {
		s = calculateState(newHead.L2Txs(), s)
	} else {
		newHead = parentState.Head
	}

	bs := BlockState{
		Head:  newHead,
		State: s,
	}
	return bs
}

func processDeposits(b *common.Block, s State) State {
	if b == nil {
		return emptyState()
	}
	for _, tx := range b.L1Txs() {
		if tx.TxType == common.DepositTx {
			v, f := s[tx.Dest]
			if f {
				s[tx.Dest] = v + tx.Amount
			} else {
				s[tx.Dest] = tx.Amount
			}
		}
	}
	return s
}

func (a L2Agg) findRoundWinner(receivedRollups []*common.Rollup, head *common.Rollup, state State, b *common.Block) BlockState {
	var win *common.Rollup
	for _, r := range receivedRollups {
		if r.Parent().RootHash() != head.RootHash() {
			continue
		}
		if win == nil || r.L1Proof.Height() > win.L1Proof.Height() || (r.L1Proof.Height() == win.L1Proof.Height() && r.Nonce < win.Nonce) {
			win = r
		}
	}

	s := copyState(state)
	s = processDeposits(b, s)
	s = calculateState(win.L2Txs(), s)
	// todo - check that s is valid against the State in the rollup, if not - call the function again with this tx excluded
	return BlockState{
		Head:  win,
		State: s,
	}
}

func (a L2Agg) Stop() {
	a.runCh1 <- false
	a.runCh2 <- false
}
