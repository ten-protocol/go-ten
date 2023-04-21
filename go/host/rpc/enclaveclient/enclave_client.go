package enclaveclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/common/retry"
	"github.com/obscuronet/go-obscuro/go/common/rpc"
	"github.com/obscuronet/go-obscuro/go/common/rpc/generated"
	"github.com/obscuronet/go-obscuro/go/common/tracers"
	"github.com/obscuronet/go-obscuro/go/config"
	"github.com/obscuronet/go-obscuro/go/responses"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
	gethrpc "github.com/ethereum/go-ethereum/rpc"
	grpc "google.golang.org/grpc"
)

func NewEnclaveRPCClient(config *config.HostConfig, logger gethlog.Logger) common.Enclave {
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

	protoClient := generated.NewEnclaveProtoClient(connection)
	return &EnclaveRPCClient{
		protoClient: protoClient,
		connection:  connection,
		config:      config,
		logger:      logger,
	}
}

type EnclaveRPCClient struct {
	protoClient generated.EnclaveProtoClient
	connection  *grpc.ClientConn
	config      *config.HostConfig
	logger      gethlog.Logger
}

func (c *EnclaveRPCClient) SubmitTx(tx common.EncryptedTx) responses.RawTx {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.SubmitTx(timeoutCtx, &generated.SubmitTxRequest{EncryptedTx: tx})
	if err != nil {
		return responses.AsPlaintextError(err)
	}

	return *responses.ToEnclaveResponse(response.EncodedEnclaveResponse)
}

func (c *EnclaveRPCClient) ObsCall(encryptedParams common.EncryptedParamsCall) responses.Call {
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

func (c *EnclaveRPCClient) GetTransactionCount(encryptedParams common.EncryptedParamsGetTxCount) responses.TxCount {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.GetTransactionCount(timeoutCtx, &generated.GetTransactionCountRequest{EncryptedParams: encryptedParams})
	if err != nil {
		return responses.AsPlaintextError(err)
	}
	return *responses.ToEnclaveResponse(response.EncodedEnclaveResponse)
}

func (c *EnclaveRPCClient) GetTransaction(encryptedParams common.EncryptedParamsGetTxByHash) responses.TxByHash {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	resp, err := c.protoClient.GetTransaction(timeoutCtx, &generated.GetTransactionRequest{EncryptedParams: encryptedParams})
	if err != nil {
		return responses.AsPlaintextError(err)
	}
	return *responses.ToEnclaveResponse(resp.EncodedEnclaveResponse)
}

func (c *EnclaveRPCClient) GetTransactionReceipt(encryptedParams common.EncryptedParamsGetTxReceipt) responses.TxReceipt {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.GetTransactionReceipt(timeoutCtx, &generated.GetTransactionReceiptRequest{EncryptedParams: encryptedParams})
	if err != nil {
		return responses.AsPlaintextError(err)
	}
	return *responses.ToEnclaveResponse(response.EncodedEnclaveResponse)
}

func (c *EnclaveRPCClient) AddViewingKey(viewingKeyBytes []byte, signature []byte) error {
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

func (c *EnclaveRPCClient) GetBalance(encryptedParams common.EncryptedParamsGetBalance) responses.Balance {
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

func (c *EnclaveRPCClient) GetCode(address gethcommon.Address, batchHash *gethcommon.Hash) ([]byte, error) {
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

func (c *EnclaveRPCClient) Subscribe(id gethrpc.ID, encryptedParams common.EncryptedParamsLogSubscription) error {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	_, err := c.protoClient.Subscribe(timeoutCtx, &generated.SubscribeRequest{
		Id:                    []byte(id),
		EncryptedSubscription: encryptedParams,
	})
	return err
}

func (c *EnclaveRPCClient) Unsubscribe(id gethrpc.ID) error {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	_, err := c.protoClient.Unsubscribe(timeoutCtx, &generated.UnsubscribeRequest{
		Id: []byte(id),
	})
	return err
}

func (c *EnclaveRPCClient) EstimateGas(encryptedParams common.EncryptedParamsEstimateGas) responses.Gas {
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

func (c *EnclaveRPCClient) GetLogs(encryptedParams common.EncryptedParamsGetLogs) responses.Logs {
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

func (c *EnclaveRPCClient) HealthCheck() (bool, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	resp, err := c.protoClient.HealthCheck(timeoutCtx, &generated.EmptyArgs{})
	if err != nil {
		return false, err
	}
	return resp.Status, nil
}

func (c *EnclaveRPCClient) DebugTraceTransaction(hash gethcommon.Hash, config *tracers.TraceConfig) (json.RawMessage, error) {
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

func (c *EnclaveRPCClient) DebugEventLogRelevancy(hash gethcommon.Hash) (json.RawMessage, error) {
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

func (c *EnclaveRPCClient) StopClient() error {
	return c.connection.Close()
}

func (c *EnclaveRPCClient) Status() (common.Status, error) {
	if c.connection.GetState() != connectivity.Ready {
		return common.Unavailable, errors.New("RPC connection is not ready")
	}

	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	resp, err := c.protoClient.Status(timeoutCtx, &generated.StatusRequest{})
	if err != nil {
		return common.Unavailable, err
	}
	if resp.GetError() != "" {
		return common.Unavailable, errors.New(resp.GetError())
	}
	return common.Status(resp.GetStatus()), nil
}

func (c *EnclaveRPCClient) Attestation() (*common.AttestationReport, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.Attestation(timeoutCtx, &generated.AttestationRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve attestation. Cause: %w", err)
	}
	return rpc.FromAttestationReportMsg(response.AttestationReportMsg), nil
}

func (c *EnclaveRPCClient) GenerateSecret() (common.EncryptedSharedEnclaveSecret, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.GenerateSecret(timeoutCtx, &generated.GenerateSecretRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to generate secret. Cause: %w", err)
	}
	return response.EncryptedSharedEnclaveSecret, nil
}

func (c *EnclaveRPCClient) InitEnclave(secret common.EncryptedSharedEnclaveSecret) error {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	resp, err := c.protoClient.InitEnclave(timeoutCtx, &generated.InitEnclaveRequest{EncryptedSharedEnclaveSecret: secret})
	if err != nil {
		return err
	}
	if resp.GetError() != "" {
		return errors.New(resp.GetError())
	}
	return nil
}

func (c *EnclaveRPCClient) SubmitL1Block(block types.Block, receipts types.Receipts, isLatest bool) (*common.BlockSubmissionResponse, error) {
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

func (c *EnclaveRPCClient) SubmitBatch(batch *common.ExtBatch) error {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	batchMsg := rpc.ToExtBatchMsg(batch)
	_, err := c.protoClient.SubmitBatch(timeoutCtx, &generated.SubmitBatchRequest{Batch: &batchMsg})
	if err != nil {
		return err
	}
	return nil
}

func (c *EnclaveRPCClient) Stop() error {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	_, err := c.protoClient.Stop(timeoutCtx, &generated.StopRequest{})
	if err != nil {
		return fmt.Errorf("could not stop enclave: %w", err)
	}
	return nil
}

func (c *EnclaveRPCClient) GenerateRollup() (*common.ExtRollup, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	resp, err := c.protoClient.CreateRollup(timeoutCtx, &generated.CreateRollupRequest{})
	if err != nil {
		return nil, err
	}
	return rpc.FromExtRollupMsg(resp.Msg), nil
}
