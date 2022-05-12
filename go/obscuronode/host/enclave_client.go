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
	nodeID      common.Address
}

func NewEnclaveRPCClient(address string, timeout time.Duration, nodeID common.Address) *EnclaveRPCClient {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	connection, err := grpc.Dial(address, opts...)
	if err != nil {
		panic(fmt.Errorf(">   Agg%d: Failed to connect to enclave RPC service: %w", obscurocommon.ShortAddress(nodeID), err))
	}
	return &EnclaveRPCClient{generated.NewEnclaveProtoClient(connection), connection, timeout, nodeID}
}

func (c *EnclaveRPCClient) StopClient() {
	if err := c.connection.Close(); err != nil {
		panic(fmt.Errorf(">   Agg%d: Failed to stop enclave RPC service: %w", obscurocommon.ShortAddress(c.nodeID), err))
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
		panic(fmt.Errorf(">   Agg%d: Failed to retrieve attestation: %w", obscurocommon.ShortAddress(c.nodeID), err))
	}
	return rpc.FromAttestationReportMsg(response.AttestationReportMsg)
}

func (c *EnclaveRPCClient) GenerateSecret() obscurocommon.EncryptedSharedEnclaveSecret {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	response, err := c.protoClient.GenerateSecret(timeoutCtx, &generated.GenerateSecretRequest{})
	if err != nil {
		panic(fmt.Errorf(">   Agg%d: Failed to generate secret: %w", obscurocommon.ShortAddress(c.nodeID), err))
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
		panic(fmt.Errorf(">   Agg%d: Failed to fetch secret: %w", obscurocommon.ShortAddress(c.nodeID), err))
	}
	return response.EncryptedSharedEnclaveSecret
}

func (c *EnclaveRPCClient) InitEnclave(secret obscurocommon.EncryptedSharedEnclaveSecret) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	_, err := c.protoClient.InitEnclave(timeoutCtx, &generated.InitEnclaveRequest{EncryptedSharedEnclaveSecret: secret})
	if err != nil {
		panic(fmt.Errorf(">   Agg%d: Failed to initialise enclave: %w", obscurocommon.ShortAddress(c.nodeID), err))
	}
}

func (c *EnclaveRPCClient) IsInitialised() bool {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	response, err := c.protoClient.IsInitialised(timeoutCtx, &generated.IsInitialisedRequest{})
	if err != nil {
		panic(fmt.Errorf(">   Agg%d: Failed to establish enclave initialisation status: %w", obscurocommon.ShortAddress(c.nodeID), err))
	}
	return response.IsInitialised
}

func (c *EnclaveRPCClient) ProduceGenesis(blkHash common.Hash) nodecommon.BlockSubmissionResponse {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	response, err := c.protoClient.ProduceGenesis(timeoutCtx, &generated.ProduceGenesisRequest{BlockHash: blkHash.Bytes()})
	if err != nil {
		panic(fmt.Errorf(">   Agg%d: Failed to produce genesis: %w", obscurocommon.ShortAddress(c.nodeID), err))
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
		panic(fmt.Errorf(">   Agg%d: Failed to ingest blocks: %w", obscurocommon.ShortAddress(c.nodeID), err))
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
		panic(fmt.Errorf(">   Agg%d: Failed to encode block: %w", obscurocommon.ShortAddress(c.nodeID), err))
	}
	_, err := c.protoClient.Start(timeoutCtx, &generated.StartRequest{EncodedBlock: buffer.Bytes()})
	if err != nil {
		panic(fmt.Errorf(">   Agg%d: Failed to start enclave: %w", obscurocommon.ShortAddress(c.nodeID), err))
	}
}

func (c *EnclaveRPCClient) SubmitBlock(block types.Block) nodecommon.BlockSubmissionResponse {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	var buffer bytes.Buffer
	if err := block.EncodeRLP(&buffer); err != nil {
		panic(fmt.Errorf(">   Agg%d: Failed to encode block: %w", obscurocommon.ShortAddress(c.nodeID), err))
	}

	response, err := c.protoClient.SubmitBlock(timeoutCtx, &generated.SubmitBlockRequest{EncodedBlock: buffer.Bytes()})
	if err != nil {
		panic(fmt.Errorf(">   Agg%d: Failed to submit block: %w", obscurocommon.ShortAddress(c.nodeID), err))
	}
	return rpc.FromBlockSubmissionResponseMsg(response.BlockSubmissionResponse)
}

func (c *EnclaveRPCClient) SubmitRollup(rollup nodecommon.ExtRollup) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	extRollupMsg := rpc.ToExtRollupMsg(&rollup)
	_, err := c.protoClient.SubmitRollup(timeoutCtx, &generated.SubmitRollupRequest{ExtRollup: &extRollupMsg})
	if err != nil {
		panic(fmt.Errorf(">   Agg%d: Failed to submit rollup: %w", obscurocommon.ShortAddress(c.nodeID), err))
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
		panic(fmt.Errorf(">   Agg%d: Failed to retrieve balance: %w", obscurocommon.ShortAddress(c.nodeID), err))
	}
	return response.Balance
}

func (c *EnclaveRPCClient) RoundWinner(parent obscurocommon.L2RootHash) (nodecommon.ExtRollup, bool, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	response, err := c.protoClient.RoundWinner(timeoutCtx, &generated.RoundWinnerRequest{Parent: parent.Bytes()})
	if err != nil {
		panic(fmt.Errorf(">   Agg%d: Failed to determine round winner: %w", obscurocommon.ShortAddress(c.nodeID), err))
	}

	if response.Winner {
		return rpc.FromExtRollupMsg(response.ExtRollup), true, nil
	}
	return nodecommon.ExtRollup{}, false, nil
}

func (c *EnclaveRPCClient) Stop() error {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	_, err := c.protoClient.Stop(timeoutCtx, &generated.StopRequest{})
	if err != nil {
		return fmt.Errorf("failed to stop enclave: %w", err)
	}
	return nil
}

func (c *EnclaveRPCClient) GetTransaction(txHash common.Hash) *nodecommon.L2Tx {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	response, err := c.protoClient.GetTransaction(timeoutCtx, &generated.GetTransactionRequest{TxHash: txHash.Bytes()})
	if err != nil {
		panic(fmt.Errorf(">   Agg%d: Failed to get transaction: %w", obscurocommon.ShortAddress(c.nodeID), err))
	}

	if !response.Known {
		return nil
	}

	l2Tx := nodecommon.L2Tx{}
	err = l2Tx.DecodeRLP(rlp.NewStream(bytes.NewReader(response.EncodedTransaction), 0))
	if err != nil {
		panic(fmt.Errorf(">   Agg%d: Failed to decode transaction: %w", obscurocommon.ShortAddress(c.nodeID), err))
	}

	return &l2Tx
}
