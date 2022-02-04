package obscuro

import (
	"fmt"
	"simulation/common"
	"simulation/wallet-mock"
)

type State = map[wallet_mock.Address]uint64

// internal structure to pass information. todo - prob not necessary
type currentWork struct {
	r   common.Rollup
	s   ProcessedState
	txs []common.L2Tx
}

// BlockState - Represents the state after an L1 block was processed.
type BlockState struct {
	Block          common.Block
	Head           common.Rollup
	State          State
	foundNewRollup bool
}

type ProcessedState struct {
	s State
	w []common.Withdrawal
}

func newProcessedState(s State) ProcessedState {
	return ProcessedState{
		s: copyState(s),
		w: []common.Withdrawal{},
	}
}

func copyState(state State) State {
	s := make(State)
	for address, balance := range state {
		s[address] = balance
	}
	return s
}

func copyProcessedState(s ProcessedState) ProcessedState {
	return ProcessedState{
		s: copyState(s.s),
		w: s.w,
	}
}

func serialize(state State) string {
	return fmt.Sprintf("%v", state)
}

// returns a modified copy of the State
func executeTransactions(txs []common.L2Tx, state ProcessedState) ProcessedState {
	is := copyProcessedState(state)
	for _, tx := range txs {
		executeTx(&is, tx)
	}
	//fmt.Printf("w1: %v\n", is.w)
	return is
}

// mutates the State
func executeTx(s *ProcessedState, tx common.L2Tx) {
	switch tx.TxType {
	case common.TransferTx:
		executeTransfer(s, tx)
	case common.WithdrawalTx:
		executeWithdrawal(s, tx)
	default:
		panic("Invalid transaction type")
	}
}

func executeWithdrawal(s *ProcessedState, tx common.L2Tx) {
	bal := s.s[tx.From]
	if bal >= tx.Amount {
		s.s[tx.From] -= tx.Amount
		s.w = append(s.w, common.Withdrawal{
			Amount:  tx.Amount,
			Address: tx.From,
		})
		//fmt.Printf("w: %v\n", s.w)
	}
}

func executeTransfer(s *ProcessedState, tx common.L2Tx) {
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
// Uses cache-ing to map the PeekHead rollup and the State to each L1Node block.
func updateState(b common.Block, db Db) BlockState {

	// This method is called recursively in case of Re-orgs. Stop when state was calculated already.
	val, found := db.FetchState(b.RootHash)
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
		db.SetState(b.RootHash, bs)
		return bs
	}

	// To calculate the state after the current block, we need the state after the parent.
	parentState, parentFound := db.FetchState(b.Parent().Root())
	if !parentFound {
		// go back and calculate the State of the Parent
		parentState = updateState(*b.ParentBlock(), db)
	}

	bs := calculateBlockState(b, parentState)

	db.SetState(b.RootHash, bs)

	return bs
}

// Calculate transactions to be included in the current rollup
func currentTxs(head common.Rollup, mempool []common.L2Tx) []common.L2Tx {
	mempoolCopy := make([]common.Tx, len(mempool))
	for i, tx := range mempool {
		mempoolCopy[i] = tx
	}
	toInclude := common.FindNotIncludedTxs(head, mempoolCopy)
	txsCopy := make([]common.L2Tx, len(toInclude))
	for i, tx := range toInclude {
		txsCopy[i] = tx.(common.L2Tx)
	}

	return txsCopy
}

func findRoundWinner(receivedRollups []common.Rollup, parent common.Rollup, parentState State) (common.Rollup, State) {
	var win = -1
	for i, r := range receivedRollups {
		switch {
		case r.Parent().Root() != parent.Root(): // ignore rollups from forks
		case win == -1:
			win = i
		case r.Proof().Height() < receivedRollups[win].Proof().Height(): // ignore rollups generated with an older proof
		case r.Proof().Height() > receivedRollups[win].Proof().Height(): // newer rollups win
			win = i
		case r.Nonce < receivedRollups[win].Nonce: // for rollups with the same proof, base on the nonce
			win = i
		}
		//common.Log(fmt.Sprintf("win: r_%d, r: %d", receivedRollups[win].RootHash.ID(), r.RootHash.ID()))
	}

	// calculate the state to compare with what is in the Rollup
	s := newProcessedState(parentState)
	p := receivedRollups[win].ParentRollup().Proof()
	s = processDeposits(&p, receivedRollups[win].Proof(), s)
	s = executeTransactions(receivedRollups[win].L2Txs(), s)
	// todo - check that s is valid against the State in the rollup, if not - call the function again with this tx excluded
	return receivedRollups[win], s.s
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

// given an L1 block, and the State as it was in the Parent block, calculates the State after the current block.
func calculateBlockState(b common.Block, parentState BlockState) BlockState {
	var newHead *common.Rollup = nil

	// find the PeekHead rollup according to the rules
	for _, t := range b.L1Txs() {
		// go through all rollup transactions
		if t.TxType == common.RollupTx {
			r := common.DecodeRollup(t.Rollup)
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
	// only change the state if there is a new l2 PeekHead in the current block
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
