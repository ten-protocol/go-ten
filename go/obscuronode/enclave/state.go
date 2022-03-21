package enclave

import (
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
)

type State = map[common.Address]uint64

// BlockState - Represents the state after an L1 block was processed.
type BlockState struct {
	Block          *types.Block
	Head           *Rollup
	State          State
	foundNewRollup bool
}

// RollupState - state after an L2 rollups was processed
type RollupState struct {
	s State
	w []nodecommon.Withdrawal
}

func newProcessedState(s State) RollupState {
	return RollupState{
		s: copyState(s),
		w: []nodecommon.Withdrawal{},
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
func executeTransactions(txs []L2Tx, state RollupState) RollupState {
	ps := copyProcessedState(state)
	for _, tx := range txs {
		executeTx(&ps, tx)
	}
	// fmt.Printf("w1: %v\n", is.w)
	return ps
}

// mutates the State
func executeTx(s *RollupState, tx L2Tx) {
	switch TxData(&tx).Type {
	case TransferTx:
		executeTransfer(s, tx)
	case WithdrawalTx:
		executeWithdrawal(s, tx)
	default:
		panic("Invalid transaction type")
	}
}

func executeWithdrawal(s *RollupState, tx L2Tx) {
	if txData := TxData(&tx); s.s[txData.From] >= txData.Amount {
		s.s[txData.From] -= txData.Amount
		s.w = append(s.w, nodecommon.Withdrawal{
			Amount:  txData.Amount,
			Address: txData.From,
		})
	}
}

func executeTransfer(s *RollupState, tx L2Tx) {
	if txData := TxData(&tx); s.s[txData.From] >= txData.Amount {
		s.s[txData.From] -= txData.Amount
		s.s[txData.To] += txData.Amount
	}
}

func emptyState() State {
	return make(State)
}

// Determine the new canonical L2 head and calculate the State
// Uses cache-ing to map the Head rollup and the State to each L1Node block.
func updateState(b *types.Block, db DB, blockResolver obscurocommon.BlockResolver) BlockState {
	// This method is called recursively in case of Re-orgs. Stop when state was calculated already.
	val, found := db.FetchState(b.Hash())
	if found {
		return val
	}

	// The genesis rollup is part of the canonical chain and will be included in an L1 block by the first Aggregator.
	if b.Hash() == obscurocommon.GenesisBlock.Hash() {
		bs := BlockState{
			Block:          b,
			Head:           &GenesisRollup,
			State:          emptyState(),
			foundNewRollup: true,
		}
		db.SetState(b.Hash(), bs)
		return bs
	}

	// To calculate the state after the current block, we need the state after the parent.
	parentState, parentFound := db.FetchState(b.ParentHash())
	if !parentFound {
		// go back and calculate the State of the Parent
		p, f := db.ParentBlock(b)
		if !f {
			panic("wtf")
		}
		parentState = updateState(p, db, blockResolver)
	}

	bs := calculateBlockState(b, parentState, db, blockResolver)

	db.SetState(b.Hash(), bs)

	return bs
}

// Calculate transactions to be included in the current rollup
func currentTxs(head *Rollup, mempool []L2Tx, db DB) []L2Tx {
	return findTxsNotIncluded(head, mempool, db)
}

func FindWinner(parent *Rollup, rollups []*Rollup, db DB, blockResolver obscurocommon.BlockResolver) (*Rollup, bool) {
	win := -1
	// todo - add statistics to determine why there are conflicts.
	for i, r := range rollups {
		switch {
		case r.Header.ParentHash != parent.Hash(): // ignore rollups from L2 forks
		case db.HeightRollup(r) <= db.HeightRollup(parent): // ignore rollups that are older than the parent
		case win == -1:
			win = i
		case r.ProofHeight(blockResolver) < rollups[win].ProofHeight(blockResolver): // ignore rollups generated with an older proof
		case r.ProofHeight(blockResolver) > rollups[win].ProofHeight(blockResolver): // newer rollups win
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

func findRoundWinner(receivedRollups []*Rollup, parent *Rollup, parentState State, db DB, blockResolver obscurocommon.BlockResolver) (*Rollup, State) {
	win, found := FindWinner(parent, receivedRollups, db, blockResolver)
	if !found {
		panic("This should not happen for gossip rounds.")
	}
	// calculate the state to compare with what is in the Rollup
	s := newProcessedState(parentState)
	s = executeTransactions(win.Transactions, s)

	p := db.ParentRollup(win).Proof(blockResolver)
	s = processDeposits(p, win.Proof(blockResolver), s, blockResolver)

	if serialize(s.s) != win.Header.State {
		panic(fmt.Sprintf("Calculated a different state. This should not happen as there are no malicious actors yet. \nGot: %s\nExp: %s\nParent state:%v\nParent state:%s\nTxs:%v",
			serialize(s.s),
			win.Header.State,
			parentState,
			parent.Header.State,
			printTxs(win.Transactions)),
		)
	}
	// todo - check that the withdrawals in the header match the withdrawals as calculated

	return win, s.s
}

// mutates the state
// process deposits from the proof of the parent rollup(exclusive) to the proof of the current rollup
func processDeposits(fromBlock *types.Block, toBlock *types.Block, s RollupState, blockResolver obscurocommon.BlockResolver) RollupState {
	from := obscurocommon.GenesisBlock.Hash()
	height := obscurocommon.L1GenesisHeight
	if fromBlock != nil {
		from = fromBlock.Hash()
		height = blockResolver.HeightBlock(fromBlock)
		if !obscurocommon.IsAncestor(fromBlock, toBlock, blockResolver) {
			panic("wtf")
		}
	}

	b := toBlock
	for {
		if b.Hash() == from {
			break
		}
		for _, tx := range b.Transactions() {
			t := obscurocommon.TxData(tx)
			// transactions to a hardcoded bridge address
			if t.TxType == obscurocommon.DepositTx {
				v, f := s.s[t.Dest]
				if f {
					s.s[t.Dest] = v + t.Amount
				} else {
					s.s[t.Dest] = t.Amount
				}
			}
		}
		if blockResolver.HeightBlock(b) < height {
			panic("something went wrong")
		}
		p, f := blockResolver.ParentBlock(b)
		if !f {
			panic("wtf")
		}
		b = p
	}
	return s
}

// given an L1 block, and the State as it was in the Parent block, calculates the State after the current block.
func calculateBlockState(b *types.Block, parentState BlockState, db DB, blockResolver obscurocommon.BlockResolver) BlockState {
	rollups := extractRollups(b, blockResolver)
	newHead, found := FindWinner(parentState.Head, rollups, db, blockResolver)

	s := newProcessedState(parentState.State)

	// only change the state if there is a new l2 Head in the current block
	if found {
		s = executeTransactions(newHead.Transactions, s)
		p := db.ParentRollup(newHead).Proof(blockResolver)
		s = processDeposits(p, newHead.Proof(blockResolver), s, blockResolver)
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

func extractRollups(b *types.Block, blockResolver obscurocommon.BlockResolver) []*Rollup {
	rollups := make([]*Rollup, 0)
	for _, t := range b.Transactions() {
		// go through all rollup transactions
		data := obscurocommon.TxData(t)
		if data.TxType == obscurocommon.RollupTx {
			r := nodecommon.DecodeRollupOrPanic(obscurocommon.TxData(t).Rollup)

			// Ignore rollups created with proofs from different L1 blocks
			// In case of L1 reorgs, rollups may end published on a fork
			if obscurocommon.IsBlockAncestor(r.Header.L1Proof, b, blockResolver) {
				rollups = append(rollups, toEnclaveRollup(r))
			}
		}
	}
	return rollups
}

func toEnclaveRollup(r *nodecommon.Rollup) *Rollup {
	return &Rollup{
		Header:       r.Header,
		Transactions: decryptTransactions(r.Transactions),
	}
}
