package common

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
)

// TenTransaction is an abstraction that transforms an Ethereum transaction into a format that can be consumed more
// easily by TEN.
type TenTransaction interface{}

type L1TxType int

const (
	RollupTx L1TxType = iota
	SecretRequestTx
	InitialiseSecretTx
	CrossChainMessageTx
	CrossChainValueTranserTx
	SequencerAddedTx
	SetImportantContractsTx
)

// ProcessedL1Data is submitted to the enclave by the guardian
type ProcessedL1Data struct {
	BlockHeader *types.Header
	Events      map[L1TxType][]*L1TxData
}

// L1TxData represents an L1 transaction that's relevant to us
type L1TxData struct {
	Type               TenTransaction
	Transaction        *types.Transaction
	Receipt            *types.Receipt
	Blobs              []*kzg4844.Blob      // Only populated for blob transactions
	CrossChainMessages *CrossChainMessages  // Only populated for xchain messages
	ValueTransfers     *ValueTransferEvents // Only populated for xchain transfers
}
