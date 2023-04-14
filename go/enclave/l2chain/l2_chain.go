package l2chain

import (
	"errors"
	"fmt"
	"sort"

	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/errutil"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/enclave/core"
	"github.com/obscuronet/go-obscuro/go/enclave/crosschain"
	"github.com/obscuronet/go-obscuro/go/enclave/db"
	"github.com/obscuronet/go-obscuro/go/enclave/evm"
	"github.com/obscuronet/go-obscuro/go/enclave/genesis"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
)

// ObscuroChain represents the canonical L2 chain, and manages the state.
type ObscuroChain struct {
	chainConfig *params.ChainConfig

	storage              db.Storage
	genesis              *genesis.Genesis
	crossChainProcessors *crosschain.Processors

	logger gethlog.Logger
}

func New(
	storage db.Storage,
	crossChainProcessors *crosschain.Processors,
	chainConfig *params.ChainConfig,
	genesis *genesis.Genesis,
	logger gethlog.Logger,
) *ObscuroChain {
	return &ObscuroChain{
		storage:              storage,
		crossChainProcessors: crossChainProcessors,
		chainConfig:          chainConfig,
		logger:               logger,
		genesis:              genesis,
	}
}

// This is where transactions are executed and the state is calculated.
// Obscuro includes a message bus embedded in the platform, and this method is responsible for transferring messages as well.
// The batch can be a final batch as received from peers or the batch under construction.
func (oc *ObscuroChain) processState(batch *core.Batch, txs []*common.L2Tx, stateDB *state.StateDB) (common.L2BatchHash, []*common.L2Tx, []*types.Receipt, []*types.Receipt) {
	var executedTransactions []*common.L2Tx
	var txReceipts []*types.Receipt

	txResults := evm.ExecuteTransactions(txs, stateDB, batch.Header, oc.storage, oc.chainConfig, 0, oc.logger)
	for _, tx := range txs {
		result, f := txResults[tx.Hash()]
		if !f {
			oc.logger.Crit("There should be an entry for each transaction ")
		}
		rec, foundReceipt := result.(*types.Receipt)
		if foundReceipt {
			executedTransactions = append(executedTransactions, tx)
			txReceipts = append(txReceipts, rec)
		} else {
			// Exclude all errors
			oc.logger.Info(fmt.Sprintf("Excluding transaction %s from batch b_%d. Cause: %s", tx.Hash().Hex(), common.ShortHash(*batch.Hash()), result))
		}
	}

	// always process deposits last, either on top of the rollup produced speculatively or the newly created rollup
	// process deposits from the fromBlock of the parent to the current block (which is the fromBlock of the new rollup)
	parent, err := oc.storage.FetchBatch(batch.Header.ParentHash)
	if err != nil {
		oc.logger.Crit("Sanity check. Rollup has no parent.", log.ErrKey, err)
	}

	parentProof, err := oc.storage.FetchBlock(parent.Header.L1Proof)
	if err != nil {
		oc.logger.Crit(fmt.Sprintf("Could not retrieve a proof for batch %s", batch.Hash()), log.ErrKey, err)
	}
	batchProof, err := oc.storage.FetchBlock(batch.Header.L1Proof)
	if err != nil {
		oc.logger.Crit(fmt.Sprintf("Could not retrieve a proof for batch %s", batch.Hash()), log.ErrKey, err)
	}

	messages := oc.crossChainProcessors.Local.RetrieveInboundMessages(parentProof, batchProof, stateDB)
	transactions := oc.crossChainProcessors.Local.CreateSyntheticTransactions(messages, stateDB)
	syntheticTransactionsResponses := evm.ExecuteTransactions(transactions, stateDB, batch.Header, oc.storage, oc.chainConfig, len(executedTransactions), oc.logger)
	synthReceipts := make([]*types.Receipt, len(syntheticTransactionsResponses))
	if len(syntheticTransactionsResponses) != len(transactions) {
		oc.logger.Crit("Sanity check. Some synthetic transactions failed.")
	}

	i := 0
	for _, resp := range syntheticTransactionsResponses {
		rec, ok := resp.(*types.Receipt)
		if !ok { // Еxtract reason for failing deposit.
			// todo (#1578) - handle the case of an error (e.g. insufficient funds)
			oc.logger.Crit("Sanity check. Expected a receipt", log.ErrKey, resp)
		}

		if rec.Status == 0 { // Synthetic transactions should not fail. In case of failure get the revert reason.
			failingTx := transactions[i]
			txCallMessage := types.NewMessage(
				oc.crossChainProcessors.Local.GetOwner(),
				failingTx.To(),
				stateDB.GetNonce(oc.crossChainProcessors.Local.GetOwner()),
				failingTx.Value(),
				failingTx.Gas(),
				gethcommon.Big0,
				gethcommon.Big0,
				gethcommon.Big0,
				failingTx.Data(),
				failingTx.AccessList(),
				false)

			clonedDB := stateDB.Copy()
			res, err := evm.ExecuteObsCall(&txCallMessage, clonedDB, batch.Header, oc.storage, oc.chainConfig, oc.logger)
			oc.logger.Crit("Synthetic transaction failed!", log.ErrKey, err, "result", res)
		}

		synthReceipts[i] = rec
		i++
	}

	rootHash, err := stateDB.Commit(true)
	if err != nil {
		oc.logger.Crit("could not commit to state DB. ", log.ErrKey, err)
	}

	sort.Sort(sortByTxIndex(txReceipts))

	return rootHash, executedTransactions, txReceipts, synthReceipts
}

// ResyncStateDB can be called to ensure stateDB data is available for the canonical L2 batch chain
// After an (ungraceful) shutdown this method must be called to rebuild the stateDB data based on the persisted batches
func (oc *ObscuroChain) ResyncStateDB() error {
	batch, err := oc.storage.FetchHeadBatch()
	if err != nil {
		if errors.Is(err, errutil.ErrNotFound) {
			// there is no head batch, this is probably a new node - there is no state to rebuild
			oc.logger.Info("no head batch found in DB after restart", log.ErrKey, err)
			return nil
		}
		return fmt.Errorf("unexpected error fetching head batch to resync- %w", err)
	}
	if !stateDBAvailableForBatch(oc.storage, batch.Hash()) {
		oc.logger.Info("state not available for latest batch after restart - rebuilding stateDB cache from batches")
		err = oc.replayBatchesToValidState()
		if err != nil {
			return fmt.Errorf("unable to replay batches to restore valid state - %w", err)
		}
	}
	return nil
}

// replayBatchesToValidState is used to repopulate the stateDB cache with data from persisted batches. Two step process:
// 1. step backwards from head batch until we find a batch that is already in stateDB cache, builds list of batches to replay
// 2. iterate that list of batches from the earliest, process the transactions to calculate and cache the stateDB
// todo (#1416) - get unit test coverage around this (and L2 Chain code more widely, see ticket #1416 )
func (oc *ObscuroChain) replayBatchesToValidState() error {
	// this slice will be a stack of batches to replay as we walk backwards in search of latest valid state
	// todo - consider capping the size of this batch list using FIFO to avoid memory issues, and then repeating as necessary
	var batchesToReplay []*core.Batch
	// `batchToReplayFrom` variable will eventually be the latest batch for which we are able to produce a StateDB
	// - we will then set that as the head of the L2 so that this node can rebuild its missing state
	batchToReplayFrom, err := oc.storage.FetchHeadBatch()
	if err != nil {
		return fmt.Errorf("no head batch found in DB but expected to replay batches - %w", err)
	}
	// loop backwards building a slice of all batches that don't have cached stateDB data available
	for !stateDBAvailableForBatch(oc.storage, batchToReplayFrom.Hash()) {
		batchesToReplay = append(batchesToReplay, batchToReplayFrom)
		if batchToReplayFrom.NumberU64() == 0 {
			// no more parents to check, replaying from genesis
			break
		}
		batchToReplayFrom, err = oc.storage.FetchBatch(batchToReplayFrom.Header.ParentHash)
		if err != nil {
			return fmt.Errorf("unable to fetch previous batch while rolling back to stable state - %w", err)
		}
	}
	oc.logger.Info("replaying batch data into stateDB cache", "fromBatch", batchesToReplay[len(batchesToReplay)-1].NumberU64(),
		"toBatch", batchesToReplay[0].NumberU64())
	// loop through the slice of batches without stateDB data (starting with the oldest) and reprocess them to update cache
	for i := len(batchesToReplay) - 1; i >= 0; i-- {
		batch := batchesToReplay[i]

		// if genesis batch then create the genesis state before continuing on with remaining batches
		if batch.NumberU64() == 0 {
			err := oc.genesis.CommitGenesisState(oc.storage)
			if err != nil {
				return err
			}
			continue
		}

		prevState, err := oc.storage.CreateStateDB(batch.Header.ParentHash)
		if err != nil {
			return err
		}
		// we don't need the return values, just want the post-batch state to be cached
		oc.processState(batch, batch.Transactions, prevState)
	}

	return nil
}

// The enclave caches a stateDB instance against each batch hash, this is the input state when producing the following
// batch in the chain and is used to query state at a certain height.
//
// This method checks if the stateDB data is available for a given batch hash (so it can be restored if not)
func stateDBAvailableForBatch(storage db.Storage, hash *common.L2BatchHash) bool {
	_, err := storage.CreateStateDB(*hash)
	return err == nil
}

type sortByTxIndex []*types.Receipt

func (c sortByTxIndex) Len() int           { return len(c) }
func (c sortByTxIndex) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c sortByTxIndex) Less(i, j int) bool { return c[i].TransactionIndex < c[j].TransactionIndex }
