package common

import (
	"fmt"
	"math/big"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/rlp"
)

// ExtBatch is an encrypted form of batch used when passing the batch around outside of an enclave.
type ExtBatch struct {
	Header *BatchHeader
	// todo - remove and replace with enclave API
	TxHashes        []TxHash // The hashes of the transactions included in the batch.
	EncryptedTxBlob EncryptedTransactions
	hash            atomic.Value
}

// Hash returns the keccak256 hash of the batch's header.
// The hash is computed on the first call and cached thereafter.
func (b *ExtBatch) Hash() L2BatchHash {
	if hash := b.hash.Load(); hash != nil {
		return hash.(L2BatchHash)
	}
	v := b.Header.Hash()
	b.hash.Store(v)
	return v
}

func (b *ExtBatch) Encoded() ([]byte, error) {
	return rlp.EncodeToBytes(b)
}
func (b *ExtBatch) SeqNo() *big.Int { return new(big.Int).Set(b.Header.SequencerOrderNo) }
func DecodeExtBatch(encoded []byte) (*ExtBatch, error) {
	var batch ExtBatch
	if err := rlp.DecodeBytes(encoded, &batch); err != nil {
		return nil, err
	}
	return &batch, nil
}

func (b *ExtBatch) SDump() string {
	return fmt.Sprintf("Tx_Len=%d, encrypted_blob_len=%d", len(b.TxHashes), len(b.EncryptedTxBlob))
}

// BatchRequest is used when requesting a range of batches from a peer.
type BatchRequest struct {
	Requester string   // The address of the requester, used to direct the response
	FromSeqNo *big.Int // The requester's view of the current head seq no, or nil if they haven't stored any batches.
}

// AttestationRequest is used when requesting attestation reports from the sequencer.
type AttestationRequest struct {
	Requester  string      // The address of the requester, used to direct the response
	EnclaveIDs []EnclaveID // The list of enclave IDs to fetch attestations for
}

// AttestationResponse contains attestation reports from the sequencer's enclaves.
type AttestationResponse struct {
	Attestations []AttestationReport
}
