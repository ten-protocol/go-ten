package core

import (
	"bytes"
	"compress/gzip"
	"io"
	"math/big"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/rlp"

	"github.com/obscuronet/go-obscuro/go/enclave/crypto"

	"github.com/google/brotli/go/cbrotli"
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

func (r *Rollup) ToExtRollup(txBlobCrypto crypto.TransactionBlobCrypto) (*common.ExtRollup, error) {
	plaintextBatchesBlob, err := rlp.EncodeToBytes(r.Batches)
	if err != nil {
		return nil, err
	}

	compressedBatchesBlob, err := brotliCompress(plaintextBatchesBlob)
	if err != nil {
		return nil, err
	}

	return &common.ExtRollup{
		Header:  r.Header,
		Batches: txBlobCrypto.Encrypt(compressedBatchesBlob),
	}, nil
}

func ToRollup(encryptedRollup *common.ExtRollup, txBlobCrypto crypto.TransactionBlobCrypto) (*Rollup, error) {
	decryptedTxs := txBlobCrypto.Decrypt(encryptedRollup.Batches)
	encryptedBatches, err := brotliDecompress(decryptedTxs)
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

// todo - move these to a service
// gzipCompress the byte array using gzip
func gzipCompress(in []byte) ([]byte, error) {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write(in); err != nil {
		return nil, err
	}
	if err := gz.Close(); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// Decompress the byte array using gzip
func gzipDecompress(in []byte) ([]byte, error) {
	reader := bytes.NewReader(in)
	gz, err := gzip.NewReader(reader)
	if err != nil {
		return nil, err
	}
	defer gz.Close()

	return io.ReadAll(gz)
}

func brotliCompress(in []byte) ([]byte, error) {
	return cbrotli.Encode(in, cbrotli.WriterOptions{Quality: 11})
}

func brotliDecompress(in []byte) ([]byte, error) {
	return cbrotli.Decode(in)
}
