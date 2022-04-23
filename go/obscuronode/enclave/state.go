package enclave

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"

	"github.com/obscuronet/obscuro-playground/go/l1client/rollupcontractlib"

	"github.com/obscuronet/obscuro-playground/go/hashing"

	"github.com/ethereum/go-ethereum/rlp"

	"github.com/obscuronet/obscuro-playground/go/log"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
)

// State - this is a placeholder for the real Trie based state
//- people send transactions to an ObsERC20 that was a withdraw(amount, from, destination) method
//In the EVM, there will be a smart contract that does the following:
//- the tokens are deducted from the "from" address , and burned
//- add to the "withdrawals" transactions - this info will be taken from the state
//Post processing, outside the evm:
//- generate withdrawal instructions (amount, destination), based on which withdrawal transaction were executed successfully
type State struct {
	balances    map[common.Address]uint64
	withdrawals []obscurocommon.TxHash
}

// blockState - Represents the state after an L1 block was processed.
type blockState struct {
	block          *types.Block
	head           *Rollup
	state          *State
	foundNewRollup bool
}

func copyState(state *State) *State {
	s := emptyState()
	if state == nil {
		return s
	}
	for address, balance := range state.balances {
		s.balances[address] = balance
	}
	s.withdrawals = append(s.withdrawals, state.withdrawals...)
	return s
}

func serialize(state *State) nodecommon.StateRoot {
	hash, err := hashing.RLPHash(fmt.Sprintf("%v", state))
	if err != nil {
		panic(err)
	}
	return hash
}

// returns a modified copy of the State
// header - the header of the rollup where this transaction will be included
// todo - remove nolint after the header starts being used
func (e *enclaveImpl) executeTransactions(
	txs []nodecommon.L2Tx,
	state *State,
	header *nodecommon.Header, //nolint
) *State {
	s := copyState(state)
	for _, tx := range txs {
		executeTx(s, tx)
	}
	// fmt.Printf("w1: %v\n", is.w)
	return s
}

// mutates the state
func executeTx(s *State, tx nodecommon.L2Tx) {
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

func executeWithdrawal(s *State, tx nodecommon.L2Tx) {
	if txData := TxData(&tx); s.balances[txData.From] >= txData.Amount {
		s.balances[txData.From] -= txData.Amount
		s.withdrawals = append(s.withdrawals, tx.Hash())
	}
}

func executeTransfer(s *State, tx nodecommon.L2Tx) {
	if txData := TxData(&tx); s.balances[txData.From] >= txData.Amount {
		s.balances[txData.From] -= txData.Amount
		s.balances[txData.To] += txData.Amount
	}
}

func executeDeposit(s *State, tx nodecommon.L2Tx) {
	t := TxData(&tx)
	v, f := s.balances[t.To]
	if f {
		s.balances[t.To] = v + t.Amount
	} else {
		s.balances[t.To] = t.Amount
	}
}

func emptyState() *State {
	return &State{
		balances: map[common.Address]uint64{},
	}
}

// Determine the new canonical L2 head and calculate the State
// Uses cache-ing to map the Head rollup and the State to each L1Node block.
func (e *enclaveImpl) updateState(b *types.Block, blockResolver BlockResolver) *blockState {
	// This method is called recursively in case of Re-orgs. Stop when state was calculated already.
	val, found := e.storage.FetchBlockState(b.Hash())
	if found {
		return val
	}

	if blockResolver.HeightBlock(b) == 0 {
		return nil
	}

	rollups := extractRollups(b, blockResolver, e.txHandler)
	genesisRollup := e.storage.FetchGenesisRollup()

	// processing blocks before genesis, so there is nothing to do
	if genesisRollup == nil && len(rollups) == 0 {
		return nil
	}

	// Detect if the incoming block contains the genesis rollup, and generate an updated state.
	// Handle the case of the block containing the genesis being processed multiple times.
	genesisState, isGenesis := handleGenesisRollup(b, e.storage, rollups, genesisRollup)
	if isGenesis {
		return genesisState
	}

	// To calculate the state after the current block, we need the state after the parent.
	// If this point is reached, there is a parent state guaranteed, because the genesis is handled above
	parentState, parentFound := e.storage.FetchBlockState(b.ParentHash())
	if !parentFound {
		// go back and calculate the State of the Parent
		p, f := e.storage.FetchBlock(b.ParentHash())
		if !f {
			panic("Could not find block parent. This should not happen.")
		}
		parentState = e.updateState(p, blockResolver)
	}

	if parentState == nil {
		panic("Something went wrong. There should be parent here.")
	}

	bs := e.calculateBlockState(b, parentState, e.storage, blockResolver, rollups)

	e.storage.SetBlockState(b.Hash(), bs)

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

func FindWinner(parent *Rollup, rollups []*Rollup, blockResolver BlockResolver) (*Rollup, bool) {
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

func (e *enclaveImpl) findRoundWinner(receivedRollups []*Rollup, parent *Rollup, parentState *State, s Storage, blockResolver BlockResolver) (*Rollup, *State, bool) {
	headRollup, found := FindWinner(parent, receivedRollups, blockResolver)
	if !found {
		panic("This should not happen for gossip rounds.")
	}
	// calculate the state to compare with what is in the Rollup
	p := s.ParentRollup(headRollup).Proof(blockResolver)
	depositTxs := processDeposits(p, headRollup.Proof(blockResolver), blockResolver, e.txHandler)

	state := e.executeTransactions(append(headRollup.Transactions, depositTxs...), parentState, headRollup.Header)

	if serialize(state) != headRollup.Header.State {
		panic(fmt.Sprintf("Calculated a different state. This should not happen as there are no malicious actors yet. \nGot: %s\nExp: %s\nParent state:%v\nParent state:%s\nTxs:%v",
			serialize(state),
			headRollup.Header.State,
			parentState,
			parent.Header.State,
			printTxs(headRollup.Transactions)),
		)
	}
	// todo - check that the withdrawals in the header match the withdrawals as calculated

	return headRollup, state, len(headRollup.Transactions) > 0 || len(depositTxs) > 0
}

// returns a list of L2 deposit transactions generated from the L1 deposit transactions
// starting with the proof of the parent rollup(exclusive) to the proof of the current rollup
func processDeposits(fromBlock *types.Block, toBlock *types.Block, blockResolver BlockResolver, txHandler rollupcontractlib.TxHandler) []nodecommon.L2Tx {
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
			t := txHandler.UnPackTx(tx)
			// transactions to a hardcoded bridge address
			if t != nil && t.TxType == obscurocommon.DepositTx {
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
func (e *enclaveImpl) calculateBlockState(b *types.Block, parentState *blockState, s Storage, blockResolver BlockResolver, rollups []*Rollup) *blockState {
	currentHead := parentState.head
	newHeadRollup, found := FindWinner(currentHead, rollups, blockResolver)
	newState := parentState.state
	// only change the state if there is a new l2 Head in the current block
	if found {
		// Preprocessing before passing to the vm
		// todo transform into an eth block structure
		parentRollup := s.ParentRollup(newHeadRollup)
		p := parentRollup.Proof(blockResolver)
		depositTxs := processDeposits(p, newHeadRollup.Proof(blockResolver), blockResolver, e.txHandler)

		// deposits have to be processed after the normal transactions were executed because during speculative execution they are not available
		txsToProcess := append(newHeadRollup.Transactions, depositTxs...)
		newState = e.executeTransactions(txsToProcess, parentState.state, newHeadRollup.Header)
	} else {
		newHeadRollup = parentState.head
	}

	bs := blockState{
		block:          b,
		head:           newHeadRollup,
		state:          newState,
		foundNewRollup: found,
	}
	return &bs
}

// Todo - this has to be implemented differently based on how we define the ObsERC20
func rollupPostProcessingWithdrawals(newHeadRollup *Rollup, newState *State) []nodecommon.Withdrawal {
	w := make([]nodecommon.Withdrawal, 0)

	// go through each transaction and check if the withdrawal was processed correctly
	for i, t := range newHeadRollup.Transactions {
		txData := TxData(&newHeadRollup.Transactions[i])
		if txData.Type == WithdrawalTx && contains(newState.withdrawals, t.Hash()) {
			w = append(w, nodecommon.Withdrawal{
				Amount:  txData.Amount,
				Address: txData.From,
			})
		}
	}
	return w
}

func extractRollups(b *types.Block, blockResolver BlockResolver, handler rollupcontractlib.TxHandler) []*Rollup {
	rollups := make([]*Rollup, 0)
	for _, tx := range b.Transactions() {
		// go through all rollup transactions
		t := handler.UnPackTx(tx)
		if t != nil && t.TxType == obscurocommon.RollupTx {
			r := nodecommon.DecodeRollupOrPanic(t.Rollup)

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
