package enclaveclient

import (
	"context"
	"encoding/json"

	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/rpc/generated"
	"github.com/obscuronet/go-obscuro/go/common/tracers"
	"github.com/obscuronet/go-obscuro/go/config"
	"github.com/obscuronet/go-obscuro/go/responses"
	"google.golang.org/grpc"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
	gethrpc "github.com/ethereum/go-ethereum/rpc"
)

// EnclaveExternalClient implements the common.EnclaveExternal interface for external facing requests (RPC)
type EnclaveExternalClient struct {
	protoClient generated.EnclaveProtoClient
	connection  *grpc.ClientConn
	config      *config.HostConfig
	logger      gethlog.Logger
}

func NewEnclaveExternalClient(
	protoClient generated.EnclaveProtoClient,
	connection *grpc.ClientConn,
	config *config.HostConfig,
	logger gethlog.Logger,
) common.EnclaveExternal {
	return &EnclaveExternalClient{
		protoClient: protoClient,
		connection:  connection,
		config:      config,
		logger:      logger,
	}
}

func (c *EnclaveExternalClient) SubmitTx(tx common.EncryptedTx) responses.RawTx {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.SubmitTx(timeoutCtx, &generated.SubmitTxRequest{EncryptedTx: tx})
	if err != nil {
		return responses.AsPlaintextError(err)
	}

	return *responses.ToEnclaveResponse(response.EncodedEnclaveResponse)
}

func (c *EnclaveExternalClient) ObsCall(encryptedParams common.EncryptedParamsCall) responses.Call {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.ObsCall(timeoutCtx, &generated.ObsCallRequest{
		EncryptedParams: encryptedParams,
	})
	if err != nil {
		return responses.AsPlaintextError(err)
	}
	return *responses.ToEnclaveResponse(response.EncodedEnclaveResponse)
}

func (c *EnclaveExternalClient) GetTransactionCount(encryptedParams common.EncryptedParamsGetTxCount) responses.TxCount {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.GetTransactionCount(timeoutCtx, &generated.GetTransactionCountRequest{EncryptedParams: encryptedParams})
	if err != nil {
		return responses.AsPlaintextError(err)
	}
	return *responses.ToEnclaveResponse(response.EncodedEnclaveResponse)
}

func (c *EnclaveExternalClient) GetTransaction(encryptedParams common.EncryptedParamsGetTxByHash) responses.TxByHash {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	resp, err := c.protoClient.GetTransaction(timeoutCtx, &generated.GetTransactionRequest{EncryptedParams: encryptedParams})
	if err != nil {
		return responses.AsPlaintextError(err)
	}
	return *responses.ToEnclaveResponse(resp.EncodedEnclaveResponse)
}

func (c *EnclaveExternalClient) GetTransactionReceipt(encryptedParams common.EncryptedParamsGetTxReceipt) responses.TxReceipt {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.GetTransactionReceipt(timeoutCtx, &generated.GetTransactionReceiptRequest{EncryptedParams: encryptedParams})
	if err != nil {
		return responses.AsPlaintextError(err)
	}
	return *responses.ToEnclaveResponse(response.EncodedEnclaveResponse)
}

func (c *EnclaveExternalClient) AddViewingKey(viewingKeyBytes []byte, signature []byte) error {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	_, err := c.protoClient.AddViewingKey(timeoutCtx, &generated.AddViewingKeyRequest{
		ViewingKey: viewingKeyBytes,
		Signature:  signature,
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *EnclaveExternalClient) GetBalance(encryptedParams common.EncryptedParamsGetBalance) responses.Balance {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	resp, err := c.protoClient.GetBalance(timeoutCtx, &generated.GetBalanceRequest{
		EncryptedParams: encryptedParams,
	})
	if err != nil {
		return responses.AsPlaintextError(err)
	}
	return *responses.ToEnclaveResponse(resp.EncodedEnclaveResponse)
}

func (c *EnclaveExternalClient) GetCode(address gethcommon.Address, batchHash *gethcommon.Hash) ([]byte, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	resp, err := c.protoClient.GetCode(timeoutCtx, &generated.GetCodeRequest{
		Address:    address.Bytes(),
		RollupHash: batchHash.Bytes(),
	})
	if err != nil {
		return nil, err
	}
	return resp.Code, nil
}

func (c *EnclaveExternalClient) Subscribe(id gethrpc.ID, encryptedParams common.EncryptedParamsLogSubscription) error {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	_, err := c.protoClient.Subscribe(timeoutCtx, &generated.SubscribeRequest{
		Id:                    []byte(id),
		EncryptedSubscription: encryptedParams,
	})
	return err
}

func (c *EnclaveExternalClient) Unsubscribe(id gethrpc.ID) error {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	_, err := c.protoClient.Unsubscribe(timeoutCtx, &generated.UnsubscribeRequest{
		Id: []byte(id),
	})
	return err
}

func (c *EnclaveExternalClient) EstimateGas(encryptedParams common.EncryptedParamsEstimateGas) responses.Gas {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	resp, err := c.protoClient.EstimateGas(timeoutCtx, &generated.EstimateGasRequest{
		EncryptedParams: encryptedParams,
	})
	if err != nil {
		return responses.AsPlaintextError(err)
	}

	return *responses.ToEnclaveResponse(resp.EncodedEnclaveResponse)
}

func (c *EnclaveExternalClient) GetLogs(encryptedParams common.EncryptedParamsGetLogs) responses.Logs {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	resp, err := c.protoClient.GetLogs(timeoutCtx, &generated.GetLogsRequest{
		EncryptedParams: encryptedParams,
	})
	if err != nil {
		return responses.AsPlaintextError(err)
	}
	return *responses.ToEnclaveResponse(resp.EncodedEnclaveResponse)
}

func (c *EnclaveExternalClient) HealthCheck() (bool, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	resp, err := c.protoClient.HealthCheck(timeoutCtx, &generated.EmptyArgs{})
	if err != nil {
		return false, err
	}
	return resp.Status, nil
}

func (c *EnclaveExternalClient) DebugTraceTransaction(hash gethcommon.Hash, config *tracers.TraceConfig) (json.RawMessage, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	confBytes, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	resp, err := c.protoClient.DebugTraceTransaction(timeoutCtx, &generated.DebugTraceTransactionRequest{
		TxHash: hash.Bytes(),
		Config: confBytes,
	})
	if err != nil {
		return nil, err
	}
	return json.RawMessage(resp.Msg), nil
}

func (c *EnclaveExternalClient) DebugEventLogRelevancy(hash gethcommon.Hash) (json.RawMessage, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	resp, err := c.protoClient.DebugEventLogRelevancy(timeoutCtx, &generated.DebugEventLogRelevancyRequest{
		TxHash: hash.Bytes(),
	})
	if err != nil {
		return nil, err
	}
	return json.RawMessage(resp.Msg), nil
}
