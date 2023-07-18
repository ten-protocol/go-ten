package core

import "C"
import (
	"sync/atomic"

	"github.com/obscuronet/go-obscuro/go/common/compression"

	"github.com/ethereum/go-ethereum/rlp"

	"github.com/obscuronet/go-obscuro/go/enclave/crypto"

	"github.com/obscuronet/go-obscuro/go/common"
)

// todo - This should be a synthetic datastructure
type Rollup struct {
	Header  *common.RollupHeader
	Batches []*Batch
	hash    atomic.Value
}

// Hash returns the keccak256 hash of b's header.
// The hash is computed on the first call and cached thereafter.
func (r *Rollup) Hash() common.L2BatchHash {
	if hash := r.hash.Load(); hash != nil {
		return hash.(common.L2BatchHash)
	}
	v := r.Header.Hash()
	r.hash.Store(v)
	return v
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

	compressedTransactionsBlob, err := compression.CompressRollup(plaintextTransactionsBlob)
	if err != nil {
		return nil, err
	}

	compressedHeadersBlob, err := compression.CompressRollup(headersBlob)
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
