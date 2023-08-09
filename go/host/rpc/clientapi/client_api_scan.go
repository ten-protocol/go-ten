package clientapi

import (
	"math/big"

	"github.com/ethereum/go-ethereum/log"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/host"
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
	return s.host.DB().GetTotalTransactions()
}

func (s *ScanAPI) GetLatestRollupHeader() (*common.RollupHeader, error) {
	return s.host.DB().GetTipRollupHeader()
}

func (s *ScanAPI) GetPublicTransactionData(pagination *common.QueryPagination) (*common.PublicQueryResponse, error) {
	return s.host.EnclaveClient().GetPublicTransactionData(pagination)
}
