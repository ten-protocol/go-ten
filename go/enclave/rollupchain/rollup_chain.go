package rollupchain

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/obscuronet/go-obscuro/go/enclave/events"

	"github.com/google/uuid"

	gethrpc "github.com/ethereum/go-ethereum/rpc"

	"github.com/obscuronet/go-obscuro/go/enclave/rpc"

	"github.com/ethereum/go-ethereum/params"

	"github.com/status-im/keycard-go/hexutils"

	"github.com/obscuronet/go-obscuro/go/common/log"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/enclave/bridge"
	obscurocore "github.com/obscuronet/go-obscuro/go/enclave/core"
	"github.com/obscuronet/go-obscuro/go/enclave/crypto"
	"github.com/obscuronet/go-obscuro/go/enclave/db"
	"github.com/obscuronet/go-obscuro/go/enclave/evm"
	"github.com/obscuronet/go-obscuro/go/enclave/mempool"
)

const (
	msgNoRollup = "could not fetch rollup"
)

// RollupChain represents the canonical chain, and manages the state.
type RollupChain struct {
	hostID          gethcommon.Address
	nodeID          uint64
	ethereumChainID int64
	chainConfig     *params.ChainConfig

	storage               db.Storage
	l1Blockchain          *core.BlockChain
	bridge                *bridge.Bridge
	transactionBlobCrypto crypto.TransactionBlobCrypto // todo - remove
	rpcEncryptionManager  rpc.EncryptionManager
	mempool               mempool.Manager
	faucet                Faucet
	subscriptionManager   *events.SubscriptionManager

	enclavePrivateKey    *ecdsa.PrivateKey // this is a key known only to the current enclave, and the public key was shared with everyone during attestation
	blockProcessingMutex sync.Mutex
}

func New(nodeID uint64, hostID gethcommon.Address, storage db.Storage, l1Blockchain *core.BlockChain, bridge *bridge.Bridge, subscriptionManager *events.SubscriptionManager, txCrypto crypto.TransactionBlobCrypto, mempool mempool.Manager, rpcem rpc.EncryptionManager, privateKey *ecdsa.PrivateKey, ethereumChainID int64, chainConfig *params.ChainConfig) *RollupChain {
	return &RollupChain{
		nodeID:                nodeID,
		hostID:                hostID,
		storage:               storage,
		l1Blockchain:          l1Blockchain,
		bridge:                bridge,
		transactionBlobCrypto: txCrypto,
		mempool:               mempool,
		faucet:                NewFaucet(storage),
		subscriptionManager:   subscriptionManager,
		enclavePrivateKey:     privateKey,
		rpcEncryptionManager:  rpcem,
		ethereumChainID:       ethereumChainID,
		chainConfig:           chainConfig,
		blockProcessingMutex:  sync.Mutex{},
	}
}

func (rc *RollupChain) ProduceGenesis(blkHash gethcommon.Hash) (*obscurocore.Rollup, *types.Block) {
	b, f := rc.storage.FetchBlock(blkHash)
	if !f {
		log.Panic("Could not find the block used as proof for the genesis rollup.")
	}

	rolGenesis := obscurocore.NewRollup(
		blkHash,
		nil,
		common.L2GenesisHeight,
		gethcommon.HexToAddress("0x0"),
		[]*common.L2Tx{},
		[]common.Withdrawal{},
		common.GenerateNonce(),
		rc.faucet.GetGenesisRoot(rc.storage),
	)
	rc.signRollup(rolGenesis)

	return rolGenesis, b
}

func (rc *RollupChain) IngestBlock(block *types.Block) common.BlockSubmissionResponse {
	// We ignore a failure on the genesis block, since insertion of the genesis also produces a failure in Geth
	// (at least with Clique, where it fails with a `vote nonce not 0x00..0 or 0xff..f`).
	if ingestionFailedResponse := rc.insertBlockIntoL1Chain(block); !rc.isGenesisBlock(block) && ingestionFailedResponse != nil {
		return *ingestionFailedResponse
	}

	rc.storage.StoreBlock(block)
	bs, subscribedLogs := rc.updateState(block)
	if bs == nil {
		return rc.noBlockStateBlockSubmissionResponse(block)
	}

	var rollup common.ExtRollup
	if bs.FoundNewRollup {
		hr, f := rc.storage.FetchRollup(bs.HeadRollup)
		if !f {
			log.Panic(msgNoRollup)
		}

		rollup = hr.ToExtRollup(rc.transactionBlobCrypto)
	}

	encryptedLogs, err := rc.subscriptionManager.EncryptLogs(subscribedLogs)
	if err != nil {
		log.Panic("Could not encrypt logs. Cause: %s", err)
	}

	return rc.newBlockSubmissionResponse(bs, rollup, encryptedLogs)
}

// Inserts the block into the L1 chain if it exists and the block is not the genesis block. Returns a non-nil
// BlockSubmissionResponse if the insertion failed.
func (rc *RollupChain) insertBlockIntoL1Chain(block *types.Block) *common.BlockSubmissionResponse {
	if rc.l1Blockchain != nil {
		_, err := rc.l1Blockchain.InsertChain(types.Blocks{block})
		if err != nil {
			causeMsg := fmt.Sprintf("Block was invalid: %v", err)
			return &common.BlockSubmissionResponse{IngestedBlock: false, BlockNotIngestedCause: causeMsg}
		}
	}
	return nil
}

func (rc *RollupChain) noBlockStateBlockSubmissionResponse(block *types.Block) common.BlockSubmissionResponse {
	return common.BlockSubmissionResponse{
		BlockHeader:   block.Header(),
		IngestedBlock: true,
		FoundNewHead:  false,
	}
}

func (rc *RollupChain) newBlockSubmissionResponse(bs *obscurocore.BlockState, rollup common.ExtRollup, logs map[uuid.UUID]common.EncryptedLogs) common.BlockSubmissionResponse {
	headRollup, f := rc.storage.FetchRollup(bs.HeadRollup)
	if !f {
		log.Panic(msgNoRollup)
	}

	headBlock, f := rc.storage.FetchBlock(bs.Block)
	if !f {
		log.Panic("could not fetch block")
	}

	var head *common.Header
	if bs.FoundNewRollup {
		head = headRollup.Header
	}
	return common.BlockSubmissionResponse{
		BlockHeader:    headBlock.Header(),
		ProducedRollup: rollup,
		IngestedBlock:  true,
		FoundNewHead:   bs.FoundNewRollup,
		RollupHead:     head,
		SubscribedLogs: logs,
	}
}

func (rc *RollupChain) isGenesisBlock(block *types.Block) bool {
	return rc.l1Blockchain != nil && bytes.Equal(block.Hash().Bytes(), rc.l1Blockchain.Genesis().Hash().Bytes())
}

//  STATE

// Recursively calculates the state and logs for the given block.
func (rc *RollupChain) updateState(b *types.Block) (*obscurocore.BlockState, map[uuid.UUID][]*types.Log) {
	// This method is called recursively in case of Re-orgs. Stop when state was calculated already.
	blockState, _, found := rc.storage.FetchBlockState(b.Hash())
	if found {
		return blockState, nil
	}

	rollups := rc.bridge.ExtractRollups(b, rc.storage)
	genesisRollup := rc.storage.FetchGenesisRollup()

	// processing blocks before genesis, so there is nothing to do
	if genesisRollup == nil && len(rollups) == 0 {
		return nil, nil
	}

	// Detect if the incoming block contains the genesis rollup, and generate an updated state.
	// Handles the case of the block containing the genesis being processed multiple times.
	genesisState, isGenesis := rc.handleGenesisRollup(b, rollups, genesisRollup)
	if isGenesis {
		return genesisState, nil
	}

	// To calculate the state after the current block, we need the state after the parent.
	// If this point is reached, there is a parent state guaranteed, because the genesis is handled above
	parentState, parentLogs, parentFound := rc.storage.FetchBlockState(b.ParentHash())
	if !parentFound {
		// go back and calculate the Root of the Parent
		parent, found := rc.storage.FetchBlock(b.ParentHash())
		if !found {
			common.LogWithID(rc.nodeID, "Could not find block parent. This should not happen.")
			return nil, nil
		}
		parentState, parentLogs = rc.updateState(parent)
	}

	if parentState == nil {
		common.LogWithID(rc.nodeID, "Something went wrong. There should be a parent here, blockNum=%d. \n Block: %d, Block Parent: %d ",
			b.Number(),
			common.ShortHash(b.Hash()),
			common.ShortHash(b.Header().ParentHash),
		)
		return nil, nil
	}

	bs, head, receipts := rc.calculateBlockState(b, parentState, rollups)
	common.TraceWithID(rc.nodeID, "Calc block state b_%d: Found: %t - r_%d, ",
		common.ShortHash(b.Hash()),
		bs.FoundNewRollup,
		common.ShortHash(bs.HeadRollup))

	var logs []*types.Log
	for _, receipt := range receipts {
		logs = append(logs, receipt.Logs...)
	}
	// todo - joel - the creation of the state DB doesn't seem right
	subscribedLogs := rc.subscriptionManager.FilterRelevantLogs(logs, rc.storage.CreateStateDB(head.Header.ParentHash))

	// TODO - #453 - Check this recursive logic works correctly (i.e. each block submission response contains the logs
	//  of all its ancestors as well).
	// TODO - #453 - Check whether this recursive logic is desirable. Won't it balloon over time?
	// TODO - #453 - Should we be storing the parent logs based on the subscriptions at the time, or recalculating
	//  based on the current subscriptions?

	// We append the rollup's logs to the logs of the parent rollup. This is to ensure events are not missed if a
	// block is missed.
	for subscriptionID, logs := range subscribedLogs {
		logsForID, found := parentLogs[subscriptionID]
		if !found {
			logsForID = make([]*types.Log, 0)
		}
		subscribedLogs[subscriptionID] = append(logsForID, logs...)
	}

	rc.storage.SaveNewHead(bs, head, receipts, subscribedLogs)

	return bs, subscribedLogs
}

func (rc *RollupChain) handleGenesisRollup(b *types.Block, rollups []*obscurocore.Rollup, genesisRollup *obscurocore.Rollup) (genesisState *obscurocore.BlockState, isGenesis bool) {
	// the incoming block holds the genesis rollup
	// calculate and return the new block state
	// todo change this to an hardcoded hash on testnet/mainnet
	if genesisRollup == nil && len(rollups) == 1 {
		common.LogWithID(rc.nodeID, "Found genesis rollup")

		genesis := rollups[0]
		rc.storage.StoreGenesisRollup(genesis)

		// The genesis rollup is part of the canonical chain and will be included in an L1 block by the first Aggregator.
		bs := obscurocore.BlockState{
			Block:          b.Hash(),
			HeadRollup:     genesis.Hash(),
			FoundNewRollup: true,
		}
		rc.storage.SaveNewHead(&bs, genesis, nil, nil)
		err := rc.faucet.CalculateGenesisState(rc.storage)
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

type sortByTxIndex []*types.Receipt

func (c sortByTxIndex) Len() int           { return len(c) }
func (c sortByTxIndex) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c sortByTxIndex) Less(i, j int) bool { return c[i].TransactionIndex < c[j].TransactionIndex }

// This is where transactions are executed and the state is calculated.
// Obscuro includes a bridge embedded in the platform, and this method is responsible for processing deposits as well.
// The rollup can be a final rollup as received from peers or the rollup under construction.
func (rc *RollupChain) processState(rollup *obscurocore.Rollup, txs []*common.L2Tx, stateDB *state.StateDB) (gethcommon.Hash, []*common.L2Tx, []*types.Receipt, []*types.Receipt) {
	var executedTransactions []*common.L2Tx
	var txReceipts []*types.Receipt

	txResults := evm.ExecuteTransactions(txs, stateDB, rollup.Header, rc.storage, rc.chainConfig, 0)
	for _, tx := range txs {
		result, f := txResults[tx.Hash()]
		if !f {
			log.Panic("There should be an entry for each transaction ")
		}
		rec, foundReceipt := result.(*types.Receipt)
		if foundReceipt {
			executedTransactions = append(executedTransactions, tx)
			txReceipts = append(txReceipts, rec)
		} else {
			// Exclude all errors
			common.LogWithID(common.ShortAddress(rc.hostID), "Excluding transaction %s from rollup r_%d. Cause: %s", tx.Hash().Hex(), common.ShortHash(rollup.Hash()), result)
		}
	}

	// always process deposits last, either on top of the rollup produced speculatively or the newly created rollup
	// process deposits from the fromBlock of the parent to the current block (which is the fromBlock of the new rollup)
	depositTxs := rc.bridge.ExtractDeposits(
		rc.storage.Proof(rc.storage.ParentRollup(rollup)),
		rc.storage.Proof(rollup),
		rc.storage,
		stateDB,
	)

	depositResponses := evm.ExecuteTransactions(depositTxs, stateDB, rollup.Header, rc.storage, rc.chainConfig, len(executedTransactions))
	depositReceipts := make([]*types.Receipt, len(depositResponses))
	if len(depositResponses) != len(depositTxs) {
		log.Panic("Sanity check. Some deposit transactions failed.")
	}
	i := 0
	for _, resp := range depositResponses {
		rec, ok := resp.(*types.Receipt)
		if !ok {
			// TODO - Handle the case of an error (e.g. insufficient funds).
			log.Panic("Sanity check. Expected a receipt, got %v", resp)
		}
		depositReceipts[i] = rec
		i++
	}

	rootHash, err := stateDB.Commit(true)
	if err != nil {
		log.Panic("could not commit to state DB. Cause: %s", err)
	}

	sort.Sort(sortByTxIndex(txReceipts))
	sort.Sort(sortByTxIndex(depositReceipts))

	// todo - handle the tx execution logs
	return rootHash, executedTransactions, txReceipts, depositReceipts
}

func (rc *RollupChain) validateRollup(rollup *obscurocore.Rollup, rootHash gethcommon.Hash, txReceipts []*types.Receipt, depositReceipts []*types.Receipt, stateDB *state.StateDB) bool {
	h := rollup.Header
	if !bytes.Equal(rootHash.Bytes(), h.Root.Bytes()) {
		dump := strings.Replace(string(stateDB.Dump(&state.DumpConfig{})), "\n", "", -1)
		log.Error("Verify rollup r_%d: Calculated a different state. This should not happen as there are no malicious actors yet. \nGot: %s\nExp: %s\nHeight:%d\nTxs:%v\nState: %s.\nDeposits: %+v",
			common.ShortHash(rollup.Hash()), rootHash, h.Root, h.Number, obscurocore.PrintTxs(rollup.Transactions), dump, depositReceipts)
		return false
	}

	//  check that the withdrawals in the header match the withdrawals as calculated
	withdrawals := rc.bridge.RollupPostProcessingWithdrawals(rollup, stateDB, toReceiptMap(txReceipts))
	for i, w := range withdrawals {
		hw := h.Withdrawals[i]
		if hw.Amount.Cmp(w.Amount) != 0 || hw.Recipient != w.Recipient || hw.Contract != w.Contract {
			log.Error("Verify rollup r_%d: Withdrawals don't match", common.ShortHash(rollup.Hash()))
			return false
		}
	}

	rec := allReceipts(txReceipts, depositReceipts)
	rbloom := types.CreateBloom(rec)
	if !bytes.Equal(rbloom.Bytes(), h.Bloom.Bytes()) {
		log.Error("Verify rollup r_%d: Invalid bloom (remote: %x  local: %x)", common.ShortHash(rollup.Hash()), h.Bloom, rbloom)
		return false
	}

	receiptSha := types.DeriveSha(rec, trie.NewStackTrie(nil))
	if !bytes.Equal(receiptSha.Bytes(), h.ReceiptHash.Bytes()) {
		log.Error("Verify rollup r_%d: invalid receipt root hash (remote: %x local: %x)", common.ShortHash(rollup.Hash()), h.ReceiptHash, receiptSha)
		return false
	}

	return true
}

// given an L1 block, and the State as it was in the Parent block, calculates the State after the current block.
func (rc *RollupChain) calculateBlockState(b *types.Block, parentState *obscurocore.BlockState, rollups []*obscurocore.Rollup) (*obscurocore.BlockState, *obscurocore.Rollup, []*types.Receipt) {
	currentHead, found := rc.storage.FetchRollup(parentState.HeadRollup)
	if !found {
		log.Panic("could not fetch parent rollup")
	}
	newHeadRollup, found := FindWinner(currentHead, rollups, rc.storage)
	var rollupTxReceipts []*types.Receipt
	// only change the state if there is a new l2 HeadRollup in the current block
	if found {
		rollupTxReceipts, _ = rc.checkRollup(newHeadRollup)
	} else {
		newHeadRollup = currentHead
	}

	bs := obscurocore.BlockState{
		Block:          b.Hash(),
		HeadRollup:     newHeadRollup.Hash(),
		FoundNewRollup: found,
	}
	return &bs, newHeadRollup, rollupTxReceipts
}

// verifies that the headers of the rollup match the results of executing the transactions
func (rc *RollupChain) checkRollup(r *obscurocore.Rollup) ([]*types.Receipt, []*types.Receipt) {
	stateDB := rc.storage.CreateStateDB(r.Header.ParentHash)
	// calculate the state to compare with what is in the Rollup
	rootHash, successfulTxs, txReceipts, depositReceipts := rc.processState(r, r.Transactions, stateDB)
	if len(successfulTxs) != len(r.Transactions) {
		panic("Sanity check. All transactions that are included in a rollup must be executed.")
	}

	isValid := rc.validateRollup(r, rootHash, txReceipts, depositReceipts, stateDB)
	if !isValid {
		log.Panic("Should only happen once we start including malicious actors. Until then, an invalid rollup means there is a bug.")
	}

	// todo - check that the transactions hash to the header.txHash

	// verify the signature
	isValid = rc.verifySig(r)
	if !isValid {
		log.Panic("Should only happen once we start including malicious actors. Until then, a rollup with an invalid signature is a bug.")
	}

	return txReceipts, depositReceipts
}

func toReceiptMap(txReceipts []*types.Receipt) map[gethcommon.Hash]*types.Receipt {
	result := make(map[gethcommon.Hash]*types.Receipt, 0)
	for _, r := range txReceipts {
		result[r.TxHash] = r
	}
	return result
}

func allReceipts(txReceipts []*types.Receipt, depositReceipts []*types.Receipt) types.Receipts {
	receipts := make([]*types.Receipt, 0)
	receipts = append(receipts, txReceipts...)
	receipts = append(receipts, depositReceipts...)
	return receipts
}

// SubmitBlock is used to update the enclave with an additional L1 block.
func (rc *RollupChain) SubmitBlock(block types.Block) common.BlockSubmissionResponse {
	rc.blockProcessingMutex.Lock()
	defer rc.blockProcessingMutex.Unlock()

	// The genesis block should always be ingested, not submitted, so we ignore it if it's passed in here.
	if rc.isGenesisBlock(&block) {
		return common.BlockSubmissionResponse{IngestedBlock: false, BlockNotIngestedCause: "Block was genesis block."}
	}

	_, foundBlock := rc.storage.FetchBlock(block.Hash())
	if foundBlock {
		return common.BlockSubmissionResponse{IngestedBlock: false, BlockNotIngestedCause: "Block already ingested."}
	}

	if ingestionFailedResponse := rc.insertBlockIntoL1Chain(&block); ingestionFailedResponse != nil {
		return *ingestionFailedResponse
	}

	_, f := rc.storage.FetchBlock(block.Header().ParentHash)
	if !f && block.NumberU64() > common.L1GenesisHeight {
		return common.BlockSubmissionResponse{IngestedBlock: false, BlockNotIngestedCause: "Block parent not stored."}
	}

	// Only store the block if the parent is available.
	stored := rc.storage.StoreBlock(&block)
	if !stored {
		return common.BlockSubmissionResponse{IngestedBlock: false}
	}

	common.TraceWithID(rc.nodeID, "Update state: b_%d", common.ShortHash(block.Hash()))
	blockState, logs := rc.updateState(&block)
	if blockState == nil {
		return rc.noBlockStateBlockSubmissionResponse(&block)
	}

	// todo - A verifier node will not produce rollups, we can check the e.mining to get the node behaviour
	r := rc.produceRollup(&block, blockState)
	rc.signRollup(r)
	rc.checkRollup(r) // Sanity check the produced rollup
	// todo - should store proposal rollups in a different storage as they are ephemeral (round based)
	rc.storage.StoreRollup(r)

	common.TraceWithID(rc.nodeID, "Processed block: b_%d(%d). Produced rollup r_%d", common.ShortHash(block.Hash()), block.NumberU64(), common.ShortHash(r.Hash()))

	encryptedLogs, err := rc.subscriptionManager.EncryptLogs(logs)
	if err != nil {
		log.Panic("Could not encrypt logs. Cause: %s", err)
	}

	return rc.newBlockSubmissionResponse(blockState, r.ToExtRollup(rc.transactionBlobCrypto), encryptedLogs)
}

func (rc *RollupChain) produceRollup(b *types.Block, bs *obscurocore.BlockState) *obscurocore.Rollup {
	headRollup, f := rc.storage.FetchRollup(bs.HeadRollup)
	if !f {
		log.Panic(msgNoRollup)
	}

	// These variables will be used to create the new rollup
	var newRollupTxs []*common.L2Tx
	var newRollupState *state.StateDB

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

	// if !speculativeExecutionSucceeded {
	// In case the speculative execution thread has not succeeded in producing a valid rollup
	// we have to create a new one from the mempool transactions
	// Create a new rollup based on the fromBlock of inclusion of the previous, including all new transactions
	nonce := common.GenerateNonce()
	r := obscurocore.EmptyRollup(rc.hostID, headRollup.Header, b.Hash(), nonce)

	newRollupTxs = rc.mempool.CurrentTxs(headRollup, rc.storage)
	newRollupState = rc.storage.CreateStateDB(r.Header.ParentHash)

	rootHash, successfulTxs, txReceipts, depositReceipts := rc.processState(r, newRollupTxs, newRollupState)

	r.Header.Root = rootHash
	r.Transactions = successfulTxs

	// Postprocessing - withdrawals
	txReceiptsMap := toReceiptMap(txReceipts)
	r.Header.Withdrawals = rc.bridge.RollupPostProcessingWithdrawals(r, newRollupState, txReceiptsMap)

	receipts := allReceipts(txReceipts, depositReceipts)
	if len(receipts) == 0 {
		r.Header.ReceiptHash = types.EmptyRootHash
	} else {
		r.Header.ReceiptHash = types.DeriveSha(receipts, trie.NewStackTrie(nil))
		r.Header.Bloom = types.CreateBloom(receipts)
	}

	if len(successfulTxs) == 0 {
		r.Header.TxHash = types.EmptyRootHash
	} else {
		r.Header.TxHash = types.DeriveSha(types.Transactions(successfulTxs), trie.NewStackTrie(nil))
	}

	// Uncomment this if you want to debug issues related to root state hashes not matching
	dump := strings.Replace(string(newRollupState.Dump(&state.DumpConfig{})), "\n", "", -1)
	log.Trace("Create rollup. State: %s", dump)

	return r
}

// TODO - this belongs in the protocol

func (rc *RollupChain) RoundWinner(parent common.L2RootHash) (common.ExtRollup, bool, error) {
	head, found := rc.storage.FetchRollup(parent)
	if !found {
		return common.ExtRollup{}, false, fmt.Errorf("rollup not found: r_%s", parent)
	}

	headState := rc.storage.FetchHeadState()
	currentHeadRollup, found := rc.storage.FetchRollup(headState.HeadRollup)
	if !found {
		panic("Should not happen since the header hash and the rollup are stored in a batch.")
	}
	// Check if round.winner is being called on an old rollup
	if !bytes.Equal(currentHeadRollup.Hash().Bytes(), parent.Bytes()) {
		return common.ExtRollup{}, false, nil
	}

	common.TraceWithID(rc.nodeID, "Round winner height: %d", head.Header.Number)
	rollupsReceivedFromPeers := rc.storage.FetchRollups(head.NumberU64() + 1)
	// filter out rollups with a different Parent
	var usefulRollups []*obscurocore.Rollup
	for _, rol := range rollupsReceivedFromPeers {
		p := rc.storage.ParentRollup(rol)
		if p == nil {
			common.LogWithID(rc.nodeID, "Received rollup from peer but don't have parent rollup - discarding...")
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
		common.TraceWithID(rc.nodeID, "Publish rollup=r_%d(%d)[r_%d]{proof=b_%d(%d)}. Num Txs: %d. Txs: %v.  Root=%v. ",
			common.ShortHash(winnerRollup.Hash()), winnerRollup.Header.Number,
			common.ShortHash(w.Hash()),
			common.ShortHash(v.Hash()),
			v.NumberU64(),
			len(winnerRollup.Transactions),
			obscurocore.PrintTxs(winnerRollup.Transactions),
			winnerRollup.Header.Root,
		)
		return winnerRollup.ToExtRollup(rc.transactionBlobCrypto), true, nil
	}
	return common.ExtRollup{}, false, nil
}

func (rc *RollupChain) ExecuteOffChainTransaction(encryptedParams common.EncryptedParamsCall) (common.EncryptedResponseCall, error) {
	paramBytes, err := rc.rpcEncryptionManager.DecryptBytes(encryptedParams)
	if err != nil {
		return nil, fmt.Errorf("could not decrypt params in eth_call request. Cause: %w", err)
	}

	contractAddress, err := rpc.ExtractCallParamTo(paramBytes)
	if err != nil {
		return nil, fmt.Errorf("could not extract `to` param from eth_call request. Cause: %w", err)
	}
	from, err := rpc.ExtractCallParamFrom(paramBytes)
	if err != nil {
		return nil, fmt.Errorf("could not extract `from` param from eth_call request. Cause: %w", err)
	}
	data, err := rpc.ExtractCallParamData(paramBytes)
	if err != nil {
		return nil, fmt.Errorf("could not extract `data` param from eth_call request. Cause: %w", err)
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
	log.Trace("!OffChain call: contractAddress=%s, from=%s, data=%s, rollup=r_%d, state=%s", contractAddress, from.Hex(), hexutils.BytesToHex(data), common.ShortHash(r.Hash()), r.Header.Root.Hex())
	s := rc.storage.CreateStateDB(hs.HeadRollup)
	result, err := evm.ExecuteOffChainCall(from, contractAddress, data, s, r.Header, rc.storage, rc.chainConfig)
	if err != nil {
		return nil, err
	}
	if result.Failed() {
		log.Error("!OffChain: Failed to execute contract %s: %s\n", contractAddress.Hex(), result.Err)
		return nil, result.Err
	}

	log.Trace("!OffChain result: %s", hexutils.BytesToHex(result.ReturnData))

	var encodedResult string
	if len(result.ReturnData) != 0 {
		encodedResult = hexutil.Encode(result.ReturnData)
	}
	encryptedResult, err := rc.rpcEncryptionManager.EncryptWithViewingKey(from, []byte(encodedResult))
	if err != nil {
		return nil, fmt.Errorf("enclave could not respond securely to eth_call request. Cause: %w", err)
	}

	return encryptedResult, nil
}

func (rc *RollupChain) GetBalance(encryptedParams common.EncryptedParamsGetBalance) (common.EncryptedResponseGetBalance, error) {
	// We decrypt the request.
	paramBytes, err := rc.rpcEncryptionManager.DecryptBytes(encryptedParams)
	if err != nil {
		return nil, fmt.Errorf("could not decrypt params in eth_getBalance request. Cause: %w", err)
	}

	// We extract the params from the request.
	var paramList []string
	err = json.Unmarshal(paramBytes, &paramList)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal RPC request params from JSON. Cause: %w", err)
	}
	if len(paramList) < 2 {
		return nil, fmt.Errorf("required exactly two params, but received zero")
	}
	// TODO - Replace all usages of `HexToAddress` with a `SafeHexToAddress` that checks that the string does not exceed 20 bytes.
	address := gethcommon.HexToAddress(paramList[0])
	blockNumber := gethrpc.BlockNumber(0)
	err = blockNumber.UnmarshalJSON([]byte(paramList[1]))
	if err != nil {
		return nil, fmt.Errorf("could not parse requested rollup number")
	}

	// We retrieve the rollup of interest.
	rollup, err := rc.getRollup(blockNumber)
	if err != nil {
		return nil, err
	}

	// We get the balance at that rollup.
	blockchainState := rc.storage.CreateStateDB(rollup.Hash())
	if blockchainState == nil || err != nil {
		return nil, err
	}
	balance := (*hexutil.Big)(blockchainState.GetBalance(address))

	// We encrypt the result.
	encryptedBalance, err := rc.rpcEncryptionManager.EncryptWithViewingKey(address, []byte(balance.String()))
	if err != nil {
		return nil, fmt.Errorf("enclave could not respond securely to eth_getBalance request. Cause: %w", err)
	}
	return encryptedBalance, blockchainState.Error()
}

func (rc *RollupChain) signRollup(r *obscurocore.Rollup) {
	var err error
	h := r.Hash()
	r.Header.R, r.Header.S, err = ecdsa.Sign(rand.Reader, rc.enclavePrivateKey, h[:])
	if err != nil {
		log.Panic("Could not sign rollup. Cause: %s", err)
	}
}

func (rc *RollupChain) verifySig(r *obscurocore.Rollup) bool {
	// If this rollup is generated by the current enclave skip the sig verification
	if bytes.Equal(r.Header.Agg.Bytes(), rc.hostID.Bytes()) {
		return true
	}

	h := r.Hash()
	if r.Header.R == nil || r.Header.S == nil {
		panic("Missing signature on rollup")
	}
	pubKey := rc.storage.FetchAttestedKey(r.Header.Agg)
	return ecdsa.Verify(pubKey, h[:], r.Header.R, r.Header.S)
}

// Retrieves the rollup with the given height, with special handling for earliest/latest/pending .
func (rc *RollupChain) getRollup(height gethrpc.BlockNumber) (*obscurocore.Rollup, error) {
	var rollup *obscurocore.Rollup
	switch height {
	case gethrpc.EarliestBlockNumber:
		rollup = rc.storage.FetchGenesisRollup()
	case gethrpc.PendingBlockNumber:
		// TODO - Depends on the current pending rollup; leaving it for a different iteration as it will need more thought.
		return nil, fmt.Errorf("requested balance for pending block. This is not handled currently")
	case gethrpc.LatestBlockNumber:
		rollupHash := rc.storage.FetchHeadState().HeadRollup
		var found bool
		rollup, found = rc.storage.FetchRollup(rollupHash)
		if !found {
			return nil, fmt.Errorf("rollup with requested height %d was not found", height)
		}
	default:
		maybeRollup, found := rc.storage.FetchRollupByHeight(uint64(height))
		if !found {
			return nil, fmt.Errorf("rollup with requested height %d was not found", height)
		}
		rollup = maybeRollup
	}
	return rollup, nil
}
