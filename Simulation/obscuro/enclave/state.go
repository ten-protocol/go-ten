package enclave

import (
	"fmt"
	"simulation/common"
	common2 "simulation/obscuro/common"
)

type State = map[common.Address]uint64

// BlockState - Represents the state after an L1 block was processed.
type BlockState struct {
	Block          *common.Block
	Head           *common2.Rollup
	State          State
	foundNewRollup bool
}

// RollupState - state after an L2 rollups was processed
type RollupState struct {
	s State
	w []common2.Withdrawal
}

func newProcessedState(s State) RollupState {
	return RollupState{
		s: copyState(s),
		w: []common2.Withdrawal{},
	}
}

func copyState(state State) State {
	s := make(State)
	for address, balance := range state {
		s[address] = balance
	}
	return s
}

func copyProcessedState(s RollupState) RollupState {
	return RollupState{
		s: copyState(s.s),
		w: s.w,
	}
}

func serialize(state State) string {
	return fmt.Sprintf("%v", state)
}

// returns a modified copy of the State
func executeTransactions(txs []common2.L2Tx, state RollupState) RollupState {
	ps := copyProcessedState(state)
	for _, tx := range txs {
		executeTx(&ps, tx)
	}
	//fmt.Printf("w1: %v\n", is.w)
	return ps
}

// mutates the State
func executeTx(s *RollupState, tx common2.L2Tx) {
	switch tx.TxType {
	case common2.TransferTx:
		executeTransfer(s, tx)
	case common2.WithdrawalTx:
		executeWithdrawal(s, tx)
	default:
		panic("Invalid transaction type")
	}
}

func executeWithdrawal(s *RollupState, tx common2.L2Tx) {
	bal := s.s[tx.From]
	if bal >= tx.Amount {
		s.s[tx.From] -= tx.Amount
		s.w = append(s.w, common2.Withdrawal{
			Amount:  tx.Amount,
			Address: tx.From,
		})
		//fmt.Printf("w: %v\n", s.w)
	}
}

func executeTransfer(s *RollupState, tx common2.L2Tx) {
	bal := s.s[tx.From]
	if bal >= tx.Amount {
		s.s[tx.From] -= tx.Amount
		s.s[tx.To] += tx.Amount
	}
}

func emptyState() State {
	return make(State)
}

// Determine the new canonical L2 head and calculate the State
// Uses cache-ing to map the Head rollup and the State to each L1Node block.
func updateState(b *common.Block, db Db) BlockState {

	// This method is called recursively in case of Re-orgs. Stop when state was calculated already.
	val, found := db.FetchState(b.Hash())
	if found {
		return val
	}

	// The genesis rollup is part of the canonical chain and will be included in an L1 block by the first Aggregator.
	if b.Hash() == common.GenesisBlock.Hash() {
		bs := BlockState{
			Block:          b,
			Head:           &common2.GenesisRollup,
			State:          emptyState(),
			foundNewRollup: true,
		}
		db.SetState(b.Hash(), bs)
		return bs
	}

	// To calculate the state after the current block, we need the state after the parent.
	parentState, parentFound := db.FetchState(b.Header.ParentHash)
	if !parentFound {
		// go back and calculate the State of the Parent
		p, f := b.Parent(db)
		if !f {
			panic("wtf")
		}
		parentState = updateState(p, db)
	}

	bs := calculateBlockState(b, parentState, db)

	db.SetState(b.Hash(), bs)

	return bs
}

// Calculate transactions to be included in the current rollup
func currentTxs(head *common2.Rollup, mempool []common2.L2Tx, db Db) []common2.L2Tx {
	return findTxsNotIncluded(head, mempool, db)
}

func FindWinner(parent *common2.Rollup, rollups []*common2.Rollup, db Db) (*common2.Rollup, bool) {
	var win = -1
	// todo - add statistics to determine why there are conflicts.
	for i, r := range rollups {
		switch {
		case r.Header.ParentHash != parent.Hash(): // ignore rollups from L2 forks
		case db.Height(r) <= db.Height(parent): // ignore rollups that are older than the parent
		case win == -1:
			win = i
		case r.ProofHeight(db) < rollups[win].ProofHeight(db): // ignore rollups generated with an older proof
		case r.ProofHeight(db) > rollups[win].ProofHeight(db): // newer rollups win
			win = i
		case r.Header.Nonce < rollups[win].Header.Nonce: // for rollups with the same proof, base on the nonce
			win = i
		}
	}
	if win == -1 {
		return nil, false
	}
	return rollups[win], true
}

func findRoundWinner(receivedRollups []*common2.Rollup, parent *common2.Rollup, parentState State, db Db) (*common2.Rollup, State) {
	win, found := FindWinner(parent, receivedRollups, db)
	if !found {
		panic("This should not happen for gossip rounds.")
	}
	// calculate the state to compare with what is in the Rollup
	s := newProcessedState(parentState)
	s = executeTransactions(win.Transactions, s)

	p := db.Parent(win).Proof(db)
	s = processDeposits(p, win.Proof(db), s, db)

	rootState := serialize(s.s)
	if rootState != win.Header.State {
		panic(fmt.Sprintf("Calculated a different state. This should not happen as there are no malicious actors yet. \nGot: %s\nExp: %s\nParent state:%v\nParent state:%s\nTxs:%v", rootState, win.Header.State, parentState, parent.Header.State, printTxs(win.Transactions)))
	}
	//todo - we need another root hash for withdrawals

	return win, s.s
}

// mutates the state
// process deposits from the proof of the parent rollup(exclusive) to the proof of the current rollup
func processDeposits(fromBlock *common.Block, toBlock *common.Block, s RollupState, db Db) RollupState {
	from := common.GenesisBlock.Hash()
	height := common.L1GenesisHeight
	if fromBlock != nil {
		from = fromBlock.Hash()
		height = fromBlock.Height(db)
		if !common.IsAncestor(fromBlock, toBlock, db) {
			panic("wtf")
		}

	}

	b := toBlock
	for {
		if b.Hash() == from {
			break
		}
		for _, tx := range b.Transactions {
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
		if b.Height(db) < height {
			panic("something went wrong")
		}
		p, f := b.Parent(db)
		if !f {
			panic("wtf")
		}
		b = p
	}
	return s
}

// given an L1 block, and the State as it was in the Parent block, calculates the State after the current block.
func calculateBlockState(b *common.Block, parentState BlockState, db Db) BlockState {
	rollups := extractRollups(b, db)
	newHead, found := FindWinner(parentState.Head, rollups, db)

	s := newProcessedState(parentState.State)

	// only change the state if there is a new l2 Head in the current block
	if found {
		s = executeTransactions(newHead.Transactions, s)
		p := db.Parent(newHead).Proof(db)
		s = processDeposits(p, newHead.Proof(db), s, db)
	} else {
		newHead = parentState.Head
	}

	bs := BlockState{
		Block:          b,
		Head:           newHead,
		State:          s.s,
		foundNewRollup: found,
	}
	return bs
}

func extractRollups(b *common.Block, db Db) []*common2.Rollup {
	rollups := make([]*common2.Rollup, 0)
	for _, t := range b.Transactions {
		// go through all rollup transactions
		if t.TxType == common.RollupTx {
			r := common2.DecodeRollup(t.Rollup)

			// Ignore rollups created with proofs from different L1 blocks
			// In case of L1 reorgs, rollups may end published on a fork
			if common.IsBlockAncestor(r.Header.L1Proof, b, db) {
				rollups = append(rollups, r)
			}
		}
	}
	return rollups
}
