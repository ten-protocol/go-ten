package host

import (
	"bytes"
	"context"
	"fmt"
	"time"

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
	timeout     time.Duration
}

func NewEnclaveRPCClient(address string, timeout time.Duration) *EnclaveRPCClient {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	connection, err := grpc.Dial(address, opts...)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to enclave RPC service: %v", err))
	}
	return &EnclaveRPCClient{generated.NewEnclaveProtoClient(connection), connection, timeout}
}

func (c *EnclaveRPCClient) StopClient() {
	if err := c.connection.Close(); err != nil {
		panic(fmt.Sprintf("failed to stop enclave RPC service: %v", err))
	}
}

func (c *EnclaveRPCClient) IsReady() error {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	_, err := c.protoClient.IsReady(timeoutCtx, &generated.IsReadyRequest{})
	return err
}

func (c *EnclaveRPCClient) Attestation() obscurocommon.AttestationReport {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	response, err := c.protoClient.Attestation(timeoutCtx, &generated.AttestationRequest{})
	if err != nil {
		panic(fmt.Sprintf("failed to retrieve attestation: %v", err))
	}
	return rpc.FromAttestationReportMsg(response.AttestationReportMsg)
}

func (c *EnclaveRPCClient) GenerateSecret() obscurocommon.EncryptedSharedEnclaveSecret {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	response, err := c.protoClient.GenerateSecret(timeoutCtx, &generated.GenerateSecretRequest{})
	if err != nil {
		panic(fmt.Sprintf("failed to generate secret: %v", err))
	}
	return response.EncryptedSharedEnclaveSecret
}

func (c *EnclaveRPCClient) FetchSecret(report obscurocommon.AttestationReport) obscurocommon.EncryptedSharedEnclaveSecret {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	attestationReportMsg := rpc.ToAttestationReportMsg(report)
	request := generated.FetchSecretRequest{AttestationReportMsg: &attestationReportMsg}
	response, err := c.protoClient.FetchSecret(timeoutCtx, &request)
	if err != nil {
		panic(fmt.Sprintf("failed to fetch secret: %v", err))
	}
	return response.EncryptedSharedEnclaveSecret
}

func (c *EnclaveRPCClient) InitEnclave(secret obscurocommon.EncryptedSharedEnclaveSecret) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	_, err := c.protoClient.InitEnclave(timeoutCtx, &generated.InitEnclaveRequest{EncryptedSharedEnclaveSecret: secret})
	if err != nil {
		panic(fmt.Sprintf("failed to initialise enclave: %v", err))
	}
}

func (c *EnclaveRPCClient) IsInitialised() bool {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	response, err := c.protoClient.IsInitialised(timeoutCtx, &generated.IsInitialisedRequest{})
	if err != nil {
		panic(fmt.Sprintf("failed to establish enclave initialisation status: %v", err))
	}
	return response.IsInitialised
}

func (c *EnclaveRPCClient) ProduceGenesis() nodecommon.BlockSubmissionResponse {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	response, err := c.protoClient.ProduceGenesis(timeoutCtx, &generated.ProduceGenesisRequest{})
	if err != nil {
		panic(fmt.Sprintf("failed to produce genesis: %v", err))
	}
	return rpc.FromBlockSubmissionResponseMsg(response.BlockSubmissionResponse)
}

func (c *EnclaveRPCClient) IngestBlocks(blocks []*types.Block) []nodecommon.BlockSubmissionResponse {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	encodedBlocks := make([][]byte, 0)
	for _, block := range blocks {
		encodedBlock := obscurocommon.EncodeBlock(block)
		encodedBlocks = append(encodedBlocks, encodedBlock)
	}
	response, err := c.protoClient.IngestBlocks(timeoutCtx, &generated.IngestBlocksRequest{EncodedBlocks: encodedBlocks})
	if err != nil {
		panic(fmt.Sprintf("failed to ingest blocks: %v", err))
	}
	responses := response.GetBlockSubmissionResponses()
	result := make([]nodecommon.BlockSubmissionResponse, len(responses))
	for i, r := range responses {
		result[i] = rpc.FromBlockSubmissionResponseMsg(r)
	}
	return result
}

func (c *EnclaveRPCClient) Start(block types.Block) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	var buffer bytes.Buffer
	if err := block.EncodeRLP(&buffer); err != nil {
		panic(fmt.Sprintf("failed to encode block: %v", err))
	}
	_, err := c.protoClient.Start(timeoutCtx, &generated.StartRequest{EncodedBlock: buffer.Bytes()})
	if err != nil {
		panic(fmt.Sprintf("failed to start enclave: %v", err))
	}
}

func (c *EnclaveRPCClient) SubmitBlock(block types.Block) nodecommon.BlockSubmissionResponse {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	var buffer bytes.Buffer
	if err := block.EncodeRLP(&buffer); err != nil {
		panic(fmt.Sprintf("failed to encode block: %v", err))
	}

	response, err := c.protoClient.SubmitBlock(timeoutCtx, &generated.SubmitBlockRequest{EncodedBlock: buffer.Bytes()})
	if err != nil {
		panic(fmt.Sprintf("failed to submit block: %v", err))
	}
	return rpc.FromBlockSubmissionResponseMsg(response.BlockSubmissionResponse)
}

func (c *EnclaveRPCClient) SubmitRollup(rollup nodecommon.ExtRollup) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	extRollupMsg := rpc.ToExtRollupMsg(&rollup)
	_, err := c.protoClient.SubmitRollup(timeoutCtx, &generated.SubmitRollupRequest{ExtRollup: &extRollupMsg})
	if err != nil {
		panic(fmt.Sprintf("failed to submit rollup: %v", err))
	}
}

func (c *EnclaveRPCClient) SubmitTx(tx nodecommon.EncryptedTx) error {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	_, err := c.protoClient.SubmitTx(timeoutCtx, &generated.SubmitTxRequest{EncryptedTx: tx})
	return err
}

func (c *EnclaveRPCClient) Balance(address common.Address) uint64 {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	response, err := c.protoClient.Balance(timeoutCtx, &generated.BalanceRequest{Address: address.Bytes()})
	if err != nil {
		panic(fmt.Sprintf("failed to retrieve balance: %v", err))
	}
	return response.Balance
}

func (c *EnclaveRPCClient) RoundWinner(parent obscurocommon.L2RootHash) (nodecommon.ExtRollup, bool, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	response, err := c.protoClient.RoundWinner(timeoutCtx, &generated.RoundWinnerRequest{Parent: parent.Bytes()})
	if err != nil {
		panic(fmt.Sprintf("failed to determine round winner: %v", err))
	}

	if response.Winner {
		return rpc.FromExtRollupMsg(response.ExtRollup), true, nil
	}
	return nodecommon.ExtRollup{}, false, nil
}

func (c *EnclaveRPCClient) Stop() {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	_, err := c.protoClient.Stop(timeoutCtx, &generated.StopRequest{})
	if err != nil {
		panic(fmt.Sprintf("failed to stop enclave: %v", err))
	}
}

func (c *EnclaveRPCClient) GetTransaction(txHash common.Hash) *nodecommon.L2Tx {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	response, err := c.protoClient.GetTransaction(timeoutCtx, &generated.GetTransactionRequest{TxHash: txHash.Bytes()})
	if err != nil {
		panic(fmt.Sprintf("failed to get transaction: %v", err))
	}

	if !response.Known {
		return nil
	}

	l2Tx := nodecommon.L2Tx{}
	err = l2Tx.DecodeRLP(rlp.NewStream(bytes.NewReader(response.EncodedTransaction), 0))
	if err != nil {
		panic(fmt.Sprintf("failed to decode transaction: %v", err))
	}

	return &l2Tx
}
