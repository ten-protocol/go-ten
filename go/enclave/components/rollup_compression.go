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
## Requirements:
1. recreate the exact batch headers as the live ones
2. security - make sure it all chains up cryptographically so it can't be gamed

### Ideas
- Exclude empty batches from publishing
- p2p - send an empty payload chained

### Required Fields
- Time - delta from the previous. If expressed in seconds, we could use int16 (2 bytes)
- BaseFee - delta again. Probably can be in8

- L1Proof? - do we need it

MixDigest - required only for the first batch, the rest will derive it by sha256 ( mixDigest || batch_ids)

All the rest of the fields can be secured
- Root
- TxHash
- ReceiptHash
- by publishing in the rollup the Mtree root of all

GasUsed - the sum of all batches published in the rollup
GasLimit - ??

Signatures:
- not required, because the aggregator will sign over the rollup, so there is nothing extra that they add

CrossChainMessages:
- not needed, because they are added to the rollup
*/
/*func createCalldataRollupHeader(batches []*core.Batch) (*common.CalldataRollupHeader, error) {
	if len(batches) == 0 {
		return nil, fmt.Errorf("cannot create a rollup without batches")
	}
	first := batches[0]

	crh := common.CalldataRollupHeader{
		FirstParentHash: first.Header.ParentHash,
		FirstNumber:     first.Number(),
		L1Proofs:        make([]gethcommon.Hash, len(batches)),
		Hashes:          make([]gethcommon.Hash, len(batches)),
		XChainHash:      make([]gethcommon.Hash, len(batches)),
		XChainHeight:    make([]*big.Int, len(batches)),

		StartTime:  first.Header.Time,
		DeltaTimes: make([]uint8, len(batches)-1),
		//baseFeesDelta: make([]uint8, len(batches)-1),
	}

	prev := first.Header.Time
	for i, h := range batches {
		if i > 0 {
			// todo - handle the conversion not going well
			crh.DeltaTimes[i-1] = uint8(h.Header.Time - prev)
			prev = h.Header.Time
		}
		crh.Hashes[i] = h.Hash()
		crh.L1Proofs[i] = h.Header.L1Proof
		crh.XChainHash[i] = h.Header.LatestInboundCrossChainHash
		crh.XChainHeight[i] = h.Header.LatestInboundCrossChainHeight
	}

	return &crh, nil
}

func reconstructBatches(calldataRollupHeader *common.CalldataRollupHeader, batchTransactions [][]*common.L2Tx) ([]*core.Batch, error) {
	var ch common.CalldataRollupHeader
	if err := rlp.DecodeBytes(headersBlob, &ch); err != nil {
		return nil, err
	}

	nrHeaders := len(ch.deltaTimes) + 1

	headers := make([]*common.BatchHeader, nrHeaders)
	headers[0] = &common.BatchHeader{
		ParentHash:                    ch.firstParentHash,
		Number:                        ch.firstNumber,
		Time:                          ch.startTime,
		L1Proof:                       ch.l1Proofs[0],
		LatestInboundCrossChainHash:   ch.xChainHash[0],
		LatestInboundCrossChainHeight: ch.xChainHeight[0],
	}

	prevTime := ch.startTime
	prevNumber := ch.firstNumber

	for i := range ch.deltaTimes {
		headers[i+1] = &common.BatchHeader{
			ParentHash:                    ch.hashes[i],
			Number:                        prevNumber.Add(prevNumber, big.NewInt(1)),
			Time:                          prevTime + headers[i+1].Time,
			L1Proof:                       ch.l1Proofs[i+1],
			LatestInboundCrossChainHash:   ch.xChainHash[i+1],
			LatestInboundCrossChainHeight: ch.xChainHeight[i+1],
		}
		prevTime = headers[i+1].Time
	}

	return headers, nil
}
*/

// responsible for
// given a list of batches - create the CalldataRollupHeader
// given a CalldataRollupHeader and transactions - go to executed batches
//

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

func (rc *RollupCompression) CreateExtRollup(r *core.Rollup) (*common.ExtRollup, error) {
	transactions := make([][]*common.L2Tx, len(r.Batches))
	for i, batch := range r.Batches {
		transactions[i] = batch.Transactions
	}
	encryptedTransactions, err := rc.serialiseCompressAndEncrypt(transactions)
	if err != nil {
		return nil, err
	}

	reorgs := make([]*common.BatchHeader, len(r.Batches))

	deltaTimes := make([]uint64, len(r.Batches))
	startTime := r.Batches[0].Header.Time
	prev := startTime

	l1Proofs := make([]*big.Int, len(r.Batches))
	var currentL1Proof *big.Int = nil

	batchHashes := make([]common.L2BatchHash, len(r.Batches))
	batchHeaders := make([]*common.BatchHeader, len(r.Batches))

	isReorg := false
	for i, batch := range r.Batches {

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

		deltaTimes[i] = batch.Header.Time - prev
		prev = batch.Header.Time

		// since this is the sequencer, it must have all the blocks, because it created the batches in the first place
		block, err := rc.storage.FetchBlock(batch.Header.L1Proof)
		if err != nil {
			return nil, err
		}
		l1Proofs[i] = nil
		// only add an entry in the l1Proofs array when the value changes
		if currentL1Proof == nil || block.NumberU64() != currentL1Proof.Uint64() {
			l1Proofs[i] = block.Number()
			currentL1Proof = block.Number()
		}
		if i == 0 && currentL1Proof == nil {
			return nil, fmt.Errorf("faile sanity check. the l1Proofs array must start with a non-nil value")
		}
	}

	if !isReorg {
		reorgs = nil
	}

	calldataRollupHeader := &common.CalldataRollupHeader{
		FirstBatchSequence: r.Batches[0].SeqNo(),
		FirstBatchHeight:   r.Batches[0].Number(), // todo - has to be canonical
		FirstParentHash:    r.Batches[0].Header.ParentHash,
		StartTime:          startTime,
		DeltaTimes:         deltaTimes,
		ReOrgs:             reorgs,
		L1Proofs:           l1Proofs,
		BatchHashes:        batchHashes,
		// BatchHeaders:       batchHeaders,
	}

	encryptedHeader, err := rc.serialiseCompressAndEncrypt(calldataRollupHeader)
	if err != nil {
		return nil, err
	}

	return &common.ExtRollup{
		Header:               r.Header,
		BatchPayloads:        encryptedTransactions,
		CalldataRollupHeader: encryptedHeader,
	}, nil
}

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
	// adjust the L1Proofs array to contains nils
	for i, proof := range calldataRollupHeader.L1Proofs {
		if len(proof.Bytes()) == 0 {
			calldataRollupHeader.L1Proofs[i] = nil
		}
	}

	// The recreation of batches is a 2-step process:
	// 1. calculate fields like: sequence, height, time, l1Proof
	// 2. execute each batch to be able to calculate the hash which is necessary for the next batch as it is the parent.

	incompleteBatches, err := rc.createIncompleteBatches(calldataRollupHeader, transactionsPerBatch)
	if err != nil {
		return err
	}

	// Step 2: execute the canonical batches one by one
	// Note: To keep the logic simple, the execution is not persisted. The Batches are stored in the db to be handled by the usual mechanisms.
	// This is slightly inefficient, but it won't be a bottleneck during live running, because most nodes will receive the batches via p2p.
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

func (rc *RollupCompression) createIncompleteBatches(calldataRollupHeader *common.CalldataRollupHeader, transactionsPerBatch [][]*common.L2Tx) ([]*batchFromRollup, error) {
	incompleteBatches := make([]*batchFromRollup, 0)

	// Step 1
	startAtSeq := calldataRollupHeader.FirstBatchSequence.Int64()
	currentHeight := calldataRollupHeader.FirstBatchHeight.Int64() - 1
	currentTime := calldataRollupHeader.StartTime
	var currentL1Proof *gethcommon.Hash = nil

	for currentBatchIdx, batchTransactions := range transactionsPerBatch {

		// the l1 proofs are stored as heights in a sparse array. There is a value only when the l1 proof changes
		l1HeightOfProof := calldataRollupHeader.L1Proofs[currentBatchIdx]
		if l1HeightOfProof != nil {
			block, err := rc.storage.FetchCanonicaBlockByHeight(l1HeightOfProof)
			if err != nil {
				rc.logger.Error("Error decompressing rollup. Did not find l1 block", log.ErrKey, err)
				return nil, err
			}
			h := block.Hash()
			currentL1Proof = &h
		}

		// todo - this should be 1 second
		// todo - multiply delta by something?
		currentTime += calldataRollupHeader.DeltaTimes[currentBatchIdx]

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

		if currentL1Proof == nil {
			return nil, fmt.Errorf("invalid rollup. The l1 proofs array is not well formed.")
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
			l1Proof:      *currentL1Proof,
			header:       h,
		})
		rc.logger.Info("Add canon batch", log.BatchSeqNoKey, currentSeqNo, log.BatchHeightKey, currentHeight)
	}
	return incompleteBatches, nil
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
