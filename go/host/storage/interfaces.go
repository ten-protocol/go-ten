package storage

import (
	"io"
	"math/big"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/host/storage/hostdb"
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
	// FetchBatchHashByHeight returns the batch hash given the batch number
	FetchBatchHashByHeight(number *big.Int) (*gethcommon.Hash, error)
	// FetchBatchHeaderByHash returns the batch header given its hash
	FetchBatchHeaderByHash(hash gethcommon.Hash) (*common.BatchHeader, error)
	// FetchHeadBatchHeader returns the latest batch header
	FetchHeadBatchHeader() (*common.BatchHeader, error)
	// FetchPublicBatchByHash returns the public batch
	FetchPublicBatchByHash(batchHash common.L2BatchHash) (*common.PublicBatch, error)
	// FetchBatch returns the `ExtBatch` with the given hash
	FetchBatch(batchHash gethcommon.Hash) (*common.ExtBatch, error)
	// FetchBatchByTx returns the `ExtBatch` with the given tx hash
	FetchBatchByTx(txHash gethcommon.Hash) (*common.ExtBatch, error)
	// FetchLatestBatch returns the head `BatchHeader`
	FetchLatestBatch() (*common.BatchHeader, error)
	// FetchBatchListing returns a paginated list of the public batch data
	FetchBatchListing(pagination *common.QueryPagination) (*common.BatchListingResponse, error)
	// FetchBatchListingDeprecated backwards compatible API to return batch data
	FetchBatchListingDeprecated(pagination *common.QueryPagination) (*common.BatchListingResponseDeprecated, error)
	// FetchBatchHeaderByHeight returns the `BatchHeader` with the given height
	FetchBatchHeaderByHeight(height *big.Int) (*common.BatchHeader, error)
	// FetchLatestRollupHeader returns the head `RollupHeader`
	FetchLatestRollupHeader() (*common.RollupHeader, error)
	// FetchRollupListing returns a paginated list of rollups
	FetchRollupListing(pagination *common.QueryPagination) (*common.RollupListingResponse, error)
	// FetchBlockListing returns a paginated list of blocks that include rollups
	FetchBlockListing(pagination *common.QueryPagination) (*common.BlockListingResponse, error)
}

type BlockResolver interface {
	AddBlock(b *types.Header, rollupHash common.L2RollupHash) error
	AddRollup(rollup *common.ExtRollup, metadata *common.PublicRollupMetadata, block *common.L1Block) error
}

type DatabaseResolver interface {
	GetDB() hostdb.HostDB
}
