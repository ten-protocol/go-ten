package clientapi

import (
	"math/big"

	"github.com/obscuronet/go-obscuro/go/host"

	"github.com/ethereum/go-ethereum/log"
	"github.com/obscuronet/go-obscuro/go/common"
)

type scanAPIServiceLocator interface {
	host.EnclaveLocator
	host.DBLocator
}

// ScanAPI implements metric specific RPC endpoints
type ScanAPI struct {
	sl     scanAPIServiceLocator
	logger log.Logger
}

func NewScanAPI(serviceLocator scanAPIServiceLocator, logger log.Logger) *ScanAPI {
	return &ScanAPI{
		sl:     serviceLocator,
		logger: logger,
	}
}

// GetTotalContractCount returns the number of recorded contracts on the network.
func (s *ScanAPI) GetTotalContractCount() (*big.Int, error) {
	return s.sl.Enclave().GetEnclaveClient().GetTotalContractCount()
}

// GetTotalTransactionCount returns the number of recorded transactions on the network.
func (s *ScanAPI) GetTotalTransactionCount() (*big.Int, error) {
	return s.sl.DB().GetTotalTransactions()
}

func (s *ScanAPI) GetLatestRollupHeader() (*common.RollupHeader, error) {
	return s.sl.DB().GetTipRollupHeader()
}

func (s *ScanAPI) GetPublicTransactionData(pagination *common.QueryPagination) (*common.PublicQueryResponse, error) {
	return s.sl.Enclave().GetEnclaveClient().GetPublicTransactionData(pagination)
}
