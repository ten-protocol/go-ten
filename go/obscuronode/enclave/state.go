package enclave

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/obscuro-playground/go/ethclient/txhandler"
	"github.com/obscuronet/obscuro-playground/go/log"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/core"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/db"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
)

// mutates the State
// header - the header of the rollup where this transaction will be included
// todo - remove nolint after the header starts being used
func executeTransactions(
	txs []nodecommon.L2Tx,
	s db.StateDB,
	header *nodecommon.Header, //nolint
) {
	for _, tx := range txs {
		executeTx(s, tx)
	}
	// fmt.Printf("w1: %v\n", is.w)
}

// mutates the state
func executeTx(s db.StateDB, tx nodecommon.L2Tx) {
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

func executeWithdrawal(s db.StateDB, tx nodecommon.L2Tx) {
	txData := core.TxData(&tx)
	from := s.GetBalance(txData.From)
	if from >= txData.Amount {
		s.SetBalance(txData.From, from-txData.Amount)
		s.AddWithdrawal(tx.Hash())
	}
}

func executeTransfer(s db.StateDB, tx nodecommon.L2Tx) {
	txData := core.TxData(&tx)
	from := s.GetBalance(txData.From)
	to := s.GetBalance(txData.To)

	if from >= txData.Amount {
		s.SetBalance(txData.From, from-txData.Amount)
		s.SetBalance(txData.To, to+txData.Amount)
	}
}

func executeDeposit(s db.StateDB, tx nodecommon.L2Tx) {
	txData := core.TxData(&tx)
	to := s.GetBalance(txData.To)
	s.SetBalance(txData.To, to+txData.Amount)
}

// Determine the new canonical L2 head and calculate the State
// Uses cache-ing to map the Head rollup and the State to each L1Node block.
func updateState(b *types.Block, blockResolver db.BlockResolver, txHandler txhandler.TxHandler, rollupResolver db.RollupResolver, bss db.BlockStateStorage) *core.BlockState {
	// This method is called recursively in case of Re-orgs. Stop when state was calculated already.
	val, found := bss.FetchBlockState(b.Hash())
	if found {
		return val
	}

	rollups := extractRollups(b, blockResolver, txHandler)
	genesisRollup := rollupResolver.FetchGenesisRollup()

	// processing blocks before genesis, so there is nothing to do
	if genesisRollup == nil && len(rollups) == 0 {
		return nil
	}

	// Detect if the incoming block contains the genesis rollup, and generate an updated state.
	// Handle the case of the block containing the genesis being processed multiple times.
	genesisState, isGenesis := handleGenesisRollup(b, rollups, genesisRollup, rollupResolver, bss)
	if isGenesis {
		return genesisState
	}

	// To calculate the state after the current block, we need the state after the parent.
	// If this point is reached, there is a parent state guaranteed, because the genesis is handled above
	parentState, parentFound := bss.FetchBlockState(b.ParentHash())
	if !parentFound {
		// go back and calculate the State of the Parent
		p, f := blockResolver.FetchBlock(b.ParentHash())
		if !f {
			log.Log("Could not find block parent. This should not happen.")
			return nil
		}
		parentState = updateState(p, blockResolver, txHandler, rollupResolver, bss)
	}

	if parentState == nil {
		log.Log(fmt.Sprintf("Something went wrong. There should be parent here. \n Block: %d - Block Parent: %d - Header: %+v",
			obscurocommon.ShortHash(b.Hash()),
			obscurocommon.ShortHash(b.Header().ParentHash),
			b.Header(),
		))
		return nil
	}

	bs, stateDB, head := calculateBlockState(b, parentState, blockResolver, rollups, txHandler, rollupResolver, bss)
	log.Trace(fmt.Sprintf("- Calc block state b_%d: Found: %t - r_%d, ",
		obscurocommon.ShortHash(b.Hash()),
		bs.FoundNewRollup,
		obscurocommon.ShortHash(bs.HeadRollup),
	))

	bss.SetBlockState(b.Hash(), bs, head)
	if bs.FoundNewRollup {
		stateDB.Commit(bs.HeadRollup)
	}

	return bs
}

func handleGenesisRollup(b *types.Block, rollups []*core.Rollup, genesisRollup *core.Rollup, resolver db.RollupResolver, bss db.BlockStateStorage) (genesisState *core.BlockState, isGenesis bool) {
	// the incoming block holds the genesis rollup
	// calculate and return the new block state
	// todo change this to an hardcoded hash on testnet/mainnet
	if genesisRollup == nil && len(rollups) == 1 {
		log.Log("Found genesis rollup")

		genesis := rollups[0]
		resolver.StoreGenesisRollup(genesis)

		// The genesis rollup is part of the canonical chain and will be included in an L1 block by the first Aggregator.
		bs := core.BlockState{
			Block:          b.Hash(),
			HeadRollup:     genesis.Hash(),
			FoundNewRollup: true,
		}
		bss.SetBlockState(b.Hash(), &bs, genesis)
		state := bss.GenesisStateDB()
		state.Commit(genesis.Hash())
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
func currentTxs(head *core.Rollup, mempool []nodecommon.L2Tx, resolver db.RollupResolver) []nodecommon.L2Tx {
	return findTxsNotIncluded(head, mempool, resolver)
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

func (e *enclaveImpl) findRoundWinner(receivedRollups []*core.Rollup, parent *core.Rollup, stateDB db.StateDB, blockResolver db.BlockResolver, rollupResolver db.RollupResolver) (*core.Rollup, db.StateDB) {
	headRollup, found := FindWinner(parent, receivedRollups, blockResolver)
	if !found {
		panic("This should not happen for gossip rounds.")
	}
	// calculate the state to compare with what is in the Rollup
	p := blockResolver.Proof(rollupResolver.ParentRollup(headRollup))
	depositTxs := processDeposits(p, blockResolver.Proof(headRollup), blockResolver, e.txHandler)

	executeTransactions(append(headRollup.Transactions, depositTxs...), stateDB, headRollup.Header)

	if stateDB.StateRoot() != headRollup.Header.State {
		panic(fmt.Sprintf("Calculated a different state. This should not happen as there are no malicious actors yet. \nGot: %s\nExp: %s\nParent state:%v\nParent state:%s\nTxs:%v",
			stateDB.StateRoot(),
			headRollup.Header.State,
			stateDB,
			parent.Header.State,
			printTxs(headRollup.Transactions)),
		)
	}
	// todo - check that the withdrawals in the header match the withdrawals as calculated

	return headRollup, stateDB
}

// returns a list of L2 deposit transactions generated from the L1 deposit transactions
// starting with the proof of the parent rollup(exclusive) to the proof of the current rollup
func processDeposits(fromBlock *types.Block, toBlock *types.Block, blockResolver db.BlockResolver, txHandler txhandler.TxHandler) []nodecommon.L2Tx {
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
			if t == nil {
				continue
			}

			if depositTx, ok := t.(*obscurocommon.L1DepositTx); ok {
				depL2TxData := core.L2TxData{
					Type:   core.DepositTx,
					To:     depositTx.To,
					Amount: depositTx.Amount,
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
func calculateBlockState(b *types.Block, parentState *core.BlockState, blockResolver db.BlockResolver, rollups []*core.Rollup, txHandler txhandler.TxHandler, rollupResolver db.RollupResolver, bss db.BlockStateStorage) (*core.BlockState, db.StateDB, *core.Rollup) {
	currentHead, found := rollupResolver.FetchRollup(parentState.HeadRollup)
	if !found {
		panic("should not happen")
	}
	newHeadRollup, found := FindWinner(currentHead, rollups, blockResolver)
	stateDB := bss.CreateStateDB(parentState.HeadRollup)
	// only change the state if there is a new l2 HeadRollup in the current block
	if found {
		// Preprocessing before passing to the vm
		// todo transform into an eth block structure
		parentRollup := rollupResolver.ParentRollup(newHeadRollup)
		p := blockResolver.Proof(parentRollup)
		depositTxs := processDeposits(p, blockResolver.Proof(newHeadRollup), blockResolver, txHandler)

		// deposits have to be processed after the normal transactions were executed because during speculative execution they are not available
		txsToProcess := append(newHeadRollup.Transactions, depositTxs...)
		executeTransactions(txsToProcess, stateDB, newHeadRollup.Header)
		// todo - handle failure , which means a new winner must be selected
	} else {
		newHeadRollup = currentHead
	}

	bs := core.BlockState{
		Block:          b.Hash(),
		HeadRollup:     newHeadRollup.Hash(),
		FoundNewRollup: found,
	}
	return &bs, stateDB, newHeadRollup
}

// Todo - this has to be implemented differently based on how we define the ObsERC20
func rollupPostProcessingWithdrawals(newHeadRollup *core.Rollup, state db.StateDB) []nodecommon.Withdrawal {
	w := make([]nodecommon.Withdrawal, 0)
	// go through each transaction and check if the withdrawal was processed correctly
	for i, t := range newHeadRollup.Transactions {
		txData := core.TxData(&newHeadRollup.Transactions[i])
		if txData.Type == core.WithdrawalTx && contains(state.Withdrawals(), t.Hash()) {
			w = append(w, nodecommon.Withdrawal{
				Amount:  txData.Amount,
				Address: txData.From,
			})
		}
	}

	return w
}

func extractRollups(b *types.Block, blockResolver db.BlockResolver, handler txhandler.TxHandler) []*core.Rollup {
	rollups := make([]*core.Rollup, 0)
	for _, tx := range b.Transactions() {
		// go through all rollup transactions
		t := handler.UnPackTx(tx)
		if t == nil {
			continue
		}

		if rolTx, ok := t.(*obscurocommon.L1RollupTx); ok {
			r := nodecommon.DecodeRollupOrPanic(rolTx.Rollup)

			// Ignore rollups created with proofs from different L1 blocks
			// In case of L1 reorgs, rollups may end published on a fork
			if blockResolver.IsBlockAncestor(b, r.Header.L1Proof) {
				rollups = append(rollups, toEnclaveRollup(r))
				log.Log(fmt.Sprintf("Extracted Rollup r_%d from block b_%d",
					obscurocommon.ShortHash(r.Hash()),
					obscurocommon.ShortHash(b.Hash()),
				))
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
