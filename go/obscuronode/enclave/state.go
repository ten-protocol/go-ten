package enclave

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/core"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/db"

	"github.com/obscuronet/obscuro-playground/go/ethclient/mgmtcontractlib"

	"github.com/ethereum/go-ethereum/rlp"

	"github.com/obscuronet/obscuro-playground/go/log"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
)

// mutates the State
// header - the header of the rollup where this transaction will be included
// todo - remove nolint after the header starts being used
func executeTransactions(
	txs []nodecommon.L2Tx,
	s *db.State,
	header *nodecommon.Header, //nolint
) {
	for _, tx := range txs {
		executeTx(s, tx)
	}
	// fmt.Printf("w1: %v\n", is.w)
}

// mutates the state
func executeTx(s *db.State, tx nodecommon.L2Tx) {
	switch core.TxData(&tx).Type {
	case core.TransferTx:
		executeTransfer(s, tx)
	case core.WithdrawalTx:
		executeWithdrawal(s, tx)
	case core.DepositTx:
		executeDeposit(s, tx)
	default:
		panic("Invalid transaction type")
	}
}

func executeWithdrawal(s *db.State, tx nodecommon.L2Tx) {
	if txData := core.TxData(&tx); s.Balances[txData.From] >= txData.Amount {
		s.Balances[txData.From] -= txData.Amount
		s.Withdrawals = append(s.Withdrawals, tx.Hash())
	}
}

func executeTransfer(s *db.State, tx nodecommon.L2Tx) {
	if txData := core.TxData(&tx); s.Balances[txData.From] >= txData.Amount {
		s.Balances[txData.From] -= txData.Amount
		s.Balances[txData.To] += txData.Amount
	}
}

func executeDeposit(s *db.State, tx nodecommon.L2Tx) {
	t := core.TxData(&tx)
	v, f := s.Balances[t.To]
	if f {
		s.Balances[t.To] = v + t.Amount
	} else {
		s.Balances[t.To] = t.Amount
	}
}

// Determine the new canonical L2 head and calculate the State
// Uses cache-ing to map the Head rollup and the State to each L1Node block.
func updateState(b *types.Block, blockResolver db.BlockResolver, storage db.Storage, txHandler mgmtcontractlib.TxHandler) *db.BlockState {
	// This method is called recursively in case of Re-orgs. Stop when state was calculated already.
	val, found := storage.FetchBlockState(b.Hash())
	if found {
		return val
	}

	if b.NumberU64() == obscurocommon.L2GenesisHeight {
		return nil
	}

	rollups := extractRollups(b, blockResolver, txHandler)
	genesisRollup := storage.FetchGenesisRollup()

	// processing blocks before genesis, so there is nothing to do
	if genesisRollup == nil && len(rollups) == 0 {
		return nil
	}

	// Detect if the incoming block contains the genesis rollup, and generate an updated state.
	// Handle the case of the block containing the genesis being processed multiple times.
	genesisState, isGenesis := handleGenesisRollup(b, storage, rollups, genesisRollup)
	if isGenesis {
		return genesisState
	}

	// To calculate the state after the current block, we need the state after the parent.
	// If this point is reached, there is a parent state guaranteed, because the genesis is handled above
	parentState, parentFound := storage.FetchBlockState(b.ParentHash())
	if !parentFound {
		// go back and calculate the State of the Parent
		p, f := storage.FetchBlock(b.ParentHash())
		if !f {
			log.Log("Could not find block parent. This should not happen.")
			return nil
		}
		parentState = updateState(p, blockResolver, storage, txHandler)
	}

	if parentState == nil {
		log.Log(fmt.Sprintf("Something went wrong. There should be parent here. \n Block: %d - Block Parent: %d - Header: %+v",
			obscurocommon.ShortHash(b.Hash()),
			obscurocommon.ShortHash(b.Header().ParentHash),
			b.Header(),
		))
		return nil
	}

	bs := calculateBlockState(b, parentState, storage, blockResolver, rollups, txHandler)

	storage.SetBlockState(b.Hash(), bs)

	return bs
}

func handleGenesisRollup(b *types.Block, s db.Storage, rollups []*core.Rollup, genesisRollup *core.Rollup) (genesisState *db.BlockState, isGenesis bool) {
	// the incoming block holds the genesis rollup
	// calculate and return the new block state
	// todo change this to an hardcoded hash on testnet/mainnet
	if genesisRollup == nil && len(rollups) == 1 {
		log.Log("Found genesis rollup")

		genesis := rollups[0]
		s.StoreGenesisRollup(genesis)

		// The genesis rollup is part of the canonical chain and will be included in an L1 block by the first Aggregator.
		bs := db.BlockState{
			Block:          b,
			Head:           genesis,
			State:          db.EmptyState(),
			FoundNewRollup: true,
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
func currentTxs(head *core.Rollup, mempool []nodecommon.L2Tx, s db.Storage) []nodecommon.L2Tx {
	return findTxsNotIncluded(head, mempool, s)
}

func FindWinner(parent *core.Rollup, rollups []*core.Rollup, blockResolver db.BlockResolver) (*core.Rollup, bool) {
	win := -1
	// todo - add statistics to determine why there are conflicts.
	for i, r := range rollups {
		switch {
		case r.Header.ParentHash != parent.Hash(): // ignore rollups from L2 forks
		case r.Header.Number <= parent.Header.Number: // ignore rollups that are older than the parent
		case win == -1:
			win = i
		case blockResolver.ProofHeight(r) < blockResolver.ProofHeight(rollups[win]): // ignore rollups generated with an older proof
		case blockResolver.ProofHeight(r) > blockResolver.ProofHeight(rollups[win]): // newer rollups win
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

func (e *enclaveImpl) findRoundWinner(receivedRollups []*core.Rollup, parent *core.Rollup, parentState *db.State, s db.Storage, blockResolver db.BlockResolver) (*core.Rollup, *db.State) {
	headRollup, found := FindWinner(parent, receivedRollups, blockResolver)
	if !found {
		panic("This should not happen for gossip rounds.")
	}
	// calculate the state to compare with what is in the Rollup
	p := blockResolver.Proof(s.ParentRollup(headRollup))
	depositTxs := processDeposits(p, blockResolver.Proof(headRollup), blockResolver, e.txHandler)

	state := db.CopyStateNoWithdrawals(parentState)
	executeTransactions(append(headRollup.Transactions, depositTxs...), state, headRollup.Header)

	if db.Serialize(state) != headRollup.Header.State {
		panic(fmt.Sprintf("Calculated a different state. This should not happen as there are no malicious actors yet. \nGot: %s\nExp: %s\nParent state:%v\nParent state:%s\nTxs:%v",
			db.Serialize(state),
			headRollup.Header.State,
			parentState,
			parent.Header.State,
			printTxs(headRollup.Transactions)),
		)
	}
	// todo - check that the withdrawals in the header match the withdrawals as calculated

	return headRollup, state
}

// returns a list of L2 deposit transactions generated from the L1 deposit transactions
// starting with the proof of the parent rollup(exclusive) to the proof of the current rollup
func processDeposits(fromBlock *types.Block, toBlock *types.Block, blockResolver db.BlockResolver, txHandler mgmtcontractlib.TxHandler) []nodecommon.L2Tx {
	from := obscurocommon.GenesisBlock.Hash()
	height := obscurocommon.L1GenesisHeight
	if fromBlock != nil {
		from = fromBlock.Hash()
		height = fromBlock.NumberU64()
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
				depL2TxData := core.L2TxData{
					Type:   core.DepositTx,
					To:     t.Dest,
					Amount: t.Amount,
				}
				allDeposits = append(allDeposits, *newL2Tx(depL2TxData))
			}
		}
		if b.NumberU64() < height {
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
func calculateBlockState(b *types.Block, parentState *db.BlockState, s db.Storage, blockResolver db.BlockResolver, rollups []*core.Rollup, txHandler mgmtcontractlib.TxHandler) *db.BlockState {
	currentHead := parentState.Head
	newHeadRollup, found := FindWinner(currentHead, rollups, blockResolver)
	newState := parentState.State
	// only change the state if there is a new l2 Head in the current block
	if found {
		// Preprocessing before passing to the vm
		// todo transform into an eth block structure
		parentRollup := s.ParentRollup(newHeadRollup)
		p := blockResolver.Proof(parentRollup)
		depositTxs := processDeposits(p, blockResolver.Proof(newHeadRollup), blockResolver, txHandler)

		// deposits have to be processed after the normal transactions were executed because during speculative execution they are not available
		txsToProcess := append(newHeadRollup.Transactions, depositTxs...)
		newState = db.CopyStateNoWithdrawals(parentState.State)
		executeTransactions(txsToProcess, newState, newHeadRollup.Header)
	} else {
		newHeadRollup = parentState.Head
	}

	bs := db.BlockState{
		Block:          b,
		Head:           newHeadRollup,
		State:          newState,
		FoundNewRollup: found,
	}
	return &bs
}

// Todo - this has to be implemented differently based on how we define the ObsERC20
func rollupPostProcessingWithdrawals(newHeadRollup *core.Rollup, newState *db.State) []nodecommon.Withdrawal {
	w := make([]nodecommon.Withdrawal, 0)
	// go through each transaction and check if the withdrawal was processed correctly
	for i, t := range newHeadRollup.Transactions {
		txData := core.TxData(&newHeadRollup.Transactions[i])
		if txData.Type == core.WithdrawalTx && contains(newState.Withdrawals, t.Hash()) {
			w = append(w, nodecommon.Withdrawal{
				Amount:  txData.Amount,
				Address: txData.From,
			})
		}
	}

	return w
}

func extractRollups(b *types.Block, blockResolver db.BlockResolver, handler mgmtcontractlib.TxHandler) []*core.Rollup {
	rollups := make([]*core.Rollup, 0)
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

func toEnclaveRollup(r *nodecommon.Rollup) *core.Rollup {
	return &core.Rollup{
		Header:       r.Header,
		Transactions: core.DecryptTransactions(r.Transactions),
	}
}

func newL2Tx(data core.L2TxData) *nodecommon.L2Tx {
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
