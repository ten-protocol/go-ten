package clientapi

import (
	"math/big"

	"github.com/ethereum/go-ethereum/log"
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

//
//func (s *ScanAPI) GetL1RollupHeaderByHash() (*common.RollupHeader, error) {
//	//return s.host.DB().GetBatchHeader()
//	return nil, nil
//}

//func (s *ScanAPI) ActiveNodeCount() (int, error) {
//	//return s.host.DB().GetBatchHeader()
//	return 0, nil
//}

func (s *ScanAPI) GetTotalContractCount() (*big.Int, error) {
	return s.host.EnclaveClient().GetTotalContractCount()
}

//
//func (s *ScanAPI) GetTotalTransactionCount() (*big.Int, error) {
//	//return s.host.DB().GetTotalTransactionCount()
//	return s.host.DB().GetTotalTransactions()
//}
