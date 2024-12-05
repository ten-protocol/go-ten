package common

import (
	"crypto/ecdsa"
	"fmt"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
	"math/big"
)

// TenTransaction is an abstraction that transforms an Ethereum transaction into a format that can be consumed more
// easily by TEN.
type TenTransaction interface{}

type L1RollupTx struct {
	Rollup EncodedRollup
}

type L1RollupHashes struct {
	BlobHashes []gethcommon.Hash
}

type L1DepositTx struct {
	Amount        *big.Int            // Amount to be deposited
	To            *gethcommon.Address // Address the ERC20 Transfer was made to (always be the Management Contract Addr)
	Sender        *gethcommon.Address // Address that issued the ERC20, the token holder or tx.origin
	TokenContract *gethcommon.Address // Address of the ERC20 Contract address that was executed
}

type L1RespondSecretTx struct {
	Secret      []byte
	RequesterID gethcommon.Address
	AttesterID  gethcommon.Address
	AttesterSig []byte
}

type L1SetImportantContractsTx struct {
	Key        string
	NewAddress gethcommon.Address
}

type L1RequestSecretTx struct {
	Attestation EncodedAttestationReport
}

type L1InitializeSecretTx struct {
	EnclaveID     *gethcommon.Address
	InitialSecret []byte
	Attestation   EncodedAttestationReport
}

// Sign signs the payload with a given private key
func (l *L1RespondSecretTx) Sign(privateKey *ecdsa.PrivateKey) *L1RespondSecretTx {
	var data []byte
	data = append(data, l.AttesterID.Bytes()...)
	data = append(data, l.RequesterID.Bytes()...)
	data = append(data, string(l.Secret)...)

	ethereumMessageHash := func(data []byte) []byte {
		prefix := fmt.Sprintf("\x19Ethereum Signed Message:\n%d", len(data))
		return crypto.Keccak256([]byte(prefix), data)
	}

	hashedData := ethereumMessageHash(data)
	// sign the hash
	signedHash, err := crypto.Sign(hashedData, privateKey)
	if err != nil {
		return nil
	}

	// set recovery id to 27; prevent malleable signatures
	signedHash[64] += 27
	l.AttesterSig = signedHash
	return l
}

// L1TxType represents different types of L1 transactions
type L1TxType uint8 // Change to uint8 for RLP serialization

const (
	RollupTx L1TxType = iota
	InitialiseSecretTx
	SecretRequestTx
	SecretResponseTx
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

//// TenTransactionWrapper wraps a TenTransaction with its concrete type
//type TenTransactionWrapper struct {
//	TypeName string // The concrete type name
//	Data     []byte // The encoded transaction data
//}

// L1TxData represents an L1 transaction that's relevant to us
type L1TxData struct {
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

//func WrapTenTransaction(tx TenTransaction) (*TenTransactionWrapper, error) {
//	if tx == nil {
//		return nil, nil
//	}
//
//	data, err := rlp.EncodeToBytes(tx)
//	if err != nil {
//		return nil, err
//	}
//
//	return &TenTransactionWrapper{
//		TypeName: fmt.Sprintf("%T", tx),
//		Data:     data,
//	}, nil
//}
//
//func (w *TenTransactionWrapper) UnwrapTransaction() (TenTransaction, error) {
//	if w == nil {
//		return nil, nil
//	}
//
//	var result TenTransaction
//	switch w.TypeName {
//	case "*L1InitializeSecretTx":
//		var tx L1InitializeSecretTx
//		if err := rlp.DecodeBytes(w.Data, &tx); err != nil {
//			return nil, err
//		}
//		result = &tx
//	case "*L1RequestSecretTx":
//		var tx L1RequestSecretTx
//		if err := rlp.DecodeBytes(w.Data, &tx); err != nil {
//			return nil, err
//		}
//		result = &tx
//
//	case "*L1SetImportantContractsTx":
//		var tx L1SetImportantContractsTx
//		if err := rlp.DecodeBytes(w.Data, &tx); err != nil {
//			return nil, err
//		}
//		result = &tx
//	case "*L1RespondSecretTx":
//		var tx L1RespondSecretTx
//		if err := rlp.DecodeBytes(w.Data, &tx); err != nil {
//			return nil, err
//		}
//		result = &tx
//	case "*L1DepositTx":
//		var tx L1DepositTx
//		if err := rlp.DecodeBytes(w.Data, &tx); err != nil {
//			return nil, err
//		}
//		result = &tx
//	case "*L1RollupHashes":
//		var tx L1RollupHashes
//		if err := rlp.DecodeBytes(w.Data, &tx); err != nil {
//			return nil, err
//		}
//		result = &tx
//	case "*L1RollupTx":
//		var tx L1RollupTx
//		if err := rlp.DecodeBytes(w.Data, &tx); err != nil {
//			return nil, err
//		}
//		result = &tx
//	default:
//		return nil, fmt.Errorf("unknown transaction type: %s", w.TypeName)
//	}
//
//	return result, nil
//}
