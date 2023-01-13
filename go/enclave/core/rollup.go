package core

import (
	"math/big"
	"sync/atomic"

	"github.com/obscuronet/go-obscuro/go/enclave/crypto"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/go-obscuro/go/common"
)

// Rollup Data structure only for the internal use of the enclave since transactions are in clear
// Making changes to this struct will require GRPC + GRPC Converters regen
type Rollup struct {
	Header *common.RollupHeader

	hash atomic.Value
	// size   atomic.Value

	Transactions []*common.L2Tx
}

// Hash returns the keccak256 hash of b's header.
// The hash is computed on the first call and cached thereafter.
func (r *Rollup) Hash() *common.L2RootHash {
	// Temporarily disabling the caching of the hash because it's causing bugs.
	// Transforming a Rollup to an ExtRollup and then back to a Rollup will generate a different hash if caching is enabled.
	// Todo - re-enable
	//if hash := r.hash.Load(); hash != nil {
	//	return hash.(common.L2RootHash)
	//}
	v := r.Header.Hash()
	r.hash.Store(v)
	return &v
}

func (r *Rollup) NumberU64() uint64 { return r.Header.Number.Uint64() }
func (r *Rollup) Number() *big.Int  { return new(big.Int).Set(r.Header.Number) }

// IsGenesis indicates whether the rollup is the genesis rollup.
// TODO - #718 - Change this to a check against a hardcoded genesis hash.
func (r *Rollup) IsGenesis() bool {
	return r.Header.Number.Cmp(big.NewInt(int64(common.L2GenesisHeight))) == 0
}

func (r *Rollup) ToExtRollup(transactionBlobCrypto crypto.TransactionBlobCrypto) *common.ExtRollup {
	txHashes := make([]gethcommon.Hash, len(r.Transactions))
	for idx, tx := range r.Transactions {
		txHashes[idx] = tx.Hash()
	}

	return &common.ExtRollup{
		Header:          r.Header,
		TxHashes:        txHashes,
		EncryptedTxBlob: transactionBlobCrypto.Encrypt(r.Transactions),
	}
}

func ToRollup(encryptedRollup *common.ExtRollup, transactionBlobCrypto crypto.TransactionBlobCrypto) *Rollup {
	return &Rollup{
		Header:       encryptedRollup.Header,
		Transactions: transactionBlobCrypto.Decrypt(encryptedRollup.EncryptedTxBlob),
	}
}

func (r *Rollup) ToBatch() *Batch {
	return &Batch{
		Header:       r.Header.ToBatchHeader(),
		Transactions: r.Transactions,
	}
}
