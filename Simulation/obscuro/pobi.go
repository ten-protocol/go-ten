package obscuro

import (
	"fmt"
	"github.com/google/uuid"
	"simulation/common"
	"sync/atomic"
)

// called when a new L1 block was mined
func (a Node) newPobiRound(b common.Block, doneCh *chan bool) {
	// A Pobi round starts when a new canonical L1 block was produced

	// Find the new canonical L2 chain and calculate the State
	// todo - a bit inefficient - cache per rollup?. The state might have been calculated already in the previous round
	blockState := a.updateState(b)

	r := a.produceRollup(b, blockState)
	a.l2Network.BroadcastRollup(encodeRollup(r))

	// wait to receive rollups From peers
	// todo - make this smarter. e.g: if 90% of the peers have sent rollups, proceed. Or if a Nonce is very low and probabilistically there is no chance, etc
	common.ScheduleInterrupt(a.cfg.GossipPeriod, doneCh, func() {
		if atomic.LoadInt32(a.interrupt) == 1 {
			return
		}

		// request the new generation rollups
		a.rollupInCh <- blockState.Head.Height() + 1
		rollupsReceivedFromPeers := <-a.rollupOutCh

		// filter out rollups with a different Parent
		var usefulRollups = []*common.Rollup{&r}
		for _, rol := range rollupsReceivedFromPeers {
			if rol.Parent().Root() == blockState.Head.Root() {
				usefulRollups = append(usefulRollups, rol)
			}
		}

		// determine the winner of the round
		winnerRollup, winnerState := a.findRoundWinner(usefulRollups, &blockState.Head, blockState.State)

		// we are the winner
		if winnerRollup.Agg == a.Id {
			var txsString []string
			for _, t := range winnerRollup.Txs() {
				t1 := t.(common.L2Tx)
				txsString = append(txsString, fmt.Sprintf("%v->%v(%d)", t1.From, t1.Dest, t1.Amount))
			}
			common.Log(fmt.Sprintf(">   Agg%d: (b_%d) create rollup=r_%d(%d)[r_%d]{poof=b_%d}. Txs: %v. State=%v.", a.Id, b.RootHash.ID(), winnerRollup.Root().ID(), winnerRollup.Height(), winnerRollup.Parent().Root().ID(), winnerRollup.Proof().RootHash.ID(), txsString, winnerRollup.State))
			// build a L1 tx with the rollup and send it to the L1 node for further broadcase
			tx := common.L1Tx{Id: uuid.New(), TxType: common.RollupTx, Rollup: winnerRollup}
			t, err := tx.Encode()
			if err != nil {
				panic(err)
			}
			a.L1Node.BroadcastTx(t)
		}

		a.roundWinnerCh <- winner{winnerRollup, winnerState}
	})
}

//
func (a Node) produceRollup(b common.Block, bs BlockState) common.Rollup {
	// retrieve the speculatively calculated State based on the previous winner and the incoming transactions
	a.speculativeWorkInCh <- true
	speculativeRollup := <-a.speculativeWorkOutCh

	newRollupTxs := speculativeRollup.txs
	newRollupState := speculativeRollup.s

	// the speculative execution has been processing on top of the wrong parent - due to failure in gossip
	if speculativeRollup.r.Root() != bs.Head.Root() {
		common.Log(fmt.Sprintf(">   Agg%d: Recalculate. speculative=r_%d(%d), published=r_%d(%d)", a.Id, speculativeRollup.r.Root().ID(), speculativeRollup.r.Height(), bs.Head.Root().ID(), bs.Head.Height()))
		a.statsCollector.L2Recalc(a.Id)

		// determine transactions to include in new rollup and process them
		newRollupTxs = a.currentTxs(bs.Head)
		newRollupState = executeTransactions(newRollupTxs, newProcessedState(bs.State))
	}

	// always process deposits last
	// process deposits from the proof of the parent to the current block (which is the proof of the new rollup)
	proof := bs.Head.Proof()
	newRollupState = processDeposits(&proof, b, copyProcessedState(newRollupState))

	// Create a new rollup based on the proof of inclusion of the previous, including all new transactions
	return common.NewRollup(&b, &bs.Head, a.Id, newRollupTxs, newRollupState.w, serialize(newRollupState.s))
}

// Determine the new canonical L2 head and calculate the State
// Uses cache-ing to map the Head rollup and the State to each L1Node block.
func (a Node) updateState(b common.Block) BlockState {

	// This method is called recursively in case of Re-orgs. Stop when state was calculated already.
	val, found := a.Db.Fetch(b.RootHash)
	if found {
		return val
	}

	// The genesis rollup is part of the canonical chain and will be included in an L1 block by the first Aggregator.
	if b.RootHash == common.GenesisBlock.RootHash {
		bs := BlockState{
			Block: b,
			Head:  common.GenesisRollup,
			State: emptyState(),
		}
		a.Db.Set(b.RootHash, bs)
		return bs
	}

	// To calculate the state after the current block, we need the state after the parent.
	parentState, parentFound := a.Db.Fetch(b.Parent().Root())
	if !parentFound {
		// go back and calculate the State of the Parent
		parentState = a.updateState(*b.ParentBlock())
	}

	bs := calculateBlockState(b, parentState)

	a.Db.Set(b.RootHash, bs)

	return bs
}

// given an L1 block, and the State as it was in the Parent block, calculates the State after the current block.
func calculateBlockState(b common.Block, parentState BlockState) BlockState {
	var newHead *common.Rollup = nil

	// find the Head rollup according to the rules
	for _, t := range b.L1Txs() {
		// go through all rollup transactions
		if t.TxType == common.RollupTx {
			r := t.Rollup
			// only consider rollups if they advance the chain
			if (r.Height() > parentState.Head.Height()) && common.IsAncestor(r.Proof(), b) {
				if newHead == nil || r.Height() > newHead.Height() || r.Proof().Height() > newHead.Proof().Height() || (r.Proof().Height() == newHead.Proof().Height() && r.Nonce < newHead.Nonce) {
					newHead = &r
				}
			}
		}
	}

	s := newProcessedState(parentState.State)

	var found bool
	// only change the state if there is a new l2 Head in the current block
	if newHead != nil {
		s = executeTransactions(newHead.L2Txs(), s)
		p := newHead.ParentRollup().Proof()
		s = processDeposits(&p, newHead.Proof(), s)
		found = true
	} else {
		newHead = &parentState.Head
		found = false
	}

	bs := BlockState{
		Block:          b,
		Head:           *newHead,
		State:          s.s,
		foundNewRollup: found,
	}
	return bs
}

// mutates the state
// process deposits from the proof of the parent rollup(exclusive) to the proof of the current rollup
func processDeposits(fromBlock *common.Block, toBlock common.Block, s ProcessedState) ProcessedState {
	from := common.GenesisBlock.RootHash
	height := common.GenesisHeight
	if fromBlock != nil {
		from = fromBlock.RootHash
		height = fromBlock.Height()
		if !common.IsAncestor(fromBlock, toBlock) {
			panic("wtf")
		}

	}

	b := &toBlock
	for {
		if b.RootHash == from {
			break
		}
		for _, tx := range b.L1Txs() {
			// transactions to a hardcoded bridge address
			if tx.TxType == common.DepositTx {
				v, f := s.s[tx.Dest]
				if f {
					s.s[tx.Dest] = v + tx.Amount
				} else {
					s.s[tx.Dest] = tx.Amount
				}
			}
		}
		if b.Height() < height {
			panic("something went wrong")
		}
		b = b.ParentBlock()
	}
	return s
}

func (a Node) findRoundWinner(receivedRollups []*common.Rollup, parent *common.Rollup, parentState State) (common.Rollup, State) {
	var win *common.Rollup
	for _, r := range receivedRollups {
		if r.Parent().Root() != parent.Root() {
			continue
		}
		if win == nil || r.Proof().Height() > win.Proof().Height() || (r.Proof().Height() == win.Proof().Height() && r.Nonce < win.Nonce) {
			win = r
		}
	}

	// calculate the state to compare with what is in the Rollup
	s := newProcessedState(parentState)
	p := win.ParentRollup().Proof()
	s = processDeposits(&p, win.Proof(), s)
	s = executeTransactions(win.L2Txs(), s)
	// todo - check that s is valid against the State in the rollup, if not - call the function again with this tx excluded
	return *win, s.s
}

// Calculate transactions to be included in the current rollup
func (a Node) currentTxs(head common.Rollup) []common.L2Tx {
	// Requests all l2 transactions received over gossip
	a.txsInCh <- true
	txs := <-a.txsOutCh

	resultTxs := common.FindNotIncludedTxs(head, txs)
	txsCopy := make([]common.L2Tx, len(resultTxs))
	for i, tx := range resultTxs {
		txsCopy[i] = tx.(common.L2Tx)
	}

	// and return only the ones not included in any rollup so far
	return txsCopy
}
