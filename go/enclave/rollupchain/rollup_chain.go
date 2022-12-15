package rollupchain

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/errutil"
	"github.com/obscuronet/go-obscuro/go/common/gethapi"
	"github.com/obscuronet/go-obscuro/go/common/gethutil"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/enclave/bridge"
	"github.com/obscuronet/go-obscuro/go/enclave/core"
	"github.com/obscuronet/go-obscuro/go/enclave/crosschain"
	"github.com/obscuronet/go-obscuro/go/enclave/db"
	"github.com/obscuronet/go-obscuro/go/enclave/evm"
	"github.com/obscuronet/go-obscuro/go/enclave/mempool"
	"github.com/status-im/keycard-go/hexutils"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethcore "github.com/ethereum/go-ethereum/core"
	gethlog "github.com/ethereum/go-ethereum/log"
	gethrpc "github.com/ethereum/go-ethereum/rpc"
)

// RollupChain represents the canonical chain, and manages the state.
type RollupChain struct {
	hostID      gethcommon.Address
	nodeType    common.NodeType
	chainConfig *params.ChainConfig
	sequencerID gethcommon.Address

	storage              db.Storage
	l1Blockchain         *gethcore.BlockChain
	bridge               *bridge.Bridge
	mempool              mempool.Manager
	faucet               Faucet
	crossChainProcessors *crosschain.Processors

	enclavePrivateKey    *ecdsa.PrivateKey // this is a key known only to the current enclave, and the public key was shared with everyone during attestation
	blockProcessingMutex sync.Mutex
	logger               gethlog.Logger

	// Gas usage values
	// TODO use the ethconfig.Config instead
	GlobalGasCap uint64
	BaseFee      *big.Int
}

func New(
	hostID gethcommon.Address,
	nodeType common.NodeType,
	storage db.Storage,
	l1Blockchain *gethcore.BlockChain,
	bridge *bridge.Bridge,
	crossChainProcessors *crosschain.Processors,
	mempool mempool.Manager,
	privateKey *ecdsa.PrivateKey,
	chainConfig *params.ChainConfig,
	sequencerID gethcommon.Address,
	logger gethlog.Logger,
) *RollupChain {
	return &RollupChain{
		hostID:               hostID,
		nodeType:             nodeType,
		storage:              storage,
		l1Blockchain:         l1Blockchain,
		bridge:               bridge,
		mempool:              mempool,
		faucet:               NewFaucet(),
		crossChainProcessors: crossChainProcessors,
		enclavePrivateKey:    privateKey,
		chainConfig:          chainConfig,
		blockProcessingMutex: sync.Mutex{},
		logger:               logger,
		GlobalGasCap:         5_000_000_000,
		BaseFee:              gethcommon.Big0,
		sequencerID:          sequencerID,
	}
}

// ProcessL1Block is used to update the enclave with an additional L1 block.
func (rc *RollupChain) ProcessL1Block(block types.Block, receipts types.Receipts, isLatest bool) (*common.L2RootHash, *core.Batch, error) {
	rc.blockProcessingMutex.Lock()
	defer rc.blockProcessingMutex.Unlock()

	// We update the L1 chain state.
	err := rc.updateL1State(block, receipts, isLatest)
	if err != nil {
		return nil, nil, err
	}

	// We update the L1 and L2 chain heads.
	newL2Head, producedBatch, err := rc.updateL1AndL2Heads(&block)
	if err != nil {
		return nil, nil, err
	}
	return newL2Head, producedBatch, nil
}

// UpdateL2Chain updates the L2 chain based on the received batch.
func (rc *RollupChain) UpdateL2Chain(batch *core.Batch) (*common.BatchHeader, error) {
	rc.blockProcessingMutex.Lock()
	defer rc.blockProcessingMutex.Unlock()

	isGenesis := false
	if batch.Number().Cmp(big.NewInt(0)) == 0 {
		isGenesis = true
	}

	// todo - joel - this checking is very ugly
	var batchTxReceipts types.Receipts
	var err error
	if !isGenesis {
		if batchTxReceipts, err = rc.checkBatch(batch); err != nil {
			return nil, fmt.Errorf("failed to check batch. Cause: %w", err)
		}
	}
	if err = rc.storage.StoreBatch(batch, batchTxReceipts); err != nil {
		return nil, fmt.Errorf("failed to store batch. Cause: %w", err)
	}
	if err = rc.storage.UpdateHeadBatch(batch.Header.L1Proof, batch, batchTxReceipts); err != nil {
		return nil, fmt.Errorf("could not store new L2 head. Cause: %w", err)
	}
	if isGenesis {
		if err = rc.faucet.CommitGenesisState(rc.storage); err != nil {
			return nil, fmt.Errorf("could not apply faucet preallocation. Cause: %w", err)
		}
	}

	return batch.Header, nil
}

func (rc *RollupChain) GetBalance(accountAddress gethcommon.Address, blockNumber *gethrpc.BlockNumber) (*gethcommon.Address, *hexutil.Big, error) {
	// get account balance at certain block/height
	balance, err := rc.GetBalanceAtBlock(accountAddress, blockNumber)
	if err != nil {
		return nil, nil, err
	}

	// check if account is a contract
	isAddrContract, err := rc.isAccountContractAtBlock(accountAddress, blockNumber)
	if err != nil {
		return nil, nil, err
	}

	// Decide which address to encrypt the result with
	address := accountAddress
	// If the accountAddress is a contract, encrypt with the address of the contract owner
	if isAddrContract {
		txHash, err := rc.storage.GetContractCreationTx(accountAddress)
		if err != nil {
			return nil, nil, err
		}
		transaction, _, _, _, err := rc.storage.GetTransaction(*txHash)
		if err != nil {
			return nil, nil, err
		}
		signer := types.NewLondonSigner(rc.chainConfig.ChainID)

		sender, err := signer.Sender(transaction)
		if err != nil {
			return nil, nil, err
		}
		address = sender
	}

	return &address, balance, nil
}

// GetBalanceAtBlock returns the balance of an account at a certain height
func (rc *RollupChain) GetBalanceAtBlock(accountAddr gethcommon.Address, blockNumber *gethrpc.BlockNumber) (*hexutil.Big, error) {
	chainState, err := rc.getChainStateAtBlock(blockNumber)
	if err != nil {
		return nil, fmt.Errorf("unable to get blockchain state - %w", err)
	}

	return (*hexutil.Big)(chainState.GetBalance(accountAddr)), nil
}

// ExecuteOffChainTransaction executes non-state changing transactions at a given block height (eth_call)
func (rc *RollupChain) ExecuteOffChainTransaction(apiArgs *gethapi.TransactionArgs, blockNumber *gethrpc.BlockNumber) (*gethcore.ExecutionResult, error) {
	result, err := rc.ExecuteOffChainTransactionAtBlock(apiArgs, blockNumber)
	if err != nil {
		rc.logger.Error(fmt.Sprintf("!OffChain: Failed to execute contract %s.", apiArgs.To), log.ErrKey, err.Error())
		return nil, err
	}

	// the execution might have succeeded (err == nil) but the evm contract logic might have failed (result.Failed() == true)
	if result.Failed() {
		rc.logger.Error(fmt.Sprintf("!OffChain: Failed to execute contract %s.", apiArgs.To), log.ErrKey, result.Err)
		return nil, result.Err
	}

	rc.logger.Trace(fmt.Sprintf("!OffChain result: %s", hexutils.BytesToHex(result.ReturnData)))

	return result, nil
}

func (rc *RollupChain) ExecuteOffChainTransactionAtBlock(apiArgs *gethapi.TransactionArgs, blockNumber *gethrpc.BlockNumber) (*gethcore.ExecutionResult, error) {
	// TODO review this during gas mechanics implementation
	callMsg, err := apiArgs.ToMessage(rc.GlobalGasCap, rc.BaseFee)
	if err != nil {
		return nil, fmt.Errorf("unable to convert TransactionArgs to Message - %w", err)
	}

	// fetch the chain state at given batch
	blockState, err := rc.getChainStateAtBlock(blockNumber)
	if err != nil {
		return nil, err
	}

	batch, err := rc.getBatch(*blockNumber)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch head state batch. Cause: %w", err)
	}

	rc.logger.Trace(
		fmt.Sprintf("!OffChain call: contractAddress=%s, from=%s, data=%s, batch=b_%d, state=%s",
			callMsg.To(),
			callMsg.From(),
			hexutils.BytesToHex(callMsg.Data()),
			common.ShortHash(*batch.Hash()),
			batch.Header.Root.Hex()),
	)

	result, err := evm.ExecuteOffChainCall(&callMsg, blockState, batch.Header, rc.storage, rc.chainConfig, rc.logger)
	if err != nil {
		// also return the result as the result can be evaluated on some errors like ErrIntrinsicGas
		return result, err
	}

	// the execution outcome was unsuccessful, but it was able to execute the call
	if result.Failed() {
		// do not return an error
		// the result object should be evaluated upstream
		rc.logger.Error(fmt.Sprintf("!OffChain: Failed to execute contract %s.", callMsg.To()), log.ErrKey, result.Err)
	}

	return result, nil
}

func (rc *RollupChain) updateL1State(block types.Block, receipts types.Receipts, isLatest bool) error {
	// We check whether we've already processed the block.
	_, err := rc.storage.FetchBlock(block.Hash())
	if err == nil {
		return common.ErrBlockAlreadyProcessed
	}
	if !errors.Is(err, errutil.ErrNotFound) {
		return fmt.Errorf("could not retrieve block. Cause: %w", err)
	}

	// Reject block if not provided with matching receipts.
	// This needs to happen before saving the block as otherwise it will be considered as processed.
	if rc.crossChainProcessors.Enabled() && !crosschain.VerifyReceiptHash(&block, receipts) {
		return errors.New("receipts do not match receipt_root in block")
	}

	// We insert the block into the L1 chain and store it.
	ingestionType, err := rc.insertBlockIntoL1Chain(&block, isLatest)
	if err != nil {
		// Do not store the block if the L1 chain insertion failed
		return err
	}
	rc.logger.Trace("block inserted successfully",
		"height", block.NumberU64(), "hash", block.Hash(), "ingestionType", ingestionType)

	rc.storage.StoreBlock(&block)

	// This requires block to be stored first ... but can permanently fail a block
	err = rc.crossChainProcessors.Remote.StoreCrossChainMessages(&block, receipts)
	if err != nil {
		return errors.New("failed to process cross chain messages")
	}

	return nil
}

// Inserts the block into the L1 chain if it exists and the block is not the genesis block
// note: this method shouldn't be called for blocks we've seen before
func (rc *RollupChain) insertBlockIntoL1Chain(block *types.Block, isLatest bool) (*blockIngestionType, error) {
	if rc.l1Blockchain != nil {
		_, err := rc.l1Blockchain.InsertChain(types.Blocks{block})
		if err != nil {
			return nil, fmt.Errorf("block was invalid: %w", err)
		}
	}
	// todo: this is minimal L1 tracking/validation, and should be removed when we are using geth's blockchain or lightchain structures for validation
	prevL1Head, err := rc.storage.FetchHeadBlock()

	if err != nil {
		if errors.Is(err, errutil.ErrNotFound) {
			// todo: we should enforce that this block is a configured hash (e.g. the L1 management contract deployment block)
			return &blockIngestionType{latest: isLatest, fork: false, preGenesis: true}, nil
		}
		return nil, fmt.Errorf("could not retrieve head block. Cause: %w", err)

		// we do a basic sanity check, comparing the received block to the head block on the chain
	} else if block.ParentHash() != prevL1Head.Hash() {
		lcaBlock, err := gethutil.LCA(block, prevL1Head, rc.storage)
		if err != nil {
			return nil, common.ErrBlockAncestorNotFound
		}
		rc.logger.Trace("parent not found",
			"blkHeight", block.NumberU64(), "blkHash", block.Hash(),
			"l1HeadHeight", prevL1Head.NumberU64(), "l1HeadHash", prevL1Head.Hash(),
			"lcaHeight", lcaBlock.NumberU64(), "lcaHash", lcaBlock.Hash(),
		)
		if lcaBlock.NumberU64() >= prevL1Head.NumberU64() {
			// This is an unexpected error scenario (a bug) because if:
			// lca == prevL1Head:
			//   if prev L1 head is (e.g) a grandfather of ingested block, and block's parent has been seen (else LCA would error),
			//   then why is ingested block's parent not the prev l1 head
			// lca > prevL1Head:
			//   this would imply ingested block is earlier on the same branch as l1 head, but ingested block should not have been seen before
			rc.logger.Error("unexpected blockchain state, incoming block is not child of L1 head and not an earlier fork of L1 head",
				"blkHeight", block.NumberU64(), "blkHash", block.Hash(),
				"l1HeadHeight", prevL1Head.NumberU64(), "l1HeadHash", prevL1Head.Hash(),
				"lcaHeight", lcaBlock.NumberU64(), "lcaHash", lcaBlock.Hash(),
			)
			return nil, errors.New("unexpected blockchain state")
		}

		// ingested block is on a different branch to the previously ingested block - we may have to rewind L2 state
		return &blockIngestionType{latest: isLatest, fork: true, preGenesis: false}, nil
	}

	// this is the typical, happy-path case. The ingested block's parent was the previously ingested block.
	return &blockIngestionType{latest: isLatest, fork: false, preGenesis: false}, nil
}

// Updates the L1 and L2 chain heads, and returns the new L2 head hash and the produced batch, if there is one.
func (rc *RollupChain) updateL1AndL2Heads(block *types.Block) (*common.L2RootHash, *core.Batch, error) {
	err := rc.processRollups(block)
	if err != nil {
		return nil, nil, fmt.Errorf("could not process rollup in block. Cause: %w", err)
	}

	// TODO - #718 - Cannot assume that the most recent rollup is on the previous block anymore. May be on the same block.
	// We retrieve the current L2 head.
	// todo - joel - address v ugly handling of no L2 head existing
	isGenesis := false
	l2Head, err := rc.storage.FetchHeadBatchForBlock(block.ParentHash())
	if err != nil {
		if !errors.Is(err, errutil.ErrNotFound) {
			return nil, nil, fmt.Errorf("could not retrieve current head batch. Cause: %w", err)
		}
		isGenesis = true
	}
	var l2HeadTxReceipts types.Receipts
	if !isGenesis {
		if l2HeadTxReceipts, err = rc.storage.GetReceiptsByHash(*l2Head.Hash()); err != nil {
			return nil, nil, fmt.Errorf("could not fetch batch receipts. Cause: %w", err)
		}
	}

	// If we're the sequencer, we produce a new L2 head to replace the old one.
	if rc.nodeType == common.Sequencer {
		if isGenesis {
			if l2Head, err = rc.produceGenesisBatch(block.Hash()); err != nil {
				return nil, nil, fmt.Errorf("could not produce batch. Cause: %w", err)
			}
		} else {
			if l2Head, err = rc.produceBatch(block); err != nil {
				return nil, nil, fmt.Errorf("could not produce batch. Cause: %w", err)
			}
		}
		if err = rc.signBatch(l2Head); err != nil {
			return nil, nil, fmt.Errorf("could not sign batch. Cause: %w", err)
		}
		if l2HeadTxReceipts, err = rc.getTxReceipts(l2Head); err != nil {
			return nil, nil, fmt.Errorf("could not get batch transaction receipts. Cause: %w", err)
		}
		if err = rc.storage.StoreBatch(l2Head, l2HeadTxReceipts); err != nil {
			return nil, nil, fmt.Errorf("failed to store batch. Cause: %w", err)
		}
	}

	// We update the chain heads.
	if l2Head != nil {
		if err = rc.storage.UpdateHeadBatch(block.Hash(), l2Head, l2HeadTxReceipts); err != nil {
			return nil, nil, fmt.Errorf("could not store new head. Cause: %w", err)
		}
		// todo - joel - this branch only applies to sequencer - can we move it up?
		if isGenesis {
			if err = rc.faucet.CommitGenesisState(rc.storage); err != nil {
				return nil, nil, fmt.Errorf("could not apply faucet preallocation. Cause: %w", err)
			}
		}
		if err = rc.storage.UpdateL1Head(block.Hash()); err != nil {
			return nil, nil, fmt.Errorf("could not store new L1 head. Cause: %w", err)
		}
	}

	// We return the produced batch, if we've produced one.
	var producedBatch *core.Batch
	if rc.nodeType == common.Sequencer {
		producedBatch = l2Head
	}
	if l2Head != nil {
		return l2Head.Hash(), producedBatch, nil
	}
	return nil, producedBatch, nil
}

// Creates a genesis batch linked to the provided L1 block and signs it.
func (rc *RollupChain) produceGenesisBatch(blkHash common.L1RootHash) (*core.Batch, error) {
	preFundGenesisState, err := rc.faucet.GetGenesisRoot(rc.storage)
	if err != nil {
		return nil, err
	}

	genesisBatch := &core.Batch{
		Header: &common.BatchHeader{
			Agg:         gethcommon.HexToAddress("0x0"),
			ParentHash:  common.L2RootHash{},
			L1Proof:     blkHash,
			Root:        *preFundGenesisState,
			TxHash:      types.EmptyRootHash,
			Number:      big.NewInt(int64(0)),
			Withdrawals: []common.Withdrawal{},
			ReceiptHash: types.EmptyRootHash,
			Time:        uint64(time.Now().Unix()),
		},
		Transactions: []*common.L2Tx{},
	}

	// TODO: Figure out a better way to bootstrap the system contracts.
	deployTx, err := rc.crossChainProcessors.Local.GenerateMessageBusDeployTx()
	if err != nil {
		rc.logger.Crit("Could not create message bus deployment transaction", "Error", err)
	}

	// Add transaction to mempool so it gets processed when it can.
	// Should be the first transaction to be processed.
	if err := rc.mempool.AddMempoolTx(deployTx); err != nil {
		rc.logger.Crit("Cannot create synthetic transaction for deploying the message bus contract on :|")
	}

	return genesisBatch, nil
}

// This is where transactions are executed and the state is calculated.
// Obscuro includes a message bus embedded in the platform, and this method is responsible for transferring messages as well.
// The batch can be a final batch as received from peers or the batch under construction.
func (rc *RollupChain) processState(batch *core.Batch, txs []*common.L2Tx, stateDB *state.StateDB) (common.L2RootHash, []*common.L2Tx, []*types.Receipt, []*types.Receipt) {
	var executedTransactions []*common.L2Tx
	var txReceipts []*types.Receipt

	txResults := evm.ExecuteTransactions(txs, stateDB, batch.Header, rc.storage, rc.chainConfig, 0, rc.logger)
	for _, tx := range txs {
		result, f := txResults[tx.Hash()]
		if !f {
			rc.logger.Crit("There should be an entry for each transaction ")
		}
		rec, foundReceipt := result.(*types.Receipt)
		if foundReceipt {
			executedTransactions = append(executedTransactions, tx)
			txReceipts = append(txReceipts, rec)
		} else {
			// Exclude all errors
			rc.logger.Info(fmt.Sprintf("Excluding transaction %s from batch b_%d. Cause: %s", tx.Hash().Hex(), common.ShortHash(*batch.Hash()), result))
		}
	}

	// always process deposits last, either on top of the rollup produced speculatively or the newly created rollup
	// process deposits from the fromBlock of the parent to the current block (which is the fromBlock of the new rollup)
	parent, err := rc.storage.FetchBatch(batch.Header.ParentHash)
	if err != nil {
		rc.logger.Crit("Sanity check. Rollup has no parent.", log.ErrKey, err)
	}

	parentProof, err := rc.storage.FetchBlock(parent.Header.L1Proof)
	if err != nil {
		rc.logger.Crit(fmt.Sprintf("Could not retrieve a proof for batch %s", batch.Hash()), log.ErrKey, err)
	}
	batchProof, err := rc.storage.FetchBlock(batch.Header.L1Proof)
	if err != nil {
		rc.logger.Crit(fmt.Sprintf("Could not retrieve a proof for batch %s", batch.Hash()), log.ErrKey, err)
	}

	// TODO: Remove this depositing logic once the new bridge is added.
	depositTxs := rc.bridge.ExtractDeposits(
		parentProof,
		batchProof,
		rc.storage,
		stateDB,
	)

	depositResponses := evm.ExecuteTransactions(depositTxs, stateDB, batch.Header, rc.storage, rc.chainConfig, len(executedTransactions), rc.logger)
	depositReceipts := make([]*types.Receipt, len(depositResponses))
	if len(depositResponses) != len(depositTxs) {
		rc.logger.Crit("Sanity check. Some deposit transactions failed.")
	}
	i := 0
	for _, resp := range depositResponses {
		rec, ok := resp.(*types.Receipt)
		if !ok {
			// TODO - Handle the case of an error (e.g. insufficient funds).
			rc.logger.Crit("Sanity check. Expected a receipt", log.ErrKey, resp)
		}
		depositReceipts[i] = rec
		i++
	}

	messages := rc.crossChainProcessors.Local.RetrieveInboundMessages(parentProof, batchProof, stateDB)
	transactions := rc.crossChainProcessors.Local.CreateSyntheticTransactions(messages, stateDB)
	syntheticTransactionsResponses := evm.ExecuteTransactions(transactions, stateDB, batch.Header, rc.storage, rc.chainConfig, len(executedTransactions), rc.logger)
	synthReceipts := make([]*types.Receipt, len(syntheticTransactionsResponses))
	if len(syntheticTransactionsResponses) != len(transactions) {
		rc.logger.Crit("Sanity check. Some synthetic transactions failed.")
	}

	i = 0
	for _, resp := range syntheticTransactionsResponses {
		rec, ok := resp.(*types.Receipt)
		if !ok { // Еxtract reason for failing deposit.
			// TODO - Handle the case of an error (e.g. insufficient funds).
			rc.logger.Crit("Sanity check. Expected a receipt", log.ErrKey, resp)
		}

		if rec.Status == 0 { // Synthetic transactions should not fail. In case of failure get the revert reason.
			failingTx := transactions[i]
			txCallMessage := types.NewMessage(
				rc.crossChainProcessors.Local.GetOwner(),
				failingTx.To(),
				stateDB.GetNonce(rc.crossChainProcessors.Local.GetOwner()),
				failingTx.Value(),
				failingTx.Gas(),
				gethcommon.Big0,
				gethcommon.Big0,
				gethcommon.Big0,
				failingTx.Data(),
				failingTx.AccessList(),
				false)

			clonedDB := stateDB.Copy()
			res, err := evm.ExecuteOffChainCall(&txCallMessage, clonedDB, batch.Header, rc.storage, rc.chainConfig, rc.logger)
			rc.logger.Crit("Synthetic transaction failed!", log.ErrKey, err, "result", res)
		}

		synthReceipts[i] = rec
		i++
	}

	rootHash, err := stateDB.Commit(true)
	if err != nil {
		rc.logger.Crit("could not commit to state DB. ", log.ErrKey, err)
	}

	sort.Sort(sortByTxIndex(txReceipts))
	sort.Sort(sortByTxIndex(depositReceipts))

	// todo - handle the tx execution logs
	return rootHash, executedTransactions, txReceipts, depositReceipts
}

// Checks the internal validity of the batch.
func (rc *RollupChain) isInternallyValidBatch(batch *core.Batch) (types.Receipts, error) {
	stateDB, err := rc.storage.CreateStateDB(batch.Header.ParentHash)
	if err != nil {
		return nil, fmt.Errorf("could not create stateDB. Cause: %w", err)
	}

	// calculate the state to compare with what is in the batch
	rootHash, executedTxs, txReceipts, depositReceipts := rc.processState(batch, batch.Transactions, stateDB)
	if len(executedTxs) != len(batch.Transactions) {
		return nil, fmt.Errorf("all transactions that are included in a batch must be executed")
	}

	// Check that the root hash in the header matches the root hash as calculated.
	if !bytes.Equal(rootHash.Bytes(), batch.Header.Root.Bytes()) {
		dump := strings.Replace(string(stateDB.Dump(&state.DumpConfig{})), "\n", "", -1)
		return nil, fmt.Errorf("verify batch b_%d: Calculated a different state. This should not happen as there are no malicious actors yet. \nGot: %s\nExp: %s\nHeight:%d\nTxs:%v\nState: %s.\nDeposits: %+v",
			common.ShortHash(*batch.Hash()), rootHash, batch.Header.Root, batch.Header.Number, core.PrintTxs(batch.Transactions), dump, depositReceipts)
	}

	// Check that the withdrawals in the header match the withdrawals as calculated.
	withdrawals := rc.bridge.BatchPostProcessingWithdrawals(batch, stateDB, toReceiptMap(txReceipts))
	for i, w := range withdrawals {
		hw := batch.Header.Withdrawals[i]
		if hw.Amount.Cmp(w.Amount) != 0 || hw.Recipient != w.Recipient || hw.Contract != w.Contract {
			return nil, fmt.Errorf("verify batch r_%d: Withdrawals don't match", common.ShortHash(*batch.Hash()))
		}
	}

	// Check that the receipts bloom in the header matches the receipts bloom as calculated.
	receipts := allReceipts(txReceipts, depositReceipts)
	receiptBloom := types.CreateBloom(receipts)
	if !bytes.Equal(receiptBloom.Bytes(), batch.Header.Bloom.Bytes()) {
		return nil, fmt.Errorf("verify batch r_%d: Invalid bloom (remote: %x  local: %x)", common.ShortHash(*batch.Hash()), batch.Header.Bloom, receiptBloom)
	}

	// Check that the receipts SHA in the header matches the receipts SHA as calculated.
	receiptSha := types.DeriveSha(receipts, trie.NewStackTrie(nil))
	if !bytes.Equal(receiptSha.Bytes(), batch.Header.ReceiptHash.Bytes()) {
		return nil, fmt.Errorf("verify batch r_%d: invalid receipt root hash (remote: %x local: %x)", common.ShortHash(*batch.Hash()), batch.Header.ReceiptHash, receiptSha)
	}

	// Check that the signature is valid.
	// todo: #1297 re-enable seq sig validation ASAP - once the testnet nodes all have access to the sequencer ID
	//if err := rc.validateSequencerSig(batch.Hash(), &batch.Header.Agg, batch.Header.R, batch.Header.S); err != nil {
	//	return nil, fmt.Errorf("verify batch r_%d: invalid signature. Cause: %s", common.ShortHash(*batch.Hash()), err.Error())
	//}

	// todo - check that the transactions hash to the header.txHash

	return txReceipts, nil
}

// Returns the receipts for the transactions in the batch.
func (rc *RollupChain) getTxReceipts(batch *core.Batch) ([]*types.Receipt, error) {
	if batch.Header.Number.Cmp(big.NewInt(int64(common.L2GenesisHeight))) == 0 {
		return nil, nil //nolint:nilnil
	}

	stateDB, err := rc.storage.CreateStateDB(batch.Header.ParentHash)
	if err != nil {
		return nil, fmt.Errorf("could not create stateDB. Cause: %w", err)
	}

	// calculate the state to compare with what is in the batch
	_, _, txReceipts, _ := rc.processState(batch, batch.Transactions, stateDB) //nolint:dogsled
	return txReceipts, nil
}

// Checks that the batch is valid, both internally and relative to the previously-stored batches.
func (rc *RollupChain) checkBatch(batch *core.Batch) ([]*types.Receipt, error) {
	// We check that the batch is internally valid (e.g. its header matches its contents).
	txReceipts, err := rc.isInternallyValidBatch(batch)
	if err != nil {
		return nil, fmt.Errorf("batch was invalid. Cause: %w", err)
	}

	// We check that we've stored the batch's parent.
	if _, err = rc.storage.FetchBatch(batch.Header.ParentHash); err != nil {
		return nil, fmt.Errorf("could not retrieve parent batch. Cause: %w", err)
	}

	// TODO - #718 - Check that the transactions in the batch are unique

	return txReceipts, nil
}

func (rc *RollupChain) signBatch(batch *core.Batch) error {
	var err error
	h := batch.Hash()
	batch.Header.R, batch.Header.S, err = ecdsa.Sign(rand.Reader, rc.enclavePrivateKey, h[:])
	if err != nil {
		return fmt.Errorf("could not sign batch. Cause: %w", err)
	}
	return nil
}

// Checks that the header is signed validly by the sequencer.
// todo: #1297 remove the nolint:unused here when validation usage is re-enabled
func (rc *RollupChain) validateSequencerSig(headerHash *gethcommon.Hash, aggregator *gethcommon.Address, sigR *big.Int, sigS *big.Int) error { //nolint:unused
	// Batches and rollups should only be produced by the sequencer.
	// TODO - #718 - Sequencer identities should be retrieved from the L1 management contract.
	if !bytes.Equal(aggregator.Bytes(), rc.sequencerID.Bytes()) {
		return fmt.Errorf("expected batch to be produced by sequencer %s, but was produced by %s", rc.sequencerID.Hex(), aggregator.Hex())
	}

	if sigR == nil || sigS == nil {
		return fmt.Errorf("missing signature on batch")
	}

	pubKey, err := rc.storage.FetchAttestedKey(*aggregator)
	if err != nil {
		return fmt.Errorf("could not retrieve attested key for aggregator %s. Cause: %w", aggregator, err)
	}

	if !ecdsa.Verify(pubKey, headerHash.Bytes(), sigR, sigS) {
		return fmt.Errorf("could not verify ECDSA signature")
	}
	return nil
}

// Retrieves the batch with the given height, with special handling for earliest/latest/pending .
func (rc *RollupChain) getBatch(height gethrpc.BlockNumber) (*core.Batch, error) {
	var batch *core.Batch
	switch height {
	case gethrpc.EarliestBlockNumber:
		genesisBatch, err := rc.storage.FetchBatchByHeight(0)
		if err != nil {
			return nil, fmt.Errorf("could not retrieve genesis rollup. Cause: %w", err)
		}
		batch = genesisBatch
	case gethrpc.PendingBlockNumber:
		// TODO - Depends on the current pending rollup; leaving it for a different iteration as it will need more thought.
		return nil, fmt.Errorf("requested balance for pending block. This is not handled currently")
	case gethrpc.LatestBlockNumber:
		headBatch, err := rc.storage.FetchHeadBatch()
		if err != nil {
			return nil, fmt.Errorf("batch with requested height %d was not found. Cause: %w", height, err)
		}
		batch = headBatch
	default:
		maybeBatch, err := rc.storage.FetchBatchByHeight(uint64(height))
		if err != nil {
			return nil, fmt.Errorf("batch with requested height %d could not be retrieved. Cause: %w", height, err)
		}
		batch = maybeBatch
	}
	return batch, nil
}

// Creates a batch.
func (rc *RollupChain) produceBatch(block *types.Block) (*core.Batch, error) {
	headBatch, err := rc.storage.FetchHeadBatchForBlock(block.ParentHash())
	if err != nil {
		return nil, fmt.Errorf("could not retrieve head batch. Cause: %w", err)
	}

	// These variables will be used to create the new batch
	var newBatchTxs []*common.L2Tx
	var newBatchState *state.StateDB

	// Create a new batch based on the fromBlock of inclusion of the previous, including all new transactions
	batch, err := core.EmptyBatch(rc.hostID, headBatch.Header, block.Hash())
	if err != nil {
		return nil, fmt.Errorf("could not create batch. Cause: %w", err)
	}

	newBatchTxs, err = rc.mempool.CurrentTxs(headBatch, rc.storage)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve current transactions. Cause: %w", err)
	}

	newBatchState, err = rc.storage.CreateStateDB(batch.Header.ParentHash)
	if err != nil {
		return nil, fmt.Errorf("could not create stateDB. Cause: %w", err)
	}

	rootHash, successfulTxs, txReceipts, depositReceipts := rc.processState(batch, newBatchTxs, newBatchState)

	batch.Header.Root = rootHash
	batch.Transactions = successfulTxs

	// Postprocessing - withdrawals
	txReceiptsMap := toReceiptMap(txReceipts)
	batch.Header.Withdrawals = rc.bridge.BatchPostProcessingWithdrawals(batch, newBatchState, txReceiptsMap)
	crossChainMessages, err := rc.crossChainProcessors.Local.ExtractOutboundMessages(txReceipts)
	if err != nil {
		rc.logger.Crit("Extracting messages L2->L1 failed", err, log.CmpKey, log.CrossChainCmp)
	}

	batch.Header.CrossChainMessages = crossChainMessages

	rc.logger.Trace(fmt.Sprintf("Added %d cross chain messages to batch. Equivalent withdrawals in header - %d",
		len(batch.Header.CrossChainMessages), len(batch.Header.Withdrawals)), log.CmpKey, log.CrossChainCmp)

	crossChainBind, err := rc.storage.FetchBlock(batch.Header.L1Proof)
	if err != nil {
		rc.logger.Crit("Failed to extract batch proof that should exist!")
	}

	batch.Header.LatestInboudCrossChainHash = crossChainBind.Hash()
	batch.Header.LatestInboundCrossChainHeight = crossChainBind.Number()

	receipts := allReceipts(txReceipts, depositReceipts)
	if len(receipts) == 0 {
		batch.Header.ReceiptHash = types.EmptyRootHash
	} else {
		batch.Header.ReceiptHash = types.DeriveSha(receipts, trie.NewStackTrie(nil))
		batch.Header.Bloom = types.CreateBloom(receipts)
	}

	if len(successfulTxs) == 0 {
		batch.Header.TxHash = types.EmptyRootHash
	} else {
		batch.Header.TxHash = types.DeriveSha(types.Transactions(successfulTxs), trie.NewStackTrie(nil))
	}

	rc.logger.Trace("Create batch.",
		"State", gethlog.Lazy{Fn: func() string {
			return strings.Replace(string(newBatchState.Dump(&state.DumpConfig{})), "\n", "", -1)
		}},
	)

	return batch, nil
}

// Returns the state of the chain at height
// TODO make this cacheable
func (rc *RollupChain) getChainStateAtBlock(blockNumber *gethrpc.BlockNumber) (*state.StateDB, error) {
	// We retrieve the batch of interest.
	batch, err := rc.getBatch(*blockNumber)
	if err != nil {
		return nil, err
	}

	// We get that of the chain at that height
	blockchainState, err := rc.storage.CreateStateDB(*batch.Hash())
	if err != nil {
		return nil, fmt.Errorf("could not create stateDB. Cause: %w", err)
	}

	if blockchainState == nil {
		return nil, fmt.Errorf("unable to fetch chain state for batch %s", batch.Hash().Hex())
	}

	return blockchainState, err
}

// Returns the whether the account is a contract or not at a certain height
func (rc *RollupChain) isAccountContractAtBlock(accountAddr gethcommon.Address, blockNumber *gethrpc.BlockNumber) (bool, error) {
	chainState, err := rc.getChainStateAtBlock(blockNumber)
	if err != nil {
		return false, fmt.Errorf("unable to get blockchain state - %w", err)
	}

	return len(chainState.GetCode(accountAddr)) > 0, nil
}

// Validates and stores the rollup in a given block.
func (rc *RollupChain) processRollups(block *common.L1Block) error {
	// We extract the rollups from the block.
	rollups := rc.bridge.ExtractRollups(block, rc.storage)

	// We sort the rollups by number in ascending order, in order to process them in the correct order.
	sort.Slice(rollups, func(i, j int) bool {
		return rollups[i].Header.Number.Cmp(rollups[j].Header.Number) < 0
	})

	// We check if there are any duplicates.
	for idx, rollup := range rollups {
		if idx+1 >= len(rollups) {
			break
		}
		if rollup.Header.Number.Cmp(rollups[idx+1].Header.Number) == 0 {
			return fmt.Errorf("duplicates rollups found in block; two rollups with number %d", rollup.Header.Number)
		}
	}

	// We retrieve the current head rollup.
	l1ParentHash := block.ParentHash()
	currentHeadRollup, err := rc.storage.FetchHeadRollupForBlock(&l1ParentHash)
	if err != nil && !errors.Is(err, errutil.ErrNotFound) {
		return fmt.Errorf("could not fetch current L2 head rollup")
	}

	// We check each rollup.
	for _, rollup := range rollups {
		// todo: #1297 re-enable seq sig validation ASAP - once the testnet nodes all have access to the sequencer ID
		//if err := rc.validateSequencerSig(rollup.Hash(), &rollup.Header.Agg, rollup.Header.R, rollup.Header.S); err != nil {
		//	return fmt.Errorf("rollup signature was invalid. Cause: %w", err)
		//}

		// We check that the rollups are sequential.
		if currentHeadRollup == nil {
			if rollup.Number().Cmp(big.NewInt(0)) != 0 {
				return fmt.Errorf("received batch with number %d but no genesis batch is stored", rollup.Number())
			}
		} else {
			if rollup.Header.ParentHash.Hex() != currentHeadRollup.Hash().Hex() {
				return fmt.Errorf("found gap in rollup chain. Rollup %s's parent was %s instead of %s",
					rollup.Header.Hash(), rollup.Header.ParentHash, currentHeadRollup.Header.Hash())
			}
		}

		// TODO - #718 - Validate the rollups in the block against the stored batches.

		if err = rc.storage.StoreRollup(rollup); err != nil {
			return fmt.Errorf("could not store rollup. Cause: %w", err)
		}

		currentHeadRollup = rollup
	}

	// We update the current head rollup to the latest processed rollup (unless no head has been stored yet).
	if currentHeadRollup != nil {
		l1Head := block.Hash()
		if err = rc.storage.UpdateHeadRollup(&l1Head, currentHeadRollup.Hash()); err != nil {
			return fmt.Errorf("could not update L2 head rollup. Cause: %w", err)
		}
	}

	return nil
}

type sortByTxIndex []*types.Receipt

func (c sortByTxIndex) Len() int           { return len(c) }
func (c sortByTxIndex) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c sortByTxIndex) Less(i, j int) bool { return c[i].TransactionIndex < c[j].TransactionIndex }

func toReceiptMap(txReceipts []*types.Receipt) map[gethcommon.Hash]*types.Receipt {
	receiptMap := make(map[gethcommon.Hash]*types.Receipt, 0)
	for _, receipt := range txReceipts {
		receiptMap[receipt.TxHash] = receipt
	}
	return receiptMap
}

func allReceipts(txReceipts []*types.Receipt, depositReceipts []*types.Receipt) types.Receipts {
	return append(txReceipts, depositReceipts...)
}
