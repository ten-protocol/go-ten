package rpc

import (
	"bytes"
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// TODO - Joel - Introduce client request timeouts.

// EnclaveClient is the implementation of EnclaveClient that should be used by the host when communicating with the enclave.
// Calls are proxied to the enclave process over RPC.
type EnclaveClient struct {
	protoClient EnclaveProtoClient
	connection  *grpc.ClientConn
}

func NewEnclaveClient(port uint64) (EnclaveClient, error) {
	connection, err := getConnection(port)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to connect to enclave RPC service: %w", err)
		return EnclaveClient{}, wrappedErr
	}
	return EnclaveClient{NewEnclaveProtoClient(connection), connection}, nil
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
	response, err := c.protoClient.Attestation(context.Background(), &AttestationRequest{})
	return fromAttestationReportMsg(response.AttestationReportMsg), err
}

func (c *EnclaveClient) GenerateSecret() (obscurocommon.EncryptedSharedEnclaveSecret, error) {
	response, err := c.protoClient.GenerateSecret(context.Background(), &GenerateSecretRequest{})
	return response.EncryptedSharedEnclaveSecret, err
}

func (c *EnclaveClient) FetchSecret(report obscurocommon.AttestationReport) (obscurocommon.EncryptedSharedEnclaveSecret, error) {
	attestationReportMsg := toAttestationReportMsg(report)
	request := FetchSecretRequest{AttestationReportMsg: &attestationReportMsg}
	response, err := c.protoClient.FetchSecret(context.Background(), &request)
	return response.EncryptedSharedEnclaveSecret, err
}

func (c *EnclaveClient) Init(secret obscurocommon.EncryptedSharedEnclaveSecret) error {
	_, err := c.protoClient.Init(context.Background(), &InitRequest{EncryptedSharedEnclaveSecret: secret})
	return err
}

func (c *EnclaveClient) IsInitialised() (bool, error) {
	response, err := c.protoClient.IsInitialised(context.Background(), &IsInitialisedRequest{})
	return response.IsInitialised, err
}

func (c *EnclaveClient) ProduceGenesis() (enclave.BlockSubmissionResponse, error) {
	response, err := c.protoClient.ProduceGenesis(context.Background(), &ProduceGenesisRequest{})
	return fromBlockSubmissionResponseMsg(response.BlockSubmissionResponse), err
}

func (c *EnclaveClient) IngestBlocks(blocks []*types.Block) error {
	encodedBlocks := make([][]byte, 0)
	for _, block := range blocks {
		encodedBlock := obscurocommon.EncodeBlock(block)
		encodedBlocks = append(encodedBlocks, encodedBlock)
	}
	_, err := c.protoClient.IngestBlocks(context.Background(), &IngestBlocksRequest{EncodedBlocks: encodedBlocks})
	return err
}

func (c *EnclaveClient) Start(block types.Block) error {
	var buffer bytes.Buffer
	if err := block.EncodeRLP(&buffer); err != nil {
		return err
	}
	_, err := c.protoClient.Start(context.Background(), &StartRequest{EncodedBlock: buffer.Bytes()})
	return err
}

func (c *EnclaveClient) SubmitBlock(block types.Block) (enclave.BlockSubmissionResponse, error) {
	var buffer bytes.Buffer
	if err := block.EncodeRLP(&buffer); err != nil {
		return enclave.BlockSubmissionResponse{}, err
	}

	response, err := c.protoClient.SubmitBlock(context.Background(), &SubmitBlockRequest{EncodedBlock: buffer.Bytes()})
	return fromBlockSubmissionResponseMsg(response.BlockSubmissionResponse), err
}

func (c *EnclaveClient) SubmitRollup(rollup nodecommon.ExtRollup) error {
	extRollupMsg := toExtRollupMsg(&rollup)
	_, err := c.protoClient.SubmitRollup(context.Background(), &SubmitRollupRequest{ExtRollup: &extRollupMsg})
	return err
}

func (c *EnclaveClient) SubmitTx(tx nodecommon.EncryptedTx) error {
	_, err := c.protoClient.SubmitTx(context.Background(), &SubmitTxRequest{EncryptedTx: tx})
	return err
}

func (c *EnclaveClient) Balance(address common.Address) (uint64, error) {
	response, err := c.protoClient.Balance(context.Background(), &BalanceRequest{Address: address.Bytes()})
	return response.Balance, err
}

func (c *EnclaveClient) RoundWinner(parent obscurocommon.L2RootHash) (nodecommon.ExtRollup, bool, error) {
	response, err := c.protoClient.RoundWinner(context.Background(), &RoundWinnerRequest{Parent: parent.Bytes()})
	if err == nil && response.Winner {
		return fromExtRollupMsg(response.ExtRollup), true, err
	}
	return nodecommon.ExtRollup{}, false, err
}

func (c *EnclaveClient) Stop() error {
	_, err := c.protoClient.Stop(context.Background(), &StopRequest{})
	return err
}

func (c *EnclaveClient) GetTransaction(txHash common.Hash) (*enclave.L2Tx, bool, error) {
	response, err := c.protoClient.GetTransaction(context.Background(), &GetTransactionRequest{TxHash: txHash.Bytes()})
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
