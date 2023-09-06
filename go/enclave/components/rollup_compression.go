package components

import (
	"errors"
	"fmt"
	"math/big"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/compression"
	"github.com/obscuronet/go-obscuro/go/common/errutil"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/enclave/core"
	"github.com/obscuronet/go-obscuro/go/enclave/crypto"
	"github.com/obscuronet/go-obscuro/go/enclave/storage"
)

/*
RollupCompression - responsible for the compression logic

## Problem
The main overhead (after the tx payloads), are the batch headers.

## Requirements:
1. recreate the exact batch headers as the live ones
2. security - make sure it all chains up cryptographically, so it can't be gamed

## Solution elements:
1. Add another compressed and encrypted metadata blob to the ExtRollup. This is the "CalldataRollupHeader".
It will also be published as calldata, together with the transactions. The role of this header is to contain the bare minimum
information required to recreate the batches.
2. Use implicit position information, deltas, and exceptions to minimise size.
Eg. If the time between 2 batches is always 1second, there is no need to store any extra information.
3. To avoid storing hashes, which don't compress at all, we execute each batch to be able to populate the parent hash.
4. The Signatures over the batches are not stored, since the rollup is itself signed.
5. The cross chain messages are calculated.
*/
type RollupCompression struct {
	dataEncryptionService  crypto.DataEncryptionService
	dataCompressionService compression.DataCompressionService
	batchRegistry          BatchRegistry
	batchExecutor          BatchExecutor
	logger                 gethlog.Logger
	storage                storage.Storage
}

func NewRollupCompression(
	batchRegistry BatchRegistry,
	batchExecutor BatchExecutor,
	dataEncryptionService crypto.DataEncryptionService,
	dataCompressionService compression.DataCompressionService,
	storage storage.Storage,
	logger gethlog.Logger,
) *RollupCompression {
	return &RollupCompression{
		batchRegistry:          batchRegistry,
		batchExecutor:          batchExecutor,
		dataEncryptionService:  dataEncryptionService,
		dataCompressionService: dataCompressionService,
		logger:                 logger,
		storage:                storage,
	}
}

// temporary data structure to help build a batch from the information found in the rollup
type batchFromRollup struct {
	transactions []*common.L2Tx
	seqNo        *big.Int
	height       *big.Int
	txHash       gethcommon.Hash
	time         uint64
	l1Proof      common.L1BlockHash

	header *common.BatchHeader // for reorgs
}

// CreateExtRollup - creates a compressed and encrypted External rollup from the internal data structure
func (rc *RollupCompression) CreateExtRollup(r *core.Rollup) (*common.ExtRollup, error) {
	header, err := rc.createRollupHeader(r.Batches)
	if err != nil {
		return nil, err
	}
	encryptedHeader, err := rc.serialiseCompressAndEncrypt(header)
	if err != nil {
		return nil, err
	}

	transactions := make([][]*common.L2Tx, len(r.Batches))
	for i, batch := range r.Batches {
		transactions[i] = batch.Transactions
	}
	encryptedTransactions, err := rc.serialiseCompressAndEncrypt(transactions)
	if err != nil {
		return nil, err
	}

	return &common.ExtRollup{
		Header:               r.Header,
		BatchPayloads:        encryptedTransactions,
		CalldataRollupHeader: encryptedHeader,
	}, nil
}

// ProcessExtRollup - given an External rollup, responsible with checking and saving all batches found inside
func (rc *RollupCompression) ProcessExtRollup(rollup *common.ExtRollup) error {
	transactionsPerBatch := make([][]*common.L2Tx, 0)
	err := rc.decryptDecompressAndSerialise(rollup.BatchPayloads, &transactionsPerBatch)
	if err != nil {
		return err
	}

	calldataRollupHeader := new(common.CalldataRollupHeader)
	err = rc.decryptDecompressAndSerialise(rollup.CalldataRollupHeader, calldataRollupHeader)
	if err != nil {
		return err
	}

	// The recreation of batches is a 2-step process:

	// 1. calculate fields like: sequence, height, time, l1Proof, from the implicit and explicit information from the metadata
	incompleteBatches, err := rc.createIncompleteBatches(calldataRollupHeader, transactionsPerBatch)
	if err != nil {
		return err
	}

	// 2. execute each batch to be able to calculate the hash which is necessary for the next batch as it is the parent.
	err = rc.executeAndSaveIncompleteBatches(calldataRollupHeader, incompleteBatches)
	if err != nil {
		return err
	}

	return nil
}

// the main logic that goes from a list of batches to the rollup header
func (rc *RollupCompression) createRollupHeader(batches []*core.Batch) (*common.CalldataRollupHeader, error) {
	reorgs := make([]*common.BatchHeader, len(batches))

	deltaTimes := make([]*big.Int, len(batches))
	startTime := batches[0].Header.Time
	prev := startTime

	l1Proofs := make([]*big.Int, len(batches))
	var prevL1Height *big.Int

	batchHashes := make([]common.L2BatchHash, len(batches))
	batchHeaders := make([]*common.BatchHeader, len(batches))

	isReorg := false
	for i, batch := range batches {
		rc.logger.Info("Add batch to rollup", log.BatchSeqNoKey, batch.SeqNo(), log.BatchHeightKey, batch.Number(), log.BatchHashKey, batch.Hash())
		// determine whether the batch is canonical
		can, err := rc.storage.FetchBatchByHeight(batch.NumberU64())
		if err != nil {
			return nil, err
		}
		if can.Hash() != batch.Hash() {
			// if the canonical batch of the same height is different from the current batch
			// then add the entire header to a "reorgs" array
			reorgs[i] = batch.Header
			isReorg = true
		} else {
			reorgs[i] = nil
		}
		batchHashes[i] = batch.Hash()
		batchHeaders[i] = batch.Header

		deltaTimes[i] = big.NewInt(int64(batch.Header.Time - prev))
		prev = batch.Header.Time

		// since this is the sequencer, it must have all the blocks, because it created the batches in the first place
		block, err := rc.storage.FetchBlock(batch.Header.L1Proof)
		if err != nil {
			return nil, err
		}

		// only add an entry in the l1Proofs array when the value changes
		if i == 0 {
			l1Proofs[i] = block.Number()
		} else {
			l1Proofs[i] = big.NewInt(block.Number().Int64() - prevL1Height.Int64())
		}
		prevL1Height = block.Number()
	}

	if !isReorg {
		reorgs = nil
	}

	calldataRollupHeader := &common.CalldataRollupHeader{
		FirstBatchSequence: batches[0].SeqNo(),
		FirstBatchHeight:   batches[0].Number(), // todo - has to be canonical
		FirstParentHash:    batches[0].Header.ParentHash,
		StartTime:          startTime,
		DeltaTimes:         deltaTimes,
		ReOrgs:             reorgs,
		L1Proofs:           l1Proofs,
		BatchHashes:        batchHashes,
		// BatchHeaders:       batchHeaders,
	}

	return calldataRollupHeader, nil
}

// the main logic to recreate the batches from the header. The logical pair of: `createRollupHeader`
func (rc *RollupCompression) createIncompleteBatches(calldataRollupHeader *common.CalldataRollupHeader, transactionsPerBatch [][]*common.L2Tx) ([]*batchFromRollup, error) {
	incompleteBatches := make([]*batchFromRollup, 0)

	startAtSeq := calldataRollupHeader.FirstBatchSequence.Int64()
	currentHeight := calldataRollupHeader.FirstBatchHeight.Int64() - 1
	currentTime := calldataRollupHeader.StartTime
	var currentL1Height *big.Int

	for currentBatchIdx, batchTransactions := range transactionsPerBatch {
		// the l1 proofs are stored as deltas, which compress well as it should be a series of 1s and 0s
		// the first element is the actual height
		l1Delta := calldataRollupHeader.L1Proofs[currentBatchIdx]
		if currentBatchIdx == 0 {
			currentL1Height = l1Delta
		} else {
			currentL1Height = big.NewInt(l1Delta.Int64() + currentL1Height.Int64())
		}
		block, err := rc.storage.FetchCanonicaBlockByHeight(currentL1Height)
		if err != nil {
			rc.logger.Error("Error decompressing rollup. Did not find l1 block", log.ErrKey, err)
			return nil, err
		}

		// todo - this should be 1 second
		// todo - multiply delta by something?
		currentTime += calldataRollupHeader.DeltaTimes[currentBatchIdx].Uint64()

		// the transactions stored in a valid rollup belong to sequential batches
		currentSeqNo := big.NewInt(startAtSeq + int64(currentBatchIdx))

		// check whether the batch is stored already in the database
		b, err := rc.storage.FetchBatchBySeqNo(currentSeqNo.Uint64())
		if err == nil {
			if len(b.Transactions) != len(batchTransactions) {
				return nil, fmt.Errorf("sanity check failed")
			}
			continue
		}
		if !errors.Is(err, errutil.ErrNotFound) {
			return nil, err
		}

		var h *common.BatchHeader
		if len(calldataRollupHeader.ReOrgs) > 0 {
			// the ReOrgs data structure contains an entire Header
			// for the batches that got re-orged.
			// the assumption is that it can't be computed because the L1 block won't be available.
			h = calldataRollupHeader.ReOrgs[currentBatchIdx]
			if h == nil {
				// only if the batch is canonical, increment the height
				currentHeight = currentHeight + 1
			}
		} else {
			currentHeight = currentHeight + 1
		}

		// calculate the hash of the txs
		var txHash gethcommon.Hash
		if len(batchTransactions) == 0 {
			txHash = types.EmptyRootHash
		} else {
			txHash = types.DeriveSha(types.Transactions(batchTransactions), trie.NewStackTrie(nil))
		}

		incompleteBatches = append(incompleteBatches, &batchFromRollup{
			transactions: batchTransactions,
			seqNo:        currentSeqNo,
			height:       big.NewInt(currentHeight),
			txHash:       txHash,
			time:         currentTime,
			l1Proof:      block.Hash(),
			header:       h,
		})
		rc.logger.Info("Add canon batch", log.BatchSeqNoKey, currentSeqNo, log.BatchHeightKey, currentHeight)
	}
	return incompleteBatches, nil
}

func (rc *RollupCompression) executeAndSaveIncompleteBatches(calldataRollupHeader *common.CalldataRollupHeader, incompleteBatches []*batchFromRollup) error {
	parentHash := calldataRollupHeader.FirstParentHash
	for i, incompleteBatch := range incompleteBatches {
		if incompleteBatch.seqNo.Uint64() == common.L2GenesisSeqNo {
			genBatch, _, err := rc.batchExecutor.CreateGenesisState(incompleteBatch.l1Proof, incompleteBatch.time)
			if err != nil {
				return err
			}
			err = rc.storage.StoreBatch(genBatch)
			if err != nil {
				return err
			}
			err = rc.storage.StoreExecutedBatch(genBatch, nil)
			if err != nil {
				return err
			}
			rc.batchRegistry.OnBatchExecuted(genBatch, nil)

			rc.logger.Info("Stored genesis", log.BatchHashKey, genBatch.Hash())
			parentHash = genBatch.Hash()
			continue
		}

		if incompleteBatch.header != nil {
			err := rc.storage.StoreBatch(&core.Batch{
				Header:       incompleteBatch.header,
				Transactions: incompleteBatch.transactions,
			})
			if err != nil {
				return err
			}
			continue
		}

		// transforms the incompleteBatch into a BatchHeader by executing the transactions
		// and then the info can be used to fill in the parent
		computedBatch, err := rc.batchExecutor.ComputeBatchLight(incompleteBatch.l1Proof,
			parentHash,
			incompleteBatch.transactions,
			incompleteBatch.time,
			incompleteBatch.seqNo,
		)
		if err != nil {
			return err
		}
		if _, err := computedBatch.Commit(true); err != nil {
			return fmt.Errorf("cannot commit stateDB for incoming valid batch seq=%d. Cause: %w", incompleteBatch.seqNo, err)
		}

		err = rc.storage.StoreBatch(computedBatch.Batch)
		if err != nil {
			return err
		}
		err = rc.storage.StoreExecutedBatch(computedBatch.Batch, computedBatch.Receipts)
		if err != nil {
			return err
		}

		parentHash = computedBatch.Batch.Hash()
		if parentHash != calldataRollupHeader.BatchHashes[i] {
			// rc.logger.Info(fmt.Sprintf("Good %+v\nCalc %+v", calldataRollupHeader.BatchHeaders[i], computedBatch.Batch.Header))
			rc.logger.Crit("Rollup decompression failure")
		}
	}
	return nil
}

func (rc *RollupCompression) serialiseCompressAndEncrypt(obj any) ([]byte, error) {
	serialised, err := rlp.EncodeToBytes(obj)
	if err != nil {
		return nil, err
	}
	compressed, err := rc.dataCompressionService.CompressRollup(serialised)
	if err != nil {
		return nil, err
	}
	encrypted, err := rc.dataEncryptionService.Encrypt(compressed)
	if err != nil {
		return nil, err
	}
	return encrypted, nil
}

func (rc *RollupCompression) decryptDecompressAndSerialise(blob []byte, obj any) error {
	plaintextBlob, err := rc.dataEncryptionService.Decrypt(blob)
	if err != nil {
		return err
	}
	serialisedBlob, err := rc.dataCompressionService.Decompress(plaintextBlob)
	if err != nil {
		return err
	}
	err = rlp.DecodeBytes(serialisedBlob, obj)
	if err != nil {
		return err
	}
	return nil
}
