package clientapi

import (
	"context"
	"math/big"

	gethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/host"
)

// ScanAPI implements metric specific RPC endpoints
type ScanAPI struct {
	host   host.Host
	logger log.Logger
}

func NewScanAPI(host host.Host, logger log.Logger) *ScanAPI {
	return &ScanAPI{
		host:   host,
		logger: logger,
	}
}

// GetTotalContractCount returns the number of recorded contracts on the network.
func (s *ScanAPI) GetTotalContractCount(ctx context.Context) (*big.Int, error) {
	return s.host.EnclaveClient().GetTotalContractCount(ctx)
}

// GetTotalTxCount returns the number of recorded transactions on the network.
func (s *ScanAPI) GetTotalTransactionCount() (*big.Int, error) {
	return s.host.Storage().FetchTotalTxCount()
}

// GetBatchListingNew returns a paginated list of batches
func (s *ScanAPI) GetBatchListingNew(pagination *common.QueryPagination) (*common.BatchListingResponse, error) {
	return s.host.Storage().FetchBatchListing(pagination)
}

// GetBatchListing returns the deprecated version of batch listing
func (s *ScanAPI) GetBatchListing(pagination *common.QueryPagination) (*common.BatchListingResponseDeprecated, error) {
	return s.host.Storage().FetchBatchListingDeprecated(pagination)
}

// GetPublicBatchByHash returns the public batch
func (s *ScanAPI) GetPublicBatchByHash(hash common.L2BatchHash) (*common.PublicBatch, error) {
	return s.host.Storage().FetchPublicBatchByHash(hash)
}

// GetBatch returns the `ExtBatch` with the given hash
func (s *ScanAPI) GetBatch(batchHash gethcommon.Hash) (*common.ExtBatch, error) {
	return s.host.Storage().FetchBatch(batchHash)
}

// GetBatchByTx returns the `ExtBatch` with the given tx hash
func (s *ScanAPI) GetBatchByTx(txHash gethcommon.Hash) (*common.ExtBatch, error) {
	return s.host.Storage().FetchBatchByTx(txHash)
}

// GetLatestBatch returns the head `BatchHeader`
func (s *ScanAPI) GetLatestBatch() (*common.BatchHeader, error) {
	return s.host.Storage().FetchLatestBatch()
}

// GetBatchByHeight returns the `BatchHeader` with the given height
func (s *ScanAPI) GetBatchByHeight(height *big.Int) (*common.PublicBatch, error) {
	return s.host.Storage().FetchBatchByHeight(height)
}

// GetRollupListing returns a paginated list of Rollups
func (s *ScanAPI) GetRollupListing(pagination *common.QueryPagination) (*common.RollupListingResponse, error) {
	return s.host.Storage().FetchRollupListing(pagination)
}

// GetLatestRollupHeader returns the head `RollupHeader`
func (s *ScanAPI) GetLatestRollupHeader() (*common.RollupHeader, error) {
	return s.host.Storage().FetchLatestRollupHeader()
}

// GetPublicTransactionData returns a paginated list of transaction data
func (s *ScanAPI) GetPublicTransactionData(ctx context.Context, pagination *common.QueryPagination) (*common.TransactionListingResponse, error) {
	return s.host.EnclaveClient().GetPublicTransactionData(ctx, pagination)
}

// GetBlockListing returns a paginated list of blocks that include rollups
func (s *ScanAPI) GetBlockListing(pagination *common.QueryPagination) (*common.BlockListingResponse, error) {
	return s.host.Storage().FetchBlockListing(pagination)
}

// GetRollupByHash returns the public rollup data given its hash
func (s *ScanAPI) GetRollupByHash(rollupHash gethcommon.Hash) (*common.PublicRollup, error) {
	return s.host.Storage().FetchRollupByHash(rollupHash)
}

// GetRollupBatches returns the list of batches included in a rollup given its hash
func (s *ScanAPI) GetRollupBatches(rollupHash gethcommon.Hash) (*common.BatchListingResponse, error) {
	return s.host.Storage().FetchRollupBatches(rollupHash)
}

// GetBatchTransactions returns the public tx data of all txs present in a rollup given its hash
func (s *ScanAPI) GetBatchTransactions(batchHash gethcommon.Hash) (*common.TransactionListingResponse, error) {
	return s.host.Storage().FetchBatchTransactions(batchHash)
}
