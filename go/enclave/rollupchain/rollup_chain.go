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
	"github.com/obscuronet/go-obscuro/go/enclave/crypto"
	"github.com/obscuronet/go-obscuro/go/enclave/db"
	"github.com/obscuronet/go-obscuro/go/enclave/evm"
	"github.com/obscuronet/go-obscuro/go/enclave/mempool"
	"github.com/status-im/keycard-go/hexutils"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethcore "github.com/ethereum/go-ethereum/core"
	gethlog "github.com/ethereum/go-ethereum/log"
	gethrpc "github.com/ethereum/go-ethereum/rpc"
)

const (
	PreGenesis BlockStage = iota
	Genesis
	PostGenesis
)

// BlockStage represents where the block falls in the L2 chain's lifecycle - whether it's a pre-genesis, genesis or
// post-genesis L2 block.
type BlockStage int64

// RollupChain represents the canonical chain, and manages the state.
type RollupChain struct {
	hostID      gethcommon.Address
	nodeType    common.NodeType
	chainConfig *params.ChainConfig

	storage               db.Storage
	l1Blockchain          *gethcore.BlockChain
	bridge                *bridge.Bridge
	transactionBlobCrypto crypto.TransactionBlobCrypto // todo - remove
	mempool               mempool.Manager
	faucet                Faucet
	crossChainProcessors  *crosschain.Processors

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
	txCrypto crypto.TransactionBlobCrypto,
	mempool mempool.Manager,
	privateKey *ecdsa.PrivateKey,
	chainConfig *params.ChainConfig,
	logger gethlog.Logger,
) *RollupChain {
	return &RollupChain{
		hostID:                hostID,
		nodeType:              nodeType,
		storage:               storage,
		l1Blockchain:          l1Blockchain,
		bridge:                bridge,
		transactionBlobCrypto: txCrypto,
		mempool:               mempool,
		faucet:                NewFaucet(),
		crossChainProcessors:  crossChainProcessors,
		enclavePrivateKey:     privateKey,
		chainConfig:           chainConfig,
		blockProcessingMutex:  sync.Mutex{},
		logger:                logger,
		GlobalGasCap:          5_000_000_000,
		BaseFee:               gethcommon.Big0,
	}
}

// ProduceGenesisRollup creates a genesis rollup linked to the provided L1 block and signs it.
func (rc *RollupChain) ProduceGenesisRollup(blkHash common.L1RootHash) (*core.Rollup, error) {
	preFundGenesisState, err := rc.faucet.GetGenesisRoot(rc.storage)
	if err != nil {
		return nil, err
	}

	h := common.Header{
		Agg:         gethcommon.HexToAddress("0x0"),
		ParentHash:  common.L2RootHash{},
		L1Proof:     blkHash,
		Root:        *preFundGenesisState,
		TxHash:      types.EmptyRootHash,
		Number:      big.NewInt(int64(0)),
		Withdrawals: []common.Withdrawal{},
		ReceiptHash: types.EmptyRootHash,
		Time:        uint64(time.Now().Unix()),
	}
	rolGenesis := &core.Rollup{
		Header:       &h,
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

	err = rc.signRollup(rolGenesis)
	if err != nil {
		return nil, fmt.Errorf("could not sign genesis rollup. Cause: %w", err)
	}

	return rolGenesis, nil
}

// ProcessL1Block is used to update the enclave with an additional L1 block.
func (rc *RollupChain) ProcessL1Block(block types.Block, receipts types.Receipts, isLatest bool) (*common.L2RootHash, *core.Rollup, error) {
	rc.blockProcessingMutex.Lock()
	defer rc.blockProcessingMutex.Unlock()

	// We update the L1 chain state.
	err := rc.updateL1State(block, receipts, isLatest)
	if err != nil {
		return nil, nil, err
	}

	// We update the L1 and L2 chain heads.
	newL2Head, producedRollup, err := rc.updateL1AndL2Heads(&block)
	if err != nil {
		return nil, nil, err
	}
	return newL2Head, producedRollup, nil
}

// UpdateL2Chain updates the L2 chain based on the received batch.
func (rc *RollupChain) UpdateL2Chain(batch *common.ExtBatch) (*common.Header, error) {
	rc.blockProcessingMutex.Lock()
	defer rc.blockProcessingMutex.Unlock()

	extRollup := common.ExtRollup{
		Header:          batch.Header,
		TxHashes:        batch.TxHashes,
		EncryptedTxBlob: batch.EncryptedTxBlob,
	}
	rollup := core.ToEnclaveRollup(&extRollup, rc.transactionBlobCrypto)

	// We retrieve the genesis batch from the L1 chain instead.
	if rollup.Header.Number.Cmp(big.NewInt(int64(common.L2GenesisHeight))) == 0 {
		return nil, nil //nolint:nilnil
	}

	rollupTxReceipts, err := rc.checkRollup(rollup)
	if err != nil {
		return nil, fmt.Errorf("failed to check rollup. Cause: %w", err)
	}
	if err = rc.storage.StoreRollup(rollup, rollupTxReceipts); err != nil {
		return nil, fmt.Errorf("failed to store rollup. Cause: %w", err)
	}
	if err = rc.storage.UpdateL2HeadForL1Block(batch.Header.L1Proof, rollup, rollupTxReceipts); err != nil {
		return nil, fmt.Errorf("could not store new L2 head. Cause: %w", err)
	}

	return rollup.Header, nil
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

	// fetch the chain state at given rollup
	blockState, err := rc.getChainStateAtBlock(blockNumber)
	if err != nil {
		return nil, err
	}

	r, err := rc.getRollup(*blockNumber)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch head state rollup. Cause: %w", err)
	}

	rc.logger.Trace(
		fmt.Sprintf("!OffChain call: contractAddress=%s, from=%s, data=%s, rollup=r_%d, state=%s",
			callMsg.To(),
			callMsg.From(),
			hexutils.BytesToHex(callMsg.Data()),
			common.ShortHash(*r.Hash()),
			r.Header.Root.Hex()),
	)

	result, err := evm.ExecuteOffChainCall(&callMsg, blockState, r.Header, rc.storage, rc.chainConfig, rc.logger)
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

// Updates the L1 and L2 chain heads, and returns the new L2 head hash and the produced rollup, if there is one.
func (rc *RollupChain) updateL1AndL2Heads(block *types.Block) (*common.L2RootHash, *core.Rollup, error) {
	// We extract the rollups from the block.
	rollupsInBlock := rc.bridge.ExtractRollups(block, rc.storage)

	// We determine whether this is a pre-genesis, genesis, or post-genesis block, and handle the block.
	blockStage, err := rc.getBlockStage(rollupsInBlock)
	if err != nil {
		return nil, nil, fmt.Errorf("could not determine block type. Cause: %w", err)
	}

	switch blockStage {
	case PreGenesis:
		return nil, nil, nil
	case Genesis:
		l2Head, err := rc.handleGenesisBlock(block, rollupsInBlock)
		return l2Head, nil, err
	case PostGenesis:
		return rc.handlePostGenesisBlock(block)
	}
	return nil, nil, fmt.Errorf("unrecognised block stage")
}

// Determines if this is a pre-genesis L2 block, the genesis L2 block, or a post-genesis L2 block.
func (rc *RollupChain) getBlockStage(rollupsInBlock []*core.Rollup) (BlockStage, error) {
	_, err := rc.storage.FetchGenesisRollup()
	if err != nil {
		if !errors.Is(err, errutil.ErrNotFound) {
			return -1, fmt.Errorf("could not retrieve genesis rollup. Cause: %w", err)
		}

		// If we haven't stored the genesis rollup and there are no rollups in this block, it cannot be the L2
		// genesis block.
		if len(rollupsInBlock) == 0 {
			return PreGenesis, nil
		}

		// If we haven't stored the genesis rollup before and this block contains rollups, it must be the L2 genesis
		// block.
		return Genesis, nil
	}

	return PostGenesis, nil
}

// We process the genesis block.
func (rc *RollupChain) handleGenesisBlock(block *types.Block, rollupsInBlock []*core.Rollup) (*common.L2RootHash, error) {
	// todo change this to a hardcoded hash on testnet/mainnet
	genesisRollup := rollupsInBlock[0]

	rc.logger.Info("Found genesis rollup", "l1Height", block.NumberU64(), "l1Hash", block.Hash())
	if err := rc.storage.StoreRollup(genesisRollup, nil); err != nil {
		return nil, fmt.Errorf("failed to store rollup. Cause: %w", err)
	}
	if err := rc.storage.StoreGenesisRollupHash(*genesisRollup.Hash()); err != nil {
		return nil, fmt.Errorf("could not store genesis rollup. Cause: %w", err)
	}
	if err := rc.storage.UpdateL2HeadForL1Block(block.Hash(), genesisRollup, nil); err != nil {
		return nil, fmt.Errorf("could not store new chain heads. Cause: %w", err)
	}
	if err := rc.storage.UpdateL1Head(block.Hash()); err != nil {
		return nil, fmt.Errorf("could not store new L1 head. Cause: %w", err)
	}
	if err := rc.faucet.CommitGenesisState(rc.storage); err != nil {
		return nil, fmt.Errorf("could not apply faucet preallocation. Cause: %w", err)
	}

	return genesisRollup.Hash(), nil
}

// This is where transactions are executed and the state is calculated.
// Obscuro includes a message bus embedded in the platform, and this method is responsible for transferring messages as well.
// The rollup can be a final rollup as received from peers or the rollup under construction.
func (rc *RollupChain) processState(rollup *core.Rollup, txs []*common.L2Tx, stateDB *state.StateDB) (common.L2RootHash, []*common.L2Tx, []*types.Receipt, []*types.Receipt) {
	var executedTransactions []*common.L2Tx
	var txReceipts []*types.Receipt

	txResults := evm.ExecuteTransactions(txs, stateDB, rollup.Header, rc.storage, rc.chainConfig, 0, rc.logger)
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
			rc.logger.Info(fmt.Sprintf("Excluding transaction %s from rollup r_%d. Cause: %s", tx.Hash().Hex(), common.ShortHash(*rollup.Hash()), result))
		}
	}

	// always process deposits last, either on top of the rollup produced speculatively or the newly created rollup
	// process deposits from the fromBlock of the parent to the current block (which is the fromBlock of the new rollup)
	parent, err := rc.storage.ParentRollup(rollup)
	if err != nil {
		rc.logger.Crit("Sanity check. Rollup has no parent.", log.ErrKey, err)
	}

	parentProof, err := rc.storage.Proof(parent)
	if err != nil {
		rc.logger.Crit(fmt.Sprintf("Could not retrieve a proof for rollup %s", rollup.Hash()), log.ErrKey, err)
	}
	rollupProof, err := rc.storage.Proof(rollup)
	if err != nil {
		rc.logger.Crit(fmt.Sprintf("Could not retrieve a proof for rollup %s", rollup.Hash()), log.ErrKey, err)
	}

	// TODO: Remove this depositing logic once the new bridge is added.
	depositTxs := rc.bridge.ExtractDeposits(
		parentProof,
		rollupProof,
		rc.storage,
		stateDB,
	)

	depositResponses := evm.ExecuteTransactions(depositTxs, stateDB, rollup.Header, rc.storage, rc.chainConfig, len(executedTransactions), rc.logger)
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

	messages := rc.crossChainProcessors.Local.RetrieveInboundMessages(parentProof, rollupProof, stateDB)
	transactions := rc.crossChainProcessors.Local.CreateSyntheticTransactions(messages, stateDB)
	syntheticTransactionsResponses := evm.ExecuteTransactions(transactions, stateDB, rollup.Header, rc.storage, rc.chainConfig, len(executedTransactions), rc.logger)
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
			res, err := evm.ExecuteOffChainCall(&txCallMessage, clonedDB, rollup.Header, rc.storage, rc.chainConfig, rc.logger)
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

func (rc *RollupChain) validateRollup(rollup *core.Rollup, rootHash common.L2RootHash, txReceipts []*types.Receipt, depositReceipts []*types.Receipt, stateDB *state.StateDB) bool {
	h := rollup.Header
	if !bytes.Equal(rootHash.Bytes(), h.Root.Bytes()) {
		dump := strings.Replace(string(stateDB.Dump(&state.DumpConfig{})), "\n", "", -1)

		rc.logger.Error(fmt.Sprintf("Verify rollup r_%d: Calculated a different state. This should not happen as there are no malicious actors yet. \nGot: %s\nExp: %s\nHeight:%d\nTxs:%v\nState: %s.\nDeposits: %+v",
			common.ShortHash(*rollup.Hash()), rootHash, h.Root, h.Number, core.PrintTxs(rollup.Transactions), dump, depositReceipts))
		return false
	}

	//  check that the withdrawals in the header match the withdrawals as calculated
	withdrawals := rc.bridge.RollupPostProcessingWithdrawals(rollup, stateDB, toReceiptMap(txReceipts))
	for i, w := range withdrawals {
		hw := h.Withdrawals[i]
		if hw.Amount.Cmp(w.Amount) != 0 || hw.Recipient != w.Recipient || hw.Contract != w.Contract {
			rc.logger.Error(fmt.Sprintf("Verify rollup r_%d: Withdrawals don't match", common.ShortHash(*rollup.Hash())))
			return false
		}
	}

	rec := allReceipts(txReceipts, depositReceipts)
	rbloom := types.CreateBloom(rec)
	if !bytes.Equal(rbloom.Bytes(), h.Bloom.Bytes()) {
		rc.logger.Error(fmt.Sprintf("Verify rollup r_%d: Invalid bloom (remote: %x  local: %x)", common.ShortHash(*rollup.Hash()), h.Bloom, rbloom))
		return false
	}

	receiptSha := types.DeriveSha(rec, trie.NewStackTrie(nil))
	if !bytes.Equal(receiptSha.Bytes(), h.ReceiptHash.Bytes()) {
		rc.logger.Error(fmt.Sprintf("Verify rollup r_%d: invalid receipt root hash (remote: %x local: %x)", common.ShortHash(*rollup.Hash()), h.ReceiptHash, receiptSha))
		return false
	}

	return true
}

// Calculates the state after processing the provided block.
func (rc *RollupChain) handlePostGenesisBlock(block *types.Block) (*common.L2RootHash, *core.Rollup, error) {
	// TODO - #718 - Cannot assume that the most recent rollup is on the previous block anymore. May be on the same block.
	// We retrieve the current L2 head.
	l2HeadHash, err := rc.storage.FetchHeadRollupForL1Block(block.ParentHash())
	if err != nil {
		return nil, nil, fmt.Errorf("could not retrieve current head rollup hash. Cause: %w", err)
	}
	l2Head, err := rc.storage.FetchRollup(*l2HeadHash)
	if err != nil {
		return nil, nil, fmt.Errorf("could not fetch parent rollup. Cause: %w", err)
	}

	// TODO - #718 - Validate any rollups in the block against the stored batches.

	// If we're the sequencer, we produce a new L2 head to replace the old one.
	var rollupTxReceipts types.Receipts
	if rc.nodeType == common.Sequencer {
		if l2Head, err = rc.produceRollup(block); err != nil {
			return nil, nil, fmt.Errorf("could not produce rollup. Cause: %w", err)
		}
		if err = rc.signRollup(l2Head); err != nil {
			return nil, nil, fmt.Errorf("could not sign rollup. Cause: %w", err)
		}
		if rollupTxReceipts, err = rc.checkRollup(l2Head); err != nil {
			return nil, nil, fmt.Errorf("could not check rollup. Cause: %w", err)
		}
		if err = rc.storage.StoreRollup(l2Head, rollupTxReceipts); err != nil {
			return nil, nil, fmt.Errorf("failed to store rollup. Cause: %w", err)
		}
	}

	// We update the chain heads.
	if err = rc.storage.UpdateL2HeadForL1Block(block.Hash(), l2Head, rollupTxReceipts); err != nil {
		return nil, nil, fmt.Errorf("could not store new head. Cause: %w", err)
	}
	if err = rc.storage.UpdateL1Head(block.Hash()); err != nil {
		return nil, nil, fmt.Errorf("could not store new L1 head. Cause: %w", err)
	}

	// We return the produced rollup, if we've produced one.
	var producedRollup *core.Rollup
	if rc.nodeType == common.Sequencer {
		producedRollup = l2Head
	}
	return l2Head.Hash(), producedRollup, nil
}

// verifies that the headers of the rollup match the results of executing the transactions
func (rc *RollupChain) checkRollup(rollup *core.Rollup) ([]*types.Receipt, error) {
	stateDB, err := rc.storage.CreateStateDB(rollup.Header.ParentHash)
	if err != nil {
		return nil, fmt.Errorf("could not create stateDB. Cause: %w", err)
	}

	// calculate the state to compare with what is in the Rollup
	rootHash, successfulTxs, txReceipts, depositReceipts := rc.processState(rollup, rollup.Transactions, stateDB)
	if len(successfulTxs) != len(rollup.Transactions) {
		return nil, fmt.Errorf("all transactions that are included in a rollup must be executed")
	}

	if !rc.validateRollup(rollup, rootHash, txReceipts, depositReceipts, stateDB) {
		return nil, fmt.Errorf("invalid rollup")
	}

	// todo - check that the transactions hash to the header.txHash

	// verify the signature
	if !rc.verifySig(rollup) {
		return nil, fmt.Errorf("invalid signature")
	}

	return txReceipts, nil
}

func (rc *RollupChain) signRollup(r *core.Rollup) error {
	var err error
	h := r.Hash()
	r.Header.R, r.Header.S, err = ecdsa.Sign(rand.Reader, rc.enclavePrivateKey, h[:])
	if err != nil {
		return fmt.Errorf("could not sign rollup. Cause: %w", err)
	}
	return nil
}

func (rc *RollupChain) verifySig(r *core.Rollup) bool {
	// If this rollup is generated by the current enclave skip the sig verification
	if bytes.Equal(r.Header.Agg.Bytes(), rc.hostID.Bytes()) {
		return true
	}

	h := r.Hash()
	if r.Header.R == nil || r.Header.S == nil {
		rc.logger.Error("Missing signature on rollup")
		return false
	}

	pubKey, err := rc.storage.FetchAttestedKey(r.Header.Agg)
	if err != nil {
		rc.logger.Error("Could not retrieve attested key for aggregator %s. Cause: %w", r.Header.Agg, err)
		return false
	}

	return ecdsa.Verify(pubKey, h[:], r.Header.R, r.Header.S)
}

// Retrieves the rollup with the given height, with special handling for earliest/latest/pending .
func (rc *RollupChain) getRollup(height gethrpc.BlockNumber) (*core.Rollup, error) {
	var rollup *core.Rollup
	switch height {
	case gethrpc.EarliestBlockNumber:
		genesisRollup, err := rc.storage.FetchGenesisRollup()
		if err != nil {
			return nil, fmt.Errorf("could not retrieve genesis rollup. Cause: %w", err)
		}
		rollup = genesisRollup
	case gethrpc.PendingBlockNumber:
		// TODO - Depends on the current pending rollup; leaving it for a different iteration as it will need more thought.
		return nil, fmt.Errorf("requested balance for pending block. This is not handled currently")
	case gethrpc.LatestBlockNumber:
		l2Head, err := rc.storage.FetchL2Head()
		if err != nil {
			return nil, fmt.Errorf("could not retrieve head state. Cause: %w", err)
		}
		rollup, err = rc.storage.FetchRollup(*l2Head)
		if err != nil {
			return nil, fmt.Errorf("rollup with requested height %d was not found. Cause: %w", height, err)
		}
	default:
		maybeRollup, err := rc.storage.FetchRollupByHeight(uint64(height))
		if err != nil {
			return nil, fmt.Errorf("rollup with requested height %d could not be retrieved. Cause: %w", height, err)
		}
		rollup = maybeRollup
	}
	return rollup, nil
}

// Creates a rollup.
func (rc *RollupChain) produceRollup(block *types.Block) (*core.Rollup, error) {
	headRollupHash, err := rc.storage.FetchHeadRollupForL1Block(block.ParentHash())
	if err != nil {
		return nil, fmt.Errorf("could not retrieve head rollup hash. Cause: %w", err)
	}
	headRollup, err := rc.storage.FetchRollup(*headRollupHash)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve head rollup. Cause: %w", err)
	}

	// These variables will be used to create the new rollup
	var newRollupTxs []*common.L2Tx
	var newRollupState *state.StateDB

	// Create a new rollup based on the fromBlock of inclusion of the previous, including all new transactions
	nonce := common.GenerateNonce()
	r, err := core.EmptyRollup(rc.hostID, headRollup.Header, block.Hash(), nonce)
	if err != nil {
		return nil, fmt.Errorf("could not create rollup. Cause: %w", err)
	}

	newRollupTxs, err = rc.mempool.CurrentTxs(headRollup, rc.storage)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve current transactions. Cause: %w", err)
	}

	newRollupState, err = rc.storage.CreateStateDB(r.Header.ParentHash)
	if err != nil {
		return nil, fmt.Errorf("could not create stateDB. Cause: %w", err)
	}

	rootHash, successfulTxs, txReceipts, depositReceipts := rc.processState(r, newRollupTxs, newRollupState)

	r.Header.Root = rootHash
	r.Transactions = successfulTxs

	// Postprocessing - withdrawals
	txReceiptsMap := toReceiptMap(txReceipts)
	r.Header.Withdrawals = rc.bridge.RollupPostProcessingWithdrawals(r, newRollupState, txReceiptsMap)
	crossChainMessages, err := rc.crossChainProcessors.Local.ExtractOutboundMessages(txReceipts)
	if err != nil {
		rc.logger.Crit("Extracting messages L2->L1 failed", err, log.CmpKey, log.CrossChainCmp)
	}

	r.Header.CrossChainMessages = crossChainMessages

	rc.logger.Trace(fmt.Sprintf("Added %d cross chain messages to rollup. Equivalent withdrawals in header - %d", len(r.Header.CrossChainMessages), len(r.Header.Withdrawals)), log.CmpKey, log.CrossChainCmp)

	crossChainBind, err := rc.storage.Proof(r)
	if err != nil {
		rc.logger.Crit("Failed to extract rollup proof that should exist!")
	}

	r.Header.LatestInboudCrossChainHash = crossChainBind.Hash()
	r.Header.LatestInboundCrossChainHeight = crossChainBind.Number()

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

	rc.logger.Trace("Create rollup.",
		"State", gethlog.Lazy{Fn: func() string {
			return strings.Replace(string(newRollupState.Dump(&state.DumpConfig{})), "\n", "", -1)
		}},
	)

	return r, nil
}

// Returns the state of the chain at height
// TODO make this cacheable
func (rc *RollupChain) getChainStateAtBlock(blockNumber *gethrpc.BlockNumber) (*state.StateDB, error) {
	// We retrieve the rollup of interest.
	rollup, err := rc.getRollup(*blockNumber)
	if err != nil {
		return nil, err
	}

	// We get that of the chain at that height
	blockchainState, err := rc.storage.CreateStateDB(*rollup.Hash())
	if err != nil {
		return nil, fmt.Errorf("could not create stateDB. Cause: %w", err)
	}

	if blockchainState == nil {
		return nil, fmt.Errorf("unable to fetch chain state for rollup %s", rollup.Hash().Hex())
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
