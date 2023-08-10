package host

import (
	"encoding/json"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/tracers"
	"github.com/obscuronet/go-obscuro/go/config"
	"github.com/obscuronet/go-obscuro/go/responses"
	"math/big"
)

// Host is the half of the Obscuro node that lives outside the enclave.
type Host interface {
	APIDBRepository
	APIEnclaveClient

	Config() *config.HostConfig

	// Start initializes the main loop of the host.
	Start() error
	// SubmitAndBroadcastTx submits an encrypted transaction to the enclave, and broadcasts it to the other hosts on the network.
	SubmitAndBroadcastTx(encryptedParams common.EncryptedParamsSendRawTx) (*responses.RawTx, error)
	// Subscribe feeds logs matching the encrypted log subscription to the matchedLogs channel.
	Subscribe(id rpc.ID, encryptedLogSubscription common.EncryptedParamsLogSubscription, matchedLogs chan []byte) error
	// Unsubscribe terminates a log subscription between the host and the enclave.
	Unsubscribe(id rpc.ID)
	// Stop gracefully stops the host execution.
	Stop() error

	// HealthCheck returns the health status of the host + enclave + db
	HealthCheck() (*HealthCheck, error)
}

type APIDBRepository interface {
	GetHeadBatchHeader() (*common.BatchHeader, error)
	GetBatchHeader(hash gethcommon.Hash) (*common.BatchHeader, error)
	GetBatchHash(number *big.Int) (*gethcommon.Hash, error)
	GetBlockHeader(hash gethcommon.Hash) (*types.Header, error)
	GetBatch(batchHash gethcommon.Hash) (*common.ExtBatch, error)
	GetBatchNumber(txHash gethcommon.Hash) (*big.Int, error)
	GetBatchTxs(batchHash gethcommon.Hash) ([]gethcommon.Hash, error)
	GetTotalTransactions() (*big.Int, error)
	GetTipRollupHeader() (*common.RollupHeader, error)
}

type APIEnclaveClient interface {
	DebugEventLogRelevancy(hash gethcommon.Hash) (json.RawMessage, common.SystemError)
	DebugTraceTransaction(hash gethcommon.Hash, config *tracers.TraceConfig) (json.RawMessage, common.SystemError)
	ObsCall(encryptedParams common.EncryptedParamsCall) (*responses.Call, common.SystemError)
	EstimateGas(encryptedParams common.EncryptedParamsEstimateGas) (*responses.Gas, common.SystemError)
	GetBalance(encryptedParams common.EncryptedParamsGetBalance) (*responses.Balance, common.SystemError)
	GetCode(address gethcommon.Address, rollupHash *gethcommon.Hash) ([]byte, common.SystemError)
	GetCustomQuery(encryptedParams common.EncryptedParamsGetStorageAt) (*responses.PrivateQueryResponse, common.SystemError)
	GetTransaction(encryptedParams common.EncryptedParamsGetTxByHash) (*responses.TxByHash, common.SystemError)
	GetTransactionCount(encryptedParams common.EncryptedParamsGetTxCount) (*responses.TxCount, common.SystemError)
	GetTransactionReceipt(encryptedParams common.EncryptedParamsGetTxReceipt) (*responses.TxReceipt, common.SystemError)
	GetLogs(encryptedParams common.EncryptedParamsGetLogs) (*responses.Logs, common.SystemError)
	Attestation() (*common.AttestationReport, common.SystemError)
	GetPublicTransactionData(pagination *common.QueryPagination) (*common.PublicQueryResponse, common.SystemError)
	GetTotalContractCount() (*big.Int, common.SystemError)
}
type BatchMsg struct {
	Batches []*common.ExtBatch // The batches being sent.
	IsLive  bool               // true if these batches are being sent as new, false if in response to a p2p request
}
