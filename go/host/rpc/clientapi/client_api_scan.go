package clientapi

import (
	gethcommon "github.com/ethereum/go-ethereum/common"
	"math/big"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/host"
	hostdb "github.com/ten-protocol/go-ten/go/host/storage/hostdb"
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
func (s *ScanAPI) GetTotalContractCount() (*big.Int, error) {
	return s.host.EnclaveClient().GetTotalContractCount()
}

// GetTotalTransactionCount returns the number of recorded transactions on the network.
func (s *ScanAPI) GetTotalTransactionCount() (*big.Int, error) {
	return hostdb.GetTotalTransactions(s.host.DB())
}

func (s *ScanAPI) GetBatchListing(pagination *common.QueryPagination) (*common.BatchListingResponse, error) {
	return hostdb.GetBatchListing(s.host.DB(), pagination)
}

func (s *ScanAPI) GetPublicBatchByHash(hash common.L2BatchHash) (*common.PublicBatch, error) {
	return hostdb.GetPublicBatch(s.host.DB(), hash)
}

// GetFullBatchByHash returns the full `ExtBatch` with the given hash.
func (s *ScanAPI) GetFullBatchByHash(batchHash gethcommon.Hash) (*common.ExtBatch, error) {
	return hostdb.GetFullBatch(s.host.DB(), batchHash)
}

// GetFullBatchByTxHash returns the full `ExtBatch` with the given hash.
func (s *ScanAPI) GetFullBatchByTxHash(txHash gethcommon.Hash) (*common.ExtBatch, error) {
	return hostdb.GetFullBatchByTx(s.host.DB(), txHash)
}

func (s *ScanAPI) GetLatestBatch() (*common.BatchHeader, error) {
	return hostdb.GetLatestBatch(s.host.DB())
}
func (s *ScanAPI) GetBatchByHeight(height *big.Int) (*common.BatchHeader, error) {
	return hostdb.GetBatchByHeight(s.host.DB(), height)
}

func (s *ScanAPI) GetRollupListing(pagination *common.QueryPagination) (*common.RollupListingResponse, error) {
	return hostdb.GetRollupListing(s.host.DB(), pagination)
}

func (s *ScanAPI) GetLatestRollupHeader() (*common.RollupHeader, error) {
	return hostdb.GetLatestRollup(s.host.DB())
}

func (s *ScanAPI) GetPublicTransactionData(pagination *common.QueryPagination) (*common.TransactionListingResponse, error) {
	return s.host.EnclaveClient().GetPublicTransactionData(pagination)
}

//func (s *ScanAPI) GetLatestTransactionHashes(pagination *common.QueryPagination) ([]gethcommon.Hash, error) {
//	panic("I don't think this is even needed")
//}

func (s *ScanAPI) GetBlockListing(pagination *common.QueryPagination) (*common.BlockListingResponse, error) {
	return hostdb.GetBlockListing(s.host.DB(), pagination)
}
