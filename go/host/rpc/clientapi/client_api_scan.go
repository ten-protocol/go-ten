package clientapi

import (
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/obscuronet/go-obscuro/go/common/host"

	gethcommon "github.com/ethereum/go-ethereum/common"
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

func (s *ScanAPI) GetLatestBlockHeader() (*types.Header, error) {
	// Place holder until the rollups are in
	return &types.Header{
		ParentHash:  gethcommon.Hash{},
		UncleHash:   gethcommon.Hash{},
		Coinbase:    gethcommon.Address{},
		Root:        gethcommon.Hash{},
		TxHash:      gethcommon.Hash{},
		ReceiptHash: gethcommon.Hash{},
		Bloom:       types.Bloom{},
		Difficulty:  gethcommon.Big0,
		Number:      gethcommon.Big0,
		GasLimit:    0,
		GasUsed:     0,
		Time:        0,
		Extra:       nil,
		MixDigest:   gethcommon.Hash{},
		Nonce:       types.BlockNonce{},
		BaseFee:     nil,
	}, nil
}
