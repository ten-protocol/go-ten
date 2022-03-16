package rpc

import (
	"bytes"
	"context"
	"fmt"
	"time"

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
	protoClient EnclaveProtoClient
	connection  *grpc.ClientConn
	timeout     time.Duration
}

func NewEnclaveClient(port uint64, timeout time.Duration) (EnclaveClient, error) {
	connection, err := getConnection(port)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to connect to enclave RPC service: %w", err)
		return EnclaveClient{}, wrappedErr
	}
	return EnclaveClient{NewEnclaveProtoClient(connection), connection, timeout}, nil
}

func (c *EnclaveClient) StopClient() error {
	return c.connection.Close()
}

// Returns an unsecured connection to use for communicating with the enclave over RPC.
func getConnection(port uint64) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	return grpc.Dial(fmt.Sprintf("localhost:%d", port), opts...)
}

func (c *EnclaveClient) Attestation() (obscurocommon.AttestationReport, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	response, err := c.protoClient.Attestation(timeoutCtx, &AttestationRequest{})
	return fromAttestationReportMsg(response.AttestationReportMsg), err
}

func (c *EnclaveClient) GenerateSecret() (obscurocommon.EncryptedSharedEnclaveSecret, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	response, err := c.protoClient.GenerateSecret(timeoutCtx, &GenerateSecretRequest{})
	return response.EncryptedSharedEnclaveSecret, err
}

func (c *EnclaveClient) FetchSecret(report obscurocommon.AttestationReport) (obscurocommon.EncryptedSharedEnclaveSecret, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	attestationReportMsg := toAttestationReportMsg(report)
	request := FetchSecretRequest{AttestationReportMsg: &attestationReportMsg}
	response, err := c.protoClient.FetchSecret(timeoutCtx, &request)
	return response.EncryptedSharedEnclaveSecret, err
}

func (c *EnclaveClient) Init(secret obscurocommon.EncryptedSharedEnclaveSecret) error {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	_, err := c.protoClient.Init(timeoutCtx, &InitRequest{EncryptedSharedEnclaveSecret: secret})
	return err
}

func (c *EnclaveClient) IsInitialised() (bool, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	response, err := c.protoClient.IsInitialised(timeoutCtx, &IsInitialisedRequest{})
	return response.IsInitialised, err
}

func (c *EnclaveClient) ProduceGenesis() (enclave.BlockSubmissionResponse, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	response, err := c.protoClient.ProduceGenesis(timeoutCtx, &ProduceGenesisRequest{})
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
	_, err := c.protoClient.IngestBlocks(timeoutCtx, &IngestBlocksRequest{EncodedBlocks: encodedBlocks})
	return err
}

func (c *EnclaveClient) Start(block types.Block) error {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	var buffer bytes.Buffer
	if err := block.EncodeRLP(&buffer); err != nil {
		return err
	}
	_, err := c.protoClient.Start(timeoutCtx, &StartRequest{EncodedBlock: buffer.Bytes()})
	return err
}

func (c *EnclaveClient) SubmitBlock(block types.Block) (enclave.BlockSubmissionResponse, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	var buffer bytes.Buffer
	if err := block.EncodeRLP(&buffer); err != nil {
		return enclave.BlockSubmissionResponse{}, err
	}

	response, err := c.protoClient.SubmitBlock(timeoutCtx, &SubmitBlockRequest{EncodedBlock: buffer.Bytes()})
	return fromBlockSubmissionResponseMsg(response.BlockSubmissionResponse), err
}

func (c *EnclaveClient) SubmitRollup(rollup nodecommon.ExtRollup) error {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	extRollupMsg := toExtRollupMsg(&rollup)
	_, err := c.protoClient.SubmitRollup(timeoutCtx, &SubmitRollupRequest{ExtRollup: &extRollupMsg})
	return err
}

func (c *EnclaveClient) SubmitTx(tx nodecommon.EncryptedTx) error {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	_, err := c.protoClient.SubmitTx(timeoutCtx, &SubmitTxRequest{EncryptedTx: tx})
	return err
}

func (c *EnclaveClient) Balance(address common.Address) (uint64, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	response, err := c.protoClient.Balance(timeoutCtx, &BalanceRequest{Address: address.Bytes()})
	return response.Balance, err
}

func (c *EnclaveClient) RoundWinner(parent obscurocommon.L2RootHash) (nodecommon.ExtRollup, bool, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	response, err := c.protoClient.RoundWinner(timeoutCtx, &RoundWinnerRequest{Parent: parent.Bytes()})
	if err == nil && response.Winner {
		return fromExtRollupMsg(response.ExtRollup), true, err
	}
	return nodecommon.ExtRollup{}, false, err
}

func (c *EnclaveClient) Stop() error {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	_, err := c.protoClient.Stop(timeoutCtx, &StopRequest{})
	return err
}

func (c *EnclaveClient) GetTransaction(txHash common.Hash) (*enclave.L2Tx, bool, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	response, err := c.protoClient.GetTransaction(timeoutCtx, &GetTransactionRequest{TxHash: txHash.Bytes()})
	if err != nil || !response.Known {
		return &enclave.L2Tx{}, false, err
	}

	l2Tx := enclave.L2Tx{}
	err = l2Tx.DecodeRLP(rlp.NewStream(bytes.NewReader(response.EncodedTransaction), 0))
	if err != nil {
		return &enclave.L2Tx{}, false, err
	}

	return &l2Tx, true, err
}
