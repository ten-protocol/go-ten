package rollupchain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/crypto"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/rpcencryptionmanager"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/obscuronet/obscuro-playground/go/log"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/bridge"
	obscurocore "github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/core"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/db"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/evm"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/mempool"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
)

const (
	msgNoRollup = "could not fetch rollup"

	// The relevant fields in a Call request's params.
	CallFieldTo   = "to"
	CallFieldFrom = "from"
	CallFieldData = "data"
	DummyBalance  = "0x0"
)

// RollupChain represents the canonical chain, and manages the state.
type RollupChain struct {
	hostID          common.Address
	nodeID          uint64
	obscuroChainID  int64
	ethereumChainID int64

	storage               db.Storage
	l1Blockchain          *core.BlockChain
	bridge                *bridge.Bridge
	transactionBlobCrypto crypto.TransactionBlobCrypto // todo - remove
	rpcEncryptionManager  rpcencryptionmanager.RPCEncryptionManager
	mempool               mempool.Manager

	blockProcessingMutex sync.Mutex
}

func New(nodeID uint64, hostId common.Address, storage db.Storage, l1Blockchain *core.BlockChain, bridge *bridge.Bridge, txCrypto crypto.TransactionBlobCrypto, mempool mempool.Manager, rpcem rpcencryptionmanager.RPCEncryptionManager, obscuroChainID int64, ethereumChainID int64) *RollupChain {
	return &RollupChain{
		nodeID:                nodeID,
		hostID:                hostId,
		storage:               storage,
		l1Blockchain:          l1Blockchain,
		bridge:                bridge,
		transactionBlobCrypto: txCrypto,
		mempool:               mempool,
		rpcEncryptionManager:  rpcem,
		obscuroChainID:        obscuroChainID,
		ethereumChainID:       ethereumChainID,
		blockProcessingMutex:  sync.Mutex{},
	}
}

func (rc *RollupChain) ProduceGenesis(blkHash common.Hash) (*obscurocore.Rollup, *types.Block) {
	b, f := rc.storage.FetchBlock(blkHash)
	if !f {
		log.Panic("Could not find the block used as proof for the genesis rollup.")
	}

	rolGenesis := obscurocore.NewRollup(
		blkHash,
		nil,
		obscurocommon.L2GenesisHeight,
		common.HexToAddress("0x0"),
		[]*nodecommon.L2Tx{},
		[]nodecommon.Withdrawal{},
		obscurocommon.GenerateNonce(),
		common.BigToHash(big.NewInt(0)))
	return &rolGenesis, b
}

func (rc *RollupChain) IngestBlock(block *types.Block) nodecommon.BlockSubmissionResponse {
	// We ignore a failure on the genesis block, since insertion of the genesis also produces a failure in Geth
	// (at least with Clique, where it fails with a `vote nonce not 0x00..0 or 0xff..f`).
	if ingestionFailedResponse := rc.insertBlockIntoL1Chain(block); !rc.isGenesisBlock(block) && ingestionFailedResponse != nil {
		return *ingestionFailedResponse
	}

	rc.storage.StoreBlock(block)
	bs := rc.updateState(block)
	if bs == nil {
		return rc.noBlockStateBlockSubmissionResponse(block)
	}

	var rollup nodecommon.ExtRollup
	if bs.FoundNewRollup {
		hr, f := rc.storage.FetchRollup(bs.HeadRollup)
		if !f {
			log.Panic(msgNoRollup)
		}

		rollup = rc.transactionBlobCrypto.ToExtRollup(hr)
	}
	return rc.newBlockSubmissionResponse(bs, rollup)
}

// Inserts the block into the L1 chain if it exists and the block is not the genesis block. Returns a non-nil
// BlockSubmissionResponse if the insertion failed.
func (rc *RollupChain) insertBlockIntoL1Chain(block *types.Block) *nodecommon.BlockSubmissionResponse {
	if rc.l1Blockchain != nil {
		_, err := rc.l1Blockchain.InsertChain(types.Blocks{block})
		if err != nil {
			causeMsg := fmt.Sprintf("Block was invalid: %v", err)
			return &nodecommon.BlockSubmissionResponse{IngestedBlock: false, BlockNotIngestedCause: causeMsg}
		}
	}
	return nil
}

func (rc *RollupChain) noBlockStateBlockSubmissionResponse(block *types.Block) nodecommon.BlockSubmissionResponse {
	return nodecommon.BlockSubmissionResponse{
		BlockHeader:   block.Header(),
		IngestedBlock: true,
		FoundNewHead:  false,
	}
}

func (rc *RollupChain) newBlockSubmissionResponse(bs *obscurocore.BlockState, rollup nodecommon.ExtRollup) nodecommon.BlockSubmissionResponse {
	headRollup, f := rc.storage.FetchRollup(bs.HeadRollup)
	if !f {
		log.Panic(msgNoRollup)
	}

	headBlock, f := rc.storage.FetchBlock(bs.Block)
	if !f {
		log.Panic("could not fetch block")
	}

	var head *nodecommon.Header
	if bs.FoundNewRollup {
		head = headRollup.Header
	}
	return nodecommon.BlockSubmissionResponse{
		BlockHeader:    headBlock.Header(),
		ProducedRollup: rollup,
		IngestedBlock:  true,
		FoundNewHead:   bs.FoundNewRollup,
		RollupHead:     head,
	}
}

func (rc *RollupChain) isGenesisBlock(block *types.Block) bool {
	return rc.l1Blockchain != nil && bytes.Equal(block.Hash().Bytes(), rc.l1Blockchain.Genesis().Hash().Bytes())
}

//  STATE

// Determine the new canonical L2 head and calculate the State
func (rc *RollupChain) updateState(b *types.Block) *obscurocore.BlockState {
	// This method is called recursively in case of Re-orgs. Stop when state was calculated already.
	val, found := rc.storage.FetchBlockState(b.Hash())
	if found {
		return val
	}

	rollups := rc.bridge.ExtractRollups(b, rc.storage)
	genesisRollup := rc.storage.FetchGenesisRollup()

	// processing blocks before genesis, so there is nothing to do
	if genesisRollup == nil && len(rollups) == 0 {
		return nil
	}

	// Detect if the incoming block contains the genesis rollup, and generate an updated state.
	// Handle the case of the block containing the genesis being processed multiple times.
	genesisState, isGenesis := rc.handleGenesisRollup(b, rollups, genesisRollup)
	if isGenesis {
		return genesisState
	}

	// To calculate the state after the current block, we need the state after the parent.
	// If this point is reached, there is a parent state guaranteed, because the genesis is handled above
	parentState, parentFound := rc.storage.FetchBlockState(b.ParentHash())
	if !parentFound {
		// go back and calculate the Root of the Parent
		p, f := rc.storage.FetchBlock(b.ParentHash())
		if !f {
			nodecommon.LogWithID(rc.nodeID, "Could not find block parent. This should not happen.")
			return nil
		}
		parentState = rc.updateState(p)
	}

	if parentState == nil {
		nodecommon.LogWithID(rc.nodeID, "Something went wrong. There should be a parent here, blockNum=%d. \n Block: %d, Block Parent: %d ",
			b.Number(),
			obscurocommon.ShortHash(b.Hash()),
			obscurocommon.ShortHash(b.Header().ParentHash),
		)
		return nil
	}

	bs, stateDB, head, receipts := rc.calculateBlockState(b, parentState, rollups)
	log.Trace(fmt.Sprintf(">   Agg%d: Calc block state b_%d: Found: %t - r_%d, ",
		rc.nodeID,
		obscurocommon.ShortHash(b.Hash()),
		bs.FoundNewRollup,
		obscurocommon.ShortHash(bs.HeadRollup)))

	if bs.FoundNewRollup {
		// todo - root
		_, err := stateDB.Commit(true)
		if err != nil {
			log.Panic("could not commit new rollup to state DB. Cause: %s", err)
		}
	}
	rc.storage.SaveNewHead(bs, head, receipts)

	return bs
}

func (rc *RollupChain) handleGenesisRollup(b *types.Block, rollups []*obscurocore.Rollup, genesisRollup *obscurocore.Rollup) (genesisState *obscurocore.BlockState, isGenesis bool) {
	// the incoming block holds the genesis rollup
	// calculate and return the new block state
	// todo change this to an hardcoded hash on testnet/mainnet
	if genesisRollup == nil && len(rollups) == 1 {
		nodecommon.LogWithID(rc.nodeID, "Found genesis rollup")

		genesis := rollups[0]
		rc.storage.StoreGenesisRollup(genesis)

		// The genesis rollup is part of the canonical chain and will be included in an L1 block by the first Aggregator.
		bs := obscurocore.BlockState{
			Block:          b.Hash(),
			HeadRollup:     genesis.Hash(),
			FoundNewRollup: true,
		}
		rc.storage.SaveNewHead(&bs, genesis, nil)
		s := rc.storage.GenesisStateDB()
		_, err := s.Commit(true)
		if err != nil {
			return nil, false
		}
		return &bs, true
	}

	// Re-processing the block that contains the rollup. This can happen as blocks can be fed to the enclave multiple times.
	// In this case we don't update the state and move on.
	if genesisRollup != nil && len(rollups) == 1 && bytes.Equal(rollups[0].Header.Hash().Bytes(), genesisRollup.Hash().Bytes()) {
		return nil, true
	}
	return nil, false
}

func (rc *RollupChain) findRoundWinner(receivedRollups []*obscurocore.Rollup, parent *obscurocore.Rollup) *obscurocore.Rollup {
	headRollup, found := FindWinner(parent, receivedRollups, rc.storage)
	if !found {
		log.Panic("could not find winner. This should not happen for gossip rounds")
	}
	rc.checkRollup(headRollup)
	return headRollup
}

func (rc *RollupChain) process(rollup *obscurocore.Rollup, stateDB *state.StateDB, depositTxs []*nodecommon.L2Tx) (common.Hash, []*types.Receipt, []*types.Receipt) {
	txReceipts := evm.ExecuteTransactions(rollup.Transactions, stateDB, rollup.Header, rc.storage, rc.obscuroChainID, 0)
	if len(rollup.Transactions) != len(txReceipts) {
		panic("Sanity check. All transactions that are included in a rollup must be executed and produce a receipt.")
	}
	depositReceipts := evm.ExecuteTransactions(depositTxs, stateDB, rollup.Header, rc.storage, rc.obscuroChainID, len(rollup.Transactions))
	rootHash, err := stateDB.Commit(true)
	if err != nil {
		log.Panic("could not commit to state DB. Cause: %s", err)
	}
	return rootHash, txReceipts, depositReceipts
}

func (rc *RollupChain) validateRollup(rollup *obscurocore.Rollup, rootHash common.Hash, txReceipts []*types.Receipt, depositReceipts []*types.Receipt, stateDB *state.StateDB) bool {
	h := rollup.Header
	if !bytes.Equal(rootHash.Bytes(), h.Root.Bytes()) {
		// dump := stateDB.Dump(&state.DumpConfig{})
		dump := ""
		log.Info("Calculated a different state. This should not happen as there are no malicious actors yet. Rollup: r_%d, \nGot: %s\nExp: %s\nHeight:%d\nTxs:%v\nState: %s.\nDeposits: %+v",
			obscurocommon.ShortHash(rollup.Hash()), rootHash, h.Root, h.Number, obscurocore.PrintTxs(rollup.Transactions), dump, depositReceipts)
		return false
	}

	//  check that the withdrawals in the header match the withdrawals as calculated
	withdrawals := rc.bridge.RollupPostProcessingWithdrawals(rollup, stateDB, toReceiptMap(txReceipts))
	for i, w := range withdrawals {
		hw := h.Withdrawals[i]
		if hw.Amount != w.Amount || hw.Recipient != w.Recipient || hw.Contract != w.Contract {
			log.Panic("Withdrawals don't match")
		}
	}

	rec := getReceipts(txReceipts, depositReceipts)
	rbloom := types.CreateBloom(rec)
	if !bytes.Equal(rbloom.Bytes(), h.Bloom.Bytes()) {
		log.Info("invalid bloom (remote: %x  local: %x)", h.Bloom, rbloom)
		return false
	}

	receiptSha := types.DeriveSha(rec, trie.NewStackTrie(nil))
	if !bytes.Equal(receiptSha.Bytes(), h.ReceiptHash.Bytes()) {
		log.Info("invalid receipt root hash (remote: %x local: %x)", h.ReceiptHash, receiptSha)
		return false
	}

	return true
}

// given an L1 block, and the State as it was in the Parent block, calculates the State after the current block.
func (rc *RollupChain) calculateBlockState(b *types.Block, parentState *obscurocore.BlockState, rollups []*obscurocore.Rollup) (*obscurocore.BlockState, *state.StateDB, *obscurocore.Rollup, []*types.Receipt) {
	currentHead, found := rc.storage.FetchRollup(parentState.HeadRollup)
	if !found {
		log.Panic("could not fetch parent rollup")
	}
	newHeadRollup, found := FindWinner(currentHead, rollups, rc.storage)
	stateDB := rc.storage.CreateStateDB(parentState.HeadRollup)
	var rollupTxReceipts []*types.Receipt
	// only change the state if there is a new l2 HeadRollup in the current block
	if found {
		// Preprocessing before passing to the vm
		// todo transform into an eth block structure
		parentRollup := rc.storage.ParentRollup(newHeadRollup)
		depositTxs := rc.bridge.ExtractDeposits(rc.storage.Proof(parentRollup), rc.storage.Proof(newHeadRollup), rc.storage, stateDB)

		// deposits have to be processed after the normal transactions were executed because during speculative execution they are not available
		txsToProcess := append(newHeadRollup.Transactions, depositTxs...)
		receipts := evm.ExecuteTransactions(txsToProcess, stateDB, newHeadRollup.Header, rc.storage, rc.obscuroChainID, 0)
		rootHash := stateDB.IntermediateRoot(true)

		// todo - this should call the "checkRollup" function, but that's difficult to achieve now without some refactoring. This will be done in the next PR.
		if !bytes.Equal(newHeadRollup.Header.Root.Bytes(), rootHash.Bytes()) {
			// dump := stateDB.Dump(&state.DumpConfig{})
			dump := ""
			log.Error("Calculated a different state. This should not happen as there are no malicious actors yet. Rollup: r_%d, \nGot: %s\nExp: %s\nHeight:%d\nTxs:%v\nState: %s",
				obscurocommon.ShortHash(newHeadRollup.Hash()), rootHash, newHeadRollup.Header.Root, newHeadRollup.Header.Number, obscurocore.PrintTxs(newHeadRollup.Transactions), dump)
		}

		// todo - handle failure , which means a new winner must be selected

		// We only return the receipts for the rollup transactions, and not for deposits.
		for _, receipt := range receipts {
			for _, tx := range newHeadRollup.Transactions {
				if receipt.TxHash == tx.Hash() {
					rollupTxReceipts = append(rollupTxReceipts, receipt)
					break
				}
			}
		}
	} else {
		newHeadRollup = currentHead
	}

	bs := obscurocore.BlockState{
		Block:          b.Hash(),
		HeadRollup:     newHeadRollup.Hash(),
		FoundNewRollup: found,
	}
	return &bs, stateDB, newHeadRollup, rollupTxReceipts
}

// verifies that the headers of the rollup match the results of executing the transactions
func (rc *RollupChain) checkRollup(r *obscurocore.Rollup) {
	stateDB := rc.storage.CreateStateDB(r.Header.ParentHash)
	// calculate the state to compare with what is in the Rollup
	depositTxs := rc.bridge.ExtractDeposits(
		rc.storage.Proof(rc.storage.ParentRollup(r)),
		rc.storage.Proof(r),
		rc.storage,
		stateDB,
	)

	rootHash, txReceipts, depositReceipts := rc.process(r, stateDB, depositTxs)
	// dump := stateDB.Dump(&state.DumpConfig{})
	dump := ""

	log.Info("State rollup: r_%d. State: %s", obscurocommon.ShortHash(r.Hash()), dump)
	isValid := rc.validateRollup(r, rootHash, txReceipts, depositReceipts, stateDB)
	if !isValid {
		log.Error("Should only happen once we start including malicious actors. Until then, an invalid rollup means there is a bug.")
	}
}

func toReceiptMap(txReceipts []*types.Receipt) map[common.Hash]*types.Receipt {
	result := make(map[common.Hash]*types.Receipt, 0)
	for _, r := range txReceipts {
		result[r.TxHash] = r
	}
	return result
}

func getReceipts(txReceipts []*types.Receipt, depositReceipts []*types.Receipt) types.Receipts {
	receipts := make([]*types.Receipt, 0)
	receipts = append(receipts, txReceipts...)
	receipts = append(receipts, depositReceipts...)
	return receipts
}

// SubmitBlock is used to update the enclave with an additional L1 block.
func (rc *RollupChain) SubmitBlock(block types.Block) nodecommon.BlockSubmissionResponse {
	rc.blockProcessingMutex.Lock()
	defer rc.blockProcessingMutex.Unlock()

	// The genesis block should always be ingested, not submitted, so we ignore it if it's passed in here.
	if rc.isGenesisBlock(&block) {
		return nodecommon.BlockSubmissionResponse{IngestedBlock: false, BlockNotIngestedCause: "Block was genesis block."}
	}

	_, foundBlock := rc.storage.FetchBlock(block.Hash())
	if foundBlock {
		return nodecommon.BlockSubmissionResponse{IngestedBlock: false, BlockNotIngestedCause: "Block already ingested."}
	}

	if ingestionFailedResponse := rc.insertBlockIntoL1Chain(&block); ingestionFailedResponse != nil {
		return *ingestionFailedResponse
	}

	_, f := rc.storage.FetchBlock(block.Header().ParentHash)
	if !f && block.NumberU64() > obscurocommon.L1GenesisHeight {
		return nodecommon.BlockSubmissionResponse{IngestedBlock: false, BlockNotIngestedCause: "Block parent not stored."}
	}

	// Only store the block if the parent is available.
	stored := rc.storage.StoreBlock(&block)
	if !stored {
		return nodecommon.BlockSubmissionResponse{IngestedBlock: false}
	}

	nodecommon.LogWithID(rc.nodeID, "Update state: %d", obscurocommon.ShortHash(block.Hash()))
	blockState := rc.updateState(&block)
	if blockState == nil {
		return rc.noBlockStateBlockSubmissionResponse(&block)
	}

	// todo - A verifier node will not produce rollups, we can check the e.mining to get the node behaviour
	r := rc.produceRollup(&block, blockState)

	rc.checkRollup(r)

	// todo - should store proposal rollups in a different storage as they are ephemeral (round based)
	rc.storage.StoreRollup(r)

	nodecommon.LogWithID(rc.nodeID, "Processed block: b_%d(%d)", obscurocommon.ShortHash(block.Hash()), block.NumberU64())

	return rc.newBlockSubmissionResponse(blockState, rc.transactionBlobCrypto.ToExtRollup(r))
}

func (rc *RollupChain) produceRollup(b *types.Block, bs *obscurocore.BlockState) *obscurocore.Rollup {
	headRollup, f := rc.storage.FetchRollup(bs.HeadRollup)
	if !f {
		log.Panic(msgNoRollup)
	}

	// These variables will be used to create the new rollup
	var newRollupTxs obscurocore.L2Txs
	var newRollupState *state.StateDB
	var newRollupHeader *nodecommon.Header

	/*
			speculativeExecutionSucceeded := false
		   todo - reenable
			if e.speculativeExecutionEnabled {
				// retrieve the speculatively calculated State based on the previous winner and the incoming transactions
				e.speculativeWorkInCh <- true
				speculativeRollup := <-e.speculativeWorkOutCh

				newRollupTxs = speculativeRollup.txs
				newRollupState = speculativeRollup.s
				newRollupHeader = speculativeRollup.h

				// the speculative execution has been processing on top of the wrong parent - due to failure in gossip or publishing to L1
				// or speculative execution is disabled
				speculativeExecutionSucceeded = speculativeRollup.found && (speculativeRollup.r.Hash() == bs.HeadRollup)

				if !speculativeExecutionSucceeded && speculativeRollup.r != nil {
					nodecommon.LogWithID(e.nodeShortID, "Recalculate. speculative=r_%d(%d), published=r_%d(%d)",
						obscurocommon.ShortHash(speculativeRollup.r.Hash()),
						speculativeRollup.r.Header.Number,
						obscurocommon.ShortHash(bs.HeadRollup),
						headRollup.Header.Number)
					if e.statsCollector != nil {
						e.statsCollector.L2Recalc(e.nodeID)
					}
				}
			}
	*/

	successfulTransactions := make([]*nodecommon.L2Tx, 0)
	// if !speculativeExecutionSucceeded {
	// In case the speculative execution thread has not succeeded in producing a valid rollup
	// we have to create a new one from the mempool transactions
	newRollupHeader = obscurocore.NewHeader(&bs.HeadRollup, headRollup.NumberU64()+1, rc.hostID)
	newRollupTxs = rc.mempool.CurrentTxs(headRollup, rc.storage)
	newRollupState = rc.storage.CreateStateDB(bs.HeadRollup)
	txReceipts := evm.ExecuteTransactions(newRollupTxs, newRollupState, newRollupHeader, rc.storage, rc.obscuroChainID, 0)
	txReceiptsMap := toReceiptMap(txReceipts)
	// todo - only transactions that fail because of the nonce should be excluded
	for i, tx := range newRollupTxs {
		_, f := txReceiptsMap[tx.Hash()]
		if f {
			successfulTransactions = append(successfulTransactions, newRollupTxs[i])
		} else {
			log.Info(">   Agg%d: Excluding transaction %d", obscurocommon.ShortAddress(rc.hostID), obscurocommon.ShortHash(tx.Hash()))
		}
	}

	// always process deposits last, either on top of the rollup produced speculatively or the newly created rollup
	// process deposits from the fromBlock of the parent to the current block (which is the fromBlock of the new rollup)
	fromBlock := rc.storage.Proof(headRollup)
	depositTxs := rc.bridge.ExtractDeposits(fromBlock, b, rc.storage, newRollupState)
	depositReceipts := evm.ExecuteTransactions(depositTxs, newRollupState, newRollupHeader, rc.storage, rc.obscuroChainID, len(successfulTransactions))
	depositReceiptsMap := toReceiptMap(depositReceipts)
	// deposits should not fail
	for _, tx := range depositTxs {
		if depositReceiptsMap[tx.Hash()] == nil || depositReceiptsMap[tx.Hash()].Status == types.ReceiptStatusFailed {
			panic("Should not happen")
		}
	}

	// Create a new rollup based on the fromBlock of inclusion of the previous, including all new transactions
	rootHash, err := newRollupState.Commit(true)
	if err != nil {
		panic(err)
	}
	r := obscurocore.NewRollupFromHeader(newRollupHeader, b.Hash(), successfulTransactions, obscurocommon.GenerateNonce(), rootHash)

	// Postprocessing - withdrawals
	r.Header.Withdrawals = rc.bridge.RollupPostProcessingWithdrawals(&r, newRollupState, txReceiptsMap)
	receipts := getReceipts(txReceipts, depositReceipts)

	if len(receipts) == 0 {
		r.Header.ReceiptHash = types.EmptyRootHash
	} else {
		r.Header.ReceiptHash = types.DeriveSha(receipts, trie.NewStackTrie(nil))
		r.Header.Bloom = types.CreateBloom(receipts)
	}

	return &r
}

// TODO - this belongs in the protocol

func (rc *RollupChain) RoundWinner(parent obscurocommon.L2RootHash) (nodecommon.ExtRollup, bool, error) {
	head, found := rc.storage.FetchRollup(parent)
	if !found {
		return nodecommon.ExtRollup{}, false, fmt.Errorf("rollup not found: r_%s", parent)
	}

	headState := rc.storage.FetchHeadState()
	currentHeadRollup, found := rc.storage.FetchRollup(headState.HeadRollup)
	if !found {
		panic("Should not happen since the header hash and the rollup are stored in a batch.")
	}
	// Check if round.winner is being called on an old rollup
	if !bytes.Equal(currentHeadRollup.Hash().Bytes(), parent.Bytes()) {
		return nodecommon.ExtRollup{}, false, nil
	}

	nodecommon.LogWithID(rc.nodeID, "Round winner height: %d", head.Header.Number)
	rollupsReceivedFromPeers := rc.storage.FetchRollups(head.NumberU64() + 1)
	// filter out rollups with a different Parent
	var usefulRollups []*obscurocore.Rollup
	for _, rol := range rollupsReceivedFromPeers {
		p := rc.storage.ParentRollup(rol)
		if p == nil {
			nodecommon.LogWithID(rc.nodeID, "Received rollup from peer but don't have parent rollup - discarding...")
			continue
		}
		if bytes.Equal(p.Hash().Bytes(), head.Hash().Bytes()) {
			usefulRollups = append(usefulRollups, rol)
		}
	}

	// determine the winner of the round
	winnerRollup := rc.findRoundWinner(usefulRollups, head)
	//if rc.config.SpeculativeExecution {
	//	go rc.notifySpeculative(winnerRollup)
	//}

	// we are the winner
	if bytes.Equal(winnerRollup.Header.Agg.Bytes(), rc.hostID.Bytes()) {
		v := rc.storage.Proof(winnerRollup)
		w := rc.storage.ParentRollup(winnerRollup)
		nodecommon.LogWithID(rc.nodeID, "Publish rollup=r_%d(%d)[r_%d]{proof=b_%d(%d)}. Num Txs: %d. Txs: %v.  Root=%v. ",
			obscurocommon.ShortHash(winnerRollup.Hash()), winnerRollup.Header.Number,
			obscurocommon.ShortHash(w.Hash()),
			obscurocommon.ShortHash(v.Hash()),
			v.NumberU64(),
			len(winnerRollup.Transactions),
			obscurocore.PrintTxs(winnerRollup.Transactions),
			winnerRollup.Header.Root,
		)
		return rc.transactionBlobCrypto.ToExtRollup(winnerRollup), true, nil
	}
	return nodecommon.ExtRollup{}, false, nil
}

func (rc *RollupChain) ExecuteOffChainTransaction(encryptedParams nodecommon.EncryptedParamsCall) (nodecommon.EncryptedResponseCall, error) {
	paramBytes, err := rc.rpcEncryptionManager.DecryptRPCCall(encryptedParams)
	if err != nil {
		return nil, fmt.Errorf("could not decrypt params in eth_call request. Cause: %w", err)
	}

	contractAddress, from, data, err := extractCallParams(paramBytes)
	if err != nil {
		return nil, err
	}

	hs := rc.storage.FetchHeadState()
	if hs == nil {
		panic("Not initialised")
	}
	// todo - get the parent
	r, f := rc.storage.FetchRollup(hs.HeadRollup)
	if !f {
		panic("not found")
	}
	s := rc.storage.CreateStateDB(hs.HeadRollup)
	result, err := evm.ExecuteOffChainCall(from, contractAddress, data, s, r.Header, rc.storage, rc.obscuroChainID)
	if err != nil {
		return nil, err
	}
	if result.Failed() {
		log.Info("Failed to execute contract %s: %s\n", contractAddress.Hex(), result.Err)
		return nil, result.Err
	}

	encryptedResult, err := rc.rpcEncryptionManager.EncryptWithViewingKey(from, result.ReturnData)
	if err != nil {
		return nil, fmt.Errorf("enclave could not respond securely to eth_call request. Cause: %w", err)
	}

	return encryptedResult, nil
}

// Extracts and validates the relevant parameters in a Call request.
func extractCallParams(decryptedParams []byte) (common.Address, common.Address, []byte, error) {
	var paramsJSONMap []interface{}
	err := json.Unmarshal(decryptedParams, &paramsJSONMap)
	if err != nil {
		return common.Address{}, common.Address{}, nil, fmt.Errorf("could not parse JSON params in Call request. Cause: %w", err)
	}

	txArgs := paramsJSONMap[0] // The first argument is the transaction arguments, the second the block, the third the state overrides.
	contractAddressString, ok := txArgs.(map[string]interface{})[CallFieldTo].(string)
	if !ok {
		return common.Address{}, common.Address{}, nil, fmt.Errorf("to field in Call request params was not of expected type string")
	}
	fromString, ok := txArgs.(map[string]interface{})[CallFieldFrom].(string)
	if !ok {
		return common.Address{}, common.Address{}, nil, fmt.Errorf("from field in Call request params was not of expected type string")
	}
	dataString, ok := txArgs.(map[string]interface{})[CallFieldData].(string)
	if !ok {
		return common.Address{}, common.Address{}, nil, fmt.Errorf("data field in Call request params was not of expected type string")
	}

	contractAddress := common.HexToAddress(contractAddressString)
	from := common.HexToAddress(fromString)
	data, err := hexutil.Decode(dataString)
	if err != nil {
		return common.Address{}, common.Address{}, nil, fmt.Errorf("could not decode data in Call request. Cause: %w", err)
	}
	return contractAddress, from, data, nil
}

func (rc *RollupChain) GetBalance(encryptedParams nodecommon.EncryptedParamsGetBalance) (nodecommon.EncryptedResponseGetBalance, error) {
	paramBytes, err := rc.rpcEncryptionManager.DecryptRPCCall(encryptedParams)
	if err != nil {
		return nil, fmt.Errorf("could not decrypt params in eth_getBalance request. Cause: %w", err)
	}

	var paramsJSONMap []string
	err = json.Unmarshal(paramBytes, &paramsJSONMap)
	if err != nil {
		return nil, fmt.Errorf("could not parse JSON params in GetBalance request. Cause: %w", err)
	}
	address := common.HexToAddress(paramsJSONMap[0]) // The first argument is the address, the second the block.

	// TODO - Calculate balance correctly, rather than returning this dummy value.
	balance := DummyBalance // The Ethereum API is to return the balance in hex.

	encryptedBalance, err := rc.rpcEncryptionManager.EncryptWithViewingKey(address, []byte(balance))
	if err != nil {
		return nil, fmt.Errorf("enclave could not respond securely to eth_getBalance request. Cause: %w", err)
	}

	return encryptedBalance, nil
}
