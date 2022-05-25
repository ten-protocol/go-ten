package enclave

import (
	"fmt"
	"math/big"
	"sort"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/obscuro-playground/contracts"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/evm"

	"github.com/ethereum/go-ethereum/core/state"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/core"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/db"

	"github.com/obscuronet/obscuro-playground/go/ethclient/mgmtcontractlib"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/go/ethclient/erc20contractlib"
	"github.com/obscuronet/obscuro-playground/go/log"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
)

// Determine the new canonical L2 head and calculate the State
// Uses cache-ing to map the Head rollup and the State to each L1Node block.
func updateState(
	b *types.Block,
	blockResolver db.BlockResolver,
	mgmtContractLib mgmtcontractlib.MgmtContractLib,
	erc20ContractLib erc20contractlib.ERC20ContractLib,
	rollupResolver db.RollupResolver,
	bss db.BlockStateStorage,
	nodeID uint64,
) *core.BlockState {
	// This method is called recursively in case of Re-orgs. Stop when state was calculated already.
	val, found := bss.FetchBlockState(b.Hash())
	if found {
		return val
	}

	rollups := extractRollups(b, blockResolver, mgmtContractLib, nodeID)
	genesisRollup := rollupResolver.FetchGenesisRollup()

	// processing blocks before genesis, so there is nothing to do
	if genesisRollup == nil && len(rollups) == 0 {
		return nil
	}

	// Detect if the incoming block contains the genesis rollup, and generate an updated state.
	// Handle the case of the block containing the genesis being processed multiple times.
	genesisState, isGenesis := handleGenesisRollup(b, rollups, genesisRollup, rollupResolver, bss, nodeID)
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
			nodecommon.LogWithID(nodeID, "Could not find block parent. This should not happen.")
			return nil
		}
		parentState = updateState(p, blockResolver, mgmtContractLib, erc20ContractLib, rollupResolver, bss, nodeID)
	}

	if parentState == nil {
		nodecommon.LogWithID(nodeID, "Something went wrong. There should be parent here. \n Block: %d - Block Parent: %d - Header: %+v",
			obscurocommon.ShortHash(b.Hash()),
			obscurocommon.ShortHash(b.Header().ParentHash),
			b.Header(),
		)
		return nil
	}

	bs, stateDB, head := calculateBlockState(b, parentState, blockResolver, rollups, erc20ContractLib, rollupResolver, bss)
	log.Trace(fmt.Sprintf(">   Agg%d: Calc block state b_%d: Found: %t - r_%d, ",
		nodeID,
		obscurocommon.ShortHash(b.Hash()),
		bs.FoundNewRollup,
		obscurocommon.ShortHash(bs.HeadRollup),
	))

	bss.SetBlockState(b.Hash(), bs, head)
	if bs.FoundNewRollup {
		// todo - root
		_, err := stateDB.Commit(true)
		if err != nil {
			log.Panic("could not commit new rollup to state DB. Cause: %s", err)
		}
	}

	return bs
}

func handleGenesisRollup(b *types.Block, rollups []*core.Rollup, genesisRollup *core.Rollup, resolver db.RollupResolver, bss db.BlockStateStorage, nodeID uint64) (genesisState *core.BlockState, isGenesis bool) {
	// the incoming block holds the genesis rollup
	// calculate and return the new block state
	// todo change this to an hardcoded hash on testnet/mainnet
	if genesisRollup == nil && len(rollups) == 1 {
		nodecommon.LogWithID(nodeID, "Found genesis rollup")

		genesis := rollups[0]
		resolver.StoreGenesisRollup(genesis)

		// The genesis rollup is part of the canonical chain and will be included in an L1 block by the first Aggregator.
		bs := core.BlockState{
			Block:          b.Hash(),
			HeadRollup:     genesis.Hash(),
			FoundNewRollup: true,
		}
		bss.SetBlockState(b.Hash(), &bs, genesis)
		s := bss.GenesisStateDB()
		_, err := s.Commit(true)
		if err != nil {
			return nil, false
		}
		return &bs, true
	}

	// Re-processing the block that contains the rollup. This can happen as blocks can be fed to the enclave multiple times.
	// In this case we don't update the state and move on.
	if genesisRollup != nil && len(rollups) == 1 && rollups[0].Header.Hash() == genesisRollup.Hash() {
		return nil, true
	}
	return nil, false
}

type SortByNonce []nodecommon.L2Tx

func (c SortByNonce) Len() int           { return len(c) }
func (c SortByNonce) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c SortByNonce) Less(i, j int) bool { return c[i].Nonce() < c[j].Nonce() }

// Calculate transactions to be included in the current rollup
func currentTxs(head *core.Rollup, mempool []nodecommon.L2Tx, resolver db.RollupResolver) []nodecommon.L2Tx {
	txs := findTxsNotIncluded(head, mempool, resolver)
	sort.Sort(SortByNonce(txs))
	return txs
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

func (e *enclaveImpl) findRoundWinner(receivedRollups []*core.Rollup, parent *core.Rollup, stateDB *state.StateDB, blockResolver db.BlockResolver, rollupResolver db.RollupResolver) (*core.Rollup, *state.StateDB) {
	headRollup, found := FindWinner(parent, receivedRollups, blockResolver)
	if !found {
		log.Panic("could not find winner. This should not happen for gossip rounds")
	}
	// calculate the state to compare with what is in the Rollup
	p := blockResolver.Proof(rollupResolver.ParentRollup(headRollup))
	depositTxs := extractDeposits(p, blockResolver.Proof(headRollup), blockResolver, e.erc20ContractLib, stateDB)
	log.Info(fmt.Sprintf(">   Agg%d: Deposits:%d", obscurocommon.ShortAddress(e.nodeID), len(depositTxs)))

	evm.ExecuteTransactions(headRollup.Transactions, stateDB, headRollup.Header, e.storage)
	evm.ExecuteTransactions(depositTxs, stateDB, headRollup.Header, e.storage)
	rootHash, err := stateDB.Commit(true)
	if err != nil {
		log.Panic("could not commit to state DB. Cause: %s", err)
	}

	if rootHash != headRollup.Header.State {
		// dump := stateDB.Dump(&state.DumpConfig{})
		dump := ""
		log.Panic("Calculated a different state. This should not happen as there are no malicious actors yet. \nGot: %s\nExp: %s\nHeight:%d\nTxs:%v\nState: %s",
			rootHash, headRollup.Header.State, headRollup.Header.Number, printTxs(headRollup.Transactions), dump)
	}
	// todo - check that the withdrawals in the header match the withdrawals as calculated

	return headRollup, stateDB
}

// returns a list of L2 deposit transactions generated from the L1 deposit transactions
// starting with the proof of the parent rollup(exclusive) to the proof of the current rollup
func extractDeposits(
	fromBlock *types.Block,
	toBlock *types.Block,
	blockResolver db.BlockResolver,
	erc20ContractLib erc20contractlib.ERC20ContractLib,
	rollupState *state.StateDB,
) []nodecommon.L2Tx {
	from := obscurocommon.GenesisBlock.Hash()
	height := obscurocommon.L1GenesisHeight
	if fromBlock != nil {
		from = fromBlock.Hash()
		height = fromBlock.NumberU64()
		if !blockResolver.IsAncestor(toBlock, fromBlock) {
			log.Panic("Deposits can't be processed because the rollups are not on the same Ethereum fork. This should not happen.")
		}
	}

	allDeposits := make([]nodecommon.L2Tx, 0)
	b := toBlock
	zero := uint64(0)
	i := &zero
	for {
		if b.Hash() == from {
			break
		}
		for _, tx := range b.Transactions() {
			t := erc20ContractLib.DecodeTx(tx)
			if t == nil {
				continue
			}

			if depositTx, ok := t.(*obscurocommon.L1DepositTx); ok {
				depL2Tx := newDepositTx(*depositTx.Sender, depositTx.Amount, rollupState, atomic.LoadUint64(i))
				allDeposits = append(allDeposits, depL2Tx)
				atomic.AddUint64(i, 1)
			}
		}
		if b.NumberU64() < height {
			log.Panic("block height is less than genesis height")
		}
		p, f := blockResolver.ParentBlock(b)
		if !f {
			log.Panic("deposits can't be processed because the rollups are not on the same Ethereum fork")
		}
		b = p
	}

	return allDeposits
}

// given an L1 block, and the State as it was in the Parent block, calculates the State after the current block.
func calculateBlockState(
	b *types.Block,
	parentState *core.BlockState,
	blockResolver db.BlockResolver,
	rollups []*core.Rollup,
	erc20ContractLib erc20contractlib.ERC20ContractLib,
	rollupResolver db.RollupResolver,
	bss db.BlockStateStorage,
) (*core.BlockState, *state.StateDB, *core.Rollup) {
	currentHead, found := rollupResolver.FetchRollup(parentState.HeadRollup)
	if !found {
		log.Panic("could not fetch parent rollup")
	}
	newHeadRollup, found := FindWinner(currentHead, rollups, blockResolver)
	stateDB := bss.CreateStateDB(parentState.HeadRollup)
	// only change the state if there is a new l2 HeadRollup in the current block
	if found {
		// Preprocessing before passing to the vm
		// todo transform into an eth block structure
		parentRollup := rollupResolver.ParentRollup(newHeadRollup)
		p := blockResolver.Proof(parentRollup)
		depositTxs := extractDeposits(p, blockResolver.Proof(newHeadRollup), blockResolver, erc20ContractLib, stateDB)

		// deposits have to be processed after the normal transactions were executed because during speculative execution they are not available
		txsToProcess := append(newHeadRollup.Transactions, depositTxs...)
		evm.ExecuteTransactions(txsToProcess, stateDB, newHeadRollup.Header, rollupResolver)
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
func rollupPostProcessingWithdrawals(newHeadRollup *core.Rollup, state *state.StateDB, receipts map[common.Hash]*types.Receipt) []nodecommon.Withdrawal {
	w := make([]nodecommon.Withdrawal, 0)
	// go through each transaction and check if the withdrawal was processed correctly
	for i, t := range newHeadRollup.Transactions {
		method, err := contracts.PedroERC20ContractABIJSON.MethodById(t.Data()[:4])
		if err != nil {
			// fmt.Printf("%s\n", err)
			continue
		}

		if method.Name != "transfer" {
			continue
		}

		args := map[string]interface{}{}
		if err := method.Inputs.UnpackIntoMap(args, t.Data()[4:]); err != nil {
			panic(err)
		}

		if *t.To() == evm.Erc20ContractAddress && args["to"].(common.Address) == evm.WithdrawalAddress {
			receipt := receipts[t.Hash()]
			if receipt != nil && receipt.Status == 1 {
				signer := types.NewLondonSigner(obscurocommon.ChainID)
				from, err1 := types.Sender(signer, &newHeadRollup.Transactions[i])
				if err1 != nil {
					panic(err1)
				}
				state.Logs()
				w = append(w, nodecommon.Withdrawal{
					Amount:  args["amount"].(*big.Int).Uint64(),
					Address: from,
				})
			}
		}
	}

	// TODO - fix the withdrawals logic
	// clearWithdrawals(state, withdrawalTxs)
	return w
}

func extractRollups(b *types.Block, blockResolver db.BlockResolver, mgmtContractLib mgmtcontractlib.MgmtContractLib, nodeID uint64) []*core.Rollup {
	rollups := make([]*core.Rollup, 0)
	for _, tx := range b.Transactions() {
		// go through all rollup transactions
		t := mgmtContractLib.DecodeTx(tx)
		if t == nil {
			continue
		}

		if rolTx, ok := t.(*obscurocommon.L1RollupTx); ok {
			r := nodecommon.DecodeRollupOrPanic(rolTx.Rollup)

			// Ignore rollups created with proofs from different L1 blocks
			// In case of L1 reorgs, rollups may end published on a fork
			if blockResolver.IsBlockAncestor(b, r.Header.L1Proof) {
				rollups = append(rollups, toEnclaveRollup(r))
				nodecommon.LogWithID(nodeID, "Extracted Rollup r_%d from block b_%d",
					obscurocommon.ShortHash(r.Hash()),
					obscurocommon.ShortHash(b.Hash()),
				)
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

// todo - nonce
func newDepositTx(address common.Address, amount uint64, rollupState *state.StateDB, i uint64) nodecommon.L2Tx {
	transferERC20data, err := contracts.PedroERC20ContractABIJSON.Pack("transfer", address, big.NewInt(int64(amount)))
	if err != nil {
		panic(err)
	}
	signer := types.NewLondonSigner(obscurocommon.ChainID)

	nonce := rollupState.GetNonce(evm.Erc20OwnerAddress) + i

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		Value:    common.Big0,
		Gas:      1_000_000,
		GasPrice: common.Big0,
		Data:     transferERC20data,
		To:       &evm.Erc20ContractAddress,
	})

	newTx, err := types.SignTx(tx, signer, evm.Erc20OwnerKey)
	if err != nil {
		log.Panic("could not encode L2 transaction data. Cause: %s", err)
	}
	return *newTx
}
