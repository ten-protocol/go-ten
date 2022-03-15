package rpc

import (
	"bytes"
	"context"
	"flag"
	"fmt"
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

type EnclaveClient struct {
	clientInternal EnclaveInternalClient
}

func NewEnclaveClient(port uint64) EnclaveClient {
	connection := enclaveClientConn(port)
	client := EnclaveClient{NewEnclaveInternalClient(connection)}
	return client
}

// TODO - Joel - Need to use same port as server, and not randomly allocate it in the server.
func enclaveClientConn(port uint64) *grpc.ClientConn {
	flag.Parse()
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", port), opts...)
	// TODO - Joel - Better error handling.
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	return conn
}

// TODO - Joel - Handle the errors coming back from the client requests.

func (c *EnclaveClient) Attestation() obscurocommon.AttestationReport {
	response, e := c.clientInternal.Attestation(context.Background(), &AttestationRequest{})
	if e != nil {
		panic(e)
	}
	return toAttestationReport(response.AttestationReportMsg)
}

func (c *EnclaveClient) GenerateSecret() obscurocommon.EncryptedSharedEnclaveSecret {
	response, e := c.clientInternal.GenerateSecret(context.Background(), &GenerateSecretRequest{})
	if e != nil {
		panic(e)
	}
	return response.EncryptedSharedEnclaveSecret
}

func (c *EnclaveClient) FetchSecret(report obscurocommon.AttestationReport) obscurocommon.EncryptedSharedEnclaveSecret {
	attestationReportMsg := toAttestationReportMsg(report)
	request := FetchSecretRequest{AttestationReportMsg: &attestationReportMsg}
	response, e := c.clientInternal.FetchSecret(context.Background(), &request)
	if e != nil {
		panic(e)
	}
	return response.EncryptedSharedEnclaveSecret
}

func (c *EnclaveClient) Init(secret obscurocommon.EncryptedSharedEnclaveSecret) {
	_, e := c.clientInternal.Init(context.Background(), &InitRequest{EncryptedSharedEnclaveSecret: secret})
	if e != nil {
		panic(e)
	}
}

func (c *EnclaveClient) IsInitialised() bool {
	response, e := c.clientInternal.IsInitialised(context.Background(), &IsInitialisedRequest{})
	if e != nil {
		panic(e)
	}
	return response.IsInitialised
}

func (c *EnclaveClient) ProduceGenesis() enclave.BlockSubmissionResponse {
	response, e := c.clientInternal.ProduceGenesis(context.Background(), &ProduceGenesisRequest{})
	if e != nil {
		panic(e)
	}
	return toBlockSubmissionResponse(response.BlockSubmissionResponse)
}

func (c *EnclaveClient) IngestBlocks(blocks []*types.Block) {
	encodedBlocks := make([][]byte, 0)
	for _, block := range blocks {
		encodedBlock := obscurocommon.EncodeBlock(block)
		encodedBlocks = append(encodedBlocks, encodedBlock)
	}
	_, e := c.clientInternal.IngestBlocks(context.Background(), &IngestBlocksRequest{EncodedBlocks: encodedBlocks})
	if e != nil {
		panic(e)
	}
}

func (c *EnclaveClient) Start(block types.Block) {
	var buffer bytes.Buffer
	block.EncodeRLP(&buffer)
	_, e := c.clientInternal.Start(context.Background(), &StartRequest{EncodedBlock: buffer.Bytes()})
	if e != nil {
		panic(e)
	}
}

func (c *EnclaveClient) SubmitBlock(block types.Block) enclave.BlockSubmissionResponse {
	var buffer bytes.Buffer
	block.EncodeRLP(&buffer)
	response, e := c.clientInternal.SubmitBlock(context.Background(), &SubmitBlockRequest{EncodedBlock: buffer.Bytes()})
	if e != nil {
		panic(e)
	}
	return toBlockSubmissionResponse(response.BlockSubmissionResponse)
}

func (c *EnclaveClient) SubmitRollup(rollup nodecommon.ExtRollup) {
	extRollupMsg := toExtRollupMsg(&rollup)
	_, e := c.clientInternal.SubmitRollup(context.Background(), &SubmitRollupRequest{ExtRollup: &extRollupMsg})
	if e != nil {
		panic(e)
	}
}

func (c *EnclaveClient) SubmitTx(tx nodecommon.EncryptedTx) error {
	_, e := c.clientInternal.SubmitTx(context.Background(), &SubmitTxRequest{EncryptedTx: tx})
	if e != nil {
		panic(e)
	}
	return e
}

func (c *EnclaveClient) Balance(address common.Address) uint64 {
	response, e := c.clientInternal.Balance(context.Background(), &BalanceRequest{Address: address.Bytes()})
	if e != nil {
		panic(e)
	}
	return response.Balance
}

func (c *EnclaveClient) RoundWinner(parent obscurocommon.L2RootHash) (nodecommon.ExtRollup, bool) {
	response, e := c.clientInternal.RoundWinner(context.Background(), &RoundWinnerRequest{Parent: parent.Bytes()})
	if e != nil {
		panic(e)
	}

	if response.Winner {
		return toExtRollup(response.ExtRollup), true
	}
	return nodecommon.ExtRollup{}, false
}

func (c *EnclaveClient) Stop() {
	_, e := c.clientInternal.Stop(context.Background(), &StopRequest{})
	if e != nil {
		panic(e)
	}
}

func (c *EnclaveClient) GetTransaction(txHash common.Hash) (*enclave.L2Tx, bool) {
	response, e := c.clientInternal.GetTransaction(context.Background(), &GetTransactionRequest{TxHash: txHash.Bytes()})
	if e != nil {
		panic(e)
	}

	if response.Unknown {
		return nil, false
	}
	l2Tx := enclave.L2Tx{}
	l2Tx.DecodeRLP(rlp.NewStream(bytes.NewReader(response.EncodedTransaction), 0))
	return &l2Tx, true
}

// Converters between RPC and regular types.

func toAttestationReport(msg *AttestationReportMsg) obscurocommon.AttestationReport {
	return obscurocommon.AttestationReport{Owner: common.BytesToAddress(msg.Owner)}
}

func toAttestationReportMsg(report obscurocommon.AttestationReport) AttestationReportMsg {
	return AttestationReportMsg{Owner: report.Owner.Bytes()}
}

func toBlockSubmissionResponse(msg *BlockSubmissionResponseMsg) enclave.BlockSubmissionResponse {
	withdrawals := make([]nodecommon.Withdrawal, 0)
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

func toBlockSubmissionResponseMsg(response enclave.BlockSubmissionResponse) BlockSubmissionResponseMsg {
	withdrawalMsgs := make([]*WithdrawalMsg, 0)
	for _, withdrawal := range response.Withdrawals {
		withdrawalMsg := WithdrawalMsg{Amount: withdrawal.Amount, Address: withdrawal.Address.Bytes()}
		withdrawalMsgs = append(withdrawalMsgs, &withdrawalMsg)
	}

	producedRollupMsg := toExtRollupMsg(&response.ProducedRollup)

	return BlockSubmissionResponseMsg{
		L1Hash:            response.L1Hash.Bytes(),
		L1Height:          response.L1Height,
		L1Parent:          response.L1Parent.Bytes(),
		L2Hash:            response.L2Hash.Bytes(),
		L2Height:          response.L2Height,
		L2Parent:          response.L2Parent.Bytes(),
		Withdrawals:       withdrawalMsgs,
		ProducedRollup:    &producedRollupMsg,
		IngestedBlock:     response.IngestedBlock,
		IngestedNewRollup: response.IngestedNewRollup,
	}
}

func toExtRollup(msg *ExtRollupMsg) nodecommon.ExtRollup {
	withdrawals := make([]nodecommon.Withdrawal, 0)
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

	txs := make([]nodecommon.EncryptedTx, 0)
	for _, tx := range msg.Txs {
		txs = append(txs, tx)
	}

	return nodecommon.ExtRollup{
		Header: &header,
		Txs:    txs,
	}
}

func toExtRollupMsg(rollup *nodecommon.ExtRollup) ExtRollupMsg {
	// TODO - Joel - See if I can avoid the case of a nil header. It's caused by `enclave.go/SubmitBlock` returning early.
	var headerMsg HeaderMsg
	if rollup.Header != nil {
		withdrawalMsgs := make([]*WithdrawalMsg, 0)
		for _, withdrawal := range rollup.Header.Withdrawals {
			withdrawalMsg := WithdrawalMsg{Amount: withdrawal.Amount, Address: withdrawal.Address.Bytes()}
			withdrawalMsgs = append(withdrawalMsgs, &withdrawalMsg)
		}

		headerMsg = HeaderMsg{
			ParentHash:  rollup.Header.ParentHash.Bytes(),
			Agg:         rollup.Header.Agg.Bytes(),
			Nonce:       rollup.Header.Nonce,
			L1Proof:     rollup.Header.L1Proof.Bytes(),
			StateRoot:   rollup.Header.State,
			Withdrawals: withdrawalMsgs,
		}
	}

	txs := make([][]byte, 0)
	for _, tx := range rollup.Txs {
		txs = append(txs, tx)
	}

	return ExtRollupMsg{Header: &headerMsg, Txs: txs}
}
