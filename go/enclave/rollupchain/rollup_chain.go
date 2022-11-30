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

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/obscuronet/go-obscuro/go/common/gethapi"
	"github.com/status-im/keycard-go/hexutils"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	gethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/errutil"
	"github.com/obscuronet/go-obscuro/go/common/gethutil"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/enclave/bridge"
	obscurocore "github.com/obscuronet/go-obscuro/go/enclave/core"
	"github.com/obscuronet/go-obscuro/go/enclave/crypto"
	"github.com/obscuronet/go-obscuro/go/enclave/db"
	"github.com/obscuronet/go-obscuro/go/enclave/events"
	"github.com/obscuronet/go-obscuro/go/enclave/evm"
	"github.com/obscuronet/go-obscuro/go/enclave/mempool"
)

var (
	errBlockAlreadyProcessed  = errors.New("block already processed")
	errBlockAncestorNotFound  = errors.New("block ancestor not found")
	errIsPreGenesis           = errors.New("genesis rollup has not yet been received")
	errIsGenesisRollupInBlock = errors.New("block contains genesis rollup")
)

// RollupChain represents the canonical chain, and manages the state.
type RollupChain struct {
	hostID      gethcommon.Address
	nodeType    common.NodeType
	chainConfig *params.ChainConfig

	storage               db.Storage
	l1Blockchain          *core.BlockChain
	bridge                *bridge.Bridge
	transactionBlobCrypto crypto.TransactionBlobCrypto // todo - remove
	mempool               mempool.Manager
	faucet                Faucet
	subscriptionManager   *events.SubscriptionManager

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
	l1Blockchain *core.BlockChain,
	bridge *bridge.Bridge,
	subscriptionManager *events.SubscriptionManager,
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
		subscriptionManager:   subscriptionManager,
		enclavePrivateKey:     privateKey,
		chainConfig:           chainConfig,
		blockProcessingMutex:  sync.Mutex{},
		logger:                logger,
		GlobalGasCap:          5_000_000_000,
		BaseFee:               gethcommon.Big0,
	}
}

// ProduceNewRollup creates a rollup, signs it, checks it, and stores it.
func (rc *RollupChain) ProduceNewRollup(blockHash *common.L1RootHash) (*common.ExtRollup, error) {
	block, err := rc.storage.FetchBlock(*blockHash)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve block to produce rollup. Cause: %w", err)
	}
	blockState, err := rc.storage.FetchHeadsAfterL1Block(*blockHash)
	if err != nil {
		if errors.Is(err, errutil.ErrNotFound) {
			// not an error state, we ingested a block but no rollup head found
			return nil, nil //nolint:nilnil
		}
		return nil, fmt.Errorf("could not retrieve block state. Cause: %w", err)
	}

	rollup, err := rc.produceRollup(block, blockState)
	if err != nil {
		return nil, fmt.Errorf("could not produce rollup. Cause: %w", err)
	}

	err = rc.signRollup(rollup)
	if err != nil {
		return nil, fmt.Errorf("could not sign rollup. Cause: %w", err)
	}

	// Sanity check the produced rollup
	_, err = rc.checkRollup(rollup)
	if err != nil {
		return nil, err
	}

	err = rc.storage.StoreRollup(rollup)
	if err != nil {
		return nil, err
	}

	// TODO - #718 - This rollup should be stored as the new head.

	extRollup := rollup.ToExtRollup(rc.transactionBlobCrypto)
	rc.logger.Trace(fmt.Sprintf("Processed block: b_%d (%d). Produced rollup r_%d",
		common.ShortHash(block.Hash()), block.NumberU64(), common.ShortHash(extRollup.Hash())))

	return &extRollup, nil
}

// ProduceGenesisRollup creates a genesis rollup linked to the provided L1 block and signs it.
func (rc *RollupChain) ProduceGenesisRollup(blkHash gethcommon.Hash) (*obscurocore.Rollup, error) {
	preFundGenesisState, err := rc.faucet.GetGenesisRoot(rc.storage)
	if err != nil {
		return nil, err
	}

	rolGenesis := obscurocore.NewRollup(
		blkHash,
		nil,
		common.L2GenesisHeight,
		gethcommon.HexToAddress("0x0"),
		[]*common.L2Tx{},
		[]common.Withdrawal{},
		common.GenerateNonce(),
		*preFundGenesisState,
	)

	err = rc.signRollup(rolGenesis)
	if err != nil {
		return nil, fmt.Errorf("could not sign genesis rollup. Cause: %w", err)
	}

	return rolGenesis, nil
}

// ProcessL1Block is used to update the enclave with an additional L1 block.
func (rc *RollupChain) ProcessL1Block(block types.Block, isLatest bool) (*common.BlockSubmissionResponse, error) {
	rc.blockProcessingMutex.Lock()
	defer rc.blockProcessingMutex.Unlock()

	// We check whether we've already processed the block.
	_, err := rc.storage.FetchBlock(block.Hash())
	if err == nil {
		return nil, rc.rejectBlockErr(errBlockAlreadyProcessed)
	}
	if !errors.Is(err, errutil.ErrNotFound) {
		return nil, fmt.Errorf("could not retrieve block. Cause: %w", err)
	}

	// We insert the block into the L1 chain and store it.
	ingestionType, err := rc.insertBlockIntoL1Chain(&block, isLatest)
	if err != nil {
		// Do not store the block if the L1 chain insertion failed
		return nil, rc.rejectBlockErr(err)
	}
	rc.logger.Trace("block inserted successfully",
		"height", block.NumberU64(), "hash", block.Hash(), "ingestionType", ingestionType)
	rc.storage.StoreBlock(&block)

	rc.logger.Trace(fmt.Sprintf("Update state: b_%d", common.ShortHash(block.Hash())))
	headsAfterL1Block, err := rc.updateHeads(&block)
	if err != nil {
		return nil, rc.rejectBlockErr(err)
	}
	return rc.produceBlockSubmissionResponse(&block, headsAfterL1Block)
}

func (rc *RollupChain) GetBalance(accountAddress gethcommon.Address, blockNumber gethrpc.BlockNumber) (*gethcommon.Address, *hexutil.Big, error) {
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
func (rc *RollupChain) GetBalanceAtBlock(accountAddr gethcommon.Address, blockNumber gethrpc.BlockNumber) (*hexutil.Big, error) {
	chainState, err := rc.getChainStateAtBlock(blockNumber)
	if err != nil {
		return nil, fmt.Errorf("unable to get blockchain state - %w", err)
	}

	return (*hexutil.Big)(chainState.GetBalance(accountAddr)), nil
}

// ExecuteOffChainTransaction executes non-state changing transactions at a given block height (eth_call)
func (rc *RollupChain) ExecuteOffChainTransaction(apiArgs *gethapi.TransactionArgs) (*core.ExecutionResult, error) {
	// TODO Hook up the blockNumber
	result, err := rc.ExecuteOffChainTransactionAtBlock(apiArgs, gethrpc.BlockNumber(0))
	if err != nil {
		rc.logger.Error(fmt.Sprintf("!OffChain: Failed to execute contract %s.", apiArgs.To), log.ErrKey, err.Error())
		return nil, err
	}

	// the execution might have succeeded but the evm contract logic might have failed
	if result.Failed() {
		rc.logger.Error(fmt.Sprintf("!OffChain: Failed to execute contract %s.", apiArgs.To), log.ErrKey, result.Err)
		return nil, result.Err
	}

	rc.logger.Trace(fmt.Sprintf("!OffChain result: %s", hexutils.BytesToHex(result.ReturnData)))

	return result, nil
}

func (rc *RollupChain) ExecuteOffChainTransactionAtBlock(apiArgs *gethapi.TransactionArgs, blockNumber gethrpc.BlockNumber) (*core.ExecutionResult, error) {
	// TODO review this during gas mechanics implementation
	callMsg, err := apiArgs.ToMessage(rc.GlobalGasCap, rc.BaseFee)
	if err != nil {
		return nil, fmt.Errorf("unable to convert TransactionArgs to Message - %w", err)
	}

	headsAfterL1Block, err := rc.storage.FetchCurrentHeadsAfterL1Block()
	if err != nil {
		return nil, fmt.Errorf("unable to fetch head state. Cause: %w", err)
	}
	// todo - get the parent
	r, err := rc.storage.FetchRollup(headsAfterL1Block.HeadRollup)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch head state rollup. Cause: %w", err)
	}

	rc.logger.Trace(fmt.Sprintf("!OffChain call: contractAddress=%s, from=%s, data=%s, rollup=r_%d, state=%s", callMsg.To(), callMsg.From(), hexutils.BytesToHex(callMsg.Data()), common.ShortHash(r.Hash()), r.Header.Root.Hex()))
	s, err := rc.storage.CreateStateDB(headsAfterL1Block.HeadRollup)
	if err != nil {
		return nil, fmt.Errorf("could not create stateDB. Cause: %w", err)
	}

	result, err := evm.ExecuteOffChainCall(&callMsg, s, r.Header, rc.storage, rc.chainConfig, rc.logger)
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
			return nil, errBlockAncestorNotFound
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

func (rc *RollupChain) produceBlockSubmissionResponse(block *types.Block, headsAfterL1Block *obscurocore.HeadsAfterL1Block) (*common.BlockSubmissionResponse, error) {
	if headsAfterL1Block == nil {
		// not an error state, we ingested a block but no rollup head found
		return &common.BlockSubmissionResponse{
			UpdatedHeadRollup: false,
		}, nil
	}

	headRollup, err := rc.storage.FetchRollup(headsAfterL1Block.HeadRollup)
	if err != nil {
		rc.logger.Crit("Could not fetch rollup", log.ErrKey, err)
	}
	var ingestedRollupHeader *common.Header
	if headsAfterL1Block.UpdatedHeadRollup {
		ingestedRollupHeader = headRollup.Header
	}

	return &common.BlockSubmissionResponse{
		UpdatedHeadRollup:    headsAfterL1Block.UpdatedHeadRollup,
		IngestedRollupHeader: ingestedRollupHeader,
		SubscribedLogs:       rc.getEncryptedLogs(*block, headsAfterL1Block),
	}, nil
}

// Recursively calculates and stores the block state, receipts and logs for the given block.
func (rc *RollupChain) updateHeads(block *types.Block) (*obscurocore.HeadsAfterL1Block, error) {
	// This method is called recursively in case of re-orgs. Stop when state was calculated already.
	headsAfterL1Block, err := rc.storage.FetchHeadsAfterL1Block(block.Hash())
	if err != nil && !errors.Is(err, errutil.ErrNotFound) {
		return nil, fmt.Errorf("could not retrieve block state. Cause: %w", err)
	}
	if err == nil {
		// The state has already been calculated, so we return it.
		return headsAfterL1Block, nil
	}

	rollupsInBlock := rc.bridge.ExtractRollups(block, rc.storage)

	// Detect if the incoming block contains the genesis rollup, and generate an updated state.
	// Handles the case of the block containing the genesis being processed multiple times.
	headsAfterGenesisRollup, err := rc.handleGenesisRollup(block, rollupsInBlock)
	if err != nil {
		if errors.Is(err, errIsPreGenesis) {
			// We wait for the genesis rollup.
			return nil, nil //nolint:nilnil
		}
		if errors.Is(err, errIsGenesisRollupInBlock) {
			// We can return the genesis state immediately.
			return headsAfterGenesisRollup, nil
		}
		return nil, fmt.Errorf("could not handle genesis rollup. Cause: %w", err)
	}

	// To calculate the state after the current block, we need the state after the parent.
	// If this point is reached, there is a parent state guaranteed, because the genesis is handled above
	headsAfterParentBlock, err := rc.getHeadsAfterParentBlock(block.ParentHash())
	if err != nil {
		return nil, fmt.Errorf("could not retrieve heads after parent block. Cause: %w", err)
	}

	headsAfterL1Block, headRollup, receipts, err := rc.calculateNewHeads(block, headsAfterParentBlock, rollupsInBlock)
	if err != nil {
		return nil, fmt.Errorf("could not calculate heads after L1 block. Cause: %w", err)
	}

	rc.logger.Trace(fmt.Sprintf("Calc block state b_%d: Found: %t - r_%d, ",
		common.ShortHash(block.Hash()), headsAfterL1Block.UpdatedHeadRollup, common.ShortHash(headsAfterL1Block.HeadRollup)))

	err = rc.storage.StoreNewHeads(headsAfterL1Block, headRollup, receipts)
	if err != nil {
		return nil, fmt.Errorf("could not store new head. Cause: %w", err)
	}

	return headsAfterL1Block, nil
}

func (rc *RollupChain) handleGenesisRollup(b *types.Block, rollupsInBlock []*obscurocore.Rollup) (*obscurocore.HeadsAfterL1Block, error) {
	genesisRollup, err := rc.storage.FetchGenesisRollup()
	if err != nil {
		// If there is no genesis yet and no rollups have arrived, there is nothing to do
		if errors.Is(err, errutil.ErrNotFound) {
			if len(rollupsInBlock) == 0 {
				return nil, errIsPreGenesis
			}
		} else {
			return nil, fmt.Errorf("could not retrieve genesis rollup. Cause: %w", err)
		}
	}

	// We've already found the genesis rollup.
	if genesisRollup != nil {
		// Re-processing the block that contains the genesis rollup. This can happen as blocks can be fed to the enclave
		// multiple times. We don't update the state and move on.
		if len(rollupsInBlock) == 1 && bytes.Equal(rollupsInBlock[0].Header.Hash().Bytes(), genesisRollup.Hash().Bytes()) {
			return nil, errIsGenesisRollupInBlock
		}
		// We return - we do not need to handle the genesis rollup, since we already have.
		return nil, nil //nolint:nilnil
	}

	// The incoming block holds the genesis rollup. Calculate and return the new block state.
	// todo change this to an hardcoded hash on testnet/mainnet
	rc.logger.Info("Found genesis rollup", "l1Height", b.NumberU64(), "l1Hash", b.Hash())
	genesis := rollupsInBlock[0]
	err = rc.storage.StoreGenesisRollup(genesis)
	if err != nil {
		return nil, fmt.Errorf("could not store genesis rollup. Cause: %w", err)
	}

	// The genesis rollup is part of the canonical chain and will be included in an L1 block by the first Aggregator.
	headsAfterL1Block := obscurocore.HeadsAfterL1Block{
		HeadBlock:         b.Hash(),
		HeadRollup:        genesis.Hash(),
		UpdatedHeadRollup: true,
	}
	err = rc.storage.StoreNewHeads(&headsAfterL1Block, genesis, nil)
	if err != nil {
		return nil, fmt.Errorf("could not store new chain heads. Cause: %w", err)
	}

	_, err = rc.faucet.CommitGenesisState(rc.storage)
	if err != nil {
		return nil, fmt.Errorf("could not apply faucet preallocation. Cause: %w", err)
	}

	return &headsAfterL1Block, errIsGenesisRollupInBlock
}

func (rc *RollupChain) getHeadsAfterParentBlock(parentBlockHash gethcommon.Hash) (*obscurocore.HeadsAfterL1Block, error) {
	headsAfterParentBlock, err := rc.storage.FetchHeadsAfterL1Block(parentBlockHash)
	if err != nil {
		if !errors.Is(err, errutil.ErrNotFound) {
			return nil, fmt.Errorf("could not retrieve parent block state. Cause: %w", err)
		}

		// We don't have the state after the parent block stored. We recursively calculate it.
		parentBlock, err := rc.storage.FetchBlock(parentBlockHash)
		if err != nil {
			return nil, fmt.Errorf("could not retrieve parent block when calculating block state. Cause: %w", err)
		}
		headsAfterParentBlock, err = rc.updateHeads(parentBlock)
		if err != nil {
			return nil, err
		}
	}

	if headsAfterParentBlock == nil {
		return nil, fmt.Errorf("could not calculate parent block state when calculating block state. Cause: %w", err)
	}
	return headsAfterParentBlock, nil
}

// This is where transactions are executed and the state is calculated.
// Obscuro includes a bridge embedded in the platform, and this method is responsible for processing deposits as well.
// The rollup can be a final rollup as received from peers or the rollup under construction.
func (rc *RollupChain) processState(rollup *obscurocore.Rollup, txs []*common.L2Tx, stateDB *state.StateDB) (gethcommon.Hash, []*common.L2Tx, []*types.Receipt, []*types.Receipt) {
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
			rc.logger.Info(fmt.Sprintf("Excluding transaction %s from rollup r_%d. Cause: %s", tx.Hash().Hex(), common.ShortHash(rollup.Hash()), result))
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

	rootHash, err := stateDB.Commit(true)
	if err != nil {
		rc.logger.Crit("could not commit to state DB. ", log.ErrKey, err)
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
		rc.logger.Error(fmt.Sprintf("Verify rollup r_%d: Calculated a different state. This should not happen as there are no malicious actors yet. \nGot: %s\nExp: %s\nHeight:%d\nTxs:%v\nState: %s.\nDeposits: %+v",
			common.ShortHash(rollup.Hash()), rootHash, h.Root, h.Number, obscurocore.PrintTxs(rollup.Transactions), dump, depositReceipts))
		return false
	}

	//  check that the withdrawals in the header match the withdrawals as calculated
	withdrawals := rc.bridge.RollupPostProcessingWithdrawals(rollup, stateDB, toReceiptMap(txReceipts))
	for i, w := range withdrawals {
		hw := h.Withdrawals[i]
		if hw.Amount.Cmp(w.Amount) != 0 || hw.Recipient != w.Recipient || hw.Contract != w.Contract {
			rc.logger.Error(fmt.Sprintf("Verify rollup r_%d: Withdrawals don't match", common.ShortHash(rollup.Hash())))
			return false
		}
	}

	rec := allReceipts(txReceipts, depositReceipts)
	rbloom := types.CreateBloom(rec)
	if !bytes.Equal(rbloom.Bytes(), h.Bloom.Bytes()) {
		rc.logger.Error(fmt.Sprintf("Verify rollup r_%d: Invalid bloom (remote: %x  local: %x)", common.ShortHash(rollup.Hash()), h.Bloom, rbloom))
		return false
	}

	receiptSha := types.DeriveSha(rec, trie.NewStackTrie(nil))
	if !bytes.Equal(receiptSha.Bytes(), h.ReceiptHash.Bytes()) {
		rc.logger.Error(fmt.Sprintf("Verify rollup r_%d: invalid receipt root hash (remote: %x local: %x)", common.ShortHash(rollup.Hash()), h.ReceiptHash, receiptSha))
		return false
	}

	return true
}

// given an L1 block, and the State as it was in the Parent block, calculates the State after the current block.
func (rc *RollupChain) calculateNewHeads(
	block *types.Block,
	headsAfterParentBlock *obscurocore.HeadsAfterL1Block,
	rollupsInBlock []*obscurocore.Rollup,
) (*obscurocore.HeadsAfterL1Block, *obscurocore.Rollup, []*types.Receipt, error) {
	currentHeadRollup, err := rc.storage.FetchRollup(headsAfterParentBlock.HeadRollup)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("could not fetch parent rollup. Cause: %w", err)
	}

	newHeadRollup, found := selectNextRollup(currentHeadRollup, rollupsInBlock, rc.storage)
	var rollupTxReceipts []*types.Receipt
	if found {
		rollupTxReceipts, err = rc.checkRollup(newHeadRollup)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to check rollup. Cause: %w", err)
		}
	} else {
		newHeadRollup = currentHeadRollup
	}

	// TODO - #718 - Instead of updating the rollup head, we should validate the stored batches against the winning
	//  rollup. We should still update the block head.

	headsAfterL1Block := obscurocore.HeadsAfterL1Block{
		HeadBlock:         block.Hash(),
		HeadRollup:        newHeadRollup.Hash(),
		UpdatedHeadRollup: found,
	}
	return &headsAfterL1Block, newHeadRollup, rollupTxReceipts, nil
}

// verifies that the headers of the rollup match the results of executing the transactions
func (rc *RollupChain) checkRollup(rollup *obscurocore.Rollup) ([]*types.Receipt, error) {
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

func (rc *RollupChain) signRollup(r *obscurocore.Rollup) error {
	var err error
	h := r.Hash()
	r.Header.R, r.Header.S, err = ecdsa.Sign(rand.Reader, rc.enclavePrivateKey, h[:])
	if err != nil {
		return fmt.Errorf("could not sign rollup. Cause: %w", err)
	}
	return nil
}

func (rc *RollupChain) verifySig(r *obscurocore.Rollup) bool {
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
func (rc *RollupChain) getRollup(height gethrpc.BlockNumber) (*obscurocore.Rollup, error) {
	var rollup *obscurocore.Rollup
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
		headsAfterL1Block, err := rc.storage.FetchCurrentHeadsAfterL1Block()
		if err != nil {
			return nil, fmt.Errorf("could not retrieve head state. Cause: %w", err)
		}
		rollup, err = rc.storage.FetchRollup(headsAfterL1Block.HeadRollup)
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

// Retrieves and encrypts the logs for the block.
func (rc *RollupChain) getEncryptedLogs(block types.Block, headsAfterL1Block *obscurocore.HeadsAfterL1Block) map[gethrpc.ID][]byte {
	logs := []*types.Log{}
	fetchedLogs, err := rc.storage.FetchLogs(block.Hash())
	if err == nil {
		logs = fetchedLogs
	} else {
		rc.logger.Error("Could not retrieve logs for stored block state; returning no logs. Cause: %w", err)
	}
	encryptedLogs, err := rc.subscriptionManager.GetSubscribedLogsEncrypted(logs, headsAfterL1Block.HeadRollup)
	if err != nil {
		rc.logger.Crit("Could not get subscribed logs in encrypted form. ", log.ErrKey, err)
	}
	return encryptedLogs
}

// Creates a rollup.
func (rc *RollupChain) produceRollup(b *types.Block, bs *obscurocore.HeadsAfterL1Block) (*obscurocore.Rollup, error) {
	headRollup, err := rc.storage.FetchRollup(bs.HeadRollup)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve head rollup. Cause: %w", err)
	}

	// These variables will be used to create the new rollup
	var newRollupTxs []*common.L2Tx
	var newRollupState *state.StateDB

	// Create a new rollup based on the fromBlock of inclusion of the previous, including all new transactions
	nonce := common.GenerateNonce()
	r, err := obscurocore.EmptyRollup(rc.hostID, headRollup.Header, b.Hash(), nonce)
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

func (rc *RollupChain) rejectBlockErr(cause error) *common.BlockRejectError {
	var hash gethcommon.Hash
	l1Head, err := rc.storage.FetchHeadBlock()
	// TODO - Handle error.
	if err == nil {
		hash = l1Head.Hash()
	}
	return &common.BlockRejectError{
		L1Head:  hash,
		Wrapped: cause,
	}
}

// Returns the state of the chain at height
// TODO make this cacheable
func (rc *RollupChain) getChainStateAtBlock(blockNumber gethrpc.BlockNumber) (*state.StateDB, error) {
	// We retrieve the rollup of interest.
	rollup, err := rc.getRollup(blockNumber)
	if err != nil {
		return nil, err
	}

	// We get that of the chain at that height
	blockchainState, err := rc.storage.CreateStateDB(rollup.Hash())
	if err != nil {
		return nil, fmt.Errorf("could not create stateDB. Cause: %w", err)
	}

	if blockchainState == nil {
		return nil, fmt.Errorf("unable to fetch chain state for rollup %s", rollup.Hash().Hex())
	}

	return blockchainState, err
}

// Returns the whether the account is a contract or not at a certain height
func (rc *RollupChain) isAccountContractAtBlock(accountAddr gethcommon.Address, blockNumber gethrpc.BlockNumber) (bool, error) {
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
