package core

import (
	"math/big"
	"sync/atomic"

	"github.com/obscuronet/go-obscuro/go/common/compression"

	"github.com/obscuronet/go-obscuro/go/enclave/crypto"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/go-obscuro/go/common"
)

// Batch Data structure only for the internal use of the enclave since transactions are in clear
// Making changes to this struct will require GRPC + GRPC Converters regen
type Batch struct {
	Header       *common.BatchHeader
	hash         atomic.Value
	Transactions []*common.L2Tx
}

// Hash returns the keccak256 hash of b's header.
// The hash is computed on the first call and cached thereafter.
func (b *Batch) Hash() common.L2BatchHash {
	if hash := b.hash.Load(); hash != nil {
		return hash.(common.L2BatchHash)
	}
	v := b.Header.Hash()
	b.hash.Store(v)
	return v
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

func (b *Batch) ToExtBatch(transactionBlobCrypto crypto.DataEncryptionService, compression compression.DataCompressionService) (*common.ExtBatch, error) {
	txHashes := make([]gethcommon.Hash, len(b.Transactions))
	for idx, tx := range b.Transactions {
		txHashes[idx] = tx.Hash()
	}

	bytes, err := rlp.EncodeToBytes(b.Transactions)
	if err != nil {
		return nil, err
	}
	compressed, err := compression.CompressBatch(bytes)
	if err != nil {
		return nil, err
	}
	return &common.ExtBatch{
		Header:          b.Header,
		TxHashes:        txHashes,
		EncryptedTxBlob: transactionBlobCrypto.Encrypt(compressed),
	}, nil
}

func ToBatch(extBatch *common.ExtBatch, transactionBlobCrypto crypto.DataEncryptionService, compression compression.DataCompressionService) (*Batch, error) {
	compressed := transactionBlobCrypto.Decrypt(extBatch.EncryptedTxBlob)
	encoded, err := compression.Decompress(compressed)
	if err != nil {
		return nil, err
	}
	var txs []*common.L2Tx
	err = rlp.DecodeBytes(encoded, &txs)
	if err != nil {
		return nil, err
	}
	return &Batch{
		Header:       extBatch.Header,
		Transactions: txs,
	}, nil
}

func DeterministicEmptyBatch(
	agg gethcommon.Address,
	parent *common.BatchHeader,
	block *types.Block,
	rand gethcommon.Hash,
	time uint64,
) *Batch {
	h := common.BatchHeader{
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
