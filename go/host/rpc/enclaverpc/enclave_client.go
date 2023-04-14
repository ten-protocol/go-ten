package enclaverpc

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
	"github.com/obscuronet/go-obscuro/go/common/rpc"
	"github.com/obscuronet/go-obscuro/go/common/rpc/generated"
	"github.com/obscuronet/go-obscuro/go/common/tracers"
	"github.com/obscuronet/go-obscuro/go/config"
	"github.com/obscuronet/go-obscuro/go/enclave/evm"

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

	// We wait for the RPC connection to be ready.
	currentTime := time.Now()
	deadline := currentTime.Add(60 * time.Second)
	currentState := connection.GetState()
	for currentState == connectivity.Idle || currentState == connectivity.Connecting || currentState == connectivity.TransientFailure {
		connection.Connect()
		if time.Now().After(deadline) {
			break
		}
		time.Sleep(500 * time.Millisecond)
		currentState = connection.GetState()
	}

	if currentState != connectivity.Ready {
		logger.Crit(fmt.Sprintf("RPC connection failed to establish. Current state is %s", currentState))
	}

	return &Client{
		protoClient: generated.NewEnclaveProtoClient(connection),
		connection:  connection,
		config:      config,
		logger:      logger,
	}
}

func (c *Client) StopClient() error {
	return c.connection.Close()
}

func (c *Client) Status() (common.Status, error) {
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

func (c *Client) Attestation() (*common.AttestationReport, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.Attestation(timeoutCtx, &generated.AttestationRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve attestation. Cause: %w", err)
	}
	return rpc.FromAttestationReportMsg(response.AttestationReportMsg), nil
}

func (c *Client) GenerateSecret() (common.EncryptedSharedEnclaveSecret, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.GenerateSecret(timeoutCtx, &generated.GenerateSecretRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to generate secret. Cause: %w", err)
	}
	return response.EncryptedSharedEnclaveSecret, nil
}

func (c *Client) InitEnclave(secret common.EncryptedSharedEnclaveSecret) error {
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

func (c *Client) SubmitL1Block(block types.Block, receipts types.Receipts, isLatest bool) (*common.BlockSubmissionResponse, error) {
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

func (c *Client) SubmitTx(tx common.EncryptedTx) (common.EncryptedResponseSendRawTx, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.SubmitTx(timeoutCtx, &generated.SubmitTxRequest{EncryptedTx: tx})
	if err != nil {
		return nil, err
	}
	return response.EncryptedHash, err
}

func (c *Client) SubmitBatch(batch *common.ExtBatch) error {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	batchMsg := rpc.ToExtBatchMsg(batch)
	_, err := c.protoClient.SubmitBatch(timeoutCtx, &generated.SubmitBatchRequest{Batch: &batchMsg})
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) ObsCall(encryptedParams common.EncryptedParamsCall) (common.EncryptedResponseCall, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.ObsCall(timeoutCtx, &generated.ObsCallRequest{
		EncryptedParams: encryptedParams,
	})
	if err != nil {
		return nil, err
	}
	if len(response.Error) > 0 {
		// The enclave always returns a SerialisableError
		var result evm.SerialisableError
		err = json.Unmarshal(response.Error, &result)
		if err != nil {
			return nil, err
		}
		return nil, result
	}
	return response.Result, nil
}

func (c *Client) GetTransactionCount(encryptedParams common.EncryptedParamsGetTxCount) (common.EncryptedResponseGetTxCount, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.GetTransactionCount(timeoutCtx, &generated.GetTransactionCountRequest{EncryptedParams: encryptedParams})
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}
	return response.Result, nil
}

func (c *Client) Stop() error {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	_, err := c.protoClient.Stop(timeoutCtx, &generated.StopRequest{})
	if err != nil {
		return fmt.Errorf("could not stop enclave: %w", err)
	}
	return nil
}

func (c *Client) GetTransaction(encryptedParams common.EncryptedParamsGetTxByHash) (common.EncryptedResponseGetTxByHash, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	resp, err := c.protoClient.GetTransaction(timeoutCtx, &generated.GetTransactionRequest{EncryptedParams: encryptedParams})
	if err != nil {
		return nil, err
	}
	return resp.EncryptedTx, nil
}

func (c *Client) GetTransactionReceipt(encryptedParams common.EncryptedParamsGetTxReceipt) (common.EncryptedResponseGetTxReceipt, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.GetTransactionReceipt(timeoutCtx, &generated.GetTransactionReceiptRequest{EncryptedParams: encryptedParams})
	if err != nil {
		return nil, err
	}
	return response.EncryptedTxReceipt, nil
}

func (c *Client) AddViewingKey(viewingKeyBytes []byte, signature []byte) error {
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

func (c *Client) GetBalance(encryptedParams common.EncryptedParamsGetBalance) (common.EncryptedResponseGetBalance, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	resp, err := c.protoClient.GetBalance(timeoutCtx, &generated.GetBalanceRequest{
		EncryptedParams: encryptedParams,
	})
	if err != nil {
		return nil, err
	}
	return resp.EncryptedBalance, nil
}

func (c *Client) GetCode(address gethcommon.Address, batchHash *gethcommon.Hash) ([]byte, error) {
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

func (c *Client) Subscribe(id gethrpc.ID, encryptedParams common.EncryptedParamsLogSubscription) error {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	_, err := c.protoClient.Subscribe(timeoutCtx, &generated.SubscribeRequest{
		Id:                    []byte(id),
		EncryptedSubscription: encryptedParams,
	})
	return err
}

func (c *Client) Unsubscribe(id gethrpc.ID) error {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	_, err := c.protoClient.Unsubscribe(timeoutCtx, &generated.UnsubscribeRequest{
		Id: []byte(id),
	})
	return err
}

func (c *Client) EstimateGas(encryptedParams common.EncryptedParamsEstimateGas) (common.EncryptedResponseEstimateGas, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	resp, err := c.protoClient.EstimateGas(timeoutCtx, &generated.EstimateGasRequest{
		EncryptedParams: encryptedParams,
	})
	if err != nil {
		return nil, err
	}
	return resp.EncryptedResponse, nil
}

func (c *Client) GetLogs(encryptedParams common.EncryptedParamsGetLogs) (common.EncryptedResponseGetLogs, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	resp, err := c.protoClient.GetLogs(timeoutCtx, &generated.GetLogsRequest{
		EncryptedParams: encryptedParams,
	})
	if err != nil {
		return nil, err
	}
	return resp.EncryptedResponse, nil
}

func (c *Client) HealthCheck() (bool, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	resp, err := c.protoClient.HealthCheck(timeoutCtx, &generated.EmptyArgs{})
	if err != nil {
		return false, err
	}
	return resp.Status, nil
}

func (c *Client) CreateBatch() (*common.ExtBatch, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout+10*time.Second)
	defer cancel()

	resp, err := c.protoClient.CreateBatch(timeoutCtx, &generated.CreateBatchRequest{})
	if err != nil {
		return nil, err
	}
	return rpc.FromExtBatchMsg(resp.Msg), nil
}

func (c *Client) CreateRollup() (*common.ExtRollup, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout+10*time.Second)
	defer cancel()

	resp, err := c.protoClient.CreateRollup(timeoutCtx, &generated.CreateRollupRequest{})
	if err != nil {
		return nil, err
	}
	return rpc.FromExtRollupMsg(resp.Msg), nil
}

func (c *Client) DebugTraceTransaction(hash gethcommon.Hash, config *tracers.TraceConfig) (json.RawMessage, error) {
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

func (c *Client) StreamBatches(from *common.L2BatchHash) chan common.StreamBatchResponse {
	batchChan := make(chan common.StreamBatchResponse, 10)
	timeoutCtx, cancel := context.WithCancel(context.Background())

	request := &generated.StreamBatchesRequest{}
	if from != nil {
		request.KnownHead = from.Bytes()
	}

	stream, err := c.protoClient.StreamBatches(timeoutCtx, request)

	if err != nil {
		c.logger.Error("Error opening batch stream.", log.ErrKey, err)
		close(batchChan)
		cancel()
		return batchChan
	}

	go func() {
		defer cancel()
		for {
			batchMsg, err := stream.Recv()
			if err != nil {
				c.logger.Error("Error receving batch from stream.", log.ErrKey, err)
				//log error?
				close(batchChan)
				break
			}

			var decoded common.StreamBatchResponse
			if err := json.Unmarshal(batchMsg.Batch, &decoded); err != nil {
				c.logger.Error("Error unmarshling batch from stream.", log.ErrKey, err)
				close(batchChan)
				break
			}

			batchChan <- decoded
		}
	}()

	return batchChan
}
