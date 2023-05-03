package enclaverpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/common/retry"
	"github.com/obscuronet/go-obscuro/go/common/rpc"
	"github.com/obscuronet/go-obscuro/go/common/rpc/generated"
	"github.com/obscuronet/go-obscuro/go/common/syserr"
	"github.com/obscuronet/go-obscuro/go/common/tracers"
	"github.com/obscuronet/go-obscuro/go/config"
	"github.com/obscuronet/go-obscuro/go/responses"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
	gethrpc "github.com/ethereum/go-ethereum/rpc"
)

// Client implements enclave.Enclave and should be used by the host when communicating with the enclave via RPC.
type Client struct {
	protoClient generated.EnclaveProtoClient
	connection  *grpc.ClientConn
	config      *config.HostConfig
	logger      gethlog.Logger
}

func NewClient(config *config.HostConfig, logger gethlog.Logger) *Client {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	connection, err := grpc.Dial(config.EnclaveRPCAddress, opts...)
	if err != nil {
		logger.Crit("Failed to connect to enclave RPC service.", log.ErrKey, err)
	}
	connection.Connect()
	// perform an initial sleep because that Connect() method is not blocking and the retry immediately checks the status
	time.Sleep(500 * time.Millisecond)

	// We wait for the RPC connection to be ready.
	err = retry.Do(func() error {
		currState := connection.GetState()
		if currState != connectivity.Ready {
			logger.Info("retrying connection until enclave is available", "status", currState.String())
			connection.Connect()
			return fmt.Errorf("connection is not ready, status=%s", currState)
		}
		// connection is ready, break out of the loop
		return nil
	}, retry.NewBackoffAndRetryForeverStrategy([]time.Duration{500 * time.Millisecond, 1 * time.Second, 5 * time.Second}, 10*time.Second))

	if err != nil {
		// this should not happen as we retry forever...
		logger.Crit("failed to connect to enclave", log.ErrKey, err)
	}

	return &Client{
		protoClient: generated.NewEnclaveProtoClient(connection),
		connection:  connection,
		config:      config,
		logger:      logger,
	}
}

func (c *Client) StopClient() common.SystemError {
	return c.connection.Close()
}

func (c *Client) Status() (common.Status, common.SystemError) {
	if c.connection.GetState() != connectivity.Ready {
		return common.Unavailable, syserr.NewInternalError(fmt.Errorf("RPC connection is not ready"))
	}

	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	status, err := c.protoClient.Status(timeoutCtx, &generated.StatusRequest{})
	if sysErr := c.handleErr(status, err); sysErr != nil {
		return common.Unavailable, sysErr
	}

	return common.Status(status.GetStatus()), nil
}

func (c *Client) Attestation() (*common.AttestationReport, common.SystemError) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	attestation, err := c.protoClient.Attestation(timeoutCtx, &generated.AttestationRequest{})
	if sysErr := c.handleErr(attestation, err); sysErr != nil {
		return nil, sysErr
	}

	return rpc.FromAttestationReportMsg(attestation.AttestationReportMsg), nil
}

func (c *Client) GenerateSecret() (common.EncryptedSharedEnclaveSecret, common.SystemError) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.GenerateSecret(timeoutCtx, &generated.GenerateSecretRequest{})
	if sysErr := c.handleErr(response, err); sysErr != nil {
		return nil, sysErr
	}

	return response.EncryptedSharedEnclaveSecret, nil
}

func (c *Client) InitEnclave(secret common.EncryptedSharedEnclaveSecret) common.SystemError {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	_, err := c.protoClient.InitEnclave(timeoutCtx, &generated.InitEnclaveRequest{EncryptedSharedEnclaveSecret: secret})
	if err != nil {
		return syserr.NewRPCError(err)
	}

	return nil
}

func (c *Client) SubmitL1Block(block types.Block, receipts types.Receipts, isLatest bool) (*common.BlockSubmissionResponse, common.SystemError) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	var buffer bytes.Buffer
	if err := block.EncodeRLP(&buffer); err != nil {
		return nil, fmt.Errorf("could not encode block. Cause: %w", err)
	}

	serialized, err := rlp.EncodeToBytes(receipts)
	if err != nil {
		return nil, fmt.Errorf("could not encode receipts. Cause: %w", err)
	}

	response, err := c.protoClient.SubmitL1Block(timeoutCtx, &generated.SubmitBlockRequest{EncodedBlock: buffer.Bytes(), EncodedReceipts: serialized, IsLatest: isLatest})
	if err != nil {
		return nil, fmt.Errorf("could not submit block. Cause: %w", err)
	}

	blockSubmissionResponse, err := rpc.FromBlockSubmissionResponseMsg(response.BlockSubmissionResponse)
	if err != nil {
		return nil, err
	}
	return blockSubmissionResponse, nil
}

func (c *Client) SubmitTx(tx common.EncryptedTx) (*responses.RawTx, common.SystemError) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.SubmitTx(timeoutCtx, &generated.SubmitTxRequest{EncryptedTx: tx})
	if sysErr := c.handleErr(response, err); sysErr != nil {
		return nil, sysErr
	}

	return responses.ToEnclaveResponse(response.EncodedEnclaveResponse), nil
}

func (c *Client) SubmitBatch(batch *common.ExtBatch) common.SystemError {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	batchMsg := rpc.ToExtBatchMsg(batch)
	_, err := c.protoClient.SubmitBatch(timeoutCtx, &generated.SubmitBatchRequest{Batch: &batchMsg})
	if err != nil {
		return syserr.NewRPCError(err)
	}
	return nil
}

func (c *Client) ObsCall(encryptedParams common.EncryptedParamsCall) (*responses.Call, common.SystemError) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.ObsCall(timeoutCtx, &generated.ObsCallRequest{
		EncryptedParams: encryptedParams,
	})
	if sysErr := c.handleErr(response, err); sysErr != nil {
		return nil, sysErr
	}
	return responses.ToEnclaveResponse(response.EncodedEnclaveResponse), nil
}

func (c *Client) GetTransactionCount(encryptedParams common.EncryptedParamsGetTxCount) (*responses.TxCount, common.SystemError) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.GetTransactionCount(timeoutCtx, &generated.GetTransactionCountRequest{EncryptedParams: encryptedParams})
	if sysErr := c.handleErr(response, err); sysErr != nil {
		return nil, sysErr
	}
	return responses.ToEnclaveResponse(response.EncodedEnclaveResponse), nil
}

func (c *Client) Stop() common.SystemError {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	_, err := c.protoClient.Stop(timeoutCtx, &generated.StopRequest{})
	if err != nil {
		return syserr.NewRPCError(fmt.Errorf("could not stop enclave: %w", err))
	}
	return nil
}

func (c *Client) GetTransaction(encryptedParams common.EncryptedParamsGetTxByHash) (*responses.TxByHash, common.SystemError) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.GetTransaction(timeoutCtx, &generated.GetTransactionRequest{EncryptedParams: encryptedParams})
	if sysErr := c.handleErr(response, err); sysErr != nil {
		return nil, sysErr
	}
	return responses.ToEnclaveResponse(response.EncodedEnclaveResponse), nil
}

func (c *Client) GetTransactionReceipt(encryptedParams common.EncryptedParamsGetTxReceipt) (*responses.TxReceipt, common.SystemError) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.GetTransactionReceipt(timeoutCtx, &generated.GetTransactionReceiptRequest{EncryptedParams: encryptedParams})
	if sysErr := c.handleErr(response, err); sysErr != nil {
		return nil, sysErr
	}
	return responses.ToEnclaveResponse(response.EncodedEnclaveResponse), nil
}

func (c *Client) AddViewingKey(viewingKeyBytes []byte, signature []byte) common.SystemError {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	_, err := c.protoClient.AddViewingKey(timeoutCtx, &generated.AddViewingKeyRequest{
		ViewingKey: viewingKeyBytes,
		Signature:  signature,
	})
	if err != nil {
		return syserr.NewRPCError(err)
	}
	return nil
}

func (c *Client) GetBalance(encryptedParams common.EncryptedParamsGetBalance) (*responses.Balance, common.SystemError) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.GetBalance(timeoutCtx, &generated.GetBalanceRequest{
		EncryptedParams: encryptedParams,
	})
	if sysErr := c.handleErr(response, err); sysErr != nil {
		return nil, sysErr
	}
	return responses.ToEnclaveResponse(response.EncodedEnclaveResponse), nil
}

func (c *Client) GetCode(address gethcommon.Address, batchHash *gethcommon.Hash) ([]byte, common.SystemError) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.GetCode(timeoutCtx, &generated.GetCodeRequest{
		Address:    address.Bytes(),
		RollupHash: batchHash.Bytes(),
	})
	if sysErr := c.handleErr(response, err); sysErr != nil {
		return nil, sysErr
	}
	return response.Code, nil
}

func (c *Client) Subscribe(id gethrpc.ID, encryptedParams common.EncryptedParamsLogSubscription) common.SystemError {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	_, err := c.protoClient.Subscribe(timeoutCtx, &generated.SubscribeRequest{
		Id:                    []byte(id),
		EncryptedSubscription: encryptedParams,
	})
	if err != nil {
		return syserr.NewRPCError(err)
	}
	return nil
}

func (c *Client) Unsubscribe(id gethrpc.ID) common.SystemError {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	_, err := c.protoClient.Unsubscribe(timeoutCtx, &generated.UnsubscribeRequest{
		Id: []byte(id),
	})
	if err != nil {
		return syserr.NewRPCError(err)
	}
	return nil
}

func (c *Client) EstimateGas(encryptedParams common.EncryptedParamsEstimateGas) (*responses.Gas, common.SystemError) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.EstimateGas(timeoutCtx, &generated.EstimateGasRequest{
		EncryptedParams: encryptedParams,
	})
	if sysErr := c.handleErr(response, err); sysErr != nil {
		return nil, sysErr
	}

	return responses.ToEnclaveResponse(response.EncodedEnclaveResponse), nil
}

func (c *Client) GetLogs(encryptedParams common.EncryptedParamsGetLogs) (*responses.Logs, common.SystemError) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.GetLogs(timeoutCtx, &generated.GetLogsRequest{
		EncryptedParams: encryptedParams,
	})
	if sysErr := c.handleErr(response, err); sysErr != nil {
		return nil, sysErr
	}
	return responses.ToEnclaveResponse(response.EncodedEnclaveResponse), nil
}

func (c *Client) HealthCheck() (bool, common.SystemError) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.HealthCheck(timeoutCtx, &generated.EmptyArgs{})
	if sysErr := c.handleErr(response, err); sysErr != nil {
		return false, sysErr
	}
	return response.Status, nil
}

func (c *Client) GenerateRollup() (*common.ExtRollup, common.SystemError) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.CreateRollup(timeoutCtx, &generated.CreateRollupRequest{})
	if sysErr := c.handleErr(response, err); sysErr != nil {
		return nil, sysErr
	}
	return rpc.FromExtRollupMsg(response.Msg), nil
}

func (c *Client) DebugTraceTransaction(hash gethcommon.Hash, config *tracers.TraceConfig) (json.RawMessage, common.SystemError) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	confBytes, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	response, err := c.protoClient.DebugTraceTransaction(timeoutCtx, &generated.DebugTraceTransactionRequest{
		TxHash: hash.Bytes(),
		Config: confBytes,
	})
	if sysErr := c.handleErr(response, err); sysErr != nil {
		return nil, sysErr
	}
	return json.RawMessage(response.Msg), nil
}

func (c *Client) DebugEventLogRelevancy(hash gethcommon.Hash) (json.RawMessage, common.SystemError) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.DebugEventLogRelevancy(timeoutCtx, &generated.DebugEventLogRelevancyRequest{
		TxHash: hash.Bytes(),
	})
	if sysErr := c.handleErr(response, err); sysErr != nil {
		return nil, sysErr
	}
	return json.RawMessage(response.Msg), nil
}

// handleErr ensures to check for errors in both the RPC connection and in the RPC response object
// creates different errors depending on the error found
func (c *Client) handleErr(rpcResp interface{}, connError error) error {
	if connError != nil {
		c.logger.Error("enclave client RPC Connection err - ", log.ErrKey, connError)
		return syserr.NewRPCError(connError)
	}

	if rpcResp != nil {
		var internalErr *generated.SystemError
		switch typedErr := rpcResp.(type) {
		case *generated.AttestationResponse:
			internalErr = typedErr.GetSystemError()
		case *generated.StatusResponse:
			internalErr = typedErr.GetSystemError()
		case *generated.GenerateSecretResponse:
			internalErr = typedErr.GetSystemError()
		case *generated.SubmitTxResponse:
			internalErr = typedErr.GetSystemError()
		case *generated.ObsCallResponse:
			internalErr = typedErr.GetSystemError()
		case *generated.GetTransactionCountResponse:
			internalErr = typedErr.GetSystemError()
		case *generated.GetTransactionResponse:
			internalErr = typedErr.GetSystemError()
		case *generated.GetTransactionReceiptResponse:
			internalErr = typedErr.GetSystemError()
		case *generated.GetBalanceResponse:
			internalErr = typedErr.GetSystemError()
		case *generated.GetCodeResponse:
			internalErr = typedErr.GetSystemError()
		case *generated.EstimateGasResponse:
			internalErr = typedErr.GetSystemError()
		case *generated.GetLogsResponse:
			internalErr = typedErr.GetSystemError()
		case *generated.CreateRollupResponse:
			internalErr = typedErr.GetSystemError()
		case *generated.DebugTraceTransactionResponse:
			internalErr = typedErr.GetSystemError()
		case *generated.DebugEventLogRelevancyResponse:
			internalErr = typedErr.GetSystemError()
		case *generated.HealthCheckResponse:
			internalErr = typedErr.GetSystemError()
		default:
			panic("Unexpected RPC error - this should never happen")
		}

		if internalErr != nil {
			c.logger.Error("enclave client RPC Message err - ", log.ErrKey, internalErr.GetErrorString())
			return syserr.NewInternalError(fmt.Errorf("%s", internalErr))
		}
	}
	return nil
}
