package common

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
)

// TenTransaction is an abstraction that transforms an Ethereum transaction into a format that can be consumed more
// easily by TEN.
type TenTransaction interface{}

// L1TxType represents different types of L1 transactions
type L1TxType uint8 // Change to uint8 for RLP serialization

const (
	RollupTx L1TxType = iota
	SecretRequestTx
	InitialiseSecretTx
	CrossChainMessageTx
	CrossChainValueTranserTx
	SequencerAddedTx
	SetImportantContractsTx
)

// L1Event represents a single event type and its associated transactions
type L1Event struct {
	Type uint8
	Txs  []*L1TxData
}

// ProcessedL1Data is submitted to the enclave by the guardian
type ProcessedL1Data struct {
	BlockHeader *types.Header
	Events      []L1Event // Changed from map to slice of L1Event
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

func (p *ProcessedL1Data) AddEvent(txType L1TxType, tx *L1TxData) {
	for i := range p.Events {
		if p.Events[i].Type == uint8(txType) {
			p.Events[i].Txs = append(p.Events[i].Txs, tx)
			return
		}
	}
	p.Events = append(p.Events, L1Event{
		Type: uint8(txType), // Convert to uint8 when storing
		Txs:  []*L1TxData{tx},
	})
}

func (p *ProcessedL1Data) GetEvents(txType L1TxType) []*L1TxData {
	if p == nil || len(p.Events) == 0 {
		return nil
	}

	for _, event := range p.Events {
		if event.Type == uint8(txType) {
			if event.Txs == nil {
				return nil
			}
			return event.Txs
		}
	}
	return nil
}
