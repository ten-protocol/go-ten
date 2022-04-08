package enclave

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/rlp"

	"github.com/obscuronet/obscuro-playground/go/log"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
)

type State = map[common.Address]uint64

// blockState - Represents the state after an L1 block was processed.
type blockState struct {
	block          *types.Block
	head           *Rollup
	state          State
	foundNewRollup bool
}

// RollupState - State after an L2 rollup was processed
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
func executeTransactions(txs []nodecommon.L2Tx, state RollupState) RollupState {
	ps := copyProcessedState(state)
	for _, tx := range txs {
		executeTx(&ps, tx)
	}
	// fmt.Printf("w1: %v\n", is.w)
	return ps
}

// mutates the state
func executeTx(s *RollupState, tx nodecommon.L2Tx) {
	switch TxData(&tx).Type {
	case TransferTx:
		executeTransfer(s, tx)
	case WithdrawalTx:
		executeWithdrawal(s, tx)
	case DepositTx:
		executeDeposit(s, tx)
	default:
		panic("Invalid transaction type")
	}
}

func executeWithdrawal(s *RollupState, tx nodecommon.L2Tx) {
	if txData := TxData(&tx); s.s[txData.From] >= txData.Amount {
		s.s[txData.From] -= txData.Amount
		s.w = append(s.w, nodecommon.Withdrawal{
			Amount:  txData.Amount,
			Address: txData.From,
		})
	}
}

func executeTransfer(s *RollupState, tx nodecommon.L2Tx) {
	if txData := TxData(&tx); s.s[txData.From] >= txData.Amount {
		s.s[txData.From] -= txData.Amount
		s.s[txData.To] += txData.Amount
	}
}

func executeDeposit(s *RollupState, tx nodecommon.L2Tx) {
	t := TxData(&tx)
	v, f := s.s[t.To]
	if f {
		s.s[t.To] = v + t.Amount
	} else {
		s.s[t.To] = t.Amount
	}
}

func emptyState() State {
	return make(State)
}

// Determine the new canonical L2 head and calculate the State
// Uses cache-ing to map the Head rollup and the State to each L1Node block.
func updateState(b *types.Block, s Storage, blockResolver BlockResolver) *blockState {
	// This method is called recursively in case of Re-orgs. Stop when state was calculated already.
	val, found := s.FetchBlockState(b.Hash())
	if found {
		return val
	}

	if blockResolver.HeightBlock(b) == 0 {
		return nil
	}

	rollups := extractRollups(b, blockResolver)
	genesisRollup := s.FetchGenesisRollup()

	// processing blocks before genesis, so there is nothing to do
	if genesisRollup == nil && len(rollups) == 0 {
		return nil
	}

	// Detect if the incoming block contains the genesis rollup, and generate an updated state.
	// Handle the case of the block containing the genesis being processed multiple times.
	genesisState, isGenesis := handleGenesisRollup(b, s, rollups, genesisRollup)
	if isGenesis {
		return genesisState
	}

	// To calculate the state after the current block, we need the state after the parent.
	// If this point is reached, there is a parent state guaranteed, because the genesis is handled above
	parentState, parentFound := s.FetchBlockState(b.ParentHash())
	if !parentFound {
		// go back and calculate the State of the Parent
		p, f := s.FetchBlock(b.ParentHash())
		if !f {
			panic("Could not find block parent. This should not happen.")
		}
		parentState = updateState(p, s, blockResolver)
	}

	bs := calculateBlockState(b, parentState, s, blockResolver, rollups)

	s.SetBlockState(b.Hash(), bs)

	return bs
}

func handleGenesisRollup(b *types.Block, s Storage, rollups []*Rollup, genesisRollup *Rollup) (genesisState *blockState, isGenesis bool) {
	// the incoming block holds the genesis rollup
	// calculate and return the new block state
	// todo change this to an hardcoded hash on testnet/mainnet
	if genesisRollup == nil && len(rollups) == 1 {
		log.Log("Found genesis rollup")

		genesis := rollups[0]
		s.StoreGenesisRollup(genesis)

		// The genesis rollup is part of the canonical chain and will be included in an L1 block by the first Aggregator.
		bs := blockState{
			block:          b,
			head:           genesis,
			state:          emptyState(),
			foundNewRollup: true,
		}
		s.SetBlockState(b.Hash(), &bs)
		return &bs, true
	}

	// Re-processing the block that contains the rollup. This can happen as blocks can be fed to the enclave multiple times.
	// In this case we don't update the state and move on.
	if genesisRollup != nil && len(rollups) == 1 && rollups[0].Header.Hash() == genesisRollup.Hash() {
		return nil, true
	}
	return nil, false
}

// Calculate transactions to be included in the current rollup
func currentTxs(head *Rollup, mempool []nodecommon.L2Tx, s Storage) []nodecommon.L2Tx {
	return findTxsNotIncluded(head, mempool, s)
}

func FindWinner(parent *Rollup, rollups []*Rollup, s Storage, blockResolver BlockResolver) (*Rollup, bool) {
	win := -1
	// todo - add statistics to determine why there are conflicts.
	for i, r := range rollups {
		switch {
		case r.Header.ParentHash != parent.Hash(): // ignore rollups from L2 forks
		case r.Header.Height <= parent.Header.Height: // ignore rollups that are older than the parent
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

func findRoundWinner(receivedRollups []*Rollup, parent *Rollup, parentState State, s Storage, blockResolver BlockResolver) (*Rollup, State) {
	win, found := FindWinner(parent, receivedRollups, s, blockResolver)
	if !found {
		panic("This should not happen for gossip rounds.")
	}
	// calculate the state to compare with what is in the Rollup
	state := newProcessedState(parentState)
	p := s.ParentRollup(win).Proof(blockResolver)
	depositTxs := processDeposits(p, win.Proof(blockResolver), blockResolver)

	state = executeTransactions(append(win.Transactions, depositTxs...), state)

	if serialize(state.s) != win.Header.State {
		panic(fmt.Sprintf("Calculated a different state. This should not happen as there are no malicious actors yet. \nGot: %s\nExp: %s\nParent state:%v\nParent state:%s\nTxs:%v",
			serialize(state.s),
			win.Header.State,
			parentState,
			parent.Header.State,
			printTxs(win.Transactions)),
		)
	}
	// todo - check that the withdrawals in the header match the withdrawals as calculated

	return win, state.s
}

// mutates the state
// process deposits from the proof of the parent rollup(exclusive) to the proof of the current rollup
func processDeposits(fromBlock *types.Block, toBlock *types.Block, blockResolver BlockResolver) []nodecommon.L2Tx {
	from := obscurocommon.GenesisBlock.Hash()
	height := obscurocommon.L1GenesisHeight
	if fromBlock != nil {
		from = fromBlock.Hash()
		height = blockResolver.HeightBlock(fromBlock)
		if !blockResolver.IsAncestor(toBlock, fromBlock) {
			panic("Deposits can't be processed because the rollups are not on the same Ethereum fork. This should not happen.")
		}
	}

	allDeposits := make([]nodecommon.L2Tx, 0)
	b := toBlock
	for {
		if b.Hash() == from {
			break
		}
		for _, tx := range b.Transactions() {
			t := obscurocommon.TxData(tx)
			// transactions to a hardcoded bridge address
			if t.TxType == obscurocommon.DepositTx {
				depL2TxData := L2TxData{
					Type:   DepositTx,
					To:     t.Dest,
					Amount: t.Amount,
				}
				allDeposits = append(allDeposits, *newL2Tx(depL2TxData))
			}
		}
		if blockResolver.HeightBlock(b) < height {
			panic("something went wrong")
		}
		p, f := blockResolver.ParentBlock(b)
		if !f {
			panic("Deposits can't be processed because the rollups are not on the same Ethereum fork. This should not happen.")
		}
		b = p
	}

	return allDeposits
}

// given an L1 block, and the State as it was in the Parent block, calculates the State after the current block.
func calculateBlockState(b *types.Block, parentState *blockState, s Storage, blockResolver BlockResolver, rollups []*Rollup) *blockState {
	currentHead := parentState.head
	newHead, found := FindWinner(currentHead, rollups, s, blockResolver)

	state := newProcessedState(parentState.state)

	// only change the state if there is a new l2 Head in the current block
	if found {
		// Preprocessing before passing to the vm
		// todo transform into an eth block structure

		p := s.ParentRollup(newHead).Proof(blockResolver)
		depositTxs := processDeposits(p, newHead.Proof(blockResolver), blockResolver)

		// deposits have to be processed after the normal transactions were executed because during speculative execution they are not available
		txsToProcess := append(newHead.Transactions, depositTxs...)
		state = executeTransactions(txsToProcess, state)

		// Postprocessing - witdrawals todo
	} else {
		newHead = parentState.head
	}

	bs := blockState{
		block:          b,
		head:           newHead,
		state:          state.s,
		foundNewRollup: found,
	}
	return &bs
}

func extractRollups(b *types.Block, blockResolver BlockResolver) []*Rollup {
	rollups := make([]*Rollup, 0)
	for _, t := range b.Transactions() {
		// go through all rollup transactions
		data := obscurocommon.TxData(t)
		if data.TxType == obscurocommon.RollupTx {
			r := nodecommon.DecodeRollupOrPanic(obscurocommon.TxData(t).Rollup)

			// Ignore rollups created with proofs from different L1 blocks
			// In case of L1 reorgs, rollups may end published on a fork
			if blockResolver.IsBlockAncestor(b, r.Header.L1Proof) {
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

func newL2Tx(data L2TxData) *nodecommon.L2Tx {
	// We should probably use a deterministic nonce instead, as in the L1.
	nonce, _ := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))

	enc, err := rlp.EncodeToBytes(data)
	if err != nil {
		// TODO - Surface this error properly.
		panic(err)
	}

	return types.NewTx(&types.LegacyTx{
		Nonce:    nonce.Uint64(),
		Value:    big.NewInt(1),
		Gas:      1,
		GasPrice: big.NewInt(1),
		Data:     enc,
	})
}
