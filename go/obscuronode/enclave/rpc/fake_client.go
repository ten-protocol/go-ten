package rpc

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
)

// EnclaveFakeRPCClient implements enclave.enclave.Enclave and should be used by the host when communicating with the enclave via RPC.
type EnclaveFakeRPCClient struct {
	enclave enclave.Enclave
}

func NewEnclaveFakeRPCClient(nodeID common.Address, collector enclave.StatsCollector) enclave.Enclave {
	return &EnclaveFakeRPCClient{enclave: enclave.NewEnclave(nodeID, true, collector)}
}

func (e *EnclaveFakeRPCClient) IsReady() error {
	return e.enclave.IsReady()
}

func (e *EnclaveFakeRPCClient) Attestation() obscurocommon.AttestationReport {
	return e.enclave.Attestation()
}

func (e *EnclaveFakeRPCClient) GenerateSecret() obscurocommon.EncryptedSharedEnclaveSecret {
	return e.enclave.GenerateSecret()
}

func (e *EnclaveFakeRPCClient) FetchSecret(report obscurocommon.AttestationReport) obscurocommon.EncryptedSharedEnclaveSecret {
	return e.enclave.FetchSecret(report)
}

func (e *EnclaveFakeRPCClient) InitEnclave(secret obscurocommon.EncryptedSharedEnclaveSecret) {
	e.enclave.InitEnclave(secret)
}

func (e *EnclaveFakeRPCClient) IsInitialised() bool {
	return e.enclave.IsInitialised()
}

func (e *EnclaveFakeRPCClient) ProduceGenesis() enclave.BlockSubmissionResponse {
	return e.enclave.ProduceGenesis()
}

func (e *EnclaveFakeRPCClient) IngestBlocks(blocks []*types.Block) {
	e.enclave.IngestBlocks(blocks)
}

func (e *EnclaveFakeRPCClient) Start(block types.Block) {
	e.enclave.Start(block)
}

func (e *EnclaveFakeRPCClient) SubmitBlock(block types.Block) enclave.BlockSubmissionResponse {
	return e.enclave.SubmitBlock(block)
}

func (e *EnclaveFakeRPCClient) SubmitRollup(rollup nodecommon.ExtRollup) {
	e.enclave.SubmitRollup(rollup)
}

func (e *EnclaveFakeRPCClient) SubmitTx(tx nodecommon.EncryptedTx) error {
	return e.enclave.SubmitTx(tx)
}

func (e *EnclaveFakeRPCClient) Balance(address common.Address) uint64 {
	return e.enclave.Balance(address)
}

func (e *EnclaveFakeRPCClient) RoundWinner(parent obscurocommon.L2RootHash) (nodecommon.ExtRollup, bool) {
	return e.enclave.RoundWinner(parent)
}

func (e *EnclaveFakeRPCClient) Stop() {
	e.enclave.Stop()
}

func (e *EnclaveFakeRPCClient) GetTransaction(txHash common.Hash) *enclave.L2Tx {
	return e.enclave.GetTransaction(txHash)
}

func (e *EnclaveFakeRPCClient) StopClient() {
	e.enclave.StopClient()
}
