package storage

import (
	"database/sql"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ten-protocol/go-ten/go/common"
	"io"
)

type Storage interface {
	BatchResolver
	BlockResolver
	DatabaseResolver
	io.Closer
}

type BatchResolver interface {
	// AddBatch stores the batch
	AddBatch(batch *common.ExtBatch) error
	// FetchBatchBySeqNo returns the batch with the given seq number.
	FetchBatchBySeqNo(seqNum uint64) (*common.ExtBatch, error)
}

type BlockResolver interface {
	AddBlock(b *types.Header, rollupHash common.L2RollupHash) error
	AddRollup(rollup *common.ExtRollup, metadata *common.PublicRollupMetadata, block *common.L1Block) error
}

type DatabaseResolver interface {
	GetDB() *sql.DB
}
