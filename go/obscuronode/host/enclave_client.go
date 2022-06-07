package host

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"time"

	"google.golang.org/grpc/connectivity"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/config"

	"github.com/obscuronet/obscuro-playground/go/log"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon/rpc"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon/rpc/generated"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// EnclaveRPCClient implements enclave.Enclave and should be used by the host when communicating with the enclave via RPC.
type EnclaveRPCClient struct {
	protoClient generated.EnclaveProtoClient
	connection  *grpc.ClientConn
	config      config.HostConfig
}

func NewEnclaveRPCClient(config config.HostConfig) *EnclaveRPCClient {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	connection, err := grpc.Dial(config.EnclaveRPCAddress, opts...)
	if err != nil {
		log.Panic(">   Agg%d: Failed to connect to enclave RPC service. Cause: %s", obscurocommon.ShortAddress(config.ID), err)
	}

	// We wait for the RPC connection to be ready.
	currentTime := time.Now()
	deadline := currentTime.Add(30 * time.Second)
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
		log.Panic(">   Agg%d: RPC connection failed to establish. Current state is %s", obscurocommon.ShortAddress(config.ID), currentState)
	}

	return &EnclaveRPCClient{generated.NewEnclaveProtoClient(connection), connection, config}
}

func (c *EnclaveRPCClient) StopClient() error {
	return c.connection.Close()
}

func (c *EnclaveRPCClient) IsReady() error {
	if c.connection.GetState() != connectivity.Ready {
		return errors.New("RPC connection is not ready")
	}

	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	resp, err := c.protoClient.IsReady(timeoutCtx, &generated.IsReadyRequest{})
	if err != nil {
		return err
	}
	if resp.GetError() != "" {
		return errors.New(resp.GetError())
	}
	return nil
}

func (c *EnclaveRPCClient) Attestation() *obscurocommon.AttestationReport {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.Attestation(timeoutCtx, &generated.AttestationRequest{})
	if err != nil {
		log.Panic(">   Agg%d: Failed to retrieve attestation. Cause: %s", obscurocommon.ShortAddress(c.config.ID), err)
	}
	return rpc.FromAttestationReportMsg(response.AttestationReportMsg)
}

func (c *EnclaveRPCClient) GenerateSecret() obscurocommon.EncryptedSharedEnclaveSecret {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.GenerateSecret(timeoutCtx, &generated.GenerateSecretRequest{})
	if err != nil {
		log.Panic(">   Agg%d: Failed to generate secret. Cause: %s", obscurocommon.ShortAddress(c.config.ID), err)
	}
	return response.EncryptedSharedEnclaveSecret
}

func (c *EnclaveRPCClient) ShareSecret(report *obscurocommon.AttestationReport) (obscurocommon.EncryptedSharedEnclaveSecret, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	attestationReportMsg := rpc.ToAttestationReportMsg(report)
	request := generated.FetchSecretRequest{AttestationReportMsg: &attestationReportMsg}
	response, err := c.protoClient.ShareSecret(timeoutCtx, &request)
	if err != nil {
		return nil, err
	}
	if response.GetError() != "" {
		return nil, errors.New(response.GetError())
	}
	return response.EncryptedSharedEnclaveSecret, nil
}

func (c *EnclaveRPCClient) InitEnclave(secret obscurocommon.EncryptedSharedEnclaveSecret) error {
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

func (c *EnclaveRPCClient) IsInitialised() bool {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.IsInitialised(timeoutCtx, &generated.IsInitialisedRequest{})
	if err != nil {
		log.Panic(">   Agg%d: Failed to establish enclave initialisation status. Cause: %s", obscurocommon.ShortAddress(c.config.ID), err)
	}
	return response.IsInitialised
}

func (c *EnclaveRPCClient) ProduceGenesis(blkHash common.Hash) nodecommon.BlockSubmissionResponse {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()
	response, err := c.protoClient.ProduceGenesis(timeoutCtx, &generated.ProduceGenesisRequest{BlockHash: blkHash.Bytes()})
	if err != nil {
		log.Panic(">   Agg%d: Failed to produce genesis. Cause: %s", obscurocommon.ShortAddress(c.config.ID), err)
	}
	return rpc.FromBlockSubmissionResponseMsg(response.BlockSubmissionResponse)
}

func (c *EnclaveRPCClient) IngestBlocks(blocks []*types.Block) []nodecommon.BlockSubmissionResponse {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	encodedBlocks := make([][]byte, 0)
	for _, block := range blocks {
		encodedBlock := obscurocommon.EncodeBlock(block)
		encodedBlocks = append(encodedBlocks, encodedBlock)
	}
	response, err := c.protoClient.IngestBlocks(timeoutCtx, &generated.IngestBlocksRequest{EncodedBlocks: encodedBlocks})
	if err != nil {
		log.Panic(">   Agg%d: Failed to ingest blocks. Cause: %s", obscurocommon.ShortAddress(c.config.ID), err)
	}
	responses := response.GetBlockSubmissionResponses()
	result := make([]nodecommon.BlockSubmissionResponse, len(responses))
	for i, r := range responses {
		result[i] = rpc.FromBlockSubmissionResponseMsg(r)
	}
	return result
}

func (c *EnclaveRPCClient) Start(block types.Block) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	var buffer bytes.Buffer
	if err := block.EncodeRLP(&buffer); err != nil {
		log.Panic(">   Agg%d: Failed to encode block. Cause: %s", obscurocommon.ShortAddress(c.config.ID), err)
	}
	_, err := c.protoClient.Start(timeoutCtx, &generated.StartRequest{EncodedBlock: buffer.Bytes()})
	if err != nil {
		log.Panic(">   Agg%d: Failed to start enclave. Cause: %s", obscurocommon.ShortAddress(c.config.ID), err)
	}
}

func (c *EnclaveRPCClient) SubmitBlock(block types.Block) nodecommon.BlockSubmissionResponse {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	var buffer bytes.Buffer
	if err := block.EncodeRLP(&buffer); err != nil {
		log.Panic(">   Agg%d: Failed to encode block. Cause: %s", obscurocommon.ShortAddress(c.config.ID), err)
	}

	response, err := c.protoClient.SubmitBlock(timeoutCtx, &generated.SubmitBlockRequest{EncodedBlock: buffer.Bytes()})
	if err != nil {
		log.Panic(">   Agg%d: Failed to submit block. Cause: %s", obscurocommon.ShortAddress(c.config.ID), err)
	}
	return rpc.FromBlockSubmissionResponseMsg(response.BlockSubmissionResponse)
}

func (c *EnclaveRPCClient) SubmitRollup(rollup nodecommon.ExtRollup) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	extRollupMsg := rpc.ToExtRollupMsg(&rollup)
	_, err := c.protoClient.SubmitRollup(timeoutCtx, &generated.SubmitRollupRequest{ExtRollup: &extRollupMsg})
	if err != nil {
		log.Panic(">   Agg%d: Failed to submit rollup. Cause: %s", obscurocommon.ShortAddress(c.config.ID), err)
	}
}

func (c *EnclaveRPCClient) SubmitTx(tx nodecommon.EncryptedTx) error {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	_, err := c.protoClient.SubmitTx(timeoutCtx, &generated.SubmitTxRequest{EncryptedTx: tx})
	return err
}

func (c *EnclaveRPCClient) ExecuteOffChainTransaction(from common.Address, contractAddress common.Address, data []byte) (nodecommon.EncryptedResponse, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.ExecuteOffChainTransaction(timeoutCtx, &generated.OffChainRequest{
		From:            from.Bytes(),
		ContractAddress: contractAddress.Bytes(),
		Data:            data,
	})
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}
	return response.Result, nil
}

func (c *EnclaveRPCClient) Nonce(address common.Address) uint64 {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.Nonce(timeoutCtx, &generated.NonceRequest{Address: address.Bytes()})
	if err != nil {
		panic(fmt.Errorf(">   Agg%d: Failed to retrieve nonce: %w", obscurocommon.ShortAddress(c.config.ID), err))
	}
	return response.Nonce
}

func (c *EnclaveRPCClient) RoundWinner(parent obscurocommon.L2RootHash) (nodecommon.ExtRollup, bool, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.RoundWinner(timeoutCtx, &generated.RoundWinnerRequest{Parent: parent.Bytes()})
	if err != nil {
		log.Panic(">   Agg%d: Failed to determine round winner. Cause: %s", obscurocommon.ShortAddress(c.config.ID), err)
	}

	if response.Winner {
		return rpc.FromExtRollupMsg(response.ExtRollup), true, nil
	}
	return nodecommon.ExtRollup{}, false, nil
}

func (c *EnclaveRPCClient) Stop() error {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	_, err := c.protoClient.Stop(timeoutCtx, &generated.StopRequest{})
	if err != nil {
		return fmt.Errorf("failed to stop enclave: %w", err)
	}
	return nil
}

func (c *EnclaveRPCClient) GetTransaction(txHash common.Hash) *nodecommon.L2Tx {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.GetTransaction(timeoutCtx, &generated.GetTransactionRequest{TxHash: txHash.Bytes()})
	if err != nil {
		log.Panic(">   Agg%d: Failed to retrieve transaction. Cause: %s", obscurocommon.ShortAddress(c.config.ID), err)
	}

	if !response.Known {
		return nil
	}

	l2Tx := nodecommon.L2Tx{}
	err = l2Tx.DecodeRLP(rlp.NewStream(bytes.NewReader(response.EncodedTransaction), 0))
	if err != nil {
		log.Panic(">   Agg%d: Failed to decode transaction. Cause: %s", obscurocommon.ShortAddress(c.config.ID), err)
	}

	return &l2Tx
}

func (c *EnclaveRPCClient) GetRollup(rollupHash obscurocommon.L2RootHash) *nodecommon.ExtRollup {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.GetRollup(timeoutCtx, &generated.GetRollupRequest{RollupHash: rollupHash.Bytes()})
	if err != nil {
		log.Panic(">   Agg%d: Failed to retrieve rollup. Cause: %s", obscurocommon.ShortAddress(c.config.ID), err)
	}

	if !response.Known {
		return nil
	}

	extRollup := rpc.FromExtRollupMsg(response.ExtRollup)
	return &extRollup
}
