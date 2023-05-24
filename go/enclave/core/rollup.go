package core

import "C"
import (
	"math/big"
	"sync/atomic"

	"github.com/obscuronet/go-obscuro/go/common/compression"

	"github.com/ethereum/go-ethereum/rlp"

	"github.com/obscuronet/go-obscuro/go/enclave/crypto"

	"github.com/obscuronet/go-obscuro/go/common"
)

type Rollup struct {
	Header  *common.RollupHeader
	Batches []*Batch
	hash    atomic.Value
}

// Hash returns the keccak256 hash of b's header.
// The hash is computed on the first call and cached thereafter.
func (r *Rollup) Hash() *common.L2BatchHash {
	// Temporarily disabling the caching of the hash because it's causing bugs.
	// Transforming a Rollup to an ExtRollup and then back to a Rollup will generate a different hash if caching is enabled.
	// todo (#1547) - re-enable
	//if hash := r.hash.Load(); hash != nil {
	//	return hash.(common.L2BatchHash)
	//}
	v := r.Header.Hash()
	r.hash.Store(v)
	return &v
}

func (r *Rollup) NumberU64() uint64 { return r.Header.Number.Uint64() }
func (r *Rollup) Number() *big.Int  { return new(big.Int).Set(r.Header.Number) }

// IsGenesis indicates whether the rollup is the genesis rollup.
// todo (#718) - Change this to a check against a hardcoded genesis hash.
func (r *Rollup) IsGenesis() bool {
	return r.Header.Number.Cmp(big.NewInt(int64(common.L2GenesisHeight))) == 0
}

func (r *Rollup) ToExtRollup(dataEncryptionService crypto.DataEncryptionService, compression compression.DataCompressionService) (*common.ExtRollup, error) {
	plaintextBatchesBlob, err := rlp.EncodeToBytes(r.Batches)
	if err != nil {
		return nil, err
	}

	compressedBatchesBlob, err := compression.Compress(plaintextBatchesBlob)
	if err != nil {
		return nil, err
	}

	return &common.ExtRollup{
		Header:  r.Header,
		Batches: dataEncryptionService.Encrypt(compressedBatchesBlob),
	}, nil
}

func ToRollup(encryptedRollup *common.ExtRollup, txBlobCrypto crypto.DataEncryptionService, compression compression.DataCompressionService) (*Rollup, error) {
	decryptedTxs := txBlobCrypto.Decrypt(encryptedRollup.Batches)
	encryptedBatches, err := compression.Decompress(decryptedTxs)
	if err != nil {
		return nil, err
	}

	batches := make([]*Batch, 0)
	err = rlp.DecodeBytes(encryptedBatches, &batches)
	if err != nil {
		return nil, err
	}

	return &Rollup{
		Header:  encryptedRollup.Header,
		Batches: batches,
	}, nil
}
