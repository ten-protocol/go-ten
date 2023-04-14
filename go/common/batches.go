package common

import (
	"fmt"
	"sync/atomic"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
)

// ExtBatch is an encrypted form of batch used when passing the batch around outside of an enclave.
// todo (#718) - expand this structure to contain the required fields.
type ExtBatch struct {
	Header          *BatchHeader
	TxHashes        []TxHash // The hashes of the transactions included in the batch.
	EncryptedTxBlob EncryptedTransactions
	hash            atomic.Value
}

// Hash returns the keccak256 hash of the batch's header.
// The hash is computed on the first call and cached thereafter.
func (b *ExtBatch) Hash() L2BatchHash {
	//if hash := b.hash.Load(); hash != nil {
	//	return hash.(L2BatchHash)
	//}
	v := b.Header.Hash()
	b.hash.Store(v)
	return v
}

func (b *ExtBatch) Size() (int, error) {
	bytes, err := rlp.EncodeToBytes(b)
	return len(bytes), err
}

func (b *ExtBatch) SDump() string {
	return fmt.Sprintf("Tx_Len=%d, encrypted_blob_len=%d", len(b.TxHashes), len(b.EncryptedTxBlob))
}

// BatchRequest is used when requesting a range of batches from a peer.
type BatchRequest struct {
	Requester        string
	CurrentHeadBatch *gethcommon.Hash // The requester's view of the current head batch, or nil if they haven't stored any batches.
}
