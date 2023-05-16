package core

import (
	"math/big"
	"sync/atomic"
	"time"

	"github.com/obscuronet/go-obscuro/go/enclave/crypto"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/go-obscuro/go/common"
)

// Batch Data structure only for the internal use of the enclave since transactions are in clear
// Making changes to this struct will require GRPC + GRPC Converters regen
type Batch struct {
	Header *common.BatchHeader

	hash atomic.Value
	// size   atomic.Value

	Transactions []*common.L2Tx
}

// Hash returns the keccak256 hash of b's header.
// The hash is computed on the first call and cached thereafter.
func (b *Batch) Hash() *common.L2BatchHash {
	// Temporarily disabling the caching of the hash because it's causing bugs.
	// Transforming a Batch to an ExtBatch and then back to a Batch will generate a different hash if caching is enabled.
	// todo (#1547) - re-enable
	//if hash := b.hash.Load(); hash != nil {
	//	return hash.(common.L2BatchHash)
	//}
	v := b.Header.Hash()
	b.hash.Store(v)
	return &v
}

func (b *Batch) Size() (int, error) {
	bytes, err := rlp.EncodeToBytes(b)
	return len(bytes), err
}

func (b *Batch) Encode() ([]byte, error) {
	return rlp.EncodeToBytes(b)
}

func (b *Batch) NumberU64() uint64 { return b.Header.Number.Uint64() }
func (b *Batch) Number() *big.Int  { return new(big.Int).Set(b.Header.Number) }

// IsGenesis indicates whether the batch is the genesis batch.
// todo (#718) - Change this to a check against a hardcoded genesis hash.
func (b *Batch) IsGenesis() bool {
	return b.Header.Number.Cmp(big.NewInt(int64(common.L2GenesisHeight))) == 0
}

func (b *Batch) ToExtBatch(transactionBlobCrypto crypto.TransactionBlobCrypto) *common.ExtBatch {
	txHashes := make([]gethcommon.Hash, len(b.Transactions))
	for idx, tx := range b.Transactions {
		txHashes[idx] = tx.Hash()
	}

	return &common.ExtBatch{
		Header:          b.Header,
		TxHashes:        txHashes,
		EncryptedTxBlob: transactionBlobCrypto.Encrypt(b.Transactions),
	}
}

func ToBatch(extBatch *common.ExtBatch, transactionBlobCrypto crypto.TransactionBlobCrypto) *Batch {
	return &Batch{
		Header:       extBatch.Header,
		Transactions: transactionBlobCrypto.Decrypt(extBatch.EncryptedTxBlob),
	}
}

func EmptyBatch(agg gethcommon.Address, parent *common.BatchHeader, blkHash gethcommon.Hash) (*Batch, error) {
	rand, err := crypto.GeneratePublicRandomness()
	if err != nil {
		return nil, err
	}
	h := common.BatchHeader{
		Agg:        agg,
		ParentHash: parent.Hash(),
		L1Proof:    blkHash,
		Number:     big.NewInt(0).Add(parent.Number, big.NewInt(1)),
		// todo (#1548) - Consider how this time should align with the time of the L1 block used as proof.
		Time: uint64(time.Now().Unix()),
		// generate true randomness inside the enclave.
		// note that this randomness will be published in the header of the batch.
		// the randomness exposed to smart contract is combining this with the shared secret.
		MixDigest: gethcommon.BytesToHash(rand),
	}
	b := Batch{
		Header: &h,
	}
	return &b, nil
}

func DeterministicEmptyBatch(
	agg gethcommon.Address,
	parent *common.BatchHeader,
	block *types.Block,
	rand gethcommon.Hash,
	time uint64,
) *Batch {
	h := common.BatchHeader{
		Agg:        agg,
		ParentHash: parent.Hash(),
		L1Proof:    block.Hash(),
		Number:     big.NewInt(0).Add(parent.Number, big.NewInt(1)),
		// todo (#1548) - Consider how this time should align with the time of the L1 block used as proof.
		Time: time,
		// generate true randomness inside the enclave.
		// note that this randomness will be published in the header of the batch.
		// the randomness exposed to smart contract is combining this with the shared secret.
		MixDigest: rand,
	}
	b := Batch{
		Header: &h,
	}
	return &b
}
