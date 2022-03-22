package rpc

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
)

// EnclaveTestRPCClient implements enclave.Enclave with a direct instance of enclave.NewEnclave
// can emulate specific rpc conditions.
type EnclaveTestRPCClient struct {
	enclave enclave.Enclave
}

func NewEnclaveTestRPCClient(nodeID common.Address, collector enclave.StatsCollector) enclave.Enclave {
	return &EnclaveTestRPCClient{enclave: enclave.NewEnclave(nodeID, true, collector)}
}

func (e *EnclaveTestRPCClient) IsReady() error {
	return e.enclave.IsReady()
}

func (e *EnclaveTestRPCClient) Attestation() obscurocommon.AttestationReport {
	return e.enclave.Attestation()
}

func (e *EnclaveTestRPCClient) GenerateSecret() obscurocommon.EncryptedSharedEnclaveSecret {
	return e.enclave.GenerateSecret()
}

func (e *EnclaveTestRPCClient) FetchSecret(report obscurocommon.AttestationReport) obscurocommon.EncryptedSharedEnclaveSecret {
	return e.enclave.FetchSecret(report)
}

func (e *EnclaveTestRPCClient) InitEnclave(secret obscurocommon.EncryptedSharedEnclaveSecret) {
	e.enclave.InitEnclave(secret)
}

func (e *EnclaveTestRPCClient) IsInitialised() bool {
	return e.enclave.IsInitialised()
}

func (e *EnclaveTestRPCClient) ProduceGenesis() enclave.BlockSubmissionResponse {
	return e.enclave.ProduceGenesis()
}

func (e *EnclaveTestRPCClient) IngestBlocks(blocks []*types.Block) {
	e.enclave.IngestBlocks(blocks)
}

func (e *EnclaveTestRPCClient) Start(block types.Block) {
	e.enclave.Start(block)
}

func (e *EnclaveTestRPCClient) SubmitBlock(block types.Block) enclave.BlockSubmissionResponse {
	return e.enclave.SubmitBlock(block)
}

func (e *EnclaveTestRPCClient) SubmitRollup(rollup nodecommon.ExtRollup) {
	e.enclave.SubmitRollup(rollup)
}

func (e *EnclaveTestRPCClient) SubmitTx(tx nodecommon.EncryptedTx) error {
	return e.enclave.SubmitTx(tx)
}

func (e *EnclaveTestRPCClient) Balance(address common.Address) uint64 {
	return e.enclave.Balance(address)
}

func (e *EnclaveTestRPCClient) RoundWinner(parent obscurocommon.L2RootHash) (nodecommon.ExtRollup, bool) {
	return e.enclave.RoundWinner(parent)
}

func (e *EnclaveTestRPCClient) Stop() {
	e.enclave.Stop()
}

func (e *EnclaveTestRPCClient) GetTransaction(txHash common.Hash) *enclave.L2Tx {
	return e.enclave.GetTransaction(txHash)
}

func (e *EnclaveTestRPCClient) StopClient() {
	e.enclave.StopClient()
}
