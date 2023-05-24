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
	headers := make([]*common.BatchHeader, len(r.Batches))
	transactions := make([][]*common.L2Tx, len(r.Batches))
	for i, batch := range r.Batches {
		headers[i] = batch.Header
		transactions[i] = batch.Transactions
	}

	plaintextTransactionsBlob, err := rlp.EncodeToBytes(transactions)
	if err != nil {
		return nil, err
	}

	headersBlob, err := rlp.EncodeToBytes(headers)
	if err != nil {
		return nil, err
	}

	compressedTransactionsBlob, err := compression.Compress(plaintextTransactionsBlob)
	if err != nil {
		return nil, err
	}

	compressedHeadersBlob, err := compression.Compress(headersBlob)
	if err != nil {
		return nil, err
	}

	return &common.ExtRollup{
		Header:        r.Header,
		BatchPayloads: dataEncryptionService.Encrypt(compressedTransactionsBlob),
		BatchHeaders:  compressedHeadersBlob,
	}, nil
}

func ToRollup(encryptedRollup *common.ExtRollup, dataEncryptionService crypto.DataEncryptionService, dataCompressionService compression.DataCompressionService) (*Rollup, error) {
	headers := make([]common.BatchHeader, 0)
	headersBlob, err := dataCompressionService.Decompress(encryptedRollup.BatchHeaders)
	if err != nil {
		return nil, err
	}
	err = rlp.DecodeBytes(headersBlob, &headers)
	if err != nil {
		return nil, err
	}

	transactions := make([][]*common.L2Tx, 0)
	decryptedTxs := dataEncryptionService.Decrypt(encryptedRollup.BatchPayloads)
	encryptedTransactions, err := dataCompressionService.Decompress(decryptedTxs)
	if err != nil {
		return nil, err
	}
	err = rlp.DecodeBytes(encryptedTransactions, &transactions)
	if err != nil {
		return nil, err
	}

	batches := make([]*Batch, len(headers))
	for i := range headers {
		batches[i] = &Batch{
			Header:       &(headers[i]),
			Transactions: transactions[i],
		}
	}

	return &Rollup{
		Header:  encryptedRollup.Header,
		Batches: batches,
	}, nil
}
