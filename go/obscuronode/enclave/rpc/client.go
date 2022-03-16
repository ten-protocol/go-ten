package rpc

import (
	"bytes"
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/obscuro-playground/go/log"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// TODO - Joel - Introduce client request timeouts.

// EnclaveClient is the implementation of Enclave that should be used by the host when communicating with the enclave.
// Calls are proxied to the enclave process over RPC.
type EnclaveClient struct {
	protoClient EnclaveProtoClient
	connection  *grpc.ClientConn
}

func NewEnclaveClient(port uint64) EnclaveClient {
	// TODO - Joel - Handle error.
	connection, _ := getConnection(port)
	client := EnclaveClient{NewEnclaveProtoClient(connection), connection}
	return client
}

func (c *EnclaveClient) StopClient() error {
	return c.connection.Close()
}

// Returns an unsecured connection to use for communicating with the enclave over RPC.
func getConnection(port uint64) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	return grpc.Dial(fmt.Sprintf("localhost:%d", port), opts...)
}

func (c *EnclaveClient) Attestation() obscurocommon.AttestationReport {
	response, err := c.protoClient.Attestation(context.Background(), &AttestationRequest{})
	if err != nil {
		log.Log(fmt.Sprintf("failed to retrieve attestion report from enclave: %v", err))
	}
	return fromAttestationReportMsg(response.AttestationReportMsg)
}

func (c *EnclaveClient) GenerateSecret() obscurocommon.EncryptedSharedEnclaveSecret {
	response, err := c.protoClient.GenerateSecret(context.Background(), &GenerateSecretRequest{})
	if err != nil {
		log.Log(fmt.Sprintf("failed to generate enclave secret: %v", err))
	}
	return response.EncryptedSharedEnclaveSecret
}

func (c *EnclaveClient) FetchSecret(report obscurocommon.AttestationReport) obscurocommon.EncryptedSharedEnclaveSecret {
	attestationReportMsg := toAttestationReportMsg(report)
	request := FetchSecretRequest{AttestationReportMsg: &attestationReportMsg}
	response, err := c.protoClient.FetchSecret(context.Background(), &request)
	if err != nil {
		log.Log(fmt.Sprintf("failed to fetch enclave secret: %v", err))
	}
	return response.EncryptedSharedEnclaveSecret
}

func (c *EnclaveClient) Init(secret obscurocommon.EncryptedSharedEnclaveSecret) {
	_, err := c.protoClient.Init(context.Background(), &InitRequest{EncryptedSharedEnclaveSecret: secret})
	if err != nil {
		log.Log(fmt.Sprintf("failed to initialise enclave: %v", err))
	}
}

func (c *EnclaveClient) IsInitialised() bool {
	response, err := c.protoClient.IsInitialised(context.Background(), &IsInitialisedRequest{})
	if err != nil {
		log.Log(fmt.Sprintf("failed to check whether enclave is initialised: %v", err))
	}
	return response.IsInitialised
}

func (c *EnclaveClient) ProduceGenesis() enclave.BlockSubmissionResponse {
	response, err := c.protoClient.ProduceGenesis(context.Background(), &ProduceGenesisRequest{})
	if err != nil {
		log.Log(fmt.Sprintf("failed to produce genesis block: %v", err))
	}
	return fromBlockSubmissionResponseMsg(response.BlockSubmissionResponse)
}

func (c *EnclaveClient) IngestBlocks(blocks []*types.Block) {
	encodedBlocks := make([][]byte, 0)
	for _, block := range blocks {
		encodedBlock := obscurocommon.EncodeBlock(block)
		encodedBlocks = append(encodedBlocks, encodedBlock)
	}
	_, err := c.protoClient.IngestBlocks(context.Background(), &IngestBlocksRequest{EncodedBlocks: encodedBlocks})
	if err != nil {
		log.Log(fmt.Sprintf("failed to ingest blocks: %v", err))
	}
}

func (c *EnclaveClient) Start(block types.Block) {
	var buffer bytes.Buffer
	err := block.EncodeRLP(&buffer)
	if err != nil {
		log.Log(fmt.Sprintf("failed to encode block to send to enclave: %v", err))
	}

	_, err = c.protoClient.Start(context.Background(), &StartRequest{EncodedBlock: buffer.Bytes()})
	if err != nil {
		log.Log(fmt.Sprintf("failed to start enclave: %v", err))
	}
}

func (c *EnclaveClient) SubmitBlock(block types.Block) enclave.BlockSubmissionResponse {
	var buffer bytes.Buffer
	err := block.EncodeRLP(&buffer)
	if err != nil {
		log.Log(fmt.Sprintf("failed to encode block to send to enclave: %v", err))
	}

	response, err := c.protoClient.SubmitBlock(context.Background(), &SubmitBlockRequest{EncodedBlock: buffer.Bytes()})
	if err != nil {
		log.Log(fmt.Sprintf("failed to submit block: %v", err))
	}
	return fromBlockSubmissionResponseMsg(response.BlockSubmissionResponse)
}

func (c *EnclaveClient) SubmitRollup(rollup nodecommon.ExtRollup) {
	extRollupMsg := toExtRollupMsg(&rollup)
	_, err := c.protoClient.SubmitRollup(context.Background(), &SubmitRollupRequest{ExtRollup: &extRollupMsg})
	if err != nil {
		log.Log(fmt.Sprintf("failed to submit rollup: %v", err))
	}
}

func (c *EnclaveClient) SubmitTx(tx nodecommon.EncryptedTx) error {
	_, err := c.protoClient.SubmitTx(context.Background(), &SubmitTxRequest{EncryptedTx: tx})
	return err
}

func (c *EnclaveClient) Balance(address common.Address) uint64 {
	response, err := c.protoClient.Balance(context.Background(), &BalanceRequest{Address: address.Bytes()})
	if err != nil {
		log.Log(fmt.Sprintf("failed to retrieve account balance: %v", err))
	}
	return response.Balance
}

func (c *EnclaveClient) RoundWinner(parent obscurocommon.L2RootHash) (nodecommon.ExtRollup, bool) {
	response, err := c.protoClient.RoundWinner(context.Background(), &RoundWinnerRequest{Parent: parent.Bytes()})
	if err != nil {
		log.Log(fmt.Sprintf("failed to determine round winner: %v", err))
	}

	if response.Winner {
		return fromExtRollupMsg(response.ExtRollup), true
	}
	return nodecommon.ExtRollup{}, false
}

func (c *EnclaveClient) Stop() {
	_, err := c.protoClient.Stop(context.Background(), &StopRequest{})
	if err != nil {
		log.Log(fmt.Sprintf("failed to stop enclave: %v", err))
	}
}

func (c *EnclaveClient) GetTransaction(txHash common.Hash) (*enclave.L2Tx, bool) {
	response, err := c.protoClient.GetTransaction(context.Background(), &GetTransactionRequest{TxHash: txHash.Bytes()})
	if err != nil {
		log.Log(fmt.Sprintf("failed to retrieve transaction: %v", err))
	}

	if !response.Known {
		return nil, false
	}

	l2Tx := enclave.L2Tx{}
	err = l2Tx.DecodeRLP(rlp.NewStream(bytes.NewReader(response.EncodedTransaction), 0))
	if err != nil {
		log.Log(fmt.Sprintf("failed to decode block received from enclave: %v", err))
	}

	return &l2Tx, true
}
