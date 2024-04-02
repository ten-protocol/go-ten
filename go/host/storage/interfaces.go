package storage

import (
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/host/storage/hostdb"
	"io"
	"math/big"
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
	// FetchBatchBySeqNo returns the batch with the given seq number
	FetchBatchBySeqNo(seqNum uint64) (*common.ExtBatch, error)
	// FetchBatchHashByNumber returns the batch hash given the batch number
	FetchBatchHashByNumber(number *big.Int) (*gethcommon.Hash, error)
	// FetchBatchHeaderByHash returns the batch header given its hash
	FetchBatchHeaderByHash(hash gethcommon.Hash) (*common.BatchHeader, error)
	// FetchHeadBatchHeader returns the latest batch header
	FetchHeadBatchHeader() (*common.BatchHeader, error)
}

type BlockResolver interface {
	AddBlock(b *types.Header, rollupHash common.L2RollupHash) error
	AddRollup(rollup *common.ExtRollup, metadata *common.PublicRollupMetadata, block *common.L1Block) error
}

type DatabaseResolver interface {
	GetDB() *hostdb.HostDB
}
