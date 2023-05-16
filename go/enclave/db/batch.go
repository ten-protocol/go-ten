package db

import (
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/enclave/core"
	obscurorawdb "github.com/obscuronet/go-obscuro/go/enclave/db/rawdb"
	"github.com/obscuronet/go-obscuro/go/enclave/db/sql"
)

type storageUpdaterImpl struct {
	sqlBatch *sql.Batch
	storage  *storageImpl
}

func NewStorageUpdater(batch *sql.Batch, storage *storageImpl) StorageUpdater {
	return &storageUpdaterImpl{
		sqlBatch: batch,
		storage:  storage,
	}
}

func (s *storageUpdaterImpl) Commit() error {
	if err := s.sqlBatch.Write(); err != nil {
		return fmt.Errorf("could not write batch to storage. Cause: %w", err)
	}
	return nil
}

func (s *storageUpdaterImpl) StoreBatch(batch *core.Batch, receipts []*types.Receipt) error {
	if err := obscurorawdb.WriteBatch(s.sqlBatch, batch); err != nil {
		return fmt.Errorf("could not write batch. Cause: %w", err)
	}
	if err := obscurorawdb.WriteTxLookupEntriesByBatch(s.sqlBatch, batch); err != nil {
		return fmt.Errorf("could not write transaction lookup entries by batch. Cause: %w", err)
	}
	if err := obscurorawdb.WriteReceipts(s.sqlBatch, *batch.Hash(), receipts); err != nil {
		return fmt.Errorf("could not write transaction receipts. Cause: %w", err)
	}
	if err := obscurorawdb.WriteContractCreationTxs(s.sqlBatch, receipts); err != nil {
		return fmt.Errorf("could not save contract creation transaction. Cause: %w", err)
	}
	return nil
}

// UpdateHeadBatch updates the canonical L2 head batch for a given L1 block.
func (s *storageUpdaterImpl) UpdateHeadBatch(l1Head common.L1BlockHash, l2Head *core.Batch, receipts []*types.Receipt) error {
	if err := obscurorawdb.SetL2HeadBatch(s.sqlBatch, *l2Head.Hash()); err != nil {
		return fmt.Errorf("could not write block state. Cause: %w", err)
	}
	if err := obscurorawdb.WriteL1ToL2BatchMapping(s.sqlBatch, l1Head, *l2Head.Hash()); err != nil {
		return fmt.Errorf("could not write block state. Cause: %w", err)
	}

	// We update the canonical hash of the batch at this height.
	if err := obscurorawdb.WriteCanonicalHash(s.sqlBatch, l2Head); err != nil {
		return fmt.Errorf("could not write canonical hash. Cause: %w", err)
	}

	if l2Head.Number().Int64() > 1 {
		err2 := s.storage.writeLogs(l2Head.Header.ParentHash, receipts, s.sqlBatch)
		if err2 != nil {
			return fmt.Errorf("could not save logs %w", err2)
		}
	}
	return nil
}

// SetHeadBatchPointer updates the canonical L2 head batch for a given L1 block.
func (s *storageUpdaterImpl) SetHeadBatchPointer(l2Head *core.Batch) error {
	// We update the canonical hash of the batch at this height.
	if err := obscurorawdb.SetL2HeadBatch(s.sqlBatch, *l2Head.Hash()); err != nil {
		return fmt.Errorf("could not write canonical hash. Cause: %w", err)
	}
	return nil
}
