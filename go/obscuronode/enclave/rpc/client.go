package rpc

import (
	"context"
	"flag"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var serverAddr = flag.String("addr", "localhost:50051", "The server address in the format of host:port")

type EnclaveClient struct {
	clientInternal EnclaveInternalClient
}

func NewEnclaveClient() EnclaveClient {
	connection := enclaveClientConn()
	client := EnclaveClient{NewEnclaveInternalClient(connection)}
	return client
}

func enclaveClientConn() *grpc.ClientConn {
	flag.Parse()
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.Dial(*serverAddr, opts...)
	// TODO - Joel - Better error handling.
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	return conn
}

// Attestation - Produces an attestation report which will be used to request the shared secret from another enclave.
func (c *EnclaveClient) Attestation() obscurocommon.AttestationReport {
	// TODO - Joel - Handle error.
	response, _ := c.clientInternal.Attestation(context.Background(), &AttestationRequest{})
	return obscurocommon.AttestationReport{Owner: common.BytesToAddress(response.Owner)}
}

// GenerateSecret - the genesis enclave is responsible with generating the secret entropy
func (c *EnclaveClient) GenerateSecret() obscurocommon.EncryptedSharedEnclaveSecret {
	// TODO - Joel - Handle error.
	response, _ := c.clientInternal.GenerateSecret(context.Background(), &GenerateSecretRequest{})
	return response.EncryptedSharedEnclaveSecret
}

// FetchSecret - return the shared secret encrypted with the key from the attestation
func (c *EnclaveClient) FetchSecret(report obscurocommon.AttestationReport) obscurocommon.EncryptedSharedEnclaveSecret {
	// TODO - Joel - Handle error.
	// TODO - Joel - Pass `report` into request.
	response, _ := c.clientInternal.FetchSecret(context.Background(), &FetchSecretRequest{})
	return response.EncryptedSharedEnclaveSecret
}

// Init - initialise an enclave with a seed received by another enclave
func (c *EnclaveClient) Init(secret obscurocommon.EncryptedSharedEnclaveSecret) {
	// TODO - Joel - Handle error.
	// TODO - Joel - Pass `secret` into request.
	c.clientInternal.Init(context.Background(), &InitRequest{})
}

// IsInitialised - true if the shared secret is available
func (c *EnclaveClient) IsInitialised() bool {
	// TODO - Joel - Handle error.
	response, _ := c.clientInternal.IsInitialised(context.Background(), &IsInitialisedRequest{})
	return response.IsInitialised
}

// ProduceGenesis - the genesis enclave produces the genesis rollup
func (c *EnclaveClient) ProduceGenesis() enclave.BlockSubmissionResponse {
	// TODO - Joel - Handle error.
	response, _ := c.clientInternal.ProduceGenesis(context.Background(), &ProduceGenesisRequest{})
	// TODO - Joel - Fill this response object in.
	return enclave.BlockSubmissionResponse{}
}

// IngestBlocks - feed L1 blocks into the enclave to catch up
func (c *EnclaveClient) IngestBlocks(blocks []*types.Block) {
	// TODO - Joel - Handle error.
	// TODO - Joel - Pass `blocks` into request.
	c.clientInternal.IngestBlocks(context.Background(), &IngestBlocksRequest{})
}

// Start - start speculative execution
func (c *EnclaveClient) Start(block types.Block) {
	// TODO - Joel - Handle error.
	// TODO - Joel - Pass `block` into request.
	c.clientInternal.Start(context.Background(), &StartRequest{})
}

// SubmitBlock - When a new POBI round starts, the host submits a block to the enclave, which responds with a rollup
// it is the responsibility of the host to gossip the returned rollup
// For good functioning the caller should always submit blocks ordered by height
// submitting a block before receiving a parent of it, will result in it being ignored
func (c *EnclaveClient) SubmitBlock(block types.Block) enclave.BlockSubmissionResponse {
	// TODO - Joel - Handle error.
	// TODO - Joel - Pass `block` into request.
	response, _ := c.clientInternal.SubmitBlock(context.Background(), &SubmitBlockRequest{})
	// TODO - Joel - Fill this response object in.
	return enclave.BlockSubmissionResponse{}
}

// SubmitRollup - receive gossiped rollups
func (c *EnclaveClient) SubmitRollup(rollup nodecommon.ExtRollup) {
	// TODO - Joel - Handle error.
	// TODO - Joel - Pass `rollup` into request.
	c.clientInternal.SubmitRollup(context.Background(), &SubmitRollupRequest{})
}

// SubmitTx - user transactions
func (c *EnclaveClient) SubmitTx(tx nodecommon.EncryptedTx) error {
	// TODO - Joel - Pass `tx` into request.
	_, err := c.clientInternal.SubmitTx(context.Background(), &SubmitTxRequest{})
	return err
}

// Balance - returns the balance of an address with a block delay
func (c *EnclaveClient) Balance(address common.Address) uint64 {
	// TODO - Joel - Handle error.
	// TODO - Joel - Pass `address` into request.
	response, _ := c.clientInternal.Balance(context.Background(), &BalanceRequest{})
	return response.Balance
}

// RoundWinner - calculates and returns the winner for a round
func (c *EnclaveClient) RoundWinner(parent obscurocommon.L2RootHash) (nodecommon.ExtRollup, bool) {
	// TODO - Joel - Handle error.
	// TODO - Joel - Pass `parent` into request.
	response, _ := c.clientInternal.RoundWinner(context.Background(), &RoundWinnerRequest{})
	// TODO - Joel - Work out when to return `false`.
	// TODO - Joel - Fill this response object in.
	return nodecommon.ExtRollup{}, true
}

// Stop gracefully stops the enclave
func (c *EnclaveClient) Stop() {
	// TODO - Joel - Handle error.
	c.clientInternal.Stop(context.Background(), &StopRequest{})
}

// GetTransaction returns a transaction given its Signed Hash, returns nil, false when Transaction is unknown
func (c *EnclaveClient) GetTransaction(txHash common.Hash) (*enclave.L2Tx, bool) {
	// TODO - Joel - Handle error.
	// TODO - Joel - Pass `txHash` into request.
	response, _ := c.clientInternal.RoundWinner(context.Background(), &RoundWinnerRequest{})
	// TODO - Joel - Work out when to return `false`.
	// TODO - Joel - Fill this response object in.
	return &enclave.L2Tx{}, true
}
