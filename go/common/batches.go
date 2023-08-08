package common

import (
	"fmt"
	"math/big"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/rlp"
)

// ExtBatch is an encrypted form of batch used when passing the batch around outside of an enclave.
// todo (#718) - expand this structure to contain the required fields.
type ExtBatch struct {
	Header *BatchHeader
	// todo - remove
	TxHashes        []TxHash // The hashes of the transactions included in the batch.
	EncryptedTxBlob EncryptedTransactions
	hash            atomic.Value
}

// Hash returns the keccak256 hash of the batch's header.
// The hash is computed on the first call and cached thereafter.
func (b *ExtBatch) Hash() L2BatchHash {
	if hash := b.hash.Load(); hash != nil {
		// todo (tudor) - remove this
		v := b.Header.Hash()
		cv := hash.(L2BatchHash)
		if v != cv {
			panic("cached ExtBatch hash is wrong!")
		}
		return v
	}
	v := b.Header.Hash()
	b.hash.Store(v)
	return v
}

func (b *ExtBatch) Size() (int, error) {
	bytes, err := rlp.EncodeToBytes(b)
	return len(bytes), err
}

func (b *ExtBatch) Encoded() ([]byte, error) {
	return rlp.EncodeToBytes(b)
}

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
