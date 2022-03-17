package rpc

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/rpc/generated"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// EnclaveClient is the implementation of EnclaveClient that should be used by the host when communicating with the enclave.
// Calls are proxied to the enclave process over RPC.
type EnclaveClient struct {
	protoClient generated.EnclaveProtoClient
	connection  *grpc.ClientConn
	timeout     time.Duration
}

func NewEnclaveClient(port uint64, timeout time.Duration) (*EnclaveClient, error) {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	connection, err := grpc.Dial(fmt.Sprintf("localhost:%d", port), opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to enclave RPC service: %w", err)
	}
	return &EnclaveClient{generated.NewEnclaveProtoClient(connection), connection, timeout}, nil
}

func (c *EnclaveClient) StopClient() error {
	return c.connection.Close()
}

func (c *EnclaveClient) IsReady() error {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	_, err := c.protoClient.IsReady(timeoutCtx, &generated.IsReadyRequest{})
	return err
}

func (c *EnclaveClient) Attestation() (obscurocommon.AttestationReport, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	response, err := c.protoClient.Attestation(timeoutCtx, &generated.AttestationRequest{})
	return fromAttestationReportMsg(response.AttestationReportMsg), err
}

func (c *EnclaveClient) GenerateSecret() (obscurocommon.EncryptedSharedEnclaveSecret, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	response, err := c.protoClient.GenerateSecret(timeoutCtx, &generated.GenerateSecretRequest{})
	return response.EncryptedSharedEnclaveSecret, err
}

func (c *EnclaveClient) FetchSecret(report obscurocommon.AttestationReport) (obscurocommon.EncryptedSharedEnclaveSecret, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	attestationReportMsg := toAttestationReportMsg(report)
	request := generated.FetchSecretRequest{AttestationReportMsg: &attestationReportMsg}
	response, err := c.protoClient.FetchSecret(timeoutCtx, &request)
	return response.EncryptedSharedEnclaveSecret, err
}

func (c *EnclaveClient) InitEnclave(secret obscurocommon.EncryptedSharedEnclaveSecret) error {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	_, err := c.protoClient.InitEnclave(timeoutCtx, &generated.InitEnclaveRequest{EncryptedSharedEnclaveSecret: secret})
	return err
}

func (c *EnclaveClient) IsInitialised() (bool, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	response, err := c.protoClient.IsInitialised(timeoutCtx, &generated.IsInitialisedRequest{})
	return response.IsInitialised, err
}

func (c *EnclaveClient) ProduceGenesis() (enclave.BlockSubmissionResponse, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	response, err := c.protoClient.ProduceGenesis(timeoutCtx, &generated.ProduceGenesisRequest{})
	return fromBlockSubmissionResponseMsg(response.BlockSubmissionResponse), err
}

func (c *EnclaveClient) IngestBlocks(blocks []*types.Block) error {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	encodedBlocks := make([][]byte, 0)
	for _, block := range blocks {
		encodedBlock := obscurocommon.EncodeBlock(block)
		encodedBlocks = append(encodedBlocks, encodedBlock)
	}
	_, err := c.protoClient.IngestBlocks(timeoutCtx, &generated.IngestBlocksRequest{EncodedBlocks: encodedBlocks})
	return err
}

func (c *EnclaveClient) Start(block types.Block) error {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	var buffer bytes.Buffer
	if err := block.EncodeRLP(&buffer); err != nil {
		return err
	}
	_, err := c.protoClient.Start(timeoutCtx, &generated.StartRequest{EncodedBlock: buffer.Bytes()})
	return err
}

func (c *EnclaveClient) SubmitBlock(block types.Block) (enclave.BlockSubmissionResponse, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	var buffer bytes.Buffer
	if err := block.EncodeRLP(&buffer); err != nil {
		return enclave.BlockSubmissionResponse{}, err
	}

	response, err := c.protoClient.SubmitBlock(timeoutCtx, &generated.SubmitBlockRequest{EncodedBlock: buffer.Bytes()})
	return fromBlockSubmissionResponseMsg(response.BlockSubmissionResponse), err
}

func (c *EnclaveClient) SubmitRollup(rollup nodecommon.ExtRollup) error {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	extRollupMsg := toExtRollupMsg(&rollup)
	_, err := c.protoClient.SubmitRollup(timeoutCtx, &generated.SubmitRollupRequest{ExtRollup: &extRollupMsg})
	return err
}

func (c *EnclaveClient) SubmitTx(tx nodecommon.EncryptedTx) error {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	_, err := c.protoClient.SubmitTx(timeoutCtx, &generated.SubmitTxRequest{EncryptedTx: tx})
	return err
}

func (c *EnclaveClient) Balance(address common.Address) (uint64, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	response, err := c.protoClient.Balance(timeoutCtx, &generated.BalanceRequest{Address: address.Bytes()})
	return response.Balance, err
}

func (c *EnclaveClient) RoundWinner(parent obscurocommon.L2RootHash) (nodecommon.ExtRollup, bool, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	response, err := c.protoClient.RoundWinner(timeoutCtx, &generated.RoundWinnerRequest{Parent: parent.Bytes()})
	if err == nil && response.Winner {
		return fromExtRollupMsg(response.ExtRollup), true, err
	}
	return nodecommon.ExtRollup{}, false, err
}

func (c *EnclaveClient) Stop() error {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	_, err := c.protoClient.Stop(timeoutCtx, &generated.StopRequest{})
	return err
}

func (c *EnclaveClient) GetTransaction(txHash common.Hash) (*enclave.L2Tx, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	response, err := c.protoClient.GetTransaction(timeoutCtx, &generated.GetTransactionRequest{TxHash: txHash.Bytes()})
	if err != nil || !response.Known {
		return nil, err
	}

	l2Tx := enclave.L2Tx{}
	err = l2Tx.DecodeRLP(rlp.NewStream(bytes.NewReader(response.EncodedTransaction), 0))
	if err != nil {
		return nil, err
	}

	return &l2Tx, err
}
