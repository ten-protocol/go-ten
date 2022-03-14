package rpc

import (
	"bytes"
	"context"
	"flag"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"

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

// TODO - Joel - Handle the errors coming back from the client requests.

// Attestation - Produces an attestation report which will be used to request the shared secret from another enclave.
func (c *EnclaveClient) Attestation() obscurocommon.AttestationReport {
	response, _ := c.clientInternal.Attestation(context.Background(), &AttestationRequest{})
	return toAttestationReport(response.AttestationReportMsg)
}

// GenerateSecret - the genesis enclave is responsible with generating the secret entropy
func (c *EnclaveClient) GenerateSecret() obscurocommon.EncryptedSharedEnclaveSecret {
	response, _ := c.clientInternal.GenerateSecret(context.Background(), &GenerateSecretRequest{})
	return response.EncryptedSharedEnclaveSecret
}

// FetchSecret - return the shared secret encrypted with the key from the attestation
func (c *EnclaveClient) FetchSecret(report obscurocommon.AttestationReport) obscurocommon.EncryptedSharedEnclaveSecret {
	attestationReportMsg := toAttestationReportMsg(report)
	request := FetchSecretRequest{AttestationReportMsg: &attestationReportMsg}
	response, _ := c.clientInternal.FetchSecret(context.Background(), &request)
	return response.EncryptedSharedEnclaveSecret
}

// Init - initialise an enclave with a seed received by another enclave
func (c *EnclaveClient) Init(secret obscurocommon.EncryptedSharedEnclaveSecret) {
	c.clientInternal.Init(context.Background(), &InitRequest{EncryptedSharedEnclaveSecret: secret})
}

// IsInitialised - true if the shared secret is available
func (c *EnclaveClient) IsInitialised() bool {
	response, _ := c.clientInternal.IsInitialised(context.Background(), &IsInitialisedRequest{})
	return response.IsInitialised
}

// ProduceGenesis - the genesis enclave produces the genesis rollup
func (c *EnclaveClient) ProduceGenesis() enclave.BlockSubmissionResponse {
	response, _ := c.clientInternal.ProduceGenesis(context.Background(), &ProduceGenesisRequest{})
	return toBlockSubmissionResponse(response.BlockSubmissionResponse)
}

// IngestBlocks - feed L1 blocks into the enclave to catch up
func (c *EnclaveClient) IngestBlocks(blocks []*types.Block) {
	var blockRlps [][]byte
	for _, block := range blocks {
		var buffer bytes.Buffer
		block.EncodeRLP(&buffer)
		blockRlps = append(blockRlps, buffer.Bytes())
	}
	c.clientInternal.IngestBlocks(context.Background(), &IngestBlocksRequest{BlockRlps: blockRlps})
}

// Start - start speculative execution
func (c *EnclaveClient) Start(block types.Block) {
	var buffer bytes.Buffer
	block.EncodeRLP(&buffer)
	c.clientInternal.Start(context.Background(), &StartRequest{BlockRlp: buffer.Bytes()})
}

// SubmitBlock - When a new POBI round starts, the host submits a block to the enclave, which responds with a rollup
// it is the responsibility of the host to gossip the returned rollup
// For good functioning the caller should always submit blocks ordered by height
// submitting a block before receiving a parent of it, will result in it being ignored
func (c *EnclaveClient) SubmitBlock(block types.Block) enclave.BlockSubmissionResponse {
	var buffer bytes.Buffer
	block.EncodeRLP(&buffer)
	response, _ := c.clientInternal.SubmitBlock(context.Background(), &SubmitBlockRequest{BlockRlp: buffer.Bytes()})
	return toBlockSubmissionResponse(response.BlockSubmissionResponse)
}

// SubmitRollup - receive gossiped rollups
func (c *EnclaveClient) SubmitRollup(rollup nodecommon.ExtRollup) {
	extRollupMsg := toExtRollupMsg(&rollup)
	c.clientInternal.SubmitRollup(context.Background(), &SubmitRollupRequest{ExtRollup: &extRollupMsg})
}

// SubmitTx - user transactions
func (c *EnclaveClient) SubmitTx(tx nodecommon.EncryptedTx) error {
	_, err := c.clientInternal.SubmitTx(context.Background(), &SubmitTxRequest{EncryptedTx: tx})
	return err
}

// Balance - returns the balance of an address with a block delay
func (c *EnclaveClient) Balance(address common.Address) uint64 {
	response, _ := c.clientInternal.Balance(context.Background(), &BalanceRequest{Address: address.Bytes()})
	return response.Balance
}

// RoundWinner - calculates and returns the winner for a round
func (c *EnclaveClient) RoundWinner(parent obscurocommon.L2RootHash) (nodecommon.ExtRollup, bool) {
	response, _ := c.clientInternal.RoundWinner(context.Background(), &RoundWinnerRequest{Parent: parent.Bytes()})

	if response.Winner {
		return toExtRollup(response.ExtRollup), true
	} else {
		return nodecommon.ExtRollup{}, false
	}
}

// Stop gracefully stops the enclave
func (c *EnclaveClient) Stop() {
	c.clientInternal.Stop(context.Background(), &StopRequest{})
}

// GetTransaction returns a transaction given its Signed Hash, returns nil, false when Transaction is unknown
func (c *EnclaveClient) GetTransaction(txHash common.Hash) (*enclave.L2Tx, bool) {
	response, _ := c.clientInternal.GetTransaction(context.Background(), &GetTransactionRequest{TxHash: txHash.Bytes()})

	if response.Unknown {
		return nil, false
	} else {
		l2Tx := enclave.L2Tx{}
		l2Tx.DecodeRLP(rlp.NewStream(bytes.NewReader(response.BlockRlp), 0))
		return &l2Tx, true
	}
}

// Converters between RPC and regular types.

func toAttestationReport(msg *AttestationReportMsg) obscurocommon.AttestationReport {
	return obscurocommon.AttestationReport{Owner: common.BytesToAddress(msg.Owner)}
}

func toAttestationReportMsg(report obscurocommon.AttestationReport) AttestationReportMsg {
	return AttestationReportMsg{Owner: report.Owner.Bytes()}
}

func toBlockSubmissionResponse(msg *BlockSubmissionResponseMsg) enclave.BlockSubmissionResponse {
	var withdrawals []nodecommon.Withdrawal
	for _, withdrawalMsg := range msg.Withdrawals {
		withdrawal := nodecommon.Withdrawal{Amount: withdrawalMsg.Amount, Address: common.BytesToAddress(withdrawalMsg.Address)}
		withdrawals = append(withdrawals, withdrawal)
	}

	return enclave.BlockSubmissionResponse{
		L1Hash:            common.BytesToHash(msg.L1Hash),
		L1Height:          msg.L1Height,
		L1Parent:          common.BytesToHash(msg.L1Parent),
		L2Hash:            common.BytesToHash(msg.L2Hash),
		L2Height:          msg.L2Height,
		L2Parent:          common.BytesToHash(msg.L2Parent),
		Withdrawals:       withdrawals,
		ProducedRollup:    toExtRollup(msg.ProducedRollup),
		IngestedBlock:     msg.IngestedBlock,
		IngestedNewRollup: msg.IngestedNewRollup,
	}
}

func toExtRollup(msg *ExtRollupMsg) nodecommon.ExtRollup {
	var withdrawals []nodecommon.Withdrawal
	for _, withdrawalMsg := range msg.Header.Withdrawals {
		withdrawal := nodecommon.Withdrawal{Amount: withdrawalMsg.Amount, Address: common.BytesToAddress(withdrawalMsg.Address)}
		withdrawals = append(withdrawals, withdrawal)
	}
	header := nodecommon.Header{
		ParentHash:  common.BytesToHash(msg.Header.ParentHash),
		Agg:         common.BytesToAddress(msg.Header.Agg),
		Nonce:       msg.Header.Nonce,
		L1Proof:     common.BytesToHash(msg.Header.L1Proof),
		State:       msg.Header.StateRoot,
		Withdrawals: withdrawals,
	}

	var txs []nodecommon.EncryptedTx
	for _, tx := range msg.Txs {
		txs = append(txs, tx)
	}

	return nodecommon.ExtRollup{
		Header: &header,
		Txs:    txs,
	}
}

func toExtRollupMsg(rollup *nodecommon.ExtRollup) ExtRollupMsg {
	var withdrawalMsgs []*WithdrawalMsg
	for _, withdrawal := range rollup.Header.Withdrawals {
		withdrawalMsg := WithdrawalMsg{Amount: withdrawal.Amount, Address: withdrawal.Address.Bytes()}
		withdrawalMsgs = append(withdrawalMsgs, &withdrawalMsg)
	}

	headerMsg := HeaderMsg{
		ParentHash:  rollup.Header.ParentHash.Bytes(),
		Agg:         rollup.Header.Agg.Bytes(),
		Nonce:       rollup.Header.Nonce,
		L1Proof:     rollup.Header.L1Proof.Bytes(),
		StateRoot:   rollup.Header.State,
		Withdrawals: withdrawalMsgs,
	}

	var txs [][]byte
	for _, tx := range rollup.Txs {
		txs = append(txs, tx)
	}

	return ExtRollupMsg{Header: &headerMsg, Txs: txs}
}
