package actors

import (
	"bytes"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/errutil"
	"github.com/obscuronet/go-obscuro/go/enclave/components"
	"github.com/obscuronet/go-obscuro/go/enclave/core"
	"github.com/obscuronet/go-obscuro/go/enclave/db"
)

type obsValidator struct {
	consumer       components.BlockConsumer
	producer       components.BatchProducer
	registry       components.BatchRegistry
	rollupConsumer components.RollupConsumer

	chainConfig *params.ChainConfig

	sequencerID gethcommon.Address
	storage     db.Storage
}

func NewValidator(
	consumer components.BlockConsumer,
	producer components.BatchProducer,
	registry components.BatchRegistry,
	rollupConsumer components.RollupConsumer,

	chainConfig *params.ChainConfig,

	sequencerID gethcommon.Address,
	storage db.Storage,
) ObsValidator {
	return &obsValidator{
		consumer:       consumer,
		producer:       producer,
		registry:       registry,
		rollupConsumer: rollupConsumer,
		chainConfig:    chainConfig,
		sequencerID:    sequencerID,
		storage:        storage,
	}
}

func (ov *obsValidator) GetLatestHead() (*core.Batch, error) {
	return ov.registry.GetHeadBatch()
}

func (ov *obsValidator) handleGenesisBatch(incomingBatch *core.Batch) (bool, error) {
	// genesis
	if incomingBatch.NumberU64() != 0 {
		return false, nil
	}

	batch, _, err := ov.producer.CreateGenesisState(incomingBatch.Header.L1Proof, ov.sequencerID, incomingBatch.Header.Time)
	if err != nil {
		return true, err
	}

	if !bytes.Equal(incomingBatch.Hash().Bytes(), batch.Hash().Bytes()) {
		return true, fmt.Errorf("received bad genesis batch")
	}

	return true, ov.registry.StoreBatch(incomingBatch, nil)
}

func (ov *obsValidator) ValidateAndStoreBatch(incomingBatch *core.Batch) error {
	if handled, err := ov.handleGenesisBatch(incomingBatch); handled {
		return err
	}

	if batch, err := ov.registry.GetBatch(*incomingBatch.Hash()); err != nil && !errors.Is(err, errutil.ErrNotFound) {
		return err
	} else if batch != nil {
		return nil // already know about this one
	}

	if err := ov.CheckSequencerSignature(incomingBatch.Hash(), &incomingBatch.Header.Agg, incomingBatch.Header.R, incomingBatch.Header.S); err != nil {
		return err
	}

	cb, err := ov.producer.ComputeBatch(&components.BatchContext{
		BlockPtr:     incomingBatch.Header.L1Proof,
		ParentPtr:    incomingBatch.Header.ParentHash,
		Transactions: incomingBatch.Transactions,
		AtTime:       incomingBatch.Header.Time,
		Randomness:   incomingBatch.Header.MixDigest,
		Creator:      incomingBatch.Header.Agg,
		ChainConfig:  ov.chainConfig,
	})

	if err != nil {
		return fmt.Errorf("failed computing batch. Cause: %w", err)
	}

	if !bytes.Equal(cb.Batch.Hash().Bytes(), incomingBatch.Hash().Bytes()) {
		return fmt.Errorf("batch is in invalid state. Incoming hash: %s  Computed hash: %s", incomingBatch.Hash().Hex(), cb.Batch.Hash().Hex())
	}

	if _, err := cb.Commit(true); err != nil {
		return fmt.Errorf("cannot commit stateDB for incoming valid batch. Cause: %w", err)
	}

	return ov.registry.StoreBatch(incomingBatch, cb.Receipts)
}

func (ov *obsValidator) ReceiveBlock(br *common.BlockAndReceipts, isLatest bool) (*components.BlockIngestionType, error) {
	ingestion, err := ov.consumer.ConsumeBlock(br, isLatest)
	if err != nil {
		return nil, err
	}

	rollups, err := ov.rollupConsumer.ProcessL1Block(br)
	if err != nil {
		//todo - log err?
		return ingestion, nil
	}

	for _, rollup := range rollups {
		for _, batch := range rollup.Batches {
			ov.ValidateAndStoreBatch(batch)
		}
	}

	return ingestion, nil
}

func (ov *obsValidator) CheckSequencerSignature(headerHash *gethcommon.Hash, aggregator *gethcommon.Address, sigR *big.Int, sigS *big.Int) error {
	// Batches and rollups should only be produced by the sequencer.
	// todo (#718) - sequencer identities should be retrieved from the L1 management contract
	if !bytes.Equal(aggregator.Bytes(), ov.sequencerID.Bytes()) {
		return fmt.Errorf("expected batch to be produced by sequencer %s, but was produced by %s", ov.sequencerID.Hex(), aggregator.Hex())
	}

	if sigR == nil || sigS == nil {
		return fmt.Errorf("missing signature on batch")
	}

	pubKey, err := ov.storage.FetchAttestedKey(*aggregator)
	if err != nil {
		return fmt.Errorf("could not retrieve attested key for aggregator %s. Cause: %w", aggregator, err)
	}

	if !ecdsa.Verify(pubKey, headerHash.Bytes(), sigR, sigS) {
		return fmt.Errorf("could not verify ECDSA signature")
	}
	return nil
}
